<!-- components/menu/index.vue - 侧边菜单 -->
<template>
  <div class="menu-wrapper">
    <!-- Logo 区域 -->
    <div class="logo-wrapper">
      <img v-if="!collapsed" src="@/assets/logo.svg" alt="logo" class="logo" />
      <span v-if="!collapsed" class="logo-text">MDM 控制台</span>
      <span v-else class="logo-collapsed">M</span>
    </div>

    <!-- 菜单 -->
    <a-menu
      v-model:selected-keys="selectedKey"
      v-model:open-keys="openKeys"
      :mode="'vertical'"
      :collapsed="collapsed"
      :accordion="true"
      :tooltip-show-to-parent="false"
    >
      <a-menu-item key="dashboard" @click="navigate('dashboard')">
        <template #icon><icon-dashboard /></template>
        <span>设备大盘</span>
      </a-menu-item>

      <a-menu-item key="devices" @click="navigate('devices')">
        <template #icon><icon-mobile /></template>
        <span>设备列表</span>
      </a-menu-item>

      <!-- 设备增强 -->
      <a-sub-menu key="devices-enhanced">
        <template #icon><icon-safe /></template>
        <template #title>设备增强</template>
        <a-menu-item key="devices-certificates" @click="navigate('devices-certificates')">证书管理</a-menu-item>
        <a-menu-item key="devices-remote" @click="navigate('devices-remote')">远程控制</a-menu-item>
        <a-menu-item key="devices-geofence" @click="navigate('devices-geofence')">地理围栏</a-menu-item>
        <a-menu-item key="devices-monitor" @click="navigate('devices-monitor')">监控面板</a-menu-item>
        <a-menu-item key="devices-pairing" @click="navigate('devices-pairing')">配对管理</a-menu-item>
      </a-sub-menu>

      <!-- OTA升级 -->
      <a-sub-menu key="ota">
        <template #icon><icon-upload /></template>
        <template #title>OTA升级</template>
        <a-menu-item key="ota-packages" @click="navigate('ota-packages')">固件包管理</a-menu-item>
        <a-menu-item key="ota-deployments" @click="navigate('ota-deployments')">部署任务</a-menu-item>
      </a-sub-menu>

      <a-menu-item key="miniclaw-firmwares" @click="navigate('miniclaw-firmwares')">
        <template #icon><icon-storage /></template>
        <span>固件管理</span>
      </a-menu-item>

      <!-- 告警中心 -->
      <a-sub-menu key="alerts">
        <template #icon><icon-mind-mapping /></template>
        <template #title>告警中心</template>
        <a-menu-item key="alert-rules" @click="navigate('alert-rules')">告警规则</a-menu-item>
        <a-menu-item key="alert-list" @click="navigate('alert-list')">告警列表</a-menu-item>
        <a-menu-item key="alert-settings" @click="navigate('alert-settings')">告警设置</a-menu-item>
      </a-sub-menu>

      <!-- 宠物管理 -->
      <a-sub-menu key="pet">
        <template #icon><icon-star /></template>
        <template #title>宠物管理</template>
        <a-menu-item key="pet-config" @click="navigate('pet-config')">宠物配置</a-menu-item>
        <a-menu-item key="pet-console" @click="navigate('pet-console')">宠物控制台</a-menu-item>
        <a-menu-item key="pet-conversations" @click="navigate('pet-conversations')">对话记录</a-menu-item>
      </a-sub-menu>

      <!-- 会员管理 -->
      <a-sub-menu key="members">
        <template #icon><icon-user-group /></template>
        <template #title>会员管理</template>
        <a-menu-item key="members" @click="navigate('members')">会员列表</a-menu-item>
        <a-menu-item key="members-cards" @click="navigate('members-cards')">会员卡管理</a-menu-item>
        <a-menu-item key="members-coupons" @click="navigate('members-coupons')">优惠券</a-menu-item>
        <a-menu-item key="members-stores" @click="navigate('members-stores')">门店管理</a-menu-item>
        <a-menu-item key="members-levels" @click="navigate('members-levels')">会员等级</a-menu-item>
        <a-menu-item key="members-tags" @click="navigate('members-tags')">标签管理</a-menu-item>
        <a-menu-item key="members-points" @click="navigate('members-points')">积分管理</a-menu-item>
        <a-menu-item key="members-promotions" @click="navigate('members-promotions')">促销活动</a-menu-item>
      </a-sub-menu>

      <a-menu-item key="knowledge" @click="navigate('knowledge')">
        <template #icon><icon-book /></template>
        <span>知识库</span>
      </a-menu-item>

      <a-menu-item key="owner-profile" @click="navigate('owner-profile')">
        <template #icon><icon-user /></template>
        <span>主人档案</span>
      </a-menu-item>

      <!-- 健康医疗 -->
      <a-sub-menu key="health">
        <template #icon><icon-heart /></template>
        <template #title>健康医疗</template>
        <a-menu-item key="health-sports" @click="navigate('health-sports')">运动统计</a-menu-item>
        <a-menu-item key="health-reports" @click="navigate('health-reports')">健康报告</a-menu-item>
        <a-menu-item key="health-warning" @click="navigate('health-warning')">健康预警</a-menu-item>
        <a-menu-item key="health-sleep" @click="navigate('health-sleep')">睡眠分析</a-menu-item>
      </a-sub-menu>

      <!-- AI 功能 -->
      <a-sub-menu key="ai">
        <template #icon><icon-robot /></template>
        <template #title>AI 功能</template>
        <a-menu-item key="ai-behavior" @click="navigate('ai-behavior')">行为分析</a-menu-item>
        <a-menu-item key="ai-emotion" @click="navigate('ai-emotion')">情感识别</a-menu-item>
      </a-sub-menu>

      <!-- 情感计算 -->
      <a-sub-menu key="emotion">
        <template #icon><icon-heart /></template>
        <template #title>情感计算</template>
        <a-menu-item key="emotion-recognition" @click="navigate('emotion-recognition')">情绪识别</a-menu-item>
        <a-menu-item key="emotion-response" @click="navigate('emotion-response')">情绪响应</a-menu-item>
        <a-menu-item key="emotion-logs" @click="navigate('emotion-logs')">情绪日志</a-menu-item>
        <a-menu-item key="emotion-family-map" @click="navigate('emotion-family-map')">家庭情绪地图</a-menu-item>
      </a-sub-menu>

      <!-- 数字孪生 -->
      <a-sub-menu key="digital-twin">
        <template #icon><icon-heart /></template>
        <template #title>数字孪生</template>
        <a-menu-item key="digital-twin-vitals" @click="navigate('digital-twin-vitals')">实时生命体征</a-menu-item>
        <a-menu-item key="digital-twin-prediction" @click="navigate('digital-twin-prediction')">行为预测</a-menu-item>
        <a-menu-item key="digital-twin-playback" @click="navigate('digital-twin-playback')">历史回放</a-menu-item>
        <a-menu-item key="digital-twin-moments" @click="navigate('digital-twin-moments')">精彩瞬间</a-menu-item>
      </a-sub-menu>

      <!-- 具身智能 -->
      <a-sub-menu key="embodied">
        <template #icon><icon-mind-mapping /></template>
        <template #title>具身智能</template>
        <a-menu-item key="embodied-perception" @click="navigate('embodied-perception')">环境感知</a-menu-item>
        <a-menu-item key="embodied-spatial" @click="navigate('embodied-spatial')">空间认知</a-menu-item>
        <a-menu-item key="embodied-motion" @click="navigate('embodied-motion')">动作模仿</a-menu-item>
        <a-menu-item key="embodied-decision" @click="navigate('embodied-decision')">AI决策引擎</a-menu-item>
      </a-sub-menu>

      <!-- 仿真测试 -->
      <a-sub-menu key="simulation">
        <template #icon><icon-tool /></template>
        <template #title>仿真测试</template>
        <a-menu-item key="simulation-pet" @click="navigate('simulation-pet')">虚拟宠物仿真</a-menu-item>
        <a-menu-item key="simulation-test" @click="navigate('simulation-test')">自动化测试框架</a-menu-item>
        <a-menu-item key="simulation-replay" @click="navigate('simulation-replay')">系统回放</a-menu-item>
      </a-sub-menu>

      <!-- 应用管理 -->
      <a-sub-menu key="apps">
        <template #icon><icon-apps /></template>
        <template #title>应用管理</template>
        <a-menu-item key="apps-list" @click="navigate('apps-list')">应用列表</a-menu-item>
        <a-menu-item key="apps-versions" @click="navigate('apps-versions')">应用版本</a-menu-item>
        <a-menu-item key="apps-distributions" @click="navigate('apps-distributions')">应用分发</a-menu-item>
      </a-sub-menu>

      <!-- 系统管理 -->
      <a-sub-menu key="system">
        <template #icon><icon-settings /></template>
        <template #title>系统管理</template>
        <a-menu-item key="system-monitor" @click="navigate('system-monitor')">服务监控</a-menu-item>
        <a-menu-item key="system-logs" @click="navigate('system-logs')">操作日志</a-menu-item>
        <a-menu-item key="policies-list" @click="navigate('policies-list')">策略列表</a-menu-item>
        <a-menu-item key="policies-configs" @click="navigate('policies-configs')">策略配置</a-menu-item>
        <a-menu-item key="compliance-rules" @click="navigate('compliance-rules')">合规规则</a-menu-item>
      </a-sub-menu>

      <!-- 租户管理 -->
      <a-sub-menu key="tenants">
        <template #icon><icon-tool /></template>
        <template #title>租户管理</template>
        <a-menu-item key="tenants-approval" @click="navigate('tenants-approval')">租户入驻审核</a-menu-item>
        <a-menu-item key="tenants-management" @click="navigate('tenants-management')">租户系统管理</a-menu-item>
        <a-menu-item key="tenants-public-archives" @click="navigate('tenants-public-archives')">公共档案</a-menu-item>
        <a-menu-item key="tenants-system-info" @click="navigate('tenants-system-info')">系统信息</a-menu-item>
      </a-sub-menu>

      <!-- 组织管理 -->
      <a-sub-menu key="org">
        <template #icon><icon-tool /></template>
        <template #title>组织管理</template>
        <a-menu-item key="org-companies" @click="navigate('org-companies')">公司管理</a-menu-item>
        <a-menu-item key="org-departments" @click="navigate('org-departments')">部门管理</a-menu-item>
        <a-menu-item key="org-posts" @click="navigate('org-posts')">岗位管理</a-menu-item>
        <a-menu-item key="org-employees" @click="navigate('org-employees')">员工管理</a-menu-item>
        <a-menu-item key="org-standard-positions" @click="navigate('org-standard-positions')">标准岗位</a-menu-item>
      </a-sub-menu>

      <!-- 多维权限 -->
      <a-sub-menu key="permissions">
        <template #icon><icon-safe /></template>
        <template #title>多维权限</template>
        <a-menu-item key="permissions-roles" @click="navigate('permissions-roles')">角色管理</a-menu-item>
        <a-menu-item key="permissions-menus" @click="navigate('permissions-menus')">菜单管理</a-menu-item>
        <a-menu-item key="permissions-groups" @click="navigate('permissions-groups')">权限组</a-menu-item>
        <a-menu-item key="permissions-data-config" @click="navigate('permissions-data-config')">数据权限</a-menu-item>
        <a-menu-item key="permissions-api" @click="navigate('permissions-api')">API权限</a-menu-item>
      </a-sub-menu>

      <!-- 通知管理 -->
      <a-sub-menu key="notifications">
        <template #icon><icon-message /></template>
        <template #title>通知管理</template>
        <a-menu-item key="notifications-list" @click="navigate('notifications-list')">推送通知</a-menu-item>
        <a-menu-item key="notifications-announcements" @click="navigate('notifications-announcements')">公告管理</a-menu-item>
        <a-menu-item key="notifications-templates" @click="navigate('notifications-templates')">通知模板</a-menu-item>
      </a-sub-menu>

      <!-- 第三方接入 -->
      <a-sub-menu key="integration">
        <template #icon><icon-link /></template>
        <template #title>第三方接入</template>
        <a-menu-item key="integration-pet-hospitals" @click="navigate('integration-pet-hospitals')">宠物医院</a-menu-item>
      </a-sub-menu>

      <a-menu-item key="portal" @click="navigate('portal')">
        <template #icon><icon-dashboard /></template>
        <span>工作台门户</span>
      </a-menu-item>

      <a-menu-item key="subscription" @click="navigate('subscription')">
        <template #icon><icon-subscribe /></template>
        <span>订阅管理</span>
      </a-menu-item>

      <!-- Webhook -->
      <a-sub-menu key="webhooks">
        <template #icon><icon-link /></template>
        <template #title>Webhook</template>
        <a-menu-item key="webhooks" @click="navigate('webhooks')">Webhook列表</a-menu-item>
        <a-menu-item key="webhooks-logs" @click="navigate('webhooks-logs')">调用日志</a-menu-item>
      </a-sub-menu>

      <a-menu-item key="developer" @click="navigate('developer')">
        <template #icon><icon-code /></template>
        <span>开发者API</span>
      </a-menu-item>

      <a-menu-item key="billing" @click="navigate('billing')">
        <template #icon><icon-tool /></template>
        <span>用量计费</span>
      </a-menu-item>
    </a-menu>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAppStore } from '@/store';
