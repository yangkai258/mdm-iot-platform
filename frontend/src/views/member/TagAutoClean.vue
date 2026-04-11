<template>
  <div class="tag-auto-clean-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>
        <router-link to="/member/tags">会员标签</router-link>
      </a-breadcrumb-item>
      <a-breadcrumb-item>标签自动清除设置</a-breadcrumb-item>
    </a-breadcrumb>

    <a-alert style="margin-bottom: 16px;">
      <template #title>自动清除说明</template>
      配置标签的自动清除规则，当会员在指定时间内未满足标签条件时，系统自动移除该标签。支持按时间过期和条件失效两种模式。
    </a-alert>

    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card hoverable>
          <a-statistic title="已配置规则" :value="ruleList.length" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card hoverable>
          <a-statistic title="本周清除数" :value="stats.weeklyCleaned || 0" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card hoverable>
          <a-statistic title="自动清除总次数" :value="stats.totalCleaned || 0" />
        </a-card>
      </a-col>
    </a-row>

    <a-card class="action-card">
      <a-space>
        <a-button type="primary" @click="showCreate">新增规则</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </a-card>

    <a-card>
      <a-table :columns="columns" :data="ruleList" :loading="loading" row-key="id" :pagination="false">
        <template #tagName="{ record }">
          <a-tag :color="record.tagColor || 'arcoblue'">{{ record.tagName || '标签ID:' + record.tagId }}</a-tag>
        </template>
      </a-table>
        <template #cleanMode="{ record }">
          <a-tag :color="record.cleanMode === 'time' ? 'blue' : 'orange'">
            {{ record.cleanMode === 'time' ? '时间过期' : '条件失效' }}
          </a-tag>
        </template>
        <template #cleanDays="{ record }">
          <span v-if="record.cleanMode === 'time'">{{ record.cleanDays }} 天后清除</span>
          <span v-else>条件失效即清除</span>
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
      
    </a-card>

    <a-modal
      v-model:visible="formVisible"
      :title="isEdit ? '编辑规则' : '新增规则'"
      @before-ok="handleFormSubmit"
      @cancel="formVisible = false"
      :width="520"
      :loading="formLoading"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="目标标签" :rules="[{ required: true, message: '请选择目标标签' }]">
          <a-select v-model="form.tagId" placeholder="选择标签" allow-create>
            <a-option v-for="tag in availableTags" :key="tag.id" :value="tag.id">{{ tag.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="清除模式" :rules="[{ required: true, message: '请选择清除模式' }]">
          <a-radio-group v-model="form.cleanMode">
            <a-radio value="time">时间过期</a-radio>
            <a-radio value="condition">条件失效</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item v-if="form.cleanMode === 'time'" label="过期天数">
          <a-input-number v-model="form.cleanDays" :min="1" :step="30" style="width: 100%" />
        </a-form-item>
        <a-form-item label="状态">
          <a-switch v-model="form.status" checked-value="active" unchecked-value="inactive" />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="form.remark" :rows="2" placeholder="规则备注" />
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

const form = reactive({
  tagId: null,
  cleanMode: 'time',
  cleanDays: 90,
  status: 'active',
  remark: ''
})

const columns = [
  { title: '目标标签', slotName: 'tagName', width: 180 },
  { title: '清除模式', slotName: 'cleanMode', width: 120 },
  { title: '清除条件', slotName: 'cleanDays', width: 160 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '备注', dataIndex: 'remark', ellipsis: true },
  { title: '操作', slotName: 'actions', width: 140 }
]

const loadData = async () => {
  loading.value = true
  try {
    const [configRes, tagsRes] = await Promise.all([
      api.getTagAutoCleanConfig(),
      api.getTagList()
    ])
    ruleList.value = configRes.data?.rules || []
    availableTags.value = tagsRes.data || []
    stats.value = configRes.data?.stats || {}
  } catch (err) {
    Message.error('加载数据失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

const showCreate = () => {
  isEdit.value = false
  Object.assign(form, { tagId: null, cleanMode: 'time', cleanDays: 90, status: 'active', remark: '' })
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
    const payload = { ...form }
    if (isEdit.value) payload.id = currentId.value
    await api.saveTagAutoCleanConfig(payload)
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
    content: '确定要删除该清除规则吗？',
    okText: '确认删除',
    onOk: async () => {
      try {
        await api.saveTagAutoCleanConfig({ id: record.id, _delete: true })
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
.tag-auto-clean-page {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
