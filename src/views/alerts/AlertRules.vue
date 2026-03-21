<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>告警管理</a-breadcrumb-item>
      <a-breadcrumb-item>告警规则</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 操作按钮组 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showAddModal">
          <template #icon><icon-plus /></template>
          新建规则
        </a-button>
        <a-button @click="loadRules">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
      </a-space>
    </div>

    <!-- 数据表格 -->
    <div class="pro-content-area">
      <a-table
        :columns="columns"
        :data="rules"
        :loading="loading"
        :pagination="{ pageSize: 20 }"
        row-key="id"
      >
        <template #name="{ record }">
          <a-space direction="vertical" :size="0">
            <span class="rule-name">{{ record.name }}</span>
            <span v-if="record.remark" class="rule-remark">{{ record.remark }}</span>
          </a-space>
        </template>

        <template #device_id="{ record }">
          <span v-if="record.device_id">{{ record.device_id }}</span>
          <a-tag v-else color="arcoblue" size="small">全局</a-tag>
        </template>

        <template #alert_type="{ record }">
          <a-tag :color="getAlertTypeColor(record.alert_type)">
            {{ getAlertTypeText(record.alert_type) }}
          </a-tag>
        </template>

        <template #condition="{ record }">
          <span class="condition-text">
            {{ getConditionSymbol(record.condition) }} {{ record.threshold }}
          </span>
        </template>

        <template #notify_ways="{ record }">
          <a-space>
            <a-tag v-if="hasNotifyWay(record.notify_ways, 'email')" size="small">
              <template #icon><icon-email /></template> 邮件
            </a-tag>
            <a-tag v-if="hasNotifyWay(record.notify_ways, 'webhook')" size="small">
              <template #icon><icon-link /></template> Webhook
            </a-tag>
            <a-tag v-if="hasNotifyWay(record.notify_ways, 'inapp') || !record.notify_ways" size="small" color="arcoblue">
              <template #icon><icon-bell /></template> 站内
            </a-tag>
          </a-space>
        </template>

        <template #severity="{ record }">
          <a-tag :color="getSeverityColor(record.severity)">
            {{ getSeverityText(record.severity) }}
          </a-tag>
        </template>

        <template #enabled="{ record }">
          <a-switch v-model="record.enabled" @change="handleToggle(record)" />
        </template>

        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="editRule(record)">编辑</a-button>
            <a-popconfirm content="确定删除该规则？" @ok="handleDelete(record.id)">
              <a-button type="text" size="small" status="danger">删除</a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 添加/编辑规则弹窗 -->
    <a-modal
      v-model:visible="modalVisible"
      :title="isEdit ? '编辑规则' : '新建规则'"
      @ok="handleSubmit"
      :width="560"
      @cancel="modalVisible = false"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="规则名称" required>
          <a-input v-model="form.name" placeholder="如：电量过低告警" />
        </a-form-item>

        <a-form-item label="告警类型" required>
          <a-select v-model="form.alert_type" placeholder="选择告警类型">
            <a-option value="battery_low">
              <a-space>🔋 电量过低</a-space>
            </a-option>
            <a-option value="offline">📡 设备离线</a-option>
            <a-option value="temperature_high">🌡️ 温度过高</a-option>
            <a-option value="geofence_violation">📍 地理围栏违规</a-option>
            <a-option value="jailbreak_detected">🔓 越狱/Root检测</a-option>
            <a-option value="storage_low">💾 存储空间不足</a-option>
            <a-option value="network_unreachable">🌐 网络不可达</a-option>
          </a-select>
        </a-form-item>

        <a-form-item label="设备ID（留空表示所有设备）">
          <a-input v-model="form.device_id" placeholder="留空则对所有设备生效" />
        </a-form-item>

        <a-form-item label="触发条件" required>
          <a-space>
            <a-select v-model="form.condition" style="width: 90px">
              <a-option value="<">&lt; 小于</a-option>
              <a-option value=">">&gt; 大于</a-option>
              <a-option value="=">= 等于</a-option>
              <a-option value="<=">&lt;= 小于等于</a-option>
              <a-option value=">=">&gt;= 大于等于</a-option>
            </a-select>
            <a-input-number v-model="form.threshold" :min="0" style="width: 140px" />
          </a-space>
        </a-form-item>

        <a-form-item label="严重程度" required>
          <a-radio-group v-model="form.severity">
            <a-radio :value="1"><a-tag color="green">低</a-tag></a-radio>
            <a-radio :value="2"><a-tag color="arcoblue">中</a-tag></a-radio>
            <a-radio :value="3"><a-tag color="orange">高</a-tag></a-radio>
            <a-radio :value="4"><a-tag color="red">严重</a-tag></a-radio>
          </a-radio-group>
        </a-form-item>

        <a-form-item label="通知方式">
          <a-checkbox-group v-model="form.notifyWays">
            <a-checkbox value="email">
              <template #icon><icon-email /></template> 邮件
            </a-checkbox>
            <a-checkbox value="webhook">Webhook</a-checkbox>
            <a-checkbox value="inapp">站内通知</a-checkbox>
          </a-checkbox-group>
        </a-form-item>

        <a-form-item label="备注">
          <a-textarea v-model="form.remark" :rows="2" placeholder="可选备注信息" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as api from '@/api/alerts'

