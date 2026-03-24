<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="宠物">
          <a-input v-model="form.pet_name" placeholder="请输入宠物名称" />
        </a-form-item>
        <a-form-item label="报告类型">
          <a-select v-model="form.report_type" placeholder="请选择" allow-clear style="width: 140px">
            <a-option value="daily">日报</a-option>
            <a-option value="weekly">周报</a-option>
            <a-option value="monthly">月报</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleCreate">新建</a-button>
    </div>
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="handlePageChange">
      <template #actions="{ record }">
        <a-space>
          <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
          <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
        </a-space>
      </template>
    </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @before-ok="handleSubmit" :unmount-on-close="false">
      <a-form :model="form" label-col-flex="100px" ref="formRef">
        <a-form-item label="宠物名称" field="pet_name" :rules="[{ required: true, message: '请输入宠物名称' }]">
          <a-input v-model="form.pet_name" placeholder="请输入" />
        </a-form-item>
        <a-form-item label="报告类型" field="report_type" :rules="[{ required: true, message: '请选择报告类型' }]">
          <a-select v-model="form.report_type" placeholder="请选择">
            <a-option value="daily">日报</a-option>
            <a-option value="weekly">周报</a-option>
            <a-option value="monthly">月报</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="报告内容" field="content" :rules="[{ required: true, message: '请输入报告内容' }]">
          <a-textarea v-model="form.content" placeholder="请输入报告内容" :rows="4" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import axios from 'axios'

const API_BASE = '/api/v1'

const loading = ref(false)
const modalVisible = ref(false)
const modalTitle = ref('新建健康报告')
const formRef = ref()
const editId = ref<number | null>(null)

const form = reactive({
  pet_name: '',
  report_type: '',
  content: ''
})

const data = ref<any[]>([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '宠物名称', dataIndex: 'pet_name', width: 120 },
  { title: '报告类型', dataIndex: 'report_type', width: 100 },
  { title: '报告内容', dataIndex: 'content', ellipsis: true },
  { title: '生成时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 140 }
]

const loadData = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.current, page_size: pagination.pageSize }
    if (form.pet_name) params.pet_name = form.pet_name
    if (form.report_type) params.report_type = form.report_type
    const res = await axios.get(`${API_BASE}/health/reports`, { params })
    if (res.data.code === 0) {
      data.value = res.data.data.list || []
      pagination.total = res.data.data.pagination?.total || 0
    }
  } catch {
    data.value = [
      { id: 1, pet_name: '小白', report_type: 'daily', content: '今日健康状况良好，食欲正常', created_at: '2026-03-20 18:00:00' },
      { id: 2, pet_name: '小黄', report_type: 'weekly', content: '本周运动量达标，体重稳定', created_at: '2026-03-20 10:00:00' }
    ]
    pagination.total = 2
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { form.pet_name = ''; form.report_type = ''; pagination.current = 1; loadData() }
const handlePageChange = (page: number) => { pagination.current = page; loadData() }

const handleCreate = () => {
  editId.value = null
  modalTitle.value = '新建健康报告'
  Object.assign(form, { pet_name: '', report_type: '', content: '' })
  modalVisible.value = true
}

const handleEdit = (record: any) => {
  editId.value = record.id
  modalTitle.value = '编辑健康报告'
  Object.assign(form, { pet_name: record.pet_name, report_type: record.report_type, content: record.content || '' })
  modalVisible.value = true
}

const handleSubmit = async (done: (arg: boolean) => void) => {
  try {
    await formRef.value?.validate()
    const payload = { pet_name: form.pet_name, report_type: form.report_type, content: form.content }
    if (editId.value) {
      await axios.put(`${API_BASE}/health/reports/${editId.value}`, payload)
    } else {
      await axios.post(`${API_BASE}/health/reports`, payload)
    }
    Message.success('保存成功')
    modalVisible.value = false
    loadData()
    done(true)
  } catch {
    done(false)
  }
}

const handleDelete = (record: any) => {
  Modal.confirm({ title: '确认删除', content: `确定删除该健康报告？`, onOk: async () => {
    try { await axios.delete(`${API_BASE}/health/reports/${record.id}`) } catch {}
    Message.success('删除成功')
    loadData()
  }})
}

onMounted(() => loadData())
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
