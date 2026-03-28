<template>
  <div class="container">
    <a-card>
      <template #title><icon-folder /> 设备分组管理</template>
      <a-tree-select v-model="selectedGroup" :data="treeData" placeholder="选择分组" style="width: 300px" />
      <a-button type="primary" style="margin-left: 16px">查看分组设备</a-button>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const selectedGroup = ref('')
const treeData = [{ title: '全部设备', value: 'all' }, { title: '测试设备', value: 'test' }]
</script>
<style scoped>.container { padding: 16px; }</style>
