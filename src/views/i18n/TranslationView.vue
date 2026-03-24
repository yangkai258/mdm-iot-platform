<template>
  <div class="page-container">
    <!-- 顶部：标题 + 操作按钮 -->
    <a-card class="header-card">
      <div class="header-row">
        <div class="page-title">翻译管理</div>
        <a-space :size="12">
          <a-button @click="handleExport">「导出」</a-button>
          <a-button @click="showImportModal = true">「导入」</a-button>
          <a-button type="primary" @click="openAddModal">「添加翻译」</a-button>
        </a-space>
      </div>
    </a-card>

    <!-- 筛选栏 -->
    <a-card class="filter-card">
      <div class="filter-row">
        <!-- 语言切换 -->
        <a-select
          v-model="filters.locale"
          placeholder="选择语言"
          style="width: 160px"
          allow-clear
          @change="handleFilterChange"
        >
          <a-option v-for="lang in languages" :key="lang.code" :value="lang.code">
            {{ lang.native_name }} ({{ lang.code }})
          </a-option>
        </a-select>
        <a-input
          v-model="filters.key"
          placeholder="搜索 Key..."
          style="width: 200px"
          allow-clear
          @change="handleFilterChange"
          @press-enter="handleFilterChange"
        />
        <a-input
          v-model="filters.search"
          placeholder="搜索译文..."
          style="width: 200px"
          allow-clear
          @change="handleFilterChange"
          @press-enter="handleFilterChange"
        />
        <a-button @click="handleFilterChange">「查询」</a-button>
        <a-button @click="resetFilters">「重置」</a-button>
      </div>
    </a-card>

    <!-- 翻译列表 -->
    <a-card class="table-card">
      <a-table
        :columns="columns"
        :data="filteredData"
        :loading="loading"
        :pagination="{ pageSize: 20, showTotal: true, showPageSize: true }"
        :scroll="{ x: 900 }"
        row-key="id"
        @page-change="handlePageChange"
        @page-size-change="handlePageSizeChange"
      >
        <template #key="{ record }">
          <span class="translation-key">{{ record.key }}</span>
        </template>
        <template #locale="{ record }">
          <a-tag :color="localeColor(record.locale)">{{ localeLabel(record.locale) }}</a-tag>
        </template>
        <template #value="{ record }">
          <span class="translation-value" :title="record.value">{{ record.value }}</span>
        </template>
        <template #tags="{ record }">
          <a-tag v-for="tag in (record.tags || [])" :key="tag" size="small">{{ tag }}</a-tag>
        </template>
        <template #updatedAt="{ record }">
          <span class="date-text">{{ formatDate(record.updated_at) }}</span>
        </template>
        <template #actions="{ record }">
          <div class="action-btns">
            <a-button type="text" size="small" @click="openEditModal(record)">「编辑」</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">「删除」</a-button>
          </div>
        </template>
      </a-table>
    </a-card>

    <!-- 添加/编辑翻译弹窗 -->
    <a-modal
      v-model:visible="formVisible"
      :title="editingItem ? '编辑翻译' : '添加翻译'"
      :width="560"
      @before-ok="handleSubmit"
      @cancel="formVisible = false"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="翻译 Key" required>
          <a-input
            v-model="form.key"
            placeholder="如：dashboard.title"
            :disabled="!!editingItem"
          />
        </a-form-item>
        <a-form-item label="语言" required>
          <a-select v-model="form.locale" placeholder="选择语言" :disabled="!!editingItem">
            <a-option v-for="lang in languages" :key="lang.code" :value="lang.code">
              {{ lang.native_name }} ({{ lang.code }})
            </a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="译文" required>
          <a-textarea
            v-model="form.value"
            placeholder="请输入翻译内容..."
            :rows="3"
          />
        </a-form-item>
        <a-form-item label="描述">
          <a-input v-model="form.description" placeholder="可选：描述该翻译的用途" />
        </a-form-item>
        <a-form-item label="标签">
          <a-select
            v-model="form.tags"
            multiple
            placeholder="选择或输入标签"
            allow-create
            :style="{ width: '100%' }"
          >
            <a-option v-for="tag in commonTags" :key="tag" :value="tag">{{ tag }}</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 导入弹窗 -->
    <a-modal
      v-model:visible="showImportModal"
      title="导入翻译"
      :width="480"
      @before-ok="handleImport"
      @cancel="showImportModal = false"
    >
      <a-form :model="importForm" layout="vertical">
        <a-form-item label="目标语言" required>
          <a-select v-model="importForm.locale" placeholder="选择语言">
            <a-option v-for="lang in languages" :key="lang.code" :value="lang.code">
              {{ lang.native_name }} ({{ lang.code }})
            </a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="JSON 数据" required>
          <a-textarea
            v-model="importForm.jsonData"
            placeholder='[{"key":"hello","value":"你好"},{"key":"world","value":"世界"}]'
            :rows="8"
            style="font-family: monospace; font-size: 12px"
          />
        </a-form-item>
        <div class="import-tip">支持 JSON 数组格式，每项包含 key 和 value 字段</div>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import {
  getTranslations,
  createTranslation,
  updateTranslation,
  deleteTranslation,
  importTranslations,
  getLanguages
} from '@/api/i18n'
import dayjs from 'dayjs'

