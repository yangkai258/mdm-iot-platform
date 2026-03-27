# -*- coding: utf-8 -*-
# Fix Roles.vue

path = "C:/Users/YKing/.openclaw/workspace/mdm-project/frontend/frontend/src/views/permissions/Roles.vue"

with open(path, 'r', encoding='utf-8') as f:
    content = f.read()

# Count occurrences
template_open = content.count('<template>')
template_close = content.count('</template>')
style_open = content.count('<style')
style_close = content.count('</style>')

print(f"template: {template_open} open, {template_close} close")
print(f"style: {style_open} open, {style_close} close")

# Find the first </style> and truncate after it
style_close_pos = content.find('</style>')
if style_close_pos != -1:
    remaining = content[style_close_pos + 8:]
    if remaining.strip():
        print(f"Content after first </style>: {remaining[:100]}...")
        content = content[:style_close_pos + 8]
        print("Truncated")
else:
    print("No </style> found!")

with open(path, 'w', encoding='utf-8') as f:
    f.write(content)

print(f"Fixed. New length: {len(content)}")
