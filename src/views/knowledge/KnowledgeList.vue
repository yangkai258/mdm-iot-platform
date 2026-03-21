<template>
  <div class="knowledge-list">
    <div class="page-content">
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
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Message } from '@arco-design/web-vue'

const columns = [
  { title: '分类', dataIndex: 'category' },
  { title: '问题', dataIndex: 'question' },
  { title: '答案', dataIndex: 'answer' },
  { title: '操作', slotName: 'actions' }
]

const knowledgeList = ref([])

const addKnowledge = () => { Message.info('添加知识') }
const editKnowledge = (record) => { Message.info('编辑知识') }
const deleteKnowledge = (record) => { Message.warn('删除知识') }
</script>

<style scoped>
.knowledge-list {
  padding: 16px;
}
</style>
