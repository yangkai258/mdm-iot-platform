<template>
  <div class="container">
    <a-card>
      <template #title><icon-home /> 宠物领养</template>
      <a-list>
        <a-list-item v-for="pet in pets" :key="pet.id">
          <a-list-item-meta :title="pet.name" :description="pet.breed + ' | ' + pet.age + '岁'" />
          <template #actions><a-button size="small" type="primary">申请领养</a-button></template>
        </a-list-item>
      </a-list>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const pets = ref([{ id: 1, name: '小黄', breed: '金毛', age: 1 }])
</script>
<style scoped>.container { padding: 16px; }</style>
