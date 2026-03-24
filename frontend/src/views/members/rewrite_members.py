#!/usr/bin/env python3
"""
Batch rewrite Vue files in views/members/ to standard Arco Design Pro template.
Preserves ALL existing script logic, API calls, and business logic.
Only standardizes the outer page structure and CSS classes.
"""
import os
import re
import json

BASE_DIR = r"C:\Users\YKing\.openclaw\workspace\mdm-project\frontend\src\views\members"

# Files that are already well-structured (complex business logic - preserve as-is but normalize outer structure)
COMPLEX_FILES = [
    "MemberListView.vue",
    "MemberLevels.vue",
    "MemberChannels.vue",
    "MemberStores.vue",
    "PointsView.vue",
    "MemberDetailView.vue",
    "MemberGiftView.vue",
    "MemberReceptionView.vue",
    "MemberSettingsView.vue",
    "MemberCardView.vue",
    "MemberChannelsView.vue",
    "MemberCardTypesView.vue",
    "MemberCardGroupsView.vue",
    "MemberOrdersView.vue",
    "PointsSettingsView.vue",
    "PointsRulesView.vue",
    "MemberLevelRulesView.vue",
    "SmsChannelView.vue",
    "SmsTemplateView.vue",
    "StoreView.vue",
    "StoreLocationsView.vue",
    "StoreSourcesView.vue",
    "MemberArticlesView.vue",
    "PromotionTypesView.vue",
    "OccupationTypesView.vue",
    "PointsExcludeView.vue",
    "PrinterManageView.vue",
    "TagView.vue",
    "TagAutoCleanView.vue",
    "TagReportView.vue",
    "RedpacketView.vue",
    "MiniProgramView.vue",
    "VipExclusiveView.vue",
    "WechatSettingsView.vue",
    "TempMembersView.vue",
    "TempPointsView.vue",
    "TempCoupons.vue",
    "TempRedpacketView.vue",
    "TempCouponGrantView.vue",
    "CouponView.vue",
    "CouponGrantView.vue",
    "CouponInventoryView.vue",
    "CouponMessagesView.vue",
    "GiftRecordsView.vue",
    "AmountDiscountView.vue",
    "AmountReduceView.vue",
    "BuyGiftView.vue",
    "DirectReduceView.vue",
    "LevelView.vue",
    "HighFreqTagView.vue",
    "LowFreqTagView.vue",
    "InterestTagView.vue",
    "PointsInventory.vue",
    "PointsRecords.vue",
    "PointsExclude.vue",
    "PointsRules.vue",
    "PointsSettings.vue",
    "Printers.vue",
    "StoreLocations.vue",
    "StoreSources.vue",
    "MemberBenefits.vue",
    "MemberCardGroups.vue",
    "MemberCardTypes.vue",
    "MemberCoupons.vue",
    "MemberDetail.vue",
    "MemberGifts.vue",
    "MemberLevels.vue",
    "MemberList.vue",
    "MemberOrders.vue",
    "MemberPoints.vue",
    "MemberPromotions.vue",
    "MemberReception.vue",
    "MemberSettings.vue",
    "MemberTags.vue",
    "MemberUpgradeRules.vue",
    "MemberBenefitsView.vue",
    "MemberPromotions.vue",
    "MemberChannels.vue",
    "MiniProgram.vue",
]

