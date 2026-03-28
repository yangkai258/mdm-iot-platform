<template>
  <div class="container">
    <a-card>
      <template #title><icon-shield /> 宠物保险</template>
      <template #extra><a-button type="primary" @click="handleBuy"><template #icon><icon-plus /></template>购买保险</a-button></template>
      <a-table :columns="columns" :data="insurances">
        <template #status="{ record }"><a-tag :color="record.status === 'active' ? 'green' : 'gray'">{{ record.status === 'active' ? '生效中' : '已过期' }}</a-tag></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '保单号', dataIndex: 'policyNo' }, { title: '宠物', dataIndex: 'petName' }, { title: '保险公司', dataIndex: 'company' }, { title: '状态', slotName: 'status' }]
const insurances = ref([{ id: 1, policyNo: 'POL001', petName: '小白', company: 'XX保险', status: 'active' }])
const handleBuy = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
