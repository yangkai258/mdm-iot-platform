<template>
  <div class="container">
    <a-card>
      <template #title><icon-language /> 国际化配置</template>
      <a-form layout="vertical" style="max-width: 500px">
        <a-form-item label="默认语言"><a-select v-model="form.defaultLang"><a-option value="zh-CN">中文</a-option><a-option value="en-US">English</a-option></a-select></a-form-item>
        <a-form-item><a-button type="primary" @click="handleSave">保存</a-button></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const form = reactive({ defaultLang: 'zh-CN' })
const handleSave = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
