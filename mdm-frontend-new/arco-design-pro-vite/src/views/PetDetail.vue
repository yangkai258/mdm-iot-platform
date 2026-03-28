<template>
  <div class="pet-config-container">
    <a-card>
      <template #title>
        <span>宠物配置</span>
      </template>
      
      <a-tabs>
        <!-- 基本信息 -->
        <a-tab-pane key="basic" title="基本信息">
          <a-form :model="basicForm" layout="vertical">
            <a-form-item label="宠物名称">
              <a-input v-model="basicForm.name" placeholder="请输入宠物名称" />
            </a-form-item>
            <a-form-item label="宠物类型">
              <a-select v-model="basicForm.type" placeholder="选择宠物类型">
                <a-option value="dog">狗狗</a-option>
                <a-option value="cat">猫咪</a-option>
                <a-option value="bird">小鸟</a-option>
                <a-option value="other">其他</a-option>
              </a-select>
            </a-form-item>
            <a-form-item label="品种">
              <a-input v-model="basicForm.breed" placeholder="请输入品种" />
            </a-form-item>
            <a-form-item label="年龄">
              <a-input-number v-model="basicForm.age" :min="0" :max="30" />
            </a-form-item>
            <a-form-item label="性别">
              <a-radio-group v-model="basicForm.gender">
                <a-radio value="male">公</a-radio>
                <a-radio value="female">母</a-radio>
              </a-radio-group>
            </a-form-item>
            <a-form-item>
              <a-space>
                <a-button type="primary" @click="handleSaveBasic">保存</a-button>
                <a-button @click="handleResetBasic">重置</a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <!-- 性格配置 -->
        <a-tab-pane key="personality" title="性格配置">
          <a-form :model="personalityForm" layout="vertical">
            <a-form-item label="性格类型">
              <a-select v-model="personalityForm.type" placeholder="选择性格类型">
                <a-option value="active">活泼好动</a-option>
                <a-option value="calm">温顺安静</a-option>
                <a-option value="curious">好奇探索</a-option>
                <a-option value="playful">顽皮爱玩</a-option>
                <a-option value="gentle">温柔体贴</a-option>
              </a-select>
            </a-form-item>
            <a-form-item label="活跃程度">
              <a-slider v-model="personalityForm.activityLevel" :marks="activityMarks" :max="100" />
            </a-form-item>
            <a-form-item label="亲人程度">
              <a-slider v-model="personalityForm.affectionLevel" :marks="affectionMarks" :max="100" />
            </a-form-item>
            <a-form-item label="训练难度">
              <a-slider v-model="personalityForm.trainingDifficulty" :marks="difficultyMarks" :max="100" />
            </a-form-item>
            <a-form-item label="性格描述">
              <a-textarea v-model="personalityForm.description" placeholder="描述宠物性格特点" :rows="3" />
            </a-form-item>
            <a-form-item>
              <a-space>
                <a-button type="primary" @click="handleSavePersonality">保存</a-button>
                <a-button @click="handleResetPersonality">重置</a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <!-- 交互设置 -->
        <a-tab-pane key="interaction" title="交互设置">
          <a-form :model="interactionForm" layout="vertical">
            <a-form-item label="每日最大交互次数">
              <a-input-number v-model="interactionForm.maxDailyInteractions" :min="1" :max="100" />
            </a-form-item>
            <a-form-item label="交互频率">
              <a-select v-model="interactionForm.frequency" placeholder="选择交互频率">
                <a-option value="high">高频 (随时响应)</a-option>
                <a-option value="normal">普通 (适度响应)</a-option>
                <a-option value="low">低频 (减少打扰)</a-option>
              </a-select>
            </a-form-item>
            <a-form-item label="免打扰时段">
              <a-space direction="vertical" fill>
                <a-switch v-model="interactionForm.dndEnabled" />
                <a-space v-if="interactionForm.dndEnabled">
                  <a-time-picker v-model="interactionForm.dndStart" format="HH:mm" placeholder="开始时间" />
                  <span>至</span>
                  <a-time-picker v-model="interactionForm.dndEnd" format="HH:mm" placeholder="结束时间" />
                </a-space>
              </a-space>
            </a-form-item>
            <a-form-item label="响应灵敏度">
              <a-slider v-model="interactionForm.sensitivity" :marks="sensitivityMarks" :max="100" />
            </a-form-item>
            <a-form-item label="声音大小">
              <a-slider v-model="interactionForm.volume" :max="100" />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="handleSaveInteraction">保存</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <!-- 宠物记忆 -->
        <a-tab-pane key="memory" title="宠物记忆">
          <a-tabs>
            <a-tab-pane key="short" title="短期记忆">
              <a-table :columns="shortTermColumns" :data="shortTermMemory" :pagination="pagination">
                <template #type="{ record }">
                  <a-tag :color="getMemoryTypeColor(record.type)">{{ record.typeText }}</a-tag>
                </template>
              </a-table>
            </a-tab-pane>
            <a-tab-pane key="long" title="长期记忆">
              <a-table :columns="longTermColumns" :data="longTermMemory" :pagination="pagination">
                <template #type="{ record }">
                  <a-tag :color="getMemoryTypeColor(record.type)">{{ record.typeText }}</a-tag>
                </template>
              </a-table>
            </a-tab-pane>
            <a-tab-pane key="learning" title="学习进度">
              <a-list :data-source="learningProgress">
                <template #item="{ item }">
                  <a-list-item-meta :title="item.action" :description="`掌握度: ${item.mastery}%`">
                    <template #avatar>
                      <a-progress type="circle" :percent="item.mastery" :size="40" />
                    </template>
                  </a-list-item-meta>
                </template>
              </a-list>
            </a-tab-pane>
          </a-tabs>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 5 });

