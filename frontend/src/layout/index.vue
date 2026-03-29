<!-- layout/index.vue - 主布局（Arco Design） -->
<template>
  <a-layout class="layout" :class="{ mobile: isMobile }">
    <!-- 侧边栏 -->
    <a-layout-sider
      v-if="renderMenu"
      class="layout-sider"
      :collapsed="collapsed"
      :collapsible="true"
      :width="220"
      :collapsed-width="64"
      :hide-trigger="true"
      breakpoint="xl"
      @collapse="handleCollapse"
    >
      <Menu />
    </a-layout-sider>

    <!-- 移动端抽屉菜单 -->
    <a-drawer
      v-if="isMobile"
      :visible="drawerVisible"
      placement="left"
      :footer="false"
      mask-closable
      :closable="false"
      @cancel="drawerVisible = false"
    >
      <Menu />
    </a-drawer>

    <!-- 主内容区 -->
    <a-layout class="layout-content-wrapper">
      <!-- 顶部导航 -->
      <a-layout-header class="layout-header">
        <NavBar />
      </a-layout-header>

      <!-- 面包屑 -->
      <div class="layout-breadcrumb" v-if="renderBreadcrumb">
        <Breadcrumb :items="breadcrumbItems" />
      </div>

      <!-- 内容 -->
      <a-layout-content class="layout-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </a-layout-content>

      <!-- 页脚 -->
      <a-layout-footer v-if="appStore.footer">
        <Footer />
      </a-layout-footer>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, computed, provide, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAppStore } from '@/store';
import { useUserStore } from '@/store';
import NavBar from '@/components/navbar/index.vue';
import Menu from '@/components/menu/index.vue';
import Footer from '@/components/footer/index.vue';
import Breadcrumb from '@/components/Breadcrumb.vue';

const appStore = useAppStore();
const route = useRoute();

// 响应式
const isMobile = ref(false);
const drawerVisible = ref(false);

// 计算属性
const collapsed = computed(() => appStore.menuCollapse);
const renderMenu = computed(() => appStore.menu);
const renderBreadcrumb = computed(() => appStore.tabBar);

// 面包屑映射
const routeBreadcrumbMap: Record<string, { label: string; parent?: string }> = {
  '/dashboard': { label: '设备大盘' },
  '/devices': { label: '设备列表' },
  '/ota/packages': { label: '固件包管理', parent: 'OTA升级' },
  '/ota/deployments': { label: '部署任务', parent: 'OTA升级' },
  '/miniclaw/firmwares': { label: '固件管理' },
  '/alerts/rules': { label: '告警规则', parent: '告警中心' },
  '/alerts/list': { label: '告警列表', parent: '告警中心' },
  '/alerts/settings': { label: '告警设置', parent: '告警中心' },
  '/pet': { label: '宠物配置', parent: '宠物管理' },
  '/pet/console': { label: '宠物控制台', parent: '宠物管理' },
  '/pet/conversations': { label: '对话记录', parent: '宠物管理' },
  '/members': { label: '会员列表' },
  '/knowledge': { label: '知识库' },
  '/owner/profile': { label: '主人档案' },
  '/health/exercise-stats': { label: '运动统计', parent: '健康医疗' },
  '/health/reports': { label: '健康报告', parent: '健康医疗' },
  '/health/warnings': { label: '健康预警', parent: '健康医疗' },
  '/health/sleep': { label: '睡眠分析', parent: '健康医疗' },
  '/ai/behavior': { label: '行为分析', parent: 'AI功能' },
  '/ai/emotion': { label: '情感识别', parent: 'AI功能' },
  '/emotion/recognition': { label: '情绪识别', parent: '情感计算' },
  '/emotion/response': { label: '情绪响应', parent: '情感计算' },
  '/emotion/logs': { label: '情绪日志', parent: '情感计算' },
  '/emotion/family-map': { label: '家庭情绪地图', parent: '情感计算' },
  '/digital-twin/vitals': { label: '实时生命体征', parent: '数字孪生' },
  '/digital-twin/prediction': { label: '行为预测', parent: '数字孪生' },
  '/digital-twin/playback': { label: '历史回放', parent: '数字孪生' },
  '/digital-twin/moments': { label: '精彩瞬间', parent: '数字孪生' },
  '/embodied/perception': { label: '环境感知', parent: '具身智能' },
  '/embodied/spatial': { label: '空间认知', parent: '具身智能' },
  '/embodied/motion': { label: '动作模仿', parent: '具身智能' },
  '/embodied/decision': { label: 'AI决策引擎', parent: '具身智能' },
  '/simulation/pet': { label: '虚拟宠物仿真', parent: '仿真测试' },
  '/simulation/test': { label: '自动化测试框架', parent: '仿真测试' },
  '/simulation/replay': { label: '系统回放', parent: '仿真测试' },
  '/apps/list': { label: '应用列表', parent: '应用管理' },
  '/apps/versions': { label: '应用版本', parent: '应用管理' },
  '/apps/distributions': { label: '应用分发', parent: '应用管理' },
  '/system/monitor': { label: '服务监控', parent: '系统管理' },
  '/system/logs': { label: '操作日志', parent: '系统管理' },
  '/policies/list': { label: '策略列表', parent: '系统管理' },
  '/policies/configs': { label: '策略配置', parent: '系统管理' },
  '/policies/compliance': { label: '合规规则', parent: '系统管理' },
  '/tenants/approval': { label: '租户入驻审核', parent: '租户管理' },
  '/tenants/management': { label: '租户系统管理', parent: '租户管理' },
  '/tenants/public-archives': { label: '公共档案', parent: '租户管理' },
  '/tenants/system-info': { label: '系统信息', parent: '租户管理' },
  '/org/companies': { label: '公司管理', parent: '组织管理' },
  '/org/departments': { label: '部门管理', parent: '组织管理' },
  '/org/posts': { label: '岗位管理', parent: '组织管理' },
  '/org/employees': { label: '员工管理', parent: '组织管理' },
  '/org/standard-positions': { label: '标准岗位', parent: '组织管理' },
  '/permissions/roles': { label: '角色管理', parent: '多维权限' },
  '/permissions/menus': { label: '菜单管理', parent: '多维权限' },
  '/permissions/groups': { label: '权限组', parent: '多维权限' },
  '/permissions/data-config': { label: '数据权限', parent: '多维权限' },
  '/permissions/api': { label: 'API权限', parent: '多维权限' },
  '/notifications/list': { label: '推送通知', parent: '通知管理' },
  '/notifications/announcements': { label: '公告管理', parent: '通知管理' },
  '/notifications/templates': { label: '通知模板', parent: '通知管理' },
  '/integration/pet-hospitals': { label: '宠物医院', parent: '第三方接入' },
  '/portal': { label: '工作台门户' },
  '/subscription': { label: '订阅管理' },
  '/webhooks': { label: 'Webhook列表', parent: 'Webhook' },
  '/webhooks/logs': { label: '调用日志', parent: 'Webhook' },
  '/developer': { label: '开发者API' },
  '/billing': { label: '用量计费' },
  '/devices/certificates': { label: '证书管理', parent: '设备增强' },
  '/devices/remote': { label: '远程控制', parent: '设备增强' },
  '/devices/geofence': { label: '地理围栏', parent: '设备增强' },
  '/devices/monitor': { label: '监控面板', parent: '设备增强' },
  '/devices/pairing': { label: '配对管理', parent: '设备增强' },
};

