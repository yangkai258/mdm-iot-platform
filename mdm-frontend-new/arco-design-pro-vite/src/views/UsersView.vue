<template>
  <div class="users-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>用户管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新增用户
          </a-button>
        </div>
      </template>
      
      <div class="search-area">
        <a-row :gutter="16">
          <a-col :span="6">
            <a-input v-model="searchForm.keyword" placeholder="搜索用户名/手机/邮箱" allow-clear />
          </a-col>
          <a-col :span="4">
            <a-select v-model="searchForm.role" placeholder="角色" allow-clear>
              <a-option value="admin">管理员</a-option>
              <a-option value="operator">运营</a-option>
              <a-option value="user">普通用户</a-option>
            </a-select>
          </a-col>
          <a-col :span="4">
            <a-select v-model="searchForm.status" placeholder="状态" allow-clear>
              <a-option value="active">正常</a-option>
              <a-option value="disabled">禁用</a-option>
            </a-select>
          </a-col>
          <a-col :span="2">
            <a-button type="primary" @click="handleSearch">筛选</a-button>
          </a-col>
        </a-row>
      </div>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="handlePageChange">
        <template #avatar="{ record }">
          <a-avatar :size="32" :style="{ backgroundColor: record.avatar ? 'transparent' : '#165DFF' }">
            <img v-if="record.avatar" :src="record.avatar" />
            <span v-else>{{ record.username?.charAt(0)?.toUpperCase() }}</span>
          </a-avatar>
        </template>
        <template #role="{ record }">
          <a-tag :color="getRoleColor(record.role)">{{ record.roleText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 'active' ? 'green' : 'red'">
            {{ record.status === 'active' ? '正常' : '禁用' }}
          </a-tag>
        </template>
        <template #lastLogin="{ record }">
          {{ record.lastLogin || '-' }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleView(record)">详情</a-link>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link @click="handleResetPassword(record)">重置密码</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 用户详情弹窗 -->
    <a-modal v-model:visible="detailVisible" title="用户详情" :width="600" :footer="null">
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="用户ID">{{ currentUser.id }}</a-descriptions-item>
        <a-descriptions-item label="用户名">{{ currentUser.username }}</a-descriptions-item>
        <a-descriptions-item label="手机号">{{ currentUser.phone }}</a-descriptions-item>
        <a-descriptions-item label="邮箱">{{ currentUser.email }}</a-descriptions-item>
        <a-descriptions-item label="角色">
          <a-tag :color="getRoleColor(currentUser.role)">{{ currentUser.roleText }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="currentUser.status === 'active' ? 'green' : 'red'">
            {{ currentUser.status === 'active' ? '正常' : '禁用' }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="最后登录">{{ currentUser.lastLogin || '-' }}</a-descriptions-item>
        <a-descriptions-item label="创建时间">{{ currentUser.createdAt }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>

    <!-- 编辑弹窗 -->
    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑用户' : '新增用户'" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="用户名" required>
          <a-input v-model="form.username" placeholder="请输入用户名" />
        </a-form-item>
        <a-form-item v-if="!isEdit" label="密码" required>
          <a-input-password v-model="form.password" placeholder="请输入密码" />
        </a-form-item>
        <a-form-item label="手机号">
          <a-input v-model="form.phone" placeholder="请输入手机号" />
        </a-form-item>
        <a-form-item label="邮箱">
          <a-input v-model="form.email" placeholder="请输入邮箱" />
        </a-form-item>
        <a-form-item label="角色">
          <a-select v-model="form.role" placeholder="选择角色">
            <a-option value="admin">管理员</a-option>
            <a-option value="operator">运营</a-option>
            <a-option value="user">普通用户</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const data = ref([
  { id: 1, username: 'admin', phone: '138****1234', email: 'admin@example.com', role: 'admin', roleText: '管理员', status: 'active', lastLogin: '2026-03-28 18:00:00', createdAt: '2026-01-01 00:00:00' },
  { id: 2, username: 'operator1', phone: '139****5678', email: 'op1@example.com', role: 'operator', roleText: '运营', status: 'active', lastLogin: '2026-03-28 10:30:00', createdAt: '2026-02-15 09:00:00' },
  { id: 3, username: 'zhangsan', phone: '137****9012', email: 'zhang@example.com', role: 'user', roleText: '普通用户', status: 'active', lastLogin: '2026-03-27 16:00:00', createdAt: '2026-02-20 14:00:00' },
  { id: 4, username: 'lisi', phone: '136****3456', email: 'li@example.com', role: 'user', roleText: '普通用户', status: 'disabled', lastLogin: '2026-03-20 10:00:00', createdAt: '2026-03-01 11:00:00' },
]);

const searchForm = reactive({ keyword: '', role: '', status: '' });
const pagination = reactive({ current: 1, pageSize: 20, total: 4 });

const columns = [
  { title: '', slotName: 'avatar', width: 60 },
  { title: '用户ID', dataIndex: 'id', width: 80 },
  { title: '用户名', dataIndex: 'username', width: 120 },
  { title: '手机号', dataIndex: 'phone', width: 130 },
  { title: '邮箱', dataIndex: 'email', width: 180 },
  { title: '角色', slotName: 'role', width: 100 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '最后登录', slotName: 'lastLogin', width: 160 },
  { title: '操作', slotName: 'actions', width: 220, fixed: 'right' },
];

const detailVisible = ref(false);
const editVisible = ref(false);
const isEdit = ref(false);
const currentUser = ref<any>({});

const form = reactive({
  username: '',
  password: '',
  phone: '',
  email: '',
  role: 'user',
});

const getRoleColor = (role: string) => {
  const map: Record<string, string> = { admin: 'red', operator: 'blue', user: 'green' };
  return map[role] || 'default';
};

const handleSearch = () => {};
const handlePageChange = (page: number) => { pagination.current = page; };
const handleCreate = () => {
  isEdit.value = false;
  Object.assign(form, { username: '', password: '', phone: '', email: '', role: 'user' });
  editVisible.value = true;
};
const handleView = (record: any) => {
  currentUser.value = record;
  detailVisible.value = true;
};
const handleEdit = (record: any) => {
  isEdit.value = true;
  Object.assign(form, record);
  editVisible.value = true;
};
const handleDelete = (record: any) => {};
const handleResetPassword = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.users-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
.search-area { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
</style>
