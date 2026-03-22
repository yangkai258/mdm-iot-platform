<template>
  <div class="page-container">
    <!-- 顶部：标题 + 操作按钮 -->
    <a-card class="header-card">
      <div class="header-row">
        <div class="page-title">区域设置</div>
        <a-button type="primary" @click="openAddModal">「添加区域」</a-button>
      </div>
    </a-card>

    <!-- 筛选栏 -->
    <a-card class="filter-card">
      <div class="filter-row">
        <a-select
          v-model="filters.status"
          placeholder="状态"
          style="width: 140px"
          allow-clear
          @change="handleFilterChange"
        >
          <a-option value="active">活跃</a-option>
          <a-option value="inactive">停用</a-option>
        </a-select>
        <a-select
          v-model="filters.locale"
          placeholder="默认语言"
          style="width: 160px"
          allow-clear
          @change="handleFilterChange"
        >
          <a-option v-for="lang in availableLocales" :key="lang" :value="lang">{{ localeLabel(lang) }}</a-option>
        </a-select>
        <a-input
          v-model="filters.keyword"
          placeholder="搜索区域名称..."
          style="width: 200px"
          allow-clear
          @change="handleFilterChange"
          @press-enter="handleFilterChange"
        />
        <a-button @click="handleFilterChange">「查询」</a-button>
        <a-button @click="resetFilters">「重置」</a-button>
      </div>
    </a-card>

    <!-- 区域列表 -->
    <a-card class="table-card">
      <a-table
        :columns="columns"
        :data="filteredData"
        :loading="loading"
        :pagination="{ pageSize: 20, showTotal: true }"
        :scroll="{ x: 1000 }"
        row-key="id"
        @page-change="handlePageChange"
      >
        <template #regionName="{ record }">
          <div class="region-name-cell">
            <span class="region-name">{{ record.region_name }}</span>
          </div>
        </template>
        <template #regionKey="{ record }">
          <span class="mono">{{ record.region_key }}</span>
        </template>
        <template #localeDefault="{ record }">
          <a-tag>{{ localeLabel(record.locale_default) }}</a-tag>
        </template>
        <template #localeSupported="{ record }">
          <div class="locale-tags">
            <a-tag v-for="loc in (record.locale_supported || []).slice(0, 3)" :key="loc" size="small">
              {{ localeShort(loc) }}
            </a-tag>
            <a-tag v-if="(record.locale_supported || []).length > 3" size="small" color="arcoblue">
              +{{ record.locale_supported.length - 3 }}
            </a-tag>
          </div>
        </template>
        <template #timezone="{ record }">
          <span class="mono tz-text">{{ record.timezone }}</span>
        </template>
        <template #status="{ record }">
          <a-badge
            :color="record.is_active ? 'green' : 'gray'"
            :text="record.is_active ? '活跃' : '停用'"
          />
        </template>
        <template #updatedAt="{ record }">
          <span class="date-text">{{ formatDate(record.updated_at) }}</span>
        </template>
        <template #actions="{ record }">
          <div class="action-btns">
            <a-button type="text" size="small" @click="openEditModal(record)">「编辑」</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">「删除」</a-button>
          </div>
        </template>
      </a-table>
    </a-card>

    <!-- 添加/编辑区域弹窗 -->
    <a-modal
      v-model:visible="formVisible"
      :title="editingItem ? '编辑区域' : '添加区域'"
      :width="600"
      @before-ok="handleSubmit"
      @cancel="formVisible = false"
    >
      <a-form :model="form" layout="vertical">
        <div class="form-row">
          <a-form-item label="区域标识" required class="half">
            <a-input v-model="form.region_key" placeholder="如：cn-east" :disabled="!!editingItem" />
          </a-form-item>
          <a-form-item label="区域名称" required class="half">
            <a-input v-model="form.region_name" placeholder="如：中国东部" />
          </a-form-item>
        </div>

        <div class="form-row">
          <a-form-item label="默认语言" required class="half">
            <a-select v-model="form.locale_default" placeholder="选择默认语言">
              <a-option v-for="lang in availableLocales" :key="lang" :value="lang">{{ localeLabel(lang) }}</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="支持语言" class="half">
            <a-select v-model="form.locale_supported" multiple placeholder="选择支持的语言">
              <a-option v-for="lang in availableLocales" :key="lang" :value="lang">{{ localeLabel(lang) }}</a-option>
            </a-select>
          </a-form-item>
        </div>

        <div class="form-row">
          <a-form-item label="时区" required class="third">
            <a-select v-model="form.timezone" placeholder="选择时区" show-search>
              <a-option v-for="tz in timezones" :key="tz.value" :value="tz.value">{{ tz.label }}</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="日期格式" class="third">
            <a-select v-model="form.date_format" placeholder="日期格式">
              <a-option value="YYYY-MM-DD">YYYY-MM-DD</a-option>
              <a-option value="DD/MM/YYYY">DD/MM/YYYY</a-option>
              <a-option value="MM/DD/YYYY">MM/DD/YYYY</a-option>
              <a-option value="YYYY/MM/DD">YYYY/MM/DD</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="时间格式" class="third">
            <a-select v-model="form.time_format" placeholder="时间格式">
              <a-option value="HH:mm:ss">24小时 (HH:mm:ss)</a-option>
              <a-option value="hh:mm:ss A">12小时 (hh:mm:ss A)</a-option>
              <a-option value="HH:mm">24小时 (HH:mm)</a-option>
            </a-select>
          </a-form-item>
        </div>

        <div class="form-row">
          <a-form-item label="数字格式" class="half">
            <a-select v-model="form.number_format" placeholder="数字格式">
              <a-option value="1,234.56">1,234.56 (1,000.00)</a-option>
              <a-option value="1.234,56">1.234,56 (1.000,00)</a-option>
              <a-option value="1 234.56">1 234,56 (法国)</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="货币代码" class="half">
            <a-select v-model="form.currency_code" placeholder="选择货币">
              <a-option value="CNY">CNY (人民币)</a-option>
              <a-option value="USD">USD (美元)</a-option>
              <a-option value="EUR">EUR (欧元)</a-option>
              <a-option value="JPY">JPY (日元)</a-option>
              <a-option value="GBP">GBP (英镑)</a-option>
              <a-option value="KRW">KRW (韩元)</a-option>
            </a-select>
          </a-form-item>
        </div>

        <a-form-item label="启用状态">
          <a-switch v-model="form.is_active" />
          <span class="switch-label">{{ form.is_active ? '区域已启用' : '区域已停用' }}</span>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import {
  getRegionSettings,
  createRegionSetting,
  updateRegionSetting,
  deleteRegionSetting
} from '@/api/i18n'
import dayjs from 'dayjs'

