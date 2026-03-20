<template>
  <div class="firmware-list">
    <a-layout-header class="header">
      <div class="header-title">MiniClaw固件</div>
    </a-layout-header>
    <a-layout-content class="content">
      <a-card>
        <template #title>
          <a-space>
            <span>固件列表</span>
            <a-button type="primary" @click="uploadFirmware">上传固件</a-button>
          </a-space>
        </template>
        <a-table :columns="columns" :data="firmwares" row-key="id">
          <template #actions="{ record }">
            <a-space>
              <a-button size="small" @click="viewDetail(record)">详情</a-button>
              <a-button size="small" type="primary" @click="deploy(record)">下发</a-button>
              <a-button size="small" type="primary" danger @click="deleteFirmware(record)">删除</a-button>
            </a-space>
          </template>
        </a-table>
      </a-card>
    </a-layout-content>
  </div>
</template>

<script>
export default {
  name: 'FirmwareList',
  data() {
    return {
      columns: [
        { title: 'ID', dataIndex: 'id' },
        { title: '版本', dataIndex: 'version' },
        { title: '文件大小', dataIndex: 'size' },
        { title: 'MD5', dataIndex: 'md5' },
        { title: '创建时间', dataIndex: 'created_at' },
        { title: '操作', slotName: 'actions' }
      ],
      firmwares: []
    }
  },
  methods: {
    uploadFirmware() { this.$message.info('上传固件') },
    viewDetail(record) { this.$message.info('查看详情') },
    deploy(record) { this.$message.success('下发固件') },
    deleteFirmware(record) { this.$message.warn('删除固件') }
  }
}
</script>
