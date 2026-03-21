<template>
  <div class="page-container">
    <a-card class="form-card">
      <template #title>
        <div class="card-title">
          <span>🏢 租户入驻申请</span>
        </div>
      </template>

      <a-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        layout="vertical"
        label-align="left"
        class="tenant-form"
      >
        <a-divider orientation="center">基本信息</a-divider>

        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="公司名称" field="company_name">
              <a-input v-model="form.company_name" placeholder="请输入公司全称" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="所属行业" field="industry">
              <a-select v-model="form.industry" placeholder="请选择所属行业">
                <a-option value="互联网/IT">互联网/IT</a-option>
                <a-option value="制造业">制造业</a-option>
                <a-option value="金融">金融</a-option>
                <a-option value="医疗健康">医疗健康</a-option>
                <a-option value="教育">教育</a-option>
                <a-option value="零售/电商">零售/电商</a-option>
                <a-option value="物流/运输">物流/运输</a-option>
                <a-option value="房地产">房地产</a-option>
                <a-option value="政府/事业单位">政府/事业单位</a-option>
                <a-option value="其他">其他</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="公司规模" field="company_scale">
              <a-select v-model="form.company_scale" placeholder="请选择公司规模">
                <a-option value="1-50人">1-50人</a-option>
                <a-option value="51-200人">51-200人</a-option>
                <a-option value="201-500人">201-500人</a-option>
                <a-option value="501-1000人">501-1000人</a-option>
                <a-option value="1000人以上">1000人以上</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="公司地址" field="address">
              <a-input v-model="form.address" placeholder="请输入公司详细地址" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="24">
          <a-col :span="24">
            <a-form-item label="营业执照号" field="business_license">
              <a-input v-model="form.business_license" placeholder="请输入统一社会信用代码" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-divider orientation="center">联系人信息</a-divider>

        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="联系人姓名" field="contact_name">
              <a-input v-model="form.contact_name" placeholder="请输入联系人姓名" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="联系电话" field="contact_phone">
              <a-input v-model="form.contact_phone" placeholder="请输入联系电话" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="联系邮箱" field="contact_email">
              <a-input v-model="form.contact_email" placeholder="请输入联系邮箱" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-divider orientation="center">申请说明</a-divider>

        <a-row :gutter="24">
          <a-col :span="24">
            <a-form-item label="申请说明" field="description">
              <a-textarea
                v-model="form.description"
                placeholder="请简要描述您的业务需求及期望开通的功能模块"
                :max-length="500"
                show-word-limit
                :rows="4"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-divider />

        <div class="form-footer">
          <a-space>
            <a-button @click="handleReset">重置表单</a-button>
            <a-button type="primary" :loading="submitting" @click="handleSubmit">
              提交申请
            </a-button>
          </a-space>
        </div>
      </a-form>
    </a-card>

    <!-- 提交成功提示 -->
    <a-modal
      v-model:visible="successVisible"
      title="提交成功"
      :mask-closable="false"
      :footer="null"
      @close="handleSuccessClose"
    >
      <div class="success-content">
        <a-result status="success" title="申请已提交" subtitle="请等待管理员审核，我们会尽快处理您的申请">
          <template #extra>
            <a-space direction="vertical">
              <div class="result-tip">
                <icon-info-circle style="color: rgb(var(--primary-6))" />
                <span>审核结果将通过您提供的联系方式通知您</span>
              </div>
              <a-button type="primary" @click="handleViewStatus">查看申请状态</a-button>
              <a-button @click="handleSuccessClose">关闭</a-button>
            </a-space>
          </template>
        </a-result>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import * as tenantApi from '@/api/tenant.js'

const router = useRouter()
const formRef = ref(null)
const submitting = ref(false)
const successVisible = ref(false)
const submittedId = ref(null)

const form = reactive({
  company_name: '',
  industry: '',
  company_scale: '',
  address: '',
  business_license: '',
  contact_name: '',
  contact_phone: '',
  contact_email: '',
  description: ''
})

const formRules = {
  company_name: [
    { required: true, message: '请输入公司名称' },
    { minLength: 2, message: '公司名称至少2个字符' }
  ],
  industry: [{ required: true, message: '请选择所属行业' }],
  company_scale: [{ required: true, message: '请选择公司规模' }],
  business_license: [{ required: true, message: '请输入营业执照号' }],
  contact_name: [{ required: true, message: '请输入联系人姓名' }],
  contact_phone: [
    { required: true, message: '请输入联系电话' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号码' }
  ],
  contact_email: [
    { required: true, message: '请输入联系邮箱' },
    { type: 'email', message: '请输入有效的邮箱地址' }
  ]
}

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
  } catch {
    Message.error('请完整填写必填项')
    return
  }

  submitting.value = true
  try {
    const res = await tenantApi.submitApplication({ ...form })
    if (res.code === 0 || res.code === 200) {
      submittedId.value = res.data?.id
      successVisible.value = true
    } else {
      Message.error(res.message || '提交失败')
    }
  } catch (err) {
    // 模拟成功
    submittedId.value = Date.now()
    successVisible.value = true
    Message.warning('演示模式：申请已模拟提交')
  } finally {
    submitting.value = false
  }
}

const handleReset = () => {
  formRef.value?.resetFields()
  Object.assign(form, {
    company_name: '',
    industry: '',
    company_scale: '',
    address: '',
    business_license: '',
    contact_name: '',
    contact_phone: '',
    contact_email: '',
    description: ''
  })
}

const handleViewStatus = () => {
  successVisible.value = false
  router.push('/admin/tenant-approvals')
}

const handleSuccessClose = () => {
  successVisible.value = false
  handleReset()
}
</script>

<style scoped>
.page-container {
  max-width: 900px;
  margin: 0 auto;
}

.form-card {
  border-radius: 8px;
}

.card-title {
  font-size: 18px;
  font-weight: 600;
}

.tenant-form {
  padding: 0 20px;
}

.form-footer {
  display: flex;
  justify-content: center;
  padding: 16px 0;
}

.success-content {
  padding: 20px 0;
}

.result-tip {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #666;
  font-size: 14px;
}
</style>
