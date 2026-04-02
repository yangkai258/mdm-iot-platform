<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="pro-page-container">

    <!-- 筛选区 -->
    <div class="pro-filter-bar" v-if="showFilter">
      <a-card class="filter-card">
        <a-space wrap>
          <a-select v-model="filterType" placeholder="场景类型" allow-clear style="width: 140px" @change="loadScenarios">
            <a-option value="preset">预置场景</a-option>
            <a-option value="custom">自定义场景</a-option>
          </a-select>
          <a-select v-model="filterPublic" placeholder="是否公开" allow-clear style="width: 120px" @change="loadScenarios">
            <a-option :value="true">公开</a-option>
            <a-option :value="false">私有</a-option>
          </a-select>
          <a-input-search v-model="searchKeyword" placeholder="搜索场景名称" style="width: 240px" search-button @search="loadScenarios" />
        </a-space>
      </a-card>
    </div>

    <!-- 操作栏 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="openCreateDrawer">创建场景</a-button>
        <a-button @click="openImportDialog">导入</a-button>
      </a-space>
      <a-space>
        <a-button :disabled="!selectedScenario" @click="handleExport">导出</a-button>
      </a-space>
    </div>

    <!-- 场景列表 -->
    <div class="pro-content-area">
      <a-row :gutter="[16, 16]" v-if="scenarios.length">
        <a-col :xs="24" :sm="12" :md="8" :lg="6" v-for="scene in scenarios" :key="scene.id">
          <a-card class="scenario-card" :class="{ 'scenario-card--selected': selectedScenario?.id === scene.id }" @click="selectScenario(scene)">
            <template #title>
              <div class="scenario-title">
                <span>{{ scene.scenario_name }}</span>
                <a-tag :color="scene.scenario_type === 'preset' ? 'arcoblue' : 'green'" size="small">
                  {{ scene.scenario_type === 'preset' ? '预置' : '自定义' }}
                </a-tag>
              </div>
            </template>
            <template #extra>
              <a-tag :color="scene.is_public ? 'green' : 'gray'" size="small">
                {{ scene.is_public ? '公开' : '私有' }}
              </a-tag>
            </template>
            <a-descriptions :column="1" size="small">
              <a-descriptions-item label="下载量">{{ scene.downloads || 0 }}</a-descriptions-item>
              <a-descriptions-item label="评分">
                <a-rate :value="scene.score || 0" disabled allow-half size="small" />
              </a-descriptions-item>
              <a-descriptions-item label="标签">
                <a-tag v-for="tag in (scene.tags || []).slice(0, 3)" :key="tag" size="small">{{ tag }}</a-tag>
              </a-descriptions-item>
            </a-descriptions>
            <template #actions>
              <a-button type="text" size="small" @click.stop="openEditDrawer(scene)">编辑</a-button>
              <a-button type="text" size="small" status="success" @click.stop="handleRun(scene)">运行</a-button>
              <a-button type="text" size="small" status="danger" @click.stop="handleDelete(scene)">删除</a-button>
            </template>
          </a-card>
        </a-col>
      </a-row>
      <a-empty v-else description="暂无场景" />

      <!-- 分页 -->
      <div class="pro-pagination" v-if="total > 0">
        <a-pagination :total="total" :current="page" :page-size="pageSize" show-total show-page-size @page-size-change="onPageSizeChange" @change="onPageChange" />
      </div>
    </div>

    <!-- 创建/编辑场景抽屉 -->
    <a-drawer v-model:visible="drawerVisible" :title="isEditing ? '编辑场景' : '创建场景'" :width="560" @before-ok="handleSaveScenario">
      <a-form :model="scenarioForm" layout="vertical" ref="formRef">
        <a-form-item label="场景名称" field="scenario_name" required>
          <a-input v-model="scenarioForm.scenario_name" placeholder="请输入场景名称" />
        </a-form-item>
        <a-form-item label="场景类型" field="scenario_type" required>
          <a-select v-model="scenarioForm.scenario_type" placeholder="请选择场景类型">
            <a-option value="preset">预置场景</a-option>
            <a-option value="custom">自定义场景</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="是否公开" field="is_public">
          <a-switch v-model="scenarioForm.is_public" />
        </a-form-item>
        <a-form-item label="标签" field="tags">
          <a-select v-model="scenarioForm.tags" multiple placeholder="请选择或输入标签" allow-create>
            <a-option value="日常">日常</a-option>
            <a-option value="客厅">客厅</a-option>
            <a-option value="餐厅">餐厅</a-option>
            <a-option value="夜间">夜间</a-option>
            <a-option value="户外">户外</a-option>
            <a-option value="紧急">紧急</a-option>
          </a-select>
        </a-form-item>
        <a-divider>场景配置</a-divider>
        <a-form-item label="持续时间（分钟）" field="config_duration">
          <a-input-number v-model="scenarioForm.config_duration" :min="1" :max="1440" placeholder="30" />
        </a-form-item>
        <a-form-item label="AI 启用" field="ai_enabled">
          <a-switch v-model="scenarioForm.ai_enabled" />
        </a-form-item>
      </a-form>
    </a-drawer>

    <!-- 导入对话框 -->
    <a-modal v-model:visible="importDialogVisible" title="导入场景" @before-ok="handleImport">
      <a-upload draggable :limit="1" accept=".json" @change="handleFileChange" />
      <template #footer>
        <a-button @click="importDialogVisible = false">取消</a-button>
        <a-button type="primary" :disabled="!importFile" @click="handleImport">导入</a-button>
      </template>
    </a-modal>

    <!-- 运行确认对话框 -->
    <a-modal v-model:visible="runDialogVisible" title="运行场景" @before-ok="confirmRun">
      <a-form :model="runForm" layout="vertical">
        <a-form-item label="选择宠物" required>
          <a-select v-model="runForm.pet_id" placeholder="请选择虚拟宠物">
            <a-option v-for="pet in availablePets" :key="pet.id" :value="pet.id">{{ pet.pet_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="回放速度">
          <a-input-number v-model="runForm.parameters.speed" :min="0.1" :max="5" :step="0.1" />
        </a-form-item>
        <a-form-item label="启用录制">
          <a-switch v-model="runForm.parameters.record_enabled" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { getScenarios, createScenario, updateScenario, deleteScenario, runScenario, importScenarios, exportScenario, getSimulationPets } from '@/api/simulation'

const showFilter = ref(true)
const scenarios = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const filterType = ref('')
const filterPublic = ref('')
const searchKeyword = ref('')
const selectedScenario = ref(null)

const drawerVisible = ref(false)
const isEditing = ref(false)
const editingId = ref(null)
const scenarioForm = reactive({
  scenario_name: '',
  scenario_type: 'custom',
  is_public: false,
  tags: [],
  config_duration: 30,
  ai_enabled: true
})

const importDialogVisible = ref(false)
const importFile = ref(null)

const runDialogVisible = ref(false)
const runForm = reactive({
  pet_id: null,
  parameters: {
    speed: 1.0,
    record_enabled: true
  }
})
const availablePets = ref([])

async function loadScenarios() {
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value
    }
    if (filterType.value) params.scenario_type = filterType.value
    if (filterPublic.value !== '') params.is_public = filterPublic.value
    const res = await getScenarios(params)
    scenarios.value = res.data?.items || res.data || []
    total.value = res.data?.total || 0
  } catch (e) {
    Message.error('加载场景列表失败')
  }
}

