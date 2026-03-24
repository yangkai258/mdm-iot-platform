<template>
  <div class="member-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员权益管理</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="权益总数" :value="stats.total || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="启用中" :value="stats.enabled || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="关联会员" :value="stats.members || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="本月使用" :value="stats.usedThisMonth || 0" />
        </a-card>
      </a-col>
    </a-row>

    <a-card class="table-card">
      <a-space direction="vertical" :size="12" style="width: 100%">
        <a-space wrap>
          <a-input-search v-model="filters.keyword" placeholder="搜索权益名称" style="width: 200px" search-button @search="loadData" />
          <a-select v-model="filters.level" placeholder="适用等级" allow-clear style="width: 140px" @change="loadData">
            <a-option value="silver">银卡</a-option>
            <a-option value="gold">金卡</a-option>
            <a-option value="platinum">铂金卡</a-option>
            <a-option value="diamond">钻石卡</a-option>
          </a-select>
          <a-button type="primary" @click="showCreate">新建</a-button>
          <a-button @click="loadData">刷新</a-button>
        </a-space>

        <a-table :columns="columns" :data="dataList" :loading="loading" :pagination="paginationConfig" @page-change="onPageChange" row-key="id" :scroll="{ x: 1000 }">
          <template #level="{ record }">
            <a-tag :color="record.level === 'diamond' ? 'purple' : record.level === 'platinum' ? 'gold' : record.level === 'gold' ? 'orange' : 'gray'">
              {{ record.levelName }}
            </a-tag>
          </template>
          <template #status="{ record }">
            <a-switch :model-value="record.status" checked-value="1" unchecked-value="0" @change="toggleStatus(record)" />
          </template>
          <template #actions="{ record }">
            <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
            <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
          </template>
        </a-table>
      </a-space>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑权益' : '新建权益'" :width="520" @before-ok="handleSubmit" @cancel="formVisible = false">
      <a-form :model="form" layout="vertical">
        <a-form-item label="权益名称" required>
          <a-input v-model="form.name" placeholder="请输入权益名称" />
        </a-form-item>
        <a-form-item label="适用等级" required>
          <a-select v-model="form.level" placeholder="请选择适用等级">
            <a-option value="silver">银卡</a-option>
            <a-option value="gold">金卡</a-option>
            <a-option value="platinum">铂金卡</a-option>
            <a-option value="diamond">钻石卡</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="权益内容">
          <a-textarea v-model="form.content" placeholder="请输入权益内容描述" :rows="3" />
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
const stats = ref({ total: 0, enabled: 0, members: 0, usedThisMonth: 0 });
const filters = reactive({ keyword: '', level: '' });
const paginationConfig = reactive({ current: 1, pageSize: 10, total: 0 });
const form = reactive({ id: null, name: '', level: 'silver', content: '', status: '1' });
const formVisible = ref(false);
const isEdit = ref(false);

const columns = [
  { title: '权益名称', dataIndex: 'name', width: 180 },
  { title: '适用等级', slotName: 'level', width: 100 },
  { title: '权益内容', dataIndex: 'content', ellipsis: true },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' }
];

const mockData = () => [
  { id: 1, name: '生日专属折扣', level: 'silver', levelName: '银卡', content: '生日当月享受9折优惠', status: '1' },
  { id: 2, name: '免费停车', level: 'gold', levelName: '金卡', content: '每月免费停车2小时', status: '1' },
  { id: 3, name: 'VIP客服', level: 'platinum', levelName: '铂金卡', content: '专属客服热线优先接入', status: '1' },
  { id: 4, name: '免费洗车', level: 'diamond', levelName: '钻石卡', content: '每月免费洗车1次', status: '1' }
];

const loadData = () => {
  loading.value = true;
  setTimeout(() => {
    const data = mockData();
    dataList.value = data.filter(item => {
      if (filters.keyword && !item.name.includes(filters.keyword)) return false;
      if (filters.level && item.level !== filters.level) return false;
      return true;
    });
    paginationConfig.total = dataList.value.length;
    stats.value = { total: 8, enabled: 6, members: 1234, usedThisMonth: 456 };
    loading.value = false;
  }, 300);
};

const onPageChange = (page) => {
  paginationConfig.current = page;
  loadData();
};

const showCreate = () => {
  isEdit.value = false;
  Object.assign(form, { id: null, name: '', level: 'silver', content: '', status: '1' });
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

const toggleStatus = (record) => {
  record.status = record.status === '1' ? '0' : '1';
  Message.success('状态已更新');
};

const handleSubmit = (done) => {
  if (!form.name) {
    Message.error('请输入权益名称');
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
