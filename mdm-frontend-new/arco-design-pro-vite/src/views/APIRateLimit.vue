<template>
  <div class="rate-limit-container">
    <a-card>
      <template #title>
        <span>API配额管理</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="plans" title="配额套餐">
          <a-table :columns="planColumns" :data="plans" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="record.enabled ? 'green' : 'gray'">{{ record.enabled ? '启用' : '禁用' }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-link @click="handleEditPlan(record)">编辑</a-link>
                <a-link status="danger" @click="handleDeletePlan(record)">删除</a-link>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="usage" title="实时用量">
          <a-row :gutter="16" style="margin-bottom: 16px;">
            <a-col :span="6">
              <a-card>
                <a-statistic title="今日API调用" :value="todayStats.calls" suffix="次" />
              </a-card>
            </a-col>
            <a-col :span="6">
              <a-card>
                <a-statistic title="配额使用率" :value="todayStats.usageRate" suffix="%" />
              </a-card>
            </a-col>
            <a-col :span="6">
              <a-card>
                <a-statistic title="限流次数" :value="todayStats.throttled" status="warning" />
              </a-card>
            </a-col>
            <a-col :span="6">
              <a-card>
                <a-statistic title="错误率" :value="todayStats.errorRate" suffix="%" />
              </a-card>
            </a-col>
          </a-row>
          <a-table :columns="usageColumns" :data="usageData" :pagination="pagination" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 4 });

const plans = ref([
  { id: 1, name: '免费版', apiCalls: 1000, rateLimit: '60次/分钟', quotaType: 'basic', enabled: true },
  { id: 2, name: '基础版', apiCalls: 10000, rateLimit: '300次/分钟', quotaType: 'standard', enabled: true },
  { id: 3, name: '高级版', apiCalls: 100000, rateLimit: '1000次/分钟', quotaType: 'premium', enabled: true },
  { id: 4, name: '企业版', apiCalls: -1, rateLimit: '无限制', quotaType: 'enterprise', enabled: true },
]);

const usageData = ref([
  { appName: '我的IoT应用', plan: '基础版', todayCalls: 8563, quota: 10000, usageRate: 85, throttled: 12 },
  { appName: '宠物健康App', plan: '高级版', todayCalls: 45230, quota: 100000, usageRate: 45, throttled: 0 },
  { appName: '测试应用', plan: '免费版', todayCalls: 1200, quota: 1000, usageRate: 120, throttled: 89 },
]);

const todayStats = reactive({
  calls: 158320,
  usageRate: 62,
  throttled: 156,
  errorRate: 0.5,
});

const planColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '套餐名称', dataIndex: 'name', width: 120 },
  { title: 'API调用量/天', dataIndex: 'apiCalls', width: 120 },
  { title: '速率限制', dataIndex: 'rateLimit', width: 150 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const usageColumns = [
  { title: '应用', dataIndex: 'appName', width: 150 },
  { title: '套餐', dataIndex: 'plan', width: 100 },
  { title: '今日调用', dataIndex: 'todayCalls', width: 100 },
  { title: '配额', dataIndex: 'quota', width: 100 },
  { title: '使用率', dataIndex: 'usageRate', width: 100 },
  { title: '限流次数', dataIndex: 'throttled', width: 100 },
];

const handleEditPlan = (record: any) => {};
const handleDeletePlan = (record: any) => {};
</script>

<style scoped>
.rate-limit-container { padding: 20px; }
</style>
