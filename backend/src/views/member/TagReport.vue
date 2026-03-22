<template>
  <div class="tag-report-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>
        <router-link to="/member/tags">会员标签</router-link>
      </a-breadcrumb-item>
      <a-breadcrumb-item>标签报表</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计概览 -->
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card hoverable>
          <a-statistic title="标签总数" :value="stats.totalTags || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card hoverable>
          <a-statistic title="已打标签会员" :value="stats.taggedMembers || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card hoverable>
          <a-statistic title="本月新增" :value="stats.monthlyAdded || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card hoverable>
          <a-statistic title="本月清除" :value="stats.monthlyCleaned || 0" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 标签分布 -->
    <a-card title="标签分布（会员数 TOP 10）" style="margin-bottom: 16px;">
      <a-table :columns="tagColumns" :data="tagDistribution" :loading="loading" row-key="tagId" :pagination="false">
        <template #tagName="{ record }">
          <a-tag :color="record.tagColor || 'arcoblue'">{{ record.tagName || '标签' + record.tagId }}</a-tag>
        </template>
        <template #memberCount="{ record }">
          <span style="color: #165DFF; font-weight: 600;">{{ record.memberCount || 0 }}</span>
        </template>
        <template #percentage="{ record }">
          <a-progress :percent="record.percentage || 0" :show-text="true" :format="(p) => p.toFixed(1) + '%'" size="small" />
        </template>
      </a-table>
    </a-card>

    <!-- 标签操作流水 -->
    <a-card title="标签操作流水">
      <template #extra>
        <a-space>
          <a-select v-model="filterTagId" placeholder="筛选标签" style="width: 160px;" allow-clear>
            <a-option v-for="tag in tagList" :key="tag.id" :value="tag.id">{{ tag.name }}</a-option>
          </a-select>
          <a-range-picker v-model="filterDate" style="width: 240px;" />
          <a-button @click="loadFlow">筛选</a-button>
          <a-button @click="exportReport">导出</a-button>
        </a-space>
      </template>
      <a-table :columns="flowColumns" :data="flowList" :loading="flowLoading" row-key="id" :pagination="{ pageSize: 10 }">
        <template #actionType="{ record }">
          <a-tag :color="getActionColor(record.actionType)">{{ getActionName(record.actionType) }}</a-tag>
        </template>
        <template #createTime="{ record }">
          {{ formatTime(record.createTime) }}
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as api from '@/api/marketing'

const tagDistribution = ref([])
const flowList = ref([])
const tagList = ref([])
const stats = ref({})
const loading = ref(false)
const flowLoading = ref(false)
const filterTagId = ref(null)
const filterDate = ref([])

const tagColumns = [
  { title: '标签名称', slotName: 'tagName', width: 180 },
  { title: '分类', dataIndex: 'category', width: 100 },
  { title: '会员数', slotName: 'memberCount', width: 120 },
  { title: '占比', slotName: 'percentage' }
]

const flowColumns = [
  { title: '操作类型', slotName: 'actionType', width: 120 },
  { title: '会员', dataIndex: 'memberName', width: 120 },
  { title: '标签', dataIndex: 'tagName', width: 150 },
  { title: '操作人', dataIndex: 'operator', width: 120 },
  { title: '时间', slotName: 'createTime', width: 180 }
]

const getActionColor = (type) => {
  const map = { add: 'green', remove: 'red', auto_add: 'blue', auto_remove: 'orange' }
  return map[type] || 'gray'
}

const getActionName = (type) => {
  const map = { add: '手动打标签', remove: '手动移除', auto_add: '自动打标签', auto_remove: '自动清除' }
  return map[type] || type
}

const formatTime = (ts) => {
  if (!ts) return '-'
  return new Date(ts).toLocaleString('zh-CN')
}

const loadDistribution = async () => {
  loading.value = true
  try {
    const res = await api.getTagReport({ type: 'distribution' })
    tagDistribution.value = res.data?.list || []
    stats.value = res.data?.stats || {}
  } catch (err) {
    Message.error('加载报表失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

const loadFlow = async () => {
  flowLoading.value = true
  try {
    const params = {}
    if (filterTagId.value) params.tagId = filterTagId.value
    if (filterDate.value?.length === 2) {
      params.startDate = filterDate.value[0]
      params.endDate = filterDate.value[1]
    }
    const res = await api.getTagReport({ type: 'flow', ...params })
    flowList.value = res.data?.list || []
  } catch (err) {
    Message.error('加载流水失败: ' + err.message)
  } finally {
    flowLoading.value = false
  }
}

const exportReport = () => {
  Message.info('正在导出报表...')
}

onMounted(async () => {
  await loadDistribution()
  try {
    const tagsRes = await api.getTagList()
    tagList.value = tagsRes.data || []
  } catch (err) { /* ignore */ }
  await loadFlow()
})
</script>

<style scoped>
.tag-report-page {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}
.breadcrumb { margin-bottom: 16px; }
</style>
