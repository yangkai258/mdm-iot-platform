<template>
  <div class="offline-container">
    <a-card>
      <template #title>
        <span>离线同步管理</span>
      </template>
      
      <template #extra>
        <a-button type="primary" @click="syncAll">
          <icon-sync :spin="syncing" />
          全部同步
        </a-button>
      </template>

      <!-- 设备缓存状态 -->
      <a-tabs default-active-key="cache">
        <a-tab key="cache" title="设备缓存">
          <a-table :data="cacheList" :loading="loading" :pagination="false">
            <a-table-column title="设备ID" dataIndex="device_id"></a-table-column>
            <a-table-column title="缓存Key" dataIndex="cache_key"></a-table-column>
            <a-table-column title="状态" dataIndex="sync_status">
              <template #cell="{ record }">
                <a-tag :color="record.sync_status === 'synced' ? 'green' : 'orange'">
                  {{ record.sync_status === 'synced' ? '已同步' : '待同步' }}
                </a-tag>
              </template>
            </a-table-column>
            <a-table-column title="创建时间" dataIndex="created_at"></a-table-column>
            <a-table-column title="操作">
              <template #cell="{ record }">
                <a-button size="small" @click="syncCache(record)">同步</a-button>
              </template>
            </a-table-column>
          </a-table>
        </a-tab>

        <a-tab key="queue" title="离线队列">
          <a-table :data="queueList" :loading="loading" :pagination="pagination">
            <a-table-column title="设备ID" dataIndex="device_id"></a-table-column>
            <a-table-column title="动作类型" dataIndex="action_type"></a-table-column>
            <a-table-column title="状态" dataIndex="status">
              <template #cell="{ record }">
                <a-tag :color="getStatusColor(record.status)">
                  {{ getStatusText(record.status) }}
                </a-tag>
              </template>
            </a-table-column>
            <a-table-column title="创建时间" dataIndex="created_at"></a-table-column>
            <a-table-column title="操作">
              <template #cell="{ record }">
                <a-button v-if="record.status === 'pending'" size="small" type="primary" @click="sendAction(record)">
                  发送
                </a-button>
              </template>
            </a-table-column>
          </a-table>
        </a-tab>

        <a-tab key="settings" title="同步设置">
          <a-form :model="syncSettings" layout="vertical">
            <a-form-item label="自动同步间隔">
              <a-select v-model="syncSettings.interval">
                <a-option value="5">5分钟</a-option>
                <a-option value="15">15分钟</a-option>
                <a-option value="30">30分钟</a-option>
                <a-option value="60">1小时</a-option>
              </a-select>
            </a-form-item>
            <a-form-item label="Wi-Fi下自动同步">
              <a-switch v-model="syncSettings.wifiOnly" />
            </a-form-item>
            <a-form-item label="同步失败自动重试">
              <a-switch v-model="syncSettings.autoRetry" />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="saveSettings">保存设置</a-button>
            </a-form-item>
          </a-form>
        </a-tab>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const syncing = ref(false)
const cacheList = ref([])
const queueList = ref([])

const pagination = ref({
  current: 1,
  pageSize: 20,
  total: 0
})

const syncSettings = ref({
  interval: '15',
  wifiOnly: true,
  autoRetry: true
})

const getStatusColor = (status) => {
  const colors = { pending: 'orange', sent: 'green', failed: 'red' }
  return colors[status] || 'gray'
}

const getStatusText = (status) => {
  const texts = { pending: '待发送', sent: '已发送', failed: '失败' }
  return texts[status] || status
}

const loadCacheList = async () => {
  try {
    const res = await fetch('/api/v1/offline/cache/device-1')
    const data = await res.json()
    if (data.code === 0) {
      cacheList.value = data.data || []
    }
  } catch (e) {
    Message.error('加载缓存列表失败')
  }
}

const loadQueueList = async () => {
  try {
    const res = await fetch('/api/v1/offline/queue/device-1')
    const data = await res.json()
    if (data.code === 0) {
      queueList.value = data.data || []
    }
  } catch (e) {
    Message.error('加载队列失败')
  }
}

const syncCache = async (record) => {
  try {
    const res = await fetch(`/api/v1/offline/cache/${record.id}/confirm`, { method: 'POST' })
    const data = await res.json()
    if (data.code === 0) {
      Message.success('同步成功')
      loadCacheList()
    }
  } catch (e) {
    Message.error('同步失败')
  }
}

const syncAll = async () => {
  syncing.value = true
  try {
    const res = await fetch('/api/v1/offline/cache/sync?device_id=device-1')
    const data = await res.json()
    if (data.code === 0) {
      Message.success(`已同步 ${data.data.pending_items} 项`)
      loadCacheList()
    }
  } finally {
    syncing.value = false
  }
}

const sendAction = async (record) => {
  try {
    const res = await fetch(`/api/v1/offline/queue/${record.id}/send`, { method: 'POST' })
    const data = await res.json()
    if (data.code === 0) {
      Message.success('发送成功')
      loadQueueList()
    }
  } catch (e) {
    Message.error('发送失败')
  }
}

const saveSettings = () => {
  Message.success('设置已保存')
}

onMounted(() => {
  loadCacheList()
  loadQueueList()
})
</script>

<style scoped>
.offline-container {
  padding: 16px;
}
</style>
