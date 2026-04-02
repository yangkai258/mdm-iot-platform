<template>
  <div class="device-status-container">
    <Breadcrumb :items="['首页', '设备管理', '设备状态']" />

    <!-- 状态统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="4">
        <a-card hoverable @click="filterByStatus(1)">
          <a-statistic title="待激活" :value="statusStats[1] || 0" :value-style="{ color: '#1890ff' }" />
        </a-card>
      </a-col>
      <a-col :span="4">
        <a-card hoverable @click="filterByStatus(2)">
          <a-statistic title="服役中" :value="statusStats[2] || 0" :value-style="{ color: '#52c41a' }" />
        </a-card>
      </a-col>
      <a-col :span="4">
        <a-card hoverable @click="filterByStatus(3)">
          <a-statistic title="维修中" :value="statusStats[3] || 0" :value-style="{ color: '#faad14' }" />
        </a-card>
      </a-col>
      <a-col :span="4">
        <a-card hoverable @click="filterByStatus(4)">
          <a-statistic title="已挂失" :value="statusStats[4] || 0" :value-style="{ color: '#ff4d4f' }" />
        </a-card>
      </a-col>
      <a-col :span="4">
        <a-card hoverable @click="filterByStatus(5)">
          <a-statistic title="已报废" :value="statusStats[5] || 0" :value-style="{ color: '#8c8c8c' }" />
        </a-card>
      </a-col>
      <a-col :span="4">
        <a-card hoverable @click="filterByStatus(null)">
          <a-statistic title="全部设备" :value="statusStats.total || 0" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 搜索筛选区 -->
    <a-card class="general-card" style="margin-top: 16px">
      <a-space>
        <a-input-search
          v-model="searchKeyword"
          placeholder="搜索设备ID或MAC地址"
          style="width: 280px"
          @search="handleSearch"
          search-button
        />
        <a-select v-model="filterStatus" placeholder="筛选状态" style="width: 150px" allow-clear @change="handleFilterChange">
          <a-option :value="1">待激活</a-option>
          <a-option :value="2">服役中</a-option>
          <a-option :value="3">维修中</a-option>
          <a-option :value="4">已挂失</a-option>
          <a-option :value="5">已报废</a-option>
        </a-select>
        <a-button @click="loadDevices">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
        <a-button v-if="filterStatus !== null || searchKeyword" @click="clearFilters">清除筛选</a-button>
      </a-space>
    </a-card>

    <!-- 数据表格 -->
    <a-card class="general-card" style="margin-top: 16px">
      <template #title>
        <span class="card-title">设备列表</span>
      </template>

      <a-table
        :columns="columns"
        :data="filteredDevices"
        :loading="loading"
        :pagination="pagination"
        @change="handleTableChange"
        row-key="device_id"
      >
        <template #deviceId="{ record }">
          <a-space>
            <a-avatar :size="24" :style="{ backgroundColor: '#165dff' }">
              {{ record.device_id.charAt(0) }}
            </a-avatar>
            <span>{{ record.device_id }}</span>
          </a-space>
        </template>
        <template #isOnline="{ record }">
          <a-badge :status="record.is_online ? 'success' : 'default'" :text="record.is_online ? '在线' : '离线'" />
        </template>
        <template #lifecycleStatus="{ record }">
          <a-tag :color="getStatusColor(record.lifecycle_status)" style="cursor: pointer" @click="showStatusModal(record)">
            {{ getStatusText(record.lifecycle_status) }}
          </a-tag>
        </template>
        <template #lastUpdate="{ record }">
          {{ formatTime(record.last_update_time) }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showStatusModal(record)">状态</a-button>
            <a-button type="text" size="small" @click="viewDevice(record)">详情</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 状态变更弹窗 -->
    <a-modal v-model:visible="statusModalVisible" title="变更设备状态" @ok="handleStatusChange" :confirm-loading="statusChanging">
      <a-descriptions :column="1" bordered size="small">
        <a-descriptions-item label="设备ID">{{ currentDevice?.device_id }}</a-descriptions-item>
        <a-descriptions-item label="MAC地址">{{ currentDevice?.mac_address }}</a-descriptions-item>
        <a-descriptions-item label="硬件型号">{{ currentDevice?.hardware_model }}</a-descriptions-item>
        <a-descriptions-item label="当前状态">
          <a-tag :color="getStatusColor(currentDevice?.lifecycle_status)">
            {{ getStatusText(currentDevice?.lifecycle_status) }}
          </a-tag>
        </a-descriptions-item>
      </a-descriptions>

      <a-divider>选择新状态</a-divider>

      <a-radio-group v-model="newStatus" direction="vertical">
        <a-radio :value="1">
          <a-tag color="#1890ff">待激活</a-tag>
          <span style="margin-left: 8px; color: #666">设备尚未激活</span>
        </a-radio>
        <a-radio :value="2">
          <a-tag color="#52c41a">服役中</a-tag>
          <span style="margin-left: 8px; color: #666">设备正常运行</span>
        </a-radio>
        <a-radio :value="3">
          <a-tag color="#faad14">维修中</a-tag>
          <span style="margin-left: 8px; color: #666">设备正在维修中</span>
        </a-radio>
        <a-radio :value="4">
          <a-tag color="#ff4d4f">已挂失</a-tag>
          <span style="margin-left: 8px; color: #666">设备已挂失，需找回</span>
        </a-radio>
        <a-radio :value="5">
          <a-tag color="#8c8c8c">已报废</a-tag>
          <span style="margin-left: 8px; color: #666">设备已报废，不可使用</span>
        </a-radio>
      </a-radio-group>

      <a-divider>变更原因（可选）</a-divider>
      <a-textarea v-model="statusRemark" placeholder="请输入状态变更原因" :rows="2" />
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue';
import { Message } from '@arco-design/web-vue';
import axios from 'axios';
import { useRouter } from 'vue-router';

const router = useRouter();
const loading = ref(false);
const searchKeyword = ref('');
const filterStatus = ref<number | null>(null);

const statusModalVisible = ref(false);
const statusChanging = ref(false);
const currentDevice = ref<any>(null);
const newStatus = ref<number | null>(null);
const statusRemark = ref('');

const devices = ref<any[]>([]);

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
});

