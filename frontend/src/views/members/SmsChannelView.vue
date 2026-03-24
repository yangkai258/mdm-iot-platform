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
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
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
const stats = ref({ total: 0, enabled: 0, sent: 0, rate: 0 });
const filters = reactive({ keyword: '', status: '' });
const paginationConfig = reactive({ current: 1, pageSize: 10, total: 0 });
const form = reactive({ id: null, name: '', type: 'aliyun', appKey: '', appSecret: '', status: '1' });
const formVisible = ref(false);
const isEdit = ref(false);

const columns = [
  { title: '通道名称', dataIndex: 'name', width: 180 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: 'AppKey', dataIndex: 'appKey', ellipsis: true },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' }
];

const mockData = () => [
  { id: 1, name: '阿里云主通道', type: 'aliyun', typeName: '阿里云', appKey: 'LTAI****************', status: '1' },
  { id: 2, name: '腾讯云备用', type: 'tencent', typeName: '腾讯云', appKey: '1400****************', status: '1' },
  { id: 3, name: '华为云营销', type: 'huawei', typeName: '华为云', appKey: 'HWA****************', status: '0' }
];

const loadData = () => {
  loading.value = true;
  setTimeout(() => {
    const data = mockData();
    dataList.value = data.filter(item => {
      if (filters.keyword && !item.name.includes(filters.keyword)) return false;
      if (filters.status && item.status !== filters.status) return false;
      return true;
    });
    paginationConfig.total = dataList.value.length;
    stats.value = { total: 3, enabled: 2, sent: 125680, rate: 98.5 };
    loading.value = false;
  }, 300);
};

const onPageChange = (page) => {
  paginationConfig.current = page;
  loadData();
};

const showCreate = () => {
  isEdit.value = false;
  Object.assign(form, { id: null, name: '', type: 'aliyun', appKey: '', appSecret: '', status: '1' });
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
  if (!form.name) {
    Message.error('请输入通道名称');
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
