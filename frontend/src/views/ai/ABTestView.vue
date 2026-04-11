<template>
  <div class="page-container">
    <a-card class="general-card" title="A/B娴嬭瘯妗嗘灦">
      <template #extra>
        <a-button type="primary" @click="openCreateModal"><icon-plus />鏂板缓瀹為獙</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="瀹為獙鍚嶇О"><a-input v-model="form.experiment_name" placeholder="璇疯緭鍏" /></a-form-item>
          <a-form-item label="鐘舵€?>
            <a-select v-model="form.status" placeholder="閫夋嫨鐘舵€" allow-clear style="width: 140px">
              <a-option value="draft">鑽夌</a-option>
              <a-option value="running">杩愯涓?/a-option>
              <a-option value="paused">宸叉殏鍋?/a-option>
              <a-option value="completed">宸插畬鎴?/a-option>
            </a-select>
          </a-form-item>
          <a-form-item><a-button type="primary" @click="loadData">鏌ヨ</a-button><a-button @click="resetForm">閲嶇疆</a-button></a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
        <template #result="{ record }">
          <a-space direction="vertical" :size="0">
            <span>A缁?<strong>{{ record.group_a_conversion || 0 }}%</strong></span>
            <span>B缁?<strong>{{ record.group_b_conversion || 0 }}%</strong></span>
          </a-space>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="viewResult(record)">缁撴灉</a-button>
          <a-button v-if="record.status === 'draft'" type="text" size="small" @click="startExperiment(record)">鍚姩</a-button>
          <a-button v-if="record.status === 'running'" type="text" size="small" @click="pauseExperiment(record)">鏆傚仠</a-button>
          <a-button type="text" size="small" @click="editExperiment(record)">缂栬緫</a-button>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model:visible="createModalVisible" :title="isEdit ? '缂栬緫瀹為獙' : '鏂板缓瀹為獙'" @before-ok="handleSubmit" :loading="submitting" :width="600">
      <a-form :model="expForm" layout="vertical">
        <a-form-item label="瀹為獙鍚嶇О" required><a-input v-model="expForm.experiment_name" placeholder="璇疯緭鍏ュ疄楠屽悕绉" /></a-form-item>
        <a-form-item label="瀹為獙鎻忚堪"><a-textarea v-model="expForm.description" :rows="2" placeholder="瀹為獙鎻忚堪" /></a-form-item>
        <a-form-item label="鍒嗘祦姣斾緥">
          <a-space>
            <span>A缁?</span><a-input-number v-model="expForm.group_a_ratio" :min="0" :max="100" :step="5" style="width: 80px" /><span>%</span>
            <span style="margin-left: 16px">B缁?</span><a-input-number v-model="expForm.group_b_ratio" :min="0" :max="100" :step="5" style="width: 80px" /><span>%</span>
          </a-space>
        </a-form-item>
        <a-form-item label="A缁勭瓥鐣?>
          <a-select v-model="expForm.group_a_strategy" placeholder="閫夋嫨A缁勭瓥鐣">
            <a-option value="control">瀵圭収缁勶紙鍘熸湁绛栫暐锛?/a-option>
            <a-option value="strategy_1">绛栫暐涓€</a-option>
            <a-option value="strategy_2">绛栫暐浜?/a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="B缁勭瓥鐣?>
          <a-select v-model="expForm.group_b_strategy" placeholder="閫夋嫨B缁勭瓥鐣">
            <a-option value="control">瀵圭収缁勶紙鍘熸湁绛栫暐锛?/a-option>
            <a-option value="strategy_1">绛栫暐涓€</a-option>
            <a-option value="strategy_2">绛栫暐浜?/a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="鐩爣鎸囨爣">
          <a-select v-model="expForm.metric" placeholder="閫夋嫨鐩爣鎸囨爣">
            <a-option value="conversion">杞寲鐜?/a-option>
            <a-option value="retention">鐣欏瓨鐜?/a-option>
            <a-option value="engagement">鐢ㄦ埛娲昏穬</a-option>
            <a-option value="satisfaction">婊℃剰搴?/a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="棰勬湡鏍锋湰閲?>
          <a-input-number v-model="expForm.expected_sample_size" :min="100" placeholder="姣忕粍棰勬湡鏍锋湰閲" style="width: 200px" />
        </a-form-item>
      </a-form>
    </a-modal>
    <a-modal v-model:visible="resultModalVisible" title="瀹為獙缁撴灉" :width="700" :footer="null">
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="瀹為獙鍚嶇О" :span="2">{{ selectedExp?.experiment_name }}</a-descriptions-item>
        <a-descriptions-item label="瀹為獙鐘舵€?><a-tag :color="getStatusColor(selectedExp?.status)">{{ getStatusText(selectedExp?.status) }}</a-tag></a-descriptions-item>
        <a-descriptions-item label="鐩爣鎸囨爣">{{ selectedExp?.metric }}</a-descriptions-item>
        <a-descriptions-item label="鎬绘牱鏈噺">{{ selectedExp?.total_samples || 0 }}</a-descriptions-item>
        <a-descriptions-item label="A缁勬牱鏈?>{{ selectedExp?.group_a_samples || 0 }}</a-descriptions-item>
        <a-descriptions-item label="B缁勬牱鏈?>{{ selectedExp?.group_b_samples || 0 }}</a-descriptions-item>
        <a-descriptions-item label="A缁勮浆鍖栫巼">{{ selectedExp?.group_a_conversion || 0 }}%</a-descriptions-item>
        <a-descriptions-item label="B缁勮浆鍖栫巼">{{ selectedExp?.group_b_conversion || 0 }}%</a-descriptions-item>
        <a-descriptions-item label="缁熻鏄捐憲鎬? :span="2">{{ selectedExp?.significance || 'N/A' }}</a-descriptions-item>
        <a-descriptions-item label="缁撹" :span="2">{{ selectedExp?.conclusion || '瀹為獙杩涜涓? }}</a-descriptions-item>
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
  { title: '瀹為獙鍚嶇О', dataIndex: 'experiment_name', width: 180 },
  { title: '鐘舵€?, slotName: 'status', width: 90 },
  { title: 'A/B缁勮浆鍖栫巼', slotName: 'result', width: 140 },
  { title: '鎬绘牱鏈?, dataIndex: 'total_samples', width: 90 },
  { title: '寮€濮嬫椂闂?, dataIndex: 'started_at', width: 170 },
  { title: '鎿嶄綔', slotName: 'actions', width: 180, fixed: 'right' }
]

const getStatusColor = (s) => ({ draft: 'gray', running: 'arcoblue', paused: 'orange', completed: 'green' }[s] || 'gray')
const getStatusText = (s) => ({ draft: '鑽夌', running: '杩愯涓?, paused: '宸叉殏鍋?, completed: '宸插畬鎴? }[s] || s)

const resetForm = () => {
  Object.keys(form).forEach(k => form[k] = '')
  loadData()
}

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
  } catch (e) { Message.error('鍔犺浇澶辫触') } finally { loading.value = false }
}

const openCreateModal = () => {
  isEdit.value = false
  Object.assign(expForm, { id: null, experiment_name: '', description: '', group_a_ratio: 50, group_b_ratio: 50, group_a_strategy: 'control', group_b_strategy: 'strategy_1', metric: 'conversion', expected_sample_size: 1000 })
  createModalVisible.value = true
}
const editExperiment = (record) => {
  isEdit.value = true
  Object.assign(expForm, record)
  createModalVisible.value = true
}

const handleSubmit = async (done) => {
  if (!expForm.experiment_name) { Message.warning('璇疯緭鍏ュ疄楠屽悕绉?); done(false); return }
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const url = isEdit.value ? `/api/v1/ai/ab-test/experiments/${expForm.id}` : '/api/v1/ai/ab-test/experiments'
    const res = await fetch(url, { method: isEdit.value ? 'PUT' : 'POST', headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(expForm) }).then(r => r.json())
    if (res.code === 0) { Message.success(isEdit.value ? '鏇存柊鎴愬姛' : '鍒涘缓鎴愬姛'); createModalVisible.value = false; loadData() }
    else { Message.error(res.message || '鎿嶄綔澶辫触') }
    done(true)
  } catch (e) { Message.error('鎿嶄綔澶辫触'); done(false) } finally { submitting.value = false }
}

const viewResult = (record) => { selectedExp.value = record; resultModalVisible.value = true }

const startExperiment = async (record) => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/ai/ab-test/experiments/${record.id}/start`, { method: 'POST', headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { Message.success('瀹為獙宸插惎鍔?); loadData() }
    else Message.error('鍚姩澶辫触')
  } catch (e) { Message.error('鍚姩澶辫触') }
}

const pauseExperiment = async (record) => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/ai/ab-test/experiments/${record.id}/pause`, { method: 'POST', headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { Message.success('瀹為獙宸叉殏鍋?); loadData() }
    else Message.error('鏆傚仠澶辫触')
  } catch (e) { Message.error('鏆傚仠澶辫触') }
}

const onPageChange = (page) => { pagination.current = page; loadData() }

onMounted(() => loadData())
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>
