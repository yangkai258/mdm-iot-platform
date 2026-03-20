<template>
  <div class="owner-profile">
    
    <div class="page-content"><a-row :gutter="16">
        <a-col :span="12">
          <a-card title="鍩烘湰淇℃伅">
            <a-form :model="profile" layout="vertical">
              <a-form-item label="鏄电О">
                <a-input v-model="profile.nickname" />
              </a-form-item>
              <a-form-item label="绉板懠鍋忓ソ">
                <a-input v-model="profile.preferred_name" placeholder="濡傦細涓讳汉銆佽€佹澘" />
              </a-form-item>
              <a-form-item label="鑱旂郴鏂瑰紡">
                <a-input v-model="profile.contact" />
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="saveProfile">淇濆瓨</a-button>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card title="鍋忓ソ璁剧疆">
            <a-form :model="preferences" layout="vertical">
              <a-form-item label="甯哥敤鏃堕棿娈?>
                <a-select v-model="preferences.active_hours" mode="multiple" placeholder="閫夋嫨娲昏穬鏃堕棿娈?>
                  <a-option value="morning">鏃╀笂</a-option>
                  <a-option value="afternoon">涓嬪崍</a-option>
                  <a-option value="evening">鏅氫笂</a-option>
                </a-select>
              </a-form-item>
              <a-form-item label="鍏磋叮鏍囩">
                <a-select v-model="preferences.interests" mode="tags" placeholder="娣诲姞鍏磋叮鏍囩">
                </a-select>
              </a-form-item>
              <a-form-item label="浜や簰椋庢牸">
                <a-radio-group v-model="preferences.interaction_style">
                  <a-radio value="formal">姝ｅ紡</a-radio>
                  <a-radio value="casual">杞绘澗</a-radio>
                  <a-radio value="humorous">骞介粯</a-radio>
                </a-radio-group>
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="savePreferences">淇濆瓨鍋忓ソ</a-button>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>
      </a-row></div>
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
      // TODO: 璋冪敤鍚庣API
    },
    async saveProfile() {
      this.$message.success('淇濆瓨鎴愬姛')
    },
    async savePreferences() {
      this.$message.success('鍋忓ソ淇濆瓨鎴愬姛')
    }
  }
}
</script>
