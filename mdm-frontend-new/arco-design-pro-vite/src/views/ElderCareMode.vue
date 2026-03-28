<template>
  <div class="elder-care-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="老人陪伴模式设备" :value="5" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="今日互动" :value="28" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="健康提醒" :value="3" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <span>老人陪伴模式</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="settings" title="模式设置">
          <a-form layout="vertical" style="max-width: 600px;">
            <a-form-item label="语速">
              <a-slider v-model="settings.speechSpeed" :marks="{0:'慢',50:'适中',100:'快'}" />
            </a-form-item>
            <a-form-item label="音量">
              <a-slider v-model="settings.volume" :marks="{0:'低',50:'中',100:'高'}" />
            </a-form-item>
            <a-form-item label="界面字体大小">
              <a-slider v-model="settings.fontSize" :marks="{0:'小',50:'中',100:'大'}" />
            </a-form-item>
            <a-form-item label="定时提醒">
              <a-space direction="vertical" fill>
                <a-checkbox v-model="settings.reminders.medication">用药提醒</a-checkbox>
                <a-checkbox v-model="settings.reminders.water">饮水提醒</a-checkbox>
                <a-checkbox v-model="settings.reminders.exercise">运动提醒</a-checkbox>
              </a-space>
            </a-form-item>
            <a-form-item label="情感陪伴">
              <a-switch v-model="settings.emotionSupport" />
              <span style="margin-left: 8px; color: #86909c;">启用主动情感陪伴功能</span>
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="handleSave">保存配置</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>
        
        <a-tab-pane key="reminders" title="提醒记录">
          <a-table :columns="reminderColumns" :data="reminders" :pagination="pagination" />
        </a-tab-pane>
        
        <a-tab-pane key="devices" title="绑定设备">
          <a-table :columns="deviceColumns" :data="devices" :pagination="pagination" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 5 });

const settings = reactive({
  speechSpeed: 30,
  volume: 70,
  fontSize: 60,
  reminders: { medication: true, water: true, exercise: false },
  emotionSupport: true,
});

const reminders = ref([
  { id: 1, type: 'medication', content: '该吃药了', petName: '小黄', time: '08:00', status: 'completed', completedAt: '08:05' },
  { id: 2, type: 'water', content: '该喝水了', petName: '小黄', time: '10:00', status: 'pending', completedAt: null },
]);

const devices = ref([
  { deviceId: 'DEV001', name: '小黄-客厅', elderName: '爷爷', linkedAt: '2026-03-01' },
  { deviceId: 'DEV002', name: '小红-卧室', elderName: '奶奶', linkedAt: '2026-03-01' },
]);

const reminderColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '类型', dataIndex: 'type', width: 120 },
  { title: '内容', dataIndex: 'content', width: 200 },
  { title: '关联宠物', dataIndex: 'petName', width: 100 },
  { title: '提醒时间', dataIndex: 'time', width: 100 },
  { title: '状态', dataIndex: 'status', width: 100 },
  { title: '完成时间', dataIndex: 'completedAt', width: 120 },
];

const deviceColumns = [
  { title: '设备ID', dataIndex: 'deviceId', width: 120 },
  { title: '设备名称', dataIndex: 'name', width: 150 },
  { title: '关联老人', dataIndex: 'elderName', width: 100 },
  { title: '绑定时间', dataIndex: 'linkedAt', width: 120 },
];

const handleSave = () => {};
</script>

<style scoped>
.elder-care-container { padding: 20px; }
</style>
