<template>
  <div class="container">
    <a-card>
      <template #title><icon-star /> 会员忠诚度计划</template>
      <a-row :gutter="16">
        <a-col :span="6"><a-card><a-statistic title="总积分" :value="stats.totalPoints" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="活跃会员" :value="stats.activeMembers" /></a-card></a-col>
      </a-row>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const stats = reactive({ totalPoints: 1256800, activeMembers: 856 })
</script>
<style scoped>.container { padding: 16px; }</style>
