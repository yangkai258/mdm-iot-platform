<template>
  <div class="edge-nodes-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card>
          <a-statistic title="边缘节点数" :value="stats.totalNodes" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="在线节点" :value="stats.onlineNodes" status="success" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="离线节点" :value="stats.offlineNodes" status="error" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="平均负载" :value="stats.avgLoad" suffix="%" />
        </a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>边缘节点管理</span>
          <a-button type="primary" @click="handleAdd">
            <template #icon><icon-plus /></template>
            添加节点
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #status="{ record }">
          <a-tag :color="record.online ? 'green' : 'red'">{{ record.online ? '在线' : '离线' }}</a-tag>
        </template>
        <template #load="{ record }">
          <a-progress :percent="record.load" :color="getLoadColor(record.load)" size="small" />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleView(record)">详情</a-link>
            <a-link @click="handleConfig(record)">配置</a-link>
            <a-link v-if="!record.online" @click="handleOnline(record)">上线</a-link>
            <a-link v-if="record.online" status="danger" @click="handleOffline(record)">下线</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 添加节点弹窗 -->
    <a-modal v-model:visible="addVisible" title="添加边缘节点" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="节点名称" required>
          <a-input v-model="form.name" placeholder="如: Edge-CN-01" />
        </a-form-item>
        <a-form-item label="节点地址" required>
          <a-input v-model="form.address" placeholder="IP或域名" />
        </a-form-item>
        <a-form-item label="所属区域">
          <a-select v-model="form.region" placeholder="选择区域">
            <a-option value="CN">中国大陆</a-option>
            <a-option value="US">北美</a-option>
            <a-option value="EU">欧洲</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="节点类型">
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="compute">计算节点</a-option>
            <a-option value="storage">存储节点</a-option>
            <a-option value="mixed">混合节点</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);

const stats = reactive({
  totalNodes: 25,
  onlineNodes: 22,
  offlineNodes: 3,
  avgLoad: 45,
});

const data = ref([
  { id: 'N001', name: 'Edge-CN-01', address: '10.0.1.1', region: 'CN', type: 'mixed', online: true, load: 45, cpu: '12%', memory: '8GB', lastSeen: '2026-03-28 18:00:00' },
  { id: 'N002', name: 'Edge-CN-02', address: '10.0.1.2', region: 'CN', type: 'compute', online: true, load: 68, cpu: '25%', memory: '16GB', lastSeen: '2026-03-28 18:00:00' },
  { id: 'N003', name: 'Edge-US-01', address: '10.0.2.1', region: 'US', type: 'mixed', online: false, load: 0, cpu: '0%', memory: '0GB', lastSeen: '2026-03-27 10:00:00' },
  { id: 'N004', name: 'Edge-EU-01', address: '10.0.3.1', region: 'EU', type: 'storage', online: true, load: 32, cpu: '8%', memory: '4GB', lastSeen: '2026-03-28 18:00:00' },
]);

const pagination = reactive({ current: 1, pageSize: 10, total: 4 });
const addVisible = ref(false);

const form = reactive({ name: '', address: '', region: 'CN', type: 'mixed' });

const columns = [
  { title: '节点ID', dataIndex: 'id', width: 100 },
  { title: '节点名称', dataIndex: 'name', width: 120 },
  { title: '地址', dataIndex: 'address', width: 120 },
  { title: '区域', dataIndex: 'region', width: 100 },
  { title: '类型', dataIndex: 'type', width: 100 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '负载', slotName: 'load', width: 120 },
  { title: 'CPU', dataIndex: 'cpu', width: 80 },
  { title: '内存', dataIndex: 'memory', width: 80 },
  { title: '最后活跃', dataIndex: 'lastSeen', width: 160 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' },
];

const getLoadColor = (load: number) => load < 60 ? '#00B42A' : load < 80 ? '#FF7D00' : '#F53F3F';

const handleAdd = () => { addVisible.value = true; };
const handleView = (record: any) => {};
const handleConfig = (record: any) => {};
const handleOnline = (record: any) => {};
const handleOffline = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); addVisible.value = false; };
</script>

<style scoped>
.edge-nodes-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
