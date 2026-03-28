<template>
  <div class="member-recharge-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="今日充值" :value="stats.todayRecharge" prefix="¥" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="本月充值" :value="stats.monthRecharge" prefix="¥" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="储值会员数" :value="stats.rechargeMembers" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="平均储值" :value="stats.avgBalance" prefix="¥" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>会员储值管理</span>
          <a-button type="primary" @click="handleRecharge">
            <template #icon><icon-plus /></template>
            充值
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #level="{ record }">
          <a-tag :color="getLevelColor(record.level)">{{ record.levelText }}</a-tag>
        </template>
        <template #balance="{ record }">
          <span style="color: #F53F3F; font-weight: bold;">¥{{ record.balance }}</span>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleDetail(record)">明细</a-link>
            <a-link @click="handleRechargeMember(record)">充值</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 充值弹窗 -->
    <a-modal v-model:visible="rechargeVisible" title="会员充值" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="会员" required>
          <a-select v-model="form.memberId" placeholder="选择会员" show-search>
            <a-option value="1">张三 - 138****1234</a-option>
            <a-option value="2">李四 - 139****5678</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="充值金额" required>
          <a-input-number v-model="form.amount" :min="0" :precision="2" style="width: 100%;" />
        </a-form-item>
        <a-form-item label="赠送金额">
          <a-input-number v-model="form.bonus" :min="0" :precision="2" style="width: 100%;" />
        </a-form-item>
        <a-form-item label="支付方式">
          <a-select v-model="form.paymentMethod" placeholder="选择支付方式">
            <a-option value="wechat">微信支付</a-option>
            <a-option value="alipay">支付宝</a-option>
            <a-option value="cash">现金</a-option>
            <a-option value="card">银行卡</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="form.note" :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 储值明细弹窗 -->
    <a-modal v-model:visible="detailVisible" title="储值明细" :width="800">
      <a-table :columns="detailColumns" :data="rechargeDetails" :pagination="paginationSmall">
        <template #type="{ record }">
          <a-tag :color="record.type === 'recharge' ? 'green' : 'blue'">{{ record.typeText }}</a-tag>
        </template>
      </a-table>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);

const stats = reactive({ todayRecharge: 5680, monthRecharge: 156780, rechargeMembers: 156, avgBalance: 856 });

const data = ref([
  { id: 1, name: '张三', phone: '138****1234', level: 3, levelText: '金卡', balance: 2500, rechargeCount: 5, lastRechargeAt: '2026-03-28' },
  { id: 2, name: '李四', phone: '139****5678', level: 2, levelText: '银卡', balance: 1200, rechargeCount: 3, lastRechargeAt: '2026-03-25' },
  { id: 3, name: '王五', phone: '137****9012', level: 1, levelText: '普通', balance: 500, rechargeCount: 1, lastRechargeAt: '2026-03-20' },
]);

const rechargeDetails = ref([
  { id: 1, type: 'recharge', typeText: '充值', amount: 1000, bonus: 50, balance: 1050, time: '2026-03-28 10:00:00', paymentMethod: '微信支付' },
  { id: 2, type: 'consume', typeText: '消费', amount: -200, bonus: 0, balance: 850, time: '2026-03-28 15:00:00', paymentMethod: '-' },
  { id: 3, type: 'recharge', typeText: '充值', amount: 500, bonus: 0, balance: 1050, time: '2026-03-27 10:00:00', paymentMethod: '现金' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 3 });
const paginationSmall = reactive({ current: 1, pageSize: 10, total: 3 });
const rechargeVisible = ref(false);
const detailVisible = ref(false);

const form = reactive({ memberId: '', amount: 0, bonus: 0, paymentMethod: 'wechat', note: '' });

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '姓名', dataIndex: 'name', width: 100 },
  { title: '手机', dataIndex: 'phone', width: 130 },
  { title: '等级', slotName: 'level', width: 80 },
  { title: '储值余额', slotName: 'balance', width: 120 },
  { title: '充值次数', dataIndex: 'rechargeCount', width: 100 },
  { title: '最后充值', dataIndex: 'lastRechargeAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const detailColumns = [
  { title: '时间', dataIndex: 'time', width: 160 },
  { title: '类型', slotName: 'type', width: 80 },
  { title: '金额', dataIndex: 'amount', width: 100 },
  { title: '赠送', dataIndex: 'bonus', width: 80 },
  { title: '余额', dataIndex: 'balance', width: 100 },
  { title: '支付方式', dataIndex: 'paymentMethod', width: 120 },
];

const getLevelColor = (l: number) => l >= 3 ? 'gold' : l >= 2 ? 'silver' : 'gray';

const handleRecharge = () => { rechargeVisible.value = true; };
const handleRechargeMember = (record: any) => { form.memberId = String(record.id); rechargeVisible.value = true; };
const handleDetail = (record: any) => { detailVisible.value = true; };
const handleSubmit = (done: boolean) => { done(true); rechargeVisible.value = false; };
</script>

<style scoped>
.member-recharge-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
