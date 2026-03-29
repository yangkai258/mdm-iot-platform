<template>
  <div class="interest-tags-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>
        <router-link to="/member/tags">会员标签</router-link>
      </a-breadcrumb-item>
      <a-breadcrumb-item>兴趣分类标签</a-breadcrumb-item>
    </a-breadcrumb>

    <a-alert style="margin-bottom: 16px;">
      <template #title>兴趣分类标签说明</template>
      根据会员的购买品类偏好自动打上兴趣标签，如：数码爱好者、美妆达人等。支持手动维护和自动生成两种模式。
    </a-alert>

    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card hoverable>
          <a-statistic title="已配置品类" :value="categoryList.length" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card hoverable>
          <a-statistic title="已打标签会员" :value="stats.taggedCount || 0" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card hoverable>
          <a-statistic title="规则状态">
            <template #suffix>
              <a-tag :color="ruleEnabled ? 'green' : 'gray'">{{ ruleEnabled ? '已启用' : '已禁用' }}</a-tag>
            </template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <a-card class="action-card">
      <a-space>
        <a-button type="primary" @click="showCreate">新增品类</a-button>
        <a-switch v-model="ruleEnabled" @change="handleToggleEnable" />
        <span>启用自动打标签</span>
        <a-divider direction="vertical" />
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </a-card>

    <a-card>
      <a-table :columns="columns" :data="categoryList" :loading="loading" row-key="id" :pagination="false">
        <template #categoryName="{ record }">
          <a-tag :color="record.tagColor || 'purple'">{{ record.categoryName }}</a-tag>
        </template>
      </a-table>
        <template #productTypes="{ record }">
          <a-space wrap>
            <a-tag v-for="(t, i) in (record.productTypes || [])" :key="i" size="small">{{ t }}</a-tag>
          </a-space>
        </template>
        <template #threshold="{ record }">
          购买 ≥ {{ record.threshold || 1 }} 次
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 'active' ? 'green' : 'gray'">
            {{ record.status === 'active' ? '生效中' : '已禁用' }}
          </a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal
      v-model:visible="formVisible"
      :title="isEdit ? '编辑品类' : '新增品类'"
      @before-ok="handleFormSubmit"
      @cancel="formVisible = false"
      :width="560"
      :loading="formLoading"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="品类名称" :rules="[{ required: true, message: '请输入品类名称' }]">
          <a-input v-model="form.categoryName" placeholder="如：数码爱好者" />
        </a-form-item>
        <a-form-item label="关联商品类型" :rules="[{ required: true, message: '请选择关联商品类型' }]">
          <a-select v-model="form.productTypes" multiple placeholder="选择商品类型（可多选）" allow-create>
            <a-option value="手机数码">手机数码</a-option>
            <a-option value="电脑办公">电脑办公</a-option>
            <a-option value="家电">家电</a-option>
            <a-option value="服装">服装</a-option>
            <a-option value="美妆护肤">美妆护肤</a-option>
            <a-option value="食品饮料">食品饮料</a-option>
            <a-option value="母婴用品">母婴用品</a-option>
            <a-option value="图书">图书</a-option>
          </a-select>
        </a-form-item>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="触发阈值（购买次数）">
              <a-input-number v-model="form.threshold" :min="1" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="统计周期（天）">
              <a-input-number v-model="form.periodDays" :min="1" style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="关联标签">
          <a-select v-model="form.tagId" placeholder="选择已创建的标签" allow-create>
            <a-option v-for="tag in availableTags" :key="tag.id" :value="tag.id">{{ tag.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="标签颜色">
          <a-color-picker v-model="form.tagColor" :show-alpha="false" :presets="colorPresets" />
        </a-form-item>
        <a-form-item label="状态">
          <a-switch v-model="form.status" checked-value="active" unchecked-value="inactive" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import * as api from '@/api/marketing'

const categoryList = ref([])
const availableTags = ref([])
const stats = ref({})
const loading = ref(false)
const formLoading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const ruleEnabled = ref(false)

const colorPresets = ['#165DFF', '#00D084', '#FF6B6B', '#FFC107', '#9C27B0', '#FF9800', '#E91E63', '#3F51B5']

const form = reactive({
  categoryName: '',
  productTypes: [],
  threshold: 1,
  periodDays: 90,
  tagId: null,
  tagName: '',
  tagColor: '#9C27B0',
  status: 'active'
})

const columns = [
  { title: '品类名称', slotName: 'categoryName', width: 150 },
  { title: '关联商品类型', slotName: 'productTypes' },
  { title: '触发阈值', slotName: 'threshold', width: 140 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 140 }
]

const loadData = async () => {
  loading.value = true
  try {
    const [listRes, tagsRes] = await Promise.all([
      api.getInterestTagList(),
      api.getTagList({ category: 'interest' })
    ])
    categoryList.value = listRes.data || []
    availableTags.value = tagsRes.data || []
    ruleEnabled.value = listRes.enabled || false
    stats.value = listRes.stats || {}
  } catch (err) {
    Message.error('加载数据失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

const handleToggleEnable = async (val) => {
  try {
    await api.saveInterestTag({ enabled: val })
    Message.success(val ? '已启用' : '已禁用')
  } catch (err) {
    Message.error(err.message)
    ruleEnabled.value = !val
  }
}

const showCreate = () => {
  isEdit.value = false
  Object.assign(form, { categoryName: '', productTypes: [], threshold: 1, periodDays: 90, tagId: null, tagName: '', tagColor: '#9C27B0', status: 'active' })
  formVisible.value = true
}

const showEdit = (record) => {
  isEdit.value = true
  currentId.value = record.id
  Object.assign(form, { ...record })
  formVisible.value = true
}

const handleFormSubmit = async (done) => {
  formLoading.value = true
  try {
    if (isEdit.value) {
      await api.saveInterestTag({ id: currentId.value, ...form })
    } else {
      await api.saveInterestTag({ ...form })
    }
    Message.success('保存成功')
    formVisible.value = false
    loadData()
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
    content: `确定要删除品类「${record.categoryName}」吗？`,
    okText: '确认删除',
    onOk: async () => {
      try {
        await api.saveInterestTag({ id: record.id, _delete: true })
        Message.success('删除成功')
        loadData()
      } catch (err) {
        Message.error(err.message || '删除失败')
      }
    }
  })
}

onMounted(() => loadData())
</script>

<style scoped>
.interest-tags-page {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
