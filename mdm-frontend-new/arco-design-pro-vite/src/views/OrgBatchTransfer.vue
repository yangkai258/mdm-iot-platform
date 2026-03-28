<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-swap /> 批量调岗操作</a-space>
      </template>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-card title="选择员工">
            <a-form layout="vertical">
              <a-form-item label="选择部门">
                <a-tree-select v-model="formData.sourceDepartment" :data="deptTree" placeholder="选择源部门" />
              </a-form-item>
              <a-form-item label="选择员工">
                <a-transfer :data="employeeList" mode="multiple" />
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>

        <a-col :span="12">
          <a-card title="调岗配置">
            <a-form :model="formData" layout="vertical">
              <a-form-item label="目标部门" required>
                <a-tree-select v-model="formData.targetDepartment" :data="deptTree" placeholder="选择目标部门" />
              </a-form-item>
              <a-form-item label="目标岗位">
                <a-select v-model="formData.targetPosition" placeholder="选择岗位">
                  <a-option value="p1">初级</a-option>
                  <a-option value="p2">中级</a-option>
                  <a-option value="p3">高级</a-option>
                </a-select>
              </a-form-item>
              <a-form-item label="生效日期">
                <a-date-picker v-model="formData.effectiveDate" style="width: 100%" />
              </a-form-item>
            </a-form>
          </a-card>

          <a-card title="影响预览" style="margin-top: 16px">
            <a-descriptions :column="1" bordered>
              <a-descriptions-item label="调岗人数">{{ selectedEmployees.length }} 人</a-descriptions-item>
              <a-descriptions-item label="原部门">{{ formData.sourceDepartment || '-' }}</a-descriptions-item>
              <a-descriptions-item label="目标部门">{{ formData.targetDepartment || '-' }}</a-descriptions-item>
            </a-descriptions>
          </a-card>
        </a-col>
      </a-row>

      <a-space style="margin-top: 16px">
        <a-button type="primary" @click="handleSubmit">执行调岗</a-button>
        <a-button @click="handleCancel">取消</a-button>
      </a-space>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const formData = reactive({
  sourceDepartment: '', targetDepartment: '', targetPosition: '', effectiveDate: ''
})
const selectedEmployees = ref([])
const employeeList = []

const deptTree = [
  { value: 'd1', label: '研发部', children: [{ value: 'd1-1', label: '前端组' }, { value: 'd1-2', label: '后端组' }] },
  { value: 'd2', label: '运维部' }
]

const handleSubmit = () => { }
const handleCancel = () => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