import {
  IconDashboard,
  IconMobile,
  IconUpload,
  IconMindMapping,
  IconUserGroup,
  IconSettings,
  IconApps,
  IconSafe,
  IconMessage,
  IconTool,
  IconRobot,
  IconHeart,
  IconStar,
  IconBook,
  IconUser,
  IconStorage,
  IconLink,
  IconSubscribe,
  IconCode,
} from '@arco-design/web-vue/es/icon';

const router = useRouter();
const route = useRoute();
const appStore = useAppStore();

const selectedKey = ref<string[]>([]);
const openKeys = ref<string[]>([]);
const collapsed = computed(() => appStore.menuCollapse);

// 路由到菜单 key 的映射
const routeNameToMenuKey: Record<string, string> = {
  'Dashboard': 'dashboard',
  'Devices': 'devices',
  'OtaPackages': 'ota-packages',
  'OtaDeployments': 'ota-deployments',
  'FirmwareList': 'miniclaw-firmwares',
  'Alert': 'alerts',
  'AlertRules': 'alert-rules',
  'AlertList': 'alert-list',
  'AlertSettings': 'alert-settings',
  'Pet': 'pet-config',
  'PetConsole': 'pet-console',
  'PetConversations': 'pet-conversations',
  'Members': 'members',
  'MemberCards': 'members-cards',
  'MemberCoupons': 'members-coupons',
  'MemberStores': 'members-stores',
  'MemberLevels': 'members-levels',
  'MemberTags': 'members-tags',
  'MemberPoints': 'members-points',
  'MemberPromotions': 'members-promotions',
  'KnowledgeList': 'knowledge',
  'OwnerProfile': 'owner-profile',
  'Monitor': 'system-monitor',
  'Logs': 'system-logs',
  'PolicyList': 'policies-list',
  'PolicyConfigs': 'policies-configs',
  'ComplianceRules': 'compliance-rules',
  'TenantApproval': 'tenants-approval',
  'TenantManagement': 'tenants-management',
  'PublicArchives': 'tenants-public-archives',
  'SystemInfo': 'tenants-system-info',
  'OrgCompanies': 'org-companies',
  'OrgDepartments': 'org-departments',
  'OrgPosts': 'org-posts',
  'OrgEmployees': 'org-employees',
  'StandardPositions': 'org-standard-positions',
  'PermissionRoles': 'permissions-roles',
  'PermissionMenus': 'permissions-menus',
  'PermissionGroups': 'permissions-groups',
  'DataPermissionConfig': 'permissions-data-config',
  'ApiPermissions': 'permissions-api',
  'NotificationList': 'notifications-list',
  'Announcements': 'notifications-announcements',
  'NotificationTemplates': 'notifications-templates',
  'AppList': 'apps-list',
  'AppVersions': 'apps-versions',
  'AppDistributions': 'apps-distributions',
  'PetHospitals': 'integration-pet-hospitals',
  'AiBehavior': 'ai-behavior',
  'AiEmotion': 'ai-emotion',
  'EmotionRecognition': 'emotion-recognition',
  'EmotionResponse': 'emotion-response',
  'EmotionLogs': 'emotion-logs',
  'EmotionFamilyMap': 'emotion-family-map',
  'HealthWarnings': 'health-warning',
  'ExerciseStats': 'health-sports',
  'HealthReports': 'health-reports',
  'SleepAnalysis': 'health-sleep',
  'Portal': 'portal',
  'SubscriptionList': 'subscription',
  'WebhookList': 'webhooks',
  'WebhookLogs': 'webhooks-logs',
  'DeveloperApi': 'developer',
  'BillingList': 'billing',
  'DeviceCertificates': 'devices-certificates',
  'DeviceRemoteControl': 'devices-remote',
  'DeviceGeofence': 'devices-geofence',
  'DeviceMonitorPanel': 'devices-monitor',
  'DevicePairing': 'devices-pairing',
  'DigitalTwinVitals': 'digital-twin-vitals',
  'DigitalTwinPrediction': 'digital-twin-prediction',
  'DigitalTwinPlayback': 'digital-twin-playback',
  'DigitalTwinMoments': 'digital-twin-moments',
  'EmbodiedPerception': 'embodied-perception',
  'EmbodiedSpatial': 'embodied-spatial',
  'EmbodiedMotion': 'embodied-motion',
  'EmbodiedDecision': 'embodied-decision',
  'SimulationPet': 'simulation-pet',
  'SimulationTest': 'simulation-test',
  'SimulationReplay': 'simulation-replay',
};

