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
    {
      path: 'content/comments',
      name: 'ContentComments',
      component: () => import('@/views/ContentComments.vue'),
      meta: { locale: 'menu.content.comments', requiresAuth: true },
    },
    {
      path: 'content/versions',
      name: 'ContentVersions',
      component: () => import('@/views/ContentVersions.vue'),
      meta: { locale: 'menu.content.versions', requiresAuth: true },
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
      path: 'org/positions',
      name: 'OrgPositions',
      component: () => import('@/views/org/PositionList.vue'),
      meta: { locale: 'menu.org.positions', requiresAuth: true },
    },
    {
      path: 'org/standard',
      name: 'OrgStandard',
      component: () => import('@/views/org/StandardPositions.vue'),
      meta: { locale: 'menu.org.standard', requiresAuth: true },
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

    // ========== 会员管理 ==========
    {
      path: 'member/list',
      name: 'MemberList',
      component: () => import('@/views/members/MemberListView.vue'),
      meta: { locale: 'menu.members.list', requiresAuth: true },
    },
    {
      path: 'member/levels',
      name: 'MemberLevels',
      component: () => import('@/views/MemberLevels.vue'),
      meta: { locale: 'menu.members.levels', requiresAuth: true },
    },
    {
      path: 'member/points',
      name: 'MemberPoints',
      component: () => import('@/views/MemberPoints.vue'),
      meta: { locale: 'menu.members.points', requiresAuth: true },
    },
    {
      path: 'member/coupons',
      name: 'MemberCoupons',
      component: () => import('@/views/members/MemberCoupons.vue'),
      meta: { locale: 'menu.members.coupons', requiresAuth: true },
    },
    {
      path: 'member/tags',
      name: 'MemberTags',
      component: () => import('@/views/MemberTags.vue'),
      meta: { locale: 'menu.members.tags', requiresAuth: true },
    },
    {
      path: 'member/orders',
      name: 'MemberOrders',
      component: () => import('@/views/MemberOrders.vue'),
      meta: { locale: 'menu.members.orders', requiresAuth: true },
    },
    {
      path: 'member/stores',
      name: 'MemberStores',
      component: () => import('@/views/members/MemberStores.vue'),
      meta: { locale: 'menu.members.stores', requiresAuth: true },
    },
    {
      path: 'member/growth',
      name: 'MemberGrowth',
      component: () => import('@/views/MemberGrowth.vue'),
      meta: { locale: 'menu.members.growth', requiresAuth: true },
    },
    {
      path: 'member/recharge',
      name: 'MemberRecharge',
      component: () => import('@/views/MemberRecharge.vue'),
      meta: { locale: 'menu.members.recharge', requiresAuth: true },
    },
    {
      path: 'member/loyalty',
      name: 'MemberLoyalty',
      component: () => import('@/views/LoyaltyProgram.vue'),
      meta: { locale: 'menu.members.loyalty', requiresAuth: true },
    },
    {
      path: 'member/card',
      name: 'MemberCard',
      component: () => import('@/views/members/MemberCardView.vue'),
      meta: { locale: 'menu.members.card', requiresAuth: true },
    },
    {
      path: 'member/benefit',
      name: 'MemberBenefit',
      component: () => import('@/views/members/MemberBenefitsView.vue'),
      meta: { locale: 'menu.members.benefit', requiresAuth: true },
    },
    {
      path: 'member/article',
      name: 'MemberArticle',
      component: () => import('@/views/members/MemberArticlesView.vue'),
      meta: { locale: 'menu.members.article', requiresAuth: true },
    },
    {
      path: 'member/gift',
      name: 'MemberGift',
      component: () => import('@/views/members/MemberGiftView.vue'),
      meta: { locale: 'menu.members.gift', requiresAuth: true },
    },
    {
      path: 'member/reception',
      name: 'MemberReception',
      component: () => import('@/views/members/MemberReceptionView.vue'),
      meta: { locale: 'menu.members.reception', requiresAuth: true },
    },
    {
      path: 'member/settings',
      name: 'MemberSettings',
      component: () => import('@/views/members/MemberSettingsView.vue'),
      meta: { locale: 'menu.members.settings', requiresAuth: true },
    },
    {
      path: 'member/channels',
      name: 'MemberChannels',
      component: () => import('@/views/members/MemberChannelsView.vue'),
      meta: { locale: 'menu.members.channels', requiresAuth: true },
    },
    {
      path: 'member/promotion',
      name: 'MemberPromotion',
      component: () => import('@/views/members/MemberPromotions.vue'),
      meta: { locale: 'menu.members.promotion', requiresAuth: true },
    },

    // ========== 积分规则 ==========
    {
      path: 'member/points-rules',
      name: 'PointsRules',
      component: () => import('@/views/members/PointsRulesView.vue'),
      meta: { locale: 'menu.points.rules', requiresAuth: true },
    },
    {
      path: 'member/points-inventory',
      name: 'PointsInventory',
      component: () => import('@/views/members/PointsInventory.vue'),
      meta: { locale: 'menu.points.inventory', requiresAuth: true },
    },
    {
      path: 'member/points-records',
      name: 'PointsRecords',
      component: () => import('@/views/members/PointsRecords.vue'),
      meta: { locale: 'menu.points.records', requiresAuth: true },
    },
    {
      path: 'member/points-settings',
      name: 'PointsSettings',
      component: () => import('@/views/members/PointsSettingsView.vue'),
      meta: { locale: 'menu.points.settings', requiresAuth: true },
    },
    {
      path: 'member/points-exclude',
      name: 'PointsExclude',
      component: () => import('@/views/members/PointsExcludeView.vue'),
      meta: { locale: 'menu.points.exclude', requiresAuth: true },
    },

    // ========== 优惠券详情 ==========
    {
      path: 'coupon/coupon-view',
      name: 'CouponView',
      component: () => import('@/views/members/CouponView.vue'),
      meta: { locale: 'menu.coupon.view', requiresAuth: true },
    },
    {
      path: 'coupon/amount-discount',
      name: 'AmountDiscount',
      component: () => import('@/views/members/AmountDiscountView.vue'),
      meta: { locale: 'menu.coupon.amountDiscount', requiresAuth: true },
    },
    {
      path: 'coupon/amount-reduce',
      name: 'AmountReduce',
      component: () => import('@/views/members/AmountReduceView.vue'),
      meta: { locale: 'menu.coupon.amountReduce', requiresAuth: true },
    },
    {
      path: 'coupon/direct-reduce',
      name: 'DirectReduce',
      component: () => import('@/views/members/DirectReduceView.vue'),
      meta: { locale: 'menu.coupon.directReduce', requiresAuth: true },
    },
    {
      path: 'coupon/buy-gift',
      name: 'BuyGift',
      component: () => import('@/views/members/BuyGiftView.vue'),
      meta: { locale: 'menu.coupon.buyGift', requiresAuth: true },
    },
    {
      path: 'coupon/gift-records',
      name: 'GiftRecords',
      component: () => import('@/views/members/GiftRecordsView.vue'),
      meta: { locale: 'menu.coupon.giftRecords', requiresAuth: true },
    },
    {
      path: 'coupon/inventory',
      name: 'CouponInventory',
      component: () => import('@/views/members/CouponInventoryView.vue'),
      meta: { locale: 'menu.coupon.inventory', requiresAuth: true },
    },
    {
      path: 'coupon/messages',
      name: 'CouponMessages',
      component: () => import('@/views/members/CouponMessagesView.vue'),
      meta: { locale: 'menu.coupon.messages', requiresAuth: true },
    },
    {
      path: 'coupon/grant',
      name: 'CouponGrant',
      component: () => import('@/views/members/CouponGrantView.vue'),
      meta: { locale: 'menu.coupon.grant', requiresAuth: true },
    },

    // ========== 会员标签 ==========
    {
      path: 'tag/high-freq',
      name: 'HighFreqTag',
      component: () => import('@/views/members/HighFreqTagView.vue'),
      meta: { locale: 'menu.tag.highFreq', requiresAuth: true },
    },
    {
      path: 'tag/low-freq',
      name: 'LowFreqTag',
      component: () => import('@/views/members/LowFreqTagView.vue'),
      meta: { locale: 'menu.tag.lowFreq', requiresAuth: true },
    },
    {
      path: 'tag/interest',
      name: 'InterestTag',
      component: () => import('@/views/members/InterestTagView.vue'),
      meta: { locale: 'menu.tag.interest', requiresAuth: true },
    },
    {
      path: 'tag/auto-clean',
      name: 'TagAutoClean',
      component: () => import('@/views/members/TagAutoCleanView.vue'),
      meta: { locale: 'menu.tag.autoClean', requiresAuth: true },
    },
    {
      path: 'tag/report',
      name: 'TagReport',
      component: () => import('@/views/members/TagReportView.vue'),
      meta: { locale: 'menu.tag.report', requiresAuth: true },
    },

    // ========== 会员等级 ==========
    {
      path: 'member/level-rules',
      name: 'MemberLevelRules',
      component: () => import('@/views/members/MemberLevelRulesView.vue'),
      meta: { locale: 'menu.level.rules', requiresAuth: true },
    },
    {
      path: 'member/level-view',
      name: 'LevelView',
      component: () => import('@/views/members/LevelView.vue'),
      meta: { locale: 'menu.level.view', requiresAuth: true },
    },

    // ========== VIP专属 ==========
    {
      path: 'vip/exclusive',
      name: 'VipExclusive',
      component: () => import('@/views/members/VipExclusiveView.vue'),
      meta: { locale: 'menu.vip.exclusive', requiresAuth: true },
    },

    // ========== 微信小程序 ==========
    {
      path: 'wechat/settings',
      name: 'WechatSettings',
      component: () => import('@/views/members/WechatSettingsView.vue'),
      meta: { locale: 'menu.wechat.settings', requiresAuth: true },
    },
    {
      path: 'wechat/mini-program',
      name: 'MiniProgram',
      component: () => import('@/views/members/MiniProgramView.vue'),
      meta: { locale: 'menu.wechat.miniProgram', requiresAuth: true },
    },

    // ========== 打印机 ==========
    {
      path: 'printer/manage',
      name: 'PrinterManage',
      component: () => import('@/views/members/PrinterManageView.vue'),
      meta: { locale: 'menu.printer.manage', requiresAuth: true },
    },

    // ========== 岗位类型 ==========
    {
      path: 'occupation/types',
      name: 'OccupationTypes',
      component: () => import('@/views/members/OccupationTypesView.vue'),
      meta: { locale: 'menu.occupation.types', requiresAuth: true },
    },

    // ========== 签到记录 ==========
    {
      path: 'signin/records',
      name: 'SignInRecords',
      component: () => import('@/views/members/SignInRecordView.vue'),
      meta: { locale: 'menu.signin.records', requiresAuth: true },
    },

    // ========== 短信渠道 ==========
    {
      path: 'sms/channel',
      name: 'SmsChannel',
      component: () => import('@/views/members/SmsChannelView.vue'),
      meta: { locale: 'menu.sms.channel', requiresAuth: true },
    },
    {
      path: 'sms/template',
      name: 'SmsTemplate',
      component: () => import('@/views/members/SmsTemplateView.vue'),
      meta: { locale: 'menu.sms.template', requiresAuth: true },
    },

    // ========== 租户管理 ==========
    {
      path: 'tenant/application',
      name: 'TenantApplication',
      component: () => import('@/views/tenants/TenantApplication.vue'),
      meta: { locale: 'menu.tenant.application', requiresAuth: true },
    },
    {
      path: 'tenant/approval',
      name: 'TenantApproval',
      component: () => import('@/views/tenants/TenantApproval.vue'),
      meta: { locale: 'menu.tenant.approval', requiresAuth: true },
    },
    {
      path: 'tenant/management',
      name: 'TenantManagement',
      component: () => import('@/views/tenants/TenantManagement.vue'),
      meta: { locale: 'menu.tenant.management', requiresAuth: true },
    },
    {
      path: 'tenant/settings',
      name: 'TenantSettings',
      component: () => import('@/views/tenants/TenantSettings.vue'),
      meta: { locale: 'menu.tenant.settings', requiresAuth: true },
    },

    // ========== 系统管理 ==========
    {
      path: 'system/email-templates',
      name: 'SystemEmailTemplates',
      component: () => import('@/views/system/EmailTemplates.vue'),
      meta: { locale: 'menu.system.emailTemplates', requiresAuth: true },
    },
    {
      path: 'system/logs',
      name: 'SystemLogs',
      component: () => import('@/views/system/Logs.vue'),
      meta: { locale: 'menu.system.logs', requiresAuth: true },
    },
    {
      path: 'system/monitor',
      name: 'SystemMonitor',
      component: () => import('@/views/system/Monitor.vue'),
      meta: { locale: 'menu.system.monitor', requiresAuth: true },
    },

    // ========== 订阅管理 ==========
    {
      path: 'subscription/list',
      name: 'SubscriptionList',
      component: () => import('@/views/subscription/SubscriptionList.vue'),
      meta: { locale: 'menu.subscription.list', requiresAuth: true },
    },
  ],
};

export default BUSINESS;
