<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-shield /> 具身AI安全审计日志</a-space>
      </template>

      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6">
          <a-card>
            <a-statistic title="安全事件" :value="stats.total" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="边界触发" :value="stats.boundaryTriggered" :value-style="{ color: '#F56C6C' }" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="已处理" :value="stats.resolved" :value-style="{ color: '#67C23A' }" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="安全边界数" :value="stats.boundaries" />
          </a-card>
        </a-col>
      </a-row>

      <a-tabs default-active-key="logs">
        <a-tab-pane key="logs" title="安全日志">
          <a-space style="margin-bottom: 16px">
            <a-select v-model="filterType" placeholder="事件类型" allow-clear style="width: 150px">
              <a-option value="boundary">边界触发</a-option>
              <a-option value="collision">碰撞检测</a-option>
              <a-option value="speed">速度异常</a-option>
            </a-select>
            <a-select v-model="filterSeverity" placeholder="严重程度" allow-clear style="width: 120px">
              <a-option value="low">低</a-option>
              <a-option value="medium">中</a-option>
              <a-option value="high">高</a-option>
            </a-select>
            <a-button type="primary">查询</a-button>
          </a-space>

          <a-table :columns="columns" :data="logs">
            <template #severity="{ record }">
              <a-tag :color="getSeverityColor(record.severity)">{{ record.severity }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-link @click="handleViewDetail(record)">详情</a-link>
              <a-link @click="handleLocate(record)">定位</a-link>
            </template>
          </a-table>
        </a-tab-pane>

        <a-tab-pane key="boundaries" title="安全边界配置">
          <a-space style="margin-bottom: 16px">
            <a-button type="primary" @click="handleAddBoundary">
              <template #icon><icon-plus /></template>
              添加边界
            </a-button>
          </a-space>

          <a-table :columns="boundaryColumns" :data="boundaries">
            <template #type="{ record }">
              <a-tag :color="record.type === 'safe' ? 'green' : 'red'">{{ record.type === 'safe' ? '安全区' : '危险区' }}</a-tag>
            </template>
            <template #enabled="{ record }">
              <a-switch v-model="record.enabled" />
            </template>
            <template #actions="{ record }">
              <a-link @click="handleEditBoundary(record)">编辑</a-link>
              <a-link @click="handleDeleteBoundary(record)">删除</a-link>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const stats = reactive({ total: 25, boundaryTriggered: 8, resolved: 23, boundaries: 5 })
const filterType = ref('')
const filterSeverity = ref('')

const columns = [
  { title: '时间', dataIndex: 'time', width: 180 },
  { title: '事件类型', dataIndex: 'eventType' },
  { title: '严重程度', slotName: 'severity', width: 100 },
  { title: '位置', dataIndex: 'location' },
  { title: '触发动作', dataIndex: 'action' },
  { title: '操作', slotName: 'actions', width: 120 }
]

const logs = ref([
  { id: 1, time: '2026-03-28 10:00', eventType: '边界触发', severity: 'medium', location: '厨房区域', action: '停止前进' },
  { id: 2, time: '2026-03-28 09:30', eventType: '碰撞检测', severity: 'low', location: '客厅角落', action: '后退避让' }
])

const boundaryColumns = [
  { title: '边界名称', dataIndex: 'name' },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '范围', dataIndex: 'range' },
  { title: '触发动作', dataIndex: 'triggerAction' },
  { title: '启用', slotName: 'enabled', width: 80 },
  { title: '操作', slotName: 'actions', width: 120 }
]

const boundaries = ref([
  { id: 1, name: '厨房危险区', type: 'danger', range: '半径1米', triggerAction: '禁止进入', enabled: true }
])

const getSeverityColor = (s) => ({ low: 'blue', medium: 'orange', high: 'red' }[s] || 'gray')
const handleViewDetail = (r) => { }
const handleLocate = (r) => { }
const handleAddBoundary = () => { }
const handleEditBoundary = (r) => { }
const handleDeleteBoundary = (r) => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
