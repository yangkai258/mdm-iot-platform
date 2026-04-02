<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="amount-discount-promo-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item><router-link to="/member/promotions/types">促销活动</router-link></a-breadcrumb-item>
      <a-breadcrumb-item>满额折促销</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8"><a-card hoverable><a-statistic title="进行中活动" :value="stats.activeCount || 0" /></a-card></a-col>
      <a-col :span="8"><a-card hoverable><a-statistic title="参与会员数" :value="stats.joinedCount || 0" /></a-card></a-col>
      <a-col :span="8"><a-card hoverable><a-statistic title="优惠总金额" :value="(stats.totalSaved || 0).toLocaleString() + ' 元'" /></a-card></a-col>
    </a-row>

    <a-card class="action-card">
      <a-space>
        <a-button type="primary" @click="showCreate">新建活动</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </a-card>

    <a-card>
      <a-table :columns="columns" :data="promoList" :loading="loading" row-key="id" :pagination="{ pageSize: 10 }">
        <template #name="{ record }"><a-tag color="cyan">{{ record.name }}</a-tag></template>
        <template #thresholds="{ record }">
          <a-space wrap>
            <a-tag v-for="(t, i) in (record.thresholds || [])" :key="i" size="small" color="cyan">
              满{{ t.threshold }}享{{ (t.discountRate * 10).toFixed(1) }}折
            </a-tag>
          </a-space>
        </template>
      </a-table>
        <template #validPeriod="{ record }">{{ record.startTime?.slice(0,10) }} ~ {{ record.endTime?.slice(0,10) }}</template>
        <template #status="{ record }"><a-tag :color="getStatusColor(record.status)">{{ getStatusName(record.status) }}</a-tag></template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑活动' : '新建活动'"
      @before-ok="handleFormSubmit" @cancel="formVisible = false" :width="580" :loading="formLoading">
      <a-form :model="form" layout="vertical">
        <a-form-item label="活动名称" :rules="[{ required: true, message: '请输入活动名称' }]">
          <a-input v-model="form.name" placeholder="如：夏装满500享8折" />
        </a-form-item>
        <a-form-item label="满折档位设置">
          <a-space direction="vertical" style="width: 100%;">
            <a-row v-for="(t, i) in form.thresholds" :key="i" :gutter="8" align="center">
              <a-col :span="10"><a-input-number v-model="t.threshold" :min="1" placeholder="满多少" style="width: 100%;" /></a-col>
              <a-col :span="10"><a-input-number v-model="t.discountRate" :min="0.1" :max="1" :step="0.05" placeholder="折扣率(0.8)" style="width: 100%;" /></a-col>
              <a-col :span="4"><a-button type="text" status="danger" @click="form.thresholds.splice(i, 1)">删除</a-button></a-col>
            </a-row>
            <a-button type="dashed" @click="form.thresholds.push({ threshold: 100, discountRate: 0.9 })">+ 添加档位</a-button>
          </a-space>
        </a-form-item>
        <a-row :gutter="12">
          <a-col :span="12"><a-form-item label="开始时间"><a-date-picker v-model="form.startTime" style="width:100%;" /></a-form-item></a-col>
          <a-col :span="12"><a-form-item label="结束时间"><a-date-picker v-model="form.endTime" style="width:100%;" /></a-form-item></a-col>
        </a-row>
        <a-form-item label="状态"><a-switch v-model="form.status" checked-value="active" unchecked-value="paused" /></a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import * as api from '@/api/marketing'

const promoList = ref([])
const stats = ref({})
const loading = ref(false)
const formLoading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)

const form = reactive({ name: '', thresholds: [{ threshold: 100, discountRate: 0.9 }], startTime: null, endTime: null, status: 'active' })

const columns = [
  { title: '活动名称', slotName: 'name', width: 180 },
  { title: '满折档位', slotName: 'thresholds' },
  { title: '有效期', slotName: 'validPeriod', width: 220 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 140 }
]

const getStatusColor = (s) => ({ active: 'green', paused: 'orange', finished: 'gray' }[s] || 'gray')
const getStatusName = (s) => ({ active: '进行中', paused: '已暂停', finished: '已结束' }[s] || s)

const loadData = async () => {
  loading.value = true
  try {
    const res = await api.getAmountDiscountPromoList()
    promoList.value = res.data?.list || []
    stats.value = res.data?.stats || {}
  } catch (err) { Message.error('加载数据失败: ' + err.message) } finally { loading.value = false }
}

const showCreate = () => {
  isEdit.value = false
  Object.assign(form, { name: '', thresholds: [{ threshold: 100, discountRate: 0.9 }], startTime: null, endTime: null, status: 'active' })
  formVisible.value = true
}
const showEdit = (record) => {
  isEdit.value = true; currentId.value = record.id
  Object.assign(form, { ...record, thresholds: record.thresholds ? [...record.thresholds] : [{ threshold: 100, discountRate: 0.9 }] })
  formVisible.value = true
}

const handleFormSubmit = async (done) => {
  formLoading.value = true
  try {
    if (isEdit.value) await api.updateAmountDiscountPromo(currentId.value, { ...form })
    else await api.createAmountDiscountPromo({ ...form })
    Message.success('保存成功'); formVisible.value = false; loadData(); done(true)
  } catch (err) { Message.error(err.message || '操作失败'); done(false) } finally { formLoading.value = false }
}

const handleDelete = (record) => {
  Modal.warning({ title: '确认删除', content: `确定删除活动「${record.name}」吗？`, okText: '确认删除',
    onOk: async () => { try { await api.deleteAmountDiscountPromo(record.id); Message.success('删除成功'); loadData() } catch (err) { Message.error(err.message || '删除失败') } }
  })
}

onMounted(() => loadData())
</script>

<style scoped>
.amount-discount-promo-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