const loading = ref(false)
const translations = ref([])
const languages = ref([])

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const filters = reactive({
  locale: '',
  key: '',
  search: ''
})

const editingItem = ref(null)
const formVisible = ref(false)
const showImportModal = ref(false)

const form = reactive({
  key: '',
  locale: '',
  value: '',
  description: '',
  tags: [] as string[]
})

const importForm = reactive({
  locale: '',
  jsonData: ''
})

const commonTags = ['common', 'menu', 'button', 'message', 'validation', 'placeholder']

const columns = [
  { title: 'Key', slotName: 'key', width: 200 },
  { title: '语言', slotName: 'locale', width: 120 },
  { title: '译文', slotName: 'value', width: 280 },
  { title: '标签', slotName: 'tags', width: 150 },
  { title: '更新时间', slotName: 'updatedAt', width: 160 },
  { title: '操作', slotName: 'actions', width: 140, fixed: 'right' }
]

const filteredData = computed(() => {
  let data = translations.value

  if (filters.locale) {
    data = data.filter(t => t.locale === filters.locale)
  }
  if (filters.key) {
    data = data.filter(t => t.key.toLowerCase().includes(filters.key.toLowerCase()))
  }
  if (filters.search) {
    data = data.filter(t => t.value.toLowerCase().includes(filters.search.toLowerCase()))
  }

  return data
})

function localeColor(locale) {
  const map: Record<string, string> = {
    'zh-CN': 'arcoblue',
    'en-US': 'orangered',
    'ja-JP': 'pink',
    'ko-KR': 'purple',
    'fr-FR': 'geekblue',
    'de-DE': 'gold',
    'es-ES': 'sunrise',
    'ru-RU': 'volcano'
  }
  return map[locale] || 'gray'
}

function localeLabel(locale) {
  const map: Record<string, string> = {
    'zh-CN': '简体中文',
    'en-US': 'English',
    'ja-JP': '日本語',
    'ko-KR': '한국어',
    'fr-FR': 'Français',
    'de-DE': 'Deutsch',
    'es-ES': 'Español',
    'ru-RU': 'Русский'
  }
  return map[locale] || locale
}

