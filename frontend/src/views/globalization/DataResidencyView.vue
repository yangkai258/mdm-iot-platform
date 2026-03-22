<template>
  <div class="page-container">
    <!-- 顶部：标题 + 添加按钮 -->
    <a-card class="header-card">
      <div class="header-row">
        <div class="page-title">数据驻留规则</div>
        <a-button type="primary" @click="openAddModal">添加规则</a-button>
      </div>
    </a-card>

    <!-- 筛选栏 -->
    <a-card class="filter-card">
      <div class="filter-row">
        <a-select
          v-model="filters.dataType"
          placeholder="数据类型"
          style="width: 160px"
          allow-clear
          @change="handleFilterChange"
        >
          <a-option v-for="t in dataTypes" :key="t.value" :value="t.value">{{ t.label }}</a-option>
        </a-select>
        <a-select
          v-model="filters.status"
          placeholder="状态"
          style="width: 140px"
          allow-clear
          @change="handleFilterChange"
        >
          <a-option value="active">生效</a-option>
          <a-option value="pending">待审批</a-option>
          <a-option value="inactive">停用</a-option>
        </a-select>
        <a-select
          v-model="filters.tenantId"
          placeholder="租户"
          style="width: 160px"
          allow-clear
          allow-search
          @change="handleFilterChange"
        >
          <a-option v-for="t in tenants" :key="t.id" :value="t.id">{{ t.name }}</a-option>
        </a-select>
      </div>
    </a-card>

    <!-- 列表 -->
    <a-card class="table-card">
      <a-table
        :columns="columns"
        :data="filteredRules"
        :loading="loading"
        :pagination="{ pageSize: 20 }"
        row-key="id"
      >
        <template #dataType="{ record }">
          <a-tag :color="dataTypeColor(record.data_type)">{{ dataTypeLabel(record.data_type) }}</a-tag>
        </template>
        <template #storageRegion="{ record }">
          <span>{{ record.target_region }}</span>
        </template>
        <template #status="{ record }">
          <a-badge :color="statusColor(record.is_active ? 'active' : 'inactive')" :text="record.is_active ? '生效' : '停用'" />
        </template>
        <template #description="{ record }">
          <span class="desc-text">{{ record.description || '-' }}</span>
        </template>
        <template #actions="{ record }">
          <div class="action-btns">
            <a-button type="text" size="small" @click="openEditModal(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </div>
        </template>
      </a-table>
    </a-card>

    <!-- 提示 -->
    <a-card class="tip-card">
      <div class="tip-content">
        <icon-info-circle class="tip-icon" />
        <span>提示：数据驻留规则变更需要重新部署相关服务</span>
      </div>
    </a-card>

    <!-- 添加/编辑弹窗 -->
    <a-modal
      v-model:visible="formVisible"
      :title="editingRule ? '编辑规则' : '添加规则'"
      :width="520"
      @before-ok="handleSubmit"
      @cancel="formVisible = false"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="规则名称" required>
          <a-input v-model="form.rule_name" placeholder="如：用户数据-中国东部" />
        </a-form-item>
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
        <a-form-item label="存储 Schema">
          <a-input v-model="form.storage_schema" placeholder="如：public_cn_east" />
        </a-form-item>
        <a-form-item label="数据保留天数">
          <a-input-number v-model="form.retention_days" placeholder="如：365" :min="1" style="width: 100%" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" placeholder="请输入描述" :rows="3" />
        </a-form-item>
        <a-form-item label="启用状态">
          <a-switch v-model="form.is_active" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import { getDataResidencyRules, createDataResidencyRule, updateDataResidencyRule, deleteDataResidencyRule, getRegions } from '@/api/globalization'
import { DATA_TYPES } from '@/composables/useGlobalization'

const dataTypes = DATA_TYPES
const loading = ref(false)
const rules = ref([])
const regions = ref([])
const tenants = ref([])
const filters = reactive({ dataType: '', status: '', tenantId: '' })

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

function handleFilterChange() {}

onMounted(async () => {
  loading.value = true
  try {
    const [rulesRes, regionsRes] = await Promise.all([
      getDataResidencyRules(),
      getRegions()
    ])
    rules.value = rulesRes.data || rulesRes || []
    regions.value = regionsRes.data || regionsRes || []
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

async function handleSubmit(done) {
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
    done(false)
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
.page-container {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  height: 100%;
  box-sizing: border-box;
}

.header-card { flex-shrink: 0; }

.header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.page-title { font-size: 16px; font-weight: 600; }

.filter-card { flex-shrink: 0; }

.filter-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.table-card { flex: 1; overflow: auto; }

.desc-text { color: var(--color-text-3); font-size: 13px; }

.action-btns { display: flex; gap: 4px; }

.tip-card { flex-shrink: 0; background: #f9f0ff; border-color: #722ed1; }

.tip-content {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #722ed1;
  font-size: 13px;
}

.tip-icon { font-size: 16px; }
</style>
