<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-warning /> 情感预警配置</a-space>
      </template>

      <a-row :gutter="16">
        <a-col :span="14">
          <a-card title="预警规则配置">
            <a-form :model="formData" layout="vertical">
              <a-form-item label="预警开关">
                <a-switch v-model="formData.enabled" />
              </a-form-item>
              <a-divider>预警条件</a-divider>
              <a-form-item label="持续低落次数">
                <a-space>
                  <a-input-number v-model="formData.consecutiveCount" :min="2" :max="10" />
                  <span>次</span>
                </a-space>
              </a-form-item>
              <a-form-item label="时间窗口">
                <a-space>
                  <a-input-number v-model="formData.timeWindow" :min="1" :max="24" />
                  <span>小时内</span>
                </a-space>
              </a-form-item>
              <a-form-item label="严重程度阈值">
                <a-slider v-model="formData.severityThreshold" :min="0" :max="100" :show-input="true" />
              </a-form-item>
              <a-divider>通知配置</a-divider>
              <a-form-item label="通知家属">
                <a-switch v-model="formData.notifyFamily" />
              </a-form-item>
              <a-form-item v-if="formData.notifyFamily" label="通知方式">
                <a-checkbox-group v-model="formData.notifyChannels">
                  <a-checkbox value="push">推送</a-checkbox>
                  <a-checkbox value="sms">短信</a-checkbox>
                  <a-checkbox value="email">邮件</a-checkbox>
                </a-checkbox-group>
              </a-form-item>
              <a-form-item label="通知模板">
                <a-textarea v-model="formData.template" placeholder="请输入通知模板" />
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="handleSave">保存规则</a-button>
                <a-button @click="handleTest">测试通知</a-button>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>

        <a-col :span="10">
          <a-card title="预警记录">
            <a-table :columns="alertColumns" :data="alerts" size="small" :pagination="pagination">
              <template #severity="{ record }">
                <a-tag :color="getSeverityColor(record.severity)">{{ record.severity }}</a-tag>
              </template>
            </a-table>
          </a-card>

          <a-card title="家庭成员" style="margin-top: 16px">
            <a-list>
              <a-list-item v-for="member in familyMembers" :key="member.id">
                <a-list-item-meta :title="member.name" :description="member.phone" />
                <template #actions>
                  <a-switch v-model="member.notifyEnabled" size="small" />
                </template>
              </a-list-item>
            </a-list>
          </a-card>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const formData = reactive({
  enabled: true, consecutiveCount: 3, timeWindow: 12, severityThreshold: 40,
  notifyFamily: true, notifyChannels: ['push'], template: '您的宠物情绪持续低落，请关注'
})

const pagination = reactive({ current: 1, pageSize: 5, total: 10 })
const alertColumns = [
  { title: '时间', dataIndex: 'time', width: 150 },
  { title: '严重程度', slotName: 'severity', width: 100 },
  { title: '持续时间', dataIndex: 'duration' },
  { title: '状态', dataIndex: 'status' }
]
const alerts = ref([
  { time: '2026-03-28 10:00', severity: 'warning', duration: '2小时', status: '已通知' },
  { time: '2026-03-27 15:00', severity: 'info', duration: '1小时', status: '已处理' }
])

const familyMembers = ref([
  { id: 1, name: '张三', phone: '138****8888', notifyEnabled: true },
  { id: 2, name: '李四', phone: '139****6666', notifyEnabled: true }
])

const getSeverityColor = (s) => ({ warning: 'orange', major: 'red', info: 'blue' }[s] || 'gray')
const handleSave = () => { }
const handleTest = () => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
