<template>
  <div class="widget-config-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="组件总数" :value="86" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="使用次数" :value="25600" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="开发者" :value="12" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>自定义组件配置</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建组件
          </a-button>
        </div>
      </template>
      
      <a-row :gutter="16">
        <a-col :span="6" v-for="widget in widgets" :key="widget.id">
          <a-card size="small" class="widget-card">
            <div class="widget-preview">{{ widget.preview }}</div>
            <div class="widget-name">{{ widget.name }}</div>
            <div class="widget-desc">{{ widget.description }}</div>
            <div class="widget-actions">
              <a-link @click="handleEdit(widget)">编辑</a-link>
              <a-link @click="handlePreview(widget)">预览</a-link>
            </div>
          </a-card>
        </a-col>
      </a-row>
    </a-card>

    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑组件' : '创建组件'" :width="800" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="组件名称" required>
          <a-input v-model="form.name" placeholder="请输入组件名称" />
        </a-form-item>
        <a-form-item label="组件类型">
          <a-select v-model="form.type">
            <a-option value="chart">图表</a-option>
            <a-option value="table">表格</a-option>
            <a-option value="form">表单</a-option>
            <a-option value="custom">自定义</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="组件配置">
          <a-textarea v-model="form.config" :rows="8" placeholder="JSON格式配置" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 6 });
const editVisible = ref(false);
const isEdit = ref(false);

const form = reactive({ name: '', type: 'chart', config: '' });

const widgets = ref([
  { id: 1, name: '设备状态图', preview: '📊', description: '展示设备在线状态分布', usageCount: 2560 },
  { id: 2, name: '会员趋势表', preview: '📈', description: '展示会员增长趋势', usageCount: 1280 },
  { id: 3, name: '销售排行', preview: '🏆', description: '展示销售排行榜', usageCount: 860 },
  { id: 4, name: '实时数据', preview: '⚡', description: '实时数据监控', usageCount: 520 },
]);

const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleEdit = (widget: any) => { isEdit.value = true; Object.assign(form, widget); editVisible.value = true; };
const handlePreview = (widget: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.widget-config-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
.widget-card { text-align: center; margin-bottom: 16px; }
.widget-preview { font-size: 48px; margin-bottom: 8px; }
.widget-name { font-weight: bold; margin-bottom: 4px; }
.widget-desc { color: #86909c; font-size: 12px; margin-bottom: 8px; }
.widget-actions { display: flex; justify-content: center; gap: 8px; }
</style>
