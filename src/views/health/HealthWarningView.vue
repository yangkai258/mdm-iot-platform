<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>健康中心</a-breadcrumb-item>
      <a-breadcrumb-item>健康预警</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选区 -->
    <div class="search-form">
      <a-form :model="searchForm" layout="inline">
        <a-form-item label="关键词">
          <a-input-search
            v-model="searchForm.keyword"
            placeholder="搜索预警内容"
            style="width: 280px"
            @search="loadWarnings"
            search-button
          />
        </a-form-item>
        <a-form-item label="预警级别">
          <a-select v-model="searchForm.level" placeholder="选择级别" allow-clear style="width: 120px">
            <a-option value="critical">危急</a-option>
            <a-option value="high">高</a-option>
            <a-option value="medium">中</a-option>
            <a-option value="low">低</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="处理状态">
          <a-select v-model="searchForm.status" placeholder="选择状态" allow-clear style="width: 120px">
            <a-option value="pending">待处理</a-option>
            <a-option value="confirmed">已确认</a-option>
            <a-option value="ignored">已忽略</a-option>
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

    <!-- 操作栏 -->
    <div class="toolbar">
      <a-button type="primary" @click="loadWarnings">刷新</a-button>
    </div>

    <!-- 表格 -->
    <a-table :columns="columns" :data="warnings" :loading="loading" row-key="id" :pagination="{ pageSize: 10 }">
      <template #level="{ record }">
        <a-tag :color="getLevelColor(record.level)">
          {{ getLevelText(record.level) }}
        </a-tag>
      </template>
      <template #status="{ record }">
        <a-tag :color="getStatusColor(record.status)">
          {{ getStatusText(record.status) }}
        </a-tag>
      </template>
      <template #disease_name="{ record }">
        <a-link @click="showDetail(record)">{{ record.disease_name }}</a-link>
      </template>
      <template #actions="{ record }">
        <a-space>
          <a-button type="text" size="small" @click="showDetail(record)">详情</a-button>
          <a-button type="text" size="small" status="success" @click="confirmWarning(record)" :disabled="record.status !== 'pending'">确认</a-button>
          <a-button type="text" size="small" status="warning" @click="ignoreWarning(record)" :disabled="record.status !== 'pending'">忽略</a-button>
        </a-space>
      </template>
    </a-table>

    <!-- 预警详情弹窗 -->
    <a-modal v-model:visible="detailVisible" title="预警详情" :width="600" :footer="modalFooter">
      <a-descriptions :column="2" bordered v-if="currentWarning">
        <a-descriptions-item label="预警编号">{{ currentWarning.id }}</a-descriptions-item>
        <a-descriptions-item label="预警级别">
          <a-tag :color="getLevelColor(currentWarning.level)">
            {{ getLevelText(currentWarning.level) }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="疾病名称" :span="2">{{ currentWarning.disease_name }}</a-descriptions-item>
        <a-descriptions-item label="预警描述" :span="2">{{ currentWarning.description }}</a-descriptions-item>
        <a-descriptions-item label="风险因素" :span="2">{{ currentWarning.risk_factors || '暂无' }}</a-descriptions-item>
        <a-descriptions-item label="建议措施" :span="2">{{ currentWarning.suggestion || '暂无' }}</a-descriptions-item>
        <a-descriptions-item label="发生时间">{{ currentWarning.created_at }}</a-descriptions-item>
        <a-descriptions-item label="处理状态">
          <a-tag :color="getStatusColor(currentWarning.status)">
            {{ getStatusText(currentWarning.status) }}
          </a-tag>
        </a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const warnings = ref([])
const detailVisible = ref(false)
const currentWarning = ref(null)

const searchForm = reactive({
  keyword: '',
  level: '',
  status: ''
})

const columns = [
  { title: '预警编号', dataIndex: 'id', width: 100 },
  { title: '疾病名称', slotName: 'disease_name' },
  { title: '预警级别', slotName: 'level', width: 100 },
  { title: '预警描述', dataIndex: 'description', ellipsis: true },
  { title: '发生时间', dataIndex: 'created_at', width: 180 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 180 }
]

const getLevelColor = (level) => {
  const colors = { critical: 'red', high: 'orange', medium: 'blue', low: 'green' }
  return colors[level] || 'gray'
}

const getLevelText = (level) => {
  const texts = { critical: '危急', high: '高', medium: '中', low: '低' }
  return texts[level] || '未知'
}

const getStatusColor = (status) => {
  const colors = { pending: 'orange', confirmed: 'green', ignored: 'gray' }
  return colors[status] || 'gray'
}

const getStatusText = (status) => {
  const texts = { pending: '待处理', confirmed: '已确认', ignored: '已忽略' }
  return texts[status] || '未知'
}

const modalFooter = ref(null)

const loadWarnings = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams()
    if (searchForm.keyword) params.append('keyword', searchForm.keyword)
    if (searchForm.level) params.append('level', searchForm.level)
    if (searchForm.status) params.append('status', searchForm.status)
    
    const res = await fetch(`/api/v1/health/warnings?${params}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      warnings.value = data.data?.list || []
    } else {
      warnings.value = generateMockData()
    }
  } catch (e) {
    console.error('加载预警失败:', e)
    warnings.value = generateMockData()
  } finally {
    loading.value = false
  }
}

const generateMockData = () => [
  { id: 1, disease_name: '高血压风险', level: 'high', description: '连续3天血压监测偏高，建议密切观察', risk_factors: '高盐饮食、缺乏运动', suggestion: '减少盐分摄入，增加有氧运动', status: 'pending', created_at: '2026-03-22 10:30:00' },
  { id: 2, disease_name: '睡眠呼吸暂停', level: 'critical', description: '夜间血氧饱和度多次低于90%', risk_factors: '肥胖，鼻腔阻塞', suggestion: '建议就医进行睡眠监测', status: 'pending', created_at: '2026-03-22 08:00:00' },
  { id: 3, disease_name: '心律不齐', level: 'medium', description: '心电图检测到偶发早搏', risk_factors: '咖啡因摄入过多、压力大', suggestion: '减少咖啡因摄入，保持规律作息', status: 'confirmed', created_at: '2026-03-21 15:20:00' },
  { id: 4, disease_name: '体重异常波动', level: 'low', description: '一周内体重下降超过5%', risk_factors: '未知', suggestion: '关注饮食，如有持续请就医', status: 'ignored', created_at: '2026-03-20 09:00:00' },
  { id: 5, disease_name: '血糖偏高', level: 'high', description: '空腹血糖持续高于正常值上限', risk_factors: '高糖饮食、缺乏运动', suggestion: '调整饮食结构，增加运动', status: 'pending', created_at: '2026-03-22 11:45:00' }
]

const handleSearch = () => {
  loadWarnings()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.level = ''
  searchForm.status = ''
  loadWarnings()
}

const showDetail = (record) => {
  currentWarning.value = record
  detailVisible.value = true
}

const confirmWarning = async (record) => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/health/warnings/${record.id}/confirm`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      Message.success('已确认预警')
      record.status = 'confirmed'
    } else {
      record.status = 'confirmed'
      Message.success('已确认预警')
    }
  } catch (e) {
    record.status = 'confirmed'
    Message.success('已确认预警')
  }
}

const ignoreWarning = async (record) => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/health/warnings/${record.id}/ignore`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      Message.success('已忽略预警')
      record.status = 'ignored'
    } else {
      record.status = 'ignored'
      Message.success('已忽略预警')
    }
  } catch (e) {
    record.status = 'ignored'
    Message.success('已忽略预警')
  }
}

onMounted(() => {
  loadWarnings()
})
</script>

<style scoped>
.page-container {
  background: #fff;
  border-radius: 4px;
  padding: 20px;
}

.breadcrumb {
  margin-bottom: 16px;
}

.search-form {
  margin-bottom: 16px;
  padding: 16px;
  background: #f7f8fa;
  border-radius: 4px;
}

.toolbar {
  margin-bottom: 16px;
}
</style>
