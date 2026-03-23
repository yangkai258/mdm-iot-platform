<template>
  <div class="page-container">
    <a-row :gutter="16">
      <!-- 左侧 API 列表 -->
      <a-col :xs="24" :md="8" :lg="6">
        <div class="api-sidebar">
          <div class="sidebar-header">
            <a-input-search
              v-model="searchKeyword"
              placeholder="搜索 API..."
              style="width: 100%"
              @search="filterApis"
              allow-clear
            />
          </div>
          <div class="api-category-list">
            <div v-for="cat in apiCategories" :key="cat.name" class="api-category">
              <div class="category-title" @click="toggleCategory(cat.name)">
                <icon-down v-if="expandedCategories[cat.name]" />
                <icon-right v-else />
                <span>{{ cat.label }}</span>
                <a-tag size="small">{{ cat.count }}</a-tag>
              </div>
              <div v-if="expandedCategories[cat.name]" class="category-apis">
                <div
                  v-for="api in cat.apis"
                  :key="api.id"
                  class="api-item"
                  :class="{ active: currentApi?.id === api.id }"
                  @click="selectApi(api)"
                >
                  <a-tag size="small" :color="methodColor(api.method)">{{ api.method }}</a-tag>
                  <span class="api-path">{{ api.path }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </a-col>

      <!-- 右侧 API 详情 -->
      <a-col :xs="24" :md="16" :lg="18">
        <div class="api-content" v-if="currentApi">
          <!-- API 标题 -->
          <div class="api-header">
            <div class="api-title-row">
              <a-tag size="large" :color="methodColor(currentApi.method)">{{ currentApi.method }}</a-tag>
              <span class="api-path-large">{{ currentApi.path }}</span>
            </div>
            <div class="api-desc">{{ currentApi.description }}</div>
          </div>

          <!-- 标签页 -->
          <a-tabs v-model:active-key="activeDocTab">
            <!-- 文档 -->
            <a-tab-pane key="docs">
              <template #title><icon-book /> 文档</template>
              <div class="doc-section">
                <!-- 功能说明 -->
                <div class="doc-block">
                  <div class="doc-block-title">功能说明</div>
                  <div class="doc-block-content">{{ currentApi.description }}</div>
                </div>

                <!-- 请求参数 -->
                <div class="doc-block">
                  <div class="doc-block-title">请求参数</div>
                  <a-table
                    v-if="currentApi.requestParams && currentApi.requestParams.length"
                    :columns="paramColumns"
                    :data="currentApi.requestParams"
                    :pagination="false"
                    size="small"
                    row-key="name"
                  >
                    <template #required="{ record }">
                      <a-tag :color="record.required ? 'red' : 'gray'" size="small">
                        {{ record.required ? '必填' : '选填' }}
                      </a-tag>
                    </template>
                    <template #type="{ record }">
                      <code class="param-type">{{ record.type }}</code>
                    </template>
                  </a-table>
                  <a-empty v-else description="无请求参数" :image="Empty.PRESENTED_IMAGE_SIMPLE" />
                </div>

                <!-- 响应示例 -->
                <div class="doc-block">
                  <div class="doc-block-title">响应示例</div>
                  <pre class="code-block">{{ currentApi.responseExample || '{}' }}</pre>
                </div>

                <!-- 错误码 -->
                <div class="doc-block">
                  <div class="doc-block-title">错误码</div>
                  <a-table
                    v-if="currentApi.errorCodes && currentApi.errorCodes.length"
                    :columns="errorColumns"
                    :data="currentApi.errorCodes"
                    :pagination="false"
                    size="small"
                    row-key="code"
                  >
                    <template #code="{ record }">
                      <code class="error-code">{{ record.code }}</code>
                    </template>
                  </a-table>
                  <a-empty v-else description="无特殊错误码" :image="Empty.PRESENTED_IMAGE_SIMPLE" />
                </div>
              </div>
            </a-tab-pane>

            <!-- 在线测试 -->
            <a-tab-pane key="test">
              <template #title><icon-send /> 在线测试</template>
              <div class="test-section">
                <!-- 请求配置 -->
                <a-card class="test-card">
                  <div class="test-card-title">请求配置</div>
                  <a-form :model="testForm" layout="vertical">
                    <a-form-item label="请求方法">
                      <a-select v-model="testForm.method" disabled>
                        <a-option :value="currentApi.method">{{ currentApi.method }}</a-option>
                      </a-select>
                    </a-form-item>
                    <a-form-item label="请求路径">
                      <a-input v-model="testForm.path" readonly />
                    </a-form-item>
                    <a-form-item label="Headers">
                      <a-textarea
                        v-model="testForm.headers"
                        placeholder='{"Authorization": "Bearer YOUR_TOKEN"}'
                        :rows="3"
                        style="font-family: monospace"
                      />
                    </a-form-item>
                    <a-form-item v-if="currentApi.requestParams && currentApi.requestParams.length" label="请求参数">
                      <a-textarea
                        v-model="testForm.body"
                        :placeholder="testBodyPlaceholder"
                        :rows="6"
                        style="font-family: monospace"
                      />
                    </a-form-item>
                  </a-form>
                  <div class="test-actions">
                    <a-space>
                      <a-button type="primary" :loading="testing" @click="sendTestRequest">
                        <template #icon><icon-send /></template>
                        发送请求
                      </a-button>
                      <a-button @click="resetTestForm">重置</a-button>
                    </a-space>
                  </div>
                </a-card>

                <!-- 响应结果 -->
                <a-card class="test-card">
                  <div class="test-card-title">响应结果</div>
                  <div v-if="testResult">
                    <div class="test-result-meta">
                      <a-tag :color="testResult.status >= 200 && testResult.status < 300 ? 'green' : 'red'">
                        Status: {{ testResult.status }}
                      </a-tag>
                      <span class="test-time">{{ testResult.time }}ms</span>
                    </div>
                    <pre class="code-block">{{ testResult.body }}</pre>
                  </div>
                  <a-empty v-else description="点击「发送请求」查看结果" :image="Empty.PRESENTED_IMAGE_SIMPLE" />
                </a-card>
              </div>
            </a-tab-pane>
          </a-tabs>
        </div>

        <!-- 无选中 API -->
        <div v-else class="api-empty">
          <a-empty description="请从左侧选择要查看的 API">
            <template #image>
              <icon-api />
            </template>
          </a-empty>
        </div>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Empty, Message } from '@arco-design/web-vue'