const loading = ref(false)
const rules = ref([])
const modalVisible = ref(false)
const isEdit = ref(false)

const form = reactive({
  name: '',
  device_id: '',
  alert_type: 'battery_low',
  condition: '<',
  threshold: 20,
  severity: 2,
  notifyWays: ['inapp'],
  remark: ''
})

const columns = [
  { title: '规则名称', slotName: 'name', minWidth: 180 },
  { title: '设备', slotName: 'device_id', width: 160 },
  { title: '告警类型', slotName: 'alert_type', width: 130 },
  { title: '触发条件', slotName: 'condition', width: 120 },
  { title: '通知方式', slotName: 'notify_ways', width: 200 },
  { title: '严重程度', slotName: 'severity', width: 100 },
  { title: '启用', slotName: 'enabled', width: 80 },
  { title: '操作', slotName: 'actions', width: 140, fixed: 'right' }
]

const getSeverityColor = (s) => ({ 1: 'green', 2: 'arcoblue', 3: 'orange', 4: 'red' })[s] || 'gray'
const getSeverityText = (s) => ({ 1: '低', 2: '中', 3: '高', 4: '严重' })[s] || '未知'

const getAlertTypeColor = (t) => ({
  battery_low: 'green', offline: 'orange', temperature_high: 'red',
  geofence_violation: 'purple', jailbreak_detected: 'red', storage_low: 'orange', network_unreachable: 'gray'
})[t] || 'arcoblue'

const getAlertTypeText = (t) => ({
  battery_low: '电量过低', offline: '设备离线', temperature_high: '温度过高',
  geofence_violation: '地理围栏', jailbreak_detected: '越狱/Root', storage_low: '存储不足', network_unreachable: '网络不可达'
})[t] || t

const getConditionSymbol = (c) => ({
  '<': '<', '>': '>', '=': '=', '<=': '≤', '>=': '≥'
})[c] || c

const hasNotifyWay = (ways, way) => {
  if (!ways) return false
  return ways.split(',').includes(way)
}

async function loadRules() {
  loading.value = true
  try {
    const data = await api.getAlertRules()
    rules.value = data.data?.list || []
  } catch (e) {
    Message.error('加载规则失败')
  } finally {
    loading.value = false
  }
}

function showAddModal() {
  isEdit.value = false
  Object.assign(form, {
    name: '', device_id: '', alert_type: 'battery_low',
    condition: '<', threshold: 20, severity: 2,
    notifyWays: ['inapp'], remark: ''
  })
  modalVisible.value = true
}

function editRule(record) {
  isEdit.value = true
  Object.assign(form, {
    name: record.name,
    device_id: record.device_id || '',
    alert_type: record.alert_type,
    condition: record.condition,
    threshold: record.threshold,
    severity: record.severity,
    notifyWays: record.notify_ways ? record.notify_ways.split(',') : ['inapp'],
    remark: record.remark || ''
  })
  modalVisible.value = true
}

async function handleSubmit() {
  if (!form.name) {
    Message.warning('请输入规则名称')
    return
  }
  const payload = {
    ...form,
    notify_ways: form.notifyWays.join(',')
  }
  try {
    if (isEdit.value) {
      const id = rules.value.find(r => r.name === form.name && r.alert_type === form.alert_type)?.id
      if (!id) { Message.error('无法确定要编辑的规则'); return }
      await api.updateAlertRule(id, payload)
      Message.success('更新成功')
    } else {
      await api.createAlertRule(payload)
      Message.success('创建成功')
    }
    modalVisible.value = false
    loadRules()
  } catch (e) {
    Message.error(isEdit.value ? '更新失败' : '创建失败')
  }
}

async function handleToggle(record) {
  try {
    await api.updateAlertRule(record.id, { ...record, notify_ways: record.notify_ways || 'inapp' })
    Message.success(record.enabled ? '已启用' : '已停用')
  } catch (e) {
    record.enabled = !record.enabled
    Message.error('操作失败')
  }
}

async function handleDelete(id) {
  try {
    await api.deleteAlertRule(id)
    Message.success('删除成功')
    loadRules()
  } catch (e) {
    Message.error('删除失败')
  }
}

onMounted(() => { loadRules() })
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.pro-action-bar { margin-bottom: 16px; }
.pro-content-area {
  background: #fff; border-radius: 8px; padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}
.rule-name { font-weight: 500; }
.rule-remark { font-size: 12px; color: #888; }
.condition-text { font-family: monospace; }
</style>
