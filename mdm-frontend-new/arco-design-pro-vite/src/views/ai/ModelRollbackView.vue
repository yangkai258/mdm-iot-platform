<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>AI 功能</a-breadcrumb-item>
      <a-breadcrumb-item>模型热回滚</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 提示信息 -->
    <a-alert class="rollback-alert" type="warning">
      <template #title>热回滚说明</template>
      <template #content>
        热回滚可在不中断服务的情况下切换到历史版本。回滚过程通常在 30 秒内完成，期间系统仍可正常处理请求。
      </template>
    </a-alert>

    <!-- 操作栏 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showAddModal = true">创建回滚任务</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </div>

    <!-- 搜索表单 -->
    <div class="pro-search-bar">
      <a-form :model="searchForm" layout="inline">
        <a-form-item label="模型">
          <a-select v-model="searchForm.model_name" placeholder="选择模型" allow-clear style="width: 180px">
            <a-option value="behavior">行为识别模型</a-option>
            <a-option value="emotion">情感分析模型</a-option>
            <a-option value="voice">语音合成模型</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model="searchForm.status" placeholder="选择状态" allow-clear style="width: 120px">
            <a-option value="success">成功</a-option>
            <a-option value="failed">失败</a-option>
            <a-option value="pending">进行中</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSearch">搜索</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </div>

    <!-- 回滚历史列表 -->
    <div class="pro-content-area">
      <a-table
        :columns="columns"
        :data="dataList"
        :loading="loading"
        :pagination="pagination"
        :scroll="{ x: 1300 }"
        @change="handleTableChange"
        row-key="id"
      >
        <template #model_name="{ record }">
          <span class="model-name">{{ record.model_name }}</span>
        </template>
        <template #from_version="{ record }">
          <a-tag>v{{ record.from_version }}</a-tag>
        </template>
        <template #to_version="{ record }">
          <a-tag color="arcoblue">v{{ record.to_version }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-button
            v-if="record.status === 'success' && !record.is_current"
            type="primary"
            status="warning"
            size="small"
            @click="handleRollback(record)"
          >
            回滚此版本
          </a-button>
          <a-button v-else-if="record.is_current" type="text" size="small" disabled>
            当前版本
          </a-button>
          <a-button type="text" size="small" @click="viewDetail(record)">详情</a-button>
        </template>
      </a-table>
    </div>

    <!-- 回滚确认弹窗 -->
    <a-modal v-model:visible="confirmVisible" title="确认回滚" :width="500">
      <a-result status="warning" title="即将回滚模型版本">
        <template #content>
          <a-descriptions :column="1" bordered size="small">
            <a-descriptions-item label="模型">{{ rollbackTarget?.model_name }}</a-descriptions-item>
            <a-descriptions-item label="当前版本">v{{ rollbackTarget?.from_version }}</a-descriptions-item>
            <a-descriptions-item label="回滚版本">v{{ rollbackTarget?.to_version }}</a-descriptions-item>
          </a-descriptions>
          <p style="margin-top: 12px; color: var(--color-text-3)">
            回滚后新请求将使用旧版本模型，历史请求不受影响。确认执行回滚？
          </p>
        </template>
      </a-result>
      <template #footer>
        <a-button @click="confirmVisible = false">取消</a-button>
        <a-button type="primary" status="warning" @click="executeRollback" :loading="rolling">
          确认回滚
        </a-button>
      </template>
    </a-modal>

    <!-- 详情弹窗 -->
    <a-modal v-model:visible="detailVisible" title="回滚详情" :width="640" footer="null">
      <a-descriptions :column="2" bordered size="large">
        <a-descriptions-item label="任务ID">{{ currentRecord?.id }}</a-descriptions-item>
        <a-descriptions-item label="模型">{{ currentRecord?.model_name }}</a-descriptions-item>
        <a-descriptions-item label="源版本">v{{ currentRecord?.from_version }}</a-descriptions-item>
        <a-descriptions-item label="目标版本">v{{ currentRecord?.to_version }}</a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="getStatusColor(currentRecord?.status)">{{ getStatusText(currentRecord?.status) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="执行人">{{ currentRecord?.operator }}</a-descriptions-item>
        <a-descriptions-item label="开始时间">{{ currentRecord?.started_at }}</a-descriptions-item>
        <a-descriptions-item label="完成时间">{{ currentRecord?.finished_at }}</a-descriptions-item>
        <a-descriptions-item label="回滚日志" :span="2">
          <pre style="margin: 0; font-size: 12px; white-space: pre-wrap">{{ currentRecord?.log }}</pre>
        </a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const dataList = ref([])
const confirmVisible = ref(false)
const detailVisible = ref(false)
const rolling = ref(false)
const currentRecord = ref(null)
const rollbackTarget = ref(null)

const searchForm = reactive({
  model_name: '',
  status: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '任务ID', dataIndex: 'id', width: 160 },
  { title: '模型', slotName: 'model_name', width: 160 },
  { title: '源版本', slotName: 'from_version', width: 100 },
  { title: '目标版本', slotName: 'to_version', width: 110 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '执行人', dataIndex: 'operator', width: 100 },
  { title: '开始时间', dataIndex: 'started_at', width: 170 },
  { title: '完成时间', dataIndex: 'finished_at', width: 170 },
  { title: '操作', width: 140, slotName: 'actions', fixed: 'right' }
]

const getStatusColor = (s) => ({
  success: 'green', failed: 'red', pending: 'arcoblue'
}[s] || 'gray')

const getStatusText = (s) => ({
  success: '成功', failed: '失败', pending: '进行中'
}[s] || s)

const handleSearch = () => {
  pagination.current = 1
  loadData()
}

const handleReset = () => {
  searchForm.model_name = ''
  searchForm.status = ''
  pagination.current = 1
  loadData()
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadData()
}

const handleRollback = (record) => {
  rollbackTarget.value = record
  confirmVisible.value = true
}

const executeRollback = () => {
  rolling.value = true
  setTimeout(() => {
    rolling.value = false
    confirmVisible.value = false
    Message.success('回滚成功')
    loadData()
  }, 1500)
}

const viewDetail = (record) => {
  currentRecord.value = record
  detailVisible.value = true
}

const loadData = async () => {
  loading.value = true
  try {
    dataList.value = [
      { id: 'rb_001', model_name: '行为识别模型', from_version: '2.1.0', to_version: '2.0.5', status: 'success', operator: 'admin', started_at: '2026-03-23 14:00:00', finished_at: '2026-03-23 14:00:28', is_current: false, log: '[14:00:00] 开始回滚\n[14:00:05] 备份当前版本配置\n[14:00:12] 切换模型文件\n[14:00:20] 更新路由配置\n[14:00:28] 回滚完成' },
      { id: 'rb_002', model_name: '情感分析模型', from_version: '2.0.5', to_version: '2.0.0', status: 'success', operator: 'admin', started_at: '2026-03-20 10:00:00', finished_at: '2026-03-20 10:00:25', is_current: false, log: '[10:00:00] 开始回滚\n[10:00:10] 切换模型文件\n[10:00:20] 更新路由配置\n[10:00:25] 回滚完成' },
      { id: 'rb_003', model_name: '行为识别模型', from_version: '2.0.5', to_version: '2.1.0', status: 'pending', operator: 'admin', started_at: '-', finished_at: '-', is_current: true, log: '' },
      { id: 'rb_004', model_name: '语音合成模型', from_version: '2.1.0', to_version: '2.0.0', status: 'failed', operator: 'dev', started_at: '2026-03-18 09:00:00', finished_at: '2026-03-18 09:00:45', is_current: false, log: '[09:00:00] 开始回滚\n[09:00:30] 切换模型文件\n[09:00:45] 错误: 模型文件损坏\n[09:00:45] 回滚失败，已恢复原版本' }
    ]
    pagination.total = dataList.value.length
  } catch (e) {
    Message.error('加载失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped lang="less">
.rollback-alert { margin-bottom: 16px; }
.model-name { font-weight: 500; }
</style>
