<template>
  <div class="feature-config">
    <div class="page-header">
      <a-breadcrumb>
        <a-breadcrumb-item>绯荤粺绠＄悊</a-breadcrumb-item>
        <a-breadcrumb-item>鍔熻兘閰嶇疆</a-breadcrumb-item>
      </a-breadcrumb>
      <h2>鍔熻兘閰嶇疆绠＄悊</h2>
      <p class="subtitle">鎷栨嫿璋冩暣鍒嗙粍鍜屽姛鑳界殑椤哄簭锛屾敮鎸佽法鍒嗙粍鎷栨嫿</p>
    </div>

    <div class="content-wrapper">
      <!-- 鎿嶄綔鏍?-->
      <div class="toolbar">
        <a-space>
          <a-button type="primary" @click="showGroupModal()">
            <template #icon><PlusOutlined /></template>
            鏂板鍒嗙粍
          </a-button>
          <a-button @click="showFeatureModal()">
            <template #icon><AppstoreAddOutlined /></template>
            鏂板鍔熻兘
          </a-button>
          <a-divider type="vertical" />
          <a-button @click="loadData">
            <template #icon><ReloadOutlined /></template>
            鍒锋柊
          </a-button>
        </a-space>
        <a-space>
          <a-switch v-model:checked="showDisabled" checked-children="鏄剧ず绂佺敤" un-checked-children="闅愯棌绂佺敤" />
        </a-space>
      </div>

      <!-- 鎷栨嫿鍖哄煙 -->
      <div class="drag-container">
        <a-row :gutter="16">
          <!-- 鍒嗙粍鍒楄〃 -->
          <a-col :span="8">
            <div class="section-title">
              <FolderOutlined /> 鍒嗙粍鍒楄〃
              <span class="count">({{ groups.length }})</span>
            </div>
            <draggable
              v-model="groups"
              item-key="id"
              class="group-list"
              ghost-class="ghost"
              @end="onGroupReorder"
            >
              <template #item="{ element }">
                <div
                  class="group-card"
                  :class="{ 'is-disabled': element.status === 0, 'is-active': selectedGroup?.id === element.id }"
                  @click="selectGroup(element)"
                  @dragstart="onGroupDragStart(element, 'group')"
                >
                  <div class="group-icon" :style="{ backgroundColor: element.color || '#165dff' }">
                    <component :is="getIcon(element.icon)" />
                  </div>
                  <div class="group-info">
                    <div class="group-name">{{ element.group_name }}</div>
                    <div class="group-meta">
                      <Tag v-if="element.status === 0" color="default">宸茬鐢?/Tag>
                      <span class="feature-count">{{ element.features?.length || 0 }} 涓姛鑳?/span>
                    </div>
                  </div>
                  <div class="group-actions" @click.stop>
                    <a-space>
                      <a-button type="text" size="small" @click="showGroupModal(element)">
                        <EditOutlined />
                      </a-button>
                      <a-popconfirm
                        title="鍒犻櫎鍒嗙粍灏嗗悓鏃跺垹闄ょ粍鍐呮墍鏈夊姛鑳斤紝纭畾鍒犻櫎锛?
                        @confirm="deleteGroup(element.id)"
                      >
                        <a-button type="text" danger size="small">
                          <DeleteOutlined />
                        </a-button>
                      </a-popconfirm>
                    </a-space>
                  </div>
                  <div class="drag-handle">
                    <DragOutlined />
                  </div>
                </div>
              </template>
            </draggable>
          </a-col>

          <!-- 鍔熻兘鍒楄〃 -->
          <a-col :span="16">
            <div class="section-title">
              <AppstoreOutlined /> 鍔熻兘鍒楄〃
              <span class="count">({{ currentFeatures.length }})</span>
              <span v-if="selectedGroup" class="current-group">- {{ selectedGroup.group_name }}</span>
            </div>

            <div v-if="!selectedGroup" class="empty-hint">
              <InfoCircleOutlined />
              璇蜂粠宸︿晶閫夋嫨涓€涓垎缁勶紝鎴栫洿鎺ユ嫋鎷藉姛鑳藉埌鍒嗙粍鍗＄墖涓?
            </div>

            <draggable
              v-else
              v-model="currentFeatures"
              item-key="id"
              class="feature-list"
              group="features"
              ghost-class="ghost"
              @end="onFeatureReorder"
            >
              <template #item="{ element }">
                <div
                  class="feature-card"
                  :class="{ 'is-disabled': element.status === 0, 'is-default': element.is_default === 1 }"
                  @dragstart="onFeatureDragStart(element, 'feature')"
                >
                  <div class="feature-icon">
                    <component :is="getIcon(element.icon)" />
                  </div>
                  <div class="feature-info">
                    <div class="feature-name">
                      {{ element.feature_name }}
                      <a-tag v-if="element.badge" color="processing" size="small">{{ element.badge }}</a-tag>
                      <Tag v-if="element.is_default === 1" color="gold" size="small">榛樿</Tag>
                      <Tag v-if="element.status === 0" color="default" size="small">绂佺敤</Tag>
                    </div>
                    <div class="feature-meta">
                      <span class="route">{{ element.route_path || '鏃犺矾鐢? }}</span>
                      <span v-if="element.permission" class="permission">{{ element.permission }}</span>
                    </div>
                  </div>
                  <div class="feature-actions" @click.stop>
                    <a-space>
                      <a-switch
                        :checked="element.status === 1"
                        size="small"
                        @change="(checked) => toggleFeatureStatus(element, checked)"
                      />
                      <a-button type="text" size="small" @click="showFeatureModal(element)">
                        <EditOutlined />
                      </a-button>
                      <a-popconfirm
                        title="纭畾鍒犻櫎姝ゅ姛鑳斤紵"
                        @confirm="deleteFeature(element.id)"
                      >
                        <a-button type="text" danger size="small">
                          <DeleteOutlined />
                        </a-button>
                      </a-popconfirm>
                    </a-space>
                  </div>
                  <div class="drag-handle">
                    <DragOutlined />
                  </div>
                </div>
              </template>
            </draggable>

            <!-- 鏈垎缁勫姛鑳?-->
            <div v-if="ungroupedFeatures.length > 0" class="ungrouped-section">
              <div class="section-title">
                <UnorderedListOutlined /> 鏈垎缁勫姛鑳?
                <span class="count">({{ ungroupedFeatures.length }})</span>
              </div>
              <draggable
                v-model="ungroupedFeatures"
                item-key="id"
                class="feature-list"
                group="features"
                ghost-class="ghost"
              >
                <template #item="{ element }">
                  <div
                    class="feature-card"
                    :class="{ 'is-disabled': element.status === 0, 'is-default': element.is_default === 1 }"
                  >
                    <div class="feature-icon">
                      <component :is="getIcon(element.icon)" />
                    </div>
                    <div class="feature-info">
                      <div class="feature-name">
                        {{ element.feature_name }}
                        <a-tag v-if="element.badge" color="processing" size="small">{{ element.badge }}</a-tag>
                      </div>
                      <div class="feature-meta">
                        <span class="route">{{ element.route_path || '鏃犺矾鐢? }}</span>
                      </div>
                    </div>
                    <div class="feature-actions" @click.stop>
                      <a-button type="text" size="small" @click="showFeatureModal(element)">
                        <EditOutlined />
                      </a-button>
                      <a-popconfirm
                        title="纭畾鍒犻櫎姝ゅ姛鑳斤紵"
                        @confirm="deleteFeature(element.id)"
                      >
                        <a-button type="text" danger size="small">
                          <DeleteOutlined />
                        </a-button>
                      </a-popconfirm>
                    </div>
                    <div class="drag-handle">
                      <DragOutlined />
                    </div>
                  </div>
                </template>
              </draggable>
            </div>
          </a-col>
        </a-row>
      </div>
    </div>

    <!-- 鍒嗙粍缂栬緫寮圭獥 -->
    <a-modal
      v-model:open="groupModalVisible"
      :title="editingGroup ? '缂栬緫鍒嗙粍' : '鏂板鍒嗙粍'"
      @ok="saveGroup"
      width="500px"
    >
      <a-form :model="groupForm" layout="vertical">
        <a-form-item label="鍒嗙粍鍚嶇О" required>
          <a-input v-model:value="groupForm.group_name" placeholder="璇疯緭鍏ュ垎缁勫悕绉" />
        </a-form-item>
        <a-form-item label="鍒嗙粍缂栫爜">
          <a-input v-model:value="groupForm.group_code" placeholder="鍞竴缂栫爜锛屼笉濉嚜鍔ㄧ敓鎴" />
        </a-form-item>
        <a-form-item label="鍥炬爣">
          <a-select v-model:value="groupForm.icon" placeholder="閫夋嫨鍥炬爣" allow-clear>
            <a-select-option value="Folder">Folder</a-select-option>
            <a-select-option value="FolderOpen">FolderOpen</a-select-option>
            <a-select-option value="Cpu">Cpu</a-select-option>
            <a-select-option value="User">User</a-select-option>
            <a-select-option value="Robot">Robot</a-select-option>
            <a-select-option value="Bell">Bell</a-select-option>
            <a-select-option value="Setting">Setting</a-select-option>
            <a-select-option value="Dashboard">Dashboard</a-select-option>
            <a-select-option value="Appstore">Appstore</a-select-option>
            <a-select-option value="Cloud">Cloud</a-select-option>
            <a-select-option value="Lock">Lock</a-select-option>
            <a-select-option value="Key">Key</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="棰滆壊">
          <color-picker v-model="groupForm.color" />
        </a-form-item>
        <a-form-item label="鎻忚堪">
          <a-textarea v-model:value="groupForm.description" :rows="2" />
        </a-form-item>
        <a-form-item label="鐘舵€?>
          <a-radio-group v-model:value="groupForm.status">
            <a-radio :value="1">鍚敤</a-radio>
            <a-radio :value="0">绂佺敤</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 鍔熻兘缂栬緫寮圭獥 -->
    <a-modal
      v-model:open="featureModalVisible"
      :title="editingFeature ? '缂栬緫鍔熻兘' : '鏂板鍔熻兘'"
      @ok="saveFeature"
      width="600px"
    >
      <a-form :model="featureForm" layout="vertical">
        <a-form-item label="鎵€灞炲垎缁?>
          <a-select v-model:value="featureForm.group_id" placeholder="閫夋嫨鍒嗙粍">
            <a-select-option :value="undefined">鏃狅紙鏈垎缁勶級</a-select-option>
            <a-select-option v-for="g in groups" :key="g.id" :value="g.id">
              {{ g.group_name }}
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="鍔熻兘鍚嶇О" required>
          <a-input v-model:value="featureForm.feature_name" placeholder="璇疯緭鍏ュ姛鑳藉悕绉" />
        </a-form-item>
        <a-form-item label="鍔熻兘缂栫爜">
          <a-input v-model:value="featureForm.feature_code" placeholder="鍞竴缂栫爜" />
        </a-form-item>
        <a-row :gutter="8">
          <a-col :span="12">
            <a-form-item label="鍥炬爣">
              <a-select v-model:value="featureForm.icon" placeholder="閫夋嫨鍥炬爣" allow-clear>
                <a-select-option value="Home">Home</a-select-option>
                <a-select-option value="User">User</a-select-option>
                <a-select-option value="Setting">Setting</a-select-option>
                <a-select-option value="Dashboard">Dashboard</a-select-option>
                <a-select-option value="Device">Device</a-select-option>
                <a-select-option value="Bell">Bell</a-select-option>
                <a-select-option value="Chat">Chat</a-select-option>
                <a-select-option value="Heart">Heart</a-select-option>
                <a-select-option value="Star">Star</a-select-option>
                <a-select-option value="Plus">Plus</a-select-option>
                <a-select-option value="Edit">Edit</a-select-option>
                <a-select-option value="Delete">Delete</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="鎺掑簭">
              <a-input-number v-model:value="featureForm.sort" :min="0" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="璺敱璺緞">
          <a-input v-model:value="featureForm.route_path" placeholder="/path/to/page" />
        </a-form-item>
        <a-form-item label="缁勪欢璺緞">
          <a-input v-model:value="featureForm.component" placeholder="ComponentNameView" />
        </a-form-item>
        <a-form-item label="鏉冮檺缂栫爜">
          <a-input v-model:value="featureForm.permission" placeholder="module:action" />
        </a-form-item>
        <a-form-item label="寰界珷">
          <a-input v-model:value="featureForm.badge" placeholder="濡傦細鏂般€丅eta銆丠ot" />
        </a-form-item>
        <a-form-item label="鎻忚堪">
          <a-textarea v-model:value="featureForm.description" :rows="2" />
        </a-form-item>
        <a-row :gutter="8">
          <a-col :span="8">
            <a-form-item label="鐘舵€?>
              <a-radio-group v-model:value="featureForm.status">
                <a-radio :value="1">鍚敤</a-radio>
                <a-radio :value="0">绂佺敤</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="榛樿閫変腑">
              <a-switch v-model:checked="featureForm.is_default" />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import draggable from 'vuedraggable'
import {
  PlusOutlined,
  AppstoreAddOutlined,
  ReloadOutlined,
  EditOutlined,
  DeleteOutlined,
  DragOutlined,
  FolderOutlined,
  AppstoreOutlined,
  InfoCircleOutlined,
  UnorderedListOutlined,
  // Icon components
  HomeOutlined,
  UserOutlined,
  SettingOutlined,
  DashboardOutlined,
  BellOutlined,
  RobotOutlined,
  StarOutlined,
  PlusOutlined as PlusIcon,
  CloudOutlined,
  LockOutlined,
  KeyOutlined,
  FolderOpenOutlined,
  DeviceOutlined,
  ChatOutlined,
  HeartOutlined,
  EditOutlined as EditIcon,
  DeleteOutlined as DeleteIcon,
} from '@ant-design/icons-vue'

const iconMap = {
  Home: HomeOutlined,
  User: UserOutlined,
  Setting: SettingOutlined,
  Dashboard: DashboardOutlined,
  Bell: BellOutlined,
  Robot: RobotOutlined,
  Star: StarOutlined,
  Plus: PlusIcon,
  Cloud: CloudOutlined,
  Lock: LockOutlined,
  Key: KeyOutlined,
  FolderOpen: FolderOpenOutlined,
  Device: DeviceOutlined,
  Chat: ChatOutlined,
  Heart: HeartOutlined,
  Edit: EditIcon,
  Delete: DeleteIcon,
  Folder: FolderOutlined,
  Appstore: AppstoreOutlined,
}

const getIcon = (iconName) => iconMap[iconName] || AppstoreOutlined

// State
const groups = ref([])
const selectedGroup = ref(null)
const showDisabled = ref(true)
const loading = ref(false)

// Drag state
const draggedItem = ref(null)
const dragType = ref(null)

// Modal state
const groupModalVisible = ref(false)
const featureModalVisible = ref(false)
const editingGroup = ref(null)
const editingFeature = ref(null)

// Forms
const groupForm = reactive({
  group_name: '',
  group_code: '',
  icon: '',
  color: '#165dff',
  description: '',
  status: 1,
})

const featureForm = reactive({
  group_id: undefined,
  feature_name: '',
  feature_code: '',
  icon: '',
  route_path: '',
  component: '',
  permission: '',
  sort: 0,
  status: 1,
  is_default: false,
  badge: '',
  description: '',
})

// Computed
const currentFeatures = computed({
  get: () => {
    if (!selectedGroup.value) return []
    return selectedGroup.value.features || []
  },
  set: (val) => {
    if (selectedGroup.value) {
      selectedGroup.value.features = val
    }
  }
})

const ungroupedFeatures = computed(() => {
  return groups.value.reduce((acc, g) => acc.concat(g.features || []), [])
    .filter(f => !f.group_id)
})

// Methods
const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/v1/feature-config/groups', {
      headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      groups.value = data.data || []
      if (groups.value.length > 0 && !selectedGroup.value) {
        selectedGroup.value = groups.value[0]
      }
    } else {
      message.error(data.message || '鍔犺浇鏁版嵁澶辫触')
    }
  } catch (err) {
    message.error('鍔犺浇鏁版嵁澶辫触: ' + err.message)
  } finally {
    loading.value = false
  }
}

const selectGroup = (group) => {
  selectedGroup.value = group
}

const showGroupModal = (group = null) => {
  editingGroup.value = group
  if (group) {
    Object.assign(groupForm, {
      group_name: group.group_name,
      group_code: group.group_code,
      icon: group.icon,
      color: group.color,
      description: group.description,
      status: group.status,
    })
  } else {
    Object.assign(groupForm, {
      group_name: '',
      group_code: '',
      icon: '',
      color: '#165dff',
      description: '',
      status: 1,
    })
  }
  groupModalVisible.value = true
}

const saveGroup = async () => {
  if (!groupForm.group_name) {
    message.warning('璇疯緭鍏ュ垎缁勫悕绉?)
    return
  }

  try {
    const method = editingGroup.value ? 'PUT' : 'POST'
    const url = editingGroup.value
      ? `/api/v1/feature-config/groups/${editingGroup.value.id}`
      : '/api/v1/feature-config/groups'

    const res = await fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify(groupForm)
    })
    const data = await res.json()

    if (data.code === 0) {
      message.success('淇濆瓨鎴愬姛')
      groupModalVisible.value = false
      loadData()
    } else {
      message.error(data.message || '淇濆瓨澶辫触')
    }
  } catch (err) {
    message.error('淇濆瓨澶辫触: ' + err.message)
  }
}

