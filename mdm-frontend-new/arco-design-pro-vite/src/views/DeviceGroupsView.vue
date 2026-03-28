<template>
  <div class="device-groups-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>设备分组</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新增分组
          </a-button>
        </div>
      </template>
      
      <div class="search-area">
        <a-row :gutter="16">
          <a-col :span="6">
            <a-input v-model="searchForm.keyword" placeholder="搜索分组名称" allow-clear />
          </a-col>
          <a-col :span="2">
            <a-button type="primary" @click="handleSearch">筛选</a-button>
          </a-col>
        </a-row>
      </div>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag :color="record.type === 'static' ? 'blue' : 'green'">
            {{ record.type === 'static' ? '静态分组' : '动态分组' }}
          </a-tag>
        </template>
        <template #deviceCount="{ record }">
          <a-badge :count="record.deviceCount" :max-count="999" />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleViewDevices(record)">查看设备</a-link>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link @click="handleDelete(record)" status="danger">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 设备列表弹窗 -->
    <a-modal v-model:visible="devicesVisible" title="分组设备" :width="800">
      <a-table :columns="deviceColumns" :data="groupDevices" :pagination="pagination">
        <template #status="{ record }">
          <a-tag :color="record.online ? 'green' : 'red'">
            {{ record.online ? '在线' : '离线' }}
          </a-tag>
        </template>
      </a-table>
    </a-modal>

    <!-- 创建/编辑弹窗 -->
    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑分组' : '新增分组'" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="分组名称" required>
          <a-input v-model="form.name" placeholder="请输入分组名称" />
        </a-form-item>
        <a-form-item label="分组类型">
          <a-radio-group v-model="form.type">
            <a-radio value="static">静态分组</a-radio>
            <a-radio value="dynamic">动态分组</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item v-if="form.type === 'dynamic'" label="动态规则">
          <a-textarea v-model="form.rule" placeholder="如: firmware = 'v2.0.0' AND status = 'online'" :rows="3" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" placeholder="请输入描述" :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const data = ref([
  { id: 'G001', name: '在线设备', type: 'dynamic', rule: "status = 'online'", deviceCount: 45, description: '所有在线设备', createdAt: '2026-03-01' },
  { id: 'G002', name: 'M5Stack设备', type: 'dynamic', rule: "type = 'M5Stack'", deviceCount: 38, description: 'M5Stack系列设备', createdAt: '2026-03-01' },
  { id: 'G003', name: '低电量告警', type: 'dynamic', rule: "battery < 20", deviceCount: 5, description: '电量低于20%的设备', createdAt: '2026-03-05' },
  { id: 'G004', name: '客厅设备', type: 'static', rule: '', deviceCount: 12, description: '客厅中的设备', createdAt: '2026-03-10' },
  { id: 'G005', name: '卧室设备', type: 'static', rule: '', deviceCount: 8, description: '卧室中的设备', createdAt: '2026-03-10' },
]);

const groupDevices = ref([
  { deviceId: 'DEV001', name: '小黄-客厅', type: 'M5Stack', online: true, battery: 85, firmware: 'v2.0.0' },
  { deviceId: 'DEV002', name: '小红-客厅', type: 'M5Stack', online: true, battery: 62, firmware: 'v2.0.0' },
]);

const searchForm = reactive({ keyword: '' });
const pagination = reactive({ current: 1, pageSize: 20, total: 5 });

const columns = [
  { title: '分组ID', dataIndex: 'id', width: 100 },
  { title: '分组名称', dataIndex: 'name', width: 150 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '设备数', slotName: 'deviceCount', width: 100 },
  { title: '描述', dataIndex: 'description' },
  { title: '创建时间', dataIndex: 'createdAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' },
];

const deviceColumns = [
  { title: '设备ID', dataIndex: 'deviceId', width: 120 },
  { title: '设备名称', dataIndex: 'name', width: 150 },
  { title: '类型', dataIndex: 'type', width: 100 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '电量', dataIndex: 'battery', width: 80 },
  { title: '固件', dataIndex: 'firmware', width: 100 },
];

const devicesVisible = ref(false);
const editVisible = ref(false);
const isEdit = ref(false);

const form = reactive({
  name: '',
  type: 'static',
  rule: '',
  description: '',
});

const handleSearch = () => {};
const handleCreate = () => {
  isEdit.value = false;
  Object.assign(form, { name: '', type: 'static', rule: '', description: '' });
  editVisible.value = true;
};
const handleEdit = (record: any) => {
  isEdit.value = true;
  Object.assign(form, record);
  editVisible.value = true;
};
const handleDelete = (record: any) => {};
const handleViewDevices = (record: any) => { devicesVisible.value = true; };
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.device-groups-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
.search-area { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
</style>
