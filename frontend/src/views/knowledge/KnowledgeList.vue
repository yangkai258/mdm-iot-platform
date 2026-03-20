<template>
  <div class="knowledge-list">
    <a-layout-header class="header">
      <div class="header-title">知识库</div>
    </a-layout-header>
    <a-layout-content class="content">
      <a-card>
        <template #title>
          <a-space>
            <span>知识库</span>
            <a-button type="primary" @click="addKnowledge">添加知识</a-button>
          </a-space>
        </template>
        <a-table :columns="columns" :data="knowledgeList" row-key="id">
          <template #actions="{ record }">
            <a-space>
              <a-button size="small" @click="editKnowledge(record)">编辑</a-button>
              <a-button size="small" type="primary" danger @click="deleteKnowledge(record)">删除</a-button>
            </a-space>
          </template>
        </a-table>
      </a-card>
    </a-layout-content>
  </div>
</template>

<script>
export default {
  name: 'KnowledgeList',
  data() {
    return {
      columns: [
        { title: '分类', dataIndex: 'category' },
        { title: '问题', dataIndex: 'question' },
        { title: '答案', dataIndex: 'answer' },
        { title: '操作', slotName: 'actions' }
      ],
      knowledgeList: []
    }
  },
  methods: {
    addKnowledge() { this.$message.info('添加知识') },
    editKnowledge(record) { this.$message.info('编辑知识') },
    deleteKnowledge(record) { this.$message.warn('删除知识') }
  }
}
</script>
