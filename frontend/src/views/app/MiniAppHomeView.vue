<template>
  <div class="container">
    <a-card class="general-card" title="小程序首页配置">
      <template #extra>
        <a-button type="primary" @click="handleSave"><icon-save />保存</a-button>
      </template>
      <a-form :model="form" layout="vertical" style="max-width: 600px">
        <a-form-item label="首页标题"><a-input v-model="form.title" /></a-form-item>
        <a-form-item label="轮播图启用">
          <a-switch v-model="form.banner_enabled" />
        </a-form-item>
        <a-form-item label="快捷入口配置"><a-textarea v-model="form.shortcuts" :rows="3" /></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const form = reactive({ title: '', banner_enabled: true, shortcuts: '' })
const handleSave = () => Message.success('保存成功')
</script>