const breadcrumbItems = computed(() => {
  const path = route.path;
  const info = routeBreadcrumbMap[path];
  if (!info) return [{ label: '首页', href: '/' }];
  const items = [{ label: '首页', href: '/' }];
  if (info.parent) items.push({ label: info.parent });
  items.push({ label: info.label });
  return items;
});

// 折叠/展开
const handleCollapse = (val: boolean) => {
  appStore.updateSettings({ menuCollapse: val });
};

// 切换抽屉菜单
const toggleDrawer = () => {
  drawerVisible.value = !drawerVisible.value;
};

// 提供给子组件
provide('toggleDrawerMenu', toggleDrawer);

// 初始化
onMounted(() => {
  const checkMobile = () => {
    isMobile.value = window.innerWidth < 992;
  };

  checkMobile();
  window.addEventListener('resize', checkMobile);

  // 初始化用户信息
  const userStore = useUserStore();
  userStore.initFromStorage();
});
</script>

<style scoped lang="less">
@nav-size-height: 60px;

.layout {
  width: 100%;
  height: 100vh;

  &.mobile {
    .layout-sider {
      display: none;
    }
  }
}

.layout-sider {
  position: fixed;
  top: 0;
  left: 0;
  z-index: 100;
  height: 100%;
  background: var(--color-bg-2);
  border-right: 1px solid var(--color-border);
  transition: all 0.2s cubic-bezier(0.34, 0.69, 0.1, 1);

  :deep(.arco-layout-sider-children) {
    display: flex;
    flex-direction: column;
    height: 100%;
  }
}

.layout-content-wrapper {
  min-height: 100vh;
  transition: padding-left 0.2s cubic-bezier(0.34, 0.69, 0.1, 1);

  &:not(.mobile) {
    margin-left: 220px;

    .layout-sider:not(.arco-layout-sider-collapsed) ~ & {
      margin-left: 220px;
    }

    .layout-sider.arco-layout-sider-collapsed ~ & {
      margin-left: 64px;
    }
  }
}

.layout-header {
  position: sticky;
  top: 0;
  z-index: 99;
  height: @nav-size-height;
  background: var(--color-bg-2);
  border-bottom: 1px solid var(--color-border);
}

.layout-breadcrumb {
  padding: 8px 16px;
  background: var(--color-bg-2);
  border-bottom: 1px solid var(--color-border);
}

.layout-content {
  min-height: calc(100vh - @nav-size-height);
  padding: 16px;
  background: var(--color-bg-1);
  overflow-y: auto;
}

// 路由切换动画
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
