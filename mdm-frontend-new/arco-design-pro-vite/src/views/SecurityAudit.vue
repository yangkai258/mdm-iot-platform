<template>
  <div class="security-audit-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="审计记录" :value="1256" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="安全评分" :value="95" suffix="/100" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="风险事件" :value="3" status="warning" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <span>具身AI安全审计</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="events" title="安全事件">
          <a-table :columns="eventColumns" :data="events" :pagination="pagination">
            <template #level="{ record }">
              <a-tag :color="getLevelColor(record.level)">{{ record.levelText }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="logs" title="操作日志">
          <a-table :columns="logColumns" :data="logs" :pagination="paginationSmall" />
        </a-tab-pane>
        
        <a-tab-pane key="settings" title="安全设置">
          <a-form layout="vertical" style="max-width: 600px;">
            <a-form-item label="动作安全检查">
              <a-switch v-model="settings.actionSafety" />
            </a-form-item>
            <a-form-item label="实时风险评估">
              <a-switch v-model="settings.realtimeRisk" />
            </a-form-item>
            <a-form-item label="自动熔断">
              <a-switch v-model="settings.autoCircuitBreaker" />
            </a-form-item>
            <a-form-item label="紧急停止阈值">
              <a-input-number v-model="settings.emergencyThreshold" :min="1" />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="handleSave">保存设置</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 3 });
const paginationSmall = reactive({ current: 1, pageSize: 10, total: 20 });

const settings = reactive({ actionSafety: true, realtimeRisk: true, autoCircuitBreaker: true, emergencyThreshold: 3 });

const events = ref([
  { id: 1, type: 'abnormal_behavior', typeText: '异常行为', level: 'warning', levelText: '警告', deviceId: 'DEV001', description: '检测到不寻常的动作模式', time: '2026-03-28 18:00:00' },
  { id: 2, type: 'boundary_violation', typeText: '边界违规', level: 'error', levelText: '错误', deviceId: 'DEV002', description: '动作超出安全边界', time: '2026-03-28 17:00:00' },
]);

const logs = ref([
  { id: 1, user: 'admin', action: '修改安全设置', ip: '192.168.1.100', time: '2026-03-28 18:00:00' },
  { id: 2, user: 'admin', action: '手动停止设备DEV001', ip: '192.168.1.100', time: '2026-03-28 17:30:00' },
]);

const eventColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '类型', dataIndex: 'typeText', width: 120 },
  { title: '级别', slotName: 'level', width: 100 },
  { title: '设备ID', dataIndex: 'deviceId', width: 120 },
  { title: '描述', dataIndex: 'description' },
  { title: '时间', dataIndex: 'time', width: 160 },
];

const logColumns = [
  { title: '用户', dataIndex: 'user', width: 120 },
  { title: '操作', dataIndex: 'action', width: 200 },
  { title: 'IP', dataIndex: 'ip', width: 150 },
  { title: '时间', dataIndex: 'time', width: 160 },
];

const getLevelColor = (l: string) => ({ info: 'blue', warning: 'orange', error: 'red', critical: '#F53F3F' }[l] || 'default');

const handleSave = () => {};
</script>

<style scoped>
.security-audit-container { padding: 20px; }
</style>
