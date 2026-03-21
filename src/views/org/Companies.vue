<template>
  <div class="company-list-page">
    <!-- 面包屑 -->
    <div class="breadcrumb-wrapper">
      <a-breadcrumb>
        <a-breadcrumb-item><a href="#/dashboard">首页</a></a-breadcrumb-item>
        <a-breadcrumb-item>组织管理</a-breadcrumb-item>
        <a-breadcrumb-item>公司管理</a-breadcrumb-item>
      </a-breadcrumb>
    </div>

    <!-- 操作区 -->
    <div class="toolbar">
      <div class="toolbar-left">
        <a-input-search 
          v-model="searchKey" 
          placeholder="搜索公司名称/编码" 
          style="width: 260px" 
          @search="loadCompanies"
        />
      </div>
      <div class="toolbar-right">
        <a-button type="primary" @click="openCreateModal(null)">「新建」</a-button>
        <a-button @click="showFilter = !showFilter">「筛选」</a-button>
        <a-button @click="loadCompanies">「刷新」</a-button>
      </div>
    </div>

    <!-- 筛选面板 -->
    <div v-if="showFilter" class="filter-panel">
      <a-card :bordered="false" size="small">
        <div class="filter-row">
          <span>状态：</span>
          <a-select v-model="filterStatus" placeholder="请选择" style="width: 120px" allow-clear>
            <a-option :value="1">正常</a-option>
            <a-option :value="0">禁用</a-option>
          </a-select>
          <a-button type="primary" size="small" @click="loadCompanies">查询</a-button>
          <a-button size="small" @click="resetFilter">重置</a-button>
        </div>
      </a-card>
    </div>

    <!-- 数据表格 -->
    <a-card :bordered="false" class="table-card">
      <a-table
        :columns="columns"
        :data="filteredData"
        :loading="loading"
        :pagination="{ pageSize: 10, showTotal: true }"
        row-key="id"
      >
        <template #companyCode="{ record }">
          <a-link @click="openCreateModal(record)">{{ record.company_code }}</a-link>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'gray'">
            {{ record.status === 1 ? '正常' : '禁用' }}
          </a-tag>
        </template>
        <template #createdAt="{ record }">
          {{ formatDate(record.created_at) }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openCreateModal(record)">「编辑」</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">「删除」</a-button>
            <a-button type="text" size="small" @click="goDepartments(record)">「部门」</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 创建/编辑弹窗 -->
    <a-modal
      v-model:visible="formVisible"
      :title="isEdit ? '编辑公司' : '创建公司'"
      :width="640"
      :mask-closable="false"
      @before-ok="submitForm"
      @cancel="formVisible = false"
    >
      <a-form ref="formRef" :model="formData" :rules="formRules" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="公司编码" field="company_code">
              <a-input v-model="formData.company_code" placeholder="请输入" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="公司名称" field="company_name">
              <a-input v-model="formData.company_name" placeholder="请输入" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="简称" field="short_name">
              <a-input v-model="formData.short_name" placeholder="请输入" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="法人代表" field="legal_person">
              <a-input v-model="formData.legal_person" placeholder="请输入" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="联系人" field="contact">
              <a-input v-model="formData.contact" placeholder="请输入" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="联系电话" field="phone">
              <a-input v-model="formData.phone" placeholder="请输入" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item label="地址">
              <a-space>
                <a-input v-model="formData.province" placeholder="省" style="width: 100px" />
                <a-input v-model="formData.city" placeholder="市" style="width: 100px" />
                <a-input v-model="formData.district" placeholder="区" style="width: 100px" />
                <a-input v-model="formData.address" placeholder="详细地址" style="width: 200px" />
              </a-space>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="状态" field="status">
              <a-radio-group v-model="formData.status">
                <a-radio :value="1">正常</a-radio>
                <a-radio :value="0">禁用</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>

    <!-- 删除确认 -->
    <a-modal
      v-model:visible="deleteVisible"
      title="删除确认"
      @before-ok="submitDelete"
    >
      <a-result status="warning" title="确定要删除该公司吗？">
        <template #subtitle>删除后无法恢复</template>
      </a-result>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'

const API_BASE = '/api/v1'
const searchKey = ref('')
const showFilter = ref(false)
const filterStatus = ref<number | ''>('')
const loading = ref(false)
const formVisible = ref(false)
const deleteVisible = ref(false)
const isEdit = ref(false)
const currentRecord = ref<any>(null)
const currentDeleteRecord = ref<any>(null)
const formRef = ref()

const formData = ref({
  company_code: '',
  company_name: '',
  short_name: '',
  legal_person: '',
  contact: '',
  phone: '',
  province: '',
  city: '',
  district: '',
  address: '',
  status: 1,
})

const formRules = {
  company_code: [{ required: true, message: '请输入公司编码' }],
  company_name: [{ required: true, message: '请输入公司名称' }],
}

const columns = [
  { title: '公司编码', slotName: 'companyCode', width: 130 },
  { title: '公司名称', dataIndex: 'company_name', width: 200 },
  { title: '简称', dataIndex: 'short_name', width: 100 },
  { title: '联系人', dataIndex: 'contact', width: 100 },
  { title: '联系电话', dataIndex: 'phone', width: 130 },
  { title: '省', dataIndex: 'province', width: 70 },
  { title: '市', dataIndex: 'city', width: 70 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '创建时间', slotName: 'createdAt', width: 110 },
  { title: '操作', slotName: 'actions', width: 160, fixed: 'right' },
]

const mockData = ref([
  { id: 1, company_code: 'C001', company_name: '深圳市智联科技有限公司', short_name: '智联科技', legal_person: '张总', contact: '张经理', phone: '13800138000', province: '广东省', city: '深圳市', district: '南山区', address: '科技园路1号', status: 1, created_at: '2026-01-15T08:00:00Z' },
  { id: 2, company_code: 'C002', company_name: '广州云智数据服务公司', short_name: '云智数据', legal_person: '李总', contact: '李总监', phone: '13900139000', province: '广东省', city: '广州市', district: '天河区', address: '天河路100号', status: 1, created_at: '2026-02-01T10:00:00Z' },
  { id: 3, company_code: 'C003', company_name: '杭州创新科技集团', short_name: '创新科技', legal_person: '王总', contact: '王总', phone: '13700137000', province: '浙江省', city: '杭州市', district: '西湖区', address: '文三路200号', status: 0, created_at: '2026-03-01T14:00:00Z' },
])

const filteredData = computed(() => {
  let data = mockData.value
  if (searchKey.value) {
    const kw = searchKey.value.toLowerCase()
    data = data.filter(item => 
      item.company_name.toLowerCase().includes(kw) || 
      item.company_code.toLowerCase().includes(kw)
    )
  }
  if (filterStatus.value !== '') {
    data = data.filter(item => item.status === filterStatus.value)
  }
  return data
})

const loadCompanies = async () => {
  loading.value = true
  await new Promise(r => setTimeout(r, 300))
  loading.value = false
}

const resetFilter = () => {
  filterStatus.value = ''
  loadCompanies()
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

const openCreateModal = (record: any) => {
  isEdit.value = !!record
  currentRecord.value = record
  if (record) {
    Object.assign(formData.value, record)
  } else {
    formData.value = { company_code: '', company_name: '', short_name: '', legal_person: '', contact: '', phone: '', province: '', city: '', district: '', address: '', status: 1 }
  }
  formRef.value?.clearValidate()
  formVisible.value = true
}

const submitForm = async (done: (val: boolean) => void) => {
  try {
    if (isEdit.value) {
      Object.assign(currentRecord.value, formData.value)
      Message.success('保存成功')
    } else {
      mockData.value.unshift({ ...formData.value, id: Date.now(), created_at: new Date().toISOString() })
      Message.success('创建成功')
    }
    formVisible.value = false
    done(true)
  } catch {
    Message.error('操作失败')
    done(false)
  }
}

const handleDelete = (record: any) => {
  currentDeleteRecord.value = record
  deleteVisible.value = true
}

const submitDelete = async (done: (val: boolean) => void) => {
  try {
    await axios.delete(`${API_BASE}/org/companies/${currentDeleteRecord.value.id}`)
    mockData.value = mockData.value.filter(item => item.id !== currentDeleteRecord.value.id)
    Message.success('删除成功')
    deleteVisible.value = false
    done(true)
  } catch (e: any) {
    Message.error(e?.response?.data?.message || '删除失败')
    done(false)
  }
}

const goDepartments = (record: any) => {
  window.location.hash = '#/org/departments?company_id=' + record.id
}

onMounted(() => loadCompanies())
</script>

<style>
/* 页面容器 */
.company-list-page {
  padding: 24px;
  min-height: 100%;
  background: #f2f3f5;
}

/* 面包屑 */
.breadcrumb-wrapper {
  margin-bottom: 16px;
}

/* 工具栏 */
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  background: #fff;
  padding: 16px;
  border-radius: 4px;
}

.toolbar-left {
  display: flex;
  align-items: center;
}

.toolbar-right {
  display: flex;
  gap: 8px;
}

/* 筛选面板 */
.filter-panel {
  margin-bottom: 16px;
}

.filter-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* 表格卡片 */
.table-card {
  background: #fff;
  border-radius: 4px;
}

/* 覆盖ArcoDesign默认样式 */
.arco-btn-primary {
  background-color: #165dff !important;
  border-color: #165dff !important;
}
</style>
