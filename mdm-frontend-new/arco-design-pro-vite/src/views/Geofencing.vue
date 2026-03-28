<template>
  <div class="geofencing-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="电子围栏" :value="12" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="告警次数" :value="36" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="保护设备" :value="8" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>地理围栏管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建围栏
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="fences" :pagination="pagination">
        <template #type="{ record }">
          <a-tag>{{ record.typeText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.enabled ? 'green' : 'gray'">{{ record.enabled ? '启用' : '禁用' }}</a-tag>
        </template>
        <template #alertCount="{ record }">
          <a-badge :value="record.alertCount" :max="99" :color="record.alertCount > 0 ? 'red' : 'green'" />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link @click="handleViewMap(record)">地图</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑围栏' : '创建围栏'" :width="600" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="围栏名称" required>
          <a-input v-model="form.name" placeholder="请输入围栏名称" />
        </a-form-item>
        <a-form-item label="围栏类型">
          <a-radio-group v-model="form.type">
            <a-radio value="circle">圆形</a-radio>
            <a-radio value="polygon">多边形</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="中心点">
          <a-space>
            <a-input-number v-model="form.lat" placeholder="纬度" style="width: 120px;" />
            <a-input-number v-model="form.lng" placeholder="经度" style="width: 120px;" />
          </a-space>
        </a-form-item>
        <a-form-item label="半径" v-if="form.type === 'circle'">
          <a-input-number v-model="form.radius" :min="10" style="width: 200px;" />
          <span style="margin-left: 8px;">米</span>
        </a-form-item>
        <a-form-item label="关联设备">
          <a-select v-model="form.deviceIds" multiple placeholder="选择设备">
            <a-option value="DEV001">小黄-客厅</a-option>
            <a-option value="DEV002">小红-卧室</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="触发动作">
          <a-checkbox-group v-model="form.actions">
            <a-checkbox value="alert">发送告警</a-checkbox>
            <a-checkbox value="notify">推送通知</a-checkbox>
            <a-checkbox value="record">记录日志</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
        <a-form-item label="启用">
          <a-switch v-model="form.enabled" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 4 });
const editVisible = ref(false);
const isEdit = ref(false);

const form = reactive({ name: '', type: 'circle', lat: 0, lng: 0, radius: 100, deviceIds: [], actions: ['alert'], enabled: true });

const fences = ref([
  { id: 1, name: '客厅区域', type: 'circle', typeText: '圆形', lat: 30.5728, lng: 114.2669, radius: 50, deviceNames: ['小黄'], alertCount: 5, enabled: true },
  { id: 2, name: '阳台区域', type: 'circle', typeText: '圆形', lat: 30.5730, lng: 114.2670, radius: 30, deviceNames: ['小红'], alertCount: 0, enabled: true },
  { id: 3, name: '家庭区域', type: 'polygon', typeText: '多边形', points: '多点', deviceNames: ['小黄', '小红'], alertCount: 12, enabled: true },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '围栏名称', dataIndex: 'name', width: 120 },
  { title: '类型', slotName: 'type', width: 80 },
  { title: '位置', dataIndex: 'lat', width: 150, ellipsis: true },
  { title: '关联设备', dataIndex: 'deviceNames', width: 150 },
  { title: '告警次数', slotName: 'alertCount', width: 100 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' },
];

const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleViewMap = (record: any) => {};
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.geofencing-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
