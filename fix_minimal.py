#!/usr/bin/env python3
"""Minimal Vue fixer - only fix specific known issues"""

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
    
    # Fix only self-closing <a-table .../> 
    content = re.sub(r'<a-table(\s+[^>/]*?)/>', r'<a-table\1></a-table>', content)
    
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
        pass

print(f"Fixed {len(fixed)} files")