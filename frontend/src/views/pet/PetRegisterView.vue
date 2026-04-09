<template>
  <div class="page-container">
    <a-card class="general-card" title="宠物登记">
      <template #extra>
        <a-button @click="loadData"><icon-refresh />刷新</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="宠物名称"><a-input v-model="form.pet_name" placeholder="请输�? /></a-form-item>
          <a-form-item label="宠物类型">
            <a-select v-model="form.pet_type" placeholder="选择类型" allow-clear style="width: 120px">
              <a-option value="dog">�?/a-option>
              <a-option value="cat">�?/a-option>
              <a-option value="bird">�?/a-option>
              <a-option value="other">其他</a-option>
            </a-select>
          </a-form-item>
          <a-form-item><a-button type="primary" @click="handleSearch">查询</a-button><a-button @click="Object.keys(form).forEach(k => form[k] = ''); loadData()">重置</a-button></a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #pet_photo="{ record }">
          <a-avatar v-if="record.pet_photo" :src="record.pet_photo" size="large" />
          <a-avatar v-else size="large">{{ record.pet_name?.[0] || '?' }}</a-avatar>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 'active' ? 'green' : 'gray'">{{ record.status === 'active' ? '已激�? : '未激�? }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="viewDetail(record)">详情</a-button>
          <a-button type="text" size="small" @click="editPet(record)">编辑</a-button>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑宠物' : '宠物登记'" @before-ok="handleSubmit" :loading="submitting" :width="640">
      <a-form :model="petForm" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="宠物名称" required><a-input v-model="petForm.pet_name" placeholder="请输入宠物名�? /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="宠物类型" required>
              <a-select v-model="petForm.pet_type" placeholder="选择类型">
                <a-option value="dog">�?/a-option>
                <a-option value="cat">�?/a-option>
                <a-option value="bird">�?/a-option>
                <a-option value="other">其他</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="品种"><a-input v-model="petForm.breed" placeholder="请输入品�? /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="性别">
              <a-select v-model="petForm.gender" placeholder="选择性别">
                <a-option value="male">�?/a-option>
                <a-option value="female">�?/a-option>
                <a-option value="unknown">未知</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="出生日期"><a-date-picker v-model="petForm.birth_date" style="width: 100%" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="体重(kg)"><a-input-number v-model="petForm.weight" :min="0" :max="500" style="width: 100%" /></a-form-item>
          </a-col>
        </a-row>
        <a-divider>主人信息</a-divider>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="主人姓名" required><a-input v-model="petForm.owner_name" placeholder="请输入主人姓�? /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="联系电话" required><a-input v-model="petForm.owner_phone" placeholder="请输入联系电�? /></a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="联系地址"><a-input v-model="petForm.owner_address" placeholder="请输入联系地址" /></a-form-item>
        <a-form-item label="宠物照片">
          <a-upload :custom-request="handleUpload" :show-file-list="false">
            <a-button v-if="!petForm.pet_photo"><icon-upload />上传照片</a-button>
            <a-image v-else :src="petForm.pet_photo" width="120" height="120" :preview="true" />
          </a-upload>
        </a-form-item>
        <a-form-item label="备注"><a-textarea v-model="petForm.notes" :rows="2" placeholder="备注信息" /></a-form-item>
      </a-form>
    </a-modal>
    <a-modal v-model:visible="detailVisible" title="宠物详情" :footer="null" :width="560">
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="宠物名称" :span="2">{{ currentPet?.pet_name }}</a-descriptions-item>
        <a-descriptions-item label="宠物类型">{{ currentPet?.pet_type }}</a-descriptions-item>
        <a-descriptions-item label="品种">{{ currentPet?.breed }}</a-descriptions-item>
        <a-descriptions-item label="性别">{{ currentPet?.gender }}</a-descriptions-item>
        <a-descriptions-item label="主人">{{ currentPet?.owner_name }}</a-descriptions-item>
        <a-descriptions-item label="电话">{{ currentPet?.owner_phone }}</a-descriptions-item>
        <a-descriptions-item label="地址" :span="2">{{ currentPet?.owner_address }}</a-descriptions-item>
        <a-descriptions-item label="登记时间" :span="2">{{ currentPet?.created_at }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const submitting = ref(false)
const data = ref([])
const formVisible = ref(false)
const detailVisible = ref(false)
const isEdit = ref(false)
const currentPet = ref(null)
const form = reactive({ pet_name: '', pet_type: '' })
const petForm = reactive({ id: null, pet_name: '', pet_type: '', breed: '', gender: '', birth_date: null, weight: null, owner_name: '', owner_phone: '', owner_address: '', pet_photo: '', notes: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: '照片', slotName: 'pet_photo', width: 80 },
  { title: '宠物名称', dataIndex: 'pet_name', width: 120 },
  { title: '类型', dataIndex: 'pet_type', width: 80 },
  { title: '品种', dataIndex: 'breed', width: 100 },
  { title: '主人', dataIndex: 'owner_name', width: 100 },
  { title: '电话', dataIndex: 'owner_phone', width: 130 },
  { title: '状�?, slotName: 'status', width: 80 },
  { title: '登记时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 120 }
]

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (form.pet_name) params.append('pet_name', form.pet_name)
    if (form.pet_type) params.append('pet_type', form.pet_type)
    const res = await fetch(`/api/v1/pet/register?${params}`, { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { data.value = res.data?.list || []; pagination.total = res.data?.total || 0 }
    else { data.value = [] }
  } catch (e) { Message.error('加载失败') } finally { loading.value = false }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }

const viewDetail = (record) => { currentPet.value = record; detailVisible.value = true }

const editPet = (record) => { isEdit.value = true; Object.assign(petForm, record); formVisible.value = true }

const handleSubmit = async (done) => {
  if (!petForm.pet_name || !petForm.owner_name || !petForm.owner_phone) { Message.warning('请填写必填项'); done(false); return }
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const url = isEdit.value ? `/api/v1/pet/register/${petForm.id}` : '/api/v1/pet/register'
    const res = await fetch(url, { method: isEdit.value ? 'PUT' : 'POST', headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(petForm) }).then(r => r.json())
    if (res.code === 0) { Message.success(isEdit.value ? '更新成功' : '登记成功'); formVisible.value = false; loadData() }
    else { Message.error(res.message || '操作失败') }
    done(true)
  } catch (e) { Message.error('操作失败'); done(false) } finally { submitting.value = false }
}

const handleUpload = async ({ file, onSuccess, onError }) => {
  try {
    const formData = new FormData()
    formData.append('file', file)
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/upload', { method: 'POST', headers: { 'Authorization': `Bearer ${token}` }, body: formData }).then(r => r.json())
    if (res.code === 0) { petForm.pet_photo = res.data?.url; onSuccess() }
    else onError()
  } catch (e) { onError() }
}

onMounted(() => loadData())
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>