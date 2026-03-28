<template>
  <div class="system-health-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="系统状态" :value="status" :value-style="{ color: statusColor }" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="在线设备" :value="onlineDevices" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="API响应" :value="98" suffix="%" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="内存使用" :value="64" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <span>系统健康监控</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="overview" title="健康概览">
          <a-row :gutter="16">
            <a-col :span="12">
              <a-card title="服务状态">
                <a-list>
                  <a-list-item v-for="svc in services" :key="svc.name">
                    <a-list-item-meta :title="svc.name" :description="svc.desc">
                      <template #avatar>
                        <a-badge :status="svc.status === 'running' ? 'success' : 'error'" />
                      </template>
                    </a-list-item-meta>
                  </a-list-item>
                </a-list>
              </a-card>
            </a-col>
            <a-col :span="12">
              <a-card title="资源使用">
                <a-progress :percent="64" label="内存" />
                <a-progress :percent="45" label="CPU" />
                <a-progress :percent="51" label="磁盘" />
                <a-progress :percent="28" label="网络" />
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
        
        <a-tab-pane key="metrics" title="性能指标">
          <a-table :columns="metricColumns" :data="metrics" :pagination="false" />
        </a-tab-pane>
        
        <a-tab-pane key="alerts" title="告警记录">
          <a-table :columns="alertColumns" :data="alerts" :pagination="pagination" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue';

const status = computed(() => '运行正常');
const statusColor = computed(() => '#00b42a');
const onlineDevices = ref(1256);
const pagination = reactive({ current: 1, pageSize: 10, total: 5 });

const services = ref([
  { name: 'API服务', desc: 'http://localhost:8080', status: 'running' },
  { name: '前端服务', desc: 'http://localhost:3002', status: 'running' },
  { name: '数据库', desc: 'PostgreSQL:5432', status: 'running' },
  { name: '缓存服务', desc: 'Redis:6379', status: 'running' },
  { name: '消息队列', desc: 'EMQX:1883', status: 'running' },
]);

const metrics = ref([
  { id: 1, name: 'QPS', value: '2,580', avg: '2,500', max: '5,200' },
  { id: 2, name: '响应时间', value: '45ms', avg: '50ms', max: '120ms' },
  { id: 3, name: '并发连接', value: '856', avg: '800', max: '1,200' },
  { id: 4, name: '错误率', value: '0.1%', avg: '0.2%', max: '0.5%' },
]);

const alerts = ref([
  { id: 1, level: 'warning', levelText: '警告', message: '内存使用率超过80%', time: '2026-03-28 18:00:00' },
  { id: 2, level: 'info', levelText: '通知', message: '数据库备份完成', time: '2026-03-28 02:00:00' },
]);

const metricColumns = [
  { title: '指标', dataIndex: 'name', width: 150 },
  { title: '当前值', dataIndex: 'value', width: 120 },
  { title: '平均值', dataIndex: 'avg', width: 120 },
  { title: '最大值', dataIndex: 'max', width: 120 },
];

const alertColumns = [
  { title: '级别', dataIndex: 'levelText', width: 100 },
  { title: '消息', dataIndex: 'message' },
  { title: '时间', dataIndex: 'time', width: 160 },
];
</script>

<style scoped>
.system-health-container { padding: 20px; }
</style>
