Agent CS - 测试工程师任务

你是高级自动化测试工程师 (agentcs)

## 当前任务：P0问题自动化测试

### 任务完成状态：✅ 完成

已在 `testing/p0_tests/` 目录创建以下测试文件：

1. **test_jwt_config.py** - JWT密钥配置验证
   - 静态检查 `middleware/jwt.go` 中是否有硬编码密钥
   - 验证源码中是否使用 `os.Getenv("JWT_SECRET")`
   - 功能验证：自定义密钥签发token应被服务器拒绝

2. **test_cors.py** - CORS配置安全验证
   - HTTP响应检查：验证 `Access-Control-Allow-Origin` 不是 `*`
   - 源码检查：`main.go` 中不应硬编码 wildcard CORS
   - 预检请求验证

3. **test_mqtt_command.py** - MQTT指令下发功能验证
   - 认证检查：指令端点需JWT认证
   - 离线设备返回错误码（不应200）
   - 参数校验：必填字段 `cmd_type`
   - 响应格式：返回 `cmd_id`
   - 指令历史API可访问

4. **test_login_enter_key.py** - 登录页Enter键功能验证
   - 源码分析：检查 Vue 组件事件绑定
   - `@submit`、`@keyup.enter` 等事件
   - 登录按钮 `type="submit"`
   - 表单验证和加载状态

### 测试文件结构
```
testing/p0_tests/
├── conftest.py              # pytest配置和共享fixtures
├── requirements.txt        # 依赖包
├── README.md               # 测试说明文档
├── test_jwt_config.py      # JWT配置测试
├── test_cors.py            # CORS安全测试
├── test_mqtt_command.py    # MQTT指令测试
└── test_login_enter_key.py # 登录Enter键测试
```

### 运行测试
```bash
cd testing
pip install -r p0_tests/requirements.txt
pytest p0_tests/ -v
```

### 发现的问题（待修复）

| Bug | 文件 | 问题描述 |
|-----|------|----------|
| JWT硬编码 | `backend/middleware/jwt.go` | `jwtSecret = []byte("mdm-secret-key-change-in-production")` |
| CORS开放 | `backend/main.go` | `Access-Control-Allow-Origin: *` 允许所有来源 |
| 离线返回200 | `backend/controllers/command_controller.go` | 设备离线时仍返回200 |
| Enter键失效 | `frontend/src/views/Login.vue` | 密码输入框未绑定Enter键事件 |

---

## 其他任务（进行中）

### mqtt_stress_test.py - MQTT压测脚本

待完成后补充到 `test_scripts/` 目录。
