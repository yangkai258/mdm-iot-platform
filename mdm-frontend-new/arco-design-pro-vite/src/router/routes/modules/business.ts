import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

// 业务路由模块 - 涵盖所有菜单视图
const BUSINESS: AppRouteRecordRaw = {
  path: '/business',
  name: 'Business',
  component: DEFAULT_LAYOUT,
  meta: {
    requiresAuth: true,
    order: 2,
  },
  children: [
    // ========== 内容市场 ==========
    {
      path: 'content/app-store',
      name: 'ContentAppStore',
      component: () => import('@/views/content/AppStoreView.vue'),
      meta: { locale: 'menu.content.appStore', requiresAuth: true },
    },
    {
      path: 'content/library',
      name: 'ContentLibrary',
      component: () => import('@/views/content/ContentLibraryView.vue'),
      meta: { locale: 'menu.content.library', requiresAuth: true },
    },

    // ========== 动作市场 ==========
    {
      path: 'market/action',
      name: 'MarketAction',
      component: () => import('@/views/market/ActionMarketView.vue'),
      meta: { locale: 'menu.market.action', requiresAuth: true },
    },
    // ========== 广告活动 ==========
    {
      path: 'market/ad-campaign',
      name: 'MarketAdCampaign',
      component: () => import('@/views/market/AdCampaignView.vue'),
      meta: { locale: 'menu.market.adCampaign', requiresAuth: true },
    },
    // ========== 内容审核 ==========
    {
      path: 'market/content-review',
      name: 'MarketContentReview',
      component: () => import('@/views/market/ContentReviewView.vue'),
      meta: { locale: 'menu.market.contentReview', requiresAuth: true },
    },
    // ========== 优惠券池 ==========
    {
      path: 'market/coupon-pool',
      name: 'MarketCouponPool',
      component: () => import('@/views/market/CouponPoolView.vue'),
      meta: { locale: 'menu.market.couponPool', requiresAuth: true },
    },
    // ========== 表情包市场 ==========
    {
      path: 'market/emoticon',
      name: 'MarketEmoticon',
      component: () => import('@/views/market/EmoticonMarketView.vue'),
      meta: { locale: 'menu.market.emoticon', requiresAuth: true },
    },
    // ========== 充值规则 ==========
    {
      path: 'market/recharge',
      name: 'MarketRecharge',
      component: () => import('@/views/market/RechargeRulesView.vue'),
      meta: { locale: 'menu.market.recharge', requiresAuth: true },
    },
    // ========== 声音配置 ==========
    {
      path: 'market/voice',
      name: 'MarketVoice',
      component: () => import('@/views/market/VoiceConfigView.vue'),
      meta: { locale: 'menu.market.voice', requiresAuth: true },
    },

    // ========== 知识库 ==========
    {
      path: 'knowledge/base',
      name: 'KnowledgeBase',
      component: () => import('@/views/knowledge/KnowledgeList.vue'),
      meta: { locale: 'menu.knowledge', requiresAuth: true },
    },
    {
      path: 'knowledge/list',
      name: 'KnowledgeList',
      component: () => import('@/views/knowledge/KnowledgeList.vue'),
      meta: { locale: 'menu.knowledge.list', requiresAuth: true },
    },

    // ========== 组织管理 ==========
    {
      path: 'org/companies',
      name: 'OrgCompanies',
      component: () => import('@/views/org/CompanyList.vue'),
      meta: { locale: 'menu.org.companies', requiresAuth: true },
    },
    {
      path: 'org/departments',
      name: 'OrgDepartments',
      component: () => import('@/views/org/DepartmentList.vue'),
      meta: { locale: 'menu.org.departments', requiresAuth: true },
    },
    {
      path: 'org/employees',
      name: 'OrgEmployees',
      component: () => import('@/views/org/EmployeeList.vue'),
      meta: { locale: 'menu.org.employees', requiresAuth: true },
    },
    {
      path: 'org/posts',
      name: 'OrgPosts',
      component: () => import('@/views/org/PostList.vue'),
      meta: { locale: 'menu.org.posts', requiresAuth: true },
    },
    {
      path: 'org/chart',
      name: 'OrgChart',
      component: () => import('@/views/OrganizationChart.vue'),
      meta: { locale: 'menu.org', requiresAuth: true },
    },

    // ========== 权限管理 ==========
    {
      path: 'permission/api',
      name: 'PermissionApi',
      component: () => import('@/views/permissions/ApiPermissions.vue'),
      meta: { locale: 'menu.permissionManage.api', requiresAuth: true },
    },
    {
      path: 'permission/data-config',
      name: 'PermissionDataConfig',
      component: () => import('@/views/permissions/DataPermissionConfig.vue'),
      meta: { locale: 'menu.permissionManage.dataConfig', requiresAuth: true },
    },
    {
      path: 'permission/menus',
      name: 'PermissionMenus',
      component: () => import('@/views/permissions/Menus.vue'),
      meta: { locale: 'menu.permissionManage.menus', requiresAuth: true },
    },
    {
      path: 'permission/groups',
      name: 'PermissionGroups',
      component: () => import('@/views/permissions/PermissionGroups.vue'),
      meta: { locale: 'menu.permissionManage.groups', requiresAuth: true },
    },
    {
      path: 'permission/roles',
      name: 'PermissionRoles',
      component: () => import('@/views/permissions/Roles.vue'),
      meta: { locale: 'menu.permissionManage.roles', requiresAuth: true },
    },

    // ========== 店铺管理 ==========
    {
      path: 'store/list',
      name: 'StoreList',
      component: () => import('@/views/members/StoreView.vue'),
      meta: { locale: 'menu.store', requiresAuth: true },
    },
    {
      path: 'store/promotion',
      name: 'StorePromotion',
      component: () => import('@/views/members/PromotionTypesView.vue'),
      meta: { locale: 'menu.store.promotion', requiresAuth: true },
    },
    {
      path: 'store/locations',
      name: 'StoreLocations',
      component: () => import('@/views/members/StoreLocationsView.vue'),
      meta: { locale: 'menu.store.locations', requiresAuth: true },
    },
    {
      path: 'store/sources',
      name: 'StoreSources',
      component: () => import('@/views/members/StoreSourcesView.vue'),
      meta: { locale: 'menu.store.sources', requiresAuth: true },
    },
  ],
};

export default BUSINESS;
