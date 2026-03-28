<template>
  <div class="member-growth-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="新增会员" :value="156" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="会员增长率" :value="12.5" suffix="%" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="会员活跃度" :value="78" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <span>会员增长分析</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="overview" title="增长概览">
          <a-row :gutter="16">
            <a-col :span="12">
              <a-card title="会员趋势">
                <a-timeline>
                  <a-timeline-item v-for="item in trendData" :key="item.date" :color="item.growth > 0 ? 'green' : 'red'">
                    {{ item.date }}: 新增{{ item.new }}人, 增长率{{ item.growth }}%
                  </a-timeline-item>
                </a-timeline>
              </a-card>
            </a-col>
            <a-col :span="12">
              <a-card title="来源分布">
                <a-progress :percent="40" label="自然增长" />
                <a-progress :percent="30" label="口碑推荐" />
                <a-progress :percent="20" label="营销活动" />
                <a-progress :percent="10" label="其他渠道" />
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
        
        <a-tab-pane key="funnel" title="转化漏斗">
          <a-card title="会员转化漏斗">
            <a-steps :current="3" style="max-width: 600px; margin: 0 auto;">
              <a-step title="访问" description="10,000人" />
              <a-step title="注册" description="2,000人" />
              <a-step title="活跃" description="500人" />
              <a-step title="付费" description="156人" />
            </a-steps>
          </a-card>
        </a-tab-pane>
        
        <a-tab-pane key="churn" title="流失分析">
          <a-table :columns="churnColumns" :data="churnData" :pagination="pagination" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 5 });

const trendData = ref([
  { date: '2026-03-28', new: 28, growth: 15.2 },
  { date: '2026-03-27', new: 32, growth: 8.5 },
  { date: '2026-03-26', new: 25, growth: -5.2 },
  { date: '2026-03-25', new: 30, growth: 12.8 },
]);

const churnData = ref([
  { id: 1, reason: '不需要了', count: 45, percent: 35 },
  { id: 2, reason: '价格太高', count: 32, percent: 25 },
  { id: 3, reason: '质量问题', count: 20, percent: 15 },
  { id: 4, reason: '转用竞品', count: 15, percent: 12 },
]);

const churnColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '流失原因', dataIndex: 'reason', width: 200 },
  { title: '人数', dataIndex: 'count', width: 100 },
  { title: '占比', dataIndex: 'percent', width: 100 },
];
</script>

<style scoped>
.member-growth-container { padding: 20px; }
</style>
