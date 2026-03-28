<template>
  <div class="container">
    <a-card>
      <template #title><icon-cloud-download /> 固件仓库</template>
      <template #extra><a-button type="primary" @click="handleUpload"><template #icon><icon-upload /></template>上传固件</a-button></template>
      <a-table :columns="columns" :data="firmwares">
        <template #status="{ record }"><a-tag :color="record.recommended ? 'green' : 'gray'">{{ record.recommended ? '推荐' : '' }}</a-tag></template>
        <template #actions="{ record }"><a-button size="small">下载</a-button></template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '版本', dataIndex: 'version' }, { title: '设备型号', dataIndex: 'model' }, { title: '大小', dataIndex: 'size' }, { title: '状态', slotName: 'status' }, { title: '操作', slotName: 'actions' }]
const firmwares = ref([{ id: 1, version: 'v2.1.0', model: 'M5Stack-1', size: '15MB', recommended: true }])
const handleUpload = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
