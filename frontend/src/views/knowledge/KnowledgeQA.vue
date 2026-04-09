<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>知识库</a-breadcrumb-item>
      <a-breadcrumb-item>问答对管理</a-breadcrumb-item>
    </a-breadcrumb>

    <div class="pro-search-bar">
      <a-space>
        <a-input-search v-model="searchKeyword" placeholder="搜索问题/答案" style="width: 280px" @search="loadQA" search-button />
        <a-select v-model="filterTag" placeholder="标签" allow-clear style="width: 160px" @change="loadQA">
          <a-option v-for="tag in tagOptions" :key="tag" :value="tag">{{ tag }}</a-option>
        </a-select>
      </a-space>
    </div>

    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showCreateModal">新建问答</a-button>
        <a-button @click="handleExport">导出</a-button>
        <a-button @click="showImportModal">导入</a-button>
        <a-button @click="loadQA">刷新</a-button>
      </a-space>
    </div>

    <div class="pro-content-area">
      <a-table :columns="columns" :data="qaList" :loading="loading" :pagination="pagination" row-key="id" @page-change="handlePageChange">
        <template #tags="{ record }">
          <a-tag v-for="tag in record.tags" :key="tag" style="margin-right: 4px">{{ tag }}</a-tag>
        </template>
        <template #enabled="{ record }">
          <a-switch v-model="record.enabled" @change="toggleQA(record)" />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="editQA(record)">编辑</a-button>
            <a-button type="text" size="small" @click="previewQA(record)">预览</a-button>
            <a-button type="text" size="small" status="danger" @click="deleteQA(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 新建/编辑弹窗 -->
    <a-modal v-model:visible="modalVisible" :title="isEdit ? '编辑问答' : '新建问答'" @ok="submitQA" :width="680" :loading="submitting">
      <a-form :model="qaForm" layout="vertical">
        <a-form-item label="问题" required>
          <a-textarea v-model="qaForm.question" placeholder="输入用户可能问的问题" :rows="2" />
        </a-form-item>
        <a-form-item label="答案" required>
          <a-textarea v-model="qaForm.answer" placeholder="输入标准答案" :rows="4" />
        </a-form-item>
        <a-form-item label="标签">
          <a-select v-model="qaForm.tags" multiple placeholder="选择或输入标签" allow-create :style="{ width: '100%' }">
            <a-option v-for="tag in tagOptions" :key="tag" :value="tag">{{ tag }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="相似问题（可选）">
          <a-textarea v-model="qaForm.alternatives" placeholder="每行一条相似问题" :rows="3" />
        </a-form-item>
        <a-form-item label="启用">
          <a-switch v-model="qaForm.enabled" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 预览弹窗 -->
    <a-modal v-model:visible="previewVisible" title="问答预览" :width="560" :footer="null">
      <a-result v-if="!previewRecord" status="info" title="请选择一条记录" />
      <template v-else>
        <a-alert type="info" style="margin-bottom: 12px">
          <template #title>用户问</template>
          <div>{{ previewRecord.question }}</div>
        </a-alert>
        <a-alert type="success">
          <template #title>系统答</template>
          <div>{{ previewRecord.answer }}</div>
        </a-alert>
        <a-divider>相似问题</a-divider>
        <a-tag v-for="alt in previewRecord.alternatives" :key="alt" style="margin: 4px">{{ alt }}</a-tag>
      </template>
    </a-modal>

    <!-- 导入弹窗 -->
    <a-modal v-model:visible="importVisible" title="导入问答" @ok="submitImport" :width="480" :loading="importing">
      <a-form-item label="导入方式">
        <a-radio-group v-model="importMode">
          <a-radio value="merge">合并（追加）</a-radio>
          <a-radio value="replace">覆盖（清空后导入）</a-radio>
        </a-radio-group>
      </a-form-item>
      <a-form-item label="选择文件">
        <a-upload :limit="1" accept=".json,.csv,.xlsx" :custom-request="handleFileChange" />
      </a-form-item>
      <a-alert type="info" message="支持 JSON、CSV、XLSX 格式，每条记录需包含 question、answer 字段" />
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'

const loading = ref(false)
const modalVisible = ref(false)
const previewVisible = ref(false)
const importVisible = ref(false)
const submitting = ref(false)
const importing = ref(false)
const isEdit = ref(false)
const searchKeyword = ref('')
const filterTag = ref('')
const previewRecord = ref<any>(null)
const importMode = ref('merge')
const importFile = ref<any>(null)

const tagOptions = ['宠物喂食', '宠物健康', '设备使用', '账户问题', '固件升级', '常见问题']
const qaList = ref<any[]>([])
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })

