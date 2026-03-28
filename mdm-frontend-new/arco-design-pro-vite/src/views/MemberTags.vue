<template>
  <div class="container">
    <a-card>
      <template #title><icon-bookmark /> 会员标签管理</template>
      <template #extra>
        <a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>新建标签</a-button>
      </template>
      <a-list>
        <a-list-item v-for="tag in tags" :key="tag.id">
          <a-list-item-meta :title="tag.name" :description="'会员数: ' + tag.count" />
          <template #actions>
            <a-link @click="handleEdit(tag)">编辑</a-link>
            <a-link @click="handleDelete(tag)">删除</a-link>
          </template>
        </a-list-item>
      </a-list>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const tags = ref([{ id: 1, name: '高价值', count: 50 }, { id: 2, name: '活跃用户', count: 120 }])
const handleCreate = () => {}
const handleEdit = (r) => {}
const handleDelete = (r) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
