<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>研究平台</a-breadcrumb-item>
      <a-breadcrumb-item>
        <a @click="goBack">数据集开放平台</a>
      </a-breadcrumb-item>
      <a-breadcrumb-item>数据集详情</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 返回按钮 -->
    <div class="mb-4">
      <a-button @click="goBack">
        <template #icon><icon-left /></template>
        返回列表
      </a-button>
    </div>

    <a-spin :spinning="loading">
      <a-row :gutter="16">
        <!-- 基本信息 -->
        <a-col :span="16">
          <a-card title="基本信息" class="mb-4">
            <a-descriptions :column="2" bordered>
              <a-descriptions-item label="数据集名称">{{ dataset.name }}</a-descriptions-item>
              <a-descriptions-item label="数据类型">
                <a-tag :color="getCategoryColor(dataset.category)">
                  {{ getCategoryText(dataset.category) }}
                </a-tag>
              </a-descriptions-item>
              <a-descriptions-item label="数据格式">{{ dataset.data_format?.toUpperCase() }}</a-descriptions-item>
              <a-descriptions-item label="数据大小">{{ formatSize(dataset.data_size) }}</a-descriptions-item>
              <a-descriptions-item label="样本数量">{{ dataset.record_count }}</a-descriptions-item>
              <a-descriptions-item label="访问权限">
                <a-tag :color="getAccessLevelColor(dataset.access_level)">
                  {{ getAccessLevelText(dataset.access_level) }}
                </a-tag>
              </a-descriptions-item>
              <a-descriptions-item label="许可证">{{ getLicenseText(dataset.license) }}</a-descriptions-item>
              <a-descriptions-item label="DOI">{{ dataset.doi }}</a-descriptions-item>
              <a-descriptions-item label="下载量">{{ dataset.download_count }}</a-descriptions-item>
              <a-descriptions-item label="引用数">{{ dataset.citation_count }}</a-descriptions-item>
              <a-descriptions-item label="发布时间" :span="2">
                {{ formatTime(dataset.published_at) }}
              </a-descriptions-item>
              <a-descriptions-item label="描述" :span="2">
                {{ dataset.description || '暂无描述' }}
              </a-descriptions-item>
            </a-descriptions>

            <div class="mt-4">
              <a-space>
                <a-button type="primary" @click="handleDownload">
                  <template #icon><icon-download /></template>
                  下载数据集
                </a-button>
                <a-button @click="handleCite">
                  <template #icon><icon-file /></template>
                  引用数据集
                </a-button>
                <a-button type="primary" status="warning" @click="handleEdit">
                  <template #icon><icon-edit /></template>
                  编辑
                </a-button>
                <a-button status="danger" @click="handleDelete">
                  <template #icon><icon-delete /></template>
                  删除
                </a-button>
              </a-space>
            </div>
          </a-card>
        </a-col>

        <!-- 统计信息 -->
        <a-col :span="8">
          <a-card title="统计信息" class="mb-4">
            <a-statistic title="下载次数" :value="dataset.download_count">
              <template #icon><icon-download /></template>
            </a-statistic>
            <a-divider />
            <a-statistic title="引用次数" :value="dataset.citation_count">
              <template #icon><icon-file /></template>
            </a-statistic>
            <a-divider />
            <a-statistic title="数据记录" :value="dataset.record_count">
              <template #icon><icon-database /></template>
            </a-statistic>
          </a-card>

          <!-- 标签 -->
          <a-card title="标签" class="mb-4">
            <a-space wrap>
              <a-tag v-for="tag in tags" :key="tag" color="blue">{{ tag }}</a-tag>
              <a-tag v-if="!tags.length" color="gray">暂无标签</a-tag>
            </a-space>
          </a-card>
        </a-col>
      </a-row>

      <!-- 版本历史 -->
      <a-card title="版本历史" class="mb-4">
        <a-table
          :columns="versionColumns"
          :data="versions"
          :loading="versionLoading"
          :pagination="false"
          row-key="id"
        >
          <template #version="{ record }">
            <a-tag color="blue">{{ record.version }}</a-tag>
          </template>
          <template #published_at="{ record }">
            {{ formatTime(record.published_at) }}
          </template>
          <template #actions="{ record }">
            <a-button type="text" size="small" @click="handleDownloadVersion(record)">下载</a-button>
          </template>
        </a-table>
        <div class="mt-4">
          <a-button type="primary" @click="showVersionModal = true">
            <template #icon><icon-plus /></template>
            创建新版本
          </a-button>
        </div>
      </a-card>
    </a-spin>

    <!-- 创建版本弹窗 -->
    <a-modal
      v-model:visible="showVersionModal"
      title="创建新版本"
      @ok="handleCreateVersion"
      @cancel="showVersionModal = false"
    >
      <a-form :model="versionForm" layout="vertical">
        <a-form-item label="版本号" required>
          <a-input v-model="versionForm.version" placeholder="例如: v1.1" />
        </a-form-item>
        <a-form-item label="变更说明">
          <a-textarea v-model="versionForm.changes" placeholder="请描述本次更新内容" :rows="3" />
        </a-form-item>
        <a-form-item label="文件URL">
          <a-input v-model="versionForm.file_url" placeholder="请输入文件URL" />
        </a-form-item>
        <a-form-item label="记录数">
          <a-input-number v-model="versionForm.record_count" :min="0" style="width: 100%" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 引用信息弹窗 -->
    <a-modal
      v-model:visible="showCiteModal"
      title="引用数据集"
      @ok="showCiteModal = false"
      :footer="false"
    >
      <a-alert type="info" class="mb-4">
        您可以引用此数据集用于您的研究论文或报告中。
      </a-alert>
      <a-descriptions :column="1" bordered>
        <a-descriptions-item label="DOI">{{ citeInfo.doi }}</a-descriptions-item>
        <a-descriptions-item label="数据集名称">{{ citeInfo.name }}</a-descriptions-item>
        <a-descriptions-item label="发布日期">{{ citeInfo.published }}</a-descriptions-item>
        <a-descriptions-item label="访问权限">{{ citeInfo.access_level }}</a-descriptions-item>
      </a-descriptions>
      <a-divider />
      <a-form-item label="引用格式">
        <a-textarea :model-value="citationText" :rows="4" readonly />
      </a-form-item>
      <a-button type="primary" @click="copyCitation">
        <template #icon><icon-copy /></template>
        复制引用
      </a-button>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Message } from '@arco-design/web-vue'

