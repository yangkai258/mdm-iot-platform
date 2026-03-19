# TOOLS.md - agenthd 工具箱

## 开发工具

### IDE
- **GoLand** - 推荐，功能强大
- **VSCode** - 轻量级，需要插件
- **Sublime** - 极快

### 必备插件
- Go Extension
- Go Doc
- File Icons
- GitLens

### API 测试
- **Postman** - 功能完整
- **Insomnia** - 轻量美观
- **curl** - 命令行快速测试

## 数据库工具

### PostgreSQL
- **pgAdmin** - 官方客户端
- **DBeaver** - 通用数据库工具
- **DataGrip** - JetBrains 全家桶

### Redis
- **Redis Desktop Manager** - 图形化
- **RedisInsight** - 官方工具

## 调试工具

### 日志
```bash
# 实时查看日志
tail -f logs/app.log

# 搜索错误
grep -i error logs/app.log
```

### 性能分析
```bash
# CPU Profile
go tool pprof http://localhost:8080/debug/pprof/profile

# Memory Profile
go tool pprof http://localhost:8080/debug/pprof/heap
```

## 部署工具

- **Docker** - 容器化
- **Docker Compose** - 多容器编排
- **Nginx** - 反向代理

---

_选对工具，事半功倍。_
