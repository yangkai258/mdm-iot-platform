<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>设备管理</a-breadcrumb-item>
      <a-breadcrumb-item>证书管理</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stat-row">
      <a-col :span="6">
        <a-card class="stat-card" hoverable>
          <a-statistic :value="stats.total" title="证书总数">
            <template #icon><icon-certificate style="font-size: 24px; color: #1650d8" /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card" hoverable>
          <a-statistic :value="stats.active" title="有效证书">
            <template #icon><icon-check-circle style="font-size: 24px; color: #00b42a" /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card" hoverable>
          <a-statistic :value="stats.expiring" title="即将到期">
            <template #icon><icon-clock style="font-size: 24px; color: #ff7d00" /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card" hoverable>
          <a-statistic :value="stats.revoked" title="已吊销">
            <template #icon><icon-close-circle style="font-size: 24px; color: #f53f3f" /></template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- 搜索栏 -->
    <div class="pro-search-bar">
      <a-space wrap>
        <a-input-search
          v-model="searchKeyword"
          placeholder="搜索证书名称/序列号/主题"
          style="width: 260px"
          @search="loadCertificates"
          search-button
        />
        <a-select v-model="filterType" placeholder="证书类型" allow-clear style="width: 140px" @change="loadCertificates">
          <a-option value="device">设备证书</a-option>
          <a-option value="client">客户端证书</a-option>
          <a-option value="server">服务器证书</a-option>
          <a-option value="ca">CA证书</a-option>
        </a-select>
        <a-select v-model="filterStatus" placeholder="证书状态" allow-clear style="width: 130px" @change="loadCertificates">
          <a-option value="active">有效</a-option>
          <a-option value="expired">已过期</a-option>
          <a-option value="revoked">已吊销</a-option>
          <a-option value="pending">待激活</a-option>
        </a-select>
      </a-space>
    </div>

    <!-- 操作按钮组 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showCreateModal">新建证书</a-button>
        <a-button @click="loadCertificates">刷新</a-button>
      </a-space>
    </div>

    <!-- 证书列表 -->
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
            <a-button type="text" size="small" @click="showDetailModal(record)">详情</a-button>
            <a-button type="text" size="small" @click="downloadCert(record)">下载</a-button>
            <a-button type="text" size="small" status="danger" @click="revokeCert(record)" :disabled="record.status === 'revoked'">吊销</a-button>
            <a-button type="text" size="small" status="danger" @click="deleteCert(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 新建证书弹窗 -->
    <a-modal v-model:visible="createModalVisible" title="新建证书" :width="560" :loading="submitting" @before-ok="handleCreate" @cancel="createModalVisible = false">
      <a-form :model="createForm" layout="vertical">
        <a-form-item label="证书名称" required>
          <a-input v-model="createForm.cert_name" placeholder="请输入证书名称" />
        </a-form-item>
        <a-form-item label="证书类型" required>
          <a-select v-model="createForm.cert_type" placeholder="选择证书类型">
            <a-option value="device">设备证书</a-option>
            <a-option value="client">客户端证书</a-option>
            <a-option value="server">服务器证书</a-option>
            <a-option value="ca">CA证书</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="证书文件 (PEM)" required>
          <a-upload
            :limit="1"
            :auto-upload="false"
            accept=".pem,.crt,.cer"
            @change="(files) => handleFileChange(files, 'cert')"
          >
            <template #upload-button>
              <div class="upload-trigger">
                <icon-upload /> 点击上传 PEM 证书文件
              </div>
            </template>
          </a-upload>
        </a-form-item>
        <a-form-item label="私钥文件 (PEM)" required>
          <a-upload
            :limit="1"
            :auto-upload="false"
            accept=".pem,.key"
            @change="(files) => handleFileChange(files, 'key')"
          >
            <template #upload-button>
              <div class="upload-trigger">
                <icon-upload /> 点击上传 PEM 私钥文件
              </div>
            </template>
          </a-upload>
        </a-form-item>
        <a-form-item label="到期提醒天数">
          <a-input-number v-model="createForm.notify_days" :min="1" :max="365" placeholder="默认30天" style="width: 200px" />
          <span class="form-tip">证书到期前多少天发送提醒</span>
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="createForm.description" placeholder="可选描述信息" :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 证书详情弹窗 -->
    <a-modal v-model:visible="detailModalVisible" title="证书详情" :width="600" :footer="null">
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="证书名称">{{ currentCert?.cert_name }}</a-descriptions-item>
        <a-descriptions-item label="证书类型">
          <a-tag :color="getTypeColor(currentCert?.cert_type)">{{ getTypeText(currentCert?.cert_type) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="序列号" :span="2">
          <a-tooltip :content="currentCert?.serial_number">{{ currentCert?.serial_number }}</a-tooltip>
        </a-descriptions-item>
        <a-descriptions-item label="主题 (Subject)" :span="2">{{ currentCert?.subject }}</a-descriptions-item>
        <a-descriptions-item label="颁发者 (Issuer)" :span="2">{{ currentCert?.issuer }}</a-descriptions-item>
        <a-descriptions-item label="SHA1指纹" :span="2">
          <a-tooltip :content="currentCert?.thumbprint">{{ currentCert?.thumbprint }}</a-tooltip>
        </a-descriptions-item>
        <a-descriptions-item label="生效时间">{{ formatDate(currentCert?.not_before) }}</a-descriptions-item>
        <a-descriptions-item label="到期时间">{{ formatDate(currentCert?.not_after) }}</a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="getStatusColor(currentCert?.status)">{{ getStatusText(currentCert?.status) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="描述" :span="2">{{ currentCert?.description || '-' }}</a-descriptions-item>
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
  { title: '证书名称', slotName: 'cert_name', ellipsis: true },
  { title: '证书类型', slotName: 'cert_type', width: 110 },
  { title: '序列号', dataIndex: 'serial_number', ellipsis: true, width: 160 },
  { title: '主题', dataIndex: 'subject', ellipsis: true, width: 140 },
  { title: '有效期', slotName: 'validity', width: 220 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '操作', slotName: 'actions', width: 220, fixed: 'right' }
]

const getTypeColor = (t) => ({ device: 'blue', client: 'green', server: 'orange', ca: 'purple' }[t] || 'gray')
const getTypeText = (t) => ({ device: '设备证书', client: '客户端证书', server: '服务器证书', ca: 'CA证书' }[t] || t)
const getStatusColor = (s) => ({ active: 'green', expired: 'red', revoked: 'orange', pending: 'blue' }[s] || 'gray')
const getStatusText = (s) => ({ active: '有效', expired: '已过期', revoked: '已吊销', pending: '待激活' }[s] || s)
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
    // 加载统计
    loadStats()
  } catch (e) {
    Message.error('加载证书列表失败')
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
    Message.warning('请填写证书名称和类型')
    done(false)
    return
  }
  if (!certFile.value || !keyFile.value) {
    Message.warning('请上传证书文件和私钥文件')
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
      Message.success('证书创建成功')
      createModalVisible.value = false
      loadCertificates()
      done(true)
    } else {
      Message.error(json.message || '创建失败')
      done(false)
    }
  } catch (e) {
    Message.error('创建失败')
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
    Message.success('下载成功')
  } catch (e) {
    Message.error('下载失败')
  }
}

const revokeCert = async (record) => {
  Modal.warning({
    title: '确认吊销证书',
    content: `确定要吊销证书「${record.cert_name}」吗？此操作不可逆。`,
    okText: '确认吊销',
    onOk: async () => {
      try {
        const res = await revokeCertificate(record.id)
        if (res.code === 0) {
          Message.success('证书已吊销')
          loadCertificates()
        } else {
          Message.error(res.message || '吊销失败')
        }
      } catch (e) {
        Message.error('吊销失败')
      }
    }
  })
}

const deleteCert = async (record) => {
  Modal.warning({
    title: '确认删除证书',
    content: `确定要删除证书「${record.cert_name}」吗？此操作不可逆。`,
    okText: '确认删除',
    onOk: async () => {
      try {
        const res = await deleteCertificate(record.id)
        if (res.code === 0) {
          Message.success('证书已删除')
          loadCertificates()
        } else {
          Message.error(res.message || '删除失败')
        }
      } catch (e) {
        Message.error('删除失败')
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
