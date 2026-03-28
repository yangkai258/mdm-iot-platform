<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-map-location /> 空间认知地图</a-space>
      </template>
      <template #extra>
        <a-space>
          <a-button @click="handleEdit">
            <template #icon><icon-edit /></template>
            编辑地图
          </a-button>
          <a-button @click="handleFullscreen">
            <template #icon><icon-fullscreen /></template>
            全屏
          </a-button>
          <a-button type="primary" @click="handleSaveMap">
            <template #icon><icon-save /></template>
            保存地图
          </a-button>
        </a-space>
      </template>

      <a-row :gutter="16">
        <a-col :span="18">
          <a-card size="small">
            <div class="map-container">
              <svg width="800" height="500" viewBox="0 0 800 500">
                <rect x="0" y="0" width="800" height="500" fill="#e8e8e8" />
                <rect x="50" y="50" width="200" height="150" fill="#b3d9ff" stroke="#333" />
                <text x="150" y="130" text-anchor="middle">客厅</text>
                <rect x="270" y="50" width="150" height="150" fill="#d9ffb3" stroke="#333" />
                <text x="345" y="130" text-anchor="middle">厨房</text>
                <rect x="50" y="220" width="200" height="200" fill="#ffe6b3" stroke="#333" />
                <text x="150" y="330" text-anchor="middle">卧室</text>
                <rect x="270" y="220" width="150" height="200" fill="#ffb3d9" stroke="#333" />
                <text x="345" y="330" text-anchor="middle">书房</text>
                <circle cx="400" cy="250" r="15" fill="#409EFF" />
                <text x="400" y="280" text-anchor="middle" font-size="10">当前位置</text>
                <circle v-for="poi in pois" :key="poi.id" :cx="poi.x" :cy="poi.y" r="8" :fill="getPoiColor(poi.type)" />
              </svg>
            </div>
          </a-card>
        </a-col>

        <a-col :span="6">
          <a-card title="地图配置">
            <a-form :model="configForm" layout="vertical">
              <a-form-item label="地图名称">
                <a-input v-model="configForm.name" />
              </a-form-item>
              <a-form-item label="地图缩放">
                <a-slider v-model="configForm.zoom" :min="0.5" :max="2" :step="0.1" />
              </a-form-item>
            </a-form>
          </a-card>

          <a-card title="兴趣点" style="margin-top: 16px">
            <a-space direction="vertical" fill>
              <a-tag v-for="poi in pois" :key="poi.id" :color="getPoiColor(poi.type)" :closable>
                {{ poi.name }}
              </a-tag>
              <a-button type="text" @click="handleAddPoi">
                <template #icon><icon-plus /></template>
                添加兴趣点
              </a-button>
            </a-space>
          </a-card>

          <a-card title="安全区/危险区" style="margin-top: 16px">
            <a-space direction="vertical" fill>
              <a-tag v-for="zone in zones" :key="zone.id" :color="zone.type === 'safe' ? 'green' : 'red'" :closable>
                {{ zone.name }}
              </a-tag>
              <a-button type="text" @click="handleAddZone">
                <template #icon><icon-plus /></template>
                添加区域
              </a-button>
            </a-space>
          </a-card>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const configForm = reactive({ name: '我家', zoom: 1 })

const pois = ref([
  { id: 1, name: '食盆', x: 300, y: 100, type: 'feeding' },
  { id: 2, name: '水盆', x: 350, y: 100, type: 'feeding' },
  { id: 3, name: '窝', x: 150, y: 350, type: 'rest' }
])

const zones = ref([
  { id: 1, name: '客厅', type: 'safe' },
  { id: 2, name: '厨房', type: 'danger' }
])

const getPoiColor = (type) => ({ feeding: 'orange', rest: 'blue', play: 'green' }[type] || 'gray')

const handleEdit = () => { }
const handleFullscreen = () => { }
const handleSaveMap = () => { }
const handleAddPoi = () => { }
const handleAddZone = () => { }
</script>

<style scoped>
.container { padding: 16px; }
.map-container { background: #f0f0f0; border-radius: 8px; overflow: hidden; }
</style>
