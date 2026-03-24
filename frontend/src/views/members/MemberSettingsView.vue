<template>
  <div class="page-container">
    <div class="toolbar">
      <a-button type="primary" @click="handleSave" :loading="saving">保存设置</a-button>
      <a-button @click="handleReset">重置</a-button>
    </div>
    <a-form :model="form" layout="vertical" class="settings-form">
      <a-divider>积分设置</a-divider>
      <a-row :gutter="16">
        <a-col :span="12">
          <a-form-item label="注册积分">
            <a-input-number v-model="form.registerPoints" :min="0" style="width: 100%" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="积分有效期（天）">
            <a-input-number v-model="form.pointsExpireDays" :min="0" style="width: 100%" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="积分抵扣比例">
            <a-input-number v-model="form.pointsDiscountRatio" :min="0" :precision="2" style="width: 100%" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="积分抵现比例">
            <a-input-number v-model="form.pointsToMoneyRatio" :min="1" style="width: 100%" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="会员有效期（月）">
            <a-input-number v-model="form.memberValidMonths" :min="0" style="width: 100%" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="到期提醒天数">
            <a-input-number v-model="form.expireReminderDays" :min="1" style="width: 100%" />
          </a-form-item>
        </a-col>
      </a-row>
      <a-divider>规则设置</a-divider>
      <a-row :gutter="16">
        <a-col :span="12">
          <a-form-item label="最低积分余额">
            <a-input-number v-model="form.minPointsBalance" :min="0" style="width: 100%" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="每单最高积分">
            <a-input-number v-model="form.maxPointsPerOrder" :min="0" style="width: 100%" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="允许自行注销">
            <a-switch v-model="form.allowSelfCancel" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="手机验证">
            <a-switch v-model="form.requireMobileVerify" />
          </a-form-item>
        </a-col>
      </a-row>
      <a-form-item label="备注">
        <a-textarea v-model="form.remark" :rows="3" placeholder="请输入备注" />
      </a-form-item>
    </a-form>
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
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.toolbar { margin-bottom: 16px; }
.settings-form { max-width: 800px; }
</style>
