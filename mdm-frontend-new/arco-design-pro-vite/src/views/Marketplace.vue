<template>
  <div class="marketplace-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card>
          <a-statistic title="表情包" :value="stats.stickers" suffix="个" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="动作" :value="stats.actions" suffix="个" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="声音" :value="stats.voices" suffix="个" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="下载量" :value="stats.downloads" />
        </a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>内容市场</span>
          <a-space>
            <a-button @click="handleAudit">审核</a-button>
            <a-button type="primary" @click="handleCreate">
              <template #icon><icon-plus /></template>
              上传内容
            </a-button>
          </a-space>
        </div>
      </template>
      
      <a-tabs v-model="activeTab">
        <a-tab-pane key="stickers" title="表情包">
          <a-row :gutter="16">
            <a-col :span="4" v-for="item in stickers" :key="item.id">
              <a-card size="small" class="market-card">
                <div class="preview">{{ item.preview }}</div>
                <div class="name">{{ item.name }}</div>
                <div class="author">by {{ item.author }}</div>
                <div class="stats">
                  <span>⬇️ {{ item.downloads }}</span>
                  <span>❤️ {{ item.likes }}</span>
                </div>
                <a-button size="small" type="primary" @click="handleDownload(item)">下载</a-button>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
        
        <a-tab-pane key="actions" title="动作">
          <a-table :columns="actionColumns" :data="actions" :pagination="pagination">
            <template #type="{ record }">
              <a-tag>{{ record.typeText }}</a-tag>
            </template>
            <template #status="{ record }">
              <a-tag :color="record.published ? 'green' : 'gray'">{{ record.published ? '已发布' : '草稿' }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-link @click="handlePreview(record)">预览</a-link>
                <a-link @click="handleEditContent(record)">编辑</a-link>
                <a-link status="danger" @click="handleDeleteContent(record)">删除</a-link>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="voices" title="声音定制">
          <a-row :gutter="16">
            <a-col :span="6" v-for="voice in voices" :key="voice.id">
              <a-card size="small" class="voice-card">
                <div class="voice-icon">🔊</div>
                <div class="name">{{ voice.name }}</div>
                <div class="type">{{ voice.typeText }}</div>
                <div class="price">{{ voice.price === 0 ? '免费' : '¥' + voice.price }}</div>
                <a-button size="small" type="primary" @click="handleDownloadVoice(voice)">下载</a-button>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <!-- 上传弹窗 -->
    <a-modal v-model:visible="uploadVisible" title="上传内容" :width="600" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="内容类型">
          <a-radio-group v-model="form.type">
            <a-radio value="sticker">表情包</a-radio>
            <a-radio value="action">动作</a-radio>
            <a-radio value="voice">声音</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="名称">
          <a-input v-model="form.name" placeholder="请输入内容名称" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" :rows="2" />
        </a-form-item>
        <a-form-item label="上传文件">
          <a-upload action="#" :limit="1" />
        </a-form-item>
        <a-form-item label="价格">
          <a-input-number v-model="form.price" :min="0" :precision="2" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const activeTab = ref('stickers');

const stats = reactive({ stickers: 256, actions: 128, voices: 64, downloads: 15832 });

const stickers = ref([
  { id: 1, name: '开心大笑', preview: '😄', author: 'User001', downloads: 1256, likes: 328 },
  { id: 2, name: '爱心发射', preview: '💕', author: 'User002', downloads: 980, likes: 256 },
  { id: 3, name: '666', preview: '6️⃣', author: 'User003', downloads: 856, likes: 198 },
  { id: 4, name: '点赞', preview: '👍', author: 'User001', downloads: 2340, likes: 512 },
]);

const actions = ref([
  { id: 1, name: '后空翻', type: 'action', typeText: '动作', author: 'User001', price: 0, downloads: 580, published: true },
  { id: 2, name: '撒娇', type: 'emotion', typeText: '情感', author: 'User002', price: 9, downloads: 320, published: true },
  { id: 3, name: '跳舞', type: 'action', typeText: '动作', author: 'User003', price: 0, downloads: 890, published: false },
]);

const voices = ref([
  { id: 1, name: '甜美声音', type: 'female', typeText: '女声', price: 0 },
  { id: 2, name: '磁性男声', type: 'male', typeText: '男声', price: 9 },
  { id: 3, name: '萌系童声', type: 'child', typeText: '童声', price: 19 },
]);

const pagination = reactive({ current: 1, pageSize: 10, total: 3 });
const uploadVisible = ref(false);

const form = reactive({ type: 'sticker', name: '', description: '', price: 0 });

const actionColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '名称', dataIndex: 'name', width: 150 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '作者', dataIndex: 'author', width: 100 },
  { title: '价格', dataIndex: 'price', width: 80 },
  { title: '下载', dataIndex: 'downloads', width: 80 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' },
];

const handleAudit = () => {};
const handleCreate = () => { uploadVisible.value = true; };
const handleDownload = (item: any) => {};
const handlePreview = (record: any) => {};
const handleEditContent = (record: any) => {};
const handleDeleteContent = (record: any) => {};
const handleDownloadVoice = (voice: any) => {};
const handleSubmit = (done: boolean) => { done(true); uploadVisible.value = false; };
</script>

<style scoped>
.marketplace-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
.market-card { text-align: center; margin-bottom: 8px; }
.preview { font-size: 48px; }
.name { font-weight: bold; margin-top: 8px; }
.author { color: #86909c; font-size: 12px; }
.stats { margin: 8px 0; }
.voice-card { text-align: center; }
.voice-icon { font-size: 48px; }
.type { color: #86909c; }
.price { color: #F53F3F; font-weight: bold; margin: 8px 0; }
</style>
