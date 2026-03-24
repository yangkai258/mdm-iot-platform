<template>
  <a-modal
    :visible="visible"
    title="宠物设置"
    :width="480"
    @ok="handleSave"
    @cancel="handleClose"
    :mask-closable="false"
  >
    <a-form :model="form" layout="vertical" class="settings-form">
      <!-- 宠物名称 -->
      <a-form-item label="宠物名称">
        <a-input
          v-model="form.name"
          placeholder="请输入宠物名称"
          :max-length="20"
          show-word-limit
        />
      </a-form-item>

      <!-- 宠物类型 -->
      <a-form-item label="宠物类型">
        <a-select v-model="form.type" placeholder="选择宠物类型">
          <a-option value="cat">🐱 猫咪</a-option>
          <a-option value="dog">🐶 狗狗</a-option>
          <a-option value="bird">🐦 小鸟</a-option>
          <a-option value="rabbit">🐰 兔子</a-option>
          <a-option value="hamster">🐹 仓鼠</a-option>
          <a-option value="fish">🐟 鱼儿</a-option>
          <a-option value="turtle">🐢 乌龟</a-option>
        </a-select>
      </a-form-item>

      <a-divider>调试参数（仅测试用）</a-divider>

      <!-- 心情值 -->
      <a-form-item label="心情值">
        <div class="slider-with-input">
          <a-slider v-model="form.mood" :min="0" :max="100" />
          <a-input-number v-model="form.mood" :min="0" :max="100" style="width: 80px" />
        </div>
      </a-form-item>

      <!-- 能量值 -->
      <a-form-item label="能量值">
        <div class="slider-with-input">
          <a-slider v-model="form.energy" :min="0" :max="100" />
          <a-input-number v-model="form.energy" :min="0" :max="100" style="width: 80px" />
        </div>
      </a-form-item>

      <!-- 饱食度 -->
      <a-form-item label="饱食度">
        <div class="slider-with-input">
          <a-slider v-model="form.hunger" :min="0" :max="100" />
          <a-input-number v-model="form.hunger" :min="0" :max="100" style="width: 80px" />
        </div>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  deviceId: {
    type: String,
    default: ''
  },
  settings: {
    type: Object,
    default: () => ({
      name: '',
      type: 'cat',
      mood: 50,
      energy: 50,
      hunger: 50
    })
  }
})

const emit = defineEmits(['update:visible', 'save'])

const form = reactive({
  name: '',
  type: 'cat',
  mood: 50,
  energy: 50,
  hunger: 50
})

// 同步外部 settings 到表单
watch(() => props.settings, (newSettings) => {
  if (newSettings) {
    form.name = newSettings.name || ''
    form.type = newSettings.type || 'cat'
    form.mood = newSettings.mood ?? 50
    form.energy = newSettings.energy ?? 50
    form.hunger = newSettings.hunger ?? 50
  }
}, { immediate: true, deep: true })

function handleSave() {
  emit('save', { ...form })
}

function handleClose() {
  emit('update:visible', false)
}
</script>

<style scoped>
.settings-form {
  padding: 8px 0;
}

.slider-with-input {
  display: flex;
  align-items: center;
  gap: 16px;
}

.slider-with-input :deep(.arco-slider) {
  flex: 1;
}
</style>
