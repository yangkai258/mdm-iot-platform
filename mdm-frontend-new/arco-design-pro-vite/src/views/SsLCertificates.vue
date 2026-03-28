<template>
  <div class="container">
    <a-card>
      <template #title><icon-safe /> SSL证书</template>
      <template #extra><a-button type="primary" @click="handleUpload"><template #icon><icon-upload /></template>上传证书</a-button></template>
      <a-table :columns="columns" :data="certificates">
        <template #status="{ record }"><a-tag :color="record.valid ? 'green' : 'red'">{{ record.valid ? '有效' : '过期' }}</a-tag></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '域名', dataIndex: 'domain' }, { title: '有效期', dataIndex: 'expireDate' }, { title: '状态', slotName: 'status' }]
const certificates = ref([{ id: 1, domain: 'example.com', expireDate: '2027-01-01', valid: true }])
const handleUpload = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
