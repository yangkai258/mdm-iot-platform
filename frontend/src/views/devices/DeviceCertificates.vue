<template>
  <div class="pro-page-container">
    <!-- 闈㈠寘灞?-->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>棣栭〉</a-breadcrumb-item>
      <a-breadcrumb-item>璁惧绠＄悊</a-breadcrumb-item>
      <a-breadcrumb-item>璇佷功绠＄悊</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 缁熻鍗＄墖 -->
    <a-row :gutter="16" class="stat-row">
      <a-col :span="6">
        <a-card class="stat-card" hoverable>
          <a-statistic :value="stats.total" title="璇佷功鎬绘暟">
            <template #icon><icon-safe style="font-size: 24px; color: #1650d8" /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card" hoverable>
          <a-statistic :value="stats.active" title="鏈夋晥璇佷功">
            <template #icon><icon-check-circle style="font-size: 24px; color: #00b42a" /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card" hoverable>
          <a-statistic :value="stats.expiring" title="鍗冲皢鍒版湡">
            <template #icon><icon-clock-circle style="font-size: 24px; color: #ff7d00" /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card" hoverable>
          <a-statistic :value="stats.revoked" title="宸插悐閿€">
            <template #icon><icon-close-circle style="font-size: 24px; color: #f53f3f" /></template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- 鎼滅储鏍?-->
    <div class="pro-search-bar">
      <a-space wrap>
        <a-input-search
          v-model="searchKeyword"
          placeholder="鎼滅储璇佷功鍚嶇О/搴忓垪鍙?涓婚"
          style="width: 260px"
          @search="loadCertificates"
          search-button
        />
        <a-select v-model="filterType" placeholder="璇佷功绫诲瀷" allow-clear style="width: 140px" @change="loadCertificates">
          <a-option value="device">璁惧璇佷功</a-option>
          <a-option value="client">瀹㈡埛绔瘉涔?/a-option>
          <a-option value="server">鏈嶅姟鍣ㄨ瘉涔?/a-option>
          <a-option value="ca">CA璇佷功</a-option>
        </a-select>
        <a-select v-model="filterStatus" placeholder="璇佷功鐘舵€? allow-clear style="width: 130px" @change="loadCertificates">
          <a-option value="active">鏈夋晥</a-option>
          <a-option value="expired">宸茶繃鏈?/a-option>
          <a-option value="revoked">宸插悐閿€</a-option>
          <a-option value="pending">寰呮縺娲?/a-option>
        </a-select>
      </a-space>
    </div>

    <!-- 鎿嶄綔鎸夐挳缁?-->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showCreateModal">鏂板缓璇佷功</a-button>
        <a-button @click="loadCertificates">鍒锋柊</a-button>
      </a-space>
    </div>

    <!-- 璇佷功鍒楄〃 -->
    <div class="pro-content-area">
      <a-table
        :columns="columns"
        :data="certificates"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @page-change="handlePageChange"
      >
        <template #cert_name="{ record }">
          <a-tooltip :content="record.subject">
            <span class="cert-name">{{ record.cert_name }}</span>
          </a-tooltip>
        </template>
        <template #cert_type="{ record }">
          <a-tag :color="getTypeColor(record.cert_type)">{{ getTypeText(record.cert_type) }}</a-tag>
        </template>
        <template #validity="{ record }">
          <span :class="{ 'text-expired': record.status === 'expired', 'text-expiring': isExpiringSoon(record.not_after) }">
            {{ formatDate(record.not_before) }} ~ {{ formatDate(record.not_after) }}
          </span>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showDetailModal(record)">璇︽儏</a-button>
            <a-button type="text" size="small" @click="downloadCert(record)">涓嬭浇</a-button>
            <a-button type="text" size="small" status="danger" @click="revokeCert(record)" :disabled="record.status === 'revoked'">鍚婇攢</a-button>
            <a-button type="text" size="small" status="danger" @click="deleteCert(record)">鍒犻櫎</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 鏂板缓璇佷功寮圭獥 -->
    <a-modal v-model:visible="createModalVisible" title="鏂板缓璇佷功" :width="560" :loading="submitting" @before-ok="handleCreate" @cancel="createModalVisible = false">
      <a-form :model="createForm" layout="vertical">
        <a-form-item label="璇佷功鍚嶇О" required>
          <a-input v-model="createForm.cert_name" placeholder="璇疯緭鍏ヨ瘉涔﹀悕绉? />
        </a-form-item>
        <a-form-item label="璇佷功绫诲瀷" required>
          <a-select v-model="createForm.cert_type" placeholder="閫夋嫨璇佷功绫诲瀷">
            <a-option value="device">璁惧璇佷功</a-option>
            <a-option value="client">瀹㈡埛绔瘉涔?/a-option>
            <a-option value="server">鏈嶅姟鍣ㄨ瘉涔?/a-option>
            <a-option value="ca">CA璇佷功</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="璇佷功鏂囦欢 (PEM)" required>
          <a-upload
            :limit="1"
            :auto-upload="false"
            accept=".pem,.crt,.cer"
            @change="(files) => handleFileChange(files, 'cert')"
          >
            <template #upload-button>
              <div class="upload-trigger">
                <icon-upload /> 鐐瑰嚮涓婁紶 PEM 璇佷功鏂囦欢
              </div>
            </template>
          </a-upload>
        </a-form-item>
        <a-form-item label="绉侀挜鏂囦欢 (PEM)" required>
          <a-upload
            :limit="1"
            :auto-upload="false"
            accept=".pem,.key"
            @change="(files) => handleFileChange(files, 'key')"
          >
            <template #upload-button>
              <div class="upload-trigger">
                <icon-upload /> 鐐瑰嚮涓婁紶 PEM 绉侀挜鏂囦欢
              </div>
            </template>
          </a-upload>
        </a-form-item>
        <a-form-item label="鍒版湡鎻愰啋澶╂暟">
          <a-input-number v-model="createForm.notify_days" :min="1" :max="365" placeholder="榛樿30澶? style="width: 200px" />
          <span class="form-tip">璇佷功鍒版湡鍓嶅灏戝ぉ鍙戦€佹彁閱?/span>
        </a-form-item>
        <a-form-item label="鎻忚堪">
          <a-textarea v-model="createForm.description" placeholder="鍙€夋弿杩颁俊鎭? :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 璇佷功璇︽儏寮圭獥 -->
    <a-modal v-model:visible="detailModalVisible" title="璇佷功璇︽儏" :width="600" :footer="null">
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="璇佷功鍚嶇О">{{ currentCert?.cert_name }}</a-descriptions-item>
        <a-descriptions-item label="璇佷功绫诲瀷">
          <a-tag :color="getTypeColor(currentCert?.cert_type)">{{ getTypeText(currentCert?.cert_type) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="搴忓垪鍙? :span="2">
          <a-tooltip :content="currentCert?.serial_number">{{ currentCert?.serial_number }}</a-tooltip>
        </a-descriptions-item>
        <a-descriptions-item label="涓婚 (Subject)" :span="2">{{ currentCert?.subject }}</a-descriptions-item>
        <a-descriptions-item label="棰佸彂鑰?(Issuer)" :span="2">{{ currentCert?.issuer }}</a-descriptions-item>
        <a-descriptions-item label="SHA1鎸囩汗" :span="2">
          <a-tooltip :content="currentCert?.thumbprint">{{ currentCert?.thumbprint }}</a-tooltip>
        </a-descriptions-item>
        <a-descriptions-item label="鐢熸晥鏃堕棿">{{ formatDate(currentCert?.not_before) }}</a-descriptions-item>
        <a-descriptions-item label="鍒版湡鏃堕棿">{{ formatDate(currentCert?.not_after) }}</a-descriptions-item>
        <a-descriptions-item label="鐘舵€?>
          <a-tag :color="getStatusColor(currentCert?.status)">{{ getStatusText(currentCert?.status) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="鎻忚堪" :span="2">{{ currentCert?.description || '-' }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import { getCertificates, getCertificate, uploadCertificate, deleteCertificate, revokeCertificate } from '@/api/security'

const loading = ref(false)
const submitting = ref(false)
const certificates = ref([])
const searchKeyword = ref('')
const filterType = ref('')
const filterStatus = ref('')
const createModalVisible = ref(false)
const detailModalVisible = ref(false)
const currentCert = ref(null)
const stats = reactive({ total: 0, active: 0, expired: 0, revoked: 0, expiring: 0 })

const certFile = ref(null)
const keyFile = ref(null)

const createForm = reactive({
  cert_name: '',
  cert_type: 'device',
  notify_days: 30,
  description: ''
})

const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: '璇佷功鍚嶇О', slotName: 'cert_name', ellipsis: true },
  { title: '璇佷功绫诲瀷', slotName: 'cert_type', width: 110 },
  { title: '搴忓垪鍙?, dataIndex: 'serial_number', ellipsis: true, width: 160 },
  { title: '涓婚', dataIndex: 'subject', ellipsis: true, width: 140 },
  { title: '鏈夋晥鏈?, slotName: 'validity', width: 220 },
  { title: '鐘舵€?, slotName: 'status', width: 90 },
  { title: '鎿嶄綔', slotName: 'actions', width: 220, fixed: 'right' }
]

const getTypeColor = (t) => ({ device: 'blue', client: 'green', server: 'orange', ca: 'purple' }[t] || 'gray')
const getTypeText = (t) => ({ device: '璁惧璇佷功', client: '瀹㈡埛绔瘉涔?, server: '鏈嶅姟鍣ㄨ瘉涔?, ca: 'CA璇佷功' }[t] || t)
const getStatusColor = (s) => ({ active: 'green', expired: 'red', revoked: 'orange', pending: 'blue' }[s] || 'gray')
const getStatusText = (s) => ({ active: '鏈夋晥', expired: '宸茶繃鏈?, revoked: '宸插悐閿€', pending: '寰呮縺娲? }[s] || s)
const isExpiringSoon = (dateStr) => {
  if (!dateStr) return false
  const expiry = new Date(dateStr)
  const now = new Date()
  const diffDays = Math.ceil((expiry - now) / (1000 * 60 * 60 * 24))
  return diffDays > 0 && diffDays <= 30
}
const formatDate = (d) => d ? new Date(d).toLocaleString('zh-CN') : '-'

const loadCertificates = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (searchKeyword.value) params.keyword = searchKeyword.value
    if (filterType.value) params.cert_type = filterType.value
    if (filterStatus.value) params.status = filterStatus.value

    const res = await fetch(`/api/v1/certificates?${new URLSearchParams(params)}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const json = await res.json()
    if (json.code === 0) {
      certificates.value = json.data.list || []
      pagination.total = json.data.total || 0
    }
    // 鍔犺浇缁熻
    loadStats()
  } catch (e) {
    Message.error('鍔犺浇璇佷功鍒楄〃澶辫触')
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/certificates/stats', { headers: { 'Authorization': `Bearer ${token}` } })
    const json = await res.json()
    if (json.code === 0) Object.assign(stats, json.data)
  } catch (e) { /* silent */ }
}

const handlePageChange = (page) => { pagination.current = page; loadCertificates() }

const showCreateModal = () => {
  Object.assign(createForm, { cert_name: '', cert_type: 'device', notify_days: 30, description: '' })
  certFile.value = null
  keyFile.value = null
  createModalVisible.value = true
}

const handleFileChange = (files, type) => {
  if (files.length > 0) {
    if (type === 'cert') certFile.value = files[0].file
    else keyFile.value = files[0].file
  }
}

const handleCreate = async (done) => {
  if (!createForm.cert_name || !createForm.cert_type) {
    Message.warning('璇峰～鍐欒瘉涔﹀悕绉板拰绫诲瀷')
    done(false)
    return
  }
  if (!certFile.value || !keyFile.value) {
    Message.warning('璇蜂笂浼犺瘉涔︽枃浠跺拰绉侀挜鏂囦欢')
    done(false)
    return
  }
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const formData = new FormData()
    formData.append('cert_name', createForm.cert_name)
    formData.append('cert_type', createForm.cert_type)
    formData.append('cert_file_data', await readFileBase64(certFile.value))
    formData.append('key_file_data', await readFileBase64(keyFile.value))
    formData.append('description', createForm.description)

    const res = await fetch('/api/v1/certificates', {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}` },
      body: formData
    })
    const json = await res.json()
    if (json.code === 0) {
      Message.success('璇佷功鍒涘缓鎴愬姛')
      createModalVisible.value = false
      loadCertificates()
      done(true)
    } else {
      Message.error(json.message || '鍒涘缓澶辫触')
      done(false)
    }
  } catch (e) {
    Message.error('鍒涘缓澶辫触')
    done(false)
  } finally {
    submitting.value = false
  }
}

const readFileBase64 = (file) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(reader.result.split(',')[1])
    reader.onerror = reject
    reader.readAsDataURL(file)
  })
}

