<template>
  <div class="container">
    <a-card>
      <template #title><icon-user /> 会员详情</template>
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="会员ID">{{ member.id }}</a-descriptions-item>
        <a-descriptions-item label="昵称">{{ member.nickname }}</a-descriptions-item>
        <a-descriptions-item label="手机">{{ member.phone }}</a-descriptions-item>
        <a-descriptions-item label="等级"><a-tag>{{ member.level }}</a-tag></a-descriptions-item>
        <a-descriptions-item label="注册时间">{{ member.registerTime }}</a-descriptions-item>
        <a-descriptions-item label="状态"><a-tag color="green">正常</a-tag></a-descriptions-item>
      </a-descriptions>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const member = reactive({ id: 'M001', nickname: '张三', phone: '138****8888', level: '黄金', registerTime: '2025-01-15' })
</script>
<style scoped>.container { padding: 16px; }</style>
