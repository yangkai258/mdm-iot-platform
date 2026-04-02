<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="微信设置">
      <template #extra>
        <a-button type="primary" @click="handleSave"><icon-save />保存</a-button>
      </template>
      <a-form :model="form" layout="vertical" style="max-width: 600px">
        <a-form-item label="AppID"><a-input v-model="form.app_id" /></a-form-item>
        <a-form-item label="AppSecret"><a-input-password v-model="form.app_secret" /></a-form-item>
        <a-form-item label="消息推送启用">
          <a-switch v-model="form.push_enabled" />
        </a-form-item>
        <a-form-item label="模板消息ID"><a-input v-model="form.template_id" /></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const form = reactive({ app_id: '', app_secret: '', push_enabled: true, template_id: '' })
const handleSave = () => Message.success('保存成功')
</script>
