<template>
  <div class="member-coupons">
    <!-- 统计卡片 -->
        <a-row :gutter="16" class="stats-row">
          <a-col :span="6">
            <a-card>
              <a-statistic title="优惠券总数" :value="stats.total" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="已发放" :value="stats.issued" :value-style="{ color: '#1890ff' }" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="已使用" :value="stats.used" :value-style="{ color: '#52c41a' }" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="使用率" :value="stats.usageRate" :suffix="'%'" :value-style="{ color: '#faad14' }" />
            </a-card>
          </a-col>
        </a-row>

        <!-- 操作栏 -->
        <a-card class="action-card">
          <a-space wrap>
            <a-select v-model="filters.couponType" placeholder="优惠券类型" allow-clear style="width: 140px" @change="loadCoupons">
              <a-option value="discount">折扣券</a-option>
              <a-option value="cash">代金券</a-option>
              <a-option value="gift">礼品券</a-option>
            </a-select>
            <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 120px" @change="loadCoupons">
              <a-option value="active">生效中</a-option>
              <a-option value="paused">已暂停</a-option>
              <a-option value="expired">已过期</a-option>
            </a-select>
            <a-input-search v-model="filters.keyword" placeholder="搜索优惠券名称/编码" style="width: 200px" search-button @search="loadCoupons" />
            <a-button type="primary" @click="showCreateDrawer = true">创建优惠券</a-button>
            <a-button @click="loadCoupons">刷新</a-button>
          </a-space>
        </a-card>

        <!-- 优惠券列表 -->
        <a-card class="coupon-card">
          <a-table
            :columns="columns"
            :data="couponList"
            :loading="loading"
            :pagination="pagination"
            row-key="id"
            @page-change="handlePageChange"
            @page-size-change="handlePageSizeChange"
          >
            <template #couponType="{ record }">
              <a-tag :color="getTypeColor(record.coupon_type)">{{ getTypeText(record.coupon_type) }}</a-tag>
            </template>
            <template #discountValue="{ record }">
              <span v-if="record.coupon_type === 'discount'">{{ record.discount_value }}折</span>
              <span v-else-if="record.coupon_type === 'cash'">¥{{ record.cash_value }}</span>
              <span v-else>{{ record.gift_name }}</span>
            </template>
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
            </template>
            <template #usageRate="{ record }">
              {{ record.usage_rate || 0 }}%
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
                <a-button type="text" size="small" @click="distributeCoupon(record)">发放</a-button>
                <a-button v-if="record.status === 'active'" type="text" size="small" status="warning" @click="pauseCoupon(record)">暂停</a-button>
                <a-button v-else-if="record.status === 'paused'" type="text" size="small" @click="resumeCoupon(record)">恢复</a-button>
              </a-space>
            </template>
          </a-table>
        </a-card>
