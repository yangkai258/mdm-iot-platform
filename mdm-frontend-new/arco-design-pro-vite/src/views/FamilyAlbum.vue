<template>
  <div class="family-album-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="相册数量" :value="5" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="照片数量" :value="128" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="视频数量" :value="12" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="共享成员" :value="4" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>家庭相册</span>
          <a-space>
            <a-button @click="handleUpload">
              <template #icon><icon-upload /></template>
              上传
            </a-button>
            <a-button type="primary" @click="handleCreateAlbum">
              <template #icon><icon-plus /></template>
              新建相册
            </a-button>
          </a-space>
        </div>
      </template>
      
      <a-row :gutter="16">
        <a-col :span="6" v-for="album in albums" :key="album.id">
          <a-card size="small" class="album-card">
            <div class="album-cover">
              <icon-image :size="64" />
            </div>
            <div class="album-name">{{ album.name }}</div>
            <div class="album-meta">{{ album.photoCount }}张照片</div>
          </a-card>
        </a-col>
      </a-row>
    </a-card>

    <a-modal v-model:visible="createVisible" title="新建相册" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="相册名称" required>
          <a-input v-model="form.name" placeholder="请输入相册名称" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" :rows="2" />
        </a-form-item>
        <a-form-item label="共享成员">
          <a-select v-model="form.members" multiple placeholder="选择成员">
            <a-option value="1">爸爸</a-option>
            <a-option value="2">妈妈</a-option>
            <a-option value="3">哥哥</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const albums = ref([
  { id: 1, name: '小黄成长记', photoCount: 45, description: '记录小黄的成长历程' },
  { id: 2, name: '公园游玩', photoCount: 28, description: '周末公园玩耍照片' },
  { id: 3, name: '生日派对', photoCount: 35, description: '小黄生日会' },
  { id: 4, name: '日常萌照', photoCount: 20, description: '日常可爱瞬间' },
]);

const createVisible = ref(false);
const form = reactive({ name: '', description: '', members: [] });

const handleUpload = () => {};
const handleCreateAlbum = () => { createVisible.value = true; };
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.family-album-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
.album-card { cursor: pointer; margin-bottom: 12px; }
.album-cover { text-align: center; padding: 24px; background: #f7f8fa; border-radius: 4px; margin-bottom: 8px; }
.album-name { font-weight: bold; }
.album-meta { color: #86909c; font-size: 12px; }
</style>
