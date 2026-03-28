<template>
  <div class="widget-market-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="应用/插件" :value="45" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="安装次数" :value="1280" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="开发者" :value="28" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="收入分成" :value="15680" prefix="¥" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>应用与插件市场</span>
          <a-button type="primary" @click="handlePublish">
            <template #icon><icon-plus /></template>
            发布应用
          </a-button>
        </div>
      </template>
      
      <a-tabs>
        <a-tab-pane key="apps" title="应用">
          <a-row :gutter="16">
            <a-col :span="6" v-for="app in apps" :key="app.id">
              <a-card size="small" class="app-card">
                <div class="app-icon">{{ app.icon }}</div>
                <div class="app-name">{{ app.name }}</div>
                <div class="app-developer">by {{ app.developer }}</div>
                <div class="app-price">{{ app.price === 0 ? '免费' : '¥' + app.price }}</div>
                <a-button type="primary" size="small" @click="handleInstall(app)">安装</a-button>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
        
        <a-tab-pane key="my-apps" title="我的应用">
          <a-table :columns="myAppColumns" :data="myApps" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="record.installed ? 'green' : 'gray'">{{ record.installed ? '已安装' : '未安装' }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <a-modal v-model:visible="publishVisible" title="发布应用" :width="700" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="应用名称" required>
          <a-input v-model="form.name" placeholder="请输入应用名称" />
        </a-form-item>
        <a-form-item label="应用图标">
          <a-upload action="#" :limit="1" />
        </a-form-item>
        <a-form-item label="应用描述">
          <a-textarea v-model="form.description" :rows="3" />
        </a-form-item>
        <a-form-item label="价格">
          <a-input-number v-model="form.price" :min="0" :precision="2" />
        </a-form-item>
        <a-form-item label="分成比例">
          <a-input-number v-model="form.revenueShare" :min="0" :max="100" suffix="%" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 5 });

const apps = ref([
  { id: 1, name: '宠物健康助手', icon: '🏥', developer: 'Dev001', price: 0, installs: 520 },
  { id: 2, name: 'AI训练师', icon: '🤖', developer: 'Dev002', price: 29, installs: 280 },
  { id: 3, name: '行为分析器', icon: '📊', developer: 'Dev001', price: 19, installs: 156 },
  { id: 4, name: '睡眠监测', icon: '😴', developer: 'Dev003', price: 0, installs: 430 },
]);

const myApps = ref([
  { id: 1, name: '宠物健康助手', developer: 'Dev001', price: 0, installed: true, installedAt: '2026-03-20' },
  { id: 2, name: 'AI训练师', developer: 'Dev002', price: 29, installed: true, installedAt: '2026-03-15' },
  { id: 3, name: '新应用', developer: 'Me', price: 0, installed: false, installedAt: null },
]);

const publishVisible = ref(false);

const form = reactive({ name: '', description: '', price: 0, revenueShare: 70 });

const myAppColumns = [
  { title: '应用名称', dataIndex: 'name', width: 150 },
  { title: '开发者', dataIndex: 'developer', width: 120 },
  { title: '价格', dataIndex: 'price', width: 80 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '安装时间', dataIndex: 'installedAt', width: 120 },
];

const handlePublish = () => { publishVisible.value = true; };
const handleInstall = (app: any) => {};
const handleSubmit = (done: boolean) => { done(true); publishVisible.value = false; };
</script>

<style scoped>
.widget-market-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
.app-card { text-align: center; margin-bottom: 12px; }
.app-icon { font-size: 48px; margin-bottom: 8px; }
.app-name { font-weight: bold; }
.app-developer { color: #86909c; font-size: 12px; }
.app-price { color: #F53F3F; font-weight: bold; margin: 8px 0; }
</style>
