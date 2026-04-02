<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="pro-page-container">
    <!-- 闈㈠寘灞?-->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>棣栭〉</a-breadcrumb-item>
      <a-breadcrumb-item>楂樼骇鍔熻兘</a-breadcrumb-item>
      <a-breadcrumb-item>鍎跨妯″紡</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 椤甸潰鏍囬 -->
    <div class="pro-page-header">
      <h2 class="pro-page-title">鍎跨妯″紡</h2>
      <p class="pro-page-desc">閰嶇疆瀹夊叏鍐呭杩囨护銆佷娇鐢ㄦ椂闀块檺鍒朵笌浣跨敤鎶ュ憡鍒嗘瀽</p>
    </div>

    <!-- 鎼滅储绛涢€夊尯 -->
    <div class="pro-search-bar">
      <a-space>
        <a-select v-model="childFilter" placeholder="閫夋嫨鍎跨璐﹀彿" allow-clear style="width: 200px" @change="loadChildModes">
          <a-option v-for="c in children" :key="c.id" :value="c.id">{{ c.name }}</a-option>
        </a-select>
        <a-select v-model="statusFilter" placeholder="妯″紡鐘舵€? allow-clear style="width: 140px" @change="loadChildModes">
          <a-option value="enabled">宸插惎鐢?/a-option>
          <a-option value="disabled">宸茬鐢?/a-option>
        </a-select>
      </a-space>
    </div>

    <!-- 鎿嶄綔鎸夐挳鍖?-->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showConfigModal(null)">
          <template #icon><icon-settings /></template>
          鏂板閰嶇疆
        </a-button>
        <a-button @click="loadChildModes">
          <template #icon><icon-refresh /></template>
          鍒锋柊
        </a-button>
      </a-space>
    </div>

    <!-- 鏁版嵁鍐呭鍖?-->
    <div class="pro-content-area">
      <a-table
        :columns="columns"
        :data="childModes"
        :loading="loading"
        :pagination="pagination"
        @page-change="onPageChange"
        row-key="id"
      >
        <template #child_name="{ record }">
          <a-avatar :style="{ backgroundColor: '#ff7d00' }" :size="32">
            {{ record.child_name?.charAt(0) || '?' }}
          </a-avatar>
          <span style="margin-left: 8px">{{ record.child_name }}</span>
        </template>
        <template #enabled="{ record }">
          <a-switch v-model="record.enabled" @change="toggleMode(record)" :disabled="savingId === record.id" />
        </template>
        <template #content_filter="{ record }">
          <a-tag :color="record.content_filter_enabled ? 'green' : 'gray'">
            {{ record.content_filter_enabled ? '宸插惎鐢? : '宸茬鐢? }}
          </a-tag>
        </template>
        <template #time_limit="{ record }">
          <span v-if="record.time_limit_enabled">
            {{ record.daily_time_limit }}鍒嗛挓/澶?          </span>
          <a-tag v-else color="gray">鏈檺鍒?/a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showConfigModal(record)">閰嶇疆</a-button>
            <a-button type="text" size="small" @click="showUsageReport(record)">鎶ュ憡</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">鍒犻櫎</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 浣跨敤鎶ュ憡鎶藉眽 -->
    <a-drawer v-model:visible="reportDrawerVisible" :width="680" :title="`浣跨敤鎶ュ憡 - ${reportChildName}`" @close="reportDrawerVisible = false">
      <div class="report-section">
        <div class="report-header">
          <a-radio-group v-model="reportPeriod" type="button" @change="loadUsageReport">
            <a-radio value="day">浠婃棩</a-radio>
            <a-radio value="week">鏈懆</a-radio>
            <a-radio value="month">鏈湀</a-radio>
          </a-radio-group>
        </div>

        <!-- 浣跨敤鏃堕暱缁熻鍗＄墖 -->
        <a-row :gutter="12" style="margin-bottom: 16px">
          <a-col :span="8">
            <a-card class="stat-card" hoverable>
              <a-statistic :value="reportData.total_minutes" :precision="0" suffix="鍒嗛挓">
                <template #prefix><icon-clock-circle :size="20" style="color:#1659f5"/></template>
                <template #title>鎬讳娇鐢ㄦ椂闀?/template>
              </a-statistic>
            </a-card>
          </a-col>
          <a-col :span="8">
            <a-card class="stat-card" hoverable>
              <a-statistic :value="reportData.avg_daily_minutes" :precision="0" suffix="鍒嗛挓">
                <template #prefix><icon-history :size="20" style="color:#0fc6c2"/></template>
                <template #title>鏃ュ潎鏃堕暱</template>
              </a-statistic>
            </a-card>
          </a-col>
          <a-col :span="8">
            <a-card class="stat-card" hoverable>
              <a-statistic :value="reportData.limit_compliance_rate" :precision="0" suffix="%">
                <template #prefix><icon-check-circle :size="20" style="color:#00b42a"/></template>
                <template #title>闄愭椂閬靛畧鐜?/template>
              </a-statistic>
            </a-card>
          </a-col>
        </a-row>

        <!-- 浣跨敤瓒嬪娍鍥?-->
        <a-card title="姣忔棩浣跨敤鏃堕暱瓒嬪娍" style="margin-bottom: 16px">
          <div ref="usageChartRef" style="height: 200px"></div>
        </a-card>

        <!-- 搴旂敤浣跨敤鍒嗗竷 -->
        <a-card title="鍐呭鍒嗙被浣跨敤鍒嗗竷">
          <div ref="categoryChartRef" style="height: 200px"></div>
        </a-card>
      </div>
    </a-drawer>

    <!-- 閰嶇疆寮圭獥 -->
    <a-modal v-model:visible="configModalVisible" :title="isEdit ? '缂栬緫鍎跨妯″紡閰嶇疆' : '鏂板鍎跨妯″紡閰嶇疆'" @ok="handleSave" :width="560" @close="resetForm">
      <a-form :model="form" layout="vertical">
        <a-form-item label="鍎跨璐﹀彿" required>
          <a-select v-model="form.child_id" placeholder="璇烽€夋嫨鍎跨璐﹀彿" :disabled="isEdit">
            <a-option v-for="c in children" :key="c.id" :value="c.id">{{ c.name }} ({{ c.phone }})</a-option>
          </a-select>
        </a-form-item>

        <a-divider>鍐呭杩囨护</a-divider>

        <a-form-item label="鍚敤鍐呭杩囨护">
          <a-switch v-model="form.content_filter_enabled" />
        </a-form-item>
        <a-form-item label="杩囨护绾у埆" v-if="form.content_filter_enabled">
          <a-radio-group v-model="form.filter_level">
            <a-radio value="strict">涓ユ牸锛堥€傚悎浣庨緞鍎跨锛?/a-radio>
            <a-radio value="moderate">涓瓑锛堥€傚悎瀛﹂緞鍎跨锛?/a-radio>
            <a-radio value="light">杞诲害锛堥€傚悎闈掑皯骞达級</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="绂佹璁块棶鐨勫垎绫? v-if="form.content_filter_enabled">
          <a-checkbox-group v-model="form.blocked_categories">
            <a-checkbox value="violence">鏆村姏鍐呭</a-checkbox>
            <a-checkbox value="adult">鎴愪汉鍐呭</a-checkbox>
            <a-checkbox value="gambling">璧屽崥</a-checkbox>
            <a-checkbox value="social">绀句氦濯掍綋</a-checkbox>
            <a-checkbox value="games">娓告垙</a-checkbox>
          </a-checkbox-group>
        </a-form-item>

        <a-divider>浣跨敤鏃堕棿闄愬埗</a-divider>

        <a-form-item label="鍚敤鏃堕棿闄愬埗">
          <a-switch v-model="form.time_limit_enabled" />
        </a-form-item>
        <a-form-item label="姣忔棩鏃堕暱闄愬埗锛堝垎閽燂級" v-if="form.time_limit_enabled">
          <a-input-number v-model="form.daily_time_limit" :min="15" :max="480" :step="15" style="width: 200px" />
        </a-form-item>
        <a-form-item label="鍏佽浣跨敤鐨勬椂闂存" v-if="form.time_limit_enabled">
          <a-time-picker-range v-model="form.allowed_time_range" format="HH:mm" style="width: 300px" />
        </a-form-item>
        <a-form-item label="绂佺敤鏃ユ湡" v-if="form.time_limit_enabled">
          <a-checkbox-group v-model="form.disabled_days">
            <a-checkbox :value="1">鍛ㄤ竴</a-checkbox>
            <a-checkbox :value="2">鍛ㄤ簩</a-checkbox>
            <a-checkbox :value="3">鍛ㄤ笁</a-checkbox>
            <a-checkbox :value="4">鍛ㄥ洓</a-checkbox>
            <a-checkbox :value="5">鍛ㄤ簲</a-checkbox>
            <a-checkbox :value="6">鍛ㄥ叚</a-checkbox>
            <a-checkbox :value="0">鍛ㄦ棩</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const childModes = ref<any[]>([])
