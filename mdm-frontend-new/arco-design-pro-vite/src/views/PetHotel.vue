<template>
  <div class="container">
    <a-card>
      <template #title><icon-home /> 宠物酒店</template>
      <a-list>
        <a-list-item v-for="room in rooms" :key="room.id">
          <a-list-item-meta :title="room.name" :description="room.price + '/天'" />
          <template #actions><a-button size="small" type="primary">预订</a-button></template>
        </a-list-item>
      </a-list>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const rooms = ref([{ id: 1, name: '豪华单间', price: '¥200' }])
</script>
<style scoped>.container { padding: 16px; }</style>
