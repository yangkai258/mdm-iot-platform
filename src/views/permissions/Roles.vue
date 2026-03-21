<template>
  <div class="page-container">
    <!-- 闈㈠寘灞?+ 鎼滅储鏍?-->
    <div class="page-header">
      <div class="header-left">
        <a-breadcrumb>
          <a-breadcrumb-item><a href="#/dashboard">棣栭〉</a></a-breadcrumb-item>
          <a-breadcrumb-item>鏉冮檺绠＄悊</a-breadcrumb-item>
          <a-breadcrumb-item>瑙掕壊绠＄悊</a-breadcrumb-item>
        </a-breadcrumb>
      </div>
      <div class="header-right">
        <a-input-search
          v-model="searchForm.keyword"
          placeholder="瑙掕壊鍚嶇О/缂栫爜"
          style="width: 240px"
          @search="doSearch"
        />
      </div>
    </div>

    <!-- 鎿嶄綔鎸夐挳鏍忥紙闈犲乏锛?-->
    <div class="action-bar">
      <a-space :size="12">
        <a-button type="primary" @click="openCreateModal">銆屾柊寤恒€?/a-button>
        <a-button @click="toggleSearch">銆岀瓫閫夈€?/a-button>
        <a-button @click="loadRoles">銆屽埛鏂般€?/a-button>
      </a-space>
    </div>

    <!-- 绛涢€夐潰鏉?-->
    <a-card v-if="showSearch" :bordered="false" style="margin-bottom: 12px">
      <a-form :model="searchForm" layout="inline">
        <a-form-item field="keyword" label="瑙掕壊鍚嶇О/缂栫爜">
          <a-input v-model="searchForm.keyword" placeholder="璇疯緭鍏ュ叧閿瘝" allow-clear style="width: 200px" />
        </a-form-item>
        <a-form-item field="status" label="鐘舵€?>
          <a-select v-model="searchForm.status" placeholder="閫夋嫨鐘舵€? allow-clear style="width: 120px">
            <a-option :value="1">鍚敤</a-option>
            <a-option :value="0">绂佺敤</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="doSearch">鏌ヨ</a-button>
          <a-button style="margin-left: 8px" @click="resetSearch">閲嶇疆</a-button>
        </a-form-item>
      </a-form>
    </a-card>

      <!-- 瑙掕壊鍒楄〃 -->
      <a-table
        :columns="columns"
        :data="roles"
        :loading="loading"
        :pagination="pagination"
        @change="handleTableChange"
        row-key="id"
        style="margin-top: 16px"
      >
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'gray'">
            {{ record.status === 1 ? '鍚敤' : '绂佺敤' }}
          </a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openEditModal(record)">缂栬緫</a-button>
            <a-button type="text" size="small" @click="openPermModal(record)">鏉冮檺</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">鍒犻櫎</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 鍒涘缓/缂栬緫瑙掕壊鍏ㄥ睆妯℃€?-->
    <a-modal
      v-model:visible="formVisible"
      :title="isEdit ? '缂栬緫瑙掕壊' : '鍒涘缓瑙掕壊'"
      :width="600"
      :mask-closable="false"
      @ok="submitForm"
      @cancel="formVisible = false"
    >
      <a-form :model="formData" layout="vertical">
        <a-form-item field="name" label="瑙掕壊鍚嶇О" required>
          <a-input v-model="formData.name" placeholder="璇疯緭鍏ヨ鑹插悕绉? />
        </a-form-item>
        <a-form-item field="code" label="瑙掕壊缂栫爜" :required="!isEdit">
          <a-input v-model="formData.code" placeholder="璇疯緭鍏ヨ鑹茬紪鐮侊紝濡?admin" :disabled="isEdit" />
        </a-form-item>
        <a-form-item field="description" label="鎻忚堪">
          <a-textarea v-model="formData.description" placeholder="璇疯緭鍏ヨ鑹叉弿杩? :rows="3" />
        </a-form-item>
        <a-form-item field="sort" label="鎺掑簭">
          <a-input-number v-model="formData.sort" :min="0" :max="9999" />
        </a-form-item>
        <a-form-item field="status" label="鐘舵€?>
          <a-switch v-model="formData.status" :checked-value="1" :unchecked-value="0" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 鏉冮檺鍒嗛厤寮圭獥 -->
    <a-modal
      v-model:visible="permVisible"
      title="鍒嗛厤鏉冮檺"
      :width="700"
      :mask-closable="false"
      @ok="submitPerms"
      @cancel="permVisible = false"
    >
      <a-alert type="info" style="margin-bottom: 16px">
        褰撳墠瑙掕壊: <b>{{ currentRole?.name }}</b>
      </a-alert>
      <a-tabs default-active-tab="tree">
        <a-tab-pane key="tree" title="鏉冮檺鏍?>
          <a-tree
            v-model:selected-keys="selectedPerms"
            v-model:checked-keys="selectedPerms"
            :data="permTreeData"
            :selectable="false"
            checkable
            :default-expand-all="true"
            field-names="{ key: 'id', title: 'name', children: 'children' }"
          />
        </a-tab-pane>
        <a-tab-pane key="list" title="鏉冮檺鍒楄〃">
          <a-checkbox-group v-model="selectedPerms">
            <a-space direction="vertical" fill>
              <div v-for="group in permGroups" :key="group.id">
                <div class="perm-group-title">{{ group.name }}</div>
                <a-space wrap>
                  <a-checkbox v-for="p in group.permissions" :key="p" :value="p">{{ p }}</a-checkbox>
                </a-space>
              </div>
            </a-space>
          </a-checkbox-group>
        </a-tab-pane>
      </a-tabs>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'

