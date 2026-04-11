<template>
  <div class="page-container">
    <a-card class="general-card" title="搴旂敤甯傚満">
      <template #extra>
        <a-space>
          <a-button @click="loadData"><icon-refresh />鍒锋柊</a-button>
        </a-space>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="搴旂敤鍚嶇О">
            <a-input v-model="form.name" placeholder="璇疯緭鍏ュ簲鐢ㄥ悕锟" style="width: 160px" />
          </a-form-item>
          <a-form-item label="鍒嗙被">
            <a-select v-model="form.category" placeholder="璇烽€夋嫨" allow-clear style="width: 140px">
              <a-option value="tool">宸ュ叿锟?/a-option>
              <a-option value="game">娓告垙锟?/a-option>
              <a-option value="education">鏁欒偛锟?/a-option>
              <a-option value="social">绀句氦锟?/a-option>
            </a-select>
          </a-form-item>
          <a-form-item>
            <a-button type="primary" @click="loadData">鏌ヨ</a-button>
            <a-button @click="handleReset">閲嶇疆</a-button>
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
          <span v-if="record.price === 0" class="free-tag">鍏嶈垂</span>
          <span v-else>楼{{ record.price }}</span>
        </template>
        <template #actions="{ record }">
          <a-button type="primary" size="small" @click="handleInstall(record)">瀹夎</a-button>
          <a-button size="small" @click="handleViewDetail(record)">璇︽儏</a-button>
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
  { title: '鍥炬爣', slotName: 'icon', width: 80 },
  { title: '搴旂敤鍚嶇О', dataIndex: 'name', width: 180 },
  { title: '鍒嗙被', dataIndex: 'category_name', width: 100 },
  { title: '鐗堟湰', dataIndex: 'version', width: 80 },
  { title: '璇勫垎', slotName: 'rating', width: 160 },
  { title: '涓嬭浇锟?, dataIndex: 'downloads', width: 100 },
  { title: '浠锋牸', slotName: 'price', width: 80 },
  { title: '鎿嶄綔', slotName: 'actions', width: 160 }
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
    { id: 1, name: '瀹犵墿鎶曞杺鍔╂墜', category: 'tool', category_name: '宸ュ叿锟?, version: '1.2.0', rating: 4.5, downloads: 12580, price: 0, icon: '' },
    { id: 2, name: '瓒ｅ懗璁粌娓告垙', category: 'game', category_name: '娓告垙锟?, version: '2.0.1', rating: 4.8, downloads: 8960, price: 6, icon: '' },
    { id: 3, name: '瀹犵墿鍋ュ悍璇惧爞', category: 'education', category_name: '鏁欒偛锟?, version: '1.0.5', rating: 4.2, downloads: 5620, price: 0, icon: '' },
    { id: 4, name: '钀屽疇绀惧尯', category: 'social', category_name: '绀句氦锟?, version: '1.5.2', rating: 4.6, downloads: 23400, price: 0, icon: '' },
    { id: 5, name: '鏅鸿兘鎷嶇収', category: 'tool', category_name: '宸ュ叿锟?, version: '1.1.0', rating: 4.3, downloads: 9870, price: 3, icon: '' }
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
  Message.success(`宸插紑濮嬪畨锟? ${record.name}`)
}

const handleViewDetail = (record) => {
  Message.info(`鏌ョ湅璇︽儏: ${record.name}`)
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