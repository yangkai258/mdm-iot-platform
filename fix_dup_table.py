#!/usr/bin/env python3
"""Fix orphaned template slots after </a-table>"""

import re
from pathlib import Path

VIEWS = Path(r"C:\Users\YKing\.openclaw\workspace\mdm-project\frontend\src\views")

def fix_dup_table(path):
    with open(path, 'r', encoding='utf-8', errors='replace') as f:
        content = f.read()
    
    if content.startswith('\ufeff'):
        content = content[1:]
    
    original = content
    
    # Pattern 1: </a-table>\n        <template #xxx>\n      </a-table>
    # Fix: remove first </a-table>, keep template between
    match = re.search(r'(</a-table>)\s*(<template #[^>]*>.*?</template>)\s*(</a-table>)', content, re.DOTALL)
    if match:
        content = match.group(1) + match.group(2) + match.group(3)
    
    if content != original:
        with open(path, 'w', encoding='utf-8', newline='\n') as f:
            f.write(content)
        return True
    return False

fixed = []
for vue in VIEWS.rglob("*.vue"):
    try:
        if fix_dup_table(vue):
            fixed.append(str(vue.name))
    except:
        pass

print(f"Fixed {len(fixed)} files")
for f in fixed[:15]:
    print(f"  {f}")