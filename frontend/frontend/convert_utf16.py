# -*- coding: utf-8 -*-
# Convert UTF-16 file to UTF-8
path = "C:/Users/YKing/.openclaw/workspace/mdm-project/frontend/frontend/src/views/permissions/Roles.vue"

with open(path, 'r', encoding='utf-16') as f:
    content = f.read()

# Write as UTF-8
with open(path, 'w', encoding='utf-8') as f:
    f.write(content)

print(f"Converted to UTF-8. Length: {len(content)}")
