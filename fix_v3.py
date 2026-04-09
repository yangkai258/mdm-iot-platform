#!/usr/bin/env python3
"""Comprehensive fixer for all Vue template issues"""
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
        
        # Fix 1: Self-closing <a-table .../>
        content = re.sub(r'<a-table(\s+[^>/]*?)/>', r'<a-table\1></a-table>', content)
        
        # Fix 2: Orphaned </a-table> before <template #slot>
        # Pattern: </a-table>\n        <template #xxx>...</template>\n      </a-table>
        content = re.sub(r'(</a-table>)\s*(<template #[^>]*>)', r'\1\n      \2', content, flags=re.DOTALL)
        
        # Fix 3: Remove duplicate </a-table> with orphaned template between
        content = re.sub(r'(</a-table>)\s*(<template #[^>]*>.*?</template>)\s*(</a-table>)', r'\1\2\3', content, flags=re.DOTALL)
        
        # Fix 4: Fix <a-form-item> without closing that comes after a-table 
        # content = re.sub(r'(</a-table>)\s*(</a-form-item>)', r'\1\2', content)
        
        # Fix 5: Extra closing tags after main template
        if content.count('</template>') > 1:
            # Find the last </template> that's the main one
            last_template_close = content.rfind('</template>')
            # Check if there are more </template> tags after the last one
            remaining = content[last_template_close + 11:]
            if '</template>' in remaining:
                # Remove extra </template> tags
                parts = content[:last_template_close].rsplit('</template>', 1)
                content = parts[0] + '</template>' + remaining.replace('</template>', '')
        
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