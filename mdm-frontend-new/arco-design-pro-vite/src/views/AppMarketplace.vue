<template>
  <div class="app-marketplace-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="应用总数" :value="128" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="安装量" :value="25600" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="开发者" :value="56" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="收入" :value="128500" prefix="¥" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <span>应用市场</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="apps" title="应用">
          <a-row :gutter="16">
            <a-col :span="6" v-for="app in apps" :key="app.id">
              <a-card size="small" class="app-card" hoverable>
                <div class="app-icon">{{ app.icon }}</div>
                <div class="app-name">{{ app.name }}</div>
                <div class="app-developer">开发者: {{ app.developer }}</div>
                <div class="app-desc">{{ app.description }}</div>
                <div class="app-stats">
                  <span>⭐ {{ app.rating }}</span>
                  <span>📥 {{ app.installs }}</span>
                </div>
                <div class="app-price">
                  <span v-if="app.price === 0" style="color: green;">免费</span>
                  <span v-else style="color: #F53F3F;">¥{{ app.price }}</span>
                </div>
                <a-button type="primary" long size="small" @click="handleInstall(app)">安装</a-button>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
        
        <a-tab-pane key="categories" title="分类">
          <a-row :gutter="16">
            <a-col :span="4" v-for="cat in categories" :key="cat.id">
              <a-card size="small" class="cat-card" hoverable>
                <div class="cat-icon">{{ cat.icon }}</div>
                <div class="cat-name">{{ cat.name }}</div>
                <div class="cat-count">{{ cat.count }}个应用</div>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
        
        <a-tab-pane key="my-apps" title="我的应用">
          <a-table :columns="myAppColumns" :data="myApps" :pagination="pagination" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 5 });

const apps = ref([
  { id: 1, name: '宠物健康助手', icon: '🏥', developer: 'DevTeam1', description: '监测宠物健康状况', price: 0, rating: 4.8, installs: 5600 },
  { id: 2, name: 'AI训练师', icon: '🤖', developer: 'DevTeam2', description: '智能训练建议', price: 29, rating: 4.6, installs: 3200 },
  { id: 3, name: '行为分析器', icon: '📊', developer: 'DevTeam1', description: '分析宠物行为模式', price: 19, rating: 4.5, installs: 2100 },
  { id: 4, name: '睡眠监测', icon: '😴', developer: 'DevTeam3', description: '监测睡眠质量', price: 0, rating: 4.7, installs: 4800 },
  { id: 5, name: '运动计步', icon: '👟', developer: 'DevTeam4', description: '记录运动数据', price: 9, rating: 4.3, installs: 1800 },
  { id: 6, name: '饮食管理', icon: '🍖', developer: 'DevTeam2', description: '科学饮食计划', price: 0, rating: 4.4, installs: 2900 },
]);

const categories = ref([
  { id: 1, name: '健康', icon: '❤️', count: 28 },
  { id: 2, name: '训练', icon: '🎯', count: 35 },
  { id: 3, name: '社交', icon: '👥', count: 22 },
  { id: 4, name: '工具', icon: '🔧', count: 45 },
]);

const myApps = ref([
  { id: 1, name: '宠物健康助手', developer: 'DevTeam1', installedAt: '2026-03-20', status: 'active' },
  { id: 2, name: 'AI训练师', developer: 'DevTeam2', installedAt: '2026-03-15', status: 'active' },
]);

const myAppColumns = [
  { title: '应用名称', dataIndex: 'name', width: 200 },
  { title: '开发者', dataIndex: 'developer', width: 150 },
  { title: '安装时间', dataIndex: 'installedAt', width: 120 },
  { title: '状态', dataIndex: 'status', width: 100 },
];

const handleInstall = (app: any) => {};
</script>

<style scoped>
.app-marketplace-container { padding: 20px; }
.app-card { text-align: center; margin-bottom: 12px; }
.app-icon { font-size: 48px; margin-bottom: 8px; }
.app-name { font-weight: bold; margin-bottom: 4px; }
.app-developer { color: #86909c; font-size: 12px; margin-bottom: 8px; }
.app-desc { color: #4a4a4a; font-size: 12px; margin-bottom: 8px; min-height: 32px; }
.app-stats { margin-bottom: 8px; }
.app-stats span { margin: 0 4px; }
.app-price { font-weight: bold; margin-bottom: 8px; }
.cat-card { text-align: center; margin-bottom: 12px; cursor: pointer; }
.cat-icon { font-size: 36px; margin-bottom: 8px; }
.cat-name { font-weight: bold; margin-bottom: 4px; }
.cat-count { color: #86909c; font-size: 12px; }
</style>
