<template>
  <Breadcrumb :items="['Home','Advanced','FeatureConfig','']" />
  <div class="page-container">
    <a-card class="general-card" title="功能配置">
      <template #extra>
        <a-button type="primary" @click="handleCreate"><icon-plus />新建分组</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="分组名称"><a-input v-model="form.group_name" placeholder="请输入分组名称" /></a-form-item>
          <a-form-item>
            <a-button type="primary" @click="loadGroups">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" />
      <a-modal v-model:visible="modalVisible" :title="modalTitle" @before-ok="handleSubmit">
        <a-form :model="form" layout="vertical">
          <a-form-item label="分组名称" required>
            <a-input v-model="form.group_name" placeholder="请输入分组名称" />
          </a-form-item>
          <a-form-item label="分组编码">
            <a-input v-model="form.group_code" placeholder="唯一编码，不填自动生成" />
          </a-form-item>
          <a-form-item label="图标">
            <a-select v-model="form.icon" placeholder="选择图标" allow-clear>
              <a-option value="icon-star">icon-star</a-option>
              <a-option value="icon-setting">icon-setting</a-option>
              <a-option value="icon-home">icon-home</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="描述">
            <a-textarea v-model="form.description" :rows="2" placeholder="分组描述" />
          </a-form-item>
        </a-form>
      </a-modal>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const data = ref<any[]>([])
const modalVisible = ref(false)
const modalTitle = ref('新建分组')
const form = ref<any>({ group_name: '', group_code: '', icon: '', description: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '分组名称', dataIndex: 'group_name', width: 180 },
  { title: '分组编码', dataIndex: 'group_code', width: 160 },
  { title: '图标', dataIndex: 'icon', width: 120 },
  { title: '状态', dataIndex: 'is_enabled', width: 80 },
  { title: '创建时间', dataIndex: 'created_at', width: 170 }
]

const pagination = ref({ current: 1, pageSize: 20, total: 0, showTotal: true })

async function loadGroups() {
  try {
    loading.value = true
    data.value = [{ id: 1, group_name: '基础功能', group_code: 'basic', icon: 'icon-star', is_enabled: '启用', created_at: '2026-03-28 10:00:00' }]
    pagination.value.total = data.value.length
  } catch (err: any) {
    Message.error('加载失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

function handleCreate() {
  modalTitle.value = '新建分组'
  form.value = { group_name: '', group_code: '', icon: '', description: '' }
  modalVisible.value = true
}

async function handleSubmit(done: (val: boolean) => void) {
  Message.success('保存成功')
  modalVisible.value = false
  loadGroups()
  done(true)
}

function handleReset() {
  form.value = { group_name: '', group_code: '' }
  loadGroups()
}

onMounted(() => { loadGroups() })
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>
