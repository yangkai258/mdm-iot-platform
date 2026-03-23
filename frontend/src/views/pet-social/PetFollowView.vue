<template>
  <div class="pet-follow">
    <a-tabs>
      <a-tab-pane key="following" title="关注">
        <a-list :data="following" :loading="loading">
          <template #item="{ item }">
            <a-list-item>
              <a-avatar>{{ item.followee_id }}</a-avatar>
              <span style="margin-left: 12px">{{ item.follow_type }} #{{ item.followee_id }}</span>
              <template #actions>
                <a-button size="small" type="text" @click="unfollow(item)">取消关注</a-button>
              </template>
            </a-list-item>
          </template>
        </a-list>
      </a-tab-pane>
      <a-tab-pane key="followers" title="粉丝">
        <a-list :data="followers" :loading="loading2">
          <template #item="{ item }">
            <a-list-item>
              <a-avatar>{{ item.follower_id }}</a-avatar>
              <span style="margin-left: 12px">{{ item.follow_type }} #{{ item.follower_id }}</span>
            </a-list-item>
          </template>
        </a-list>
      </a-tab-pane>
    </a-tabs>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'
const loading = ref(false)
const loading2 = ref(false)
const following = ref([])
const followers = ref([])

const loadFollowing = async () => {
  loading.value = true
  const res = await fetch(`${API_BASE}/pet-social/following?user_id=1`)
  following.value = await res.json()
  loading.value = false
}

const loadFollowers = async () => {
  loading2.value = true
  const res = await fetch(`${API_BASE}/pet-social/followers?user_id=1`)
  followers.value = await res.json()
  loading2.value = false
}

const unfollow = async (item) => {
  await fetch(`${API_BASE}/pet-social/follow/${item.id}`, { method: 'DELETE' })
  Message.success('已取消关注')
  loadFollowing()
}

onMounted(() => { loadFollowing(); loadFollowers() })
</script>

<style scoped>
.pet-follow {
  padding: 16px;
}
</style>
