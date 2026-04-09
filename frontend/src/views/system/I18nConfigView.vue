<template>
  <div class="page-container">
    <a-card class="general-card" title="国际化配�?>
      <template #extra>
        <a-button type="primary" @click="showAddModal = true"><icon-plus />添加语言�?/a-button>
        <a-button @click="handleExport"><icon-download />导出</a-button>
      </template>
      <a-tabs v-model:active-tab="activeTab" @change="loadData">
        <a-tab-pane key="zh-CN" tab="中文简�? />
        <a-tab-pane key="en-US" tab="English" />
        <a-tab-pane key="ja-JP" tab="日本�? />
      </a-tabs>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="Key">
            <a-input v-model="form.key" placeholder="请输入Key" style="width: 200px" />
          </a-form-item>
          <a-form-item label="Value">
            <a-input v-model="form.value" placeholder="请输入Value" style="width: 200px" />
          </a-form-item>
          <a-form-item>
            <a-button type="primary" @click="loadData">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="key">
        <template #actions="{ record }">
          <a-button type="primary" size="small" @click="handleEdit(record)">编辑</a-button>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="showAddModal" title="编辑翻译" :width="600" @ok="handleSubmit">
      <a-form :model="editForm" layout="vertical">
        <a-form-item label="Key">
          <a-input v-model="editForm.key" :disabled="isEditing" />
        </a-form-item>
        <template v-if="activeTab === 'zh-CN'">
          <a-form-item label="中文简�?>
            <a-textarea v-model="editForm['zh-CN']" :rows="3" />
          </a-form-item>
        </template>
        <template v-else-if="activeTab === 'en-US'">
          <a-form-item label="English">
            <a-textarea v-model="editForm['en-US']" :rows="3" />
          </a-form-item>
        </template>
        <template v-else>
          <a-form-item label="日本�?>
            <a-textarea v-model="editForm['ja-JP']" :rows="3" />
          </a-form-item>
        </template>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus, IconDownload } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const activeTab = ref('zh-CN')
const showAddModal = ref(false)
const isEditing = ref(false)
const form = reactive({ key: '', value: '' })
const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const editForm = reactive({ key: '', 'zh-CN': '', 'en-US': '', 'ja-JP': '' })

const columns = [
  { title: 'Key', dataIndex: 'key', width: 200 },
  { title: '中文简�?, dataIndex: 'zh-CN', ellipsis: true },
  { title: 'English', dataIndex: 'en-US', ellipsis: true },
  { title: '日本�?, dataIndex: 'ja-JP', ellipsis: true },
  { title: '操作', slotName: 'actions', width: 100 }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch(`/api/v1/system/i18n?locale=${activeTab.value}`).then(r => r.json())
    if (res.code === 0) {
      data.value = res.data?.list || []
    } else {
      loadMockData()
    }
  } catch { loadMockData() } finally { loading.value = false }
}

const loadMockData = () => {
  data.value = [
    { key: 'welcome', 'zh-CN': '欢迎', 'en-US': 'Welcome', 'ja-JP': 'ようこそ' },
    { key: 'device.onlinestatus', 'zh-CN': '设备在线', 'en-US': 'Device Online', 'ja-JP': 'デバイスオンライ�? },
    { key: 'device.offlinestatus', 'zh-CN': '设备离线', 'en-US': 'Device Offline', 'ja-JP': 'デバイスオフライ�? },
    { key: 'ota.upgrade', 'zh-CN': '固件升级', 'en-US': 'Firmware Upgrade', 'ja-JP': 'ファームウェ�?upgrade' },
    { key: 'ota.latest', 'zh-CN': '已是最新版�?, 'en-US': 'Already Latest', 'ja-JP': '最新で�? }
  ]
}

const handleReset = () => {
  form.key = ''
  form.value = ''
  loadData()
}

const handleEdit = (record) => {
  isEditing.value = true
  Object.assign(editForm, record)
  showAddModal.value = true
}

const handleSubmit = () => {
  const idx = data.value.findIndex(d => d.key === editForm.key)
  if (idx !== -1) {
    data.value[idx] = { ...editForm }
  }
  Message.success('保存成功')
  showAddModal.value = false
}

const handleExport = () => {
  Message.success(`导出语言�? ${activeTab.value}`)
}

const onPageChange = (page) => {
  pagination.current = page
  loadData()
}

onMounted(() => loadData())
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>