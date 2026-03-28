<template>
  <div class="points-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>积分规则</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新增规则
          </a-button>
        </div>
      </template>
      
      <a-tabs default-active-key="rules">
        <a-tab-pane key="rules" title="积分规则">
          <a-table :columns="ruleColumns" :data="rules" :pagination="pagination">
            <template #action="{ record }">
              <a-tag :color="record.action === 'earn' ? 'green' : 'orange'">
                {{ record.action === 'earn' ? '获取' : '抵扣' }}
              </a-tag>
            </template>
            <template #status="{ record }">
              <a-switch :checked="record.enabled" disabled />
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-link @click="handleEdit(record)">编辑</a-link>
                <a-link status="danger" @click="handleDelete(record)">删除</a-link>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
        <a-tab-pane key="records" title="积分记录">
          <a-table :columns="recordColumns" :data="records" :pagination="pagination">
            <template #type="{ record }">
              <a-tag :color="record.type === 'earn' ? 'green' : 'orange'">
                {{ record.type === 'earn' ? '获取' : '抵扣' }}
              </a-tag>
            </template>
            <template #amount="{ record }">
              <span :style="{ color: record.type === 'earn' ? '#00B42A' : '#FF7D00' }">
                {{ record.type === 'earn' ? '+' : '-' }}{{ record.amount }}
              </span>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const rules = ref([
  { id: 'PR001', name: '每日签到', action: 'earn', points: 10, limit: '1次/天', enabled: true, description: '每日签到获取积分' },
  { id: 'PR002', name: '消费返积分', action: 'earn', points: 1, limit: '1积分/元', enabled: true, description: '每消费1元返1积分' },
  { id: 'PR003', name: '兑换优惠券', action: 'redeem', points: 100, limit: '-', enabled: true, description: '100积分兑换优惠券' },
  { id: 'PR004', name: '积分抽奖', action: 'redeem', points: 50, limit: '-', enabled: false, description: '50积分参与抽奖' },
]);

const records = ref([
  { id: 'REC001', memberName: '张三', type: 'earn', points: 100, amount: 10, source: '消费返积分', createdAt: '2026-03-28 10:00:00' },
  { id: 'REC002', memberName: '张三', type: 'earn', points: 100, amount: 10, source: '每日签到', createdAt: '2026-03-28 09:00:00' },
  { id: 'REC003', memberName: '李四', type: 'redeem', points: 50, amount: 50, source: '兑换优惠券', createdAt: '2026-03-27 16:00:00' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 4 });

const ruleColumns = [
  { title: '规则ID', dataIndex: 'id', width: 80 },
  { title: '规则名称', dataIndex: 'name', width: 150 },
  { title: '类型', slotName: 'action', width: 80 },
  { title: '积分', dataIndex: 'points', width: 80 },
  { title: '限制', dataIndex: 'limit', width: 100 },
  { title: '描述', dataIndex: 'description' },
  { title: '启用', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' },
];

const recordColumns = [
  { title: '记录ID', dataIndex: 'id', width: 80 },
  { title: '会员', dataIndex: 'memberName', width: 100 },
  { title: '类型', slotName: 'type', width: 80 },
  { title: '积分', slotName: 'amount', width: 80 },
  { title: '来源', dataIndex: 'source', width: 150 },
  { title: '时间', dataIndex: 'createdAt', width: 160 },
];

const handleCreate = () => {};
const handleEdit = (record: any) => {};
const handleDelete = (record: any) => {};
</script>

<style scoped>
.points-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
