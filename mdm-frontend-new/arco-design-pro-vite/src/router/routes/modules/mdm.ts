import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

// ========== 会员管理 ==========
const MemberManage: AppRouteRecordRaw = {
  path: '/member',
  name: 'MemberManage',
  component: DEFAULT_LAYOUT,
  redirect: '/member/list',
  meta: { locale: 'menu.members', requiresAuth: true, order: 10 },
  children: [
    { path: 'list', name: 'MemberList', component: () => import('@/views/members/MemberListView.vue'), meta: { locale: 'menu.members.list', requiresAuth: true } },
    { path: 'levels', name: 'MemberLevels', component: () => import('@/views/MemberLevels.vue'), meta: { locale: 'menu.members.levels', requiresAuth: true } },
    { path: 'points', name: 'MemberPoints', component: () => import('@/views/MemberPoints.vue'), meta: { locale: 'menu.members.points', requiresAuth: true } },
    { path: 'coupons', name: 'MemberCoupons', component: () => import('@/views/members/MemberCoupons.vue'), meta: { locale: 'menu.members.coupons', requiresAuth: true } },
    { path: 'tags', name: 'MemberTags', component: () => import('@/views/MemberTags.vue'), meta: { locale: 'menu.members.tags', requiresAuth: true } },
    { path: 'orders', name: 'MemberOrders', component: () => import('@/views/MemberOrders.vue'), meta: { locale: 'menu.members.orders', requiresAuth: true } },
    { path: 'stores', name: 'MemberStores', component: () => import('@/views/members/MemberStores.vue'), meta: { locale: 'menu.members.stores', requiresAuth: true } },
    { path: 'growth', name: 'MemberGrowth', component: () => import('@/views/MemberGrowth.vue'), meta: { locale: 'menu.members.growth', requiresAuth: true } },
    { path: 'recharge', name: 'MemberRecharge', component: () => import('@/views/MemberRecharge.vue'), meta: { locale: 'menu.members.recharge', requiresAuth: true } },
    { path: 'loyalty', name: 'MemberLoyalty', component: () => import('@/views/LoyaltyProgram.vue'), meta: { locale: 'menu.members.loyalty', requiresAuth: true } },
    { path: 'card', name: 'MemberCard', component: () => import('@/views/members/MemberCardView.vue'), meta: { locale: 'menu.members.card', requiresAuth: true } },
    { path: 'benefit', name: 'MemberBenefit', component: () => import('@/views/members/MemberBenefitsView.vue'), meta: { locale: 'menu.members.benefit', requiresAuth: true } },
    { path: 'article', name: 'MemberArticle', component: () => import('@/views/members/MemberArticlesView.vue'), meta: { locale: 'menu.members.article', requiresAuth: true } },
    { path: 'gift', name: 'MemberGift', component: () => import('@/views/members/MemberGiftView.vue'), meta: { locale: 'menu.members.gift', requiresAuth: true } },
    { path: 'reception', name: 'MemberReception', component: () => import('@/views/members/MemberReceptionView.vue'), meta: { locale: 'menu.members.reception', requiresAuth: true } },
    { path: 'settings', name: 'MemberSettings', component: () => import('@/views/members/MemberSettingsView.vue'), meta: { locale: 'menu.members.settings', requiresAuth: true } },
    { path: 'channels', name: 'MemberChannels', component: () => import('@/views/members/MemberChannelsView.vue'), meta: { locale: 'menu.members.channels', requiresAuth: true } },
    { path: 'promotion', name: 'MemberPromotion', component: () => import('@/views/members/MemberPromotions.vue'), meta: { locale: 'menu.members.promotion', requiresAuth: true } },
    // 积分规则
    { path: 'points-rules', name: 'PointsRules', component: () => import('@/views/members/PointsRulesView.vue'), meta: { locale: 'menu.points.rules', requiresAuth: true } },
    { path: 'points-inventory', name: 'PointsInventory', component: () => import('@/views/members/PointsInventory.vue'), meta: { locale: 'menu.points.inventory', requiresAuth: true } },
    { path: 'points-records', name: 'PointsRecords', component: () => import('@/views/members/PointsRecords.vue'), meta: { locale: 'menu.points.records', requiresAuth: true } },
    { path: 'points-settings', name: 'PointsSettings', component: () => import('@/views/members/PointsSettingsView.vue'), meta: { locale: 'menu.points.settings', requiresAuth: true } },
    { path: 'points-exclude', name: 'PointsExclude', component: () => import('@/views/members/PointsExcludeView.vue'), meta: { locale: 'menu.points.exclude', requiresAuth: true } },
    // 会员等级
    { path: 'level-rules', name: 'MemberLevelRules', component: () => import('@/views/members/MemberLevelRulesView.vue'), meta: { locale: 'menu.level.rules', requiresAuth: true } },
    { path: 'level-view', name: 'LevelView', component: () => import('@/views/members/LevelView.vue'), meta: { locale: 'menu.level.view', requiresAuth: true } },
    // 会员标签
    { path: 'tag/high-freq', name: 'HighFreqTag', component: () => import('@/views/members/HighFreqTagView.vue'), meta: { locale: 'menu.tag.highFreq', requiresAuth: true } },
    { path: 'tag/low-freq', name: 'LowFreqTag', component: () => import('@/views/members/LowFreqTagView.vue'), meta: { locale: 'menu.tag.lowFreq', requiresAuth: true } },
    { path: 'tag/interest', name: 'InterestTag', component: () => import('@/views/members/InterestTagView.vue'), meta: { locale: 'menu.tag.interest', requiresAuth: true } },
    { path: 'tag/auto-clean', name: 'TagAutoClean', component: () => import('@/views/members/TagAutoCleanView.vue'), meta: { locale: 'menu.tag.autoClean', requiresAuth: true } },
    { path: 'tag/report', name: 'TagReport', component: () => import('@/views/members/TagReportView.vue'), meta: { locale: 'menu.tag.report', requiresAuth: true } },
  ],
};

