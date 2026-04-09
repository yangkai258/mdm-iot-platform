#!/usr/bin/env python3
"""Fix corrupted Vue files that have content outside <template> tags."""

import os
import re
from pathlib import Path

VIEWS = Path(r"C:\Users\YKing\.openclaw\workspace\mdm-project\frontend\src\views")

def fix_file(path):
    with open(path, 'r', encoding='utf-8', errors='replace') as f:
        content = f.read()
    
    original = content
    
    # Remove BOM
    if content.startswith('\ufeff'):
        content = content[1:]
    
    # Remove content that appears AFTER the closing </template> tag
    # Pattern: </template>\n...content outside template...
    if '</template>' in content:
        parts = content.split('</template>', 1)
        template_part = parts[0] + '</template>'
        # Keep only the <template>...</template> part plus <style> and <script> inside script tags
        content = template_part
        # Re-attach script and style sections if they exist (they're already inside the template split)
        # Actually template part already has everything, just make sure we don't have garbage after
    
    # Fix self-closing <a-table .../> followed by </a-table>
    content = re.sub(r'<a-table(\s+[^>]*?)/>\s*</a-table>', r'<a-table\1></a-table>', content)
    
    # Fix orphaned </a-table> followed by </a-card> then another </a-table>
    # Pattern: </a-table>\n    </a-card>\n      </a-table>
    content = re.sub(r'(\n\s*</a-card>)\s*\n\s*</a-table>', r'\1', content)
    
    # Fix duplicate </a-table> patterns
    content = re.sub(r'(</a-table>)\s*\1', r'\1', content)
    
    if content != original:
        with open(path, 'w', encoding='utf-8', newline='\n') as f:
            f.write(content)
        return True
    return False

fixed = []
errors = []
for vue in VIEWS.rglob("*.vue"):
    try:
        if fix_file(vue):
            fixed.append(str(vue.name))
    except Exception as e:
        errors.append(f"{vue.name}: {e}")

print(f"Fixed {len(fixed)} files")
for f in fixed:
    print(f"  {f}")
if errors:
    print(f"\nErrors:")
    for e in errors:
        print(f"  {e}")
