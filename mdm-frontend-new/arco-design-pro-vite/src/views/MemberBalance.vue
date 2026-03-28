<template>
  <div class="container">
    <a-row :gutter="16">
      <a-col :span="8">
        <a-card>
          <template #title>当前余额</template>
          <a-statistic :value="balance" :precision="2" prefix="¥" />
          <a-divider />
          <a-space vertical fill>
            <a-button type="primary" long @click="handleRecharge">
              <template #icon><icon-plus /></template>
              充值
            </a-button>
            <a-button long @click="handleDeduct">
              <template #icon><icon-minus /></template>
              扣款
            </a-button>
          </a-space>
        </a-card>
      </a-col>
      <a-col :span="16">
        <a-card title="快速充值">
          <a-space wrap>
            <a-button v-for="amount in quickAmounts" :key="amount" @click="handleQuickRecharge(amount)">
              ¥{{ amount }}
            </a-button>
            <a-button @click="handleCustomAmount">
              <template #icon><icon-edit /></template>
              自定义
            </a-button>
          </a-space>
        </a-card>

        <a-card title="交易流水" style="margin-top: 16px">
          <template #extra>
            <a-space>
              <a-input-search v-model="keyword" placeholder="搜索" style="width: 200px" />
              <a-button @click="handleExport">
                <template #icon><icon-download /></template>
                导出
              </a-button>
            </a-space>
          </template>
          <a-table :columns="columns" :data="transactions" :pagination="pagination" />
        </a-card>
      </a-col>
    </a-row>

    <a-modal v-model:visible="rechargeVisible" title="充值" @ok="handleRechargeSubmit">
      <a-form :model="rechargeForm" layout="vertical">
        <a-form-item label="充值金额" required>
          <a-input-number v-model="rechargeForm.amount" :min="0.01" :precision="2" style="width: 100%" />
        </a-form-item>
        <a-form-item label="支付方式">
          <a-radio-group v-model="rechargeForm.paymentMethod">
            <a-radio value="alipay">支付宝</a-radio>
            <a-radio value="wechat">微信支付</a-radio>
            <a-radio value="bankcard">银行卡</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal v-model:visible="deductVisible" title="扣款" @ok="handleDeductSubmit">
      <a-form :model="deductForm" layout="vertical">
        <a-form-item label="扣款金额" required>
          <a-input-number v-model="deductForm.amount" :min="0.01" :precision="2" style="width: 100%" />
        </a-form-item>
        <a-form-item label="扣款原因" required>
          <a-textarea v-model="deductForm.reason" placeholder="请输入扣款原因" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const balance = ref(2580.50)
const keyword = ref('')
const pagination = reactive({ current: 1, pageSize: 10, total: 50 })
const quickAmounts = [50, 100, 200, 500, 1000, 2000]
const rechargeVisible = ref(false)
const deductVisible = ref(false)

const rechargeForm = reactive({ amount: 0, paymentMethod: 'alipay' })
const deductForm = reactive({ amount: 0, reason: '' })

const columns = [
  { title: '时间', dataIndex: 'time', width: 180 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '金额', slotName: 'amount', width: 120 },
  { title: '余额', dataIndex: 'balanceAfter' },
  { title: '备注', dataIndex: 'remark' }
]
const transactions = ref([
  { time: '2026-03-28 10:00:00', type: 'recharge', amount: '+200.00', balanceAfter: '2580.50', remark: '支付宝充值' },
  { time: '2026-03-27 15:30:00', type: 'deduct', amount: '-50.00', balanceAfter: '2380.50', remark: '订阅扣费' }
])

const handleRecharge = () => { rechargeVisible.value = true }
const handleDeduct = () => { deductVisible.value = true }
const handleQuickRecharge = (amount) => { rechargeForm.amount = amount; rechargeVisible.value = true }
const handleCustomAmount = () => { rechargeVisible.value = true }
const handleExport = () => { }
const handleRechargeSubmit = () => { rechargeVisible.value = false }
const handleDeductSubmit = () => { deductVisible.value = false }
</script>

<style scoped>
.container { padding: 16px; }
</style>
