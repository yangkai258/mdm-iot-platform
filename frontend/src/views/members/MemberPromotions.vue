<template>
  <div class="page-container">
    <Breadcrumb :items="['menu.member', 'menu.member.promotions']" />
    <a-card class="general-card" title="营销活动">
      <template #extra>
        <a-button type="primary" @click="handleCreate"><icon-plus />新建</a-button>
      </template>
      <div class="search-form">
        <a-form :model="filters" layout="inline">
          <a-form-item label="关键词"><a-input v-model="filters.keyword" placeholder="请输入关键词" /></a-form-item>
          <a-form-item>
            <a-button type="primary" @click="handleSearch">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="promotionList" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #type="{ record }">
          <a-tag :color="getTypeColor(record.type)">{{ getTypeText(record.type) }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
          <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
        </template>
      </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @before-ok="handleSubmit" @cancel="modalVisible = false">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="活动名称"><a-input v-model="form.name" placeholder="请输入活动名称" /></a-form-item>
        <a-form-item label="活动类型">
          <a-select v-model="form.type" placeholder="请选择">
            <a-option value="points_double">双倍积分</a-option>
            <a-option value="discount">折扣活动</a-option>
            <a-option value="gift">赠品活动</a-option>
          </a-select>
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit(() => {})">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>

import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const promotionList = ref([])
const showCreateDrawer = ref(false)
const modalVisible = ref(false)
const modalTitle = ref('新建活动')
const isEdit = ref(false)
const currentId = ref(null)

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
const getTypeColor = (type) => ({ points_double: 'blue', discount: 'orange', gift: 'purple' }[type] || 'gray')
const getStatusText = (s) => ({ active: '进行中', ended: '已结束', draft: '草稿' }[s] || s)
const getStatusColor = (s) => ({ active: 'green', ended: 'gray', draft: 'orange' }[s] || 'gray')

const loadPromotions = () => {
  loading.value = true
  setTimeout(() => {
    const mock = [
      { id: 1, name: '国庆双倍积分', type: 'points_double', start_time: '2026-10-01', end_time: '2026-10-07', status: 'ended' },
      { id: 2, name: '会员日折扣', type: 'discount', start_time: '2026-03-15', end_time: '2026-03-31', status: 'active' }
    ]
    promotionList.value = mock.filter(item => {
      if (filters.keyword && !item.name.includes(filters.keyword)) return false
      return true
    })
    pagination.total = promotionList.value.length
    loading.value = false
  }, 300)
}

const handleSearch = () => { pagination.current = 1; loadPromotions() }
const handleReset = () => { Object.assign(filters, { status: undefined, keyword: '' }); pagination.current = 1; loadPromotions() }
const handleCreate = () => { isEdit.value = false; currentId.value = null; modalTitle.value = '新建活动'; Object.assign(form, { name: '', type: '' }); modalVisible.value = true }
const handleEdit = (record) => { isEdit.value = true; currentId.value = record.id; modalTitle.value = '编辑活动'; Object.assign(form, record); modalVisible.value = true }
const handleDelete = (record) => { Message.success('删除成功'); loadPromotions() }
const handleSubmit = (done) => {
  if (!form.name) { Message.error('请输入活动名称'); done && done(false); return }
  setTimeout(() => { Message.success(isEdit.value ? '更新成功' : '创建成功'); modalVisible.value = false; loadPromotions(); done && done(true) }, 400)
}
const onPageChange = (page) => { pagination.current = page; loadPromotions() }

onMounted(() => loadPromotions())

</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>
