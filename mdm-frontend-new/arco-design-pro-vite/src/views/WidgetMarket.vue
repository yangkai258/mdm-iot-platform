<template>
  <div class="container">
    <a-card>
      <template #title><icon-app-store /> 小部件市场</template>
      <a-list>
        <a-list-item v-for="widget in widgets" :key="widget.id">
          <a-list-item-meta :title="widget.name" :description="widget.description" />
          <template #actions><a-button size="small">安装</a-button></template>
        </a-list-item>
      </a-list>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const widgets = ref([
  { id: 1, name: '设备状态卡片', description: '显示设备实时状态' },
  { id: 2, name: '告警统计', description: '显示告警统计数据' }
])
</script>
<style scoped>.container { padding: 16px; }</style>
