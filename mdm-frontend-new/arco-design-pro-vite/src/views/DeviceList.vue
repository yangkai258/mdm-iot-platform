<template>
  <div class="container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>设备管理</span>
          <a-button type="primary">注册设备</a-button>
        </div>
      </template>
      
      <a-table 
        :columns="columns" 
        :data="data" 
        :loading="loading"
        :pagination="{ total: total, current: current, pageSize: pageSize }"
        @page-change="handlePageChange"
      >
        <template #status="{ record }">
          <a-tag :color="record.online ? 'green' : 'red'">
            {{ record.online ? '在线' : '离线' }}
          </a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleView(record)">详情</a-link>
            <a-link @click="handleCommand(record)">指令</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';

const loading = ref(false);
const data = ref([]);
const total = ref(0);
const current = ref(1);
const pageSize = ref(20);

const columns = [
  { title: '设备ID', dataIndex: 'device_id' },
  { title: '设备名称', dataIndex: 'device_name' },
  { title: '型号', dataIndex: 'model' },
  { title: '状态', dataIndex: 'status', slotName: 'status' },
  { title: '最后活跃', dataIndex: 'last_seen' },
  { title: '操作', slotName: 'actions' },
];

const fetchData = async () => {
  loading.value = true;
  try {
    const response = await fetch('/api/v1/devices', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token') || ''}`
      }
    });
    const result = await response.json();
    if (result.code === 0) {
      data.value = result.data?.list || [];
      total.value = result.data?.total || 0;
    }
  } catch (error) {
    console.error('获取设备列表失败:', error);
  } finally {
    loading.value = false;
  }
};

const handlePageChange = (page: number) => {
  current.value = page;
  fetchData();
};

const handleView = (record: any) => {
  console.log('查看设备:', record.device_id);
};

const handleCommand = (record: any) => {
  console.log('发送指令:', record.device_id);
};

onMounted(() => {
  fetchData();
});
</script>

<style scoped>
.container {
  padding: 20px;
}
.card-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
