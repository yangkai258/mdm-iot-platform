<template>
  <div class="pet-timeline-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="成长记录" :value="5680" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="照片数量" :value="25600" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="视频数量" :value="1280" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <span>宠物成长时间线</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="timeline" title="成长时间线">
          <a-timeline>
            <a-timeline-item v-for="item in timeline" :key="item.id" :color="item.type === 'milestone' ? 'red' : 'blue'">
              <div class="timeline-item">
                <div class="timeline-date">{{ item.date }}</div>
                <div class="timeline-title">{{ item.title }}</div>
                <div class="timeline-desc">{{ item.description }}</div>
                <div v-if="item.media" class="timeline-media">
                  <a-image v-if="item.mediaType === 'image'" :src="item.media" width="200" />
                  <span v-else>{{ item.media }}</span>
                </div>
              </div>
            </a-timeline-item>
          </a-timeline>
        </a-tab-pane>
        
        <a-tab-pane key="milestones" title="里程碑">
          <a-table :columns="milestoneColumns" :data="milestones" :pagination="pagination" />
        </a-tab-pane>
        
        <a-tab-pane key="growth" title="成长报告">
          <a-card title="小黄的成长报告">
            <a-statistic title="年龄" :value="2" suffix="岁" />
            <a-statistic title="体重增长" :value="120" suffix="%" />
            <a-statistic title="技能数" :value="15" />
          </a-card>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 8 });

const timeline = ref([
  { id: 1, date: '2026-03-28', title: '学会了握手', description: '今天小黄学会了握手技能', type: 'skill', media: null, mediaType: null },
  { id: 2, date: '2026-03-20', title: '体重记录', description: '体重增长到5.2kg', type: 'health', media: null, mediaType: null },
  { id: 3, date: '2026-03-15', title: '里程碑: 2岁生日', description: '小黄2岁啦！', type: 'milestone', media: null, mediaType: null },
]);

const milestones = ref([
  { id: 1, title: '出生', date: '2024-03-15', description: '小黄出生' },
  { id: 2, title: '第一次走路', date: '2024-05-01', description: '学会走路' },
  { id: 3, title: '2岁生日', date: '2026-03-15', description: '2岁生日派对' },
]);

const milestoneColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '标题', dataIndex: 'title', width: 150 },
  { title: '日期', dataIndex: 'date', width: 120 },
  { title: '描述', dataIndex: 'description' },
];

const handleAddMilestone = () => {};
</script>

<style scoped>
.pet-timeline-container { padding: 20px; }
.timeline-item { padding: 8px 0; }
.timeline-date { color: #86909c; font-size: 12px; }
.timeline-title { font-weight: bold; margin: 4px 0; }
.timeline-desc { color: #4a4a4a; }
.timeline-media { margin-top: 8px; }
</style>
