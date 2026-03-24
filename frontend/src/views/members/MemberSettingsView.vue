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
import * as api from '@/api/member'

const saving = ref(false)
const originalForm = {}

const form = reactive({
  registerPoints: 0,
  pointsExpireDays: 0,
  pointsDiscountRatio: 0.01,
  pointsToMoneyRatio: 100,
  memberValidMonths: 0,
  expireReminderDays: 7,
  allowSelfCancel: false,
  requireMobileVerify: true,
  minPointsBalance: 0,
  maxPointsPerOrder: 0,
  remark: ''
})

const loadSettings = async () => {
  try {
    const res = await api.getMemberSettings()
    const data = res.data || {}
    Object.assign(form, data)
    Object.assign(originalForm, data)
  } catch (err) {
    // ignore, use defaults
  }
}

const handleSave = async () => {
  saving.value = true
  try {
    await api.updateMemberSettings({ ...form })
    Message.success('保存成功')
    Object.assign(originalForm, form)
  } catch (err) {
    Message.error(err.message || '保存失败')
  } finally {
    saving.value = false
  }
}

const handleReset = () => {
  Object.assign(form, originalForm)
  Message.info('已重置为保存的值')
}

onMounted(() => {
  loadSettings()
})

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
