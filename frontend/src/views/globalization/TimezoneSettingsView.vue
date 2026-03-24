<template>
  <div class="page-container">
    <div class="current-tz-section">
      <div class="tz-label">当前时区</div>
      <div class="tz-value">{{ currentTz ? `${currentTz.timezone} (UTC${currentTz.offset >= 0 ? '+' : ''}${currentTz.offset})` : '-' }}</div>
      <a-button type="primary" @click="openSwitchModal">切换时区</a-button>
    </div>
    <div class="search-form">
      <div class="section-title">组织时区设置</div>
      <a-form :model="orgSettings" layout="inline">
        <a-form-item label="租户级别时区">
          <a-select v-model="orgSettings.tenantTimezone" placeholder="继承系统设置" style="width: 280px">
            <a-option value="inherit">继承系统设置</a-option>
            <a-option v-for="tz in timezoneList" :key="tz.value" :value="tz.value">{{ tz.label }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="部门">
          <a-select v-model="orgSettings.departmentId" placeholder="使用部门设置" style="width: 200px" allow-search>
            <a-option value="">使用部门设置</a-option>
            <a-option v-for="dept in departments" :key="dept.id" :value="dept.id">{{ dept.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" :loading="saving" @click="handleSaveOrgTz">保存设置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <div class="section-title">时区列表</div>
    </div>
    <div class="tz-display-list">
      <div class="tz-display-item">
        <span class="tz-display-label">显示格式</span>
        <span class="tz-display-value">{{ nowLocal }}</span>
      </div>
      <div class="tz-display-item">
        <span class="tz-display-label">UTC</span>
        <span class="tz-display-value">{{ nowUtc }}</span>
      </div>
    </div>
    <a-table :columns="tzColumns" :data="timezoneList" :pagination="pagination" row-key="value" size="small">
      <template #offset="{ record }">UTC{{ record.offset >= 0 ? '+' : '' }}{{ record.offset }}</template>
    </a-table>
    <a-modal v-model:visible="switchModalVisible" title="切换时区" :width="480" @before-ok="handleSwitchTz">
      <a-form :model="tzForm" label-col-flex="100px">
        <a-form-item label="选择时区" required>
          <a-select v-model="tzForm.timezone" placeholder="请选择时区" style="width: 100%" filterable>
            <a-option v-for="tz in timezoneList" :key="tz.value" :value="tz.value">{{ tz.label }}</a-option>
          </a-select>
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="switchModalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSwitchTz">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { getTimezone, updateTimezone } from '@/api/globalization'
import { TIMEZONE_LIST } from '@/composables/useGlobalization'
import dayjs from 'dayjs'
import utc from 'dayjs/plugin/utc'
import timezone from 'dayjs/plugin/timezone'

dayjs.extend(utc)
dayjs.extend(timezone)

const timezoneList = TIMEZONE_LIST
const currentTz = ref(null)
const saving = ref(false)
const switchModalVisible = ref(false)
const departments = ref([])
const pagination = reactive({ pageSize: 15, current: 1, total: 0 })

const nowLocal = computed(() => {
  if (!currentTz.value) return '-'
  return dayjs().tz(currentTz.value.timezone).format('YYYY-MM-DD HH:mm:ss') + ` (${currentTz.value.timezone})`
})

const nowUtc = computed(() => {
  return dayjs().utc().format('YYYY-MM-DD HH:mm:ss') + ' (UTC)'
})

const orgSettings = reactive({ tenantTimezone: 'inherit', departmentId: '' })
const tzForm = reactive({ timezone: '' })

const tzColumns = [
  { title: '时区标识', dataIndex: 'value', width: 220 },
  { title: 'UTC 偏移', slotName: 'offset', width: 100 },
  { title: '标签', dataIndex: 'label', ellipsis: true }
]

onMounted(async () => {
  try {
    const res = await getTimezone()
    const data = res.data || res
    currentTz.value = { timezone: data.timezone || 'Asia/Shanghai', offset: data.offset || 8 }
    tzForm.timezone = currentTz.value.timezone
    if (data.tenant_timezone) orgSettings.tenantTimezone = data.tenant_timezone
    if (data.department_id) orgSettings.departmentId = data.department_id
    pagination.total = timezoneList.length
  } catch (e) {
    console.error('加载时区失败', e)
  }
})

function openSwitchModal() {
  tzForm.timezone = currentTz.value?.timezone || 'Asia/Shanghai'
  switchModalVisible.value = true
}

async function handleSwitchTz() {
  if (!tzForm.timezone) {
    Message.warning('请选择时区')
    return
  }
  try {
    await updateTimezone({ timezone: tzForm.timezone })
    currentTz.value = { timezone: tzForm.timezone, offset: timezoneList.find(t => t.value === tzForm.timezone)?.offset || 0 }
    Message.success('时区切换成功')
    switchModalVisible.value = false
  } catch (e) {
    Message.error('切换失败')
  }
}

async function handleSaveOrgTz() {
  saving.value = true
  try {
    await updateTimezone({ tenant_timezone: orgSettings.tenantTimezone, department_id: orgSettings.departmentId || null })
    Message.success('保存成功')
  } catch (e) {
    Message.error('保存失败')
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.current-tz-section { display: flex; align-items: center; gap: 16px; margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.tz-label { font-size: 14px; color: var(--color-text-3); width: 80px; }
.tz-value { font-size: 16px; font-weight: 600; flex: 1; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.section-title { font-size: 14px; font-weight: 600; margin-bottom: 12px; }
.toolbar { margin-bottom: 16px; }
.tz-display-list { display: flex; flex-direction: column; gap: 8px; margin-bottom: 16px; }
.tz-display-item { display: flex; align-items: center; gap: 12px; }
.tz-display-label { font-size: 13px; color: var(--color-text-3); width: 80px; }
.tz-display-value { font-size: 14px; font-family: monospace; }
</style>
