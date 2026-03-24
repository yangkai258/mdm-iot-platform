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

import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)
const dataList = ref([])
const stats = ref({ total: 0, issued: 0, claimed: 0, totalValue: 0 })
const filters = reactive({ keyword: '', status: '' })
const paginationConfig = reactive({ current: 1, pageSize: 10, total: 0 })
const form = reactive({ id: null, name: '', content: '', value: 0, validDays: 30, status: '1' })

const columns = [
  { title: '礼包名称', dataIndex: 'name', width: 200 },
  { title: '包含内容', dataIndex: 'content', ellipsis: true },
  { title: '价值（元）', dataIndex: 'value', width: 120 },
  { title: '已发放数', dataIndex: 'issuedCount', width: 120 },
  { title: '已领取数', dataIndex: 'claimedCount', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' }
]

const mockData = () => [
  { id: 1, name: '新会员礼包', content: '新人优惠券满100减20、100积分', value: 50, issuedCount: 234, claimedCount: 210, status: 1 },
  { id: 2, name: '生日专属礼包', content: '生日优惠券满200减50、双倍积分卡', value: 100, issuedCount: 567, claimedCount: 520, status: 1 },
  { id: 3, name: '节日礼包', content: '节日专属优惠券包', value: 80, issuedCount: 1234, claimedCount: 1100, status: 1 }
]

const loadData = () => {
  loading.value = true
  setTimeout(() => {
    dataList.value = mockData()
    const total = mockData()
    stats.value = {
      total: total.length,
      issued: total.reduce((s, d) => s + d.issuedCount, 0),
      claimed: total.reduce((s, d) => s + d.claimedCount, 0),
      totalValue: total.reduce((s, d) => s + d.value * d.issuedCount, 0)
    }
    paginationConfig.total = dataList.value.length
    loading.value = false
  }, 400)
}

const onPageChange = (page) => { paginationConfig.current = page; loadData() }
const showCreate = () => { isEdit.value = false; Object.assign(form, { id: null, name: '', content: '', value: 0, validDays: 30, status: '1' }); formVisible.value = true }
const showEdit = (record) => { isEdit.value = true; Object.assign(form, { ...record, status: String(record.status) }); formVisible.value = true }
const handleSubmit = (done) => {
  if (!form.name) { Message.error('请填写必填项'); done(false); return }
  setTimeout(() => { Message.success(isEdit.value ? '更新成功' : '创建成功'); formVisible.value = false; loadData(); done(true) }, 400)
}
const handleGrant = (record) => { Message.success(`已向符合条件的会员发放「${record.name}」`) }
const handleDelete = () => { Message.success('删除成功'); loadData() }

loadData()

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
