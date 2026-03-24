<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="名称"><a-input v-model="form.name" placeholder="请输入" /></a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleCreate">新建</a-button>
    </div>
    <a-table :columns="columns" :data="dataList" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
      <template #actions="{ record }">
        <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
        <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
      </template>
    </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @before-ok="handleSubmit" @cancel="modalVisible = false">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="名称"><a-input v-model="form.name" placeholder="请输入" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>

import { ref, reactive, onMounted } from 'vue';
import { Message } from '@arco-design/web-vue';

const loading = ref(false);
const dataList = ref([]);
const stats = ref({ total: 0, today: 0, views: 0 });
const filters = reactive({ keyword: '', status: '' });
const paginationConfig = reactive({ current: 1, pageSize: 10, total: 0 });
const form = reactive({ id: null, title: '', summary: '', status: 1 });
const formVisible = ref(false);
const isEdit = ref(false);

const columns = [
  { title: '文章标题', dataIndex: 'title', width: 220 },
  { title: '内容摘要', dataIndex: 'summary', ellipsis: true },
  { title: '发布时间', dataIndex: 'publishTime', width: 180 },
  { title: '阅读量', dataIndex: 'views', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' }
];

const mockData = () => [
  { id: 1, title: '新会员专享福利来袭', summary: '新注册会员可领取专属优惠券', publishTime: '2026-03-20 10:00:00', views: 2345, status: 1 },
  { id: 2, title: '生日月双倍福利', summary: '3月生日会员享受双倍积分', publishTime: '2026-03-15 10:00:00', views: 1890, status: 1 },
  { id: 3, title: '新品上市抢先看', summary: '春季新品预告，会员优先购买', publishTime: '2026-03-10 10:00:00', views: 3200, status: 1 }
];

const loadData = () => {
  loading.value = true;
  setTimeout(() => {
    const data = mockData();
    dataList.value = data.filter(item => {
      if (filters.keyword && !item.title.includes(filters.keyword)) return false;
      if (filters.status && String(item.status) !== filters.status) return false;
      return true;
    });
    paginationConfig.total = dataList.value.length;
    stats.value = { total: 12, today: 2, views: 15234 };
    loading.value = false;
  }, 300);
};

const onPageChange = (page) => {
  paginationConfig.current = page;
  loadData();
};

const showCreate = () => {
  isEdit.value = false;
  Object.assign(form, { id: null, title: '', summary: '', status: 1 });
  formVisible.value = true;
};

const handleEdit = (record) => {
  isEdit.value = true;
  Object.assign(form, record);
  formVisible.value = true;
};

const handleDelete = () => {
  Message.success('删除成功');
  loadData();
};

const handleSubmit = (done) => {
  if (!form.title) {
    Message.error('请输入文章标题');
    done(false);
    return;
  }
  setTimeout(() => {
    Message.success(isEdit.value ? '更新成功' : '创建成功');
    formVisible.value = false;
    loadData();
    done(true);
  }, 400);
};

onMounted(() => {
  loadData();
});

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
