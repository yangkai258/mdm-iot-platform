<template>
  <div class="container">
    <a-card>
      <template #title><icon-migrate /> 数据迁移</template>
      <a-steps :current="step" style="margin-bottom: 24px">
        <a-step title="选择源" />
        <a-step title="配置映射" />
        <a-step title="执行迁移" />
      </a-steps>
      <a-space>
        <a-button @click="step = Math.max(0, step - 1)">上一步</a-button>
        <a-button type="primary" @click="step = Math.min(2, step + 1)">下一步</a-button>
      </a-space>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const step = ref(0)
</script>
<style scoped>.container { padding: 16px; }</style>
