<template>
  <div class="page-container">
    <a-card class="general-card" title="会员积分">
      <template #extra>
        <a-button type="primary" @click="handleCreate"><icon-plus />新建</a-button>
      </template>
      <div class="search-form">
        <a-form :model="filters" layout="inline">
          <a-form-item label="关键词"><a-input v-model="filters.keyword" placeholder="请输入" /></a-form-item>
          <a-form-item>
            <a-button type="primary" @click="loadMembers">查询</a-button>
            <a-button @click="filters.keyword = ''; loadMembers()">重置</a-button>
          </a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="memberList" :loading="loading" :pagination="pagination" row-key="id">
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
          <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
        </template>
      </a-table>
    </a-card>
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

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const memberList = ref([])
const modalVisible = ref(false)
const modalTitle = ref('新建')
const isEdit = ref(false)
const currentId = ref(null)

const filters = reactive({ keyword: '' })
const form = reactive({ name: '' })

const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: '会员名称', dataIndex: 'member_name', width: 150 },
  { title: '手机号', dataIndex: 'phone', width: 130 },
  { title: '当前积分', dataIndex: 'points', width: 120 },
  { title: '注册时间', dataIndex: 'created_at', width: 160 },
  { title: '操作', slotName: 'actions', width: 160 }
]

const loadMembers = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (filters.keyword) params.keyword = filters.keyword
    const res = await fetch(`/api/v1/members/points?${new URLSearchParams(params)}`, {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    memberList.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch { memberList.value = [] } finally { loading.value = false }
}

const handleCreate = () => { isEdit.value = false; modalTitle.value = '新建'; Object.assign(form, { name: '' }); modalVisible.value = true }
const handleEdit = (record) => { isEdit.value = true; modalTitle.value = '编辑'; Object.assign(form, record); modalVisible.value = true }
const handleSubmit = () => { modalVisible.value = false; Message.success(isEdit.value ? '更新成功' : '创建成功'); loadMembers() }
const handleDelete = () => { Message.success('删除成功'); loadMembers() }
const handlePageChange = (page) => { pagination.current = page; loadMembers() }

onMounted(() => loadMembers())
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>

