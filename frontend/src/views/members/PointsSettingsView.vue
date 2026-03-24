<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>积分规则设置</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="抵扣比例" :value="settings.deduct_ratio" suffix="积分 = 1元" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="单笔最高抵扣" :value="settings.max_deduct_percent" suffix="%" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="最低抵扣门槛" :value="settings.min_points" suffix="积分" />
        </a-card>
      </a-col>
    </a-row>

    <div class="pro-content-area">
      <a-form :model="settings" layout="vertical" style="max-width: 640px;">
        <a-divider>积分抵扣规则</a-divider>

        <a-form-item label="积分抵扣比例">
          <a-space>
            <a-input-number v-model="settings.deduct_ratio" :min="1" style="width: 100px;" />
            <span>积分 = </span>
            <span>1 元</span>
          </a-space>
          <div style="color: #999; font-size: 12px; margin-top: 4px;">例：100积分 = 1元</div>
        </a-form-item>

        <a-form-item label="单笔订单最高抵扣比例">
          <a-slider v-model="settings.max_deduct_percent" :min="0" :max="100" :step="5" />
          <div style="color: #999; font-size: 12px;">当前值：{{ settings.max_deduct_percent }}%（0=不限制）</div>
        </a-form-item>

        <a-form-item label="最低抵扣门槛">
          <a-input-number v-model="settings.min_points" :min="0" :step="100" style="width: 200px;" />
          <div style="color: #999; font-size: 12px; margin-top: 4px;">积分不足此值时不可抵扣，0=不限制</div>
        </a-form-item>

        <a-divider>积分获取规则</a-divider>

        <a-form-item label="消费积分基础倍率">
          <a-space>
            <span>每消费</span>
            <a-input-number v-model="settings.points_per_yuan_base" :min="1" style="width: 80px;" />
            <span>元获得</span>
            <a-input-number v-model="settings.points_per_yuan" :min="0" :precision="1" style="width: 80px;" />
            <span>积分</span>
          </a-space>
        </a-form-item>

        <a-form-item label="积分到期时间">
          <a-select v-model="settings.expire_type" style="width: 300px;">
            <a-option value="never">永不过期</a-option>
            <a-option value="year">每年12月31日过期</a-option>
            <a-option value="monthly">获取后12个月过期</a-option>
          </a-select>
        </a-form-item>

        <a-form-item label="积分获取上限">
          <a-input-number v-model="settings.max_points_per_day" :min="0" style="width: 200px;" />
          <div style="color: #999; font-size: 12px; margin-top: 4px;">每日积分获取上限，0=不限制</div>
        </a-form-item>

        <a-form-item label="生日双倍积分">
          <a-switch v-model="settings.birthday_double" checked-value="1" unchecked-value="0" />
          <span style="margin-left: 8px; color: #999;">开启后会员生日当天消费获得双倍积分</span>
        </a-form-item>

        <a-divider>其他设置</a-divider>

        <a-form-item label="允许积分互转">
          <a-switch v-model="settings.allow_transfer" checked-value="1" unchecked-value="0" />
          <span style="margin-left: 8px; color: #999;">允许会员之间互相转让积分</span>
        </a-form-item>

        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSave">保存</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </div>
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
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { border-radius: 8px; text-align: center; }
.pro-content-area { background: #fff; border-radius: 8px; padding: 24px; box-shadow: 0 1px 3px rgba(0,0,0,0.04); }
</style>
