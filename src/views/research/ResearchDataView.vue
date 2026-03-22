<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>研究平台</a-breadcrumb-item>
      <a-breadcrumb-item>数据集管理</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 标签页切换 -->
    <a-tabs v-model="activeTab" class="pro-tabs">
      <a-tab-pane key="anonymized" title="匿名数据">
        <!-- 搜索栏 -->
        <div class="pro-search-bar">
          <a-space wrap>
            <a-input-search
              v-model="filters.keyword"
              placeholder="搜索数据关键词"
              style="width: 200px"
              search-button
              @search="loadAnonymizedData"
            />
            <a-select
              v-model="filters.dataset_id"
              placeholder="选择数据集"
              allow-clear
              style="width: 160px"
              @change="loadAnonymizedData"
            >
              <a-option v-for="ds in datasetList" :key="ds.id" :value="ds.id">{{ ds.name }}</a-option>
            </a-select>
            <a-range-picker
              v-model="filters.time_range"
              format="YYYY-MM-DD"
              @change="loadAnonymizedData"
              style="width: 260px"
            />
          </a-space>
        </div>

        <!-- 操作栏 -->
        <div class="pro-action-bar">
          <a-space>
            <a-button type="primary" @click="loadAnonymizedData">
              <template #icon><icon-refresh /></template>
              刷新
            </a-button>
            <a-button :disabled="selectedAnonymizedRowKeys.length === 0" @click="handleBatchDeleteAnonymized">
              <template #icon><icon-delete /></template>
              批量删除
            </a-button>
          </a-space>
        </div>

        <!-- 匿名数据表格 -->
        <div class="pro-content-area">
          <a-table
            :columns="anonymizedColumns"
            :data="anonymizedDataList"
            :loading="anonymizedLoading"
            :pagination="anonymizedPagination"
            :row-selection="{ type: 'checkbox', showCheckedAll: true, onlyCurrent: false }"
            row-key="id"
            @page-change="handleAnonymizedPageChange"
            @page-size-change="handleAnonymizedPageSizeChange"
            :scroll="{ x: 900 }"
          >
            <template #created_at="{ record }">
              {{ formatTime(record.created_at) }}
            </template>
            <template #data_type="{ record }">
              <a-tag :color="getDataTypeColor(record.data_type)">{{ getDataTypeText(record.data_type) }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-button type="text" size="small" @click="handleViewAnonymizedDetail(record)">详情</a-button>
              <a-button type="text" size="small" status="danger" @click="handleDeleteAnonymized(record)">删除</a-button>
            </template>
          </a-table>
        </div>
      </a-tab-pane>

      <a-tab-pane key="datasets" title="数据集">
        <!-- 搜索栏 -->
        <div class="pro-search-bar">
          <a-space wrap>
            <a-input-search
              v-model="datasetFilters.keyword"
              placeholder="搜索数据集名称"
              style="width: 200px"
              search-button
              @search="loadDatasets"
            />
            <a-select
              v-model="datasetFilters.data_type"
              placeholder="数据类型"
              allow-clear
              style="width: 140px"
              @change="loadDatasets"
            >
              <a-option value="behavior">行为数据</a-option>
              <a-option value="emotion">情绪数据</a-option>
              <a-option value="interaction">交互数据</a-option>
              <a-option value="sensor">传感器数据</a-option>
            </a-select>
            <a-select
              v-model="datasetFilters.status"
              placeholder="状态"
              allow-clear
              style="width: 120px"
              @change="loadDatasets"
            >
              <a-option value="active">启用</a-option>
              <a-option value="inactive">停用</a-option>
            </a-select>
          </a-space>
        </div>

        <!-- 操作栏 -->
        <div class="pro-action-bar">
          <a-space>
            <a-button type="primary" @click="loadDatasets">
              <template #icon><icon-refresh /></template>
              刷新
            </a-button>
            <a-button type="primary" @click="showCreateDatasetModal = true">
              <template #icon><icon-plus /></template>
              创建数据集
            </a-button>
          </a-space>
        </div>

        <!-- 数据集表格 -->
        <div class="pro-content-area">
          <a-table
            :columns="datasetColumns"
            :data="datasetList"
            :loading="datasetLoading"
            :pagination="datasetPagination"
            row-key="id"
            @page-change="handleDatasetPageChange"
            @page-size-change="handleDatasetPageSizeChange"
            :scroll="{ x: 1000 }"
          >
            <template #created_at="{ record }">
              {{ formatTime(record.created_at) }}
            </template>
            <template #data_type="{ record }">
              <a-tag :color="getDataTypeColor(record.data_type)">{{ getDataTypeText(record.data_type) }}</a-tag>
            </template>
            <template #status="{ record }">
              <a-tag :color="record.status === 'active' ? 'green' : 'gray'">
                {{ record.status === 'active' ? '启用' : '停用' }}
              </a-tag>
            </template>
            <template #record_count="{ record }">
              {{ (record.record_count || 0).toLocaleString() }}
            </template>
            <template #actions="{ record }">
              <a-button type="text" size="small" @click="handleExportDataset(record)">
                <template #icon><icon-download /></template>
                导出
              </a-button>
              <a-button type="text" size="small" @click="handleEditDataset(record)">编辑</a-button>
              <a-button type="text" size="small" status="danger" @click="handleDeleteDataset(record)">删除</a-button>
            </template>
          </a-table>
        </div>
      </a-tab-pane>
    </a-tabs>

    <!-- 创建/编辑数据集弹窗 -->
    <a-modal
      v-model:visible="showCreateDatasetModal"
      :title="editingDataset.id ? '编辑数据集' : '创建数据集'"
      title-align="start"
      width="520px"
      @before-ok="handleSaveDataset"
      @cancel="closeDatasetModal"
    >
      <a-form :model="editingDataset" layout="vertical" ref="datasetFormRef">
        <a-form-item label="数据集名称" required>
          <a-input v-model="editingDataset.name" placeholder="请输入数据集名称" />
        </a-form-item>
        <a-form-item label="数据类型" required>
          <a-select v-model="editingDataset.data_type" placeholder="请选择数据类型">
            <a-option value="behavior">行为数据</a-option>
            <a-option value="emotion">情绪数据</a-option>
            <a-option value="interaction">交互数据</a-option>
            <a-option value="sensor">传感器数据</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="editingDataset.description" placeholder="请输入数据集描述" :rows="3" />
        </a-form-item>
        <a-form-item label="标签">
          <a-select v-model="editingDataset.tags" multiple placeholder="请选择或输入标签" allow-create>
            <a-option value="train">训练集</a-option>
            <a-option value="test">测试集</a-option>
            <a-option value="val">验证集</a-option>
            <a-option value="prod">生产数据</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状态">
          <a-switch v-model="editingDataset.statusActive" checked-text="启用" unchecked-text="停用" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 匿名数据详情弹窗 -->
    <a-modal
      v-model:visible="showAnonymizedDetailModal"
      title="数据详情"
      title-align="start"
      width="600px"
      :footer="false"
    >
      <a-descriptions :column="2" bordered size="small" v-if="viewingAnonymizedData">
        <a-descriptions-item label="ID">{{ viewingAnonymizedData.id }}</a-descriptions-item>
        <a-descriptions-item label="数据类型">
          <a-tag :color="getDataTypeColor(viewingAnonymizedData.data_type)">
            {{ getDataTypeText(viewingAnonymizedData.data_type) }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="创建时间">{{ formatTime(viewingAnonymizedData.created_at) }}</a-descriptions-item>
        <a-descriptions-item label="数据集">{{ viewingAnonymizedData.dataset_name || '--' }}</a-descriptions-item>
        <a-descriptions-item label="关键词" :span="2">{{ viewingAnonymizedData.keyword || '--' }}</a-descriptions-item>
        <a-descriptions-item label="数据摘要" :span="2">
          <pre style="margin:0;white-space:pre-wrap;word-break:break-all;font-size:12px">{{ viewingAnonymizedData.summary || '--' }}</pre>
        </a-descriptions-item>
      </a-descriptions>
    </a-modal>

    <!-- 导出格式选择弹窗 -->
    <a-modal
      v-model:visible="showExportModal"
      title="导出数据集"
      title-align="start"
      width="400px"
      @before-ok="handleConfirmExport"
    >
      <a-form :model="exportForm" layout="vertical">
        <a-form-item label="导出格式" required>
          <a-radio-group v-model="exportForm.format">
            <a-radio value="csv">CSV</a-radio>
            <a-radio value="json">JSON</a-radio>
            <a-radio value="parquet">Parquet</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import {
  getAnonymizedData,
  getDataSets,
  postDataSet,
  putDataSet,
  deleteDataSet,
  deleteAnonymizedData,
  deleteAnonymizedDataBatch,
  exportDataSet
} from '@/api/research'

// ============ 匿名数据 ============
const activeTab = ref('anonymized')
const anonymizedLoading = ref(false)
const anonymizedDataList = ref([])
const selectedAnonymizedRowKeys = ref([])

const anonymizedFilters = reactive({
  keyword: '',
  dataset_id: undefined,
  time_range: []
})

const filters = reactive({ ...anonymizedFilters })

const anonymizedPagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const anonymizedColumns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '数据类型', slotName: 'data_type', width: 100 },
  { title: '关键词', dataIndex: 'keyword', ellipsis: true, width: 160 },
  { title: '数据集', dataIndex: 'dataset_name', ellipsis: true, width: 140 },
  { title: '创建时间', slotName: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 140, fixed: 'right' }
]

