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
      @click="handleMenuClick"
    >
      <!-- Dashboard -->
      <a-menu-item key="dashboard">
        <template #icon><icon-dashboard /></template>
        <span>设备大盘</span>
      </a-menu-item>

      <!-- 设备管理 -->
      <a-menu-item key="devices">
        <template #icon><icon-mobile /></template>
        <span>设备列表</span>
      </a-menu-item>

      <!-- OTA升级 -->
      <a-sub-menu key="ota">
        <template #icon><icon-upload /></template>
        <template #title>OTA升级</template>
        <a-menu-item key="ota-packages">固件包管理</a-menu-item>
        <a-menu-item key="ota-deployments">部署任务</a-menu-item>
      </a-sub-menu>

      <!-- 固件管理 -->
      <a-menu-item key="miniclaw-firmwares">
        <template #icon><icon-storage /></template>
        <span>固件管理</span>
      </a-menu-item>

      <!-- 告警中心 -->
      <a-sub-menu key="alerts">
        <template #icon><icon-mind-mapping /></template>
        <template #title>告警中心</template>
        <a-menu-item key="alert-rules">告警规则</a-menu-item>
        <a-menu-item key="alert-list">告警列表</a-menu-item>
        <a-menu-item key="alert-settings">告警设置</a-menu-item>
      </a-sub-menu>

      <!-- 宠物管理 -->
      <a-sub-menu key="pet">
        <template #icon><icon-star /></template>
        <template #title>宠物管理</template>
        <a-menu-item key="pet-config">宠物配置</a-menu-item>
        <a-menu-item key="pet-console">宠物控制台</a-menu-item>
        <a-menu-item key="pet-conversations">对话记录</a-menu-item>
      </a-sub-menu>

      <!-- 会员管理 -->
      <a-sub-menu key="members">
        <template #icon><icon-user-group /></template>
        <template #title>会员管理</template>
        <a-menu-item key="members">会员列表</a-menu-item>
        <a-menu-item key="members-cards">会员卡管理</a-menu-item>
        <a-menu-item key="members-coupons">优惠券</a-menu-item>
        <a-menu-item key="members-stores">门店管理</a-menu-item>
        <a-menu-item key="members-levels">会员等级</a-menu-item>
        <a-menu-item key="members-tags">标签管理</a-menu-item>
        <a-menu-item key="members-points">积分管理</a-menu-item>
        <a-menu-item key="members-promotions">促销活动</a-menu-item>
      </a-sub-menu>

      <!-- 知识库 -->
      <a-menu-item key="knowledge">
        <template #icon><icon-book /></template>
        <span>知识库</span>
      </a-menu-item>

      <!-- 主人档案 -->
      <a-menu-item key="owner-profile">
        <template #icon><icon-user /></template>
        <span>主人档案</span>
      </a-menu-item>

      <!-- 健康医疗 -->
      <a-sub-menu key="health">
        <template #icon><icon-heart /></template>
        <template #title>健康医疗</template>
        <a-menu-item key="health-sports">运动统计</a-menu-item>
        <a-menu-item key="health-reports">健康报告</a-menu-item>
        <a-menu-item key="health-warning">健康预警</a-menu-item>
        <a-menu-item key="health-sleep">睡眠分析</a-menu-item>
      </a-sub-menu>

      <!-- AI 功能 -->
      <a-sub-menu key="ai">
        <template #icon><icon-robot /></template>
        <template #title>AI 功能</template>
        <a-menu-item key="ai-behavior">行为分析</a-menu-item>
        <a-menu-item key="ai-emotion">情感识别</a-menu-item>
      </a-sub-menu>

      <!-- 应用管理 -->
      <a-sub-menu key="apps">
        <template #icon><icon-apps /></template>
        <template #title>应用管理</template>
        <a-menu-item key="apps-list">应用列表</a-menu-item>
        <a-menu-item key="apps-versions">应用版本</a-menu-item>
        <a-menu-item key="apps-distributions">应用分发</a-menu-item>
      </a-sub-menu>

      <!-- 系统管理 -->
      <a-sub-menu key="system">
        <template #icon><icon-settings /></template>
        <template #title>系统管理</template>
        <a-menu-item key="system-monitor">服务监控</a-menu-item>
        <a-menu-item key="system-logs">操作日志</a-menu-item>
        <a-menu-item key="policies-list">策略列表</a-menu-item>
        <a-menu-item key="policies-configs">策略配置</a-menu-item>
        <a-menu-item key="compliance-rules">合规规则</a-menu-item>
      </a-sub-menu>

      <!-- 租户管理 -->
      <a-sub-menu key="tenants">
        <template #icon><icon-tool /></template>
        <template #title>租户管理</template>
        <a-menu-item key="tenants-approval">租户入驻审核</a-menu-item>
        <a-menu-item key="tenants-management">租户系统管理</a-menu-item>
        <a-menu-item key="tenants-public-archives">公共档案</a-menu-item>
        <a-menu-item key="tenants-system-info">系统信息</a-menu-item>
      </a-sub-menu>

      <!-- 组织管理 -->
      <a-sub-menu key="org">
        <template #icon><icon-tree-list /></template>
        <template #title>组织管理</template>
        <a-menu-item key="org-companies">公司管理</a-menu-item>
        <a-menu-item key="org-departments">部门管理</a-menu-item>
        <a-menu-item key="org-posts">岗位管理</a-menu-item>
        <a-menu-item key="org-employees">员工管理</a-menu-item>
        <a-menu-item key="org-standard-positions">标准岗位</a-menu-item>
      </a-sub-menu>

      <!-- 多维权限 -->
      <a-sub-menu key="permissions">
        <template #icon><icon-safe /></template>
        <template #title>多维权限</template>
        <a-menu-item key="permissions-roles">角色管理</a-menu-item>
        <a-menu-item key="permissions-menus">菜单管理</a-menu-item>
        <a-menu-item key="permissions-groups">权限组</a-menu-item>
        <a-menu-item key="permissions-data-config">数据权限</a-menu-item>
        <a-menu-item key="permissions-api">API权限</a-menu-item>
      </a-sub-menu>

      <!-- 通知管理 -->
      <a-sub-menu key="notifications">
        <template #icon><icon-message /></template>
        <template #title>通知管理</template>
        <a-menu-item key="notifications-list">推送通知</a-menu-item>
        <a-menu-item key="notifications-announcements">公告管理</a-menu-item>
        <a-menu-item key="notifications-templates">通知模板</a-menu-item>
      </a-sub-menu>

      <!-- 第三方接入 -->
      <a-sub-menu key="integration">
        <template #icon><icon-link /></template>
        <template #title>第三方接入</template>
        <a-menu-item key="integration-pet-hospitals">宠物医院</a-menu-item>
      </a-sub-menu>

      <!-- 门户 -->
      <a-menu-item key="portal">
        <template #icon><icon-dashboard /></template>
        <span>工作台门户</span>
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
  IconTreeList,
  IconLink,
} from '@arco-design/web-vue/es/icon';

