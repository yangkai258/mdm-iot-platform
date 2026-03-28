<template>
  <div class="container">
    <a-card>
      <template #title><icon-lock /> 角色权限</template>
      <a-form layout="vertical" style="max-width: 500px">
        <a-form-item label="角色"><a-select v-model="form.role"><a-option value="admin">管理员</a-option></a-select></a-form-item>
        <a-form-item label="权限"><a-checkbox-group v-model="form.permissions"><a-checkbox value="read">读取</a-checkbox><a-checkbox value="write">写入</a-checkbox></a-checkbox-group></a-form-item>
        <a-form-item><a-button type="primary" @click="handleSave">保存</a-button></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const form = reactive({ role: '', permissions: [] })
const handleSave = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
