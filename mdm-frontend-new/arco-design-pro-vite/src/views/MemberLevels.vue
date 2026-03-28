<template>
  <div class="member-levels-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="等级数量" :value="4" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="升级会员数" :value="28" suffix="本月" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="降级会员数" :value="5" suffix="本月" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>会员等级权益配置</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新增等级
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #icon="{ record }">
          <div class="level-icon" :style="{ backgroundColor: record.color }">{{ record.icon }}</div>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.enabled ? 'green' : 'gray'">{{ record.enabled ? '启用' : '禁用' }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link @click="handleViewMembers(record)">查看会员</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 创建/编辑弹窗 -->
    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑等级' : '新增等级'" :width="700" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="等级名称" required>
          <a-input v-model="form.name" placeholder="如: 金卡会员" />
        </a-form-item>
        <a-form-item label="等级图标">
          <a-input v-model="form.icon" placeholder="图标emoji" />
        </a-form-item>
        <a-form-item label="等级颜色">
          <a-color-picker v-model="form.color" />
        </a-form-item>
        <a-form-item label="升级门槛">
          <a-space vertical>
            <span>成长值 >= </span>
            <a-input-number v-model="form.threshold" :min="0" />
          </a-space>
        </a-form-item>
        <a-form-item label="享受折扣">
          <a-input-number v-model="form.discount" :min="0" :max="100" :precision="0" suffix="%" />
        </a-form-item>
        <a-form-item label="专属权益">
          <a-checkbox-group v-model="form.benefits">
            <a-checkbox value="free_shipping">免运费</a-checkbox>
            <a-checkbox value="priority_service">优先客服</a-checkbox>
            <a-checkbox value="exclusive_events">专属活动</a-checkbox>
            <a-checkbox value="points_bonus">积分加倍</a-checkbox>
            <a-checkbox value="birthday_gift">生日礼品</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
        <a-form-item label="启用">
          <a-switch v-model="form.enabled" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);

const data = ref([
  { id: 1, name: '普通会员', icon: '🪪', color: '#86909c', threshold: 0, discount: 100, benefits: [], memberCount: 1500, enabled: true },
  { id: 2, name: '银卡会员', icon: '🥈', color: '#8c8c8c', threshold: 1000, discount: 95, benefits: ['free_shipping', 'points_bonus'], memberCount: 580, enabled: true },
  { id: 3, name: '金卡会员', icon: '🥇', color: '#ffc532', threshold: 5000, discount: 90, benefits: ['free_shipping', 'priority_service', 'points_bonus'], memberCount: 156, enabled: true },
  { id: 4, name: '黑金会员', icon: '👑', color: '#1a1a1a', threshold: 20000, discount: 85, benefits: ['free_shipping', 'priority_service', 'exclusive_events', 'points_bonus', 'birthday_gift'], memberCount: 23, enabled: true },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 4 });
const editVisible = ref(false);
const isEdit = ref(false);

const form = reactive({ name: '', icon: '', color: '#165DFF', threshold: 0, discount: 100, benefits: [], enabled: true });

const columns = [
  { title: '图标', slotName: 'icon', width: 80 },
  { title: '等级名称', dataIndex: 'name', width: 120 },
  { title: '升级门槛', dataIndex: 'threshold', width: 120 },
  { title: '折扣', dataIndex: 'discount', width: 80 },
  { title: '会员数', dataIndex: 'memberCount', width: 100 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleViewMembers = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.member-levels-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
.level-icon { width: 40px; height: 40px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 20px; }
</style>
