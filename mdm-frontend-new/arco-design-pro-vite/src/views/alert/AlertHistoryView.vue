<template>
  <Breadcrumb :items="['Home','Alert','AlertHistory','']" />
  <div class="page-container">
    <a-card class="general-card" title="告警历史">
      <template #extra>
        <a-button @click="loadData"><icon-refresh />刷新</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="设备ID"><a-input v-model="form.device_id" placeholder="设备ID" /></a-form-item>
          <a-form-item label="告警类型">
            <a-select v-model="form.alert_type" placeholder="选择类型" allow-clear style="width: 120px">
              <a-option value="temperature">温度告警</a-option>
              <a-option value="battery">电量告警</a-option>
              <a-option value="offline">离线告警</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="级别">
            <a-select v-model="form.severity" placeholder="选择级别" allow-clear style="width: 100px">
              <a-option value="info">提示</a-option>
              <a-option value="warning">警告</a-option>
              <a-option value="error">错误</a-option>
            </a-select>
          </a-form-item>
          <a-form-item>
            <a-button type="primary" @click="loadData">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" />
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconRefresh } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const data = ref<any[]>([])
const form = ref<any>({ device_id: '', alert_type: '', severity: '' })

const columns = [
  { title: '时间', dataIndex: 'created_at', width: 170 },
  { title: '设备', dataIndex: 'device_id', width: 140 },
  { title: '告警类型', dataIndex: 'alert_type', width: 120 },
  { title: '级别', dataIndex: 'severity', width: 80 },
  { title: '告警消息', dataIndex: 'message', ellipsis: true },
  { title: '触发值', dataIndex: 'trigger_value', width: 100 },
  { title: '状态', dataIndex: 'status', width: 90 }
]

const pagination = ref({ current: 1, pageSize: 20, total: 0, showTotal: true })

async function loadData() {
  try {
    loading.value = true
    data.value = []
    pagination.value.total = 0
  } catch (err: any) {
    Message.error('加载失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

function handleReset() {
  form.value = { device_id: '', alert_type: '', severity: '' }
  loadData()
}

onMounted(() => { loadData() })
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>
