<template>
  <div class="page-container">
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
              <a-statistic title="参与人数" :value="stats.participants" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="累计优惠" :value="stats.totalDiscount" :value-style="{ color: '#ff6b00' }" />
            </a-card>
          </a-col>
        </a-row>

        <!-- 操作栏 -->
        <a-card class="action-card">
          <a-space wrap>
            <a-select v-model="filters.promotionType" placeholder="活动类型" allow-clear style="width: 140px" @change="loadPromotions">
              <a-option value="points_multi">双倍积分</a-option>
              <a-option value="discount">折扣活动</a-option>
              <a-option value="gift">赠品活动</a-option>
              <a-option value="coupon">优惠券发放</a-option>
            </a-select>
            <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 120px" @change="loadPromotions">
              <a-option value="active">进行中</a-option>
              <a-option value="upcoming">即将开始</a-option>
              <a-option value="ended">已结束</a-option>
            </a-select>
            <a-input-search v-model="filters.keyword" placeholder="搜索活动名称" style="width: 200px" search-button @search="loadPromotions" />
            <a-button type="primary" @click="showCreateDrawer = true">创建活动</a-button>
            <a-button @click="loadPromotions">刷新</a-button>
          </a-space>
        </a-card>

        <!-- 活动列表 -->
        <a-card class="promotion-card">
          <a-table
            :columns="columns"
            :data="promotionList"
            :loading="loading"
            :pagination="pagination"
            row-key="id"
            @page-change="handlePageChange"
            @page-size-change="handlePageSizeChange"
          >
            <template #promotionType="{ record }">
              <a-tag :color="getTypeColor(record.promotion_type)">{{ getTypeText(record.promotion_type) }}</a-tag>
            </template>
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
            </template>
            <template #participants="{ record }">
              {{ record.participant_count || 0 }}
            </template>
            <template #timeRange="{ record }">
              {{ formatTime(record.start_time) }} ~ {{ formatTime(record.end_time) }}
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
                <a-button type="text" size="small" @click="editPromotion(record)">编辑</a-button>
                <a-button type="text" size="small" @click="viewParticipants(record)">参与者</a-button>
              </a-space>
            </template>
          </a-table>
        </a-card>
