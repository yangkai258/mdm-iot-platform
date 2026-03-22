<template>
  <div class="buy-gift-promo-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>
        <router-link to="/member/promotions/types">促销活动</router-link>
      </a-breadcrumb-item>
      <a-breadcrumb-item>买赠促销</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card hoverable><a-statistic title="进行中活动" :value="stats.activeCount || 0" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card hoverable><a-statistic title="本月参与会员" :value="stats.joinedCount || 0" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card hoverable><a-statistic title="本月发放赠品" :value="stats.giftedCount || 0" /></a-card>
      </a-col>
    </a-row>

    <a-card class="action-card">
      <a-space>
        <a-button type="primary" @click="showCreate">新建活动</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </a-card>

    <a-card>
      <a-table :columns="columns" :data="promoList" :loading="loading" row-key="id" :pagination="{ pageSize: 10 }">
        <template #name="{ record }">
          <a-tag color="blue">{{ record.name }}</a-tag>
        </template>
        <template #buyCondition="{ record }">
          购买 {{ record.buyProduct }} ≥ {{ record.buyCount }} 件
        </template>
        <template #gift="{ record }">
          赠 {{ record.giftName }} × {{ record.giftCount }}
        </template>
        <template #validPeriod="{ record }">
          {{ record.startTime?.slice(0,10) }} ~ {{ record.endTime?.slice(0,10) }}
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusName(record.status) }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑活动' : '新建活动'"
      @before-ok="handleFormSubmit" @cancel="formVisible = false" :width="560" :loading="formLoading">
      <a-form :model="form" layout="vertical">
        <a-form-item label="活动名称" :rules="[{ required: true, message: '请输入活动名称' }]">
          <a-input v-model="form.name" placeholder="如：买三送一" />
        </a-form-item>
        <a-form-item label="购买条件 - 商品">
          <a-input v-model="form.buyProduct" placeholder="商品名称或品类" />
        </a-form-item>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="购买数量">
              <a-input-number v-model="form.buyCount" :min="1" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="赠送数量">
              <a-input-number v-model="form.giftCount" :min="1" style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="赠品名称">
          <a-input v-model="form.giftName" placeholder="赠品名称" />
        </a-form-item>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="开始时间">
              <a-date-picker v-model="form.startTime" style="width: 100%;" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="结束时间">
              <a-date-picker v-model="form.endTime" style="width: 100%;" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="状态">
          <a-switch v-model="form.status" checked-value="active" unchecked-value="paused" />
        </a-form-item>
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

const form = reactive({
  name: '', buyProduct: '', buyCount: 3, giftName: '', giftCount: 1,
  startTime: null, endTime: null, status: 'active'
})

const columns = [
  { title: '活动名称', slotName: 'name', width: 160 },
  { title: '购买条件', slotName: 'buyCondition', width: 200 },
  { title: '赠品', slotName: 'gift', width: 180 },
  { title: '有效期', slotName: 'validPeriod', width: 220 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 140 }
]

const getStatusColor = (s) => ({ active: 'green', paused: 'orange', finished: 'gray' }[s] || 'gray')
const getStatusName = (s) => ({ active: '进行中', paused: '已暂停', finished: '已结束' }[s] || s)

const loadData = async () => {
  loading.value = true
  try {
    const res = await api.getBuyGiftPromoList()
    promoList.value = res.data?.list || []
    stats.value = res.data?.stats || {}
  } catch (err) {
    Message.error('加载数据失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

const showCreate = () => {
  isEdit.value = false
  Object.assign(form, { name: '', buyProduct: '', buyCount: 3, giftName: '', giftCount: 1, startTime: null, endTime: null, status: 'active' })
  formVisible.value = true
}

const showEdit = (record) => {
  isEdit.value = true
  currentId.value = record.id
  Object.assign(form, { ...record })
  formVisible.value = true
}

const handleFormSubmit = async (done) => {
  formLoading.value = true
  try {
    if (isEdit.value) await api.updateBuyGiftPromo(currentId.value, { ...form })
    else await api.createBuyGiftPromo({ ...form })
    Message.success('保存成功')
    formVisible.value = false
    loadData()
    done(true)
  } catch (err) {
    Message.error(err.message || '操作失败')
    done(false)
  } finally {
    formLoading.value = false
  }
}

const handleDelete = (record) => {
  Modal.warning({ title: '确认删除', content: `确定删除活动「${record.name}」吗？`, okText: '确认删除',
    onOk: async () => {
      try { await api.deleteBuyGiftPromo(record.id); Message.success('删除成功'); loadData() }
      catch (err) { Message.error(err.message || '删除失败') }
    }
  })
}

onMounted(() => loadData())
</script>

<style scoped>
.buy-gift-promo-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
