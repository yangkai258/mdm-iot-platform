<template>
  <div class="page-container">
    <Breadcrumb :items="['menu.security', 'menu.security.certificate']" />
    <!-- 顶部工具栏 -->
    <a-card class="toolbar-card">
      <div class="toolbar-row">
        <div class="toolbar-left">
          <a-select
            v-model="filter.type"
            placeholder="类型"
            style="width: 120px"
            allow-clear
            @change="loadCertificates"
          >
            <a-option value="device">设备</a-option>
            <a-option value="user">用户</a-option>
            <a-option value="ca">CA</a-option>
          </a-select>
          <a-select
            v-model="filter.status"
            placeholder="状态"
            style="width: 130px"
            allow-clear
            @change="loadCertificates"
          >
            <a-option value="active">有效</a-option>
            <a-option value="expiring">即将到期</a-option>
            <a-option value="expired">已过期</a-option>
            <a-option value="revoked">已吊销</a-option>
          </a-select>
          <a-select
            v-model="filter.expiry"
            placeholder="到期时间"
            style="width: 140px"
            allow-clear
            @change="loadCertificates"
          >
            <a-option value="30">30天内到期</a-option>
            <a-option value="90">90天内到期</a-option>
            <a-option value="180">180天内到期</a-option>
            <a-option value="expired">已过期</a-option>
          </a-select>
          <a-input-search
            v-model="filter.keyword"
            placeholder="搜索证书名称..."
            style="width: 200px"
            @search="loadCertificates"
            @press-enter="loadCertificates"
          />
        </div>
        <div class="toolbar-right">
          <a-button type="primary" @click="openUploadModal">
            <template #icon><icon-upload /></template>
            上传证书
          </a-button>
        </div>
      </div>
    </a-card>

    <!-- 证书列表 -->
    <a-card class="table-card">
      <a-table
        :columns="columns"
        :data="certificates"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @change="handleTableChange"
      >
        <template #name="{ record }">
          <a-link @click="openDetail(record)">{{ record.subject_cn || record.name }}</a-link>
        </template>
        <template #type="{ record }">
          <a-tag :color="typeColor(record.cert_type)">{{ typeLabel(record.cert_type) }}</a-tag>
        </template>
        <template #serial="{ record }">
          <span class="serial-text">{{ truncateSerial(record.cert_serial) }}</span>
        </template>
        <template #expiry="{ record }">
          <span :class="expiryClass(record.not_after)">{{ formatDate(record.not_after) }}</span>
        </template>
        <template #status="{ record }">
          <a-tag :color="statusColor(record)">{{ statusLabel(record) }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
          <a-divider direction="vertical" />
          <a-button type="text" size="small" @click="handleDownload(record)">下载</a-button>
          <a-divider direction="vertical" />
          <a-popconfirm content="确定吊销该证书？" @ok="handleRevoke(record)">
            <a-button type="text" size="small" status="danger">吊销</a-button>
          </a-popconfirm>
          <a-divider direction="vertical" />
          <a-popconfirm content="确定续期该证书？" @ok="handleRenew(record)">
            <a-button type="text" size="small">续期</a-button>
          </a-popconfirm>
        </template>
      </a-table>
    </a-card>

    <!-- 证书详情弹窗 -->
    <a-modal
      v-model:visible="detailVisible"
      title="证书详情"
      width="560px"
      :footer="null"
    >
      <div class="cert-detail" v-if="currentCert">
        <div class="detail-row">
          <span class="detail-label">证书名称</span>
          <span class="detail-value">{{ currentCert.subject_cn || '-' }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">类型</span>
          <span class="detail-value">{{ typeLabel(currentCert.cert_type) }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">序列号</span>
          <span class="detail-value mono">{{ currentCert.cert_serial }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">颁发者</span>
          <span class="detail-value">{{ currentCert.issuer_dn || '-' }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">生效时间</span>
          <span class="detail-value">{{ formatDate(currentCert.not_before) }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">到期时间</span>
          <span class="detail-value" :class="expiryClass(currentCert.not_after)">{{ formatDate(currentCert.not_after) }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">指纹</span>
          <span class="detail-value mono">{{ currentCert.fingerprint || '-' }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">状态</span>
          <span class="detail-value">
            <a-tag :color="statusColor(currentCert)">{{ statusLabel(currentCert) }}</a-tag>
          </span>
        </div>
      </div>
    </a-modal>

    <!-- 上传证书弹窗 -->
    <a-modal
      v-model:visible="uploadVisible"
      title="上传证书"
      width="480px"
      @ok="handleUpload"
      :ok-loading="uploading"
    >
      <a-form :model="uploadForm" layout="vertical">
        <a-form-item label="证书名称" required>
          <a-input v-model="uploadForm.name" placeholder="请输入证书名称" />
        </a-form-item>
        <a-form-item label="证书类型" required>
          <a-select v-model="uploadForm.type" placeholder="请选择类型">
            <a-option value="device">设备证书</a-option>
            <a-option value="user">用户证书</a-option>
            <a-option value="ca">CA证书</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="证书文件" required>
          <a-upload
            draggable
            accept=".pem,.crt,.pfx"
            :limit="1"
            @change="handleFileChange"
          >
            <div class="upload-hint">
              <icon-upload />
              <div>点击或拖拽上传证书文件</div>
              <div class="upload-formats">支持 .pem / .crt / .pfx 格式</div>
            </div>
          </a-upload>
        </a-form-item>
        <a-form-item label="私钥文件（可选）" v-if="uploadForm.type !== 'ca'">
          <a-upload
            draggable
            accept=".pem,.key"
            :limit="1"
            @change="handleKeyChange"
          >
            <div class="upload-hint">
              <icon-upload />
              <div>点击或拖拽上传私钥文件</div>
            </div>
          </a-upload>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { Message } from '@arco-design/web-vue'
import { getCertificates, uploadCertificate, revokeCertificate, renewCertificate, downloadCertificate } from '@/api/security'
import dayjs from 'dayjs'

const loading = ref(false)
const uploading = ref(false)
const certificates = ref([])
const detailVisible = ref(false)
const uploadVisible = ref(false)
const currentCert = ref(null)

const filter = reactive({
  type: '',
  status: '',
  expiry: '',
  keyword: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const uploadForm = reactive({
  name: '',
  type: 'device',
  certFile: null,
  keyFile: null
})

const columns = [
  { title: '证书名称', slotName: 'name', minWidth: 160 },
  { title: '类型', slotName: 'type', width: 90 },
  { title: '序列号', slotName: 'serial', width: 140 },
  { title: '到期时间', slotName: 'expiry', width: 140 },
  { title: '状态', slotName: 'status', width: 110 },
  { title: '操作', slotName: 'actions', width: 260, fixed: 'right' }
]

onMounted(() => {
  loadCertificates()
})

async function loadCertificates() {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
      ...filter
    }
    const res = await getCertificates(params)
    const data = res.data || res
    certificates.value = Array.isArray(data) ? data : (data.list || data.records || [])
    pagination.total = data.total || certificates.value.length
  } catch (e) {
    console.error('加载证书列表失败', e)
  } finally {
    loading.value = false
  }
}

function handleTableChange(pag) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadCertificates()
}

function typeColor(type) {
  const map = { device: 'blue', user: 'green', ca: 'orange' }
  return map[type] || 'default'
}

function typeLabel(type) {
  const map = { device: '设备', user: '用户', ca: 'CA' }
  return map[type] || type
}

function statusColor(cert) {
  if (cert.status === 'revoked') return 'red'
  const now = dayjs()
  const expiry = dayjs(cert.not_after)
  if (expiry.isBefore(now)) return 'gray'
  if (expiry.isBefore(now.add(30, 'day'))) return 'yellow'
  return 'green'
}

function statusLabel(cert) {
  if (cert.status === 'revoked') return '已吊销'
  const now = dayjs()
  const expiry = dayjs(cert.not_after)
  if (expiry.isBefore(now)) return '已过期'
  if (expiry.isBefore(now.add(30, 'day'))) return '即将到期'
  return '有效'
}

function expiryClass(date) {
  const now = dayjs()
  const expiry = dayjs(date)
  if (expiry.isBefore(now)) return 'text-expired'
  if (expiry.isBefore(now.add(30, 'day'))) return 'text-warning'
  return ''
}

function truncateSerial(serial) {
  if (!serial) return '-'
  return serial.length > 16 ? serial.substring(0, 8) + ':...' : serial
}

function formatDate(date) {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD')
}

function openDetail(cert) {
  currentCert.value = cert
  detailVisible.value = true
}

function openUploadModal() {
  uploadForm.name = ''
  uploadForm.type = 'device'
  uploadForm.certFile = null
  uploadForm.keyFile = null
  uploadVisible.value = true
}

function handleFileChange(file) {
  uploadForm.certFile = file
}

function handleKeyChange(file) {
  uploadForm.keyFile = file
}

async function handleUpload() {
  if (!uploadForm.name || !uploadForm.certFile) {
    Message.warning('请填写证书名称并上传证书文件')
    return
  }
  uploading.value = true
  try {
    const formData = new FormData()
    formData.append('name', uploadForm.name)
    formData.append('type', uploadForm.type)
    formData.append('cert_file', uploadForm.certFile)
    if (uploadForm.keyFile) {
      formData.append('key_file', uploadForm.keyFile)
    }
    await uploadCertificate(formData)
    Message.success('上传成功')
    uploadVisible.value = false
    loadCertificates()
  } catch (e) {
    Message.error('上传失败')
  } finally {
    uploading.value = false
  }
}

async function handleDownload(cert) {
  try {
    const blob = await downloadCertificate(cert.id)
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${cert.subject_cn || cert.name}.pem`
    a.click()
    URL.revokeObjectURL(url)
    Message.success('下载成功')
  } catch (e) {
    Message.error('下载失败')
  }
}

async function handleRevoke(cert) {
  try {
    await revokeCertificate(cert.id)
    Message.success('吊销成功')
    loadCertificates()
  } catch (e) {
    Message.error('吊销失败')
  }
}

async function handleRenew(cert) {
  try {
    await renewCertificate(cert.id)
    Message.success('续期成功')
    loadCertificates()
  } catch (e) {
    Message.error('续期失败')
  }
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

.toolbar-card {
  flex-shrink: 0;
}

.toolbar-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.toolbar-right {
  display: flex;
  gap: 8px;
}

.table-card {
  flex: 1;
  overflow: auto;
}

.serial-text {
  font-family: monospace;
  font-size: 12px;
  color: var(--color-text-3);
}

.text-expired {
  color: var(--color-text-3);
  text-decoration: line-through;
}

.text-warning {
  color: #fe9c3e;
}

.cert-detail {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.detail-row {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.detail-label {
  width: 90px;
  color: var(--color-text-3);
  font-size: 13px;
  flex-shrink: 0;
}

.detail-value {
  flex: 1;
  font-size: 13px;
  word-break: break-all;
}

.mono {
  font-family: monospace;
  font-size: 12px;
}

.upload-hint {
  padding: 20px;
  text-align: center;
  color: var(--color-text-3);
}

.upload-formats {
  font-size: 12px;
  margin-top: 4px;
}
</style>
