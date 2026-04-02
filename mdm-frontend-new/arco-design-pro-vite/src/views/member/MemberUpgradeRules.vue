<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="member-upgrade-rules-page">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员升级规则</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索规则名称" style="width: 240px" search-button @search="handleSearch" />
        <a-select v-model="filters.targetLevelId" placeholder="目标等级" allow-clear style="width: 140px" @change="handleSearch">
          <a-option v-for="lv in levelList" :key="lv.id" :value="lv.id">{{ lv.name }}</a-option>
        </a-select>
        <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 120px" @change="handleSearch">
          <a-option :value="1">启用</a-option>
          <a-option :value="0">禁用</a-option>
        </a-select>
        <a-button @click="handleSearch">「搜索」</a-button>
        <a-button @click="resetFilters">「重置」</a-button>
      </a-space>
    </a-card>

    <!-- 操作栏 -->
    <a-card class="table-card">
      <template #extra>
        <a-space>
          <a-button type="primary" @click="showCreate">「新建」</a-button>
          <a-button @click="handleExport">「导出」</a-button>
          <a-button @click="loadData">🔄</a-button>
        </a-space>
      </template>

      <a-table
        :columns="columns"
        :data="dataList"
        :loading="loading"
        :pagination="paginationConfig"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
        row-key="id"
        :scroll="{ x: 1100 }"
      >
        <template #targetLevel="{ record }">
          <a-tag :color="getLevelColor(record.targetLevelId)">{{ record.targetLevelName || '-' }}</a-tag>
        </template>
        <template #conditionType="{ record }">
          <a-tag>{{ conditionTypeMap[record.conditionType] || record.conditionType }}</a-tag>
        </template>
        <template #conditionValue="{ record }">
          <span>{{ formatConditionValue(record) }}</span>
        </template>
        <template #pointsBonus="{ record }">
          <span style="color: #ff6b00; font-weight: 600;">{{ record.pointsBonus || 0 }}分</span>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '启用' : '禁用' }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showDetail(record)">详情</a-button>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 详情抽屉 -->
    <a-drawer v-model:visible="detailVisible" title="升级规则详情" :width="520">
      <template v-if="currentRecord">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="规则名称">{{ currentRecord.name }}</a-descriptions-item>
          <a-descriptions-item label="目标等级"><a-tag :color="getLevelColor(currentRecord.targetLevelId)">{{ currentRecord.targetLevelName }}</a-tag></a-descriptions-item>
          <a-descriptions-item label="条件类型">{{ conditionTypeMap[currentRecord.conditionType] || currentRecord.conditionType }}</a-descriptions-item>
          <a-descriptions-item label="条件值">{{ formatConditionValue(currentRecord) }}</a-descriptions-item>
          <a-descriptions-item label="升级赠送积分">{{ currentRecord.pointsBonus || 0 }}分</a-descriptions-item>
          <a-descriptions-item label="有效期">{{ currentRecord.effectiveDate || '-' }} 至 {{ currentRecord.expireDate || '永久' }}</a-descriptions-item>
          <a-descriptions-item label="状态"><a-tag :color="currentRecord.status === 1 ? 'green' : 'gray'">{{ currentRecord.status === 1 ? '启用' : '禁用' }}</a-tag></a-descriptions-item>
          <a-descriptions-item label="描述">{{ currentRecord.description || '-' }}</a-descriptions-item>
        </a-descriptions>
      </template>
    </a-drawer>

    <!-- 新增/编辑弹窗 -->
    <a-modal
      v-model:visible="formVisible"
      :title="isEdit ? '编辑升级规则' : '新建升级规则'"
      @before-ok="handleFormSubmit"
      @cancel="formVisible = false"
      :width="560"
      :loading="formLoading"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="规则名称" field="name" :rules="[{ required: true, message: '请输入规则名称' }]">
          <a-input v-model="form.name" placeholder="如：消费满5000升级黄金会员" />
        </a-form-item>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="目标等级" field="targetLevelId" :rules="[{ required: true, message: '请选择目标等级' }]">
              <a-select v-model="form.targetLevelId" placeholder="选择目标等级">
                <a-option v-for="lv in levelList" :key="lv.id" :value="lv.id">{{ lv.name }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="条件类型" field="conditionType" :rules="[{ required: true, message: '请选择条件类型' }]">
              <a-select v-model="form.conditionType" placeholder="选择条件类型">
                <a-option value="consume">累计消费金额</a-option>
                <a-option value="points">累计积分达到</a-option>
                <a-option value="orders">累计订单数</a-option>
                <a-option value="days">注册天数</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="条件值" field="conditionValue" :rules="[{ required: true, message: '请输入条件值' }]">
              <a-input-number v-model="form.conditionValue" :min="0" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="升级赠送积分" field="pointsBonus">
              <a-input-number v-model="form.pointsBonus" :min="0" style="width: 100%">
                <template #suffix>分</template>
              </a-input-number>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="生效日期" field="effectiveDate">
              <a-date-picker v-model="form.effectiveDate" format="YYYY-MM-DD" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="失效日期" field="expireDate">
              <a-date-picker v-model="form.expireDate" format="YYYY-MM-DD" style="width: 100%" placeholder="不填表示永久有效" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="状态" field="status">
          <a-radio-group v-model="form.status">
            <a-radio :value="1">启用</a-radio>
            <a-radio :value="0">禁用</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" :rows="2" placeholder="简要描述该规则" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import * as api from '@/api/member'

const dataList = ref([])
const levelList = ref([])
const loading = ref(false)
const formLoading = ref(false)
const formVisible = ref(false)
const detailVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const currentRecord = ref(null)

const conditionTypeMap = { consume: '累计消费金额', points: '累计积分', orders: '累计订单数', days: '注册天数' }

const filters = reactive({ keyword: '', targetLevelId: undefined, status: undefined })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current, pageSize: pagination.pageSize, total: pagination.total,
  showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50, 100]
}))

