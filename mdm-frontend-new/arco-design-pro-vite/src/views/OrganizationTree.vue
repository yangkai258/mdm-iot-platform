<template>
  <div class="container">
    <a-card>
      <template #title><icon-scope /> 组织架构树</template>
      <a-tree :data="treeData" @select="handleSelect">
        <template #title="{ data }">{{ data.title }}</template>
      </a-tree>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const treeData = [{ title: '总公司', key: '1', children: [{ title: '研发部', key: '2' }, { title: '运维部', key: '3' }] }]
const handleSelect = (keys) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
