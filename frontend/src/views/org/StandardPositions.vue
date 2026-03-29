<template>
  <div class="company-list-page">
    <a-card class="general-card" title="基准岗位">
      <template #extra>
        <a-button type="primary" @click="openFormModal(null)"><icon-plus />新建</a-button>
      </template>
      <div class="search-bar">
        <a-input-search v-model="searchKey" placeholder="搜索编码或名称..." style="width: 280px" @search="loadData" allow-clear />
      </div>
      <a-spin :loading="loading">
        <a-table :columns="columns" :data="filteredData" :pagination="{ pageSize: 10 }" row-key="id">
          <template #status="{ record }">
            <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '启用' : '禁用' }}</a-tag>
          </template>
          <template #actions="{ record }">
            <a-space>
              <a-button type="text" size="small" @click="openFormModal(record)">「编辑」</a-button>
              <a-button type="text" size="small" @click="handleCopy(record)">「复制」</a-button>
              <a-button type="text" size="small" status="danger" @click="handleDelete(record)">「删除」</a-button>
            </a-space>
          </template>
        </a-table>
      </a-spin>
    </a-card>

    <!-- 创建/编辑弹窗 -->
    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑基准岗位' : '新建基准岗位'" @before-ok="submitForm" @cancel="formVisible = false" :width="560">
      <a-form :model="formData" layout="vertical">
        <a-form-item label="编码" required>
          <a-input v-model="formData.sp_code" placeholder="请输入编码，如 SP001" />
        </a-form-item>
        <a-form-item label="名称" required>
          <a-input v-model="formData.sp_name" placeholder="请输入岗位名称" />
        </a-form-item>
        <a-form-item label="类别" required>
          <a-select v-model="formData.category" placeholder="请选择类别">
            <a-option value="技术">技术</a-option>
            <a-option value="产品">产品</a-option>
            <a-option value="运营">运营</a-option>
            <a-option value="销售">销售</a-option>
            <a-option value="职能">职能</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="级别" required>
          <a-select v-model="formData.level" placeholder="请选择级别">
            <a-option value="P1">P1</a-option>
            <a-option value="P2">P2</a-option>
            <a-option value="P3">P3</a-option>
            <a-option value="P4">P4</a-option>
            <a-option value="P5">P5</a-option>
            <a-option value="P6">P6</a-option>
            <a-option value="P7">P7</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="formData.description" placeholder="请输入岗位描述" :rows="2" />
        </a-form-item>
        <a-form-item label="职责">
          <a-textarea v-model="formData.responsibility" placeholder="请输入岗位职责" :rows="3" />
        </a-form-item>
        <a-form-item label="任职要求">
          <a-textarea v-model="formData.requirement" placeholder="请输入任职要求" :rows="3" />
        </a-form-item>
        <a-form-item label="技能要求">
          <a-textarea v-model="formData.skills" placeholder="请输入技能要求，多个用逗号分隔" :rows="2" />
        </a-form-item>
        <a-form-item label="状态">
          <a-switch v-model="formData.status" :checked-value="1" :unchecked-value="0" />
          <span style="margin-left: 8px">{{ formData.status === 1 ? '启用' : '禁用' }}</span>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 删除确认弹窗 -->
    <a-modal v-model:visible="deleteVisible" title="确认删除" @before-ok="confirmDelete" @cancel="deleteVisible = false">
      <div>确定要删除基准岗位「{{ deleteTarget?.sp_name }}」吗？删除后不可恢复。</div>
    </a-modal>

    <!-- 复制到部门弹窗 -->
    <a-modal v-model:visible="copyVisible" title="复制到部门" @before-ok="confirmCopy" @cancel="copyVisible = false">
      <a-form layout="vertical">
        <a-form-item label="选择部门" required>
          <a-select v-model="copyDeptId" placeholder="请选择要复制到的部门" :options="deptOptions" />
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
const formVisible = ref(false)
const deleteVisible = ref(false)
const copyVisible = ref(false)
const isEdit = ref(false)
const deleteTarget = ref<any>(null)
const copyTarget = ref<any>(null)
const copyDeptId = ref<number | undefined>(undefined)