const router = useRouter();
const route = useRoute();
const appStore = useAppStore();

// 菜单状态
const selectedKey = ref<string[]>([]);
const openKeys = ref<string[]>([]);

// 折叠状态
const collapsed = computed(() => appStore.menuCollapse);

// 路由映射到菜单 key
const routeNameToMenuKey: Record<string, string> = {
  'Dashboard': 'dashboard',
  'Devices': 'devices',
  'OtaPackages': 'ota-packages',
  'OtaDeployments': 'ota-deployments',
  'MiniclawFirmwares': 'miniclaw-firmwares',
  'Alert': 'alerts',
  'AlertRules': 'alert-rules',
  'AlertList': 'alert-list',
  'AlertSettings': 'alert-settings',
  'PetConfig': 'pet-config',
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
  'Knowledge': 'knowledge',
  'OwnerProfile': 'owner-profile',
  'SystemMonitor': 'system-monitor',
  'SystemLogs': 'system-logs',
  'PolicyList': 'policies-list',
  'PolicyConfigs': 'policies-configs',
  'ComplianceRules': 'compliance-rules',
  'TenantsApproval': 'tenants-approval',
  'TenantsManagement': 'tenants-management',
  'TenantsPublicArchives': 'tenants-public-archives',
  'TenantsSystemInfo': 'tenants-system-info',
  'OrgCompanies': 'org-companies',
  'OrgDepartments': 'org-departments',
  'OrgPosts': 'org-posts',
  'OrgEmployees': 'org-employees',
  'OrgStandardPositions': 'org-standard-positions',
  'PermissionsRoles': 'permissions-roles',
  'PermissionsMenus': 'permissions-menus',
  'PermissionsGroups': 'permissions-groups',
  'PermissionsDataConfig': 'permissions-data-config',
  'PermissionsApi': 'permissions-api',
  'NotificationsList': 'notifications-list',
  'NotificationsAnnouncements': 'notifications-announcements',
  'NotificationsTemplates': 'notifications-templates',
  'AppsList': 'apps-list',
  'AppsVersions': 'apps-versions',
  'AppsDistributions': 'apps-distributions',
  'IntegrationPetHospitals': 'integration-pet-hospitals',
  'AiBehavior': 'ai-behavior',
  'AiEmotion': 'ai-emotion',
  'HealthWarning': 'health-warning',
  'HealthSports': 'health-sports',
  'HealthReports': 'health-reports',
  'HealthSleep': 'health-sleep',
  'Portal': 'portal',
};

// 监听路由变化，更新选中状态
watch(
  () => route.name,
  (name) => {
    if (name) {
      const menuKey = routeNameToMenuKey[name as string] || name as string;
      selectedKey.value = [menuKey];

      // 自动展开父级菜单
      const keyParts = menuKey.split('-');
      if (keyParts.length > 1) {
        openKeys.value = [keyParts[0]];
      }
    }
  },
  { immediate: true }
);

// 点击菜单
const handleMenuClick = ({ key }: { key: string }) => {
  router.push({ name: key });
};
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