const qaForm = reactive({
  id: 0, question: '', answer: '', tags: [] as string[], alternatives: '', enabled: true,
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '问题', dataIndex: 'question', ellipsis: true },
  { title: '答案预览', dataIndex: 'answer', ellipsis: true },
  { title: '标签', slotName: 'tags', width: 200 },
  { title: '启用', slotName: 'enabled', width: 80 },
  { title: '操作', slotName: 'actions', fixed: 'right', width: 200 },
]

const loadQA = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/knowledge/qa', {
      params: { page: pagination.current, page_size: pagination.pageSize, keyword: searchKeyword.value, tag: filterTag.value },
    })
    qaList.value = res.data.items || []
    pagination.total = res.data.total || 0
  } catch {
    qaList.value = [
      { id: 1, question: '如何给宠物添加新设备？', answer: '打开设备管理页面，点击添加设备按钮...', tags: ['设备使用'], alternatives: ['怎么绑定设备', '设备怎么连接'], enabled: true },
      { id: 2, question: '固件升级失败怎么办？', answer: '请检查网络连接，确保设备在线...', tags: ['固件升级', '常见问题'], alternatives: ['OTA升级不了'], enabled: true },
      { id: 3, question: '宠物健康数据在哪看？', answer: '在健康医疗菜单中查看健康报告...', tags: ['宠物健康'], alternatives: ['健康报告怎么查'], enabled: false },
    ]
    pagination.total = 3
  } finally { loading.value = false }
}

const showCreateModal = () => {
  isEdit.value = false
  Object.assign(qaForm, { id: 0, question: '', answer: '', tags: [], alternatives: '', enabled: true })
  modalVisible.value = true
}

const editQA = (record: any) => {
  isEdit.value = true
  Object.assign(qaForm, { ...record, alternatives: (record.alternatives || []).join('\n') })
  modalVisible.value = true
}

const previewQA = (record: any) => {
  previewRecord.value = record
  previewVisible.value = true
}

const submitQA = async () => {
  submitting.value = true
  try {
    const payload = { ...qaForm, alternatives: qaForm.alternatives.split('\n').filter(Boolean) }
    if (isEdit.value) {
      await axios.put(`/api/v1/knowledge/qa/${qaForm.id}`, payload)
      Message.success('更新成功')
    } else {
      await axios.post('/api/v1/knowledge/qa', payload)
      Message.success('创建成功')
    }
    modalVisible.value = false
    loadQA()
  } catch { Message.error('操作失败') } finally { submitting.value = false }
}

const toggleQA = async (record: any) => {
  try { await axios.put(`/api/v1/knowledge/qa/${record.id}`, { enabled: record.enabled }); Message.success('更新成功') }
  catch { record.enabled = !record.enabled; Message.error('更新失败') }
}

const deleteQA = async (record: any) => {
  try { await axios.delete(`/api/v1/knowledge/qa/${record.id}`); Message.success('删除成功'); loadQA() }
  catch { Message.error('删除失败') }
}

const handleExport = async () => {
  try {
    const res = await axios.get('/api/v1/knowledge/qa/export', { params: { keyword: searchKeyword.value } })
    const blob = new Blob([JSON.stringify(res.data, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a'); a.href = url; a.download = 'knowledge-qa.json'; a.click()
    URL.revokeObjectURL(url)
    Message.success('导出成功')
  } catch { Message.error('导出失败') }
}

const showImportModal = () => { importVisible.value = true }
const handleFileChange = (options: any) => { importFile.value = options.file; options.onSuccess() }

const submitImport = async () => {
  if (!importFile.value) { Message.warning('请选择文件'); return }
  importing.value = true
  try {
    const formData = new FormData()
    formData.append('file', importFile.value)
    formData.append('mode', importMode.value)
    await axios.post('/api/v1/knowledge/qa/import', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
    Message.success('导入成功')
    importVisible.value = false
    loadQA()
  } catch { Message.error('导入失败') } finally { importing.value = false }
}

const handlePageChange = (page: number) => { pagination.current = page; loadQA() }

onMounted(() => loadQA())
</script>
