# AGENTS.md - 后端开发团队协作规范

## 代码规范

### Go 语言规范
- 使用 `go fmt` 格式化代码
- 变量命名遵循 Go 惯例
- 错误处理：`if err != nil` 必须检查
- 注释：导出函数必须注释

### API 设计规范
- 路径：`/api/v1/module/action`
- 方法：RESTful
- 认证：JWT Bearer Token
- 限流：重要接口加限流中间件

### Git 规范
- 分支命名：`feature/xxx` / `fix/xxx`
- Commit Message：`feat: | fix: | docs: | refactor:`
- PR 必须有代码审查

## 数据库规范

### 表设计
- 主键：`id` 使用 uint 自动递增
- 软删除：`deleted_at` timestamp
- 时间戳：`created_at`, `updated_at`
- 索引：高频查询字段必须加索引

### 命名
- 表名单数：`member` 不是 `members`
- 字段下划线：`member_name` 不是 `memberName`
- 布尔字段：`is_xxx` / `has_xxx`

## 中间件顺序
```
1. Logger
2. Recovery
3. CORS
4. JWT Auth
5. Rate Limit
6. Router
```

## 依赖管理
- 使用 Go Modules
- 定期 `go mod tidy`
- 锁定版本号

---

_代码是给人读的，顺便能在机器上运行。_
