<template>
  <div class="geofence-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="12">
        <a-card title="地理围栏地图">
          <div class="map-placeholder">
            <icon-global :size="64" />
            <div style="margin-top: 16px;">地图组件占位</div>
            <div style="color: #86909c; font-size: 12px; margin-top: 8px;">显示围栏区域和设备位置</div>
          </div>
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="围栏统计">
          <a-statistic title="活跃围栏数" :value="stats.activeFences" style="margin-bottom: 16px;" />
          <a-descriptions :column="1" size="small" bordered>
            <a-descriptions-item label="设备在围栏内">12台</a-descriptions-item>
            <a-descriptions-item label="设备离开围栏">2台</a-descriptions-item>
            <a-descriptions-item label="告警触发">3次</a-descriptions-item>
          </a-descriptions>
        </a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>地理围栏配置</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新增围栏
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag :color="record.type === 'circle' ? 'blue' : 'green'">{{ record.type === 'circle' ? '圆形' : '多边形' }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.enabled ? 'green' : 'gray'">{{ record.enabled ? '启用' : '禁用' }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link @click="handleViewDevices(record)">查看设备</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 创建/编辑弹窗 -->
    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑围栏' : '新增围栏'" :width="600" @before-ok="handleSubmit">
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
        <a-form-item v-if="form.type === 'circle'" label="中心点坐标">
          <a-space>
            <a-input-number v-model="form.lat" placeholder="纬度" style="width: 120px;" />
            <span>经度</span>
            <a-input-number v-model="form.lng" placeholder="经度" style="width: 120px;" />
          </a-space>
        </a-form-item>
        <a-form-item v-if="form.type === 'circle'" label="半径">
          <a-space>
            <a-input-number v-model="form.radius" :min="10" style="width: 120px;" />
            <span>米</span>
          </a-space>
        </a-form-item>
        <a-form-item label="触发动作">
          <a-checkbox-group v-model="form.actions">
            <a-checkbox value="enter">进入告警</a-checkbox>
            <a-checkbox value="exit">离开告警</a-checkbox>
            <a-checkbox value="dwell">停留超时告警</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
        <a-form-item v-if="form.actions.includes('dwell')" label="停留超时时间">
          <a-input-number v-model="form.dwellTime" :min="1" /> 分钟
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

const loading = ref(false);

const stats = reactive({ activeFences: 5, devicesInside: 12, devicesOutside: 2, alerts: 3 });

const data = ref([
  { id: 1, name: 'home-fence', type: 'circle', lat: 39.9042, lng: 116.4074, radius: 100, actions: ['enter', 'exit'], dwellTime: 0, deviceCount: 5, alertCount: 2, enabled: true },
  { id: 2, name: 'park-fence', type: 'circle', lat: 39.9142, lng: 116.4174, radius: 200, actions: ['exit'], dwellTime: 0, deviceCount: 3, alertCount: 0, enabled: true },
  { id: 3, name: 'office-fence', type: 'polygon', coordinates: [], actions: ['enter', 'dwell'], dwellTime: 30, deviceCount: 4, alertCount: 1, enabled: true },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 3 });
const editVisible = ref(false);
const isEdit = ref(false);

const form = reactive({ name: '', type: 'circle', lat: 39.9042, lng: 116.4074, radius: 100, actions: [], dwellTime: 0, enabled: true });

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '围栏名称', dataIndex: 'name', width: 150 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '设备数', dataIndex: 'deviceCount', width: 80 },
  { title: '告警数', dataIndex: 'alertCount', width: 80 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 220, fixed: 'right' },
];

const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleViewDevices = (record: any) => {};
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.geofence-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
.map-placeholder { text-align: center; padding: 48px; color: #86909c; background: #f7f8fa; border-radius: 8px; }
</style>