const loading = ref(false)
const showSearch = ref(false)
const formVisible = ref(false)
const permVisible = ref(false)
const isEdit = ref(false)
const roles = ref([])
const permTreeData = ref([])
const selectedPerms = ref([])
const currentRole = ref(null)

const searchForm = reactive({
  keyword: '',
  status: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const formData = reactive({
  id: null,
  name: '',
  code: '',
  description: '',
  sort: 0,
  status: 1
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '瑙掕壊鍚嶇О', dataIndex: 'name', width: 150 },
  { title: '瑙掕壊缂栫爜', dataIndex: 'code', width: 150 },
  { title: '鎻忚堪', dataIndex: 'description', ellipsis: true },
  { title: '鎺掑簭', dataIndex: 'sort', width: 80 },
  { title: '鐘舵€?, slotName: 'status', width: 100 },
  { title: '鍒涘缓鏃堕棿', dataIndex: 'created_at', width: 180 },
  { title: '鎿嶄綔', slotName: 'actions', width: 200, fixed: 'right' }
]

const permGroups = ref([
  { id: 1, name: '绉熸埛绠＄悊', permissions: ['tenant:view', 'tenant:manage'] },
  { id: 2, name: '鐢ㄦ埛绠＄悊', permissions: ['user:view', 'user:manage'] },
  { id: 3, name: '璁惧绠＄悊', permissions: ['device:view', 'device:manage', 'device:control'] },
  { id: 4, name: 'OTA绠＄悊', permissions: ['ota:view', 'ota:deploy'] },
  { id: 5, name: '鍛婅绠＄悊', permissions: ['alert:view', 'alert:manage'] },
  { id: 6, name: '浼氬憳绠＄悊', permissions: ['member:view', 'member:manage'] },
  { id: 7, name: '绛栫暐绠＄悊', permissions: ['policy:view', 'policy:manage'] },
  { id: 8, name: '閫氱煡绠＄悊', permissions: ['notification:view', 'notification:manage'] },
  { id: 9, name: '搴旂敤绠＄悊', permissions: ['app:view', 'app:manage'] },
  { id: 10, name: '绯荤粺绠＄悊', permissions: ['system:view', 'system:manage', 'role:manage'] },
  { id: 11, name: '鏁版嵁鎿嶄綔', permissions: ['data:export', 'data:import'] },
])

const apiBase = '/api/v1'

const getToken = () => localStorage.getItem('token')

const fetchApi = async (url, options = {}) => {
  const res = await fetch(`${apiBase}${url}`, {
    ...options,
    headers: {
      'Authorization': `Bearer ${getToken()}`,
      'Content-Type': 'application/json',
      ...options.headers
    }
  })
  return res.json()
}

const loadRoles = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams({
      page: pagination.current,
      page_size: pagination.pageSize
    })
    if (searchForm.keyword) params.set('keyword', searchForm.keyword)
    if (searchForm.status !== '') params.set('status', searchForm.status)

    const data = await fetchApi(`/roles?${params}`)
    if (data.code === 0) {
      roles.value = data.data?.list || []
      pagination.total = data.data?.total || 0
    }
  } catch (e) {
    Message.error('鍔犺浇瑙掕壊鍒楄〃澶辫触')
  } finally {
    loading.value = false
  }
}

const loadPermTree = async () => {
  try {
    const data = await fetchApi('/permissions')
    if (data.code === 0) {
      permTreeData.value = data.data || []
    }
  } catch (e) {
    console.error('鍔犺浇鏉冮檺鏍戝け璐?, e)
  }
}

const toggleSearch = () => {
  showSearch.value = !showSearch.value
}

const doSearch = () => {
  pagination.current = 1
  loadRoles()
}

const resetSearch = () => {
  searchForm.keyword = ''
  searchForm.status = ''
  pagination.current = 1
  loadRoles()
}

const openCreateModal = () => {
  isEdit.value = false
  formData.id = null
  formData.name = ''
  formData.code = ''
  formData.description = ''
  formData.sort = 0
  formData.status = 1
  formVisible.value = true
}

const openEditModal = (record) => {
  isEdit.value = true
  Object.assign(formData, {
    id: record.id,
    name: record.name,
    code: record.code,
    description: record.description,
    sort: record.sort,
    status: record.status
  })
  formVisible.value = true
}

const submitForm = async () => {
  if (!formData.name) {
    Message.warning('璇疯緭鍏ヨ鑹插悕绉?)
    return
  }
  if (!isEdit.value && !formData.code) {
    Message.warning('璇疯緭鍏ヨ鑹茬紪鐮?)
    return
  }

  try {
    const url = isEdit.value ? `/roles/${formData.id}` : '/roles'
    const method = isEdit.value ? 'PUT' : 'POST'
    const data = await fetchApi(url, {
      method,
      body: JSON.stringify(formData)
    })
    if (data.code === 0) {
      Message.success(isEdit.value ? '鏇存柊鎴愬姛' : '鍒涘缓鎴愬姛')
      formVisible.value = false
      loadRoles()
    } else {
      Message.error(data.message || '鎿嶄綔澶辫触')
    }
  } catch (e) {
    Message.error('鎿嶄綔澶辫触')
  }
}

const openPermModal = async (record) => {
  currentRole.value = record
  permVisible.value = true
  selectedPerms.value = []

  // 鍔犺浇瑙掕壊宸叉湁鏉冮檺
  try {
    const data = await fetchApi(`/roles/${record.id}/permissions`)
    if (data.code === 0) {
      selectedPerms.value = data.data?.permission_ids || []
    }
  } catch (e) {
    console.error('鍔犺浇瑙掕壊鏉冮檺澶辫触', e)
  }

  // 纭繚鏉冮檺鏍戝凡鍔犺浇
  if (permTreeData.value.length === 0) {
    await loadPermTree()
  }
}

const submitPerms = async () => {
  if (!currentRole.value) return

  try {
    // 灏嗛€変腑鏉冮檺鐨刢ode杞负ID锛堣繖閲岀畝鍖栧鐞嗭紝鍓嶇浼燾ode鍒楄〃锛屽悗绔鐞嗭級
    const permIds = selectedPerms.value.map(p => typeof p === 'number' ? p : p)
    const data = await fetchApi(`/roles/${currentRole.value.id}/permissions`, {
      method: 'POST',
      body: JSON.stringify({ permission_ids: permIds })
    })
    if (data.code === 0) {
      Message.success('鏉冮檺鍒嗛厤鎴愬姛')
      permVisible.value = false
    } else {
      Message.error(data.message || '鍒嗛厤澶辫触')
    }
  } catch (e) {
    Message.error('鏉冮檺鍒嗛厤澶辫触')
  }
}

const handleDelete = (record) => {
  Modal.warning({
    title: '纭鍒犻櫎',
    content: `纭畾瑕佸垹闄よ鑹层€?{record.name}銆嶅悧锛熷垹闄ゅ悗涓嶅彲鎭㈠銆俙,
    okText: '鍒犻櫎',
    onOk: async () => {
      try {
        const data = await fetchApi(`/roles/${record.id}`, { method: 'DELETE' })
        if (data.code === 0) {
          Message.success('鍒犻櫎鎴愬姛')
          loadRoles()
        } else {
          Message.error(data.message || '鍒犻櫎澶辫触')
        }
      } catch (e) {
        Message.error('鍒犻櫎澶辫触')
      }
    }
  })
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadRoles()
}

onMounted(() => {
  loadRoles()
  loadPermTree()
})
</script>

<style scoped>
.page-container {
  padding: 20px;
}
.breadcrumb {
  margin-bottom: 12px;
}
.search-bar {
  margin-bottom: 12px;
}
.action-bar {
  margin-bottom: 16px;
}
.search-area {
  padding: 16px;
  background: #f7f8fa;
  border-radius: 4px;
  margin-bottom: 16px;
}
.perm-group-title {
  font-weight: bold;
  margin-bottom: 8px;
  color: var(--color-text-1);
}
</style>
