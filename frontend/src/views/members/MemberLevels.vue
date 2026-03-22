<template>
  <div class="member-levels-page">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员等级</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6" v-for="lv in levelList" :key="lv.id">
        <a-card class="stat-card" hoverable>
          <a-statistic :title="lv.name" :value="lv.memberCount || 0">
            <template #prefix>
              <a-avatar :style="{ backgroundColor: getLevelColor(lv.id) }" :size="24">
                {{ lv.name?.charAt(0) }}
              </a-avatar>
            </template>
          </a-statistic>
          <div class="level-info">
            <span>折扣: {{ ((lv.discountRate || 1) * 100).toFixed(0) }}%</span>
            <span>积分倍率: {{ (lv.pointsRate || 1).toFixed(1) }}x</span>
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 操作栏 -->
    <a-card class="action-card">
      <a-space>
        <a-button type="primary" @click="showCreate">新增等级</a-button>
        <a-button @click="loadLevels">刷新</a-button>
      </a-space>
    </a-card>

    <!-- 等级列表 -->
    <a-card>
      <a-table :columns="columns" :data="levelList" :loading="loading" row-key="id" :pagination="false">
        <template #levelColor="{ record }">
          <a-badge :color="getLevelColor(record.id)" :text="record.name" />
        </template>
        <template #discountRate="{ record }">
          <a-tag color="orange">{{ ((record.discountRate || 1) * 100).toFixed(0) }}%</a-tag>
        </template>
        <template #pointsRate="{ record }">
          <span style="color: #ff6b00; font-weight: 600;">{{ (record.pointsRate || 1).toFixed(1) }}x</span>
        </template>
        <template #minAmount="{ record }">
          ¥{{ (record.minAmount || 0).toLocaleString() }}
        </template>
        <template #maxAmount="{ record }">
          {{ record.maxAmount ? '¥' + record.maxAmount.toLocaleString() : '无上限' }}
        </template>
        <template #benefits="{ record }">
          <a-space wrap>
            <a-tag v-for="(b, i) in (record.benefits || [])" :key="i" size="small">{{ b }}</a-tag>
            <span v-if="!record.benefits?.length" style="color: #999">无</span>
          </a-space>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 新增/编辑等级弹窗 -->
    <a-modal
      v-model:visible="formVisible"
      :title="isEdit ? '编辑等级' : '新增等级'"
      @before-ok="handleFormSubmit"
      @cancel="formVisible = false"
      :width="560"
      :loading="formLoading"
    >
      <a-form :model="form" layout="vertical" ref="formRef">
        <a-form-item label="等级名称" field="name" :rules="[{ required: true, message: '请输入等级名称' }]">
          <a-input v-model="form.name" placeholder="如：黄金会员" />
        </a-form-item>
        <a-form-item label="等级编码" field="code" :rules="[{ required: true, message: '请输入等级编码' }]">
          <a-input v-model="form.code" placeholder="如：GOLD" :disabled="isEdit" />
        </a-form-item>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="最低门槛(累计消费)" field="minAmount" :rules="[{ required: true, message: '请输入最低门槛' }]">
              <a-input-number v-model="form.minAmount" :min="0" :step="100" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="最高门槛(空=无上限)">
              <a-input-number v-model="form.maxAmount" :min="0" :step="100" style="width: 100%" placeholder="无上限" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="折扣率" field="discountRate" :rules="[{ required: true, message: '请输入折扣率' }]">
              <a-input-number v-model="form.discountRate" :min="0" :max="1" :step="0.05" style="width: 100%">
                <template #suffix>如 0.95</template>
              </a-input-number>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="积分倍率" field="pointsRate" :rules="[{ required: true, message: '请输入积分倍率' }]">
              <a-input-number v-model="form.pointsRate" :min="0" :step="0.5" style="width: 100%">
                <template #suffix>x</template>
              </a-input-number>
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="等级权益">
          <a-select v-model="form.benefits" multiple placeholder="选择权益（可多选）" allow-create>
            <a-option value="新人礼包">新人礼包</a-option>
            <a-option value="95折优惠">95折优惠</a-option>
            <a-option value="9折优惠">9折优惠</a-option>
            <a-option value="85折优惠">85折优惠</a-option>
            <a-option value="专属客服">专属客服</a-option>
            <a-option value="优先发货">优先发货</a-option>
            <a-option value="专属活动">专属活动</a-option>
            <a-option value="免费配送">免费配送</a-option>
            <a-option value="生日礼包">生日礼包</a-option>
            <a-option value="积分兑换">积分兑换</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="等级描述">
          <a-textarea v-model="form.description" :rows="2" placeholder="简要描述该等级" />
        </a-form-item>
        <a-form-item label="排序值">
          <a-input-number v-model="form.sort" :min="0" style="width: 100%" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import * as api from '@/api/member'

