# -*- coding: utf-8 -*-
with open('C:/Users/YKing/.openclaw/workspace/mdm-project/frontend/frontend/original_roles.vue', 'r', encoding='utf-16') as f:
    orig = f.read()
with open('C:/Users/YKing/.openclaw/workspace/mdm-project/frontend/frontend/src/views/permissions/Roles.vue', 'r', encoding='utf-8') as f:
    curr = f.read()

print(f'Original: {len(orig)} chars')
print(f'Current: {len(curr)} chars')
print(f'--- Original file ---')
print(f'template: {orig.count("<template")} open, {orig.count("</template>")} close')
print(f'script: {orig.count("<script")} open, {orig.count("</script>")} close')
print(f'style: {orig.count("<style")} open, {orig.count("</style>")} close')
print(f'--- Current file ---')
print(f'template: {curr.count("<template")} open, {curr.count("</template>")} close')
print(f'script: {curr.count("<script")} open, {curr.count("</script>")} close')
print(f'style: {curr.count("<style")} open, {curr.count("</style>")} close')
