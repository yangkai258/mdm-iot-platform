#!/usr/bin/env python3
"""
Fix corrupted Vue files by:
1. Removing content between first </template> and second </template> (orphaned/corrupted template content)
2. Removing orphaned </a-table> tags after second </template>
3. Reconstructing valid Vue SFC
"""

import re
from pathlib import Path

VIEWS = Path(r"C:\Users\YKing\.openclaw\workspace\mdm-project\frontend\src\views")

def fix_corrupted_vue(content):
    original = content
    
    # Remove BOM
    if content.startswith('\ufeff'):
        content = content[1:]
    
    # Find all </template> positions
    template_positions = [m.start() for m in re.finditer(r'</template>', content)]
    
    if len(template_positions) >= 2:
        # Two or more </template> tags - likely corruption
        first_end = template_positions[0] + len('</template>')
        second_end = template_positions[1] + len('</template>')
        
        # The content between first and second </template> is corrupted
        before = content[:first_end]
        after = content[second_end:]
        
        # Remove orphaned </a-table> from the after part
        after = re.sub(r'\s*</a-table>\s*', '', after)
        
        content = before + after
    
    # Fix self-closing <a-table .../> followed by </a-table>
    content = re.sub(r'<a-table(\s+[^>]*?)/>\s*</a-table>', r'<a-table\1></a-table>', content)
    
    # Fix duplicate </a-table></a-table>
    content = re.sub(r'(</a-table>)\s*\1', r'\1', content)
    
    # Fix orphaned </a-table> after </a-card> (content outside template)
    # Pattern: </a-card>\n      </a-table>
    content = re.sub(r'(</a-card>)\s*\n\s*</a-table>', r'\1', content)
    
    return content

def is_valid_vue(content):
    """Basic validation - has template, script, and no obvious corruption markers."""
    has_template = '<template>' in content or '<template ' in content
    has_script = '<script' in content
    # Check for content after </template> that looks like template HTML (not script/style)
    parts = content.split('</template>', 1)
    if len(parts) > 1:
        after = parts[1].strip()
        # Should start with <script or <style or be empty
        if after and not after.startswith('<') and not after.startswith('\n<script') and not after.startswith('\n<style'):
            return False
    return has_template and has_script

fixed_files = []
error_files = []

for vue in VIEWS.rglob("*.vue"):
    try:
        with open(vue, 'r', encoding='utf-8', errors='replace') as f:
            content = f.read()
        
        if not content.strip():
            continue
            
        # Check if file is corrupted (has multiple </template> tags)
        template_count = content.count('</template>')
        
        if template_count >= 2 or '</a-table>' in content.split('</template>', 1)[-1] if '</template>' in content else False:
            new_content = fix_corrupted_vue(content)
            if new_content != content:
                with open(vue, 'w', encoding='utf-8', newline='\n') as f:
                    f.write(new_content)
                fixed_files.append(vue.name)
        elif not is_valid_vue(content):
            # Try to fix anyway
            new_content = fix_corrupted_vue(content)
            if new_content != content:
                with open(vue, 'w', encoding='utf-8', newline='\n') as f:
                    f.write(new_content)
                fixed_files.append(vue.name)
    except Exception as e:
        error_files.append(f"{vue.name}: {e}")

print(f"Fixed: {len(fixed_files)}")
for f in fixed_files:
    print(f"  {f}")
if error_files:
    print(f"\nErrors ({len(error_files)}):")
    for e in error_files[:10]:
        print(f"  {e}")
