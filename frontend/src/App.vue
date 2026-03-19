<template>
  <div id="app">
    <!-- 登录后显示 -->
    <div v-if="isLoggedIn" class="layout">
      <div class="sidebar">
        <div class="logo">MDM 控制台</div>
        <div class="menu">
          <div 
            v-for="item in menuItems" 
            :key="item.key"
            :class="['menu-item', { active: selectedKeys === item.key }]"
            @click="handleMenuClick(item.key)"
          >
            {{ item.label }}
          </div>
        </div>
      </div>
      <div class="main">
        <div class="header">
          <span class="username">{{ username }}</span>
          <button @click="handleLogout">退出</button>
        </div>
        <div class="content">
          <router-view />
        </div>
      </div>
    </div>
    
    <!-- 未登录 -->
    <router-view v-else />
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()
const selectedKeys = ref('/dashboard')
const username = ref('Admin')
const isLoggedIn = ref(false)

const menuItems = [
  { key: '/dashboard', label: '设备大盘' },
  { key: '/devices', label: '设备列表' },
  { key: '/ota', label: 'OTA固件' },
  { key: '/alert', label: '告警管理' },
  { key: '/system/monitor', label: '服务监控' },
  { key: '/system/logs', label: '操作日志' },
]

onMounted(() => {
  checkLogin()
})

const checkLogin = () => {
  const token = localStorage.getItem('token')
  isLoggedIn.value = !!token
  if (!token) {
    router.push('/login')
    return
  }
  const user = localStorage.getItem('user')
  if (user) {
    try {
      const userData = JSON.parse(user)
      username.value = userData.nickname || userData.username || 'Admin'
    } catch (e) {}
  }
  selectedKeys.value = route.path
}

watch(() => route.path, (newPath) => {
  selectedKeys.value = newPath
})

const handleMenuClick = (key) => {
  selectedKeys.value = key
  router.push(key)
}

const handleLogout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  isLoggedIn.value = false
  router.push('/login')
}
</script>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: Arial, sans-serif; }
#app { min-height: 100vh; }
.layout { display: flex; min-height: 100vh; }
.sidebar { width: 220px; background: #001529; color: #fff; }
.logo { height: 64px; display: flex; align-items: center; justify-content: center; font-size: 18px; font-weight: bold; }
.menu-item { padding: 14px 24px; cursor: pointer; }
.menu-item:hover { background: rgba(255,255,255,0.1); }
.menu-item.active { background: #1890ff; }
.main { flex: 1; }
.header { height: 64px; background: #fff; padding: 0 16px; display: flex; justify-content: flex-end; align-items: center; gap: 16px; box-shadow: 0 1px 4px rgba(0,0,0,0.1); }
.content { margin: 16px; }
button { padding: 6px 16px; cursor: pointer; }
</style>
