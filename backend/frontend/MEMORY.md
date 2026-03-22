# MEMORY.md - agentqd 知识库

## 技术栈

### 核心技能
- Vue 3 (Composition API)
- Vite
- Arco Design
- Vue Router
- Pinia

### 了解工具
- Chrome DevTools
- Vue DevTools
- Figma (读懂设计稿)

## 项目经验

### MDM 控制中台
- 路径: frontend/dist/
- 框架: Vue 3 + Vite
- UI库: Arco Design
- 端口: 3000

## 已完成页面

| 页面 | 路径 | 说明 |
|------|------|------|
| 登录 | login.html | JWT登录 |
| 仪表盘 | dashboard.html | 统计概览 |
| 设备管理 | devices.html | 设备列表 |
| 会员管理 | member.html | 9个子模块 |
| ... | ... | ... |

## 组件库

### 常用组件
- `a-table` - 表格
- `a-form` - 表单
- `a-modal` - 弹窗
- `a-input` - 输入框
- `a-select` - 下拉选择
- `a-button` - 按钮
- `a-menu` - 菜单

## 踩坑记录

### 1. Arco 表格列宽
```js
// 不要写死宽度，用 minWidth
columns = [{ title: '名称', dataIndex: 'name', minWidth: 150 }]
```

### 2. Vue 响应式
```js
// 数组更新
array.value.push(newItem) // OK
array.value = [...array, newItem] // OK
array.value[index] = newItem // 不触发响应
```

### 3. API 跨域
```js
// vite.config.js
server: {
  proxy: {
    '/api': 'http://localhost:8080'
  }
}
```

## 学习计划

- [ ] 深入 Vue 3 响应式原理
- [ ] 掌握 TypeScript
- [ ] 学习性能优化技巧
- [ ] 提升 UI/UX 设计能力

---

_每一个像素都有意义，每一个交互都值得优化。_
