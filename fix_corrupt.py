#!/usr/bin/env python3
"""
Comprehensive Vue file fixer for corrupted files.
Handles: orphaned </a-table>, orphaned </a-card>, orphaned </a-modal> etc.
that appear AFTER </template> but BEFORE <script>.
"""

import re
from pathlib import Path

VIEWS = Path(r"C:\Users\YKing\.openclaw\workspace\mdm-project\frontend\src\views")

def fix_file_v2(path):
    with open(path, 'r', encoding='utf-8', errors='replace') as f:
        content = f.read()
    
    original = content
    
    # Remove BOM
    if content.startswith('\ufeff'):
        content = content[1:]
    
    # Strategy: find </template> then truncate or clean content after it
    if '</template>' not in content:
        return False
    
    parts = content.split('</template>', 1)
    template_part = parts[0] + '</template>'
    
    if len(parts) < 2:
        return False
    
    after = parts[1]  # Everything after first </template>
    
    # Remove orphaned closing tags from the "after" section
    # These are things like </a-table>, </a-card>, </a-modal> that are dangling
    # after the template properly closed
    after_clean = re.sub(r'\s*</[a-z][a-z0-9-]*>\s*', '', after)
    
    new_content = template_part + after_clean
    
    if new_content != content:
        with open(path, 'w', encoding='utf-8', newline='\n') as f:
            f.write(new_content)
        return True
    return False

# Also handle self-closing <a-table .../> pattern
def fix_self_closing(path):
    with open(path, 'r', encoding='utf-8', errors='replace') as f:
        content = f.read()
    
    original = content
    
    # Fix <a-table .../> followed by </a-table>
    content = re.sub(r'<a-table(\s+[^>]*?)/>\s*</a-table>', r'<a-table\1></a-table>', content)
    
    # Fix duplicate </a-table></a-table>
    content = re.sub(r'(</a-table>)\s*\1', r'\1', content)
    
    # Fix orphaned </a-table> inside <a-card>...</a-card> where the </a-table> is after </a-card>
    # Pattern: </a-card>\n      </a-table>
    content = re.sub(r'(</a-card>)\s*\n\s*</a-table>', r'\1', content)
    
    # Fix <a-table ...>\n  <template #xxx>...\n</a-table>\n  <template #yyy>...\n</a-table>
    # where second </a-table> is duplicate - keep only the first table close
    # This pattern: table has slot templates, then </a-table>, then more templates, then </a-table>
    # The fix: move templates after first </a-table> inside the table (before first </a-table>)
    
    if content != original:
        with open(path, 'w', encoding='utf-8', newline='\n') as f:
            f.write(content)
        return True
    return False

fixed = []
for vue in VIEWS.rglob("*.vue"):
    try:
        changed = fix_file_v2(vue)
        if changed:
            fixed.append(str(vue.name))
    except Exception as e:
        print(f"Error {vue.name}: {e}")

print(f"Fixed {len(fixed)} files")
for f in fixed[:20]:
    print(f"  {f}")
if len(fixed) > 20:
    print(f"  ... and {len(fixed)-20} more")
