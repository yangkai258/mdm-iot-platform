<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>高级功能</a-breadcrumb-item>
      <a-breadcrumb-item>老人陪伴模式</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 -->
    <div class="pro-page-header">
      <h2 class="pro-page-title">老人陪伴模式</h2>
      <p class="pro-page-desc">为老年用户配置简化界面、主动问候与紧急求助功能</p>
    </div>

    <!-- 搜索筛选区 -->
    <div class="pro-search-bar">
      <a-space>
        <a-select v-model="elderFilter" placeholder="选择老人账号" allow-clear style="width: 200px" @change="loadElderModes">
          <a-option v-for="e in elders" :key="e.id" :value="e.id">{{ e.name }}</a-option>
        </a-select>
        <a-select v-model="statusFilter" placeholder="模式状态" allow-clear style="width: 140px" @change="loadElderModes">
          <a-option value="enabled">已启用</a-option>
          <a-option value="disabled">已禁用</a-option>
        </a-select>
      </a-space>
    </div>

    <!-- 操作按钮区 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showConfigModal(null)">
          <template #icon><icon-settings /></template>
          新增配置
        </a-button>
        <a-button @click="loadElderModes">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
      </a-space>
    </div>

    <!-- 数据内容区 -->
    <div class="pro-content-area">
      <a-table
        :columns="columns"
        :data="elderModes"
        :loading="loading"
        :pagination="pagination"
        @page-change="onPageChange"
        row-key="id"
      >
        <template #elder_name="{ record }">
          <a-avatar :style="{ backgroundColor: '#8e4e9c' }" :size="32">
            {{ record.elder_name?.charAt(0) || '?' }}
          </a-avatar>
          <span style="margin-left: 8px">{{ record.elder_name }}</span>
        </template>
        <template #enabled="{ record }">
          <a-switch v-model="record.enabled" @change="toggleMode(record)" :disabled="savingId === record.id" />
        </template>
        <template #simplified_ui="{ record }">
          <a-tag :color="record.simplified_ui ? 'green' : 'gray'">
            {{ record.simplified_ui ? '已启用' : '已禁用' }}
          </a-tag>
        </template>
        <template #font_size="{ record }">
          {{ getFontLabel(record.font_size) }}
        </template>
        <template #high_contrast="{ record }">
          <a-tag :color="record.high_contrast ? 'green' : 'gray'">
            {{ record.high_contrast ? '已启用' : '已禁用' }}
          </a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showConfigModal(record)">配置</a-button>
            <a-button type="text" size="small" @click="showHealthData(record)">健康数据</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 健康数据抽屉 -->
    <a-drawer v-model:visible="healthDrawerVisible" :width="720" :title="`健康数据 - ${healthElderName}`" @close="healthDrawerVisible = false">
      <a-row :gutter="12" style="margin-bottom: 16px">
        <a-col :span="8">
          <a-card class="stat-card" hoverable>
            <a-statistic :value="healthData.heart_rate" :precision="0" suffix="bpm">
              <template #prefix><icon-heart :size="20" style="color:#f53f3f"/></template>
              <template #title>心率</template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :span="8">
          <a-card class="stat-card" hoverable>
            <a-statistic :value="healthData.blood_pressure_systolic" suffix="/" :precision="0" :value-from="0">
              <template #prefix><icon-activity :size="20" style="color:#f53f3f"/></template>
              <template #title>血压</template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :span="8">
          <a-card class="stat-card" hoverable>
            <a-statistic :value="healthData.sleep_hours" :precision="1" suffix="小时">
              <template #prefix><icon-moon :size="20" style="color:#722ed1"/></template>
              <template #title>睡眠时长</template>
            </a-statistic>
          </a-card>
        </a-col>
      </a-row>

      <a-row :gutter="12" style="margin-bottom: 16px">
        <a-col :span="8">
          <a-card class="stat-card" hoverable>
            <a-statistic :value="healthData.step_count" :precision="0" suffix="步">
              <template #prefix><icon-fire :size="20" style="color:#ff6700"/></template>
              <template #title>日行走步数</template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :span="8">
          <a-card class="stat-card" hoverable>
            <a-statistic :value="healthData.calories_burned" :precision="0" suffix="kcal">
              <template #prefix><icon-thumb-up :size="20" style="color:#0fc6c2"/></template>
              <template #title>消耗卡路里</template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :span="8">
          <a-card class="stat-card" hoverable>
            <a-statistic :value="healthData.oxygen_saturation" :precision="0" suffix="%">
              <template #prefix><icon-air :size="20" style="color:#1659f5"/></template>
              <template #title>血氧饱和度</template>
            </a-statistic>
          </a-card>
        </a-col>
      </a-row>

      <a-card title="近7天活动趋势">
        <div ref="healthChartRef" style="height: 200px"></div>
      </a-card>
    </a-drawer>

    <!-- 配置弹窗 -->
    <a-modal v-model:visible="configModalVisible" :title="isEdit ? '编辑老人陪伴模式配置' : '新增老人陪伴模式配置'" @ok="handleSave" :width="600" @close="resetForm">
      <a-form :model="form" layout="vertical">

        <a-form-item label="老人账号" required>
          <a-select v-model="form.elder_id" placeholder="请选择老人账号" :disabled="isEdit">
            <a-option v-for="e in elders" :key="e.id" :value="e.id">{{ e.name }} ({{ e.phone }})</a-option>
          </a-select>
        </a-form-item>

        <a-divider>简化界面设置</a-divider>

        <a-form-item label="启用简化界面">
          <a-switch v-model="form.simplified_ui" />
        </a-form-item>
        <a-form-item label="字体大小" v-if="form.simplified_ui">
          <a-radio-group v-model="form.font_size">
            <a-radio value="large">大（18px）</a-radio>
            <a-radio value="extra_large">特大（22px）</a-radio>
            <a-radio value="jumbo">超大（26px）</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="高对比度模式" v-if="form.simplified_ui">
          <a-switch v-model="form.high_contrast" />
        </a-form-item>

        <a-divider>主动问候设置</a-divider>

        <a-form-item label="启用主动问候">
          <a-switch v-model="form.greeting_enabled" />
        </a-form-item>
        <a-form-item label="问候时间" v-if="form.greeting_enabled">
          <a-space direction="vertical">
            <a-time-picker-range v-model="form.greeting_time_range" format="HH:mm" style="width: 280px" placeholder="选择问候时段" />
          </a-space>
        </a-form-item>
        <a-form-item label="问候语风格" v-if="form.greeting_enabled">
          <a-select v-model="form.greeting_style" placeholder="选择问候语风格">
            <a-option value="warm">温馨关怀型</a-option>
            <a-option value="humorous">轻松幽默型</a-option>
            <a-option value="formal">正式礼貌型</a-option>
          </a-select>
        </a-form-item>

        <a-divider>紧急求助配置</a-divider>

        <a-form-item label="启用紧急求助">
          <a-switch v-model="form.emergency_enabled" />
        </a-form-item>
        <a-form-item label="紧急联系电话" v-if="form.emergency_enabled">
          <a-input v-model="form.emergency_phone" placeholder="请输入紧急联系电话" style="width: 240px" />
        </a-form-item>
        <a-form-item label="紧急联系人关系" v-if="form.emergency_enabled">
          <a-select v-model="form.emergency_relation" placeholder="选择关系" style="width: 160px">
            <a-option value="son">儿子</a-option>
            <a-option value="daughter">女儿</a-option>
            <a-option value="spouse">配偶</a-option>
            <a-option value="caregiver">护工</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="SOS 触发方式" v-if="form.emergency_enabled">
          <a-checkbox-group v-model="form.emergency_trigger">
            <a-checkbox value="button">物理按钮</a-checkbox>
            <a-checkbox value="voice">语音指令</a-checkbox>
            <a-checkbox value="gesture">手势识别</a-checkbox>
          </a-checkbox-group>
        </a-form-item>

      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const elderModes = ref<any[]>([])
