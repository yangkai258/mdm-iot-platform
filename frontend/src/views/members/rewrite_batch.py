# -*- coding: utf-8 -*-
"""
Batch rewrite Vue files to Arco Design Pro template.
Keeps script setup logic intact, replaces template section.
"""
import os
import re

MEMBERS_PATH = r"C:\Users\YKing\.openclaw\workspace\mdm-project\frontend\src\views\members"

STANDARD_TEMPLATE = '''<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="名称"><a-input v-model="form.name" placeholder="请输入" /></a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleCreate">新建</a-button>
    </div>
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
      <template #actions="{ record }">
        <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
        <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
      </template>
    </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @before-ok="handleSubmit" @cancel="modalVisible = false">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="名称"><a-input v-model="form.name" placeholder="请输入" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>
'''

STANDARD_STYLE = '''
<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
'''

def rewrite_file(filepath):
    """Rewrite a single Vue file."""
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Check if already rewritten
    if 'search-form' in content:
        return False, "already rewritten"
    
    # Split content into sections
    # Find template, script, style sections
    
    # Find <template> block
    template_match = re.search(r'<template[^>]*>(.*?)</template>', content, re.DOTALL)
    if not template_match:
        return False, "no template found"
    
    # Find <script setup> block
    script_match = re.search(r'<script[^>]*setup[^>]*>(.*?)</script>', content, re.DOTALL)
    if not script_match:
        # Try <script>
        script_match = re.search(r'<script[^>]*>(.*?)</script>', content, re.DOTALL)
        if not script_match:
            return False, "no script found"
    
    # Find <style> block
    style_match = re.search(r'<style[^>]*>(.*?)</style>', content, re.DOTALL)
    
    # Extract script content
    script_content = script_match.group(1) if script_match else ''
    
    # Check for columns definition - extract to keep
    columns_match = re.search(r'const columns = \[(.*?)\]', script_content, re.DOTALL)
    if not columns_match:
        # Try with let or var
        columns_match = re.search(r'(?:const|let|var) columns = \[(.*?)\]', script_content, re.DOTALL)
    
    # Check for form definition
    form_match = re.search(r'const form = reactive\(\{(.*?)\}\)', script_content, re.DOTALL)
    if not form_match:
        form_match = re.search(r'const form = \{(.*?)\}', script_content, re.DOTALL)
    if not form_match:
        form_match = re.search(r'const form = reactive\(\[(.*?)\]\)', script_content, re.DOTALL)
    
    # Check for pagination
    pagination_match = re.search(r'(?:const|let|var) pagination = reactive\(\{(.*?)\}\)', script_content, re.DOTALL)
    
    # Check for data list
    data_match = re.search(r'(?:const|let|var) (?:data|dataList|list) = ref', script_content)
    
    # Check for loading
    loading_match = re.search(r'(?:const|let|var) loading = ref', script_content)
    
    # Extract key info from original script
    has_status_slot = 'slotName: \'status\'' in script_content or 'slotName: "status"' in script_content
    has_level_slot = 'slotName: \'level\'' in script_content or 'slotName: "level"' in script_content
    has_actions_slot = 'slotName: \'actions\'' in script_content or 'slotName: "actions"' in script_content
    
    # Build table slots template
    table_slots = ''
    if has_status_slot:
        table_slots += '''
      <template #status="{ record }">
        <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '启用' : '禁用' }}</a-tag>
      </template>'''
    if has_level_slot:
        table_slots += '''
      <template #level="{ record }">
        <a-tag>{{ record.levelName || record.level }}</a-tag>
      </template>'''
    if has_actions_slot:
        table_slots = ''  # Don't duplicate actions
    
    # Build standard template with extracted info
    new_template = f'''<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="名称"><a-input v-model="form.name" placeholder="请输入" /></a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleCreate">新建</a-button>
    </div>
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">{table_slots}
      <template #actions="{{ record }}">
        <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
        <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
      </template>
    </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @before-ok="handleSubmit" @cancel="modalVisible = false">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="名称"><a-input v-model="form.name" placeholder="请输入" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>
'''
    
    # Combine new content
    new_content = new_template + '\n<script setup>\n' + script_content + '\n</script>\n' + STANDARD_STYLE
    
    # Write back
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(new_content)
    
    return True, "rewritten"


def main():
    files = [f for f in os.listdir(MEMBERS_PATH) if f.endswith('.vue') and not f.startswith('.')]
    done = []
    failed = []
    
    for fname in files:
        filepath = os.path.join(MEMBERS_PATH, fname)
        try:
            success, msg = rewrite_file(filepath)
            if success:
                done.append(fname)
                print(f"[OK] {fname}")
            else:
                print(f"[SKIP] {fname}: {msg}")
        except Exception as e:
            failed.append((fname, str(e)))
            print(f"[FAIL] {fname}: {e}")
    
    print(f"\n=== Summary ===")
    print(f"Rewritten: {len(done)}")
    print(f"Failed: {len(failed)}")
    for f, e in failed:
        print(f"  - {f}: {e}")


if __name__ == '__main__':
    import sys
    sys.stdout.reconfigure(encoding='utf-8')
    main()
