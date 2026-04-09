#!/usr/bin/env python3
"""Comprehensive Vue fixer for multiple patterns"""

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
    
    # Fix 1: Self-closing <a-table .../> followed by </a-table>
    content = re.sub(r'<a-table(\s+[^>/]*?)/>\s*(</a-table>)', r'<a-table\1></a-table>', content)
    
    # Fix 2: Missing </a-table> inside </a-tab-pane> where a-table slots exist
    # Pattern: </template> inside <a-table>...</a-table> then </a-tab-pane>
    # Need </a-table> before </a-tab-pane>
    content = re.sub(r'(</a-tab-pane>)(?!\s*</a-table>)', r'</a-table>\n    \1', content)
    
    # Fix 3: Missing </a-table> inside </a-form-item> where a-table slots exist
    content = re.sub(r'(</a-form-item>)(?!\s*</a-table>)', r'</a-table>\n    \1', content)
    
    if content != original:
        with open(path, 'w', encoding='utf-8', newline='\n') as f:
            f.write(content)
        return True
    return False

fixed = []
for vue in VIEWS.rglob("*.vue"):
    try:
        if fix_file(vue):
            fixed.append(str(vue.name))
    except Exception as e:
        print(f"Error {vue.name}: {e}")

print(f"Fixed {len(fixed)} files")
for f in fixed[:20]:
    print(f"  {f}")