const loading = ref(false)
const regionSettings = ref([])

const filters = reactive({
  status: '',
  locale: '',
  keyword: ''
})

const editingItem = ref(null)
const formVisible = ref(false)

const form = reactive({
  region_key: '',
  region_name: '',
  locale_default: 'zh-CN',
  locale_supported: [] as string[],
  timezone: 'Asia/Shanghai',
  date_format: 'YYYY-MM-DD',
  time_format: 'HH:mm:ss',
  number_format: '1,234.56',
  currency_code: 'CNY',
  is_active: true
})

const availableLocales = [
  'zh-CN', 'en-US', 'ja-JP', 'ko-KR', 'fr-FR', 'de-DE', 'es-ES', 'ru-RU'
]

const timezones = [
  { value: 'Asia/Shanghai', label: 'Asia/Shanghai (UTC+8)' },
  { value: 'Asia/Tokyo', label: 'Asia/Tokyo (UTC+9)' },
  { value: 'Asia/Seoul', label: 'Asia/Seoul (UTC+9)' },
  { value: 'America/New_York', label: 'America/New_York (UTC-5)' },
  { value: 'America/Los_Angeles', label: 'America/Los_Angeles (UTC-8)' },
  { value: 'Europe/London', label: 'Europe/London (UTC+0)' },
  { value: 'Europe/Paris', label: 'Europe/Paris (UTC+1)' },
  { value: 'Europe/Berlin', label: 'Europe/Berlin (UTC+1)' },
  { value: 'Europe/Moscow', label: 'Europe/Moscow (UTC+3)' },
  { value: 'Asia/Singapore', label: 'Asia/Singapore (UTC+8)' },
  { value: 'Australia/Sydney', label: 'Australia/Sydney (UTC+10)' }
]

const columns = [
  { title: '区域名称', slotName: 'regionName', width: 160 },
  { title: '区域标识', slotName: 'regionKey', width: 120 },
  { title: '默认语言', slotName: 'localeDefault', width: 120 },
  { title: '支持语言', slotName: 'localeSupported', width: 180 },
  { title: '时区', slotName: 'timezone', width: 160 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '更新时间', slotName: 'updatedAt', width: 160 },
  { title: '操作', slotName: 'actions', width: 140, fixed: 'right' }
]

const filteredData = computed(() => {
  return regionSettings.value.filter(r => {
    if (filters.status) {
      const isActive = filters.status === 'active'
      if (r.is_active !== isActive) return false
    }
    if (filters.locale && r.locale_default !== filters.locale) return false
    if (filters.keyword && !r.region_name.toLowerCase().includes(filters.keyword.toLowerCase())) return false
    return true
  })
})

