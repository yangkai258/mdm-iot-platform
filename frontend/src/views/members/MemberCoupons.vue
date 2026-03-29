<template>
  <div class="page-container">
    <a-card class="general-card" title="优惠券管理">
      <template #extra>
        <a-button type="primary" @click="handleCreate"><icon-plus />新建</a-button>
      </template>
      <div class="search-form">
        <a-form :model="filters" layout="inline">
          <a-form-item label="名称"><a-input v-model="filters.keyword" placeholder="请输入" /></a-form-item>
          <a-form-item>
            <a-button type="primary" @click="loadCoupons">查询</a-button>
            <a-button @click="filters.keyword = ''; loadCoupons()">重置</a-button>
          </a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="couponList" :loading="loading" :pagination="pagination" row-key="id">
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
          <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
        </template>
      </a-table>
    </a-card>
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
import { IconPlus } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const couponList = ref([])
const showCreateDrawer = ref(false)

const filters = reactive({ status: undefined, keyword: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const stats = reactive({ total: 0, issued: 0, used: 0, unused: 0 })
const form = reactive({ name: '', type: '', value: 0 })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '优惠券名称', dataIndex: 'name' },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '面值', dataIndex: 'value', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '有效期', dataIndex: 'expire_time', width: 160 },
  { title: '操作', slotName: 'actions', width: 160 }
]

const getTypeText = (type) => ({ discount: '折扣券', cash: '现金券', gift: '礼品券' }[type] || type)
const getStatusText = (s) => ({ active: '有效', expired: '已过期', used: '已使用' }[s] || s)
const getStatusColor = (s) => ({ active: 'green', expired: 'gray', used: 'blue' }[s] || 'gray')

const loadCoupons = () => { loading.value = true; setTimeout(() => { loading.value = false }, 300) }
const viewDetail = (r) => Message.info('查看详情')
const distributeCoupon = (r) => Message.info('发放优惠券')
const handleCreate = () => { Message.success('创建成功'); showCreateDrawer.value = false }

onMounted(() => loadCoupons())

</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>


