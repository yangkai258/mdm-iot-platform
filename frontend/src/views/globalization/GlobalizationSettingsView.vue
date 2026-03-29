<template>
  <div class="container">
    <a-card class="general-card" title="全球化设置">
      <template #extra>
        <a-button type="primary" @click="handleSave"><icon-save />保存</a-button>
      </template>
      <a-form :model="form" layout="vertical" style="max-width: 600px">
        <a-form-item label="默认语言">
          <a-select v-model="form.default_language" style="width: 100%">
            <a-option value="zh-CN">简体中文</a-option>
            <a-option value="en-US">English</a-option>
            <a-option value="ja-JP">日本語</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="默认时区">
          <a-select v-model="form.default_timezone" style="width: 100%">
            <a-option value="Asia/Shanghai">Asia/Shanghai (UTC+8)</a-option>
            <a-option value="America/New_York">America/New_York (UTC-5)</a-option>
            <a-option value="Europe/London">Europe/London (UTC+0)</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="默认货币">
          <a-input v-model="form.default_currency" />
        </a-form-item>
        <a-form-item label="数据驻留区域">
          <a-select v-model="form.data_region" style="width: 100%">
            <a-option value="cn">中国</a-option>
            <a-option value="us">美国</a-option>
            <a-option value="eu">欧洲</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const form = reactive({ default_language: 'zh-CN', default_timezone: 'Asia/Shanghai', default_currency: 'CNY', data_region: 'cn' })
const handleSave = () => Message.success('保存成功')
</script>
