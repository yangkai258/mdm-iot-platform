#!/usr/bin/env python3
"""Fix self-closing a-table AND orphaned template patterns"""
import re
from pathlib import Path

VIEWS = Path(r"C:\Users\YKing\.openclaw\workspace\mdm-project\frontend\src\views")

def fix_vue(path):
    try:
        with open(path, 'r', encoding='utf-8', errors='replace') as f:
            content = f.read()
        
        original = content
        if content.startswith('\ufeff'):
            content = content[1:]
        
        # Fix 1: Self-closing a-table
        content = re.sub(r'<a-table(\s+[^>/]*?)/>', r'<a-table\1></a-table>', content)
        
        # Fix 2: Duplicate </a-table> with orphaned template - remove extra </a-table>
        # Pattern: </a-table>\n        <template #xxx>...\n      </a-table>
        content = re.sub(r'(</a-table>)\s*(<template #[^>]*>.*?</template>)\s*(</a-table>)', r'\1\2', content, flags=re.DOTALL)
        
        if content != original:
            with open(path, 'w', encoding='utf-8', newline='\n') as f:
                f.write(content)
            return True
    except:
        pass
    return False

fixed = []
for vue in VIEWS.rglob("*.vue"):
    if fix_vue(vue):
        fixed.append(vue.name)

print(f"Fixed {len(fixed)} files")
for f in fixed[:15]:
    print(f"  {f}")