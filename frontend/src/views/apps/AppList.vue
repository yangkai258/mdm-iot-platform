<template>
  <div class="page-container">
<!-- 统计卡片 -->
        <a-row :gutter="16" class="stats-row">
          <a-col :span="6">
            <a-card>
              <a-statistic title="总应用数" :value="stats.total" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="已发布" :value="stats.published" :value-style="{ color: '#52c41a' }" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="待审核" :value="stats.pending" :value-style="{ color: '#faad14' }" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="安装次数" :value="stats.installCount" />
            </a-card>
          </a-col>
        </a-row>

        <!-- 操作栏 -->
        <a-card class="action-card">
          <a-space wrap>
            <a-select v-model="filters.appType" placeholder="应用类型" allow-clear style="width: 140px" @change="loadApps">
              <a-option value="ios">iOS</a-option>
              <a-option value="android">Android</a-option>
              <a-option value="windows">Windows</a-option>
              <a-option value="macos">macOS</a-option>
            </a-select>
            <a-select v-model="filters.status" placeholder="应用状态" allow-clear style="width: 140px" @change="loadApps">
              <a-option value="pending">待审核</a-option>
              <a-option value="approved">已发布</a-option>
              <a-option value="rejected">已拒绝</a-option>
              <a-option value="archived">已归档</a-option>
            </a-select>
            <a-input-search v-model="filters.keyword" placeholder="搜索应用名称/编码" style="width: 200px" search-button @search="loadApps" />
            <a-button type="primary" @click="showCreateDrawer = true">上传应用</a-button>
            <a-button @click="loadApps">刷新</a-button>
          </a-space>
        </a-card>

        <!-- 应用列表 -->
        <a-card class="app-card">
          <a-table
            :columns="columns"
            :data="appList"
            :loading="loading"
            :pagination="pagination"
            row-key="id"
            @page-change="handlePageChange"
            @page-size-change="handlePageSizeChange"
          >
            <template #appType="{ record }">
              <a-tag :color="getAppTypeColor(record.app_type)">{{ record.app_type.toUpperCase() }}</a-tag>
            </template>
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
            </template>
            <template #installCount="{ record }">
              {{ record.install_count || 0 }}
            </template>
            <template #createdAt="{ record }">
              {{ formatTime(record.created_at) }}
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
                <a-button type="text" size="small" @click="openVersions(record)">版本</a-button>
                <a-button type="text" size="small" @click="openDistribute(record)">分发</a-button>
                <a-button v-if="record.status === 'pending'" type="text" size="small" @click="approveApp(record)">审核</a-button>
                <a-button v-if="record.status === 'approved'" type="text" size="small" status="warning" @click="rejectApp(record)">拒绝</a-button>
              </a-space>
            </template>
          </a-table>
        </a-card>