const router = useRouter()
const route = useRoute()

const loading = ref(false)
const dataset = ref({})
const versions = ref([])
const versionLoading = ref(false)
const showVersionModal = ref(false)
const showCiteModal = ref(false)

const versionForm = reactive({
  version: '',
  changes: '',
  file_url: '',
  record_count: 0,
})

const citeInfo = ref({})

const versionColumns = [
  { title: '版本', dataIndex: 'version', width: 100 },
  { title: '变更说明', dataIndex: 'changes', ellipsis: true },
  { title: '记录数', dataIndex: 'record_count', width: 100 },
  { title: '发布时间', dataIndex: 'published_at', width: 150 },
  { title: '操作', dataIndex: 'actions', width: 100 },
]

const tags = computed(() => {
  if (!dataset.value.tags) return []
  try {
    return JSON.parse(dataset.value.tags)
  } catch {
    return []
  }
})

const citationText = computed(() => {
  return `${citeInfo.value.name} (${citeInfo.value.published}). DOI: ${citeInfo.value.doi}`
})

onMounted(() => {
  loadDataset()
  loadVersions()
})

const loadDataset = async () => {
  loading.value = true
  try {
    const res = await fetch(`/api/research/datasets/${route.params.id}`)
    const data = await res.json()
    if (data.status === 'success') {
      dataset.value = data.data
    }
  } catch (error) {
    Message.error('加载数据集失败')
  } finally {
    loading.value = false
  }
}

const loadVersions = async () => {
  versionLoading.value = true
  try {
    const res = await fetch(`/api/research/datasets/${route.params.id}/versions`)
    const data = await res.json()
    if (data.status === 'success') {
      versions.value = data.data
    }
  } catch (error) {
    Message.error('加载版本失败')
  } finally {
    versionLoading.value = false
  }
}

const goBack = () => {
  router.push('/research/datasets')
}

const handleDownload = async () => {
  try {
    const res = await fetch(`/api/research/datasets/${route.params.id}/download`, { method: 'POST' })
    const data = await res.json()
    if (data.status === 'success' && data.download_url) {
      window.open(data.download_url, '_blank')
      Message.success('开始下载')
    }
  } catch (error) {
    Message.error('下载失败')
  }
}

const handleCite = async () => {
  try {
    const res = await fetch(`/api/research/datasets/${route.params.id}/cite`, { method: 'POST' })
    const data = await res.json()
    if (data.status === 'success') {
      citeInfo.value = data.citation
      showCiteModal.value = true
    }
  } catch (error) {
    Message.error('获取引用信息失败')
  }
}

const copyCitation = () => {
  navigator.clipboard.writeText(citationText.value)
  Message.success('已复制到剪贴板')
}

const handleCreateVersion = async () => {
  try {
    const res = await fetch(`/api/research/datasets/${route.params.id}/versions`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(versionForm),
    })
    const data = await res.json()
    if (data.status === 'success') {
      Message.success('创建成功')
      showVersionModal.value = false
      loadVersions()
    }
  } catch (error) {
    Message.error('创建失败')
  }
}

const handleDownloadVersion = (version) => {
  if (version.file_url) {
    window.open(version.file_url, '_blank')
  }
}

const handleEdit = () => {
  // 跳转到编辑页面或打开编辑弹窗
  Message.info('编辑功能开发中')
}

const handleDelete = async () => {
  try {
    await new Promise((resolve) => {
      // 确认对话框
      resolve()
    })
    const res = await fetch(`/api/research/datasets/${route.params.id}`, { method: 'DELETE' })
    const data = await res.json()
    if (data.status === 'success') {
      Message.success('删除成功')
      router.push('/research/datasets')
    }
  } catch (error) {
    Message.error('删除失败')
  }
}

const getCategoryColor = (category) => {
  const colors = { behavior: 'blue', emotion: 'purple', health: 'green', vocalization: 'orange' }
  return colors[category] || 'gray'
}

const getCategoryText = (category) => {
  const texts = { behavior: '行为', emotion: '情感', health: '健康', vocalization: '声音' }
  return texts[category] || category
}

const getAccessLevelColor = (level) => {
  const colors = { public: 'green', restricted: 'orange', private: 'red' }
  return colors[level] || 'gray'
}

const getAccessLevelText = (level) => {
  const texts = { public: '公开', restricted: '受限', private: '私有' }
  return texts[level] || level
}

const getLicenseText = (license) => {
  const texts = { MIT: 'MIT', GPL: 'GPL', proprietary: '专有' }
  return texts[license] || license
}

const formatSize = (bytes) => {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  let i = 0
  while (bytes >= 1024 && i < units.length - 1) {
    bytes /= 1024
    i++
  }
  return `${bytes.toFixed(1)} ${units[i]}`
}

const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.pro-page-container {
  padding: 16px;
}
.pro-breadcrumb {
  margin-bottom: 16px;
}
.mb-4 {
  margin-bottom: 16px;
}
.mt-4 {
  margin-top: 16px;
}
</style>