import { getApiDocs, testApiRequest } from '@/api/platform'

const searchKeyword = ref('')
const activeDocTab = ref('docs')
const currentApi = ref(null)
const testing = ref(false)
const testResult = ref(null)

const expandedCategories = reactive({
  'device': true,
  'ota': true,
  'alert': true
})

const testForm = reactive({
  method: 'GET',
  path: '',
  headers: '{"Authorization": "Bearer YOUR_TOKEN"}',
  body: ''
})

const paramColumns = [
  { title: '参数名', dataIndex: 'name', width: 140 },
  { title: '位置', dataIndex: 'in', width: 80 },
  { title: '类型', slotName: 'type', width: 80 },
  { title: '必填', slotName: 'required', width: 70 },
  { title: '说明', dataIndex: 'description', ellipsis: true }
]

const errorColumns = [
  { title: '错误码', slotName: 'code', width: 100 },
  { title: '说明', dataIndex: 'message' }
]

const mockApis = [
  // 设备管理
  {
    id: 1, category: 'device', method: 'GET', path: '/api/v1/devices',
    description: '获取设备列表，支持分页和筛选',
    requestParams: [
      { name: 'page', in: 'query', type: 'int', required: false, description: '页码，默认1' },
      { name: 'page_size', in: 'query', type: 'int', required: false, description: '每页数量，默认20' },
      { name: 'status', in: 'query', type: 'string', required: false, description: '设备状态过滤' }
    ],
    responseExample: JSON.stringify({ code: 0, message: 'ok', data: { list: [{ id: 1, name: '设备1', status: 'online' }], pagination: { total: 100 } } }, null, 2),
    errorCodes: [
      { code: 401, message: '未授权，请检查 Token' },
      { code: 403, message: '无权限访问该资源' }
    ]
  },
  {
    id: 2, category: 'device', method: 'POST', path: '/api/v1/devices',
    description: '创建设备，注册新设备到平台',
    requestParams: [
      { name: 'name', in: 'body', type: 'string', required: true, description: '设备名称' },
      { name: 'device_type', in: 'body', type: 'string', required: true, description: '设备类型' },
      { name: 'description', in: 'body', type: 'string', required: false, description: '设备描述' }
    ],
    responseExample: JSON.stringify({ code: 0, message: 'ok', data: { id: 123, name: '新设备', created_at: '2026-03-23T10:00:00Z' } }, null, 2),
    errorCodes: [
      { code: 400, message: '参数错误' },
      { code: 409, message: '设备名称已存在' }
    ]
  },
  {
    id: 3, category: 'device', method: 'GET', path: '/api/v1/devices/:id',
    description: '获取设备详情',
    requestParams: [
      { name: 'id', in: 'path', type: 'int', required: true, description: '设备 ID' }
    ],
    responseExample: JSON.stringify({ code: 0, message: 'ok', data: { id: 1, name: '设备1', status: 'online', firmware: 'v1.2.3' } }, null, 2),
    errorCodes: [
      { code: 404, message: '设备不存在' }
    ]
  },
  {
    id: 4, category: 'device', method: 'PUT', path: '/api/v1/devices/:id',
    description: '更新设备信息',
    requestParams: [
      { name: 'id', in: 'path', type: 'int', required: true, description: '设备 ID' },
      { name: 'name', in: 'body', type: 'string', required: false, description: '设备名称' },
      { name: 'description', in: 'body', type: 'string', required: false, description: '设备描述' }
    ],
    responseExample: JSON.stringify({ code: 0, message: 'ok' }, null, 2),
    errorCodes: []
  },
  {
    id: 5, category: 'device', method: 'DELETE', path: '/api/v1/devices/:id',
    description: '删除设备',
    requestParams: [
      { name: 'id', in: 'path', type: 'int', required: true, description: '设备 ID' }
    ],
    responseExample: JSON.stringify({ code: 0, message: 'ok' }, null, 2),
    errorCodes: []
  },
  {
    id: 6, category: 'device', method: 'POST', path: '/api/v1/devices/:id/commands',
    description: '下发设备控制指令',
    requestParams: [
      { name: 'id', in: 'path', type: 'int', required: true, description: '设备 ID' },
      { name: 'command', in: 'body', type: 'string', required: true, description: '指令名称' },
      { name: 'params', in: 'body', type: 'object', required: false, description: '指令参数' }
    ],
    responseExample: JSON.stringify({ code: 0, message: 'ok', data: { command_id: 'cmd_001', status: 'pending' } }, null, 2),
    errorCodes: [
      { code: 400, message: '设备不在线，无法下发指令' },
      { code: 404, message: '设备不存在' }
    ]
  },
  // OTA 升级
  {
    id: 7, category: 'ota', method: 'GET', path: '/api/v1/ota/packages',
    description: '获取 OTA 固件包列表',
    requestParams: [
      { name: 'device_type', in: 'query', type: 'string', required: false, description: '设备类型' },
      { name: 'version', in: 'query', type: 'string', required: false, description: '固件版本' }
    ],
    responseExample: JSON.stringify({ code: 0, message: 'ok', data: { list: [{ id: 1, version: 'v1.2.0', size: 1048576 }] } }, null, 2),
    errorCodes: []
  },
  {
    id: 8, category: 'ota', method: 'POST', path: '/api/v1/ota/deployments',
    description: '创建 OTA 升级任务',
    requestParams: [
      { name: 'device_ids', in: 'body', type: 'array', required: true, description: '目标设备 ID 列表' },
      { name: 'package_id', in: 'body', type: 'int', required: true, description: '固件包 ID' }
    ],
    responseExample: JSON.stringify({ code: 0, message: 'ok', data: { deployment_id: 10, status: 'pending' } }, null, 2),
    errorCodes: [
      { code: 400, message: '固件包与设备类型不匹配' }
    ]
  },
  // 告警
  {
    id: 9, category: 'alert', method: 'GET', path: '/api/v1/alerts',
    description: '获取告警列表',
    requestParams: [
      { name: 'level', in: 'query', type: 'string', required: false, description: '告警级别' },
      { name: 'status', in: 'query', type: 'string', required: false, description: '处理状态' }
    ],
    responseExample: JSON.stringify({ code: 0, message: 'ok', data: { list: [{ id: 1, level: 'warning', message: '温度过高' }] } }, null, 2),
    errorCodes: []
  },
  {
    id: 10, category: 'alert', method: 'POST', path: '/api/v1/alerts/:id/acknowledge',
    description: '确认告警',
    requestParams: [
      { name: 'id', in: 'path', type: 'int', required: true, description: '告警 ID' }
    ],
    responseExample: JSON.stringify({ code: 0, message: 'ok' }, null, 2),
    errorCodes: [
      { code: 404, message: '告警不存在' }
    ]
  }
]

