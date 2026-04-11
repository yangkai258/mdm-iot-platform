<template>
  <div class="page-container">
    <a-card class="general-card" title="瀹犵墿鐧昏">
      <template #extra>
        <a-button @click="loadData"><icon-refresh />鍒锋柊</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="瀹犵墿鍚嶇О"><a-input v-model="form.pet_name" placeholder="璇疯緭锟" /></a-form-item>
          <a-form-item label="瀹犵墿绫诲瀷">
            <a-select v-model="form.pet_type" placeholder="閫夋嫨绫诲瀷" allow-clear style="width: 120px">
              <a-option value="dog">锟?/a-option>
              <a-option value="cat">锟?/a-option>
              <a-option value="bird">锟?/a-option>
              <a-option value="other">鍏朵粬</a-option>
            </a-select>
          </a-form-item>
          <a-form-item><a-button type="primary" @click="handleSearch">鏌ヨ</a-button><a-button @click="Object.keys(form).forEach(k => form[k] = ''); loadData()">閲嶇疆</a-button></a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #pet_photo="{ record }">
          <a-avatar v-if="record.pet_photo" :src="record.pet_photo" size="large" />
          <a-avatar v-else size="large">{{ record.pet_name?.[0] || '?' }}</a-avatar>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 'active' ? 'green' : 'gray'">{{ record.status === 'active' ? '宸叉縺锟? : '鏈縺锟? }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="viewDetail(record)">璇︽儏</a-button>
          <a-button type="text" size="small" @click="editPet(record)">缂栬緫</a-button>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model:visible="formVisible" :title="isEdit ? '缂栬緫瀹犵墿' : '瀹犵墿鐧昏'" @before-ok="handleSubmit" :loading="submitting" :width="640">
      <a-form :model="petForm" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="瀹犵墿鍚嶇О" required><a-input v-model="petForm.pet_name" placeholder="璇疯緭鍏ュ疇鐗╁悕锟" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="瀹犵墿绫诲瀷" required>
              <a-select v-model="petForm.pet_type" placeholder="閫夋嫨绫诲瀷">
                <a-option value="dog">锟?/a-option>
                <a-option value="cat">锟?/a-option>
                <a-option value="bird">锟?/a-option>
                <a-option value="other">鍏朵粬</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="鍝佺"><a-input v-model="petForm.breed" placeholder="璇疯緭鍏ュ搧锟" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="鎬у埆">
              <a-select v-model="petForm.gender" placeholder="閫夋嫨鎬у埆">
                <a-option value="male">锟?/a-option>
                <a-option value="female">锟?/a-option>
                <a-option value="unknown">鏈煡</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="鍑虹敓鏃ユ湡"><a-date-picker v-model="petForm.birth_date" style="width: 100%" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="浣撻噸(kg)"><a-input-number v-model="petForm.weight" :min="0" :max="500" style="width: 100%" /></a-form-item>
          </a-col>
        </a-row>
        <a-divider>涓讳汉淇℃伅</a-divider>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="涓讳汉濮撳悕" required><a-input v-model="petForm.owner_name" placeholder="璇疯緭鍏ヤ富浜哄锟" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="鑱旂郴鐢佃瘽" required><a-input v-model="petForm.owner_phone" placeholder="璇疯緭鍏ヨ仈绯荤數锟" /></a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="鑱旂郴鍦板潃"><a-input v-model="petForm.owner_address" placeholder="璇疯緭鍏ヨ仈绯诲湴鍧€" /></a-form-item>
        <a-form-item label="瀹犵墿鐓х墖">
          <a-upload :custom-request="handleUpload" :show-file-list="false">
            <a-button v-if="!petForm.pet_photo"><icon-upload />涓婁紶鐓х墖</a-button>
            <a-image v-else :src="petForm.pet_photo" width="120" height="120" :preview="true" />
          </a-upload>
        </a-form-item>
        <a-form-item label="澶囨敞"><a-textarea v-model="petForm.notes" :rows="2" placeholder="澶囨敞淇℃伅" /></a-form-item>
      </a-form>
    </a-modal>
    <a-modal v-model:visible="detailVisible" title="瀹犵墿璇︽儏" :footer="null" :width="560">
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="瀹犵墿鍚嶇О" :span="2">{{ currentPet?.pet_name }}</a-descriptions-item>
        <a-descriptions-item label="瀹犵墿绫诲瀷">{{ currentPet?.pet_type }}</a-descriptions-item>
        <a-descriptions-item label="鍝佺">{{ currentPet?.breed }}</a-descriptions-item>
        <a-descriptions-item label="鎬у埆">{{ currentPet?.gender }}</a-descriptions-item>
        <a-descriptions-item label="涓讳汉">{{ currentPet?.owner_name }}</a-descriptions-item>
        <a-descriptions-item label="鐢佃瘽">{{ currentPet?.owner_phone }}</a-descriptions-item>
        <a-descriptions-item label="鍦板潃" :span="2">{{ currentPet?.owner_address }}</a-descriptions-item>
        <a-descriptions-item label="鐧昏鏃堕棿" :span="2">{{ currentPet?.created_at }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const submitting = ref(false)
const data = ref([])
const formVisible = ref(false)
const detailVisible = ref(false)
const isEdit = ref(false)
const currentPet = ref(null)
const form = reactive({ pet_name: '', pet_type: '' })
const petForm = reactive({ id: null, pet_name: '', pet_type: '', breed: '', gender: '', birth_date: null, weight: null, owner_name: '', owner_phone: '', owner_address: '', pet_photo: '', notes: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: '鐓х墖', slotName: 'pet_photo', width: 80 },
  { title: '瀹犵墿鍚嶇О', dataIndex: 'pet_name', width: 120 },
  { title: '绫诲瀷', dataIndex: 'pet_type', width: 80 },
  { title: '鍝佺', dataIndex: 'breed', width: 100 },
  { title: '涓讳汉', dataIndex: 'owner_name', width: 100 },
  { title: '鐢佃瘽', dataIndex: 'owner_phone', width: 130 },
  { title: '鐘讹拷?, slotName: 'status', width: 80 },
  { title: '鐧昏鏃堕棿', dataIndex: 'created_at', width: 170 },
  { title: '鎿嶄綔', slotName: 'actions', width: 120 }
]

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (form.pet_name) params.append('pet_name', form.pet_name)
    if (form.pet_type) params.append('pet_type', form.pet_type)
    const res = await fetch(`/api/v1/pet/register?${params}`, { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { data.value = res.data?.list || []; pagination.total = res.data?.total || 0 }
    else { data.value = [] }
  } catch (e) { Message.error('鍔犺浇澶辫触') } finally { loading.value = false }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }

const viewDetail = (record) => { currentPet.value = record; detailVisible.value = true }

const editPet = (record) => { isEdit.value = true; Object.assign(petForm, record); formVisible.value = true }

const handleSubmit = async (done) => {
  if (!petForm.pet_name || !petForm.owner_name || !petForm.owner_phone) { Message.warning('璇峰～鍐欏繀濉」'); done(false); return }
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const url = isEdit.value ? `/api/v1/pet/register/${petForm.id}` : '/api/v1/pet/register'
    const res = await fetch(url, { method: isEdit.value ? 'PUT' : 'POST', headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(petForm) }).then(r => r.json())
    if (res.code === 0) { Message.success(isEdit.value ? '鏇存柊鎴愬姛' : '鐧昏鎴愬姛'); formVisible.value = false; loadData() }
    else { Message.error(res.message || '鎿嶄綔澶辫触') }
    done(true)
  } catch (e) { Message.error('鎿嶄綔澶辫触'); done(false) } finally { submitting.value = false }
}

const handleUpload = async ({ file, onSuccess, onError }) => {
  try {
    const formData = new FormData()
    formData.append('file', file)
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/upload', { method: 'POST', headers: { 'Authorization': `Bearer ${token}` }, body: formData }).then(r => r.json())
    if (res.code === 0) { petForm.pet_photo = res.data?.url; onSuccess() }
    else onError()
  } catch (e) { onError() }
}

onMounted(() => loadData())
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>