// ========== 设备管理 ==========
const DeviceManage: AppRouteRecordRaw = {
  path: '/device',
  name: 'DeviceManage',
  component: DEFAULT_LAYOUT,
  redirect: '/device/list',
  meta: { locale: 'menu.devices', requiresAuth: true, order: 20 },
  children: [
    { path: 'list', name: 'DeviceList', component: () => import('@/views/DeviceList.vue'), meta: { locale: 'menu.devices.list', requiresAuth: true } },
    { path: 'detail/:id', name: 'DeviceDetail', component: () => import('@/views/DeviceDetail.vue'), meta: { locale: 'menu.devices.detail', requiresAuth: true, hideInMenu: true } },
    { path: 'pairing', name: 'DevicePairing', component: () => import('@/views/DevicePairing.vue'), meta: { locale: 'menu.devices.pairing', requiresAuth: true } },
    { path: 'groups', name: 'DeviceGroups', component: () => import('@/views/DeviceGroups.vue'), meta: { locale: 'menu.devices.groups', requiresAuth: true } },
    { path: 'commands', name: 'DeviceCommands', component: () => import('@/views/DeviceCommands.vue'), meta: { locale: 'menu.devices.commands', requiresAuth: true } },
    { path: 'geofence', name: 'DeviceGeofence', component: () => import('@/views/DeviceGeofence.vue'), meta: { locale: 'menu.devices.geofence', requiresAuth: true } },
    { path: 'certificates', name: 'DeviceCertificates', component: () => import('@/views/DeviceCertificates.vue'), meta: { locale: 'menu.deviceManage.certificates', requiresAuth: true } },
    { path: 'monitor-panel', name: 'DeviceMonitorPanel', component: () => import('@/views/DeviceMonitorPanel.vue'), meta: { locale: 'menu.deviceManage.monitorPanel', requiresAuth: true } },
    { path: 'remote-control', name: 'DeviceRemoteControl', component: () => import('@/views/DeviceRemoteControl.vue'), meta: { locale: 'menu.deviceManage.remoteControl', requiresAuth: true } },
  ],
};

