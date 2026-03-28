<template>
  <div class="device-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>设备管理</span>
          <a-space>
            <a-button type="primary" @click="handleRegister">
              <template #icon><icon-plus /></template>
              注册设备
            </a-button>
            <a-button @click="handleExport">
              <template #icon><icon-download /></template>
              导出
            </a-button>
          </a-space>
        </div>
      </template>
      
      <!-- 搜索筛选区 -->
      <div class="search-area">
        <a-row :gutter="16">
          <a-col :span="6">
            <a-input v-model="searchForm.keyword" placeholder="搜索设备ID/名称" allow-clear @press-enter="handleSearch">
              <template #prefix><icon-search /></template>
            </a-input>
          </a-col>
          <a-col :span="4">
            <a-select v-model="searchForm.status" placeholder="在线状态" allow-clear>
              <a-option value="online">在线</a-option>
              <a-option value="offline">离线</a-option>
            </a-select>
          </a-col>
          <a-col :span="4">
            <a-select v-model="searchForm.type" placeholder="设备类型" allow-clear>
              <a-option value="M5Stack">M5Stack</a-option>
              <a-option value="ESP32">ESP32</a-option>
              <a-option value="RaspberryPi">树莓派</a-option>
            </a-select>
          </a-col>
          <a-col :span="4">
            <a-select v-model="searchForm.firmware" placeholder="固件版本" allow-clear>
              <a-option value="v1.0.0">v1.0.0</a-option>
              <a-option value="v1.1.0">v1.1.0</a-option>
              <a-option value="v2.0.0">v2.0.0</a-option>
            </a-select>
          </a-col>
          <a-col :span="6">
            <a-button type="primary" @click="handleSearch">
              <template #icon><icon-search /></template>
              筛选
            </a-button>
            <a-button style="margin-left: 8px;" @click="handleReset">重置</a-button>
          </a-col>
        </a-row>
      </div>
      
      <!-- 数据表格 -->
      <a-table 
        :columns="columns" 
        :data="data" 
        :loading="loading"
        :pagination="pagination"
        :row-selection="{ type: 'checkbox', showCheckedAll: true }"
        @page-change="handlePageChange"
      >
        <template #status="{ record }">
          <a-tag :color="record.online ? 'green' : 'red'">
            {{ record.online ? '在线' : '离线' }}
          </a-tag>
        </template>
        <template #battery="{ record }">
          <a-progress :percent="record.battery" :color="getBatteryColor(record.battery)" size="small" />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleView(record)">详情</a-link>
            <a-link @click="handleCommand(record)">指令</a-link>
            <a-link @click="handleOTA(record)">OTA</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 设备详情弹窗 -->
    <a-modal v-model:visible="detailVisible" title="设备详情" :width="700">
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="设备ID">{{ currentDevice.deviceId }}</a-descriptions-item>
        <a-descriptions-item label="设备名称">{{ currentDevice.name }}</a-descriptions-item>
        <a-descriptions-item label="设备类型">{{ currentDevice.type }}</a-descriptions-item>
        <a-descriptions-item label="固件版本">{{ currentDevice.firmware }}</a-descriptions-item>
        <a-descriptions-item label="在线状态">
          <a-tag :color="currentDevice.online ? 'green' : 'red'">
            {{ currentDevice.online ? '在线' : '离线' }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="电池电量">
          <a-progress :percent="currentDevice.battery || 0" :color="getBatteryColor(currentDevice.battery)" />
        </a-descriptions-item>
        <a-descriptions-item label="最后活跃">{{ currentDevice.lastSeen }}</a-descriptions-item>
        <a-descriptions-item label="运行时长">{{ currentDevice.uptime }}</a-descriptions-item>
        <a-descriptions-item label="创建时间" :span="2">{{ currentDevice.createdAt }}</a-descriptions-item>
      </a-descriptions>
      <template #footer>
        <a-button @click="detailVisible = false">关闭</a-button>
        <a-button type="primary" @click="handleCommand(currentDevice)">发送指令</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const data = ref([]);
const detailVisible = ref(false);
const currentDevice = ref<any>({});

const searchForm = reactive({
  keyword: '',
  status: '',
  type: '',
  firmware: '',
});

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
});

const columns = [
  { title: '设备ID', dataIndex: 'deviceId', width: 120 },
  { title: '设备名称', dataIndex: 'name', width: 150 },
  { title: '类型', dataIndex: 'type', width: 100 },
  { title: '固件版本', dataIndex: 'firmware', width: 100 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '电池', slotName: 'battery', width: 100 },
  { title: '最后活跃', dataIndex: 'lastSeen', width: 160 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' },
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
      data.value = result.data?.list || generateMockData();
      pagination.total = result.data?.total || data.value.length;
    } else {
      data.value = generateMockData();
      pagination.total = data.value.length;
    }
  } catch (error) {
    console.error('获取设备列表失败:', error);
    data.value = generateMockData();
    pagination.total = data.value.length;
  } finally {
    loading.value = false;
  }
};

const generateMockData = () => {
  return [
    { deviceId: 'DEV001', name: '小黄-客厅', type: 'M5Stack', firmware: 'v2.0.0', online: true, battery: 85, lastSeen: '2026-03-28 09:55:00', uptime: '15天3小时', createdAt: '2026-03-01 10:00:00' },
    { deviceId: 'DEV002', name: '小红-卧室', type: 'M5Stack', firmware: 'v1.1.0', online: true, battery: 62, lastSeen: '2026-03-28 09:50:00', uptime: '8天7小时', createdAt: '2026-03-15 14:30:00' },
    { deviceId: 'DEV003', name: '小白-书房', type: 'ESP32', firmware: 'v2.0.0', online: false, battery: 23, lastSeen: '2026-03-27 18:00:00', uptime: '离线', createdAt: '2026-02-20 09:00:00' },
    { deviceId: 'DEV004', name: '小绿-阳台', type: 'M5Stack', firmware: 'v1.0.0', online: true, battery: 100, lastSeen: '2026-03-28 09:55:00', uptime: '30天5小时', createdAt: '2026-02-25 16:00:00' },
  ];
};

const handleSearch = () => {
  pagination.current = 1;
  fetchData();
};

const handleReset = () => {
  Object.assign(searchForm, { keyword: '', status: '', type: '', firmware: '' });
  handleSearch();
};

const handlePageChange = (page: number) => {
  pagination.current = page;
  fetchData();
};

const handleRegister = () => {
  // TODO: 打开注册设备弹窗
};

const handleExport = () => {
  // TODO: 导出数据
};

const handleView = (record: any) => {
  currentDevice.value = record;
  detailVisible.value = true;
};

const handleCommand = (record: any) => {
  currentDevice.value = record;
  detailVisible.value = false;
  // TODO: 打开指令发送弹窗
};

const handleOTA = (record: any) => {
  // TODO: 打开OTA升级弹窗
};

const handleDelete = (record: any) => {
  // TODO: 删除设备
};

const getBatteryColor = (battery: number) => {
  if (battery >= 60) return '#00B42A';
  if (battery >= 30) return '#FF7D00';
  return '#F53F3F';
};

fetchData();
</script>

<style scoped>
.device-container {
  padding: 20px;
}
.card-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.search-area {
  margin-bottom: 16px;
  padding: 16px;
  background: #f7f8fa;
  border-radius: 4px;
}
</style>
