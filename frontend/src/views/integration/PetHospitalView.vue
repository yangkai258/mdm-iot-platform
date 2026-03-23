<template>
  <div class="pet-hospital-view">
    <a-card title="宠物医疗">
      <a-alert type="info" style="margin-bottom: 16px">
        第三方宠物医疗服务集成，可在紧急情况下自动联系附近宠物医院。
      </a-alert>

      <a-row :gutter="16" style="margin-bottom: 24px">
        <a-col :span="8">
          <a-statistic title="绑定医院" :value="stats.hospital_count" />
        </a-col>
        <a-col :span="8">
          <a-statistic title="紧急联系人" :value="stats.emergency_contacts" />
        </a-col>
        <a-col :span="8">
          <a-statistic title="本月问诊" :value="stats.consultations" />
        </a-col>
      </a-row>

      <a-tabs>
        <a-tab-pane key="hospitals" title="合作医院">
          <a-space style="margin-bottom: 12px">
            <a-button type="primary" @click="showAddHospital = true">添加医院</a-button>
            <a-input-search v-model="searchKeyword" placeholder="搜索医院名称" style="width: 240px" @search="loadHospitals" />
          </a-space>
          <a-table :columns="hospitalColumns" :data="hospitals" :loading="loading" :pagination="pagination" @page-change="onPageChange">
            <template #status="{ record }">
              <a-tag :color="record.status === 'active' ? 'green' : 'gray'">{{ record.status === 'active' ? '已绑定' : '未绑定' }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="viewHospital(record)">查看</a-button>
                <a-button type="text" size="small" status="danger" @click="unbindHospital(record)">解绑</a-button>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
        <a-tab-pane key="records" title="问诊记录">
          <a-table :columns="recordColumns" :data="records" :loading="loading" />
        </a-tab-pane>
        <a-tab-pane key="settings" title="急救设置">
          <a-form layout="vertical">
            <a-form-item label="自动拨打急救电话">
              <a-switch v-model="settings.auto_emergency_call" />
            </a-form-item>
            <a-form-item label="紧急联系人手机">
              <a-input v-model="settings.emergency_phone" placeholder="请输入手机号" />
            </a-form-item>
            <a-form-item label="健康异常自动通知">
              <a-switch v-model="settings.auto_notify" />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="saveSettings">保存设置</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const showAddHospital = ref(false)
const searchKeyword = ref('')
const hospitals = ref([])
const records = ref([])

const stats = reactive({ hospital_count: 3, emergency_contacts: 2, consultations: 7 })

const settings = reactive({
  auto_emergency_call: true,
  emergency_phone: '',
  auto_notify: true
})

const hospitalColumns = [
  { title: '医院名称', dataIndex: 'name' },
  { title: '地址', dataIndex: 'address' },
  { title: '电话', dataIndex: 'phone' },
  { title: '距离', dataIndex: 'distance', suffix: 'km' },
  { title: '状态', slotName: 'status' },
  { title: '操作', slotName: 'actions' }
]

const recordColumns = [
  { title: '日期', dataIndex: 'date' },
  { title: '宠物', dataIndex: 'pet_name' },
  { title: '医院', dataIndex: 'hospital' },
  { title: '诊断', dataIndex: 'diagnosis' },
  { title: '状态', dataIndex: 'status' }
]

const pagination = reactive({ total: 0, current: 1, pageSize: 10 })

async function loadHospitals() {
  loading.value = true
  try {
    hospitals.value = [
      { id: 1, name: '阳光宠物医院', address: '朝阳区建国路88号', phone: '010-12345678', distance: 2.3, status: 'active' },
      { id: 2, name: '爱康宠物诊所', address: '海淀区中关村大街1号', phone: '010-87654321', distance: 5.1, status: 'active' },
      { id: 3, name: '宠物急救中心', address: '东城区东单北大街3号', phone: '010-11223344', distance: 8.7, status: 'inactive' }
    ]
    pagination.total = 3
  } finally {
    loading.value = false
  }
}

async function loadRecords() {
  records.value = [
    { date: '2026-03-20', pet_name: '小白', hospital: '阳光宠物医院', diagnosis: '皮肤炎症', status: '已完成' },
    { date: '2026-03-15', pet_name: '小白', hospital: '爱康宠物诊所', diagnosis: '常规体检', status: '已完成' }
  ]
}

function onPageChange(page) {
  pagination.current = page
  loadHospitals()
}

function viewHospital(record) {
  Message.info('查看医院: ' + record.name)
}

function unbindHospital(record) {
  Message.warning('解绑医院: ' + record.name)
}

async function saveSettings() {
  Message.success('设置已保存')
}

onMounted(() => {
  loadHospitals()
  loadRecords()
})
</script>