const basicForm = reactive({
  name: '小黄',
  type: 'dog',
  breed: '金毛',
  age: 3,
  gender: 'male',
});

const personalityForm = reactive({
  type: 'active',
  activityLevel: 75,
  affectionLevel: 90,
  trainingDifficulty: 30,
  description: '活泼好动，喜欢和人互动，学习能力强',
});

const interactionForm = reactive({
  maxDailyInteractions: 20,
  frequency: 'normal',
  dndEnabled: true,
  dndStart: '22:00',
  dndEnd: '08:00',
  sensitivity: 80,
  volume: 70,
});

const shortTermMemory = ref([
  { id: 'M001', type: 'play', typeText: '玩耍', content: '今天和主人玩了飞盘', time: '2026-03-28 10:00:00' },
  { id: 'M002', type: 'food', typeText: '进食', content: '午餐吃了狗粮', time: '2026-03-28 12:00:00' },
  { id: 'M003', type: 'walk', typeText: '散步', content: '公园散步30分钟', time: '2026-03-28 15:00:00' },
]);

const longTermMemory = ref([
  { id: 'L001', type: 'command', typeText: '指令', content: '学会了坐下指令', learnedAt: '2026-02-15' },
  { id: 'L002', type: 'command', typeText: '指令', content: '学会了握手指令', learnedAt: '2026-02-20' },
  { id: 'L003', type: 'experience', typeText: '经历', content: '第一次去海边', learnedAt: '2026-01-10' },
]);

const learningProgress = ref([
  { action: '坐下', mastery: 95 },
  { action: '握手', mastery: 88 },
  { action: '卧倒', mastery: 72 },
  { action: '翻滚', mastery: 45 },
  { action: '装死', mastery: 30 },
]);

const activityMarks = { 0: '安静', 50: '适中', 100: '活跃' };
const affectionMarks = { 0: '独立', 50: '一般', 100: '亲人' };
const difficultyMarks = { 0: '容易', 50: '中等', 100: '困难' };
const sensitivityMarks = { 0: '迟钝', 50: '适中', 100: '灵敏' };

const shortTermColumns = [
  { title: '类型', slotName: 'type', width: 100 },
  { title: '内容', dataIndex: 'content' },
  { title: '时间', dataIndex: 'time', width: 160 },
];

const longTermColumns = [
  { title: '类型', slotName: 'type', width: 100 },
  { title: '内容', dataIndex: 'content' },
  { title: '学习时间', dataIndex: 'learnedAt', width: 120 },
];

const getMemoryTypeColor = (type: string) => {
  const map: Record<string, string> = { play: 'blue', food: 'green', walk: 'orange', command: 'purple', experience: 'cyan' };
  return map[type] || 'default';
};

const handleSaveBasic = () => {};
const handleResetBasic = () => {};
const handleSavePersonality = () => {};
const handleResetPersonality = () => {};
const handleSaveInteraction = () => {};
</script>

<style scoped>
.pet-config-container { padding: 20px; }
</style>
