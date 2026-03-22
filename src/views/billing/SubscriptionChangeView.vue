<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>订阅管理</a-breadcrumb-item>
      <a-breadcrumb-item>变更方案</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 -->
    <div class="pro-page-header">
      <h2 class="pro-page-title">订阅变更</h2>
    </div>

    <div class="pro-action-bar">
      <a-space>
        <a-button @click="loadData">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
      </a-space>
    </div>

    <a-row :gutter="24" v-loading="loading">
      <!-- 当前方案 -->
      <a-col :span="12">
        <a-card class="plan-compare-card">
          <template #title>当前方案</template>
          <div class="plan-compare-name">{{ currentPlan?.plan_name || '免费版' }}</div>
          <div class="plan-compare-price">¥{{ currentPlan?.price_monthly || 0 }}/月</div>
          <a-divider />
          <div class="plan-compare-features">
            <div v-for="(f, idx) in (currentPlan?.features || ['5台设备', '基础功能', '邮件支持'])" :key="idx" class="feature-item">
              <icon-check class="feature-icon" /> {{ f }}
            </div>
          </div>
          <div class="plan-compare-effective">已于 {{ formatDate(currentSubscription?.started_at) }} 生效</div>
        </a-card>
      </a-col>

      <!-- 新方案 -->
      <a-col :span="12">
        <a-card class="plan-compare-card plan-compare-card--new">
          <template #title>新方案</template>
          <div class="plan-select">
            <a-select v-model="newPlanId" placeholder="请选择方案" style="width: 100%;" @change="onPlanChange">
              <a-option v-for="p in allPlans" :key="p.id" :value="p.id" :label="`${p.plan_name} - ¥${p.price_monthly}/月`" />
            </a-select>
          </div>
          <div class="plan-compare-name" v-if="newPlan">{{ newPlan.plan_name }}</div>
          <div class="plan-compare-price" v-if="newPlan">¥{{ newPlan.price_monthly }}/月</div>
          <a-divider />
          <div class="plan-compare-features" v-if="newPlan">
            <div v-for="(f, idx) in newPlan.features" :key="idx" class="feature-item">
              <icon-check class="feature-icon" /> {{ f }}
            </div>
          </div>

          <!-- 价格对比 -->
          <a-divider v-if="priceDiff !== 0" />
          <div class="price-diff" v-if="priceDiff !== 0">
            <div class="price-diff-label">{{ priceDiff > 0 ? '每月需补差价' : '每月节省' }}</div>
            <div class="price-diff-value" :class="{ 'price-diff-value--save': priceDiff < 0 }">
              ¥{{ Math.abs(priceDiff) }}/月
            </div>
          </div>

          <div class="plan-compare-effective" v-if="newPlan">
            变更将在 {{ effectiveText }} 生效
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 确认操作 -->
    <div class="change-actions" v-if="newPlanId">
      <a-space>
        <a-button type="primary" size="large" @click="confirmChange">
          确认变更
        </a-button>
        <a-button size="large" @click="newPlanId = null; newPlan = null">
          取消
        </a-button>
      </a-space>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { billingApi, type SubscriptionPlan, type UserSubscription } from '@/api/billing'

const loading = ref(false)
const currentSubscription = ref<UserSubscription | null>(null)
const currentPlan = ref<SubscriptionPlan | null>(null)
const allPlans = ref<SubscriptionPlan[]>([])
const newPlanId = ref<number | null>(null)
const newPlan = ref<SubscriptionPlan | null>(null)

const priceDiff = computed(() => {
  if (!currentPlan.value || !newPlan.value) return 0
  return newPlan.value.price_monthly - currentPlan.value.price_monthly
})

const effectiveText = computed(() => {
  if (priceDiff.value > 0) return '立即'
  return '当前周期末'
})

const formatDate = (dateStr?: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('zh-CN')
}

const onPlanChange = (id: number) => {
  newPlan.value = allPlans.value.find(p => p.id === id) || null
}

