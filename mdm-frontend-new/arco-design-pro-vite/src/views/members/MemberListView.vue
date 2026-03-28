<template>
  <div class="container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>会员管理</span>
          <a-button type="primary" @click="handleCreate">新增会员</a-button>
        </div>
      </template>
      
      <a-table 
        :columns="columns" 
        :data="data" 
        :loading="loading"
        :pagination="{ total: total, current: current, pageSize: pageSize }"
        @page-change="handlePageChange"
      >
        <template #status="{ record }">
          <a-tag :color="record.status === 'active' ? 'green' : 'red'">
            {{ record.status === 'active' ? '活跃' : '不活跃' }}
          </a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link @click="handleView(record)">查看</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';

const loading = ref(false);
const data = ref([]);
const total = ref(0);
const current = ref(1);
const pageSize = ref(20);

const columns = [
  { title: 'ID', dataIndex: 'id' },
  { title: '用户名', dataIndex: 'username' },
  { title: '手机号', dataIndex: 'phone' },
  { title: '等级', dataIndex: 'level' },
  { title: '积分', dataIndex: 'points' },
  { title: '状态', dataIndex: 'status', slotName: 'status' },
  { title: '注册时间', dataIndex: 'created_at' },
  { title: '操作', slotName: 'actions' },
];

const fetchData = async () => {
  loading.value = true;
  try {
    const response = await fetch('/api/v1/members', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token') || ''}`
      }
    });
    const result = await response.json();
    if (result.code === 0) {
      data.value = result.data?.list || [];
      total.value = result.data?.total || 0;
    }
  } catch (error) {
    console.error('获取会员列表失败:', error);
  } finally {
    loading.value = false;
  }
};

const handlePageChange = (page: number) => {
  current.value = page;
  fetchData();
};

const handleCreate = () => {
  // TODO: 打开创建会员弹窗
};

const handleEdit = (record: any) => {
  // TODO: 打开编辑会员弹窗
};

const handleView = (record: any) => {
  // TODO: 跳转到会员详情页
};

const handleDelete = (record: any) => {
  // TODO: 删除会员
};

onMounted(() => {
  fetchData();
});
</script>

<style scoped>
.container {
  padding: 20px;
}
.card-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
