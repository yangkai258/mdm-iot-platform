<template>
  <div class="workbench-container">
    <!-- 欢迎区域 -->
    <a-row :gutter="[16, 16]">
      <a-col :span="24">
        <a-card class="welcome-card">
          <div class="welcome-content">
            <div class="welcome-text">
              <h2>👋 您好，{{ username }}！</h2>
              <p class="welcome-sub">{{ greeting }}，祝您工作愉快！</p>
            </div>
            <div class="welcome-info">
              <span class="info-item">📅 {{ today }}</span>
              <span class="info-item">🌤️ 晴 26°C</span>
            </div>
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 快捷操作入口 -->
    <a-row :gutter="[16, 16]" style="margin-top: 16px;">
      <a-col :span="24">
        <a-card title="快捷操作" class="section-card">
          <template #extra>
            <a-link @click="$router.push('/portal')">管理</a-link>
          </template>
          <div class="quick-grid">
            <div class="quick-item" @click="$router.push('/devices')">
              <icon-desktop class="quick-icon" />
              <span>设备列表</span>
            </div>
            <div class="quick-item" @click="$router.push('/members/list')">
              <icon-user-group class="quick-icon" />
              <span>会员管理</span>
            </div>
            <div class="quick-item" @click="$router.push('/alert')">
              <icon-alert class="quick-icon" />
              <span>告警中心</span>
            </div>
            <div class="quick-item" @click="$router.push('/notifications/list')">
              <icon-message class="quick-icon" />
              <span>消息通知</span>
            </div>
            <div class="quick-item" @click="$router.push('/ota/packages')">
              <icon-upload class="quick-icon" />
              <span>固件管理</span>
            </div>
            <div class="quick-item" @click="$router.push('/policies/list')">
              <icon-safe class="quick-icon" />
              <span>合规策略</span>
            </div>
            <div class="quick-item" @click="$router.push('/pet/config')">
              <icon-book class="quick-icon" />
              <span>宠物配置</span>
            </div>
            <div class="quick-item" @click="$router.push('/permissions/roles')">
              <icon-settings class="quick-icon" />
              <span>权限管理</span>
            </div>
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 待办事项 & 最新动态 -->
    <a-row :gutter="[16, 16]" style="margin-top: 16px;">
      <!-- 待办事项 -->
      <a-col :xs="24" :lg="12">
        <a-card title="待办事项" class="section-card">
          <template #extra>
            <a-link @click="$router.push('/process')">查看全部</a-link>
          </template>
          <a-list :data="todos" :bordered="false">
            <template #empty>
              <a-empty description="暂无待办" />
            </template>
            <template #item="{ item }">
              <a-list-item>
                <a-list-item-meta :title="item.title" :description="item.time">
                  <template #avatar>
                    <a-checkbox v-model="item.done" @change="toggleTodo(item)" />
                  </template>
                </a-list-item-meta>
                <template #actions>
                  <a-tag :color="getPriorityColor(item.priority)">{{ getPriorityText(item.priority) }}</a-tag>
                </template>
              </a-list-item>
            </template>
          </a-list>
        </a-card>
      </a-col>

      <!-- 最新动态 -->
      <a-col :xs="24" :lg="12">
        <a-card title="最新动态" class="section-card">
          <template #extra>
            <a-link @click="$router.push('/activity')">查看全部</a-link>
          </template>
          <a-timeline>
            <a-timeline-item v-for="item in activities" :key="item.id" :color="item.color">
              <div class="activity-item">
                <div class="activity-title">{{ item.title }}</div>
                <div class="activity-desc">{{ item.desc }}</div>
                <div class="activity-time">{{ item.time }}</div>
              </div>
            </a-timeline-item>
          </a-timeline>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import dayjs from 'dayjs'

const username = ref('管理员')
const todos = ref([
  { id: 1, title: '审批设备注册申请', time: '今天 10:00', priority: 1, done: false },
  { id: 2, title: '处理设备离线告警', time: '今天 09:30', priority: 2, done: false },
  { id: 3, title: '更新固件版本', time: '今天 09:00', priority: 3, done: true },
  { id: 4, title: '审核会员升级申请', time: '昨天 18:00', priority: 2, done: false }
])

const activities = ref([
  { id: 1, title: '新设备上线', desc: '设备 MDM-001 成功注册', time: '10:30', color: 'green' },
  { id: 2, title: '固件升级', desc: '10 台设备完成 OTA 升级', time: '09:45', color: 'blue' },
  { id: 3, title: '告警触发', desc: '设备 MDM-003 电量低于 15%', time: '09:20', color: 'red' },
  { id: 4, title: '会员注册', desc: '新会员 张三 注册成功', time: '08:55', color: 'gray' },
  { id: 5, title: '策略更新', desc: '合规策略 v2.1 已发布', time: '08:30', color: 'arcoblue' }
])

const today = computed(() => dayjs().format('YYYY年MM月DD日 dddd'))

const greeting = computed(() => {
  const hour = dayjs().hour()
  if (hour < 12) return '上午'
  if (hour < 14) return '中午'
  if (hour < 18) return '下午'
  return '晚上'
})

const getPriorityColor = (p) => ({ 1: 'red', 2: 'orange', 3: 'blue' }[p] || 'gray')
const getPriorityText = (p) => ({ 1: '紧急', 2: '重要', 3: '普通' }[p] || '普通')

const toggleTodo = (item) => {
  item.done = !item.done
}

onMounted(() => {
  const userStr = localStorage.getItem('user')
  if (userStr) {
    try {
      const user = JSON.parse(userStr)
      username.value = user.username || user.name || '管理员'
    } catch {}
  }
})
</script>

<style scoped>
.workbench-container {
  padding: 0;
}

.welcome-card {
  background: linear-gradient(135deg, #165dff 0%, #4080ff 100%);
  border: none;
  border-radius: 8px;
}

.welcome-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #fff;
}

.welcome-text h2 {
  margin: 0 0 4px;
  font-size: 22px;
  color: #fff;
}

.welcome-sub {
  margin: 0;
  font-size: 14px;
  opacity: 0.85;
}

.welcome-info {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
}

.info-item {
  font-size: 13px;
  opacity: 0.9;
}

.section-card {
  border-radius: 8px;
}

.quick-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
}

.quick-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px;
  border-radius: 8px;
  background: #f7f8fa;
  cursor: pointer;
  transition: all 0.3s;
}

.quick-item:hover {
  background: #e6f4ff;
  transform: translateY(-2px);
}

.quick-icon {
  font-size: 28px;
  color: #165dff;
  margin-bottom: 8px;
}

.quick-item span {
  font-size: 13px;
  color: #4e5969;
}

.activity-item {
  padding-bottom: 4px;
}

.activity-title {
  font-size: 14px;
  font-weight: 500;
  color: #1f2329;
}

.activity-desc {
  font-size: 13px;
  color: #646a73;
  margin: 2px 0;
}

.activity-time {
  font-size: 12px;
  color: #86909c;
}

@media (max-width: 768px) {
  .quick-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  .welcome-content {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  .welcome-info {
    align-items: flex-start;
  }
}
</style>
