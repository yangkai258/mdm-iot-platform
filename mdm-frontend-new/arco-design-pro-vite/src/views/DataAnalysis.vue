<template>
  <div class="data-analysis-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card>
          <a-statistic title="总用户数" :value="stats.totalUsers" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="总设备数" :value="stats.totalDevices" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="今日活跃" :value="stats.todayActive" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="转化率" :value="stats.conversionRate" suffix="%" />
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16">
      <a-col :span="12">
        <a-card title="用户趋势">
          <a-chart :option="userTrendChart" style="height: 300px;" />
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="设备状态分布">
          <a-chart :option="deviceStatusChart" style="height: 300px;" />
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16" style="margin-top: 16px;">
      <a-col :span="12">
        <a-card title="订阅收入">
          <a-chart :option="revenueChart" style="height: 300px;" />
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="地域分布">
          <a-table :columns="regionColumns" :data="regionData" :pagination="paginationSmall" />
        </a-card>
      </a-col>
    </a-row>

    <a-card title="实时数据" style="margin-top: 16px;">
      <a-tabs>
        <a-tab-pane key="active" title="实时活跃">
          <a-table :columns="activeColumns" :data="activeUsers" :pagination="paginationSmall">
            <template #device="{ record }">
              <a-tag>{{ record.deviceId }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
        <a-tab-pane key="alerts" title="实时告警">
          <a-table :columns="alertColumns" :data="realTimeAlerts" :pagination="paginationSmall">
            <template #level="{ record }">
              <a-tag :color="getAlertColor(record.level)">{{ record.levelText }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const paginationSmall = reactive({ current: 1, pageSize: 5, total: 5 });

const stats = reactive({
  totalUsers: 25680,
  totalDevices: 15832,
  todayActive: 3842,
  conversionRate: 23.5,
});

const userTrendChart = {
  xAxis: { type: 'category', data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'] },
  yAxis: { type: 'value' },
  series: [
    { name: '新增用户', type: 'bar', data: [120, 150, 180, 170, 200, 250, 220] },
    { name: '活跃用户', type: 'line', data: [3000, 3200, 3500, 3400, 3800, 4200, 3842], smooth: true },
  ],
};

const deviceStatusChart = {
  tooltip: { trigger: 'item' },
  legend: { bottom: 0 },
  series: [{ name: '设备状态', type: 'pie', radius: '60%', data: [
    { value: 12000, name: '在线' },
    { value: 3000, name: '离线' },
    { value: 500, name: '告警' },
    { value: 332, name: '升级中' },
  ]}],
};

const revenueChart = {
  xAxis: { type: 'category', data: ['1月', '2月', '3月', '4月', '5月', '6月'] },
  yAxis: { type: 'value', axisLabel: { formatter: '¥{value}' } },
  series: [{ name: '收入', type: 'bar', data: [120000, 135000, 158000, 142000, 168000, 185000] }],
};

const regionData = ref([
  { region: '华东', users: 8500, ratio: 33.1 },
  { region: '华北', users: 6200, ratio: 24.1 },
  { region: '华南', users: 5800, ratio: 22.6 },
  { region: '西南', users: 3200, ratio: 12.5 },
  { region: '其他', users: 1980, ratio: 7.7 },
]);

const activeUsers = ref([
  { userId: 'U001', deviceId: 'DEV001', action: 'AI对话', startTime: '18:50:00' },
  { userId: 'U002', deviceId: 'DEV005', action: 'OTA查询', startTime: '18:49:30' },
  { userId: 'U003', deviceId: 'DEV012', action: '固件升级', startTime: '18:49:00' },
]);

const realTimeAlerts = ref([
  { id: 'A001', deviceId: 'DEV010', level: 'critical', levelText: 'Critical', message: '设备离线超过30分钟', time: '18:50:00' },
  { id: 'A002', deviceId: 'DEV015', level: 'warning', levelText: 'Warning', message: '电量低于20%', time: '18:49:30' },
  { id: 'A003', deviceId: 'DEV020', level: 'info', levelText: 'Info', message: 'OTA升级完成', time: '18:49:00' },
]);

const regionColumns = [
  { title: '地区', dataIndex: 'region', width: 100 },
  { title: '用户数', dataIndex: 'users', width: 100 },
  { title: '占比', dataIndex: 'ratio', width: 100 },
];

const activeColumns = [
  { title: '用户ID', dataIndex: 'userId', width: 100 },
  { title: '设备', slotName: 'device', width: 100 },
  { title: '操作', dataIndex: 'action', width: 120 },
  { title: '开始时间', dataIndex: 'startTime', width: 120 },
];

const alertColumns = [
  { title: '告警ID', dataIndex: 'id', width: 80 },
  { title: '设备', dataIndex: 'deviceId', width: 100 },
  { title: '级别', slotName: 'level', width: 100 },
  { title: '消息', dataIndex: 'message' },
  { title: '时间', dataIndex: 'time', width: 100 },
];

const getAlertColor = (l: string) => ({ critical: 'red', warning: 'orange', info: 'blue' }[l] || 'default');
</script>

<style scoped>
.data-analysis-container { padding: 20px; }
</style>
