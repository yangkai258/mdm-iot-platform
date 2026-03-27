# AGENTS.md - 前端开发团队协作规范

## 代码规范

### Vue 3 规范
- 使用 Composition API (`<script setup>`)
- 组件文件：`PascalCase.vue`
- 样式：使用 Scoped CSS
- 状态：使用 Pinia

### 命名规范
- 组件：`MemberList.vue`
- 变量：`camelCase`
- 常量：`UPPER_SNAKE_CASE`
- CSS 类：`kebab-case`

### Arco Design 使用
- 优先使用 Arco 组件
- 自定义样式要克制
- 主题变量要统一

## 目录结构
```
frontend/
├── src/
│   ├── components/    # 公共组件
│   ├── views/         # 页面
│   ├── stores/        # 状态管理
│   ├── utils/         # 工具函数
│   └── api/           # API 调用
└── dist/              # 静态输出
```

## 页面开发规范

### 新增页面
1. 在 `views/` 创建 `.vue` 文件
2. 路由添加到 `router/index.js`
3. 菜单添加到侧边栏配置
4. API 接口文档同步更新

### 组件规范
```vue
<template>
  <div class="component-name">
    <!-- 模板 -->
  </div>
</template>

<script setup>
// 逻辑
</script>

<style scoped>
/* 样式 */
</style>
```

## 性能优化

### 首屏加载
- 路由懒加载
- 图片懒加载
- Gzip 压缩

### 交互响应
- 骨架屏loading
- 防抖节流
- 虚拟滚动

---

_用户体验是产品的生命线。_