const levelList = ref([])
const loading = ref(false)
const formLoading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)

const form = reactive({
  name: '',
  code: '',
  minAmount: 0,
  maxAmount: undefined,
  discountRate: 1,
  pointsRate: 1,
  benefits: [],
  description: '',
  sort: 0
})

const columns = [
  { title: '等级', slotName: 'levelColor', width: 120 },
  { title: '编码', dataIndex: 'code', width: 100 },
  { title: '折扣率', slotName: 'discountRate', width: 100 },
  { title: '积分倍率', slotName: 'pointsRate', width: 110 },
  { title: '最低门槛', slotName: 'minAmount', width: 130 },
  { title: '最高门槛', slotName: 'maxAmount', width: 130 },
  { title: '等级权益', slotName: 'benefits' },
  { title: '会员数', dataIndex: 'memberCount', width: 90 },
  { title: '操作', slotName: 'actions', width: 140 }
]

const getLevelColor = (id) => {
  const colors = { 1: '#95de64', 2: '#1890ff', 3: '#faad14', 4: '#ff4d4f', 5: '#722ed1' }
  return colors[id] || 'gray'
}

const loadLevels = async () => {
  loading.value = true
  try {
    const res = await api.getLevelList()
    levelList.value = res.data || []
  } catch (err) {
    Message.error('加载等级列表失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

const showCreate = () => {
  isEdit.value = false
  Object.assign(form, {
    name: '', code: '', minAmount: 0, maxAmount: undefined,
    discountRate: 1, pointsRate: 1, benefits: [], description: '', sort: 0
  })
  formVisible.value = true
}

const showEdit = (record) => {
  isEdit.value = true
  currentId.value = record.id
  Object.assign(form, {
    name: record.name,
    code: record.code,
    minAmount: record.minAmount || 0,
    maxAmount: record.maxAmount || undefined,
    discountRate: record.discountRate || 1,
    pointsRate: record.pointsRate || 1,
    benefits: record.benefits || [],
    description: record.description || '',
    sort: record.sort || 0
  })
  formVisible.value = true
}

const handleFormSubmit = async (done) => {
  if (!form.name || !form.code) {
    Message.warning('请填写等级名称和编码')
    done(false)
    return
  }
  formLoading.value = true
  try {
    const payload = { ...form }
    if (isEdit.value) {
      await api.updateLevel(currentId.value, payload)
      Message.success('更新成功')
    } else {
      await api.createLevel(payload)
      Message.success('创建成功')
    }
    formVisible.value = false
    loadLevels()
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
    content: `确定要删除等级「${record.name}」吗？该等级下的会员将变为无等级状态。`,
    okText: '确认删除',
    onOk: async () => {
      try {
        await api.deleteLevel(record.id)
        Message.success('删除成功')
        loadLevels()
      } catch (err) {
        Message.error(err.message || '删除失败')
      }
    }
  })
}

onMounted(() => loadLevels())
</script>

<style scoped>
.member-levels-page {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}
.breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { border-radius: 8px; }
.level-info {
  display: flex;
  gap: 12px;
  margin-top: 4px;
  font-size: 12px;
  color: #666;
}
.action-card { margin-bottom: 16px; }
</style>
