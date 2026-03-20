<template>
  <div class="knowledge-list">
    
    <div class="page-content"><a-card>
        <template #title>
          <a-space>
            <span>鐭ヨ瘑搴?/span>
            <a-button type="primary" @click="addKnowledge">娣诲姞鐭ヨ瘑</a-button>
          </a-space>
        </template>
        <a-table :columns="columns" :data="knowledgeList" row-key="id">
          <template #actions="{ record }">
            <a-space>
              <a-button size="small" @click="editKnowledge(record)">缂栬緫</a-button>
              <a-button size="small" type="primary" danger @click="deleteKnowledge(record)">鍒犻櫎</a-button>
            </a-space>
          </template>
        </a-table>
      </a-card></div>
  </div>
</template>

<script>
export default {
  name: 'KnowledgeList',
  data() {
    return {
      columns: [
        { title: '鍒嗙被', dataIndex: 'category' },
        { title: '闂', dataIndex: 'question' },
        { title: '绛旀', dataIndex: 'answer' },
        { title: '鎿嶄綔', slotName: 'actions' }
      ],
      knowledgeList: []
    }
  },
  methods: {
    addKnowledge() { this.$message.info('娣诲姞鐭ヨ瘑') },
    editKnowledge(record) { this.$message.info('缂栬緫鐭ヨ瘑') },
    deleteKnowledge(record) { this.$message.warn('鍒犻櫎鐭ヨ瘑') }
  }
}
</script>
