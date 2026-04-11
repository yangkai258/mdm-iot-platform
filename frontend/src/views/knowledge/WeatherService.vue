<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>棣栭〉</a-breadcrumb-item>
      <a-breadcrumb-item>鐭ヨ瘑搴?/a-breadcrumb-item>
      <a-breadcrumb-item>澶╂皵鏈嶅姟</a-breadcrumb-item>
    </a-breadcrumb>

    <a-tabs v-model:active-key="activeTab" class="pro-content-area">
      <!-- 鍩庡競閰嶇疆 -->
      <a-tab-pane key="cities" title="鍩庡競閰嶇疆">
        <div class="pro-action-bar">
          <a-space>
            <a-button type="primary" @click="showAddCityModal">娣诲姞鍩庡競</a-button>
            <a-button @click="loadCities">鍒锋柊</a-button>
          </a-space>
        </div>
        <a-table :columns="cityColumns" :data="cities" :loading="loading" :pagination="pagination" row-key="id" @page-change="handlePageChange">
          <template #enabled="{ record }">
            <a-switch v-model="record.enabled" @change="toggleCity(record)" />
          </template>
          <template #current_weather="{ record }">
            <a-space>
              <span>{{ record.current_temp }}掳C</span>
              <span style="color: var(--color-text-3)">{{ record.current_condition }}</span>
            </a-space>
          </template>
          <template #updated_at="{ record }">
            {{ formatDate(record.updated_at) }}
          </template>
          <template #actions="{ record }">
            <a-space>
              <a-button type="text" size="small" @click="refreshCityWeather(record)">鍒锋柊</a-button>
              <a-button type="text" size="small" status="danger" @click="deleteCity(record)">鍒犻櫎</a-button>
            </a-space>
          </template>
        </a-table>

        <!-- 娣诲姞鍩庡競寮圭獥 -->
        <a-modal v-model:visible="addCityVisible" title="娣诲姞鍩庡競" @ok="submitAddCity" :width="480" :loading="submitting">
          <a-form :model="cityForm" layout="vertical">
            <a-form-item label="鍩庡競鍚嶇О" required>
              <a-input v-model="cityForm.city_name" placeholder="濡傦細鍖椾含銆佷笂娴枫€丅eijing" />
            </a-form-item>
            <a-form-item label="鍩庡競浠ｇ爜">
              <a-input v-model="cityForm.city_code" placeholder="鍩庡競浠ｇ爜锛堝彲閫夛紝鐢ㄤ簬绮剧‘鍖归厤锛" />
            </a-form-item>
            <a-form-item label="鍥藉/鍦板尯">
              <a-input v-model="cityForm.country" placeholder="濡傦細CN銆乁S" />
            </a-form-item>
            <a-form-item label="缁忓害">
              <a-input-number v-model="cityForm.latitude" placeholder="绾害" style="width: 100%" />
            </a-form-item>
            <a-form-item label="缁忓害">
              <a-input-number v-model="cityForm.longitude" placeholder="缁忓害" style="width: 100%" />
            </a-form-item>
          </a-form>
        </a-modal>
      </a-tab-pane>

      <!-- 鏁版嵁婧愰厤缃?-->
      <a-tab-pane key="source" title="鏁版嵁婧愰厤缃?>
        <a-card title="澶╂皵鏁版嵁婧愰厤缃?>
          <a-form :model="weatherConfig" layout="vertical" ref="configFormRef">
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="鏁版嵁婧愮被鍨? required>
                  <a-select v-model="weatherConfig.provider" placeholder="閫夋嫨鏁版嵁婧">
                    <a-option value="open-meteo">Open-Meteo锛堝厤璐癸紝鏃犻渶API Key锛?/a-option>
                    <a-option value="openweather">OpenWeather</a-option>
                    <a-option value="qweather">鍜岄澶╂皵</a-option>
                    <a-option value="weatherapi">WeatherAPI</a-option>
                  </a-select>
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="API Key">
                  <a-input-password v-model="weatherConfig.api_key" placeholder="杈撳叆API Key" />
                </a-form-item>
              </a-col>
            </a-row>
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="API Endpoint">
                  <a-input v-model="weatherConfig.endpoint" placeholder="API鍦板潃" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="璇锋眰瓒呮椂锛堢锛?>
                  <a-input-number v-model="weatherConfig.timeout" :min="1" :max="60" style="width: 100%" />
                </a-form-item>
              </a-col>
            </a-row>
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="鍚敤HTTPS">
                  <a-switch v-model="weatherConfig.use_https" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="鏁版嵁婧愮姸鎬?>
                  <a-badge v-if="weatherConfig.enabled" status="success" text="宸插惎鐢? />
                  <a-badge v-else status="error" text="宸茬鐢? />
                </a-form-item>
              </a-col>
            </a-row>
            <a-form-item>
              <a-space>
                <a-button type="primary" @click="saveConfig" :loading="saving">淇濆瓨閰嶇疆</a-button>
                <a-button @click="testSource" :loading="testing">娴嬭瘯杩炴帴</a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </a-card>
      </a-tab-pane>

      <!-- 缂撳瓨绛栫暐 -->
      <a-tab-pane key="cache" title="缂撳瓨绛栫暐">
        <a-card title="澶╂皵鏁版嵁缂撳瓨绛栫暐">
          <a-form :model="cacheConfig" layout="vertical">
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="缂撳瓨绛栫暐">
                  <a-select v-model="cacheConfig.strategy" placeholder="閫夋嫨缂撳瓨绛栫暐">
                    <a-option value="no-cache">涓嶇紦瀛橈紙瀹炴椂璇锋眰锛?/a-option>
                    <a-option value="time-based">鏃堕棿杩囨湡</a-option>
                    <a-option value="stale-while-revalidate">Stale-While-Revalidate</a-option>
                  </a-select>
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="缂撳瓨TTL锛堝垎閽燂級">
                  <a-input-number v-model="cacheConfig.ttl_minutes" :min="1" :max="1440" style="width: 100%" />
                </a-form-item>
              </a-col>
            </a-row>
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="寮哄埗鍒锋柊闂撮殧">
                  <a-input-number v-model="cacheConfig.force_refresh_minutes" :min="1" :max="360" style="width: 100%" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="缂撳瓨鏈€澶ф潯鐩?>
                  <a-input-number v-model="cacheConfig.max_entries" :min="1" :max="10000" style="width: 100%" />
                </a-form-item>
              </a-col>
            </a-row>
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="褰撴棩鏁版嵁杩囨湡鏃堕棿">
                  <a-time-picker v-model="cacheConfig.daily_reset_time" format="HH:mm" style="width: 100%" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="鍏佽缂撳瓨閿欒鍝嶅簲">
                  <a-switch v-model="cacheConfig.cache_errors" />
                </a-form-item>
              </a-col>
            </a-row>
            <a-divider>缂撳瓨缁熻</a-divider>
            <a-row :gutter="16">
              <a-col :span="8">
                <a-statistic title="缂撳瓨鍛戒腑鐜? :value="cacheStats.hit_rate" suffix="%" :precision="1" />
              </a-col>
              <a-col :span="8">
                <a-statistic title="褰撳墠缂撳瓨鏉＄洰" :value="cacheStats.current_entries" />
              </a-col>
              <a-col :span="8">
                <a-statistic title="缂撳瓨鎬昏姹? :value="cacheStats.total_requests" />
              </a-col>
            </a-row>
            <a-form-item style="margin-top: 16px">
              <a-space>
                <a-button type="primary" @click="saveCacheConfig" :loading="saving">淇濆瓨缂撳瓨閰嶇疆</a-button>
                <a-button status="warning" @click="clearCache">娓呯┖缂撳瓨</a-button>
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
  { title: '鍩庡競鍚嶇О', dataIndex: 'city_name', width: 120 },
  { title: '鍩庡競浠ｇ爜', dataIndex: 'city_code', width: 100 },
  { title: '鍧愭爣', dataIndex: 'coordinates', width: 160 },
  { title: '褰撳墠澶╂皵', slotName: 'current_weather', width: 160 },
  { title: '鍚敤', slotName: 'enabled', width: 80 },
  { title: '鏇存柊鏃堕棿', dataIndex: 'updated_at', slotName: 'updated_at', width: 180 },
  { title: '鎿嶄綔', slotName: 'actions', fixed: 'right', width: 140 },
]