const deleteGroup = async (id) => {
  try {
    const res = await fetch(`/api/v1/feature-config/groups/${id}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      message.success('鍒犻櫎鎴愬姛')
      if (selectedGroup.value?.id === id) {
        selectedGroup.value = null
      }
      loadData()
    } else {
      message.error(data.message || '鍒犻櫎澶辫触')
    }
  } catch (err) {
    message.error('鍒犻櫎澶辫触: ' + err.message)
  }
}

const showFeatureModal = (feature = null) => {
  editingFeature.value = feature
  if (feature) {
    Object.assign(featureForm, {
      group_id: feature.group_id,
      feature_name: feature.feature_name,
      feature_code: feature.feature_code,
      icon: feature.icon,
      route_path: feature.route_path,
      component: feature.component,
      permission: feature.permission,
      sort: feature.sort,
      status: feature.status,
      is_default: feature.is_default === 1,
      badge: feature.badge,
      description: feature.description,
    })
  } else {
    Object.assign(featureForm, {
      group_id: selectedGroup.value?.id,
      feature_name: '',
      feature_code: '',
      icon: '',
      route_path: '',
      component: '',
      permission: '',
      sort: 0,
      status: 1,
      is_default: false,
      badge: '',
      description: '',
    })
  }
  featureModalVisible.value = true
}

const saveFeature = async () => {
  if (!featureForm.feature_name) {
    message.warning('璇疯緭鍏ュ姛鑳藉悕绉?)
    return
  }

  try {
    const method = editingFeature.value ? 'PUT' : 'POST'
    const url = editingFeature.value
      ? `/api/v1/feature-config/features/${editingFeature.value.id}`
      : '/api/v1/feature-config/features'

    const payload = {
      ...featureForm,
      is_default: featureForm.is_default ? 1 : 0,
    }

    const res = await fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify(payload)
    })
    const data = await res.json()

    if (data.code === 0) {
      message.success('淇濆瓨鎴愬姛')
      featureModalVisible.value = false
      loadData()
    } else {
      message.error(data.message || '淇濆瓨澶辫触')
    }
  } catch (err) {
    message.error('淇濆瓨澶辫触: ' + err.message)
  }
}

const deleteFeature = async (id) => {
  try {
    const res = await fetch(`/api/v1/feature-config/features/${id}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      message.success('鍒犻櫎鎴愬姛')
      loadData()
    } else {
      message.error(data.message || '鍒犻櫎澶辫触')
    }
  } catch (err) {
    message.error('鍒犻櫎澶辫触: ' + err.message)
  }
}

