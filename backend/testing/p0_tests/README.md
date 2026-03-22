# MDM P0 Tests

## 概述
P0级别自动化测试用例，针对已修复的P0问题进行验证。

## 测试用例

| 文件 | 测试目标 | Bug描述 |
|------|----------|---------|
| `test_jwt_config.py` | JWT密钥配置 | `middleware/jwt.go` 中JWT密钥硬编码为 `mdm-secret-key-change-in-production` |
| `test_cors.py` | CORS配置安全 | `main.go` 中CORS设置为 `Access-Control-Allow-Origin: *` 允许所有来源 |
| `test_mqtt_command.py` | MQTT指令下发 | 设备离线时仍返回200，指令下发缺少MQTT验证 |
| `test_login_enter_key.py` | 登录页回车功能 | Vue登录表单未正确绑定Enter键事件 |

## 安装依赖

```bash
pip install -r requirements.txt
```

## 运行测试

```bash
# 运行所有P0测试
pytest p0_tests/ -v

# 运行单个测试文件
pytest p0_tests/test_jwt_config.py -v

# 带详细输出
pytest p0_tests/ -v --tb=short
```

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `MDM_API_BASE_URL` | `http://localhost:8080` | 后端API地址 |
| `MDM_BACKEND_PATH` | 项目目录/backend | 后端源码路径 |
| `MDM_FRONTEND_PATH` | 项目目录/frontend | 前端源码路径 |
| `TEST_USERNAME` | `admin` | 测试用户名 |
| `TEST_PASSWORD` | `admin123` | 测试密码 |
| `MQTT_BROKER` | `tcp://localhost:1883` | MQTT Broker地址 |
| `MQTT_USERNAME` | `admin` | MQTT用户名 |
| `MQTT_PASSWORD` | `public` | MQTT密码 |

## 测试策略

### test_jwt_config.py
- **静态检查**: 读取 `middleware/jwt.go` 源码，正则匹配硬编码密钥模式
- **环境变量检查**: 确认源码中包含 `os.Getenv("JWT_SECRET")`
- **功能验证**: 用自定义密钥签发token，验证服务器拒绝（JWT_SECRET env var生效）

### test_cors.py
- **HTTP响应检查**: 实际请求API，验证响应头 `Access-Control-Allow-Origin` 不是 `*`
- **源码检查**: 验证 `main.go` 中没有硬编码 `Access-Control-Allow-Origin", "*"`
- **跨域预检**: 发送 OPTIONS 请求，验证预检响应正确

### test_mqtt_command.py
- **认证检查**: 验证指令端点需要JWT认证
- **离线设备检查**: 验证离线设备返回错误码（非200）
- **参数校验**: 验证必填字段 `cmd_type` 的校验
- **响应格式**: 验证返回 `cmd_id` 字段
- **指令历史**: 验证历史记录API可访问

### test_login_enter_key.py
- **源码分析**: 检查 `Login.vue` 中 `@submit`、`@keyup.enter` 等事件绑定
- **按钮类型**: 验证登录按钮 `type="submit"`
- **表单验证**: 验证 `handleLogin` 包含输入校验
- **加载状态**: 验证防重复提交机制

## 预期结果

修复完成后，所有测试应 **PASS**：
```
test_jwt_config.py     PASSED
test_cors.py           PASSED
test_mqtt_command.py   PASSED
test_login_enter_key.py PASSED
```

## 注意事项

- 测试文件仅做**静态代码检查**和**API功能验证**，不实际连接MQTT broker
- 部分测试需要后端服务运行在 `localhost:8080`
- 前端源码路径通过环境变量指定，用于静态代码分析
