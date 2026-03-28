<template>
  <div class="customer-service-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="工单总数" :value="128" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="待处理" :value="12" status="warning" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="今日响应" :value="96" suffix="%" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="满意度" :value="98" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>客服工单管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建工单
          </a-button>
        </div>
      </template>
      
      <a-tabs>
        <a-tab-pane key="tickets" title="工单列表">
          <a-table :columns="ticketColumns" :data="tickets" :pagination="pagination">
            <template #priority="{ record }">
              <a-tag :color="getPriorityColor(record.priority)">{{ record.priorityText }}</a-tag>
            </template>
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="statistics" title="统计报表">
          <a-row :gutter="16">
            <a-col :span="12">
              <a-card title="工单来源分布">
                <a-progress :percent="30" label="设备问题" />
                <a-progress :percent="25" label="会员咨询" />
                <a-progress :percent="20" label="支付问题" />
                <a-progress :percent="15" label="功能建议" />
                <a-progress :percent="10" label="其他" />
              </a-card>
            </a-col>
            <a-col :span="12">
              <a-card title="解决效率">
                <a-statistic title="平均响应时间" :value="2" suffix="小时" />
                <a-statistic title="平均解决时间" :value="24" suffix="小时" />
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
        
        <a-tab-pane key="auto-reply" title="自动回复">
          <a-form layout="vertical" style="max-width: 600px;">
            <a-form-item label="启用自动回复">
              <a-switch v-model="settings.autoReply" />
            </a-form-item>
            <a-form-item label="关键词回复规则">
              <a-textarea v-model="settings.rules" :rows="6" placeholder="格式: 关键词|回复内容" />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="handleSave">保存设置</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <a-modal v-model:visible="createVisible" title="创建工单" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="工单类型">
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="device">设备问题</a-option>
            <a-option value="account">账号问题</a-option>
            <a-option value="payment">支付问题</a-option>
            <a-option value="suggestion">功能建议</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="优先级">
          <a-select v-model="form.priority" placeholder="选择优先级">
            <a-option value="low">低</a-option>
            <a-option value="medium">中</a-option>
            <a-option value="high">高</a-option>
            <a-option value="urgent">紧急</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="问题描述">
          <a-textarea v-model="form.description" :rows="4" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 12 });
const createVisible = ref(false);
const settings = reactive({ autoReply: true, rules: '' });
const form = reactive({ type: '', priority: '', description: '' });

const tickets = ref([
  { id: 'TK001', title: '设备无法连接', type: 'device', typeText: '设备问题', priority: 'high', priorityText: '高', status: 'open', statusText: '待处理', assignee: '客服A', createdAt: '2026-03-28 10:00:00' },
  { id: 'TK002', title: '如何绑定设备', type: 'account', typeText: '账号问题', priority: 'low', priorityText: '低', status: 'resolved', statusText: '已解决', assignee: '客服B', createdAt: '2026-03-27 15:00:00' },
]);

const ticketColumns = [
  { title: '工单ID', dataIndex: 'id', width: 100 },
  { title: '标题', dataIndex: 'title', width: 200 },
  { title: '类型', dataIndex: 'typeText', width: 100 },
  { title: '优先级', slotName: 'priority', width: 80 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '负责人', dataIndex: 'assignee', width: 100 },
  { title: '创建时间', dataIndex: 'createdAt', width: 160 },
];

const getPriorityColor = (p: string) => ({ low: 'green', medium: 'blue', high: 'orange', urgent: 'red' }[p] || 'default');
const getStatusColor = (s: string) => ({ open: 'orange', 'in-progress': 'blue', resolved: 'green', closed: 'gray' }[s] || 'default');

const handleCreate = () => { createVisible.value = true; };
const handleSave = () => {};
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.customer-service-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
