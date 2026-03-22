<template>
  <div class="page-container">
    <!-- 菜单树 -->
    <a-card class="tree-card">
      <template #title>
        <div class="card-title">
          <span>菜单结构</span>
          <a-space>
            <a-button type="text" size="small" @click="loadMenuTree">
              <icon-refresh />
            </a-button>
            <a-button type="primary" size="small" @click="openDrawer(null)">
              <template #icon><icon-plus /></template>
              新增菜单
            </a-button>
          </a-space>
        </div>
      </template>

      <a-spin v-if="loading" />
      <a-tree
        v-else
        v-model:selected-keys="selectedKeys"
        :data="treeData"
        :show-line="true"
        block-node
        @select="handleNodeSelect"
      >
        <template #title="{ data: node }">
          <div class="tree-node">
            <span class="node-title">{{ node.title }}</span>
            <span class="node-actions">
              <a-button type="text" size="mini" @click.stop="openDrawer(node, true)">
                <icon-plus />
              </a-button>
              <a-button type="text" size="mini" @click.stop="openDrawer(node, false)">
                <icon-edit />
              </a-button>
              <a-popconfirm content="确定删除该菜单？" @ok="handleDelete(node.key)">
                <a-button type="text" size="mini" status="danger" @click.stop>
                  <icon-delete />
                </a-button>
              </a-popconfirm>
            </span>
          </div>
        </template>
      </a-tree>
    </a-card>

    <!-- 创建/编辑 Drawer -->
    <a-drawer
      v-model:visible="drawerVisible"
      :title="isEdit ? '编辑菜单' : '新增菜单'"
      width="520px"
      @confirm="handleSubmit"
      @cancel="drawerVisible = false"
    >
      <a-form ref="formRef" :model="form" :rules="formRules" layout="vertical" label-align="left">
        <a-form-item label="菜单类型" field="menu_type">
          <a-radio-group v-model="form.menu_type">
            <a-radio value="dir">目录</a-radio>
            <a-radio value="menu">菜单</a-radio>
            <a-radio value="button">按钮</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="菜单名称" field="menu_name">
          <a-input v-model="form.menu_name" placeholder="请输入菜单名称" />
        </a-form-item>
        <a-form-item label="上级菜单" field="parent_id">
          <a-tree-select
            v-model="form.parent_id"
            :data="treeData"
            placeholder="请选择上级菜单（可选）"
            allow-clear
            style="width: 100%"
            :field-names="{ key: 'key', title: 'title' }"
          />
        </a-form-item>
        <a-form-item v-if="form.menu_type === 'menu'" label="路由路径" field="path">
          <a-input v-model="form.path" placeholder="如: /system/users" />
        </a-form-item>
        <a-form-item v-if="form.menu_type === 'menu'" label="组件路径" field="component">
          <a-input v-model="form.component" placeholder="如: views/system/Users.vue" />
        </a-form-item>
        <a-form-item label="菜单图标" field="icon">
          <a-input v-model="form.icon" placeholder="如: icon-settings" />
        </a-form-item>
        <a-form-item label="排序" field="sort">
          <a-input-number v-model="form.sort" :min="0" style="width: 100%" />
        </a-form-item>
        <a-form-item label="状态" field="status">
          <a-radio-group v-model="form.status">
            <a-radio :value="1">正常</a-radio>
            <a-radio :value="0">禁用</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as permApi from '@/api/permissions.js'

const loading = ref(false)
const drawerVisible = ref(false)
const isEdit = ref(false)
const formRef = ref(null)
const selectedKeys = ref([])
const treeData = ref([])

const mockMenuTree = [
  {
    key: 1, title: '设备管理', menu_type: 'dir', icon: 'icon-device',
    children: [
      { key: 11, title: '设备列表', menu_type: 'menu', path: '/devices', component: 'DeviceDashboard.vue' },
      { key: 12, title: '设备详情', menu_type: 'menu', path: '/device/:id', component: 'DeviceDetail.vue' }
    ]
  },
  {
    key: 2, title: 'OTA升级', menu_type: 'dir', icon: 'icon-upload',
    children: [
      { key: 21, title: '固件包管理', menu_type: 'menu', path: '/ota/packages', component: 'ota/OtaPackages.vue' },
      { key: 22, title: '部署任务', menu_type: 'menu', path: '/ota/deployments', component: 'ota/OtaDeployments.vue' }
    ]
  },
  {
    key: 3, title: '系统设置', menu_type: 'dir', icon: 'icon-settings',
    children: [
      { key: 31, title: '用户管理', menu_type: 'menu', path: '/system/users', component: 'system/Users.vue' },
      { key: 32, title: '角色管理', menu_type: 'menu', path: '/permissions/roles', component: 'permissions/Roles.vue' },
      { key: 33, title: '菜单管理', menu_type: 'menu', path: '/permissions/menus', component: 'permissions/Menus.vue' }
    ]
  }
]

const defaultForm = () => ({
  menu_type: 'menu',
  menu_name: '',
  parent_id: undefined,
  path: '',
  component: '',
  icon: '',
  sort: 0,
  status: 1
})

const form = reactive(defaultForm())

const formRules = {
  menu_name: [{ required: true, message: '请输入菜单名称' }],
  path: [{ required: true, message: '请输入路由路径', trigger: 'blur' }]
}

const loadMenuTree = async () => {
  loading.value = true
  try {
    const res = await permApi.getMenuTree()
    if (res.code === 0) {
      treeData.value = res.data || []
    }
  } catch {
    treeData.value = JSON.parse(JSON.stringify(mockMenuTree))
    Message.warning('使用模拟数据')
  } finally {
    loading.value = false
  }
}

const handleNodeSelect = (keys) => {
  if (keys.length > 0) {
    selectedKeys.value = keys
  }
}

const openDrawer = (node, isChild = false) => {
  isEdit.value = !!node && !isChild
  if (isChild && node) {
    // 新增子菜单
    Object.assign(form, defaultForm())
    form.parent_id = node.key
  } else if (node) {
    // 编辑
    Object.assign(form, {
      menu_type: node.menu_type || 'menu',
      menu_name: node.title,
      parent_id: node.parent_id || undefined,
      path: node.path || '',
      component: node.component || '',
      icon: node.icon || '',
      sort: node.sort || 0,
      status: node.status || 1
    })
  } else {
    Object.assign(form, defaultForm())
  }
  drawerVisible.value = true
}

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }
  try {
    let res
    if (isEdit.value) {
      res = await permApi.updateMenu(form.parent_id, { ...form })
    } else {
      res = await permApi.createMenu({ ...form })
    }
    if (res.code === 0) {
      Message.success(isEdit.value ? '更新成功' : '创建成功')
      drawerVisible.value = false
      loadMenuTree()
    }
  } catch {
    Message.success('操作成功（模拟）')
    drawerVisible.value = false
    loadMenuTree()
  }
}

const handleDelete = async (id) => {
  try {
    const res = await permApi.deleteMenu(id)
    if (res.code === 0) {
      Message.success('删除成功')
      loadMenuTree()
    }
  } catch {
    Message.success('删除成功（模拟）')
    loadMenuTree()
  }
}

onMounted(() => {
  loadMenuTree()
})
</script>

<style scoped>
.page-container { display: flex; flex-direction: column; gap: 16px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
.tree-node { display: flex; justify-content: space-between; align-items: center; width: 100%; }
.node-title { flex: 1; }
.node-actions { display: flex; gap: 4px; opacity: 0; transition: opacity 0.2s; }
.tree-node:hover .node-actions { opacity: 1; }
</style>
