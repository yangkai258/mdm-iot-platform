<template>
  <div class="container">
    <a-card>
      <template #title><icon-environment /> 仿真环境</template>
      <a-form layout="vertical" style="max-width: 500px">
        <a-form-item label="环境类型"><a-select v-model="form.type"><a-option value="indoor">室内</a-option><a-option value="outdoor">户外</a-option></a-select></a-form-item>
        <a-form-item label="地形"><a-select v-model="form.terrain"><a-option value="flat">平坦</a-option><a-option value="complex">复杂</a-option></a-select></a-form-item>
        <a-form-item><a-button type="primary" @click="handleStart">启动仿真</a-button></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const form = reactive({ type: 'indoor', terrain: 'flat' })
const handleStart = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
