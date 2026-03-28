<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-global /> 数据驻留合规</a-space>
      </template>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-card title="区域配置">
            <a-form :model="formData" layout="vertical">
              <a-form-item label="数据驻留区域">
                <a-select v-model="formData.region">
                  <a-option value="cn">中国大陆</a-option>
                  <a-option value="us">美国</a-option>
                  <a-option value="eu">欧盟</a-option>
                </a-select>
              </a-form-item>
              <a-form-item label="数据分类">
                <a-checkbox-group v-model="formData.dataTypes">
                  <a-checkbox value="personal">个人信息</a-checkbox>
                  <a-checkbox value="device">设备数据</a-checkbox>
                  <a-checkbox value="behavior">行为数据</a-checkbox>
                </a-checkbox-group>
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="handleSave">保存配置</a-button>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>

        <a-col :span="12">
          <a-card title="合规状态">
            <a-descriptions :column="1" bordered>
              <a-descriptions-item label="合规状态">
                <a-badge status="success" text="符合要求" />
              </a-descriptions-item>
              <a-descriptions-item label="最后审计">2026-03-28</a-descriptions-item>
              <a-descriptions-item label="审计结果">通过</a-descriptions-item>
            </a-descriptions>
          </a-card>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const formData = reactive({ region: 'cn', dataTypes: ['personal', 'device'] })

const handleSave = () => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
