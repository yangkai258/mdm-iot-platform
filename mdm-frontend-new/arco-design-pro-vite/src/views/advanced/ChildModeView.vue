<template>
  <Breadcrumb :items="['Home','Advanced','ChildMode','']" />
  <div class="page-container">
    <a-card class="general-card" title="儿童模式">
      <template #extra>
        <a-button type="primary" @click="handleCreate"><icon-plus />新建规则</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="儿童账号">
            <a-select v-model="form.child_id" placeholder="选择儿童账号" allow-clear style="width: 200px">
              <a-option v-for="c in children" :key="c.id" :value="c.id">{{ c.name }}</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="模式状态">
            <a-select v-model="form.status" placeholder="选择状态" allow-clear style="width: 140px">
              <a-option value="enabled">已启用</a-option>
              <a-option value="disabled">已禁用</a-option>
            </a-select>
          </a-form-item>
          <a-form-item>
            <a-button type="primary" @click="loadChildModes">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" />
      <a-modal v-model:visible="modalVisible" :title="modalTitle" @before-ok="handleSubmit">
        <a-form :model="form" layout="vertical">
          <a-form-item label="儿童账号" required>
            <a-select v-model="form.child_id" placeholder="选择儿童账号">
              <a-option v-for="c in children" :key="c.id" :value="c.id">{{ c.name }}</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="内容过滤">
            <a-switch v-model="form.content_filter" checked-text="开" unchecked-text="关" />
          </a-form-item>
          <a-form-item label="使用时长限制">
            <a-input-number v-model="form.time_limit" :min="0" :max="24" placeholder="小时/天" />
          </a-form-item>
          <a-form-item label="允许应用">
            <a-textarea v-model="form.allowed_apps" :rows="2" placeholder="逗号分隔应用名" />
          </a-form-item>
          <a-form-item label="备注">
            <a-textarea v-model="form.remark" :rows="2" placeholder="备注信息" />
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
const modalTitle = ref('新建儿童模式')
const form = ref<any>({
  child_id: '', content_filter: true, time_limit: 2, allowed_apps: '', remark: ''
})
const children = ref<any[]>([
  { id: 1, name: '小明' },
  { id: 2, name: '小红' }
])

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '儿童账号', dataIndex: 'child_name', width: 120 },
  { title: '内容过滤', dataIndex: 'content_filter', width: 100 },
  { title: '时长限制(h)', dataIndex: 'time_limit', width: 120 },
  { title: '状态', dataIndex: 'status', width: 90 },
  { title: '创建时间', dataIndex: 'created_at', width: 170 }
]

const pagination = ref({ current: 1, pageSize: 20, total: 0, showTotal: true })

async function loadChildModes() {
  try {
    loading.value = true
    data.value = [
      { id: 1, child_name: '小明', content_filter: '开启', time_limit: 2, status: 'enabled', created_at: '2026-03-28 10:00:00' }
    ]
    pagination.value.total = data.value.length
  } catch (err: any) {
    Message.error('加载失败')
  } finally {
    loading.value = false
  }
}

function handleCreate() {
  modalTitle.value = '新建儿童模式'
  form.value = { child_id: '', content_filter: true, time_limit: 2, allowed_apps: '', remark: '' }
  modalVisible.value = true
}

async function handleSubmit(done: (val: boolean) => void) {
  Message.success('保存成功')
  modalVisible.value = false
  loadChildModes()
  done(true)
}

function handleReset() {
  form.value = { child_id: '', status: '' }
  loadChildModes()
}

onMounted(() => { loadChildModes() })
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>
