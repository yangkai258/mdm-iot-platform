<template>
  <div class="loyalty-program-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="会员总数" :value="8560" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="积分发放" :value="256000" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="积分消耗" :value="198000" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="活跃会员" :value="5680" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>会员积分体系</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建规则
          </a-button>
        </div>
      </template>
      
      <a-tabs>
        <a-tab-pane key="rules" title="积分规则">
          <a-table :columns="ruleColumns" :data="rules" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="record.enabled ? 'green' : 'gray'">{{ record.enabled ? '启用' : '禁用' }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="records" title="积分记录">
          <a-table :columns="recordColumns" :data="records" :pagination="pagination" />
        </a-tab-pane>
        
        <a-tab-pane key="levels" title="会员等级">
          <a-row :gutter="16">
            <a-col :span="6" v-for="level in levels" :key="level.id">
              <a-card class="level-card">
                <div class="level-name">{{ level.name }}</div>
                <div class="level-range">{{ level.min }}-{{ level.max }}积分</div>
                <div class="level-benefits">
                  <div v-for="benefit in level.benefits" :key="benefit">✓ {{ benefit }}</div>
                </div>
                <a-button type="primary" long size="small">查看详情</a-button>
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

const pagination = reactive({ current: 1, pageSize: 10, total: 6 });

const rules = ref([
  { id: 1, name: '消费得积分', action: 'consume', points: 100, description: '每消费1元得1积分', enabled: true },
  { id: 2, name: '签到得积分', action: 'checkin', points: 10, description: '每日签到得10积分', enabled: true },
  { id: 3, name: '推荐得积分', action: 'referral', points: 500, description: '推荐新会员得500积分', enabled: true },
]);

const records = ref([
  { id: 1, memberName: '张三', type: 'earn', points: 100, action: '消费', balance: 1500, time: '2026-03-28 10:00:00' },
  { id: 2, memberName: '李四', type: 'redeem', points: -200, action: '兑换', balance: 800, time: '2026-03-28 09:00:00' },
]);

const levels = ref([
  { id: 1, name: '普通会员', min: 0, max: 1000, benefits: ['1.2x积分', '专属客服'] },
  { id: 2, name: '银卡会员', min: 1000, max: 5000, benefits: ['1.5x积分', '9折优惠', '专属客服'] },
  { id: 3, name: '金卡会员', min: 5000, max: 20000, benefits: ['2x积分', '8折优惠', '优先发货', '专属客服'] },
  { id: 4, name: '钻石会员', min: 20000, max: 999999, benefits: ['3x积分', '7折优惠', '优先发货', '专属客服', '免费礼品'] },
]);

const ruleColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '规则名称', dataIndex: 'name', width: 150 },
  { title: '行为', dataIndex: 'action', width: 100 },
  { title: '积分', dataIndex: 'points', width: 80 },
  { title: '描述', dataIndex: 'description' },
  { title: '状态', slotName: 'status', width: 80 },
];

const recordColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '会员', dataIndex: 'memberName', width: 100 },
  { title: '类型', dataIndex: 'type', width: 80 },
  { title: '积分', dataIndex: 'points', width: 100 },
  { title: '行为', dataIndex: 'action', width: 100 },
  { title: '余额', dataIndex: 'balance', width: 100 },
  { title: '时间', dataIndex: 'time', width: 160 },
];

const handleCreate = () => {};
</script>

<style scoped>
.loyalty-program-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
.level-card { text-align: center; }
.level-name { font-weight: bold; font-size: 18px; color: #1659d5; }
.level-range { color: #86909c; margin: 8px 0; }
.level-benefits { text-align: left; margin: 16px 0; color: #4a4a4a; }
.level-benefits div { margin: 4px 0; }
</style>
