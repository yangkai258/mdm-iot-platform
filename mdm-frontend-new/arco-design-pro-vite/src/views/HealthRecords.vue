<template>
  <div class="health-records-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="健康记录" :value="856" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="疫苗记录" :value="128" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="体检报告" :value="56" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="异常记录" :value="12" status="warning" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>宠物健康档案</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            添加记录
          </a-button>
        </div>
      </template>
      
      <a-tabs>
        <a-tab-pane key="records" title="健康记录">
          <a-table :columns="recordColumns" :data="records" :pagination="pagination">
            <template #type="{ record }">
              <a-tag>{{ record.typeText }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="vaccines" title="疫苗记录">
          <a-table :columns="vaccineColumns" :data="vaccines" :pagination="pagination" />
        </a-tab-pane>
        
        <a-tab-pane key="exams" title="体检报告">
          <a-table :columns="examColumns" :data="exams" :pagination="pagination" />
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <a-modal v-model:visible="createVisible" title="添加健康记录" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="宠物" required>
          <a-select v-model="form.petId" placeholder="选择宠物">
            <a-option value="P001">小黄</a-option>
            <a-option value="P002">小红</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="记录类型">
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="checkup">体检</a-option>
            <a-option value="vaccine">疫苗</a-option>
            <a-option value="disease">疾病</a-option>
            <a-option value="surgery">手术</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="记录日期">
          <a-date-picker v-model="form.recordDate" style="width: 100%;" />
        </a-form-item>
        <a-form-item label="医院/诊所">
          <a-input v-model="form.hospital" placeholder="请输入医院名称" />
        </a-form-item>
        <a-form-item label="诊断结果">
          <a-textarea v-model="form.diagnosis" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 6 });
const createVisible = ref(false);

const form = reactive({ petId: '', type: '', recordDate: '', hospital: '', diagnosis: '' });

const records = ref([
  { id: 1, petName: '小黄', type: 'checkup', typeText: '体检', hospital: '宠物医院A', diagnosis: '身体健康', recordDate: '2026-03-20', doctor: '张医生' },
  { id: 2, petName: '小红', type: 'vaccine', typeText: '疫苗', hospital: '宠物医院B', diagnosis: '狂犬疫苗已接种', recordDate: '2026-03-15', doctor: '李医生' },
]);

const vaccines = ref([
  { id: 1, petName: '小黄', vaccineName: '狂犬疫苗', batchNo: 'AB123456', hospital: '宠物医院A', date: '2026-03-20', nextDate: '2027-03-20' },
]);

const exams = ref([
  { id: 1, petName: '小黄', examType: '年度体检', hospital: '宠物医院A', date: '2026-03-20', result: '正常', doctor: '张医生' },
]);

const recordColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '宠物', dataIndex: 'petName', width: 100 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '医院', dataIndex: 'hospital', width: 150 },
  { title: '诊断', dataIndex: 'diagnosis' },
  { title: '日期', dataIndex: 'recordDate', width: 120 },
];

const vaccineColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '宠物', dataIndex: 'petName', width: 100 },
  { title: '疫苗名称', dataIndex: 'vaccineName', width: 150 },
  { title: '批号', dataIndex: 'batchNo', width: 120 },
  { title: '接种日期', dataIndex: 'date', width: 120 },
  { title: '下次日期', dataIndex: 'nextDate', width: 120 },
];

const examColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '宠物', dataIndex: 'petName', width: 100 },
  { title: '体检类型', dataIndex: 'examType', width: 120 },
  { title: '医院', dataIndex: 'hospital', width: 150 },
  { title: '结果', dataIndex: 'result', width: 100 },
  { title: '日期', dataIndex: 'date', width: 120 },
];

const handleCreate = () => { createVisible.value = true; };
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.health-records-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
