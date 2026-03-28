<template>
  <div class="member-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>会员管理</span>
          <a-space>
            <a-button type="primary" @click="handleCreate">
              <template #icon><icon-plus /></template>
              新增会员
            </a-button>
            <a-button @click="handleExport">
              <template #icon><icon-download /></template>
              导出
            </a-button>
          </a-space>
        </div>
      </template>
      
      <!-- 搜索筛选区 -->
      <div class="search-area">
        <a-row :gutter="16">
          <a-col :span="6">
            <a-input v-model="searchForm.keyword" placeholder="搜索会员ID/手机号/昵称" allow-clear @press-enter="handleSearch">
              <template #prefix><icon-search /></template>
            </a-input>
          </a-col>
          <a-col :span="4">
            <a-select v-model="searchForm.level" placeholder="会员等级" allow-clear>
              <a-option value="gold">黄金会员</a-option>
              <a-option value="silver">白银会员</a-option>
              <a-option value="bronze">青铜会员</a-option>
            </a-select>
          </a-col>
          <a-col :span="4">
            <a-select v-model="searchForm.status" placeholder="状态" allow-clear>
              <a-option value="active">活跃</a-option>
              <a-option value="inactive">不活跃</a-option>
              <a-option value="frozen">已冻结</a-option>
            </a-select>
          </a-col>
          <a-col :span="4">
            <a-date-picker v-model="searchForm.startDate" placeholder="注册开始日期" style="width: 100%" />
          </a-col>
          <a-col :span="4">
            <a-date-picker v-model="searchForm.endDate" placeholder="注册结束日期" style="width: 100%" />
          </a-col>
          <a-col :span="2">
            <a-button type="primary" @click="handleSearch">
              <template #icon><icon-search /></template>
              筛选
            </a-button>
          </a-col>
        </a-row>
      </div>
      
      <!-- 数据表格 -->
      <a-table 
        :columns="columns" 
        :data="data" 
        :loading="loading"
        :pagination="pagination"
        :row-selection="{ type: 'checkbox', showCheckedAll: true }"
        @page-change="handlePageChange"
        @page-size-change="handlePageSizeChange"
      >
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">
            {{ getStatusText(record.status) }}
          </a-tag>
        </template>
        <template #level="{ record }">
          <a-tag :color="getLevelColor(record.level)">
            {{ getLevelText(record.level) }}
          </a-tag>
        </template>
        <template #avatar="{ record }">
          <a-avatar :size="32" :style="{ backgroundColor: record.color || '#165DFF' }">
            {{ record.name?.charAt(0)?.toUpperCase() || 'U' }}
          </a-avatar>
        </template>
        <template #points="{ record }">
          <span style="color: #FF7D00; font-weight: 500;">{{ record.points?.toLocaleString() || 0 }}</span>
        </template>
        <template #consumption="{ record }">
          <span>¥{{ record.totalConsumption?.toLocaleString() || 0 }}</span>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleView(record)">详情</a-link>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link @click="handlePoints(record)">积分</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 会员详情弹窗 -->
    <a-modal v-model:visible="detailVisible" title="会员详情" :width="700" :footer="null">
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="会员ID">{{ currentMember.id }}</a-descriptions-item>
        <a-descriptions-item label="昵称">{{ currentMember.name }}</a-descriptions-item>
        <a-descriptions-item label="手机号">{{ currentMember.phone }}</a-descriptions-item>
        <a-descriptions-item label="会员等级">
          <a-tag :color="getLevelColor(currentMember.level)">{{ getLevelText(currentMember.level) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="积分">
          <span style="color: #FF7D00;">{{ currentMember.points?.toLocaleString() || 0 }}</span>
        </a-descriptions-item>
        <a-descriptions-item label="累计消费">
          ¥{{ currentMember.totalConsumption?.toLocaleString() || 0 }}
        </a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="getStatusColor(currentMember.status)">{{ getStatusText(currentMember.status) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="注册时间">{{ currentMember.createdAt }}</a-descriptions-item>
        <a-descriptions-item label="最后活跃" :span="2">{{ currentMember.lastActiveAt }}</a-descriptions-item>
        <a-descriptions-item label="邮箱" :span="2">{{ currentMember.email || '-' }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>

    <!-- 新增/编辑弹窗 -->
    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑会员' : '新增会员'" :width="500" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="昵称" required>
          <a-input v-model="form.name" placeholder="请输入会员昵称" />
        </a-form-item>
        <a-form-item label="手机号" required>
          <a-input v-model="form.phone" placeholder="请输入手机号" />
        </a-form-item>
        <a-form-item label="邮箱">
          <a-input v-model="form.email" placeholder="请输入邮箱" />
        </a-form-item>
        <a-form-item label="会员等级">
          <a-select v-model="form.level" placeholder="请选择会员等级">
            <a-option value="gold">黄金会员</a-option>
            <a-option value="silver">白银会员</a-option>
            <a-option value="bronze">青铜会员</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const data = ref([]);
const detailVisible = ref(false);
const formVisible = ref(false);
const isEdit = ref(false);
const currentMember = ref<any>({});

const searchForm = reactive({
  keyword: '',
  level: '',
  status: '',
  startDate: '',
  endDate: '',
});

const form = reactive({
  name: '',
  phone: '',
  email: '',
  level: 'bronze',
});

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
});

const columns = [
  { title: '', slotName: 'avatar', width: 50 },
  { title: '会员ID', dataIndex: 'id', width: 100 },
  { title: '昵称', dataIndex: 'name', width: 120 },
  { title: '手机号', dataIndex: 'phone', width: 130 },
  { title: '等级', slotName: 'level', width: 100 },
  { title: '积分', slotName: 'points', width: 100, sortable: true },
  { title: '累计消费', slotName: 'consumption', width: 120, sortable: true },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '注册时间', dataIndex: 'createdAt', width: 160 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' },
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
      data.value = result.data?.list || generateMockData();
      pagination.total = result.data?.total || data.value.length;
    } else {
      data.value = generateMockData();
      pagination.total = data.value.length;
    }
  } catch (error) {
    console.error('获取会员列表失败:', error);
    data.value = generateMockData();
    pagination.total = data.value.length;
  } finally {
    loading.value = false;
  }
};

const generateMockData = () => {
  return [
    { id: 'M001', name: '张三', phone: '138****1234', level: 'gold', points: 15800, totalConsumption: 25800, status: 'active', createdAt: '2026-03-15 10:30:00', lastActiveAt: '2026-03-28 09:00:00', color: '#165DFF' },
    { id: 'M002', name: '李四', phone: '139****5678', level: 'silver', points: 8600, totalConsumption: 12600, status: 'active', createdAt: '2026-03-10 14:20:00', lastActiveAt: '2026-03-27 18:30:00', color: '#00B42A' },
    { id: 'M003', name: '王五', phone: '137****9012', level: 'bronze', points: 3200, totalConsumption: 4800, status: 'inactive', createdAt: '2026-02-28 09:15:00', lastActiveAt: '2026-03-20 11:00:00', color: '#FF7D00' },
    { id: 'M004', name: '赵六', phone: '136****3456', level: 'gold', points: 22500, totalConsumption: 42000, status: 'active', createdAt: '2026-01-15 16:45:00', lastActiveAt: '2026-03-28 08:00:00', color: '#722ED1' },
  ];
};

const handleSearch = () => {
  pagination.current = 1;
  fetchData();
};

const handlePageChange = (page: number) => {
  pagination.current = page;
  fetchData();
};

const handlePageSizeChange = (size: number) => {
  pagination.pageSize = size;
  fetchData();
};

const handleCreate = () => {
  isEdit.value = false;
  Object.assign(form, { name: '', phone: '', email: '', level: 'bronze' });
  formVisible.value = true;
};

const handleEdit = (record: any) => {
  isEdit.value = true;
  Object.assign(form, record);
  currentMember.value = record;
  formVisible.value = true;
};

const handleView = (record: any) => {
  currentMember.value = record;
  detailVisible.value = true;
};

const handlePoints = (record: any) => {
  currentMember.value = record;
  // TODO: 打开积分管理弹窗
};

const handleDelete = (record: any) => {
  // TODO: 删除会员
};

const handleExport = () => {
  // TODO: 导出数据
};

const handleSubmit = (done: boolean) => {
  // TODO: 提交表单
  done(true);
  formVisible.value = false;
};

const getStatusColor = (status: string) => {
  const map: Record<string, string> = { active: 'green', inactive: 'gray', frozen: 'red' };
  return map[status] || 'gray';
};

const getStatusText = (status: string) => {
  const map: Record<string, string> = { active: '活跃', inactive: '不活跃', frozen: '已冻结' };
  return map[status] || status;
};

const getLevelColor = (level: string) => {
  const map: Record<string, string> = { gold: 'gold', silver: 'silver', bronze: 'bronze' };
  return map[level] || 'default';
};

const getLevelText = (level: string) => {
  const map: Record<string, string> = { gold: '黄金会员', silver: '白银会员', bronze: '青铜会员' };
  return map[level] || level;
};

fetchData();
</script>

<style scoped>
.member-container {
  padding: 20px;
}
.card-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.search-area {
  margin-bottom: 16px;
  padding: 16px;
  background: #f7f8fa;
  border-radius: 4px;
}
</style>