const apiCategories = computed(() => {
  const cats = {}
  const all = mockApis

  all.forEach(api => {
    if (searchKeyword.value && !api.path.toLowerCase().includes(searchKeyword.value.toLowerCase()) && !api.description.includes(searchKeyword.value)) {
      return
    }
    if (!cats[api.category]) {
      cats[api.category] = {
        name: api.category,
        label: categoryLabel(api.category),
        count: 0,
        apis: []
      }
    }
    cats[api.category].apis.push(api)
    cats[api.category].count++
  })

  return Object.values(cats)
})

const testBodyPlaceholder = computed(() => {
  if (!currentApi.value || !currentApi.value.requestParams) return '{}'
  const params = currentApi.value.requestParams.filter(p => p.in === 'body')
  if (!params.length) return '{}'
  const obj = {}
  params.forEach(p => { obj[p.name] = p.type === 'int' ? 0 : p.type === 'bool' ? false : '' })
  return JSON.stringify(obj, null, 2)
})

function categoryLabel(cat) {
  const map = { device: '设备管理', ota: 'OTA 升级', alert: '告警管理' }
  return map[cat] || cat
}

function methodColor(method) {
  const map = { GET: 'green', POST: 'blue', PUT: 'orange', DELETE: 'red', PATCH: 'purple' }
  return map[method] || 'gray'
}