const elders = ref<any[]>([])
const loading = ref(false)
const savingId = ref<number | null>(null)
const elderFilter = ref('')
const statusFilter = ref('')
const configModalVisible = ref(false)
const isEdit = ref(false)
const healthDrawerVisible = ref(false)
const healthElderName = ref('')
const healthData = reactive({
  heart_rate: 0,
  blood_pressure_systolic: 0,
  blood_pressure_diastolic: 0,
  sleep_hours: 0,
  step_count: 0,
  calories_burned: 0,
  oxygen_saturation: 0
})

const form = reactive({
  id: null as number | null,
  elder_id: null as number | null,
  simplified_ui: true,
  font_size: 'large',
  high_contrast: false,
  greeting_enabled: true,
  greeting_time_range: null as any,
  greeting_style: 'warm',
  emergency_enabled: true,
  emergency_phone: '',
  emergency_relation: '',
  emergency_trigger: [] as string[]
})

const columns = [
  { title: '老人账号', dataIndex: 'elder_name', slotName: 'elder_name', width: 200 },
  { title: '模式开关', dataIndex: 'enabled', slotName: 'enabled', width: 120 },
  { title: '简化界面', dataIndex: 'simplified_ui', slotName: 'simplified_ui', width: 120 },
  { title: '字体大小', dataIndex: 'font_size', slotName: 'font_size', width: 120 },
  { title: '高对比度', dataIndex: 'high_contrast', slotName: 'high_contrast', width: 120 },
  { title: '操作', slotName: 'actions', width: 280 }
]

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

