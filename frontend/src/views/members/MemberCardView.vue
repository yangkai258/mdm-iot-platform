<template>
  <div class="container">
    <a-card class="general-card" title="会员卡管理">
      <template #extra>
        <a-space :size="12">
          <a-button type="primary" @click="showCreate"><icon-plus />新建</a-button>
          <a-button @click="loadData"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-divider style="margin: 0 0 16px 0" />
      <a-table :columns="columns" :data="cardTypeList" :loading="loading" :pagination="paginationConfig" @page-change="onPageChange" row-key="id">
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
          <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
        </template>
      </a-table>
      </a-table>
    </a-card>
    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑卡类型' : '新建卡类型'">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="名称"><a-input v-model="form.name" /></a-form-item>
        <a-form-item label="编码"><a-input v-model="form.code" /></a-form-item>
        <a-form-item label="折扣率"><a-input-number v-model="form.discountRate" :min="0" :max="1" :precision="2" style="width: 100%" /></a-form-item>
        <a-form-item label="积分倍率"><a-input-number v-model="form.pointsRate" :min="1" style="width: 100%" /></a-form-item>
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

const cardTypeList = ref([])
const loading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const paginationConfig = computed(() => ({ current: pagination.current, pageSize: pagination.pageSize, total: pagination.total, showTotal: true }))
const form = reactive({ name: '', code: '', discountRate: 1, pointsRate: 1 })
const columns = [
  { title: '卡类型名称', dataIndex: 'name', width: 180 },
  { title: '编码', dataIndex: 'code', width: 120 },
  { title: '折扣率', dataIndex: 'discountRate', width: 100 },
  { title: '积分倍率', dataIndex: 'pointsRate', width: 110 },
  { title: '操作', slotName: 'actions', width: 180 }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch(`/api/v1/members/cards/types?page=${pagination.current}&page_size=${pagination.pageSize}`, {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    cardTypeList.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch { cardTypeList.value = [] } finally { loading.value = false }
}

const showCreate = () => { isEdit.value = false; Object.assign(form, { name: '', code: '', discountRate: 1, pointsRate: 1 }); formVisible.value = true }
const showEdit = (record) => { isEdit.value = true; Object.assign(form, record); formVisible.value = true }
const handleSubmit = () => { formVisible.value = false; Message.success(isEdit.value ? '更新成功' : '创建成功'); loadData() }
const handleDelete = () => { Message.success('删除成功'); loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }

onMounted(() => loadData())
</script>
