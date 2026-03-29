<template>
  <div class="container">
    <Breadcrumb :items="['menu.marketing', 'menu.marketing.tempCouponGrant']" />
    <a-card class="general-card" title="临时优惠券发放">
      <template #extra>
        <a-button type="primary" @click="openCreate"><icon-plus />发放</a-button>
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
    </a-card>
    <a-modal v-model="formVisible" title="发放优惠券" :width="560">
      <a-form :model="form" layout="vertical">
        <a-form-item label="会员ID"><a-input v-model="form.member_id" /></a-form-item>
        <a-form-item label="优惠券ID"><a-input v-model="form.coupon_id" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="formVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
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
const form = reactive({ keyword: '', member_id: '', coupon_id: '' })
const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = [
  { title: '发放ID', dataIndex: 'id', width: 80 },
  { title: '会员ID', dataIndex: 'member_id', width: 100 },
  { title: '优惠券ID', dataIndex: 'coupon_id', width: 100 },
  { title: '发放时间', dataIndex: 'created_at', width: 170 }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/v1/marketing/temp-coupon-grants', { headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') } }).then(r => r.json())
    data.value = res.data?.list || []
    pagination.total = data.value.length
  } catch { data.value = [] } finally { loading.value = false }
}
const openCreate = () => { Object.assign(form, { member_id: '', coupon_id: '' }); formVisible.value = true }
const handleSubmit = () => { formVisible.value = false; Message.success('发放成功'); loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }
onMounted(() => loadData())
</script>
