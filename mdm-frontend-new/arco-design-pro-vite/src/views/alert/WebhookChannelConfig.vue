<template>
  <Breadcrumb :items="['Home','Alert','WebhookChannel','']" />
  <div class="page-container">
    <a-card class="general-card" title="Webhook 通知">
      <template #extra>
        <a-button type="primary" @click="handleCreate"><icon-plus />新建</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="Webhook名称"><a-input v-model="form.name" placeholder="名称" /></a-form-item>
          <a-form-item label="状态"><a-select v-model="form.enabled" placeholder="选择状态" allow-clear style="width: 120px"><a-option value="1">启用</a-option><a-option value="0">禁用</a-option></a-select></a-form-item>
          <a-form-item><a-button type="primary" @click="loadData">查询</a-button></a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" />
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const data = ref<any[]>([])
const form = ref<any>({ name: '', enabled: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '名称', dataIndex: 'name', width: 160 },
  { title: 'URL', dataIndex: 'url', ellipsis: true },
  { title: '加密密钥', dataIndex: 'secret', width: 180 },
  { title: '状态', dataIndex: 'enabled', width: 80 },
  { title: '创建时间', dataIndex: 'created_at', width: 170 }
]

const pagination = ref({ current: 1, pageSize: 20, total: 0, showTotal: true })

async function loadData() { loading.value = true; data.value = []; loading.value = false }
function handleCreate() {}
onMounted(() => { loadData() })
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>
