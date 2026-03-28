<template>
  <div class="alert-upgrade-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>告警升级配置</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新增规则
          </a-button>
        </div>
      </template>
      
      <a-alert type="warning" style="margin-bottom: 16px;">
        当告警在指定时间内未被处理，将自动升级严重程度并通知相关人员
      </a-alert>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #level="{ record }">
          <a-tag :color="getLevelColor(record.level)">{{ record.levelText }}</a-tag>
        </template>
        <template #escalateLevel="{ record }">
          <icon-arrow-right style="margin: 0 8px; color: #FF7D00;" />
          <a-tag :color="getLevelColor(record.escalateTo)">{{ record.escalateToText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-switch :checked="record.enabled" disabled />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 创建/编辑弹窗 -->
    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑规则' : '新增规则'" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="规则名称" required>
          <a-input v-model="form.name" placeholder="请输入规则名称" />
        </a-form-item>
        <a-form-item label="告警类型">
          <a-select v-model="form.alertType" placeholder="选择告警类型">
            <a-option value="device_offline">设备离线</a-option>
            <a-option value="low_battery">电量低</a-option>
            <a-option value="abnormal_behavior">异常行为</a-option>
            <a-option value="ota_failed">OTA升级失败</a-option>
            <a-option value="all">所有告警</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="触发条件">
          <a-space>
            <span>告警产生后</span>
            <a-input-number v-model="form.unattendedMinutes" :min="5" :max="1440" />
            <span>分钟未处理</span>
          </a-space>
        </a-form-item>
        <a-form-item label="升级方式">
          <a-radio-group v-model="form.escalateType">
            <a-radio value="level">升级严重程度</a-radio>
            <a-radio value="notify">增加通知人</a-radio>
            <a-radio value="both">同时升级+通知</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item v-if="form.escalateType !== 'notify'" label="升级到级别">
          <a-select v-model="form.escalateTo" placeholder="选择升级目标级别">
            <a-option value="warning">Warning</a-option>
            <a-option value="critical">Critical</a-option>
            <a-option value="fatal">Fatal</a-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="form.escalateType !== 'level'" label="追加通知人">
          <a-select v-model="form.additionalReceivers" multiple placeholder="选择通知人">
            <a-option value="manager">部门经理</a-option>
            <a-option value="tech_lead">技术负责人</a-option>
            <a-option value="ops">运维人员</a-option>
            <a-option value="admin">系统管理员</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="启用">
          <a-switch v-model="form.enabled" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const data = ref([
  { id: 1, name: 'Critical超时升级', alertType: 'all', alertTypeText: '所有告警', level: 'critical', levelText: 'Critical', unattendedMinutes: 30, escalateType: 'level', escalateTypeText: '升级级别', escalateTo: 'fatal', escalateToText: 'Fatal', additionalReceivers: [], enabled: true },
  { id: 2, name: 'Warning超时升级', alertType: 'all', alertTypeText: '所有告警', level: 'warning', levelText: 'Warning', unattendedMinutes: 60, escalateType: 'both', escalateTypeText: '同时', escalateTo: 'critical', escalateToText: 'Critical', additionalReceivers: ['manager'], enabled: true },
  { id: 3, name: '设备离线升级', alertType: 'device_offline', alertTypeText: '设备离线', level: 'critical', levelText: 'Critical', unattendedMinutes: 15, escalateType: 'notify', escalateTypeText: '通知', escalateTo: '', escalateToText: '-', additionalReceivers: ['tech_lead', 'ops'], enabled: true },
  { id: 4, name: 'OTA失败升级', alertType: 'ota_failed', alertTypeText: 'OTA升级失败', level: 'warning', levelText: 'Warning', unattendedMinutes: 120, escalateType: 'level', escalateTypeText: '升级级别', escalateTo: 'critical', escalateToText: 'Critical', additionalReceivers: [], enabled: false },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 4 });
const editVisible = ref(false);
const isEdit = ref(false);

const form = reactive({
  name: '',
  alertType: '',
  level: 'warning',
  unattendedMinutes: 30,
  escalateType: 'level',
  escalateTo: 'critical',
  additionalReceivers: [],
  enabled: true,
});

const columns = [
  { title: '规则ID', dataIndex: 'id', width: 80 },
  { title: '规则名称', dataIndex: 'name', width: 150 },
  { title: '告警类型', dataIndex: 'alertTypeText', width: 120 },
  { title: '原始级别', slotName: 'level', width: 100 },
  { title: '触发时间', width: 120 },
  { title: '升级方式', dataIndex: 'escalateTypeText', width: 100 },
  { title: '升级到', slotName: 'escalateLevel', width: 150 },
  { title: '启用', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' },
];

const getLevelColor = (level: string) => {
  const map: Record<string, string> = { warning: 'orange', critical: 'red', fatal: '#F53F3F' };
  return map[level] || 'default';
};

const handleCreate = () => {
  isEdit.value = false;
  Object.assign(form, { name: '', alertType: '', level: 'warning', unattendedMinutes: 30, escalateType: 'level', escalateTo: 'critical', additionalReceivers: [], enabled: true });
  editVisible.value = true;
};

const handleEdit = (record: any) => {
  isEdit.value = true;
  Object.assign(form, record);
  editVisible.value = true;
};

const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.alert-upgrade-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