const form = reactive({
  name: '', targetLevelId: undefined, conditionType: 'consume', conditionValue: 0,
  pointsBonus: 0, effectiveDate: '', expireDate: '', status: 1, description: ''
})

const columns = [
  { title: '规则名称', dataIndex: 'name', width: 200 },
  { title: '目标等级', slotName: 'targetLevel', width: 130 },
  { title: '条件类型', slotName: 'conditionType', width: 130 },
  { title: '条件值', slotName: 'conditionValue', width: 160 },
  { title: '升级积分', slotName: 'pointsBonus', width: 100 },
  { title: '有效期', dataIndex: 'effectiveDate', width: 170 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' }
]

const getLevelColor = (id) => ({ 1: '#95de64', 2: '#1890ff', 3: '#faad14', 4: '#ff4d4f' }[id] || 'gray')

const formatConditionValue = (record) => {
  const map = { consume: '¥', points: '积分', orders: '笔', days: '天' }
  const suffix = map[record.conditionType] || ''
  return `${record.conditionValue || 0}${suffix}`
}

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.targetLevelId) params.targetLevelId = filters.targetLevelId
    if (filters.status !== undefined) params.status = filters.status
    const res = await api.getUpgradeRules()
    const d = res.data || {}
    dataList.value = Array.isArray(d) ? d : (d.list || [])
    pagination.total = dataList.value.length
  } catch (err) { Message.error('加载失败: ' + err.message) }
  finally { loading.value = false }
}

const loadLevels = async () => {
  try {
    const res = await api.getLevelList()
    levelList.value = res.data || []
  } catch (err) { /* ignore */ }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const resetFilters = () => { filters.keyword = ''; filters.targetLevelId = undefined; filters.status = undefined; pagination.current = 1; loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }
const onPageSizeChange = (pageSize) => { pagination.pageSize = pageSize; pagination.current = 1; loadData() }

const showCreate = () => {
  isEdit.value = false
  Object.assign(form, { name: '', targetLevelId: undefined, conditionType: 'consume', conditionValue: 0, pointsBonus: 0, effectiveDate: '', expireDate: '', status: 1, description: '' })
  formVisible.value = true
}

const showEdit = (record) => {
  isEdit.value = true; currentId.value = record.id
  Object.assign(form, {
    name: record.name, targetLevelId: record.targetLevelId, conditionType: record.conditionType || 'consume',
    conditionValue: record.conditionValue || 0, pointsBonus: record.pointsBonus || 0,
    effectiveDate: record.effectiveDate || '', expireDate: record.expireDate || '',
    status: record.status ?? 1, description: record.description || ''
  })
  formVisible.value = true
}

const showDetail = (record) => { currentRecord.value = record; detailVisible.value = true }

const handleFormSubmit = async (done) => {
  if (!form.name || !form.targetLevelId) { Message.warning('请填写规则名称和目标等级'); done(false); return }
  formLoading.value = true
  try {
    const payload = { ...form }
    if (isEdit.value) { await api.updateUpgradeRules({ ...payload, id: currentId.value }); Message.success('更新成功') }
    else { await api.updateUpgradeRules(payload); Message.success('创建成功') }
    formVisible.value = false; loadData(); done(true)
  } catch (err) { Message.error(err.message || '操作失败'); done(false) }
  finally { formLoading.value = false }
}

const handleDelete = (record) => {
  Modal.warning({ title: '确认删除', content: `确定要删除规则「${record.name}」吗？`, okText: '确认删除', onOk: async () => {
    try { await api.updateUpgradeRules({ ...record, deleted: true }); Message.success('删除成功'); loadData() }
    catch (err) { Message.error(err.message || '删除失败') }
  }})
}

const handleExport = () => { Message.info('导出功能开发中') }

onMounted(() => { loadData(); loadLevels() })
</script>

<style scoped>
.member-upgrade-rules-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
.table-card { border-radius: 8px; }
</style>
