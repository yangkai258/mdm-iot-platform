<template>
  <div class="page-container">
    <a-card class="general-card" title="璁惧杩滅▼鎿﹂櫎">
      <template #extra>
        <a-button @click="loadData"><icon-refresh />鍒锋柊</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="璁惧ID"><a-input v-model="form.device_id" placeholder="璇疯緭锟" /></a-form-item>
          <a-form-item label="鐘讹拷?>
            <a-select v-model="form.status" placeholder="閫夋嫨鐘讹拷" allow-clear style="width: 140px">
              <a-option value="pending">寰呭锟?/a-option>
              <a-option value="wiping">鎿﹂櫎锟?/a-option>
              <a-option value="completed">宸插畬锟?/a-option>
              <a-option value="failed">澶辫触</a-option>
            </a-select>
          </a-form-item>
          <a-form-item><a-button type="primary" @click="handleSearch">鏌ヨ</a-button><a-button @click="handleReset">閲嶇疆</a-button></a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
        <template #progress="{ record }">
          <a-progress :percent="record.progress || 0" :status="record.status === 'failed' ? 'exception' : 'normal'" size="small" />
        </template>
        <template #actions="{ record }">
          <a-button v-if="record.status === 'pending'" type="primary" size="small" @click="confirmWipe(record)">纭鎿﹂櫎</a-button>
          <a-button v-if="record.status === 'wiping'" type="text" size="small" @click="viewProgress(record)">鏌ョ湅杩涘害</a-button>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model:visible="confirmVisible" title="杩滅▼鎿﹂櫎纭" @before-ok="executeWipe" :loading="submitting">
      <a-result status="warning" title="鍗冲皢鎵ц杩滅▼鎿﹂櫎">
        <template #subtitle>
          <div>
            <p>璁惧ID: <strong>{{ selectedDevice?.device_id }}</strong></p>
            <p>璁惧鍚嶇О: {{ selectedDevice?.device_name }}</p>
            <p style="color: #f53f3f">姝ゆ搷浣滃皢鎿﹂櫎璁惧涓婄殑鎵€鏈夋暟鎹紝涓旀棤娉曟仮澶嶏紒</p>
          </div>
        </template>
      </a-result>
      <a-form layout="vertical" style="margin-top: 16px">
        <a-form-item label="鎿﹂櫎妯″紡">
          <a-radio-group v-model="wipeMode">
            <a-radio value="fast">蹇€熸摝闄わ紙浠呯敤鎴锋暟鎹級</a-radio>
            <a-radio value="full">瀹屽叏鎿﹂櫎锛堝叏閮ㄦ暟鎹級</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>
    <a-modal v-model:visible="progressVisible" title="鎿﹂櫎杩涘害" :footer="null" :width="560">
      <a-descriptions :column="1" bordered>
        <a-descriptions-item label="璁惧ID">{{ selectedDevice?.device_id }}</a-descriptions-item>
        <a-descriptions-item label="褰撳墠鐘讹拷?><a-tag :color="getStatusColor(selectedDevice?.status)">{{ getStatusText(selectedDevice?.status) }}</a-tag></a-descriptions-item>
        <a-descriptions-item label="鎿﹂櫎杩涘害"><a-progress :percent="selectedDevice?.progress || 0" size="large" /></a-descriptions-item>
        <a-descriptions-item label="寮€濮嬫椂锟?>{{ selectedDevice?.started_at }}</a-descriptions-item>
        <a-descriptions-item label="棰勮瀹屾垚">{{ selectedDevice?.estimated_time }}</a-descriptions-item>
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
const confirmVisible = ref(false)
const progressVisible = ref(false)
const selectedDevice = ref(null)
const wipeMode = ref('fast')
const form = reactive({ device_id: '', status: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '璁惧ID', dataIndex: 'device_id', width: 120 },
  { title: '璁惧鍚嶇О', dataIndex: 'device_name', width: 140 },
  { title: '鎿﹂櫎鐘讹拷?, slotName: 'status', width: 100 },
  { title: '鎿﹂櫎杩涘害', slotName: 'progress', width: 160 },
  { title: '鍙戣捣锟?, dataIndex: 'initiator', width: 100 },
  { title: '鍙戣捣鏃堕棿', dataIndex: 'created_at', width: 170 },
  { title: '鎿嶄綔', slotName: 'actions', width: 120 }
]

const getStatusColor = (status) => ({ pending: 'orange', wiping: 'arcoblue', completed: 'green', failed: 'red' }[status] || 'gray')
const getStatusText = (status) => ({ pending: '寰呭锟?, wiping: '鎿﹂櫎锟?, completed: '宸插畬锟?, failed: '澶辫触' }[status] || status)

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (form.device_id) params.append('device_id', form.device_id)
    if (form.status) params.append('status', form.status)
    const res = await fetch(`/api/v1/device/remote-wipe?${params}`, { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { data.value = res.data?.list || []; pagination.total = res.data?.total || 0 }
    else { data.value = [] }
  } catch (e) { Message.error('鍔犺浇澶辫触') } finally { loading.value = false }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { Object.assign(form, { device_id: '', status: '' }); loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }
const confirmWipe = (record) => { selectedDevice.value = record; confirmVisible.value = true }
const viewProgress = (record) => { selectedDevice.value = record; progressVisible.value = true }

const executeWipe = async (done) => {
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/device/remote-wipe/${selectedDevice.value.id}/execute`, { method: 'POST', headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify({ mode: wipeMode.value }) }).then(r => r.json())
    if (res.code === 0) { Message.success('鎿﹂櫎浠诲姟宸蹭笅锟?); confirmVisible.value = false; loadData() }
    else { Message.error(res.message || '鎿嶄綔澶辫触') }
    done(true)
  } catch (e) { Message.error('鎿嶄綔澶辫触'); done(false) } finally { submitting.value = false }
}

onMounted(() => loadData())
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>