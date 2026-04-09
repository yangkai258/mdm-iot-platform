import os
import re
from pathlib import Path

views_dir = Path(r"C:\Users\YKing\.openclaw\workspace\mdm-project\frontend\src\views")

def find_line_containing(lines, pattern, start, end):
    for i in range(start, min(end, len(lines))):
        if pattern in lines[i]:
            return i
    return -1

def fix_file(filepath):
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    original = content
    changed = False
    
    # Pattern: <a-table .../> then </a-table> (self-closing table bug)
    content = re.sub(r'<a-table([^>]*)/>\s*</a-table>', r'<a-table\1></a-table>', content)
    
    if content != original:
        changed = True
    
    # Pattern: </a-table> followed by <template...> followed by </a-table>
    # Fix: move templates inside the first table (before first </a-table>)
    # Use iterative replacement
    max_iterations = 50
    for _ in range(max_iterations):
        # Find pattern: </a-table> followed by lines that are all <template...>...</template> then </a-table>
        match = re.search(r'(</a-table>)\s*(<template\s[^>]*>.*?</template>)\s*(</a-table>)', content, re.DOTALL)
        if not match:
            break
        
        # Replace with: first </a-table>, then the templates
        # Actually this is wrong - templates should be BEFORE first </a-table>
        # The pattern says the table closed at first </a-table>, so we need to
        # put the templates inside by removing the first </a-table>
        # Let me handle this differently
        pass
    
    # Line-by-line approach
    lines = content.split('\n')
    n = len(lines)
    new_lines = []
    i = 0
    
    while i < n:
        # Look for <a-table ...> tag (not self-closing)
        if '<a-table' in lines[i] and '/>' not in lines[i] and '</a-table>' not in lines[i]:
            # Found opening table tag, collect it and all following content until properly closed
            # Keep scanning for slots and closing tag
            new_lines.append(lines[i])
            i += 1
            
            # Collect content until we find the matching </a-table>
            # But if we find </a-table> followed by <template>...</template></a-table>
            # we need to merge
            while i < n:
                if '</a-table>' in lines[i]:
                    # Check if next lines have orphaned templates
                    j = i + 1
                    while j < n and (lines[j].strip().startswith('<template ') or 
                                     lines[j].strip().startswith('</a-table')):
                        if '</a-table>' in lines[j]:
                            break
                        j += 1
                    
                    # Check if there are orphaned templates (lines between i and j are templates)
                    has_orphans = False
                    k = i + 1
                    while k < j and k < n:
                        if lines[k].strip().startswith('<template '):
                            has_orphans = True
                            break
                        k += 1
                    
                    if has_orphans:
                        # Collect orphaned templates
                        orphan_templates = []
                        k = i + 1
                        while k < j and k < n:
                            if lines[k].strip().startswith('<template '):
                                # Collect until </template>
                                tpl = [lines[k]]
                                k += 1
                                while k < j and k < n and '</template>' not in lines[k-1]:
                                    tpl.append(lines[k])
                                    k += 1
                                if k < n and '</template>' in lines[k-1]:
                                    pass  # already included
                                elif k < j:
                                    tpl.append(lines[k])
                                    k += 1
                                orphan_templates.append(''.join(tpl))
                            else:
                                k += 1
                        
                        # Find the REAL </a-table> (the one after orphans)
                        real_close = j
                        while real_close < n and '</a-table>' not in lines[real_close]:
                            real_close += 1
                        
                        # Add the first </a-table>
                        new_lines.append(lines[i].rstrip())
                        # Add orphaned templates inside table
                        for ot in orphan_templates:
                            new_lines.append('        ' + ot)
                        # Add the real </a-table>
                        if real_close < n:
                            new_lines.append(lines[real_close].rstrip())
                        i = real_close + 1
                        changed = True
                    else:
                        new_lines.append(lines[i].rstrip())
                        i += 1
                else:
                    new_lines.append(lines[i].rstrip())
                    i += 1
        else:
            new_lines.append(lines[i].rstrip())
            i += 1
    
    result = '\n'.join(new_lines)
    
    if changed or result != original:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(result)
        return True
    return False

fixed_files = []
for vue_file in views_dir.rglob("*.vue"):
    try:
        if fix_file(vue_file):
            fixed_files.append(vue_file.name)
    except Exception as e:
        print(f"Error {vue_file.name}: {e}")

print(f"Fixed {len(fixed_files)} files:")
for f in fixed_files:
    print(f"  - {f}")
