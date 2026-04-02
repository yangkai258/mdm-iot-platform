<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="标签自动清理">
      <template #extra>
        <a-button type="primary" @click="handleSave"><icon-save />保存</a-button>
      </template>
      <a-form :model="form" layout="vertical" style="max-width: 600px">
        <a-form-item label="自动清理启用">
          <a-switch v-model="form.enabled" />
        </a-form-item>
        <a-form-item label="清理周期（天）"><a-input-number v-model="form.cleanup_days" :min="1" style="width:100%" /></a-form-item>
        <a-form-item label="低于此使用次数则清理"><a-input-number v-model="form.min_use_count" :min="0" style="width:100%" /></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const form = reactive({ enabled: false, cleanup_days: 30, min_use_count: 5 })
const handleSave = () => Message.success('保存成功')
</script>