const loadCities = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/knowledge/weather/cities', { params: { page: pagination.current, page_size: pagination.pageSize } })
    cities.value = res.data.items || []
    pagination.total = res.data.total || 0
  } catch {
    cities.value = [
      { id: 1, city_name: '鍖椾含', city_code: '101010100', country: 'CN', coordinates: '39.9N, 116.4E', current_temp: 18, current_condition: '鏅?, enabled: true, updated_at: new Date().toISOString() },
      { id: 2, city_name: '涓婃捣', city_code: '101020100', country: 'CN', coordinates: '31.2N, 121.5E', current_temp: 22, current_condition: '澶氫簯', enabled: true, updated_at: new Date().toISOString() },
      { id: 3, city_name: '骞垮窞', city_code: '101280101', country: 'CN', coordinates: '23.1N, 113.3E', current_temp: 28, current_condition: '闆烽樀闆?, enabled: false, updated_at: new Date(Date.now() - 3600000).toISOString() },
    ]
    pagination.total = 3
  } finally { loading.value = false }
}

const showAddCityModal = () => { addCityVisible.value = true }
const submitAddCity = async () => {
  submitting.value = true
  try {
    await axios.post('/api/v1/knowledge/weather/cities', cityForm)
    Message.success('娣诲姞鎴愬姛')
    addCityVisible.value = false
    loadCities()
  } catch { Message.error('娣诲姞澶辫触') } finally { submitting.value = false }
}

const toggleCity = async (record: any) => {
  try { await axios.put(`/api/v1/knowledge/weather/cities/${record.id}`, { enabled: record.enabled }); Message.success('鏇存柊鎴愬姛') }
  catch { record.enabled = !record.enabled; Message.error('鏇存柊澶辫触') }
}

const refreshCityWeather = async (record: any) => {
  try {
    await axios.post(`/api/v1/knowledge/weather/cities/${record.id}/refresh`)
    Message.success('鍒锋柊鎴愬姛')
    loadCities()
  } catch { Message.error('鍒锋柊澶辫触') }
}

const deleteCity = async (record: any) => {
  try { await axios.delete(`/api/v1/knowledge/weather/cities/${record.id}`); Message.success('鍒犻櫎鎴愬姛'); loadCities() }
  catch { Message.error('鍒犻櫎澶辫触') }
}

const saveConfig = async () => {
  saving.value = true
  try { await axios.post('/api/v1/knowledge/weather/config', weatherConfig); Message.success('閰嶇疆淇濆瓨鎴愬姛') }
  catch { Message.error('淇濆瓨澶辫触') } finally { saving.value = false }
}

const testSource = async () => {
  testing.value = true
  try {
    await axios.post('/api/v1/knowledge/weather/test', weatherConfig)
    Message.success('杩炴帴娴嬭瘯鎴愬姛')
  } catch (e: any) { Message.error('娴嬭瘯澶辫触: ' + (e.response?.data?.message || e.message)) }
  finally { testing.value = false }
}

const saveCacheConfig = async () => {
  saving.value = true
  try { await axios.post('/api/v1/knowledge/weather/cache-config', cacheConfig); Message.success('缂撳瓨閰嶇疆淇濆瓨鎴愬姛') }
  catch { Message.error('淇濆瓨澶辫触') } finally { saving.value = false }
}

const clearCache = async () => {
  try { await axios.delete('/api/v1/knowledge/weather/cache'); Message.success('缂撳瓨宸叉竻绌?) }
  catch { Message.error('娓呯┖澶辫触') }
}

const handlePageChange = (page: number) => { pagination.current = page; loadCities() }
const formatDate = (d: string) => d ? new Date(d).toLocaleString('zh-CN') : '-'

onMounted(() => { loadCities() })
</script>
