<template>
  <div class="container">
    <a-card>
      <template #title><icon-shopping-cart /> 会员订单</template>
      <template #extra>
        <a-button @click="handleExport"><template #icon><icon-download /></template>导出</a-button>
      </template>
      <a-table :columns="columns" :data="orders">
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ record.status }}</a-tag>
        </template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '订单号', dataIndex: 'orderId' }, { title: '会员', dataIndex: 'member' }, { title: '金额', dataIndex: 'amount' }, { title: '状态', slotName: 'status' }]
const orders = ref([{ orderId: 'O001', member: '张三', amount: '¥99', status: '已完成' }])
const getStatusColor = (s) => ({ '已完成': 'green', '进行中': 'blue', '已取消': 'gray' }[s] || 'gray')
const handleExport = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
