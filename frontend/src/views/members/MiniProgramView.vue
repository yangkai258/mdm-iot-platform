<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>小程序基础数据</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="8"><a-card class="stat-card"><a-statistic title="小程序名称" :value="form.miniprogram_name || '-'" /></a-card></a-col>
      <a-col :span="8"><a-card class="stat-card"><a-statistic title="AppID" :value="form.app_id ? '已配置' : '未配置'" /></a-card></a-col>
      <a-col :span="8"><a-card class="stat-card"><a-statistic title="配置状态" :value="form.app_id && form.app_secret ? '已完成' : '未完成'" :value-style="{ color: form.app_id && form.app_secret ? '#52c41a' : '#ff6b00' }" /></a-card></a-col>
    </a-row>

    <div class="pro-content-area" style="margin-top: 16px;">
      <a-form :model="form" layout="vertical" style="max-width: 600px;">
        <a-divider>小程序基本信息</a-divider>

        <a-form-item label="小程序名称">
          <a-input v-model="form.miniprogram_name" placeholder="请输入小程序名称" />
        </a-form-item>

        <a-form-item label="AppID">
          <a-input v-model="form.app_id" placeholder="微信小程序AppID" />
        </a-form-item>

        <a-form-item label="AppSecret">
          <a-input v-model="form.app_secret" placeholder="微信小程序AppSecret" type="password" />
        </a-form-item>

        <a-form-item label="小程序Logo">
          <a-space direction="vertical">
            <a-input v-model="form.logo_url" placeholder="请输入logo URL" style="width: 300px;" />
            <img v-if="form.logo_url" :src="form.logo_url" style="max-width: 80px; max-height: 80px; border-radius: 8px; border: 1px solid #eee;" />
          </a-space>
        </a-form-item>

        <a-divider>版本配置</a-divider>

        <a-form-item label="当前版本">
          <a-input v-model="form.version" placeholder="如: 1.0.0" style="width: 200px;" />
        </a-form-item>

        <a-form-item label="最低支持版本">
          <a-input v-model="form.min_version" placeholder="如: 1.