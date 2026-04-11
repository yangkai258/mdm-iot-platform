import re
import os
from pathlib import Path

VIEWS_DIR = Path("C:/Users/YKing/.openclaw/workspace/mdm-project/frontend/src/views")
FE_DIR = "C:/Users/YKing/.openclaw/workspace/mdm-project/frontend"

def find_errors():
    """Run npm build and capture errors"""
    import subprocess
    result = subprocess.run(
        ["npm", "run", "build"],
        cwd=FE_DIR,
        capture_output=True,
        text=True,
        timeout=60
    )
    if result.returncode == 0:
        return None
    
    output = result.stderr + result.stdout
    # Match: src/views/path/file.vue (LINE:COL): Error message
    match = re.search(r'src/views/([^:]+\.vue) \((\d+):(\d+)\): (Invalid end tag|Element is missing end tag)', output)
    if match:
        return match.group(1), int(match.group(2))
    return None

def fix_file_manual(filepath, error_line):
    """Manually fix specific known patterns"""
    with open(filepath, 'r', encoding='utf-8') as f:
        lines = f.readlines()
    
    modified = False
    new_lines = []
    
    # For DeviceDetail we know line 198 was orphan </a-table>
    # Pattern: orphan </a-table> that appears AFTER template slots but BEFORE final </a-table>
    if 'DeviceDetail.vue' in str(filepath):
        # Check around error line for orphan pattern
        for i, line in enumerate(lines):
            # Look for pattern: </a-table> followed by more template slots then another </a-table>
            if '</a-table>' in line and i < len(lines) - 1:
                next_line = lines[i+1].strip() if i+1 < len(lines) else ''
                # If next line starts with <template, this is orphan
                if next_line.startswith('<template'):
                    # Skip this orphan </a-table>
                    modified = True
                    print(f"  Removing orphan </a-table> at line {i+1}")
                    continue
            new_lines.append(line)
        
        if modified:
            with open(filepath, 'w', encoding='utf-8') as f:
                f.writelines(new_lines)
            return True
    
    return False

# Simple approach: iterate and fix manually
for _ in range(50):
    err = find_errors()
    if err is None:
        print("BUILD SUCCESS!")
        break
    
    filepath, line = err
    full_path = VIEWS_DIR / filepath
    print(f"Error: {filepath} line {line}")
    
    # Check and fix if known pattern
    if not fix_file_manual(full_path, line):
        print(f"  Manual fix not implemented for this error")
