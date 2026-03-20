<template>
  <a-layout class="app-distributions">
    <a-layout-sider v-model:collapsed="collapsed" :trigger="null" collapsible>
      <div class="logo">
        <span v-if="!collapsed">MDM 控制台</span>
      </div>
      <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline" @click="handleMenuClick">
        <a-menu-item key="dashboard"><span>设备大盘</span></a-menu-item>
        <a-menu-item key="status"><span>设备状态</span></a-menu-item>
        <a-menu-item key="pet"><span>宠物配置</span></a-menu-item>
        <a-menu-item key="ota"><span>OTA 固件</span></a-menu-item>
        <a-menu-item key="apps"><span>应用管理</span></a-menu-item>
      </a-menu>
    </a-layout-sider>

    <a-layout>
      <a-layout-header class="header">
        <div class="header-left">
          <a-button type="text" @click="collapsed = !collapsed">
            <span v-if="collapsed">☰</span>
            <span v-else>✕</span>
          </a-button>
        </div>
        <div class="header-title">
          <span>应用分发</span>
        </div>
        <div class="header-right"></div>
      </a-layout-header>

      <a-layout-content class="content">
        <!-- 统计卡片 -->
        <a-row :gutter="16" class="stats-row">
          <a-col :span="6">
            <a-card>
              <a-statistic title="总任务数" :value="stats.total" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="进行中" :value="stats.running" :value-style="{ color: '#1890ff' }" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="已完成" :value="stats.completed" :value-style="{ color: '#52c41a' }" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="成功率" :value="stats.successRate" suffix="%" />
            </a-card>
          </a-col>
        </a-row>

        <!-- 操作栏 -->
        <a-card class="action-card">
          <a-space wrap>
            <a-select v-model="filters.status" placeholder="任务状态" allow-clear style="width: 140px" @change="loadDistributions">
              <a-option value="pending">待执行</a-option>
              <a-option value="running">进行中</a-option>
              <a-option value="completed">已完成</a-option>
              <a-option value="failed">失败</a-option>
              <a-option value="cancelled">已取消</a-option>
            </a-select>
            <a-button type="primary" @click="showCreateDrawer = true">新建分发任务</a-button>
            <a-button @click="loadDistributions">刷新</a-button>
          </a-space>
        </a-card>

        <!-- 分发任务列表 -->
        <a-card class="distribution-card">
          <template #title><span>分发任务列表</span></template>
          <a-table
            :columns="columns"
            :data="distributionList"
            :loading="loading"
            :pagination="pagination"
            row-key="id"
            @page-change="handlePageChange"
          >
            <template #distributionType="{ record }">
              <a-tag :color="getDistTypeColor(record.distribution_type)">
                {{ getDistTypeText(record.distribution_type) }}
              </a-tag>
            </template>
            <template #status="{ record }">
              <a-tag :color="getDistStatusColor(record.status)">
                {{ getDistStatusText(record.status) }}
              </a-tag>
            </template>
            <template #progress="{ record }">
              <a-progress
                v-if="record.total_count > 0"
                :percent="Math.round((record.success_count / record.total_count) * 100)"
                :color="getProgressColor(record.status)"
                style="width: 120px"
              />
              <span v-else>-</span>
            </template>
            <template #targetType="{ record }">
              {{ record.target_type }}/{{ record.target_ids?.length || 0 }}
            </template>
            <template #createdAt="{ record }">
              {{ formatTime(record.created_at) }}
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
                <a-button v-if="record.status === 'pending' || record.status === 'running'" type="text" size="small" status="danger" @click="cancelTask(record)">取消</a-button>
              </a-space>
            </template>
          </a-table>
        </a-card>
      </a-layout-content>
    </a-layout>

    <!-- 创建分发任务抽屉 -->
    <a-drawer
      v-model:visible="showCreateDrawer"
      title="新建分发任务"
      :width="520"
    >
      <a-form :model="createForm" layout="vertical" @submit-success="handleCreateSubmit">
        <a-form-item label="选择应用" required>
          <a-select v-model="createForm.app_id" placeholder="请选择应用" style="width: 100%" @change="onAppChange">
            <a-option v-for="app in appList" :key="app.id" :value="app.id">
              {{ app.app_name }} ({{ app.app_code }})
            </a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="选择版本" required>
          <a-select v-model="createForm.app_version_id" placeholder="请先选择应用" style="width: 100%" :disabled="!createForm.app_id">
            <a-option v-for="ver in versionList" :key="ver.id" :value="ver.id">
              {{ ver.version }} {{ ver.is_latest ? '(最新)' : '' }}
            </a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="分发类型" required>
          <a-select v-model="createForm.distribution_type" placeholder="选择分发类型">
            <a-option value="install">普通安装</a-option>
            <a-option value="force_install">强制安装</a-option>
            <a-option value="uninstall">强制卸载</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="目标类型" required>
          <a-select v-model="createForm.target_type" placeholder="选择目标类型">
            <a-option value="device">设备</a-option>
            <a-option value="user">用户</a-option>
            <a-option value="group">分组</a-option>
            <a-option value="org_unit">组织单元</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="目标 ID 列表" required>
          <a-textarea v-model="targetIdsStr" :rows="3" placeholder="多个 ID 用换行分隔" />
          <div style="color: #999; font-size: 12px; margin-top: 4px">每行一个目标 ID</div>
        </a-form-item>
        <a-form-item label="计划开始时间">
          <a-date-picker v-model="createForm.scheduled_at" style="width: 100%" show-time placeholder="留空则立即执行" />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" html-type="submit">创建</a-button>
            <a-button @click="showCreateDrawer = false">取消</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-drawer>

    <!-- 任务详情抽屉 -->
    <a-drawer
      v-model:visible="showDetailDrawer"
      title="分发任务详情"
      :width="560"
    >
      <template v-if="currentTask">
        <a-descriptions :column="1" bordered size="small" style="margin-bottom: 16px">
          <a-descriptions-item label="任务 ID">{{ currentTask.id }}</a-descriptions-item>
          <a-descriptions-item label="应用">{{ currentTask.app_name }}</a-descriptions-item>
          <a-descriptions-item label="版本">{{ currentTask.version }}</a-descriptions-item>
          <a-descriptions-item label="分发类型">
            <a-tag :color="getDistTypeColor(currentTask.distribution_type)">
              {{ getDistTypeText(currentTask.distribution_type) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="目标类型">{{ currentTask.target_type }}</a-descriptions-item>
          <a-descriptions-item label="目标数量">{{ currentTask.total_count }}</a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="getDistStatusColor(currentTask.status)">
              {{ getDistStatusText(currentTask.status) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="创建人">{{ currentTask.created_by }}</a-descriptions-item>
          <a-descriptions-item label="创建时间">{{ formatTime(currentTask.created_at) }}</a-descriptions-item>
        </a-descriptions>

        <!-- 进度信息 -->
        <a-card>
          <template #title><span>分发进度</span></template>
          <a-space direction="vertical" style="width: 100%" size="large">
            <a-progress
              :percent="currentTask.total_count > 0
                ? Math.round((currentTask.success_count / currentTask.total_count) * 100)
                : 0"
              :color="getProgressColor(currentTask.status)"
            />
            <a-row :gutter="16">
              <a-col :span="8">
                <a-statistic title="目标总数" :value="currentTask.total_count" />
              </a-col>
              <a-col :span="8">
                <a-statistic title="成功" :value="currentTask.success_count" :value-style="{ color: '#52c41a' }" />
              </a-col>
              <a-col :span="8">
                <a-statistic title="失败" :value="currentTask.failed_count" :value-style="{ color: '#ff4d4f' }" />
              </a-col>
            </a-row>
          </a-space>
        </a-card>
      </template>
    </a-drawer>
  </a-layout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import axios from 'axios'
import { Message, Modal } from '@arco-design/web-vue'

const router = useRouter()
const route = useRoute()
const collapsed = ref(false)
const selectedKeys = ref(['apps'])
const loading = ref(false)
const distributionList = ref([])
const appList = ref([])
const versionList = ref([])
const showCreateDrawer = ref(false)
const showDetailDrawer = ref(false)
const currentTask = ref(null)
const targetIdsStr = ref('')

const filters = reactive({
  status: undefined
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const stats = reactive({
  total: 0,
  running: 0,
  completed: 0,
  successRate: 0
})

const createForm = reactive({
  app_id: null,
  app_version_id: null,
  distribution_type: '',
  target_type: '',
  scheduled_at: null,
  created_by: 'admin'
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '应用', dataIndex: 'app_name', width: 150, ellipsis: true },
  { title: '版本', dataIndex: 'version', width: 100 },
  { title: '分发类型', slotName: 'distributionType', width: 100 },
  { title: '目标', slotName: 'targetType', width: 100 },
  { title: '进度', slotName: 'progress', width: 140 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '创建时间', slotName: 'createdAt', width: 160 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' }
]

const handleMenuClick = ({ key }) => {
  if (key === 'apps') {
    router.push('/apps')
  } else {
    selectedKeys.value = [key]
    router.push('/' + key)
  }
}

const loadDistributions = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (filters.status) params.status = filters.status

    const res = await axios.get('/api/v1/apps/distributions', { params })
    if (res.data.code === 0) {
      distributionList.value = res.data.data.list || []
      pagination.total = res.data.data.pagination?.total || 0

      // 更新统计
      stats.total = distributionList.value.length
      stats.running = distributionList.value.filter(d => d.status === 'running').length
      stats.completed = distributionList.value.filter(d => d.status === 'completed').length
      const totalSuccess = distributionList.value.reduce((sum, d) => sum + (d.success_count || 0), 0)
      const totalAll = distributionList.value.reduce((sum, d) => sum + (d.total_count || 0), 0)
      stats.successRate = totalAll > 0 ? Math.round((totalSuccess / totalAll) * 100) : 0
    }
  } catch (err) {
    Message.error('加载分发任务列表失败')
  } finally {
    loading.value = false
  }
}

const loadApps = async () => {
  try {
    const res = await axios.get('/api/v1/apps', { params: { page: 1, page_size: 100 } })
    if (res.data.code === 0) {
      appList.value = res.data.data.list || []
    }
  } catch (err) {}
}

const onAppChange = async (appId) => {
  createForm.app_version_id = null
  if (!appId) return
  try {
    const res = await axios.get(`/api/v1/apps/${appId}/versions`)
    if (res.data.code === 0) {
      versionList.value = res.data.data.list || []
    }
  } catch (err) {}
}

const handleCreateSubmit = async () => {
  try {
    const payload = { ...createForm }
    if (targetIdsStr.value) {
      payload.target_ids = targetIdsStr.value.split('\n').map(s => s.trim()).filter(Boolean)
    }
    if (payload.scheduled_at) {
      payload.scheduled_at = new Date(payload.scheduled_at).toISOString()
    }
    const res = await axios.post('/api/v1/apps/distributions', payload)
    if (res.data.code === 0) {
      Message.success('创建成功')
      showCreateDrawer.value = false
      resetCreateForm()
      loadDistributions()
    } else {
      Message.error(res.data.message || '创建失败')
    }
  } catch (err) {
    Message.error('创建失败')
  }
}

const resetCreateForm = () => {
  createForm.app_id = null
  createForm.app_version_id = null
  createForm.distribution_type = ''
  createForm.target_type = ''
  createForm.scheduled_at = null
  targetIdsStr.value = ''
  versionList.value = []
}

const openDetail = async (record) => {
  try {
    const res = await axios.get(`/api/v1/apps/distributions/${record.id}`)
    if (res.data.code === 0) {
      currentTask.value = res.data.data
      showDetailDrawer.value = true
    }
  } catch (err) {
    Message.error('加载详情失败')
  }
}

const cancelTask = (record) => {
  Modal.confirm({
    title: '确认取消',
    content: '确定要取消该分发任务吗？',
    onOk: async () => {
      try {
        const res = await axios.post(`/api/v1/apps/distributions/${record.id}/cancel`)
        if (res.data.code === 0) {
          Message.success('已取消')
          loadDistributions()
        } else {
          Message.error(res.data.message || '取消失败')
        }
      } catch (err) {
        Message.error('取消失败')
      }
    }
  })
}

const handlePageChange = (page) => {
  pagination.current = page
  loadDistributions()
}

const getDistTypeColor = (type) => {
  const map = { install: 'blue', force_install: 'orange', uninstall: 'red' }
  return map[type] || 'default'
}

const getDistTypeText = (type) => {
  const map = { install: '普通安装', force_install: '强制安装', uninstall: '卸载' }
  return map[type] || type
}

const getDistStatusColor = (status) => {
  const map = { pending: 'yellow', running: 'blue', completed: 'green', failed: 'red', cancelled: 'gray' }
  return map[status] || 'default'
}

const getDistStatusText = (status) => {
  const map = { pending: '待执行', running: '进行中', completed: '已完成', failed: '失败', cancelled: '已取消' }
  return map[status] || status
}

const getProgressColor = (status) => {
  if (status === 'failed') return '#ff4d4f'
  if (status === 'completed') return '#52c41a'
  return '#1890ff'
}

const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

onMounted(() => {
  loadDistributions()
  loadApps()

  // 如果有 appId query 参数，自动打开创建
  if (route.query.appId) {
    createForm.app_id = parseInt(route.query.appId)
    onAppChange(createForm.app_id)
    showCreateDrawer.value = true
  }
})
</script>

<style scoped>
.app-distributions { min-height: 100vh; }
.header { background: #fff; padding: 0 16px; display: flex; align-items: center; gap: 16px; box-shadow: 0 1px 4px rgba(0,0,0,0.1); }
.header-left { display: flex; align-items: center; }
.header-title { font-size: 16px; font-weight: 500; }
.content { padding: 16px; background: #f0f2f5; }
.stats-row { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
