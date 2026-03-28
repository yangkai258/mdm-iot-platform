<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-database /> 仿真数据集管理</a-space>
      </template>
      <template #extra>
        <a-button type="primary" @click="handleImport">
          <template #icon><icon-upload /></template>
          导入数据集
        </a-button>
      </template>

      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6">
          <a-card size="small">
            <a-statistic title="数据集总数" :value="stats.total" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card size="small">
            <a-statistic title="总大小" :value="stats.totalSize" suffix="MB" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card size="small">
            <a-statistic title="版本数" :value="stats.versionCount" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card size="small">
            <a-statistic title="存储配额" :value="stats.quotaUsed" suffix="%">
              <template #prefix>
                <a-progress :percent="stats.quotaUsed" :show-text="false" size="small" />
              </template>
            </a-statistic>
          </a-card>
        </a-col>
      </a-row>

      <a-table :columns="columns" :data="datasets">
        <template #tags="{ record }">
          <a-tag v-for="tag in record.tags" :key="tag" color="blue">{{ tag }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-link @click="handleView(record)">查看</a-link>
          <a-link @click="handleVersion(record)">版本</a-link>
          <a-link @click="handleExport(record)">导出</a-link>
          <a-link @click="handleDelete(record)">删除</a-link>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="importVisible" title="导入数据集" @ok="handleImportConfirm">
      <a-form :model="importForm" layout="vertical">
        <a-form-item label="数据集名称" required>
          <a-input v-model="importForm.name" placeholder="请输入名称" />
        </a-form-item>
        <a-form-item label="版本号">
          <a-input v-model="importForm.version" placeholder="v1.0.0" />
        </a-form-item>
        <a-form-item label="数据集文件">
          <a-upload draggable accept=".zip,.tar,.tar.gz" />
        </a-form-item>
        <a-form-item label="标签">
          <a-select v-model="importForm.tags" multiple placeholder="选择标签">
            <a-option value="pet">宠物</a-option>
            <a-option value="behavior">行为</a-option>
            <a-option value="environment">环境</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const stats = reactive({ total: 15, totalSize: 2048, versionCount: 42, quotaUsed: 45 })
const importVisible = ref(false)
const importForm = reactive({ name: '', version: '', tags: [] })

const columns = [
  { title: '名称', dataIndex: 'name' },
  { title: '版本', dataIndex: 'version' },
  { title: '文件数', dataIndex: 'fileCount' },
  { title: '大小', dataIndex: 'size' },
  { title: '标签', slotName: 'tags' },
  { title: '创建时间', dataIndex: 'createdAt' },
  { title: '操作', slotName: 'actions', width: 200 }
]

const datasets = ref([
  { id: 1, name: '宠物行为数据集', version: 'v2.1.0', fileCount: 1250, size: '256MB', tags: ['pet', 'behavior'], createdAt: '2026-03-28' }
])

const handleImport = () => { importVisible.value = true }
const handleView = (r) => { }
const handleVersion = (r) => { }
const handleExport = (r) => { }
const handleDelete = (r) => { }
const handleImportConfirm = () => { importVisible.value = false }
</script>

<style scoped>
.container { padding: 16px; }
</style>
