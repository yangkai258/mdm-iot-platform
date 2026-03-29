<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>研究平台</a-breadcrumb-item>
      <a-breadcrumb-item>研究项目管理</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索栏 -->
    <div class="pro-search-bar">
      <a-space wrap>
        <a-input-search
          v-model="filters.keyword"
          placeholder="搜索项目名称/描述"
          style="width: 280px"
          search-button
          @search="loadProjects"
        />
        <a-select v-model="filters.status" placeholder="项目状态" allow-clear style="width: 140px" @change="loadProjects">
          <a-option value="draft">草稿</a-option>
          <a-option value="active">进行中</a-option>
          <a-option value="completed">已完成</a-option>
          <a-option value="archived">已归档</a-option>
        </a-select>
      </a-space>
    </div>

    <!-- 操作按钮 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="loadProjects">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
        <a-button type="primary" @click="handleCreate">
          <template #icon><icon-plus /></template>
          创建项目
        </a-button>
      </a-space>
    </div>

    <!-- 数据表格 -->
    <div class="pro-content-area">
      <a-table
        :columns="columns"
        :data="projectList"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @page-change="handlePageChange"
        @page-size-change="handlePageSizeChange"
        :scroll="{ x: 1000 }"
      >
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">
            {{ getStatusText(record.status) }}
          </a-tag>
        </template>
        <template #start_date="{ record }">
          {{ formatTime(record.start_date) }}
        </template>
        <template #end_date="{ record }">
          {{ formatTime(record.end_date) }}
        </template>
        <template #created_at="{ record }">
          {{ formatTime(record.created_at) }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click.stop="goToDetail(record.id)">详情</a-button>
            <a-button type="text" size="small" @click.stop="goToExperiments(record.id)">实验</a-button>
            <a-button type="text" size="small" @click.stop="handleEdit(record)">编辑</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 创建/编辑项目弹窗 -->
    <a-modal
      v-model:visible="modalVisible"
      :title="isEdit ? '编辑项目' : '创建项目'"
      @ok="handleSubmit"
      @cancel="modalVisible = false"
      :width="600"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="项目名称" required>
          <a-input v-model="form.name" placeholder="请输入项目名称" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" placeholder="请输入项目描述" :rows="3" />
        </a-form-item>
        <a-form-item label="开始日期">
          <a-date-picker v-model="form.start_date" style="width: 100%" />
        </a-form-item>
        <a-form-item label="关联数据集">
          <a-select v-model="form.dataset_ids" multiple placeholder="请选择数据集">
            <a-option v-for="ds in datasets" :key="ds.id" :value="ds.id">{{ ds.name }}</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 提交研究弹窗 -->
    <a-modal
      v-model:visible="submitVisible"
      title="提交研究"
      @ok="handleSubmitResearch"
      @cancel="submitVisible = false"
      :width="600"
    >
      <a-form :model="submitForm" layout="vertical">
        <a-form-item label="研究发现">
          <a-textarea v-model="submitForm.findings" placeholder="请输入研究发现" :rows="4" />
        </a-form-item>
        <a-form-item label="发表论文链接">
          <a-input v-model="submitForm.published_paper" placeholder="请输入论文链接" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'

const router = useRouter()

const columns = [
  { title: '项目名称', dataIndex: 'name', width: 200, fixed: 'left' },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '状态', dataIndex: 'status', width: 100 },
  { title: '开始日期', dataIndex: 'start_date', width: 120 },
  { title: '结束日期', dataIndex: 'end_date', width: 120 },
  { title: '创建时间', dataIndex: 'created_at', width: 120 },
  { title: '操作', dataIndex: 'actions', width: 180, fixed: 'right' },
]

const projectList = ref([])
const loading = ref(false)
const filters = reactive({
  keyword: '',
  status: '',
})
const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
})

const modalVisible = ref(false)
const isEdit = ref(false)
const form = reactive({
  id: null,
  name: '',
  description: '',
  start_date: null,
  dataset_ids: [],
})

const submitVisible = ref(false)
const submitForm = reactive({
  findings: '',
  published_paper: '',
})

const datasets = ref([])

onMounted(() => {
  loadProjects()
  loadDatasets()
})

const loadProjects = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams()
    params.append('page', pagination.current)
    params.append('page_size', pagination.pageSize)
    if (filters.keyword) params.append('keyword', filters.keyword)
    if (filters.status) params.append('status', filters.status)

    const res = await fetch(`/api/v1/research/projects?${params}`)
    const data = await res.json()
    if (data.status === 'success') {
      projectList.value = data.data
      pagination.total = data.total
    }
  } catch (error) {
    Message.error('加载项目失败')
  } finally {
    loading.value = false
  }
}

const loadDatasets = async () => {
  try {
    const res = await fetch('/api/v1/research/datasets?page_size=100')
    const data = await res.json()
    if (data.status === 'success') {
      datasets.value = data.data
    }
  } catch (error) {
    console.error('加载数据集失败', error)
  }
}

const handlePageChange = (page) => {
  pagination.current = page
  loadProjects()
}

const handlePageSizeChange = (pageSize) => {
  pagination.pageSize = pageSize
  loadProjects()
}

const goToDetail = (id) => {
  router.push(`/research/project/${id}`)
}

const goToExperiments = (id) => {
  router.push(`/research/experiments?project_id=${id}`)
}

const handleCreate = () => {
  isEdit.value = false
  Object.assign(form, {
    id: null,
    name: '',
    description: '',
    start_date: null,
    dataset_ids: [],
  })
  modalVisible.value = true
}

const handleEdit = (record) => {
  isEdit.value = true
  Object.assign(form, {
    id: record.id,
    name: record.name,
    description: record.description,
    start_date: record.start_date,
    dataset_ids: record.dataset_ids ? JSON.parse(record.dataset_ids) : [],
  })
  modalVisible.value = true
}

const handleSubmit = async () => {
  try {
    const url = isEdit.value ? `/api/v1/research/projects/${form.id}` : '/api/v1/research/projects'
    const method = isEdit.value ? 'PUT' : 'POST'
    
    const submitData = {
      ...form,
      dataset_ids: JSON.stringify(form.dataset_ids),
    }
    
    const res = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(submitData),
    })
    const data = await res.json()
    if (data.status === 'success') {
      Message.success(isEdit.value ? '更新成功' : '创建成功')
      modalVisible.value = false
      loadProjects()
    }
  } catch (error) {
    Message.error('操作失败')
  }
}

const handleSubmitResearch = async () => {
  try {
    const res = await fetch(`/api/v1/research/projects/${form.id}/submit`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(submitForm),
    })
    const data = await res.json()
    if (data.status === 'success') {
      Message.success('提交成功')
      submitVisible.value = false
      loadProjects()
    }
  } catch (error) {
    Message.error('提交失败')
  }
}

const getStatusColor = (status) => {
  const colors = { draft: 'default', active: 'blue', completed: 'green', archived: 'gray' }
  return colors[status] || 'default'
}

const getStatusText = (status) => {
  const texts = { draft: '草稿', active: '进行中', completed: '已完成', archived: '已归档' }
  return texts[status] || status
}

const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.pro-page-container {
  padding: 16px;
}
.pro-breadcrumb {
  margin-bottom: 16px;
}
.pro-search-bar {
  margin-bottom: 16px;
}
.pro-action-bar {
  margin-bottom: 16px;
}
.pro-content-area {
  background: var(--color-bg-1);
  border-radius: 4px;
}
</style>
