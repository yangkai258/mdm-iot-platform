<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space>
          <icon-layers /> 设备分组管理
        </a-space>
      </template>
      <template #extra>
        <a-space>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新建分组
          </a-button>
          <a-button @click="handleExport">
            <template #icon><icon-download /></template>
            导出
          </a-button>
        </a-space>
      </template>

      <a-row :gutter="16">
        <a-col :span="6">
          <a-card title="分组列表" size="small">
            <a-tree
              :data="treeData"
              :show-line="true"
              @select="handleSelect"
            >
              <template #extra="node">
                <a-space size="mini">
                  <a-link @click="handleEdit(node)"><icon-edit /></a-link>
                  <a-link @click="handleDelete(node)"><icon-delete /></a-link>
                </a-space>
              </template>
            </a-tree>
          </a-card>
        </a-col>
        <a-col :span="18">
          <a-table :columns="columns" :data="tableData" :pagination="false">
            <template #status="{ record }}">
              <a-tag :color="record.status === 'online' ? 'green' : 'gray'">
                {{ record.status === 'online' ? '在线' : '离线' }}
              </a-tag>
            </template>
            <template #actions="{ record }">
              <a-link @click="handleView(record)">详情</a-link>
              <a-link @click="handleEdit(record)">编辑</a-link>
              <a-link @click="handleDelete(record)">删除</a-link>
            </template>
          </a-table>
        </a-col>
      </a-row>
    </a-card>

    <a-modal v-model:visible="modalVisible" :title="modalTitle" @ok="handleSubmit" @cancel="modalVisible = false">
      <a-form :model="formData" layout="vertical">
        <a-form-item label="分组名称" required>
          <a-input v-model="formData.groupName" placeholder="请输入分组名称" />
        </a-form-item>
        <a-form-item label="上级分组">
          <a-select v-model="formData.parentId" placeholder="请选择上级分组">
            <a-option value="">无</a-option>
            <a-option v-for="g in groups" :key="g.id" :value="g.id">{{ g.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="formData.description" placeholder="请输入描述" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const modalVisible = ref(false)
const modalTitle = ref('新建分组')
const formData = reactive({ groupName: '', parentId: '', description: '' })
const groups = ref([])
const treeData = ref([
  { title: '全部设备', key: 'all', children: [] },
  { title: '华北区域', key: 'north', children: [{ title: '北京', key: 'bj' }, { title: '天津', key: 'tj' }] },
  { title: '华南区域', key: 'south', children: [{ title: '深圳', key: 'sz' }, { title: '广州', key: 'gz' }] }
])
const columns = [
  { title: '设备ID', dataIndex: 'deviceId' },
  { title: '设备名称', dataIndex: 'deviceName' },
  { title: '型号', dataIndex: 'model' },
  { title: '状态', slotName: 'status' },
  { title: '所属分组', dataIndex: 'groupName' },
  { title: '操作', slotName: 'actions', width: 180 }
]
const tableData = ref([
  { deviceId: 'D001', deviceName: '设备-001', model: 'M5Stack-1', status: 'online', groupName: '北京' },
  { deviceId: 'D002', deviceName: '设备-002', model: 'M5Stack-2', status: 'offline', groupName: '深圳' }
])

const handleCreate = () => { modalVisible.value = true; modalTitle.value = '新建分组' }
const handleEdit = (node) => { modalVisible.value = true; modalTitle.value = '编辑分组' }
const handleDelete = (node) => { }
const handleSelect = (keys) => { }
const handleExport = () => { }
const handleView = (record) => { }
const handleSubmit = () => { modalVisible.value = false }
</script>

<style scoped>
.container { padding: 16px; }
</style>