</div>

    <!-- 创建应用抽屉 -->
    <a-drawer
      v-model:visible="showCreateDrawer"
      title="上传应用"
      :width="480"
      @before-ok="handleCreate"
    >
      <a-form :model="createForm" layout="vertical" @submit-success="handleCreateSubmit">
        <a-form-item label="应用编码" required>
          <a-input v-model="createForm.app_code" placeholder="如 APP001" />
        </a-form-item>
        <a-form-item label="应用名称" required>
          <a-input v-model="createForm.app_name" placeholder="如 企业IM" />
        </a-form-item>
        <a-form-item label="应用类型" required>
          <a-select v-model="createForm.app_type" placeholder="选择应用类型">
            <a-option value="ios">iOS</a-option>
            <a-option value="android">Android</a-option>
            <a-option value="windows">Windows</a-option>
            <a-option value="macos">macOS</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="Bundle ID / Package Name" required>
          <a-input v-model="createForm.bundle_id" placeholder="如 com.company.im" />
        </a-form-item>
        <a-form-item label="开发者">
          <a-input v-model="createForm.developer" placeholder="内部开发组" />
        </a-form-item>
        <a-form-item label="应用描述">
          <a-textarea v-model="createForm.description" :rows="3" placeholder="应用功能描述" />
        </a-form-item>
        <a-form-item label="企业应用">
          <a-switch v-model="createForm.is_enterprise" />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" html-type="submit">创建</a-button>
            <a-button @click="showCreateDrawer = false">取消</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-drawer>

    <!-- 应用详情抽屉 -->
    <a-drawer
      v-model:visible="showDetailDrawer"
      title="应用详情"
      :width="520"
    >
      <template v-if="currentApp">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="应用编码">{{ currentApp.app_code }}</a-descriptions-item>
          <a-descriptions-item label="应用名称">{{ currentApp.app_name }}</a-descriptions-item>
          <a-descriptions-item label="类型">
            <a-tag :color="getAppTypeColor(currentApp.app_type)">{{ currentApp.app_type.toUpperCase() }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="Bundle ID">{{ currentApp.bundle_id }}</a-descriptions-item>
          <a-descriptions-item label="开发者">{{ currentApp.developer || '-' }}</a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="getStatusColor(currentApp.status)">{{ getStatusText(currentApp.status) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="安装次数">{{ currentApp.install_count || 0 }}</a-descriptions-item>
          <a-descriptions-item label="最新版本">{{ currentApp.latest_version || '-' }}</a-descriptions-item>
          <a-descriptions-item label="应用描述">{{ currentApp.description || '-' }}</a-descriptions-item>
          <a-descriptions-item label="创建时间">{{ formatTime(currentApp.created_at) }}</a-descriptions-item>
        </a-descriptions>
      </template>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import axios from 'axios'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const appList = ref([])
const showCreateDrawer = ref(false)
const showDetailDrawer = ref(false)
const currentApp = ref(null)

const filters = reactive({
  appType: undefined,
  status: undefined,
  keyword: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const stats = reactive({
  total: 0,
  published: 0,
  pending: 0,
  installCount: 0
})

const createForm = reactive({
  app_code: '',
  app_name: '',
  app_type: '',
  bundle_id: '',
  developer: '',
  description: '',
  is_enterprise: true,
  created_by: 'admin'
})

const columns = [
  { title: '应用名称', dataIndex: 'app_name', width: 180 },
  { title: '编码', dataIndex: 'app_code', width: 100 },
  { title: '类型', slotName: 'appType', width: 80 },
  { title: 'Bundle ID', dataIndex: 'bundle_id', width: 200, ellipsis: true },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '版本', dataIndex: 'latest_version', width: 100 },
  { title: '安装次数', slotName: 'installCount', width: 100 },
  { title: '创建时间', slotName: 'createdAt', width: 160 },
  { title: '操作', slotName: 'actions', width: 220, fixed: 'right' }
]

else {
    selectedKeys.value = [key]
    router.push('/' + key)
  }
}

const loadApps = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize
    }
    if (filters.appType) params.app_type = filters.appType
    if (filters.status) params.status = filters.status
    if (filters.keyword) params.keyword = filters.keyword

    const res = await axios.get('/api/v1/apps', { params })
    const data = res.data
    if (data.code === 0) {
      appList.value = data.data.list || []
      pagination.total = data.data.pagination?.total || 0

      // 更新统计数据
      stats.total = appList.value.length
      stats.published = appList.value.filter(a => a.status === 'approved').length
      stats.pending = appList.value.filter(a => a.status === 'pending').length
      stats.installCount = appList.value.reduce((sum, a) => sum + (a.install_count || 0), 0)
    }
  } catch (err) {
    Message.error('加载应用列表失败')
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page) => {
  pagination.current = page
  loadApps()
}

const handlePageSizeChange = (pageSize) => {
  pagination.pageSize = pageSize
  pagination.current = 1
  loadApps()
}

const handleCreate = (done) => {
  done(true)
}

const handleCreateSubmit = async () => {
  try {
    const res = await axios.post('/api/v1/apps', createForm)
    if (res.data.code === 0) {
      Message.success('创建成功')
      showCreateDrawer.value = false
      resetCreateForm()
      loadApps()
    } else {
      Message.error(res.data.message || '创建失败')
    }
  } catch (err) {
    Message.error('创建失败')
  }
}

const resetCreateForm = () => {
  createForm.app_code = ''
  createForm.app_name = ''
  createForm.app_type = ''
  createForm.bundle_id = ''
  createForm.developer = ''
  createForm.description = ''
  createForm.is_enterprise = true
}

const openDetail = (record) => {
  currentApp.value = record
  showDetailDrawer.value = true
}

const openVersions = (record) => {
  router.push(`/apps/versions/${record.id}`)
}

const openDistribute = (record) => {
  router.push({ path: '/apps/distributions', query: { appId: record.id } })
}

const approveApp = async (record) => {
  try {
    const res = await axios.post(`/api/v1/apps/${record.id}/approve`)
    if (res.data.code === 0) {
      Message.success('审核通过')
      loadApps()
    } else {
      Message.error(res.data.message || '操作失败')
    }
  } catch (err) {
    Message.error('操作失败')
  }
}

const rejectApp = async (record) => {
  try {
    const res = await axios.post(`/api/v1/apps/${record.id}/reject`, { reason: '' })
    if (res.data.code === 0) {
      Message.success('已拒绝')
      loadApps()
    } else {
      Message.error(res.data.message || '操作失败')
    }
  } catch (err) {
    Message.error('操作失败')
  }
}

const getAppTypeColor = (type) => {
  const map = { ios: 'pink', android: 'green', windows: 'blue', macos: 'purple' }
  return map[type] || 'gray'
}

const getStatusColor = (status) => {
  const map = { pending: 'yellow', approved: 'green', rejected: 'red', archived: 'gray' }
  return map[status] || 'gray'
}

const getStatusText = (status) => {
  const map = { pending: '待审核', approved: '已发布', rejected: '已拒绝', archived: '已归档' }
  return map[status] || status
}

const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

onMounted(() => {
  loadApps()
})
</script>

<style scoped>
.app-list { min-height: 100vh; }
.header { background: #fff; padding: 0 16px; display: flex; align-items: center; gap: 16px; box-shadow: 0 1px 4px rgba(0,0,0,0.1); }
.header-left { display: flex; align-items: center; }
.header-title { font-size: 16px; font-weight: 500; }
.content { padding: 16px; background: #f0f2f5; }
.stats-row { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
.app-card { }
</style>
