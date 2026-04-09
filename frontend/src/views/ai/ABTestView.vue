<template>
  <div class="page-container">
    <a-card class="general-card" title="A/B测试框架">
      <template #extra>
        <a-button type="primary" @click="openCreateModal"><icon-plus />新建实验</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="实验名称"><a-input v-model="form.experiment_name" placeholder="请输�? /></a-form-item>
          <a-form-item label="状�?>
            <a-select v-model="form.status" placeholder="选择状�? allow-clear style="width: 140px">
              <a-option value="draft">草稿</a-option>
              <a-option value="running">运行�?/a-option>
              <a-option value="paused">已暂�?/a-option>
              <a-option value="completed">已完�?/a-option>
            </a-select>
          </a-form-item>
          <a-form-item><a-button type="primary" @click="loadData">查询</a-button><a-button @click="Object.keys(form).forEach(k => form[k] = ''); loadData()">重置</a-button></a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
        <template #result="{ record }">
          <a-space direction="vertical" :size="0">
            <span>A�? <strong>{{ record.group_a_conversion || 0 }}%</strong></span>
            <span>B�? <strong>{{ record.group_b_conversion || 0 }}%</strong></span>
          </a-space>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="viewResult(record)">结果</a-button>
          <a-button v-if="record.status === 'draft'" type="text" size="small" @click="startExperiment(record)">启动</a-button>
          <a-button v-if="record.status === 'running'" type="text" size="small" @click="pauseExperiment(record)">暂停</a-button>
          <a-button type="text" size="small" @click="editExperiment(record)">编辑</a-button>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model:visible="createModalVisible" :title="isEdit ? '编辑实验' : '新建实验'" @before-ok="handleSubmit" :loading="submitting" :width="600">
      <a-form :model="expForm" layout="vertical">
        <a-form-item label="实验名称" required><a-input v-model="expForm.experiment_name" placeholder="请输入实验名�? /></a-form-item>
        <a-form-item label="实验描述"><a-textarea v-model="expForm.description" :rows="2" placeholder="实验描述" /></a-form-item>
        <a-form-item label="分流比例">
          <a-space>
            <span>A�?</span><a-input-number v-model="expForm.group_a_ratio" :min="0" :max="100" :step="5" style="width: 80px" /><span>%</span>
            <span style="margin-left: 16px">B�?</span><a-input-number v-model="expForm.group_b_ratio" :min="0" :max="100" :step="5" style="width: 80px" /><span>%</span>
          </a-space>
        </a-form-item>
        <a-form-item label="A组策�?>
          <a-select v-model="expForm.group_a_strategy" placeholder="选择A组策�?>
            <a-option value="control">对照组（原有策略�?/a-option>
            <a-option value="strategy_1">策略一</a-option>
            <a-option value="strategy_2">策略�?/a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="B组策�?>
          <a-select v-model="expForm.group_b_strategy" placeholder="选择B组策�?>
            <a-option value="control">对照组（原有策略�?/a-option>
            <a-option value="strategy_1">策略一</a-option>
            <a-option value="strategy_2">策略�?/a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="目标指标">
          <a-select v-model="expForm.metric" placeholder="选择目标指标">
            <a-option value="conversion">转化�?/a-option>
            <a-option value="retention">留存�?/a-option>
            <a-option value="engagement"> engagement">用户活跃</a-option>
            <a-option value="satisfaction">满意�?/a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="预期样本�?>
          <a-input-number v-model="expForm.expected_sample_size" :min="100" placeholder="每组预期样本�? style="width: 200px" />
        </a-form-item>
      </a-form>
    </a-modal>
    <a-modal v-model:visible="resultModalVisible" title="实验结果" :width="700" :footer="null">
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="实验名称" :span="2">{{ selectedExp?.experiment_name }}</a-descriptions-item>
        <a-descriptions-item label="实验状�?><a-tag :color="getStatusColor(selectedExp?.status)">{{ getStatusText(selectedExp?.status) }}</a-tag></a-descriptions-item>
        <a-descriptions-item label="目标指标">{{ selectedExp?.metric }}</a-descriptions-item>
        <a-descriptions-item label="总样本量">{{ selectedExp?.total_samples || 0 }}</a-descriptions-item>
        <a-descriptions-item label="A组样�?>{{ selectedExp?.group_a_samples || 0 }}</a-descriptions-item>
        <a-descriptions-item label="B组样�?>{{ selectedExp?.group_b_samples || 0 }}</a-descriptions-item>
        <a-descriptions-item label="A组转化率">{{ selectedExp?.group_a_conversion || 0 }}%</a-descriptions-item>
        <a-descriptions-item label="B组转化率">{{ selectedExp?.group_b_conversion || 0 }}%</a-descriptions-item>
        <a-descriptions-item label="统计显著�? :span="2">{{ selectedExp?.significance || 'N/A' }}</a-descriptions-item>
        <a-descriptions-item label="结论" :span="2">{{ selectedExp?.conclusion || '实验进行�? }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const submitting = ref(false)
const data = ref([])
const createModalVisible = ref(false)
const resultModalVisible = ref(false)
const isEdit = ref(false)
const selectedExp = ref(null)
const form = reactive({ experiment_name: '', status: '' })
const expForm = reactive({ id: null, experiment_name: '', description: '', group_a_ratio: 50, group_b_ratio: 50, group_a_strategy: 'control', group_b_strategy: 'strategy_1', metric: 'conversion', expected_sample_size: 1000 })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '实验名称', dataIndex: 'experiment_name', width: 180 },
  { title: '状�?, slotName: 'status', width: 90 },
  { title: 'A/B组转化率', slotName: 'result', width: 140 },
  { title: '总样�?, dataIndex: 'total_samples', width: 90 },
  { title: '开始时�?, dataIndex: 'started_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' }
]

const getStatusColor = (s) => ({ draft: 'gray', running: 'arcoblue', paused: 'orange', completed: 'green' }[s] || 'gray')
const getStatusText = (s) => ({ draft: '草稿', running: '运行�?, paused: '已暂�?, completed: '已完�? }[s] || s)

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (form.experiment_name) params.append('experiment_name', form.experiment_name)
    if (form.status) params.append('status', form.status)
    const res = await fetch(`/api/v1/ai/ab-test/experiments?${params}`, { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { data.value = res.data?.list || []; pagination.total = res.data?.total || 0 }
    else { data.value = [] }
  } catch (e) { Message.error('加载失败') } finally { loading.value = false }
}

const openCreateModal = () => { isEdit.value = false; Object.assign(expForm, { id: null, experiment_name: '', description: '', group_a_ratio: 50, group_b_ratio: 50, group_a_strategy: 'control', group_b_strategy: 'strategy_1', metric: 'conversion', expected_sample_size: 1000 }); createModalVisible.value = true }
const editExperiment = (record) => { isEdit.value = true; Object.assign(expForm, record); createModalVisible.value = true }

const handleSubmit = async (done) => {
  if (!expForm.experiment_name) { Message.warning('请输入实验名�?); done(false); return }
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const url = isEdit.value ? `/api/v1/ai/ab-test/experiments/${expForm.id}` : '/api/v1/ai/ab-test/experiments'
    const res = await fetch(url, { method: isEdit.value ? 'PUT' : 'POST', headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(expForm) }).then(r => r.json())
    if (res.code === 0) { Message.success(isEdit.value ? '更新成功' : '创建成功'); createModalVisible.value = false; loadData() }
    else { Message.error(res.message || '操作失败') }
    done(true)
  } catch (e) { Message.error('操作失败'); done(false) } finally { submitting.value = false }
}

const viewResult = (record) => { selectedExp.value = record; resultModalVisible.value = true }

const startExperiment = async (record) => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/ai/ab-test/experiments/${record.id}/start`, { method: 'POST', headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { Message.success('实验已启�?); loadData() }
    else Message.error('启动失败')
  } catch (e) { Message.error('启动失败') }
}

const pauseExperiment = async (record) => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/ai/ab-test/experiments/${record.id}/pause`, { method: 'POST', headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { Message.success('实验已暂�?); loadData() }
    else Message.error('暂停失败')
  } catch (e) { Message.error('暂停失败') }
}

const onPageChange = (page) => { pagination.current = page; loadData() }

onMounted(() => loadData())
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>