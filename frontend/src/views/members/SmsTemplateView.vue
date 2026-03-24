<template>
  <div class="member-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>短信模板设置</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="模板总数" :value="stats.total || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="营销类" :value="stats.marketing || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="通知类" :value="stats.notice || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="验证码类" :value="stats.verify || 0" />
        </a-card>
      </a-col>
    </a-row>

    <a-card class="table-card">
      <a-space direction="vertical" :size="12" style="width: 100%">
        <a-space wrap>
          <a-input-search v-model="filters.keyword" placeholder="搜索模板名称" style="width: 220px" search-button @search="loadData" />
          <a-select v-model="filters.type" placeholder="模板类型" allow-clear style="width: 140px" @change="loadData">
            <a-option value="marketing">营销类</a-option>
            <a-option value="notice">通知类</a-option>
            <a-option value="verify">验证码类</a-option>
          </a-select>
          <a-button type="primary" @click="showCreate">新建</a-button>
          <a-button @click="loadData">刷新</a-button>
        </a-space>

        <a-table :columns="columns" :data="dataList" :loading="loading" :pagination="paginationConfig" @page-change="onPageChange" row-key="id" :scroll="{ x: 900 }">
          <template #type="{ record }">
            <a-tag :color="record.type === 'marketing' ? 'orange' : record.type === 'notice' ? 'blue' : 'green'">
              {{ record.type === 'marketing' ? '营销类' : record.type === 'notice' ? '通知类' : '验证码类' }}
            </a-tag>
          </template>
          <template #actions="{ record }">
            <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
            <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
          </template>
        </a-table>
      </a-space>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑模板' : '新建模板'" :width="520" @before-ok="handleSubmit" @cancel="formVisible = false">
      <a-form :model="form" layout="vertical">
        <a-form-item label="模板名称" required>
          <a-input v-model="form.name" placeholder="请输入模板名称" />
        </a-form-item>
        <a-form-item label="模板类型" required>
          <a-select v-model="form.type" placeholder="请选择模板类型">
            <a-option value="marketing">营销类</a-option>
            <a-option value="notice">通知类</a-option>
            <a-option value="verify">验证码类</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="模板内容" required>
          <a-textarea v-model="form.content" placeholder="请输入模板内容，使用{var}表示变量" :rows="4" />
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
.member-page { padding: 16px; }
.breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { text-align: center; }
.table-card { margin-top: 16px; }
</style>
