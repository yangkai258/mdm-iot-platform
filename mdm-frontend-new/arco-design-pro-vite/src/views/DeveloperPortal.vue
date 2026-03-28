<template>
  <div class="container">
    <a-card>
      <template #title><icon-code /> 开发者门户</template>
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="API调用量">125,680</a-descriptions-item>
        <a-descriptions-item label="错误率">0.5%</a-descriptions-item>
        <a-descriptions-item label="应用数量">5</a-descriptions-item>
        <a-descriptions-item label="活跃设备">128</a-descriptions-item>
      </a-descriptions>
      <a-space style="margin-top: 16px">
        <a-button type="primary">查看API文档</a-button>
        <a-button>SDK下载</a-button>
      </a-space>
    </a-card>
  </div>
</template>
<script setup>
</script>
<style scoped>.container { padding: 16px; }</style>
