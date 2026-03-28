<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-bookmark /> 设备标签管理</a-space>
      </template>
      <template #extra>
        <a-space>
          <a-button type="primary" @click="handleCreateTag">
            <template #icon><icon-plus /></template>
            新建标签
          </a-button>
          <a-button type="primary" status="warning" @click="handleBatchTag">
            <template #icon><icon-tag /></template>
            批量打标
          </a-button>
        </a-space>
      </template>

      <a-row :gutter="16">
        <a-col :span="8">
          <a-card title="标签云" size="small">
            <a-tag v-for="tag in tags" :key="tag.key" :color="tag.color" : closable @close="handleDeleteTag(tag)">
              {{ tag.name }} ({{ tag.count }})
            </a-tag>
            <a-button type="text" @click="handleCreateTag">
              <template #icon><icon-plus /></template>
              添加标签
            </a-button>
          </a-card>
        </a-col>
        <a-col :span="16">
          <a-form layout="inline">
            <a-form-item label="标签筛选">
              <a-select v-model="filterTag" placeholder="选择标签" allow-clear style="width: 200px">
                <a-option v-for="t in tags" :key="t.key" :value="t.key">{{ t.name }}</a-option>
              </a-select>
            </a-form-item>
            <a-form-item label="设备ID">
              <a-input v-model="filterDeviceId" placeholder="输入设备ID" />
            </a-form-item>
            <a-form-item>
              <a-button type="primary">查询</a-button>
            </a-form-item>
          </a-form>

          <a-table :columns="columns" :data="tableData" :pagination="pagination">
            <template #tags="{ record }">
              <a-tag v-for="t in record.tags" :key="t" :color="getTagColor(t)">{{ t }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-link @click="handleTagDevice(record)">打标</a-link>
              <a-link @click="handleViewDevice(record)">详情</a-link>
            </template>
          </a-table>
        </a-col>
      </a-row>
    </a-card>

    <a-modal v-model:visible="tagModalVisible" title="新建标签" @ok="handleTagSubmit">
      <a-form :model="tagForm" layout="vertical">
        <a-form-item label="标签名称" required>
          <a-input v-model="tagForm.name" placeholder="请输入标签名称" />
        </a-form-item>
        <a-form-item label="标签颜色">
          <a-color-picker v-model="tagForm.color" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const tags = ref([
  { key: 'outdoor', name: '户外', color: 'green', count: 12 },
  { key: 'indoor', name: '室内', color: 'blue', count: 8 },
  { key: 'test', name: '测试', color: 'orange', count: 5 },
  { key: 'prod', name: '生产', color: 'red', count: 15 },
  { key: 'vip', name: 'VIP', color: 'purple', count: 3 }
])
const columns = [
  { title: '设备ID', dataIndex: 'deviceId' },
  { title: '设备名称', dataIndex: 'deviceName' },
  { title: '标签', slotName: 'tags' },
  { title: '在线状态', dataIndex: 'status' },
  { title: '最后活跃', dataIndex: 'lastActive' },
  { title: '操作', slotName: 'actions', width: 150 }
]
const tableData = ref([
  { deviceId: 'D001', deviceName: '设备-001', tags: ['outdoor', 'vip'], status: '在线', lastActive: '2026-03-28 10:00' },
  { deviceId: 'D002', deviceName: '设备-002', tags: ['indoor', 'prod'], status: '离线', lastActive: '2026-03-27 15:30' }
])
const filterTag = ref('')
const filterDeviceId = ref('')
const pagination = reactive({ current: 1, pageSize: 10, total: 100 })
const tagModalVisible = ref(false)
const tagForm = reactive({ name: '', color: '#1890ff' })

const getTagColor = (tagKey) => tags.value.find(t => t.key === tagKey)?.color || 'gray'
const handleCreateTag = () => { tagModalVisible.value = true }
const handleDeleteTag = (tag) => { }
const handleBatchTag = () => { }
const handleTagDevice = (record) => { }
const handleViewDevice = (record) => { }
const handleTagSubmit = () => { tagModalVisible.value = false }
</script>

<style scoped>
.container { padding: 16px; }
</style>
