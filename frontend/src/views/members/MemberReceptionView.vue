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
const detailVisible = ref(false)
const detailRecord = ref({})
const dataList = ref([])
const stats = ref({ today: 5, week: 32, month: 128 })
const filters = reactive({ keyword: '', dateRange: [] })
const paginationConfig = reactive({ current: 1, pageSize: 10, total: 0 })

const columns = [
  { title: '接待记录', dataIndex: 'recordNo', width: 180 },
  { title: '会员名称', dataIndex: 'memberName', width: 150 },
  { title: '接待时间', dataIndex: 'receptionTime', width: 180 },
  { title: '接待人', dataIndex: 'receptionist', width: 120 },
  { title: '接待类型', slotName: 'type', width: 120 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' }
]

const typeColor = (t) => ({ '到店接待': 'blue', '电话回访': 'green', '线上咨询': 'purple', '活动接待': 'orange' }[t] || 'gray')

const mockData = () => [
  { id: 1, recordNo: 'RC20260324001', memberName: '张三', receptionTime: '2026-03-24 10:00:00', receptionist: '客服小李', type: '到店接待', content: 'VIP会员到店咨询，介绍新品', remark: '意向强烈' },
  { id: 2, recordNo: 'RC20260324002', memberName: '李四', receptionTime: '2026-03-24 11:30:00', receptionist: '客服小王', type: '电话回访', content: '生日会员电话回访', remark: '满意' },
  { id: 3, recordNo: 'RC20260324003', memberName: '王五', receptionTime: '2026-03-24 14:00:00', receptionist: '客服小张', type: '线上咨询', content: '线上咨询会员权益', remark: '-' }
]

const loadData = () => {
  loading.value = true
  setTimeout(() => {
    dataList.value = mockData()
    paginationConfig.total = dataList.value.length
    loading.value = false
  }, 400)
}

const onPageChange = (page) => { paginationConfig.current = page; loadData() }
const showDetail = (record) => { detailRecord.value = record; detailVisible.value = true }

loadData()

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