// 菜单 key → 路由 path 映射
const menuKeyToPath: Record<string, string> = {
  'dashboard': '/dashboard',
  'devices': '/devices',
  'ota-packages': '/ota/packages',
  'ota-deployments': '/ota/deployments',
  'miniclaw-firmwares': '/miniclaw/firmwares',
  'alert-rules': '/alerts/rules',
  'alert-list': '/alerts/list',
  'alert-settings': '/alerts/settings',
  'pet-config': '/pet',
  'pet-console': '/pet/console',
  'pet-conversations': '/pet/conversations',
  'members': '/members',
  'members-cards': '/members/cards',
  'members-coupons': '/members/coupons',
  'members-stores': '/members/stores',
  'members-levels': '/members/levels',
  'members-tags': '/members/tags',
  'members-points': '/members/points',
  'members-promotions': '/members/promotions',
  'knowledge': '/knowledge',
  'owner-profile': '/owner/profile',
  'health-sports': '/health/exercise-stats',
  'health-reports': '/health/reports',
  'health-warning': '/health/warnings',
  'health-sleep': '/health/sleep',
  'ai-behavior': '/ai/behavior',
  'ai-emotion': '/ai/emotion',
  'emotion-recognition': '/emotion/recognition',
  'emotion-response': '/emotion/response',
  'emotion-logs': '/emotion/logs',
  'emotion-family-map': '/emotion/family-map',
  'digital-twin-vitals': '/digital-twin/vitals',
  'digital-twin-prediction': '/digital-twin/prediction',
  'digital-twin-playback': '/digital-twin/playback',
  'digital-twin-moments': '/digital-twin/moments',
  'embodied-perception': '/embodied/perception',
  'embodied-spatial': '/embodied/spatial',
  'embodied-motion': '/embodied/motion',
  'embodied-decision': '/embodied/decision',
  'simulation-pet': '/simulation/pet',
  'simulation-test': '/simulation/test',
  'simulation-replay': '/simulation/replay',
  'apps-list': '/apps/list',
  'apps-versions': '/apps/versions',
  'apps-distributions': '/apps/distributions',
  'system-monitor': '/system/monitor',
  'system-logs': '/system/logs',
  'policies-list': '/policies/list',
  'policies-configs': '/policies/configs',
  'compliance-rules': '/policies/compliance',
  'tenants-approval': '/tenants/approval',
  'tenants-management': '/tenants/management',
  'tenants-public-archives': '/tenants/public-archives',
  'tenants-system-info': '/tenants/system-info',
  'org-companies': '/org/companies',
  'org-departments': '/org/departments',
  'org-posts': '/org/posts',
  'org-employees': '/org/employees',
  'org-standard-positions': '/org/standard-positions',
  'permissions-roles': '/permissions/roles',
  'permissions-menus': '/permissions/menus',
  'permissions-groups': '/permissions/groups',
  'permissions-data-config': '/permissions/data-config',
  'permissions-api': '/permissions/api',
  'notifications-list': '/notifications/list',
  'notifications-announcements': '/notifications/announcements',
  'notifications-templates': '/notifications/templates',
  'integration-pet-hospitals': '/integration/pet-hospitals',
  'portal': '/portal',
  'subscription': '/subscription',
  'webhooks': '/webhooks',
  'webhooks-logs': '/webhooks/logs',
  'developer': '/developer',
  'billing': '/billing',
  'devices-certificates': '/devices/certificates',
  'devices-remote': '/devices/remote',
  'devices-geofence': '/devices/geofence',
  'devices-monitor': '/devices/monitor',
  'devices-pairing': '/devices/pairing',
};

