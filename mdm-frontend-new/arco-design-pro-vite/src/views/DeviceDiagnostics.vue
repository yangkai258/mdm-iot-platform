<template>
  <div class="container">
    <a-card>
      <template #title><icon-tool /> 设备诊断</template>
      <a-button type="primary" @click="handleDiagnose">开始诊断</a-button>
      <a-divider />
      <a-card v-if="result" title="诊断结果">
        <a-descriptions :column="1" bordered>
          <a-descriptions-item label="网络状态">{{ result.network }}</a-descriptions-item>
          <a-descriptions-item label="固件状态">{{ result.firmware }}</a-descriptions-item>
          <a-descriptions-item label="电池状态">{{ result.battery }}</a-descriptions-item>
        </a-descriptions>
      </a-card>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const result = ref('')
const handleDiagnose = () => { result.value = { network: '正常', firmware: '已是最新', battery: '85%' } }
</script>
<style scoped>.container { padding: 16px; }</style>
