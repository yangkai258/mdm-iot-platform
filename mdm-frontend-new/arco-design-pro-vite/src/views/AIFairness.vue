<template>
  <div class="container">
    <a-card>
      <template #title><icon-balance /> AI公平性检测</template>
      <a-row :gutter="16">
        <a-col :span="6"><a-card><a-statistic title="公平性评分" :value="stats.fairnessScore" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="偏差率" :value="stats.biasRate" suffix="%" /></a-card></a-col>
      </a-row>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const stats = reactive({ fairnessScore: 92, biasRate: 2.5 })
</script>
<style scoped>.container { padding: 16px; }</style>
