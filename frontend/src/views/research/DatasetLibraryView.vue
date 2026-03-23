<template>
  <div class="dataset-library">
    <a-card title="AI 数据集开放平台">
      <template #extra>
        <a-button type="primary">申请数据集</a-button>
      </template>
      
      <a-spin :loading="loading">
        <a-table :columns="columns" :data="datasets" :pagination="{ pageSize: 10 }">
          <template #name="{ record }">
            <div>
              <div style="font-weight: 600">{{ record.name }}</div>
              <div style="font-size: 12px; color: var(--color-text-3)">{{ record.description }}</div>
            </div>
          </template>
          <template #category="{ record }">
            <a-tag>{{ record.category }}</a-tag>
          </template>
          <template #actions="{ record }">
            <a-button size="small" @click="download(record)">下载</a-button>
            <a-button size="small" type="text" @click="cite(record)">引用</a-button>
          </template>
        </a-table>
      </a-spin>
    </a-card>
    
    <a-card title="研究项目" style="margin-top: 16px">
      <a-list :data="projects" :loading="loading2">
        <template #item="{ item }">
          <a-list-item>
            <div style="flex: 1">
              <div style="font-weight: 600">{{ item.name }}</div>
              <div style="font-size: 12px; color: var(--color-text-3)">{{ item.description }}</div>
            </div>
            <template #actions>
              <a-tag :color="item.status === 'active' ? 'green' : 'default'">{{ item.status }}</a-tag>
              <a-button size="small" type="text">查看实验</a-button>
            </template>
          </a-list-item>
        </template>
      </a-list>
    </a-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'
const loading = ref(false)
const loading2 = ref(false)
const datasets = ref([])
const projects = ref([])

const columns = [
  { title: '数据集', slotName: 'name', width: 300 },
  { title: '分类', slotName: 'category' },
  { title: '样本数', dataIndex: 'record_count' },
  { title: '格式', dataIndex: 'data_format' },
  { title: '下载量', dataIndex: 'download_count' },
  { title: '操作', slotName: 'actions', width: 150 }
]

const loadDatasets = async () => {
  loading.value = true
  const res = await fetch(`${API_BASE}/research/datasets`)
  const data = await res.json()
  datasets.value = data.datasets || []
  loading.value = false
}

const loadProjects = async () => {
  loading2.value = true
  const res = await fetch(`${API_BASE}/research/projects`)
  projects.value = await res.json()
  loading2.value = false
}

const download = async (record) => {
  const res = await fetch(`${API_BASE}/research/datasets/${record.id}/download`)
  const data = await res.json()
  window.open(data.download_url, '_blank')
}

const cite = async (record) => {
  await fetch(`${API_BASE}/research/datasets/${record.id}/cite`, { method: 'POST' })
  Message.success('引用成功')
}

onMounted(() => { loadDatasets(); loadProjects() })
</script>
