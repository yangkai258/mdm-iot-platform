<template>
  <div class="container">
    <a-card>
      <template #title><icon-team /> 多宠物协作</template>
      <a-form layout="vertical" style="max-width: 500px">
        <a-form-item label="协作模式"><a-select v-model="form.mode"><a-option value="parallel">并行</a-option><a-option value="sequential">顺序</a-option></a-select></a-form-item>
        <a-form-item><a-button type="primary" @click="handleSave">保存</a-button></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const form = reactive({ mode: 'parallel' })
const handleSave = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
