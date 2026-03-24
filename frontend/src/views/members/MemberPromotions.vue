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
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const promotionList = ref([])
const showCreateDrawer = ref(false)

const filters = reactive({ status: undefined, keyword: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const stats = reactive({ total: 0, running: 0, ended: 0, participants: 0 })
const form = reactive({ name: '', type: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '活动名称', dataIndex: 'name' },
  { title: '类型', slotName: 'type', width: 120 },
  { title: '开始时间', dataIndex: 'start_time', width: 160 },
  { title: '结束时间', dataIndex: 'end_time', width: 160 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 160 }
]

const getTypeText = (type) => ({ points_double: '双倍积分', discount: '折扣活动', gift: '赠品活动' }[type] || type)
const getStatusText = (s) => ({ active: '进行中', ended: '已结束', draft: '草稿' }[s] || s)
const getStatusColor = (s) => ({ active: 'green', ended: 'gray', draft: 'orange' }[s] || 'gray')

const loadPromotions = () => { loading.value = true; setTimeout(() => { loading.value = false }, 300) }
const viewDetail = (r) => Message.info('查看详情')
const editPromotion = (r) => Message.info('编辑活动')
const handleCreate = () => { Message.success('创建成功'); showCreateDrawer.value = false }

onMounted(() => loadPromotions())

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
