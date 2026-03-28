<template>
  <div class="knowledge-base-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>知识库管理</span>
          <a-space>
            <a-button @click="handleSync">同步</a-button>
            <a-button type="primary" @click="handleCreate">
              <template #icon><icon-plus /></template>
              新增词条
            </a-button>
          </a-space>
        </div>
      </template>
      
      <div class="search-area">
        <a-row :gutter="16">
          <a-col :span="6">
            <a-input v-model="searchForm.keyword" placeholder="搜索问题/答案" allow-clear />
          </a-col>
          <a-col :span="4">
            <a-select v-model="searchForm.category" placeholder="分类" allow-clear>
              <a-option value="behavior">行为知识</a-option>
              <a-option value="health">健康知识</a-option>
              <a-option value="training">训练知识</a-option>
              <a-option value="emotion">情感知识</a-option>
            </a-select>
          </a-col>
          <a-col :span="2">
            <a-button type="primary" @click="handleSearch">筛选</a-button>
          </a-col>
        </a-row>
      </div>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #category="{ record }">
          <a-tag :color="getCategoryColor(record.category)">{{ record.categoryText }}</a-tag>
        </template>
        <template #version="{ record }">
          <a-tag>{{ record.version }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.enabled ? 'green' : 'gray'">
            {{ record.enabled ? '启用' : '禁用' }}
          </a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleView(record)">详情</a-link>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link @click="handleVersion(record)">版本</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 词条详情 -->
    <a-modal v-model:visible="detailVisible" title="词条详情" :width="700">
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="词条ID">{{ currentItem.id }}</a-descriptions-item>
        <a-descriptions-item label="分类">
          <a-tag :color="getCategoryColor(currentItem.category)">{{ currentItem.categoryText }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="版本">{{ currentItem.version }}</a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="currentItem.enabled ? 'green' : 'gray'">{{ currentItem.enabled ? '启用' : '禁用' }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="问题" :span="2">{{ currentItem.question }}</a-descriptions-item>
        <a-descriptions-item label="标准答案" :span="2">{{ currentItem.answer }}</a-descriptions-item>
        <a-descriptions-item label="标签" :span="2">
          <a-space>
            <a-tag v-for="tag in currentItem.tags" :key="tag">{{ tag }}</a-tag>
          </a-space>
        </a-descriptions-item>
        <a-descriptions-item label="创建时间">{{ currentItem.createdAt }}</a-descriptions-item>
        <a-descriptions-item label="更新时间">{{ currentItem.updatedAt }}</a-descriptions-item>
      </a-descriptions>
      
      <a-divider>版本历史</a-divider>
      <a-timeline>
        <a-timeline-item v-for="v in currentItem.versionHistory" :key="v.version" :color="v.isCurrent ? 'green' : 'gray'">
          <b>{{ v.version }}</b> - {{ v.updatedAt }}
          <span style="color: #86909c;">{{ v.changeNote }}</span>
          <a-tag :color="v.isCurrent ? 'green' : 'gray'" size="small">{{ v.isCurrent ? '当前' : '历史' }}</a-tag>
        </a-timeline-item>
      </a-timeline>
    </a-modal>

    <!-- 创建/编辑弹窗 -->
    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑词条' : '新增词条'" :width="600" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="问题" required>
          <a-textarea v-model="form.question" placeholder="请输入标准问题" :rows="2" />
        </a-form-item>
        <a-form-item label="标准答案" required>
          <a-textarea v-model="form.answer" placeholder="请输入标准答案" :rows="4" />
        </a-form-item>
        <a-form-item label="分类">
          <a-select v-model="form.category" placeholder="选择分类">
            <a-option value="behavior">行为知识</a-option>
            <a-option value="health">健康知识</a-option>
            <a-option value="training">训练知识</a-option>
            <a-option value="emotion">情感知识</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="标签">
          <a-select v-model="form.tags" multiple placeholder="选择或输入标签">
            <a-option value="基础">基础</a-option>
            <a-option value="进阶">进阶</a-option>
            <a-option value="高级">高级</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="版本">
          <a-input v-model="form.version" placeholder="如: v1.0.0" />
        </a-form-item>
        <a-form-item label="启用">
          <a-switch v-model="form.enabled" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 版本历史弹窗 -->
    <a-modal v-model:visible="versionVisible" title="版本历史" :width="600">
      <a-timeline>
        <a-timeline-item v-for="v in versionHistory" :key="v.version" :color="v.isCurrent ? 'green' : 'gray'">
          <b>{{ v.version }}</b> - {{ v.updatedAt }}
          <div>{{ v.changeNote }}</div>
          <a-space>
            <a-link :disabled="v.isCurrent">设为当前</a-link>
            <a-link>查看</a-link>
            <a-link>回滚</a-link>
          </a-space>
        </a-timeline-item>
      </a-timeline>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const data = ref([
  { id: 'K001', question: '狗狗摇尾巴是什么意思？', answer: '狗狗摇尾巴通常表示兴奋或友好，但也要看摇动的方向和速度。', category: 'behavior', categoryText: '行为知识', version: 'v1.2.0', tags: ['基础', '行为'], enabled: true, createdAt: '2026-01-15', updatedAt: '2026-03-20' },
  { id: 'K002', question: '如何判断狗狗是否发烧？', answer: '狗狗正常体温38-39.2°C，超过39.5°C为发烧。症状包括鼻头干、精神萎靡、食欲下降。', category: 'health', categoryText: '健康知识', version: 'v1.1.0', tags: ['健康', '基础'], enabled: true, createdAt: '2026-01-10', updatedAt: '2026-03-15' },
  { id: 'K003', question: '如何训练狗狗坐下？', answer: '1.手持零食靠近狗狗鼻子 2.将零食向上移动，狗狗自然坐下 3.坐下时说"坐" 4.立即给予奖励', category: 'training', categoryText: '训练知识', version: 'v2.0.0', tags: ['训练', '进阶'], enabled: true, createdAt: '2026-02-01', updatedAt: '2026-03-25' },
  { id: 'K004', question: '狗狗不开心时有什么表现？', answer: '狗狗不开心时可能表现为：食欲下降、不愿互动、躲藏、频繁打哈欠、舔嘴唇等。', category: 'emotion', categoryText: '情感知识', version: 'v1.0.0', tags: ['情感', '基础'], enabled: true, createdAt: '2026-02-15', updatedAt: '2026-02-15' },
  { id: 'K005', question: '猫咪为什么喜欢蹭人？', answer: '猫咪蹭人是为了留下气味标记，表示亲近和信任。同时也在主人身上留下信息素。', category: 'behavior', categoryText: '行为知识', version: 'v1.1.0', tags: ['基础', '行为'], enabled: false, createdAt: '2026-01-20', updatedAt: '2026-03-10' },
]);

const versionHistory = ref([
  { version: 'v2.0.0', updatedAt: '2026-03-25', changeNote: '新增训练步骤详解', isCurrent: true },
  { version: 'v1.0.0', updatedAt: '2026-02-01', changeNote: '初始版本', isCurrent: false },
]);

const searchForm = reactive({ keyword: '', category: '' });
const pagination = reactive({ current: 1, pageSize: 20, total: 5 });

const columns = [
  { title: '词条ID', dataIndex: 'id', width: 100 },
  { title: '问题', dataIndex: 'question', width: 200 },
  { title: '分类', slotName: 'category', width: 100 },
  { title: '版本', slotName: 'version', width: 80 },
  { title: '标签', dataIndex: 'tags', width: 150 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '更新时间', dataIndex: 'updatedAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 220, fixed: 'right' },
];

const detailVisible = ref(false);
const editVisible = ref(false);
const versionVisible = ref(false);
const isEdit = ref(false);
const currentItem = ref<any>({});

const form = reactive({
  question: '',
  answer: '',
  category: 'behavior',
  tags: [],
  version: 'v1.0.0',
  enabled: true,
});

const getCategoryColor = (cat: string) => {
  const map: Record<string, string> = { behavior: 'blue', health: 'green', training: 'orange', emotion: 'purple' };
  return map[cat] || 'default';
};

const handleSearch = () => {};
const handleSync = () => {};
const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleView = (record: any) => {
  currentItem.value = {
    ...record,
    versionHistory: [
      { version: 'v1.2.0', updatedAt: '2026-03-20', changeNote: '优化描述', isCurrent: true },
      { version: 'v1.1.0', updatedAt: '2026-02-15', changeNote: '补充内容', isCurrent: false },
      { version: 'v1.0.0', updatedAt: '2026-01-15', changeNote: '初始版本', isCurrent: false },
    ],
  };
  detailVisible.value = true;
};
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleVersion = (record: any) => { versionVisible.value = true; };
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.knowledge-base-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
.search-area { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
</style>
