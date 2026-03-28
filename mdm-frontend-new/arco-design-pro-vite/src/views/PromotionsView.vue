<template>
  <div class="container">
    <a-card>
      <template #title><icon-promotion /> 促销管理</template>
      <template #extra>
        <a-space>
          <a-button type="primary" @click="handleCreate"><template #icon><icon-plus /></template>创建促销</a-button>
        </a-space>
      </template>
      <a-table :columns="columns" :data="promotions" :pagination="pagination">
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
        <template #type="{ record }">
          <a-tag>{{ getTypeText(record.type) }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-link @click="handleEdit(record)">编辑</a-link>
          <a-divider direction="vertical" />
          <a-link @click="handleDelete(record)" status="error">删除</a-link>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑促销' : '创建促销'" @ok="handleSave">
      <a-form :model="form" layout="vertical">
        <a-form-item label="促销名称" required>
          <a-input v-model="form.name" placeholder="请输入促销名称" />
        </a-form-item>
        <a-form-item label="促销类型" required>
          <a-select v-model="form.type" placeholder="请选择类型">
            <a-option value="discount">折扣</a-option>
            <a-option value="coupon">优惠券</a-option>
            <a-option value="gift">赠品</a-option>
            <a-option value="flash">限时秒杀</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="促销时间">
          <a-range-picker v-model="form.timeRange" style="width: 100%" />
        </a-form-item>
        <a-form-item label="促销内容">
          <a-textarea v-model="form.content" placeholder="请输入促销内容" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const columns = [
  { title: '促销名称', dataIndex: 'name' },
  { title: '类型', slotName: 'type' },
  { title: '开始时间', dataIndex: 'startTime' },
  { title: '结束时间', dataIndex: 'endTime' },
  { title: '状态', slotName: 'status' },
  { title: '操作', slotName: 'actions', width: 150 }
]
const pagination = { pageSize: 10 }
const promotions = ref([
  { id: 1, name: '春节大促', type: 'discount', startTime: '2026-01-20', endTime: '2026-02-10', status: 'ended' },
  { id: 2, name: '新品上市', type: 'coupon', startTime: '2026-03-01', endTime: '2026-03-31', status: 'active' },
  { id: 3, name: '会员日', type: 'flash', startTime: '2026-03-15', endTime: '2026-03-15', status: 'pending' }
])
const editVisible = ref(false)
const isEdit = ref(false)
const form = reactive({ id: null, name: '', type: '', timeRange: [], content: '' })
const getStatusColor = (s) => ({ active: 'green', pending: 'blue', ended: 'gray' }[s] || 'gray')
const getStatusText = (s) => ({ active: '进行中', pending: '未开始', ended: '已结束' }[s] || s)
const getTypeText = (t) => ({ discount: '折扣', coupon: '优惠券', gift: '赠品', flash: '限时秒杀' }[t] || t)
const handleCreate = () => { isEdit.value = false; form.id = null; Object.assign(form, { name: '', type: '', timeRange: [], content: '' }); editVisible.value = true }
const handleEdit = (r) => { isEdit.value = true; Object.assign(form, r); editVisible.value = true }
const handleSave = () => { editVisible.value = false }
const handleDelete = (r) => {}
</script>
<style scoped>
.container { padding: 16px; }
</style>
