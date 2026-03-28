<template>
  <div class="container">
    <a-card>
      <template #title><icon-pet /> 宠物详情</template>
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="宠物ID">{{ pet.id }}</a-descriptions-item>
        <a-descriptions-item label="名称">{{ pet.name }}</a-descriptions-item>
        <a-descriptions-item label="品种">{{ pet.breed }}</a-descriptions-item>
        <a-descriptions-item label="年龄">{{ pet.age }}岁</a-descriptions-item>
      </a-descriptions>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const pet = reactive({ id: 'P001', name: '小白', breed: '萨摩耶', age: 2 })
</script>
<style scoped>.container { padding: 16px; }</style>
