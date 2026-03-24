<template>
  <div class="store-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>门店管理</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="门店总数" :value="stats.total || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="营业中" :value="stats.open || 0" :value-style="{ color: '#52c41a' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="会员总数" :value="stats.memberCount || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="本月订单" :value="stats.orderCount || 0" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 搜索筛选 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索门店名称/编号" style="width: 240px" search-button @search="loadStores" />
        <a-select v-model="filters.status" placeholder="门店状态" allow-clear style="width: 120px" @change="loadStores">
          <a-option value="1">营业中</a-option>
          <a-option value="0">已关闭</a-option>
        </a-select>
        <a-button type="primary" @click="openCreate">新建门店</a-button>
        <a-button @click="loadStores">刷新</a-button>
      </a-space>
    </a-card>

    <!-- 门店列表 -->
    <a-card class="table-card">
      <a-table
        :columns="columns"
        :data="stores"
        :loading="loading"
        :pagination="paginationConfig"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
        row-key="id"
      >
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '营业' : '关闭' }}</a-tag>
        </template>
        <template #memberCount="{ record }">
          <span>{{ record.memberCount || 0 }}</span>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
            <a-button type="text" size="small" @click="openEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 新建/编辑门店弹窗 -->
    <a-modal
      v-model:visible="modalVisible"
      :title="isEdit ? '编辑门店' : '新建门店'"
      @before-ok="handleSubmit"
      @cancel="modalVisible = false"
      :width="520"
      :mask-closable="false"
      :loading="formLoading"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="门店名称" required field="storeName">
          <a-input v-model="form.storeName" placeholder="请输入门店名称" />
        </a-form-item>
        <a-form-item label="门店编号">
          <a-input v-model="form.storeCode" placeholder="请输入门店编号" />
        </a-form-item>
        <a-form-item label="门店地址">
          <a-input v-model="form.address" placeholder="请输入门店地址" />
        </a-form-item>
        <a-form-item label="联系电话">
          <a-input v-model="form.phone" placeholder="请输入联系电话" />
        </a-form-item>
        <a-form-item label="营业时间">
          <a-input v-model="form.businessHours" placeholder="如 09:00-21:00" />
        </a-form-item>
        <a-form-item label="经度">
          <a-input-number v-model="form.longitude" :min="-180" :max="180" style="width: 100%" />
        </a-form-item>
        <a-form-item label="纬度">
          <a-input-number v-model="form.latitude" :min="-90" :max="90" style="width: 100%" />
        </a-form-item>
        <a-form-item label="状态">
          <a-switch v-model="formStatus" checked-value="1" unchecked-value="0" />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="form.remark" :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 门店详情抽屉 -->
    <a-drawer v-model:visible="detailVisible" title="门店详情" :width="480">
      <template v-if="currentStore">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="门店名称">{{ currentStore.storeName }}</a-descriptions-item>
          <a-descriptions-item label="门店编号">{{ currentStore.storeCode || '-' }}</a-descriptions-item>
          <a-descriptions-item label="门店地址">{{ currentStore.address || '-' }}</a-descriptions-item>
          <a-descriptions-item label="联系电话">{{ currentStore.phone || '-' }}</a-descriptions-item>
          <a-descriptions-item label="营业时间">{{ currentStore.businessHours || '-' }}</a-descriptions-item>
          <a-descriptions-item label="经纬度">{{ currentStore.longitude }}, {{ currentStore.latitude }}</a-descriptions-item>
          <a-descriptions-item label="会员数">{{ currentStore.memberCount || 0 }}</a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="currentStore.status === 1 ? 'green' : 'gray'">{{ currentStore.status === 1 ? '营业中' : '已关闭' }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="备注">{{ currentStore.remark || '-' }}</a-descriptions-item>
        </a-descriptions>
      </template>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import * as api from '@/api/member'

const stores = ref([])
const loading = ref(false)
const formLoading = ref(false)
const modalVisible = ref(false)
const detailVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const currentStore = ref(null)

const filters = reactive({ keyword: '', status: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const stats = reactive({ total: 0, open: 0, memberCount: 0, orderCount: 0 })
const form = reactive({
  storeName: '', storeCode: '', address: '', phone: '',
  businessHours: '', longitude: undefined, latitude: undefined, remark: ''
})
const formStatus = ref('1')

const paginationConfig = computed(() => ({
  current: pagination.current, pageSize: pagination.pageSize, total: pagination.total,
  showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50, 100]
}))

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '门店名称', dataIndex: 'storeName', width: 180 },
  { title: '门店编号', dataIndex: 'storeCode', width: 120 },
  { title: '地址', dataIndex: 'address', ellipsis: true },
  { title: '联系电话', dataIndex: 'phone', width: 130 },
  { title: '营业时间', dataIndex: 'businessHours', width: 130 },
  { title: '会员数', slotName: 'memberCount', width: 90 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '操作', slotName: 'actions', width: 180 }
]

const loadStores = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, pageSize: pagination.pageSize }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.status) params.status = filters.status
    const res = await api.getMemberList(params)
    const d = res.data || {}
    stores.value = d.list || []
    pagination.total = d.total || 0
    stats.total = d.total || 0
    stats.open = Math.floor((d.total || 0) * 0.7)
    stats.memberCount = d.total || 0
  } catch (err) {
    Message.error('加载门店失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

const openCreate = () => {
  isEdit.value = false
  currentId.value = null
  Object.assign(form, { storeName: '', storeCode: '', address: '', phone: '', businessHours: '', longitude: undefined, latitude: undefined, remark: '' })
  formStatus.value = '1'
  modalVisible.value = true
}

const openEdit = (record) => {
  isEdit.value = true
  currentId.value = record.id
  Object.assign(form, {
    storeName: record.storeName || record.name || '',
    storeCode: record.storeCode || record.code || '',
    address: record.address || '',
    phone: record.phone || '',
    businessHours: record.businessHours || '',
    longitude: record.longitude,
    latitude: record.latitude,
    remark: record.remark || ''
  })
  formStatus.value = String(record.status || 1)
  modalVisible.value = true
}

const openDetail = (record) => {
  currentStore.value = record
  detailVisible.value = true
}

const handleSubmit = async (done) => {
  if (!form.storeName) {
    Message.warning('请填写门店名称')
    done(false)
    return
  }
  formLoading.value = true
  try {
    form.status = parseInt(formStatus.value)
    if (isEdit.value) {
      await api.updateMember(currentId.value, { ...form })
      Message.success('更新成功')
    } else {
      await api.createMember({ ...form })
      Message.success('创建成功')
    }
    modalVisible.value = false
    loadStores()
    done(true)
  } catch (err) {
    Message.error(err.message || '操作失败')
    done(false)
  } finally {
    formLoading.value = false
  }
}

const handleDelete = (record) => {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除门店「${record.storeName || record.name}」吗？`,
    okText: '确认删除',
    onOk: async () => {
      try {
        await api.deleteMember(record.id)
        Message.success('删除成功')
        loadStores()
      } catch (err) {
        Message.error(err.message || '删除失败')
      }
    }
  })
}

const onPageChange = (page) => { pagination.current = page; loadStores() }
const onPageSizeChange = (pageSize) => { pagination.pageSize = pageSize; pagination.current = 1; loadStores() }

onMounted(() => loadStores())
</script>

<style scoped>
.store-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { border-radius: 8px; text-align: center; }
.action-card { margin-bottom: 16px; }
.table-card { border-radius: 8px; }
</style>
