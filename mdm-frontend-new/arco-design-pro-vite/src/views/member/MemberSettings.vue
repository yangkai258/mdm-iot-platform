<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="member-settings-page">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员参数设置</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 表单区域 -->
    <a-card class="settings-card">
      <template #title>
        <span style="font-weight: 600; font-size: 15px;">会员参数配置</span>
      </template>
      <template #extra>
        <a-space>
          <a-button @click="handleReset">「重置」</a-button>
          <a-button type="primary" :loading="saveLoading" @click="handleSave">「保存」</a-button>
        </a-space>
      </template>

      <a-form :model="form" layout="vertical" ref="formRef">
        <!-- 基础设置 -->
        <a-divider>基础设置</a-divider>
        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="默认会员等级" field="defaultLevelId">
              <a-select v-model="form.defaultLevelId" placeholder="选择默认会员等级">
                <a-option v-for="lv in levelList" :key="lv.id" :value="lv.id">{{ lv.name }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="会员编号前缀" field="memberNoPrefix">
              <a-input v-model="form.memberNoPrefix" placeholder="如：MEM" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="注册方式" field="registerWays">
              <a-checkbox-group v-model="form.registerWays">
                <a-checkbox value="mobile">手机号注册</a-checkbox>
                <a-checkbox value="wechat">微信注册</a-checkbox>
                <a-checkbox value="name">姓名注册</a-checkbox>
              </a-checkbox-group>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="启用会员卡" field="enableCard">
              <a-switch v-model="form.enableCard" />
            </a-form-item>
          </a-col>
        </a-row>

        <!-- 积分设置 -->
        <a-divider>积分设置</a-divider>
        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="积分名称" field="pointsName">
              <a-input v-model="form.pointsName" placeholder="如：积分" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="消费1元兑换积分数" field="consumeToPoints">
              <a-input-number v-model="form.consumeToPoints" :min="0" :step="1" style="width: 100%">
                <template #suffix>分</template>
              </a-input-number>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="积分抵现比例" field="pointsToMoney">
              <a-input-number v-model="form.pointsToMoney" :min="0" :step="0.01" style="width: 100%">
                <template #suffix>积分=1元</template>
              </a-input-number>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="积分有效期" field="pointsExpireDays">
              <a-input-number v-model="form.pointsExpireDays" :min="0" style="width: 100%">
                <template #suffix>天（0表示永久有效）</template>
              </a-input-number>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="生日赠送积分" field="birthdayPoints">
              <a-input-number v-model="form.birthdayPoints" :min="0" style="width: 100%">
                <template #suffix>分</template>
              </a-input-number>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="生日月份赠送积分" field="birthdayMonthPoints">
              <a-input-number v-model="form.birthdayMonthPoints" :min="0" style="width: 100%">
                <template #suffix>分</template>
              </a-input-number>
            </a-form-item>
          </a-col>
        </a-row>

        <!-- 升级设置 -->
        <a-divider>升级设置</a-divider>
        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="自动升级" field="autoUpgrade">
              <a-switch v-model="form.autoUpgrade" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="升级判断周期" field="upgradeCheckCycle">
              <a-select v-model="form.upgradeCheckCycle" placeholder="选择判断周期">
                <a-option value="daily">每天</a-option>
                <a-option value="weekly">每周</a-option>
                <a-option value="monthly">每月</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <!-- 通知设置 -->
        <a-divider>通知设置</a-divider>
        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="注册成功通知" field="notifyOnRegister">
              <a-switch v-model="form.notifyOnRegister" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="升级成功通知" field="notifyOnUpgrade">
              <a-switch v-model="form.notifyOnUpgrade" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="积分变动通知" field="notifyOnPointsChange">
              <a-switch v-model="form.notifyOnPointsChange" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="优惠券到期提醒" field="notifyOnCouponExpire">
              <a-switch v-model="form.notifyOnCouponExpire" />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as api from '@/api/member'

const levelList = ref([])
const saveLoading = ref(false)

const form = reactive({
  defaultLevelId: undefined,
  memberNoPrefix: 'MEM',
  registerWays: ['mobile'],
  enableCard: true,
  pointsName: '积分',
  consumeToPoints: 1,
  pointsToMoney: 100,
  pointsExpireDays: 365,
  birthdayPoints: 100,
  birthdayMonthPoints: 50,
  autoUpgrade: true,
  upgradeCheckCycle: 'monthly',
  notifyOnRegister: true,
  notifyOnUpgrade: true,
  notifyOnPointsChange: true,
  notifyOnCouponExpire: true
})

const defaultForm = { ...form }

const loadSettings = async () => {
  try {
    const res = await api.getMemberSettings()
    const data = res.data || {}
    Object.assign(form, {
      defaultLevelId: data.defaultLevelId,
      memberNoPrefix: data.memberNoPrefix || 'MEM',
      registerWays: data.registerWays || ['mobile'],
      enableCard: data.enableCard ?? true,
      pointsName: data.pointsName || '积分',
      consumeToPoints: data.consumeToPoints ?? 1,
      pointsToMoney: data.pointsToMoney ?? 100,
      pointsExpireDays: data.pointsExpireDays ?? 365,
      birthdayPoints: data.birthdayPoints ?? 100,
      birthdayMonthPoints: data.birthdayMonthPoints ?? 50,
      autoUpgrade: data.autoUpgrade ?? true,
      upgradeCheckCycle: data.upgradeCheckCycle || 'monthly',
      notifyOnRegister: data.notifyOnRegister ?? true,
      notifyOnUpgrade: data.notifyOnUpgrade ?? true,
      notifyOnPointsChange: data.notifyOnPointsChange ?? true,
      notifyOnCouponExpire: data.notifyOnCouponExpire ?? true
    })
    Object.assign(defaultForm, { ...form })
  } catch (err) { /* use defaults */ }
}

const loadLevels = async () => {
  try {
    const res = await api.getLevelList()
    levelList.value = res.data || []
  } catch (err) { /* ignore */ }
}

const handleSave = async () => {
  saveLoading.value = true
  try {
    await api.updateMemberSettings({ ...form })
    Message.success('保存成功')
    Object.assign(defaultForm, { ...form })
  } catch (err) {
    Message.error('保存失败: ' + err.message)
  } finally {
    saveLoading.value = false
  }
}

const handleReset = () => {
  Object.assign(form, defaultForm)
  Message.info('已重置为已保存的配置')
}

onMounted(() => { loadSettings(); loadLevels() })
</script>

<style scoped>
.member-settings-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.settings-card { border-radius: 8px; }
</style>
