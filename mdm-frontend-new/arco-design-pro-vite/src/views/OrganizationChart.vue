<template>
  <div class="organization-chart-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="部门" :value="12" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="员工" :value="128" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="岗位" :value="36" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="在线" :value="98" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>组织架构管理</span>
          <a-space>
            <a-button @click="handleExpand">展开全部</a-button>
            <a-button type="primary" @click="handleCreate">
              <template #icon><icon-plus /></template>
              添加部门
            </a-button>
          </a-space>
        </div>
      </template>
      
      <a-tree :data="treeData" :show-line="true" @select="handleSelect">
        <template #title="{ data }">
          <span>{{ data.name }}</span>
          <a-link style="margin-left: 8px;" @click.stop="handleEditDept(data)">编辑</a-link>
        </template>
      </a-tree>
    </a-card>

    <a-drawer v-model:visible="detailVisible" :title="currentDept.name" :width="500">
      <a-descriptions :column="1" bordered>
        <a-descriptions-item label="部门名称">{{ currentDept.name }}</a-descriptions-item>
        <a-descriptions-item label="上级部门">{{ currentDept.parent }}</a-descriptions-item>
        <a-descriptions-item label="负责人">{{ currentDept.leader }}</a-descriptions-item>
        <a-descriptions-item label="员工数量">{{ currentDept.count }}</a-descriptions-item>
        <a-descriptions-item label="描述">{{ currentDept.description }}</a-descriptions-item>
      </a-descriptions>
      <a-divider>员工列表</a-divider>
      <a-list :data="currentEmployees" />
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const detailVisible = ref(false);
const currentDept = ref<any>({});
const currentEmployees = ref([]);

const treeData = ref([
  {
    key: '1',
    name: '技术部',
    children: [
      { key: '1-1', name: '研发组' },
      { key: '1-2', name: '测试组' },
    ],
  },
  { key: '2', name: '运营部' },
  { key: '3', name: '市场部' },
]);

const handleExpand = () => {};
const handleCreate = () => {};
const handleSelect = (selected: boolean, node: any) => {
  if (selected && node.value) {
    currentDept.value = { name: node.value.name, parent: '-', leader: '张总', count: 25, description: '技术部门负责产品研发' };
    currentEmployees.value = ['张三', '李四', '王五'];
    detailVisible.value = true;
  }
};
const handleEditDept = (data: any) => {};
</script>

<style scoped>
.organization-chart-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
