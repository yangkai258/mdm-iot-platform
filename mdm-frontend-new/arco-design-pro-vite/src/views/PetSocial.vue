<template>
  <div class="container">
    <a-card>
      <template #title><icon-share /> 宠物社交</template>
      <a-button type="primary" @click="handlePost">发布动态</a-button>
      <a-list style="margin-top: 16px">
        <a-list-item v-for="post in posts" :key="post.id">
          <a-list-item-meta :title="post.petName" :description="post.content" />
          <template #actions>
            <a-link @click="handleLike(post)">点赞 ({{ post.likes }})</a-link>
            <a-link @click="handleComment(post)">评论</a-link>
          </template>
        </a-list-item>
      </a-list>
    </a-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const posts = ref([
  { id: 1, petName: '小白', content: '今天天气真好！', likes: 12 }
])
const handlePost = () => {}
const handleLike = (p) => {}
const handleComment = (p) => {}
</script>
<style scoped>.container { padding: 16px; }</style>
