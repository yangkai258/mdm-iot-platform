<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="固件列表">
      <template #extra>
        <a-space :size="12">
          <a-button type="primary" @click="uploadFirmware"><icon-upload />上传固件</a-button>
          <a-button @click="loadData"><icon-refresh />刷新</a-button>
        </a-space>
      </template>
      <a-table :columns="columns" :data="firmwares" :loading="loading" :pagination="pagination" row-key="id">
        <template #actions="{ record }">
          <a-space>
            <a-button size="small" @click="viewDetail(record)">详情</a-button>
            <a-button size="small" type="primary" @click="deploy(record)">下发</a-button>
            <a-button size="small" type="primary" status="danger" @click="deleteFirmware(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/breadcrumb'

const loading = ref(false)
const firmwares = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '版本', dataIndex: 'version', width: 120 },
  { title: '文件大小', dataIndex: 'size', width: 120 },
  { title: 'MD5', dataIndex: 'md5', width: 200, ellipsis: true },
  { title: '创建时间', dataIndex: 'created_at', width: 170 }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/miniclaw/firmwares', {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    }).then(r => r.json())
    firmwares.value = res.data?.list || []
    pagination.total = firmwares.value.length
  } catch { firmwares.value = [] } finally { loading.value = false }
}

const uploadFirmware = () => { Message.info('上传固件') }
const viewDetail = (record) => { Message.info('查看详情') }
const deploy = (record) => { Message.success('下发固件') }
const deleteFirmware = (record) => { Message.warn('删除固件') }

onMounted(() => loadData())
</script>
