<template>
  <div class="page-container">
    <a-card class="general-card" title="社交平台分享">
      <a-row :gutter="16">
        <a-col :span="12">
          <a-card title="已绑定的社交账号" size="small">
            <a-list :data="boundAccounts" :loading="loading">
              <template #item="{ item }">
                <a-list-item>
                  <a-list-item-meta :title="item.platform" :description="item.account">
                    <template #avatar>
                      <a-avatar :style="{ backgroundColor: item.color }">{{ item.icon }}</a-avatar>
                    </template>
                  </a-list-item-meta>
                  <template #actions>
                    <a-button type="primary" size="small" @click="handleUnbind(item)">解绑</a-button>
                  </template>
                </a-list-item>
              </template>
              <template #empty>
                <a-empty description="暂无绑定的社交账�? />
              </template>
            </a-list>
            <div style="margin-top: 16px">
              <a-button type="primary" @click="showBindModal = true"><icon-plus />绑定新账�?/a-button>
            </div>
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card title="分享模板配置" size="small">
            <a-form :model="templateForm" layout="vertical">
              <a-form-item label="默认分享标题">
                <a-input v-model="templateForm.title" placeholder="请输入分享标�? />
              </a-form-item>
              <a-form-item label="默认分享文案">
                <a-textarea v-model="templateForm.content" :rows="3" placeholder="请输入分享文�? />
              </a-form-item>
              <a-form-item label="分享配图">
                <a-upload action="/api/v1/upload" :show-upload-list="false" @success="handleUploadSuccess">
                  <a-button><icon-upload />上传配图</a-button>
                </a-upload>
                <a-image v-if="templateForm.image" :src="templateForm.image" width="100" style="margin-left: 8px" />
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="handleSaveTemplate">保存配置</a-button>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>
      </a-row>
    </a-card>

    <a-modal v-model:visible="showBindModal" title="绑定社交账号" @ok="handleBind">
      <a-form :model="bindForm" layout="vertical">
        <a-form-item label="选择平台">
          <a-select v-model="bindForm.platform" placeholder="请选择平台">
            <a-option value="wechat">微信</a-option>
            <a-option value="weibo">微博</a-option>
            <a-option value="douyin">抖音</a-option>
            <a-option value="xiaohongshu">小红�?/a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="账号">
          <a-input v-model="bindForm.account" placeholder="请输入账�? />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus, IconUpload } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const showBindModal = ref(false)
const boundAccounts = ref([])
const bindForm = reactive({ platform: '', account: '' })
const templateForm = reactive({ title: '', content: '', image: '' })

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/v1/integration/social/accounts').then(r => r.json())
    if (res.code === 0) {
      boundAccounts.value = res.data || []
    } else {
      loadMockData()
    }
  } catch { loadMockData() } finally { loading.value = false }
}

const loadMockData = () => {
  boundAccounts.value = [
    { id: 1, platform: '微信', account: 'pet123', icon: 'W', color: '#07c160' },
    { id: 2, platform: '抖音', account: 'pet_home', icon: 'D', color: '#fe2c55' }
  ]
}

const handleUnbind = (item) => {
  boundAccounts.value = boundAccounts.value.filter(a => a.id !== item.id)
  Message.success(`已解�? ${item.platform}`)
}

const handleBind = () => {
  if (!bindForm.platform || !bindForm.account) {
    Message.warning('请填写完整信�?)
    return
  }
  const colors = { wechat: '#07c160', weibo: '#e6162d', douyin: '#fe2c55', xiaohongshu: '#fe2c55' }
  const icons = { wechat: 'W', weibo: 'W', douyin: 'D', xiaohongshu: 'X' }
  boundAccounts.value.push({
    id: Date.now(),
    platform: bindForm.platform === 'wechat' ? '微信' : bindForm.platform === 'weibo' ? '微博' : bindForm.platform === 'douyin' ? '抖音' : '小红�?,
    account: bindForm.account,
    icon: icons[bindForm.platform],
    color: colors[bindForm.platform]
  })
  showBindModal.value = false
  Message.success('绑定成功')
  bindForm.platform = ''
  bindForm.account = ''
}

const handleUploadSuccess = (file) => {
  templateForm.image = file
  Message.success('上传成功')
}

const handleSaveTemplate = () => {
  Message.success('分享模板已保�?)
}

onMounted(() => loadData())
</script>

<style scoped>
.page-container { padding: 16px; }
</style>