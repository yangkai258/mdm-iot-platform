<template>
  <div class="smart-home-view">
    <a-card title="智能家居">
      <a-alert type="info" style="margin-bottom: 16px">
        智能家居设备集成控制中心，可与宠物电子宠物联动，实现环境自动调节。
      </a-alert>

      <a-row :gutter="16" style="margin-bottom: 24px">
        <a-col :span="8">
          <a-statistic title="绑定设备" :value="stats.devices" />
        </a-col>
        <a-col :span="8">
          <a-statistic title="在线设备" :value="stats.online" />
        </a-col>
        <a-col :span="8">
          <a-statistic title="自动化场景" :value="stats.scenes" />
        </a-col>
      </a-row>

      <a-tabs>
        <a-tab-pane key="devices" title="设备列表">
          <a-space style="margin-bottom: 12px">
            <a-button type="primary" @click="showAddDevice = true">添加设备</a-button>
            <a-select v-model="deviceType" placeholder="设备类型" style="width: 140px" allow-clear>
              <a-option value="light">灯光</a-option>
              <a-option value="ac">空调</a-option>
              <a-option value="humidifier">加湿器</a-option>
              <a-option value="speaker">音响</a-option>
            </a-select>
          </a-space>
          <a-row :gutter="16">
            <a-col :span="6" v-for="device in filteredDevices" :key="device.id">
              <a-card class="device-card" hoverable>
                <template #cover>
                  <div class="device-icon">
                    <icon-home v-if="device.type === 'light'" />
                    <icon-sun v-else-if="device.type === 'ac'" />
                    <icon-water v-else-if="device.type === 'humidifier'" />
                    <icon-music v-else />
                  </div>
                </template>
                <a-card-meta :title="device.name" :description="device.room" />
                <div style="margin-top: 8px">
                  <a-switch v-model="device.status" @change="toggleDevice(device)" />
                </div>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>

        <a-tab-pane key="scenes" title="自动化场景">
          <a-space style="margin-bottom: 12px">
            <a-button type="primary" @click="showAddScene = true">新建场景</a-button>
          </a-space>
          <a-table :columns="sceneColumns" :data="scenes" :loading="loading">
            <template #status="{ record }">
              <a-tag :color="record.enabled ? 'green' : 'gray'">{{ record.enabled ? '已启用' : '已禁用' }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="editScene(record)">编辑</a-button>
                <a-button type="text" size="small" status="danger" @click="deleteScene(record)">删除</a-button>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>

        <a-tab-pane key="linkage" title="宠物联动">
          <a-form layout="vertical">
            <a-form-item label="宠物情绪 → 环境调节">
              <a-textarea v-model="linkageConfig" :rows="6" placeholder="配置宠物情绪与家居设备的联动规则" />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="saveLinkage">保存联动配置</a-button>
            </a-form-item>
          </a-form>
          <a-divider>示例联动</a-divider>
          <a-list size="small" :data="linkageExamples">
            <template #render="{ item }">
              <a-list-item>
                <icon-check-circle-fill style="color: green; margin-right: 8px" />
                {{ item }}
              </a-list-item>
            </template>
          </a-list>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const showAddDevice = ref(false)
const showAddScene = ref(false)
const deviceType = ref('')
const linkageConfig = ref('')
const linkageExamples = [
  '宠物开心时 → 灯光调为暖色调',
  '宠物睡觉时 → 空调调至静音模式',
  '室内温度过高 → 自动开启加湿器',
  '宠物情绪低落 → 播放舒缓音乐'
]

const stats = reactive({ devices: 8, online: 6, scenes: 5 })

const devices = ref([
  { id: 1, name: '客厅大灯', type: 'light', room: '客厅', status: true },
  { id: 2, name: '卧室空调', type: 'ac', room: '卧室', status: false },
  { id: 3, name: '客厅加湿器', type: 'humidifier', room: '客厅', status: true },
  { id: 4, name: '背景音乐', type: 'speaker', room: '全屋', status: false }
])

const scenes = ref([
  { id: 1, name: '宠物睡眠模式', trigger: '宠物入睡', actions: '关灯+静音空调', enabled: true },
  { id: 2, name: '宠物活动模式', trigger: '宠物活跃', actions: '开灯+播放音乐', enabled: true },
  { id: 3, name: '离家模式', trigger: '主人外出', actions: '全部关闭', enabled: false }
])

const filteredDevices = computed(() => {
  if (!deviceType.value) return devices.value
  return devices.value.filter(d => d.type === deviceType.value)
})

const sceneColumns = [
  { title: '场景名称', dataIndex: 'name' },
  { title: '触发条件', dataIndex: 'trigger' },
  { title: '执行动作', dataIndex: 'actions' },
  { title: '状态', slotName: 'status' },
  { title: '操作', slotName: 'actions' }
]

function toggleDevice(device) {
  Message.success((device.status ? '开启' : '关闭') + device.name)
}

function editScene(record) {
  Message.info('编辑场景: ' + record.name)
}

function deleteScene(record) {
  scenes.value = scenes.value.filter(s => s.id !== record.id)
  Message.success('已删除场景')
}

async function saveLinkage() {
  Message.success('联动配置已保存')
}

onMounted(() => {
  linkageConfig.value = '# 宠物情绪联动配置\nhappy: 灯光=暖色, 音乐=欢快\nsleep: 灯光=关闭, 空调=静音\nsad: 音乐=舒缓'
})
</script>

<style scoped>
.device-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100px;
  font-size: 48px;
  color: var(--color-primary);
  background: var(--color-fill-1);
}
</style>
