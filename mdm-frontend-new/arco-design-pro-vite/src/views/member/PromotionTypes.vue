<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="promotion-types-page">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>促销活动类型</a-breadcrumb-item>
    </a-breadcrumb>

    <a-alert style="margin-bottom: 16px;">
      <template #title>促销活动类型说明</template>
      管理系统中已配置的所有促销活动类型及其统计信息。支持按类型快速跳转查看具体活动。
    </a-alert>

    <a-row :gutter="16">
      <a-col :span="8" v-for="type in promoTypes" :key="type.code">
        <a-card hoverable class="type-card" @click="goToPromo(type.code)">
          <a-space direction="vertical" :size="8" style="width: 100%;">
            <a-space>
              <a-avatar :style="{ backgroundColor: type.color }" :size="40">{{ type.name.charAt(0) }}</a-avatar>
              <a-statistic :title="type.name" :value="type.count || 0" />
            </a-space>
            <a-divider style="margin: 8px 0;" />
            <a-space wrap>
              <a-tag v-for="(tag, i) in type.tags" :key="i" size="small">{{ tag }}</a-tag>
            </a-space>
          </a-space>
        </a-card>
      </a-col>
    </a-row>

    <!-- 类型概览表格 -->
    <a-card title="促销活动概览" style="margin-top: 16px;">
      <a-table :columns="columns" :data="promoTypes" :loading="loading" row-key="code" :pagination="false">
        <template #typeName="{ record }">
          <a-space>
            <a-avatar :style="{ backgroundColor: record.color, fontSize: '12px' }" :size="28">{{ record.name.charAt(0) }}</a-avatar>
            {{ record.name }}
          </a-space>
        </template>
      </a-table>
        <template #status="{ record }"><a-tag :color="record.enabled ? 'green' : 'gray'">{{ record.enabled ? '已启用' : '已禁用' }}</a-tag></template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="goToPromo(record.code)">查看活动</a-button>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import * as api from '@/api/marketing'

const router = useRouter()
const promoTypes = ref([])
const loading = ref(false)

const columns = [
  { title: '类型', slotName: 'typeName', width: 200 },
  { title: '活动数', dataIndex: 'count', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 120 }
]

const promoTypeMeta = [
  { code: 'buy_gift', name: '买赠促销', color: '#165DFF', route: '/member/promotions/buy-gift', tags: ['买商品赠赠品', '多买多赠'] },
  { code: 'direct_reduce', name: '直减促销', color: '#00D084', route: '/member/promotions/direct-reduce', tags: ['立减优惠', '无门槛'] },
  { code: 'amount_reduce', name: '满额减促销', color: '#9C27B0', route: '/member/promotions/amount-reduce', tags: ['满减优惠', '阶梯档位'] },
  { code: 'amount_discount', name: '满额折促销', color: '#00BCD4', route: '/member/promotions/amount-discount', tags: ['满折优惠', '阶梯折扣'] },
  { code: 'vip_exclusive', name: '最高等级促销', color: '#FF9800', route: '/member/promotions/vip-exclusive', tags: ['VIP专享', '专属优惠'] },
  { code: 'redpacket', name: '红包促销', color: '#FF4D4F', route: '/member/redpackets', tags: ['现金红包', '发放核销'] }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await api.getPromotionTypes()
    const list = res.data || []
    promoTypes.value = promoTypeMeta.map(meta => {
      const found = list.find(l => l.code === meta.code) || {}
      return { ...meta, count: found.count || 0, enabled: found.enabled !== false }
    })
  } catch (err) {
    promoTypes.value = promoTypeMeta.map(meta => ({ ...meta, count: 0, enabled: true }))
  } finally {
    loading.value = false
  }
}

const goToPromo = (code) => {
  const meta = promoTypeMeta.find(m => m.code === code)
  if (meta) router.push(meta.route)
  else Message.info('该类型活动页面正在开发中')
}

onMounted(() => loadData())
</script>

<style scoped>
.promotion-types-page { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.breadcrumb { margin-bottom: 16px; }
.type-card { cursor: pointer; border-radius: 8px; }
.type-card:hover { box-shadow: 0 4px 12px rgba(0,0,0,0.1); }
</style>
