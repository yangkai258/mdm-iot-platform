path = r"C:\Users\YKing\.openclaw\workspace\mdm-project\frontend\src\views\system\Monitor.vue"
with open(path, 'r', encoding='utf-8', errors='replace') as f:
    content = f.read()

old = '''      </a-table>
        <template #uptime="{ record }">{{ record.uptime || '-' }}</template>
      </a-table>'''

new = '''        <template #uptime="{ record }">{{ record.uptime || '-' }}</template>
      </a-table>'''

if old in content:
    content = content.replace(old, new)
    with open(path, 'w', encoding='utf-8', newline='\n') as f:
        f.write(content)
    print('Fixed Monitor.vue')
else:
    print('Pattern not found')