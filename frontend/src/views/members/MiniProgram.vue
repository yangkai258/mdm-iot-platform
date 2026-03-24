<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="名称"><a-input v-model="form.name" placeholder="请输入" /></a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleCreate">新建</a-button>
    </div>
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
      <template #actions="{ record }">
        <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
        <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
      </template>
    </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @before-ok="handleSubmit" @cancel="modalVisible = false">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="名称"><a-input v-model="form.name" placeholder="请输入" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>

import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'
const previewVisible = ref(false)

const form = reactive({
  miniprogram_name: '',
  app_id: '',
  app_secret: '',
  logo_url: '',
  card_bg: '',
  card_activate_url: '',
  show_points: '1',
  show_level: '1',
  show_coupons: '1',
  subscribe_reply: '',
  birthday_wish: ''
})

const loadData = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/miniprogram`, { headers: { 'Authorization': `Bearer ${token}` } })
    const data = await res.json()
    if (data.code === 0 && data.data) Object.assign(form, data.data)
  } catch (e) {}
}

const handleSave = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/miniprogram`, {
      method: 'PUT',
      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(form)
    })
    const data = await res.json()
    if (data.code === 0) Message.success('保存成功')
    else Message.error(data.message || '保存失败')
  } catch (e) { Message.error('保存失败') }
}

const handleReset = () => {
  Object.assign(form, {
    miniprogram_name: '', app_id: '', app_secret: '', logo_url: '',
    card_bg: '', card_activate_url: '', show_points: '1', show_level: '1',
    show_coupons: '1', subscribe_reply: '', birthday_wish: ''
  })
  Message.info('已重置')
}

const previewLogo = () => { previewVisible.value = true }

onMounted(() => loadData())

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
