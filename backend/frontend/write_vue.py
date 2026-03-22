# -*- coding: utf-8 -*-
import os

content = """<template>
  <div class="apply-page">
    <div class="apply-banner">
      <div class="banner-content">
        <h1>MDM 设备管理平台 — 租户入驻申请</h1>
        <p>填写以下信息提交入驻申请，我们的团队将在 1-3 个工作日内完成审核</p>
      </div>
    </div>

    <a-card class="apply-card" :bordered="false">
      <a-steps :current="currentStep" style="margin-bottom: 32px">
        <a-step title="填写企业信息" />
        <a-step title="选择服务套餐" />
        <a-step title="提交申请" />
      </a-steps>

      <a-form ref="formRef" :model="form" :rules="formRules" layout="vertical" style="max-width: 720px; margin: 0 auto">
        <div v-show="currentStep === 0">
          <a-divider orientation="center">企业基本信息</a-divider>
          <a-row :gutter="24">
            <a-col :span="24">
              <a-form-item label="公司名称" field="company_name" required>
                <a-input v-model="form.company_name" placeholder="请输入完整的公司营业执照名称" size="large" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="所属行业" field="industry" required>
                <a-select v-model="form.industry" placeholder="请选择所属行业" size="large">
                  <a-option value="互联网">互联网</a-option>
                  <a-option value="金融科技">金融科技</a-option>
                  <a-option value="制造业">制造业</a-option>
                  <a-option value="零售">零售</a-option>
                  <a-option value="医疗">医疗</a-option>
                  <a-option value="教育">教育</a-option>
                  <a-option value="物流">物流</a-option>
                  <a-option value="能源">能源</a-option>
                  <a-option value="政府/公共事业">政府/公共事业</a-option>
                  <a-option value="其他">其他</a-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="企业规模" field="company_size" required>
                <a-select v-model="form.company_size" placeholder="请选择企业规模" size="large">
                  <a-option value="1-50">1-50 人</a-option>
                  <a-option value="51-200">51-200 人</a-option>
                  <a-option value="201-500">201-500 人</a-option>
                  <a-option value="500+">500 人以上</a-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="24">
              <a-form-item label="公司地址" field="company_address">
                <a-input v-model="form.company_address" placeholder="请输入公司详细地址" size="large" />
              </a-form-item>
            </a-col>
          </a-row>

          <a-divider orientation="center">联系人信息</a-divider>
          <a-row :gutter="24">
            <a-col :span="12">
              <a-form-item label="联系人姓名" field="contact_name" required>
                <a-input v-model="form.contact_name" placeholder="请输入联系人姓名" size="large" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="联系人职务" field="contact_title">
                <a-input v-model="form.contact_title" placeholder="如：技术总监 / 运维负责人" size="large" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="联系电话" field="contact_phone" required>
                <a-input v-model="form.contact_phone" placeholder="请输入手机号或座机" size="large" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="联系邮箱" field="contact_email" required>
                <a-input v-model="form.contact_email" placeholder="请输入工作邮箱" size="large" />
              </a-form-item>
            </a-col>
          </a-row>

          <div class="step-actions">
            <a-button type="primary" size="large" @click="nextStep">下一步：选择套餐</a-button>
          </div>
        </div>

        <div v-show="currentStep === 1">
          <a-divider orientation="center">选择服务套餐</a-divider>
          <a-alert type="info" style="margin-bottom: 20px">选择您需要的套餐版本，所有套餐均支持设备管理核心功能。</a-alert>

          <a-radio-group v-model="form.plan_id" direction="vertical" style="width: 100%">
            <a-card v-for="plan in plans" :key="plan.id" :class="['plan-card', { 'plan-card-selected': form.plan_id === plan.id }]" hoverable @click="form.plan_id = plan.id">
              <a-radio :value="plan.id">
                <div class="plan-row">
                  <div class="plan-info">
                    <h3>{{ plan.plan_name }}</h3>
                    <p class="plan-desc">{{ plan.description }}</p>
                    <div class="plan-tags">
                      <a-tag v-for="tag in plan.tags" :key="tag" size="small">{{ tag }}</a-tag>
                    </div>
                  </div>
                  <div class="plan-pricing">
                    <span class="price">¥{{ plan.price_yearly }}</span>
                    <span class="period">/ 年</span>
                    <div class="plan-quota">
                      <span>设备上限: {{ plan.device_quota === -1 ? '不限' : plan.device_quota + ' 台' }}</span>
                      <span>用户上限: {{ plan.user_quota === -1 ? '不限' : plan.user_quota + ' 人' }}</span>
                    </div>
                  </div>
                </div>
              </a-radio>
            </a-card>
          </a-radio-group>

          <a-divider orientation="center">使用场景（选填）</a-divider>
          <a-form-item label="使用场景描述" field="use_case">
            <a-textarea v-model="form.use_case" placeholder="请简要描述您的业务场景和需求" :rows="4" />
          </a-form-item>

          <div class="step-actions">
            <a-button size="large" @click="currentStep = 0">上一步</a-button>
            <a-button type="primary" size="large" @click="nextStep2">下一步：确认提交</a-button>
          </div>
        </div>

        <div v-show="currentStep === 2">
          <a-divider orientation="center">申请信息确认</a-divider>
          <a-card class="confirm-card" :bordered="false">
            <a-descriptions :column="2" bordered size="large">
              <a-descriptions-item label="公司名称" :span="2">{{ form.company_name }}</a-descriptions-item>
              <a-descriptions-item label="所属行业">{{ form.industry }}</a-descriptions-item>
              <a-descriptions-item label="企业规模">{{ form.company_size }}</a-descriptions-item>
              <a-descriptions-item label="公司地址" :span="2">{{ form.company_address || '-' }}</a-descriptions-item>
              <a-descriptions-item label="联系人">{{ form.contact_name }}</a-descriptions-item>
              <a-descriptions-item label="联系人职务">{{ form.contact_title || '-' }}</a-descriptions-item>
              <a-descriptions-item label="联系电话">{{ form.contact_phone }}</a-descriptions-item>
              <a-descriptions-item label="联系邮箱">{{ form.contact_email }}</a-descriptions-item>
              <a-descriptions-item label="申请套餐" :span="2">
                <a-tag color="arcoblue" size="large">{{ selectedPlan?.plan_name }} — ¥{{ selectedPlan?.price_yearly }}/年</a-tag>
              </a-descriptions-item>
              <a-descriptions-item label="使用场景" :span="2">{{ form.use_case || '-' }}</a-descriptions-item>
            </a-descriptions>
          </a-card>

          <a-alert type="warning" style="margin-top: 16px">请确认以上信息准确无误。提交后，审核团队将在 1-3 个工作日内完成审核，并通过邮箱通知结果。</a-alert>

          <div class="step-actions">
            <a-button size="large" @click="currentStep = 1">修改信息</a-button>
            <a-button type="primary" status="success" size="large" :loading="submitting" @click="handleSubmit">提交申请</a-button>
          </div>
        </div>
      </a-form>
    </a-card>

    <a-modal v-model:visible="successVisible" :title="null" :footer="null" :mask-closable="false" :closable="false" width="480px">
      <a-result status="success" title="申请已提交！">
        <template #subtitle>
          <p>您的入驻申请已成功提交，申请编号：<strong>{{ applicationCode }}</strong></p>
          <p>审核结果将在 1-3 个工作日内发送至您的邮箱：<br /><strong>{{ form.contact_email }}</strong></p>
        </template>
        <template #extra>
          <a-space direction="vertical" fill>
            <a-button type="primary" long @click="goToLogin">返回登录页</a-button>
            <a-button long @click="resetForm">继续申请（测试）</a-button>
          </a-space>
        </template>
      </a-result>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'

const router = useRouter()
const currentStep = ref(0)
const submitting = ref(false)
const successVisible = ref(false)
const applicationCode = ref('')
const formRef = ref()

const form = ref({
  company_name: '',
  industry: '',
  company_size: '',
  company_address: '',
  contact_name: '',
  contact_title: '',
  contact_phone: '',
  contact_email: '',
  plan_id: 2,
  use_case: '',
})

const formRules = {
  company_name: [{ required: true, message: '请输入公司名称' }],
  industry: [{ required: true, message: '请选择所属行业' }],
  company_size: [{ required: true, message: '请选择企业规模' }],
  contact_name: [{ required: true, message: '请输入联系人姓名' }],
  contact_phone: [{ required: true, message: '请输入联系电话' }],
  contact_email: [{ required: true, message: '请输入联系邮箱' }],
  plan_id: [{ required: true, message: '请选择套餐' }],
}

const plans = [
  { id: 1, plan_name: '免费版', description: '适合个人或小团队试用，支持基础设备管理功能', price_yearly: 0, device_quota: 10, user_quota: 5, tags: ['基础设备管理', '7天数据存储'] },
  { id: 2, plan_name: '基础版', description: '适合小微企业，支持 OTA 固件升级和远程控制', price_yearly: 299, device_quota: 50, user_quota: 10, tags: ['OTA 固件升级', '远程控制', '90天数据存储'] },
  { id: 3, plan_name: '专业版', description: '适合成长期企业，支持更多设备配额和数据分析', price_yearly: 799, device_quota: 200, user_quota: 20, tags: ['OTA 固件升级', '远程控制', '数据导出', '自定义Logo', '1年数据存储'] },
  { id: 4, plan_name: '企业版', description: '适合大型企业，支持无限配额和全部高级功能，专属客服', price_yearly: 1999, device_quota: 500, user_quota: 50, tags: ['全部功能', 'API开放接口', '专属客服', '无限数据存储'] },
]

const selectedPlan = computed(() => plans.find(p => p.id === form.value.plan_id))

const fieldLabels = {
  company_name: '公司名称',
  industry: '所属行业',
  company_size: '企业规模',
  contact_name: '联系人姓名',
  contact_phone: '联系电话',
  contact_email: '联系邮箱',
}

const nextStep = () => {
  const required = ['company_name', 'industry', 'company_size', 'contact_name', 'contact_phone', 'contact_email']
  for (const key of required) {
    if (!form.value[key]) {
      Message.warning('请填写：' + fieldLabels[key])
      return
    }
  }
  currentStep.value = 1
}

const nextStep2 = () => {
  if (!form.value.plan_id) {
    Message.warning('请选择一个套餐')
    return
  }
  currentStep.value = 2
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    await new Promise(resolve => setTimeout(resolve, 1500))
    const now = new Date()
    applicationCode.value = 'APP-' + now.getFullYear() + String(now.getMonth() + 1).padStart(2, '0') + String(now.getDate()).padStart(2, '0') + '-' + String(Math.floor(Math.random() * 900) + 100)
    successVisible.value = true
  } catch (e) {
    Message.error('提交失败，请稍后重试')
  } finally {
    submitting.value = false
  }
}

const goToLogin = () => { router.push('/login') }

const resetForm = () => {
  form.value = { company_name: '', industry: '', company_size: '', company_address: '', contact_name: '', contact_title: '', contact_phone: '', contact_email: '', plan_id: 2, use_case: '' }
  currentStep.value = 0
  successVisible.value = false
}
</script>

<style scoped>
.apply-page { min-height: 100vh; background: #f5f7fa; padding-bottom: 40px; }
.apply-banner { background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%); color: #fff; padding: 40px 24px; margin-bottom: 0; }
.banner-content { max-width: 720px; margin: 0 auto; text-align: center; }
.banner-content h1 { font-size: 24px; font-weight: 700; margin-bottom: 12px; letter-spacing: 1px; }
.banner-content p { font-size: 15px; color: rgba(255,255,255,0.75); }
.apply-card { max-width: 820px; margin: -24px auto 0; border-radius: 12px; box-shadow: 0 4px 20px rgba(0,0,0,0.08); padding: 8px; }
.step-actions { display: flex; justify-content: center; gap: 16px; margin-top: 32px; }
.plan-card { margin-bottom: 16px; cursor: pointer; border: 2px solid transparent; transition: all 0.2s; }
.plan-card:hover { border-color: #c0d8ff; }
.plan-card-selected { border-color: #1650ff !important; background: #f2f5ff; }
.plan-row { display: flex; align-items: flex-start; justify-content: space-between; width: 100%; padding: 4px 0; }
.plan-info { flex: 1; }
.plan-info h3 { margin: 0 0 6px; font-size: 16px; font-weight: 600; }
.plan-desc { color: #4e5969; font-size: 13px; margin: 0 0 8px; }
.plan-tags { display: flex; flex-wrap: wrap; gap: 6px; }
.plan-pricing { text-align: right; flex-shrink: 0; margin-left: 24px; }
.price { font-size: 26px; font-weight: 700; color: #1650ff; }
.period { font-size: 13px; color: #86909c; }
.plan-quota { margin-top: 6px; display: flex; flex-direction: column; align-items: flex-end; font-size: 12px; color: #4e5969; }
.confirm-card { background: #fafafa; margin-bottom: 0; }
</style>
"""

path = "C:/Users/YKing/.openclaw/workspace/mdm-project/frontend/frontend/src/views/tenants/TenantApplication.vue"
with open(path, 'w', encoding='utf-8') as f:
    f.write(content)
print(f"Written {len(content)} bytes to {path}")
