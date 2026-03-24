<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="名称"><a-input v-model="form.name" placeholder="请输入" /></a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleCreate">新建</a-button>
    </div>
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
      <template #actions="{ record }">
        <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
        <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
      </template>
    </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @before-ok="handleSubmit" @cancel="modalVisible = false">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="名称"><a-input v-model="form.name" placeholder="请输入" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>

import { reactive } from 'vue'
import { Message } from '@arco-design/web-vue'

const form = reactive({
  appId: 'wx1234567890abcdef',
  appSecret: '************************',
  serverUrl: 'https://api.example.com/wechat/callback',
  token: 'my_wechat_token_2026',
  encodingAesKey: '',
  encryptType: 'safe',
  autoReply: '1',
  templateMsg: '1',
  customMenu: '1',
  connected: true
})

const copyUrl = () => {
  navigator.clipboard.writeText(form.serverUrl).then(() => {
    Message.success('已复制到剪贴板')
  })
}

const testConnection = () => {
  Message.loading({ content: '正在测试连接...', duration: 1500 })
  setTimeout(() => {
    Message.success('连接正常')
  }, 1600)
}

const handleSave = () => {
  if (!form.appId || !form.appSecret || !form.token) {
    Message.error('请填写必填项')
    return
  }
  Message.success('保存成功')
}

const handleReset = () => {
  Object.assign(form, {
    appId: '',
    appSecret: '',
    serverUrl: 'https://api.example.com/wechat/callback',
    token: '',
    encodingAesKey: '',
    encryptType: 'safe',
    autoReply: '1',
    templateMsg: '1',
    customMenu: '1',
    connected: false
  })
  Message.info('已重置')
}

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
