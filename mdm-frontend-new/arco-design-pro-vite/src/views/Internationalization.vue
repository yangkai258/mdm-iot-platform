<template>
  <div class="i18n-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>国际化配置</span>
          <a-space>
            <a-button @click="handleSync">同步翻译</a-button>
            <a-button type="primary" @click="handleAdd">
              <template #icon><icon-plus /></template>
              添加语言
            </a-button>
          </a-space>
        </div>
      </template>
      
      <a-tabs>
        <a-tab-pane key="languages" title="语言配置">
          <a-table :columns="columns" :data="languages" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="record.enabled ? 'green' : 'gray'">{{ record.enabled ? '启用' : '禁用' }}</a-tag>
            </template>
            <template #default="{ record }">
              <a-tag v-if="record.isDefault" color="blue">默认</a-tag>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-link @click="handleEditLang(record)">编辑</a-link>
                <a-link @click="handleTranslations(record)">翻译</a-link>
                <a-link v-if="!record.isDefault" status="danger" @click="handleDeleteLang(record)">删除</a-link>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="translations" title="翻译管理">
          <a-alert style="margin-bottom: 16px;">
            当前编辑: <a-tag>中文(简体) - zh-CN</a-tag>
          </a-alert>
          <a-table :columns="transColumns" :data="translations" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="record.translated ? 'green' : 'orange'">
                {{ record.translated ? '已翻译' : '未翻译' }}
              </a-tag>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <!-- 添加语言弹窗 -->
    <a-modal v-model:visible="addVisible" title="添加语言" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="语言名称" required>
          <a-input v-model="form.name" placeholder="如: English" />
        </a-form-item>
        <a-form-item label="语言代码" required>
          <a-input v-model="form.code" placeholder="如: en-US" />
        </a-form-item>
        <a-form-item label="设为默认语言">
          <a-switch v-model="form.isDefault" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 5 });

const languages = ref([
  { id: 1, name: '中文(简体)', code: 'zh-CN', isDefault: true, enabled: true, progress: 100 },
  { id: 2, name: 'English', code: 'en-US', isDefault: false, enabled: true, progress: 95 },
  { id: 3, name: '日本語', code: 'ja-JP', isDefault: false, enabled: true, progress: 80 },
  { id: 4, name: '한국어', code: 'ko-KR', isDefault: false, enabled: false, progress: 60 },
]);

const translations = ref([
  { id: 1, key: 'menu.dashboard', 'zh-CN': '仪表盘', 'en-US': 'Dashboard', 'ja-JP': 'ダッシュボード', translated: true },
  { id: 2, key: 'menu.devices', 'zh-CN': '设备管理', 'en-US': 'Device Management', 'ja-JP': 'デバイス管理', translated: true },
  { id: 3, key: 'menu.pets', 'zh-CN': '宠物管理', 'en-US': 'Pet Management', 'ja-JP': '', translated: false },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '语言名称', dataIndex: 'name', width: 120 },
  { title: '代码', dataIndex: 'code', width: 100 },
  { title: '默认', slotName: 'default', width: 80 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '翻译进度', dataIndex: 'progress', width: 120 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' },
];

const transColumns = [
  { title: 'Key', dataIndex: 'key', width: 200 },
  { title: '中文', dataIndex: 'zh-CN', width: 120 },
  { title: 'English', dataIndex: 'en-US', width: 120 },
  { title: '日本語', dataIndex: 'ja-JP', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
];

const addVisible = ref(false);
const form = reactive({ name: '', code: '', isDefault: false });

const handleSync = () => {};
const handleAdd = () => { addVisible.value = true; };
const handleEditLang = (record: any) => {};
const handleTranslations = (record: any) => {};
const handleDeleteLang = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); addVisible.value = false; };
</script>

<style scoped>
.i18n-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
