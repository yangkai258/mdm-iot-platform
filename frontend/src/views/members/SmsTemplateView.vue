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
const stats = ref({ total: 0, marketing: 0, notice: 0, verify: 0 });
const filters = reactive({ keyword: '', type: '' });
const paginationConfig = reactive({ current: 1, pageSize: 10, total: 0 });
const form = reactive({ id: null, name: '', type: 'notice', content: '' });
const formVisible = ref(false);
const isEdit = ref(false);

const columns = [
  { title: '模板名称', dataIndex: 'name', width: 180 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '模板内容', dataIndex: 'content', ellipsis: true },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' }
];

const mockData = () => [
  { id: 1, name: '注册验证码', type: 'verify', content: '您的验证码是{code}，5分钟内有效' },
  { id: 2, name: '生日祝福', type: 'marketing', content: '亲爱的{name}，祝您生日快乐！' },
  { id: 3, name: '订单通知', type: 'notice', content: '您的订单{orderno}已发货' },
  { id: 4, name: '积分到期', type: 'marketing', content: '您的{points}积分即将过期' }
];

const loadData = () => {
  loading.value = true;
  setTimeout(() => {
    const data = mockData();
    dataList.value = data.filter(item => {
      if (filters.keyword && !item.name.includes(filters.keyword)) return false;
      if (filters.type && item.type !== filters.type) return false;
      return true;
    });
    paginationConfig.total = dataList.value.length;
    stats.value = { total: 15, marketing: 6, notice: 5, verify: 4 };
    loading.value = false;
  }, 300);
};

const onPageChange = (page) => {
  paginationConfig.current = page;
  loadData();
};

const showCreate = () => {
  isEdit.value = false;
  Object.assign(form, { id: null, name: '', type: 'notice', content: '' });
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
  if (!form.name || !form.content) {
    Message.error('请填写完整信息');
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
