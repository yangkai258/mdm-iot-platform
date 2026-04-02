<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="company-list-page">
    <a-card class="general-card" title="权限组管理">
      <template #extra>
        <a-button type="primary" @click="addGroup"><icon-plus />新建</a-button>
      </template>
      <div class="search-bar">
        <a-input-search v-model="searchKey" placeholder="搜索..." style="width: 260px" @search="loadGroups" />
      </div>
      <a-table :columns="columns" :data="groups" :loading="loading" :pagination="{ pageSize: 10 }" row-key="id">
        <template #status="{ record }">
          <a-tag :color="record.status === 'active' ? 'green' : 'gray'">{{ record.status === 'active' ? '启用' : '禁用' }}</a-tag>
        </template>
      </a-table>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="editGroup(record)">「编辑」</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">「删除」</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑权限组' : '新建权限组'" @before-ok="submitForm">
      <a-form :model="formData" layout="vertical">
        <a-form-item label="权限组名称" required>
          <a-input v-model="formData.group_name" placeholder="请输入" />
        </a-form-item>
        <a-form-item label="权限组编码">
          <a-input v-model="formData.group_code" placeholder="请输入" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'

const searchKey = ref('')
const loading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)
const formData = ref({ group_name: '', group_code: '' })

const columns = [
  { title: '权限组编码', dataIndex: 'group_code', width: 150 },
  { title: '权限组名称', dataIndex: 'group_name', width: 200 },
  { title: '描述', dataIndex: 'description', width: 200 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' },
]

const groups = ref([
  { id: 1, group_code: 'G001', group_name: '管理员', description: '系统管理员', status: 'active' },
  { id: 2, group_code: 'G002', group_name: '普通用户', description: '普通用户权限', status: 'active' },
])

const loadGroups = async () => {
  loading.value = true
  await new Promise(r => setTimeout(r, 300))
  loading.value = false
}

const addGroup = () => {
  isEdit.value = false
  formData.value = { group_name: '', group_code: '' }
  formVisible.value = true
}

const editGroup = (record: any) => {
  isEdit.value = true
  formData.value = { ...record }
  formVisible.value = true
}

const submitForm = async (done: (val: boolean) => void) => {
  Message.success(isEdit.value ? '保存成功' : '创建成功')
  done(true)
}

const handleDelete = (record: any) => {
  Message.success('删除成功')
}

onMounted(() => loadGroups())
</script>

<style scoped>
.company-list-page { padding: 16px; }
.search-bar { margin-bottom: 16px; }
</style>