</div>

    <!-- 创建/编辑活动抽屉 -->
    <a-drawer
      v-model:visible="showCreateDrawer"
      :title="isEdit ? '编辑活动' : '创建活动'"
      :width="560"
      @before-ok="handleSave"
    >
      <a-form :model="form" layout="vertical" @submit-success="handleSaveSubmit">
        <a-form-item label="活动编码" required>
          <a-input v-model="form.promotion_code" placeholder="如 ACT001" :disabled="isEdit" />
        </a-form-item>
        <a-form-item label="活动名称" required>
          <a-input v-model="form.promotion_name" placeholder="如 周年庆双倍积分" />
        </a-form-item>
        <a-form-item label="活动类型" required>
          <a-select v-model="form.promotion_type" placeholder="选择活动类型">
            <a-option value="points_multi">双倍积分</a-option>
            <a-option value="discount">折扣活动</a-option>
            <a-option value="gift">赠品活动</a-option>
            <a-option value="coupon">优惠券发放</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="活动描述">
          <a-textarea v-model="form.description" :rows="3" placeholder="活动内容描述" />
        </a-form-item>
        <a-form-item v-if="form.promotion_type === 'points_multi'" label="积分倍数">
          <a-input-number v-model="form.points_multiplier" :min="1" :max="10" placeholder="如 2 表示双倍" style="width: 100%" />
        </a-form-item>
        <a-form-item v-if="form.promotion_type === 'discount'" label="折扣比例">
          <a-input-number v-model="form.discount_rate" :min="0.1" :max="1" :step="0.1" :precision="1" placeholder="如 0.8 表示 8 折" style="width: 100%" />
        </a-form-item>
        <a-form-item v-if="form.promotion_type === 'gift'" label="赠品名称">
          <a-input v-model="form.gift_name" placeholder="如 小风扇" />
        </a-form-item>
        <a-form-item v-if="form.promotion_type === 'coupon'" label="关联优惠券">
          <a-select v-model="form.coupon_id" placeholder="选择优惠券">
            <a-option v-for="c in couponList" :key="c.id" :value="c.id">{{ c.coupon_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="开始时间" required>
          <a-date-picker v-model="form.start_time" show-time style="width: 100%" />
        </a-form-item>
        <a-form-item label="结束时间" required>
          <a-date-picker v-model="form.end_time" show-time style="width: 100%" />
        </a-form-item>
        <a-form-item label="参与条件">
          <a-textarea v-model="form.conditions" :rows="2" placeholder="如：黄金会员及以上" />
        </a-form-item>
        <a-form-item label="状态">
          <a-radio-group v-model="form.status">
            <a-radio value="active">进行中</a-radio>
            <a-radio value="upcoming">即将开始</a-radio>
            <a-radio value="ended">已结束</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" html-type="submit">{{ isEdit ? '保存' : '创建' }}</a-button>
            <a-button @click="showCreateDrawer = false">取消</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-drawer>

    <!-- 活动详情抽屉 -->
    <a-drawer
      v-model:visible="showDetailDrawer"
      title="活动详情"
      :width="560"
    >
      <template v-if="currentPromotion">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="活动编码">{{ currentPromotion.promotion_code }}</a-descriptions-item>
          <a-descriptions-item label="活动名称">{{ currentPromotion.promotion_name }}</a-descriptions-item>
          <a-descriptions-item label="活动类型">
            <a-tag :color="getTypeColor(currentPromotion.promotion_type)">{{ getTypeText(currentPromotion.promotion_type) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="getStatusColor(currentPromotion.status)">{{ getStatusText(currentPromotion.status) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="参与人数">{{ currentPromotion.participant_count || 0 }}</a-descriptions-item>
          <a-descriptions-item label="活动描述">{{ currentPromotion.description || '-' }}</a-descriptions-item>
          <a-descriptions-item label="参与条件">{{ currentPromotion.conditions || '-' }}</a-descriptions-item>
          <a-descriptions-item label="活动时间">{{ formatTime(currentPromotion.start_time) }} ~ {{ formatTime(currentPromotion.end_time) }}</a-descriptions-item>
          <a-descriptions-item label="创建时间">{{ formatTime(currentPromotion.created_at) }}</a-descriptions-item>
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
const promotionList = ref([])
const couponList = ref([])
const showCreateDrawer = ref(false)
const showDetailDrawer = ref(false)
const isEdit = ref(false)
const currentPromotion = ref(null)

const filters = reactive({
  promotionType: undefined,
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
  running: 0,
  participants: 0,
  totalDiscount: 0
})

const form = reactive({
  promotion_code: '',
  promotion_name: '',
  promotion_type: 'points_multi',
  description: '',
  points_multiplier: 2,
  discount_rate: 0.8,
  gift_name: '',
  coupon_id: null,
  start_time: '',
  end_time: '',
  conditions: '',
  status: 'upcoming',
  created_by: 'admin'
})

const columns = [
  { title: '活动名称', dataIndex: 'promotion_name', width: 200 },
  { title: '编码', dataIndex: 'promotion_code', width: 100 },
  { title: '类型', slotName: 'promotionType', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '参与人数', slotName: 'participants', width: 100 },
  { title: '活动时间', slotName: 'timeRange', width: 280 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' }
]

if (routes[key]) router.push(routes[key])
  selectedKeys.value = [key]
}

const loadPromotions = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (filters.promotionType) params.promotion_type = filters.promotionType
    if (filters.status) params.status = filters.status
    if (filters.keyword) params.keyword = filters.keyword

    const res = await axios.get('/api/v1/members/promotions', { params })
    const data = res.data
    if (data.code === 0) {
      promotionList.value = data.data?.list || []
      pagination.total = data.data?.pagination?.total || 0
      updateStats()
    }
  } catch (err) {
    Message.error('加载促销活动列表失败')
  } finally {
    loading.value = false
  }
}

const updateStats = () => {
  stats.total = promotionList.value.length
  stats.running = promotionList.value.filter(p => p.status === 'active').length
  stats.participants = promotionList.value.reduce((sum, p) => sum + (p.participant_count || 0), 0)
  stats.totalDiscount = promotionList.value.reduce((sum, p) => sum + (p.total_discount || 0), 0)
}

const handlePageChange = (page) => {
  pagination.current = page
  loadPromotions()
}

const handlePageSizeChange = (pageSize) => {
  pagination.pageSize = pageSize
  pagination.current = 1
  loadPromotions()
}

const handleSave = (done) => { done(true) }

const handleSaveSubmit = async () => {
  try {
    const url = isEdit.value ? `/api/v1/members/promotions/${form.id}` : '/api/v1/members/promotions'
    const method = isEdit.value ? 'put' : 'post'
    const res = await axios[method](url, form)
    if (res.data.code === 0) {
      Message.success(isEdit.value ? '保存成功' : '创建成功')
      showCreateDrawer.value = false
      resetForm()
      loadPromotions()
    } else {
      Message.error(res.data.message || '操作失败')
    }
  } catch (err) {
    Message.error('操作失败')
  }
}

const resetForm = () => {
  form.promotion_code = ''
  form.promotion_name = ''
  form.promotion_type = 'points_multi'
  form.description = ''
  form.points_multiplier = 2
  form.discount_rate = 0.8
  form.gift_name = ''
  form.coupon_id = null
  form.start_time = ''
  form.end_time = ''
  form.conditions = ''
  form.status = 'upcoming'
  form.id = null
}

const openDetail = (record) => {
  currentPromotion.value = record
  showDetailDrawer.value = true
}

const editPromotion = (record) => {
  isEdit.value = true
  Object.assign(form, record)
  showCreateDrawer.value = true
}

const viewParticipants = (record) => {
  router.push({ path: '/members/promotions/participants', query: { promotionId: record.id } })
}

const getTypeColor = (type) => {
  const map = { points_multi: 'blue', discount: 'green', gift: 'purple', coupon: 'orange' }
  return map[type] || 'gray'
}

const getTypeText = (type) => {
  const map = { points_multi: '双倍积分', discount: '折扣活动', gift: '赠品活动', coupon: '优惠券发放' }
  return map[type] || type
}

const getStatusColor = (status) => {
  const map = { active: 'green', upcoming: 'blue', ended: 'gray' }
  return map[status] || 'gray'
}

const getStatusText = (status) => {
  const map = { active: '进行中', upcoming: '即将开始', ended: '已结束' }
  return map[status] || status
}

const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

onMounted(() => {
  loadPromotions()
})
</script>

<style scoped>
.member-promotions { min-height: 100vh; }
.header { background: #fff; padding: 0 16px; display: flex; align-items: center; gap: 16px; box-shadow: 0 1px 4px rgba(0,0,0,0.1); }
.header-left { display: flex; align-items: center; }
.header-title { font-size: 16px; font-weight: 500; }
.content { padding: 16px; background: #f0f2f5; }
.stats-row { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
