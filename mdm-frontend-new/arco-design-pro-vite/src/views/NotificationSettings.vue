<template>
  <div class="notification-settings-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="通知规则" :value="36" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="已启用" :value="28" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="今日发送" :value="568" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="送达率" :value="98" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>通知渠道配置</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建规则
          </a-button>
        </div>
      </template>
      
      <a-tabs>
        <a-tab-pane key="rules" title="通知规则">
          <a-table :columns="ruleColumns" :data="rules" :pagination="pagination">
            <template #channel="{ record }">
              <a-tag>{{ record.channelText }}</a-tag>
            </template>
            <template #status="{ record }">
              <a-switch v-model="record.enabled" @change="handleToggle(record)" />
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="channels" title="渠道管理">
          <a-row :gutter="16">
            <a-col :span="8" v-for="ch in channels" :key="ch.id">
              <a-card class="channel-card">
                <div class="channel-icon">{{ ch.icon }}</div>
                <div class="channel-name">{{ ch.name }}</div>
                <div class="channel-status">
                  <a-badge :status="ch.enabled ? 'success' : 'default'" :text="ch.enabled ? '已启用' : '已禁用'" />
                </div>
                <a-button @click="handleConfigChannel(ch)">配置</a-button>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
        
        <a-tab-pane key="history" title="发送历史">
          <a-table :columns="historyColumns" :data="history" :pagination="pagination" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 8 });

const rules = ref([
  { id: 1, name: '设备离线告警', channel: 'sms', channelText: '短信', event: 'device.offline', enabled: true },
  { id: 2, name: '会员注册通知', channel: 'email', channelText: '邮件', event: 'member.registered', enabled: true },
  { id: 3, name: '订单创建提醒', channel: 'app', channelText: 'App推送', event: 'order.created', enabled: false },
]);

const channels = ref([
  { id: 1, name: '短信', icon: '📱', enabled: true },
  { id: 2, name: '邮件', icon: '📧', enabled: true },
  { id: 3, name: 'App推送', icon: '🔔', enabled: true },
  { id: 4, name: '微信', icon: '💬', enabled: false },
]);

const history = ref([
  { id: 1, channel: '短信', template: '设备离线告警', recipient: '138****8888', status: 'sent', time: '2026-03-28 18:00:00' },
]);

const ruleColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '规则名称', dataIndex: 'name', width: 150 },
  { title: '渠道', slotName: 'channel', width: 100 },
  { title: '触发事件', dataIndex: 'event', width: 150 },
  { title: '状态', slotName: 'status', width: 100 },
];

const historyColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '渠道', dataIndex: 'channel', width: 100 },
  { title: '模板', dataIndex: 'template', width: 150 },
  { title: '接收人', dataIndex: 'recipient', width: 120 },
  { title: '状态', dataIndex: 'status', width: 100 },
  { title: '时间', dataIndex: 'time', width: 160 },
];

const handleCreate = () => {};
const handleToggle = (record: any) => {};
const handleConfigChannel = (ch: any) => {};
</script>

<style scoped>
.notification-settings-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
.channel-card { text-align: center; }
.channel-icon { font-size: 36px; margin-bottom: 8px; }
.channel-name { font-weight: bold; margin-bottom: 8px; }
.channel-status { margin-bottom: 12px; }
</style>
