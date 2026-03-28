<template>
  <div class="container">
    <a-card>
      <template #title><icon-up /> 告警升级配置</template>
      <a-form layout="vertical" style="max-width: 500px">
        <a-form-item label="升级条件"><a-input v-model="form.condition" placeholder="如：持续30分钟" /></a-form-item>
        <a-form-item label="升级级别"><a-select v-model="form.level"><a-option value="warning">警告</a-option><a-option value="major">重要</a-option></a-select></a-form-item>
        <a-form-item><a-button type="primary" @click="handleSave">保存</a-button></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const form = reactive({ condition: '', level: '' })
const handleSave = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