const statusStats = reactive({
  1: 0, 2: 0, 3: 0, 4: 0, 5: 0, total: 0
});

const columns = [
  { title: '设备ID', slotName: 'deviceId', width: 180 },
  { title: 'MAC地址', dataIndex: 'mac_address', width: 160 },
  { title: '硬件型号', dataIndex: 'hardware_model', width: 120 },
  { title: '固件版本', dataIndex: 'firmware_version', width: 100 },
  { title: '在线状态', slotName: 'isOnline', width: 100 },
  { title: '电量', dataIndex: 'battery_level', width: 80 },
  { title: '状态', slotName: 'lifecycleStatus', width: 100 },
  { title: '最后更新', slotName: 'lastUpdate', width: 180 },
  { title: '操作', slotName: 'actions', width: 150 }
];

const API_BASE = '/api';

const filteredDevices = computed(() => {
  let result = devices.value;
  if (filterStatus.value !== null) {
    result = result.filter(d => d.lifecycle_status === filterStatus.value);
  }
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase();
    result = result.filter(d =>
      d.device_id.toLowerCase().includes(keyword) ||
      d.mac_address.toLowerCase().includes(keyword)
    );
  }
  return result;
});

const loadDevices = async () => {
  loading.value = true;
  try {
    const token = localStorage.getItem('token');
    const res = await axios.get(`${API_BASE}/devices`, {
      params: { page: pagination.current, page_size: pagination.pageSize },
      headers: { 'Authorization': `Bearer ${token}` }
    });
    if (res.data.code === 0) {
      devices.value = res.data.data?.list || [];
      pagination.total = res.data.data?.pagination?.total || 0;
      calcStatusStats();
    }
  } catch (err) {
    devices.value = [
      { device_id: 'DEV001', mac_address: '00:11:22:33:44:55', hardware_model: 'MDM-Pro-200', firmware_version: 'v1.2.0', is_online: true, battery_level: 85, lifecycle_status: 2, last_update_time: '2026-03-19 10:30:00' },
      { device_id: 'DEV002', mac_address: '00:11:22:33:44:56', hardware_model: 'MDM-Mini-100', firmware_version: 'v1.1.5', is_online: true, battery_level: 72, lifecycle_status: 2, last_update_time: '2026-03-19 10:25:00' },
      { device_id: 'DEV003', mac_address: '00:11:22:33:44:57', hardware_model: 'MDM-Lite-50', firmware_version: 'v1.0.0', is_online: false, battery_level: 0, lifecycle_status: 1, last_update_time: '2026-03-18 14:20:00' },
      { device_id: 'DEV004', mac_address: '00:11:22:33:44:58', hardware_model: 'MDM-Pro-200', firmware_version: 'v1.2.0', is_online: true, battery_level: 95, lifecycle_status: 3, last_update_time: '2026-03-19 09:00:00' },
      { device_id: 'DEV005', mac_address: '00:11:22:33:44:59', hardware_model: 'MDM-Mini-100', firmware_version: 'v1.1.5', is_online: false, battery_level: 0, lifecycle_status: 4, last_update_time: '2026-03-17 18:30:00' }
    ];
    pagination.total = devices.value.length;
    calcStatusStats();
    Message.warning('使用模拟数据');
  } finally {
    loading.value = false;
  }
};

