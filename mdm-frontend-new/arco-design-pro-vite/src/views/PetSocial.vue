<template>
  <div class="pet-social-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>宠物社区</span>
          <a-space>
            <a-button @click="handleRefresh">刷新</a-button>
            <a-button type="primary" @click="handleCreate">
              <template #icon><icon-plus /></template>
              发帖
            </a-button>
          </a-space>
        </div>
      </template>
      
      <a-tabs v-model="activeTab">
        <a-tab-pane key="feed" title="动态">
          <a-list :data-source="posts" :pagination="pagination">
            <template #item="{ item }">
              <a-list-item-meta>
                <template #avatar>
                  <a-avatar :size="48" :style="{ backgroundColor: item.avatarColor }">
                    {{ item.petName?.charAt(0) }}
                  </a-avatar>
                </template>
                <template #title>
                  <span>{{ item.petName }}</span>
                  <a-tag :color="item.petType === 'dog' ? 'blue' : 'orange'" size="small">{{ item.petType === 'dog' ? '狗狗' : '猫咪' }}</a-tag>
                  <span style="float: right; color: #86909c; font-size: 12px;">{{ item.time }}</span>
                </template>
                <template #description>
                  <div class="post-content">{{ item.content }}</div>
                  <div v-if="item.images?.length" class="post-images">
                    <a-image-preview-group>
                      <a-image v-for="(img, i) in item.images" :key="i" :src="img" :width="100" />
                    </a-image-preview-group>
                  </div>
                  <div class="post-actions">
                    <a-space>
                      <a-link @click="handleLike(item)">
                        <icon-thumb-up :size="14" /> {{ item.likes }}
                      </a-link>
                      <a-link @click="handleComment(item)">
                        <icon-message :size="14" /> {{ item.comments }}
                      </a-link>
                      <a-link><icon-share :size="14" /> 分享</a-link>
                    </a-space>
                  </div>
                </template>
              </a-list-item-meta>
            </template>
          </a-list>
        </a-tab-pane>
        
        <a-tab-pane key="nearby" title="附近宠物">
          <a-row :gutter="16">
            <a-col :span="6" v-for="pet in nearbyPets" :key="pet.id">
              <a-card size="small" class="pet-card">
                <a-avatar :size="80" :style="{ backgroundColor: pet.avatarColor }">
                  {{ pet.name?.charAt(0) }}
                </a-avatar>
                <div class="pet-name">{{ pet.name }}</div>
                <div class="pet-breed">{{ pet.breed }}</div>
                <div class="pet-distance">{{ pet.distance }}m</div>
                <a-button size="small" type="primary" @click="handleInvitePlay(pet)">约玩</a-button>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
        
        <a-tab-pane key="lost" title="寻回网络">
          <a-alert type="warning" style="margin-bottom: 16px;">
            有3只宠物正在寻回中，查看详情或帮助提供线索
          </a-alert>
          <a-list :data-source="lostPets" :pagination="paginationSmall">
            <template #item="{ item }">
              <a-card size="small" class="lost-card">
                <a-row :gutter="16">
                  <a-col :span="4">
                    <a-avatar :size="80" :style="{ backgroundColor: item.avatarColor }">
                      {{ item.name?.charAt(0) }}
                    </a-avatar>
                  </a-col>
                  <a-col :span="20">
                    <a-descriptions :column="2" size="small">
                      <a-descriptions-item label="宠物名">{{ item.name }}</a-descriptions-item>
                      <a-descriptions-item label="品种">{{ item.breed }}</a-descriptions-item>
                      <a-descriptions-item label="丢失地点">{{ item.lostLocation }}</a-descriptions-item>
                      <a-descriptions-item label="丢失时间">{{ item.lostTime }}</a-descriptions-item>
                      <a-descriptions-item label="联系方式">{{ item.contact }}</a-descriptions-item>
                      <a-descriptions-item label="悬赏">{{ item.reward }}</a-descriptions-item>
                    </a-descriptions>
                    <a-button type="primary" size="small" style="margin-top: 8px;">提供线索</a-button>
                  </a-col>
                </a-row>
              </a-card>
            </template>
          </a-list>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <!-- 发帖弹窗 -->
    <a-modal v-model:visible="createVisible" title="发布动态" @before-ok="handleSubmit">
      <a-form :model="postForm" layout="vertical">
        <a-form-item label="选择宠物">
          <a-select v-model="postForm.petId" placeholder="选择要发帖的宠物">
            <a-option value="P001">小黄 - 金毛</a-option>
            <a-option value="P002">小红 - 柯基</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="内容">
          <a-textarea v-model="postForm.content" placeholder="分享你家宠物的精彩瞬间..." :rows="4" />
        </a-form-item>
        <a-form-item label="上传图片">
          <a-upload action="#" :limit="9" multiple />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const activeTab = ref('feed');

const posts = ref([
  { id: 1, petName: '小黄', petType: 'dog', avatarColor: '#165DFF', content: '今天天气真好，和主人去公园玩飞盘啦！', images: [], likes: 128, comments: 23, time: '2小时前' },
  { id: 2, petName: '小红', petType: 'dog', avatarColor: '#00B42A', content: '第一次吃冰淇淋，太好吃了！', images: [], likes: 256, comments: 45, time: '4小时前' },
  { id: 3, petName: '咪咪', petType: 'cat', avatarColor: '#FF7D00', content: '晒太阳真舒服~', images: [], likes: 89, comments: 12, time: '6小时前' },
]);

const nearbyPets = ref([
  { id: 1, name: '旺财', breed: '柴犬', avatarColor: '#FF7D00', distance: 320 },
  { id: 2, name: '球球', breed: '柯基', avatarColor: '#00B42A', distance: 580 },
  { id: 3, name: '小白', breed: '萨摩耶', avatarColor: '#722ED1', distance: 1200 },
  { id: 4, name: '花花', breed: '边牧', avatarColor: '#165DFF', distance: 2300 },
]);

const lostPets = ref([
  { id: 1, name: '豆豆', breed: '金毛', lostLocation: '朝阳区公园', lostTime: '2026-03-25', contact: '138****1234', reward: '5000元', avatarColor: '#165DFF' },
  { id: 2, name: '小虎', breed: '橘猫', lostLocation: '海淀区小区', lostTime: '2026-03-26', contact: '139****5678', reward: '2000元', avatarColor: '#FF7D00' },
  { id: 3, name: '贝贝', breed: '比熊', lostLocation: '东城区商场', lostTime: '2026-03-27', contact: '137****9012', reward: '3000元', avatarColor: '#00B42A' },
]);

const pagination = reactive({ current: 1, pageSize: 10, total: 3 });
const paginationSmall = reactive({ current: 1, pageSize: 5, total: 3 });
const createVisible = ref(false);

const postForm = reactive({ petId: '', content: '', images: [] });

const handleRefresh = () => {};
const handleCreate = () => { createVisible.value = true; };
const handleLike = (item: any) => { item.likes++; };
const handleComment = (item: any) => {};
const handleInvitePlay = (pet: any) => {};
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.pet-social-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
.post-content { margin: 8px 0; }
.post-images { margin-top: 8px; }
.post-actions { margin-top: 8px; color: #86909c; }
.pet-card { text-align: center; }
.pet-name { font-weight: bold; margin-top: 8px; }
.pet-breed { color: #86909c; font-size: 12px; }
.pet-distance { color: #165DFF; font-size: 12px; margin: 4px 0; }
.lost-card { margin-bottom: 8px; }
</style>