const children = ref<any[]>([])
const loading = ref(false)
const savingId = ref<number | null>(null)
const childFilter = ref('')
const statusFilter = ref('')
const configModalVisible = ref(false)
const isEdit = ref(false)
const reportDrawerVisible = ref(false)
const reportChildName = ref('')
const reportPeriod = ref('week')
const reportData = reactive({
  total_minutes: 0,
  avg_daily_minutes: 0,
  limit_compliance_rate: 0,
  daily_trend: [] as any[],
  category_usage: [] as any[]
})

const form = reactive({
  id: null as number | null,
  child_id: null as number | null,
  content_filter_enabled: false,
  filter_level: 'moderate',
  blocked_categories: [] as string[],
  time_limit_enabled: false,
  daily_time_limit: 60,
  allowed_time_range: null as any,
  disabled_days: [] as number[]
})

const columns = [
  { title: '鍎跨璐﹀彿', dataIndex: 'child_name', slotName: 'child_name', width: 200 },
  { title: '妯″紡寮€鍏?, dataIndex: 'enabled', slotName: 'enabled', width: 120 },
  { title: '鍐呭杩囨护', dataIndex: 'content_filter', slotName: 'content_filter', width: 120 },
  { title: '鏃堕棿闄愬埗', dataIndex: 'time_limit', slotName: 'time_limit', width: 160 },
  { title: '鎿嶄綔', slotName: 'actions', width: 220 }
]

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

async function loadChildren() {
  try {
    const res = await fetch('/api/family/members?role=child', { credentials: 'include' })
    const data = await res.json()
    children.value = data.data?.list || data.data?.members || []
  } catch { /* ignore */ }
}

async function loadChildModes() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (childFilter.value) params.append('child_id', childFilter.value)
    if (statusFilter.value) params.append('status', statusFilter.value)
    params.append('page', String(pagination.current))
    params.append('page_size', String(pagination.pageSize))

    const res = await fetch(`/api/advanced/child-mode?${params}`, { credentials: 'include' })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      childModes.value = data.data?.list || data.data || []
      pagination.total = data.data?.total || 0
    } else {
      Message.error(data.message || '鍔犺浇澶辫触')
    }
  } catch {
    Message.error('缃戠粶閿欒')
  } finally {
    loading.value = false
  }
}