watch(
  () => route.name,
  (name) => {
    if (name) {
      const menuKey = routeNameToMenuKey[name as string] || name as string;
      selectedKey.value = [menuKey];
      const keyParts = menuKey.split('-');
      if (keyParts.length > 1) {
        openKeys.value = [keyParts[0]];
      }
    }
  },
  { immediate: true }
);

// 导航方法
const navigate = (key: string) => {
  const path = menuKeyToPath[key] || '/' + key.split('-').join('/');
  router.push(path);
};

// 全局暴露
if (typeof window !== 'undefined') {
  (window as any).menuNavigate = navigate;
  (window as any).vueRouter = router;
}
</script>

<style scoped lang="less">
.menu-wrapper {
  height: 100%;
  display: flex;
  flex-direction: column;

  .logo-wrapper {
    height: 60px;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    padding: 0 16px;
    background: var(--color-bg-2);
    border-bottom: 1px solid var(--color-border);

    .logo {
      width: 32px;
      height: 32px;
    }

    .logo-text {
      font-size: 16px;
      font-weight: 600;
      color: var(--color-text-1);
      white-space: nowrap;
    }

    .logo-collapsed {
      font-size: 20px;
      font-weight: 700;
      color: var(--color-primary);
    }
  }

  :deep(.arco-menu) {
    flex: 1;
    overflow-y: auto;

    &::-webkit-scrollbar {
      width: 4px;
    }

    &::-webkit-scrollbar-thumb {
      background: rgba(255, 255, 255, 0.2);
      border-radius: 2px;
    }
  }
}
</style>
