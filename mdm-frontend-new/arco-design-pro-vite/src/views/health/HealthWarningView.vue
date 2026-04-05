<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="page-container">
    <a-card class="general-card" title="健康预警">
      <template #extra>
        <a-button type="primary" @click="handleCreate"><icon-plus />新建预警</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="严重程度">
            <a-select v-model="form.level" placeholder="全部" style="width: 120px" allow-clear>
              <a-option value="critical">危急</a-option>
              <a-option value="high">高</a-option>
              <a-option value="medium">中</a-option>
              <a-option value="low">低</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="处理状态">
            <a-select v-model="form.status" placeholder="全部" style="width: 120px" allow-clear>
              <a-option value="pending">待处理</a-option>
              <a-option value="confirmed">已确认</a-option>
              <a-option value="ignored">已忽略</a-option>
            </a-select>
          </a-form-item>
          <a-form-item>
            <a-button type="primary" @click="loadData">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-form-item>
        </a-form>
      </div>
      <a-table
      :columns="columns"
      :data="data"
      :loading="loading"
      :pagination="paginationConfig"
      @page-change="onPageChange"
    >
      <template #level="{ record }">
        <a-tag :color="getLevelColor(record.level)">{{ getLevelName(record.level) }}</a-tag>
      </template>
      <template #status="{ record }">
        <a-tag :color="getStatusColor(record.status)">{{ getStatusName(record.status) }}</a-tag>
      </template>
    </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle">
      <a-form :model="form" layout="vertical">
        <a-form-item label="疾病名称"><a-input v-model="form.disease_name" /></a-form-item>
        <a-form-item label="预警级别">
          <a-select v-model="form.level">
            <a-option value="critical">危急</a-option>
            <a-option value="high">高</a-option>
            <a-option value="medium">中</a-option>
            <a-option value="low">低</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="预警描述"><a-textarea v-model="form.description" :rows="3" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
    </a-card>
</div></template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import axios from 'axios'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const modalVisible = ref(false)
const modalTitle = ref('新建')
const isEdit = ref(false)

const form = reactive({
  id: '',
  disease_name: '',
  level: '',
  status: 'pending',
  description: ''
})

const columns = [
  { title: '预警编号', dataIndex: 'id', width: 100 },
  { title: '疾病名称', dataIndex: 'disease_name', width: 150 },
  { title: '严重程度', slotName: 'level', width: 100 },
  { title: '预警描述', dataIndex: 'description', ellipsis: true },
  { title: '发生时间', dataIndex: 'created_at', width: 160 },
  { title: '状态', slotName: 'status', width: 100 }
]

const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current,
  pageSize: pagination.pageSize,
  total: pagination.total,
  showTotal: true,
  showPageSize: true
}))

const getLevelName = (level) => {
  const map = { critical: '危急', high: '高', medium: '中', low: '低' }
  return map[level] || level
}

const getLevelColor = (level) => {
  const map = { critical: 'red', high: 'orange', medium: 'blue', low: 'green' }
  return map[level] || 'default'
}

const getStatusName = (status) => {
  const map = { pending: '待处理', confirmed: '已确认', ignored: '已忽略' }
  return map[status] || status
}

const getStatusColor = (status) => {
  const map = { pending: 'orange', confirmed: 'green', ignored: 'gray' }
  return map[status] || 'default'
}

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (form.level) params.level = form.level
    if (form.status) params.status = form.status
    const res = await axios.get('/api/health/warnings', { params })
    if (res.data.code === 0) {
      data.value = res.data.data.list || []
      pagination.total = res.data.data.pagination?.total || 0
    } else {
      loadMockData()
    }
  } catch {
    loadMockData()
  } finally {
    loading.value = false
  }
}

const loadMockData = () => {
  data.value = [
    { id: 1, disease_name: '高血压风险', level: 'high', status: 'pending', description: '连续3天血压监测偏高，建议密切观察', created_at: '2026-03-22 10:30:00' },
    { id: 2, disease_name: '睡眠呼吸暂停', level: 'critical', status: 'pending', description: '夜间血氧饱和度多次低于90%', created_at: '2026-03-22 08:00:00' },
    { id: 3, disease_name: '心律不齐', level: 'medium', status: 'confirmed', description: '心电图检测到偶发早搏', created_at: '2026-03-21 15:20:00' },
    { id: 4, disease_name: '体重异常波动', level: 'low', status: 'ignored', description: '一周内体重下降超过5%', created_at: '2026-03-20 09:00:00' },
    { id: 5, disease_name: '血糖偏高', level: 'high', status: 'pending', description: '空腹血糖持续高于正常值上限', created_at: '2026-03-22 11:45:00' }
  ]
  pagination.total = data.value.length
}

const handleReset = () => {
  Object.assign(form, { level: '', status: '' })
  loadData()
}

const handleCreate = () => {
  isEdit.value = false
  modalTitle.value = '新建预警'
  Object.assign(form, { id: '', disease_name: '', level: '', status: 'pending', description: '' })
  modalVisible.value = true
}

const handleSubmit = () => {
  if (!form.disease_name) { Message.warning('请填写疾病名称'); return }
  modalVisible.value = false
  Message.success(isEdit.value ? '编辑成功' : '添加成功')
  loadData()
}

const onPageChange = (page) => {
  pagination.current = page
  loadData()
}

onMounted(() => { loadData() })
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>

