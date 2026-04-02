<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="家庭成员">
      <template #extra>
        <a-space :size="12">
          <a-button type="primary" @click="handleInvite"><icon-plus />邀请成员</a-button>
          <a-button @click="loadData"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-row :gutter="16">
        <a-col :span="6">
          <a-form-item label="成员姓名">
            <a-input v-model="form.keyword" placeholder="请输入" @pressEnter="handleSearch" style="width: 100%" />
          </a-form-item>
        </a-col>
        <a-col :span="6">
          <a-form-item label="成员角色">
            <a-select v-model="form.role" placeholder="请选择" allow-clear style="width: 100%">
              <a-option value="owner">户主</a-option>
              <a-option value="adult">成人</a-option>
              <a-option value="child">儿童</a-option>
              <a-option value="elder">老人</a-option>
              <a-option value="guest">访客</a-option>
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
    <a-modal v-model:visible="modalVisible" :title="modalTitle" :width="480">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item v-if="!editingId" label="手机号码">
          <a-input v-model="form.phone" placeholder="请输入" />
        </a-form-item>
        <a-form-item label="成员角色">
          <a-select v-model="form.role" placeholder="请选择">
            <a-option value="adult">成人</a-option>
            <a-option value="child">儿童</a-option>
            <a-option value="elder">老人</a-option>
            <a-option value="guest">访客</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="邀请留言">
          <a-textarea v-model="form.message" :rows="3" />
        </a-form-item>
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
const modalVisible = ref(false)
const editingId = ref<number | null>(null)

const form = reactive({
  phone: '',
  role: 'adult',
  message: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const modalTitle = computed(() => editingId.value ? '编辑角色' : '邀请成员')

const columns = [
  { title: '姓名', dataIndex: 'name', width: 140 },
  { title: '手机号', dataIndex: 'phone', width: 140 },
  { title: '角色', dataIndex: 'role', width: 120 },
  { title: '状态', dataIndex: 'status', width: 100 },
  { title: '加入时间', dataIndex: 'joined_at', width: 180 },
  { title: '操作', slotName: 'actions', width: 180 }
]

async function loadData() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (form.keyword) params.append('keyword', form.keyword)
    if (form.role) params.append('role', form.role)
    if (form.status) params.append('status', form.status)
    params.append('page', String(pagination.current))
    params.append('page_size', String(pagination.pageSize))

    const res = await fetch(`/api/v1/family/members?${params}`)
    const json = await res.json()
    data.value = json.data?.list || []
    pagination.total = json.data?.total || 0
  } catch {
    Message.error('加载失败')
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.current = 1
  loadData()
}

function handleReset() {
  form.keyword = ''
  form.role = ''
  form.status = ''
  pagination.current = 1
  loadData()
}

function handleInvite() {
  editingId.value = null
  form.phone = ''
  form.role = 'adult'
  form.message = ''
  modalVisible.value = true
}

async function handleSubmit() {
  try {
    if (!editingId.value) {
      await fetch('/api/v1/family/members/invite', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(form)
      })
      Message.success('邀请已发送')
    } else {
      await fetch(`/api/v1/family/members/${editingId.value}/role`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ role: form.role })
      })
      Message.success('角色更新成功')
    }
    modalVisible.value = false
    loadData()
  } catch {
    Message.error('操作失败')
  }
}

function onPageChange(page: number) {
  pagination.current = page
  loadData()
}

onMounted(() => loadData())
</script>

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
</style>
