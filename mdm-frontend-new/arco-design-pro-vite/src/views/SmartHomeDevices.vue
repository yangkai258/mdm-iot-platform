<template>
  <div class="smart-home-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="已连接设备" :value="12" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="支持品牌" :value="8" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="今日联动" :value="45" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>智能家居对接</span>
          <a-button type="primary" @click="handleAdd">
            <template #icon><icon-plus /></template>
            添加设备
          </a-button>
        </div>
      </template>
      
      <a-tabs>
        <a-tab-pane key="devices" title="设备列表">
          <a-row :gutter="16">
            <a-col :span="6" v-for="device in devices" :key="device.id">
              <a-card size="small" class="device-card">
                <div class="device-icon">{{ device.icon }}</div>
                <div class="device-name">{{ device.name }}</div>
                <div class="device-brand">{{ device.brand }}</div>
                <a-switch :checked="device.online" disabled />
                <a-button type="text" size="small">配置</a-button>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
        
        <a-tab-pane key="brands" title="支持品牌">
          <a-table :columns="brandColumns" :data="brands" :pagination="pagination" />
        </a-tab-pane>
        
        <a-tab-pane key="scenes" title="联动场景">
          <a-table :columns="sceneColumns" :data="scenes" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="record.enabled ? 'green' : 'gray'">{{ record.enabled ? '启用' : '禁用' }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 8 });

const devices = ref([
  { id: 1, name: '智能灯泡', brand: '小米', icon: '💡', online: true },
  { id: 2, name: '智能音箱', brand: '小米', icon: '🔊', online: true },
  { id: 3, name: '空调', brand: '格力', icon: '❄️', online: true },
  { id: 4, name: '扫地机器人', brand: '科沃斯', icon: '🤖', online: false },
]);

const brands = ref([
  { id: 1, name: '小米', logo: '', deviceCount: 5, supported: true },
  { id: 2, name: '格力', logo: '', deviceCount: 2, supported: true },
  { id: 3, name: '科沃斯', logo: '', deviceCount: 1, supported: true },
  { id: 4, name: '华为', logo: '', deviceCount: 0, supported: false },
]);

const scenes = ref([
  { id: 1, name: '晚安模式', trigger: '时间: 22:00', actions: '关闭灯光、调低音量', enabled: true },
  { id: 2, name: '离家模式', trigger: '设备离开', actions: '开启监控、关闭电器', enabled: true },
]);

const brandColumns = [
  { title: '品牌', dataIndex: 'name', width: 150 },
  { title: '设备数', dataIndex: 'deviceCount', width: 100 },
  { title: '支持状态', dataIndex: 'supported', width: 100 },
];

const sceneColumns = [
  { title: '场景名称', dataIndex: 'name', width: 150 },
  { title: '触发条件', dataIndex: 'trigger', width: 200 },
  { title: '执行动作', dataIndex: 'actions' },
  { title: '状态', slotName: 'status', width: 100 },
];

const handleAdd = () => {};
</script>

<style scoped>
.smart-home-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
.device-card { text-align: center; margin-bottom: 12px; }
.device-icon { font-size: 48px; margin-bottom: 8px; }
.device-name { font-weight: bold; }
.device-brand { color: #86909c; font-size: 12px; }
</style>
