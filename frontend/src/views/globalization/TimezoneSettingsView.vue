<template>
  <div class="page-container">
    <!-- 当前时区 -->
    <a-card class="current-tz-card">
      <div class="current-tz-section">
        <div class="tz-label">当前时区</div>
        <div class="tz-value">
          {{ currentTz ? `${currentTz.timezone} (UTC${currentTz.offset >= 0 ? '+' : ''}${currentTz.offset})` : '-' }}
        </div>
        <a-button type="primary" @click="openSwitchModal">切换时区</a-button>
      </div>
    </a-card>

    <!-- 组织时区 -->
    <a-card class="org-tz-card">
      <template #title>
        <div class="section-title">组织时区</div>
      </template>
      <div class="org-tz-content">
        <div class="org-tz-row">
          <div class="org-tz-label">租户级别时区</div>
          <div class="org-tz-value">
            <a-select
              v-model="orgSettings.tenantTimezone"
              placeholder="继承系统设置"
              style="width: 280px"
              :disabled="saving"
            >
              <a-option value="inherit">继承系统设置</a-option>
              <a-option v-for="tz in timezoneList" :key="tz.value" :value="tz.value">{{ tz.label }}</a-option>
            </a-select>
          </div>
        </div>
        <div class="org-tz-row">
          <div class="org-tz-label">部门级别时区</div>
          <div class="org-tz-value">
            <a-select
              v-model="orgSettings.departmentId"
              placeholder="使用部门设置"
              style="width: 200px"
              allow-search
              :disabled="saving"
            >
              <a-option value="">使用部门设置</a-option>
              <a-option v-for="dept in departments" :key="dept.id" :value="dept.id">{{ dept.name }}</a-option>
            </a-select>
            <a-tag v-if="orgSettings.departmentId" color="arcoblue" style="margin-left: 8px">使用部门设置</a-tag>
          </div>
        </div>
        <div class="org-tz-save">
          <a-button type="primary" :loading="saving" @click="handleSaveOrgTz">保存设置</a-button>
        </div>
      </div>
    </a-card>

    <!-- 时区列表 -->
    <a-card class="tz-list-card">
      <template #title>
        <div class="section-title">时区列表</div>
      </template>
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
      <a-divider />
      <div class="tz-table-section">
        <a-table
          :columns="tzColumns"
          :data="timezoneList"
          :pagination="{ pageSize: 15 }"
          row-key="value"
          size="small"
        >
          <template #offset="{ record }">
            UTC{{ record.offset >= 0 ? '+' : '' }}{{ record.offset }}
          </template>
        </a-table>
      </div>
    </a-card>

    <!-- 切换时区弹窗 -->
    <a-modal
      v-model:visible="switchModalVisible"
      title="切换时区"
      :width="480"
      @before-ok="handleSwitchTz"
      @cancel="switchModalVisible = false"
    >
      <a-form :model="tzForm" layout="vertical">
        <a-form-item label="选择时区" required>
          <a-select
            v-model="tzForm.timezone"
            placeholder="请选择时区"
            style="width: 100%"
            filterable
          >
            <a-option v-for="tz in timezoneList" :key="tz.value" :value="tz.value">{{ tz.label }}</a-option>
          </a-select>
        </a-form-item>
      </a-form>
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

const nowLocal = computed(() => {
  if (!currentTz.value) return '-'
  return dayjs().tz(currentTz.value.timezone).format('YYYY-MM-DD HH:mm:ss') + ` (${currentTz.value.timezone})`
})

const nowUtc = computed(() => {
  return dayjs().utc().format('YYYY-MM-DD HH:mm:ss') + ' (UTC)'
})

const orgSettings = reactive({
  tenantTimezone: 'inherit',
  departmentId: ''
})

const tzForm = reactive({
  timezone: ''
})

const tzColumns = [
  { title: '时区标识', dataIndex: 'value', width: 220 },
  { title: 'UTC 偏移', slotName: 'offset', width: 100 },
  { title: '标签', dataIndex: 'label', ellipsis: true }
]

onMounted(async () => {
  try {
    const res = await getTimezone()
    const data = res.data || res
    currentTz.value = {
      timezone: data.timezone || 'Asia/Shanghai',
      offset: data.offset || 8
    }
    tzForm.timezone = currentTz.value.timezone
    if (data.tenant_timezone) orgSettings.tenantTimezone = data.tenant_timezone
    if (data.department_id) orgSettings.departmentId = data.department_id
  } catch (e) {
    console.error('加载时区失败', e)
  }
})

function openSwitchModal() {
  tzForm.timezone = currentTz.value?.timezone || 'Asia/Shanghai'
  switchModalVisible.value = true
}

async function handleSwitchTz(done) {
  if (!tzForm.timezone) {
    Message.warning('请选择时区')
    done(false)
    return
  }
  try {
    await updateTimezone({ timezone: tzForm.timezone })
    currentTz.value = {
      timezone: tzForm.timezone,
      offset: timezoneList.find(t => t.value === tzForm.timezone)?.offset || 0
    }
    Message.success('时区切换成功')
    switchModalVisible.value = false
  } catch (e) {
    Message.error('切换失败')
    done(false)
  }
}

async function handleSaveOrgTz() {
  saving.value = true
  try {
    await updateTimezone({
      tenant_timezone: orgSettings.tenantTimezone,
      department_id: orgSettings.departmentId || null
    })
    Message.success('保存成功')
  } catch (e) {
    Message.error('保存失败')
  } finally {
    saving.value = false
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

.current-tz-card {
  flex-shrink: 0;
}

.current-tz-section {
  display: flex;
  align-items: center;
  gap: 16px;
}

.tz-label {
  font-size: 14px;
  color: var(--color-text-3);
  width: 80px;
}

.tz-value {
  font-size: 16px;
  font-weight: 600;
  flex: 1;
}

.org-tz-card,
.tz-list-card {
  flex-shrink: 0;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
}

.org-tz-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.org-tz-row {
  display: flex;
  align-items: center;
  gap: 16px;
}

.org-tz-label {
  font-size: 13px;
  color: var(--color-text-2);
  width: 120px;
  flex-shrink: 0;
}

.org-tz-value {
  display: flex;
  align-items: center;
}

.org-tz-save {
  margin-top: 8px;
}

.tz-display-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.tz-display-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.tz-display-label {
  font-size: 13px;
  color: var(--color-text-3);
  width: 80px;
}

.tz-display-value {
  font-size: 14px;
  font-family: monospace;
}

.tz-table-section {
  max-height: 400px;
  overflow: auto;
}
</style>
