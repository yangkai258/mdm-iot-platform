<template>
  <div class="device-provisioning-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="配置模板" :value="15" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="已配置设备" :value="856" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="配置成功率" :value="98" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>设备批量配置</span>
          <a-button type="primary" @click="handleCreateTemplate">
            <template #icon><icon-plus /></template>
            创建模板
          </a-button>
        </div>
      </template>
      
      <a-tabs>
        <a-tab-pane key="templates" title="配置模板">
          <a-table :columns="templateColumns" :data="templates" :pagination="pagination">
            <template #actions="{ record }">
              <a-space>
                <a-link @click="handleApply(record)">应用</a-link>
                <a-link @click="handleEdit(record)">编辑</a-link>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="batch" title="批量配置">
          <a-form layout="vertical" style="max-width: 800px;">
            <a-form-item label="选择设备">
              <a-select v-model="selectedDevices" multiple placeholder="选择要配置的设备">
                <a-option value="DEV001">小黄-客厅</a-option>
                <a-option value="DEV002">小红-卧室</a-option>
                <a-option value="DEV003">咪咪-阳台</a-option>
              </a-select>
            </a-form-item>
            <a-form-item label="选择模板">
              <a-select v-model="selectedTemplate" placeholder="选择配置模板">
                <a-option v-for="tpl in templates" :key="tpl.id" :value="tpl.id">{{ tpl.name }}</a-option>
              </a-select>
            </a-form-item>
            <a-form-item label="配置参数">
              <a-form layout="inline">
                <a-form-item label="WiFi SSID">
                  <a-input v-model="config.wifiSsid" />
                </a-form-item>
                <a-form-item label="WiFi密码">
                  <a-input-password v-model="config.wifiPassword" />
                </a-form-item>
                <a-form-item label="服务器地址">
                  <a-input v-model="config.serverUrl" />
                </a-form-item>
              </a-form>
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="handleBatchApply">批量应用配置</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>
        
        <a-tab-pane key="history" title="配置历史">
          <a-table :columns="historyColumns" :data="history" :pagination="pagination" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 6 });
const selectedDevices = ref([]);
const selectedTemplate = ref('');
const config = reactive({ wifiSsid: '', wifiPassword: '', serverUrl: '' });

const templates = ref([
  { id: 1, name: '家庭网络配置', description: '适用于家庭WiFi环境', deviceCount: 520, createdAt: '2026-03-20' },
  { id: 2, name: '企业网络配置', description: '适用于企业WiFi环境', deviceCount: 280, createdAt: '2026-03-15' },
  { id: 3, name: '离线模式配置', description: '适用于无网络环境', deviceCount: 56, createdAt: '2026-03-10' },
]);

const history = ref([
  { id: 1, deviceId: 'DEV001', template: '家庭网络配置', status: 'success', statusText: '成功', time: '2026-03-28 18:00:00' },
  { id: 2, deviceId: 'DEV002', template: '企业网络配置', status: 'failed', statusText: '失败', reason: 'WiFi密码错误', time: '2026-03-28 17:00:00' },
]);

const templateColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '模板名称', dataIndex: 'name', width: 150 },
  { title: '描述', dataIndex: 'description' },
  { title: '设备数', dataIndex: 'deviceCount', width: 100 },
  { title: '创建时间', dataIndex: 'createdAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const historyColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '设备ID', dataIndex: 'deviceId', width: 120 },
  { title: '模板', dataIndex: 'template', width: 150 },
  { title: '状态', dataIndex: 'statusText', width: 100 },
  { title: '原因', dataIndex: 'reason', width: 150 },
  { title: '时间', dataIndex: 'time', width: 160 },
];

const handleCreateTemplate = () => {};
const handleApply = (record: any) => { selectedTemplate.value = record.id; };
const handleEdit = (record: any) => {};
const handleBatchApply = () => {};
</script>

<style scoped>
.device-provisioning-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
