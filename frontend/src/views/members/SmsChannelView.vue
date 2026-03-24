<template>
  <div class="member-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员短信通道</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="通道总数" :value="stats.total || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="启用中" :value="stats.enabled || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="本月发送" :value="stats.sent || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="成功率" :value="stats.rate || 0" suffix="%" />
        </a-card>
      </a-col>
    </a-row>

    <a-card class="table-card">
      <a-space direction="vertical" :size="12" style="width: 100%">
        <a-space wrap>
          <a-input-search v-model="filters.keyword" placeholder="搜索通道名称" style="width: 200px" search-button @search="loadData" />
          <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 120px" @change="loadData">
            <a-option value="1">启用</a-option>
            <a-option value="0">禁用</a-option>
          </a-select>
          <a-button type="primary" @click="showCreate">新建</a-button>
          <a-button @click="loadData">刷新</a-button>
        </a-space>

        <a-table :columns="columns" :data="dataList" :loading="loading" :pagination="paginationConfig" @page-change="onPageChange" row-key="id" :scroll="{ x: 900 }">
          <template #type="{ record }">
            <a-tag :color="record.type === 'aliyun' ? 'blue' : record.type === 'tencent' ? 'red' : 'green'">
              {{ record.typeName }}
            </a-tag>
          </template>
          <template #status="{ record }">
            <a-tag :color="record.status === '1' ? 'green' : 'gray'">
              {{ record.status === '1' ? '启用' : '禁用' }}
            </a-tag>
          </template>
          <template #actions="{ record }">
            <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
            <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
          </template>
        </a-table>
      </a-space>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑通道' : '新建通道'" :width="520" @before-ok="handleSubmit" @cancel="formVisible = false">
      <a-form :model="form" layout="vertical">
        <a-form-item label="通道名称" required>
          <a-input v-model="form.name" placeholder="请输入通道名称" />
        </a-form-item>
        <a-form-item label="通道类型" required>
          <a-select v-model="form.type" placeholder="请选择通道类型">
            <a-option value="aliyun">阿里云</a-option>
            <a-option value="tencent">腾讯云</a-option>
            <a-option value="huawei">华为云</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="AppKey">
          <a-input v-model="form.appKey" placeholder="请输入AppKey" />
        </a-form-item>
        <a-form-item label="AppSecret">
          <a-input v-model="form.appSecret" placeholder="请输入AppSecret" type="password" />
        </a-form-item>
        <a-form-item label="状态">
          <a-radio-group v-model="form.status">
            <a-radio value="1">启用</a-radio>
            <a-radio value="0">禁用</a-radio>
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
.member-page { padding: 16px; }
.breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { text-align: center; }
.table-card { margin-top: 16px; }
</style>
