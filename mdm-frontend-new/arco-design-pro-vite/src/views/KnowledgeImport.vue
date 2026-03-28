<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-upload /> 知识批量导入</a-space>
      </template>

      <a-steps :current="currentStep">
        <a-step title="上传文件" />
        <a-step title="数据预览" />
        <a-step title="导入配置" />
        <a-step title="导入完成" />
      </a-steps>

      <div class="step-content">
        <a-card v-if="currentStep === 0" title="上传文件">
          <a-space direction="vertical" fill style="width: 100%">
            <a-button>
              <template #icon><icon-download /></template>
              下载模板
            </a-button>
            <a-upload draggable accept=".xlsx,.csv" :limit="1" />
            <a-form :model="uploadForm" layout="vertical" style="margin-top: 16px">
              <a-form-item label="文件编码">
                <a-select v-model="uploadForm.encoding">
                  <a-option value="utf8">UTF-8</a-option>
                  <a-option value="gbk">GBK</a-option>
                </a-select>
              </a-form-item>
            </a-form>
            <a-button type="primary" @click="handleNext">下一步</a-button>
          </a-space>
        </a-card>

        <a-card v-if="currentStep === 1" title="数据预览">
          <a-table :columns="previewColumns" :data="previewData" size="small" :pagination="false" />
          <a-space style="margin-top: 16px">
            <a-button @click="handlePrev">上一步</a-button>
            <a-button type="primary" @click="handleNext">下一步</a-button>
          </a-space>
        </a-card>

        <a-card v-if="currentStep === 2" title="导入配置">
          <a-form :model="importConfig" layout="vertical">
            <a-form-item label="知识分类">
              <a-select v-model="importConfig.category">
                <a-option value="pet">宠物知识</a-option>
                <a-option value="health">健康知识</a-option>
              </a-select>
            </a-form-item>
            <a-form-item label="标签">
              <a-select v-model="importConfig.tags" multiple />
            </a-form-item>
            <a-form-item label="是否覆盖已有">
              <a-switch v-model="importConfig.overwrite" />
            </a-form-item>
          </a-form>
          <a-space style="margin-top: 16px">
            <a-button @click="handlePrev">上一步</a-button>
            <a-button type="primary" @click="handleImport">开始导入</a-button>
          </a-space>
        </a-card>

        <a-card v-if="currentStep === 3" title="导入结果">
          <a-result status="success" title="导入成功">
            <template #sub-title>
              成功导入 {{ importResult.success }} 条，失败 {{ importResult.failed }} 条
            </template>
            <template #extra>
              <a-button @click="handleExportFailed">导出失败记录</a-button>
              <a-button type="primary" @click="handleFinish">完成</a-button>
            </template>
          </a-result>
        </a-card>
      </div>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const currentStep = ref(0)
const uploadForm = reactive({ encoding: 'utf8' })
const importConfig = reactive({ category: '', tags: [], overwrite: false })
const importResult = reactive({ success: 95, failed: 5 })

const previewColumns = [
  { title: '标题', dataIndex: 'title' },
  { title: '内容', dataIndex: 'content' },
  { title: '标签', dataIndex: 'tags' }
]
const previewData = ref([
  { title: '宠物饮食注意', content: '...', tags: '饮食' }
])

const handleNext = () => { currentStep.value++ }
const handlePrev = () => { currentStep.value-- }
const handleImport = () => { currentStep.value = 3 }
const handleExportFailed = () => { }
const handleFinish = () => { currentStep.value = 0 }
</script>

<style scoped>
.container { padding: 16px; }
.step-content { margin-top: 24px; }
</style>
