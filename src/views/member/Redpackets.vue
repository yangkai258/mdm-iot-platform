<template>
  <div class="redpackets-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员红包</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card hoverable>
          <a-statistic title="红包活动数" :value="stats.totalCount || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card hoverable>
          <a-statistic title="已发放红包" :value="stats.grantedCount || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card hoverable>
          <a-statistic title="已核销红包" :value="stats.usedCount || 0" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card hoverable>
          <a-statistic title="红包金额总计">
            <template #suffix>元</template>
            <template #value>{{ (stats.totalAmount || 0).toLocaleString() }}</template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <a-card class="action-card">
      <a-space>
        <a-button type="primary" @click="showCreate">新建红包</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
      <template #extra>
        <a-radio-group v-model="listType" type="button" size="small">
          <a-radio value="all">全部</a-radio>
          <a-radio value="active">进行中</a-radio>
          <a-radio value="paused">已暂停</a-radio>
          <a-radio value="finished">已结束</a-radio>
        </a-radio-group>
      </template>
    </a-card>

    <a-card>
      <a-table :columns="columns" :data="filteredList" :loading="loading" row-key="id" :pagination="{ pageSize: 10 }">
        <template #name="{ record }">
          <a-tag color="red">{{ record.name }}</a-tag>
        </template>
        <template #amount="{ record }">
          <span style="color: #ff4d4f; font-weight: 600;">¥{{ record.amount }}</span>
        </template>
        <template #totalCount="{ record }">
          {{ record.totalCount || '无限制' }}
        </template>
        <template #usedCount="{ record }">
          {{ record.usedCount || 0 }} / {{ record.totalCount || '∞' }}
        </template>
        <template #validPeriod="{ record }">
          {{ record.startTime?.slice(0, 10) }} ~ {{ record.endTime?.slice(0, 10) }}
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusName(record.status) }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showGrant(record)">发放</a-button>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 新增/编辑红包弹窗 -->
    <a-modal
      v-model:visible="formVisible"
      :title="isEdit ? '编辑红包' : '新建红包'"
      @before-ok="handleFormSubmit"
      @cancel="formVisible = false"
      :width="560"
      :loading="formLoading"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="红包名称" :rules="[{ required: true, message: '请输入红包名称' }]">
          <a-input v-model="form.name" placeholder="如：新人专享红包" />
        </a-form-item>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="红包金额(元)" :rules="[{ required: true, message: '请输入金额' }]">
              <a-input-number v-model="form.amount" :min="0.01" :precision="2" :step="1" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="发放总量">
              <a-input-number v-model="form.totalCount" :min="0" placeholder="0=无限制" style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="使用门槛">
          <a-input-number v-model="form.minAmount" :min="0" :precision="2" placeholder="满多少可用（0=无门槛）" style="width: 100%" />
        </a-form-item>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="开始时间">
              <a-date-picker v-model="form.startTime" style="width: 100%;" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="结束时间">
              <a-date-picker v-model="form.endTime" style="width: 100%;" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="状态">
          <a-radio-group v-model="form.status">
            <a-radio value="active">进行中</a-radio>
            <a-radio value="paused">暂停</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="form.remark" :rows="2" placeholder="红包说明" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 发放红包弹窗 -->
    <a-modal
      v-model:visible="grantVisible"
      title="发放红包"
      @before-ok="handleGrant"
      @cancel="grantVisible = false"
      :width="480"
      :loading="grantLoading"
    >
      <a-form :model="grantForm" layout="vertical">
        <a-form-item label="发放方式">
          <a-radio-group v-model="grantForm.type">
            <a-radio value="member">指定会员</a-radio>
            <a-radio value="level">按等级发放</a-radio>
            <a-radio value="tag">按标签发放</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item v-if="grantForm.type === 'member'" label="会员ID/手机号">
          <a-input v-model="grantForm.memberId" placeholder="输入会员ID或手机号" />
        </a-form-item>
        <a-form-item v-if="grantForm.type === 'level'" label="会员等级">
          <a-select v-model="grantForm.levelId" placeholder="选择等级" allow-create>
            <a-option v-for="lv in levelList" :key="lv.id" :value="lv.id">{{ lv.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="grantForm.type === 'tag'" label="会员标签">
          <a-select v-model="grantForm.tagId" placeholder="选择标签" allow-create>
            <a-option v-for="tag in tagList" :key="tag.id" :value="tag.id">{{ tag.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="发放数量">
          <a-input-number v-model="grantForm.count" :min="1" style="width: 100%" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import * as api from '@/api/marketing'
import * as memberApi from '@/api/member'

const redpacketList = ref([])
const tagList = ref([])
const levelList = ref([])
const stats = ref({})
const loading = ref(false)
const formLoading = ref(false)
const grantLoading = ref(false)
const formVisible = ref(false)
const grantVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const listType = ref('all')

const form = reactive({
  name: '', amount: 5, totalCount: 0, minAmount: 0,
  startTime: null, endTime: null, status: 'active', remark: ''
})

const grantForm = reactive({
  type: 'member', memberId: '', levelId: null, tagId: null, count: 1
})

const columns = [
  { title: '红包名称', slotName: 'name', width: 160 },
  { title: '金额', slotName: 'amount', width: 100 },
  { title: '发放总量', slotName: 'totalCount', width: 120 },
  { title: '已发放/核销', slotName: 'usedCount', width: 140 },
  { title: '有效期', slotName: 'validPeriod', width: 220 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 180 }
]

const filteredList = computed(() => {
  if (listType.value === 'all') return redpacketList.value
  return redpacketList.value.filter(r => r.status === listType.value)
})

const getStatusColor = (s) => ({ active: 'green', paused: 'orange', finished: 'gray' }[s] || 'gray')
const getStatusName = (s) => ({ active: '进行中', paused: '已暂停', finished: '已结束' }[s] || s)

const loadData = async () => {
  loading.value = true
  try {
    const [rpRes, tagsRes, lvRes] = await Promise.all([
      api.getRedpacketList(),
      api.getTagList(),
      memberApi.getLevelList()
    ])
    redpacketList.value = rpRes.data?.list || []
    stats.value = rpRes.data?.stats || {}
    tagList.value = tagsRes.data || []
    levelList.value = lvRes.data || []
  } catch (err) {
    Message.error('加载数据失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

const showCreate = () => {
  isEdit.value = false
  Object.assign(form, { name: '', amount: 5, totalCount: 0, minAmount: 0, startTime: null, endTime: null, status: 'active', remark: '' })
  formVisible.value = true
}

const showEdit = (record) => {
  isEdit.value = true
  currentId.value = record.id
  Object.assign(form, { ...record })
  formVisible.value = true
}

const showGrant = (record) => {
  currentId.value = record.id
  Object.assign(grantForm, { type: 'member', memberId: '', levelId: null, tagId: null, count: 1 })
  grantVisible.value = true
}

const handleFormSubmit = async (done) => {
  formLoading.value = true
  try {
    if (isEdit.value) {
      await api.updateRedpacket(currentId.value, { ...form })
      Message.success('更新成功')
    } else {
      await api.createRedpacket({ ...form })
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

const handleGrant = async (done) => {
  grantLoading.value = true
  try {
    await api.grantRedpacket({ redpacketId: currentId.value, ...grantForm })
    Message.success('发放成功')
    grantVisible.value = false
    loadData()
    done(true)
  } catch (err) {
    Message.error(err.message || '发放失败')
    done(false)
  } finally {
    grantLoading.value = false
  }
}

const handleDelete = (record) => {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除红包「${record.name}」吗？`,
    okText: '确认删除',
    onOk: async () => {
      try {
        await api.deleteRedpacket(record.id)
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
.redpackets-page {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
