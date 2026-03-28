<template>
  <div class="child-mode-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="儿童模式设备" :value="8" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="使用中儿童" :value="12" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="使用时长限制" :value="2" suffix="小时/天" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <span>儿童模式配置</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="settings" title="模式设置">
          <a-form layout="vertical" style="max-width: 600px;">
            <a-form-item label="内容过滤">
              <a-checkbox-group v-model="settings.contentFilter">
                <a-checkbox value="violence">暴力内容</a-checkbox>
                <a-checkbox value="adult">成人内容</a-checkbox>
                <a-checkbox value="advertising">广告</a-checkbox>
              </a-checkbox-group>
            </a-form-item>
            <a-form-item label="每日使用时长限制">
              <a-slider v-model="settings.dailyLimit" :marks="{1:'1h',3:'3h',5:'5h',8:'8h'}" :max="8" />
            </a-form-item>
            <a-form-item label="单次使用时长限制">
              <a-select v-model="settings.sessionLimit">
                <a-option value="15">15分钟</a-option>
                <a-option value="30">30分钟</a-option>
                <a-option value="60">1小时</a-option>
                <a-option value="unlimited">不限制</a-option>
              </a-select>
            </a-form-item>
            <a-form-item label="休息间隔">
              <a-input-number v-model="settings.breakInterval" :min="5" /> 分钟
            </a-form-item>
            <a-form-item label="允许语音互动">
              <a-switch v-model="settings.allowVoice" />
            </a-form-item>
            <a-form-item label="允许视频通话">
              <a-switch v-model="settings.allowVideo" />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="handleSave">保存配置</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>
        
        <a-tab-pane key="devices" title="儿童模式设备">
          <a-table :columns="deviceColumns" :data="devices" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="record.childModeEnabled ? 'green' : 'gray'">{{ record.childModeEnabled ? '已启用' : '未启用' }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-link @click="handleToggle(record)">{{ record.childModeEnabled ? '关闭' : '启用' }}</a-link>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="usage" title="使用统计">
          <a-row :gutter="16">
            <a-col :span="12">
              <a-card title="今日使用时长">
                <a-chart :option="usageChart" style="height: 200px;" />
              </a-card>
            </a-col>
            <a-col :span="12">
              <a-card title="使用分布">
                <a-chart :option="distributionChart" style="height: 200px;" />
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 8 });

const settings = reactive({
  contentFilter: ['violence', 'advertising'],
  dailyLimit: 3,
  sessionLimit: '30',
  breakInterval: 15,
  allowVoice: true,
  allowVideo: false,
});

const devices = ref([
  { deviceId: 'DEV001', name: '小黄-客厅', childModeEnabled: true, todayUsage: 120, weeklyUsage: 450 },
  { deviceId: 'DEV002', name: '小红-卧室', childModeEnabled: true, todayUsage: 45, weeklyUsage: 180 },
  { deviceId: 'DEV003', name: '小绿-书房', childModeEnabled: false, todayUsage: 0, weeklyUsage: 60 },
]);

const deviceColumns = [
  { title: '设备ID', dataIndex: 'deviceId', width: 120 },
  { title: '设备名称', dataIndex: 'name', width: 150 },
  { title: '今日使用', dataIndex: 'todayUsage', width: 100 },
  { title: '本周使用(分钟)', dataIndex: 'weeklyUsage', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 100 },
];

const usageChart = {
  xAxis: { type: 'category', data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'] },
  yAxis: { type: 'value' },
  series: [{ type: 'bar', data: [60, 90, 75, 120, 180, 240, 120] }],
};

const distributionChart = {
  tooltip: { trigger: 'item' },
  series: [{ name: '使用分布', type: 'pie', radius: '60%', data: [
    { value: 35, name: '学习' },
    { value: 25, name: '游戏' },
    { value: 20, name: '视频' },
    { value: 20, name: '互动' },
  ]}],
};

const handleSave = () => {};
const handleToggle = (record: any) => { record.childModeEnabled = !record.childModeEnabled; };
</script>

<style scoped>
.child-mode-container { padding: 20px; }
</style>
