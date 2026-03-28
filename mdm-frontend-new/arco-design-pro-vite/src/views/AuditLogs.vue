<template>
  <div class="audit-logs-container">
    <a-card>
      <template #title>
        <span>合规审计日志</span>
      </template>
      
      <div class="search-area">
        <a-row :gutter="16">
          <a-col :span="5">
            <a-input v-model="searchForm.keyword" placeholder="搜索操作人/内容" allow-clear />
          </a-col>
          <a-col :span="4">
            <a-select v-model="searchForm.category" placeholder="操作类别" allow-clear>
              <a-option value="data_access">数据访问</a-option>
              <a-option value="data_export">数据导出</a-option>
              <a-option value="data_delete">数据删除</a-option>
              <a-option value="permission">权限变更</a-option>
              <a-option value="login">登录</a-option>
            </a-select>
          </a-col>
          <a-col :span="4">
            <a-select v-model="searchForm.riskLevel" placeholder="风险级别" allow-clear>
              <a-option value="low">低</a-option>
              <a-option value="medium">中</a-option>
              <a-option value="high">高</a-option>
              <a-option value="critical">严重</a-option>
            </a-select>
          </a-col>
          <a-col :span="2">
            <a-button type="primary" @click="handleSearch">筛选</a-button>
          </a-col>
          <a-col :span="2">
            <a-button @click="handleExport">导出</a-button>
          </a-col>
        </a-row>
      </div>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #category="{ record }">
          <a-tag :color="getCategoryColor(record.category)">{{ record.categoryText }}</a-tag>
        </template>
        <template #riskLevel="{ record }">
          <a-tag :color="getRiskColor(record.riskLevel)">{{ record.riskLevelText }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleView(record)">详情</a-link>
            <a-link v-if="record.riskLevel === 'high' || record.riskLevel === 'critical'" status="danger" @click="handleAlert">
              告警
            </a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 详情弹窗 -->
    <a-modal v-model:visible="detailVisible" title="审计详情" :width="700">
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="日志ID">{{ currentLog.id }}</a-descriptions-item>
        <a-descriptions-item label="操作时间">{{ currentLog.time }}</a-descriptions-item>
        <a-descriptions-item label="操作用户">{{ currentLog.user }}</a-descriptions-item>
        <a-descriptions-item label="用户角色">{{ currentLog.role }}</a-descriptions-item>
        <a-descriptions-item label="操作类别">
          <a-tag :color="getCategoryColor(currentLog.category)">{{ currentLog.categoryText }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="风险级别">
          <a-tag :color="getRiskColor(currentLog.riskLevel)">{{ currentLog.riskLevelText }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="操作IP" :span="2">{{ currentLog.ip }}</a-descriptions-item>
        <a-descriptions-item label="操作描述" :span="2">{{ currentLog.description }}</a-descriptions-item>
        <a-descriptions-item label="原始数据" :span="2">
          <a-textarea :value="currentLog.rawData" :rows="3" readonly />
        </a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);

const searchForm = reactive({ keyword: '', category: '', riskLevel: '' });
const pagination = reactive({ current: 1, pageSize: 20, total: 5 });

const data = ref([
  { id: 'A001', time: '2026-03-28 18:50:00', user: 'admin', role: '管理员', category: 'data_export', categoryText: '数据导出', ip: '192.168.1.100', description: '导出全部用户数据', riskLevel: 'high', riskLevelText: '高风险', rawData: '{"action":"export","target":"users","format":"csv"}' },
  { id: 'A002', time: '2026-03-28 18:45:00', user: 'operator1', role: '运营', category: 'permission', categoryText: '权限变更', ip: '192.168.1.101', description: '修改用户角色权限', riskLevel: 'medium', riskLevelText: '中风险', rawData: '{"action":"update_role","user_id":5,"old_role":"user","new_role":"admin"}' },
  { id: 'A003', time: '2026-03-28 18:40:00', user: 'user123', role: '普通用户', category: 'data_access', categoryText: '数据访问', ip: '192.168.1.102', description: '查询其他用户订单信息', riskLevel: 'low', riskLevelText: '低风险', rawData: '{"action":"query","table":"orders","filter":"user_id=3"}' },
  { id: 'A004', time: '2026-03-28 18:35:00', user: 'admin', role: '管理员', category: 'data_delete', categoryText: '数据删除', ip: '192.168.1.100', description: '删除30天前设备日志', riskLevel: 'medium', riskLevelText: '中风险', rawData: '{"action":"delete","table":"device_logs","condition":"created_at < 30_days_ago"}' },
]);

const columns = [
  { title: '日志ID', dataIndex: 'id', width: 100 },
  { title: '时间', dataIndex: 'time', width: 160 },
  { title: '用户', dataIndex: 'user', width: 100 },
  { title: '角色', dataIndex: 'role', width: 100 },
  { title: '类别', slotName: 'category', width: 100 },
  { title: '风险级别', slotName: 'riskLevel', width: 100 },
  { title: '描述', dataIndex: 'description' },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const detailVisible = ref(false);
const currentLog = ref<any>({});

const getCategoryColor = (c: string) => ({ data_access: 'blue', data_export: 'orange', data_delete: 'red', permission: 'purple', login: 'green' }[c] || 'default');
const getRiskColor = (r: string) => ({ low: 'green', medium: 'orange', high: 'red', critical: '#F53F3F' }[r] || 'default');

const handleSearch = () => {};
const handleExport = () => {};
const handleView = (record: any) => { currentLog.value = record; detailVisible.value = true; };
const handleAlert = (record: any) => {};
</script>

<style scoped>
.audit-logs-container { padding: 20px; }
.search-area { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
</style>
