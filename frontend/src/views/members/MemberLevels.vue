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
    <a-table :columns="columns" :data="levelList" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
      <template #actions="{ record }">
        <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
        <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
      </template>
    </a-table>
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
import { Message, Modal } from '@arco-design/web-vue'
import * as api from '@/api/member'

const levelList = ref([])
const loading = ref(false)
const formLoading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)

const form = reactive({
  name: '',
  code: '',
  minAmount: 0,
  maxAmount: undefined,
  discountRate: 1,
  pointsRate: 1,
  benefits: [],
  description: '',
  sort: 0
})

const columns = [
  { title: '等级', slotName: 'levelColor', width: 120 },
  { title: '编码', dataIndex: 'code', width: 100 },
  { title: '折扣率', slotName: 'discountRate', width: 100 },
  { title: '积分倍率', slotName: 'pointsRate', width: 110 },
  { title: '最低门槛', slotName: 'minAmount', width: 130 },
  { title: '最高门槛', slotName: 'maxAmount', width: 130 },
  { title: '等级权益', slotName: 'benefits' },
  { title: '会员数', dataIndex: 'memberCount', width: 90 },
  { title: '操作', slotName: 'actions', width: 140 }
]

const getLevelColor = (id) => {
  const colors = { 1: '#95de64', 2: '#1890ff', 3: '#faad14', 4: '#ff4d4f', 5: '#722ed1' }
  return colors[id] || 'gray'
}

const loadLevels = async () => {
  loading.value = true
  try {
    const res = await api.getLevelList()
    levelList.value = res.data || []
  } catch (err) {
    Message.error('加载等级列表失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

const showCreate = () => {
  isEdit.value = false
  Object.assign(form, {
    name: '', code: '', minAmount: 0, maxAmount: undefined,
    discountRate: 1, pointsRate: 1, benefits: [], description: '', sort: 0
  })
  formVisible.value = true
}

const showEdit = (record) => {
  isEdit.value = true
  currentId.value = record.id
  Object.assign(form, {
    name: record.name,
    code: record.code,
    minAmount: record.minAmount || 0,
    maxAmount: record.maxAmount || undefined,
    discountRate: record.discountRate || 1,
    pointsRate: record.pointsRate || 1,
    benefits: record.benefits || [],
    description: record.description || '',
    sort: record.sort || 0
  })
  formVisible.value = true
}

const handleFormSubmit = async (done) => {
  if (!form.name || !form.code) {
    Message.warning('请填写等级名称和编码')
    done(false)
    return
  }
  formLoading.value = true
  try {
    const payload = { ...form }
    if (isEdit.value) {
      await api.updateLevel(currentId.value, payload)
      Message.success('更新成功')
    } else {
      await api.createLevel(payload)
      Message.success('创建成功')
    }
    formVisible.value = false
    loadLevels()
    done(true)
  } catch (err) {
    Message.error(err.message || '操作失败')
    done(false)
  } finally {
    formLoading.value = false
  }
}

const handleDelete = (record) => {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除等级「${record.name}」吗？该等级下的会员将变为无等级状态。`,
    okText: '确认删除',
    onOk: async () => {
      try {
        await api.deleteLevel(record.id)
        Message.success('删除成功')
        loadLevels()
      } catch (err) {
        Message.error(err.message || '删除失败')
      }
    }
  })
}

onMounted(() => loadLevels())

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
