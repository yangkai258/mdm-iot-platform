<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-filter /> 告警抑制规则</a-space>
      </template>
      <template #extra>
        <a-button type="primary" @click="handleCreate">
          <template #icon><icon-plus /></template>
          新建规则
        </a-button>
      </template>

      <a-table :columns="columns" :data="tableData">
        <template #enabled="{ record }">
          <a-switch v-model="record.enabled" />
        </template>
        <template #mergeMode="{ record }">
          <a-tag>{{ { count: '计数合并', summary: '摘要合并', latest: '保留最新' }[record.mergeMode] }}</a-tag>
        </template>
        <template #hitCount="{ record }">
          <a-statistic :value="record.hitCount" :precision="0" />
        </template>
        <template #actions="{ record }">
          <a-link @click="handleEdit(record)">编辑</a-link>
          <a-link @click="handleCopy(record)">复制</a-link>
          <a-link @click="handleDelete(record)">删除</a-link>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="modalVisible" :title="modalTitle" @ok="handleSubmit" width="600px">
      <a-form :model="formData" layout="vertical">
        <a-form-item label="规则名称" required>
          <a-input v-model="formData.ruleName" placeholder="请输入规则名称" />
        </a-form-item>
        <a-form-item label="去重关键字段">
          <a-checkbox-group v-model="formData.dedupFields">
            <a-checkbox value="alert_type">告警类型</a-checkbox>
            <a-checkbox value="device_id">设备ID</a-checkbox>
            <a-checkbox value="metric">指标</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
        <a-form-item label="抑制时长(分钟)">
          <a-input-number v-model="formData.suppressMinutes" :min="1" :max="1440" />
        </a-form-item>
        <a-form-item label="合并模式">
          <a-radio-group v-model="formData.mergeMode">
            <a-radio value="count">计数合并</a-radio>
            <a-radio value="summary">摘要合并</a-radio>
            <a-radio value="latest">保留最新</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const modalVisible = ref(false)
const modalTitle = ref('新建规则')
const formData = reactive({
  ruleName: '', dedupFields: ['alert_type', 'device_id'], suppressMinutes: 30, mergeMode: 'count'
})
const columns = [
  { title: '规则名称', dataIndex: 'ruleName' },
  { title: '去重字段', dataIndex: 'dedupFields' },
  { title: '抑制时长', dataIndex: 'suppressMinutes' },
  { title: '合并模式', slotName: 'mergeMode' },
  { title: '命中次数', slotName: 'hitCount' },
  { title: '启用', slotName: 'enabled' },
  { title: '操作', slotName: 'actions', width: 180 }
]
const tableData = ref([
  { id: 1, ruleName: '同设备同类型告警5分钟内合并', dedupFields: 'alert_type, device_id', suppressMinutes: 5, mergeMode: 'latest', hitCount: 156, enabled: true },
  { id: 2, ruleName: '温度告警30分钟内抑制', dedupFields: 'alert_type, device_id, metric', suppressMinutes: 30, mergeMode: 'count', hitCount: 89, enabled: true }
])

const handleCreate = () => { modalVisible.value = true; modalTitle.value = '新建规则' }
const handleEdit = (r) => { modalVisible.value = true; modalTitle.value = '编辑规则' }
const handleCopy = (r) => { }
const handleDelete = (r) => { }
const handleSubmit = () => { modalVisible.value = false }
</script>

<style scoped>
.container { padding: 16px; }
</style>
