#!/usr/bin/env python3
"""Comprehensive fix for all known patterns"""
import re
from pathlib import Path
VIEWS = Path(r"C:\Users\YKing\.openclaw\workspace\mdm-project\frontend\src\views")

def comprehensive_fix(path):
    with open(path, 'r', encoding='utf-8', errors='replace') as f:
        content = f.read()
    original = content
    if content.startswith('\ufeff'):
        content = content[1:]
    
    # Fix self-closing <a-table .../>
    content = re.sub(r'<a-table(\s+[^>/]*?)/>', r'<a-table\1></a-table>', content)
    
    # Fix duplicate </a-table> after orphaned templates
    content = re.sub(r'(</a-table>)\s*(<template #[^>]*>.*?</template>)\s*(</a-table>)', r'\1\2\3', content, flags=re.DOTALL)
    
    # Fix missing </a-table> before </a-tab-pane>
    content = re.sub(r'(</a-tab-pane>)(?!\s*</a-table>)', r'</a-table>\n    \1', content)
    
    # Fix missing </a-table> before </a-form-item> where table slots exist
    content = re.sub(r'(</a-form-item>)(?=\s*</a-form>', r'</a-table>\n    \1', content)
    
    if content != original:
        with open(path, 'w', encoding='utf-8', newline='\n') as f:
            f.write(content)
        return True
    return False

fixed = sum(1 for v in VIEWS.rglob("*.vue") if comprehensive_fix(v))
print(f"Fixed {fixed} files")