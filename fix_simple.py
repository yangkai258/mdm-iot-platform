#!/usr/bin/env python3
"""Simple fix for self-closing tables"""
import re
from pathlib import Path
VIEWS = Path(r"C:\Users\YKing\.openclaw\workspace\mdm-project\frontend\src\views")
count = 0
for vue in VIEWS.rglob("*.vue"):
    try:
        with open(vue, 'r', encoding='utf-8', errors='replace') as f:
            content = f.read()
        original = content
        if content.startswith('\ufeff'):
            content = content[1:]
        # Fix self-closing tables
        content = re.sub(r'<a-table(\s+[^>/]*?)/>', r'<a-table\1></a-table>', content)
        if content != original:
            with open(vue, 'w', encoding='utf-8', newline='\n') as f:
                f.write(content)
            count += 1
    except: pass
print(f"Fixed {count} files")