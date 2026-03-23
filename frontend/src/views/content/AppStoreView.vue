<template>
  <div class="app-store">
    <a-card title="应用商店">
      <a-tabs>
        <a-tab-pane key="packages" title="应用包">
          <a-spin :loading="loading">
            <div class="app-grid">
              <a-card v-for="pkg in packages" :key="pkg.id" class="app-card" hoverable>
                <div class="app-icon">
                  <icon-app :size="40" />
                </div>
                <div class="app-name">{{ pkg.name }}</div>
                <div class="app-version">v{{ pkg.version }}</div>
                <a-tag :color="getStatusColor(pkg.status)">{{ pkg.status }}</a-tag>
                <div class="app-actions">
                  <a-button size="small" v-if="pkg.status === 'approved'" type="primary">安装</a-button>
                  <a-button size="small" v-else-if="pkg.status === 'draft'" @click="submitReview(pkg)">提交审核</a-button>
                </div>
              </a-card>
            </div>
          </a-spin>
        </a-tab-pane>
        <a-tab-pane key="installed" title="已安装">
          <a-list :data="installs">
            <template #item="{ item }">
              <a-list-item>
                <span>App #{{ item.app_id }} - v{{ item.version }}</span>
                <template #actions>
                  <a-tag :color="item.status === 'installed' ? 'green' : 'orange'">{{ item.status }}</a-tag>
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
const packages = ref([])
const installs = ref([])

const getStatusColor = (s) => ({ approved: 'green', pending_review: 'orange', rejected: 'red', draft: 'default' }[s] || 'default')

const loadPackages = async () => {
  loading.value = true
  const res = await fetch(`${API_BASE}/apps/packages`)
  packages.value = await res.json()
  loading.value = false
}

const loadInstalls = async () => {
  const res = await fetch(`${API_BASE}/apps/installs`)
  installs.value = await res.json()
}

const submitReview = async (pkg) => {
  await fetch(`${API_BASE}/apps/packages/${pkg.id}/submit-review`, { method: 'POST' })
  Message.success('已提交审核')
  loadPackages()
}

onMounted(() => { loadPackages(); loadInstalls() })
</script>

<style scoped>
.app-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(180px, 1fr)); gap: 16px; }
.app-card { text-align: center; }
.app-icon { margin-bottom: 8px; }
.app-name { font-weight: 600; margin-bottom: 4px; }
.app-version { font-size: 12px; color: var(--color-text-3); margin-bottom: 8px; }
.app-actions { margin-top: 8px; }
</style>
