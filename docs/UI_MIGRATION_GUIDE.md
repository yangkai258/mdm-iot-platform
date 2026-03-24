# UI 迁移指南 - Arco Design Pro

## ⚠️ 重要说明（必读）

**现状**：当前 `frontend/src/App.vue` 使用纯 CSS 自定义布局，没有使用 Arco Design 组件。
**目标**：使用 Arco Design Pro 官方组件和布局系统。

---

## 一、Arco Design Pro 核心架构

```
frontend/src/
├── App.vue                    # 入口：只用 a-config-provider + router-view
├── main.js                    # 入口JS
├── layout/                    # ⭐ 核心：布局组件
│   ├── index.vue              # 默认布局（侧边栏+Header+内容）
│   └── page-layout.vue        # 页面内容布局
├── components/                # ⭐ 核心组件
│   ├── navbar/                # 顶部导航栏
│   │   └── index.vue
│   ├── menu/                  # 侧边菜单
│   │   ├── index.vue
│   │   └── use-menu-tree.ts
│   ├── footer/                 # 页脚
│   │   └── index.vue
│   └── tab-bar/               # 标签页
│       └── index.vue
├── store/                     # ⭐ 状态管理（Pinia）
│   ├── index.ts
│   ├── modules/
│   │   ├── app.ts             # 应用状态（菜单折叠等）
│   │   └── user.ts            # 用户状态
├── router/                    # 路由（需要 meta.icon）
│   └── index.ts
├── hooks/                     # 工具钩子
│   ├── locale.ts
│   ├── permission.ts
│   └── responsive.ts
└── views/                     # 页面（不变）
```

---

## 二、正确写法 vs 错误写法

### ❌ 错误写法（当前）

```vue
<!-- App.vue - 不要这样做！ -->
<template>
  <div id="app">
    <div class="layout">
      <div class="sidebar">自定义侧边栏</div>
      <div class="main">
        <router-view />
      </div>
    </div>
  </div>
</template>

<style>
.sidebar { background: linear-gradient(...); }
.menu-item { padding: 12px; ... }
</style>
```

### ✅ 正确写法

```vue
<!-- App.vue - 只需这样 -->
<template>
  <a-config-provider :locale="locale">
    <router-view />
  </a-config-provider>
</template>

<script setup>
import { computed } from 'vue';
import zhCN from '@arco-design/web-vue/es/locale/lang/zh-cn';
import enUS from '@arco-design/web-vue/es/locale/lang/en-us';
import useLocale from '@/hooks/locale';

const { currentLocale } = useLocale();
const locale = computed(() => currentLocale.value === 'zh-CN' ? zhCN : enUS);
</script>
```

```vue
<!-- layout/index.vue - 布局组件 -->
<template>
  <a-layout class="layout">
    <a-layout-sider v-if="renderMenu" :collapsed="collapsed" :width="220">
      <Menu />
    </a-layout-sider>
    <a-layout>
      <a-layout-header>
        <NavBar />
      </a-layout-header>
      <a-layout-content>
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>
```

---

## 三、路由 meta 字段规范

```ts
// router/index.ts
{
  path: '/dashboard',
  name: 'Dashboard',
  component: () => import('@/views/dashboard/index.vue'),
  meta: {
    title: '仪表盘',
    icon: 'icon-dashboard',        // Arco 图标
    requiresAuth: true
  }
}
```

**图标格式**：
- Arco 内置图标：`icon-dashboard`, `icon-user`, `icon-settings`
- 或者使用组件：`<component :is="customIcon" />`

---

## 四、Menu 组件正确用法

```vue
<!-- components/menu/index.vue -->
<template>
  <a-menu 
    v-model:selected-keys="selectedKey"
    v-model:open-keys="openKeys"
    :mode="'vertical'"
    :collapsed="collapsed"
  >
    <a-menu-item key="dashboard">
      <template #icon><icon-dashboard /></template>
      <span>仪表盘</span>
    </a-menu-item>
    
    <a-sub-menu key="members">
      <template #icon><icon-user /></template>
      <template #title>会员管理</template>
      <a-menu-item key="member-list">会员列表</a-menu-item>
      <a-menu-item key="member-cards">会员卡</a-menu-item>
    </a-sub-menu>
  </a-menu>
</template>
```

---

## 五、样式变量（Dark Theme）

```css
/* 在 App.vue 或全局样式中 */
:root {
  --color-bg-1: #0f0f0f;
  --color-bg-2: #1a1a1a;
  --color-bg-3: #252525;
  --color-text-1: #e5e5e5;
  --color-text-2: #a3a3a3;
  --color-border: #2d2d2d;
}
```

---

## 六、迁移步骤清单

| 步骤 | 操作 | 状态 |
|------|------|------|
| 1 | 创建 `src/store/index.ts` + `modules/app.ts` + `modules/user.ts` | ⬜ |
| 2 | 创建 `src/hooks/locale.ts` | ⬜ |
| 3 | 创建 `src/layout/index.vue`（使用 a-layout） | ⬜ |
| 4 | 创建 `src/components/navbar/index.vue` | ⬜ |
| 5 | 创建 `src/components/menu/index.vue`（使用 a-menu） | ⬜ |
| 6 | 重写 `src/App.vue`（只用 a-config-provider） | ⬜ |
| 7 | 更新路由，添加 `meta.icon` | ⬜ |
| 8 | 测试 `npm run build` | ⬜ |

---

## 七、常见错误

### 1. 菜单不显示
- 检查路由是否添加了 `meta.title`
- 检查 `store/modules/app.ts` 是否正确设置 `menu: true`

### 2. 样式不对
- 确保引入了 Arco Design CSS：`import '@arco-design/web-vue/dist/arco.css'`

### 3. 图标不显示
- 使用 `<icon-dashboard />` 而不是 `icon-dashboard` 字符串
- 或者引入图标：`import { IconDashboard } from '@arco-design/web-vue/es/icon'`

---

## 八、参考资料

- Arco Design Vue: https://arco.design/vue/component/menu
- Arco Design Pro: https://arco.design/vue/pro/
- 图标库: https://arco.design/vue/component/icon

---

**最后更新**: 2026-03-24