const calcStatusStats = () => {
  statusStats[1] = devices.value.filter(d => d.lifecycle_status === 1).length;
  statusStats[2] = devices.value.filter(d => d.lifecycle_status === 2).length;
  statusStats[3] = devices.value.filter(d => d.lifecycle_status === 3).length;
  statusStats[4] = devices.value.filter(d => d.lifecycle_status === 4).length;
  statusStats[5] = devices.value.filter(d => d.lifecycle_status === 5).length;
  statusStats.total = devices.value.length;
};

const handleTableChange = (pag: any) => {
  pagination.current = pag.current;
  loadDevices();
};

const handleSearch = () => {};
const handleFilterChange = () => {};

const filterByStatus = (status: number | null) => {
  filterStatus.value = status;
};

const clearFilters = () => {
  filterStatus.value = null;
  searchKeyword.value = '';
};

const showStatusModal = (record: any) => {
  currentDevice.value = { ...record };
  newStatus.value = record.lifecycle_status;
  statusRemark.value = '';
  statusModalVisible.value = true;
};

const handleStatusChange = async () => {
  if (newStatus.value === currentDevice.value?.lifecycle_status) {
    Message.warning('状态未变更');
    return;
  }
  statusChanging.value = true;
  try {
    const token = localStorage.getItem('token');
    await axios.put(`${API_BASE}/devices/${currentDevice.value.device_id}/status`, {
      status: newStatus.value,
      remark: statusRemark.value
    }, {
      headers: { 'Authorization': `Bearer ${token}` }
    });
    Message.success('设备状态已更新');
  } catch (err) {
    setTimeout(() => {
      const idx = devices.value.findIndex(d => d.device_id === currentDevice.value?.device_id);
      if (idx !== -1) devices.value[idx].lifecycle_status = newStatus.value;
      calcStatusStats();
      Message.success('设备状态已更新（模拟）');
    }, 500);
  } finally {
    statusChanging.value = false;
    statusModalVisible.value = false;
  }
};

const viewDevice = (record: any) => {
  router.push(`/devices/detail/${record.device_id}`);
};

const getStatusColor = (status: number) => {
  const colors: Record<number, string> = { 1: 'blue', 2: 'green', 3: 'orange', 4: 'red', 5: 'gray' };
  return colors[status] || 'default';
};

const getStatusText = (status: number) => {
  const texts: Record<number, string> = { 1: '待激活', 2: '服役中', 3: '维修中', 4: '已挂失', 5: '已报废' };
  return texts[status] || '未知';
};

const formatTime = (time: string) => time || '-';

onMounted(() => {
  loadDevices();
});
</script>

<style scoped lang="less">
.device-status-container {
  padding: 0 20px 20px 20px;
}

.general-card {
  border-radius: 8px;
}

.card-title {
  font-weight: 600;
  font-size: 15px;
}

.stats-row {
  margin-bottom: 0;

  :deep(.arco-card) {
    cursor: pointer;
    transition: all 0.3s;
    text-align: center;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    }
  }
}
</style>
