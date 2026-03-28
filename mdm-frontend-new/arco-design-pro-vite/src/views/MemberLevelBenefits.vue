<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-gift /> 会员等级权益配置</a-space>
      </template>
      <template #extra>
        <a-button type="primary" @click="handleSave">
          <template #icon><icon-save /></template>
          保存配置
        </a-button>
      </template>

      <a-table :columns="columns" :data="benefitData" :bordered="true" :pagination="false">
        <template #level="{ record }">
          <a-tag :color="record.color">{{ record.level }}</a-tag>
        </template>
        <template #discount="{ record }">
          <a-input-number v-model="record.discount" :min="0" :max="100" :precision="0" suffix="%" />
        </template>
        <template #pointRate="{ record }">
          <a-input-number v-model="record.pointRate" :min="1" :max="10" :precision="1" suffix="x" />
        </template>
        <template #privileges="{ record }">
          <a-checkbox-group v-model="record.privileges">
            <a-checkbox value="ai_unlimited">无限AI对话</a-checkbox>
            <a-checkbox value="priority">优先客服</a-checkbox>
            <a-checkbox value="exclusive">专属内容</a-checkbox>
          </a-checkbox-group>
        </template>
      </a-table>

      <a-divider>权益预览面板</a-divider>
      <a-row :gutter="16">
        <a-col :span="6" v-for="level in benefitData" :key="level.level">
          <a-card :class="'level-card level-' + level.level.toLowerCase()">
            <template #title>
              <a-tag :color="level.color">{{ level.level }}</a-tag>
            </template>
            <a-statistic title="折扣" :value="level.discount" suffix="%" />
            <a-divider />
            <a-list size="small">
              <a-list-item v-for="p in level.privileges" :key="p">{{ p }}</a-list-item>
            </a-list>
          </a-card>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const columns = [
  { title: '等级', slotName: 'level', width: 100 },
  { title: '升级门槛(成长值)', dataIndex: 'threshold', width: 150 },
  { title: '折扣比例', slotName: 'discount', width: 150 },
  { title: '积分倍率', slotName: 'pointRate', width: 150 },
  { title: '特权', slotName: 'privileges' }
]

const benefitData = ref([
  { level: '普通', color: 'gray', threshold: 0, discount: 100, pointRate: 1, privileges: [] },
  { level: '白银', color: 'blue', threshold: 1000, discount: 95, pointRate: 1.2, privileges: ['ai_unlimited'] },
  { level: '黄金', color: 'orange', threshold: 5000, discount: 90, pointRate: 1.5, privileges: ['ai_unlimited', 'priority'] },
  { level: '铂金', color: 'purple', threshold: 20000, discount: 85, pointRate: 2, privileges: ['ai_unlimited', 'priority', 'exclusive'] },
  { level: '钻石', color: 'red', threshold: 50000, discount: 80, pointRate: 3, privileges: ['ai_unlimited', 'priority', 'exclusive'] }
])

const handleSave = () => { }
</script>

<style scoped>
.container { padding: 16px; }
.level-card { border-left: 3px solid; }
.level-普通 { border-color: gray; }
.level-白银 { border-color: #409EFF; }
.level-黄金 { border-color: #E6A23C; }
.level-铂金 { border-color: #9C27B0; }
.level-钻石 { border-color: #F56C6C; }
</style>
