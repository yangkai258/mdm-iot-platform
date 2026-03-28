<template>
  <div class="feedback-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="反馈总数" :value="456" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="待处理" :value="28" status="warning" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="解决率" :value="94" suffix="%" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="平均响应" :value="4.5" suffix="小时" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <span>用户反馈管理</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="list" title="反馈列表">
          <a-table :columns="columns" :data="feedbacks" :pagination="pagination">
            <template #type="{ record }">
              <a-tag>{{ record.typeText }}</a-tag>
            </template>
            <template #priority="{ record }">
              <a-tag :color="getPriorityColor(record.priority)">{{ record.priorityText }}</a-tag>
            </template>
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="statistics" title="统计分析">
          <a-row :gutter="16">
            <a-col :span="12">
              <a-card title="反馈类型分布">
                <a-progress :percent="40" label="功能建议" />
                <a-progress :percent="30" label="问题反馈" />
                <a-progress :percent="20" label="体验投诉" />
                <a-progress :percent="10" label="其他" />
              </a-card>
            </a-col>
            <a-col :span="12">
              <a-card title="解决效率">
                <a-statistic title="24小时内解决" :value="78" suffix="%" />
                <a-statistic title="平均解决时间" :value="4.5" suffix="小时" />
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 8 });

const feedbacks = ref([
  { id: 1, userName: '张三', type: 'suggestion', typeText: '功能建议', priority: 'medium', priorityText: '中', title: '建议增加定时喂食功能', status: 'pending', statusText: '待处理', createdAt: '2026-03-28 10:00:00' },
  { id: 2, userName: '李四', type: 'bug', typeText: '问题反馈', priority: 'high', priorityText: '高', title: '设备经常掉线', status: 'resolved', statusText: '已解决', createdAt: '2026-03-27 15:00:00' },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '用户', dataIndex: 'userName', width: 100 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '优先级', slotName: 'priority', width: 80 },
  { title: '标题', dataIndex: 'title' },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '时间', dataIndex: 'createdAt', width: 160 },
];

const getPriorityColor = (p: string) => ({ low: 'green', medium: 'blue', high: 'orange', urgent: 'red' }[p] || 'default');
const getStatusColor = (s: string) => ({ pending: 'orange', processing: 'blue', resolved: 'green', rejected: 'gray' }[s] || 'default');
</script>

<style scoped>
.feedback-container { padding: 20px; }
</style>
