<template>
  <div class="owner-profile">
    <a-layout-header class="header">
      <div class="header-title">主人画像</div>
    </a-layout-header>
    <a-layout-content class="content">
      <a-row :gutter="16">
        <a-col :span="12">
          <a-card title="基本信息">
            <a-form :model="profile" layout="vertical">
              <a-form-item label="昵称">
                <a-input v-model="profile.nickname" />
              </a-form-item>
              <a-form-item label="称呼偏好">
                <a-input v-model="profile.preferred_name" placeholder="如：主人、老板" />
              </a-form-item>
              <a-form-item label="联系方式">
                <a-input v-model="profile.contact" />
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="saveProfile">保存</a-button>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card title="偏好设置">
            <a-form :model="preferences" layout="vertical">
              <a-form-item label="常用时间段">
                <a-select v-model="preferences.active_hours" mode="multiple" placeholder="选择活跃时间段">
                  <a-option value="morning">早上</a-option>
                  <a-option value="afternoon">下午</a-option>
                  <a-option value="evening">晚上</a-option>
                </a-select>
              </a-form-item>
              <a-form-item label="兴趣标签">
                <a-select v-model="preferences.interests" mode="tags" placeholder="添加兴趣标签">
                </a-select>
              </a-form-item>
              <a-form-item label="交互风格">
                <a-radio-group v-model="preferences.interaction_style">
                  <a-radio value="formal">正式</a-radio>
                  <a-radio value="casual">轻松</a-radio>
                  <a-radio value="humorous">幽默</a-radio>
                </a-radio-group>
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="savePreferences">保存偏好</a-button>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>
      </a-row>
    </a-layout-content>
  </div>
</template>

<script>
export default {
  name: 'OwnerProfile',
  data() {
    return {
      profile: {
        nickname: '',
        preferred_name: '',
        contact: ''
      },
      preferences: {
        active_hours: [],
        interests: [],
        interaction_style: 'casual'
      }
    }
  },
  mounted() {
    this.loadProfile()
  },
  methods: {
    async loadProfile() {
      // TODO: 调用后端API
    },
    async saveProfile() {
      this.$message.success('保存成功')
    },
    async savePreferences() {
      this.$message.success('偏好保存成功')
    }
  }
}
</script>