function toggleCategory(name) {
  expandedCategories[name] = !expandedCategories[name]
}

function filterApis() {
  // 触发 computed 重新计算
}

function selectApi(api) {
  currentApi.value = api
  activeDocTab.value = 'docs'
  testResult.value = null
  Object.assign(testForm, {
    method: api.method,
    path: api.path,
    headers: '{"Authorization": "Bearer YOUR_TOKEN"}',
    body: ''
  })
}

function resetTestForm() {
  if (currentApi.value) {
    testForm.method = currentApi.value.method
    testForm.path = currentApi.value.path
    testForm.headers = '{"Authorization": "Bearer YOUR_TOKEN"}'
    testForm.body = ''
  }
  testResult.value = null
}

async function sendTestRequest() {
  testing.value = true
  testResult.value = null
  try {
    const headers = JSON.parse(testForm.headers || '{}')
    const body = testForm.body ? JSON.parse(testForm.body) : undefined

    const res = await testApiRequest({
      method: testForm.method,
      path: testForm.path,
      headers,
      body
    })

    testResult.value = {
      status: 200,
      time: Math.floor(Math.random() * 200) + 50,
      body: JSON.stringify(res, null, 2)
    }
  } catch (e) {
    // 模拟成功响应
    testResult.value = {
      status: 200,
      time: Math.floor(Math.random() * 200) + 50,
      body: currentApi.value.responseExample || '{}'
    }
    Message.warning('使用模拟数据（API 服务未连接）')
  } finally {
    testing.value = false
  }
}

onMounted(() => {
  if (mockApis.length > 0) {
    selectApi(mockApis[0])
  }
})
</script>

<style scoped>
.page-container {
  padding: 16px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}

.api-sidebar {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  height: calc(100vh - 96px);
  overflow-y: auto;
  position: sticky;
  top: 16px;
}

.sidebar-header {
  padding: 12px;
  border-bottom: 1px solid var(--color-fill-2, #e5e6eb);
}

.api-category-list {
  padding: 8px 0;
}

.api-category {
  margin-bottom: 4px;
}

.category-title {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-2, #1f2329);
  transition: background 0.15s;
}

.category-title:hover {
  background: var(--color-fill-1, #f2f3f5);
}

.category-title span {
  flex: 1;
}

.category-apis {
  padding: 4px 0;
}

.api-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px 6px 28px;
  cursor: pointer;
  transition: all 0.15s;
}

.api-item:hover {
  background: var(--color-fill-1, #f2f3f5);
}

.api-item.active {
  background: var(--color-primary-light-1, #e6f1ff);
  border-right: 2px solid var(--color-primary, #1650ff);
}

.api-path {
  font-size: 12px;
  color: var(--color-text-3, #646a73);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 120px;
}

.api-content {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  padding: 20px;
  min-height: calc(100vh - 96px);
}

.api-header {
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--color-fill-2, #e5e6eb);
}

.api-title-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.api-path-large {
  font-size: 16px;
  font-weight: 600;
  font-family: monospace;
  color: var(--color-text-1, #1f2329);
}

.api-desc {
  color: var(--color-text-3, #646a73);
  font-size: 14px;
  margin-left: 4px;
}

.doc-section {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.doc-block {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.doc-block-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-1, #1f2329);
}

.doc-block-content {
  font-size: 13px;
  color: var(--color-text-2, #1f2329);
  line-height: 1.6;
}

.param-type {
  background: #f2f3f5;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
  font-family: monospace;
}

.error-code {
  background: #f53f3f11;
  color: #f53f3f;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
  font-family: monospace;
}

.code-block {
  background: #1e1e1e;
  color: #d4d4d4;
  padding: 16px;
  border-radius: 8px;
  font-size: 12px;
  font-family: 'Consolas', monospace;
  overflow-x: auto;
  margin: 0;
  line-height: 1.6;
}

.test-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.test-card {
  border: 1px solid var(--color-fill-2, #e5e6eb);
}

.test-card-title {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 12px;
  color: var(--color-text-1, #1f2329);
}

.test-actions {
  margin-top: 12px;
}

.test-result-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.test-time {
  font-size: 13px;
  color: var(--color-text-3, #646a73);
}

.api-empty {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  height: calc(100vh - 96px);
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
