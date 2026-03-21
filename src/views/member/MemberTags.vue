<template>
  <div class="member-tags-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员标签</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 操作栏 -->
    <a-card class="action-card">
      <a-space>
        <a-button type="primary" @click="showCreate">新增标签</a-button>
        <a-button @click="loadTags">刷新</a-button>
      </a-space>
      <template #extra>
        <a-space>
          <span style="color: #666; font-size: 13px;">共 {{ tagList.length }} 个标签</span>
        </a-space>
      </template>
    </a-card>

    <!-- 标签分类导航 -->
    <a-card class="nav-card" style="margin-bottom: 16px;">
      <a-space wrap>
        <a-link @click="activeCategory = 'all'" :type="activeCategory === 'all' ? 'primary' : 'secondary'">全部</a-link>
        <a-link @click="activeCategory = 'buy'" :type="activeCategory === 'buy' ? 'primary' : 'secondary'">购买行为</a-link>
        <a-link @click="activeCategory = 'interest'" :type="activeCategory === 'interest' ? 'primary' : 'secondary'">兴趣偏好</a-link>
        <a-link @click="activeCategory = 'level'" :type="activeCategory === 'level' ? 'primary' : 'secondary'">等级相关</a-link>
        <a-link @click="activeCategory = 'custom'" :type="activeCategory === 'custom' ? 'primary' : 'secondary'">自定义</a-link>
        <a-divider direction="vertical" />
        <router-link to="/member/tags/high-freq">
          <a-button type="text" size="small">高频购买标签</a-button>
        </router-link>
        <router-link to="/member/tags/low-freq">
          <a-button type="text" size="small">低频购买标签</a-button>
        </router-link>
        <router-link to="/member/tags/interest">
          <a-button type="text" size="small">兴趣分类标签</a-button>
        </router-link>
        <router-link to="/member/tags/auto-clean">
          <a-button type="text" size="small">自动清除设置</a-button>
        </router-link>
        <router-link to="/member/tags/report">
          <a-button type="text" size="small">标签报表</a-button>
        </router-link>
      </a-space>
    </a-card>

    <!-- 标签列表 -->
    <a-card>
      <a-table :columns="columns" :data="filteredTags" :loading="loading" row-key="id" :pagination="{ pageSize: 10 }">
        <template #tagName="{ record }">
          <a-tag :color="record.color">{{ record.name }}</a-tag>
        </template>
        <template #category="{ record }">
          <a-tag :color="getCategoryColor(record.category)">{{ getCategoryName(record.category) }}</a-tag>
        </template>
        <template #memberCount="{ record }">
          <a-statistic :value="record.memberCount || 0" :value-style="{ fontSize: '14px' }" />
        </template>
        <template #autoClean="{ record }">
          <a-tag :color="record.autoClean ? 'green' : 'gray'">{{ record.autoClean ? '自动清除' : '手动' }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 新增/编辑标签弹窗 -->
    <a-modal
      v-model:visible="formVisible"
      :title="isEdit ? '编辑标签' : '新增标签'"
      @before-ok="handleFormSubmit"
      @cancel="formVisible = false"
      :width="480"
      :loading="formLoading"
    >
      <a-form :model="form" layout="vertical" ref="formRef">
        <a-form-item label="标签名称" field="name" :rules="[{ required: true, message: '请输入标签名称' }]">
          <a-input v-model="form.name" placeholder="如：高频购买用户" />
        </a-form-item>
        <a-form-item label="标签编码" field="code" :rules="[{ required: true, message: '请输入标签编码' }]">
          <a-input v-model="form.code" placeholder="如：HIGH_FREQ_BUYER" :disabled="isEdit" />
        </a-form-item>
        <a-form-item label="标签分类" field="category" :rules="[{ required: true, message: '请选择标签分类' }]">
          <a-select v-model="form.category" placeholder="选择分类">
            <a-option value="buy">购买行为</a-option>
            <a-option value="interest">兴趣偏好</a-option>
            <a-option value="level">等级相关</a-option>
            <a-option value="custom">自定义</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="标签颜色">
          <a-space>
            <a-color-picker v-model="form.color" :show-alpha="false" :presets="colorPresets" />
            <a-tag :color="form.color">预览</a-tag>
          </a-space>
        </a-form-item>
        <a-form-item label="自动清除">
          <a-switch v-model="form.autoClean" />
        </a-form-item>
        <a-form-item label="标签描述">
          <a-textarea v-model="form.description" :rows="2" placeholder="简要描述该标签" />
        </a-form-item>
        <a-form-item label="排序值">
          <a-input-number v-model="form.sort" :min="0" style="width: 100%" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import * as api from '@/api/marketing'

const tagList = ref([])
const loading = ref(false)
const formLoading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const activeCategory = ref('all')

const colorPresets = [
  '#165DFF', '#00D084', '#FF6B6B', '#FFC107',
  '#9C27B0', '#FF9800', '#E91E63', '#3F51B5'
]

const form = reactive({
  name: '',
  code: '',
  category: 'custom',
  color: '#165DFF',
  autoClean: false,
  description: '',
  sort: 0
})

const columns = [
  { title: '标签名称', slotName: 'tagName', width: 180 },
  { title: '编码', dataIndex: 'code', width: 160 },
  { title: '分类', slotName: 'category', width: 120 },
  { title: '会员数', slotName: 'memberCount', width: 100 },
  { title: '清除方式', slotName: 'autoClean', width: 100 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '操作', slotName: 'actions', width: 140 }
]

const filteredTags = computed(() => {
  if (activeCategory.value === 'all') return tagList.value
  return tagList.value.filter(t => t.category === activeCategory.value)
})

const getCategoryColor = (cat) => {
  const map = { buy: 'blue', interest: 'purple', level: 'orange', custom: 'gray' }
  return map[cat] || 'gray'
}

const getCategoryName = (cat) => {
  const map = { buy: '购买行为', interest: '兴趣偏好', level: '等级相关', custom: '自定义' }
  return map[cat] || cat
}

const loadTags = async () => {
  loading.value = true
  try {
    const res = await api.getTagList()
    tagList.value = res.data || []
  } catch (err) {
    Message.error('加载标签列表失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

const showCreate = () => {
  isEdit.value = false
  Object.assign(form, { name: '', code: '', category: 'custom', color: '#165DFF', autoClean: false, description: '', sort: 0 })
  formVisible.value = true
}

const showEdit = (record) => {
  isEdit.value = true
  currentId.value = record.id
  Object.assign(form, { ...record })
  formVisible.value = true
}

const handleFormSubmit = async (done) => {
  if (!form.name || !form.code) {
    Message.warning('请填写标签名称和编码')
    done(false)
    return
  }
  formLoading.value = true
  try {
    if (isEdit.value) {
      await api.updateTag(currentId.value, { ...form })
      Message.success('更新成功')
    } else {
      await api.createTag({ ...form })
      Message.success('创建成功')
    }
    formVisible.value = false
    loadTags()
    done(true)
  } catch (err) {
    Message.error(err.message || '操作失败')
    done(false)
  } finally {
    formLoading.value = false
  }
}

const handleDelete = (record) => {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除标签「${record.name}」吗？`,
    okText: '确认删除',
    onOk: async () => {
      try {
        await api.deleteTag(record.id)
        Message.success('删除成功')
        loadTags()
      } catch (err) {
        Message.error(err.message || '删除失败')
      }
    }
  })
}

onMounted(() => loadTags())
</script>

<style scoped>
.member-tags-page {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
.nav-card :deep(.arco-card-body) { padding: 12px 16px; }
</style>
