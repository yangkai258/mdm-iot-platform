<template>
  <div class="campaigns-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>促销活动管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建活动
          </a-button>
        </div>
      </template>
      
      <a-tabs v-model="activeTab">
        <a-tab-pane key="list" title="活动列表">
          <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
            <template #type="{ record }">
              <a-tag :color="getTypeColor(record.type)">{{ record.typeText }}</a-tag>
            </template>
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-link @click="handleView(record)">详情</a-link>
                <a-link @click="handleEdit(record)">编辑</a-link>
                <a-link v-if="record.status === 'draft'" @click="handlePublish(record)">发布</a-link>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="stats" title="活动统计">
          <a-row :gutter="16">
            <a-col :span="6">
              <a-card><a-statistic title="进行中活动" :value="stats.active" /></a-card>
            </a-col>
            <a-col :span="6">
              <a-card><a-statistic title="参与人数" :value="stats.participants" /></a-card>
            </a-col>
            <a-col :span="6">
              <a-card><a-statistic title="优惠券发放" :value="stats.couponsIssued" /></a-card>
            </a-col>
            <a-col :span="6">
              <a-card><a-statistic title="带动销售额" :value="stats.salesAmount" prefix="¥" /></a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <!-- 创建/编辑弹窗 -->
    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑活动' : '创建活动'" :width="700" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="活动名称" required>
          <a-input v-model="form.name" placeholder="请输入活动名称" />
        </a-form-item>
        <a-form-item label="活动类型">
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="discount">打折促销</a-option>
            <a-option value="coupon">发放优惠券</a-option>
            <a-option value="points">积分加倍</a-option>
            <a-option value="gift">赠品活动</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="活动时间">
          <a-space>
            <a-date-picker v-model="form.startAt" show-time placeholder="开始时间" />
            <span>至</span>
            <a-date-picker v-model="form.endAt" show-time placeholder="结束时间" />
          </a-space>
        </a-form-item>
        <a-form-item label="活动规则">
          <a-textarea v-model="form.rule" placeholder="请输入活动规则" :rows="3" />
        </a-form-item>
        <a-form-item label="状态">
          <a-radio-group v-model="form.status">
            <a-radio value="draft">草稿</a-radio>
            <a-radio value="published">已发布</a-radio>
            <a-radio value="paused">已暂停</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const activeTab = ref('list');

const stats = reactive({ active: 5, participants: 1234, couponsIssued: 5678, salesAmount: 156780 });

const data = ref([
  { id: 1, name: '春季大促', type: 'discount', typeText: '打折促销', rule: '全场8折', startAt: '2026-03-01', endAt: '2026-03-31', participants: 500, couponsIssued: 300, status: 'published', statusText: '进行中' },
  { id: 2, name: '新会员礼包', type: 'coupon', typeText: '发放优惠券', rule: '新用户首单立减50元', startAt: '2026-03-15', endAt: '2026-04-15', participants: 234, couponsIssued: 234, status: 'published', statusText: '进行中' },
  { id: 3, name: '会员日双倍积分', type: 'points', typeText: '积分加倍', rule: '每月15日会员日双倍积分', startAt: '2026-03-15', endAt: '2026-12-31', participants: 156, couponsIssued: 0, status: 'published', statusText: '进行中' },
  { id: 4, name: '宠物食品专区', type: 'gift', typeText: '赠品活动', rule: '购买宠物食品满200元赠零食', startAt: '2026-04-01', endAt: '2026-04-30', participants: 0, couponsIssued: 0, status: 'draft', statusText: '草稿' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 4 });
const editVisible = ref(false);
const isEdit = ref(false);

const form = reactive({ name: '', type: 'discount', rule: '', startAt: null, endAt: null, status: 'draft' });

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '活动名称', dataIndex: 'name', width: 150 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '活动时间', dataIndex: 'startAt', width: 220 },
  { title: '参与人数', dataIndex: 'participants', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' },
];

const getTypeColor = (t: string) => ({ discount: 'blue', coupon: 'green', points: 'orange', gift: 'purple' }[t] || 'default');
const getStatusColor = (s: string) => ({ published: 'green', paused: 'orange', draft: 'gray' }[s] || 'default');

const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleView = (record: any) => {};
const handlePublish = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.campaigns-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
