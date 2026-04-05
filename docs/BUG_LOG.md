# 🐛 错题本 - MDM 项目 Bug 与问题追踪

> 每次发现的问题、原因、解决方案都记录下来，避免重蹈覆辙。

---

## 如何使用

| 标记 | 含义 |
|------|------|
| 🔴 致命 | 阻断服务启动/登录不可用 |
| 🟠 严重 | 核心功能失效 |
| 🟡 一般 | 功能异常但有 workaround |
| 🟢 修复 | 已解决 |

---

## 🔴 [2026-04-05] 前端 API 路径缺少 /v1

### 问题
`/api/auth/me` 返回 401，但后端路由是 `/api/v1/auth/me`。

### 根因
前端代码调用 `/api/auth/me`，但后端 API 统一加 `v1` 版本前缀。
Vite proxy rewrite 规则 `path.replace(/^\/api/, '/api/v1')` 在前端跑 3001 端口时未生效（proxy 只监听配置的 3000）。

### 解决方案
1. 杀掉占用 3000 端口的旧前端进程
2. 重启 Arco Pro 前端自动绑定 3000
3. 验证 proxy 正常工作

### 预防
- 前端 API 统一写 `/api/xxx`，不写死 `/api/v1/xxx`
- 确保 vite proxy 配置的端口与实际运行端口一致
- 每次启动前端后验证 `/api/v1/auth/me` 返回 401（不是 404）

### 相关文件
- `mdm-frontend-new/arco-design-pro-vite/config/vite.config.dev.ts`
- `mdm-frontend-new/arco-design-pro-vite/src/api/user.ts`

---

## 🔴 [2026-04-04] Docker Desktop 无法启动 - Hyper-V 未安装

### 问题
```
WSL2 is not supported with your current machine configuration.
wslErrorCode: Wsl/Service/RegisterDistro/CreateVm/HCS/HCS_E_HYPERV_NOT_INSTALLED
```

### 根因
Docker Desktop 依赖 Hyper-V 运行，但机器上未启用该功能。

### 解决方案
用户需要手动执行（需要管理员权限）：
```powershell
Enable-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V-All
```
然后重启计算机，再启动 Docker Desktop。

### 预防
- 新机器部署前检查 Hyper-V 状态
- 在 README 或部署文档中注明前置条件

### 影响服务
PostgreSQL (5432), Redis (6379), EMQX (1883/8083/18083) 全部无法启动

---

## 🟠 [2026-03-29] 前端 git submodule 配置损坏

### 问题
`mdm-frontend-new` 目录 git submodule 配置损坏，无法正常拉取。

### 根因
 submodule 引用了一个不存在的远程仓库。

### 解决方案
删除损坏的 submodule 配置，重新初始化：
```bash
git submodule deinit --force mdm-frontend-new
git rm mdm-frontend-new
# 重新添加或手动克隆
```

### 预防
- 避免使用 git submodule，改用完整克隆
- 在 AGENTS.md 中明确标注

---

## 🟠 [2026-03-28] 后端构建目录混淆

### 问题
有两个入口点：
- `mdm-project/main.go`（不完整，缺少很多路由）
- `mdm-project/backend/main.go`（完整版本）

从错误目录构建导致 API 缺失。

### 根因
项目初期有两个独立的 main.go，共存于同一仓库。

### 解决方案
永远从 `backend/` 目录构建：
```bash
cd backend && go build -o ../mdm-server.exe .
```

### 预防
- 在 AGENTS.md 中明确标注构建路径
- 禁止从根目录构建

### 相关文件
- `mdm-project/main.go`（废弃）
- `mdm-project/backend/main.go`（正确）

---

## 🟠 [2026-03-23] JWT RefreshToken 过期检查逻辑错误

### 问题
`claims.ExpiresAt.Unix()` 用于判断 RefreshToken 是否过期，但 RefreshToken 的过期时间是绝对时间戳，无需额外计算。

### 根因
混淆了 RefreshToken 过期时间（绝对时间戳）和相对过期时长。

### 解决方案
修复 `GenerateToken` 函数，正确传递 `IsSuperAdmin` 参数：
```go
token, err := GenerateToken(userID, username, isSuperAdmin, expiresAt)
```

### 相关文件
- `backend/controllers/auth_controller.go`

---

## 🟠 [2026-03-22] Git 浅克隆导致文件丢失

