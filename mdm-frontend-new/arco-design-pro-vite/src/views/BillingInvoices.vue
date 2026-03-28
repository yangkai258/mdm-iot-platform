<template>
  <div class="container">
    <a-card>
      <template #title><icon-file-text /> 账单发票</template>
      <template #extra><a-button @click="handleExport"><template #icon><icon-download /></template>导出</a-button></template>
      <a-table :columns="columns" :data="invoices">
        <template #status="{ record }"><a-tag :color="record.status === 'paid' ? 'green' : 'orange'">{{ record.status === 'paid' ? '已支付' : '待支付' }}</a-tag></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '发票号', dataIndex: 'invoiceId' }, { title: '金额', dataIndex: 'amount' }, { title: '状态', slotName: 'status' }, { title: '时间', dataIndex: 'time' }]
const invoices = ref([{ invoiceId: 'INV001', amount: '¥99', status: 'paid', time: '2026-03-28' }])
const handleExport = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
