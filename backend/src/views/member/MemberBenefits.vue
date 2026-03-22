<template>
  <div class="member-benefits-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员服务</a-breadcrumb-item>
      <a-breadcrumb-item>会员权益管理</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8"><a-card hoverable><a-statistic title="权益总数" :value="stats.totalCount || 0" /></a-card></a-col>
      <a-col :span="8"><a-card hoverable><a-statistic title="已配置权益" :value="stats.configuredCount || 0" /></a-card></a-col>
      <a-col :span="8"><a-card hoverable><a-statistic title="本月使用" :value="stats.monthlyUsed || 0" /></a-card></a-col>
    </a-row>

    <a-card class="action-card">
      <a-space>
        <a-button type="primary" @click="showCreate">新增权益</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </a-card>

    <a-card>
      <a-table :columns="columns" :data="benefitList" :loading="loading" row-key="id" :pagination="{ pageSize: 10 }">
        <template #name="{ record }">
          <a-space>
            <a-avatar :style="{ backgroundColor: record.color || '#165DFF' }" :size="28">{{ record.name.charAt(0) }}</a-avatar>
            {{ record.name }}
          </a-space>
        </template>
        <template #type="{ record }"><a-tag>{{ record.type }}</a-tag></template>
        <template #applicableLevels="{ record }">
          <a-space wrap>
            <a-tag v-for="(lv, i) in (record.applicableLevels || [])" :key="i" size="small">{{ lv }}</a-tag>
          </a-space>
        </template>
        <template #status="{ record }"><a-tag :color="record.status === 'active' ? 'green' : 'gray'">{{ record.status === 'active' ? '已启用' : '已禁用' }}</a-tag></template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑权益' : '新增权益'"
      @before-ok="handleFormSubmit" @cancel="formVisible = false" :width="560" :loading="formLoading">
      <a-form :model="form" layout="vertical">
        <a-form-item label="权益名称" :rules="[{ required: true, message: '请输入权益名称' }]">
          <a-input v-model="form.name" placeholder="如：专属客服" />
        </a-form-item>
        <a-form-item label="权益类型">
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="折扣">折扣</a-option>
            <a-option value="礼品">礼品</a-option>
            <a-option value="服务">服务</a-option>
            <a-option value="积分">积分</a-option>
            <a-option value="优先">优先权</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="权益描述">
          <a-textarea v-model="form.description" :rows="2" placeholder="权益详细描述" />
        </a-form-item>
        <a-form-item label="权益颜色">
          <a-color-picker v-model="form.color" :show-alpha="false" :presets="colorPresets" />
        </a-form-item>
        <a-form-item label="适用等级">
          <a-select v-model="form.applicableLevels" multiple placeholder="选择适用等级（不选=全部）" allow-create>
            <a-option v-for="lv in levelList" :key="lv.id" :value="lv.name">{{ lv.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状态"><a-switch v-model="form.status" checked-value="active" unchecked-value="inactive" /></a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import * as api from '@/api/marketing'
import * as memberApi from '@/api/member'

const benefitList = ref([])
const levelList = ref([])
const stats = ref({})
const loading = ref(false)
const formLoading = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)

const colorPresets = ['#165DFF', '#00D084', '#FF6B6B', '#FFC107', '#9C27B0', '#FF9800', '#E91E63', '#3F51B5']

const form = reactive({ name: '', type: '折扣', description: '', color: '#165DFF', applicableLevels: [], status: 'active' })

const columns = [
  { title: '权益名称', slotName: 'name', width: 200 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '适用等级', slotName: 'applicableLevels', width: 200 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 140 }
]

const loadData = async () => {
  loading.value = true
  try {
    const [benefitRes, lvRes] = await Promise.all([api.getMemberBenefitList(), memberApi.getLevelList()])
    benefitList.value = benefitRes.data || []
    stats.value = benefitRes.data?.stats || {}
    levelList.value = lvRes.data || []
  } catch (err) { Message.error('加载数据失败: ' + err.message) } finally { loading.value = false }
}

const showCreate = () => { isEdit.value = false; Object.assign(form, { name: '', type: '折扣', description: '', color: '#165DFF', applicableLevels: [], status: 'active' }); formVisible.value = true }
const showEdit = (record) => { isEdit.value = true; currentId.value = record.id; Object.assign(form, { ...record, applicableLevels: record.applicableLevels ? [...record.applicableLevels] : [] }); formVisible.value = true }

const handleFormSubmit = async (done) => {
  if (!form.name) { Message.warning('请输入权益名称'); done(false); return }
  formLoading.value = true
  try {
    if (isEdit.value) await api.updateMemberBenefit(currentId.value, { ...form })
    else await api.createMemberBenefit({ ...form })
    Message.success('保存成功'); formVisible.value = false; loadData(); done(true)
  } catch (err) { Message.error(err.message || '操作失败'); done(false) } finally { formLoading.value = false }
}

const handleDelete = (record) => {
  Modal.warning({ title: '确认删除', content: `确定删除权益「${record.name}」吗？`, okText: '确认删除',
    onOk: async () => { try { await api.deleteMemberBenefit(record.id); Message.success('删除成功'); loadData() } catch (err) { Message.error(err.message || '删除失败') } }
  })
}

onMounted(() => loadData())
</script>

<style scoped>
.member-benefits-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