const toggleFeatureStatus = async (feature, checked) => {
  feature.status = checked ? 1 : 0
  try {
    const res = await fetch(`/api/v1/feature-config/features/${feature.id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({ ...feature, status: feature.status })
    })
    const data = await res.json()
    if (data.code !== 0) {
      message.error('鏇存柊鐘舵€佸け璐?)
      feature.status = checked ? 0 : 1
    }
  } catch (err) {
    message.error('鏇存柊鐘舵€佸け璐? ' + err.message)
    feature.status = checked ? 0 : 1
  }
}

// Drag handlers
const onGroupDragStart = (item) => {
  draggedItem.value = item
  dragType.value = 'group'
}

const onFeatureDragStart = (item) => {
  draggedItem.value = item
  dragType.value = 'feature'
}

const onGroupReorder = async () => {
  // Save new order
  const items = groups.value.map((g, idx) => ({ id: g.id, sort: idx }))
  try {
    await fetch('/api/v1/feature-config/sort', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({ items })
    })
  } catch (err) {
    message.error('淇濆瓨鎺掑簭澶辫触')
  }
}

const onFeatureReorder = async () => {
  // Save new order for current group
  const items = currentFeatures.value.map((f, idx) => ({
    id: f.id,
    sort: idx,
    group_id: selectedGroup.value.id
  }))

  // Also include ungrouped features
  ungroupedFeatures.value.forEach((f, idx) => {
    if (!items.find(i => i.id === f.id)) {
      items.push({ id: f.id, sort: idx, group_id: null })
    }
  })

  try {
    await fetch('/api/v1/feature-config/sort', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({ items })
    })
  } catch (err) {
    message.error('淇濆瓨鎺掑簭澶辫触')
  }
}

// Lifecycle
onMounted(() => {
  loadData()
})
</script>

<style scoped>
.feature-config {
  padding: 24px;
}

.page-header {
  margin-bottom: 24px;
}

.page-header h2 {
  margin: 8px 0 4px;
  font-size: 20px;
  font-weight: 600;
}

.subtitle {
  color: #8c8c8c;
  font-size: 14px;
  margin: 0;
}

.content-wrapper {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.drag-container {
  min-height: 500px;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #262626;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.section-title .count {
  color: #8c8c8c;
  font-weight: normal;
}

.section-title .current-group {
  color: #1890ff;
  font-weight: normal;
}

.group-list,
.feature-list {
  min-height: 100px;
}

.group-card {
  display: flex;
  align-items: center;
  padding: 12px;
  margin-bottom: 8px;
  background: #fafafa;
  border: 1px solid #f0f0f0;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.group-card:hover {
  border-color: #1890ff;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.15);
}

.group-card.is-disabled {
  opacity: 0.6;
  background: #f5f5f5;
}

.group-card.is-active {
  border-color: #1890ff;
  background: #e6f7ff;
}

.group-card.ghost {
  opacity: 0.5;
  background: #e6f7ff;
  border: 2px dashed #1890ff;
}

.group-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 18px;
  margin-right: 12px;
  flex-shrink: 0;
}

.group-info {
  flex: 1;
  min-width: 0;
}

.group-name {
  font-size: 14px;
  font-weight: 500;
  color: #262626;
}

.group-meta {
  font-size: 12px;
  color: #8c8c8c;
  display: flex;
  align-items: center;
  gap: 8px;
}

.feature-count {
  color: #8c8c8c;
}

.group-actions {
  opacity: 0;
  transition: opacity 0.2s;
}

.group-card:hover .group-actions {
  opacity: 1;
}

.drag-handle {
  color: #d9d9d9;
  cursor: grab;
  padding: 4px;
  margin-left: 8px;
}

.drag-handle:active {
  cursor: grabbing;
}

.feature-card {
  display: flex;
  align-items: center;
  padding: 10px 12px;
  margin-bottom: 6px;
  background: #fff;
  border: 1px solid #f0f0f0;
  border-radius: 6px;
  transition: all 0.2s;
}

.feature-card:hover {
  border-color: #52c41a;
  box-shadow: 0 2px 6px rgba(82, 196, 26, 0.15);
}

.feature-card.is-disabled {
  opacity: 0.6;
  background: #fafafa;
}

.feature-card.is-default {
  background: #fffbe6;
  border-color: #ffe58f;
}

.feature-card.ghost {
  opacity: 0.5;
  background: #f6ffed;
  border: 2px dashed #52c41a;
}

.feature-icon {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  background: #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #595959;
  margin-right: 10px;
  flex-shrink: 0;
}

.feature-info {
  flex: 1;
  min-width: 0;
}

.feature-name {
  font-size: 14px;
  color: #262626;
  display: flex;
  align-items: center;
  gap: 6px;
}

.feature-meta {
  font-size: 12px;
  color: #8c8c8c;
  display: flex;
  gap: 12px;
  margin-top: 2px;
}

.route {
  font-family: monospace;
}

.permission {
  color: #722ed1;
}

.feature-actions {
  opacity: 0;
  transition: opacity 0.2s;
}

.feature-card:hover .feature-actions {
  opacity: 1;
}

.empty-hint {
  text-align: center;
  padding: 40px;
  color: #8c8c8c;
  background: #fafafa;
  border-radius: 8px;
  border: 1px dashed #d9d9d9;
}

.ungrouped-section {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #f0f0f0;
}
</style>
