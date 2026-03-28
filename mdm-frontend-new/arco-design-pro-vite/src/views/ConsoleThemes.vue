<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-palette /> 控制台主题定制</a-space>
      </template>

      <a-row :gutter="16">
        <a-col :span="16">
          <a-card title="主题预览">
            <a-space wrap>
              <a-card v-for="theme in themes" :key="theme.id" :class="{ selected: selectedTheme?.id === theme.id }" class="theme-card" @click="handleSelectTheme(theme)">
                <div class="theme-preview" :style="{ background: theme.colors.primary }">
                  <div class="theme-header" :style="{ background: theme.colors.header }"></div>
                  <div class="theme-body">
                    <div class="theme-sidebar" :style="{ background: theme.colors.sidebar }"></div>
                    <div class="theme-content" :style="{ background: theme.colors.background }"></div>
                  </div>
                </div>
                <template #title>{{ theme.name }}</template>
              </a-card>
            </a-space>
          </a-card>

          <a-card title="自定义配色" style="margin-top: 16px">
            <a-form :model="colorForm" layout="vertical">
              <a-row :gutter="16">
                <a-col :span="6">
                  <a-form-item label="主色">
                    <a-color-picker v-model="colorForm.primary" />
                  </a-form-item>
                </a-col>
                <a-col :span="6">
                  <a-form-item label="侧边栏">
                    <a-color-picker v-model="colorForm.sidebar" />
                  </a-form-item>
                </a-col>
                <a-col :span="6">
                  <a-form-item label="头部">
                    <a-color-picker v-model="colorForm.header" />
                  </a-form-item>
                </a-col>
                <a-col :span="6">
                  <a-form-item label="背景">
                    <a-color-picker v-model="colorForm.background" />
                  </a-form-item>
                </a-col>
              </a-row>
            </a-form>
          </a-card>
        </a-col>

        <a-col :span="8">
          <a-card title="主题配置">
            <a-form :model="selectedTheme" layout="vertical">
              <a-form-item label="主题名称">
                <a-input v-model="selectedTheme.name" />
              </a-form-item>
              <a-form-item label="是否系统主题">
                <a-switch v-model="selectedTheme.isSystem" />
              </a-form-item>
              <a-form-item>
                <a-button type="primary" long @click="handleSave">保存主题</a-button>
              </a-form-item>
              <a-form-item>
                <a-button long @click="handleReset">恢复默认</a-button>
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

const themes = ref([
  { id: 1, name: '深蓝科技', isSystem: true, colors: { primary: '#165qff', header: '#001234', sidebar: '#001529', background: '#f0f2f5' } },
  { id: 2, name: '浅色简约', isSystem: true, colors: { primary: '#165qff', header: '#ffffff', sidebar: '#ffffff', background: '#f0f2f5' } }
])
const selectedTheme = ref(themes.value[0])
const colorForm = reactive({ primary: '#165qff', header: '#001234', sidebar: '#001529', background: '#f0f2f5' })

const handleSelectTheme = (theme) => { selectedTheme.value = theme }
const handleSave = () => { }
const handleReset = () => { }
</script>

<style scoped>
.container { padding: 16px; }
.theme-card { width: 200px; cursor: pointer; }
.theme-card.selected { border: 2px solid #165qff; }
.theme-preview { height: 120px; border-radius: 4px; overflow: hidden; }
.theme-header { height: 20px; }
.theme-body { display: flex; height: 100px; }
.theme-sidebar { width: 30%; }
.theme-content { width: 70%; }
</style>
