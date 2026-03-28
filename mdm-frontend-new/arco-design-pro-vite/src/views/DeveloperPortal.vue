<template>
  <div class="developer-portal-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card>
          <a-statistic title="注册应用" :value="stats.apps" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="API调用量" :value="stats.apiCalls" suffix="次" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="认证开发者" :value="stats.developers" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="Token数" :value="stats.tokens" />
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16">
      <a-col :span="12">
        <a-card title="应用管理">
          <template #extra>
            <a-button type="primary" @click="handleCreateApp">
              <template #icon><icon-plus /></template>
              创建应用
            </a-button>
          </template>
          <a-table :columns="appColumns" :data="apps" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="record.enabled ? 'green' : 'gray'">{{ record.enabled ? '启用' : '禁用' }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-link @click="handleViewApp(record)">详情</a-link>
                <a-link @click="handleManageToken(record)">Token</a-link>
                <a-link @click="handleAppSettings(record)">设置</a-link>
              </a-space>
            </template>
          </a-table>
        </a-card>
      </a-col>
      
      <a-col :span="12">
        <a-card title="API文档">
          <a-tabs>
            <a-tab-pane key="device" title="设备API">
              <a-descriptions :column="1" size="small">
                <a-descriptions-item label="注册设备">POST /api/v1/devices</a-descriptions-item>
                <a-descriptions-item label="查询设备">GET /api/v1/devices/:id</a-descriptions-item>
                <a-descriptions-item label="下发指令">POST /api/v1/devices/:id/commands</a-descriptions-item>
                <a-descriptions-item label="查询状态">GET /api/v1/devices/:id/status</a-descriptions-item>
              </a-descriptions>
            </a-tab-pane>
            <a-tab-pane key="pet" title="宠物API">
              <a-descriptions :column="1" size="small">
                <a-descriptions-item label="创建宠物">POST /api/v1/pets</a-descriptions-item>
                <a-descriptions-item label="查询宠物">GET /api/v1/pets/:id</a-descriptions-item>
                <a-descriptions-item label="更新宠物">PUT /api/v1/pets/:id</a-descriptions-item>
              </a-descriptions>
            </a-tab-pane>
            <a-tab-pane key="member" title="会员API">
              <a-descriptions :column="1" size="small">
                <a-descriptions-item label="创建会员">POST /api/v1/members</a-descriptions-item>
                <a-descriptions-item label="查询会员">GET /api/v1/members/:id</a-descriptions-item>
              </a-descriptions>
            </a-tab-pane>
          </a-tabs>
          <a-button type="primary" style="margin-top: 16px;">查看完整文档</a-button>
        </a-card>
      </a-col>
    </a-row>

    <a-card title="API调用日志" style="margin-top: 16px;">
      <a-table :columns="logColumns" :data="apiLogs" :pagination="paginationSmall">
        <template #status="{ record }">
          <a-tag :color="record.status < 400 ? 'green' : 'red'">{{ record.status }}</a-tag>
        </template>
      </a-table>
    </a-card>

    <!-- 创建应用弹窗 -->
    <a-modal v-model:visible="createVisible" title="创建应用" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="应用名称" required>
          <a-input v-model="form.name" placeholder="请输入应用名称" />
        </a-form-item>
        <a-form-item label="应用描述">
          <a-textarea v-model="form.description" :rows="2" />
        </a-form-item>
        <a-form-item label="应用类型">
          <a-select v-model="form.type" placeholder="选择应用类型">
            <a-option value="web">Web应用</a-option>
            <a-option value="mobile">移动应用</a-option>
            <a-option value="iot">IoT设备</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="回调URL">
          <a-input v-model="form.callbackUrl" placeholder="https://your-app.com/callback" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- Token管理弹窗 -->
    <a-modal v-model:visible="tokenVisible" title="Token管理" :width="700">
      <a-alert>请妥善保管您的Token，不要泄露给他人。</a-alert>
      <a-divider />
      <div class="token-list">
        <div v-for="token in tokens" :key="token.id" class="token-item">
          <div class="token-name">{{ token.name }}</div>
          <div class="token-value">{{ token.masked }}</div>
          <div class="token-meta">
            <span>创建于: {{ token.createdAt }}</span>
            <span>最后使用: {{ token.lastUsed }}</span>
          </div>
          <a-space>
            <a-link @click="handleRegenerateToken(token)">重新生成</a-link>
            <a-link status="danger" @click="handleRevokeToken(token)">撤销</a-link>
          </a-space>
        </div>
      </div>
      <a-button type="primary" style="margin-top: 16px;" @click="handleCreateToken">创建Token</a-button>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const stats = reactive({ apps: 156, apiCalls: 2586320, developers: 89, tokens: 234 });

const apps = ref([
  { id: 1, name: '我的IoT应用', type: 'iot', description: '智能家居控制', appKey: 'ak_****', enabled: true, calls: 156320 },
  { id: 2, name: '宠物健康App', type: 'mobile', description: '宠物健康管理', appKey: 'ak_****', enabled: true, calls: 89320 },
  { id: 3, name: '测试应用', type: 'web', description: '开发测试用', appKey: 'ak_****', enabled: false, calls: 1230 },
]);

const apiLogs = ref([
  { id: 1, appName: '我的IoT应用', method: 'POST', path: '/api/v1/devices', status: 200, latency: 45, time: '2026-03-28 18:50:00' },
  { id: 2, appName: '宠物健康App', method: 'GET', path: '/api/v1/pets/1', status: 200, latency: 32, time: '2026-03-28 18:49:30' },
  { id: 3, appName: '测试应用', method: 'POST', path: '/api/v1/commands', status: 401, latency: 15, time: '2026-03-28 18:49:00' },
]);

const tokens = ref([
  { id: 1, name: '生产环境Token', masked: 'tk_live_****xxxx', createdAt: '2026-03-01', lastUsed: '2026-03-28 18:00:00' },
  { id: 2, name: '测试环境Token', masked: 'tk_test_****xxxx', createdAt: '2026-03-15', lastUsed: '2026-03-27 10:00:00' },
]);

const pagination = reactive({ current: 1, pageSize: 10, total: 3 });
const paginationSmall = reactive({ current: 1, pageSize: 5, total: 3 });
const createVisible = ref(false);
const tokenVisible = ref(false);

const form = reactive({ name: '', description: '', type: 'web', callbackUrl: '' });

const appColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '应用名称', dataIndex: 'name', width: 150 },
  { title: '类型', dataIndex: 'type', width: 100 },
  { title: 'AppKey', dataIndex: 'appKey', width: 150 },
  { title: '调用量', dataIndex: 'calls', width: 100 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' },
];

const logColumns = [
  { title: '时间', dataIndex: 'time', width: 160 },
  { title: '应用', dataIndex: 'appName', width: 150 },
  { title: '方法', dataIndex: 'method', width: 80 },
  { title: '路径', dataIndex: 'path', width: 200 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '延迟(ms)', dataIndex: 'latency', width: 100 },
];

const handleCreateApp = () => { createVisible.value = true; };
const handleViewApp = (record: any) => {};
const handleManageToken = (record: any) => { tokenVisible.value = true; };
const handleAppSettings = (record: any) => {};
const handleCreateToken = () => {};
const handleRegenerateToken = (token: any) => {};
const handleRevokeToken = (token: any) => {};
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.developer-portal-container { padding: 20px; }
.token-list { margin-top: 16px; }
.token-item { padding: 16px; border: 1px solid #e5e6e8; border-radius: 4px; margin-bottom: 12px; }
.token-name { font-weight: bold; }
.token-value { font-family: monospace; background: #f7f8fa; padding: 8px; margin: 8px 0; }
.token-meta { font-size: 12px; color: #86909c; margin-bottom: 8px; }
</style>
