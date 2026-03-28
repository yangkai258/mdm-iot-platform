<template>
  <div class="pet-training-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="训练任务" :value="128" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="进行中" :value="12" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="完成率" :value="85" suffix="%" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="技能数" :value="25" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>宠物训练中心</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建任务
          </a-button>
        </div>
      </template>
      
      <a-tabs>
        <a-tab-pane key="tasks" title="训练任务">
          <a-table :columns="taskColumns" :data="tasks" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
            </template>
            <template #progress="{ record }">
              <a-progress :percent="record.progress" :status="record.progress === 100 ? 'success' : 'normal'" />
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="skills" title="技能库">
          <a-row :gutter="16">
            <a-col :span="6" v-for="skill in skills" :key="skill.id">
              <a-card size="small" class="skill-card">
                <div class="skill-icon">{{ skill.icon }}</div>
                <div class="skill-name">{{ skill.name }}</div>
                <div class="skill-difficulty">
                  <a-rate :model-value="skill.difficulty" disabled count="3" />
                </div>
                <div class="skill-desc">{{ skill.description }}</div>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
        
        <a-tab-pane key="history" title="训练记录">
          <a-table :columns="historyColumns" :data="history" :pagination="pagination" />
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <a-modal v-model:visible="createVisible" title="创建训练任务" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="宠物" required>
          <a-select v-model="form.petId" placeholder="选择宠物">
            <a-option value="P001">小黄</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="训练技能">
          <a-select v-model="form.skillId" placeholder="选择技能">
            <a-option value="S001">坐下</a-option>
            <a-option value="S002">握手</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="训练目标">
          <a-input-number v-model="form.targetCount" :min="1" /> 次
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="form.note" :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 6 });
const createVisible = ref(false);

const form = reactive({ petId: '', skillId: '', targetCount: 10, note: '' });

const tasks = ref([
  { id: 1, petName: '小黄', skillName: '坐下', status: 'training', statusText: '进行中', progress: 60, totalCount: 10, currentCount: 6 },
  { id: 2, petName: '小黄', skillName: '握手', status: 'completed', statusText: '已完成', progress: 100, totalCount: 20, currentCount: 20 },
]);

const skills = ref([
  { id: 1, name: '坐下', icon: '🪑', difficulty: 1, description: '基础服从命令' },
  { id: 2, name: '握手', icon: '🤝', difficulty: 1, description: '建立亲密关系' },
  { id: 3, name: '翻滚', icon: '🔄', difficulty: 2, description: '趣味技能' },
  { id: 4, name: '装死', icon: '💀', difficulty: 3, description: '高级表演技能' },
]);

const history = ref([
  { id: 1, petName: '小黄', skillName: '坐下', score: 95, duration: '15分钟', time: '2026-03-28 10:00:00' },
]);

const taskColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '宠物', dataIndex: 'petName', width: 100 },
  { title: '技能', dataIndex: 'skillName', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '进度', slotName: 'progress', width: 150 },
  { title: '完成/目标', dataIndex: 'currentCount', width: 120 },
];

const historyColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '宠物', dataIndex: 'petName', width: 100 },
  { title: '技能', dataIndex: 'skillName', width: 120 },
  { title: '评分', dataIndex: 'score', width: 80 },
  { title: '耗时', dataIndex: 'duration', width: 100 },
  { title: '时间', dataIndex: 'time', width: 160 },
];

const getStatusColor = (s: string) => ({ pending: 'gray', training: 'blue', completed: 'green' }[s] || 'default');

const handleCreate = () => { createVisible.value = true; };
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.pet-training-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
.skill-card { text-align: center; margin-bottom: 12px; }
.skill-icon { font-size: 36px; margin-bottom: 8px; }
.skill-name { font-weight: bold; margin-bottom: 4px; }
.skill-difficulty { margin-bottom: 8px; }
.skill-desc { color: #86909c; font-size: 12px; }
</style>
