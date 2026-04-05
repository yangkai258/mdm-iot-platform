<template>
  <Breadcrumb :items="['Home','Members','Tempcoupons','']" />
  <div class="page-container">
    <a-card class="general-card" title="T e m p c o u p o n s">
      <template #extra>
        <a-button type="primary" @click="handleCreate"><icon-plus />新建</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="关键词"><a-input v-model="form.keyword" placeholder="请输入" /></a-form-item>
          <a-form-item><a-button type="primary" @click="loadData">查询</a-button><a-button @click="handleReset">重置</a-button></a-form-item>
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
const form = ref<any>({ keyword: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '名称', dataIndex: 'name', width: 160 },
  { title: '状态', dataIndex: 'status', width: 90 },
  { title: '创建时间', dataIndex: 'created_at', width: 170 }
]

const pagination = ref({ current: 1, pageSize: 20, total: 0, showTotal: true })

async function loadData() {
  try {
    loading.value = true
    data.value = []
    pagination.value.total = 0
  } catch (err: any) {
    Message.error('加载失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

function handleCreate() {}
function handleReset() { form.value = { keyword: '' }; loadData() }
onMounted(() => { loadData() })
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>