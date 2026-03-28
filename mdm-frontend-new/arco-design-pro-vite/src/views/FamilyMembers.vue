<template>
  <div class="container">
    <a-card>
      <template #title><icon-team /> 家庭成员</template>
      <template #extra><a-button type="primary" @click="handleAdd"><template #icon><icon-plus /></template>添加成员</a-button></template>
      <a-list>
        <a-list-item v-for="member in members" :key="member.id">
          <a-list-item-meta :title="member.name" :description="member.role" />
          <template #actions>
            <a-link @click="handleEdit(member)">编辑</a-link>
          </template>
        </a-list-item>
      </a-list>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const members = ref([
  { id: 1, name: '张三', role: '管理员' },
  { id: 2, name: '李四', role: '成员' }
])
const handleAdd = () => {}
const handleEdit = (m) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
