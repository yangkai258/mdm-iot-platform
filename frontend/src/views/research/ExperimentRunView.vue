<template>
  <div class="experiment-run">
    <a-card title="AI 实验运行监控">
      <a-tabs>
        <a-tab-pane key="experiments" title="实验列表">
          <a-spin :loading="loading">
            <a-table :columns="columns" :data="experiments" :pagination="{ pageSize: 10 }">
              <template #status="{ record }">
                <a-tag :color="getStatusColor(record.status)">{{ record.status }}</a-tag>
              </template>
              <template #actions="{ record }">
                <a-button size="small" v-if="record.status === 'draft'" type="primary" @click="start(record)">启动</a-button>
                <a-button size="small" v-else-if="record.status === 'running'" type="danger" @click="stop(record)">停止</a-button>
                <a-button size="small" v-else @click="viewResults(record)">查看结果</a-button>
              </template>
            </a-table>
          </a-spin>
        </a-tab-pane>
        <a-tab-pane key="collaborations" title="学术合作">
          <a-list :data="collaborations">
            <template #item="{ item }">
              <a-list-item>
                <span>合作 #{{ item.id }} - {{ item.role }}</span>
                <template #actions>
                  <a-tag :color="item.status === 'accepted' ? 'green' : 'orange'">{{ item.status }}</a-tag>
                </template>
              </a-list-item>
            </template>
          </a-list>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'
const loading = ref(false)
const experiments = ref([])
const collaborations = ref([])

const columns = [
  { title: '实验名称', dataIndex: 'name' },
  { title: '项目ID', dataIndex: 'project_id' },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '开始时间', dataIndex: 'started_at', render: ({ started_at }) => started_at ? new Date(started_at).toLocaleString() : '-' },
  { title: '操作', slotName: 'actions', width: 120 }
]

const getStatusColor = (s) => ({ running: 'green', completed: 'arcoblue', failed: 'red', draft: 'default' }[s] || 'default')

const load = async () => {
  loading.value = true
  const [expRes, collabRes] = await Promise.all([
    fetch(`${API_BASE}/research/experiments`),
    fetch(`${API_BASE}/research/collaborations`)
  ])
  experiments.value = await expRes.json()
  collaborations.value = await collabRes.json()
  loading.value = false
}

const start = async (exp) => {
  await fetch(`${API_BASE}/research/experiments/${exp.id}/start`, { method: 'POST' })
  Message.success('已启动')
  load()
}

const stop = async (exp) => {
  await fetch(`${API_BASE}/research/experiments/${exp.id}/stop`, { 
    method: 'POST', 
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ status: 'completed' })
  })
  Message.success('已停止')
  load()
}

const viewResults = (exp) => Message.info('结果查看开发中')

onMounted(load)
</script>
