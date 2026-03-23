<template>
  <div class="pet-feed">
    <a-card title="宠物动态">
      <template #extra>
        <a-button type="primary" @click="showPostModal = true">发布动态</a-button>
      </template>
      
      <a-spin :loading="loading">
        <div class="feed-list">
          <a-card v-for="post in posts" :key="post.id" class="post-card" hoverable>
            <div class="post-header">
              <a-avatar :size="40">{{ post.pet_id }}</a-avatar>
              <div class="post-info">
                <div class="pet-name">宠物 #{{ post.pet_id }}</div>
                <div class="post-time">{{ formatTime(post.created_at) }}</div>
              </div>
              <a-tag :color="post.post_type === 'milestone' ? 'gold' : 'default'">
                {{ post.post_type }}
              </a-tag>
            </div>
            
            <div class="post-content">{{ post.content }}</div>
            
            <div class="post-media" v-if="post.media_urls">
              <img v-for="(url, idx) in parseMedia(post.media_urls)" :key="idx" :src="url" class="media-img" />
            </div>
            
            <div class="post-actions">
              <a-space>
                <a-button size="small" @click="likePost(post)">
                  <icon-thumb-up /> {{ post.like_count }}
                </a-button>
                <a-button size="small" @click="viewComments(post)">
                  <icon-message /> {{ post.comment_count }}
                </a-button>
              </a-space>
            </div>
          </a-card>
        </div>
        
        <div class="load-more" v-if="hasMore">
          <a-button @click="loadMore">加载更多</a-button>
        </div>
      </a-spin>
    </a-card>
    
    <a-modal v-model:visible="showPostModal" title="发布动态" @ok="submitPost">
      <a-form :model="postForm" layout="vertical">
        <a-form-item label="内容" required>
          <a-textarea v-model="postForm.content" placeholder="分享你家宠物的动态..." :max-length="500" show-word-limit />
        </a-form-item>
        <a-form-item label="类型">
          <a-select v-model="postForm.post_type">
            <a-option value="photo">照片</a-option>
            <a-option value="video">视频</a-option>
            <a-option value="milestone">里程碑</a-option>
            <a-option value="achievement">成就</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const posts = ref([])
const hasMore = ref(false)
const showPostModal = ref(false)
const postForm = ref({
  content: '',
  post_type: 'photo',
  pet_id: 1
})
const page = ref(1)

const API_BASE = '/api/v1'

const loadPosts = async () => {
  loading.value = true
  try {
    const res = await fetch(`${API_BASE}/pet-social/posts?page=${page.value}`)
    const data = await res.json()
    posts.value = page.value === 1 ? data.posts : [...posts.value, ...data.posts]
    hasMore.value = data.has_more || false
  } catch (e) {
    Message.error('加载失败')
  } finally {
    loading.value = false
  }
}

const loadMore = () => {
  page.value++
  loadPosts()
}

const likePost = async (post) => {
  await fetch(`${API_BASE}/pet-social/posts/${post.id}/like?user_id=1`, { method: 'POST' })
  post.like_count++
}

const viewComments = (post) => {
  // TODO: 实现评论功能
  console.log('view comments', post.id)
}

const submitPost = async () => {
  await fetch(`${API_BASE}/pet-social/posts`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ ...postForm.value, author_id: 1 })
  })
  showPostModal.value = false
  Message.success('发布成功')
  loadPosts()
}

const formatTime = (t) => new Date(t).toLocaleString()
const parseMedia = (urls) => { try { return JSON.parse(urls) } catch { return [] } }

onMounted(loadPosts)
</script>

<style scoped>
.post-card { margin-bottom: 16px; }
.post-header { display: flex; align-items: center; gap: 12px; margin-bottom: 12px; }
.post-info { flex: 1; }
.pet-name { font-weight: 600; }
.post-time { font-size: 12px; color: var(--color-text-3); }
.post-content { margin: 12px 0; }
.post-media { display: flex; gap: 8px; flex-wrap: wrap; }
.media-img { width: 120px; height: 120px; object-fit: cover; border-radius: 4px; }
.post-actions { margin-top: 12px; border-top: 1px solid var(--color-fill-3); padding-top: 12px; }
</style>
