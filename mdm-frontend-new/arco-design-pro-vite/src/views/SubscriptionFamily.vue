<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-home /> 家庭计划管理</a-space>
      </template>

      <a-row :gutter="16">
        <a-col :span="8">
          <a-card title="计划概览">
            <a-statistic title="共享时长" :value="plan.sharedDays" suffix="天" />
            <a-divider />
            <a-statistic title="已用时长" :value="plan.usedDays" suffix="天" />
            <a-divider />
            <a-statistic title="成员数" :value="plan.memberCount" :suffix="'/' + plan.maxMembers" />
            <a-divider />
            <a-button type="primary" long @click="handleInvite">
              <template #icon><icon-plus /></template>
              邀请成员
            </a-button>
          </a-card>
        </a-col>

        <a-col :span="16">
          <a-card title="成员列表">
            <template #extra>
              <a-button @click="handleQuotaConfig">
                <template #icon><icon-settings /></template>
                管理配额
              </a-button>
            </template>
            <a-table :columns="columns" :data="members">
              <template #role="{ record }">
                <a-tag :color="record.role === 'owner' ? 'gold' : 'blue'">{{ record.role === 'owner' ? '所有者' : '成员' }}</a-tag>
              </template>
              <template #usage="{ record }">
                <a-progress :percent="record.usedQuota / record.quotaLimit * 100" :show-text="true" size="small" />
                <span class="usage-text">{{ record.usedQuota }} / {{ record.quotaLimit }}</span>
              </template>
              <template #actions="{ record }">
                <a-link v-if="record.role !== 'owner'" @click="handleRemove(record)">移除</a-link>
              </template>
            </a-table>
          </a-card>

          <a-card title="配额分配" style="margin-top: 16px">
            <a-form :model="quotaForm" layout="inline">
              <a-form-item label="总配额">
                <a-input-number v-model="quotaForm.totalQuota" :min="0" suffix="天" />
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="handleSaveQuota">保存</a-button>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const plan = reactive({ sharedDays: 30, usedDays: 12, memberCount: 3, maxMembers: 5 })
const quotaForm = reactive({ totalQuota: 30 })

const columns = [
  { title: '用户', dataIndex: 'userName' },
  { title: '角色', slotName: 'role' },
  { title: '加入时间', dataIndex: 'joinedAt' },
  { title: '配额使用', slotName: 'usage' },
  { title: '操作', slotName: 'actions', width: 100 }
]

const members = ref([
  { userName: '张三', role: 'owner', joinedAt: '2026-03-01', usedQuota: 10, quotaLimit: 30 },
  { userName: '李四', role: 'member', joinedAt: '2026-03-15', usedQuota: 5, quotaLimit: 15 },
  { userName: '王五', role: 'member', joinedAt: '2026-03-20', usedQuota: 2, quotaLimit: 10 }
])

const handleInvite = () => { }
const handleQuotaConfig = () => { }
const handleRemove = (r) => { }
const handleSaveQuota = () => { }
</script>

<style scoped>
.container { padding: 16px; }
.usage-text { font-size: 12px; color: #909399; }
</style>
