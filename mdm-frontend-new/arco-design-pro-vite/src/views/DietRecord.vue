<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-food /> 饮食记录</a-space>
      </template>

      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6">
          <a-card>
            <a-statistic title="今日摄入" :value="stats.todayIntake" suffix="kcal" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="目标摄入" :value="stats.targetIntake" suffix="kcal" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="完成率" :value="stats.completionRate" suffix="%" />
          </a-card>
        </a-col>
      </a-row>

      <a-card title="饮食记录">
        <template #extra>
          <a-button type="primary" @click="handleAdd">
            <template #icon><icon-plus /></template>
            添加记录
          </a-button>
        </template>
        <a-table :columns="columns" :data="records">
          <template #type="{ record }">
            <a-tag :color="getTypeColor(record.type)">{{ record.type }}</a-tag>
          </template>
        </a-table>
      </a-card>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const stats = reactive({ todayIntake: 850, targetIntake: 1000, completionRate: 85 })

const columns = [
  { title: '时间', dataIndex: 'time' },
  { title: '类型', slotName: 'type' },
  { title: '食物', dataIndex: 'food' },
  { title: '热量', dataIndex: 'calories' }
]
const records = ref([
  { time: '2026-03-28 08:00', type: 'breakfast', food: '狗粮 100g', calories: '300kcal' },
  { time: '2026-03-28 12:00', type: 'lunch', food: '狗粮 150g', calories: '450kcal' }
])

const getTypeColor = (type) => ({ breakfast: 'orange', lunch: 'green', dinner: 'blue' }[type] || 'gray')
const handleAdd = () => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
