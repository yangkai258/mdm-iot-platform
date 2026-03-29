<template>
  <div class="company-list-page">
    <a-card class="general-card" title="数据权限配置">
      <template #extra>
        <a-button type="primary" @click="addRule"><icon-plus />新建</a-button>
      </template>
      <div class="search-bar">
        <a-input-search v-model="searchKey" placeholder="搜索..." style="width: 260px" @search="loadData" />
      </div>
      <a-tabs default-active-tab="dimension">
        <a-tab-pane key="dimension" title="权限维度">
          <a-row :gutter="[16, 16]">
            <a-col v-for="dim in dimensions" :key="dim.key" :span="6">
              <a-card hoverable>
                <a-space direction="vertical" fill>
                  <a-checkbox v-model="dim.enabled">{{ dim.label }}</a-checkbox>
                  <a-typography-text type="secondary">{{ dim.description }}</a-typography-text>
                </a-space>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'

const searchKey = ref('')
const loading = ref(false)

const dimensions = ref([
  { key: 'company', label: '公司维度', description: '按公司隔离数据', enabled: true },
  { key: 'department', label: '部门维度', description: '按部门隔离数据', enabled: true },
  { key: 'role', label: '角色维度', description: '按角色隔离数据', enabled: false },
])

const loadData = async () => {
  loading.value = true
  await new Promise(r => setTimeout(r, 300))
  loading.value = false
}

const addRule = () => {}

onMounted(() => loadData())
</script>

<style scoped>
.company-list-page { padding: 16px; }
.search-bar { margin-bottom: 16px; }
</style>
