#!/usr/bin/env python3
"""Keep fixing until no more changes"""
import subprocess, sys
from pathlib import Path

scripts = [
    r"C:\Users\YKing\.openclaw\workspace\mdm-project\fix_simple.py",
    r"C:\Users\YKing\.openclaw\workspace\mdm-project\fix_dup_table.py"
]

for i in range(10):
    total = 0
    for s in scripts:
        result = subprocess.run(["python", s], capture_output=True, text=True)
        # Extract count from output
        for line in result.stdout.split('\n'):
            if 'Fixed' in line and 'files' in line:
                count = int(line.split()[1])
                total += count
                print(f"Loop {i+1}: {s.split('\\')[-1]} fixed {count} files")
    if total == 0:
        print("No more files to fix")
        break
    print(f"Loop {i+1}: Total fixed {total}")
print("Done fixing")