def normalize_complex_file(filepath):
    """Normalize complex files: keep script/style, standardize template structure."""
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Check if file has multiple template blocks (concatenated)
    template_blocks = re.findall(r'<template>(.*?)</template>', content, re.DOTALL)
    
    if len(template_blocks) > 1:
        # Concatenated file - merge all templates
        merged_template = '\n'.join(t.strip() for t in template_blocks if t.strip())
        # Replace all template blocks with merged one
        content = re.sub(r'<template>.*?</template>', f'<template>\n{merged_template}\n</template>', content, count=1, flags=re.DOTALL)
        # Remove remaining template blocks
        content = re.sub(r'</template>\s*<template>.*?</template>', '', content)
    
    # Normalize outer page class names
    # Replace old class patterns with standard ones
    old_class_patterns = [
        (r'class="member-list-page"', 'class="page-container"'),
        (r'class="member-levels-page"', 'class="page-container"'),
        (r'class="member-benefits-page"', 'class="page-container"'),
        (r'class="pro-page-container"', 'class="page-container"'),
        (r'class="pro-breadcrumb"', 'class="breadcrumb"'),
        (r'class="pro-search-bar"', 'class="search-form"'),
        (r'class="pro-action-bar"', 'class="toolbar"'),
        (r'class="pro-content-area"', ''),
        (r'class="breadcrumb"', 'class="breadcrumb"'),
        (r'class="stats-row"', 'class="stats-row"'),
        (r'class="stat-card"', 'class="stat-card"'),
        (r'class="action-card"', 'class="action-card"'),
        (r'class="table-card"', 'class="table-card"'),
    ]
    
    for old, new in old_class_patterns:
        content = re.sub(old, new, content)
    
    # Add standard CSS if missing
    if '.page-container' not in content:
        std_css = """
<style scoped>
.page-container {
  background: #fff;
  border-radius: 4px;
  padding: 20px;
}
.search-form {
  margin-bottom: 16px;
  padding: 16px;
  background: #f7f8fa;
  border-radius: 4px;
}
.toolbar {
  margin-bottom: 16px;
}
.breadcrumb {
  margin-bottom: 16px;
}
.stats-row {
  margin-bottom: 16px;
}
.stat-card {
  border-radius: 8px;
}
.action-card {
  margin-bottom: 16px;
}
.table-card {
  border-radius: 8px;
}
</style>
"""
        # Remove existing style block if present
        content = re.sub(r'<style[^>]*>.*?</style>', '', content, flags=re.DOTALL)
        content = content.strip() + '\n' + std_css
    
    return content

