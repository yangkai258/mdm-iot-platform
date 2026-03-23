<template>
  <div class="family-emotion-map">
    <a-card title="家庭情绪地图">
      <a-spin :loading="loading">
        <a-row :gutter="16">
          <a-col :span="12">
            <div class="family-index">
              <a-progress type="circle" :percent="familyData.family_mood_index" :color="getIndexColor(familyData.family_mood_index)" />
              <div class="index-label">家庭情绪指数</div>
            </div>
          </a-col>
          <a-col :span="12">
            <div class="members-list">
              <div v-for="member in familyData.members" :key="member.pet_id" class="member-item">
                <a-avatar :size="36">{{ member.pet_id }}</a-avatar>
                <div class="member-info">
                  <div class="name">宠物 #{{ member.pet_id }}</div>
                  <div class="mood">
                    <a-tag :color="getEmotionColor(member.current_mood)">{{ member.current_mood }}</a-tag>
                    强度: {{ member.intensity }}
                  </div>
                </div>
              </div>
            </div>
          </a-col>
        </a-row>
        
        <a-divider>互动热点</a-divider>
        <div class="hotspots">
          <div v-for="(hs, idx) in familyData.hotspots" :key="idx" class="hotspot-item">
            宠物 #{{ hs.from_pet_id }} ↔ 宠物 #{{ hs.to_pet_id }}: {{ hs.interaction_count }} 次互动
          </div>
          <div v-if="!familyData.hotspots?.length" class="empty">暂无互动数据</div>
        </div>
      </a-spin>
    </a-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'
const loading = ref(false)
const familyData = ref({
  family_mood_index: 0,
  members: [],
  hotspots: []
})

const getIndexColor = (v) => v >= 70 ? '#52c41a' : v >= 40 ? '#faad14' : '#ff4d4f'
const getEmotionColor = (e) => ({
  happy: '#52c41a', sad: '#1890ff', angry: '#ff4d4f',
  calm: '#722ed1', excited: '#faad14', anxious: '#f5222d'
}[e] || '#1890ff')

const load = async () => {
  loading.value = true
  try {
    const res = await fetch(`${API_BASE}/emotions/family-map?family_id=1`)
    familyData.value = await res.json()
  } catch (e) {
    Message.error('加载失败')
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.family-index { text-align: center; padding: 20px; }
.index-label { margin-top: 12px; font-size: 16px; font-weight: 500; }
.members-list { display: flex; flex-direction: column; gap: 12px; }
.member-item { display: flex; align-items: center; gap: 12px; }
.member-info { flex: 1; }
.member-info .name { font-weight: 500; }
.member-info .mood { font-size: 12px; display: flex; gap: 8px; align-items: center; }
.hotspots { display: flex; flex-direction: column; gap: 8px; }
.hotspot-item { padding: 8px 12px; background: var(--color-fill-1); border-radius: 4px; }
.empty { color: var(--color-text-3); text-align: center; padding: 20px; }
</style>