// ========== 内容市场 ==========
const ContentMarket: AppRouteRecordRaw = {
  path: '/content',
  name: 'ContentMarket',
  component: DEFAULT_LAYOUT,
  redirect: '/content/app-store',
  meta: { locale: 'menu.content', requiresAuth: true, order: 30 },
  children: [
    { path: 'app-store', name: 'ContentAppStore', component: () => import('@/views/content/AppStoreView.vue'), meta: { locale: 'menu.content.appStore', requiresAuth: true } },
    { path: 'library', name: 'ContentLibrary', component: () => import('@/views/content/ContentLibraryView.vue'), meta: { locale: 'menu.content.library', requiresAuth: true } },
    { path: 'comments', name: 'ContentComments', component: () => import('@/views/ContentComments.vue'), meta: { locale: 'menu.content.comments', requiresAuth: true } },
    { path: 'versions', name: 'ContentVersions', component: () => import('@/views/ContentVersions.vue'), meta: { locale: 'menu.content.versions', requiresAuth: true } },
  ],
};

// ========== 营销管理 ==========
const Marketing: AppRouteRecordRaw = {
  path: '/marketing',
  name: 'Marketing',
  component: DEFAULT_LAYOUT,
  redirect: '/marketing/redpacket',
  meta: { locale: 'menu.marketing', requiresAuth: true, order: 40 },
  children: [
    { path: 'redpacket', name: 'RedpacketView', component: () => import('@/views/RedpacketView.vue'), meta: { locale: 'menu.marketing.redpacket', requiresAuth: true } },
    { path: 'temp-coupon', name: 'TempCouponGrantView', component: () => import('@/views/TempCouponGrantView.vue'), meta: { locale: 'menu.marketing.tempCoupon', requiresAuth: true } },
    { path: 'temp-coupons', name: 'TempCoupons', component: () => import('@/views/TempCoupons.vue'), meta: { locale: 'menu.marketing.tempCoupons', requiresAuth: true } },
    { path: 'temp-redpacket', name: 'TempRedpacketView', component: () => import('@/views/TempRedpacketView.vue'), meta: { locale: 'menu.marketing.tempRedpacket', requiresAuth: true } },
  ],
};

// ========== 动作市场 ==========
const ActionMarket: AppRouteRecordRaw = {
  path: '/market/action', name: 'MarketAction',
  component: () => import('@/views/market/ActionMarketView.vue'),
  meta: { locale: 'menu.market.action', requiresAuth: true, order: 50 },
};

// ========== 广告活动 ==========
const AdCampaign: AppRouteRecordRaw = {
  path: '/market/ad-campaign', name: 'MarketAdCampaign',
  component: () => import('@/views/market/AdCampaignView.vue'),
  meta: { locale: 'menu.market.adCampaign', requiresAuth: true, order: 51 },
};

// ========== 内容审核 ==========
const ContentReview: AppRouteRecordRaw = {
  path: '/market/content-review', name: 'MarketContentReview',
  component: () => import('@/views/market/ContentReviewView.vue'),
  meta: { locale: 'menu.market.contentReview', requiresAuth: true, order: 52 },
};

// ========== 优惠券池 ==========
const CouponPool: AppRouteRecordRaw = {
  path: '/market/coupon-pool', name: 'MarketCouponPool',
  component: () => import('@/views/market/CouponPoolView.vue'),
  meta: { locale: 'menu.market.couponPool', requiresAuth: true, order: 53 },
};

// ========== 表情包市场 ==========
const EmoticonMarket: AppRouteRecordRaw = {
  path: '/market/emoticon', name: 'MarketEmoticon',
  component: () => import('@/views/market/EmoticonMarketView.vue'),
  meta: { locale: 'menu.market.emoticon', requiresAuth: true, order: 54 },
};

// ========== 充值规则 ==========
const RechargeRules: AppRouteRecordRaw = {
  path: '/market/recharge', name: 'MarketRecharge',
  component: () => import('@/views/market/RechargeRulesView.vue'),
  meta: { locale: 'menu.market.recharge', requiresAuth: true, order: 55 },
};

// ========== 声音配置 ==========
const VoiceConfig: AppRouteRecordRaw = {
  path: '/market/voice', name: 'MarketVoice',
  component: () => import('@/views/market/VoiceConfigView.vue'),
  meta: { locale: 'menu.market.voice', requiresAuth: true, order: 56 },
};

