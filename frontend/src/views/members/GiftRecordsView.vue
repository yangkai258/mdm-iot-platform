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
const dataList = ref([])
const giftOptions = ref(['新会员礼包', '生日专属礼包', '节日礼包'])
const stats = ref({ total: 0, claimed: 0, unclaimed: 0 })
const filters = reactive({ keyword: '', giftName: '', status: '', dateRange: [] })
const paginationConfig = reactive({ current: 1, pageSize: 10, total: 0 })

const columns = [
  { title: '会员名称', dataIndex: 'memberName', width: 150 },
  { title: '手机号', dataIndex: 'phone', width: 130 },
  { title: '礼包名称', dataIndex: 'giftName', width: 180 },
  { title: '礼包价值', dataIndex: 'giftValue', width: 100 },
  { title: '发放时间', dataIndex: 'grantTime', width: 180 },
  { title: '领取时间', dataIndex: 'claimTime', width: 180 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' }
]

const mockData = () => [
  { id: 1, memberName: '张三', phone: '13800138001', giftName: '新会员礼包', giftValue: 50, grantTime: '2026-03-01 10:00:00', claimTime: '2026-03-01 11:30:00', status: 'claimed' },
  { id: 2, memberName: '李四', phone: '13800138002', giftName: '生日专属礼包', giftValue: 100, grantTime: '2026-03-05 09:00:00', claimTime: '-', status: 'unclaimed' },
  { id: 3, memberName: '王五', phone: '13800138003', giftName: '节日礼包', giftValue: 80, grantTime: '2026-02-14 08:00:00', claimTime: '-', status: 'expired' }
]

const statusColor = (s) => ({ claimed: 'green', unclaimed: 'orange', expired: 'gray' }[s] || 'gray')
const statusText = (s) => ({ claimed: '已领取', unclaimed: '未领取', expired: '已过期' }[s] || s)
const showDetail = (record) => { Message.info(`查看 ${record.memberName} 的 ${record.giftName} 明细`) }

const loadData = () => {
  loading.value = true
  setTimeout(() => {
    dataList.value = mockData()
    stats.value = { total: 3, claimed: 1, unclaimed: 1 }
    paginationConfig.total = dataList.value.length
    loading.value = false
  }, 400)
}

const onPageChange = (page) => { paginationConfig.current = page; loadData() }
const handleExport = () => { Message.success('导出成功') }

loadData()

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
