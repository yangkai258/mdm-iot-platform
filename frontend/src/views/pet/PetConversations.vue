<template>
  <div class="pet-conversations">
    <a-card class="general-card" title="会话列表">
      <template #extra>
        <a-button type="primary" @click="createConversation"><icon-plus />新建会话</a-button>
      </template>
      <a-table :columns="columns" :data="conversations" row-key="id">
        <template #status="{ record }">
          <a-tag :color="record.status === 'active' ? 'green' : 'gray'">
            {{ record.status === 'active' ? '进行中' : '已结束' }}
          </a-tag>
        </template>
      </a-table>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'

const columns = [
  { title: 'ID', dataIndex: 'id' },
  { title: '设备ID', dataIndex: 'device_id' },
  { title: '主题', dataIndex: 'title' },
  { title: '状态', slotName: 'status' },
  { title: '创建时间', dataIndex: 'created_at' }
]

const conversations = ref([])

const createConversation = () => {
  Message.info('新建会话功能')
}
</script>

<style scoped>
.pet-conversations {
  padding: 16px;
}
</style>