function formatDate(date) {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

async function loadData() {
  loading.value = true
  try {
    const [transRes, langRes] = await Promise.all([
      getTranslations({ page: pagination.page, page_size: pagination.pageSize }),
      getLanguages()
    ])
    translations.value = transRes.data || transRes || []
    pagination.total = transRes.total || translations.value.length
    languages.value = langRes.data || langRes || []
  } catch (e) {
    console.error('加载翻译列表失败', e)
    // 降级：使用默认语言
    languages.value = [
      { code: 'zh-CN', name: 'Chinese (Simplified)', native_name: '简体中文', is_default: true, is_active: true, translation_count: 0 },
      { code: 'en-US', name: 'English (US)', native_name: 'English', is_default: false, is_active: true, translation_count: 0 }
    ]
  } finally {
    loading.value = false
  }
}

function handleFilterChange() {
  pagination.page = 1
  // 本地筛选直接通过 computed filteredData，API筛选则重新请求
  // 这里采用本地筛选，computed filteredData 已处理
}

function resetFilters() {
  filters.locale = ''
  filters.key = ''
  filters.search = ''
  pagination.page = 1
}

function handlePageChange(page) {
  pagination.page = page
  loadData()
}

function handlePageSizeChange(size) {
  pagination.pageSize = size
  pagination.page = 1
  loadData()
}

function openAddModal() {
  editingItem.value = null
  Object.assign(form, { key: '', locale: '', value: '', description: '', tags: [] })
  formVisible.value = true
}

function openEditModal(record) {
  editingItem.value = record
  Object.assign(form, {
    key: record.key,
    locale: record.locale,
    value: record.value,
    description: record.description || '',
    tags: record.tags ? [...record.tags] : []
  })
  formVisible.value = true
}

async function handleSubmit(done) {
  try {
    if (!form.key || !form.locale || !form.value) {
      Message.error('请填写完整信息')
      done(false)
      return
    }
    if (editingItem.value) {
      await updateTranslation(editingItem.value.id, {
        value: form.value,
        description: form.description,
        tags: form.tags
      })
      Message.success('更新成功')
    } else {
      await createTranslation({
        key: form.key,
        locale: form.locale,
        value: form.value,
        description: form.description,
        tags: form.tags
      })
      Message.success('创建成功')
    }
    formVisible.value = false
    await loadData()
  } catch (e) {
    Message.error('操作失败')
    done(false)
  }
}

function handleDelete(record) {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除翻译「${record.key}」(${localeLabel(record.locale)}) 吗？`,
    okText: '删除',
    onOk: async () => {
      try {
        await deleteTranslation(record.id)
        Message.success('删除成功')
        await loadData()
      } catch (e) {
        Message.error('删除失败')
      }
    }
  })
}

async function handleExport() {
  try {
    const data = await import('@/api/i18n').then(m => m.exportTranslations({ locale: filters.locale || undefined }))
    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `translations-${filters.locale || 'all'}.json`
    a.click()
    URL.revokeObjectURL(url)
    Message.success('导出成功')
  } catch (e) {
    Message.error('导出失败')
  }
}

async function handleImport(done) {
  try {
    if (!importForm.locale || !importForm.jsonData) {
      Message.error('请填写完整信息')
      done(false)
      return
    }
    let translations
    try {
      translations = JSON.parse(importForm.jsonData)
    } catch {
      Message.error('JSON 格式错误')
      done(false)
      return
    }
    await importTranslations({ locale: importForm.locale, translations })
    Message.success('导入成功')
    showImportModal.value = false
    importForm.locale = ''
    importForm.jsonData = ''
    await loadData()
  } catch (e) {
    Message.error('导入失败')
    done(false)
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.page-container {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  height: 100%;
  box-sizing: border-box;
}

.header-card { flex-shrink: 0; }

.header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.page-title {
  font-size: 16px;
  font-weight: 600;
}

.filter-card { flex-shrink: 0; }

.filter-row {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.table-card {
  flex: 1;
  overflow: auto;
}

.translation-key {
  font-family: monospace;
  font-size: 12px;
  color: var(--color-text-2);
  word-break: break-all;
}

.translation-value {
  display: block;
  max-width: 260px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.date-text {
  font-size: 12px;
  color: var(--color-text-3);
}

.action-btns {
  display: flex;
  gap: 4px;
}

.import-tip {
  font-size: 12px;
  color: var(--color-text-3);
  margin-top: 4px;
}
</style>
