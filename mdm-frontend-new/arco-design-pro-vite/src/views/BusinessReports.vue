<template>
  <div class="business-reports-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="报告总数" :value="86" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="本周生成" :value="12" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="订阅用户" :value="56" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="平均生成" :value="45" suffix="秒" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>业务报表中心</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建报表
          </a-button>
        </div>
      </template>
      
      <a-tabs>
        <a-tab-pane key="reports" title="报表列表">
          <a-table :columns="reportColumns" :data="reports" :pagination="pagination">
            <template #type="{ record }">
              <a-tag>{{ record.typeText }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-link @click="handleView(record)">查看</a-link>
                <a-link @click="handleDownload(record)">下载</a-link>
                <a-link @click="handleSubscribe(record)">订阅</a-link>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="subscriptions" title="订阅管理">
          <a-table :columns="subColumns" :data="subscriptions" :pagination="pagination" />
        </a-tab-pane>
        
        <a-tab-pane key="schedules" title="定时报表">
          <a-table :columns="scheduleColumns" :data="schedules" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="record.enabled ? 'green' : 'gray'">{{ record.enabled ? '启用' : '禁用' }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 8 });

const reports = ref([
  { id: 1, name: '设备销售报表', type: 'sales', typeText: '销售', period: '2026-03', generatedAt: '2026-03-28 10:00:00', size: '2.5MB' },
  { id: 2, name: '会员活跃报表', type: 'member', typeText: '会员', period: '2026-03', generatedAt: '2026-03-28 09:00:00', size: '1.8MB' },
]);

const subscriptions = ref([
  { id: 1, reportName: '设备销售报表', userName: 'admin', frequency: 'daily', nextRun: '2026-03-29 10:00:00' },
]);

const schedules = ref([
  { id: 1, reportName: '设备销售报表', schedule: '每天 10:00', email: 'admin@example.com', enabled: true },
]);

const reportColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '报表名称', dataIndex: 'name', width: 150 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '周期', dataIndex: 'period', width: 100 },
  { title: '生成时间', dataIndex: 'generatedAt', width: 160 },
  { title: '大小', dataIndex: 'size', width: 100 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' },
];

const subColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '报表', dataIndex: 'reportName', width: 150 },
  { title: '订阅用户', dataIndex: 'userName', width: 120 },
  { title: '频率', dataIndex: 'frequency', width: 100 },
  { title: '下次执行', dataIndex: 'nextRun', width: 160 },
];

const scheduleColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '报表', dataIndex: 'reportName', width: 150 },
  { title: '执行时间', dataIndex: 'schedule', width: 120 },
  { title: '通知邮箱', dataIndex: 'email', width: 180 },
  { title: '状态', slotName: 'status', width: 80 },
];

const handleCreate = () => {};
const handleView = (record: any) => {};
const handleDownload = (record: any) => {};
const handleSubscribe = (record: any) => {};
</script>

<style scoped>
.business-reports-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
