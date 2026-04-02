<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="老人陪伴模式">
      <template #extra>
        <a-space :size="12">
          <a-button type="primary" @click="handleCreate"><icon-plus />新建</a-button>
          <a-button @click="handleSearch"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-row :gutter="16">
        <a-col :span="6">
          <a-form-item label="老人账号">
            <a-select v-model="form.elder_id" placeholder="请选择" allow-clear style="width: 100%">
              <a-option v-for="e in elders" :key="e.id" :value="e.id">{{ e.name }}</a-option>
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
      </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" :width="560">
      <a-form :model="form" label-col-flex="130px">
        <a-form-item label="老人账号">
          <a-select v-model="form.elder_id" placeholder="请选择">
            <a-option v-for="e in elders" :key="e.id" :value="e.id">{{ e.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="简化界面"><a-switch v-model="form.simplified_ui" /></a-form-item>
        <a-form-item label="字体大小">
          <a-select v-model="form.font_size" style="width: 100%">
            <a-option value="large">大（1.2倍）</a-option>
            <a-option value="xlarge">超大（1.5倍）</a-option>
            <a-option value="xxlarge">极大（2倍）</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="高对比度"><a-switch v-model="form.high_contrast" /></a-form-item>
        <a-form-item label="语音播报"><a-switch v-model="form.voice_announce" /></a-form-item>
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
const elders = ref<any[]>([])
const modalVisible = ref(false)
const editingId = ref<number | null>(null)
const form = reactive({ elder_id: null as number | null, simplified_ui: true, font_size: 'large', high_contrast: false, voice_announce: false })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const modalTitle = computed(() => editingId.value ? '编辑配置' : '新建配置')
const columns = [
  { title: '老人账号', dataIndex: 'elder_name', width: 160 },
  { title: '模式开关', dataIndex: 'enabled', width: 120 },
  { title: '简化界面', dataIndex: 'simplified_ui', width: 120 },
  { title: '字体大小', dataIndex: 'font_size', width: 140 },
  { title: '高对比度', dataIndex: 'high_contrast', width: 120 },
  { title: '操作', slotName: 'actions', width: 160 }
]

async function loadElders() {
  try {
    const res = await fetch('/api/v1/family/members?role=elder', {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    elders.value = res.data?.list || []
  } catch {}
}

async function loadData() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (form.elder_id) params.append('elder_id', String(form.elder_id))
    if (form.status) params.append('status', form.status)
    params.append('page', String(pagination.current))
    params.append('page_size', String(pagination.pageSize))
    const res = await fetch(`/api/v1/family/elder-mode?${params}`, {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    data.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch { Message.error('加载失败') } finally { loading.value = false }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { form.elder_id = null; form.status = ''; pagination.current = 1; loadData() }
const handleCreate = () => { editingId.value = null; Object.assign(form, { elder_id: null, simplified_ui: true, font_size: 'large', high_contrast: false, voice_announce: false }); modalVisible.value = true }
const handleSubmit = async () => {
  try {
    const method = editingId.value ? 'PUT' : 'POST'
    const url = editingId.value ? `/api/v1/family/elder-mode/${editingId.value}` : '/api/v1/family/elder-mode'
    await fetch(url, { method, headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token'), 'Content-Type': 'application/json' }, body: JSON.stringify(form) })
    Message.success('保存成功')
    modalVisible.value = false
    loadData()
  } catch { Message.error('保存失败') }
}
const onPageChange = (page: number) => { pagination.current = page; loadData() }

onMounted(() => { loadElders(); loadData() })
</script>
