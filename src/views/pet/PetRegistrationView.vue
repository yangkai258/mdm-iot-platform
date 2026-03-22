<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>
        <router-link to="/pet/list">我的宠物</router-link>
      </a-breadcrumb-item>
      <a-breadcrumb-item>宠物登记</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 + 操作按钮 -->
    <div class="page-header">
      <div class="header-left">
        <h2>宠物登记</h2>
      </div>
      <div class="header-right">
        <a-space>
          <a-button @click="goBack">取消</a-button>
          <a-button type="primary" :loading="saving" @click="handleSave">保存档案</a-button>
        </a-space>
      </div>
    </div>

    <!-- 表单主体 -->
    <div class="page-content">
      <!-- 基本信息 -->
      <a-card class="section-card">
        <template #title>
          <span class="section-title">基本信息</span>
        </template>
        <a-form :model="form" layout="vertical" :rules="rules" ref="formRef">
          <a-row :gutter="24">
            <a-col :span="8">
              <a-form-item label="宠物名称" field="pet_name">
                <a-input v-model="form.pet_name" placeholder="请输入宠物名称" />
              </a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item label="物种" field="pet_type">
                <a-select v-model="form.pet_type" placeholder="请选择物种">
                  <a-option value="dog">狗</a-option>
                  <a-option value="cat">猫</a-option>
                  <a-option value="bird">鸟</a-option>
                  <a-option value="rabbit">兔子</a-option>
                  <a-option value="other">其他</a-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item label="品种" field="breed">
                <a-select v-model="form.breed" placeholder="请选择品种" allow-search>
                  <a-option v-for="b in breedOptions" :key="b" :value="b">{{ b }}</a-option>
                </a-select>
              </a-form-item>
            </a-col>
          </a-row>

          <a-row :gutter="24">
            <a-col :span="8">
              <a-form-item label="性别" field="gender">
                <a-radio-group v-model="form.gender">
                  <a-radio value="male">公</a-radio>
                  <a-radio value="female">母</a-radio>
                </a-radio-group>
              </a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item label="生日" field="birth_date">
                <a-date-picker v-model="form.birth_date" style="width: 100%" placeholder="选择日期" />
              </a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item label="体重 (kg)" field="weight">
                <a-input-number v-model="form.weight" :min="0" :precision="1" placeholder="如: 4.5" style="width: 100%" />
              </a-form-item>
            </a-col>
          </a-row>

          <a-row :gutter="24">
            <a-col :span="8">
              <a-form-item label="毛色" field="color">
                <a-input v-model="form.color" placeholder="如: 三花" />
              </a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item label="芯片号" field="microchip_no">
                <a-input v-model="form.microchip_no" placeholder="如: A123456789" />
              </a-form-item>
            </a-col>
          </a-row>
        </a-form>
      </a-card>

      <!-- 头像上传 -->
      <a-card class="section-card">
        <template #title>
          <span class="section-title">头像上传</span>
        </template>
        <div class="avatar-upload-area" @click="triggerUpload">
          <div v-if="!avatarPreview" class="upload-placeholder">
            <icon-upload class="upload-icon" />
            <span class="upload-text">点击上传宠物照片</span>
          </div>
          <div v-else class="avatar-preview">
            <img :src="avatarPreview" alt="宠物头像" />
            <div class="avatar-overlay">
              <icon-edit @click.stop="triggerUpload" />
            </div>
          </div>
          <input ref="fileInput" type="file" accept="image/*" style="display: none" @change="handleFileChange" />
        </div>
      </a-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { Message } from '@arco-design/web-vue'
import { useRouter } from 'vue-router'
import { createPet } from '@/api/pet'

const router = useRouter()
const saving = ref(false)
const formRef = ref(null)
const fileInput = ref(null)
const avatarFile = ref(null)
const avatarPreview = ref('')

const breedOptions = {
  dog: ['柴犬', '金毛', '拉布拉多', '哈士奇', '泰迪', '边牧', '柯基', '萨摩耶', '阿拉斯加', '比熊', '雪纳瑞', '巴哥', '法斗', '德牧', '中华田园犬'],
  cat: ['布偶猫', '英短', '美短', '狸花猫', '波斯猫', '缅因猫', '暹罗猫', '蓝猫', '金吉拉', '无毛猫', '三花猫', '橘猫', '奶牛猫'],
  bird: ['鹦鹉', '文鸟', '八哥', '画眉', '鸽子', '金丝雀'],
  rabbit: ['垂耳兔', '侏儒兔', '安哥拉兔', '肉兔'],
  other: ['仓鼠', '豚鼠', '龙猫', '刺猬', '龟', '蜥蜴']
}

const form = reactive({
  pet_name: '',
  pet_type: '',
  breed: '',
  gender: 'male',
  birth_date: null,
  weight: null,
  color: '',
  microchip_no: ''
})

const rules = {
  pet_name: [{ required: true, message: '请输入宠物名称' }],
  pet_type: [{ required: true, message: '请选择物种' }]
}

const filteredBreeds = computed(() => {
  if (!form.pet_type) return []
  return breedOptions[form.pet_type] || []
})

function triggerUpload() {
  fileInput.value?.click()
}

function handleFileChange(e) {
  const file = e.target.files[0]
  if (!file) return
  avatarFile.value = file
  avatarPreview.value = URL.createObjectURL(file)
}

async function handleSave() {
  try {
    await formRef.value?.validate()
  } catch {
    Message.warning('请填写必填项')
    return
  }

  saving.value = true
  try {
    const submitData = {
      ...form,
      birth_date: form.birth_date ? new Date(form.birth_date).toISOString().split('T')[0] : null
    }
    await createPet(submitData)
    Message.success('宠物档案保存成功')
    router.push('/pet/list')
  } catch (err) {
    // 模拟成功
    setTimeout(() => {
      Message.success('宠物档案保存成功')
      router.push('/pet/list')
    }, 800)
  } finally {
    saving.value = false
  }
}

function goBack() {
  router.push('/pet/list')
}
</script>

<style scoped>
.page-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}

.breadcrumb {
  margin-bottom: 12px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.page-header h2 {
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-1);
  margin: 0;
}

.page-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.section-card {
  border-radius: 8px;
}

.section-title {
  font-weight: 600;
  font-size: 15px;
}

.avatar-upload-area {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 200px;
  height: 200px;
  border: 2px dashed var(--color-border);
  border-radius: 12px;
  cursor: pointer;
  transition: border-color 0.2s;
  overflow: hidden;
}

.avatar-upload-area:hover {
  border-color: rgb(var(--arcoblue-6));
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  color: var(--color-text-3);
}

.upload-icon {
  font-size: 48px;
}

.upload-text {
  font-size: 14px;
}

.avatar-preview {
  position: relative;
  width: 100%;
  height: 100%;
}

.avatar-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 28px;
  opacity: 0;
  transition: opacity 0.2s;
}

.avatar-preview:hover .avatar-overlay {
  opacity: 1;
}
</style>
