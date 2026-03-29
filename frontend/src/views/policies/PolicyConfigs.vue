<template>
  <div class="container">
    <a-card class="general-card" title="策略配置">
      <template #extra>
        <a-button type="primary" @click="handleSave"><icon-save />保存</a-button>
      </template>
      <a-form :model="form" layout="vertical" style="max-width: 600px">
        <a-form-item label="策略启用">
          <a-switch v-model="form.enabled" />
        </a-form-item>
        <a-form-item label="默认策略"><a-input v-model="form.default_policy" /></a-form-item>
        <a-form-item label="优先级"><a-input-number v-model="form.priority" :min="0" style="width:100%" /></a-form-item>
        <a-form-item label="描述"><a-textarea v-model="form.description" :rows="3" /></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const form = reactive({ enabled: true, default_policy: '', priority: 0, description: '' })
const handleSave = () => Message.success('保存成功')
</script>