const defaultFormData = () => ({
  id: undefined as number | undefined,
  sp_code: '',
  sp_name: '',
  category: '',
  level: '',
  description: '',
  responsibility: '',
  requirement: '',
  skills: '',
  status: 1 as 1 | 0,
})

const formData = ref(defaultFormData())

const deptOptions = ref([
  { label: '技术部', value: 1 },
  { label: '产品部', value: 2 },
  { label: '运营部', value: 3 },
])

const columns = [
  { title: '编码', dataIndex: 'sp_code', width: 120 },
  { title: '名称', dataIndex: 'sp_name', width: 160 },
  { title: '类别', dataIndex: 'category', width: 100 },
  { title: '级别', dataIndex: 'level', width: 80 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 160, fixed: 'right' },
]

const mockData = ref([
  { id: 1, sp_code: 'SP001', sp_name: '前端开发工程师', category: '技术', level: 'P3', description: '负责前端开发', responsibility: '完成前端开发任务', requirement: '3年以上经验', skills: 'Vue,React,TypeScript', status: 1 },
  { id: 2, sp_code: 'SP002', sp_name: '后端开发工程师', category: '技术', level: 'P3', description: '负责后端开发', responsibility: '完成后端开发任务', requirement: '3年以上经验', skills: 'Go,Python,Java', status: 1 },
  { id: 3, sp_code: 'SP003', sp_name: '产品经理', category: '产品', level: 'P2', description: '负责产品规划', responsibility: '产品需求分析', requirement: '5年以上经验', skills: 'PRD,Axure,数据分析', status: 1 },
  { id: 4, sp_code: 'SP004', sp_name: '运营专员', category: '运营', level: 'P1', description: '负责运营工作', responsibility: '运营执行', requirement: '1年以上经验', skills: '内容运营,用户运营', status: 0 },
  { id: 5, sp_code: 'SP005', sp_name: '销售经理', category: '销售', level: 'P4', description: '负责销售管理', responsibility: '销售目标达成', requirement: '5年以上经验', skills: '客户管理,谈判', status: 1 },
])

const filteredData = computed(() => {
  if (!searchKey.value) return mockData.value
  const kw = searchKey.value.toLowerCase()
  return mockData.value.filter(item =>
    item.sp_code.toLowerCase().includes(kw) || item.sp_name.toLowerCase().includes(kw)
  )
})

const loadData = async () => {
  loading.value = true
  await new Promise(r => setTimeout(r, 300))
  loading.value = false
}

const openFormModal = (record: any) => {
  isEdit.value = !!record
  if (record) {
    formData.value = { ...record }
  } else {
    formData.value = defaultFormData()
  }
  formVisible.value = true
}

const submitForm = async (done: (val: boolean) => void) => {
  if (!formData.value.sp_code || !formData.value.sp_name || !formData.value.category || !formData.value.level) {
    Message.error('请填写必填项')
    done(false)
    return
  }
  if (isEdit.value) {
    const idx = mockData.value.findIndex((item: any) => item.id === formData.value.id)
    if (idx !== -1) mockData.value[idx] = { ...formData.value }
    Message.success('保存成功')
  } else {
    const newId = Math.max(...mockData.value.map((item: any) => item.id), 0) + 1
    mockData.value.unshift({ ...formData.value, id: newId })
    Message.success('创建成功')
  }
  done(true)
}

const handleDelete = (record: any) => {
  deleteTarget.value = record
  deleteVisible.value = true
}

const confirmDelete = (done: (val: boolean) => void) => {
  mockData.value = mockData.value.filter((item: any) => item.id !== deleteTarget.value.id)
  Message.success('删除成功')
  done(true)
}

const handleCopy = (record: any) => {
  copyTarget.value = record
  copyDeptId.value = undefined
  copyVisible.value = true
}

const confirmCopy = (done: (val: boolean) => void) => {
  if (!copyDeptId.value) {
    Message.error('请选择部门')
    done(false)
    return
  }
  Message.success(`已复制「${copyTarget.value.sp_name}」到目标部门`)
  done(true)
}

onMounted(() => loadData())
</script>

<style scoped>
.company-list-page { padding: 16px; }
.search-bar { margin-bottom: 16px; }
</style>
