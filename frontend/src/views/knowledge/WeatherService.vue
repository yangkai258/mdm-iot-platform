<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>知识库</a-breadcrumb-item>
      <a-breadcrumb-item>天气服务</a-breadcrumb-item>
    </a-breadcrumb>

    <a-tabs v-model:active-key="activeTab" class="pro-content-area">
      <!-- 城市配置 -->
      <a-tab-pane key="cities" title="城市配置">
        <div class="pro-action-bar">
          <a-space>
            <a-button type="primary" @click="showAddCityModal">添加城市</a-button>
            <a-button @click="loadCities">刷新</a-button>
          </a-space>
        </div>
        <a-table :columns="cityColumns" :data="cities" :loading="loading" :pagination="pagination" row-key="id" @page-change="handlePageChange">
          <template #enabled="{ record }">
            <a-switch v-model="record.enabled" @change="toggleCity(record)" />
          </template>
          <template #current_weather="{ record }">
            <a-space>
              <span>{{ record.current_temp }}°C</span>
              <span style="color: var(--color-text-3)">{{ record.current_condition }}</span>
            </a-space>
          </template>
          <template #updated_at="{ record }">
            {{ formatDate(record.updated_at) }}
          </template>
          <template #actions="{ record }">
            <a-space>
              <a-button type="text" size="small" @click="refreshCityWeather(record)">刷新</a-button>
              <a-button type="text" size="small" status="danger" @click="deleteCity(record)">删除</a-button>
            </a-space>
          </template>
        </a-table>

        <!-- 添加城市弹窗 -->
        <a-modal v-model:visible="addCityVisible" title="添加城市" @ok="submitAddCity" :width="480" :loading="submitting">
          <a-form :model="cityForm" layout="vertical">
            <a-form-item label="城市名称" required>
              <a-input v-model="cityForm.city_name" placeholder="如：北京、上海、Beijing" />
            </a-form-item>
            <a-form-item label="城市代码">
              <a-input v-model="cityForm.city_code" placeholder="城市代码（可选，用于精确匹配）" />
            </a-form-item>
            <a-form-item label="国家/地区">
              <a-input v-model="cityForm.country" placeholder="如：CN、US" />
            </a-form-item>
            <a-form-item label="经度">
              <a-input-number v-model="cityForm.latitude" placeholder="纬度" style="width: 100%" />
            </a-form-item>
            <a-form-item label="经度">
              <a-input-number v-model="cityForm.longitude" placeholder="经度" style="width: 100%" />
            </a-form-item>
          </a-form>
        </a-modal>
      </a-tab-pane>

      <!-- 数据源配置 -->
      <a-tab-pane key="source" title="数据源配置">
        <a-card title="天气数据源配置">
          <a-form :model="weatherConfig" layout="vertical" ref="configFormRef">
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="数据源类型" required>
                  <a-select v-model="weatherConfig.provider" placeholder="选择数据源">
                    <a-option value="open-meteo">Open-Meteo（免费，无需API Key）</a-option>
                    <a-option value="openweather">OpenWeather</a-option>
                    <a-option value="qweather">和风天气</a-option>
                    <a-option value="weatherapi">WeatherAPI</a-option>
                  </a-select>
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="API Key">
                  <a-input-password v-model="weatherConfig.api_key" placeholder="输入API Key" />
                </a-form-item>
              </a-col>
            </a-row>
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="API Endpoint">
                  <a-input v-model="weatherConfig.endpoint" placeholder="API地址" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="请求超时（秒）">
                  <a-input-number v-model="weatherConfig.timeout" :min="1" :max="60" style="width: 100%" />
                </a-form-item>
              </a-col>
            </a-row>
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="启用HTTPS">
                  <a-switch v-model="weatherConfig.use_https" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="数据源状态">
                  <a-badge v-if="weatherConfig.enabled" status="success" text="已启用" />
                  <a-badge v-else status="error" text="已禁用" />
                </a-form-item>
              </a-col>
            </a-row>
            <a-form-item>
              <a-space>
                <a-button type="primary" @click="saveConfig" :loading="saving">保存配置</a-button>
                <a-button @click="testSource" :loading="testing">测试连接</a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </a-card>
      </a-tab-pane>

      <!-- 缓存策略 -->
      <a-tab-pane key="cache" title="缓存策略">
        <a-card title="天气数据缓存策略">
          <a-form :model="cacheConfig" layout="vertical">
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="缓存策略">
                  <a-select v-model="cacheConfig.strategy" placeholder="选择缓存策略">
                    <a-option value="no-cache">不缓存（实时请求）</a-option>
                    <a-option value="time-based">时间过期</a-option>
                    <a-option value="stale-while-revalidate">Stale-While-Revalidate</a-option>
                  </a-select>
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="缓存TTL（分钟）">
                  <a-input-number v-model="cacheConfig.ttl_minutes" :min="1" :max="1440" style="width: 100%" />
                </a-form-item>
              </a-col>
            </a-row>
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="强制刷新间隔">
                  <a-input-number v-model="cacheConfig.force_refresh_minutes" :min="1" :max="360" style="width: 100%" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="缓存最大条目">
                  <a-input-number v-model="cacheConfig.max_entries" :min="1" :max="10000" style="width: 100%" />
                </a-form-item>
              </a-col>
            </a-row>
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="当日数据过期时间">
                  <a-time-picker v-model="cacheConfig.daily_reset_time" format="HH:mm" style="width: 100%" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="允许缓存错误响应">
                  <a-switch v-model="cacheConfig.cache_errors" />
                </a-form-item>
              </a-col>
            </a-row>
            <a-divider>缓存统计</a-divider>
            <a-row :gutter="16">
              <a-col :span="8">
                <a-statistic title="缓存命中率" :value="cacheStats.hit_rate" suffix="%" :precision="1" />
              </a-col>
              <a-col :span="8">
                <a-statistic title="当前缓存条目" :value="cacheStats.current_entries" />
              </a-col>
              <a-col :span="8">
                <a-statistic title="缓存总请求" :value="cacheStats.total_requests" />
              </a-col>
            </a-row>
            <a-form-item style="margin-top: 16px">
              <a-space>
                <a-button type="primary" @click="saveCacheConfig" :loading="saving">保存缓存配置</a-button>
                <a-button status="warning" @click="clearCache">清空缓存</a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </a-card>
      </a-tab-pane>
    </a-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'

