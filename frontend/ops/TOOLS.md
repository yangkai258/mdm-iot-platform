# TOOLS.md - agentyw 工具箱

## 监控工具

### 应用监控
- **Prometheus** - 指标收集
- **Grafana** - 可视化
- **APM** - 应用性能监控

### 日志
- **ELK** - 日志收集分析
- **Loki** - Grafana 生态
- **EFK** - 日志收集

### 告警
- **Alertmanager** - 告警管理
- **PagerDuty** - 告警通知
- **钉钉/飞书** - 集成告警

## 部署工具

### 容器
- **Docker** - 容器化
- **Docker Compose** - 本地编排
- **Kubernetes** - 生产编排

### CI/CD
- **GitHub Actions**
- **GitLab CI**
- **Jenkins**
- **ArgoCD** - GitOps

## 诊断工具

### 网络
```bash
# 端口检查
netstat -ano | findstr "8080"

# 网络延迟
ping host

# 路由追踪
tracert host
```

### 系统
```bash
# 资源使用
top
htop
df -h
free -m

# 进程
ps aux | grep mdm
```

### 日志
```bash
# 实时日志
tail -f app.log

# 错误日志
grep -i error app.log

# 日志统计
wc -l app.log
```

## 运维脚本

### 常用脚本
- 启动服务脚本
- 日志清理脚本
- 备份脚本
- 健康检查脚本

---

_自动化是运维的生命线。_