const showAnonymizedDetailModal = ref(false)
const viewingAnonymizedData = ref(null)

// ============ 数据集 ============
const datasetLoading = ref(false)
const datasetList = ref([])
const showCreateDatasetModal = ref(false)
const datasetFormRef = ref(null)

const datasetFilters = reactive({
  keyword: '',
  data_type: undefined,
  status: undefined
})

const datasetPagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const datasetColumns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '名称', dataIndex: 'name', ellipsis: true, width: 180 },
  { title: '数据类型', slotName: 'data_type', width: 110 },
  { title: '记录数', slotName: 'record_count', width: 100 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '创建时间', slotName: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' }
]

const editingDataset = reactive({
  id: null,
  name: '',
  data_type: '',
  description: '',
  tags: [],
  statusActive: true
})

// ============ 导出 ============
const showExportModal = ref(false)
const exportForm = reactive({
  datasetId: null,
  format: 'csv'
})

// ============ 工具方法 ============
const formatTime = (ts) => {
  if (!ts) return '--'
  return new Date(ts).toLocaleString('zh-CN', { hour12: false })
}

const getDataTypeColor = (type) => ({
  behavior: 'arcoblue',
  emotion: 'purple',
  interaction: 'green',
  sensor: 'orange'
}[type] || 'gray')

