<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-language /> 多语言切换</a-space>
      </template>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-card title="选择语言">
            <a-list>
              <a-list-item v-for="lang in languages" :key="lang.code" :class="{ selected: currentLang === lang.code }" @click="handleSelectLang(lang)">
                <a-list-item-meta :title="lang.name" :description="lang.nativeName" />
                <template #actions>
                  <a-badge v-if="currentLang === lang.code" status="success" text="当前" />
                  <a-tag v-if="lang.progress < 100" color="orange">{{ lang.progress }}%</a-tag>
                </template>
              </a-list-item>
            </a-list>
          </a-card>
        </a-col>

        <a-col :span="12">
          <a-card title="翻译进度">
            <a-progress :percent="translationProgress" :show-text="true" />
            <a-divider />
            <a-space direction="vertical" fill>
              <a-descriptions :column="1" bordered>
                <a-descriptions-item label="已翻译">2,580 条</a-descriptions-item>
                <a-descriptions-item label="未翻译">420 条</a-descriptions-item>
                <a-descriptions-item label="总条目">3,000 条</a-descriptions-item>
              </a-descriptions>
            </a-space>
          </a-card>

          <a-card title="翻译贡献" style="margin-top: 16px">
            <a-space direction="vertical" fill>
              <a-button type="primary" long>
                <template #icon><icon-edit /></template>
                帮助完善翻译
              </a-button>
              <a-button long>
                <template #icon><icon-upload /></template>
                上传翻译文件
              </a-button>
            </a-space>
          </a-card>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const currentLang = ref('zh-CN')
const translationProgress = ref(86)

const languages = ref([
  { code: 'zh-CN', name: '简体中文', nativeName: '简体中文', progress: 100 },
  { code: 'en-US', name: 'English', nativeName: 'English', progress: 100 },
  { code: 'ja-JP', name: '日本語', nativeName: '日本語', progress: 86 },
  { code: 'ko-KR', name: '한국어', nativeName: '한국어', progress: 72 }
])

const handleSelectLang = (lang) => { currentLang.value = lang.code }
</script>

<style scoped>
.container { padding: 16px; }
.selected { background: #e6f7ff; cursor: pointer; }
</style>