function localeLabel(locale) {
  const map: Record<string, string> = {
    'zh-CN': '简体中文',
    'en-US': 'English',
    'ja-JP': '日本語',
    'ko-KR': '한국어',
    'fr-FR': 'Français',
    'de-DE': 'Deutsch',
    'es-ES': 'Español',
    'ru-RU': 'Русский'
  }
  return map[locale] || locale
}

function localeShort(locale) {
  const map: Record<string, string> = {
    'zh-CN': 'zh',
    'en-US': 'en',
    'ja-JP': 'ja',
    'ko-KR': 'ko',
    'fr-FR': 'fr',
    'de-DE': 'de',
    'es-ES': 'es',
    'ru-RU': 'ru'
  }
  return map[locale] || locale
}

function formatDate(date) {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

async function loadData() {
  loading.value = true
  try {
    const res = await getRegionSettings()
    regionSettings.value = res.data || res || []
  } catch (e) {
    console.error('加载区域设置失败', e)
    // 降级显示示例数据
    regionSettings.value = [
      {
        id: 1,
        region_key: 'cn-east',
        region_name: '中国东部',
        locale_default: 'zh-CN',
        locale_supported: ['zh-CN', 'en-US'],
        timezone: 'Asia/Shanghai',
        date_format: 'YYYY-MM-DD',
        time_format: 'HH:mm:ss',
        number_format: '1,234.56',
        currency_code: 'CNY',
        is_active: true,
        updated_at: new Date().toISOString()
      }
    ]
  } finally {
    loading.value = false
  }
}

function handleFilterChange() {
  // computed filteredData handles this
}

function resetFilters() {
  filters.status = ''
  filters.locale = ''
  filters.keyword = ''
}

function handlePageChange(page) {
  // pagination handled by Arco table
}

function openAddModal() {
  editingItem.value = null
  Object.assign(form, {
    region_key: '',
    region_name: '',
    locale_default: 'zh-CN',
    locale_supported: [],
    timezone: 'Asia/Shanghai',
    date_format: 'YYYY-MM-DD',
    time_format: 'HH:mm:ss',
    number_format: '1,234.56',
    currency_code: 'CNY',
    is_active: true
  })
  formVisible.value = true
}

function openEditModal(record) {
  editingItem.value = record
  Object.assign(form, {
    region_key: record.region_key,
    region_name: record.region_name,
    locale_default: record.locale_default,
    locale_supported: record.locale_supported ? [...record.locale_supported] : [],
    timezone: record.timezone,
    date_format: record.date_format || 'YYYY-MM-DD',
    time_format: record.time_format || 'HH:mm:ss',
    number_format: record.number_format || '1,234.56',
    currency_code: record.currency_code || 'CNY',
    is_active: record.is_active
  })
  formVisible.value = true
}

async function handleSubmit(done) {
  try {
    if (!form.region_key || !form.region_name || !form.locale_default || !form.timezone) {
      Message.error('请填写必填项')
      done(false)
      return
    }
    if (editingItem.value) {
      await updateRegionSetting(editingItem.value.id, form)
      Message.success('更新成功')
    } else {
      await createRegionSetting(form)
      Message.success('创建成功')
    }
    formVisible.value = false
    await loadData()
  } catch (e) {
    Message.error('操作失败')
    done(false)
  }
}

function handleDelete(record) {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除区域「${record.region_name}」吗？此操作不可恢复。`,
    okText: '删除',
    onOk: async () => {
      try {
        await deleteRegionSetting(record.id)
        Message.success('删除成功')
        await loadData()
      } catch (e) {
        Message.error('删除失败')
      }
    }
  })
}

onMounted(() => {
  loadData()
})
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

.header-card { flex-shrink: 0; }

.header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.page-title {
  font-size: 16px;
  font-weight: 600;
}

.filter-card { flex-shrink: 0; }

.filter-row {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.table-card {
  flex: 1;
  overflow: auto;
}

.region-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.region-name {
  font-weight: 500;
}

.mono {
  font-family: monospace;
  font-size: 12px;
}

.tz-text {
  color: var(--color-text-2);
}

.locale-tags {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.date-text {
  font-size: 12px;
  color: var(--color-text-3);
}

.action-btns {
  display: flex;
  gap: 4px;
}

.form-row {
  display: flex;
  gap: 12px;
}

.form-row .half { flex: 1; }

.form-row .third { flex: 1; }

.switch-label {
  margin-left: 10px;
  font-size: 13px;
  color: var(--color-text-3);
}
</style>