// ========== 知识库 ==========
const KnowledgeBase: AppRouteRecordRaw = {
  path: '/knowledge', name: 'KnowledgeBase',
  component: DEFAULT_LAYOUT, redirect: '/knowledge/base',
  meta: { locale: 'menu.knowledge', requiresAuth: true, order: 60 },
  children: [
    { path: 'base', name: 'KnowledgeBaseMain', component: () => import('@/views/knowledge/KnowledgeList.vue'), meta: { locale: 'menu.knowledge', requiresAuth: true } },
    { path: 'list', name: 'KnowledgeList', component: () => import('@/views/knowledge/KnowledgeList.vue'), meta: { locale: 'menu.knowledge.list', requiresAuth: true } },
  ],
};

// ========== 组织管理 ==========
const OrgManagement: AppRouteRecordRaw = {
  path: '/org', name: 'OrgManagement',
  component: DEFAULT_LAYOUT, redirect: '/org/chart',
  meta: { locale: 'menu.org', requiresAuth: true, order: 70 },
  children: [
    { path: 'chart', name: 'OrgChart', component: () => import('@/views/OrganizationChart.vue'), meta: { locale: 'menu.org', requiresAuth: true } },
    { path: 'companies', name: 'OrgCompanies', component: () => import('@/views/org/CompanyList.vue'), meta: { locale: 'menu.org.companies', requiresAuth: true } },
    { path: 'departments', name: 'OrgDepartments', component: () => import('@/views/org/DepartmentList.vue'), meta: { locale: 'menu.org.departments', requiresAuth: true } },
    { path: 'employees', name: 'OrgEmployees', component: () => import('@/views/org/EmployeeList.vue'), meta: { locale: 'menu.org.employees', requiresAuth: true } },
    { path: 'posts', name: 'OrgPosts', component: () => import('@/views/org/PostList.vue'), meta: { locale: 'menu.org.posts', requiresAuth: true } },
    { path: 'positions', name: 'OrgPositions', component: () => import('@/views/org/PositionList.vue'), meta: { locale: 'menu.org.positions', requiresAuth: true } },
    { path: 'standard', name: 'OrgStandard', component: () => import('@/views/org/StandardPositions.vue'), meta: { locale: 'menu.org.standard', requiresAuth: true } },
  ],
};

// ========== 权限管理 ==========
const PermissionManage: AppRouteRecordRaw = {
  path: '/permission', name: 'PermissionManage',
  component: DEFAULT_LAYOUT, redirect: '/permission/api',
  meta: { locale: 'menu.permissionManage', requiresAuth: true, order: 80 },
  children: [
    { path: 'api', name: 'PermissionApi', component: () => import('@/views/permissions/ApiPermissions.vue'), meta: { locale: 'menu.permissionManage.api', requiresAuth: true } },
    { path: 'data-config', name: 'PermissionDataConfig', component: () => import('@/views/permissions/DataPermissionConfig.vue'), meta: { locale: 'menu.permissionManage.dataConfig', requiresAuth: true } },
    { path: 'menus', name: 'PermissionMenus', component: () => import('@/views/permissions/Menus.vue'), meta: { locale: 'menu.permissionManage.menus', requiresAuth: true } },
    { path: 'groups', name: 'PermissionGroups', component: () => import('@/views/permissions/PermissionGroups.vue'), meta: { locale: 'menu.permissionManage.groups', requiresAuth: true } },
    { path: 'roles', name: 'PermissionRoles', component: () => import('@/views/permissions/Roles.vue'), meta: { locale: 'menu.permissionManage.roles', requiresAuth: true } },
  ],
};

// ========== 店铺管理 ==========
const StoreManage: AppRouteRecordRaw = {
  path: '/store', name: 'StoreManage',
  component: DEFAULT_LAYOUT, redirect: '/store/list',
  meta: { locale: 'menu.store', requiresAuth: true, order: 90 },
  children: [
    { path: 'list', name: 'StoreList', component: () => import('@/views/members/StoreView.vue'), meta: { locale: 'menu.store', requiresAuth: true } },
    { path: 'promotion', name: 'StorePromotion', component: () => import('@/views/members/PromotionTypesView.vue'), meta: { locale: 'menu.store.promotion', requiresAuth: true } },
    { path: 'locations', name: 'StoreLocations', component: () => import('@/views/members/StoreLocationsView.vue'), meta: { locale: 'menu.store.locations', requiresAuth: true } },
    { path: 'sources', name: 'StoreSources', component: () => import('@/views/members/StoreSourcesView.vue'), meta: { locale: 'menu.store.sources', requiresAuth: true } },
  ],
};

