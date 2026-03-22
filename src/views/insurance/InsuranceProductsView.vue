<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>保险服务</a-breadcrumb-item>
      <a-breadcrumb-item>保险产品</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="在售产品" :value="stats.active">
            <template #prefix>📦</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="本月新增保单" :value="stats.new_this_month">
            <template #prefix>🆕</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="累计保费" :value="stats.total_premium" prefix="¥" :precision="0" :value-style="{ color: '#0fc6c2' }">
            <template #prefix>💰</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="理赔完成率" :value="stats.claim_rate" suffix="%" :precision="1" :value-style="{ color: '#165dff' }">
            <template #prefix>🏥</template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- 筛选区 -->
    <a-card class="section-card" style="margin-top: 16px">
      <a-space wrap>
        <a-select v-model="filterForm.category" placeholder="产品类别" style="width: 160px" allow-clear @change="loadProducts">
          <a-option value="device">设备保障</a-option>
          <a-option value="health">健康医疗</a-option>
          <a-option value="accident">意外险</a-option>
          <a-option value="property">财产险</a-option>
        </a-select>
        <a-select v-model="filterForm.status" placeholder="产品状态" style="width: 140px" allow-clear @change="loadProducts">
          <a-option value="active">在售</a-option>
          <a-option value="inactive">停售</a-option>
        </a-select>
        <a-input-search v-model="filterForm.keyword" placeholder="搜索产品名称/编码" style="width: 220px" search-button @search="loadProducts" @change="e => !e.target.value && loadProducts()" />
        <a-button @click="loadProducts">
          <template #icon><icon-refresh :spin="loading" /></template>
          刷新
        </a-button>
        <a-button type="primary" @click="showCreateModal">
          <template #icon><icon-plus /></template>
          新增产品
        </a-button>
      </a-space>
    </a-card>

    <!-- 产品列表 -->
    <a-card class="section-card" style="margin-top: 16px">
      <a-spin :loading="loading" tip="加载中...">
        <div class="product-grid">
          <div v-for="p in products" :key="p.product_id" class="product-card">
            <div class="product-header">
              <span class="product-icon">{{ p.icon }}</span>
              <div class="product-title-area">
                <div class="product-name">{{ p.name }}</div>
                <div class="product-code">{{ p.product_id }}</div>
              </div>
              <a-tag :color="p.status === 'active' ? 'green' : 'gray'" size="small">{{ p.status === 'active' ? '在售' : '停售' }}</a-tag>
            </div>
            <div class="product-desc">{{ p.description }}</div>
            <div class="product-coverage">
              <div class="coverage-title">保障范围</div>
              <a-tag v-for="c in p.coverages" :key="c" size="small">{{ c }}</a-tag>
            </div>
            <div class="product-price">
              <span class="price-label">月费：</span>
              <span class="price-value">¥{{ p.monthly_premium }}</span>
              <span class="price-unit">/月</span>
            </div>
            <div class="product-claims">
              <span>已赔付 <strong>{{ p.claimed_count }}</strong> 次</span>
              <span>赔付率 <strong>{{ p.claim_rate }}%</strong></span>
            </div>
            <div class="product-actions">
              <a-button type="primary" size="small" @click="handleSubscribe(p)">立即投保</a-button>
              <a-button size="small" @click="viewDetail(p)">详情</a-button>
              <a-button type="text" size="small" status="warning" @click="showEditModal(p)">编辑</a-button>
            </div>
          </div>
        </div>
      </a-spin>
    </a-card>

    <!-- 新增/编辑产品弹窗 -->
    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑产品' : '新增产品'" :width="560" @before-ok="handleSave" @cancel="formVisible = false">
      <a-form :model="productForm" layout="vertical">
        <a-form-item label="产品名称" field="name">
          <a-input v-model="productForm.name" placeholder="请输入产品名称" />
        </a-form-item>
        <a-form-item label="产品类别" field="category">
          <a-select v-model="productForm.category" placeholder="请选择类别">
            <a-option value="device">设备保障</a-option>
            <a-option value="health">健康医疗</a-option>
            <a-option value="accident">意外险</a-option>
            <a-option value="property">财产险</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="月保费（元）" field="monthly_premium">
          <a-input-number v-model="productForm.monthly_premium" :min="0" :precision="2" style="width: 100%" />
        </a-form-item>
        <a-form-item label="保障额度（元）" field="coverage_amount">
          <a-input-number v-model="productForm.coverage_amount" :min="0" :precision="0" style="width: 100%" />
        </a-form-item>
        <a-form-item label="产品描述" field="description">
          <a-textarea v-model="productForm.description" placeholder="请输入产品描述" :rows="3" />
        </a-form-item>
        <a-form-item label="保障范围" field="coverages">
          <a-input v-model="productForm.coverages" placeholder="多个用逗号分隔" />
        </a-form-item>
        <a-form-item label="产品状态" field="status">
          <a-switch v-model="productForm.status" checked-value="active" unchecked-value="inactive" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 产品详情弹窗 -->
    <a-drawer v-model:visible="detailVisible" :title="currentProduct?.name" :width="540" unmountOnHide>
      <a-descriptions :column="2" bordered size="large" style="margin-bottom: 16px">
        <a-descriptions-item label="产品ID">{{ currentProduct?.product_id }}</a-descriptions-item>
        <a-descriptions-item label="产品状态">
          <a-tag :color="currentProduct?.status === 'active' ? 'green' : 'gray'">{{ currentProduct?.status === 'active' ? '在售' : '停售' }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="产品类别">{{ currentProduct?.category_text }}</a-descriptions-item>
        <a-descriptions-item label="月保费">¥{{ currentProduct?.monthly_premium }}</a-descriptions-item>
        <a-descriptions-item label="保障额度" :span="2">¥{{ currentProduct?.coverage_amount }}</a-descriptions-item>
        <a-descriptions-item label="已赔付次数">{{ currentProduct?.claimed_count }}</a-descriptions-item>
        <a-descriptions-item label="赔付率">{{ currentProduct?.claim_rate }}%</a-descriptions-item>
        <a-descriptions-item label="产品描述" :span="2">{{ currentProduct?.description }}</a-descriptions-item>
      </a-descriptions>
      <a-card title="保障范围明细">
        <a-list size="small">
          <a-list-item v-for="c in currentProduct?.coverages" :key="c">{{ c }}</a-list-item>
        </a-list>
      </a-card>
    </a-drawer>

    <!-- 投保确认弹窗 -->
    <a-modal v-model:visible="subscribeVisible" title="确认投保信息" @before-ok="handleSubscribeConfirm" @cancel="subscribeVisible = false">
      <a-descriptions :column="1" bordered>
        <a-descriptions-item label="产品名称">{{ currentProduct?.name }}</a-descriptions-item>
        <a-descriptions-item label="月保费">¥{{ currentProduct?.monthly_premium }} / 月</a-descriptions-item>
        <a-descriptions-item label="保障额度">¥{{ currentProduct?.coverage_amount }}</a-descriptions-item>
        <a-descriptions-item label="保障范围">
          <a-tag v-for="c in currentProduct?.coverages" :key="c" size="small" style="margin: 2px">{{ c }}</a-tag>
        </a-descriptions-item>
      </a-descriptions>
      <a-alert style="margin-top: 16px">投保成功后次日凌晨生效，请确认以上信息无误。</a-alert>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const formVisible = ref(false)
const detailVisible = ref(false)
const subscribeVisible = ref(false)
const isEdit = ref(false)
const currentProduct = ref(null)

const stats = reactive({ active: 0, new_this_month: 0, total_premium: 0, claim_rate: 0 })

const filterForm = reactive({ category: '', status: '', keyword: '' })

const productForm = reactive({
  name: '', category: '', monthly_premium: 0, coverage_amount: 0,
  description: '', coverages: '', status: 'active'
})

const products = ref([])

const loadProducts = async () => {
  loading.value = true
  await new Promise(r => setTimeout(r, 600))
  products.value = [
    { product_id: 'INS-DEV-001', name: '设备意外险', category: 'device', category_text: '设备保障', icon: '📱', description: '覆盖智能设备意外损坏、进水、摔碎等风险', monthly_premium: 29.9, coverage_amount: 3000, coverages: ['意外损坏', '屏幕碎裂', '进水保障', '免费维修'], status: 'active', claimed_count: 128, claim_rate: 8.5 },
    { product_id: 'INS-DEV-002', name: '设备延保险', category: 'device', category_text: '设备保障', icon: '🔧', description: '延长官方保修期，涵盖主板、电池等核心部件', monthly_premium: 19.9, coverage_amount: 2000, coverages: ['主板保修', '电池老化', '免费更换'], status: 'active', claimed_count: 65, claim_rate: 5.2 },
    { product_id: 'INS-HEALTH-001', name: '健康医疗险', category: 'health', category_text: '健康医疗', icon: '🏥', description: '门诊、住院全覆盖，最高可报销80%医疗费用', monthly_premium: 99.0, coverage_amount: 100000, coverages: ['门诊报销', '住院津贴', '特效药报销', '绿色通道'], status: 'active', claimed_count: 342, claim_rate: 12.3 },
    { product_id: 'INS-ACC-001', name: '综合意外险', category: 'accident', category_text: '意外险', icon: '🛡️', description: '涵盖意外身故、伤残、医疗，最高赔付50万', monthly_premium: 49.0, coverage_amount: 500000, coverages: ['意外身故', '意外伤残', '意外医疗', '住院津贴'], status: 'active', claimed_count: 89, claim_rate: 6.8 },
    { product_id: 'INS-PROP-001', name: '财产综合险', category: 'property', category_text: '财产险', icon: '🏠', description: '家庭财产被盗、火灾、水渍等综合保障', monthly_premium: 39.9, coverage_amount: 200000, coverages: ['火灾保障', '盗窃赔付', '水渍损失', '地震保障'], status: 'active', claimed_count: 45, claim_rate: 4.1 },
    { product_id: 'INS-DEV-003', name: '设备盗窃险', category: 'device', category_text: '设备保障', icon: '🔒', description: '设备被盗或抢劫，一键报案快速赔付', monthly_premium: 15.0, coverage_amount: 5000, coverages: ['被盗赔付', '抢劫保障', '定位协助'], status: 'inactive', claimed_count: 12, claim_rate: 3.2 }
  ]
  stats.active = products.value.filter(p => p.status === 'active').length
  stats.new_this_month = 23
  stats.total_premium = 128950
  stats.claim_rate = 87.5
  loading.value = false
}

const showCreateModal = () => {
  isEdit.value = false
  Object.assign(productForm, { name: '', category: '', monthly_premium: 0, coverage_amount: 0, description: '', coverages: '', status: 'active' })
  formVisible.value = true
}

const showEditModal = (p) => {
  isEdit.value = true
  currentProduct.value = p
  Object.assign(productForm, { name: p.name, category: p.category, monthly_premium: p.monthly_premium, coverage_amount: p.coverage_amount, description: p.description, coverages: p.coverages.join(','), status: p.status })
  formVisible.value = true
}

const handleSave = async (done) => {
  if (!productForm.name || !productForm.category) {
    Message.error('请填写必填字段')
    done(false)
    return
  }
  await new Promise(r => setTimeout(r, 500))
  Message.success(isEdit.value ? '产品更新成功' : '产品创建成功')
  formVisible.value = false
  loadProducts()
  done(true)
}

const viewDetail = (p) => {
  currentProduct.value = p
  detailVisible.value = true
}

const handleSubscribe = (p) => {
  currentProduct.value = p
  subscribeVisible.value = true
}

const handleSubscribeConfirm = async (done) => {
  await new Promise(r => setTimeout(r, 800))
  Message.success('投保成功！您的保单已生效')
  subscribeVisible.value = false
  done(true)
}

loadProducts()
</script>

<style scoped>
.page-container { padding: 20px; }
.breadcrumb { margin-bottom: 12px; }
.stat-card { text-align: center; }
.section-card { margin-bottom: 16px; }
.product-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(340px, 1fr)); gap: 16px; }
.product-card { border: 1px solid #e5e6e8; border-radius: 8px; padding: 16px; background: #fff; transition: box-shadow 0.2s; }
.product-card:hover { box-shadow: 0 4px 16px rgba(0,0,0,0.08); }
.product-header { display: flex; align-items: center; gap: 12px; margin-bottom: 12px; }
.product-icon { font-size: 32px; }
.product-title-area { flex: 1; }
.product-name { font-weight: 600; font-size: 16px; color: #1d2129; }
.product-code { font-size: 12px; color: #86909c; }
.product-desc { font-size: 13px; color: #4e5969; margin-bottom: 12px; line-height: 1.5; }
.product-coverage { margin-bottom: 12px; }
.coverage-title { font-size: 12px; color: #86909c; margin-bottom: 6px; }
.product-price { display: flex; align-items: baseline; margin-bottom: 8px; }
.price-label { font-size: 13px; color: #86909c; }
.price-value { font-size: 24px; font-weight: 700; color: #f53f3f; }
.price-unit { font-size: 12px; color: #86909c; margin-left: 2px; }
.product-claims { display: flex; gap: 16px; font-size: 12px; color: #86909c; margin-bottom: 12px; }
.product-claims strong { color: #1d2129; }
.product-actions { display: flex; gap: 8px; }
</style>