const activeTab = ref('cities')
const loading = ref(false)
const saving = ref(false)
const testing = ref(false)
const submitting = ref(false)
const addCityVisible = ref(false)
const configFormRef = ref()

const cities = ref<any[]>([])
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })

const cityForm = reactive({ city_name: '', city_code: '', country: 'CN', latitude: 0, longitude: 0 })

const weatherConfig = reactive({
  provider: 'open-meteo', api_key: '', endpoint: 'https://api.open-meteo.com/v1', timeout: 10, use_https: true, enabled: true,
})

const cacheConfig = reactive({
  strategy: 'time-based', ttl_minutes: 30, force_refresh_minutes: 5, max_entries: 1000, daily_reset_time: null as any, cache_errors: false,
})

const cacheStats = reactive({ hit_rate: 0, current_entries: 0, total_requests: 0 })

const cityColumns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '城市名称', dataIndex: 'city_name', width: 120 },
  { title: '城市代码', dataIndex: 'city_code', width: 100 },
  { title: '坐标', dataIndex: 'coordinates', width: 160 },
  { title: '当前天气', slotName: 'current_weather', width: 160 },
  { title: '启用', slotName: 'enabled', width: 80 },
  { title: '更新时间', dataIndex: 'updated_at', slotName: 'updated_at', width: 180 },
  { title: '操作', slotName: 'actions', fixed: 'right', width: 140 },
]

const loadCities = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/knowledge/weather/cities', { params: { page: pagination.current, page_size: pagination.pageSize } })
    cities.value = res.data.items || []
    pagination.total = res.data.total || 0
  } catch {
    cities.value = [
      { id: 1, city_name: '北京', city_code: '101010100', country: 'CN', coordinates: '39.9N, 116.4E', current_temp: 18, current_condition: '晴', enabled: true, updated_at: new Date().toISOString() },
      { id: 2, city_name: '上海', city_code: '101020100', country: 'CN', coordinates: '31.2N, 121.5E', current_temp: 22, current_condition: '多云', enabled: true, updated_at: new Date().toISOString() },
      { id: 3, city_name: '广州', city_code: '101280101', country: 'CN', coordinates: '23.1N, 113.3E', current_temp: 28, current_condition: '雷阵雨', enabled: false, updated_at: new Date(Date.now() - 3600000).toISOString() },
    ]
    pagination.total = 3
  } finally { loading.value = false }
}

const showAddCityModal = () => { addCityVisible.value = true }
const submitAddCity = async () => {
  submitting.value = true
  try {
    await axios.post('/api/v1/knowledge/weather/cities', cityForm)
    Message.success('添加成功')
    addCityVisible.value = false
    loadCities()
  } catch { Message.error('添加失败') } finally { submitting.value = false }
}

const toggleCity = async (record: any) => {
  try { await axios.put(`/api/v1/knowledge/weather/cities/${record.id}`, { enabled: record.enabled }); Message.success('更新成功') }
  catch { record.enabled = !record.enabled; Message.error('更新失败') }
}

const refreshCityWeather = async (record: any) => {
  try {
    await axios.post(`/api/v1/knowledge/weather/cities/${record.id}/refresh`)
    Message.success('刷新成功')
    loadCities()
  } catch { Message.error('刷新失败') }
}

const deleteCity = async (record: any) => {
  try { await axios.delete(`/api/v1/knowledge/weather/cities/${record.id}`); Message.success('删除成功'); loadCities() }
  catch { Message.error('删除失败') }
}

const saveConfig = async () => {
  saving.value = true
  try { await axios.post('/api/v1/knowledge/weather/config', weatherConfig); Message.success('配置保存成功') }
  catch { Message.error('保存失败') } finally { saving.value = false }
}

const testSource = async () => {
  testing.value = true
  try {
    await axios.post('/api/v1/knowledge/weather/test', weatherConfig)
    Message.success('连接测试成功')
  } catch (e: any) { Message.error('测试失败: ' + (e.response?.data?.message || e.message)) }
  finally { testing.value = false }
}

const saveCacheConfig = async () => {
  saving.value = true
  try { await axios.post('/api/v1/knowledge/weather/cache-config', cacheConfig); Message.success('缓存配置保存成功') }
  catch { Message.error('保存失败') } finally { saving.value = false }
}

const clearCache = async () => {
  try { await axios.delete('/api/v1/knowledge/weather/cache'); Message.success('缓存已清空') }
  catch { Message.error('清空失败') }
}

const handlePageChange = (page: number) => { pagination.current = page; loadCities() }
const formatDate = (d: string) => d ? new Date(d).toLocaleString('zh-CN') : '-'

onMounted(() => { loadCities() })
</script>
