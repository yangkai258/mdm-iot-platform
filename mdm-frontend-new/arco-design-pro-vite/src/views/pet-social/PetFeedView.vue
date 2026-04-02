<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="宠物喂养记录">
      <template #extra>
        <a-button @click="loadData"><icon-refresh />刷新</a-button>
      </template>
      <a-row :gutter="16">
        <a-col :span="8">
          <a-form-item label="关键词"><a-input v-model="form.keyword" placeholder="请输入" @pressEnter="loadData" /></a-form-item>
        </a-col>
        <a-col :flex="'86px'" style="display: flex; align-items: flex-end">
          <a-space direction="vertical" :size="8">
            <a-button type="primary" @click="loadData">查询</a-button>
            <a-button @click="Object.keys(form).forEach(k => form[k] = ''); loadData()">重置</a-button>
          </a-space>
        </a-col>
      </a-row>
      <a-divider style="margin: 0 0 16px 0" />
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id" />
    </a-table>
  </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const loading = ref(false)
const form = reactive({ keyword: '' })
const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = [
  { title: '宠物名称', dataIndex: 'pet_name', width: 140 },
  { title: '喂养时间', dataIndex: 'feed_time', width: 170 },
  { title: '食物类型', dataIndex: 'food_type', width: 120 },
  { title: '数量', dataIndex: 'quantity', width: 100 },
  { title: '备注', dataIndex: 'remark', ellipsis: true }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/v1/pets/feed-records', { headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') } }).then(r => r.json())
    data.value = res.data?.list || []
    pagination.total = data.value.length
  } catch { data.value = [] } finally { loading.value = false }
}
const onPageChange = (page) => { pagination.current = page; loadData() }
onMounted(() => loadData())
</script>

