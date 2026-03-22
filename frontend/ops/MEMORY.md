# MEMORY.md - agentyw 知识库

## 环境信息

### MDM 控制中台
| 服务 | 端口 | 说明 |
|------|------|------|
| 后端 | 8080 | Go 服务 |
| 前端 | 3000 | Vue 服务 |
| 数据库 | 5432 | PostgreSQL |
| 缓存 | 6379 | Redis |

### 部署路径
- 后端: C:\Users\YKing\.openclaw\workspace\mdm-project\backend
- 前端: C:\Users\YKing\.openclaw\workspace\mdm-project\frontend
- 文档: C:\Users\YKing\.openclaw\workspace\mdm-project\docs

## 常用命令

### 后端服务
```bash
# 构建
go build -o mdm-server.exe

# 启动
./mdm-server.exe

# 查看日志
tail -f logs/app.log
```

### 前端服务
```bash
# 安装依赖
npm install

# 开发模式
npm run dev

# 构建
npm run build
```

### Docker
```bash
# 构建镜像
docker build -t mdm-backend:latest .

# 运行
docker run -d -p 8080:8080 mdm-backend:latest
```

## 巡检清单

### 每日巡检
- [ ] 服务运行状态
- [ ] 日志错误数量
- [ ] 磁盘空间
- [ ] 内存使用
- [ ] CPU 负载

### 每周巡检
- [ ] 数据库增长
- [ ] 日志清理
- [ ] 备份验证
- [ ] 证书过期检查

## 故障案例

### 端口被占用
```bash
# 查找占用进程
netstat -ano | findstr "8080"

# 杀掉进程
taskkill /PID <pid> /F
```

### 数据库连接失败
1. 检查 PostgreSQL 服务
2. 检查连接数限制
3. 检查防火墙

---

_不怕出问题，怕的是出问题不知道。_
