# -*- coding: utf-8 -*-
# Fix TenantManagement.vue by removing duplicate content after </style>

path = "C:/Users/YKing/.openclaw/workspace/mdm-project/frontend/frontend/src/views/tenants/TenantManagement.vue"

with open(path, 'r', encoding='utf-8') as f:
    content = f.read()

# Find the FIRST </style> tag (the proper end) and keep everything up to and including it
# The duplicate content appears AFTER this tag
style_close = content.find('</style>')
if style_close != -1:
    # Keep everything up to and including </style>
    content = content[:style_close + 8]
    print(f"Found </style> at position {style_close}")
else:
    print("ERROR: No </style> tag found!")

with open(path, 'w', encoding='utf-8') as f:
    f.write(content)

print(f"Fixed. New length: {len(content)} characters")
