<template>
  <div class="container">
    <a-card class="general-card" title="AI 行为日志">
      <template #extra>
        <a-space :size="12">
          <a-button type="primary" @click="handleCreate"><icon-plus />新建</a-button>
          <a-button @click="handleSearch"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-row :gutter="16">
        <a-col :span="8">
          <a-form-item label="设备/用户">
            <a-input v-model="form.keyword" placeholder="搜索设备/用户" @pressEnter="handleSearch" />
          </a-form-item>
        </a-col>
        <a-col :span="8">
          <a-form-item label="模型版本">
            <a-input v-model="form.model_version" placeholder="请输入" />
          </a-form-item>
        </a-col>
        <a-col :span="8">
          <a-form-item label="状态">
            <a-select v-model="form.status" placeholder="请选择" allow-clear>
              <a-option value="success">成功</a-option>
              <a-option value="failed">失败</a-option>
            </a-select>
          </a-form-item>
        </a-col>
        <a-col :flex="'86px'" style="display: flex; align-items: flex-end">
          <a-space direction="vertical" :size="8">
            <a-button type="primary" @click="handleSearch">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-space>
        </a-col>
      </a-row>
      <a-divider style="margin: 0 0 16px 0" />
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" row-key="id" />
    </a-card>
      </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="名称"><a-input v-model="form.name" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const router = useRouter()
const loading = ref(false)
const data = ref([])
const modalVisible = ref(false)
const modalTitle = ref('新建')

const form = reactive({ keyword: '', model_version: '', behavior_type: '', status: '' })

const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: '时间', dataIndex: 'created_at', width: 170 },
  { title: '设备', dataIndex: 'device_id', width: 140 },
  { title: '用户', dataIndex: 'user_id', width: 100 },
  { title: '模型', dataIndex: 'model_version', width: 160 },
  { title: '事件类型', dataIndex: 'behavior_type', width: 120 },
  { title: '延迟(ms)', dataIndex: 'latency_ms', width: 90 },
  { title: '置信度', dataIndex: 'confidence', width: 90 },
  { title: '状态', dataIndex: 'status', width: 90 }
]

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { Object.keys(form).forEach(k => form[k] = ''); pagination.current = 1; loadData() }
const handleCreate = () => { modalTitle.value = '新建'; modalVisible.value = true }
const handleSubmit = () => { modalVisible.value = false; Message.success('保存成功') }
const goToDetail = (id) => { router.push(`/ai/behavior-detail/${id}`) }

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    Object.keys(form).forEach(k => { if (form[k]) params[k] = form[k] })
    const res = await fetch('/api/v1/ai/monitor/events?' + new URLSearchParams(params), {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    if (res.code === 0) { data.value = res.data?.list || []; pagination.total = res.data?.total || 0 }
  } catch (e) { Message.error('加载失败') } finally { loading.value = false }
}

onMounted(() => { loadData() })
</script>
