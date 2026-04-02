<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="high-freq-tags-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>
        <router-link to="/member/tags">会员标签</router-link>
      </a-breadcrumb-item>
      <a-breadcrumb-item>高频购买标签</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 提示说明 -->
    <a-alert style="margin-bottom: 16px;">
      <template #title>高频购买标签规则说明</template>
      当会员在指定时间周期内的购买次数达到阈值时，系统自动为其打上对应标签。支持多档位配置。
    </a-alert>

    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card hoverable>
          <a-statistic title="已配置档位" :value="ruleList.length" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card hoverable>
          <a-statistic title="覆盖会员数" :value="stats.memberCount || 0" />
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

    <!-- 操作栏 -->
    <a-card class="action-card">
      <a-space>
        <a-button type="primary" @click="showCreate">新增档位</a-button>
        <a-switch v-model="ruleEnabled" @change="handleToggleEnable" />
        <span>启用高频购买标签</span>
        <a-divider direction="vertical" />
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </a-card>

    <!-- 档位列表 -->
    <a-card>
      <a-table :columns="columns" :data="ruleList" :loading="loading" row-key="id" :pagination="false">
        <template #level="{ record }">
          <a-tag color="blue">{{ record.levelName || '档位' + record.level }}</a-tag>
        </template>
      </a-table>
        <template #buyCount="{ record }">
          <span style="color: #ff6b00; font-weight: 600;">≥ {{ record.buyCount }} 次</span>
        </template>
        <template #periodDays="{ record }">
          {{ record.periodDays }} 天
        </template>
        <template #tagName="{ record }">
          <a-tag :color="record.tagColor || 'blue'">{{ record.tagName || '未设置' }}</a-tag>
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

    <!-- 新增/编辑档位弹窗 -->
    <a-modal
      v-model:visible="formVisible"
      :title="isEdit ? '编辑档位' : '新增档位'"
      @before-ok="handleFormSubmit"
      @cancel="formVisible = false"
      :width="520"
      :loading="formLoading"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="档位名称" :rules="[{ required: true, message: '请输入档位名称' }]">
          <a-input v-model="form.levelName" placeholder="如：高频购买（≥10次）" />
        </a-form-item>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="购买次数阈值" :rules="[{ required: true, message: '请输入购买次数' }]">
              <a-input-number v-model="form.buyCount" :min="1" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="统计周期（天）" :rules="[{ required: true, message: '请输入统计周期' }]">
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

const ruleList = ref([])
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
  levelName: '',
  buyCount: 5,
  periodDays: 30,
  tagId: null,
  tagName: '',
  tagColor: '#165DFF',
  status: 'active'
})

const columns = [
  { title: '档位', slotName: 'level', width: 150 },
  { title: '购买次数', slotName: 'buyCount', width: 120 },
  { title: '统计周期', slotName: 'periodDays', width: 120 },
  { title: '关联标签', slotName: 'tagName', width: 160 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 140 }
]

const loadData = async () => {
  loading.value = true
  try {
    const [rulesRes, tagsRes] = await Promise.all([
      api.getHighFreqTagList(),
      api.getTagList({ category: 'buy' })
    ])
    ruleList.value = rulesRes.data || []
    availableTags.value = tagsRes.data || []
    ruleEnabled.value = rulesRes.enabled || false
    stats.value = rulesRes.stats || {}
  } catch (err) {
    Message.error('加载数据失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

const handleToggleEnable = async (val) => {
  try {
    await api.saveHighFreqTag({ enabled: val })
    Message.success(val ? '已启用' : '已禁用')
  } catch (err) {
    Message.error(err.message)
    ruleEnabled.value = !val
  }
}

const showCreate = () => {
  isEdit.value = false
  Object.assign(form, { levelName: '', buyCount: 5, periodDays: 30, tagId: null, tagName: '', tagColor: '#165DFF', status: 'active' })
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
      await api.saveHighFreqTag({ id: currentId.value, ...form })
      Message.success('更新成功')
    } else {
      await api.saveHighFreqTag({ ...form })
      Message.success('创建成功')
    }
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
    content: `确定要删除档位「${record.levelName}」吗？`,
    okText: '确认删除',
    onOk: async () => {
      try {
        await api.saveHighFreqTag({ id: record.id, _delete: true })
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
.high-freq-tags-page {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