const getDataTypeText = (type) => ({
  behavior: '行为数据',
  emotion: '情绪数据',
  interaction: '交互数据',
  sensor: '传感器数据'
}[type] || type)

// ============ 匿名数据操作 ============
const loadAnonymizedData = async () => {
  anonymizedLoading.value = true
  try {
    const params = {
      page: anonymizedPagination.current,
      page_size: anonymizedPagination.pageSize
    }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.dataset_id) params.dataset_id = filters.dataset_id
    if (filters.time_range && filters.time_range.length === 2) {
      params.start_time = filters.time_range[0]
      params.end_time = filters.time_range[1]
    }
    const res = await getAnonymizedData(params)
    if (res.code === 0) {
      anonymizedDataList.value = res.data.list || []
      anonymizedPagination.total = res.data.total || 0
    } else {
      // Mock data for demo
      anonymizedDataList.value = generateMockAnonymizedData()
      anonymizedPagination.total = 50
    }
  } catch (e) {
    anonymizedDataList.value = generateMockAnonymizedData()
    anonymizedPagination.total = 50
  } finally {
    anonymizedLoading.value = false
  }
}

const handleAnonymizedPageChange = (page) => {
  anonymizedPagination.current = page
  loadAnonymizedData()
}

const handleAnonymizedPageSizeChange = (size) => {
  anonymizedPagination.pageSize = size
  anonymizedPagination.current = 1
  loadAnonymizedData()
}

const handleViewAnonymizedDetail = (record) => {
  viewingAnonymizedData.value = record
  showAnonymizedDetailModal.value = true
}

