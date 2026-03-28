<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-sync /> 自动续费配置</a-space>
      </template>

      <a-row :gutter="16">
        <a-col :span="16">
          <a-card title="续费设置">
            <a-form :model="formData" layout="vertical">
              <a-form-item label="自动续费">
                <a-switch v-model="formData.autoRenewEnabled" />
              </a-form-item>
              <a-form-item label="扣款方式">
                <a-select v-model="formData.paymentMethod" style="width: 300px">
                  <a-option value="alipay">支付宝</a-option>
                  <a-option value="wechat">微信支付</a-option>
                  <a-option value="bankcard">银行卡</a-option>
                </a-select>
              </a-form-item>
              <a-form-item label="续费提醒">
                <a-checkbox-group v-model="formData.reminderDays">
                  <a-checkbox :value="3">提前3天</a-checkbox>
                  <a-checkbox :value="7">提前7天</a-checkbox>
                  <a-checkbox :value="15">提前15天</a-checkbox>
                </a-checkbox-group>
              </a-form-item>
              <a-form-item label="扣款失败处理策略">
                <a-radio-group v-model="formData.failStrategy">
                  <a-radio value="pause">暂停服务</a-radio>
                  <a-radio value="notify">发送通知</a-radio>
                  <a-radio value="retry">自动重试</a-radio>
                </a-radio-group>
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="handleSave">保存配置</a-button>
                <a-button @click="handleTest">测试扣款</a-button>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>

        <a-col :span="8">
          <a-card title="续费预览">
            <a-descriptions :column="1" bordered>
              <a-descriptions-item label="当前套餐">专业版</a-descriptions-item>
              <a-descriptions-item label="下次扣款日">2026-04-28</a-descriptions-item>
              <a-descriptions-item label="扣款金额">¥99.00</a-descriptions-item>
              <a-descriptions-item label="续费周期">每月</a-descriptions-item>
            </a-descriptions>
          </a-card>

          <a-card title="历史扣款记录" style="margin-top: 16px">
            <a-table :columns="historyColumns" :data="history" size="small" :pagination="false" />
          </a-card>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const formData = reactive({
  autoRenewEnabled: true, paymentMethod: 'alipay',
  reminderDays: [3, 7], failStrategy: 'retry'
})

const historyColumns = [
  { title: '时间', dataIndex: 'time' },
  { title: '金额', dataIndex: 'amount' },
  { title: '状态', dataIndex: 'status' }
]
const history = ref([
  { time: '2026-03-28 09:00', amount: '¥99.00', status: '成功' },
  { time: '2026-02-28 09:00', amount: '¥99.00', status: '成功' },
  { time: '2026-01-28 09:00', amount: '¥99.00', status: '成功' }
])

const handleSave = () => { }
const handleTest = () => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
