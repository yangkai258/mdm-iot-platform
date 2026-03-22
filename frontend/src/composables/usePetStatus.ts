import { ref, computed } from 'vue'
import { getPetStatus, updatePetSettings } from '@/api/pet'

export function usePetStatus(deviceId) {
  const petStatus = ref({
    name: '未知宠物',
    type: 'cat',
    mood: 0,
    energy: 0,
    hunger: 0,
    is_online: false,
    last_seen: null
  })
  
  const loading = ref(false)
  const error = ref<Error | null>(null)

  const isOnline = computed(() => petStatus.value.is_online)

  const moodLevel = computed(() => {
    const mood = petStatus.value.mood
    if (mood >= 70) return 'happy'
    if (mood >= 40) return 'normal'
    return 'sad'
  })

  async function fetchStatus() {
    if (!deviceId.value) return
    
    loading.value = true
    error.value = null
    
    try {
      const res = await getPetStatus(deviceId.value)
      if (res.data) {
        petStatus.value = res.data
      }
    } catch (e) {
      error.value = e as Error
      console.error('Failed to fetch pet status:', e)
    } finally {
      loading.value = false
    }
  }

  async function saveSettings(settings) {
    if (!deviceId.value) return

    loading.value = true
    error.value = null
    
    try {
      await updatePetSettings(deviceId.value, settings)
      await fetchStatus()
    } catch (e) {
      error.value = e as Error
      console.error('Failed to save settings:', e)
      throw e
    } finally {
      loading.value = false
    }
  }

  function updateStatus(data) {
    petStatus.value = { ...petStatus.value, ...data }
  }

  return {
    petStatus,
    loading,
    error,
    isOnline,
    moodLevel,
    fetchStatus,
    saveSettings,
    updateStatus
  }
}
