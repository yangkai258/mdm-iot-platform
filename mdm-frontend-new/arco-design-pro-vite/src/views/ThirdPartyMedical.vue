<template>
  <div class="medical-records-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="已对接医院" :value="5" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="健康记录" :value="128" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="待同步" :value="3" status="warning" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <span>宠物医疗记录</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="records" title="健康记录">
          <a-table :columns="columns" :data="records" :pagination="pagination">
            <template #type="{ record }">
              <a-tag>{{ record.typeText }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="hospitals" title="对接医院">
          <a-row :gutter="16">
            <a-col :span="6" v-for="h in hospitals" :key="h.id">
              <a-card size="small" class="hospital-card">
                <div class="hospital-name">{{ h.name }}</div>
                <div class="hospital-type">{{ h.type }}</div>
                <a-tag :color="h.connected ? 'green' : 'gray'">{{ h.connected ? '已连接' : '未连接' }}</a-tag>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
        
        <a-tab-pane key="vaccinations" title="疫苗记录">
          <a-table :columns="vaccineColumns" :data="vaccinations" :pagination="pagination" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 5 });

const records = ref([
  { id: 1, petName: '小黄', type: 'checkup', typeText: '体检', hospital: '宠物医院A', date: '2026-03-15', result: '正常', doctor: '张医生' },
  { id: 2, petName: '小黄', type: 'vaccine', typeText: '疫苗', hospital: '宠物医院A', date: '2026-03-01', result: '已完成', doctor: '张医生' },
  { id: 3, petName: '小红', type: 'disease', typeText: '疾病', hospital: '宠物医院B', date: '2026-02-20', result: '皮肤病', doctor: '李医生' },
]);

const hospitals = ref([
  { id: 1, name: '宠物医院A', type: '综合医院', connected: true },
  { id: 2, name: '宠物医院B', type: '专科医院', connected: true },
  { id: 3, name: '宠物医院C', type: '急诊中心', connected: false },
]);

const vaccinations = ref([
  { id: 1, petName: '小黄', vaccine: '狂犬疫苗', date: '2026-03-01', nextDate: '2027-03-01', hospital: '宠物医院A' },
  { id: 2, petName: '小黄', vaccine: '犬瘟热疫苗', date: '2026-02-15', nextDate: '2027-02-15', hospital: '宠物医院A' },
  { id: 3, petName: '小红', vaccine: '猫三联', date: '2026-01-10', nextDate: '2027-01-10', hospital: '宠物医院B' },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '宠物', dataIndex: 'petName', width: 100 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '医院', dataIndex: 'hospital', width: 150 },
  { title: '日期', dataIndex: 'date', width: 120 },
  { title: '结果', dataIndex: 'result', width: 150 },
  { title: '医生', dataIndex: 'doctor', width: 100 },
];

const vaccineColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '宠物', dataIndex: 'petName', width: 100 },
  { title: '疫苗', dataIndex: 'vaccine', width: 150 },
  { title: '接种日期', dataIndex: 'date', width: 120 },
  { title: '下次接种', dataIndex: 'nextDate', width: 120 },
  { title: '医院', dataIndex: 'hospital', width: 150 },
];
</script>

<style scoped>
.medical-records-container { padding: 20px; }
.hospital-card { text-align: center; margin-bottom: 12px; }
.hospital-name { font-weight: bold; margin-bottom: 4px; }
.hospital-type { color: #86909c; font-size: 12px; margin-bottom: 8px; }
</style>
