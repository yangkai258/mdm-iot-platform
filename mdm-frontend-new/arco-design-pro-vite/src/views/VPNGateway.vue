<template>
  <div class="container">
    <a-card>
      <template #title><icon-global /> VPN网关</template>
      <a-form layout="vertical" style="max-width: 500px">
        <a-form-item label="VPN类型"><a-select v-model="form.type"><a-option value="pptp">PPTP</a-option><a-option value="l2tp">L2TP</a-option></a-select></a-form-item>
        <a-form-item label="服务器"><a-input v-model="form.server" /></a-form-item>
        <a-form-item><a-button type="primary" @click="handleSave">保存</a-button></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const form = reactive({ type: 'pptp', server: '' })
const handleSave = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
