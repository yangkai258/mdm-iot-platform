<template>
  <div class="pet-conversations">
    <a-layout-header class="header">
      <div class="header-title">对话列表</div>
    </a-layout-header>
    <a-layout-content class="content">
      <a-card>
        <template #title>
          <a-space>
            <span>对话列表</span>
            <a-button type="primary" @click="createConversation">新建对话</a-button>
          </a-space>
        </template>
        <a-table :columns="columns" :data="conversations" row-key="id">
          <template #status="{ record }">
            <a-tag :color="record.status === 'active' ? 'green' : 'gray'">
              {{ record.status === 'active' ? '进行中' : '已结束' }}
            </a-tag>
          </template>
        </a-table>
      </a-card>
    </a-layout-content>
  </div>
</template>

<script>
export default {
  name: 'PetConversations',
  data() {
    return {
      columns: [
        { title: 'ID', dataIndex: 'id' },
        { title: '设备ID', dataIndex: 'device_id' },
        { title: '标题', dataIndex: 'title' },
        { title: '状态', slotName: 'status' },
        { title: '创建时间', dataIndex: 'created_at' }
      ],
      conversations: []
    }
  },
  methods: {
    createConversation() {
      this.$message.info('创建新对话')
    }
  }
}
</script>
