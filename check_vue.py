#!/usr/bin/env python3
import re
path = r"C:\Users\YKing\.openclaw\workspace\mdm-project\frontend\src\views\ModalTest.vue"
with open(path, 'rb') as f:
    data = f.read()
print(f'File size: {len(data)} bytes')
print(f'</template> count: {data.count(b"</template>")}')
print(f'<template count: {data.count(b"<template")}')
for m in re.finditer(b'<template|</template>', data):
    start = max(0, m.start()-20)
    end = min(len(data), m.end()+40)
    print(f'At {m.start()}: {repr(data[start:end])}')
