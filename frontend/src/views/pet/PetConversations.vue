<template>
  <div class="pet-conversations">
    
    <div class="page-content"><a-card>
        <template #title>
          <a-space>
            <span>瀵硅瘽鍒楄〃</span>
            <a-button type="primary" @click="createConversation">鏂板缓瀵硅瘽</a-button>
          </a-space>
        </template>
        <a-table :columns="columns" :data="conversations" row-key="id">
          <template #status="{ record }">
            <a-tag :color="record.status === 'active' ? 'green' : 'gray'">
              {{ record.status === 'active' ? '杩涜涓? : '宸茬粨鏉? }}
            </a-tag>
          </template>
        </a-table>
      </a-card></div>
  </div>
</template>

<script>
export default {
  name: 'PetConversations',
  data() {
    return {
      columns: [
        { title: 'ID', dataIndex: 'id' },
        { title: '璁惧ID', dataIndex: 'device_id' },
        { title: '鏍囬', dataIndex: 'title' },
        { title: '鐘舵€?, slotName: 'status' },
        { title: '鍒涘缓鏃堕棿', dataIndex: 'created_at' }
      ],
      conversations: []
    }
  },
  methods: {
    createConversation() {
      this.$message.info('鍒涘缓鏂板璇?)
    }
  }
}
</script>
