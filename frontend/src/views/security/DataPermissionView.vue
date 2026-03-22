<template>
  <div class="page-container">
    <!-- 顶部筛选栏 -->
    <a-card class="filter-card">
      <div class="filter-row">
        <a-select
          v-model="selectedRoleId"
          placeholder="角色选择"
          style="width: 180px"
          allow-search
          @change="loadDataPermissions"
        >
          <a-option v-for="r in roles" :key="r.id" :value="r.id">{{ r.name }}</a-option>
        </a-select>
        <a-button type="primary" @click="handleSave">保存配置</a-button>
      </div>
    </a-card>

    <!-- 内容区 -->
    <a-card class="content-card" v-if="selectedRoleId">
      <!-- 数据范围 -->
      <div class="section-title">数据范围</div>
      <div class="scope-section">
        <a-radio-group v-model="form.dataScope">
          <a-radio value="all">
            <span class="radio-label">全部数据</span>
            <span class="radio-desc">可访问所有数据</span>
          </a-radio>
          <a-radio value="department">
            <span class="radio-label">本部门</span>
            <span class="radio-desc">仅可访问所属部门数据</span>
          </a-radio>
          <a-radio value="self">
            <span class="radio-label">本人</span>
            <span class="radio-desc">仅可访问自己创建的数据</span>
          </a-radio>
        </a-radio-group>
      </div>

      <a-divider />

      <!-- 列级权限 -->
      <div class="section-title">字段权限</div>
      <div class="field-section">
        <div class="field-group-title">会员字段</div>
        <a-checkbox
          v-for="field in memberFields"
          :key="field.key"
          v-model="field.checked"
          class="perm-check"
        >{{ field.label }}</a-checkbox>
        <div class="field-group-title" style="margin-top:12px">设备字段</div>
        <a-checkbox
          v-for="field in deviceFields"
          :key="field.key"
          v-model="field.checked"
          class="perm-check"
        >{{ field.label }}</a-checkbox>
      </div>

      <a-divider />

      <!-- 行级权限表达式 -->
      <div class="section-title">行级权限表达式</div>
      <div class="expr-section">
        <a-input
          v-model="form.rowExpression"
          placeholder="department_id = current_user.department_id"
          style="width: 500px"
        />
        <a-popover title="表达式说明" trigger="click">
          <a-button type="text" slot="trigger">语法帮助</a-button>
          <div class="expr-help">
            <div>支持字段：department_id, organization_id, creator_id, region</div>
            <div>运算符：=, !=, IN, NOT IN, LIKE</div>
            <div>示例：department_id = current_user.department_id</div>
            <div>示例：region IN ('华北', '华东')</div>
          </div>
        </a-popover>
      </div>
    </a-card>

    <!-- 空状态 -->
    <a-card class="content-card empty-state" v-else>
      <a-empty description="请选择角色以配置数据权限" />
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { getRoles, getDataPermissionRoles, updateDataPermissionRoles } from '@/api/security'

const selectedRoleId = ref(null)
const roles = ref([])

const form = reactive({
  dataScope: 'all',
  rowExpression: ''
})

const memberFields = reactive([
  { key: 'phone', label: '手机号', checked: true },
  { key: 'address', label: '地址', checked: true },
  { key: 'email', label: '邮箱', checked: true },
  { key: 'name', label: '姓名', checked: true },
  { key: 'id_card', label: '身份证', checked: false },
  { key: 'birthday', label: '生日', checked: false }
])

const deviceFields = reactive([
  { key: 'device_name', label: '设备名称', checked: true },
  { key: 'device_location', label: '设备位置', checked: true },
  { key: 'device_ip', label: 'IP地址', checked: false },
  { key: 'device_sn', label: '序列号', checked: false }
])

onMounted(async () => {
  try {
    const res = await getRoles()
    roles.value = res.data || res || []
  } catch (e) {
    console.error('加载角色列表失败', e)
  }
})

async function loadDataPermissions(roleId) {
  if (!roleId) return
  try {
    const res = await getDataPermissionRoles(roleId)
    const data = res.data || res
    if (data) {
      form.dataScope = data.data_scope || 'all'
      form.rowExpression = data.row_expression || ''
      const allowedFields = data.allowed_fields || []
      ;[memberFields, deviceFields].forEach(group => {
        group.forEach(f => {
          f.checked = allowedFields.includes(f.key)
        })
      })
    }
  } catch (e) {
    console.error('加载数据权限失败', e)
  }
}

async function handleSave() {
  if (!selectedRoleId.value) {
    Message.warning('请先选择角色')
    return
  }
  try {
    const allowedFields = [...memberFields, ...deviceFields]
      .filter(f => f.checked)
      .map(f => f.key)
    await updateDataPermissionRoles(selectedRoleId.value, {
      data_scope: form.dataScope,
      row_expression: form.rowExpression,
      allowed_fields: allowedFields
    })
    Message.success('保存成功')
  } catch (e) {
    Message.error('保存失败')
  }
}
</script>

<style scoped>
.page-container {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  height: 100%;
  box-sizing: border-box;
}

.filter-card {
  flex-shrink: 0;
}

.filter-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.content-card {
  flex: 1;
  overflow: auto;
}

.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 300px;
}

.section-title {
  font-size: 15px;
  font-weight: 600;
  margin-bottom: 12px;
  color: var(--color-text-1);
}

.scope-section {
  display: flex;
  gap: 24px;
}

.radio-label {
  font-weight: 500;
  margin-right: 4px;
}

.radio-desc {
  color: var(--color-text-3);
  font-size: 12px;
  margin-left: 4px;
}

.field-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.field-group-title {
  font-size: 13px;
  color: var(--color-text-2);
  margin-bottom: 4px;
}

.perm-check {
  margin-right: 20px;
}

.expr-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.expr-help {
  font-size: 13px;
  line-height: 1.8;
  color: var(--color-text-2);
}
</style>
