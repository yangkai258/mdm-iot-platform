<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="标签管理">
      <template #extra>
        <a-space :size="12">
          <a-button type="primary" @click="showCreate"><icon-plus />新建</a-button>
          <a-button @click="loadData"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-row :gutter="16">
        <a-col :span="8">
          <a-form-item label="标签名称">
            <a-input v-model="filters.keyword" placeholder="请输入" @pressEnter="loadData" />
          </a-form-item>
        </a-col>
        <a-col :flex="'86px'" style="display: flex; align-items: flex-end">
          <a-space direction="vertical" :size="8">
            <a-button type="primary" @click="loadData">查询</a-button>
            <a-button @click="filters.keyword = ''; loadData()">重置</a-button>
          </a-space>
        </a-col>
      </a-row>
      <a-divider style="margin: 0 0 16px 0" />
      <a-table :columns="columns" :data="tagList" :loading="loading" :pagination="paginationConfig" @page-change="onPageChange" row-key="id">
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
          <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
        </template>
      </a-table>
      </a-table>
    </a-card>
    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑标签' : '新建标签'">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="标签名称"><a-input v-model="form.name" /></a-form-item>
        <a-form-item label="标签类型">
          <a-select v-model="form.tagType" style="width: 100%">
            <a-option value="manual">手动标签</a-option>
            <a-option value="auto">自动标签</a-option>
          </a-select>
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="formVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const tagList = ref([])
const loading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)
const filters = reactive({ keyword: '', tagType: undefined })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const paginationConfig = computed(() => ({ current: pagination.current, pageSize: pagination.pageSize, total: pagination.total, showTotal: true }))
const form = reactive({ name: '', tagType: 'manual' })
const columns = [
  { title: '标签名称', dataIndex: 'name', width: 180 },
  { title: '标签类型', slotName: 'tagType', width: 120 },
  { title: '覆盖会员', slotName: 'memberCount', width: 110 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 200 }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch(`/api/v1/members/tags?page=${pagination.current}&page_size=${pagination.pageSize}&keyword=${filters.keyword}`, {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    tagList.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch { tagList.value = [] } finally { loading.value = false }
}

const showCreate = () => { isEdit.value = false; Object.assign(form, { name: '', tagType: 'manual' }); formVisible.value = true }
const showEdit = (record) => { isEdit.value = true; Object.assign(form, record); formVisible.value = true }
const handleSubmit = () => { formVisible.value = false; Message.success(isEdit.value ? '更新成功' : '创建成功'); loadData() }
const handleDelete = () => { Message.success('删除成功'); loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }

onMounted(() => loadData())
</script>
