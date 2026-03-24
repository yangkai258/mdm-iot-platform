<template>
  <div class="member-settings-page">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员参数设置</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 表单区域 -->
    <a-card class="form-card">
      <template #title>
        <span style="font-weight: 600; font-size: 15px;">会员参数设置</span>
      </template>
      <template #extra>
        <a-space>
          <a-button @click="handleReset">重置</a-button>
          <a-button type="primary" :loading="saving" @click="handleSave">保存</a-button>
        </a-space>
      </template>

      <a-form :model="form" layout="vertical" ref="formRef" class="settings-form">
        <a-divider orientation="left">注册与积分</a-divider>
        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="注册赠送积分">
              <a-input-number v-model="form.registerPoints" :min="0" :max="999999" style="width: 100%" />
              <template #extra>新会员注册时赠送的积分数量，0表示不赠送</template>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="积分过期天数">
              <a-input-number v-model="form.pointsExpireDays" :min="0" :max="3650" style="width: 100%" />
              <template #extra>积分自获得之日起多少天后过期，0表示永不过期</template>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="积分抵扣比例">
              <a-input-number v-model="form.pointsDiscountRatio" :min="0" :max="1" :step="0.01" style="width: 100%" />
              <template #extra>积分抵扣现金的比例，如 0.01 表示100积分抵扣1元</template>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="积分与人民币比例">
              <a-input-number v-model="form.pointsToMoneyRatio" :min="0" :max="1000" :step="1" style="width: 100%" />
              <template #extra>消费1元可获得的积分数，如 100 表示1元=100积分</template>
            </a-form-item>
          </a-col>
        </a-row>

        <a-divider orientation="left">会员有效期</a-divider>
        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="会员有效期（月）">
              <a-input-number v-model="form.memberValidMonths" :min="0" :max="120" style="width: 100%" />
              <template #extra>会员注册后的有效期月数，0表示永久有效</template>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="会员到期提醒（天）">
              <a-input-number v-model="form.expireReminderDays" :min="0" :max="90" style="width: 100%" />
              <template #extra>会员到期前多少天发送提醒，0表示不提醒</template>
            </a-form-item>
          </a-col>
        </a-row>

        <a-divider orientation="left">其他设置</a-divider>
        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="允许自行注销">
              <a-radio-group v-model="form.allowSelfCancel">
                <a-radio :value="true">允许</a-radio>
                <a-radio :value="false">不允许</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="注册需手机验证">
              <a-radio-group v-model="form.requireMobileVerify">
                <a-radio :value="true">是</a-radio>
                <a-radio :value="false">否</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="积分下限">
              <a-input-number v-model="form.minPointsBalance" :min="0" :max="999999" style="width: 100%" />
              <template #extra>积分余额低于此值时不可抵扣</template>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="单次积分抵扣上限">
              <a-input-number v-model="form.maxPointsPerOrder" :min="0" :max="999999" style="width: 100%" />
              <template #extra>单次订单积分抵扣的上限，0表示不限制</template>
            </a-form-item>
          </a-col>
        </a-row>

        <a-divider orientation="left">备注</a-divider>
        <a-form-item label="系统备注">
          <a-textarea v-model="form.remark" :rows="3" placeholder="系统配置备注信息" />
        </a-form-item>
      </a-form>

      <div class="form-footer">
        <a-space>
          <a-button @click="handleReset">重置</a-button>
          <a-button type="primary" :loading="saving" @click="handleSave">保存</a-button>
        </a-space>
      </div>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as api from '@/api/member'

const saving = ref(false)
const originalForm = {}

const form = reactive({
  registerPoints: 0,
  pointsExpireDays: 0,
  pointsDiscountRatio: 0.01,
  pointsToMoneyRatio: 100,
  memberValidMonths: 0,
  expireReminderDays: 7,
  allowSelfCancel: false,
  requireMobileVerify: true,
  minPointsBalance: 0,
  maxPointsPerOrder: 0,
  remark: ''
})

const loadSettings = async () => {
  try {
    const res = await api.getMemberSettings()
    const data = res.data || {}
    Object.assign(form, data)
    Object.assign(originalForm, data)
  } catch (err) {
    // ignore, use defaults
  }
}

const handleSave = async () => {
  saving.value = true
  try {
    await api.updateMemberSettings({ ...form })
    Message.success('保存成功')
    Object.assign(originalForm, form)
  } catch (err) {
    Message.error(err.message || '保存失败')
  } finally {
    saving.value = false
  }
}

const handleReset = () => {
  Object.assign(form, originalForm)
  Message.info('已重置为保存的值')
}

onMounted(() => {
  loadSettings()
})
</script>

<style scoped>
.member-settings-page {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}
.breadcrumb { margin-bottom: 16px; }
.form-card { border-radius: 8px; }
.settings-form { max-width: 900px; }
.form-footer {
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #e5e7eb;
}
</style>
