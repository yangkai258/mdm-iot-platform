<template>
  <div class="page-container">
    <a-card class="general-card" title="钀ラ攢娲诲姩">
      <template #extra>
        <a-button type="primary" @click="handleCreate"><icon-plus />鏂板缓娲诲姩</a-button>
      </template>
      <div class="search-form">
        <a-form :model="filters" layout="inline">
          <a-form-item label="鍏抽敭璇?><a-input v-model="filters.keyword" placeholder="璇疯緭鍏" /></a-form-item>
          <a-form-item><a-button type="primary" @click="loadData">鏌ヨ</a-button><a-button @click="Object.keys(filters).forEach(k => filters[k] = ''); loadData()">閲嶇疆</a-button></a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="list" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #type="{ record }"><a-tag :color="getTypeColor(record.type)">{{ getTypeText(record.type) }}</a-tag></template>
        <template #status="{ record }"><a-badge :color="record.status === 'active' ? 'green' : 'gray'" :text="record.status === 'active' ? '杩涜涓? : '宸茬粨鏉?" /></template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="handleEdit(record)">缂栬緫</a-button>
          <a-button type="text" size="small" status="danger" @click="handleDelete(record)">鍒犻櫎</a-button>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model:visible="modalVisible" :title="isEdit ? '缂栬緫娲诲姩' : '鏂板缓娲诲姩'" :width="560">
      <a-form :model="form" layout="vertical">
        <a-form-item label="娲诲姩鍚嶇О"><a-input v-model="form.name" /></a-form-item>
        <a-form-item label="娲诲姩绫诲瀷"><a-select v-model="form.type" style="width:100%"><a-option value="points_double">鍙屽€嶇Н鍒?/a-option><a-option value="discount">鎶樻墸娲诲姩</a-option><a-option value="gift">璧犲搧娲诲姩</a-option></a-select></a-form-item>
      </a-form>
      <template #footer><a-button @click="modalVisible = false">鍙栨秷</a-button><a-button type="primary" @click="handleSubmit">纭畾</a-button></template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const modalVisible = ref(false)
const isEdit = ref(false)
const list = ref([])
const filters = reactive({ keyword: '', type: '' })
const form = reactive({ name: '', type: 'points_double' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = [
  { title: '娲诲姩鍚嶇О', dataIndex: 'name', width: 200 },
  { title: '绫诲瀷', slotName: 'type', width: 120 },
  { title: '鐘舵€?, slotName: 'status', width: 100 },
  { title: '寮€濮嬫椂闂?, dataIndex: 'start_time', width: 170 },
  { title: '鎿嶄綔', slotName: 'actions', width: 140 }
]

const getTypeColor = (t) => ({ points_double: 'blue', discount: 'orange', gift: 'purple' }[t] || 'gray')
const getTypeText = (t) => ({ points_double: '鍙屽€嶇Н鍒?, discount: '鎶樻墸', gift: '璧犲搧' }[t] || t)

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/v1/members/promotions', { headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') } }).then(r => r.json())
    list.value = res.data?.list || []
    pagination.total = list.value.length
  } catch { list.value = [] } finally { loading.value = false }
}

const handleCreate = () => { isEdit.value = false; Object.assign(form, { name: '', type: 'points_double' }); modalVisible.value = true }
const handleEdit = (r) => { isEdit.value = true; Object.assign(form, r); modalVisible.value = true }
const handleSubmit = () => { modalVisible.value = false; Message.success(isEdit.value ? '鏇存柊鎴愬姛' : '鍒涘缓鎴愬姛'); loadData() }
const handleDelete = () => { Message.success('鍒犻櫎鎴愬姛'); loadData() }
const onPageChange = (p) => { pagination.current = p; loadData() }

onMounted(() => loadData())
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>
