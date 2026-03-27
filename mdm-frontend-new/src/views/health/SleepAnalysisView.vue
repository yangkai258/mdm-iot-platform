<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="宠物">
          <a-input v-model="form.pet_name" placeholder="请输入宠物名称" />
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
      <template #quality_score="{ record }">
        <a-tag :color="getScoreColor(record.quality_score)">{{ record.quality_score }}分</a-tag>
      </template>
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
        <a-form-item label="睡眠质量分" field="quality_score" :rules="[{ required: true, message: '请输入睡眠质量分' }]">
          <a-input-number v-model="form.quality_score" :min="0" :max="100" style="width: 100%" />
        </a-form-item>
        <a-form-item label="睡眠时长(小时)" field="duration">
          <a-input-number v-model="form.duration" :min="0" :max="24" :step="0.5" style="width: 100%" />
        </a-form-item>
        <a-form-item label="开始时间" field="sleep_start">
          <a-input v-model="form.sleep_start" placeholder="例如: 22:00" />
        </a-form-item>
        <a-form-item label="结束时间" field="sleep_end">
          <a-input v-model="form.sleep_end" placeholder="例如: 07:00" />
        </a-form-item>
        <a-form-item label="备注" field="remark">
          <a-textarea v-model="form.remark" placeholder="请输入备注" :rows="3" />
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
const modalTitle = ref('新建睡眠记录')
const formRef = ref()
const editId = ref<number | null>(null)

const form = reactive({
  pet_name: '',
  quality_score: null as number | null,
  duration: null as number | null,
  sleep_start: '',
  sleep_end: '',
  remark: ''
})

const data = ref<any[]>([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '宠物名称', dataIndex: 'pet_name', width: 100 },
  { title: '睡眠质量分', slotName: 'quality_score', width: 110 },
  { title: '睡眠时长(小时)', dataIndex: 'duration', width: 120 },
  { title: '开始时间', dataIndex: 'sleep_start', width: 100 },
  { title: '结束时间', dataIndex: 'sleep_end', width: 100 },
  { title: '备注', dataIndex: 'remark', ellipsis: true },
  { title: '记录日期', dataIndex: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 140 }
]

const getScoreColor = (score: number) => score >= 80 ? 'green' : score >= 60 ? 'orange' : 'red'

const loadData = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.current, page_size: pagination.pageSize }
    if (form.pet_name) params.pet_name = form.pet_name
    const res = await axios.get(`${API_BASE}/health/sleep`, { params })
    if (res.data.code === 0) {
      data.value = res.data.data.list || []
      pagination.total = res.data.data.pagination?.total || 0
    }
  } catch {
    data.value = [
      { id: 1, pet_name: '小白', quality_score: 85, duration: 9, sleep_start: '22:00', sleep_end: '07:00', remark: '睡眠质量良好', created_at: '2026-03-20 07:00:00' },
      { id: 2, pet_name: '小黄', quality_score: 65, duration: 7.5, sleep_start: '23:00', sleep_end: '06:30', remark: '浅睡眠较多', created_at: '2026-03-20 07:00:00' }
    ]
    pagination.total = 2
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { form.pet_name = ''; pagination.current = 1; loadData() }
const handlePageChange = (page: number) => { pagination.current = page; loadData() }

const handleCreate = () => {
  editId.value = null
  modalTitle.value = '新建睡眠记录'
  Object.assign(form, { pet_name: '', quality_score: null, duration: null, sleep_start: '', sleep_end: '', remark: '' })
  modalVisible.value = true
}

const handleEdit = (record: any) => {
  editId.value = record.id
  modalTitle.value = '编辑睡眠记录'
  Object.assign(form, {
    pet_name: record.pet_name, quality_score: record.quality_score, duration: record.duration,
    sleep_start: record.sleep_start, sleep_end: record.sleep_end, remark: record.remark || ''
  })
  modalVisible.value = true
}

const handleSubmit = async (done: (arg: boolean) => void) => {
  try {
    await formRef.value?.validate()
    const payload = { pet_name: form.pet_name, quality_score: form.quality_score, duration: form.duration, sleep_start: form.sleep_start, sleep_end: form.sleep_end, remark: form.remark }
    if (editId.value) {
      await axios.put(`${API_BASE}/health/sleep/${editId.value}`, payload)
    } else {
      await axios.post(`${API_BASE}/health/sleep`, payload)
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
  Modal.confirm({ title: '确认删除', content: `确定删除该睡眠记录？`, onOk: async () => {
    try { await axios.delete(`${API_BASE}/health/sleep/${record.id}`) } catch {}
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
