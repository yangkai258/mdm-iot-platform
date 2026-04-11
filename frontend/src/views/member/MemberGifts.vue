<template>
  <div class="member-gifts-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员礼包</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6"><a-card hoverable><a-statistic title="礼包活动数" :value="stats.totalCount || 0" /></a-card></a-col>
      <a-col :span="6"><a-card hoverable><a-statistic title="已发放礼包" :value="stats.grantedCount || 0" /></a-card></a-col>
      <a-col :span="6"><a-card hoverable><a-statistic title="已领取礼包" :value="stats.claimedCount || 0" /></a-card></a-col>
      <a-col :span="6"><a-card hoverable><a-statistic title="礼包价值总计" :value="(stats.totalValue || 0).toLocaleString() + ' 元'" /></a-card></a-col>
    </a-row>

    <a-card class="action-card">
      <a-space>
        <a-button type="primary" @click="showCreate">新建礼包</a-button>
        <a-button @click="loadData">刷新</a-button>
        <router-link to="/member/gift-records">
          <a-button>发放明细</a-button>
        </router-link>
      </a-space>
    </a-card>

    <a-card>
      <a-table :columns="columns" :data="giftList" :loading="loading" row-key="id" :pagination="{ pageSize: 10 }">
        <template #name="{ record }"><a-tag color="arcoblue">{{ record.name }}</a-tag></template>
        <template #items="{ record }">
          <a-space wrap>
            <a-tag v-for="(item, i) in (record.items || [])" :key="i" size="small">{{ item.name }}×{{ item.count }}</a-tag>
          </a-space>
        </template>
      </a-table>
        <template #totalValue="{ record }"><span style="color:#ff4d4f; font-weight:600;">¥{{ record.totalValue || 0 }}</span></template>
        <template #validPeriod="{ record }">{{ record.startTime?.slice(0,10) }} ~ {{ record.endTime?.slice(0,10) }}</template>
        <template #status="{ record }"><a-tag :color="getStatusColor(record.status)">{{ getStatusName(record.status) }}</a-tag></template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showGrant(record)">发放</a-button>
            <a-button type="text" size="small" @click="showEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      
    </a-card>

    <!-- 新增/编辑礼包弹窗 -->
    <a-modal v-model:visible="formVisible" :title="isEdit ? '编辑礼包' : '新建礼包'"
      @before-ok="handleFormSubmit" @cancel="formVisible = false" :width="600" :loading="formLoading">
      <a-form :model="form" layout="vertical">
        <a-form-item label="礼包名称" :rules="[{ required: true, message: '请输入礼包名称' }]">
          <a-input v-model="form.name" placeholder="如：新人入会礼包" />
        </a-form-item>
        <a-form-item label="礼包内容">
          <a-space direction="vertical" style="width: 100%;">
            <a-row v-for="(item, i) in form.items" :key="i" :gutter="8" align="center">
              <a-col :span="8"><a-input v-model="item.name" placeholder="物品名称" /></a-col>
              <a-col :span="6"><a-input-number v-model="item.count" :min="1" placeholder="数量" style="width: 100%;" /></a-col>
              <a-col :span="6"><a-input-number v-model="item.value" :min="0" :precision="2" placeholder="价值(元)" style="width: 100%;" /></a-col>
              <a-col :span="4"><a-button type="text" status="danger" @click="form.items.splice(i, 1)">删除</a-button></a-col>
            </a-row>
            <a-button type="dashed" @click="form.items.push({ name: '', count: 1, value: 0 })">+ 添加物品</a-button>
          </a-space>
        </a-form-item>
        <a-row :gutter="12">
          <a-col :span="12"><a-form-item label="开始时间"><a-date-picker v-model="form.startTime" style="width:100%;" /></a-form-item></a-col>
          <a-col :span="12"><a-form-item label="结束时间"><a-date-picker v-model="form.endTime" style="width:100%;" /></a-form-item></a-col>
        </a-row>
        <a-form-item label="发放方式">
          <a-radio-group v-model="form.grantType">
            <a-radio value="manual">手动发放</a-radio>
            <a-radio value="level">按等级发放</a-radio>
            <a-radio value="auto">入会自动</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="状态"><a-switch v-model="form.status" checked-value="active" unchecked-value="paused" /></a-form-item>
      </a-form>
    </a-modal>

    <!-- 发放礼包弹窗 -->
    <a-modal v-model:visible="grantVisible" title="发放礼包" @before-ok="handleGrant" @cancel="grantVisible = false"
      :width="480" :loading="grantLoading">
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
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import * as api from '@/api/marketing'
import * as memberApi from '@/api/member'

