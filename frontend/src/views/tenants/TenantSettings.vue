<template>
  <div class="page-container">
    <a-card class="general-card" title="租户配置">
      <a-tabs default-active-tab="info">
        <!-- 租户信息 -->
        <a-tab-pane key="info" title="基本信息">
          <a-form
            ref="infoFormRef"
            :model="infoForm"
            :rules="infoFormRules"
            layout="vertical"
            style="max-width: 700px"
          >
            <a-row :gutter="24">
              <a-col :span="12">
                <a-form-item label="租户名称" field="tenant_name">
                  <a-input v-model="infoForm.tenant_name" placeholder="请输入租户显示名称" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="联系人" field="contact_name">
                  <a-input v-model="infoForm.contact_name" placeholder="请输入联系人姓名" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="联系电话" field="contact_phone">
                  <a-input v-model="infoForm.contact_phone" placeholder="请输入联系电话" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="联系邮箱" field="contact_email">
                  <a-input v-model="infoForm.contact_email" placeholder="请输入联系邮箱" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="自定义域名" field="custom_domain">
                  <a-input v-model="infoForm.custom_domain" placeholder="yourcompany.mdm.com">
                    <template #prefix>https://</template>
                  </a-input>
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="Logo URL">
                  <a-input v-model="infoForm.tenant_logo" placeholder="https://cdn.example.com/logo.png" />
                </a-form-item>
              </a-col>
            </a-row>

            <a-divider>当前状态</a-divider>
            <a-descriptions :column="3" bordered size="small">
              <a-descriptions-item label="租户状态">
                <a-tag :color="getStatusColor(infoForm.status)">{{ getStatusText(infoForm.status) }}</a-tag>
              </a-descriptions-item>
              <a-descriptions-item label="当前套餐">
                <a-tag color="arcoblue">{{ infoForm.plan_name }}</a-tag>
              </a-descriptions-item>
              <a-descriptions-item label="到期时间">{{ infoForm.expire_time || '未设置' }}</a-descriptions-item>
            </a-descriptions>

            <a-divider />
            <div class="form-actions">
              <a-button type="primary" @click="submitInfo">「保存」</a-button>
              <a-button @click="resetInfo">「取消」</a-button>
            </div>
          </a-form>
        </a-tab-pane>

        <!-- 套餐升级/续费 -->
        <a-tab-pane key="plan" title="套餐管理">
          <a-alert type="info" style="margin-bottom: 16px">
            当前套餐：<strong>{{ infoForm.plan_name }}</strong>，到期时间：<strong>{{ infoForm.expire_time || '未设置' }}</strong>
          </a-alert>

          <a-row :gutter="[24, 24]" class="plan-cards">
            <a-col v-for="plan in plans" :key="plan.id" :span="6">
              <a-card
                :class="{ 'plan-card-active': plan.id === infoForm.plan_id }"
                :bordered="plan.id === infoForm.plan_id"
                hoverable
              >
                <div class="plan-header">
                  <h3>{{ plan.plan_name }}</h3>
                  <p class="plan-price">
                    <span class="price">¥{{ plan.price_yearly }}</span>
                    <span class="period">/{{ plan.period }}</span>
                  </p>
                </div>
                <a-divider />
                <ul class="plan-features">
                  <li>设备上限：{{ plan.device_quota === -1 ? '不限' : plan.device_quota + '台' }}</li>
                  <li>用户上限：{{ plan.user_quota === -1 ? '不限' : plan.user_quota + '人' }}</li>
                  <li>部门/门店：{{ plan.dept_quota }}/{{ plan.store_quota }}</li>
                  <li v-for="f in plan.feature_list" :key="f">{{ f }}</li>
                </ul>
                <a-button
                  v-if="plan.id !== infoForm.plan_id"
                  :type="isUpgradable(plan) ? 'primary' : 'outline'"
                  long
                  @click="openPlanChange(plan)"
                >
                  {{ isUpgradable(plan) ? '升级' : '切换' }}
                </a-button>
                <a-button v-else type="outline" long disabled>当前套餐</a-button>
              </a-card>
            </a-col>
          </a-row>

          <a-divider>续费</a-divider>
          <a-form :model="renewForm" layout="inline">
            <a-form-item label="续费时长">
              <a-select v-model="renewForm.duration" style="width: 160px">
                <a-option :value="1">1 年</a-option>
                <a-option :value="2">2 年</a-option>
                <a-option :value="3">3 年</a-option>
              </a-select>
            </a-form-item>
            <a-form-item label="费用">
              <span class="renew-price">¥{{ renewPrice }}</span>
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="submitRenew">立即续费</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <!-- 系统功能开关 -->
        <a-tab-pane key="features" title="功能开关">
          <a-form :model="featureForm" layout="vertical" style="max-width: 700px">
            <a-alert type="warning" style="margin-bottom: 16px">
              功能开关由套餐决定，部分功能需要升级套餐后才能启用。
            </a-alert>

            <a-form-item label="OTA 固件升级">
              <a-switch v-model="featureForm.enable_ota" :disabled="!infoForm.plan_features?.ota" />
              <span style="margin-left: 12px; color: #86909c">
                {{ infoForm.plan_features?.ota ? '已包含在当前套餐' : '当前套餐不支持' }}
              </span>
            </a-form-item>

            <a-form-item label="远程设备控制">
              <a-switch v-model="featureForm.enable_remote_control" :disabled="!infoForm.plan_features?.remote_control" />
              <span style="margin-left: 12px; color: #86909c">
                {{ infoForm.plan_features?.remote_control ? '已包含在当前套餐' : '当前套餐不支持' }}
              </span>
            </a-form-item>

            <a-form-item label="电子围栏">
              <a-switch v-model="featureForm.enable_geofence" :disabled="!infoForm.plan_features?.geofence" />
              <span style="margin-left: 12px; color: #86909c">
                {{ infoForm.plan_features?.geofence ? '已包含在当前套餐' : '当前套餐不支持' }}
              </span>
            </a-form-item>

            <a-form-item label="数据导出">
              <a-switch v-model="featureForm.enable_data_export" :disabled="!infoForm.plan_features?.data_export" />
              <span style="margin-left: 12px; color: #86909c">
                {{ infoForm.plan_features?.data_export ? '已包含在当前套餐' : '当前套餐不支持' }}
              </span>
            </a-form-item>

            <a-form-item label="API 开放接口">
              <a-switch v-model="featureForm.enable_api" :disabled="!infoForm.plan_features?.api" />
              <span style="margin-left: 12px; color: #86909c">
                {{ infoForm.plan_features?.api ? '已包含在当前套餐' : '当前套餐不支持' }}
              </span>
            </a-form-item>

            <a-divider />
            <div class="form-actions">
              <a-button type="primary" @click="submitFeatures">「保存」</a-button>
              <a-button @click="resetFeatures">「取消」</a-button>
            </div>
          </a-form>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <!-- 套餐变更确认 -->
    <a-modal
      v-model:visible="planChangeVisible"
      title="套餐变更确认"
      @before-ok="submitPlanChange"
      @cancel="planChangeVisible = false"
    >
      <a-result v-if="targetPlan" status="info" :title="`确定升级至 ${targetPlan.plan_name} 吗？`">
        <template #subtitle>
          <p>年费：<strong>¥{{ targetPlan.price_yearly }}</strong></p>
          <p>设备上限：{{ targetPlan.device_quota === -1 ? '不限' : targetPlan.device_quota + '台' }}</p>
          <p>用户上限：{{ targetPlan.user_quota === -1 ? '不限' : targetPlan.user_quota + '人' }}</p>
        </template>
      </a-result>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'

