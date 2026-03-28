<template>
  <div class="channels-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>通知渠道配置</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新增渠道
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag :color="getTypeColor(record.type)">{{ record.typeText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.enabled ? 'green' : 'gray'">{{ record.enabled ? '启用' : '禁用' }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleTest(record)">测试</a-link>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 编辑弹窗 -->
    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑渠道' : '新增渠道'" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="渠道名称" required>
          <a-input v-model="form.name" placeholder="请输入渠道名称" />
        </a-form-item>
        <a-form-item label="渠道类型" required>
          <a-radio-group v-model="form.type">
            <a-radio value="email">邮件</a-radio>
            <a-radio value="sms">短信</a-radio>
            <a-radio value="webhook">Webhook</a-radio>
            <a-radio value="wechat">企业微信</a-radio>
          </a-radio-group>
        </a-form-item>
        
        <!-- 邮件配置 -->
        <template v-if="form.type === 'email'">
          <a-form-item label="SMTP服务器">
            <a-input v-model="form.config.host" placeholder="smtp.example.com" />
          </a-form-item>
          <a-form-item label="端口">
            <a-input-number v-model="form.config.port" :min="1" :max="65535" />
          </a-form-item>
          <a-form-item label="用户名">
            <a-input v-model="form.config.username" placeholder="user@example.com" />
          </a-form-item>
          <a-form-item label="密码">
            <a-input-password v-model="form.config.password" />
          </a-form-item>
          <a-form-item label="发件人">
            <a-input v-model="form.config.from" placeholder="noreply@example.com" />
          </a-form-item>
        </template>
        
        <!-- 短信配置 -->
        <template v-if="form.type === 'sms'">
          <a-form-item label="短信服务商">
            <a-select v-model="form.config.provider" placeholder="选择服务商">
              <a-option value="aliyun">阿里云</a-option>
              <a-option value="tencent">腾讯云</a-option>
              <a-option value="huawei">华为云</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="AccessKey">
            <a-input v-model="form.config.accessKey" />
          </a-form-item>
          <a-form-item label="AccessSecret">
            <a-input-password v-model="form.config.accessSecret" />
          </a-form-item>
          <a-form-item label="签名">
            <a-input v-model="form.config.signature" placeholder="MDM告警" />
          </a-form-item>
        </template>
        
        <!-- Webhook配置 -->
        <template v-if="form.type === 'webhook'">
          <a-form-item label="Webhook地址">
            <a-input v-model="form.config.url" placeholder="https://your-webhook.com/notify" />
          </a-form-item>
          <a-form-item label="请求方法">
            <a-select v-model="form.config.method">
              <a-option value="POST">POST</a-option>
              <a-option value="GET">GET</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="Secret密钥">
            <a-input-password v-model="form.config.secret" />
          </a-form-item>
          <a-form-item label="Headers">
            <a-textarea v-model="form.config.headers" placeholder='{"Content-Type": "application/json"}' :rows="2" />
          </a-form-item>
        </template>
        
        <!-- 企业微信配置 -->
        <template v-if="form.type === 'wechat'">
          <a-form-item label="企业ID">
            <a-input v-model="form.config.corpId" />
          </a-form-item>
          <a-form-item label="应用Secret">
            <a-input-password v-model="form.config.corpSecret" />
          </a-form-item>
          <a-form-item label="AgentID">
            <a-input v-model="form.config.agentId" />
          </a-form-item>
        </template>
        
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
  { id: 1, name: '管理员邮件', type: 'email', typeText: '邮件', enabled: true, config: { host: 'smtp.exmail.qq.com', port: 465 }, lastTest: '2026-03-28 10:00:00' },
  { id: 2, name: '紧急短信', type: 'sms', typeText: '短信', enabled: true, config: { provider: 'aliyun', signature: 'MDM告警' }, lastTest: '2026-03-28 09:00:00' },
  { id: 3, name: '飞书Webhook', type: 'webhook', typeText: 'Webhook', enabled: true, config: { url: 'https://open.feishu.cn', method: 'POST' }, lastTest: '2026-03-27 15:00:00' },
  { id: 4, name: '企业微信', type: 'wechat', typeText: '企业微信', enabled: false, config: { corpId: 'wwxxx' }, lastTest: '2026-03-25 10:00:00' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 4 });
const editVisible = ref(false);
const isEdit = ref(false);

const form = reactive({
  name: '',
  type: 'email',
  enabled: true,
  config: {
    host: '', port: 465, username: '', password: '', from: '',
    provider: 'aliyun', accessKey: '', accessSecret: '', signature: '',
    url: '', method: 'POST', secret: '', headers: '',
    corpId: '', corpSecret: '', agentId: '',
  },
});

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '渠道名称', dataIndex: 'name', width: 150 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '上次测试', dataIndex: 'lastTest', width: 160 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' },
];

const getTypeColor = (t: string) => ({ email: 'blue', sms: 'green', webhook: 'orange', wechat: 'cyan' }[t] || 'default');

const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleTest = (record: any) => {};
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.channels-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
