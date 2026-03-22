<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>订阅管理</a-breadcrumb-item>
      <a-breadcrumb-item>订阅计划</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 -->
    <div class="pro-page-header">
      <h2 class="pro-page-title">订阅计划</h2>
    </div>

    <!-- 操作栏 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button @click="loadPlans">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
      </a-space>
    </div>

    <!-- 订阅计划卡片 -->
    <div class="plans-grid" v-loading="loading">
      <div
        v-for="plan in plans"
        :key="plan.id"
        class="plan-card"
        :class="{ 'plan-card--recommended': plan.is_recommended, 'plan-card--current': isCurrentPlan(plan.plan_code) }"
      >
        <div v-if="plan.is_recommended" class="plan-badge">⭐ 推荐</div>
        <div class="plan-name">{{ plan.plan_name }}</div>
        <div class="plan-price">
          <span class="plan-price__amount">¥{{ plan.price_monthly }}</span>
          <span class="plan-price__unit">/月</span>
        </div>
        <div class="plan-divider"></div>
        <ul class="plan-features">
          <li v-for="(feature, idx) in plan.features" :key="idx">
            <icon-check class="feature-icon" /> {{ feature }}
          </li>
        </ul>
        <div class="plan-action">
          <a-button
            v-if="isCurrentPlan(plan.plan_code)"
            type="outline"
            disabled
          >当前版本</a-button>
          <a-button
            v-else-if="plan.plan_type === 'paid'"
            type="primary"
            @click="handleUpgrade(plan)"
          >立即升级</a-button>
          <a-button
            v-else
            type="outline"
            @click="handleSelectPlan(plan)"
          >选择此计划</a-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { billingApi, type SubscriptionPlan, type UserSubscription } from '@/api/billing'

const loading = ref(false)
const plans = ref<SubscriptionPlan[]>([])
const currentSubscription = ref<UserSubscription | null>(null)

const loadPlans = async () => {
  loading.value = true
  try {
    const [plansRes, subRes] = await Promise.all([
      billingApi.getSubscriptionPlans(),
      billingApi.getCurrentSubscription()
    ])
    if (plansRes.code === 0 || plansRes.code === 200) {
      plans.value = plansRes.data || []
    } else {
      // mock data
      plans.value = getMockPlans()
    }
    if (subRes.code === 0 || subRes.code === 200) {
      currentSubscription.value = subRes.data
    } else {
      currentSubscription.value = null
    }
  } catch (e) {
    plans.value = getMockPlans()
  } finally {
    loading.value = false
  }
}

const isCurrentPlan = (planCode: string) => {
  if (!currentSubscription.value) return planCode === 'free'
  return currentSubscription.value.plan_name.toLowerCase().includes(planCode)
}

const handleUpgrade = (plan: SubscriptionPlan) => {
  Message.info(`升级到 ${plan.plan_name} 功能开发中`)
}

const handleSelectPlan = (plan: SubscriptionPlan) => {
  Message.info(`选择 ${plan.plan_name} 功能开发中`)
}

const getMockPlans = (): SubscriptionPlan[] => [
  {
    id: 1,
    plan_code: 'free',
    plan_name: '免费版',
    plan_type: 'free',
    price_monthly: 0,
    price_yearly: 0,
    features: ['5台设备', '基础功能', '邮件支持'],
    quotas: { devices: 5, api_calls: 1000, storage_gb: 10 },
    is_active: true,
    is_recommended: false,
    sort_order: 1
  },
  {
    id: 2,
    plan_code: 'pro',
    plan_name: '专业版',
    plan_type: 'paid',
    price_monthly: 99,
    price_yearly: 990,
    features: ['50台设备', '高级功能', '优先支持', '数据分析'],
    quotas: { devices: 50, api_calls: 50000, storage_gb: 100 },
    is_active: true,
    is_recommended: true,
    sort_order: 2
  },
  {
    id: 3,
    plan_code: 'enterprise',
    plan_name: '企业版',
    plan_type: 'paid',
    price_monthly: 299,
    price_yearly: 2990,
    features: ['不限设备', '全部功能', '专属客服', 'SLA保障'],
    quotas: { devices: -1, api_calls: -1, storage_gb: -1 },
    is_active: true,
    is_recommended: false,
    sort_order: 3
  }
]

onMounted(() => {
  loadPlans()
})
</script>

<style scoped>
.plans-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  padding: 8px 0;
}

.plan-card {
  position: relative;
  background: var(--color-fill-2, #f2f3f5);
  border: 1px solid var(--color-border, #e5e6eb);
  border-radius: 8px;
  padding: 28px 24px;
  display: flex;
  flex-direction: column;
  transition: box-shadow 0.2s;
}

.plan-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.plan-card--recommended {
  border-color: rgb(var(--primary-6));
  background: linear-gradient(135deg, rgba(var(--primary-6), 0.04), rgba(var(--primary-6), 0.08));
}

.plan-card--current {
  border-color: rgb(var(--success-6));
}

.plan-badge {
  position: absolute;
  top: -1px;
  right: 20px;
  background: linear-gradient(135deg, rgb(var(--warning-6)), rgb(var(--warning-5)));
  color: #fff;
  font-size: 12px;
  padding: 2px 10px;
  border-radius: 0 0 8px 8px;
  font-weight: 600;
}

.plan-name {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-1, #1f2329);
  margin-bottom: 8px;
}

.plan-price {
  margin-bottom: 16px;
}

.plan-price__amount {
  font-size: 32px;
  font-weight: 700;
  color: var(--color-text-1, #1f2329);
}

.plan-price__unit {
  font-size: 14px;
  color: var(--color-text-3, #86909c);
}

.plan-divider {
  height: 1px;
  background: var(--color-border, #e5e6eb);
  margin-bottom: 16px;
}

.plan-features {
  list-style: none;
  padding: 0;
  margin: 0 0 24px 0;
  flex: 1;
}

.plan-features li {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 0;
  font-size: 14px;
  color: var(--color-text-2, #646a73);
}

.feature-icon {
  color: rgb(var(--success-6));
  flex-shrink: 0;
}

.plan-action {
  margin-top: auto;
}

.plan-action .arco-btn {
  width: 100%;
}
</style>
