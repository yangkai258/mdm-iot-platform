<template>
  <div class="container">
    <a-row :gutter="16">
      <a-col :span="8">
        <a-card title="基本信息">
          <a-descriptions :column="1" bordered>
            <a-descriptions-item label="会员ID">{{ profile.id }}</a-descriptions-item>
            <a-descriptions-item label="昵称">{{ profile.nickname }}</a-descriptions-item>
            <a-descriptions-item label="手机号">{{ profile.phone }}</a-descriptions-item>
            <a-descriptions-item label="注册时间">{{ profile.registerTime }}</a-descriptions-item>
            <a-descriptions-item label="会员等级">
              <a-tag :color="profile.levelColor">{{ profile.level }}</a-tag>
            </a-descriptions-item>
          </a-descriptions>
        </a-card>

        <a-card title="标签云" style="margin-top: 16px">
          <a-tag v-for="tag in profile.tags" :key="tag" :color="tag.color" : closable>{{ tag.name }}</a-tag>
          <a-button type="text" @click="handleAddTag">
            <template #icon><icon-plus /></template>
            添加标签
          </a-button>
        </a-card>

        <a-card title="价值分层" style="margin-top: 16px">
          <a-progress :percent="profile.valueScore" :stroke-color="getValueColor(profile.valueLevel)" />
          <a-space style="margin-top: 8px">
            <a-tag :color="getValueColor(profile.valueLevel)">{{ profile.valueLevel }}</a-tag>
            <span class="score">{{ profile.valueScore }}分</span>
          </a-space>
        </a-card>
      </a-col>

      <a-col :span="16">
        <a-card title="行为分析">
          <a-row :gutter="16">
            <a-col :span="8">
              <a-statistic title="活跃度" :value="profile.activity" suffix="分" />
            </a-col>
            <a-col :span="8">
              <a-statistic title="消费能力" :value="profile.spending" suffix="分" />
            </a-col>
            <a-col :span="8">
              <a-statistic title="互动频率" :value="profile.interaction" suffix="次/周" />
            </a-col>
          </a-row>
          <a-divider />
          <a-chart :option="behaviorChart" style="height: 200px" />
        </a-card>

        <a-card title="兴趣偏好" style="margin-top: 16px">
          <a-tag v-for="interest in profile.interests" :key="interest.name" :color="interest.color">
            {{ interest.name }} ({{ interest.score }}%)
          </a-tag>
        </a-card>

        <a-card title="消费记录" style="margin-top: 16px">
          <a-table :columns="spendColumns" :data="spendData" size="small" :pagination="pagination" />
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const profile = reactive({
  id: 'M001', nickname: '张三', phone: '138****8888', registerTime: '2025-01-15',
  level: '黄金', levelColor: 'orange', tags: [
    { name: '高价值', color: 'gold' }, { name: '活跃用户', color: 'green' }, { name: '宠物爱好者', color: 'blue' }
  ],
  valueLevel: '高价值用户', valueScore: 85, activity: 92, spending: 78, interaction: 15,
  interests: [
    { name: '宠物用品', color: 'orange', score: 85 },
    { name: '智能设备', color: 'blue', score: 72 },
    { name: '健康管理', color: 'green', score: 65 }
  ]
})

const pagination = reactive({ current: 1, pageSize: 5, total: 20 })
const spendColumns = [
  { title: '时间', dataIndex: 'time' },
  { title: '类型', dataIndex: 'type' },
  { title: '金额', dataIndex: 'amount' }
]
const spendData = ref([
  { time: '2026-03-28', type: '订阅续费', amount: '¥99' },
  { time: '2026-03-15', type: '积分兑换', amount: '-200积分' }
])

const behaviorChart = reactive({
  tooltip: {}, radar: {
    indicator: [
      { name: 'AI对话', max: 100 }, { name: '设备使用', max: 100 },
      { name: '社交互动', max: 100 }, { name: '内容消费', max: 100 }
    ]
  },
  series: [{ type: 'radar', data: [{ value: [85, 72, 68, 90], name: '行为分析' }] }]
})

const getValueColor = (level) => ({ '高价值用户': 'gold', '中价值用户': 'blue', '低价值用户': 'gray' }[level] || 'gray')
const handleAddTag = () => { }
</script>

<style scoped>
.container { padding: 16px; }
.score { font-size: 12px; color: #909399; }
</style>
