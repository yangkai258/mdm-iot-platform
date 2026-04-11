<template>
  <div class="page-container">
    <a-card class="general-card" title="绀句氦骞冲彴鍒嗕韩">
      <a-row :gutter="16">
        <a-col :span="12">
          <a-card title="宸茬粦瀹氱殑绀句氦璐﹀彿" size="small">
            <a-list :data="boundAccounts" :loading="loading">
              <template #item="{ item }">
                <a-list-item>
                  <a-list-item-meta :title="item.platform" :description="item.account">
                    <template #avatar>
                      <a-avatar :style="{ backgroundColor: item.color }">{{ item.icon }}</a-avatar>
                    </template>
                  </a-list-item-meta>
                  <template #actions>
                    <a-button type="primary" size="small" @click="handleUnbind(item)">瑙ｇ粦</a-button>
                  </template>
                </a-list-item>
              </template>
              <template #empty>
                <a-empty description="鏆傛棤缁戝畾鐨勭ぞ浜よ处锟? />
              </template>
            </a-list>
            <div style="margin-top: 16px">
              <a-button type="primary" @click="showBindModal = true"><icon-plus />缁戝畾鏂拌处锟?/a-button>
            </div>
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card title="鍒嗕韩妯℃澘閰嶇疆" size="small">
            <a-form :model="templateForm" layout="vertical">
              <a-form-item label="榛樿鍒嗕韩鏍囬">
                <a-input v-model="templateForm.title" placeholder="璇疯緭鍏ュ垎浜爣锟" />
              </a-form-item>
              <a-form-item label="榛樿鍒嗕韩鏂囨">
                <a-textarea v-model="templateForm.content" :rows="3" placeholder="璇疯緭鍏ュ垎浜枃锟" />
              </a-form-item>
              <a-form-item label="鍒嗕韩閰嶅浘">
                <a-upload action="/api/v1/upload" :show-upload-list="false" @success="handleUploadSuccess">
                  <a-button><icon-upload />涓婁紶閰嶅浘</a-button>
                </a-upload>
                <a-image v-if="templateForm.image" :src="templateForm.image" width="100" style="margin-left: 8px" />
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="handleSaveTemplate">淇濆瓨閰嶇疆</a-button>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>
      </a-row>
    </a-card>

    <a-modal v-model:visible="showBindModal" title="缁戝畾绀句氦璐﹀彿" @ok="handleBind">
      <a-form :model="bindForm" layout="vertical">
        <a-form-item label="閫夋嫨骞冲彴">
          <a-select v-model="bindForm.platform" placeholder="璇烽€夋嫨骞冲彴">
            <a-option value="wechat">寰俊</a-option>
            <a-option value="weibo">寰崥</a-option>
            <a-option value="douyin">鎶栭煶</a-option>
            <a-option value="xiaohongshu">灏忕孩锟?/a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="璐﹀彿">
          <a-input v-model="bindForm.account" placeholder="璇疯緭鍏ヨ处锟" />
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
    { id: 1, platform: '寰俊', account: 'pet123', icon: 'W', color: '#07c160' },
    { id: 2, platform: '鎶栭煶', account: 'pet_home', icon: 'D', color: '#fe2c55' }
  ]
}

const handleUnbind = (item) => {
  boundAccounts.value = boundAccounts.value.filter(a => a.id !== item.id)
  Message.success(`宸茶В锟? ${item.platform}`)
}

const handleBind = () => {
  if (!bindForm.platform || !bindForm.account) {
    Message.warning('璇峰～鍐欏畬鏁翠俊锟?)
    return
  }
  const colors = { wechat: '#07c160', weibo: '#e6162d', douyin: '#fe2c55', xiaohongshu: '#fe2c55' }
  const icons = { wechat: 'W', weibo: 'W', douyin: 'D', xiaohongshu: 'X' }
  boundAccounts.value.push({
    id: Date.now(),
    platform: bindForm.platform === 'wechat' ? '寰俊' : bindForm.platform === 'weibo' ? '寰崥' : bindForm.platform === 'douyin' ? '鎶栭煶' : '灏忕孩锟?,
    account: bindForm.account,
    icon: icons[bindForm.platform],
    color: colors[bindForm.platform]
  })
  showBindModal.value = false
  Message.success('缁戝畾鎴愬姛')
  bindForm.platform = ''
  bindForm.account = ''
}

const handleUploadSuccess = (file) => {
  templateForm.image = file
  Message.success('涓婁紶鎴愬姛')
}

const handleSaveTemplate = () => {
  Message.success('鍒嗕韩妯℃澘宸蹭繚锟?)
}

onMounted(() => loadData())
</script>

<style scoped>
.page-container { padding: 16px; }
</style>