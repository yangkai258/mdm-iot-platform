<template>
  <div class="login-container">
    <a-card class="login-card">
      <template #title>
        <div class="login-title">
          <span>MDM 控制中台</span>
        </div>
      </template>
      
      <a-form
        :model="form"
        :rules="rules"
        @submit="handleLogin"
      >
        <a-form-item field="username" label="用户名">
          <a-input v-model="form.username" placeholder="请输入用户名" allow-clear>
            <template #prefix>
              <span>👤</span>
            </template>
          </a-input>
        </a-form-item>
        
        <a-form-item field="password" label="密码">
          <a-input
            v-model="form.password"
            placeholder="请输入密码"
            password
            allow-clear
            @keydown.enter="handleLogin"
          >
            <template #prefix>
              <span>🔒</span>
            </template>
          </a-input>
        </a-form-item>
        
        <a-form-item>
          <a-button type="primary" html-type="submit" long :loading="loading">
            登录
          </a-button>
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'

const router = useRouter()
const loading = ref(false)

const form = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名' }],
  password: [{ required: true, message: '请输入密码' }]
}

const handleLogin = async () => {
  if (!form.username || !form.password) {
    Message.warning('请填写用户名和密码')
    return
  }
  
  loading.value = true
  try {
    const res = await fetch('/api/v1/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(form)
    })
    const data = await res.json()
    
    if (data.code === 0) {
      localStorage.setItem('token', data.data.token)
      localStorage.setItem('user', JSON.stringify(data.data.user))
      Message.success('登录成功')
      router.push('/dashboard')
    } else {
      Message.error(data.message || '登录失败')
    }
  } catch (e) {
    Message.error('网络错误')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 400px;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.login-title {
  text-align: center;
  font-size: 20px;
  font-weight: bold;
  color: #165dff;
}
</style>
