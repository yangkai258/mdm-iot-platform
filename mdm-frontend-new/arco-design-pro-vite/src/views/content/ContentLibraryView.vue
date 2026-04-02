<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="文件库">
      <template #extra>
        <a-space>
          <a-select v-model="filter.category" placeholder="选择分类" style="width: 120px" allow-clear>
            <a-option value="emoticon">表情包</a-option>
            <a-option value="action">动作</a-option>
            <a-option value="voice">声音</a-option>
            <a-option value="wallpaper">壁纸</a-option>
          </a-select>
          <a-button type="primary"><icon-upload />上传文件</a-button>
        </a-space>
      </template>
      <a-spin :loading="loading">
        <a-table :columns="columns" :data="files" :pagination="{ pageSize: 12 }" row-key="id">
          <template #file_name="{ record }">
            <div style="display: flex; align-items: center; gap: 8px">
              <icon-file :size="20" />
              {{ record.file_name }}
            </div>
          </template>
      </a-table>
          <template #file_type="{ record }"><a-tag>{{ record.file_type }}</a-tag></template>
          <template #actions="{ record }">
            <a-button size="small" type="text" @click="download(record)">下载</a-button>
            <a-button size="small" type="text" @click="distribute(record)">分发</a-button>
            <a-button size="small" type="text" status="danger" @click="remove(record)">删除</a-button>
          </template>
        </a-table>
      </a-spin>
    </a-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const loading = ref(false)
const filter = ref({ category: '' })
const files = ref([])
const columns = [
  { title: '文件名', dataIndex: 'file_name', slotName: 'file_name' },
  { title: '类型', dataIndex: 'file_type', slotName: 'file_type' },
  { title: '大小', dataIndex: 'file_size', render: ({ file_size }) => `${(file_size/1024).toFixed(1)} KB` },
  { title: '下载量', dataIndex: 'download_count' },
  { title: '操作', slotName: 'actions', width: 180 }
]

const load = async () => {
  loading.value = true
  const url = filter.value.category ? `/api/v1/content/files?category=${filter.value.category}` : '/api/v1/content/files'
  const res = await fetch(url, { headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') } })
  const data = await res.json()
  files.value = data.files || []
  loading.value = false
}

const download = async (record) => {
  const res = await fetch(`/api/v1/content/files/${record.id}/download`, { headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') } })
  const data = await res.json()
  if (data.download_url) window.open(data.download_url, '_blank')
}
const distribute = () => Message.info('分发功能开发中')
const remove = () => Message.success('删除成功')

onMounted(load)
</script>
