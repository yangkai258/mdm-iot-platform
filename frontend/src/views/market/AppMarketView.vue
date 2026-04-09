<template>
  <div class="page-container">
    <a-card class="general-card" title="应用市场">
      <template #extra>
        <a-space>
          <a-button @click="loadData"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="应用名称">
            <a-input v-model="form.name" placeholder="请输入应用名�? style="width: 160px" />
          </a-form-item>
          <a-form-item label="分类">
            <a-select v-model="form.category" placeholder="请选择" allow-clear style="width: 140px">
              <a-option value="tool">工具�?/a-option>
              <a-option value="game">游戏�?/a-option>
              <a-option value="education">教育�?/a-option>
              <a-option value="social">社交�?/a-option>
            </a-select>
          </a-form-item>
          <a-form-item>
            <a-button type="primary" @click="loadData">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #icon="{ record }">
          <a-avatar :size="40" shape="square">
            <img v-if="record.icon" :src="record.icon" />
            <icon-apps v-else />
          </a-avatar>
        </template>
        <template #rating="{ record }">
          <a-rate :model-value="record.rating" allow-half disabled :count="5" />
          <span style="margin-left: 8px">{{ record.rating?.toFixed(1) }}</span>
        </template>
        <template #price="{ record }">
          <span v-if="record.price === 0" class="free-tag">免费</span>
          <span v-else>¥{{ record.price }}</span>
        </template>
        <template #actions="{ record }">
          <a-button type="primary" size="small" @click="handleInstall(record)">安装</a-button>
          <a-button size="small" @click="handleViewDetail(record)">详情</a-button>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconApps, IconRefresh } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const form = reactive({ name: '', category: '' })
const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: '图标', slotName: 'icon', width: 80 },
  { title: '应用名称', dataIndex: 'name', width: 180 },
  { title: '分类', dataIndex: 'category_name', width: 100 },
  { title: '版本', dataIndex: 'version', width: 80 },
  { title: '评分', slotName: 'rating', width: 160 },
  { title: '下载�?, dataIndex: 'downloads', width: 100 },
  { title: '价格', slotName: 'price', width: 80 },
  { title: '操作', slotName: 'actions', width: 160 }
]

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams()
    if (form.name) params.append('name', form.name)
    if (form.category) params.append('category', form.category)
    params.append('page', pagination.current)
    params.append('page_size', pagination.pageSize)

    const res = await fetch(`/api/v1/market/apps?${params}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    }).then(r => r.json())

    if (res.code === 0) {
      data.value = res.data?.list || []
      pagination.total = res.data?.total || 0
    } else {
      loadMockData()
    }
  } catch { loadMockData() } finally { loading.value = false }
}

const loadMockData = () => {
  data.value = [
    { id: 1, name: '宠物投喂助手', category: 'tool', category_name: '工具�?, version: '1.2.0', rating: 4.5, downloads: 12580, price: 0, icon: '' },
    { id: 2, name: '趣味训练游戏', category: 'game', category_name: '游戏�?, version: '2.0.1', rating: 4.8, downloads: 8960, price: 6, icon: '' },
    { id: 3, name: '宠物健康课堂', category: 'education', category_name: '教育�?, version: '1.0.5', rating: 4.2, downloads: 5620, price: 0, icon: '' },
    { id: 4, name: '萌宠社区', category: 'social', category_name: '社交�?, version: '1.5.2', rating: 4.6, downloads: 23400, price: 0, icon: '' },
    { id: 5, name: '智能拍照', category: 'tool', category_name: '工具�?, version: '1.1.0', rating: 4.3, downloads: 9870, price: 3, icon: '' }
  ]
  pagination.total = data.value.length
}

const handleReset = () => {
  form.name = ''
  form.category = ''
  pagination.current = 1
  loadData()
}

const handleInstall = (record) => {
  Message.success(`已开始安�? ${record.name}`)
}

const handleViewDetail = (record) => {
  Message.info(`查看详情: ${record.name}`)
}

const onPageChange = (page) => {
  pagination.current = page
  loadData()
}

onMounted(() => loadData())
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
.free-tag { color: #52c41a; font-weight: 500; }
</style>