# 明日工作计划 - 2026-04-02

## 背景
新前端 `mdm-frontend-new/arco-design-pro-vite/` 已就绪，但旧前端 `frontend/` 有 **270+ 个文件** 尚未迁移。登录、菜单、面包屑已修复，当前端基本可用。

## 迁移清单文档
- 完整清单：`mdm-project/docs/MIGRATION_MISSING.md`

---

## 明日任务（按优先级）

### Phase 1：核心页面 ✅
**目标：确保基本业务功能可用**

| 优先级 | 文件 | 说明 |
|--------|------|------|
| P0 | `Dashboard.vue` → `DashboardView.vue` | 主仪表盘迁移 |
| P0 | `DeviceDashboard.vue` | 设备仪表盘 |
| P0 | `DeviceStatus.vue` | 设备状态 |
| P0 | `Member.vue` | 会员总览 |
| P0 | `PetConfig.vue` | 宠物配置 |
| P0 | `OtaFirmware.vue` | OTA固件 |
| P0 | `Alert.vue` | 告警总览 |

### Phase 2：高频业务模块
**目标：覆盖日常使用功能**

| 模块 | 文件数 | 说明 |
|------|--------|------|
| 会员详细 | ~27 | 促销、标签、卡管理 |
| AI行为 | ~12 | 行为引擎、决策日志 |
| 情感计算 | ~7 | 情绪识别、趋势、配置 |
| 设备详细 | ~10 | 证书、地理围栏、远程控制 |
| 告警详细 | ~10 | 邮件/短信/Webhook通道 |

### Phase 3：企业级功能
**目标：完整企业功能**

| 模块 | 文件数 | 说明 |
|------|--------|------|
| 系统设置 | ~50 | LDAP、审计、证书、邮件模板 |
| 合规策略 | ~7 | 策略模板、分发、审计 |
| 全球化 | ~10 | 语言包、货币、时区 |
| 数字孪生 | ~5 | 生命体征、历史回放 |
| 健康医疗 | ~4 | 健康报告、预警 |

### Phase 4：高级功能
**目标：完整功能覆盖**

| 模块 | 文件数 | 说明 |
|------|--------|------|
| 仿真测试 | ~7 | 场景、压力测试 |
| 平台增强 | ~4 | 边缘AI、模型分片 |
| 应用分发 | ~3 | App版本管理 |
| 研究平台 | ~6 | 数据集、实验 |

---

## 技术注意事项

### 迁移标准
每个 Vue 文件迁移时必须：
1. **添加面包屑** - `<Breadcrumb :items="[...]" />`
2. **三段式布局** - 面包屑 → 搜索筛选区 → 数据表格
3. **API 对接** - 请求通过 `/api` 代理到后端 `8080`
4. **Arco Design** - 使用 `a-card.general-card` 包裹内容
5. **Authorization** - axios 请求带 `Bearer token`

### API 路径规范
- 前端请求：`/api/v1/xxx`
- Vite 代理：`/api` → `http://localhost:8080/api/v1`
- 不再直连 `localhost:8080`

### 响应格式
- 成功：`{ "code": 20000, "data": {...} }`
- 失败：`{ "code": 4xx, "message": "..." }`

---

## 启动命令

```bash
# 后端
cd C:/Users/YKing/.openclaw/workspace/mdm-project
./mdm-server.exe

# 前端（新）
cd C:/Users/YKing/.openclaw/workspace/mdm-project/mdm-frontend-new/arco-design-pro-vite
npm run dev
```

## 访问信息
- 前端：http://localhost:3000 （admin / admin123）
- 后端：http://localhost:8080
- EMQX：http://localhost:18083 （admin/public）
