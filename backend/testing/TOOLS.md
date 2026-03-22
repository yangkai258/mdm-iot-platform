# TOOLS.md - agentcs 工具箱

## 测试工具

### 功能测试
- **Postman** - API 测试
- **Selenium IDE** - 浏览器自动化
- **Cypress** - 现代前端测试

### 性能测试
- **JMeter** - 负载测试
- **k6** - 轻量级性能测试
- **Lighthouse** - 页面性能

### 移动测试
- **Charles** - 抓包工具
- **Fiddler** - Web调试代理

## 缺陷管理

### 工具
- **JIRA** - 企业级
- **禅道** - 国产开源
- **Mantis** - 轻量级

### 模板
```
缺陷标题: [模块-功能] 简短描述
严重程度: P0/P1/P2/P3
优先级: 高/中/低
环境: 
版本: 
复现步骤:
  1. 
  2. 
预期结果: 
实际结果: 
截图/日志: 
```

## 自动化

### API 自动化
```bash
# Newman 运行 Postman Collection
newman run collection.json -e environment.json -r html,json
```

### CI/CD 集成
- GitHub Actions
- GitLab CI
- Jenkins

## 辅助工具

### 抓包
- **Charles** - Mac 推荐
- **Fiddler** - Windows 推荐
- **Wireshark** - 协议分析

### 虚拟化
- **Docker** - 快速环境搭建
- **VirtualBox** - 虚拟机

---

_工具是手的延伸，用好工具事半功倍。_