function selectScenario(scene) {
  selectedScenario.value = scene
}

function openCreateDrawer() {
  isEditing.value = false
  editingId.value = null
  Object.assign(scenarioForm, { scenario_name: '', scenario_type: 'custom', is_public: false, tags: [], config_duration: 30, ai_enabled: true })
  drawerVisible.value = true
}

function openEditDrawer(scene) {
  isEditing.value = true
  editingId.value = scene.id
  Object.assign(scenarioForm, {
    scenario_name: scene.scenario_name,
    scenario_type: scene.scenario_type,
    is_public: scene.is_public,
    tags: scene.tags || [],
    config_duration: scene.config?.duration_minutes || 30,
    ai_enabled: scene.config?.ai_enabled ?? true
  })
  drawerVisible.value = true
}

async function handleSaveScenario() {
  try {
    const data = {
      scenario_name: scenarioForm.scenario_name,
      scenario_type: scenarioForm.scenario_type,
      is_public: scenarioForm.is_public,
      tags: scenarioForm.tags,
      config: {
        duration_minutes: scenarioForm.config_duration,
        ai_enabled: scenarioForm.ai_enabled
      }
    }
    if (isEditing.value) {
      await updateScenario(editingId.value, data)
      Message.success('场景更新成功')
    } else {
      await createScenario(data)
      Message.success('场景创建成功')
    }
    drawerVisible.value = false
    loadScenarios()
  } catch (e) {
    Message.error('保存失败')
    return false
  }
}

async function handleDelete(scene) {
  try {
    await deleteScenario(scene.id)
    Message.success('删除成功')
    loadScenarios()
  } catch (e) {
    Message.error('删除失败')
  }
}

async function handleRun(scene) {
  selectedScenario.value = scene
  // 加载可用宠物
  try {
    const res = await getSimulationPets({ page: 1, page_size: 100 })
    availablePets.value = res.data?.items || res.data || []
  } catch {
    availablePets.value = []
  }
  runDialogVisible.value = true
}

async function confirmRun() {
  try {
    await runScenario(selectedScenario.value.id, runForm)
    Message.success('场景开始运行')
    runDialogVisible.value = false
  } catch (e) {
    Message.error('启动失败')
    return false
  }
}

function openImportDialog() {
  importFile.value = null
  importDialogVisible.value = true
}

function handleFileChange(file) {
  importFile.value = file
}

async function handleImport() {
  if (!importFile.value) return
  try {
    const formData = new FormData()
    formData.append('file', importFile.value)
    await importScenarios(formData)
    Message.success('导入成功')
    importDialogVisible.value = false
    loadScenarios()
  } catch (e) {
    Message.error('导入失败')
    return false
  }
}

async function handleExport() {
  if (!selectedScenario.value) return
  try {
    await exportScenario(selectedScenario.value.id)
    Message.success('导出成功')
  } catch (e) {
    Message.error('导出失败')
  }
}

function onPageChange(p) {
  page.value = p
  loadScenarios()
}

function onPageSizeChange(s) {
  pageSize.value = s
  page.value = 1
  loadScenarios()
}

onMounted(() => {
  loadScenarios()
})
</script>

<style scoped>
.scenario-title {
  display: flex;
  align-items: center;
  gap: 8px;
}
.scenario-card {
  cursor: pointer;
  transition: all 0.2s;
}
.scenario-card--selected {
  border-color: rgb(var(--primary-6));
}
.filter-card {
  background: #F2F3F5;
  border-radius: 4px;
}
</style>
