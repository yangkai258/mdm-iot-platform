<template>
  <div class="behavior-rewards-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="奖励总数" :value="5680" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="已发放" :value="4560" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="积分价值" :value="25600" prefix="¥" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <span>行为奖励中心</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="behaviors" title="行为规则">
          <a-table :columns="behColumns" :data="behaviors" :pagination="pagination">
            <template #enabled="{ record }">
              <a-switch v-model="record.enabled" @change="handleToggle(record)" />
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="records" title="奖励记录">
          <a-table :columns="recordColumns" :data="records" :pagination="pagination" />
        </a-tab-pane>
        
        <a-tab-pane key="leaderboard" title="排行榜">
          <a-table :columns="rankColumns" :data="leaderboard" :pagination="false">
            <template #rank="{ record }">
              <a-tag v-if="record.rank <= 3" :color="['gold', 'silver', 'bronze'][record.rank-1]">第{{ record.rank }}名</a-tag>
              <span v-else>第{{ record.rank }}名</span>
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

const behaviors = ref([
  { id: 1, name: '每日登录', points: 10, enabled: true },
  { id: 2, name: '完成训练', points: 50, enabled: true },
  { id: 3, name: '分享动态', points: 20, enabled: true },
]);

const records = ref([
  { id: 1, petName: '小黄', behavior: '完成训练', points: 50, awardedAt: '2026-03-28 10:00:00' },
]);

const leaderboard = ref([
  { rank: 1, petName: '小黄', ownerName: '张三', totalPoints: 5680 },
  { rank: 2, petName: '旺财', ownerName: '李四', totalPoints: 5200 },
  { rank: 3, petName: '咪咪', ownerName: '王五', totalPoints: 4800 },
]);

const behColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '行为名称', dataIndex: 'name', width: 150 },
  { title: '奖励积分', dataIndex: 'points', width: 100 },
  { title: '启用', slotName: 'enabled', width: 80 },
];

const recordColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '宠物', dataIndex: 'petName', width: 100 },
  { title: '行为', dataIndex: 'behavior', width: 150 },
  { title: '积分', dataIndex: 'points', width: 100 },
  { title: '时间', dataIndex: 'awardedAt', width: 160 },
];

const rankColumns = [
  { title: '排名', slotName: 'rank', width: 100 },
  { title: '宠物', dataIndex: 'petName', width: 120 },
  { title: '主人', dataIndex: 'ownerName', width: 120 },
  { title: '总积分', dataIndex: 'totalPoints', width: 120 },
];

const handleToggle = (record: any) => {};
</script>

<style scoped>
.behavior-rewards-container { padding: 20px; }
</style>
