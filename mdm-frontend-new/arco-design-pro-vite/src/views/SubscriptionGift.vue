<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-gift /> 订阅赠送</a-space>
      </template>

      <a-tabs default-active-key="gift">
        <a-tab-pane key="gift" title="赠送订阅">
          <a-row :gutter="16">
            <a-col :span="12">
              <a-card title="选择套餐">
                <a-form :model="formData" layout="vertical">
                  <a-form-item label="选择套餐" required>
                    <a-select v-model="formData.planId" placeholder="选择要赠送的套餐">
                      <a-option value="basic">基础版 - ¥29/月</a-option>
                      <a-option value="pro">专业版 - ¥99/月</a-option>
                      <a-option value="enterprise">企业版 - ¥299/月</a-option>
                    </a-select>
                  </a-form-item>
                  <a-form-item label="接收人" required>
                    <a-input v-model="formData.recipientId" placeholder="输入用户ID或手机号" />
                  </a-form-item>
                  <a-form-item label="赠送类型">
                    <a-radio-group v-model="formData.giftType">
                      <a-radio value="trial">试用赠送</a-radio>
                      <a-radio value="paid">付费赠送</a-radio>
                    </a-radio-group>
                  </a-form-item>
                  <a-form-item>
                    <a-button type="primary" long @click="handleGenerate">生成赠送链接</a-button>
                  </a-form-item>
                </a-form>
              </a-card>
            </a-col>
            <a-col :span="12">
              <a-card title="赠送记录">
                <a-table :columns="giftColumns" :data="giftRecords" :pagination="false" />
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>

        <a-tab-pane key="my" title="我的赠送">
          <a-table :columns="myColumns" :data="myGifts" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const formData = reactive({ planId: '', recipientId: '', giftType: 'trial' })

const giftColumns = [
  { title: '接收人', dataIndex: 'recipient' },
  { title: '套餐', dataIndex: 'plan' },
  { title: '类型', dataIndex: 'type' },
  { title: '状态', dataIndex: 'status' },
  { title: '过期时间', dataIndex: 'expiresAt' }
]
const giftRecords = ref([
  { recipient: '张三', plan: '专业版', type: '试用', status: '待领取', expiresAt: '2026-04-05' }
])

const myColumns = [
  { title: '套餐', dataIndex: 'plan' },
  { title: '发送时间', dataIndex: 'sentAt' },
  { title: '领取状态', dataIndex: 'status' },
  { title: '操作', slotName: 'action' }
]
const myGifts = ref([])

const handleGenerate = () => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