// ========== 优惠券 ==========
const CouponManage: AppRouteRecordRaw = {
  path: '/coupon', name: 'CouponManage',
  component: DEFAULT_LAYOUT, redirect: '/coupon/coupon-view',
  meta: { locale: 'menu.members.coupon', requiresAuth: true, order: 100 },
  children: [
    { path: 'coupon-view', name: 'CouponView', component: () => import('@/views/members/CouponView.vue'), meta: { locale: 'menu.coupon.view', requiresAuth: true } },
    { path: 'amount-discount', name: 'AmountDiscount', component: () => import('@/views/members/AmountDiscountView.vue'), meta: { locale: 'menu.coupon.amountDiscount', requiresAuth: true } },
    { path: 'amount-reduce', name: 'AmountReduce', component: () => import('@/views/members/AmountReduceView.vue'), meta: { locale: 'menu.coupon.amountReduce', requiresAuth: true } },
    { path: 'direct-reduce', name: 'DirectReduce', component: () => import('@/views/members/DirectReduceView.vue'), meta: { locale: 'menu.coupon.directReduce', requiresAuth: true } },
    { path: 'buy-gift', name: 'BuyGift', component: () => import('@/views/members/BuyGiftView.vue'), meta: { locale: 'menu.coupon.buyGift', requiresAuth: true } },
    { path: 'gift-records', name: 'GiftRecords', component: () => import('@/views/members/GiftRecordsView.vue'), meta: { locale: 'menu.coupon.giftRecords', requiresAuth: true } },
    { path: 'inventory', name: 'CouponInventory', component: () => import('@/views/members/CouponInventoryView.vue'), meta: { locale: 'menu.coupon.inventory', requiresAuth: true } },
    { path: 'messages', name: 'CouponMessages', component: () => import('@/views/members/CouponMessagesView.vue'), meta: { locale: 'menu.coupon.messages', requiresAuth: true } },
    { path: 'grant', name: 'CouponGrant', component: () => import('@/views/members/CouponGrantView.vue'), meta: { locale: 'menu.coupon.grant', requiresAuth: true } },
  ],
};

// ========== VIP专属 ==========
const VipExclusive: AppRouteRecordRaw = {
  path: '/vip/exclusive', name: 'VipExclusive',
  component: () => import('@/views/members/VipExclusiveView.vue'),
  meta: { locale: 'menu.vip.exclusive', requiresAuth: true, order: 110 },
};

// ========== 微信小程序 ==========
const WechatManage: AppRouteRecordRaw = {
  path: '/wechat', name: 'WechatManage',
  component: DEFAULT_LAYOUT, redirect: '/wechat/settings',
  meta: { locale: 'menu.wechat.settings', requiresAuth: true, order: 120 },
  children: [
    { path: 'settings', name: 'WechatSettings', component: () => import('@/views/members/WechatSettingsView.vue'), meta: { locale: 'menu.wechat.settings', requiresAuth: true } },
    { path: 'mini-program', name: 'MiniProgram', component: () => import('@/views/members/MiniProgramView.vue'), meta: { locale: 'menu.wechat.miniProgram', requiresAuth: true } },
  ],
};

// ========== 打印机 ==========
const PrinterManage: AppRouteRecordRaw = {
  path: '/printer/manage', name: 'PrinterManage',
  component: () => import('@/views/members/PrinterManageView.vue'),
  meta: { locale: 'menu.printer.manage', requiresAuth: true, order: 130 },
};

// ========== 短信 ==========
const SmsManage: AppRouteRecordRaw = {
  path: '/sms', name: 'SmsManage',
  component: DEFAULT_LAYOUT, redirect: '/sms/channel',
  meta: { locale: 'menu.sms.channel', requiresAuth: true, order: 140 },
  children: [
    { path: 'channel', name: 'SmsChannel', component: () => import('@/views/members/SmsChannelView.vue'), meta: { locale: 'menu.sms.channel', requiresAuth: true } },
    { path: 'template', name: 'SmsTemplate', component: () => import('@/views/members/SmsTemplateView.vue'), meta: { locale: 'menu.sms.template', requiresAuth: true } },
  ],
};

