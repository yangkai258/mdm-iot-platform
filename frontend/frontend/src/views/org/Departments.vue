<template>
  <div class="company-list-page">
    <div class="breadcrumb-wrapper">
      <a-breadcrumb>
        <a-breadcrumb-item><a href="#/dashboard">首页</a></a-breadcrumb-item>
        <a-breadcrumb-item>组织管理</a-breadcrumb-item>
        <a-breadcrumb-item>部门管理</a-breadcrumb-item>
      </a-breadcrumb>
    </div>

    <div class="toolbar">
      <div class="toolbar-left">
        <a-input-search v-model="searchKey" placeholder="搜索部门..." style="width: 260px" @search="loadDepartments" />
      </div>
      <div class="toolbar-right">
        <a-button type="primary" @click="openCreateModal(null)">「新建」</a-button>
        <a-button @click="loadDepartments">「刷新」</a-button>
      </div>
    </div>

    <a-card :bordered="false" class="table-card">
      <a-spin :loading="loading">
        <a-tree :data="treeData" :block-node="true" :show-line="true" row-key="id" @expand="onExpand">
          <template #title="node">
            <div class="tree-node-content">
              <span class="node-name">{{ node.dept_name }}</span>
              <span class="node-code">({{ node.dept_code }})</span>
              <a-space class="node-actions">
                <a-button type="text" size="mini" @click.stop="openCreateModal(node)">「新增」</a-button>
                <a-button type="text" size="mini" @click.stop="openEditModal(node)">「编辑」</a-button>
                <a-button type="text" size="mini" status="danger" @click.stop="handleDelete(node)">「删除」</a-button>
              </a-space>
            </div>
          </template>
        </a-tree>
      </a-spin>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑部门' : '新建部门'" @before-ok="submitForm">
      <a-form :model="formData" layout="vertical">
        <a-form-item label="部门名称" required>
          <a-input v-model="formData.dept_name" placeholder="请输入" />
        </a-form-item>
        <a-form-item label="上级部门">
          <a-tree-select v-model="formData.parent_id" :data="treeData" placeholder="请选择" allow-clear />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const searchKey = ref('')
const formVisible = ref(false)
const isEdit = ref(false)
const formData = ref({ dept_name: '', parent_id: undefined })
const treeData = ref([
  { id: 1, dept_name: '技术部', dept_code: 'D001', children: [{ id: 2, dept_name: '前端组', dept_code: 'D001-1' }, { id: 3, dept_name: '后端组', dept_code: 'D001-2' }] },
  { id: 4, dept_name: '运营部', dept_code: 'D002', children: [] },
])

const loadDepartments = async () => {
  loading.value = true
  await new Promise(r => setTimeout(r, 300))
  loading.value = false
}

const openCreateModal = (parent: any) => {
  isEdit.value = false
  formData.value = { dept_name: '', parent_id: parent?.id }
  formVisible.value = true
}

const openEditModal = (node: any) => {
  isEdit.value = true
  formData.value = { dept_name: node.dept_name, parent_id: node.parent_id }
  formVisible.value = true
}

const submitForm = async (done: (val: boolean) => void) => {
  Message.success(isEdit.value ? '保存成功' : '创建成功')
  done(true)
}

const handleDelete = (node: any) => {
  Message.success('删除成功')
}

const onExpand = (keys: string[]) => {}

onMounted(() => loadDepartments())
</script>

<style scoped>
.company-list-page { padding: 24px; min-height: 100%; background: #f2f3f5; }
.breadcrumb-wrapper { margin-bottom: 16px; }
.toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; background: #fff; padding: 16px; border-radius: 4px; }
.toolbar-left, .toolbar-right { display: flex; gap: 8px; }
.table-card { background: #fff; border-radius: 4px; }
.tree-node-content { display: flex; align-items: center; gap: 8px; }
.node-code { color: #999; font-size: 12px; }
.node-actions { margin-left: auto; }
</style>