### 问题
GitHub 上有 `docs/04_会员营销系统.md` 等文件，但本地没有。

### 根因
Clone 或 Fetch 时用了 `--depth=1`（浅克隆），只拉最新 commit，遗漏历史中的文件。

### 解决方案
```bash
git fetch origin --unshallow
```

### 预防
- 永远不用 `--depth=1`
- 已记录到 TOOLS.md

---

## 🟠 [2026-03-22] 多 Agent 并行开发导致代码冲突

### 问题
5个 Agent 同时修改相同文件，push 被拒后代码丢失。

### 根因
- 无文件范围分配机制
- 无推送协调机制
- push 失败后未及时 rebase

### 解决方案
已建立 AGENTS.md 中的规范：
1. 任务分配明确文件范围
2. 串行推送
3. 编译门禁
4. 超时监控

### 预防
- 大任务拆分（≤30分钟）
- 文件范围禁止重叠
- push 前先 fetch --rebase

---

## 🟡 [2026-04-01] 菜单 Breadcrumb 为空

### 问题
redirect 路由没有 `component` 但被菜单渲染，导致 breadcrumb 为空。

### 根因
部分路由设置了 `redirect` 但 `hideInMenu` 标记缺失。

### 解决方案
给 redirect 路由加 `hideInMenu: true`：
```ts
{
  path: '/redirect',
  component: () => import('@/views/redirect.vue'),
  hideInMenu: true,
}
```

### 相关文件
- `mdm-frontend-new/arco-design-pro-vite/src/router/routes/modules/*.ts`

---

## 🟡 [2026-03-29] GBK 编码文件导致乱码

### 问题
部分 Vue 文件保存为 GBK 编码，中文显示乱码。

### 根因
编辑器或构建工具未正确处理文件编码。

### 解决方案
用 UTF-8 重新保存以下文件：
- `FirmwareList.vue`
- `PetConsole.vue`

### 预防
- 配置编辑器默认 UTF-8
- 添加 .gitattributes 指定编码

---

## 🔴 [2026-04-05] 150+ Vue 文件模板语法错误（重复 `</a-table>` + `` `n `` 乱码）

### 问题
访问多个页面时 `Failed to load resource: 500`，Vue 动态导入失败。

### 根因
AI 代码生成时错误地在 `<a-table ... />` 后面加了一个 `</a-table>` 闭合标签，并且模板末尾出现了 `` `n ``（模板字符串字面量被嵌入 HTML）。

典型错误结构：
```html
<a-table :columns="columns" :data="data" />
</a-table>  ← 多余的闭合标签
</a-card>`n</div></template>  ← `n` 是垃圾字符
```

### 解决方案
1. 清除所有重复的 `</a-table>`：`$content -replace '(</a-table>)\s*(</a-table>)', '$1'`
2. 清除 `` `n</div></template> `` 乱码：`$content -replace '`n</div></template>', "\n</div>\n</template>"`
3. 共修复 **150+ 个 Vue 文件**

### 预防
- AI 生成 Vue 代码后，用 `Select-String -Pattern '</a-table>'` 检查
- 建议用 Vite build 预检查：`npm run build` 能提前发现动态 import 500 错误

---

## 🟢 [2026-04-05] 前端 API 路径缺少 /v1

### 问题
`/api/auth/me` 返回 401，但后端路由是 `/api/v1/auth/me`。

### 根因
Vite proxy 配置监听 3000，但旧前端占用了 3000 端口，导致 Arco Pro 前端跑到 3001，proxy 未生效。

### 解决方案
1. 杀掉占用 3000 端口的旧前端进程
2. 重启 Arco Pro 前端自动绑定 3000

---

## 🟢 [2026-03-23] Bcrypt 密码验证被注释

### 问题
`auth_controller.go` 中 bcrypt 验证被注释，导致所有密码都能登录。

### 根因
调试时临时注释，忘记恢复。

### 解决方案
恢复 bcrypt 验证逻辑。

---

## 🟢 [2026-03-23] SQL 迁移 varchar:20 语法错误

### 问题
Migration 中写了 `varchar:20` 而非 `varchar(20)`。

### 根因
GORM 语法错误。

### 解决方案
修复迁移文件中的类型声明。

---

_每解决一个 bug，就在前面加 🟢 标记，记录日期和解决方案。_
