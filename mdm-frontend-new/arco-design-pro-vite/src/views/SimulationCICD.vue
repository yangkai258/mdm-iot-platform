<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-sync-circle /> CI/CD集成配置</a-space>
      </template>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-card title="CI/CD平台配置">
            <a-form :model="platformForm" layout="vertical">
              <a-form-item label="平台类型">
                <a-select v-model="platformForm.type">
                  <a-option value="github">GitHub Actions</a-option>
                  <a-option value="gitlab">GitLab CI</a-option>
                  <a-option value="jenkins">Jenkins</a-option>
                </a-select>
              </a-form-item>
              <a-form-item label="仓库URL">
                <a-input v-model="platformForm.repoUrl" placeholder="https://github.com/xxx/repo" />
              </a-form-item>
              <a-form-item label="Webhook Secret">
                <a-input-password v-model="platformForm.webhookSecret" />
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="handleTestConnection">测试连接</a-button>
                <a-button @click="handleSave">保存配置</a-button>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>

        <a-col :span="12">
          <a-card title="连接状态">
            <a-descriptions :column="1" bordered>
              <a-descriptions-item label="连接状态">
                <a-badge status="success" text="已连接" />
              </a-descriptions-item>
              <a-descriptions-item label="最后同步">2026-03-28 10:00:00</a-descriptions-item>
              <a-descriptions-item label="同步状态">
                <a-badge status="success" text="正常" />
              </a-descriptions-item>
            </a-descriptions>
          </a-card>

          <a-card title="触发规则配置" style="margin-top: 16px">
            <a-form :model="triggerForm" layout="vertical">
              <a-form-item label="触发事件">
                <a-checkbox-group v-model="triggerForm.events">
                  <a-checkbox value="push">代码推送</a-checkbox>
                  <a-checkbox value="pr">Pull Request</a-checkbox>
                  <a-checkbox value="tag">标签发布</a-checkbox>
                </a-checkbox-group>
              </a-form-item>
              <a-form-item label="触发分支">
                <a-input v-model="triggerForm.branches" placeholder="main,develop" />
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>
      </a-row>

      <a-card title="构建历史" style="margin-top: 16px">
        <template #extra>
          <a-button @click="handleRefresh">
            <template #icon><icon-refresh /></template>
            刷新
          </a-button>
        </template>
        <a-table :columns="buildColumns" :data="buildHistory">
          <template #status="{ record }">
            <a-badge :status="getBuildStatus(record.status)" :text="getBuildStatusText(record.status)" />
          </template>
          <template #actions="{ record }">
            <a-link @click="handleViewLog(record)">查看日志</a-link>
            <a-link @click="handleRebuild(record)">重新构建</a-link>
          </template>
        </a-table>
      </a-card>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const platformForm = reactive({ type: 'github', repoUrl: '', webhookSecret: '' })
const triggerForm = reactive({ events: ['push'], branches: 'main' })

const buildColumns = [
  { title: '构建ID', dataIndex: 'id', width: 100 },
  { title: '触发方式', dataIndex: 'trigger' },
  { title: '分支', dataIndex: 'branch' },
  { title: '状态', slotName: 'status', width: 120 },
  { title: '开始时间', dataIndex: 'startTime' },
  { title: '耗时', dataIndex: 'duration' },
  { title: '操作', slotName: 'actions', width: 180 }
]

const buildHistory = ref([
  { id: 'B001', trigger: '代码推送', branch: 'main', status: 'success', startTime: '2026-03-28 10:00', duration: '5分钟' },
  { id: 'B002', trigger: 'Pull Request', branch: 'feature', status: 'failed', startTime: '2026-03-28 09:00', duration: '3分钟' }
])

const getBuildStatus = (s) => ({ success: 'success', failed: 'error', running: 'processing' }[s] || 'default')
const getBuildStatusText = (s) => ({ success: '成功', failed: '失败', running: '运行中' }[s] || s)

const handleTestConnection = () => { }
const handleSave = () => { }
const handleRefresh = () => { }
const handleViewLog = (r) => { }
const handleRebuild = (r) => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
