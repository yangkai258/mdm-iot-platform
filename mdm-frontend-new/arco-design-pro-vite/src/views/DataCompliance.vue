<template>
  <div class="container">
    <a-card>
      <template #title><icon-shield /> 数据合规中心</template>
      <a-result status="success" title="合规状态：符合要求">
        <template #sub-title>
          最近审计时间: 2026-03-28 | 审计结果: 通过
        </template>
      </a-result>
    </a-card>
  </div>
</template>
<script setup>
</script>
<style scoped>.container { padding: 16px; }</style>
