<template>
  <div class="container">
    <a-card class="general-card" title="会员设置">
      <template #extra>
        <a-button type="primary" @click="handleSave"><icon-save />保存</a-button>
      </template>
      <a-form :model="form" layout="vertical" style="max-width: 600px">
        <a-form-item label="注册赠送积分">
          <a-input-number v-model="form.register_points" :min="0" style="width: 100%" />
        </a-form-item>
        <a-form-item label="消费积分比例">
          <a-input-number v-model="form.consume_ratio" :min="0" :precision="2" style="width: 100%" />
        </a-form-item>
        <a-form-item label="积分抵扣比例">
          <a-input-number v-model="form.redeem_ratio" :min="0" :precision="2" style="width: 100%" />
        </a-form-item>
        <a-form-item label="会员有效期（天）">
          <a-input-number v-model="form.validity_days" :min="0" style="width: 100%" />
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const form = reactive({ register_points: 100, consume_ratio: 1, redeem_ratio: 0.01, validity_days: 365 })

const handleSave = () => { Message.success('保存成功') }

onMounted(() => {})
</script>
