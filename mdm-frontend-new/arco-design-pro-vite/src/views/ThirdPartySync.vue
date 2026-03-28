<template>
  <div class="third-party-sync-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="同步任务" :value="25" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="同步记录" :value="5680" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="成功同步" :value="5564" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="第三方平台" :value="8" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>第三方数据同步</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            添加同步
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="syncs" :loading="loading" :pagination="pagination">
        <template #platform="{ record }">
          <a-tag>{{ record.platformName }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-badge :status="record.status === 'active' ? 'success' : 'default'" :text="record.statusText" />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleSync(record)">同步</a-link>
            <a-link @click="handleConfig(record)">配置</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="configVisible" title="同步配置" @before-ok="handleSubmitConfig">
      <a-form :model="configForm" layout="vertical">
        <a-form-item label="API Key">
          <a-input-password v-model="configForm.apiKey" />
        </a-form-item>
        <a-form-item label="同步频率">
          <a-select v-model="configForm.frequency">
            <a-option value="realtime">实时</a-option>
            <a-option value="hourly">每小时</a-option>
            <a-option value="daily">每天</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="数据映射">
          <a-textarea v-model="configForm.mapping" :rows="4" placeholder="JSON格式映射规则" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const pagination = reactive({ current: 1, pageSize: 10, total: 6 });
const configVisible = ref(false);

const configForm = reactive({ apiKey: '', frequency: 'daily', mapping: '' });

const syncs = ref([
  { id: 1, platform: 'wechat', platformName: '微信', dataType: '会员', status: 'active', statusText: '启用', lastSync: '2026-03-28 18:00:00' },
  { id: 2, platform: 'alipay', platformName: '支付宝', dataType: '订单', status: 'active', statusText: '启用', lastSync: '2026-03-28 17:00:00' },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '平台', slotName: 'platform', width: 120 },
  { title: '数据类型', dataIndex: 'dataType', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '最后同步', dataIndex: 'lastSync', width: 160 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' },
];

const handleCreate = () => {};
const handleSync = (record: any) => {};
const handleConfig = (record: any) => { configVisible.value = true; };
const handleDelete = (record: any) => {};
const handleSubmitConfig = (done: boolean) => { done(true); configVisible.value = false; };
</script>

<style scoped>
.third-party-sync-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
