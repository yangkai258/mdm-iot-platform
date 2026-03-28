<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-history /> 影子版本历史</a-space>
      </template>
      <template #extra>
        <a-space>
          <a-button @click="handleExport">
            <template #icon><icon-download /></template>
            导出快照
          </a-button>
          <a-button type="primary" @click="handleCompare">
            <template #icon><icon-swap /></template>
            对比版本
          </a-button>
        </a-space>
      </template>

      <a-row :gutter="16">
        <a-col :span="8">
          <a-card title="版本时间轴" size="small">
            <a-timeline>
              <a-timeline-item v-for="v in versions" :key="v.version" :color="v.current ? 'green' : 'gray'">
                <template #label>{{ v.version }}</template>
                <a-space direction="vertical" size="mini">
                  <span class="version-time">{{ v.createdAt }}</span>
                  <span class="version-user">{{ v.createdBy }}</span>
                  <a-tag v-if="v.current" color="green" size="small">当前版本</a-tag>
                </a-space>
              </a-timeline-item>
            </a-timeline>
          </a-card>
        </a-col>
        <a-col :span="16">
          <a-card title="版本详情">
            <a-form :model="selectedVersion" layout="vertical">
              <a-form-item label="版本号">{{ selectedVersion.version }}</a-form-item>
              <a-form-item label="创建时间">{{ selectedVersion.createdAt }}</a-form-item>
              <a-form-item label="创建方式">{{ selectedVersion.createdBy }}</a-form-item>
              <a-form-item label="影子状态 (JSON)">
                <a-code-block :code="selectedVersion.state" language="json" />
              </a-form-item>
            </a-form>
          </a-card>

          <a-card title="与上一版本差异" style="margin-top: 16px">
            <a-diff :old-value="selectedVersion.diffFromPrevious" language="json" />
          </a-card>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const versions = ref([
  { version: 'v12', createdAt: '2026-03-28 10:00:00', createdBy: 'auto', current: true },
  { version: 'v11', createdAt: '2026-03-27 15:30:00', createdBy: 'auto', current: false },
  { version: 'v10', createdAt: '2026-03-27 09:00:00', createdBy: 'admin', current: false },
  { version: 'v9', createdAt: '2026-03-26 18:00:00', createdBy: 'auto', current: false }
])
const selectedVersion = reactive({
  version: 'v12',
  createdAt: '2026-03-28 10:00:00',
  createdBy: 'auto',
  state: JSON.stringify({ volume: 80, mode: 'active', DND: false, brightness: 75 }, null, 2),
  diffFromPrevious: JSON.stringify({ volume: { old: 75, new: 80 } }, null, 2)
})

const handleExport = () => { }
const handleCompare = () => { }
</script>

<style scoped>
.container { padding: 16px; }
.version-time { font-size: 12px; color: #909399; }
.version-user { font-size: 11px; color: #c0c4cc; }
</style>
