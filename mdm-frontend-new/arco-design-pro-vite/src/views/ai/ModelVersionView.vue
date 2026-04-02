<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>AI 功能</a-breadcrumb-item>
      <a-breadcrumb-item>模型版本管理</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 操作栏 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showAddModal = true">新建版本</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </div>

    <!-- 搜索表单 -->
    <div class="pro-search-bar">
      <a-form :model="searchForm" layout="inline">
        <a-form-item label="模型名称">
          <a-input v-model="searchForm.model_name" placeholder="搜索模型名称" allow-clear style="width: 180px" />
        </a-form-item>
        <a-form-item label="版本">
          <a-input v-model="searchForm.version" placeholder="如 2.1.0" allow-clear style="width: 120px" />
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model="searchForm.status" placeholder="选择状态" allow-clear style="width: 120px">
            <a-option value="published">已发布</a-option>
            <a-option value="testing">测试中</a-option>
            <a-option value="draft">草稿</a-option>
            <a-option value="deprecated">已废弃</a-option>
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

    <!-- 版本列表 -->
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
        <template #status="{ record }">
          <a-badge :status="getStatusBadge(record.status)" />
          <span>{{ getStatusText(record.status) }}</span>
        </template>
        <template #accuracy="{ record }">
          <a-progress :percent="record.accuracy" size="small" :show-text="true"
            :color="record.accuracy >= 95 ? '#00b42a' : record.accuracy >= 90 ? '#ff7d00' : '#f53f3f'" />
        </template>
        <template #is_current="{ record }">
          <a-tag v-if="record.is_current" color="green">当前版本</a-tag>
          <span v-else>-</span>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="viewDetail(record)">详情</a-button>
            <a-button type="text" size="small" @click="compare(record)" :disabled="record.status === 'deprecated'">对比</a-button>
            <a-button type="text" size="small" @click="rollback(record)" :disabled="record.is_current">回滚</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 新建版本弹窗 -->
    <a-modal v-model:visible="showAddModal" title="新建模型版本" :width="640" @before-ok="handleAddSubmit">
      <a-form :model="addForm" layout="vertical">
        <a-form-item label="模型名称" required>
          <a-select v-model="addForm.model_name" placeholder="选择模型">
            <a-option value="behavior">行为识别模型</a-option>
            <a-option value="emotion">情感分析模型</a-option>
            <a-option value="voice">语音合成模型</a-option>
            <a-option value="pose">姿态估计模型</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="版本号" required>
          <a-input v-model="addForm.version" placeholder="如 2.2.0" />
        </a-form-item>
        <a-form-item label="模型文件">
          <a-upload action="/" :limit="1" />
        </a-form-item>
        <a-form-item label="配置参数">
          <a-textarea v-model="addForm.config" :rows="4" placeholder="JSON 配置参数" />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="addForm.remark" :rows="2" placeholder="版本说明" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 详情弹窗 -->
    <a-modal v-model:visible="detailVisible" title="版本详情" :width="700" footer="null">
      <a-descriptions :column="2" bordered size="large">
        <a-descriptions-item label="模型ID">{{ currentRecord?.id }}</a-descriptions-item>
        <a-descriptions-item label="版本">v{{ currentRecord?.version }}</a-descriptions-item>
        <a-descriptions-item label="模型名称">{{ currentRecord?.model_name }}</a-descriptions-item>
        <a-descriptions-item label="状态">{{ getStatusText(currentRecord?.status) }}</a-descriptions-item>
        <a-descriptions-item label="准确率">{{ currentRecord?.accuracy }}%</a-descriptions-item>
        <a-descriptions-item label="是否为当前版本">{{ currentRecord?.is_current ? '是' : '否' }}</a-descriptions-item>
        <a-descriptions-item label="发布人">{{ currentRecord?.publisher }}</a-descriptions-item>
        <a-descriptions-item label="发布时间">{{ currentRecord?.published_at }}</a-descriptions-item>
        <a-descriptions-item label="变更说明" :span="2">{{ currentRecord?.changelog }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>

    <!-- 对比弹窗 -->
    <a-modal v-model:visible="compareVisible" title="版本对比" :width="900">
      <a-row :gutter="16">
        <a-col :span="12">
          <a-card title="版本 A" size="small">
            <a-descriptions :column="1" size="small">
              <a-descriptions-item label="版本">v{{ compareRecordA?.version }}</a-descriptions-item>
              <a-descriptions-item label="准确率">{{ compareRecordA?.accuracy }}%</a-descriptions-item>
              <a-descriptions-item label="延迟">{{ compareRecordA?.latency_ms }}ms</a-descriptions-item>
              <a-descriptions-item label="状态">{{ getStatusText(compareRecordA?.status) }}</a-descriptions-item>
            </a-descriptions>
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card title="版本 B" size="small">
            <a-descriptions :column="1" size="small">
              <a-descriptions-item label="版本">v{{ compareRecordB?.version }}</a-descriptions-item>
              <a-descriptions-item label="准确率">{{ compareRecordB?.accuracy }}%</a-descriptions-item>
              <a-descriptions-item label="延迟">{{ compareRecordB?.latency_ms }}ms</a-descriptions-item>
              <a-descriptions-item label="状态">{{ getStatusText(compareRecordB?.status) }}</a-descriptions-item>
            </a-descriptions>
          </a-card>
        </a-col>
      </a-row>
      <template #footer>
        <a-button @click="compareVisible = false">关闭</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const dataList = ref([])
const showAddModal = ref(false)
const detailVisible = ref(false)
const compareVisible = ref(false)
const currentRecord = ref(null)
const compareRecordA = ref(null)
const compareRecordB = ref(null)

const searchForm = reactive({
  model_name: '',
  version: '',
  status: ''
})

const addForm = reactive({
  model_name: '',
  version: '',
  config: '',
  remark: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '模型名称', dataIndex: 'model_name', width: 160 },
  { title: '版本', dataIndex: 'version', width: 100 },
  { title: '状态', width: 110, slotName: 'status' },
  { title: '准确率', dataIndex: 'accuracy', width: 160, slotName: 'accuracy' },
  { title: '延迟', dataIndex: 'latency_ms', width: 90 },
  { title: '当前版本', width: 110, slotName: 'is_current' },
  { title: '发布人', dataIndex: 'publisher', width: 100 },
  { title: '发布时间', dataIndex: 'published_at', width: 170 },
  { title: '操作', width: 160, slotName: 'actions', fixed: 'right' }
]

const getStatusBadge = (s) => ({
  published: 'normal', testing: 'processing', draft: 'default', deprecated: 'error'
}[s] || 'default')

const getStatusText = (s) => ({
  published: '已发布', testing: '测试中', draft: '草稿', deprecated: '已废弃'
}[s] || s)

const handleSearch = () => {
  pagination.current = 1
  loadData()
}

const handleReset = () => {
  searchForm.model_name = ''
  searchForm.version = ''
  searchForm.status = ''
  pagination.current = 1
  loadData()
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadData()
}

const handleAddSubmit = (done) => {
  if (!addForm.model_name || !addForm.version) {
    Message.error('请填写必填项')
    done(false)
    return
  }
  showAddModal.value = false
  Message.success('版本创建成功')
  loadData()
  done(true)
}

const viewDetail = (record) => {
  currentRecord.value = record
  detailVisible.value = true
}

const compare = (record) => {
  compareRecordA.value = record
  compareRecordB.value = dataList.value.find(v => v.version === record.version && v.id !== record.id) || dataList.value[0]
  compareVisible.value = true
}

const rollback = (record) => {
  Message.info(`回滚到 v${record.version} 功能开发中`)
}

const loadData = async () => {
  loading.value = true
  try {
    dataList.value = [
      { id: 'mv_001', model_name: '行为识别模型', version: '2.1.0', status: 'published', accuracy: 96.5, latency_ms: 98, is_current: true, publisher: 'admin', published_at: '2026-03-20 10:00:00', changelog: '优化移动检测精度' },
      { id: 'mv_002', model_name: '行为识别模型', version: '2.0.5', status: 'published', accuracy: 95.2, latency_ms: 105, is_current: false, publisher: 'admin', published_at: '2026-03-15 14:30:00', changelog: '修复边界情况' },
      { id: 'mv_003', model_name: '情感分析模型', version: '2.0.5', status: 'published', accuracy: 94.2, latency_ms: 145, is_current: true, publisher: 'admin', published_at: '2026-03-15 14:30:00', changelog: '新增恐惧情绪识别' },
      { id: 'mv_004', model_name: '情感分析模型', version: '2.0.0', status: 'deprecated', accuracy: 92.8, latency_ms: 160, is_current: false, publisher: 'admin', published_at: '2026-03-10 09:00:00', changelog: '初始版本' },
      { id: 'mv_005', model_name: '语音合成模型', version: '2.2.0', status: 'testing', accuracy: 97.1, latency_ms: 80, is_current: false, publisher: 'dev', published_at: '-', changelog: '支持多语言' },
      { id: 'mv_006', model_name: '姿态估计模型', version: '1.9.0', status: 'deprecated', accuracy: 89.5, latency_ms: 320, is_current: false, publisher: 'admin', published_at: '2026-02-20 11:00:00', changelog: '已废弃' }
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
</style>
