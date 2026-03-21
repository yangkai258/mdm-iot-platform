<template>
  <div class="member-promotions">
    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card>
          <a-statistic title="活动总数" :value="stats.total" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="进行中" :value="stats.running" :value-style="{ color: '#52c41a' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="已结束" :value="stats.ended" :value-style="{ color: '#999' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="参与人数" :value="stats.participants" :value-style="{ color: '#1890ff' }" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 操作区域 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-select v-model="filters.status" placeholder="活动状态" allow-clear style="width: 120px" @change="loadPromotions">
          <a-option value="active">进行中</a-option>
          <a-option value="ended">已结束</a-option>
          <a-option value="draft">草稿</a-option>
        </a-select>
        <a-input-search v-model="filters.keyword" placeholder="搜索活动名称" style="width: 200px" search-button @search="loadPromotions" />
        <a-button type="primary" @click="showCreateDrawer = true">创建活动</a-button>
        <a-button @click="loadPromotions">刷新</a-button>
      </a-space>
    </a-card>

    <!-- 促销活动列表 -->
    <a-card>
      <a-table :columns="columns" :data="promotionList" :loading="loading" :pagination="pagination" row-key="id">
        <template #type="{ record }">
          <a-tag>{{ getTypeText(record.type) }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="viewDetail(record)">详情</a-button>
            <a-button type="text" size="small" @click="editPromotion(record)">编辑</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 创建活动抽屉 -->
    <a-drawer v-model:visible="showCreateDrawer" title="创建活动" :width="480">
      <a-form :model="form" layout="vertical">
        <a-form-item label="活动名称" required>
          <a-input v-model="form.name" placeholder="请输入活动名称" />
        </a-form-item>
        <a-form-item label="活动类型" required>
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="points_double">双倍积分</a-option>
            <a-option value="discount">折扣活动</a-option>
            <a-option value="gift">赠品活动</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleCreate">创建</a-button>
        </a-form-item>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const promotionList = ref([])
const showCreateDrawer = ref(false)

const filters = reactive({ status: undefined, keyword: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const stats = reactive({ total: 0, running: 0, ended: 0, participants: 0 })
const form = reactive({ name: '', type: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '活动名称', dataIndex: 'name' },
  { title: '类型', slotName: 'type', width: 120 },
  { title: '开始时间', dataIndex: 'start_time', width: 160 },
  { title: '结束时间', dataIndex: 'end_time', width: 160 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 160 }
]

const getTypeText = (type) => ({ points_double: '双倍积分', discount: '折扣活动', gift: '赠品活动' }[type] || type)
const getStatusText = (s) => ({ active: '进行中', ended: '已结束', draft: '草稿' }[s] || s)
const getStatusColor = (s) => ({ active: 'green', ended: 'gray', draft: 'orange' }[s] || 'gray')

const loadPromotions = () => { loading.value = true; setTimeout(() => { loading.value = false }, 300) }
const viewDetail = (r) => Message.info('查看详情')
const editPromotion = (r) => Message.info('编辑活动')
const handleCreate = () => { Message.success('创建成功'); showCreateDrawer.value = false }

onMounted(() => loadPromotions())
</script>

<style scoped>
.member-promotions { padding: 16px; }
.stats-row { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
