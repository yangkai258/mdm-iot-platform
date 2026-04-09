#!/usr/bin/env python3
"""
Vue file fixer - use LAST </template> as the actual template end.
Everything after LAST </template> should be <script> or <style> sections.
Everything between first and last </template> (that's not script/style) is corrupted.
"""

import re
from pathlib import Path

VIEWS = Path(r"C:\Users\YKing\.openclaw\workspace\mdm-project\frontend\src\views")

def fix_vue_final(path):
    with open(path, 'r', encoding='utf-8', errors='replace') as f:
        content = f.read()
    
    original = content
    
    # Remove BOM
    if content.startswith('\ufeff'):
        content = content[1:]
    
    # Find ALL </template> positions
    positions = [m.start() for m in re.finditer(r'</template>', content)]
    
    if len(positions) <= 1:
        # Either 0 or 1 </template> - either fine or missing close (handled elsewhere)
        return False
    
    # Use the LAST </template> as the actual template end
    last_end = positions[-1] + len('</template>')
    
    template_block = content[:last_end]
    after_block = content[last_end:]
    
    # Clean up the after_block - remove any orphaned HTML closing tags
    # (anything that's not <script or <style)
    lines = after_block.split('\n')
    clean_lines = []
    for line in lines:
        stripped = line.strip()
        if stripped.startswith('<script') or stripped.startswith('<style') or stripped.startswith('//') or stripped == '' or stripped.startswith('/*') or stripped.startswith('*'):
            clean_lines.append(line)
        elif stripped == '' or stripped.startswith('<'):
            # Skip orphaned HTML tags
            pass
        else:
            # Likely script content
            clean_lines.append(line)
    
    after_clean = '\n'.join(clean_lines)
    
    new_content = template_block + '\n' + after_clean
    
    # Also fix self-closing <a-table .../> followed by </a-table>
    new_content = re.sub(r'<a-table(\s+[^>]*?)/>\s*</a-table>', r'<a-table\1></a-table>', new_content)
    
    # Fix duplicate </a-table>
    new_content = re.sub(r'(</a-table>)\s*\1', r'\1', new_content)
    
    # Fix orphaned </a-table> after </a-card>
    new_content = re.sub(r'(</a-card>)\s*\n\s*</a-table>', r'\1', new_content)
    
    if new_content != content:
        with open(path, 'w', encoding='utf-8', newline='\n') as f:
            f.write(new_content)
        return True
    return False

fixed = []
for vue in VIEWS.rglob("*.vue"):
    try:
        if fix_vue_final(vue):
            fixed.append(str(vue.name))
    except Exception as e:
        print(f"Error {vue.name}: {e}")

print(f"Fixed {len(fixed)} files")
for f in fixed[:30]:
    print(f"  {f}")
if len(fixed) > 30:
    print(f"  ... and {len(fixed)-30} more")
