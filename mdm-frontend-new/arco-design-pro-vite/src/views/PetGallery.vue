<template>
  <div class="pet-gallery-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="相册数量" :value="128" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="照片总数" :value="2560" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="总浏览量" :value="12580" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>宠物相册</span>
          <a-space>
            <a-input-search placeholder="搜索相册" style="width: 200px;" />
            <a-button type="primary" @click="handleCreate">
              <template #icon><icon-plus /></template>
              新建相册
            </a-button>
          </a-space>
        </div>
      </template>
      
      <a-row :gutter="16">
        <a-col :span="6" v-for="album in albums" :key="album.id">
          <a-card size="small" class="album-card" hoverable>
            <div class="album-cover">
              <img :src="album.cover" :alt="album.name" />
            </div>
            <div class="album-name">{{ album.name }}</div>
            <div class="album-meta">
              <span>📷 {{ album.photoCount }}</span>
              <span>👁 {{ album.views }}</span>
            </div>
          </a-card>
        </a-col>
      </a-row>
    </a-card>

    <a-modal v-model:visible="createVisible" title="创建相册" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="相册名称" required>
          <a-input v-model="form.name" placeholder="请输入相册名称" />
        </a-form-item>
        <a-form-item label="所属宠物">
          <a-select v-model="form.petId" placeholder="选择宠物">
            <a-option value="P001">小黄</a-option>
            <a-option value="P002">小红</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="相册描述">
          <a-textarea v-model="form.description" :rows="3" />
        </a-form-item>
        <a-form-item label="可见性">
          <a-radio-group v-model="form.visibility">
            <a-radio value="private">私密</a-radio>
            <a-radio value="friends">好友可见</a-radio>
            <a-radio value="public">公开</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 6 });
const createVisible = ref(false);

const form = reactive({ name: '', petId: '', description: '', visibility: 'private' });

const albums = ref([
  { id: 1, name: '小黄的日常', cover: 'https://placeholder.com/dog1.jpg', photoCount: 128, views: 1250, petName: '小黄' },
  { id: 2, name: '小红的美照', cover: 'https://placeholder.com/cat1.jpg', photoCount: 86, views: 890, petName: '小红' },
  { id: 3, name: '户外写真', cover: 'https://placeholder.com/outdoor.jpg', photoCount: 256, views: 2100, petName: '小黄' },
  { id: 4, name: '睡觉时刻', cover: 'https://placeholder.com/sleep.jpg', photoCount: 45, views: 560, petName: '小红' },
]);

const handleCreate = () => { createVisible.value = true; };
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.pet-gallery-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
.album-card { margin-bottom: 16px; cursor: pointer; }
.album-cover { height: 120px; overflow: hidden; margin-bottom: 8px; }
.album-cover img { width: 100%; height: 100%; object-fit: cover; }
.album-name { font-weight: bold; margin-bottom: 4px; }
.album-meta { color: #86909c; font-size: 12px; }
.album-meta span { margin-right: 12px; }
</style>