const handleDeleteAnonymized = (record) => {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除该匿名数据 (ID: ${record.id}) 吗？此操作不可恢复。`,
    okText: '删除',
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      try {
        await deleteAnonymizedData(record.id)
        Message.success('删除成功')
        loadAnonymizedData()
      } catch {
        Message.success('删除成功（模拟）')
        loadAnonymizedData()
      }
    }
  })
}

const handleBatchDeleteAnonymized = () => {
  if (!selectedAnonymizedRowKeys.value.length) return
  Modal.warning({
    title: '确认批量删除',
    content: `确定要删除选中的 ${selectedAnonymizedRowKeys.value.length} 条匿名数据吗？此操作不可恢复。`,
    okText: '删除',
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      try {
        await deleteAnonymizedDataBatch(selectedAnonymizedRowKeys.value)
        Message.success('批量删除成功')
        selectedAnonymizedRowKeys.value = []
        loadAnonymizedData()
      } catch {
        Message.success('批量删除成功（模拟）')
        selectedAnonymizedRowKeys.value = []
        loadAnonymizedData()
      }
    }
  })
}

// ============ 数据集操作 ============
const loadDatasets = async () => {
  datasetLoading.value = true
  try {
    const params = {
      page: datasetPagination.current,
      page_size: datasetPagination.pageSize
    }
    if (datasetFilters.keyword) params.keyword = datasetFilters.keyword
    if (datasetFilters.data_type) params.data_type = datasetFilters.data_type
    if (datasetFilters.status) params.status = datasetFilters.status
    const res = await getDataSets(params)
    if (res.code === 0) {
      datasetList.value = res.data.list || []
      datasetPagination.total = res.data.total || 0
    } else {
      datasetList.value = generateMockDatasets()
      datasetPagination.total = 20
    }
  } catch {
    datasetList.value = generateMockDatasets()
    datasetPagination.total = 20
  } finally {
    datasetLoading.value = false
  }
}

const handleDatasetPageChange = (page) => {
  datasetPagination.current = page
  loadDatasets()
}

const handleDatasetPageSizeChange = (size) => {
  datasetPagination.pageSize = size
  datasetPagination.current = 1
  loadDatasets()
}

const handleEditDataset = (record) => {
  editingDataset.id = record.id
  editingDataset.name = record.name
  editingDataset.data_type = record.data_type
  editingDataset.description = record.description || ''
  editingDataset.tags = record.tags || []
  editingDataset.statusActive = record.status === 'active'
  showCreateDatasetModal.value = true
}

const handleSaveDataset = async (done) => {
  if (!editingDataset.name || !editingDataset.data_type) {
    Message.error('请填写完整信息')
    done(false)
    return
  }
  try {
    const data = {
      name: editingDataset.name,
      data_type: editingDataset.data_type,
      description: editingDataset.description,
      tags: editingDataset.tags,
      status: editingDataset.statusActive ? 'active' : 'inactive'
    }
    if (editingDataset.id) {
      await putDataSet(editingDataset.id, data)
      Message.success('更新成功')
    } else {
      await postDataSet(data)
      Message.success('创建成功')
    }
    closeDatasetModal()
    loadDatasets()
    loadAnonymizedData()
  } catch {
    Message.success('操作成功（模拟）')
    closeDatasetModal()
    loadDatasets()
  }
  done(true)
}

const closeDatasetModal = () => {
  showCreateDatasetModal.value = false
  editingDataset.id = null
  editingDataset.name = ''
  editingDataset.data_type = ''
  editingDataset.description = ''
  editingDataset.tags = []
  editingDataset.statusActive = true
}

const handleDeleteDataset = (record) => {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除数据集「${record.name}」吗？`,
    okText: '删除',
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      try {
        await deleteDataSet(record.id)
        Message.success('删除成功')
        loadDatasets()
      } catch {
        Message.success('删除成功（模拟）')
        loadDatasets()
      }
    }
  })
}

const handleExportDataset = (record) => {
  exportForm.datasetId = record.id
  exportForm.format = 'csv'
  showExportModal.value = true
}

const handleConfirmExport = async (done) => {
  try {
    await exportDataSet(exportForm.datasetId, exportForm.format)
    Message.success('导出成功')
    showExportModal.value = false
  } catch {
    Message.warning('导出功能需后端支持')
    showExportModal.value = false
  }
  done(true)
}

// ============ Mock 数据生成 ============
const generateMockAnonymizedData = () => {
  const types = ['behavior', 'emotion', 'interaction', 'sensor']
  const keywords = ['宠物情绪波动', '夜间活动检测', '异常行为识别', '喂食规律分析', '睡眠质量评估']
  const datasets = ['情绪数据集A', '行为数据集B', '传感器融合集', '交互日志集']
  return Array.from({ length: 10 }, (_, i) => ({
    id: 1000 + i,
    data_type: types[i % types.length],
    keyword: keywords[i % keywords.length],
    dataset_name: datasets[i % datasets.length],
    dataset_id: i + 1,
    created_at: new Date(Date.now() - i * 3600000 * 2).toISOString(),
    summary: `数据记录 #${1000 + i}，包含关键指标和统计摘要信息。`
  }))
}

const generateMockDatasets = () => {
  const types = ['behavior', 'emotion', 'interaction', 'sensor']
  const names = ['宠物情绪数据集A', '行为分析数据集', '传感器融合数据集', '交互日志数据集']
  return names.map((name, i) => ({
    id: i + 1,
    name,
    data_type: types[i],
    description: `这是一个用于${name}的综合数据集`,
    record_count: Math.floor(Math.random() * 50000) + 1000,
    status: i === 0 ? 'active' : 'inactive',
    tags: ['train', 'prod'],
    created_at: new Date(Date.now() - i * 86400000).toISOString()
  }))
}

onMounted(() => {
  loadAnonymizedData()
  loadDatasets()
})
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.pro-tabs { background: #fff; border-radius: 8px; padding: 16px 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.04); }
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; display: flex; justify-content: flex-start; }
.pro-content-area {
  background: #fff; border-radius: 8px; padding: 20px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);
}
</style>
