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
const dataList = ref([])
const stats = ref({ total: 12, monthNew: 3, monthClean: 1, active: 11 })
const filters = reactive({ tagType: '', dateRange: [] })
const paginationConfig = reactive({ current: 1, pageSize: 10, total: 0 })
const pieChartRef = ref(null)
const lineChartRef = ref(null)

const columns = [
  { title: '标签名称', dataIndex: 'name', width: 180 },
  { title: '标签类型', dataIndex: 'type', width: 150 },
  { title: '当前会员数', dataIndex: 'memberCount', width: 130 },
  { title: '上月会员数', dataIndex: 'lastMonthCount', width: 130 },
  { title: '环比变化', slotName: 'trend', width: 120 },
  { title: '创建时间', dataIndex: 'createTime', width: 180 },
  { title: '状态', dataIndex: 'status', width: 100 }
]

const mockData = () => [
  { id: 1, name: '月度活跃买家', type: '高频购买', memberCount: 1234, lastMonthCount: 1100, trend: 134, createTime: '2026-01-01', status: '启用' },
  { id: 2, name: '沉睡会员', type: '低频购买', memberCount: 3456, lastMonthCount: 3200, trend: 256, createTime: '2026-01-01', status: '启用' },
  { id: 3, name: '美食爱好者', type: '兴趣分类', memberCount: 2345, lastMonthCount: 2100, trend: 245, createTime: '2026-02-01', status: '启用' }
]

const loadData = () => {
  loading.value = true
  setTimeout(() => {
    dataList.value = mockData()
    paginationConfig.total = dataList.value.length
    loading.value = false
    renderCharts()
  }, 400)
}

const renderCharts = () => {
  // simple visual placeholder - in production would usearco charts
}

const onPageChange = (page) => { paginationConfig.current = page; loadData() }

const handleExport = () => {
  Message.success('导出成功')
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