const showDetailModal = async (record) => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/certificates/${record.id}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const json = await res.json()
    currentCert.value = json.data || record
  } catch {
    currentCert.value = record
  }
  detailModalVisible.value = true
}

const downloadCert = async (record) => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/certificates/${record.id}/download`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const blob = await res.blob()
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `certificate-${record.cert_name || record.id}.pem`
    a.click()
    URL.revokeObjectURL(url)
    Message.success('涓嬭浇鎴愬姛')
  } catch (e) {
    Message.error('涓嬭浇澶辫触')
  }
}

const revokeCert = async (record) => {
  Modal.warning({
    title: '纭鍚婇攢璇佷功',
    content: `纭畾瑕佸悐閿€璇佷功銆?{record.cert_name}銆嶅悧锛熸鎿嶄綔涓嶅彲閫嗐€俙,
    okText: '纭鍚婇攢',
    onOk: async () => {
      try {
        const res = await revokeCertificate(record.id)
        if (res.code === 0) {
          Message.success('璇佷功宸插悐閿€')
          loadCertificates()
        } else {
          Message.error(res.message || '鍚婇攢澶辫触')
        }
      } catch (e) {
        Message.error('鍚婇攢澶辫触')
      }
    }
  })
}

const deleteCert = async (record) => {
  Modal.warning({
    title: '纭鍒犻櫎璇佷功',
    content: `纭畾瑕佸垹闄よ瘉涔︺€?{record.cert_name}銆嶅悧锛熸鎿嶄綔涓嶅彲閫嗐€俙,
    okText: '纭鍒犻櫎',
    onOk: async () => {
      try {
        const res = await deleteCertificate(record.id)
        if (res.code === 0) {
          Message.success('璇佷功宸插垹闄?)
          loadCertificates()
        } else {
          Message.error(res.message || '鍒犻櫎澶辫触')
        }
      } catch (e) {
        Message.error('鍒犻櫎澶辫触')
      }
    }
  })
}

onMounted(() => { loadCertificates() })
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.stat-row { margin-bottom: 16px; }
.stat-card { border-radius: 8px; }
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; }
.pro-content-area { background: #fff; border-radius: 8px; padding: 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.04); }
.cert-name { font-weight: 500; }
.text-expired { color: #f53f3f; }
.text-expiring { color: #ff7d00; }
.upload-trigger {
  border: 1px dashed #ccc;
  border-radius: 4px;
  padding: 16px;
  text-align: center;
  cursor: pointer;
  color: #666;
}
.upload-trigger:hover { border-color: #1650d8; color: #1650d8; }
.form-tip { color: #999; font-size: 12px; margin-left: 8px; }
</style>


