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
  { key: '/dashboard', label: '📊 设备大盘' },
  { key: '/devices', label: '📱 设备列表' },
  {
    key: '/ota',
    label: '🔄 OTA升级',
    children: [
      { key: '/ota/packages', label: '固件包管理' },
      { key: '/ota/deployments', label: '部署任务' }
    ]
  },
  {
    key: '/alert',
    label: '⚠️ 告警中心',
    children: [
      { key: '/alerts/rules', label: '告警规则' },
      { key: '/alerts/list', label: '告警列表' },
      { key: '/alerts/settings', label: '告警设置' },
      { key: '/alert/notification', label: '告警通知' },
      { key: '/alert/history', label: '告警历史' },
      { key: '/alert/notification-logs', label: '通知日志' }
    ]
  },
  {
    key: '/members',
    label: '👥 会员管理',
    children: [
      { key: '/members', label: '会员列表' },
      { key: '/members/points', label: '积分管理' },
      { key: '/members/coupons', label: '优惠券' },
      { key: '/members/promotions', label: '促销活动' }
    ]
  },
  {
    key: '/pet',
    label: '🐾 宠物管理',
    children: [
      { key: '/pet', label: '宠物配置' },
      { key: '/pet/console', label: '宠物控制台' },
      { key: '/pet/conversations', label: '会话记录' }
    ]
  },
  {
    key: '/ai',
    label: '🤖 AI管理',
    children: [
      { key: '/ai/quality-dashboard', label: 'AI 质量监控' },
      { key: '/ai/models', label: 'AI 模型管理' },
      { key: '/ai/behavior-log', label: 'AI 行为日志' },
      { key: '/ai/model-version', label: '模型版本管理' },
      { key: '/ai/sandbox', label: 'AI 沙箱测试' },
      { key: '/ai/training', label: 'AI 训练任务' },
      { key: '/ai/model-publish', label: '模型发布工作流' }
    ]
  },
  {
    key: '/emotion',
    label: '💝 情绪管理',
    children: [
      { key: '/emotion/recognize', label: '情绪识别配置' },
      { key: '/emotion/logs', label: '情绪日志' },
      { key: '/emotion/response-config', label: '响应配置' },
      { key: '/emotion/reports', label: '情绪报告' }
    ]
  },
  {
    key: '/digital-twin',
    label: '🪞 数字孪生',
    children: [
      { key: '/digital-twin/vitals', label: '生命体征仪表盘' },
      { key: '/digital-twin/vitals-chart', label: '实时体征曲线' },
      { key: '/digital-twin/history', label: '历史回放' },
      { key: '/digital-twin/behavior', label: '行为预测' }
    ]
  },
  {
    key: '/health',
    label: '❤️ 健康管理',
    children: [
      { key: '/health/warnings', label: '健康预警' },
      { key: '/health/exercise', label: '运动统计' },
      { key: '/health/sleep', label: '睡眠分析' },
      { key: '/health/report', label: '健康报告' }
    ]
  },
  {
    key: '/family',
    label: '👨‍👩‍👧 家庭管理',
    children: [
      { key: '/family/members', label: '家庭成员' },
      { key: '/family/child-mode', label: '儿童模式' },
      { key: '/family/elder-mode', label: '老人陪伴模式' },
      { key: '/family/album', label: '家庭相册' },
      { key: '/family/settings', label: '家庭设置' }
    ]
  },
  {
    key: '/market',
    label: '🎨 内容生态',
    children: [
      { key: '/market/emoticons', label: '表情包市场' },
      { key: '/market/actions', label: '动作资源库' },
      { key: '/market/voices', label: '声音定制' }
    ]
  },
  {
    key: '/knowledge',
    label: '📚 知识库',
    children: [
      { key: '/knowledge', label: '知识列表' }
    ]
  },
  {
    key: '/miniclaw',
    label: '⚡ MiniClaw固件',
    children: [
      { key: '/miniclaw/firmwares', label: '固件列表' }
    ]
  },
  {
    key: '/developer',
    label: '🛠️ 开发者平台',
    children: [
      { key: '/developer/apps', label: '开发者应用管理' },
      { key: '/developer/stats', label: 'API 使用统计' }
    ]
  },
  {
    key: '/integration',
    label: '🔌 第三方集成',
    children: [
      { key: '/integration/smarthome', label: '智能家居' },
      { key: '/integration/pet-hospital', label: '宠物医疗' },
      { key: '/integration/pet-shop', label: '宠物用品商城' }
    ]
  },
  {
    key: '/tech',
    label: '⚙️ 技术架构',
    children: [
      { key: '/tech/edge-models', label: '端侧推理模型' },
      { key: '/tech/mesh-network', label: 'BLE Mesh 网络' },
      { key: '/tech/device-ota', label: '设备 OTA 优化' }
    ]
  },
  {
    key: '/system-group',
    label: '⚙️ 系统管理',
    children: [
      { key: '/system/monitor', label: '服务监控' },
      { key: '/system/logs', label: '操作日志' },
      { key: '/policies/list', label: '策略列表' },
      { key: '/policies/configs', label: '策略配置' },
      { key: '/policies/compliance', label: '合规规则' },
      { key: '/performance/dashboard', label: '性能仪表盘' },
      { key: '/performance/cache', label: '缓存管理' }
    ]
  },
  {
    key: '/tenants',
    label: '🏢 租户管理',
    children: [
      { key: '/tenants/approval', label: '租户入驻审核' },
      { key: '/tenants/management', label: '租户系统管理' },
      { key: '/tenants/settings', label: '租户设置' },
      { key: '/tenants/public-archives', label: '公共档案' },
      { key: '/tenants/system-info', label: '系统信息' }
    ]
  },
  {
    key: '/permissions',
    label: '🔐 多维权限',
    children: [
      { key: '/permissions/roles', label: '角色管理' },
      { key: '/permissions/menus', label: '菜单管理' },
      { key: '/permissions/groups', label: '权限组管理' },
      { key: '/permissions/data-config', label: '数据权限配置' },
      { key: '/permissions/api', label: 'API 权限' }
    ]
  },
  {
    key: '/security',
    label: '🛡️ 安全中心',
    children: [
      { key: '/security/permission', label: '权限分配' },
      { key: '/security/data-permission', label: '数据权限' },
      { key: '/security/certificate', label: '证书管理' },
      { key: '/security/device', label: '设备安全' },
      { key: '/security/ldap', label: 'LDAP配置' },
      { key: '/security/user-sync', label: '用户同步' },
      { key: '/security/settings', label: '安全设置' },
      { key: '/security/privacy', label: '数据隐私' },
      { key: '/security/audit-log', label: '审计日志' },
      { key: '/security/two-factor', label: '双因素认证' },
      { key: '/security/sessions', label: '会话管理' },
      { key: '/security/audit', label: '安全审计' }
    ]
  },
  {
    key: '/globalization',
    label: '🌍 全球化',
    children: [
      { key: '/globalization', label: '全球化设置' },
      { key: '/globalization/region', label: '区域管理' },
      { key: '/globalization/timezone', label: '时区设置' },
      { key: '/globalization/data-residency', label: '数据驻留规则' },
      { key: '/globalization/ai-node', label: '区域AI节点' },
      { key: '/globalization/sync-status', label: '跨区域同步状态' }
    ]
  },
  {
    key: '/notifications',
    label: '📢 通知管理',
    children: [
      { key: '/notifications/list', label: '推送通知' },
      { key: '/notifications/announcements', label: '公告管理' },
      { key: '/notifications/templates', label: '通知模板' }
    ]
  },
  {
    key: '/i18n',
    label: '🌐 国际化',
    children: [
      { key: '/i18n/translations', label: '翻译管理' },
      { key: '/i18n/region-settings', label: '区域设置' }
    ]
  },
  {
    key: '/app',
    label: '📲 App管理',
    children: [
      { key: '/app/devices', label: 'App设备列表' },
      { key: '/miniapp/home', label: '小程序首页' },
      { key: '/miniapp/devices', label: '小程序设备' }
    ]
  },
  {
    key: '/owner',
    label: '👤 个人信息',
    children: [
      { key: '/owner/profile', label: '个人资料' }
    ]
  },
  {
    key: '/research',
    label: '🔬 研究平台',
    children: [
      { key: '/research/data', label: '数据集管理' },
      { key: '/research/experiment', label: 'AI 行为实验' }
    ]
  },
  { key: '/test-modals', label: '🎨 UI测试' }
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
