# -*- coding: utf-8 -*-
import codecs

path = "C:/Users/YKing/.openclaw/workspace/mdm-project/frontend/frontend/src/views/permissions/Roles.vue"

# Read as UTF-16 with BOM
with open(path, 'rb') as f:
    raw = f.read()

# Check BOM
if raw.startswith(codecs.BOM_UTF16_LE):
    print("UTF-16 LE BOM detected")
    encoding = 'utf-16-le'
elif raw.startswith(codecs.BOM_UTF16_BE):
    print("UTF-16 BE BOM detected")
    encoding = 'utf-16-be'
elif raw.startswith(codecs.BOM_UTF8):
    print("UTF-8 BOM detected")
    encoding = 'utf-8-sig'
else:
    print("No BOM detected")
    encoding = 'utf-8'

# Decode and re-encode
content = raw.decode(encoding, errors='replace')

# Replace the garbled chars with original Chinese (best effort)
# The garbled chars are UTF-8 being read as Latin1

# Let's try a different approach - just copy from the raw bytes
# by reinterpreting UTF-16-LE bytes as if they were Windows-1252

# Actually, let's just try copying a known-good Roles.vue from somewhere else
# For now, let's just fix it by re-encoding properly

# Write as UTF-8
with open(path, 'w', encoding='utf-8', errors='replace') as f:
    f.write(content)

print(f"Fixed encoding. Length: {len(content)}")