function onPageChange(page: number) {
  pagination.current = page
  loadChildModes()
}

function showConfigModal(record: any) {
  if (record) {
    isEdit.value = true
    form.id = record.id
    form.child_id = record.child_id
    form.content_filter_enabled = record.content_filter_enabled
    form.filter_level = record.filter_level || 'moderate'
    form.blocked_categories = record.blocked_categories || []
    form.time_limit_enabled = record.time_limit_enabled
    form.daily_time_limit = record.daily_time_limit || 60
    form.allowed_time_range = record.allowed_time_range
    form.disabled_days = record.disabled_days || []
  } else {
    isEdit.value = false
    resetForm()
  }
  configModalVisible.value = true
}

function resetForm() {
  form.id = null
  form.child_id = null
  form.content_filter_enabled = false
  form.filter_level = 'moderate'
  form.blocked_categories = []
  form.time_limit_enabled = false
  form.daily_time_limit = 60
  form.allowed_time_range = null
  form.disabled_days = []
}

async function handleSave() {
  if (!isEdit.value && !form.child_id) {
    Message.warning('璇烽€夋嫨鍎跨璐﹀彿')
    return
  }
  try {
    const payload = { ...form }
    if (payload.allowed_time_range && payload.allowed_time_range.length === 2) {
      payload.allowed_start_time = payload.allowed_time_range[0]
      payload.allowed_end_time = payload.allowed_time_range[1]
    }
    delete payload.allowed_time_range

    const url = isEdit.value ? `/api/advanced/child-mode/${form.id}` : '/api/advanced/child-mode'
    const method = isEdit.value ? 'PUT' : 'POST'
    const res = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(payload)
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      Message.success(isEdit.value ? '閰嶇疆宸叉洿鏂? : '閰嶇疆宸蹭繚瀛?)
      configModalVisible.value = false
      loadChildModes()
    } else {
      Message.error(data.message || '淇濆瓨澶辫触')
    }
  } catch {
    Message.error('缃戠粶閿欒')
  }
}

async function toggleMode(record: any) {
  savingId.value = record.id
  try {
    const res = await fetch(`/api/advanced/child-mode/${record.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ enabled: record.enabled })
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      Message.success(record.enabled ? '宸插惎鐢? : '宸茬鐢?)
    } else {
      record.enabled = !record.enabled
      Message.error(data.message || '鎿嶄綔澶辫触')
    }
  } catch {
    record.enabled = !record.enabled
    Message.error('缃戠粶閿欒')
  } finally {
    savingId.value = null
  }
}

async function handleDelete(record: any) {
  try {
    const res = await fetch(`/api/advanced/child-mode/${record.id}`, {
      method: 'DELETE',
      credentials: 'include'
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      Message.success('宸插垹闄?)
      loadChildModes()
    } else {
      Message.error(data.message || '鍒犻櫎澶辫触')
    }
  } catch {
    Message.error('缃戠粶閿欒')
  }
}

async function showUsageReport(record: any) {
  reportChildName.value = record.child_name
  reportDrawerVisible.value = true
  await loadUsageReport(record.id)
}

async function loadUsageReport(childId?: number) {
  try {
    const id = childId || childFilter.value
    const params = new URLSearchParams({ period: reportPeriod.value })
    if (id) params.append('child_id', String(id))
    const res = await fetch(`/api/advanced/child-mode/usage-report?${params}`, { credentials: 'include' })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      Object.assign(reportData, data.data || {})
    }
  } catch { /* ignore */ }
}

onMounted(() => {
  loadChildren()
  loadChildModes()
})
</script>

