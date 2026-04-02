<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="container">
    <a-card class="general-card" title="积分设置">
      <template #extra>
        <a-button type="primary" @click="handleSave"><icon-save />保存</a-button>
      </template>
      <a-form :model="form" layout="vertical" style="max-width: 600px">
        <a-form-item label="积分启用">
          <a-switch v-model="form.enabled" />
        </a-form-item>
        <a-form-item label="消费获得积分比"><a-input-number v-model="form.consume_ratio" :min="0" :precision="2" style="width:100%" /></a-form-item>
        <a-form-item label="积分抵现比"><a-input-number v-model="form.redeem_ratio" :min="0" :precision="2" style="width:100%" /></a-form-item>
        <a-form-item label="积分上限"><a-input-number v-model="form.max_points" :min="0" style="width:100%" /></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const form = reactive({ enabled: true, consume_ratio: 1, redeem_ratio: 0.01, max_points: 99999 })
const handleSave = () => Message.success('保存成功')
</script>
