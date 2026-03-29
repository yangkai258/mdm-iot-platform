<template>
  <div class="company-list-page">
    <a-card class="general-card" title="岗位管理">
      <template #extra>
        <a-button type="primary" @click="openCreateModal(null)"><icon-plus />新建</a-button>
      </template>
      <div class="search-bar">
        <a-input-search v-model="searchKey" placeholder="搜索岗位..." style="width: 260px" @search="loadPosts" allow-clear />
      </div>
      <a-table :columns="columns" :data="filteredData" :loading="loading" :pagination="{ pageSize: 10 }" row-key="id">
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '正常' : '禁用' }}</a-tag>
        </template>
      </a-table>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openCreateModal(record)">「编辑」</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">「删除」</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑岗位' : '新建岗位'" @before-ok="submitForm">
      <a-form :model="formData" layout="vertical">
        <a-form-item label="岗位名称" required>
          <a-input v-model="formData.pos_name" placeholder="请输入" />
        </a-form-item>
        <a-form-item label="岗位编码">
          <a-input v-model="formData.pos_code" placeholder="请输入" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const searchKey = ref('')
const showFilter = ref(false)
const filterStatus = ref<number | ''>('')
const formVisible = ref(false)
const isEdit = ref(false)
const formData = ref({ pos_name: '', pos_code: '' })

const columns = [
  { title: '岗位编码', dataIndex: 'pos_code', width: 120 },
  { title: '岗位名称', dataIndex: 'pos_name', width: 150 },
  { title: '类别', dataIndex: 'category', width: 100 },
  { title: '级别', dataIndex: 'level', width: 80 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' },
]

const mockData = ref([
  { id: 1, pos_code: 'P001', pos_name: '前端开发工程师', category: '技术', level: 'P3', status: 1 },
  { id: 2, pos_code: 'P002', pos_name: '后端开发工程师', category: '技术', level: 'P3', status: 1 },
  { id: 3, pos_code: 'P003', pos_name: '产品经理', category: '产品', level: 'P2', status: 1 },
])

const filteredData = computed(() => {
  let data = mockData.value
  if (searchKey.value) {
    const kw = searchKey.value.toLowerCase()
    data = data.filter(item => item.pos_name.toLowerCase().includes(kw) || item.pos_code.toLowerCase().includes(kw))
  }
  if (filterStatus.value !== '') {
    data = data.filter(item => item.status === filterStatus.value)
  }
  return data
})

const loadPosts = async () => {
  loading.value = true
  await new Promise(r => setTimeout(r, 300))
  loading.value = false
}

const resetFilter = () => {
  filterStatus.value = ''
  loadPosts()
}

const openCreateModal = (record: any) => {
  isEdit.value = !!record
  formData.value = record ? { ...record } : { pos_name: '', pos_code: '' }
  formVisible.value = true
}

const submitForm = async (done: (val: boolean) => void) => {
  Message.success(isEdit.value ? '保存成功' : '创建成功')
  done(true)
}

const handleDelete = (record: any) => {
  Message.success('删除成功')
}

onMounted(() => loadPosts())
</script>

<style scoped>
.company-list-page { padding: 16px; }
.search-bar { margin-bottom: 16px; }
</style>
