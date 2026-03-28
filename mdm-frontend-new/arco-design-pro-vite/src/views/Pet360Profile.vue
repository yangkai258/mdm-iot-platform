<template>
  <div class="container">
    <a-card>
      <template #title><icon-heart /> 宠物360画像</template>
      <a-descriptions :column="3" bordered>
        <a-descriptions-item label="宠物名">{{ pet.name }}</a-descriptions-item>
        <a-descriptions-item label="品种">{{ pet.breed }}</a-descriptions-item>
        <a-descriptions-item label="年龄">{{ pet.age }}岁</a-descriptions-item>
      </a-descriptions>
      <a-row :gutter="16" style="margin-top: 16px">
        <a-col :span="8"><a-card><a-statistic title="活跃度" :value="stats.activity" /></a-card></a-col>
        <a-col :span="8"><a-card><a-statistic title="健康指数" :value="stats.health" /></a-card></a-col>
        <a-col :span="8"><a-card><a-statistic title="情绪指数" :value="stats.emotion" /></a-card></a-col>
      </a-row>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const pet = reactive({ name: '小白', breed: '萨摩耶', age: 2 })
const stats = reactive({ activity: 85, health: 92, emotion: 88 })
</script>
<style scoped>.container { padding: 16px; }</style>
