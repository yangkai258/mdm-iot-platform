<template>
  <div class="personal-container">
    <a-row :gutter="[16, 16]">
      <!-- 左侧：个人设置 -->
      <a-col :xs="24" :lg="8">
        <!-- 用户信息卡片 -->
        <a-card title="个人信息" class="section-card">
          <div class="user-info">
            <a-avatar :size="72" class="user-avatar">
              <icon-user />
            </a-avatar>
            <div class="user-detail">
              <h3>{{ user.username }}</h3>
              <p class="user-role">{{ user.role }}</p>
              <p class="user-dept">{{ user.department }}</p>
            </div>
          </div>
          <a-divider />
          <div class="user-stats">
            <div class="stat-item">
              <span class="stat-value">{{ user.devices }}</span>
              <span class="stat-label">管理设备</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ user.loginDays }}</span>
              <span class="stat-label">登录天数</span>
            </div>
          </div>
        </a-card>

        <!-- 账号安全 -->
        <a-card title="账号安全" class="section-card" style="margin-top: 16px;">
          <a-list :data="securityItems" :bordered="false" size="small">
            <template #item="{ item }">
              <a-list-item>
                <template #avatar><component :is="item.icon" /></template>
                <div class="security-item">
                  <span>{{ item.title }}</span>
                  <a-switch v-model="item.enabled" size="small" />
                </div>
              </a-list-item>
            </template>
          </a-list>
        </a-card>
      </a-col>

      <!-- 右侧：消息通知 & 收藏夹 -->
      <a-col :xs="24" :lg="16">
        <!-- 消息通知 -->
        <a-card title="消息通知" class="section-card">
          <template #extra>
            <a-link @click="$router.push('/notifications/list')">查看全部</a-link>
          </template>
          <a-tabs default-active-key="all">
            <a-tab-pane key="all" title="全部">
              <a-list :data="notifications" :bordered="false">
                <template #empty>
                  <a-empty description="暂无通知" />
                </template>
                <template #item="{ item }">
                  <a-list-item>
                    <a-list-item-meta :title="item.title" :description="item.time">
                      <template #avatar>
                        <a-badge :status="item.status" />
                      </template>
                    </a-list-item-meta>
                    <template #actions>
                      <a-button type="text" size="small" @click="markAsRead(item)" v-if="!item.read">标为已读</a-button>
                    </template>
                  </a-list-item>
                </template>
              </a-list>
            </a-tab-pane>
            <a-tab-pane key="unread" title="未读">
              <a-list :data="notifications.filter(n => !n.read)" :bordered="false">
                <template #empty>
                  <a-empty description="暂无未读通知" />
                </template>
                <template #item="{ item }">
                  <a-list-item>
                    <a-list-item-meta :title="item.title" :description="item.time">
                      <template #avatar>
                        <a-badge status="error" />
                      </template>
                    </a-list-item-meta>
                  </a-list-item>
                </template>
              </a-list>
            </a-tab-pane>
          </a-tabs>
        </a-card>

        <!-- 收藏夹 -->
        <a-card title="我的收藏" class="section-card" style="margin-top: 16px;">
          <template #extra>
            <a-button type="text" @click="showAddFavorite = true">
              <icon-plus /> 添加
            </a-button>
          </template>
          <a-list :data="favorites" :bordered="false">
            <template #empty>
              <a-empty description="暂无收藏" />
            </template>
            <template #item="{ item }">
              <a-list-item>
                <a-list-item-meta
                  :title="item.title"
                  :description="item.path"
                >
                  <template #avatar>
                    <a-avatar :style="{ backgroundColor: item.color }">
                      <component :is="item.icon" />
                    </a-avatar>
                  </template>
                </a-list-item-meta>
                <template #actions>
                  <a-button type="text" size="small" @click="$router.push(item.path)">访问</a-button>
                  <a-button type="text" size="small" @click="removeFavorite(item)">
                    <icon-delete />
                  </a-button>
                </template>
              </a-list-item>
            </template>
          </a-list>
        </a-card>
      </a-col>
    </a-row>

    <!-- 添加收藏弹窗 -->
    <a-modal v-model:visible="showAddFavorite" title="添加收藏" @ok="addFavorite" @cancel="showAddFavorite = false">
      <a-form :model="favoriteForm" layout="vertical">
        <a-form-item label="名称" required>
          <a-input v-model="favoriteForm.title" placeholder="收藏名称" />
        </a-form-item>
        <a-form-item label="路径" required>
          <a-input v-model="favoriteForm.path" placeholder="/devices" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, h } from 'vue'
import dayjs from 'dayjs'

const user = ref({
  username: '管理员',
  role: '系统管理员',
  department: '技术部',
  devices: 128,
  loginDays: 36
})

const securityItems = ref([
  { icon: 'icon-lock', title: '两步验证', enabled: false },
  { icon: 'icon-message', title: '登录通知', enabled: true },
  { icon: 'icon-alert', title: '异常告警', enabled: true }
])

const notifications = ref([
  { id: 1, title: '您有 3 条告警待处理', time: '10分钟前', status: 'error', read: false },
  { id: 2, title: '设备 MDM-001 已上线', time: '30分钟前', status: 'success', read: false },
  { id: 3, title: 'OTA 升级任务已完成', time: '1小时前', status: 'warning', read: true },
  { id: 4, title: '系统将在今晚维护', time: '2小时前', status: 'warning', read: true }
])

const favorites = ref([
  { id: 1, title: '设备列表', path: '/devices', icon: 'icon-desktop', color: '#165dff' },
  { id: 2, title: '会员列表', path: '/members/list', icon: 'icon-user-group', color: '#52c41a' },
  { id: 3, title: '告警中心', path: '/alert', icon: 'icon-alert', color: '#ff4d4f' }
])

const showAddFavorite = ref(false)
const favoriteForm = reactive({ title: '', path: '' })

const markAsRead = (item) => {
  item.read = true
  item.status = 'success'
}

const addFavorite = () => {
  if (favoriteForm.title && favoriteForm.path) {
    favorites.value.push({
      id: Date.now(),
      title: favoriteForm.title,
      path: favoriteForm.path,
      icon: 'icon-star',
      color: '#faad14'
    })
    favoriteForm.title = ''
    favoriteForm.path = ''
    showAddFavorite.value = false
  }
}

const removeFavorite = (item) => {
  const idx = favorites.value.indexOf(item)
  if (idx > -1) favorites.value.splice(idx, 1)
}

onMounted(() => {
  const userStr = localStorage.getItem('user')
  if (userStr) {
    try {
      const u = JSON.parse(userStr)
      user.value.username = u.username || u.name || '管理员'
      user.value.role = u.role || '系统管理员'
    } catch {}
  }
})
</script>

<style scoped>
.personal-container {
  padding: 0;
}

.section-card {
  border-radius: 8px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-avatar {
  background: #165dff;
  font-size: 32px;
}

.user-detail h3 {
  margin: 0 0 4px;
  font-size: 18px;
}

.user-role {
  margin: 0 0 2px;
  font-size: 13px;
  color: #4e5969;
}

.user-dept {
  margin: 0;
  font-size: 12px;
  color: #86909c;
}

.user-stats {
  display: flex;
  gap: 32px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #165dff;
}

.stat-label {
  font-size: 12px;
  color: #86909c;
}

.security-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  padding: 4px 0;
}

.security-item span {
  font-size: 14px;
  color: #1f2329;
}
</style>