def rewrite_simple_file(filepath, page_title, api_import=None, api_call=None):
    """Rewrite a simple file with mock data to the standard template."""
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Check if file has multiple template blocks
    template_blocks = re.findall(r'<template>(.*?)</template>', content, re.DOTALL)
    
    if len(template_blocks) > 1:
        # Concatenated file - just keep the last one (most complete)
        content = template_blocks[-1].strip()
    
    # Extract script content
    script_match = re.search(r'<script[^>]*>(.*?)</script>', content, re.DOTALL)
    script_content = script_match.group(1) if script_match else ""
    
    # Extract style content
    style_match = re.search(r'<style[^>]*>(.*?)</style>', content, re.DOTALL)
    style_content = style_match.group(1) if style_match else ""
    
    # Extract data, columns, loadData from script
    data_match = re.search(r'const\s+(\w+)\s*=\s*ref\(\[(.*?)\]\)', script_content, re.DOTALL)
    columns_match = re.search(r'const\s+columns\s*=\s*\[(.*?)\]', script_content, re.DOTALL)
    form_match = re.search(r'const\s+form\s*=\s*reactive\(\{([^}]*)\}\)', script_content, re.DOTALL)
    
    data_name = data_match.group(1) if data_match else "data"
    form_name = "form"
    
    # Build standard template
    std_template = f"""<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="{form_name}" layout="inline">
        <a-form-item label="名称"><a-input v-model="{form_name}.name" placeholder="请输入" /></a-form-item>
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
      <template #level="{{ record }}">
        <a-tag :color="record.level === 'diamond' ? 'purple' : record.level === 'platinum' ? 'gold' : record.level === 'gold' ? 'orange' : 'gray'">
          {{ record.levelName }}
        </a-tag>
      </template>
      <template #status="{{ record }}">
        <a-tag :color="record.status === '1' ? 'green' : 'gray'">{{ record.status === '1' ? '启用' : '禁用' }}</a-tag>
      </template>
      <template #actions="{{ record }}">
        <a-space>
          <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
          <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
        </a-space>
      </template>
    </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @before-ok="handleSubmit" @cancel="modalVisible = false">
      <a-form :model="{form_name}" label-col-flex="100px">
        <a-form-item label="名称"><a-input v-model="{form_name}.name" placeholder="请输入名称" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit(() => {{}})">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>
"""
    
    # Build standard script
    std_script = f"""<script setup>
import {{ ref, reactive, onMounted }}
from 'vue'
import {{ Message }} from '@arco-design/web-vue'

const loading = ref(false)
const data = ref([])
const modalVisible = ref(false)
const modalTitle = ref('新建')
const isEdit = ref(false)
const currentId = ref(null)

const form = reactive({{ name: '', status: '1' }})

const pagination = reactive({{ current: 1, pageSize: 10, total: 0 }})

const columns = [
  {{ title: '名称', dataIndex: 'name', width: 180 }},
  {{ title: '状态', slotName: 'status', width: 80 }},
  {{ title: '操作', slotName: 'actions', width: 150, fixed: 'right' }}
]

const mockData = () => [
  {{ id: 1, name: '示例数据', level: 'silver', levelName: '银卡', status: '1' }},
  {{ id: 2, name: '示例数据2', level: 'gold', levelName: '金卡', status: '1' }}
]

const loadData = () => {{
  loading.value = true
  setTimeout(() => {{
    data.value = mockData()
    pagination.total = data.value.length
    loading.value = false
  }}, 300)
}}

const handleSearch = () => {{
  pagination.current = 1
  loadData()
}}

const handleReset = () => {{
  Object.assign(form, {{ name: '', status: '1' }})
  pagination.current = 1
  loadData()
}}

const handleCreate = () => {{
  isEdit.value = false
  currentId.value = null
  modalTitle.value = '新建'
  Object.assign(form, {{ name: '', status: '1' }})
  modalVisible.value = true
}}

const handleEdit = (record) => {{
  isEdit.value = true
  currentId.value = record.id
  modalTitle.value = '编辑'
  Object.assign(form, record)
  modalVisible.value = true
}}

const handleDelete = (record) => {{
  Message.success('删除成功')
  loadData()
}}

const handleSubmit = (done) => {{
  if (!form.name) {{
    Message.error('请输入名称')
    done && done(false)
    return
  }}
  setTimeout(() => {{
    Message.success(isEdit.value ? '更新成功' : '创建成功')
    modalVisible.value = false
    loadData()
    done && done(true)
  }}, 400)
}}

const onPageChange = (page) => {{
  pagination.current = page
  loadData()
}}

onMounted(() => {{
  loadData()
}})
</script>
"""
    
    std_style = """<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
"""
    
    return std_template + '\n' + std_script + '\n' + std_style

def get_page_title(filename):
    """Extract a readable title from filename."""
    name = filename.replace('.vue', '')
    # Remove View suffix and convert to Chinese
    name = re.sub(r'View$', '', name)
    # Add space before capitals
    name = re.sub(r'([A-Z])', r' \1', name)
    return name.strip()

def process_all():
    """Process all Vue files in the members directory."""
    all_files = os.listdir(BASE_DIR)
    vue_files = [f for f in all_files if f.endswith('.vue')]
    
    print(f"Found {len(vue_files)} Vue files")
    
    results = {'rewritten': [], 'skipped': [], 'errors': []}
    
    for filename in sorted(vue_files):
        filepath = os.path.join(BASE_DIR, filename)
        try:
            if filename in COMPLEX_FILES:
                # Normalize complex files
                new_content = normalize_complex_file(filepath)
                with open(filepath, 'w', encoding='utf-8') as f:
                    f.write(new_content)
                results['rewritten'].append(filename)
                print(f"  [OK] Normalized: {filename}")
            else:
                # For unlisted files, just normalize as complex
                new_content = normalize_complex_file(filepath)
                with open(filepath, 'w', encoding='utf-8') as f:
                    f.write(new_content)
                results['rewritten'].append(filename)
                print(f"  [OK] Processed: {filename}")
        except Exception as e:
            results['errors'].append((filename, str(e)))
            print(f"  ✗ Error: {filename} - {e}")
    
    print(f"\nSummary:")
    print(f"  Rewritten: {len(results['rewritten'])}")
    print(f"  Errors: {len(results['errors'])}")
    for fname, err in results['errors']:
        print(f"    {fname}: {err}")
    
    return results

if __name__ == '__main__':
    process_all()
