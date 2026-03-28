<template>
  <div class="pet-profile-container">
    <!-- 宠物基本信息卡片 -->
    <a-row :gutter="16">
      <a-col :span="6">
        <a-card>
          <div class="pet-avatar">
            <a-avatar :size="100" :style="{ backgroundColor: '#165DFF' }">
              <icon-robot :size="60" />
            </a-avatar>
            <div class="pet-name">{{ pet.name }}</div>
            <div class="pet-type">{{ pet.type }} - {{ pet.breed }}</div>
          </div>
          <a-descriptions :column="1" size="small">
            <a-descriptions-item label="年龄">{{ pet.age }}岁</a-descriptions-item>
            <a-descriptions-item label="性别">{{ pet.gender === 'male' ? '公' : '母' }}</a-descriptions-item>
            <a-descriptions-item label="体重">{{ pet.weight }}kg</a-descriptions-item>
            <a-descriptions-item label="绑定设备">{{ pet.deviceId }}</a-descriptions-item>
          </a-descriptions>
        </a-card>
      </a-col>

      <!-- 健康雷达图 -->
      <a-col :span="9">
        <a-card title="健康雷达">
          <a-chart :option="healthRadar" style="height: 280px;" />
        </a-card>
      </a-col>

      <!-- 行为分析 -->
      <a-col :span="9">
        <a-card title="行为分析">
          <a-descriptions :column="2" size="small">
            <a-descriptions-item label="今日活跃度">
              <a-progress :percent="pet.activities.today" :color="getActivityColor(pet.activities.today)" />
            </a-descriptions-item>
            <a-descriptions-item label="睡眠质量">
              <a-progress :percent="pet.activities.sleep" :color="getActivityColor(pet.activities.sleep)" />
            </a-descriptions-item>
            <a-descriptions-item label="食欲状态">
              <a-progress :percent="pet.activities.appetite" :color="getActivityColor(pet.activities.appetite)" />
            </a-descriptions-item>
            <a-descriptions-item label="情绪指数">
              <a-progress :percent="pet.activities.mood" :color="getActivityColor(pet.activities.mood)" />
            </a-descriptions-item>
          </a-descriptions>
          <a-divider />
          <div class="mood-trend">
            <div class="trend-label">情绪趋势</div>
            <a-chart :option="moodTrend" style="height: 100px;" />
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 详细信息Tab -->
    <a-card style="margin-top: 16px;">
      <a-tabs>
        <!-- 基础档案 -->
        <a-tab-pane key="basic" title="基础档案">
          <a-descriptions :column="3" bordered>
            <a-descriptions-item label="宠物名称">{{ pet.name }}</a-descriptions-item>
            <a-descriptions-item label="宠物类型">{{ pet.type }}</a-descriptions-item>
            <a-descriptions-item label="品种">{{ pet.breed }}</a-descriptions-item>
            <a-descriptions-item label="出生日期">{{ pet.birthDate }}</a-descriptions-item>
            <a-descriptions-item label="年龄">{{ pet.age }}岁</a-descriptions-item>
            <a-descriptions-item label="性别">{{ pet.gender === 'male' ? '公' : '母' }}</a-descriptions-item>
            <a-descriptions-item label="体重">{{ pet.weight }}kg</a-descriptions-item>
            <a-descriptions-item label="身高">{{ pet.height }}cm</a-descriptions-item>
            <a-descriptions-item label="毛色">{{ pet.color }}</a-descriptions-item>
          </a-descriptions>
          <a-divider>性格特征</a-divider>
          <a-space wrap>
            <a-tag v-for="trait in pet.personality" :key="trait" color="blue">{{ trait }}</a-tag>
          </a-space>
        </a-tab-pane>

        <!-- 健康记录 -->
        <a-tab-pane key="health" title="健康记录">
          <a-table :columns="healthColumns" :data="healthRecords" :pagination="pagination">
            <template #type="{ record }">
              <a-tag :color="getHealthTypeColor(record.type)">{{ record.typeText }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>

        <!-- 互动历史 -->
        <a-tab-pane key="interaction" title="互动历史">
          <a-table :columns="interactionColumns" :data="interactions" :pagination="pagination">
            <template #type="{ record }">
              <a-tag :color="getInteractionColor(record.type)">{{ record.typeText }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>

        <!-- 训练进度 -->
        <a-tab-pane key="training" title="训练进度">
          <a-row :gutter="16">
            <a-col :span="8" v-for="skill in trainingSkills" :key="skill.name">
              <a-card size="small">
                <div class="skill-item">
                  <div class="skill-name">{{ skill.name }}</div>
                  <a-progress :percent="skill.mastery" :color="getMasteryColor(skill.mastery)" />
                  <div class="skill-meta">
                    <span>学习次数: {{ skill.times }}</span>
                    <span>最近: {{ skill.lastPracticed }}</span>
                  </div>
                </div>
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

const pagination = reactive({ current: 1, pageSize: 10, total: 5 });

const pet = reactive({
  id: 'P001',
  name: '小黄',
  type: '狗狗',
  breed: '金毛',
  age: 3,
  gender: 'male',
  weight: 28.5,
  height: 65,
  color: '金色',
  birthDate: '2023-01-15',
  deviceId: 'DEV001',
  personality: ['活泼', '友好', '好奇', '贪吃'],
  activities: {
    today: 78,
    sleep: 85,
    appetite: 92,
    mood: 88,
  },
});

const healthRecords = ref([
  { date: '2026-03-28', type: 'checkup', typeText: '体检', doctor: '张医生', result: '正常', note: '各项指标良好' },
  { date: '2026-03-15', type: 'vaccine', typeText: '疫苗', doctor: '李医生', result: '已完成', note: '狂犬疫苗接种' },
  { date: '2026-03-01', type: 'deworming', typeText: '驱虫', doctor: '-', result: '已完成', note: '内外驱虫' },
]);

const interactions = ref([
  { time: '2026-03-28 10:00', type: 'play', typeText: '玩耍', duration: '30分钟', content: '玩飞盘' },
  { time: '2026-03-28 08:00', type: 'walk', typeText: '散步', duration: '45分钟', content: '公园散步' },
  { time: '2026-03-27 20:00', type: 'training', typeText: '训练', duration: '20分钟', content: '坐下、握手' },
  { time: '2026-03-27 18:00', type: 'feed', typeText: '喂食', duration: '-', content: '晚餐' },
]);

const trainingSkills = ref([
  { name: '坐下', mastery: 95, times: 50, lastPracticed: '2026-03-27' },
  { name: '握手', mastery: 88, times: 45, lastPracticed: '2026-03-27' },
  { name: '卧倒', mastery: 72, times: 30, lastPracticed: '2026-03-25' },
  { name: '翻滚', mastery: 45, times: 15, lastPracticed: '2026-03-20' },
  { name: '装死', mastery: 30, times: 8, lastPracticed: '2026-03-15' },
  { name: '取物', mastery: 20, times: 5, lastPracticed: '2026-03-10' },
]);

const healthColumns = [
  { title: '日期', dataIndex: 'date', width: 120 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '医生', dataIndex: 'doctor', width: 80 },
  { title: '结果', dataIndex: 'result', width: 100 },
  { title: '备注', dataIndex: 'note' },
];

const interactionColumns = [
  { title: '时间', dataIndex: 'time', width: 160 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '时长', dataIndex: 'duration', width: 100 },
  { title: '内容', dataIndex: 'content' },
];

const healthRadar = {
  radar: {
    indicator: [
      { name: '体力', max: 100 },
      { name: '食欲', max: 100 },
      { name: '睡眠', max: 100 },
      { name: '情绪', max: 100 },
      { name: '免疫', max: 100 },
    ],
  },
  series: [{
    type: 'radar',
    data: [{ value: [85, 92, 85, 88, 90], name: '健康指数' }],
  }],
};

const moodTrend = {
  xAxis: { type: 'category', data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'] },
  yAxis: { type: 'value', min: 0, max: 100 },
  series: [{ type: 'line', data: [82, 85, 88, 90, 87, 92, 88], smooth: true, areaStyle: {} }],
};

const getActivityColor = (val: number) => val >= 80 ? '#00B42A' : val >= 60 ? '#FF7D00' : '#F53F3F';
const getHealthTypeColor = (t: string) => ({ checkup: 'blue', vaccine: 'green', deworming: 'orange' }[t] || 'default');
const getInteractionColor = (t: string) => ({ play: 'blue', walk: 'green', training: 'orange', feed: 'purple' }[t] || 'default');
const getMasteryColor = (val: number) => val >= 80 ? '#00B42A' : val >= 50 ? '#FF7D00' : '#F53F3F';
</script>

<style scoped>
.pet-profile-container { padding: 20px; }
.pet-avatar { text-align: center; margin-bottom: 16px; }
.pet-name { font-size: 24px; font-weight: bold; margin-top: 12px; }
.pet-type { color: #86909c; margin-top: 4px; }
.mood-trend { margin-top: 16px; }
.trend-label { font-size: 14px; color: #4a5568; margin-bottom: 8px; }
.skill-item { text-align: center; }
.skill-name { font-weight: 500; margin-bottom: 8px; }
.skill-meta { display: flex; justify-content: space-between; font-size: 12px; color: #86909c; margin-top: 8px; }
</style>
