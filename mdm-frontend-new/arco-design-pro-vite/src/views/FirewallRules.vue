<template>
  <div class="firewall-rules-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="防火墙规则" :value="256" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="启用规则" :value="240" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="拦截次数" :value="5680" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="阻断率" :value="2.5" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>防火墙规则管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建规则
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="rules" :loading="loading" :pagination="pagination">
        <template #action="{ record }">
          <a-tag>{{ record.actionText }}</a-tag>
        </template>
        <template #enabled="{ record }">
          <a-switch v-model="record.enabled" @change="handleToggle(record)" />
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="createVisible" title="创建防火墙规则" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="规则名称" required>
          <a-input v-model="form.name" placeholder="请输入规则名称" />
        </a-form-item>
        <a-form-item label="协议">
          <a-select v-model="form.protocol">
            <a-option value="tcp">TCP</a-option>
            <a-option value="udp">UDP</a-option>
            <a-option value="all">全部</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="端口">
          <a-input v-model="form.port" placeholder="如: 80,443" />
        </a-form-item>
        <a-form-item label="动作">
          <a-select v-model="form.action">
            <a-option value="allow">允许</a-option>
            <a-option value="deny">拒绝</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="form.note" :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const pagination = reactive({ current: 1, pageSize: 10, total: 10 });
const createVisible = ref(false);

const form = reactive({ name: '', protocol: '', port: '', action: '', note: '' });

const rules = ref([
  { id: 1, name: '允许HTTP', protocol: 'tcp', port: '80', action: 'allow', actionText: '允许', enabled: true, hitCount: 25600 },
  { id: 2, name: '允许HTTPS', protocol: 'tcp', port: '443', action: 'allow', actionText: '允许', enabled: true, hitCount: 18600 },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '规则名称', dataIndex: 'name', width: 150 },
  { title: '协议', dataIndex: 'protocol', width: 80 },
  { title: '端口', dataIndex: 'port', width: 100 },
  { title: '动作', slotName: 'action', width: 80 },
  { title: '命中次数', dataIndex: 'hitCount', width: 100 },
  { title: '启用', slotName: 'enabled', width: 80 },
];

const handleCreate = () => { createVisible.value = true; };
const handleToggle = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.firewall-rules-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