const infoFormRef = ref()
const planChangeVisible = ref(false)
const targetPlan = ref<any>(null)
const searchKeyword = ref('')

const infoForm = ref({
  id: 1,
  tenant_name: '深圳市智联科技有限公司',
  contact_name: '张经理',
  contact_phone: '13800138000',
  contact_email: 'zhang@zhilian.com',
  custom_domain: 'zhilian.mdm.com',
  tenant_logo: '',
  plan_name: '企业版',
  plan_id: 4,
  plan_code: 'enterprise',
  status: 'active',
  expire_time: '2026-06-20',
  plan_features: {
    ota: true,
    remote_control: true,
    geofence: false,
    data_export: true,
    api: false,
  },
})

const infoFormRules = {
  tenant_name: [{ required: true, message: '请输入租户名称' }],
  contact_name: [{ required: true, message: '请输入联系人' }],
  contact_phone: [{ required: true, message: '请输入联系电话' }],
}

const renewForm = ref({ duration: 1 })

const featureForm = ref({
  enable_ota: true,
  enable_remote_control: true,
  enable_geofence: false,
  enable_data_export: true,
  enable_api: false,
})

const plans = [
  {
    id: 1,
    plan_name: '免费版',
    price_yearly: 0,
    period: '年',
    device_quota: 10,
    user_quota: 5,
    dept_quota: 1,
    store_quota: 1,
    feature_list: ['基础设备管理'],
  },
  {
    id: 2,
    plan_name: '基础版',
    price_yearly: 299,
    period: '年',
    device_quota: 50,
    user_quota: 10,
    dept_quota: 5,
    store_quota: 10,
    feature_list: ['OTA 固件升级', '远程控制'],
  },
  {
    id: 3,
    plan_name: '专业版',
    price_yearly: 799,
    period: '年',
    device_quota: 200,
    user_quota: 20,
    dept_quota: 20,
    store_quota: 50,
    feature_list: ['OTA 固件升级', '远程控制', '数据导出', '自定义Logo/域名'],
  },
  {
    id: 4,
    plan_name: '企业版',
    price_yearly: 1999,
    period: '年',
    device_quota: 500,
    user_quota: 50,
    dept_quota: -1,
    store_quota: -1,
    feature_list: ['OTA 固件升级', '远程控制', '数据导出', '自定义Logo/域名', 'API开放接口', '专属客服'],
  },
]

