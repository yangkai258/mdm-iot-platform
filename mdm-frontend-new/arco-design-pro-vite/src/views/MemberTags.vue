<template>
  <div class="member-tags-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>会员标签管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新增标签
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag :color="record.color">{{ record.typeText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-switch :checked="record.enabled" disabled />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link @click="handleViewMembers(record)">查看会员</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 创建/编辑弹窗 -->
    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑标签' : '新增标签'" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="标签名称" required>
          <a-input v-model="form.name" placeholder="请输入标签名称" />
        </a-form-item>
        <a-form-item label="标签类型">
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="behavior">行为标签</a-option>
            <a-option value="preference">偏好标签</a-option>
            <a-option value="attribute">属性标签</a-option>
            <a-option value="custom">自定义</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="标签颜色">
          <a-color-picker v-model="form.color" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" placeholder="标签描述" :rows="2" />
        </a-form-item>
        <a-form-item label="启用">
          <a-switch v-model="form.enabled" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 标签会员弹窗 -->
    <a-modal v-model:visible="membersVisible" title="标签会员" :width="800">
      <a-table :columns="memberColumns" :data="tagMembers" :pagination="pagination">
        <template #level="{ record }">
          <a-tag>{{ record.levelText }}</a-tag>
        </template>
      </a-table>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const data = ref([
  { id: 1, name: '高价值', type: 'attribute', typeText: '属性标签', color: 'red', memberCount: 56, enabled: true, description: '累计消费超过5000元' },
  { id: 2, name: '活跃用户', type: 'behavior', typeText: '行为标签', color: 'green', memberCount: 128, enabled: true, description: '30天内有互动' },
  { id: 3, name: '爱宠人士', type: 'preference', typeText: '偏好标签', color: 'blue', memberCount: 89, enabled: true, description: '有多只宠物' },
  { id: 4, name: '新手用户', type: 'behavior', typeText: '行为标签', color: 'orange', memberCount: 234, enabled: true, description: '注册时间不足30天' },
  { id: 5, name: 'VIP', type: 'attribute', typeText: '属性标签', color: 'gold', memberCount: 23, enabled: false, description: '年度订阅用户' },
]);

const tagMembers = ref([
  { id: 1, name: '张三', phone: '138****1234', level: 2, levelText: '银卡', taggedAt: '2026-03-28' },
  { id: 2, name: '李四', phone: '139****5678', level: 3, levelText: '金卡', taggedAt: '2026-03-27' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 5 });
const editVisible = ref(false);
const membersVisible = ref(false);
const isEdit = ref(false);

const form = reactive({ name: '', type: 'behavior', color: '#165DFF', description: '', enabled: true });

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '标签名称', dataIndex: 'name', width: 120 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '会员数', dataIndex: 'memberCount', width: 80 },
  { title: '描述', dataIndex: 'description' },
  { title: '启用', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' },
];

const memberColumns = [
  { title: '会员ID', dataIndex: 'id', width: 80 },
  { title: '姓名', dataIndex: 'name', width: 120 },
  { title: '手机', dataIndex: 'phone', width: 130 },
  { title: '等级', slotName: 'level', width: 80 },
  { title: '打标时间', dataIndex: 'taggedAt', width: 120 },
];

const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleViewMembers = (record: any) => { membersVisible.value = true; };
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.member-tags-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
