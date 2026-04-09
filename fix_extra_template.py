#!/usr/bin/env python3
"""Remove extra </template> tags - keep only the last one"""
import re
from pathlib import Path

VIEWS = Path(r"C:\Users\YKing\.openclaw\workspace\mdm-project\frontend\src\views")

def fix_extra_template(path):
    try:
        with open(path, 'r', encoding='utf-8', errors='replace') as f:
            content = f.read()
        
        if content.startswith('\ufeff'):
            content = content[1:]
        
        # Count </template> tags
        template_closes = list(re.finditer(r'</template>', content))
        
        if len(template_closes) > 1:
            # Keep only the LAST </template> as the real one
            # Remove all intermediate </template> tags
            last_pos = template_closes[-1].start()
            before_last = content[:last_pos]
            after_last = content[last_pos:]
            
            # Remove all </template> from before_last
            before_last = re.sub(r'</template>', '', before_last)
            
            # Also remove the last one from after_last, then add it back
            after_last = after_last.replace('</template>', '', 1)
            
            content = before_last + '</template>' + after_last
            
            with open(path, 'w', encoding='utf-8', newline='\n') as f:
                f.write(content)
            return True
    except:
        pass
    return False

fixed = []
for vue in VIEWS.rglob("*.vue"):
    if fix_extra_template(vue):
        fixed.append(vue.name)

print(f"Fixed {len(fixed)} files")
for f in fixed[:20]:
    print(f"  {f}")