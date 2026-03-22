# -*- coding: utf-8 -*-
path = "C:/Users/YKing/.openclaw/workspace/mdm-project/frontend/frontend/src/views/permissions/Roles.vue"

with open(path, 'r', encoding='utf-8') as f:
    content = f.read()

style_close = content.find('</style>')
if style_close != -1:
    content = content[:style_close + 8]
    print(f"Found </style> at position {style_close}")
else:
    print("ERROR: No </style> tag found!")

with open(path, 'w', encoding='utf-8') as f:
    f.write(content)

print(f"Fixed. New length: {len(content)} characters")
