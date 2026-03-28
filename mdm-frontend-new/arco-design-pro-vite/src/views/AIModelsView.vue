<template>
  <div class="ai-models-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>AI模型管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            上传模型
          </a-button>
        </div>
      </template>
      
      <div class="search-area">
        <a-row :gutter="16">
          <a-col :span="6">
            <a-input v-model="searchForm.keyword" placeholder="搜索模型名称/版本" allow-clear />
          </a-col>
          <a-col :span="4">
            <a-select v-model="searchForm.type" placeholder="模型类型" allow-clear>
              <a-option value="tts">TTS语音</a-option>
              <a-option value="asr">ASR语音识别</a-option>
              <a-option value="nlu">NLU意图识别</a-option>
              <a-option value="llm">大语言模型</a-option>
              <a-option value="cv">视觉模型</a-option>
            </a-select>
          </a-col>
          <a-col :span="4">
            <a-select v-model="searchForm.status" placeholder="状态" allow-clear>
              <a-option value="active">已上线</a-option>
              <a-option value="testing">测试中</a-option>
              <a-option value="deprecated">已下线</a-option>
            </a-select>
          </a-col>
          <a-col :span="2">
            <a-button type="primary" @click="handleSearch">筛选</a-button>
          </a-col>
        </a-row>
      </div>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag :color="getTypeColor(record.type)">{{ record.typeText }}</a-tag>
        </template>
        <template #size="{ record }">
          {{ formatSize(record.size) }}
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
        </template>
        <template #defaultVersion="{ record }">
          <a-tag :color="record.isDefault ? 'green' : 'gray'">
            {{ record.isDefault ? '默认' : '-' }}
          </a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleView(record)">详情</a-link>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link @click="handleSetDefault(record)">设为默认</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 模型详情 -->
    <a-modal v-model:visible="detailVisible" title="模型详情" :width="700">
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="模型ID">{{ currentModel.id }}</a-descriptions-item>
        <a-descriptions-item label="模型名称">{{ currentModel.name }}</a-descriptions-item>
        <a-descriptions-item label="版本">{{ currentModel.version }}</a-descriptions-item>
        <a-descriptions-item label="类型">
          <a-tag :color="getTypeColor(currentModel.type)">{{ currentModel.typeText }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="模型大小">{{ formatSize(currentModel.size) }}</a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="getStatusColor(currentModel.status)">{{ currentModel.statusText }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="上线时间">{{ currentModel.deployedAt }}</a-descriptions-item>
        <a-descriptions-item label="默认版本">
          <a-tag :color="currentModel.isDefault ? 'green' : 'gray'">{{ currentModel.isDefault ? '是' : '否' }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="描述" :span="2">{{ currentModel.description }}</a-descriptions-item>
      </a-descriptions>
      
      <a-divider>版本历史</a-divider>
      <a-timeline>
        <a-timeline-item v-for="v in currentModel.versions" :key="v.version" :color="v.isActive ? 'green' : 'gray'">
          <b>{{ v.version }}</b> - {{ v.deployedAt }}
          <a-tag :color="v.isActive ? 'green' : 'gray'" size="small">{{ v.isActive ? '当前' : '历史' }}</a-tag>
        </a-timeline-item>
      </a-timeline>
    </a-modal>

    <!-- 上传/编辑弹窗 -->
    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑模型' : '上传模型'" :width="600" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="模型名称" required>
          <a-input v-model="form.name" placeholder="请输入模型名称" />
        </a-form-item>
        <a-form-item label="模型类型" required>
          <a-select v-model="form.type" placeholder="选择模型类型">
            <a-option value="tts">TTS语音合成</a-option>
            <a-option value="asr">ASR语音识别</a-option>
            <a-option value="nlu">NLU意图识别</a-option>
            <a-option value="llm">大语言模型</a-option>
            <a-option value="cv">视觉模型</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="版本号">
          <a-input v-model="form.version" placeholder="如: v1.0.0" />
        </a-form-item>
        <a-form-item label="模型文件">
          <a-upload action="#" :limit="1" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" placeholder="请输入模型描述" :rows="3" />
        </a-form-item>
        <a-form-item label="设为默认版本">
          <a-switch v-model="form.isDefault" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const data = ref([
  { id: 'M001', name: 'MiniMax-TTS', version: 'v2.1.0', type: 'tts', typeText: 'TTS语音', size: 256000000, status: 'active', statusText: '已上线', deployedAt: '2026-03-15', isDefault: true, description: 'MiniMax高质量语音合成模型，支持中英双语' },
  { id: 'M002', name: 'MiniMax-ASR', version: 'v1.5.0', type: 'asr', typeText: 'ASR语音识别', size: 180000000, status: 'active', statusText: '已上线', deployedAt: '2026-03-10', isDefault: true, description: '语音识别模型，识别准确率98%' },
  { id: 'M003', name: 'Pet-NLU', version: 'v3.0.0', type: 'nlu', typeText: 'NLU意图识别', size: 520000000, status: 'testing', statusText: '测试中', deployedAt: '2026-03-20', isDefault: false, description: '宠物意图识别专用NLU模型' },
  { id: 'M004', name: 'Pet-LLM', version: 'v1.0.0', type: 'llm', typeText: '大语言模型', size: 2048000000, status: 'active', statusText: '已上线', deployedAt: '2026-03-01', isDefault: true, description: '宠物对话专用大语言模型' },
  { id: 'M005', name: 'Pet-CV', version: 'v2.0.0', type: 'cv', typeText: '视觉模型', size: 1024000000, status: 'deprecated', statusText: '已下线', deployedAt: '2026-02-01', isDefault: false, description: '宠物视觉识别模型，已被v3.0替代' },
]);

const searchForm = reactive({ keyword: '', type: '', status: '' });
const pagination = reactive({ current: 1, pageSize: 20, total: 5 });

const columns = [
  { title: '模型ID', dataIndex: 'id', width: 100 },
  { title: '模型名称', dataIndex: 'name', width: 150 },
  { title: '版本', dataIndex: 'version', width: 100 },
  { title: '类型', slotName: 'type', width: 120 },
  { title: '大小', slotName: 'size', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '默认', slotName: 'defaultVersion', width: 80 },
  { title: '上线时间', dataIndex: 'deployedAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 220, fixed: 'right' },
];

const detailVisible = ref(false);
const editVisible = ref(false);
const isEdit = ref(false);
const currentModel = ref<any>({});

const form = reactive({
  name: '',
  type: '',
  version: '',
  description: '',
  isDefault: false,
});

const getTypeColor = (type: string) => {
  const map: Record<string, string> = { tts: 'blue', asr: 'green', nlu: 'purple', llm: 'orange', cv: 'cyan' };
  return map[type] || 'default';
};

const getStatusColor = (status: string) => {
  const map: Record<string, string> = { active: 'green', testing: 'blue', deprecated: 'gray' };
  return map[status] || 'default';
};

const formatSize = (bytes: number) => {
  if (bytes >= 1073741824) return (bytes / 1073741824).toFixed(1) + ' GB';
  if (bytes >= 1048576) return (bytes / 1048576).toFixed(0) + ' MB';
  return (bytes / 1024).toFixed(0) + ' KB';
};

const handleSearch = () => {};
const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleView = (record: any) => {
  currentModel.value = {
    ...record,
    versions: [
      { version: 'v2.1.0', deployedAt: '2026-03-15', isActive: true },
      { version: 'v2.0.0', deployedAt: '2026-02-01', isActive: false },
      { version: 'v1.0.0', deployedAt: '2026-01-01', isActive: false },
    ],
  };
  detailVisible.value = true;
};
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleSetDefault = (record: any) => {};
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.ai-models-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
.search-area { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
</style>
