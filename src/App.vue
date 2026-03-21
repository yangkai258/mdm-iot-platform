<template>
  <div id="app">
    <!-- 登录后显示 -->
    <div v-if="isLoggedIn" class="layout">
      <div class="sidebar">
        <div class="logo">✨ MDM 控制台</div>
        <div class="menu">
          <template v-for="item in menuItems" :key="item.key">
            <!-- 有子菜单 -->
            <div v-if="item.children">
              <div 
                :class="['menu-item', { expanded: expandedKeys.includes(item.key) }]"
                @click="toggleExpand(item.key)"
              >
                <span class="icon">{{ item.icon }}</span>
                <span class="text">{{ item.label }}</span>
                <span class="arrow">▶</span>
              </div>
              <div :class="['menu-children', { expanded: expandedKeys.includes(item.key) }]">
                <div 
                  v-for="child in item.children" 
                  :key="child.key"
                  :class="['menu-item', 'menu-item-child', { active: selectedKeys === child.key }]"
                  @click="handleMenuClick(child.key)"
                >
                  {{ child.label }}
                </div>
              </div>
            </div>
            <!-- 无子菜单 -->
            <div v-else
              :class="['menu-item', { active: selectedKeys === item.key }]"
              @click="handleMenuClick(item.key)"
            >
              <span class="icon">{{ item.icon }}</span>
              <span class="text">{{ item.label }}</span>
            </div>
          </template>
        </div>
      </div>
      <div class="main">
        <div class="header">
          <span class="username">👤 {{ username }}</span>
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
const expandedKeys = ref([])

const toggleExpand = (key) => {
  const index = expandedKeys.value.indexOf(key)
  if (index > -1) {
    expandedKeys.value.splice(index, 1)
  } else {
    expandedKeys.value.push(key)
  }
}

const menuItems = [
  { key: '/dashboard', label: '设备大盘' },
  { key: '/devices', label: '设备列表' },
  {
    key: '/ota',
    label: 'OTA升级',
    children: [
      { key: '/ota/packages', label: '固件包管理' },
      { key: '/ota/deployments', label: '部署任务' }
    ]
  },
  {
    key: '/alert',
    label: '告警管理',
    children: [
      { key: '/alert/rules', label: '告警规则' },
      { key: '/alert/list', label: '告警列表' },
      { key: '/alert/settings', label: '告警设置' },
      { key: '/alert/channels', label: '通知渠道' }
    ]
  },
  {
    key: '/members',
    label: '会员管理',
    children: [
      { key: '/members/points', label: '积分管理' },
      { key: '/members/coupons', label: '优惠券' },
      { key: '/members/promotions', label: '促销活动' }
    ]
  },
  {
    key: '/system',
    label: '系统管理',
    children: [
      { key: '/system/monitor', label: '服务监控' },
      { key: '/system/logs', label: '操作日志' },
      { key: '/policies', label: '策略管理' }
    ]
  },
  {
    key: '/tenants',
    label: '系统管理',
    children: [
      { key: '/tenants/approval', label: '租户入驻审核' },
      { key: '/tenants/management', label: '租户系统管理' },
      { key: '/tenants/public-archives', label: '公共档案' },
      { key: '/tenants/system-info', label: '系统信息' }
    ]
  },
  {
    key: '/permissions',
    label: '多维权限',
    children: [
      { key: '/permissions/groups', label: '权限组管理' },
      { key: '/permissions/data-config', label: '数据权限配置' }
    ]
  },
  {
    key: '/notifications',
    label: '通知管理',
    children: [
      { key: '/notifications/list', label: '推送通知' },
      { key: '/notifications/announcements', label: '公告管理' },
      { key: '/notifications/templates', label: '通知模板' }
    ]
  },
  {
    key: '/apps',
    label: '应用管理',
    children: [
      { key: '/apps/list', label: '应用列表' },
      { key: '/apps/distributions', label: '应用分发' }
    ]
  },
  {
    key: '/policies',
    label: '合规策略',
    children: [
      { key: '/policies/list', label: '策略管理' },
      { key: '/policies/configs', label: '配置文件库' },
      { key: '/policies/compliance-rules', label: '合规规则' },
      { key: '/policies/device-compliance', label: '设备合规状态' }
    ]
  },
  { key: '/test-modals', label: '🎨 UI测试' },
]

