<template>
  <div class="container">
    <a-row :gutter="16">
      <a-col :span="12">
        <a-card title="动画效果选择">
          <a-form layout="vertical">
            <a-form-item label="动画类型">
              <a-radio-group v-model="config.type">
                <a-radio value="particle">粒子爆炸</a-radio>
                <a-radio value="light">光效渐变</a-radio>
                <a-radio value="badge">徽章升级</a-radio>
                <a-radio value="star">星星闪烁</a-radio>
              </a-radio-group>
            </a-form-item>
            <a-form-item label="动画时长">
              <a-slider v-model="config.duration" :min="1" :max="5" :step="0.5" :show-input="true" />
              <span style="color: #909399">单位：秒</span>
            </a-form-item>
            <a-form-item label="背景颜色">
              <a-color-picker v-model="config.bgColor" />
            </a-form-item>
            <a-form-item label="音效">
              <a-switch v-model="config.soundEnabled" />
              <span style="margin-left: 8px">启用音效</span>
            </a-form-item>
            <a-form-item v-if="config.soundEnabled" label="音效选择">
              <a-select v-model="config.sound" placeholder="请选择音效">
                <a-option value="fanfare">胜利号角</a-option>
                <a-option value="chime">清脆铃声</a-option>
                <a-option value="applause">掌声</a-option>
                <a-option value="custom">自定义上传</a-option>
              </a-select>
            </a-form-item>
            <a-form-item>
              <a-space>
                <a-button type="primary" @click="handleSave">保存配置</a-button>
                <a-button @click="handlePreview">预览动画</a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="动画预览">
          <div class="preview-area" :style="{ backgroundColor: config.bgColor }">
            <div class="level-badge">
              <span class="old-level">{{ previewOldLevel }}</span>
              <icon-arrow-right />
              <span class="new-level">{{ previewNewLevel }}</span>
            </div>
            <div v-if="showAnimation" class="animation-overlay" :class="config.type">
              <div v-for="i in 20" :key="i" class="particle"></div>
            </div>
          </div>
          <a-divider />
          <a-form layout="vertical">
            <a-form-item label="预览等级变化">
              <a-select v-model="previewOldLevel" style="width: 120px">
                <a-option value="白银">白银</a-option>
                <a-option value="黄金">黄金</a-option>
                <a-option value="铂金">铂金</a-option>
              </a-select>
              <icon-arrow-right style="margin: 0 8px" />
              <a-select v-model="previewNewLevel" style="width: 120px">
                <a-option value="黄金">黄金</a-option>
                <a-option value="铂金">铂金</a-option>
                <a-option value="钻石">钻石</a-option>
              </a-select>
            </a-form-item>
          </a-form>
        </a-card>
      </a-col>
    </a-row>
    <a-divider>升级动画列表</a-divider>
    <a-card>
      <template #extra><a-button type="primary" @click="handleAdd"><template #icon><icon-plus /></template>新建动画</a-button></template>
      <a-table :columns="columns" :data="animations">
        <template #type="{ record }"><a-tag>{{ getTypeText(record.type) }}</a-tag></template>
        <template #status="{ record }"><a-tag :color="record.enabled ? 'green' : 'gray'">{{ record.enabled ? '启用' : '禁用' }}</a-tag></template>
        <template #actions="{ record }">
          <a-link @click="handleEdit(record)">编辑</a-link>
          <a-divider direction="vertical" />
          <a-switch v-model="record.enabled" size="small" />
        </template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const config = reactive({
  type: 'particle',
  duration: 2,
  bgColor: '#f0f2f5',
  soundEnabled: true,
  sound: 'fanfare'
})
const previewOldLevel = ref('白银')
const previewNewLevel = ref('黄金')
const showAnimation = ref(false)
const columns = [
  { title: '动画名称', dataIndex: 'name' },
  { title: '类型', slotName: 'type' },
  { title: '时长', dataIndex: 'duration', customRender: ({ text }) => text + 's' },
  { title: '状态', slotName: 'status' },
  { title: '操作', slotName: 'actions', width: 150 }
]
const animations = ref([
  { id: 1, name: '粒子爆炸', type: 'particle', duration: 2, enabled: true },
  { id: 2, name: '光效渐变', type: 'light', duration: 1.5, enabled: true },
  { id: 3, name: '徽章升级', type: 'badge', duration: 2.5, enabled: false }
])
const getTypeText = (t) => ({ particle: '粒子爆炸', light: '光效渐变', badge: '徽章升级', star: '星星闪烁' }[t] || t)
const handleSave = () => {}
const handlePreview = () => { showAnimation.value = true; setTimeout(() => { showAnimation.value = false }, 2000) }
const handleAdd = () => {}
const handleEdit = (r) => {}
</script>
<style scoped>
.container { padding: 16px; }
.preview-area {
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  position: relative;
  overflow: hidden;
}
.level-badge {
  display: flex;
  align-items: center;
  gap: 16px;
  font-size: 24px;
  font-weight: bold;
}
.old-level { color: #c0c4cc; }
.new-level { color: #1659d5; }
.animation-overlay {
  position: absolute;
  inset: 0;
  pointer-events: none;
}
.particle {
  position: absolute;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  animation: particle-float 2s ease-out forwards;
}
@keyframes particle-float {
  0% { transform: translate(0, 0) scale(1); opacity: 1; }
  100% { transform: translate(var(--tx), var(--ty)) scale(0); opacity: 0; }
}
</style>
