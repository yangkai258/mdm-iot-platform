<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>订阅管理</a-breadcrumb-item>
      <a-breadcrumb-item>我的订阅</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 -->
    <div class="pro-page-header">
      <h2 class="pro-page-title">我的订阅</h2>
    </div>

    <!-- 操作栏 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button @click="loadSubscription">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
        <a-button v-if="subscription?.status === 'active'" @click="handleRenew">
          <template #icon><icon-refresh /></template>
          续费
        </a-button>
        <a-button v-if="subscription?.status === 'active'" @click="handleChangePlan">
          变更方案
        </a-button>
        <a-button v-if="subscription?.status === 'active'" status="danger" @click="handleCancel">
          取消订阅
        </a-button>
      </a-space>
    </div>

    <!-- 当前订阅状态卡片 -->
    <div v-loading="loading">
      <a-row :gutter="16" v-if="subscription">
        <a-col :span="8">
          <a-card class="sub-card">
            <template #title>当前计划</template>
            <div class="sub-plan-name">{{ subscription.plan_name }}</div>
            <a-tag :color="statusColor(subscription.status)">{{ subscription.status }}</a-tag>
          </a-card>
        </a-col>
        <a-col :span="8">
          <a-card class="sub-card">
            <template #title>到期时间</template>
            <div class="sub-value">{{ formatDate(subscription.expires_at) }}</div>
            <div class="sub-hint" v-if="isExpiringSoon">即将到期</div>
          </a-card>
        </a-col>
        <a-col :span="8">
          <a-card class="sub-card">
            <template #title>自动续费</template>
            <div class="sub-value">
              <a-switch v-model="autoRenew" @change="handleAutoRenewChange" />
            </div>
            <div class="sub-hint">到期前自动扣款</div>
          </a-card>
        </a-col>
      </a-row>

      <!-- 无订阅 -->
      <a-empty v-if="!loading && !subscription" description="暂无订阅记录">
        <a-button type="primary" @click="$router.push('/billing/plans')">选择订阅计划</a-button>
      </a-empty>

      <!-- 订阅详情 -->
      <a-card v-if="subscription" class="detail-card" style="margin-top: 16px;">
        <template #title>订阅详情</template>
        <a-descriptions :column="2" bordered>
          <a-descriptions-item label="订阅ID">{{ subscription.id }}</a-descriptions-item>
          <a-descriptions-item label="用户ID">{{ subscription.user_id }}</a-descriptions-item>
          <a-descriptions-item label="开始时间">{{ formatDate(subscription.started_at) }}</a-descriptions-item>
          <a-descriptions-item label="下次扣款日">{{ formatDate(subscription.next_billing_at) }}</a-descriptions-item>
          <a-descriptions-item label="订阅状态">
            <a-tag :color="statusColor(subscription.status)">{{ subscription.status }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="取消时间" v-if="subscription.cancelled_at">
            {{ formatDate(subscription.cancelled_at) }}
          </a-descriptions-item>
        </a-descriptions>
      </a-card>
    </div>

    <!-- 取消确认弹窗 -->
    <a-modal
      v-model:visible="cancelModalVisible"
      title="取消订阅"
      @before-ok="confirmCancel"
    >
      <a-form :model="cancelForm">
        <a-form-item label="取消原因">
          <a-textarea v-model="cancelForm.reason" placeholder="请输入取消原因（可选）" :rows="3" />
        </a-form-item>
      </a-form>
      <div class="cancel-warning">取消后您的订阅将在到期日终止，不再自动扣款。</div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { billingApi, type UserSubscription } from '@/api/billing'

const loading = ref(false)
const subscription = ref<UserSubscription | null>(null)
const autoRenew = ref(true)
const cancelModalVisible = ref(false)
const cancelForm = ref({ reason: '' })

const statusColor = (status: string) => {
  const map: Record<string, string> = {
    active: 'green',
    trial: 'arcoblue',
    cancelled: 'red',
    expired: 'gray'
  }
  return map[status] || 'gray'
}

const formatDate = (dateStr?: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('zh-CN')
}

const isExpiringSoon = computed(() => {
  if (!subscription.value) return false
  const days = (new Date(subscription.value.expires_at).getTime() - Date.now()) / 86400000
  return days <= 7 && days > 0
})

const loadSubscription = async () => {
  loading.value = true
  try {
    const res = await billingApi.getCurrentSubscription()
    if (res.code === 0 || res.code === 200) {
      subscription.value = res.data
      autoRenew.value = res.data?.auto_renew ?? true
    } else {
      // mock
      subscription.value = {
        id: 1,
        user_id: 1,
        plan_id: 2,
        plan_name: '专业版',
        status: 'active',
        started_at: '2026-02-01T00:00:00Z',
        expires_at: '2026-04-01T00:00:00Z',
        next_billing_at: '2026-04-01T00:00:00Z',
        auto_renew: true
      }
      autoRenew.value = true
    }
  } catch (e) {
    subscription.value = null
  } finally {
    loading.value = false
  }
}

const handleRenew = async () => {
  if (!subscription.value) return
  try {
    const res = await billingApi.renewSubscription(subscription.value.id)
    if (res.code === 0 || res.code === 200) {
      Message.success('续费成功')
      loadSubscription()
    } else {
      Message.error('续费失败')
    }
  } catch (e) {
    Message.error('续费失败，请稍后重试')
  }
}

const handleChangePlan = () => {
  Message.info('跳转至变更方案页面')
}

const handleCancel = () => {
  cancelModalVisible.value = true
}

const confirmCancel = async (done: (val: boolean) => void) => {
  if (!subscription.value) { done(false); return }
  try {
    const res = await billingApi.cancelSubscription(subscription.value.id, { reason: cancelForm.value.reason })
    if (res.code === 0 || res.code === 200) {
      Message.success('取消订阅成功')
      cancelModalVisible.value = false
      loadSubscription()
      done(true)
    } else {
      Message.error('取消失败')
      done(false)
    }
  } catch (e) {
    Message.error('取消失败，请稍后重试')
    done(false)
  }
}

const handleAutoRenewChange = (val: boolean) => {
  Message.info(`自动续费已${val ? '开启' : '关闭'}`)
}

onMounted(() => {
  loadSubscription()
})
</script>

<style scoped>
.sub-card :deep(.arco-card-header) {
  font-weight: 600;
}
.sub-plan-name {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 8px;
}
.sub-value {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-1, #1f2329);
}
.sub-hint {
  font-size: 12px;
  color: var(--color-text-3, #86909c);
  margin-top: 4px;
}
.detail-card {
  margin-top: 16px;
}
.cancel-warning {
  color: var(--color-text-3, #86909c);
  font-size: 13px;
  margin-top: 8px;
}
</style>
