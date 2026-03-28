<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-history /> 内容版本历史</a-space>
      </template>

      <a-row :gutter="16">
        <a-col :span="6">
          <a-card title="版本时间轴" size="small">
            <a-timeline>
              <a-timeline-item v-for="v in versions" :key="v.version" :color="v.current ? 'green' : 'gray'">
                <template #label>v{{ v.version }}</template>
                <a-space direction="vertical" size="mini">
                  <span class="time">{{ v.createdAt }}</span>
                  <span class="user">{{ v.createdBy }}</span>
                  <a-tag v-if="v.current" color="green" size="small">当前</a-tag>
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
                <a-tag :color="selectedVersion.current ? 'green' : 'gray'">{{ selectedVersion.current ? '当前版本' : '历史版本' }}</a-tag>
              </a-descriptions-item>
              <a-descriptions-item label="创建时间">{{ selectedVersion.createdAt }}</a-descriptions-item>
              <a-descriptions-item label="创建人">{{ selectedVersion.createdBy }}</a-descriptions-item>
              <a-descriptions-item label="变更说明" :span="2>{{ selectedVersion.changeSummary }}</a-descriptions-item>
            </a-descriptions>

            <a-divider>文件信息</a-divider>
            <a-table :columns="fileColumns" :data="selectedVersion.files" size="small" :pagination="false">
              <template #size="{ record }">{{ (record.size / 1024).toFixed(2) }} KB</template>
            </a-table>
          </a-card>

          <a-card title="版本对比" style="margin-top: 16px">
            <template #extra>
              <a-space>
                <a-select v-model="compareVersion" placeholder="选择对比版本" style="width: 150px">
                  <a-option v-for="v in versions" :key="v.version" :value="v.version">{{ v.version }}</a-option>
                </a-select>
                <a-button @click="handleCompare">对比</a-button>
              </a-space>
            </template>
            <a-diff :old-value="diffOld" language="json" />
          </a-card>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const versions = ref([
  { version: '3', createdAt: '2026-03-28 10:00', createdBy: 'admin', current: true, changeSummary: '更新UI组件', files: [{ name: 'index.vue', size: 12340 }] },
  { version: '2', createdAt: '2026-03-27 15:00', createdBy: 'dev', current: false, changeSummary: '修复bug', files: [{ name: 'index.vue', size: 12000 }] },
  { version: '1', createdAt: '2026-03-26 09:00', createdBy: 'dev', current: false, changeSummary: '初始版本', files: [{ name: 'index.vue', size: 10000 }] }
])
const selectedVersion = ref(versions.value[0])
const compareVersion = ref('')
const diffOld = ref('{}')

const fileColumns = [
  { title: '文件名', dataIndex: 'name' },
  { title: '大小', slotName: 'size' },
  { title: '哈希', dataIndex: 'hash' }
]

const handleCompare = () => { }
</script>

<style scoped>
.container { padding: 16px; }
.time { font-size: 12px; color: #909399; }
.user { font-size: 11px; color: #c0c4cc; }
</style>
