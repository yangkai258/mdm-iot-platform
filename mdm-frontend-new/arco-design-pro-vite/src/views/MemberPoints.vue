<template>
  <div class="container">
    <a-card>
      <template #title><icon-star /> 会员积分管理</template>
      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6"><a-card><a-statistic title="总积分" :value="stats.total" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="今日发放" :value="stats.todayGrant" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="今日消耗" :value="stats.todayUse" /></a-card></a-col>
      </a-row>
      <a-table :columns="columns" :data="records" />
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const stats = reactive({ total: 125680, todayGrant: 520, todayUse: 380 })
const columns = [{ title: '会员', dataIndex: 'member' }, { title: '积分变动', dataIndex: 'change' }, { title: '时间', dataIndex: 'time' }]
const records = ref([{ member: '张三', change: '+100', time: '2026-03-28' }])
</script>
<style scoped>.container { padding: 16px; }</style>
