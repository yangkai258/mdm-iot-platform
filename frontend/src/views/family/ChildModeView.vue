<template>
  <div class="container">
    <a-card class="general-card" title="儿童模式">
      <template #extra>
        <a-space :size="12">
          <a-button type="primary" @click="handleCreate"><icon-plus />新建</a-button>
          <a-button @click="handleSearch"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-row :gutter="16">
        <a-col :span="6">
          <a-form-item label="儿童账号">
            <a-select v-model="form.child_id" placeholder="请选择" allow-clear style="width: 100%">
              <a-option v-for="c in children" :key="c.id" :value="c.id">{{ c.name }}</a-option>
            </a-select>
          </a-form-item>
        </a-col>
        <a-col :span="6">
          <a-form-item label="状态">
            <a-select v-model="form.status" placeholder="请选择" allow-clear style="width: 100%">
              <a-option value="enabled">已启用</a-option>
              <a-option value="disabled">已禁用</a-option>
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
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id" />
    </a-card>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" :width="560">
      <a-form :model="form" label-col-flex="130px">
        <a-form-item label="儿童账号">
          <a-select v-model="form.child_id" placeholder="请选择">
            <a-option v-for="c in children" :key="c.id" :value="c.id">{{ c.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="内容过滤"><a-switch v-model="form.content_filter_enabled" /></a-form-item>
        <a-form-item label="时间限制"><a-switch v-model="form.time_limit_enabled" /></a-form-item>
        <a-form-item label="每日时长(分钟)"><a-input-number v-model="form.daily_time_limit" :min="15" :max="480" style="width: 100%" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const loading = ref(false)
const data = ref<any[]>([])
const children = ref<any[]>([])
const modalVisible = ref(false)
const editingId = ref<number | null>(null)
const form = reactive({ child_id: null as number | null, content_filter_enabled: false, time_limit_enabled: false, daily_time_limit: 60 })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const modalTitle = computed(() => editingId.value ? '编辑配置' : '新建配置')
const columns = [
  { title: '儿童账号', dataIndex: 'child_name', width: 160 },
  { title: '模式开关', dataIndex: 'enabled', width: 120 },
  { title: '内容过滤', dataIndex: 'content_filter_enabled', width: 120 },
  { title: '时间限制', dataIndex: 'time_limit', width: 160 },
  { title: '操作', slotName: 'actions', width: 160 }
]

async function loadChildren() {
  try {
    const res = await fetch('/api/v1/family/members?role=child', {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    children.value = res.data?.list || []
  } catch {}
}

async function loadData() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (form.child_id) params.append('child_id', String(form.child_id))
    if (form.status) params.append('status', form.status)
    params.append('page', String(pagination.current))
    params.append('page_size', String(pagination.pageSize))
    const res = await fetch(`/api/v1/family/child-mode?${params}`, {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    data.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch { Message.error('加载失败') } finally { loading.value = false }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { form.child_id = null; form.status = ''; pagination.current = 1; loadData() }
const handleCreate = () => { editingId.value = null; Object.assign(form, { child_id: null, content_filter_enabled: false, time_limit_enabled: false, daily_time_limit: 60 }); modalVisible.value = true }
const handleSubmit = async () => {
  try {
    const method = editingId.value ? 'PUT' : 'POST'
    const url = editingId.value ? `/api/v1/family/child-mode/${editingId.value}` : '/api/v1/family/child-mode'
    await fetch(url, { method, headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token'), 'Content-Type': 'application/json' }, body: JSON.stringify(form) })
    Message.success('保存成功')
    modalVisible.value = false
    loadData()
  } catch { Message.error('保存失败') }
}
const onPageChange = (page: number) => { pagination.current = page; loadData() }

onMounted(() => { loadChildren(); loadData() })
</script>