const giftList = ref([])
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

const form = reactive({
  name: '', items: [{ name: '', count: 1, value: 0 }],
  startTime: null, endTime: null, grantType: 'manual', status: 'active'
})

const grantForm = reactive({ type: 'member', memberId: '', levelId: null, tagId: null })

const columns = [
  { title: '礼包名称', slotName: 'name', width: 160 },
  { title: '礼包内容', slotName: 'items' },
  { title: '总价值', slotName: 'totalValue', width: 120 },
  { title: '发放方式', dataIndex: 'grantType', width: 120 },
  { title: '有效期', slotName: 'validPeriod', width: 220 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 180 }
]

const getStatusColor = (s) => ({ active: 'green', paused: 'orange', finished: 'gray' }[s] || 'gray')
const getStatusName = (s) => ({ active: '进行中', paused: '已暂停', finished: '已结束' }[s] || s)

const loadData = async () => {
  loading.value = true
  try {
    const [giftRes, tagsRes, lvRes] = await Promise.all([
      api.getMemberGiftList(),
      api.getTagList(),
      memberApi.getLevelList()
    ])
    giftList.value = giftRes.data?.list || []
    stats.value = giftRes.data?.stats || {}
    tagList.value = tagsRes.data || []
    levelList.value = lvRes.data || []
  } catch (err) { Message.error('加载数据失败: ' + err.message) } finally { loading.value = false }
}

const showCreate = () => {
  isEdit.value = false
  Object.assign(form, { name: '', items: [{ name: '', count: 1, value: 0 }], startTime: null, endTime: null, grantType: 'manual', status: 'active' })
  formVisible.value = true
}
const showEdit = (record) => {
  isEdit.value = true; currentId.value = record.id
  Object.assign(form, { ...record, items: record.items ? [...record.items] : [{ name: '', count: 1, value: 0 }] })
  formVisible.value = true
}
const showGrant = (record) => { currentId.value = record.id; Object.assign(grantForm, { type: 'member', memberId: '', levelId: null, tagId: null }); grantVisible.value = true }

const handleFormSubmit = async (done) => {
  formLoading.value = true
  try {
    if (isEdit.value) await api.updateMemberGift(currentId.value, { ...form })
    else await api.createMemberGift({ ...form })
    Message.success('保存成功'); formVisible.value = false; loadData(); done(true)
  } catch (err) { Message.error(err.message || '操作失败'); done(false) } finally { formLoading.value = false }
}

const handleGrant = async (done) => {
  grantLoading.value = true
  try {
    await api.grantMemberGift({ giftId: currentId.value, ...grantForm })
    Message.success('发放成功'); grantVisible.value = false; done(true)
  } catch (err) { Message.error(err.message || '发放失败'); done(false) } finally { grantLoading.value = false }
}

const handleDelete = (record) => {
  Modal.warning({ title: '确认删除', content: `确定删除礼包「${record.name}」吗？`, okText: '确认删除',
    onOk: async () => { try { await api.deleteMemberGift(record.id); Message.success('删除成功'); loadData() } catch (err) { Message.error(err.message || '删除失败') } }
  })
}

onMounted(() => loadData())
</script>

<style scoped>
.member-gifts-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
