<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="名称"><a-input v-model="form.name" placeholder="请输入" /></a-form-item>
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
      <template #level="{ record }">
        <a-tag :color="record.level === 'diamond' ? 'purple' : record.level === 'platinum' ? 'gold' : record.level === 'gold' ? 'orange' : 'gray'">
          {{ record.levelName }}
        </a-tag>
      </template>
      <template #status="{ record }">
        <a-tag :color="record.status === '1' ? 'green' : 'gray'">{{ record.status === '1' ? '启用' : '禁用' }}</a-tag>
      </template>
      <template #actions="{ record }">
        <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
        <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
      </template>
    </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @before-ok="handleSubmit" @cancel="modalVisible = false">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="权益名称"><a-input v-model="form.name" placeholder="请输入权益名称" /></a-form-item>
        <a-form-item label="适用等级">
          <a-select v-model="form.level" placeholder="请选择">
            <a-option value="silver">银卡</a-option>
            <a-option value="gold">金卡</a-option>
            <a-option value="platinum">铂金卡</a-option>
            <a-option value="diamond">钻石卡</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="权益内容"><a-textarea v-model="form.content" placeholder="请输入权益内容" :rows="3" /></a-form-item>
        <a-form-item label="状态">
          <a-radio-group v-model="form.status">
            <a-radio value="1">启用</a-radio>
            <a-radio value="0">禁用</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const data = ref([])
const modalVisible = ref(false)
const modalTitle = ref('新建权益')
const isEdit = ref(false)
const currentId = ref(null)

const form = reactive({ id: null, name: '', level: 'silver', content: '', status: '1' })

const pagination = reactive({ current: 1, pageSize: 10, total: 0 })

const columns = [
  { title: '权益名称', dataIndex: 'name', width: 180 },
  { title: '适用等级', slotName: 'level', width: 100 },
  { title: '权益内容', dataIndex: 'content', ellipsis: true },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' }
]

const mockData = () => [
  { id: 1, name: '生日专属折扣', level: 'silver', levelName: '银卡', content: '生日当月享受9折优惠', status: '1' },
  { id: 2, name: '免费停车', level: 'gold', levelName: '金卡', content: '每月免费停车2小时', status: '1' },
  { id: 3, name: 'VIP客服', level: 'platinum', levelName: '铂金卡', content: '专属客服热线优先接入', status: '1' },
  { id: 4, name: '免费洗车', level: 'diamond', levelName: '钻石卡', content: '每月免费洗车1次', status: '1' }
]

const loadData = () => {
  loading.value = true
  setTimeout(() => {
    const all = mockData()
    data.value = all.filter(item => {
      if (form.name && !item.name.includes(form.name)) return false
      return true
    })
    pagination.total = data.value.length
    loading.value = false
  }, 300)
}

const handleSearch = () => {
  pagination.current = 1
  loadData()
}

const handleReset = () => {
  Object.assign(form, { id: null, name: '', level: 'silver', content: '', status: '1' })
  pagination.current = 1
  loadData()
}

const handleCreate = () => {
  isEdit.value = false
  currentId.value = null
  modalTitle.value = '新建权益'
  Object.assign(form, { id: null, name: '', level: 'silver', content: '', status: '1' })
  modalVisible.value = true
}

const handleEdit = (record) => {
  isEdit.value = true
  currentId.value = record.id
  modalTitle.value = '编辑权益'
  Object.assign(form, record)
  modalVisible.value = true
}

const handleDelete = (record) => {
  Message.success('删除成功')
  loadData()
}

const handleSubmit = (done) => {
  if (!form.name) {
    Message.error('请输入权益名称')
    done(false)
    return
  }
  setTimeout(() => {
    Message.success(isEdit.value ? '更新成功' : '创建成功')
    modalVisible.value = false
    loadData()
    done(true)
  }, 400)
}

const onPageChange = (page) => {
  pagination.current = page
  loadData()
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
