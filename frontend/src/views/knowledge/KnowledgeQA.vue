<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>棣栭〉</a-breadcrumb-item>
      <a-breadcrumb-item>鐭ヨ瘑搴?/a-breadcrumb-item>
      <a-breadcrumb-item>闂瓟瀵圭鐞?/a-breadcrumb-item>
    </a-breadcrumb>

    <div class="pro-search-bar">
      <a-space>
        <a-input-search v-model="searchKeyword" placeholder="鎼滅储闂/绛旀" style="width: 280px" @search="loadQA" search-button />
        <a-select v-model="filterTag" placeholder="鏍囩" allow-clear style="width: 160px" @change="loadQA">
          <a-option v-for="tag in tagOptions" :key="tag" :value="tag">{{ tag }}</a-option>
        </a-select>
      </a-space>
    </div>

    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showCreateModal">鏂板缓闂瓟</a-button>
        <a-button @click="handleExport">瀵煎嚭</a-button>
        <a-button @click="showImportModal">瀵煎叆</a-button>
        <a-button @click="loadQA">鍒锋柊</a-button>
      </a-space>
    </div>

    <div class="pro-content-area">
      <a-table :columns="columns" :data="qaList" :loading="loading" :pagination="pagination" row-key="id" @page-change="handlePageChange">
        <template #tags="{ record }">
          <a-tag v-for="tag in record.tags" :key="tag" style="margin-right: 4px">{{ tag }}</a-tag>
        </template>
        <template #enabled="{ record }">
          <a-switch v-model="record.enabled" @change="toggleQA(record)" />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="editQA(record)">缂栬緫</a-button>
            <a-button type="text" size="small" @click="previewQA(record)">棰勮</a-button>
            <a-button type="text" size="small" status="danger" @click="deleteQA(record)">鍒犻櫎</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 鏂板缓/缂栬緫寮圭獥 -->
    <a-modal v-model:visible="modalVisible" :title="isEdit ? '缂栬緫闂瓟' : '鏂板缓闂瓟'" @ok="submitQA" :width="680" :loading="submitting">
      <a-form :model="qaForm" layout="vertical">
        <a-form-item label="闂" required>
          <a-textarea v-model="qaForm.question" placeholder="杈撳叆鐢ㄦ埛鍙兘闂殑闂" :rows="2" />
        </a-form-item>
        <a-form-item label="绛旀" required>
          <a-textarea v-model="qaForm.answer" placeholder="杈撳叆鏍囧噯绛旀" :rows="4" />
        </a-form-item>
        <a-form-item label="鏍囩">
          <a-select v-model="qaForm.tags" multiple placeholder="閫夋嫨鎴栬緭鍏ユ爣绛" allow-create :style="{ width: '100%' }">
            <a-option v-for="tag in tagOptions" :key="tag" :value="tag">{{ tag }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="鐩镐技闂锛堝彲閫夛級">
          <a-textarea v-model="qaForm.alternatives" placeholder="姣忚涓€鏉＄浉浼奸棶棰" :rows="3" />
        </a-form-item>
        <a-form-item label="鍚敤">
          <a-switch v-model="qaForm.enabled" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 棰勮寮圭獥 -->
    <a-modal v-model:visible="previewVisible" title="闂瓟棰勮" :width="560" :footer="null">
      <a-result v-if="!previewRecord" status="info" title="璇烽€夋嫨涓€鏉¤褰? />
      <template v-else>
        <a-alert type="info" style="margin-bottom: 12px">
          <template #title>鐢ㄦ埛闂?/template>
          <div>{{ previewRecord.question }}</div>
        </a-alert>
        <a-alert type="success">
          <template #title>绯荤粺绛?/template>
          <div>{{ previewRecord.answer }}</div>
        </a-alert>
        <a-divider>鐩镐技闂</a-divider>
        <a-tag v-for="alt in previewRecord.alternatives" :key="alt" style="margin: 4px">{{ alt }}</a-tag>
      </template>
    </a-modal>

    <!-- 瀵煎叆寮圭獥 -->
    <a-modal v-model:visible="importVisible" title="瀵煎叆闂瓟" @ok="submitImport" :width="480" :loading="importing">
      <a-form-item label="瀵煎叆鏂瑰紡">
        <a-radio-group v-model="importMode">
          <a-radio value="merge">鍚堝苟锛堣拷鍔狅級</a-radio>
          <a-radio value="replace">瑕嗙洊锛堟竻绌哄悗瀵煎叆锛?/a-radio>
        </a-radio-group>
      </a-form-item>
      <a-form-item label="閫夋嫨鏂囦欢">
        <a-upload :limit="1" accept=".json,.csv,.xlsx" :custom-request="handleFileChange" />
      </a-form-item>
      <a-alert type="info" message="鏀寔 JSON銆丆SV銆乆LSX 鏍煎紡锛屾瘡鏉¤褰曢渶鍖呭惈 question銆乤nswer 瀛楁" />
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'

const loading = ref(false)
const modalVisible = ref(false)
const previewVisible = ref(false)
const importVisible = ref(false)
const submitting = ref(false)
const importing = ref(false)
const isEdit = ref(false)
const searchKeyword = ref('')
const filterTag = ref('')
const previewRecord = ref<any>(null)
const importMode = ref('merge')
const importFile = ref<any>(null)

const tagOptions = ['瀹犵墿鍠傞', '瀹犵墿鍋ュ悍', '璁惧浣跨敤', '璐︽埛闂', '鍥轰欢鍗囩骇', '甯歌闂']
const qaList = ref<any[]>([])
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })

