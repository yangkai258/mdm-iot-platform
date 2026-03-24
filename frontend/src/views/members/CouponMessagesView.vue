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
    <a-table :columns="columns" :data="dataList" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
      <template #status="{ record }">
        <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '启用' : '禁用' }}</a-tag>
      </template>
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

import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as api from '@/api/member'

const dataList = ref([])
const loading = ref(false)

const filters = reactive({ type: undefined, memberName: '', status: undefined, dateRange: [] })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current, pageSize: pagination.pageSize, total: pagination.total,
  showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50, 100]
}))

const columns = [
  { title: '消息类型', slotName: 'type', width: 120 },
  { title: '会员名称', dataIndex: 'memberName', width: 150 },
  { title: '优惠券名称', dataIndex: 'couponName', width: 160 },
  { title: '消息内容', dataIndex: 'content', width: 280, ellipsis: true },
  { title: '发送时间', dataIndex: 'sendTime', width: 170 },
  { title: '状态', slotName: 'status', width: 100 }
]

const getTypeColor = (t) => ({ grant: 'blue', remind: 'orange', use: 'green', expire: 'purple' }[t] || 'gray')
const getTypeText = (t) => ({ grant: '发放通知', remind: '到期提醒', use: '核销通知', expire: '过期通知' }[t] || t)
const getStatusColor = (s) => ({ success: 'green', failed: 'red', pending: 'orange' }[s] || 'gray')
const getStatusText = (s) => ({ success: '已发送', failed: '发送失败', pending: '待发送' }[s] || s)

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.type) params.type = filters.type
    if (filters.memberName) params.memberName = filters.memberName
    if (filters.status) params.status = filters.status
    if (filters.dateRange && filters.dateRange.length === 2) {
      params.startDate = filters.dateRange[0]
      params.endDate = filters.dateRange[1]
    }
    const res = await api.getCouponMessageList(params)
    const d = res.data || {}
    dataList.value = d.list || []
    pagination.total = d.total || 0
  } catch (err) {
    Message.error('加载消息流水失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

const handleExport = () => {
  Message.info('导出功能开发中')
}

const onPageChange = (page) => { pagination.current = page; loadData() }
const onPageSizeChange = (pageSize) => { pagination.pageSize = pageSize; pagination.current = 1; loadData() }

onMounted(() => loadData())

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