// ========== OTA管理 ==========
const OtaManage: AppRouteRecordRaw = {
  path: '/ota', name: 'OtaManage',
  component: DEFAULT_LAYOUT, redirect: '/ota/firmware',
  meta: { locale: 'menu.otaManage', requiresAuth: true, order: 150 },
  children: [
    { path: 'firmware', name: 'OTAView', component: () => import('@/views/OTAView.vue'), meta: { locale: 'menu.otaManage.firmware', requiresAuth: true } },
    { path: 'packages', name: 'OtaPackages', component: () => import('@/views/OtaPackages.vue'), meta: { locale: 'menu.otaManage.packages', requiresAuth: true } },
    { path: 'deployments', name: 'OtaDeployments', component: () => import('@/views/OtaDeployments.vue'), meta: { locale: 'menu.otaManage.deployments', requiresAuth: true } },
    { path: 'matrix', name: 'OtaMatrix', component: () => import('@/views/OtaMatrix.vue'), meta: { locale: 'menu.ota.matrix', requiresAuth: true } },
    { path: 'gray', name: 'OTAGrayDeploy', component: () => import('@/views/OTAGrayDeploy.vue'), meta: { locale: 'menu.ota.gray', requiresAuth: true } },
  ],
};

// ========== 租户管理 ==========
const TenantManage: AppRouteRecordRaw = {
  path: '/tenant', name: 'TenantManage',
  component: DEFAULT_LAYOUT, redirect: '/tenant/application',
  meta: { locale: 'menu.tenantManage', requiresAuth: true, order: 160 },
  children: [
    { path: 'application', name: 'TenantApplication', component: () => import('@/views/tenants/TenantApplication.vue'), meta: { locale: 'menu.tenant.application', requiresAuth: true } },
    { path: 'approval', name: 'TenantApproval', component: () => import('@/views/tenants/TenantApproval.vue'), meta: { locale: 'menu.tenant.approval', requiresAuth: true } },
    { path: 'management', name: 'TenantManagement', component: () => import('@/views/tenants/TenantManagement.vue'), meta: { locale: 'menu.tenant.management', requiresAuth: true } },
    { path: 'settings', name: 'TenantSettings', component: () => import('@/views/tenants/TenantSettings.vue'), meta: { locale: 'menu.tenant.settings', requiresAuth: true } },
  ],
};

// ========== 系统管理 ==========
const SystemManage: AppRouteRecordRaw = {
  path: '/system', name: 'SystemManage',
  component: DEFAULT_LAYOUT, redirect: '/system/logs',
  meta: { locale: 'menu.system', requiresAuth: true, order: 200 },
  children: [
    { path: 'email-templates', name: 'SystemEmailTemplates', component: () => import('@/views/system/EmailTemplates.vue'), meta: { locale: 'menu.system.emailTemplates', requiresAuth: true } },
    { path: 'logs', name: 'SystemLogs', component: () => import('@/views/system/Logs.vue'), meta: { locale: 'menu.system.logs', requiresAuth: true } },
    { path: 'monitor', name: 'SystemMonitor', component: () => import('@/views/system/Monitor.vue'), meta: { locale: 'menu.system.monitor', requiresAuth: true } },
  ],
};

// ========== 订阅管理 ==========
const SubscriptionManage: AppRouteRecordRaw = {
  path: '/subscription', name: 'SubscriptionManage',
  component: DEFAULT_LAYOUT, redirect: '/subscription/list',
  meta: { locale: 'menu.subscriptionManage', requiresAuth: true, order: 210 },
  children: [
    { path: 'list', name: 'SubscriptionList', component: () => import('@/views/subscription/SubscriptionList.vue'), meta: { locale: 'menu.subscription.list', requiresAuth: true } },
  ],
};

export default [
  MemberManage,
  DeviceManage,
  ContentMarket,
  Marketing,
  ActionMarket,
  AdCampaign,
  ContentReview,
  CouponPool,
  EmoticonMarket,
  RechargeRules,
  VoiceConfig,
  KnowledgeBase,
  OrgManagement,
  PermissionManage,
  StoreManage,
  CouponManage,
  VipExclusive,
  WechatManage,
  PrinterManage,
  SmsManage,
  OtaManage,
  TenantManage,
  SystemManage,
  SubscriptionManage,
];