function getFontLabel(size: string) {
  const map: Record<string, string> = { large: '大', extra_large: '特大', jumbo: '超大' }
  return map[size] || size
}

async function loadElders() {
  try {
    const res = await fetch('/api/v1/family/members?role=elder', { credentials: 'include' })
    const data = await res.json()
    elders.value = data.data?.list || data.data?.members || []
  } catch { /* ignore */ }
}

async function loadElderModes() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (elderFilter.value) params.append('elder_id', elderFilter.value)
    if (statusFilter.value) params.append('status', statusFilter.value)
    params.append('page', String(pagination.current))
    params.append('page_size', String(pagination.pageSize))

    const res = await fetch(`/api/v1/advanced/elder-mode?${params}`, { credentials: 'include' })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      elderModes.value = data.data?.list || data.data || []
      pagination.total = data.data?.total || 0
    } else {
      Message.error(data.message || '加载失败')
    }
  } catch {
    Message.error('网络错误')
  } finally {
    loading.value = false
  }
}

function onPageChange(page: number) {
  pagination.current = page
  loadElderModes()
}

function showConfigModal(record: any) {
  if (record) {
    isEdit.value = true
    form.id = record.id
    form.elder_id = record.elder_id
    form.simplified_ui = record.simplified_ui ?? true
    form.font_size = record.font_size || 'large'
    form.high_contrast = record.high_contrast ?? false
    form.greeting_enabled = record.greeting_enabled ?? true
    form.greeting_time_range = record.greeting_time_range
    form.greeting_style = record.greeting_style || 'warm'
    form.emergency_enabled = record.emergency_enabled ?? true
    form.emergency_phone = record.emergency_phone || ''
    form.emergency_relation = record.emergency_relation || ''
    form.emergency_trigger = record.emergency_trigger || []
  } else {
    isEdit.value = false
    resetForm()
  }
  configModalVisible.value = true
}

function resetForm() {
  form.id = null
  form.elder_id = null
  form.simplified_ui = true
  form.font_size = 'large'
  form.high_contrast = false
  form.greeting_enabled = true
  form.greeting_time_range = null
  form.greeting_style = 'warm'
  form.emergency_enabled = true
  form.emergency_phone = ''
  form.emergency_relation = ''
  form.emergency_trigger = []
}

async function handleSave() {
  if (!isEdit.value && !form.elder_id) {
    Message.warning('请选择老人账号')
    return
  }
  try {
    const payload: any = { ...form }
    if (payload.greeting_time_range && payload.greeting_time_range.length === 2) {
      payload.greeting_start_time = payload.greeting_time_range[0]
      payload.greeting_end_time = payload.greeting_time_range[1]
    }
    delete payload.greeting_time_range

    const url = isEdit.value ? `/api/v1/advanced/elder-mode/${form.id}` : '/api/v1/advanced/elder-mode'
    const method = isEdit.value ? 'PUT' : 'POST'
    const res = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(payload)
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      Message.success(isEdit.value ? '配置已更新' : '配置已保存')
      configModalVisible.value = false
      loadElderModes()
    } else {
      Message.error(data.message || '保存失败')
    }
  } catch {
    Message.error('网络错误')
  }
}

async function toggleMode(record: any) {
  savingId.value = record.id
  try {
    const res = await fetch(`/api/v1/advanced/elder-mode/${record.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ enabled: record.enabled })
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      Message.success(record.enabled ? '已启用' : '已禁用')
    } else {
      record.enabled = !record.enabled
      Message.error(data.message || '操作失败')
    }
  } catch {
    record.enabled = !record.enabled
    Message.error('网络错误')
  } finally {
    savingId.value = null
  }
}

async function handleDelete(record: any) {
  try {
    const res = await fetch(`/api/v1/advanced/elder-mode/${record.id}`, {
      method: 'DELETE',
      credentials: 'include'
    })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      Message.success('已删除')
      loadElderModes()
    } else {
      Message.error(data.message || '删除失败')
    }
  } catch {
    Message.error('网络错误')
  }
}

async function showHealthData(record: any) {
  healthElderName.value = record.elder_name
  healthDrawerVisible.value = true
  try {
    const params = new URLSearchParams({ elder_id: String(record.elder_id || record.id) })
    const res = await fetch(`/api/v1/advanced/elder-mode/health-data?${params}`, { credentials: 'include' })
    const data = await res.json()
    if (data.code === 0 || data.code === 200) {
      Object.assign(healthData, data.data || {})
    }
  } catch { /* ignore */ }
}

onMounted(() => {
  loadElders()
  loadElderModes()
})
</script>
