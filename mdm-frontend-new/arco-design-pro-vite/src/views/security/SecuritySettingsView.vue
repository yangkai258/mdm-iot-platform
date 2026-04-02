<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="page-container">
    <!-- 加密配置 -->
    <a-card class="section-card">
      <template #title>
        <div class="card-title">
          <icon-lock />
          <span>加密配置</span>
        </div>
      </template>
      <a-form :model="encryptionConfig" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="传输加密">
              <a-switch v-model="encryptionConfig.tlsEnabled" />
              <span class="config-hint ml-3">{{ encryptionConfig.tlsEnabled ? '已启用' : '已禁用' }}</span>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="TLS 版本">
              <a-select v-model="encryptionConfig.tlsVersion" :disabled="!encryptionConfig.tlsEnabled">
                <a-option value="1.2">TLS 1.2</a-option>
                <a-option value="1.3">TLS 1.3</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="数据加密算法">
              <a-select v-model="encryptionConfig.algorithm">
                <a-option value="AES-256-GCM">AES-256-GCM</a-option>
                <a-option value="AES-128-GCM">AES-128-GCM</a-option>
                <a-option value="ChaCha20-Poly1305">ChaCha20-Poly1305</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="哈希算法">
              <a-select v-model="encryptionConfig.hashAlgorithm">
                <a-option value="SHA-256">SHA-256</a-option>
                <a-option value="SHA-384">SHA-384</a-option>
                <a-option value="SHA-512">SHA-512</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item>
          <a-button type="primary" :loading="saving.encryption" @click="saveEncryption">
            保存配置
          </a-button>
        </a-form-item>
      </a-form>
    </a-card>

    <!-- 密钥管理 -->
    <a-card class="section-card">
      <template #title>
        <div class="card-title">
          <icon-key />
          <span>密钥管理</span>
        </div>
      </template>
      <div class="toolbar-row">
        <div class="toolbar-left">
          <a-input-search
            v-model="keyFilter.keyword"
            placeholder="搜索密钥名称..."
            style="width: 200px"
            @search="loadKeys"
            @press-enter="loadKeys"
          />
          <a-select
            v-model="keyFilter.type"
            placeholder="密钥类型"
            style="width: 140px"
            allow-clear
            @change="loadKeys"
          >
            <a-option value="symmetric">对称密钥</a-option>
            <a-option value="asymmetric">非对称密钥</a-option>
            <a-option value="hmac">HMAC</a-option>
          </a-select>
        </div>
        <div class="toolbar-right">
          <a-button type="primary" @click="openCreateKeyModal">
            <template #icon><icon-plus /></template>
            创建密钥
          </a-button>
        </div>
      </div>
      <a-table
        :columns="keyColumns"
        :data="keys"
        :loading="loading.keys"
        :pagination="pagination"
        row-key="id"
        @change="handleKeyTableChange"
      >
        <template #name="{ record }">
          <a-link @click="openKeyDetail(record)">{{ record.name }}</a-link>
        </template>
        <template #type="{ record }">
          <a-tag :color="keyTypeColor(record.type)">{{ keyTypeLabel(record.type) }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-badge :status="record.enabled ? 'success' : 'danger'" />
          {{ record.enabled ? '启用' : '禁用' }}
        </template>
        <template #created_at="{ record }">
          {{ formatDate(record.created_at) }}
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="openKeyDetail(record)">详情</a-button>
          <a-divider direction="vertical" />
          <a-button type="text" size="small" @click="toggleKeyStatus(record)">
            {{ record.enabled ? '禁用' : '启用' }}
          </a-button>
          <a-divider direction="vertical" />
          <a-popconfirm content="确定删除该密钥？" @ok="deleteKey(record)">
            <a-button type="text" size="small" status="danger">删除</a-button>
          </a-popconfirm>
        </template>
      </a-table>
    </a-card>

    <!-- 密钥轮换 -->
    <a-card class="section-card">
      <template #title>
        <div class="card-title">
          <icon-sync />
          <span>密钥轮换</span>
        </div>
      </template>
      <a-form :model="rotationConfig" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="自动轮换">
              <a-switch v-model="rotationConfig.autoRotation" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="轮换周期">
              <a-select v-model="rotationConfig.interval" :disabled="!rotationConfig.autoRotation">
                <a-option value="30">30 天</a-option>
                <a-option value="60">60 天</a-option>
                <a-option value="90">90 天</a-option>
                <a-option value="180">180 天</a-option>
                <a-option value="365">365 天</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="提前预警天数">
              <a-input-number
                v-model="rotationConfig.warningDays"
                :min="1"
                :max="30"
                :disabled="!rotationConfig.autoRotation"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="上次轮换时间">
              <span class="info-text">{{ rotationConfig.lastRotation || '从未轮换' }}</span>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="24">
            <a-form-item label="下次轮换时间">
              <span class="info-text">{{ rotationConfig.nextRotation || '未设置' }}</span>
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item>
          <a-space>
            <a-button type="primary" :loading="saving.rotation" @click="saveRotation">
              保存配置
            </a-button>
            <a-button :loading="rotating" @click="rotateNow">
              立即轮换
            </a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>

    <!-- 创建密钥弹窗 -->
    <a-modal
      v-model:visible="createKeyVisible"
      title="创建密钥"
      width="480px"
      @ok="handleCreateKey"
      :ok-loading="creatingKey"
    >
      <a-form :model="newKeyForm" layout="vertical">
        <a-form-item label="密钥名称" required>
          <a-input v-model="newKeyForm.name" placeholder="请输入密钥名称" />
        </a-form-item>
        <a-form-item label="密钥类型" required>
          <a-select v-model="newKeyForm.type" placeholder="请选择密钥类型">
            <a-option value="symmetric">对称密钥</a-option>
            <a-option value="asymmetric">非对称密钥</a-option>
            <a-option value="hmac">HMAC</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="密钥用途" required>
          <a-select v-model="newKeyForm.purpose" placeholder="请选择密钥用途">
            <a-option value="encryption">数据加密</a-option>
            <a-option value="signing">数字签名</a-option>
            <a-option value="authentication">身份认证</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="密钥长度">
          <a-select v-model="newKeyForm.keySize" placeholder="请选择密钥长度">
            <a-option value="128">128 位</a-option>
            <a-option value="256">256 位</a-option>
            <a-option value="4096">4096 位</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 密钥详情弹窗 -->
    <a-modal
      v-model:visible="keyDetailVisible"
      title="密钥详情"
      width="560px"
      :footer="null"
    >
      <div class="detail-grid" v-if="currentKey">
        <div class="detail-item">
          <span class="detail-label">密钥名称</span>
          <span class="detail-value">{{ currentKey.name }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">密钥类型</span>
          <span class="detail-value">
            <a-tag :color="keyTypeColor(currentKey.type)">{{ keyTypeLabel(currentKey.type) }}</a-tag>
          </span>
        </div>
        <div class="detail-item">
          <span class="detail-label">密钥用途</span>
          <span class="detail-value">{{ purposeLabel(currentKey.purpose) }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">状态</span>
          <span class="detail-value">
            <a-badge :status="currentKey.enabled ? 'success' : 'danger'" />
            {{ currentKey.enabled ? '启用' : '禁用' }}
          </span>
        </div>
        <div class="detail-item">
          <span class="detail-label">密钥 ID</span>
          <span class="detail-value mono">{{ currentKey.key_id }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">创建时间</span>
          <span class="detail-value">{{ formatDate(currentKey.created_at) }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">最后使用</span>
          <span class="detail-value">{{ currentKey.last_used_at ? formatDate(currentKey.last_used_at) : '从未使用' }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">使用次数</span>
          <span class="detail-value">{{ currentKey.use_count || 0 }}</span>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import {
  getEncryptionConfig,
  updateEncryptionConfig,
  getKeys,
  createKey,
  deleteKey as deleteKeyApi,
  toggleKeyStatus as toggleKeyStatusApi,
  getRotationConfig,
  updateRotationConfig,
  rotateKeys
} from '@/api/security'
import dayjs from 'dayjs'

const loading = reactive({ keys: false, config: false })
const saving = reactive({ encryption: false, rotation: false })
const rotating = ref(false)
const keys = ref([])
const keysList = ref([])
const currentKey = ref(null)

const keyFilter = reactive({
  keyword: '',
  type: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const encryptionConfig = reactive({
  tlsEnabled: true,
  tlsVersion: '1.3',
  algorithm: 'AES-256-GCM',
  hashAlgorithm: 'SHA-256'
})

const rotationConfig = reactive({
  autoRotation: true,
  interval: '90',
  warningDays: 7,
  lastRotation: '',
  nextRotation: ''
})

const createKeyVisible = ref(false)
const keyDetailVisible = ref(false)
const creatingKey = ref(false)

const newKeyForm = reactive({
  name: '',
  type: 'symmetric',
  purpose: 'encryption',
  keySize: '256'
})

const keyColumns = [
  { title: '密钥名称', slotName: 'name', minWidth: 160 },
  { title: '类型', slotName: 'type', width: 110 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '创建时间', slotName: 'created_at', width: 140 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' }
]

onMounted(() => {
  loadEncryptionConfig()
  loadRotationConfig()
  loadKeys()
})

async function loadEncryptionConfig() {
  try {
    const res = await getEncryptionConfig()
    const data = res.data || res
    Object.assign(encryptionConfig, data)
  } catch (e) {
    console.error('加载加密配置失败', e)
  }
}

async function loadRotationConfig() {
  try {
    const res = await getRotationConfig()
    const data = res.data || res
    Object.assign(rotationConfig, data)
  } catch (e) {
    console.error('加载轮换配置失败', e)
  }
}

async function loadKeys() {
  loading.keys = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
      ...keyFilter
    }
    const res = await getKeys(params)
    const data = res.data || res
    keys.value = data.list || data.records || []
    pagination.total = data.total || keys.value.length
  } catch (e) {
    console.error('加载密钥列表失败', e)
  } finally {
    loading.keys = false
  }
}

function handleKeyTableChange(pag) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadKeys()
}

function keyTypeColor(type) {
  const map = { symmetric: 'blue', asymmetric: 'green', hmac: 'orange' }
  return map[type] || 'default'
}

function keyTypeLabel(type) {
  const map = { symmetric: '对称密钥', asymmetric: '非对称密钥', hmac: 'HMAC' }
  return map[type] || type
}

function purposeLabel(purpose) {
  const map = { encryption: '数据加密', signing: '数字签名', authentication: '身份认证' }
  return map[purpose] || purpose
}

function formatDate(date) {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

function openCreateKeyModal() {
  newKeyForm.name = ''
  newKeyForm.type = 'symmetric'
  newKeyForm.purpose = 'encryption'
  newKeyForm.keySize = '256'
  createKeyVisible.value = true
}

function openKeyDetail(key) {
  currentKey.value = key
  keyDetailVisible.value = true
}

async function handleCreateKey() {
  if (!newKeyForm.name) {
    Message.warning('请输入密钥名称')
    return
  }
  creatingKey.value = true
  try {
    await createKey(newKeyForm)
    Message.success('创建成功')
    createKeyVisible.value = false
    loadKeys()
  } catch (e) {
    Message.error('创建失败')
  } finally {
    creatingKey.value = false
  }
}

async function toggleKeyStatus(key) {
  try {
    await toggleKeyStatusApi(key.id)
    Message.success('更新成功')
    loadKeys()
  } catch (e) {
    Message.error('更新失败')
  }
}

async function deleteKey(key) {
  try {
    await deleteKeyApi(key.id)
    Message.success('删除成功')
    loadKeys()
  } catch (e) {
    Message.error('删除失败')
  }
}

async function saveEncryption() {
  saving.encryption = true
  try {
    await updateEncryptionConfig(encryptionConfig)
    Message.success('保存成功')
  } catch (e) {
    Message.error('保存失败')
  } finally {
    saving.encryption = false
  }
}

async function saveRotation() {
  saving.rotation = true
  try {
    await updateRotationConfig(rotationConfig)
    Message.success('保存成功')
    loadRotationConfig()
  } catch (e) {
    Message.error('保存失败')
  } finally {
    saving.rotation = false
  }
}

async function rotateNow() {
  rotating.value = true
  try {
    await rotateKeys()
    Message.success('轮换成功')
    loadRotationConfig()
    loadKeys()
  } catch (e) {
    Message.error('轮换失败')
  } finally {
    rotating.value = false
  }
}
</script>

<style scoped>
.page-container {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.section-card {
  flex-shrink: 0;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toolbar-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.toolbar-right {
  display: flex;
  gap: 8px;
}

.config-hint {
  margin-left: 8px;
  color: var(--color-text-3);
}

.info-text {
  color: var(--color-text-2);
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-label {
  color: var(--color-text-3);
  font-size: 13px;
}

.detail-value {
  font-size: 13px;
  word-break: break-all;
}

.mono {
  font-family: monospace;
  font-size: 12px;
}

.ml-3 {
  margin-left: 12px;
}
</style>