const planOrder = [1, 2, 3, 4]

const isUpgradable = (plan: any) => planOrder.indexOf(plan.id) > planOrder.indexOf(infoForm.value.plan_id)

const renewPrice = computed(() => {
  const basePlan = plans.find(p => p.id === infoForm.value.plan_id)
  return basePlan ? basePlan.price_yearly * renewForm.value.duration : 0
})

const getStatusColor = (status: string) => {
  const map: Record<string, string> = { active: 'green', suspended: 'yellow', expired: 'red', pending: 'orange' }
  return map[status] || 'gray'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = { active: '正常', suspended: '已禁用', expired: '已到期', pending: '待审核' }
  return map[status] || status
}

const loadSettings = async () => {
  try {
    // API: GET /api/v1/admin/tenants/{id}/settings
    // const res = await axios.get(`/api/v1/admin/tenants/${infoForm.value.id}/settings`)
    // infoForm.value = { ...infoForm.value, ...res.data.data }
    // featureForm: sync from infoForm
    featureForm.value.enable_ota = infoForm.value.plan_features?.ota ?? true
    featureForm.value.enable_remote_control = infoForm.value.plan_features?.remote_control ?? true
    featureForm.value.enable_geofence = infoForm.value.plan_features?.geofence ?? false
    featureForm.value.enable_data_export = infoForm.value.plan_features?.data_export ?? true
    featureForm.value.enable_api = infoForm.value.plan_features?.api ?? false
  } catch {
    Message.error('加载配置失败')
  }
}

const handleSearch = (value: string) => {
  // TODO: 实现租户设置搜索过滤
  console.log('search:', value)
}

const submitInfo = async () => {
  try {
    await infoFormRef.value?.validate()
    // API: PUT /api/v1/admin/tenants/{id}/settings
    // await axios.put(`/api/v1/admin/tenants/${infoForm.value.id}/settings`, infoForm.value)
    Message.success('配置已保存')
  } catch {
    // validation failed
  }
}

const resetInfo = () => {
  loadSettings()
  Message.info('已重置')
}

const openPlanChange = (plan: any) => {
  targetPlan.value = plan
  planChangeVisible.value = true
}

const submitPlanChange = async (done: (val: boolean) => void) => {
  if (!targetPlan.value) { done(false); return }
  try {
    // API: POST /api/v1/admin/tenants/{id}/change-plan
    // await axios.post(`/api/v1/admin/tenants/${infoForm.value.id}/change-plan`, {
    //   plan_id: targetPlan.value.id,
    //   effective_type: 'immediate'
    // })
    infoForm.value.plan_id = targetPlan.value.id
    infoForm.value.plan_name = targetPlan.value.plan_name
    Message.success(`已升级至 ${targetPlan.value.plan_name}`)
    planChangeVisible.value = false
    done(true)
  } catch {
    Message.error('套餐变更失败')
    done(false)
  }
}

const submitRenew = async () => {
  try {
    // API: PUT /api/v1/admin/tenants/{id}/extend
    // await axios.put(`/api/v1/admin/tenants/${infoForm.value.id}/extend`, {
    //   extend_days: renewForm.value.duration * 365
    // })
    Message.success(`已提交续费 ${renewForm.value.duration} 年的申请`)
  } catch {
    Message.error('续费失败')
  }
}

const submitFeatures = async () => {
  try {
    // API: PUT /api/v1/admin/tenants/{id}/settings/features
    // await axios.put(`/api/v1/admin/tenants/${infoForm.value.id}/settings/features`, featureForm.value)
    Message.success('功能开关已保存')
  } catch {
    Message.error('保存失败')
  }
}

const resetFeatures = () => {
  loadSettings()
  Message.info('已重置')
}

onMounted(() => {
  loadSettings()
})
</script>

<style scoped>
.page-container {
  padding: 16px;
}
.form-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}
.plan-cards {
  margin-bottom: 24px;
}
.plan-cards .plan-card-active {
  border-color: #1650ff;
  background: #f2f3f5;
}
.plan-header h3 {
  text-align: center;
  margin: 0;
}
.plan-price {
  text-align: center;
  margin: 8px 0;
}
.price {
  font-size: 28px;
  font-weight: 700;
  color: #1650ff;
}
.period {
  color: #86909c;
  font-size: 14px;
}
.plan-features {
  list-style: none;
  padding: 0;
  margin: 0 0 16px;
}
.plan-features li {
  padding: 4px 0;
  color: #4e5969;
  font-size: 13px;
}
.plan-features li::before {
  content: '✓ ';
  color: #00b42a;
}
.renew-price {
  font-size: 20px;
  font-weight: 700;
  color: #1650ff;
}
</style>
