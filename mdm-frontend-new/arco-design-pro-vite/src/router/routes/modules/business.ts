import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

// All business routes as top-level menu items
const BUSINESS: AppRouteRecordRaw[] = [
  // ========== 内容市场 ==========
  {
    path: '/content/app-store',
    name: 'ContentAppStore',
    component: () => import('@/views/content/AppStoreView.vue'),
    meta: { locale: 'menu.content.appStore', requiresAuth: true, order: 21 },
  },
  {
    path: '/content/library',
    name: 'ContentLibrary',
    component: () => import('@/views/content/ContentLibraryView.vue'),
    meta: { locale: 'menu.content.library', requiresAuth: true, order: 22 },
  },
  {
    path: '/content/comments',
    name: 'ContentComments',
    component: () => import('@/views/ContentComments.vue'),
    meta: { locale: 'menu.content.comments', requiresAuth: true, order: 23 },
  },
  {
    path: '/content/versions',
    name: 'ContentVersions',
    component: () => import('@/views/ContentVersions.vue'),
    meta: { locale: 'menu.content.versions', requiresAuth: true, order: 24 },
  },

  // ========== 动作市场 ==========
  {
    path: '/market/action',
    name: 'MarketAction',
    component: () => import('@/views/market/ActionMarketView.vue'),
    meta: { locale: 'menu.market.action', requiresAuth: true, order: 31 },
  },
  // ========== 广告活动 ==========
  {
    path: '/market/ad-campaign',
    name: 'MarketAdCampaign',
    component: () => import('@/views/market/AdCampaignView.vue'),
    meta: { locale: 'menu.market.adCampaign', requiresAuth: true, order: 32 },
  },
  // ========== 内容审核 ==========
  {
    path: '/market/content-review',
    name: 'MarketContentReview',
    component: () => import('@/views/market/ContentReviewView.vue'),
    meta: { locale: 'menu.market.contentReview', requiresAuth: true, order: 33 },
  },
  // ========== 优惠券池 ==========
  {
    path: '/market/coupon-pool',
    name: 'MarketCouponPool',
    component: () => import('@/views/market/CouponPoolView.vue'),
    meta: { locale: 'menu.market.couponPool', requiresAuth: true, order: 34 },
  },
  // ========== 表情包市场 ==========
  {
    path: '/market/emoticon',
    name: 'MarketEmoticon',
    component: () => import('@/views/market/EmoticonMarketView.vue'),
    meta: { locale: 'menu.market.emoticon', requiresAuth: true, order: 35 },
  },
  // ========== 充值规则 ==========
  {
    path: '/market/recharge',
    name: 'MarketRecharge',
    component: () => import('@/views/market/RechargeRulesView.vue'),
    meta: { locale: 'menu.market.recharge', requiresAuth: true, order: 36 },
  },
  // ========== 声音配置 ==========
  {
    path: '/market/voice',
    name: 'MarketVoice',
    component: () => import('@/views/market/VoiceConfigView.vue'),
    meta: { locale: 'menu.market.voice', requiresAuth: true, order: 37 },
  },

  // ========== 知识库 ==========
  {
    path: '/knowledge/base',
    name: 'KnowledgeBase',
    component: () => import('@/views/knowledge/KnowledgeList.vue'),
    meta: { locale: 'menu.knowledge', requiresAuth: true, order: 40 },
  },
  {
    path: '/knowledge/list',
    name: 'KnowledgeList',
    component: () => import('@/views/knowledge/KnowledgeList.vue'),
    meta: { locale: 'menu.knowledge.list', requiresAuth: true, order: 41 },
  },

  // ========== 组织管理 ==========
  {
    path: '/org/companies',
    name: 'OrgCompanies',
    component: () => import('@/views/org/CompanyList.vue'),
    meta: { locale: 'menu.org.companies', requiresAuth: true, order: 50 },
  },
  {
    path: '/org/departments',
    name: 'OrgDepartments',
    component: () => import('@/views/org/DepartmentList.vue'),
    meta: { locale: 'menu.org.departments', requiresAuth: true, order: 51 },
  },
  {
    path: '/org/employees',
    name: 'OrgEmployees',
    component: () => import('@/views/org/EmployeeList.vue'),
    meta: { locale: 'menu.org.employees', requiresAuth: true, order: 52 },
  },
  {
    path: '/org/posts',
    name: 'OrgPosts',
    component: () => import('@/views/org/PostList.vue'),
    meta: { locale: 'menu.org.posts', requiresAuth: true, order: 53 },
  },
  {
    path: '/org/positions',
    name: 'OrgPositions',
    component: () => import('@/views/org/PositionList.vue'),
    meta: { locale: 'menu.org.positions', requiresAuth: true, order: 54 },
  },
  {
    path: '/org/standard',
    name: 'OrgStandard',
    component: () => import('@/views/org/StandardPositions.vue'),
    meta: { locale: 'menu.org.standard', requiresAuth: true, order: 55 },
  },
  {
    path: '/org/chart',
    name: 'OrgChart',
    component: () => import('@/views/OrganizationChart.vue'),
    meta: { locale: 'menu.org', requiresAuth: true, order: 56 },
  },

  // ========== 权限管理 ==========
  {
    path: '/permission/api',
    name: 'PermissionApi',
    component: () => import('@/views/permissions/ApiPermissions.vue'),
    meta: { locale: 'menu.permissionManage.api', requiresAuth: true, order: 60 },
  },
  {
    path: '/permission/data-config',
    name: 'PermissionDataConfig',
    component: () => import('@/views/permissions/DataPermissionConfig.vue'),
    meta: { locale: 'menu.permissionManage.dataConfig', requiresAuth: true, order: 61 },
  },
  {
    path: '/permission/menus',
    name: 'PermissionMenus',
    component: () => import('@/views/permissions/Menus.vue'),
    meta: { locale: 'menu.permissionManage.menus', requiresAuth: true, order: 62 },
  },
  {
    path: '/permission/groups',
    name: 'PermissionGroups',
    component: () => import('@/views/permissions/PermissionGroups.vue'),
    meta: { locale: 'menu.permissionManage.groups', requiresAuth: true, order: 63 },
  },
  {
    path: '/permission/roles',
    name: 'PermissionRoles',
    component: () => import('@/views/permissions/Roles.vue'),
    meta: { locale: 'menu.permissionManage.roles', requiresAuth: true, order: 64 },
  },

  // ========== 店铺管理 ==========
  {
    path: '/store/list',
    name: 'StoreList',
    component: () => import('@/views/members/StoreView.vue'),
    meta: { locale: 'menu.store', requiresAuth: true, order: 70 },
  },
  {
    path: '/store/promotion',
    name: 'StorePromotion',
    component: () => import('@/views/members/PromotionTypesView.vue'),
    meta: { locale: 'menu.store.promotion', requiresAuth: true, order: 71 },
  },
  {
    path: '/store/locations',
    name: 'StoreLocations',
    component: () => import('@/views/members/StoreLocationsView.vue'),
    meta: { locale: 'menu.store.locations', requiresAuth: true, order: 72 },
  },
  {
    path: '/store/sources',
    name: 'StoreSources',
    component: () => import('@/views/members/StoreSourcesView.vue'),
    meta: { locale: 'menu.store.sources', requiresAuth: true, order: 73 },
  },

  // ========== 会员管理 ==========
  {
    path: '/member/list',
    name: 'MemberList',
    component: () => import('@/views/members/MemberListView.vue'),
    meta: { locale: 'menu.members.list', requiresAuth: true, order: 80 },
  },
  {
    path: '/member/levels',
    name: 'MemberLevels',
    component: () => import('@/views/MemberLevels.vue'),
    meta: { locale: 'menu.members.levels', requiresAuth: true, order: 81 },
  },
  {
    path: '/member/points',
    name: 'MemberPoints',
    component: () => import('@/views/MemberPoints.vue'),
    meta: { locale: 'menu.members.points', requiresAuth: true, order: 82 },
  },
  {
    path: '/member/coupons',
    name: 'MemberCoupons',
    component: () => import('@/views/members/MemberCoupons.vue'),
    meta: { locale: 'menu.members.coupons', requiresAuth: true, order: 83 },
  },
  {
    path: '/member/tags',
    name: 'MemberTags',
    component: () => import('@/views/MemberTags.vue'),
    meta: { locale: 'menu.members.tags', requiresAuth: true, order: 84 },
  },
  {
    path: '/member/orders',
    name: 'MemberOrders',
    component: () => import('@/views/MemberOrders.vue'),
    meta: { locale: 'menu.members.orders', requiresAuth: true, order: 85 },
  },
  {
    path: '/member/stores',
    name: 'MemberStores',
    component: () => import('@/views/members/MemberStores.vue'),
    meta: { locale: 'menu.members.stores', requiresAuth: true, order: 86 },
  },
  {
    path: '/member/growth',
    name: 'MemberGrowth',
    component: () => import('@/views/MemberGrowth.vue'),
    meta: { locale: 'menu.members.growth', requiresAuth: true, order: 87 },
  },
  {
    path: '/member/recharge',
    name: 'MemberRecharge',
    component: () => import('@/views/MemberRecharge.vue'),
    meta: { locale: 'menu.members.recharge', requiresAuth: true, order: 88 },
  },
  {
    path: '/member/loyalty',
    name: 'MemberLoyalty',
    component: () => import('@/views/LoyaltyProgram.vue'),
    meta: { locale: 'menu.members.loyalty', requiresAuth: true, order: 89 },
  },
  {
    path: '/member/card',
    name: 'MemberCard',
    component: () => import('@/views/members/MemberCardView.vue'),
    meta: { locale: 'menu.members.card', requiresAuth: true, order: 90 },
  },
  {
    path: '/member/benefit',
    name: 'MemberBenefit',
    component: () => import('@/views/members/MemberBenefitsView.vue'),
    meta: { locale: 'menu.members.benefit', requiresAuth: true, order: 91 },
  },
  {
    path: '/member/article',
    name: 'MemberArticle',
    component: () => import('@/views/members/MemberArticlesView.vue'),
    meta: { locale: 'menu.members.article', requiresAuth: true, order: 92 },
  },
  {
    path: '/member/gift',
    name: 'MemberGift',
    component: () => import('@/views/members/MemberGiftView.vue'),
    meta: { locale: 'menu.members.gift', requiresAuth: true, order: 93 },
  },
  {
    path: '/member/reception',
    name: 'MemberReception',
    component: () => import('@/views/members/MemberReceptionView.vue'),
    meta: { locale: 'menu.members.reception', requiresAuth: true, order: 94 },
  },
  {
    path: '/member/settings',
    name: 'MemberSettings',
    component: () => import('@/views/members/MemberSettingsView.vue'),
    meta: { locale: 'menu.members.settings', requiresAuth: true, order: 95 },
  },
  {
    path: '/member/channels',
    name: 'MemberChannels',
    component: () => import('@/views/members/MemberChannelsView.vue'),
    meta: { locale: 'menu.members.channels', requiresAuth: true, order: 96 },
  },
  {
    path: '/member/promotion',
    name: 'MemberPromotion',
    component: () => import('@/views/members/MemberPromotions.vue'),
    meta: { locale: 'menu.members.promotion', requiresAuth: true, order: 97 },
  },

  // ========== 积分规则 ==========
  {
    path: '/member/points-rules',
    name: 'PointsRules',
    component: () => import('@/views/members/PointsRulesView.vue'),
    meta: { locale: 'menu.points.rules', requiresAuth: true, order: 100 },
  },
  {
    path: '/member/points-inventory',
    name: 'PointsInventory',
    component: () => import('@/views/members/PointsInventory.vue'),
    meta: { locale: 'menu.points.inventory', requiresAuth: true, order: 101 },
  },
  {
    path: '/member/points-records',
    name: 'PointsRecords',
    component: () => import('@/views/members/PointsRecords.vue'),
    meta: { locale: 'menu.points.records', requiresAuth: true, order: 102 },
  },
  {
    path: '/member/points-settings',
    name: 'PointsSettings',
    component: () => import('@/views/members/PointsSettingsView.vue'),
    meta: { locale: 'menu.points.settings', requiresAuth: true, order: 103 },
  },
  {
    path: '/member/points-exclude',
    name: 'PointsExclude',
    component: () => import('@/views/members/PointsExcludeView.vue'),
    meta: { locale: 'menu.points.exclude', requiresAuth: true, order: 104 },
  },

  // ========== 优惠券详情 ==========
  {
    path: '/coupon/coupon-view',
    name: 'CouponView',
    component: () => import('@/views/members/CouponView.vue'),
    meta: { locale: 'menu.coupon.view', requiresAuth: true, order: 110 },
  },
  {
    path: '/coupon/amount-discount',
    name: 'AmountDiscount',
    component: () => import('@/views/members/AmountDiscountView.vue'),
    meta: { locale: 'menu.coupon.amountDiscount', requiresAuth: true, order: 111 },
  },
  {
    path: '/coupon/amount-reduce',
    name: 'AmountReduce',
    component: () => import('@/views/members/AmountReduceView.vue'),
    meta: { locale: 'menu.coupon.amountReduce', requiresAuth: true, order: 112 },
  },
  {
    path: '/coupon/direct-reduce',
    name: 'DirectReduce',
    component: () => import('@/views/members/DirectReduceView.vue'),
    meta: { locale: 'menu.coupon.directReduce', requiresAuth: true, order: 113 },
  },
  {
    path: '/coupon/buy-gift',
    name: 'BuyGift',
    component: () => import('@/views/members/BuyGiftView.vue'),
    meta: { locale: 'menu.coupon.buyGift', requiresAuth: true, order: 114 },
  },
  {
    path: '/coupon/gift-records',
    name: 'GiftRecords',
    component: () => import('@/views/members/GiftRecordsView.vue'),
    meta: { locale: 'menu.coupon.giftRecords', requiresAuth: true, order: 115 },
  },
  {
    path: '/coupon/inventory',
    name: 'CouponInventory',
    component: () => import('@/views/members/CouponInventoryView.vue'),
    meta: { locale: 'menu.coupon.inventory', requiresAuth: true, order: 116 },
  },
  {
    path: '/coupon/messages',
    name: 'CouponMessages',
    component: () => import('@/views/members/CouponMessagesView.vue'),
    meta: { locale: 'menu.coupon.messages', requiresAuth: true, order: 117 },
  },
  {
    path: '/coupon/grant',
    name: 'CouponGrant',
    component: () => import('@/views/members/CouponGrantView.vue'),
    meta: { locale: 'menu.coupon.grant', requiresAuth: true, order: 118 },
  },

  // ========== 会员标签 ==========
  {
    path: '/tag/high-freq',
    name: 'HighFreqTag',
    component: () => import('@/views/members/HighFreqTagView.vue'),
    meta: { locale: 'menu.tag.highFreq', requiresAuth: true, order: 120 },
  },
  {
    path: '/tag/low-freq',
    name: 'LowFreqTag',
    component: () => import('@/views/members/LowFreqTagView.vue'),
    meta: { locale: 'menu.tag.lowFreq', requiresAuth: true, order: 121 },
  },
  {
    path: '/tag/interest',
    name: 'InterestTag',
    component: () => import('@/views/members/InterestTagView.vue'),
    meta: { locale: 'menu.tag.interest', requiresAuth: true, order: 122 },
  },
  {
    path: '/tag/auto-clean',
    name: 'TagAutoClean',
    component: () => import('@/views/members/TagAutoCleanView.vue'),
    meta: { locale: 'menu.tag.autoClean', requiresAuth: true, order: 123 },
  },
  {
    path: '/tag/report',
    name: 'TagReport',
    component: () => import('@/views/members/TagReportView.vue'),
    meta: { locale: 'menu.tag.report', requiresAuth: true, order: 124 },
  },

  // ========== 会员等级 ==========
  {
    path: '/member/level-rules',
    name: 'MemberLevelRules',
    component: () => import('@/views/members/MemberLevelRulesView.vue'),
    meta: { locale: 'menu.level.rules', requiresAuth: true, order: 130 },
  },
  {
    path: '/member/level-view',
    name: 'LevelView',
    component: () => import('@/views/members/LevelView.vue'),
    meta: { locale: 'menu.level.view', requiresAuth: true, order: 131 },
  },

  // ========== VIP专属 ==========
  {
    path: '/vip/exclusive',
    name: 'VipExclusive',
    component: () => import('@/views/members/VipExclusiveView.vue'),
    meta: { locale: 'menu.vip.exclusive', requiresAuth: true, order: 140 },
  },

  // ========== 微信小程序 ==========
  {
    path: '/wechat/settings',
    name: 'WechatSettings',
    component: () => import('@/views/members/WechatSettingsView.vue'),
    meta: { locale: 'menu.wechat.settings', requiresAuth: true, order: 150 },
  },
  {
    path: '/wechat/mini-program',
    name: 'MiniProgram',
    component: () => import('@/views/members/MiniProgramView.vue'),
    meta: { locale: 'menu.wechat.miniProgram', requiresAuth: true, order: 151 },
  },

  // ========== 打印机 ==========
  {
    path: '/printer/manage',
    name: 'PrinterManage',
    component: () => import('@/views/members/PrinterManageView.vue'),
    meta: { locale: 'menu.printer.manage', requiresAuth: true, order: 160 },
  },

  // ========== 岗位类型 ==========
  {
    path: '/occupation/types',
    name: 'OccupationTypes',
    component: () => import('@/views/members/OccupationTypesView.vue'),
    meta: { locale: 'menu.occupation.types', requiresAuth: true, order: 170 },
  },

  // ========== 签到记录 ==========
  {
    path: '/signin/records',
    name: 'SignInRecords',
    component: () => import('@/views/members/SignInRecordView.vue'),
    meta: { locale: 'menu.signin.records', requiresAuth: true, order: 180 },
  },

  // ========== 短信渠道 ==========
  {
    path: '/sms/channel',
    name: 'SmsChannel',
    component: () => import('@/views/members/SmsChannelView.vue'),
    meta: { locale: 'menu.sms.channel', requiresAuth: true, order: 190 },
  },
  {
    path: '/sms/template',
    name: 'SmsTemplate',
    component: () => import('@/views/members/SmsTemplateView.vue'),
    meta: { locale: 'menu.sms.template', requiresAuth: true, order: 191 },
  },

  // ========== 租户管理 ==========
  {
    path: '/tenant/application',
    name: 'TenantApplication',
    component: () => import('@/views/tenants/TenantApplication.vue'),
    meta: { locale: 'menu.tenant.application', requiresAuth: true, order: 200 },
  },
  {
    path: '/tenant/approval',
    name: 'TenantApproval',
    component: () => import('@/views/tenants/TenantApproval.vue'),
    meta: { locale: 'menu.tenant.approval', requiresAuth: true, order: 201 },
  },
  {
    path: '/tenant/management',
    name: 'TenantManagement',
    component: () => import('@/views/tenants/TenantManagement.vue'),
    meta: { locale: 'menu.tenant.management', requiresAuth: true, order: 202 },
  },
  {
    path: '/tenant/settings',
    name: 'TenantSettings',
    component: () => import('@/views/tenants/TenantSettings.vue'),
    meta: { locale: 'menu.tenant.settings', requiresAuth: true, order: 203 },
  },

  // ========== 系统管理 ==========
  {
    path: '/system/email-templates',
    name: 'SystemEmailTemplates',
    component: () => import('@/views/system/EmailTemplates.vue'),
    meta: { locale: 'menu.system.emailTemplates', requiresAuth: true, order: 210 },
  },
  {
    path: '/system/logs',
    name: 'SystemLogs',
    component: () => import('@/views/system/Logs.vue'),
    meta: { locale: 'menu.system.logs', requiresAuth: true, order: 211 },
  },
  {
    path: '/system/monitor',
    name: 'SystemMonitor',
    component: () => import('@/views/system/Monitor.vue'),
    meta: { locale: 'menu.system.monitor', requiresAuth: true, order: 212 },
  },

  // ========== 订阅管理 ==========
  {
    path: '/subscription/list',
    name: 'SubscriptionList',
    component: () => import('@/views/subscription/SubscriptionList.vue'),
    meta: { locale: 'menu.subscription.list', requiresAuth: true, order: 220 },
  },
];

export default BUSINESS;
