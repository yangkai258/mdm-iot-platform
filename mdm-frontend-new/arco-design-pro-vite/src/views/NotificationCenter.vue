<template>
  <div class="notification-center-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="未读通知" :value="15" status="warning" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="今日通知" :value="128" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="消息模板" :value="36" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>通知中心</span>
          <a-space>
            <a-button @click="handleReadAll">全部标为已读</a-button>
            <a-button type="primary" @click="handleCreateTemplate">
              <template #icon><icon-plus /></template>
              创建模板
            </a-button>
          </a-space>
        </div>
      </template>
      
      <a-tabs>
        <a-tab-pane key="list" title="通知列表">
          <a-table :columns="columns" :data="notifications" :pagination="pagination">
            <template #type="{ record }">
              <a-tag :color="getTypeColor(record.type)">{{ record.typeText }}</a-tag>
            </template>
            <template #read="{ record }">
              <a-switch v-model="record.isRead" @change="handleReadChange(record)" />
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="templates" title="消息模板">
          <a-table :columns="templateColumns" :data="templates" :pagination="pagination">
            <template #actions="{ record }">
              <a-space>
                <a-link @click="handleEditTemplate(record)">编辑</a-link>
                <a-link status="danger" @click="handleDeleteTemplate(record)">删除</a-link>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="settings" title="通知设置">
          <a-form layout="vertical" style="max-width: 600px;">
            <a-form-item label="邮件通知">
              <a-switch v-model="settings.email" />
            </a-form-item>
            <a-form-item label="短信通知">
              <a-switch v-model="settings.sms" />
            </a-form-item>
            <a-form-item label="App推送">
              <a-switch v-model="settings.appPush" />
            </a-form-item>
            <a-form-item label="微信通知">
              <a-switch v-model="settings.wechat" />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="handleSaveSettings">保存设置</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <a-modal v-model:visible="templateVisible" title="消息模板" @before-ok="handleSubmitTemplate">
      <a-form :model="templateForm" layout="vertical">
        <a-form-item label="模板名称" required>
          <a-input v-model="templateForm.name" placeholder="请输入模板名称" />
        </a-form-item>
        <a-form-item label="通知类型">
          <a-select v-model="templateForm.type" placeholder="选择类型">
            <a-option value="system">系统通知</a-option>
            <a-option value="alert">告警通知</a-option>
            <a-option value="marketing">营销通知</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="模板内容">
          <a-textarea v-model="templateForm.content" :rows="4" placeholder="支持变量: {{name}}, {{time}}" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 15 });
const templateVisible = ref(false);

const settings = reactive({ email: true, sms: true, appPush: true, wechat: false });
const templateForm = reactive({ name: '', type: '', content: '' });

const notifications = ref([
  { id: 1, type: 'alert', typeText: '告警', title: '设备离线告警', content: '设备DEV001已离线超过5分钟', isRead: false, time: '2026-03-28 18:00:00' },
  { id: 2, type: 'system', typeText: '系统', title: '系统更新', content: '系统将于今晚23:00进行维护', isRead: true, time: '2026-03-28 10:00:00' },
  { id: 3, type: 'marketing', typeText: '营销', title: '新功能上线', content: 'AI训练师功能已上线', isRead: false, time: '2026-03-27 15:00:00' },
]);

const templates = ref([
  { id: 1, name: '设备离线提醒', type: 'alert', typeText: '告警', content: '您的设备{{device_name}}已离线{{duration}}', createdAt: '2026-03-20' },
  { id: 2, name: '会员生日祝福', type: 'marketing', typeText: '营销', content: '亲爱的{{name}},祝您生日快乐！', createdAt: '2026-03-15' },
]);

const columns = [
  { title: '类型', slotName: 'type', width: 100 },
  { title: '标题', dataIndex: 'title', width: 200 },
  { title: '内容', dataIndex: 'content' },
  { title: '已读', slotName: 'read', width: 80 },
  { title: '时间', dataIndex: 'time', width: 160 },
];

const templateColumns = [
  { title: '模板名称', dataIndex: 'name', width: 150 },
  { title: '类型', dataIndex: 'typeText', width: 100 },
  { title: '模板内容', dataIndex: 'content' },
  { title: '创建时间', dataIndex: 'createdAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 120 },
];

const getTypeColor = (t: string) => ({ alert: 'red', system: 'blue', marketing: 'green' }[t] || 'default');

const handleReadAll = () => { notifications.value.forEach(n => n.isRead = true); };
const handleReadChange = (record: any) => {};
const handleCreateTemplate = () => { templateVisible.value = true; };
const handleEditTemplate = (record: any) => { Object.assign(templateForm, record); templateVisible.value = true; };
const handleDeleteTemplate = (record: any) => {};
const handleSaveSettings = () => {};
const handleSubmitTemplate = (done: boolean) => { done(true); templateVisible.value = false; };
</script>

<style scoped>
.notification-center-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
