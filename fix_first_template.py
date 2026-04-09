#!/usr/bin/env python3
"""Keep ONLY the FIRST </template>, remove all others"""
import re
from pathlib import Path

VIEWS = Path(r"C:\Users\YKing\.openclaw\workspace\mdm-project\frontend\src\views")

def fix_template_closes(path):
    try:
        with open(path, 'r', encoding='utf-8', errors='replace') as f:
            content = f.read()
        
        if content.startswith('\ufeff'):
            content = content[1:]
        
        # Find the first </template> and keep only that
        first_template = content.find('</template>')
        if first_template == -1:
            return False
            
        # Find second </template>
        second_template = content.find('</template>', first_template + 1)
        if second_template == -1:
            return False  # Only one, ok
            
        # Remove all </template> after the first one
        before = content[:first_template + 9]  # Keep first </template>
        after = content[second_template:]  # From second onwards
        
        # Remove all </template> from 'after'
        after = re.sub(r'</template>', '', after)
        
        new_content = before + after
        
        if new_content != content:
            with open(path, 'w', encoding='utf-8', newline='\n') as f:
                f.write(new_content)
            return True
    except:
        pass
    return False

fixed = []
for vue in VIEWS.rglob("*.vue"):
    if fix_template_closes(vue):
        fixed.append(vue.name)

print(f"Fixed {len(fixed)} files")
for f in fixed[:20]:
    print(f"  {f}")