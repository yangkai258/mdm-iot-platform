<template>
  <div class="container">
    <a-card>
      <template #title><icon-location /> 电子围栏配置</template>
      <a-form layout="vertical" style="max-width: 500px">
        <a-form-item label="围栏名称"><a-input v-model="form.name" /></a-form-item>
        <a-form-item label="围栏类型"><a-select v-model="form.type"><a-option value="circle">圆形</a-option><a-option value="polygon">多边形</a-option></a-select></a-form-item>
        <a-form-item label="半径/范围"><a-input-number v-model="form.radius" /> 米</a-form-item>
        <a-form-item><a-button type="primary" @click="handleSave">保存</a-button></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const form = reactive({ name: '', type: 'circle', radius: 500 })
const handleSave = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
