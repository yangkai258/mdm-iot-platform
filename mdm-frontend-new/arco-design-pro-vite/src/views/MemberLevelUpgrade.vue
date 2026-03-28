<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-up-circle /> 会员等级自动升级配置</a-space>
      </template>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-card title="升级规则配置">
            <a-form :model="formData" layout="vertical">
              <a-form-item label="自动升级">
                <a-switch v-model="formData.autoUpgrade" />
              </a-form-item>
              <a-form-item label="升级条件">
                <a-radio-group v-model="formData.conditionType">
                  <a-radio value="growth_value">成长值达标</a-radio>
                  <a-radio value="spending">累计消费</a-radio>
                  <a-radio value="both">两者都满足</a-radio>
                </a-radio-group>
              </a-form-item>
              <a-divider>升级门槛配置</a-divider>
              <a-form-item label="白银→黄金">
                <a-space>
                  成长值 >= <a-input-number v-model="formData.goldThreshold" :min="0" style="width: 100px" />
                </a-space>
              </a-form-item>
              <a-form-item label="黄金→铂金">
                <a-space>
                  成长值 >= <a-input-number v-model="formData.platinumThreshold" :min="0" style="width: 100px" />
                </a-space>
              </a-form-item>
              <a-form-item label="铂金→钻石">
                <a-space>
                  成长值 >= <a-input-number v-model="formData.diamondThreshold" :min="0" style="width: 100px" />
                </a-space>
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="handleSave">保存规则</a-button>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card title="升级预览">
            <a-list>
              <a-list-item v-for="level in levels" :key="level.name">
                <a-card size="small">
                  <template #title>
                    <a-space>
                      <a-tag :color="level.color">{{ level.name }}</a-tag>
                      <span>成长值 >= {{ level.threshold }}</span>
                    </a-space>
                  </template>
                  <a-space>
                    <a-tag>折扣: {{ level.discount }}</a-tag>
                    <a-tag>积分倍率: {{ level.pointRate }}x</a-tag>
                  </a-space>
                </a-card>
              </a-list-item>
            </a-list>
          </a-card>

          <a-card title="升级动画配置" style="margin-top: 16px">
            <a-form :model="animationForm" layout="vertical">
              <a-form-item label="启用升级动画">
                <a-switch v-model="animationForm.enabled" />
              </a-form-item>
              <a-form-item label="动画效果">
                <a-select v-model="animationForm.effect">
                  <a-option value="sparkle">光效</a-option>
                  <a-option value="particles">粒子</a-option>
                  <a-option value="badge">徽章</a-option>
                </a-select>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const formData = reactive({
  autoUpgrade: true, conditionType: 'growth_value',
  goldThreshold: 1000, platinumThreshold: 5000, diamondThreshold: 20000
})
const animationForm = reactive({ enabled: true, effect: 'sparkle' })

const levels = ref([
  { name: '普通', color: 'gray', threshold: 0, discount: '无', pointRate: 1 },
  { name: '白银', color: 'blue', threshold: 1000, discount: '95折', pointRate: 1.2 },
  { name: '黄金', color: 'orange', threshold: 5000, discount: '9折', pointRate: 1.5 },
  { name: '铂金', color: 'purple', threshold: 20000, discount: '85折', pointRate: 2 },
  { name: '钻石', color: 'red', threshold: 50000, discount: '8折', pointRate: 3 }
])

const handleSave = () => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
