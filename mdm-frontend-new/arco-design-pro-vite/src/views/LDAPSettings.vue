<template>
  <div class="ldap-settings-container">
    <a-card>
      <template #title>
        <span>LDAP/AD 集成设置</span>
      </template>
      
      <a-tabs>
        <!-- 连接配置 -->
        <a-tab-pane key="connection" title="连接配置">
          <a-form :model="connectionForm" layout="vertical" style="max-width: 600px;">
            <a-form-item label="LDAP服务器地址" required>
              <a-input v-model="connectionForm.server" placeholder="ldap://ldap.example.com:389" />
            </a-form-item>
            <a-form-item label="端口" required>
              <a-input-number v-model="connectionForm.port" :min="1" :max="65535" />
            </a-form-item>
            <a-form-item label="Base DN" required>
              <a-input v-model="connectionForm.baseDN" placeholder="dc=example,dc=com" />
            </a-form-item>
            <a-form-item label="管理员DN">
              <a-input v-model="connectionForm.adminDN" placeholder="cn=admin,dc=example,dc=com" />
            </a-form-item>
            <a-form-item label="管理员密码">
              <a-input-password v-model="connectionForm.password" placeholder="请输入密码" />
            </a-form-item>
            <a-form-item label="使用SSL">
              <a-switch v-model="connectionForm.useSSL" />
            </a-form-item>
            <a-form-item>
              <a-space>
                <a-button type="primary" @click="handleTestConnection">测试连接</a-button>
                <a-button @click="handleSaveConnection">保存</a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <!-- 用户同步 -->
        <a-tab-pane key="sync" title="用户同步">
          <a-space style="margin-bottom: 16px;">
            <a-button type="primary" @click="handleSyncNow">立即同步</a-button>
            <a-button @click="handleSyncSettings">同步设置</a-button>
          </a-space>
          
          <a-table :columns="syncColumns" :data="syncLogs" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="getSyncStatusColor(record.status)">{{ record.statusText }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>

        <!-- 同步规则 -->
        <a-tab-pane key="rules" title="同步规则">
          <a-form :model="ruleForm" layout="vertical" style="max-width: 600px;">
            <a-form-item label="用户过滤器">
              <a-input v-model="ruleForm.userFilter" placeholder="(objectClass=person)" />
            </a-form-item>
            <a-form-item label="用户属性映射">
              <a-space direction="vertical" fill>
                <a-form-item label="用户名属性">
                  <a-select v-model="ruleForm.usernameAttr" placeholder="选择属性">
                    <a-option value="uid">uid</a-option>
                    <a-option value="sAMAccountName">sAMAccountName</a-option>
                    <a-option value="cn">cn</a-option>
                  </a-select>
                </a-form-item>
                <a-form-item label="邮箱属性">
                  <a-select v-model="ruleForm.emailAttr" placeholder="选择属性">
                    <a-option value="mail">mail</a-option>
                    <a-option value="userPrincipalName">userPrincipalName</a-option>
                  </a-select>
                </a-form-item>
                <a-form-item label="部门属性">
                  <a-select v-model="ruleForm.deptAttr" placeholder="选择属性">
                    <a-option value="department">department</a-option>
                    <a-option value="ou">ou</a-option>
                  </a-select>
                </a-form-item>
              </a-space>
            </a-form-item>
            <a-form-item label="自动同步间隔">
              <a-select v-model="ruleForm.syncInterval" placeholder="选择间隔">
                <a-option value="1h">每小时</a-option>
                <a-option value="6h">每6小时</a-option>
                <a-option value="12h">每12小时</a-option>
                <a-option value="daily">每天</a-option>
              </a-select>
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="handleSaveRules">保存规则</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <!-- 同步设置弹窗 -->
    <a-modal v-model:visible="settingsVisible" title="同步设置">
      <a-form :model="settingsForm" layout="vertical">
        <a-form-item label="同步范围">
          <a-checkbox-group v-model="settingsForm.scope">
            <a-checkbox value="users">同步用户</a-checkbox>
            <a-checkbox value="groups">同步用户组</a-checkbox>
            <a-checkbox value="ous">同步组织单位</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
        <a-form-item label="冲突处理策略">
          <a-radio-group v-model="settingsForm.conflictStrategy">
            <a-radio value="ldap_wins">LDAP优先（覆盖本地）</a-radio>
            <a-radio value="local_wins">本地优先（保留本地）</a-radio>
            <a-radio value="manual">手动处理</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="同步后自动启用">
          <a-switch v-model="settingsForm.autoEnable" />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="settingsVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSaveSettings">保存</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 3 });

const connectionForm = reactive({
  server: 'ldap://ldap.example.com',
  port: 389,
  baseDN: 'dc=example,dc=com',
  adminDN: 'cn=admin,dc=example,dc=com',
  password: '',
  useSSL: false,
});

const ruleForm = reactive({
  userFilter: '(objectClass=person)',
  usernameAttr: 'uid',
  emailAttr: 'mail',
  deptAttr: 'department',
  syncInterval: 'daily',
});

const settingsForm = reactive({
  scope: ['users'],
  conflictStrategy: 'ldap_wins',
  autoEnable: true,
});

const syncLogs = ref([
  { time: '2026-03-28 10:00:00', type: 'full', status: 'success', statusText: '成功', users: 150, groups: 12, duration: '5m30s' },
  { time: '2026-03-27 10:00:00', type: 'full', status: 'success', statusText: '成功', users: 148, groups: 12, duration: '5m15s' },
  { time: '2026-03-26 10:00:00', type: 'incremental', status: 'success', statusText: '成功', users: 3, groups: 0, duration: '45s' },
]);

const syncColumns = [
  { title: '同步时间', dataIndex: 'time', width: 160 },
  { title: '类型', dataIndex: 'type', width: 100 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '用户数', dataIndex: 'users', width: 80 },
  { title: '用户组', dataIndex: 'groups', width: 80 },
  { title: '耗时', dataIndex: 'duration', width: 80 },
];

const settingsVisible = ref(false);

const getSyncStatusColor = (s: string) => ({ success: 'green', failed: 'red', running: 'blue' }[s] || 'default');

const handleTestConnection = () => {};
const handleSaveConnection = () => {};
const handleSyncNow = () => {};
const handleSyncSettings = () => { settingsVisible.value = true; };
const handleSaveRules = () => {};
const handleSaveSettings = () => { settingsVisible.value = false; };
</script>

<style scoped>
.ldap-settings-container { padding: 20px; }
</style>
