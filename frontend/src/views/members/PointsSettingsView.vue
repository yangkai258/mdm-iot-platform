<template>
  <div class="page-container">
    <div class="toolbar">
      <a-button type="primary" @click="handleSave">保存设置</a-button>
      <a-button @click="handleReset">重置</a-button>
    </div>
    <a-form :model="settings" layout="vertical" class="settings-form">
      <a-divider>积分规则</a-divider>
      <a-row :gutter="16">
        <a-col :span="12">
          <a-form-item label="每元积分">
            <a-input-number v-model="settings.points_per_yuan" :min="0" style="width: 100%" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="积分抵现比例">
            <a-input-number v-model="settings.deduct_ratio" :min="0" :max="100" style="width: 100%" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="最高抵扣比例（%）">
            <a-input-number v-model="settings.max_deduct_percent" :min="0" :max="100" style="width: 100%" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="最低兑换积分">
            <a-input-number v-model="settings.min_points" :min="0" style="width: 100%" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="每日积分上限">
            <a-input-number v-model="settings.max_points_per_day" :min="0" style="width: 100%" />
          </a-form-item>
        </a-col>
      </a-row>
      <a-divider>有效期设置</a-divider>
      <a-row :gutter="16">
        <a-col :span="12">
          <a-form-item label="积分过期类型">
            <a-select v-model="settings.expire_type" placeholder="请选择">
              <a-option value="never">永不过期</a-option>
              <a-option value="year">每年底过期</a-option>
              <a-option value="month">每月过期</a-option>
              <a-option value="custom">自定义天数</a-option>
            </a-select>
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="生日双倍积分">
            <a-switch v-model="settings.birthday_double" true-value="1" false-value="0" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="允许积分转让">
            <a-switch v-model="settings.allow_transfer" true-value="1" false-value="0" />
          </a-form-item>
        </a-col>
      </a-row>
    </a-form>
  </div>
</template>

<script setup>

import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'

const settings = reactive({
  deduct_ratio: 100,
  max_deduct_percent: 50,
  min_points: 100,
  points_per_yuan_base: 1,
  points_per_yuan: 1,
  expire_type: 'never',
  max_points_per_day: 0,
  birthday_double: '0',
  allow_transfer: '0'
})

const loadSettings = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/points/settings`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0 && data.data) Object.assign(settings, data.data)
  } catch (e) {}
}

const handleSave = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/points/settings`, {
      method: 'PUT',
      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(settings)
    })
    const data = await res.json()
    if (data.code === 0) Message.success('保存成功')
    else Message.error(data.message || '保存失败')
  } catch (e) { Message.error('保存失败') }
}

const handleReset = () => {
  Object.assign(settings, {
    deduct_ratio: 100, max_deduct_percent: 50, min_points: 100,
    points_per_yuan_base: 1, points_per_yuan: 1, expire_type: 'never',
    max_points_per_day: 0, birthday_double: '0', allow_transfer: '0'
  })
  Message.info('已重置为默认值')
}

onMounted(() => loadSettings())

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.toolbar { margin-bottom: 16px; }
.settings-form { max-width: 800px; }
</style>
