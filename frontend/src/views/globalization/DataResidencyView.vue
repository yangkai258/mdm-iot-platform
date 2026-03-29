<template>
  <div class="container">
    <Breadcrumb :items="['menu.globalization', 'menu.globalization.dataResidency']" />
    <a-card class="general-card" title="数据驻留配置">
      <template #extra>
        <a-space :size="12">
          <a-button type="primary" @click="openAddModal"><icon-plus />添加规则</a-button>
          <a-button @click="() => {}"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-row :gutter="16">
        <a-col :span="6">
          <a-form-item label="数据类型">
            <a-select v-model="filters.dataType" placeholder="请选择" allow-clear style="width: 100%">
              <a-option v-for="t in dataTypes" :key="t.value" :value="t.value">{{ t.label }}</a-option>
            </a-select>
          </a-form-item>
        </a-col>
        <a-col :span="6">
          <a-form-item label="状态">
            <a-select v-model="filters.status" placeholder="请选择" allow-clear style="width: 100%">
              <a-option value="active">生效</a-option>
              <a-option value="inactive">停用</a-option>
            </a-select>
          </a-form-item>
        </a-col>
        <a-col :flex="'86px'" style="display: flex; align-items: flex-end">
          <a-space direction="vertical" :size="8">
            <a-button type="primary" @click="handleFilterChange">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-space>
        </a-col>
      </a-row>
      <a-divider style="margin: 0 0 16px 0" />
      <a-table :columns="columns" :data="filteredRules" :loading="loading" :pagination="pagination" row-key="id">
        <template #dataType="{ record }">
          <a-tag :color="dataTypeColor(record.data_type)">{{ dataTypeLabel(record.data_type) }}</a-tag>
        </template>
        <template #storageRegion="{ record }"><span>{{ record.target_region }}</span></template>
        <template #status="{ record }">
          <a-badge :color="statusColor(record.is_active ? 'active' : 'inactive')" :text="record.is_active ? '生效' : '停用'" />
        </template>
        <template #description="{ record }"><span class="desc-text">{{ record.description || '-' }}</span></template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="openEditModal(record)">编辑</a-button>
          <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model:visible="formVisible" :title="editingRule ? '编辑规则' : '添加规则'" :width="520">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="规则名称" required><a-input v-model="form.rule_name" placeholder="如：用户数据-中国东部" /></a-form-item>
        <a-form-item label="数据类型" required>
          <a-select v-model="form.data_type" placeholder="选择数据类型">
            <a-option v-for="t in dataTypes" :key="t.value" :value="t.value">{{ t.label }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="存储区域" required>
          <a-select v-model="form.target_region" placeholder="选择存储区域">
            <a-option v-for="r in regions" :key="r.region_code" :value="r.region_code">{{ r.region_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="存储 Schema"><a-input v-model="form.storage_schema" placeholder="如：public_cn_east" /></a-form-item>
        <a-form-item label="保留天数"><a-input-number v-model="form.retention_days" :min="1" style="width: 100%" /></a-form-item>
        <a-form-item label="描述"><a-textarea v-model="form.description" :rows="3" /></a-form-item>
        <a-form-item label="启用状态"><a-switch v-model="form.is_active" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="formVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'
import { getDataResidencyRules, createDataResidencyRule, updateDataResidencyRule, deleteDataResidencyRule, getRegions } from '@/api/globalization'
import { DATA_TYPES } from '@/composables/useGlobalization'

const dataTypes = DATA_TYPES
const loading = ref(false)
const rules = ref([])
const regions = ref([])
const filters = reactive({ dataType: '', status: '', tenantId: '' })
const pagination = reactive({ pageSize: 20, current: 1, total: 0 })

const columns = [
  { title: '数据类型', slotName: 'dataType', width: 130 },
  { title: '存储区域', slotName: 'storageRegion', width: 130 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '描述', slotName: 'description', ellipsis: true },
  { title: '操作', slotName: 'actions', width: 120 }
]

const filteredRules = computed(() => {
  return rules.value.filter(r => {
    if (filters.dataType && r.data_type !== filters.dataType) return false
    if (filters.status) {
      const isActive = filters.status === 'active'
      if (r.is_active !== isActive) return false
    }
    if (filters.tenantId && r.tenant_id !== filters.tenantId) return false
    return true
  })
})

const formVisible = ref(false)
const editingRule = ref(null)
const form = reactive({
  rule_name: '',
  data_type: '',
  target_region: '',
  storage_schema: '',
  retention_days: null,
  description: '',
  is_active: true
})

function dataTypeColor(type) {
  const map = {
    user_data: 'arcoblue',
    device_data: 'green',
    ai_training_data: 'purple',
    ai_inference_data: 'orange',
    log_data: 'gray'
  }
  return map[type] || 'default'
}

function dataTypeLabel(type) {
  const map = {
    user_data: '用户数据',
    device_data: '设备数据',
    ai_training_data: 'AI训练数据',
    ai_inference_data: 'AI推理数据',
    log_data: '日志数据'
  }
  return map[type] || type
}

function statusColor(status) {
  return status === 'active' ? 'green' : 'gray'
}

function handleFilterChange() {
  pagination.current = 1
}

function handleReset() {
  filters.dataType = ''
  filters.status = ''
  pagination.current = 1
}

onMounted(async () => {
  loading.value = true
  try {
    const [rulesRes, regionsRes] = await Promise.all([
      getDataResidencyRules(),
      getRegions()
    ])
    rules.value = rulesRes.data || rulesRes || []
    regions.value = regionsRes.data || regionsRes || []
    pagination.total = filteredRules.value.length
  } catch (e) {
    console.error('加载数据驻留规则失败', e)
  } finally {
    loading.value = false
  }
})

function openAddModal() {
  editingRule.value = null
  Object.assign(form, { rule_name: '', data_type: '', target_region: '', storage_schema: '', retention_days: null, description: '', is_active: true })
  formVisible.value = true
}

function openEditModal(record) {
  editingRule.value = record
  Object.assign(form, {
    rule_name: record.rule_name,
    data_type: record.data_type,
    target_region: record.target_region,
    storage_schema: record.storage_schema || '',
    retention_days: record.retention_days,
    description: record.description || '',
    is_active: record.is_active
  })
  formVisible.value = true
}

async function handleSubmit() {
  try {
    if (editingRule.value) {
      await updateDataResidencyRule(editingRule.value.id, form)
      Message.success('更新成功')
    } else {
      await createDataResidencyRule(form)
      Message.success('创建成功')
    }
    formVisible.value = false
    const res = await getDataResidencyRules()
    rules.value = res.data || res || []
  } catch (e) {
    Message.error('操作失败')
  }
}

function handleDelete(record) {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除规则「${record.rule_name}」吗？`,
    okText: '删除',
    onOk: async () => {
      try {
        await deleteDataResidencyRule(record.id)
        Message.success('删除成功')
        const res = await getDataResidencyRules()
        rules.value = res.data || res || []
      } catch (e) {
        Message.error('删除失败')
      }
    }
  })
}
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
.desc-text { color: var(--color-text-3); font-size: 13px; }
</style>
