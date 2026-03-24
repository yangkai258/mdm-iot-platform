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

import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'

const form = reactive({
  cleanCycle: 'monthly',
  inactiveDays: 90,
  notifyBefore: '1',
  allowRestore: '1',
  protectedTags: ['vip', 'birthday'],
  logRetentionDays: 90
})

const logLoading = ref(false)
const logList = ref([
  { id: 1, time: '2026-03-01 00:00:00', tagName: '流失风险', cleanedCount: 234, status: 1 },
  { id: 2, time: '2026-02-01 00:00:00', tagName: '沉默会员', cleanedCount: 567, status: 1 },
  { id: 3, time: '2026-01-01 00:00:00', tagName: '边缘客户', cleanedCount: 123, status: 1 }
])

const logColumns = [
  { title: '清除时间', dataIndex: 'time', width: 200 },
  { title: '清除标签', dataIndex: 'tagName', width: 200 },
  { title: '清除数量', dataIndex: 'cleanedCount', width: 150 },
  { title: '状态', slotName: 'status', width: 150 }
]

const handleSave = () => {
  Message.success('保存成功')
}

const handleReset = () => {
  Object.assign(form, {
    cleanCycle: 'monthly',
    inactiveDays: 90,
    notifyBefore: '1',
    allowRestore: '1',
    protectedTags: ['vip', 'birthday'],
    logRetentionDays: 90
  })
  Message.info('已重置')
}

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
