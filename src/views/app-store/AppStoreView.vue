<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>应用生态</a-breadcrumb-item>
      <a-breadcrumb-item>应用商店</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计概览 -->
    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="上架应用" :value="stats.total_apps">
            <template #prefix>🗂️</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="总安装量" :value="stats.total_installs">
            <template #prefix>📥</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="本周新增" :value="stats.weekly_new">
            <template #prefix>🆕</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="平均评分" :value="stats.avg_rating" suffix="⭐" :precision="1">
            <template #prefix>⭐</template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- 分类导航 + 搜索 -->
    <a-card class="section-card" style="margin-top: 16px">
      <a-space wrap align="center">
        <a-tag :class="{ 'cat-active': !filterForm.category }" class="cat-tag" @click="() => { filterForm.category = ''; loadApps() }">全部</a-tag>
        <a-tag v-for="cat in categories" :key="cat.key" :class="{ 'cat-active': filterForm.category === cat.key }" class="cat-tag" @click="() => { filterForm.category = cat.key; loadApps() }">
          {{ cat.icon }} {{ cat.name }}
        </a-tag>
      </a-space>
      <div style="margin-top: 12px; display: flex; gap: 12px; flex-wrap: wrap; align-items: center">
        <a-select v-model="filterForm.sort" placeholder="排序方式" style="width: 140px" @change="loadApps">
          <a-option value="downloads">下载量</a-option>
          <a-option value="rating">评分</a-option>
          <a-option value="newest">最新</a-option>
        </a-select>
        <a-select v-model="filterForm.price" placeholder="价格" style="width: 120px" allow-clear @change="loadApps">
          <a-option value="free">免费</a-option>
          <a-option value="paid">付费</a-option>
        </a-select>
        <a-input-search v-model="filterForm.keyword" placeholder="搜索应用名称/开发者" style="width: 240px" search-button @search="loadApps" @change="e => !e.target.value && loadApps()" />
        <a-button type="primary" @click="showSubmitModal">
          <template #icon><icon-upload /></template>
          提交应用
        </a-button>
      </div>
    </a-card>

    <!-- 应用列表 -->
    <a-card class="section-card" style="margin-top: 16px">
      <a-spin :loading="loading" tip="加载中...">
        <div class="app-grid">
          <div v-for="app in apps" :key="app.app_id" class="app-card">
            <div class="app-banner" :style="{ background: app.banner_color }">
              <span class="app-icon">{{ app.icon }}</span>
              <a-tag v-if="app.is_official" color="blue" size="small" class="app-official-tag">官方</a-tag>
            </div>
            <div class="app-info">
              <div class="app-name" :title="app.name">{{ app.name }}</div>
              <div class="app-dev">{{ app.developer }}</div>
              <div class="app-category">
                <a-tag size="small" :color="getCatColor(app.category)">{{ app.category_text }}</a-tag>
                <span v-if="app.price > 0" class="app-price">¥{{ app.price }}</span>
                <span v-else class="app-price free">免费</span>
              </div>
              <div class="app-rating">
                <icon-star-fill v-for="i in 5" :key="i" :style="{ color: i <= app.rating ? '#ffc53d' : '#d9d9d9', fontSize: '12px' }" />
                <span class="rating-num">{{ app.rating.toFixed(1) }}</span>
                <span class="download-num">{{ formatNum(app.downloads) }}次下载</span>
              </div>
              <div class="app-desc">{{ app.description }}</div>
            </div>
            <div class="app-actions">
              <a-button type="primary" size="small" @click="handleInstall(app)">
                <template #icon><icon-download /></template>
                安装
              </a-button>
              <a-button size="small" @click="viewAppDetail(app)">详情</a-button>
            </div>
          </div>
        </div>
      </a-spin>
    </a-card>

    <!-- 应用详情 -->
    <a-drawer v-model:visible="detailVisible" :title="currentApp?.name" :width="600" unmountOnHide>
      <div class="detail-banner" :style="{ background: currentApp?.banner_color }">
        <span class="detail-icon">{{ currentApp?.icon }}</span>
        <div class="detail-title">{{ currentApp?.name }}</div>
      </div>
      <a-descriptions :column="2" bordered size="large" style="margin: 16px 0">
        <a-descriptions-item label="开发者">{{ currentApp?.developer }}</a-descriptions-item>
        <a-descriptions-item label="版本">{{ currentApp?.version }}</a-descriptions-item>
        <a-descriptions-item label="分类">{{ currentApp?.category_text }}</a-descriptions-item>
        <a-descriptions-item label="价格">{{ currentApp?.price > 0 ? `¥${currentApp?.price}` : '免费' }}</a-descriptions-item>
        <a-descriptions-item label="评分">⭐ {{ currentApp?.rating?.toFixed(1) }}</a-descriptions-item>
        <a-descriptions-item label="下载量">{{ currentApp?.downloads }} 次</a-descriptions-item>
        <a-descriptions-item label="大小">{{ currentApp?.size }}</a-descriptions-item>
        <a-descriptions-item label="更新日期">{{ currentApp?.updated_at }}</a-descriptions-item>
      </a-descriptions>
      <a-card title="应用介绍" style="margin-bottom: 16px">
        <div style="line-height: 1.7; color: #4e5969">{{ currentApp?.description }}</div>
      </a-card>
      <a-card title="版本更新">
        <a-timeline>
          <a-timeline-item v-for="(v, i) in currentApp?.changelog || []" :key="i" :color="i === 0 ? 'green' : 'gray'">
            <div style="font-weight: 600">{{ v.version }}</div>
            <div style="font-size: 12px; color: #86909c">{{ v.date }}</div>
            <div style="font-size: 13px; margin-top: 4px">{{ v.desc }}</div>
          </a-timeline-item>
        </a-timeline>
      </a-card>
    </a-drawer>

    <!-- 提交应用弹窗 -->
    <a-modal v-model:visible="submitVisible" title="提交应用到商店" :width="560" @before-ok="handleSubmit" @cancel="submitVisible = false">
      <a-form :model="submitForm" layout="vertical">
        <a-form-item label="应用名称" field="name">
          <a-input v-model="submitForm.name" placeholder="请输入应用名称" />
        </a-form-item>
        <a-form-item label="开发者" field="developer">
          <a-input v-model="submitForm.developer" placeholder="请输入开发者名称" />
        </a-form-item>
        <a-form-item label="版本号" field="version">
          <a-input v-model="submitForm.version" placeholder="如 v1.0.0" />
        </a-form-item>
        <a-form-item label="应用分类" field="category">
          <a-select v-model="submitForm.category" placeholder="请选择分类">
            <a-option v-for="cat in categories" :key="cat.key" :value="cat.key">{{ cat.icon }} {{ cat.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="价格（元）" field="price">
          <a-input-number v-model="submitForm.price" :min="0" :precision="2" placeholder="0表示免费" style="width: 100%" />
        </a-form-item>
        <a-form-item label="应用描述" field="description">
          <a-textarea v-model="submitForm.description" placeholder="请输入应用描述（100-500字）" :rows="4" />
        </a-form-item>
        <a-form-item label="安装包" field="package">
          <a-upload action="/" accept=".apk,.zip" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const detailVisible = ref(false)
const submitVisible = ref(false)
const currentApp = ref(null)

const stats = reactive({ total_apps: 0, total_installs: 0, weekly_new: 0, avg_rating: 0 })

const filterForm = reactive({ category: '', sort: 'downloads', price: '', keyword: '' })

const submitForm = reactive({ name: '', developer: '', version: '', category: '', price: 0, description: '', package: '' })

const categories = [
  { key: 'tool', icon: '🔧', name: '工具类' },
  { key: 'iot', icon: '📡', name: 'IoT 控制' },
  { key: 'monitor', icon: '📊', name: '监控管理' },
  { key: 'ai', icon: '🤖', name: 'AI 智能' },
  { key: 'life', icon: '🏠', name: '生活服务' },
  { key: 'education', icon: '📚', name: '教育培训' }
]

const apps = ref([])

const getCatColor = (cat) => ({ tool: 'orange', iot: 'blue', monitor: 'green', ai: 'purple', life: 'cyan', education: 'magenta' }[cat] || 'gray')

const formatNum = (n) => n >= 10000 ? `${(n / 10000).toFixed(1)}万` : n

const loadApps = async () => {
  loading.value = true
  await new Promise(r => setTimeout(r, 600))
  apps.value = [
    { app_id: 'APP-001', name: '智能家居控制', developer: '官方团队', icon: '🏠', banner_color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)', version: 'v2.3.0', category: 'iot', category_text: 'IoT 控制', price: 0, rating: 4.8, downloads: 158200, description: '一站式智能家居控制中心，支持设备联动、场景模式、定时任务等功能', is_official: true, size: '32.5MB', updated_at: '2026-03-15', changelog: [{ version: 'v2.3.0', date: '2026-03-15', desc: '新增设备分组功能，优化响应速度' }, { version: 'v2.2.0', date: '2026-02-20', desc: '支持语音控制' }] },
    { app_id: 'APP-002', name: '设备监控大师', developer: '物联科技', icon: '📊', banner_color: 'linear-gradient(135deg, #11998e 0%, #38ef7d 100%)', version: 'v1.8.5', category: 'monitor', category_text: '监控管理', price: 29.9, rating: 4.6, downloads: 87600, description: '实时监控设备状态，支持多设备数据看板、告警推送、历史数据查询', is_official: false, size: '28.1MB', updated_at: '2026-03-10', changelog: [{ version: 'v1.8.5', date: '2026-03-10', desc: '修复数据延迟问题' }] },
    { app_id: 'APP-003', name: 'AI 助手 Pro', developer: '智云科技', icon: '🤖', banner_color: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)', version: 'v3.1.0', category: 'ai', category_text: 'AI 智能', price: 59.0, rating: 4.9, downloads: 234000, description: '基于大模型的智能助手，支持设备问答、自动化编排、数据分析', is_official: false, size: '45.2MB', updated_at: '2026-03-18', changelog: [{ version: 'v3.1.0', date: '2026-03-18', desc: '支持多模态交互' }, { version: 'v3.0.0', date: '2026-02-01', desc: '全新大模型底座' }] },
    { app_id: 'APP-004', name: '网络测速工具', developer: 'SpeedLab', icon: '⚡', banner_color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)', version: 'v1.2.0', category: 'tool', category_text: '工具类', price: 0, rating: 4.3, downloads: 45200, description: '快速检测网络速度，支持 WiFi/5G 多网络测试，生成测试报告', is_official: false, size: '12.3MB', updated_at: '2026-01-25', changelog: [] },
    { app_id: 'APP-005', name: '家庭保险柜', developer: '安心保险', icon: '🛡️', banner_color: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)', version: 'v1.0.5', category: 'life', category_text: '生活服务', price: 19.9, rating: 4.5, downloads: 32100, description: '家庭财产保障应用，支持设备绑定、理赔报案、保险方案推荐', is_official: false, size: '18.7MB', updated_at: '2026-02-28', changelog: [] },
    { app_id: 'APP-006', name: '儿童学习乐园', developer: '教育局合作', icon: '📚', banner_color: 'linear-gradient(135deg, #a8edea 0%, #fed6e3 100%)', version: 'v2.0.1', category: 'education', category_text: '教育培训', price: 0, rating: 4.7, downloads: 67800, description: '寓教于乐的儿童学习应用，涵盖语数外、科普等精品课程', is_official: true, size: '56.0MB', updated_at: '2026-03-01', changelog: [{ version: 'v2.0.1', date: '2026-03-01', desc: '新增编程启蒙课程' }] }
  ]
  stats.total_apps = apps.value.length
  stats.total_installs = apps.value.reduce((s, a) => s + a.downloads, 0)
  stats.weekly_new = 12
  stats.avg_rating = (apps.value.reduce((s, a) => s + a.rating, 0) / apps.value.length).toFixed(1)
  loading.value = false
}

const viewAppDetail = (app) => {
  currentApp.value = app
  detailVisible.value = true
}

const handleInstall = (app) => {
  Message.success(`${app.name} 安装请求已提交`)
}

const showSubmitModal = () => {
  Object.assign(submitForm, { name: '', developer: '', version: '', category: '', price: 0, description: '' })
  submitVisible.value = true
}

const handleSubmit = async (done) => {
  if (!submitForm.name || !submitForm.developer || !submitForm.category) {
    Message.error('请填写必填字段')
    done(false)
    return
  }
  await new Promise(r => setTimeout(r, 800))
  Message.success('应用已提交，审核通过后将在商店展示')
  submitVisible.value = false
  done(true)
}

loadApps()
</script>

<style scoped>
.page-container { padding: 20px; }
.breadcrumb { margin-bottom: 12px; }
.stat-card { text-align: center; }
.section-card { margin-bottom: 16px; }
.cat-tag { cursor: pointer; padding: 4px 12px; font-size: 14px; }
.cat-tag.cat-active { background: #165dff; color: #fff; border-color: #165dff; }
.app-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 16px; }
.app-card { border: 1px solid #e5e6e8; border-radius: 8px; overflow: hidden; background: #fff; transition: box-shadow 0.2s; }
.app-card:hover { box-shadow: 0 4px 16px rgba(0,0,0,0.08); }
.app-banner { height: 80px; display: flex; align-items: center; padding: 0 16px; position: relative; }
.app-icon { font-size: 36px; }
.app-official-tag { position: absolute; top: 8px; right: 8px; }
.app-info { padding: 12px 16px; }
.app-name { font-weight: 600; font-size: 15px; color: #1d2129; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.app-dev { font-size: 12px; color: #86909c; margin-bottom: 6px; }
.app-category { display: flex; align-items: center; gap: 8px; margin-bottom: 6px; }
.app-price { font-size: 14px; font-weight: 600; color: #f53f3f; }
.app-price.free { color: #0fc6c2; }
.app-rating { display: flex; align-items: center; gap: 4px; margin-bottom: 8px; }
.rating-num { font-size: 12px; font-weight: 600; margin-left: 4px; }
.download-num { font-size: 11px; color: #86909c; margin-left: 8px; }
.app-desc { font-size: 12px; color: #4e5969; line-height: 1.5; overflow: hidden; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; margin-bottom: 12px; }
.app-actions { display: flex; gap: 8px; padding: 0 16px 12px; }
.detail-banner { padding: 32px 24px; display: flex; align-items: center; gap: 16px; }
.detail-icon { font-size: 48px; }
.detail-title { font-size: 22px; font-weight: 700; color: #fff; }
</style>
