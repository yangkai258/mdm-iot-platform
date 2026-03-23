<template>
  <div class="pet-playdate">
    <a-card title="宠物约会">
      <template #extra>
        <a-button type="primary" @click="showModal = true">创建约会</a-button>
      </template>
      
      <a-list :data="playdates" :loading="loading">
        <template #item="{ item }">
          <a-list-item>
            <div class="playdate-info">
              <div class="title">{{ item.title }}</div>
              <div class="meta">
                <icon-location /> {{ item.location }} | 
                <icon-clock /> {{ formatTime(item.start_time) }}
              </div>
              <a-tag :color="item.status === 'confirmed' ? 'green' : 'orange'">{{ item.status }}</a-tag>
            </div>
            <template #actions>
              <a-button size="small" type="text" v-if="item.status === 'pending'" @click="join(item)">参加</a-button>
            </template>
          </a-list-item>
        </template>
      </a-list>
    </a-card>
    
    <a-modal v-model:visible="showModal" title="创建约会" @ok="create">
      <a-form :model="form" layout="vertical">
        <a-form-item label="标题" required><a-input v-model="form.title" /></a-form-item>
        <a-form-item label="地点"><a-input v-model="form.location" /></a-form-item>
        <a-form-item label="开始时间"><a-date-picker v-model="form.start_time" show-time /></a-form-item>
        <a-form-item label="结束时间"><a-date-picker v-model="form.end_time" show-time /></a-form-item>
        <a-form-item label="描述"><a-textarea v-model="form.description" /></a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'
const loading = ref(false)
const showModal = ref(false)
const playdates = ref([])
const form = ref({ title: '', location: '', description: '', organizer_id: 1 })

const load = async () => {
  loading.value = true
  const res = await fetch(`${API_BASE}/pet-social/playdates`)
  playdates.value = await res.json()
  loading.value = false
}

const create = async () => {
  await fetch(`${API_BASE}/pet-social/playdates`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ ...form.value, start_time: new Date(form.value.start_time).toISOString(), end_time: new Date(form.value.end_time).toISOString() })
  })
  showModal.value = false
  Message.success('创建成功')
  load()
}

const join = async (item) => {
  await fetch(`${API_BASE}/pet-social/playdates/${item.id}/join`, { method: 'POST' })
  Message.success('已参加')
  load()
}

const formatTime = (t) => new Date(t).toLocaleString()

onMounted(load)
</script>

<style scoped>
.playdate-info { flex: 1; }
.playdate-info .title { font-weight: 600; margin-bottom: 4px; }
.playdate-info .meta { font-size: 12px; color: var(--color-text-3); display: flex; gap: 8px; align-items: center; }
</style>
