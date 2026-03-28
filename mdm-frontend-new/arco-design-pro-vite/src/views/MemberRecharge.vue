<template>
  <div class="container">
    <a-card>
      <template #title><icon-wallet /> 会员充值</template>
      <a-form layout="vertical" style="max-width: 500px">
        <a-form-item label="充值金额"><a-input-number v-model="form.amount" :min="1" style="width: 100%" /></a-form-item>
        <a-form-item label="支付方式"><a-radio-group v-model="form.paymentMethod"><a-radio value="alipay">支付宝</a-radio><a-radio value="wechat">微信</a-radio></a-radio-group></a-form-item>
        <a-form-item><a-button type="primary" @click="handleRecharge">立即充值</a-button></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const form = reactive({ amount: 100, paymentMethod: 'alipay' })
const handleRecharge = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
