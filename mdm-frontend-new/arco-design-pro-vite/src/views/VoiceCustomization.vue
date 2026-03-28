<template>
  <div class="voice-custom-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="声音模板" :value="12" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="已定制" :value="3" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="使用次数" :value="1568" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <span>声音定制</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="templates" title="声音模板">
          <a-row :gutter="16">
            <a-col :span="6" v-for="voice in voices" :key="voice.id">
              <a-card size="small" class="voice-card">
                <div class="voice-icon">{{ voice.icon }}</div>
                <div class="voice-name">{{ voice.name }}</div>
                <div class="voice-type">{{ voice.typeText }}</div>
                <div class="voice-price">{{ voice.price === 0 ? '免费' : '¥' + voice.price }}</div>
                <a-button type="primary" size="small" @click="handleUse(voice)">使用</a-button>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
        
        <a-tab-pane key="custom" title="自定义录制">
          <a-form layout="vertical" style="max-width: 600px;">
            <a-form-item label="选择基础声音">
              <a-select v-model="customForm.baseVoice" placeholder="选择基础声音">
                <a-option value="male">男声基础</a-option>
                <a-option value="female">女声基础</a-option>
                <a-option value="child">童声基础</a-option>
              </a-select>
            </a-form-item>
            <a-form-item label="音调调整">
              <a-slider v-model="customForm.pitch" :marks="{0:'低',50:'中',100:'高'}" />
            </a-form-item>
            <a-form-item label="语速调整">
              <a-slider v-model="customForm.speed" :marks="{0:'慢',50:'中',100:'快'}" />
            </a-form-item>
            <a-form-item label="情感强度">
              <a-slider v-model="customForm.emotion" :marks="{0:'平静',50:'适中',100:'饱满'}" />
            </a-form-item>
            <a-form-item label="预览文本">
              <a-input v-model="customForm.previewText" placeholder="请输入要预览的文本" />
            </a-form-item>
            <a-form-item>
              <a-space>
                <a-button type="primary" @click="handlePreview">预览</a-button>
                <a-button @click="handleSave">保存</a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts'>
import { ref, reactive } from 'vue';

const voices = ref([
  { id: 1, name: '标准男声', icon: '🎤', type: 'male', typeText: '男声', price: 0, usage: 560 },
  { id: 2, name: '温柔女声', icon: '🎙️', type: 'female', typeText: '女声', price: 9, usage: 420 },
  { id: 3, name: '萌系童声', icon: '🧒', type: 'child', typeText: '童声', price: 19, usage: 280 },
  { id: 4, name: '磁性低音', icon: '🎸', type: 'male', typeText: '男声', price: 29, usage: 156 },
]);

const customForm = reactive({
  baseVoice: 'female',
  pitch: 50,
  speed: 50,
  emotion: 50,
  previewText: '你好，我是你的AI宠物伙伴',
});

const handleUse = (voice: any) => {};
const handlePreview = () => {};
const handleSave = () => {};
</script>

<style scoped>
.voice-custom-container { padding: 20px; }
.voice-card { text-align: center; margin-bottom: 12px; }
.voice-icon { font-size: 48px; margin-bottom: 8px; }
.voice-name { font-weight: bold; }
.voice-type { color: #86909c; font-size: 12px; }
.voice-price { color: #F53F3F; font-weight: bold; margin: 8px 0; }
</style>