</div>

    <!-- 创建优惠券抽屉 -->
    <a-drawer
      v-model:visible="showCreateDrawer"
      title="创建优惠券"
      :width="480"
      @before-ok="handleSave"
    >
      <a-form :model="form" layout="vertical" @submit-success="handleSaveSubmit">
        <a-form-item label="优惠券编码" required>
          <a-input v-model="form.coupon_code" placeholder="如 CPL001" />
        </a-form-item>
        <a-form-item label="优惠券名称" required>
          <a-input v-model="form.coupon_name" placeholder="如 新用户专享券" />
        </a-form-item>
        <a-form-item label="优惠券类型" required>
          <a-select v-model="form.coupon_type" placeholder="选择类型">
            <a-option value="discount">折扣券</a-option>
            <a-option value="cash">代金券</a-option>
            <a-option value="gift">礼品券</a-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="form.coupon_type === 'discount'" label="折扣值" required>
          <a-input-number v-model="form.discount_value" :min="1" :max="9.9" :step="0.1" :precision="1" placeholder="如 8.5 表示 8.5 折" style="width: 100%" />
        </a-form-item>
        <a-form-item v-if="form.coupon_type === 'cash'" label="代金金额" required>
          <a-input-number v-model="form.cash_value" :min="1" :max="10000" placeholder="如 100" style="width: 100%" />
          <template #extra>单位：元</template>
        </a-form-item>
        <a-form-item v-if="form.coupon_type === 'gift'" label="礼品名称" required>
          <a-input v-model="form.gift_name" placeholder="如 小风扇" />
        </a-form-item>
        <a-form-item label="使用条件">
          <a-textarea v-model="form.usage_condition" :rows="2" placeholder="满减条件，如：满100元可用" />
        </a-form-item>
        <a-form-item label="有效期开始" required>
          <a-date-picker v-model="form.start_time" style="width: 100%" />
        </a-form-item>
        <a-form-item label="有效期结束" required>
          <a-date-picker v-model="form.end_time" style="width: 100%" />
        </a-form-item>
        <a-form-item label="发行数量">
          <a-input-number v-model="form.total_quantity" :min="1" placeholder="不填表示不限量" style="width: 100%" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" :rows="3" placeholder="优惠券描述" />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" html-type="submit">创建</a-button>
            <a-button @click="showCreateDrawer = false">取消</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-drawer>

    <!-- 优惠券详情抽屉 -->
    <a-drawer
      v-model:visible="showDetailDrawer"
      title="优惠券详情"
      :width="480"
    >
      <template v-if="currentCoupon">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="优惠券编码">{{ currentCoupon.coupon_code }}</a-descriptions-item>
          <a-descriptions-item label="优惠券名称">{{ currentCoupon.coupon_name }}</a-descriptions-item>
          <a-descriptions-item label="类型">
            <a-tag :color="getTypeColor(currentCoupon.coupon_type)">{{ getTypeText(currentCoupon.coupon_type) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="面值">
            <span v-if="currentCoupon.coupon_type === 'discount'">{{ currentCoupon.discount_value }}折</span>
            <span v-else-if="currentCoupon.coupon_type === 'cash'">¥{{ currentCoupon.cash_value }}</span>
            <span v-else>{{ currentCoupon.gift_name }}</span>
          </a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="getStatusColor(currentCoupon.status)">{{ getStatusText(currentCoupon.status) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="使用条件">{{ currentCoupon.usage_condition || '-' }}</a-descriptions-item>
          <a-descriptions-item label="发行数量">{{ currentCoupon.total_quantity || '不限量' }}</a-descriptions-item>
          <a-descriptions-item label="已发放">{{ currentCoupon.issued_count || 0 }}</a-descriptions-item>
          <a-descriptions-item label="已使用">{{ currentCoupon.used_count || 0 }}</a-descriptions-item>
          <a-descriptions-item label="使用率">{{ currentCoupon.usage_rate || 0 }}%</a-descriptions-item>
          <a-descriptions-item label="有效期">{{ formatTime(currentCoupon.start_time) }} ~ {{ formatTime(currentCoupon.end_time) }}</a-descriptions-item>
          <a-descriptions-item label="描述">{{ currentCoupon.description || '-' }}</a-descriptions-item>
        </a-descriptions>
      </template>
    </a-drawer>

    <!-- 发放优惠券抽屉 -->
    <a-drawer
      v-model:visible="showDistributeDrawer"
      title="发放优惠券"
      :width="400"
    >
      <a-form :model="distributeForm" layout="vertical">
        <a-form-item label="发放方式">
          <a-radio-group v-model="distributeForm.distribute_type">
            <a-radio value="member">指定会员</a-radio>
            <a-radio value="all">全部会员</a-radio>
            <a-radio value="level">按等级发放</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item v-if="distributeForm.distribute_type === 'member'" label="选择会员" required>
          <a-select v-model="distributeForm.member_ids" placeholder="选择会员" multiple searchable>
            <a-option v-for="m in memberList" :key="m.id" :value="m.id">{{ m.member_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="distributeForm.distribute_type === 'level'" label="选择等级" required>
          <a-checkbox-group v-model="distributeForm.levels">
            <a-checkbox value="gold">黄金会员</a-checkbox>
            <a-checkbox value="silver">白银会员</a-checkbox>
            <a-checkbox value="bronze">青铜会员</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
        <a-form-item label="发放数量">
          <a-input-number v-model="distributeForm.quantity" :min="1" :max="10" placeholder="每人发放数量" style="width: 100%" />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleDistribute">确认发放</a-button>
            <a-button @click="showDistributeDrawer = false">取消</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import axios from 'axios'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const couponList = ref([])
const memberList = ref([])
const showCreateDrawer = ref(false)
const showDetailDrawer = ref(false)
const showDistributeDrawer = ref(false)
const currentCoupon = ref(null)

const filters = reactive({
  couponType: undefined,
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
  issued: 0,
  used: 0,
  usageRate: 0
})

const form = reactive({
  coupon_code: '',
  coupon_name: '',
  coupon_type: 'discount',
  discount_value: 0,
  cash_value: 0,
  gift_name: '',
  usage_condition: '',
  start_time: '',
  end_time: '',
  total_quantity: null,
  description: '',
  created_by: 'admin'
})

const distributeForm = reactive({
  distribute_type: 'member',
  member_ids: [],
  levels: [],
  quantity: 1,
  coupon_id: null
})

const columns = [
  { title: '优惠券名称', dataIndex: 'coupon_name', width: 180 },
  { title: '编码', dataIndex: 'coupon_code', width: 100 },
  { title: '类型', slotName: 'couponType', width: 100 },
  { title: '面值', slotName: 'discountValue', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '发行/已用', dataIndex: 'issued_count', width: 100 },
  { title: '使用率', slotName: 'usageRate', width: 80 },
  { title: '有效期', dataIndex: 'end_time', width: 160 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' }
]

if (routes[key]) router.push(routes[key])
  selectedKeys.value = [key]
}

const loadCoupons = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (filters.couponType) params.coupon_type = filters.couponType
    if (filters.status) params.status = filters.status
    if (filters.keyword) params.keyword = filters.keyword

    const res = await axios.get('/api/v1/members/coupons', { params })
    const data = res.data
    if (data.code === 0) {
      couponList.value = data.data?.list || []
      pagination.total = data.data?.pagination?.total || 0
      updateStats()
    }
  } catch (err) {
    Message.error('加载优惠券列表失败')
  } finally {
    loading.value = false
  }
}

const updateStats = () => {
  stats.total = couponList.value.length
  stats.issued = couponList.value.reduce((sum, c) => sum + (c.issued_count || 0), 0)
  stats.used = couponList.value.reduce((sum, c) => sum + (c.used_count || 0), 0)
  stats.usageRate = stats.issued > 0 ? Math.round((stats.used / stats.issued) * 100) : 0
}

const handlePageChange = (page) => {
  pagination.current = page
  loadCoupons()
}

const handlePageSizeChange = (pageSize) => {
  pagination.pageSize = pageSize
  pagination.current = 1
  loadCoupons()
}

const handleSave = (done) => { done(true) }

const handleSaveSubmit = async () => {
  try {
    const res = await axios.post('/api/v1/members/coupons', form)
    if (res.data.code === 0) {
      Message.success('创建成功')
      showCreateDrawer.value = false
      resetForm()
      loadCoupons()
    } else {
      Message.error(res.data.message || '创建失败')
    }
  } catch (err) {
    Message.error('创建失败')
  }
}

const resetForm = () => {
  form.coupon_code = ''
  form.coupon_name = ''
  form.coupon_type = 'discount'
  form.discount_value = 0
  form.cash_value = 0
  form.gift_name = ''
  form.usage_condition = ''
  form.start_time = ''
  form.end_time = ''
  form.total_quantity = null
  form.description = ''
}

const openDetail = (record) => {
  currentCoupon.value = record
  showDetailDrawer.value = true
}

const distributeCoupon = async (record) => {
  distributeForm.coupon_id = record.id
  showDistributeDrawer.value = true
}

const handleDistribute = async () => {
  try {
    const res = await axios.post('/api/v1/members/coupons/distribute', distributeForm)
    if (res.data.code === 0) {
      Message.success(`成功发放 ${res.data.data?.count || 0} 张优惠券`)
      showDistributeDrawer.value = false
      loadCoupons()
    } else {
      Message.error(res.data.message || '发放失败')
    }
  } catch (err) {
    Message.error('发放失败')
  }
}

const pauseCoupon = async (record) => {
  try {
    const res = await axios.post(`/api/v1/members/coupons/${record.id}/pause`)
    if (res.data.code === 0) {
      Message.success('已暂停')
      loadCoupons()
    } else {
      Message.error(res.data.message || '操作失败')
    }
  } catch (err) {
    Message.error('操作失败')
  }
}

const resumeCoupon = async (record) => {
  try {
    const res = await axios.post(`/api/v1/members/coupons/${record.id}/resume`)
    if (res.data.code === 0) {
      Message.success('已恢复')
      loadCoupons()
    } else {
      Message.error(res.data.message || '操作失败')
    }
  } catch (err) {
    Message.error('操作失败')
  }
}

const getTypeColor = (type) => {
  const map = { discount: 'blue', cash: 'green', gift: 'purple' }
  return map[type] || 'gray'
}

const getTypeText = (type) => {
  const map = { discount: '折扣券', cash: '代金券', gift: '礼品券' }
  return map[type] || type
}

const getStatusColor = (status) => {
  const map = { active: 'green', paused: 'yellow', expired: 'gray' }
  return map[status] || 'gray'
}

const getStatusText = (status) => {
  const map = { active: '生效中', paused: '已暂停', expired: '已过期' }
  return map[status] || status
}

const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

onMounted(() => {
  loadCoupons()
})
</script>

<style scoped>
.member-coupons {
  padding: 16px;
}
.stats-row {
  margin-bottom: 16px;
}
.action-card {
  margin-bottom: 16px;
}
</style>
