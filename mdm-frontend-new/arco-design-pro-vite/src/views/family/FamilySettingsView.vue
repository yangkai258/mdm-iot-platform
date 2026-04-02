<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="家庭设置">
      <template #extra>
        <a-button type="primary" @click="handleSave"><icon-save />保存</a-button>
      </template>
      <a-form :model="form" layout="vertical" style="max-width: 600px">
        <a-form-item label="家庭名称"><a-input v-model="form.name" /></a-form-item>
        <a-form-item label="成员数量上限"><a-input-number v-model="form.max_members" :min="1" style="width:100%" /></a-form-item>
        <a-form-item label="儿童模式">
          <a-switch v-model="form.child_mode" />
        </a-form-item>
        <a-form-item label="描述"><a-textarea v-model="form.description" :rows="3" /></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const form = reactive({ name: '', max_members: 5, child_mode: false, description: '' })
const handleSave = () => Message.success('保存成功')
</script>
