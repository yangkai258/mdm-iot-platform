<template>
  <div class="cost-analysis-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="本月成本" :value="25680" prefix="¥" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="环比增长" :value="5.2" suffix="%" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="预算使用" :value="68" suffix="%" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="成本优化" :value="12" suffix="%可优化" /></a-card>
      </a-col>
    </a-row>
    <a-card>
      <template #title>成本分析</template>
      <a-tabs>
        <a-tab-pane key="overview" title="成本概览">
          <a-row :gutter="16">
            <a-col :span="12">
              <a-card title="成本分布">
                <a-progress :percent="40" label="计算资源" />
                <a-progress :percent="25" label="存储" />
                <a-progress :percent="20" label="网络" />
                <a-progress :percent="15" label="其他" />
              </a-card>
            </a-col>
            <a-col :span="12">
              <a-card title="趋势">
                <a-timeline><a-timeline-item>2026-03: ¥25,680</a-timeline-item><a-timeline-item>2026-02: ¥24,420</a-timeline-item></a-timeline>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
        <a-tab-pane key="details" title="成本明细">
          <a-table :columns="columns" :data="details" :pagination="pagination" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
const pagination = reactive({ current: 1, pageSize: 10, total: 8 });
const details = ref([
  { id: 1, service: 'ECS', cost: 10272, percentage: 40, trend: 'up' },
  { id: 2, service: 'OSS', cost: 6420, percentage: 25, trend: 'stable' },
]);
const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '服务', dataIndex: 'service', width: 150 },
  { title: '成本', dataIndex: 'cost', width: 120 },
  { title: '占比', dataIndex: 'percentage', width: 100 },
  { title: '趋势', dataIndex: 'trend', width: 100 },
];
</script>

<style scoped>
.cost-analysis-container { padding: 20px; }
</style>
