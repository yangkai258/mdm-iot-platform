<template>
  <div class="social-sharing-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="分享记录" :value="456" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="今日分享" :value="28" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="转化率" :value="12" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <span>社交平台分享</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="records" title="分享记录">
          <a-table :columns="columns" :data="records" :pagination="pagination">
            <template #platform="{ record }">
              <a-tag>{{ record.platformText }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="templates" title="分享模板">
          <a-row :gutter="16">
            <a-col :span="6" v-for="tpl in templates" :key="tpl.id">
              <a-card size="small" class="tpl-card">
                <div class="tpl-preview">{{ tpl.preview }}</div>
                <div class="tpl-name">{{ tpl.name }}</div>
                <a-button type="primary" size="small" @click="handleUse(tpl)">使用</a-button>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 5 });

const records = ref([
  { id: 1, userName: '张三', platform: 'wechat', platformText: '微信', content: '我家宠物真可爱', petName: '小黄', clicks: 128, conversions: 12, time: '2026-03-28 10:00:00' },
  { id: 2, userName: '李四', platform: 'weibo', platformText: '微博', content: 'AI宠物太智能了', petName: '小红', clicks: 256, conversions: 30, time: '2026-03-28 09:00:00' },
]);

const templates = ref([
  { id: 1, name: '日常分享', preview: '🐾💕', description: '分享宠物日常' },
  { id: 2, name: '成就解锁', preview: '🏆', description: '分享训练成就' },
  { id: 3, name: '精彩瞬间', preview: '📸', description: '分享精彩照片' },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '用户', dataIndex: 'userName', width: 100 },
  { title: '平台', slotName: 'platform', width: 100 },
  { title: '内容', dataIndex: 'content', width: 200 },
  { title: '宠物', dataIndex: 'petName', width: 100 },
  { title: '点击', dataIndex: 'clicks', width: 80 },
  { title: '转化', dataIndex: 'conversions', width: 80 },
  { title: '时间', dataIndex: 'time', width: 160 },
];

const handleUse = (tpl: any) => {};
</script>

<style scoped>
.social-sharing-container { padding: 20px; }
.tpl-card { text-align: center; margin-bottom: 12px; }
.tpl-preview { font-size: 48px; margin-bottom: 8px; }
.tpl-name { font-weight: bold; margin-bottom: 8px; }
</style>
