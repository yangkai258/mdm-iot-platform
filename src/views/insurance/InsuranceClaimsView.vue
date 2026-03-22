<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>保险服务</a-breadcrumb-item>
      <a-breadcrumb-item>理赔申请</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="待处理" :value="stats.pending" :value-style="{ color: '#ff7d00' }">
            <template #prefix>⏳</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="处理中" :value="stats.processing">
            <template #prefix>🔄</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="已完成" :value="stats.completed" :value-style="{ color: '#0fc6c2' }">
            <template #prefix>✅</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="本月赔付金额" :value="stats.monthly_payout" prefix="¥" :precision="0" :value-style="{ color: '#f53f3f' }">
            <template #prefix>💸</template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- 筛选区 -->
    <a-card class="section-card" style="margin-top: 16px">
      <a-space wrap>
        <a-select v-model="filterForm.status" placeholder="理赔状态" style="width: 150px" allow-clear @change="loadClaims">
          <a-option value="pending">待处理</a-option>
          <a-option value="under_review">审核中</a-option>
          <a-option value="approved">已通过</a-option>
          <a-option value="rejected">已拒绝</a-option>
          <a-option value="paid">已打款</a-option>
        </a-select>
        <a-select v-model="filterForm.category" placeholder="产品类别" style="width: 150px" allow-clear @change="loadClaims">
          <a-option value="device">设备保障</a-option>
          <a-option value="health">健康医疗</a-option>
          <a-option value="accident">意外险</a-option>
          <a-option value="property">财产险</a-option>
        </a-select>
        <a-input-search v-model="filterForm.keyword" placeholder="搜索单号/投保人/设备ID" style="width: 240px" search-button @search="loadClaims" @change="e => !e.target.value && loadClaims()" />
        <a-range-picker v-model="filterForm.dateRange" style="width: 260px" @change="loadClaims" />
        <a-button @click="loadClaims">
          <template #icon><icon-refresh :spin="loading" /></template>
          刷新
        </a-button>
        <a-button type="primary" @click="showCreateModal">
          <template #icon><icon-plus /></template>
          新增理赔
        </a-button>
      </a-space>
    </a-card>

    <!-- 理赔列表 -->
    <a-card class="section-card" style="margin-top: 16px">
      <a-spin :loading="loading" tip="加载中...">
        <a-table :data="claims" :pagination="{ pageSize: 10, total: claims.length, showTotal: true }" :columns="columns" row-key="claim_id" stripe>
          <template #status="{ record }">
            <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
          </template>
          <template #amount="{ record }">
            ¥{{ record.amount?.toFixed(2) }}
          </template>
          <template #created_at="{ record }">
            {{ formatDate(record.created_at) }}
          </template>
          <template #operations="{ record }">
            <a-space>
              <a-button type="text" size="small" @click="viewDetail(record)">详情</a-button>
              <a-button v-if="record.status === 'pending'" type="text" size="small" @click="showReviewModal(record)">审核</a-button>
              <a-button v-if="record.status === 'approved'" type="text" size="small" status="success" @click="handlePay(record)">打款</a-button>
              <a-button v-if="['pending', 'under_review'].includes(record.status)" type="text" size="small" status="danger" @click="handleReject(record)">拒绝</a-button>
            </a-space>
          </template>
        </a-table>
      </a-spin>
    </a-card>

    <!-- 新增理赔弹窗 -->
    <a-modal v-model:visible="createVisible" title="新增理赔申请" :width="600" @before-ok="handleCreate" @cancel="createVisible = false">
      <a-form :model="claimForm" layout="vertical">
        <a-form-item label="产品" field="product_id">
          <a-select v-model="claimForm.product_id" placeholder="请选择投保产品">
            <a-option v-for="p in products" :key="p.product_id" :value="p.product_id">{{ p.name }}（剩余额度：¥{{ p.remaining_amount }}）</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="理赔类型" field="claim_type">
          <a-select v-model="claimForm.claim_type" placeholder="请选择理赔类型">
            <a-option value="accident">意外损坏</a-option>
            <a-option value="malfunction">设备故障</a-option>
            <a-option value="theft">被盗</a-option>
            <a-option value="medical">医疗费用</a-option>
            <a-option value="other">其他</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="理赔金额（元）" field="amount">
          <a-input-number v-model="claimForm.amount" :min="0" :precision="2" style="width: 100%" />
        </a-form-item>
        <a-form-item label="设备ID" field="device_id">
          <a-input v-model="claimForm.device_id" placeholder="请输入关联设备ID（可选）" />
        </a-form-item>
        <a-form-item label="事故描述" field="description">
          <a-textarea v-model="claimForm.description" placeholder="请详细描述事故经过" :rows="4" />
        </a-form-item>
        <a-form-item label="上传凭证" field="attachments">
          <a-upload action="/" :limit="5" accept="image/*,.pdf" />
          <div style="font-size: 12px; color: #86909c; margin-top: 4px">支持 JPG、PNG、PDF，单个文件不超过 10MB</div>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 理赔详情 -->
    <a-drawer v-model:visible="detailVisible" :title="`理赔详情: ${currentClaim?.claim_id}`" :width="640" unmountOnHide>
      <a-descriptions :column="2" bordered size="large" style="margin-bottom: 16px">
        <a-descriptions-item label="理赔单号">{{ currentClaim?.claim_id }}</a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="getStatusColor(currentClaim?.status)">{{ getStatusText(currentClaim?.status) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="产品名称">{{ currentClaim?.product_name }}</a-descriptions-item>
        <a-descriptions-item label="理赔类型">{{ getClaimTypeText(currentClaim?.claim_type) }}</a-descriptions-item>
        <a-descriptions-item label="投保人">{{ currentClaim?.policyholder }}</a-descriptions-item>
        <a-descriptions-item label="联系方式">{{ currentClaim?.contact }}</a-descriptions-item>
        <a-descriptions-item label="理赔金额">¥{{ currentClaim?.amount?.toFixed(2) }}</a-descriptions-item>
        <a-descriptions-item label="申请时间">{{ formatDate(currentClaim?.created_at) }}</a-descriptions-item>
        <a-descriptions-item label="事故描述" :span="2">{{ currentClaim?.description }}</a-descriptions-item>
      </a-descriptions>

      <a-card title="处理记录" style="margin-bottom: 16px">
        <a-timeline>
          <a-timeline-item v-for="(log, idx) in currentClaim?.logs || []" :key="idx" :color="log.color || 'gray'">
            <div class="timeline-title">{{ log.title }}</div>
            <div class="timeline-time">{{ log.time }}</div>
            <div class="timeline-desc">{{ log.desc }}</div>
          </a-timeline-item>
        </a-timeline>
      </a-card>

      <a-card title="审核意见" v-if="currentClaim?.review_note">
        <a-alert :type="currentClaim?.status === 'rejected' ? 'error' : 'success'">{{ currentClaim?.review_note }}</a-alert>
      </a-card>
    </a-drawer>

    <!-- 审核弹窗 -->
    <a-modal v-model:visible="reviewVisible" title="理赔审核" @before-ok="handleReview" @cancel="reviewVisible = false">
      <a-form :model="reviewForm" layout="vertical">
        <a-form-item label="审核结论" field="result">
          <a-radio-group v-model="reviewForm.result">
            <a-radio value="approved">通过</a-radio>
            <a-radio value="rejected">拒绝</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item v-if="reviewForm.result === 'approved'" label="赔付金额（元）" field="amount">
          <a-input-number v-model="reviewForm.amount" :min="0" :precision="2" style="width: 100%" />
        </a-form-item>
        <a-form-item label="审核意见" field="note">
          <a-textarea v-model="reviewForm.note" placeholder="请输入审核意见" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const createVisible = ref(false)
const detailVisible = ref(false)
const reviewVisible = ref(false)
const currentClaim = ref(null)

const stats = reactive({ pending: 0, processing: 0, completed: 0, monthly_payout: 0 })

const filterForm = reactive({ status: '', category: '', keyword: '', dateRange: [] })

const claimForm = reactive({ product_id: '', claim_type: '', amount: 0, device_id: '', description: '', attachments: [] })
const reviewForm = reactive({ result: 'approved', amount: 0, note: '' })

const products = ref([
  { product_id: 'INS-DEV-001', name: '设备意外险', remaining_amount: 2800 },
  { product_id: 'INS-HEALTH-001', name: '健康医疗险', remaining_amount: 85000 },
  { product_id: 'INS-ACC-001', name: '综合意外险', remaining_amount: 450000 }
])

const claims = ref([])

const columns = [
  { title: '理赔单号', dataIndex: 'claim_id', width: 150 },
  { title: '产品', dataIndex: 'product_name', width: 120 },
  { title: '投保人', dataIndex: 'policyholder', width: 100 },
  { title: '理赔类型', dataIndex: 'claim_type', width: 100 },
  { title: '金额', dataIndex: 'amount', width: 110, slotName: 'amount' },
  { title: '状态', dataIndex: 'status', width: 90, slotName: 'status' },
  { title: '申请时间', dataIndex: 'created_at', width: 120, slotName: 'created_at' },
  { title: '操作', slotName: 'operations', width: 180 }
]

const getStatusColor = (s) => ({ pending: 'orange', under_review: 'blue', approved: 'green', rejected: 'red', paid: 'cyan' }[s] || 'gray')
const getStatusText = (s) => ({ pending: '待处理', under_review: '审核中', approved: '已通过', rejected: '已拒绝', paid: '已打款' }[s] || s)
const getClaimTypeText = (t) => ({ accident: '意外损坏', malfunction: '设备故障', theft: '被盗', medical: '医疗费用', other: '其他' }[t] || t)

const formatDate = (d) => d ? new Date(d).toLocaleDateString('zh-CN') : '-'

const loadClaims = async () => {
  loading.value = true
  await new Promise(r => setTimeout(r, 600))
  claims.value = [
    { claim_id: 'CLM-202603001', product_id: 'INS-DEV-001', product_name: '设备意外险', policyholder: '张三', contact: '138****1234', claim_type: 'accident', amount: 599.00, device_id: 'DEV-1001', description: '手机意外跌落导致屏幕碎裂，已拍照留证', status: 'pending', created_at: '2026-03-22 10:30', logs: [{ title: '提交申请', time: '2026-03-22 10:30', desc: '用户提交理赔申请' }] },
    { claim_id: 'CLM-202603002', product_id: 'INS-HEALTH-001', product_name: '健康医疗险', policyholder: '李四', contact: '139****5678', claim_type: 'medical', amount: 3200.00, device_id: '', description: '住院治疗费用理赔，共住院5天', status: 'under_review', created_at: '2026-03-21 14:20', logs: [{ title: '提交申请', time: '2026-03-21 14:20', desc: '用户提交理赔申请' }, { title: '资料审核', time: '2026-03-21 16:00', desc: '理赔专员开始审核材料', color: 'blue' }] },
    { claim_id: 'CLM-202603003', product_id: 'INS-ACC-001', product_name: '综合意外险', policyholder: '王五', contact: '137****9012', claim_type: 'accident', amount: 5000.00, device_id: '', description: '交通事故受伤，申请伤残理赔', status: 'approved', created_at: '2026-03-20 09:00', review_note: '事故属实，伤残鉴定为10级，按合同赔付', logs: [{ title: '提交申请', time: '2026-03-20 09:00', desc: '用户提交理赔申请' }, { title: '资料审核', time: '2026-03-20 11:00', desc: '理赔专员审核通过' }, { title: '赔付审核', time: '2026-03-21 10:00', desc: '赔付金额确认: ¥5000', color: 'green' }] },
    { claim_id: 'CLM-202603004', product_id: 'INS-DEV-002', product_name: '设备延保险', policyholder: '赵六', contact: '136****3456', claim_type: 'malfunction', amount: 899.00, device_id: 'DEV-2003', description: '设备电池老化故障', status: 'paid', created_at: '2026-03-18 15:30', review_note: '审核通过，已完成打款', logs: [{ title: '提交申请', time: '2026-03-18 15:30', desc: '用户提交理赔申请' }, { title: '资料审核', time: '2026-03-18 17:00', desc: '审核通过' }, { title: '打款完成', time: '2026-03-19 10:00', desc: '赔付金额 ¥899 已到账', color: 'green' }] },
    { claim_id: 'CLM-202603005', product_id: 'INS-DEV-001', product_name: '设备意外险', policyholder: '钱七', contact: '135****7890', claim_type: 'accident', amount: 1200.00, device_id: 'DEV-1005', description: '设备浸水损坏', status: 'rejected', created_at: '2026-03-17 11:00', review_note: '设备浸水属于免责条款，不予赔付', logs: [{ title: '提交申请', time: '2026-03-17 11:00', desc: '用户提交理赔申请' }, { title: '资料审核', time: '2026-03-17 14:00', desc: '审核不通过', color: 'red' }, { title: '理赔拒绝', time: '2026-03-17 16:00', desc: '因属于免责条款拒绝赔付', color: 'red' }] }
  ]
  stats.pending = claims.value.filter(c => c.status === 'pending').length
  stats.processing = claims.value.filter(c => ['under_review', 'approved'].includes(c.status)).length
  stats.completed = claims.value.filter(c => c.status === 'paid').length
  stats.monthly_payout = claims.value.filter(c => c.status === 'paid').reduce((s, c) => s + c.amount, 0)
  loading.value = false
}

const showCreateModal = () => {
  Object.assign(claimForm, { product_id: '', claim_type: '', amount: 0, device_id: '', description: '' })
  createVisible.value = true
}

const handleCreate = async (done) => {
  if (!claimForm.product_id || !claimForm.claim_type || !claimForm.amount) {
    Message.error('请填写必填字段')
    done(false)
    return
  }
  await new Promise(r => setTimeout(r, 800))
  Message.success('理赔申请已提交')
  createVisible.value = false
  loadClaims()
  done(true)
}

const viewDetail = (record) => {
  currentClaim.value = record
  detailVisible.value = true
}

const showReviewModal = (record) => {
  currentClaim.value = record
  Object.assign(reviewForm, { result: 'approved', amount: record.amount, note: '' })
  reviewVisible.value = true
}

const handleReview = async (done) => {
  if (!reviewForm.note) {
    Message.error('请填写审核意见')
    done(false)
    return
  }
  await new Promise(r => setTimeout(r, 500))
  Message.success(reviewForm.result === 'approved' ? '审核通过' : '已拒绝')
  reviewVisible.value = false
  loadClaims()
  done(true)
}

const handlePay = async (record) => {
  await new Promise(r => setTimeout(r, 500))
  Message.success('打款成功')
  loadClaims()
}

const handleReject = async (record) => {
  await new Promise(r => setTimeout(r, 300))
  Message.warning('理赔已拒绝')
  loadClaims()
}

loadClaims()
</script>

<style scoped>
.page-container { padding: 20px; }
.breadcrumb { margin-bottom: 12px; }
.stat-card { text-align: center; }
.section-card { margin-bottom: 16px; }
.timeline-title { font-weight: 600; color: #1d2129; }
.timeline-time { font-size: 12px; color: #86909c; }
.timeline-desc { font-size: 13px; color: #4e5969; margin-top: 4px; }
</style>