onMounted(() => {
  checkLogin()
})

const checkLogin = () => {
  const token = localStorage.getItem('token')
  if (!token) {
    isLoggedIn.value = false
    router.push('/login')
    return
  }
  
  // Token存在，认为已登录，渲染layout
  isLoggedIn.value = true
  
  // 尝试从localStorage获取用户信息
  const user = localStorage.getItem('user')
  if (user) {
    try {
      const userData = JSON.parse(user)
      username.value = userData.nickname || userData.username || 'Admin'
    } catch (e) {
      username.value = 'Admin'
    }
  }
  selectedKeys.value = route.path
}

watch(() => route.path, (newPath) => {
  selectedKeys.value = newPath
  // 每次路由变化都重新检查登录状态
  checkLogin()
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
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; }
#app { min-height: 100vh; }
.layout { display: flex; min-height: 100vh; }

/* 侧边栏 */
.sidebar { 
  width: 240px; 
  background: linear-gradient(180deg, #1a1a2e 0%, #16213e 100%); 
  color: #fff; 
  display: flex; 
  flex-direction: column;
  box-shadow: 2px 0 10px rgba(0,0,0,0.3);
  flex-shrink: 0;
}
.logo { 
  height: 64px; 
  display: flex; 
  align-items: center; 
  justify-content: center; 
  font-size: 18px; 
  font-weight: 600;
  background: rgba(255,255,255,0.05);
  border-bottom: 1px solid rgba(255,255,255,0.1);
  letter-spacing: 1px;
  flex-shrink: 0;
}
.menu { flex: 1; overflow-y: auto; padding: 12px 0; }
.menu::-webkit-scrollbar { width: 4px; }
.menu::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.2); border-radius: 2px; }

/* 子菜单容器 - 解决展开时变窄问题 */
.menu > div { width: 100%; }

.menu-item { 
  padding: 12px 20px; 
  cursor: pointer; 
  display: flex;
  align-items: center;
  gap: 12px;
  transition: all 0.2s ease;
  margin: 2px 12px;
  border-radius: 8px;
  font-size: 14px;
  width: auto !important;
  box-sizing: border-box;
}
.menu-item:hover { 
  background: rgba(255,255,255,0.1); 
}
.menu-item:hover .text {
  transform: translateX(4px);
}
.menu-item.active { 
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}
.menu-item .icon { font-size: 18px; flex-shrink: 0; }
.menu-item .text { flex: 1; }
.menu-item .arrow {
  font-size: 12px;
  opacity: 0.6;
  transition: transform 0.2s;
}
.menu-item.expanded .arrow { transform: rotate(90deg); }

.menu-group-label { 
  padding: 16px 20px 6px; 
  color: rgba(255,255,255,0.4); 
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 1.5px;
  font-weight: 600;
}

/* 子菜单展开/收起 */
.menu-children { 
  overflow: hidden;
  max-height: 0;
  transition: max-height 0.3s ease;
}
.menu-children.expanded { max-height: 500px; }

/* 子菜单项 */
.menu-children .menu-item-child { 
  padding: 10px 16px 10px 48px; 
  font-size: 13px;
  margin: 2px 12px;
  border-radius: 6px;
  opacity: 0.85;
}
.menu-children .menu-item-child:hover { opacity: 1; }
.menu-children .menu-item-child.active { 
  background: rgba(102, 126, 234, 0.3);
  border-left: 3px solid #667eea;
}

/* 主内容区 */
.main { flex: 1; display: flex; flex-direction: column; background: #f5f7fa; }
.header { 
  height: 64px; 
  background: #fff; 
  padding: 0 24px; 
  display: flex; 
  justify-content: flex-end; 
  align-items: center; 
  gap: 20px; 
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
}
.username {
  font-weight: 500;
  color: #333;
}
.header button {
  padding: 8px 20px;
  background: #f5f5f5;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 14px;
}
.header button:hover {
  background: #ff4d4f;
  color: #fff;
}
.content { flex: 1; padding: 20px; overflow-y: auto; }
</style>