const qaForm = reactive({
  id: 0, question: '', answer: '', tags: [] as string[], alternatives: '', enabled: true,
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '闂', dataIndex: 'question', ellipsis: true },
  { title: '绛旀棰勮', dataIndex: 'answer', ellipsis: true },
  { title: '鏍囩', slotName: 'tags', width: 200 },
  { title: '鍚敤', slotName: 'enabled', width: 80 },
  { title: '鎿嶄綔', slotName: 'actions', fixed: 'right', width: 200 },
]

const loadQA = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/knowledge/qa', {
      params: { page: pagination.current, page_size: pagination.pageSize, keyword: searchKeyword.value, tag: filterTag.value },
    })
    qaList.value = res.data.items || []
    pagination.total = res.data.total || 0
  } catch {
    qaList.value = [
      { id: 1, question: '濡備綍缁欏疇鐗╂坊鍔犳柊璁惧锛?, answer: '鎵撳紑璁惧绠＄悊椤甸潰锛岀偣鍑绘坊鍔犺澶囨寜閽?..', tags: ['璁惧浣跨敤'], alternatives: ['鎬庝箞缁戝畾璁惧', '璁惧鎬庝箞杩炴帴'], enabled: true },
      { id: 2, question: '鍥轰欢鍗囩骇澶辫触鎬庝箞鍔烇紵', answer: '璇锋鏌ョ綉缁滆繛鎺ワ紝纭繚璁惧鍦ㄧ嚎...', tags: ['鍥轰欢鍗囩骇', '甯歌闂'], alternatives: ['OTA鍗囩骇涓嶄簡'], enabled: true },
      { id: 3, question: '瀹犵墿鍋ュ悍鏁版嵁鍦ㄥ摢鐪嬶紵', answer: '鍦ㄥ仴搴峰尰鐤楄彍鍗曚腑鏌ョ湅鍋ュ悍鎶ュ憡...', tags: ['瀹犵墿鍋ュ悍'], alternatives: ['鍋ュ悍鎶ュ憡鎬庝箞鏌?], enabled: false },
    ]
    pagination.total = 3
  } finally { loading.value = false }
}

const showCreateModal = () => {
  isEdit.value = false
  Object.assign(qaForm, { id: 0, question: '', answer: '', tags: [], alternatives: '', enabled: true })
  modalVisible.value = true
}

const editQA = (record: any) => {
  isEdit.value = true
  Object.assign(qaForm, { ...record, alternatives: (record.alternatives || []).join('\n') })
  modalVisible.value = true
}

const previewQA = (record: any) => {
  previewRecord.value = record
  previewVisible.value = true
}

const submitQA = async () => {
  submitting.value = true
  try {
    const payload = { ...qaForm, alternatives: qaForm.alternatives.split('\n').filter(Boolean) }
    if (isEdit.value) {
      await axios.put(`/api/v1/knowledge/qa/${qaForm.id}`, payload)
      Message.success('鏇存柊鎴愬姛')
    } else {
      await axios.post('/api/v1/knowledge/qa', payload)
      Message.success('鍒涘缓鎴愬姛')
    }
    modalVisible.value = false
    loadQA()
  } catch { Message.error('鎿嶄綔澶辫触') } finally { submitting.value = false }
}

const toggleQA = async (record: any) => {
  try { await axios.put(`/api/v1/knowledge/qa/${record.id}`, { enabled: record.enabled }); Message.success('鏇存柊鎴愬姛') }
  catch { record.enabled = !record.enabled; Message.error('鏇存柊澶辫触') }
}

const deleteQA = async (record: any) => {
  try { await axios.delete(`/api/v1/knowledge/qa/${record.id}`); Message.success('鍒犻櫎鎴愬姛'); loadQA() }
  catch { Message.error('鍒犻櫎澶辫触') }
}

const handleExport = async () => {
  try {
    const res = await axios.get('/api/v1/knowledge/qa/export', { params: { keyword: searchKeyword.value } })
    const blob = new Blob([JSON.stringify(res.data, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a'); a.href = url; a.download = 'knowledge-qa.json'; a.click()
    URL.revokeObjectURL(url)
    Message.success('瀵煎嚭鎴愬姛')
  } catch { Message.error('瀵煎嚭澶辫触') }
}

const showImportModal = () => { importVisible.value = true }
const handleFileChange = (options: any) => { importFile.value = options.file; options.onSuccess() }

const submitImport = async () => {
  if (!importFile.value) { Message.warning('璇烽€夋嫨鏂囦欢'); return }
  importing.value = true
  try {
    const formData = new FormData()
    formData.append('file', importFile.value)
    formData.append('mode', importMode.value)
    await axios.post('/api/v1/knowledge/qa/import', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
    Message.success('瀵煎叆鎴愬姛')
    importVisible.value = false
    loadQA()
  } catch { Message.error('瀵煎叆澶辫触') } finally { importing.value = false }
}

const handlePageChange = (page: number) => { pagination.current = page; loadQA() }

onMounted(() => loadQA())
</script>
