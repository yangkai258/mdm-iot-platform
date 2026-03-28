<template>
  <div class="container">
    <a-card>
      <template #title><icon-pet /> 宠物列表</template>
      <template #extra><a-button type="primary" @click="handleAdd"><template #icon><icon-plus /></template>添加宠物</a-button></template>
      <a-table :columns="columns" :data="pets">
        <template #actions="{ record }"><a-link @click="handleDetail(record)">详情</a-link></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '宠物ID', dataIndex: 'id' }, { title: '名称', dataIndex: 'name' }, { title: '品种', dataIndex: 'breed' }, { title: '操作', slotName: 'actions' }]
const pets = ref([{ id: 'P001', name: '小白', breed: '萨摩耶' }])
const handleAdd = () => {}
const handleDetail = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
