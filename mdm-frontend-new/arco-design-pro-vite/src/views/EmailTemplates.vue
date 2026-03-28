<template>
  <div class="container">
    <a-card>
      <template #title><icon-email /> 邮件模板</template>
      <template #extra>
        <a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>新建模板</a-button>
      </template>
      <a-table :columns="columns" :data="templates">
        <template #actions="{ record }">
          <a-link @click="handlePreview(record)">预览</a-link>
          <a-link @click="handleEdit(record)">编辑</a-link>
        </template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const columns = [{ title: '模板名称', dataIndex: 'name' }, { title: '类型', dataIndex: 'type' }, { title: '操作', slotName: 'actions' }]
const templates = ref([{ id: 1, name: '欢迎邮件', type: 'welcome' }])
const handleCreate = () => {}
const handlePreview = (r) => {}
const handleEdit = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
