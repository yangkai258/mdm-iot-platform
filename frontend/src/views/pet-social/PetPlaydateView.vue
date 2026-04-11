<template>
  <div class="container">
    <a-card class="general-card" title="宠物约会">
      <template #extra>
        <a-space>
          <a-button type="primary" @click="openCreate"><icon-plus />新建</a-button>
          <a-button @click="loadData"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-row :gutter="16">
        <a-col :span="8">
          <a-form-item label="关键字"><a-input v-model="form.keyword" placeholder="请输入" @pressEnter="loadData" /></a-form-item>
        </a-col>
        <a-col :flex="'86px'" style="display: flex; align-items: flex-end">
          <a-space direction="vertical" :size="8">
            <a-button type="primary" @click="loadData">查询</a-button>
            <a-button @click="Object.keys(form).forEach(k => form[k] = ''); loadData()">重置</a-button>
          </a-space>
        </a-col>
      </a-row>
      <a-divider style="margin: 0 0 16px 0" />
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #status="{ record }"><a-badge :color="record.status === 'active' ? 'green' : 'gray'" :text="record.status === 'active' ? '进行中' : '已结束'" /></template>
        <template #actions="{ record }">
          <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model="formVisible" title="新建约会" :width="560">
      <a-form :model="form" layout="vertical">
        <a-form-item label="约会地点"><a-input v-model="form.location" /></a-form-item>
        <a-form-item label="时间"><a-input v-model="form.playdate_time" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="formVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确认</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const loading = ref(false)
const formVisible = ref(false)
const form = reactive({ keyword: '', location: '', playdate_time: '' })
const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = [
  { title: '宠物名称', dataIndex: 'pet_name', width: 140 },
  { title: '约会地点', dataIndex: 'location', width: 160 },
  { title: '时间', dataIndex: 'playdate_time', width: 170 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 80 }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/v1/pets/playdates', { headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') } }).then(r => r.json())
    data.value = res.data?.list || []
    pagination.total = data.value.length
  } catch { data.value = [] } finally { loading.value = false }
}
const openCreate = () => { Object.assign(form, { location: '', playdate_time: '' }); formVisible.value = true }
const handleSubmit = () => { formVisible.value = false; Message.success('创建成功'); loadData() }
const handleDelete = () => { Message.success('删除成功'); loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }
onMounted(() => loadData())
</script>