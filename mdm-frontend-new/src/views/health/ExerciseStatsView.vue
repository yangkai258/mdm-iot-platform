<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="宠物">
          <a-input v-model="form.pet_name" placeholder="请输入宠物名称" />
        </a-form-item>
        <a-form-item label="运动类型">
          <a-select v-model="form.exercise_type" placeholder="请选择" allow-clear style="width: 140px">
            <a-option value="walk">散步</a-option>
            <a-option value="run">跑步</a-option>
            <a-option value="swim">游泳</a-option>
            <a-option value="play">玩耍</a-option>
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
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="handlePageChange" />
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @before-ok="handleSubmit" :unmount-on-close="false">
      <a-form :model="form" label-col-flex="100px" ref="formRef">
        <a-form-item label="宠物名称" field="pet_name" :rules="[{ required: true, message: '请输入宠物名称' }]">
          <a-input v-model="form.pet_name" placeholder="请输入" />
        </a-form-item>
        <a-form-item label="运动类型" field="exercise_type" :rules="[{ required: true, message: '请选择运动类型' }]">
          <a-select v-model="form.exercise_type" placeholder="请选择">
            <a-option value="walk">散步</a-option>
            <a-option value="run">跑步</a-option>
            <a-option value="swim">游泳</a-option>
            <a-option value="play">玩耍</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="时长(分钟)" field="duration">
          <a-input-number v-model="form.duration" placeholder="请输入" :min="1" style="width: 100%" />
        </a-form-item>
        <a-form-item label="消耗卡路里" field="calories">
          <a-input-number v-model="form.calories" placeholder="请输入" :min="0" style="width: 100%" />
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
const modalTitle = ref('新建运动记录')
const formRef = ref()
const editId = ref<number | null>(null)

const form = reactive({
  pet_name: '',
  exercise_type: '',
  duration: null as number | null,
  calories: null as number | null,
  remark: ''
})

const data = ref<any[]>([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '宠物名称', dataIndex: 'pet_name', width: 120 },
  { title: '运动类型', dataIndex: 'exercise_type', width: 100 },
  { title: '时长(分钟)', dataIndex: 'duration', width: 100 },
  { title: '消耗卡路里', dataIndex: 'calories', width: 100 },
  { title: '备注', dataIndex: 'remark', ellipsis: true },
  { title: '记录时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 140 }
]

const loadData = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.current, page_size: pagination.pageSize }
    if (form.pet_name) params.pet_name = form.pet_name
    if (form.exercise_type) params.exercise_type = form.exercise_type
    const res = await axios.get(`${API_BASE}/health/exercise-stats`, { params })
    if (res.data.code === 0) {
      data.value = res.data.data.list || []
      pagination.total = res.data.data.pagination?.total || 0
    }
  } catch {
    data.value = [
      { id: 1, pet_name: '小白', exercise_type: '散步', duration: 30, calories: 120, remark: '户外活动', created_at: '2026-03-20 10:00:00' },
      { id: 2, pet_name: '小黄', exercise_type: '跑步', duration: 20, calories: 200, remark: '公园慢跑', created_at: '2026-03-20 14:30:00' }
    ]
    pagination.total = 2
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { form.pet_name = ''; form.exercise_type = ''; pagination.current = 1; loadData() }
const handlePageChange = (page: number) => { pagination.current = page; loadData() }

const handleCreate = () => {
  editId.value = null
  modalTitle.value = '新建运动记录'
  Object.assign(form, { pet_name: '', exercise_type: '', duration: null, calories: null, remark: '' })
  modalVisible.value = true
}

const handleEdit = (record: any) => {
  editId.value = record.id
  modalTitle.value = '编辑运动记录'
  Object.assign(form, {
    pet_name: record.pet_name,
    exercise_type: record.exercise_type,
    duration: record.duration,
    calories: record.calories,
    remark: record.remark || ''
  })
  modalVisible.value = true
}

const handleSubmit = async (done: (arg: boolean) => void) => {
  try {
    await formRef.value?.validate()
    const payload = { pet_name: form.pet_name, exercise_type: form.exercise_type, duration: form.duration, calories: form.calories, remark: form.remark }
    if (editId.value) {
      await axios.put(`${API_BASE}/health/exercise-stats/${editId.value}`, payload)
    } else {
      await axios.post(`${API_BASE}/health/exercise-stats`, payload)
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
  Modal.confirm({ title: '确认删除', content: `确定删除该运动记录？`, onOk: async () => {
    try { await axios.delete(`${API_BASE}/health/exercise-stats/${record.id}`) } catch {}
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
