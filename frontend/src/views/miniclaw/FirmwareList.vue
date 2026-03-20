<template>
  <div class="firmware-list">
    
    <div class="page-content"><a-card>
        <template #title>
          <a-space>
            <span>鍥轰欢鍒楄〃</span>
            <a-button type="primary" @click="uploadFirmware">涓婁紶鍥轰欢</a-button>
          </a-space>
        </template>
        <a-table :columns="columns" :data="firmwares" row-key="id">
          <template #actions="{ record }">
            <a-space>
              <a-button size="small" @click="viewDetail(record)">璇︽儏</a-button>
              <a-button size="small" type="primary" @click="deploy(record)">涓嬪彂</a-button>
              <a-button size="small" type="primary" danger @click="deleteFirmware(record)">鍒犻櫎</a-button>
            </a-space>
          </template>
        </a-table>
      </a-card></div>
  </div>
</template>

<script>
export default {
  name: 'FirmwareList',
  data() {
    return {
      columns: [
        { title: 'ID', dataIndex: 'id' },
        { title: '鐗堟湰', dataIndex: 'version' },
        { title: '鏂囦欢澶у皬', dataIndex: 'size' },
        { title: 'MD5', dataIndex: 'md5' },
        { title: '鍒涘缓鏃堕棿', dataIndex: 'created_at' },
        { title: '鎿嶄綔', slotName: 'actions' }
      ],
      firmwares: []
    }
  },
  methods: {
    uploadFirmware() { this.$message.info('涓婁紶鍥轰欢') },
    viewDetail(record) { this.$message.info('鏌ョ湅璇︽儏') },
    deploy(record) { this.$message.success('涓嬪彂鍥轰欢') },
    deleteFirmware(record) { this.$message.warn('鍒犻櫎鍥轰欢') }
  }
}
</script>
