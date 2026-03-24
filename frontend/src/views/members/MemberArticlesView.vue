<template>
  <div class="member-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员推文流水</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="文章总数" :value="stats.total || 0" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="今日发布" :value="stats.today || 0" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="stat-card">
          <a-statistic title="总阅读量" :value="stats.views || 0" />
        </a-card>
      </a-col>
    </a-row>

    <a-card class="table-card">
      <a-space direction="vertical" :size="12" style="width: 100%">
        <a-space wrap>
          <a-input-search v-model="filters.keyword" placeholder="搜索文章标题" style="width: 240px" search-button @search="loadData" />
          <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 120px" @change="loadData">
            <a-option value="1">已发布</a-option>
            <a-option value="0">草稿</a-option>
          </a-select>
          <a-button type="primary" @click="showCreate">新建</a-button>
          <a-button @click="loadData">刷新</a-button>
        </a-space>

        <a-table :columns="columns" :data="dataList" :loading="loading" :pagination="paginationConfig" @page-change="onPageChange" row-key="id" :scroll="{ x: 1200 }">
          <template #status="{ record }">
            <a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '已发布' : '草稿' }}</a-tag>
          </template>
          <template #actions="{ record }">
            <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
            <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
          </template>
        </a-table>
      </a-space>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑推文' : '新建推文'" :width="520" @before-ok="handleSubmit" @cancel="formVisible = false">
      <a-form :model="form" layout="vertical">
        <a-form-item label="文章标题" required>
          <a-input v-model="form.title" placeholder="请输入文章标题" />
        </a-form-item>
        <a-form-item label="内容摘要">
          <a-textarea v-model="form.summary" placeholder="请输入内容摘要" :rows="3" />
        </a-form-item>
        <a-form-item label="发布状态">
          <a-radio-group v-model="form.status">
            <a-radio :value="1">已发布</a-radio>
            <a-radio :value="0">草稿</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
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
.member-page { padding: 16px; }
.breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { text-align: center; }
.table-card { margin-top: 16px; }
</style>
