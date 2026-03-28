<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-history /> 知识版本历史</a-space>
      </template>

      <a-row :gutter="16">
        <a-col :span="6">
          <a-card title="版本列表" size="small">
            <a-timeline>
              <a-timeline-item v-for="v in versions" :key="v.version" :color="v.current ? 'green' : 'gray'">
                <template #label>v{{ v.version }}</template>
                <a-space direction="vertical" size="mini">
                  <span>{{ v.createdAt }}</span>
                  <a-tag v-if="v.approved" color="green" size="small">已审核</a-tag>
                  <a-tag v-else color="orange" size="small">待审核</a-tag>
                </a-space>
              </a-timeline-item>
            </a-timeline>
          </a-card>
        </a-col>

        <a-col :span="18">
          <a-card title="版本详情">
            <a-descriptions :column="2" bordered>
              <a-descriptions-item label="版本号">v{{ selectedVersion.version }}</a-descriptions-item>
              <a-descriptions-item label="状态">
                <a-tag :color="selectedVersion.approved ? 'green' : 'orange'">
                  {{ selectedVersion.approved ? '已审核' : '待审核' }}
                </a-tag>
              </a-descriptions-item>
              <a-descriptions-item label="创建时间">{{ selectedVersion.createdAt }}</a-descriptions-item>
              <a-descriptions-item label="创建人">{{ selectedVersion.createdBy }}</a-descriptions-item>
              <a-descriptions-item label="变更说明" :span="2>{{ selectedVersion.changeSummary }}</a-descriptions-item>
            </a-descriptions>

            <a-divider>内容预览</a-divider>
            <a-card size="small">
              <div v-html="selectedVersion.content"></div>
            </a-card>
          </a-card>

          <a-space style="margin-top: 16px">
            <a-button v-if="!selectedVersion.approved" type="primary" @click="handleApprove">审核通过</a-button>
            <a-button v-if="!selectedVersion.approved" @click="handleReject">拒绝</a-button>
            <a-button @click="handleRestore">恢复此版本</a-button>
            <a-button @click="handleCompare">版本对比</a-button>
          </a-space>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const versions = ref([
  { version: '3', createdAt: '2026-03-28 10:00', createdBy: 'admin', approved: true, changeSummary: '更新内容', content: '知识内容...' },
  { version: '2', createdAt: '2026-03-27 15:00', createdBy: 'editor', approved: true, changeSummary: '修正错误', content: '知识内容...' },
  { version: '1', createdAt: '2026-03-26 09:00', createdBy: 'editor', approved: true, changeSummary: '初始版本', content: '知识内容...' }
])
const selectedVersion = ref(versions.value[0])

const handleApprove = () => { }
const handleReject = () => { }
const handleRestore = () => { }
const handleCompare = () => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