const loadData = async () => {
  loading.value = true
  try {
    const [plansRes, subRes] = await Promise.all([
      billingApi.getSubscriptionPlans(),
      billingApi.getCurrentSubscription()
    ])
    allPlans.value = plansRes.code === 0 || plansRes.code === 200
      ? (plansRes.data || [])
      : getMockPlans()
    currentSubscription.value = subRes.code === 0 || subRes.code === 200
      ? subRes.data
      : getMockSubscription()
    currentPlan.value = allPlans.value.find(p => p.id === currentSubscription.value?.plan_id) || allPlans.value[0]
  } catch (e) {
    allPlans.value = getMockPlans()
    currentSubscription.value = getMockSubscription()
    currentPlan.value = allPlans.value[0]
  } finally {
    loading.value = false
  }
}

const getMockPlans = (): SubscriptionPlan[] => [
  { id: 1, plan_code: 'free', plan_name: '免费版', plan_type: 'free', price_monthly: 0, price_yearly: 0, features: ['5台设备', '基础功能', '邮件支持'], quotas: { devices: 5, api_calls: 1000, storage_gb: 10 }, is_active: true, is_recommended: false, sort_order: 1 },
  { id: 2, plan_code: 'pro', plan_name: '专业版', plan_type: 'paid', price_monthly: 99, price_yearly: 990, features: ['50台设备', '高级功能', '优先支持', '数据分析'], quotas: { devices: 50, api_calls: 50000, storage_gb: 100 }, is_active: true, is_recommended: true, sort_order: 2 },
  { id: 3, plan_code: 'enterprise', plan_name: '企业版', plan_type: 'paid', price_monthly: 299, price_yearly: 2990, features: ['不限设备', '全部功能', '专属客服', 'SLA保障'], quotas: { devices: -1, api_calls: -1, storage_gb: -1 }, is_active: true, is_recommended: false, sort_order: 3 }
]

const getMockSubscription = (): UserSubscription => ({
  id: 1, user_id: 1, plan_id: 2, plan_name: '专业版', status: 'active',
  started_at: '2026-02-01T00:00:00Z', expires_at: '2026-04-01T00:00:00Z',
  next_billing_at: '2026-04-01T00:00:00Z', auto_renew: true
})

const confirmChange = async () => {
  if (!newPlan.value || !currentSubscription.value) return
  const isUpgrade = priceDiff.value > 0
  try {
    let res: any
    if (isUpgrade) {
      res = await billingApi.upgradeSubscription(currentSubscription.value.id, { new_plan_id: newPlan.value.id })
    } else {
      res = await billingApi.downgradeSubscription(currentSubscription.value.id, { new_plan_id: newPlan.value.id })
    }
    if (res.code === 0 || res.code === 200) {
      Message.success('变更成功')
      loadData()
    } else {
      Message.error('变更失败')
    }
  } catch (e) {
    Message.error('变更失败，请稍后重试')
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.plan-compare-card {
  height: 100%;
}
.plan-compare-card--new {
  border-color: rgb(var(--primary-6));
}
.plan-compare-name {
  font-size: 22px;
  font-weight: 700;
  margin-bottom: 8px;
}
.plan-compare-price {
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text-1);
}
.plan-compare-features {
  margin: 16px 0;
}
.feature-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 0;
  font-size: 14px;
  color: var(--color-text-2, #646a73);
}
.feature-icon {
  color: rgb(var(--success-6));
}
.plan-compare-effective {
  font-size: 12px;
  color: var(--color-text-3, #86909c);
  margin-top: 8px;
}
.price-diff {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 0;
}
.price-diff-label {
  font-size: 14px;
  color: var(--color-text-3, #86909c);
}
.price-diff-value {
  font-size: 20px;
  font-weight: 700;
  color: rgb(var(--danger-6));
}
.price-diff-value--save {
  color: rgb(var(--success-6));
}
.change-actions {
  margin-top: 24px;
  text-align: center;
}
</style>
