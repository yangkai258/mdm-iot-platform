<template>
  <div class="family-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>家庭成员</span>
          <a-button type="primary" @click="handleAdd">
            <template #icon><icon-plus /></template>
            添加成员
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #avatar="{ record }">
          <a-avatar :size="40" :style="{ backgroundColor: record.color }">
            {{ record.name?.charAt(0) }}
          </a-avatar>
        </template>
        <template #role="{ record }">
          <a-tag :color="getRoleColor(record.role)">{{ record.roleText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.active ? 'green' : 'gray'">
            {{ record.active ? '在线' : '离线' }}
          </a-tag>
        </template>
        <template #permissions="{ record }">
          <a-space wrap>
            <a-tag v-for="p in record.permissions" :key="p" size="small">{{ p }}</a-tag>
          </a-space>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link @click="handlePermissions(record)">权限</a-link>
            <a-link status="danger" @click="handleDelete(record)">移除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 添加/编辑弹窗 -->
    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑成员' : '添加成员'" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="成员名称" required>
          <a-input v-model="form.name" placeholder="请输入成员名称" />
        </a-form-item>
        <a-form-item label="成员角色">
          <a-select v-model="form.role" placeholder="选择角色">
            <a-option value="owner">主人</a-option>
            <a-option value="family">家庭成员</a-option>
            <a-option value="guest">访客</a-option>
            <a-option value="child">儿童</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="手机号">
          <a-input v-model="form.phone" placeholder="请输入手机号" />
        </a-form-item>
        <a-form-item label="与宠物关系">
          <a-input v-model="form.relationship" placeholder="如：爸爸、妈妈、哥哥等" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 权限配置弹窗 -->
    <a-modal v-model:visible="permVisible" title="成员权限配置">
      <a-form layout="vertical">
        <a-form-item label="设备控制权限">
          <a-checkbox-group v-model="currentMember.permissions">
            <a-checkbox value="view">查看设备</a-checkbox>
            <a-checkbox value="control">控制设备</a-checkbox>
            <a-checkbox value="ota">执行OTA</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
        <a-form-item label="宠物管理权限">
          <a-checkbox-group v-model="currentMember.permissions">
            <a-checkbox value="feed">喂食控制</a-checkbox>
            <a-checkbox value="schedule">作息安排</a-checkbox>
            <a-checkbox value="health">健康查看</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
        <a-form-item label="数据查看权限">
          <a-checkbox-group v-model="currentMember.permissions">
            <a-checkbox value="history">历史数据</a-checkbox>
            <a-checkbox value="reports">报告查看</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="permVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSavePerm">保存</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const data = ref([
  { id: 1, name: '张三', phone: '138****1234', role: 'owner', roleText: '主人', relationship: '爸爸', active: true, color: '#165DFF', permissions: ['view', 'control', 'ota', 'feed', 'schedule', 'health', 'history', 'reports'] },
  { id: 2, name: '李四', phone: '139****5678', role: 'family', roleText: '家庭成员', relationship: '妈妈', active: true, color: '#00B42A', permissions: ['view', 'control', 'feed', 'history'] },
  { id: 3, name: '王小明', phone: '137****9012', role: 'child', roleText: '儿童', relationship: '哥哥', active: false, color: '#FF7D00', permissions: ['view', 'history'] },
  { id: 4, name: '访客', phone: '-', role: 'guest', roleText: '访客', relationship: '亲戚', active: false, color: '#86909c', permissions: ['view'] },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 4 });

const columns = [
  { title: '', slotName: 'avatar', width: 60 },
  { title: '名称', dataIndex: 'name', width: 120 },
  { title: '手机', dataIndex: 'phone', width: 130 },
  { title: '角色', slotName: 'role', width: 100 },
  { title: '与宠物关系', dataIndex: 'relationship', width: 120 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '权限', slotName: 'permissions', width: 250 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' },
];

const editVisible = ref(false);
const permVisible = ref(false);
const isEdit = ref(false);
const currentMember = ref<any>({});

const form = reactive({
  name: '',
  phone: '',
  role: 'family',
  relationship: '',
});

const getRoleColor = (role: string) => {
  const map: Record<string, string> = { owner: 'red', family: 'blue', guest: 'gray', child: 'orange' };
  return map[role] || 'default';
};

const handleAdd = () => {
  isEdit.value = false;
  Object.assign(form, { name: '', phone: '', role: 'family', relationship: '' });
  editVisible.value = true;
};

const handleEdit = (record: any) => {
  isEdit.value = true;
  Object.assign(form, record);
  editVisible.value = true;
};

const handlePermissions = (record: any) => {
  currentMember.value = { ...record, permissions: [...record.permissions] };
  permVisible.value = true;
};

const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
const handleSavePerm = () => { permVisible.value = false; };
</script>

<style scoped>
.family-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
