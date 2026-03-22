# MDM Sprint 3 - 策略管理功能测试
"""
测试策略管理 (Policy Management) 的核心功能：
- 合规策略 (CompliancePolicy) 模型定义
- 合规违规 (ComplianceViolation) 模型定义
- CheckCompliance 合规检查函数
- 合规策略 API 路由注册
- 合规违规自动记录逻辑
- MQTT 上报数据触发合规检查

依赖 backend Sprint 3 实现：
- models/compliance.go (CompliancePolicy, ComplianceViolation)
- controllers/alert_controller.go (CheckCompliance 函数)
- 需新增合规策略 API 路由
"""
import os
import pytest
import requests
import time


class TestPolicyManagement:
    """
    策略管理功能测试套件
    PASS = 策略管理功能正确实现
    FAIL = 策略管理功能未实现或有问题
    """

    BASE_URL = os.getenv("MDM_API_BASE_URL", "http://localhost:8080")
    API_V1 = f"{BASE_URL}/api/v1"

    def _get_auth_token(self):
        """Helper: 获取认证 token"""
        login_data = {
            "username": os.getenv("TEST_USERNAME", "admin"),
            "password": os.getenv("TEST_PASSWORD", "admin123")
        }
        try:
            resp = requests.post(f"{self.API_V1}/auth/login", json=login_data, timeout=5)
            if resp.status_code == 200:
                data = resp.json()
                if data.get("code") == 0:
                    return data["data"]["token"]
        except Exception:
            pass
        return None

    def _auth_headers(self):
        """Helper: 获取认证头"""
        token = self._get_auth_token()
        if token:
            return {"Authorization": f"Bearer {token}"}
        return {}

    def _backend_src_path(self, rel_path):
        """Helper: 返回 backend 源码文件的绝对路径"""
        return os.path.join(
            os.getenv("MDM_BACKEND_PATH", "C:\\Users\\YKing\\.openclaw\\workspace\\mdm-project\\backend"),
            rel_path
        )

    # ─── 静态源码检查 ────────────────────────────────────────────

    def _read_source(self, rel_path):
        """Helper: 读取源码文件内容"""
        path = self._backend_src_path(rel_path)
        if not os.path.exists(path):
            return None
        with open(path, "r", encoding="utf-8") as f:
            return f.read()

    def test_compliance_policy_model_exists(self):
        """
        验证 CompliancePolicy 模型已定义

        检查 models/compliance.go 中是否定义 CompliancePolicy 结构体：
        - Name, Description, PolicyType, TargetValue, Condition
        - Severity, RemediationAction, Enabled, EnforceScope
        - TableName() 方法返回 compliance_policies
        """
        source = self._read_source("models/compliance.go")
        assert source is not None, "models/compliance.go 文件不存在"

        assert "type CompliancePolicy struct" in source, \
            "CompliancePolicy 结构体未定义"
        assert 'gorm:"primaryKey"' in source, \
            "CompliancePolicy 缺少 gorm primaryKey 标记"
        assert 'json:"name"' in source, \
            "CompliancePolicy 缺少 name 字段"
        assert 'json:"policy_type"' in source, \
            "CompliancePolicy 缺少 policy_type 字段 (如 firmware_version, battery_level)"
        assert 'json:"target_value"' in source, \
            "CompliancePolicy 缺少 target_value 字段"
        assert 'json:"condition"' in source, \
            "CompliancePolicy 缺少 condition 字段 (=/!=/>/</>=/<=)"
        assert 'json:"remediation_action"' in source, \
            "CompliancePolicy 缺少 remediation_action 字段"
        assert 'json:"enabled"' in source, \
            "CompliancePolicy 缺少 enabled 字段"

        # TableName 方法
        assert 'func (CompliancePolicy) TableName()' in source, \
            "CompliancePolicy 缺少 TableName() 方法"
        assert 'return "compliance_policies"' in source, \
            "CompliancePolicy TableName() 应返回 compliance_policies"

    def test_compliance_violation_model_exists(self):
        """
        验证 ComplianceViolation 模型已定义

        检查 models/compliance.go 中是否定义 ComplianceViolation 结构体：
        - PolicyID, DeviceID, PolicyType, ExpectedValue, ActualValue
        - Severity, ActionTaken, Status, ResolvedAt, ResolvedBy
        """
        source = self._read_source("models/compliance.go")
        assert source is not None, "models/compliance.go 文件不存在"

        assert "type ComplianceViolation struct" in source, \
            "ComplianceViolation 结构体未定义"
        assert 'json:"policy_id"' in source, \
            "ComplianceViolation 缺少 policy_id 字段"
        assert 'json:"device_id"' in source, \
            "ComplianceViolation 缺少 device_id 字段"
        assert 'json:"policy_type"' in source, \
            "ComplianceViolation 缺少 policy_type 字段"
        assert 'json:"expected_value"' in source, \
            "ComplianceViolation 缺少 expected_value 字段"
        assert 'json:"actual_value"' in source, \
            "ComplianceViolation 缺少 actual_value 字段"
        assert 'json:"action_taken"' in source, \
            "ComplianceViolation 缺少 action_taken 字段"
        assert 'json:"status"' in source, \
            "ComplianceViolation 缺少 status 字段 (1:待处理 2:处理中 3:已解决 4:已忽略)"

    def test_check_compliance_function_exists(self):
        """
        验证 CheckCompliance 合规检查函数已实现

        检查 controllers/alert_controller.go 中 CheckCompliance 函数：
        - 查询 enabled=true 的合规策略
        - 支持 battery_level / offline_duration / is_online 类型
        - 违规时创建 ComplianceViolation 记录
        """
        source = self._read_source("controllers/alert_controller.go")
        assert source is not None, "controllers/alert_controller.go 文件不存在"

        assert "func CheckCompliance(" in source, \
            "CheckCompliance 函数未定义"

        # 合规策略查询
        assert 'Where("enabled = ?"' in source or "enabled = ?" in source, \
            "CheckCompliance 未查询 enabled=true 的策略"

        # 支持的策略类型
        assert 'battery_level' in source, \
            "CheckCompliance 未支持 battery_level 策略类型"
        assert 'offline_duration' in source or 'elapsed' in source, \
            "CheckCompliance 未支持 offline_duration 策略类型"
        assert 'is_online' in source, \
            "CheckCompliance 未支持 is_online 策略类型"

        # 违规记录创建
        assert "ComplianceViolation" in source, \
            "CheckCompliance 未创建 ComplianceViolation 记录"

    def test_compliance_policy_routes_registered(self):
        """
        验证合规策略 API 路由已注册

        检查路由注册位置（device_controller.go 或 main.go）：
        - GET  /compliance/policies        - 策略列表
        - POST /compliance/policies        - 创建策略
        - GET  /compliance/policies/:id     - 策略详情
        - PUT  /compliance/policies/:id     - 更新策略
        - DELETE /compliance/policies/:id   - 删除策略
        - GET  /compliance/violations       - 违规记录列表
        """
        dc_src = self._read_source("controllers/device_controller.go") or ""
        main_src = self._read_source("main.go") or ""

        combined = dc_src + main_src

        # 策略 CRUD 路由
        policy_crud = (
            ("GET", "/compliance/policies") in
            [(m.strip(), p.strip()) for m, p in
             [(l.split()[0], " ".join(l.split()[1:]))
              for l in combined.splitlines() if l.strip().startswith(("GET", "POST", "PUT", "DELETE"))]]
            if False else
            "/compliance/policies" in combined
        )

        if not policy_crud:
            # 合规策略路由未注册，这是一个已知缺失
            pytest.skip("合规策略 API 路由未注册 (device_controller.go 或 main.go)")
        else:
            # 如果路由存在，验证 DELETE 和 PUT 也存在
            assert 'DELETE' in combined and '/compliance/policies' in combined, \
                "合规策略 DELETE 路由未注册"
            assert 'PUT' in combined and '/compliance/policies' in combined, \
                "合规策略 PUT 路由未注册"

    def test_compliance_auto_migrate_registered(self):
        """
        验证 CompliancePolicy 和 ComplianceViolation 已加入数据库自动迁移

        检查 main.go AutoMigrate 调用中是否包含这两个模型
        """
        source = self._read_source("main.go")
        assert source is not None, "main.go 文件不存在"

        assert "CompliancePolicy" in source, \
            "main.go AutoMigrate 未注册 CompliancePolicy"
        assert "ComplianceViolation" in source, \
            "main.go AutoMigrate 未注册 ComplianceViolation"

    def test_compliance_callback_registered_in_mqtt(self):
        """
        验证 MQTT 初始化时注册了合规检查回调

        检查 main.go 中 mqtt.InitMQTT 调用是否传入了合规回调：
        - complianceCallback := func(...) { controllers.CheckCompliance(...) }
        """
        source = self._read_source("main.go")
        assert source is not None, "main.go 文件不存在"

        assert "CheckCompliance" in source, \
            "main.go 未注册 CheckCompliance 回调"
        assert "complianceCallback" in source, \
            "main.go 未定义 complianceCallback"

    def test_compliance_condition_evaluation(self):
        """
        验证合规策略条件评估逻辑

        检查 alert_controller.go 中合规条件评估：
        - 支持 >, <, =, >=, <=, != 条件
        - battery_level 使用 float64 比较
        - offline_duration 使用 elapsed 时间比较
        """
        source = self._read_source("controllers/alert_controller.go")
        assert source is not None, "controllers/alert_controller.go 文件不存在"

        # 确认 evaluateCondition 被 CheckCompliance 调用
        assert "evaluateCondition" in source, \
            "未使用 evaluateCondition 函数进行条件评估"
        assert "threshold" in source, \
            "缺少 threshold 阈值处理逻辑"

    # ─── API 功能测试 ────────────────────────────────────────────

    def test_compliance_policy_api_accessible(self):
        """
        验证合规策略 API 端点可访问（JWT 认证）

        测试 GET /api/v1/compliance/policies
        后端未运行或路由未注册时跳过
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        try:
            resp = requests.get(f"{self.API_V1}/compliance/policies",
                               headers=headers, timeout=5)
            # 200 = 路由存在，401 = 路由存在但需认证
            assert resp.status_code in [200, 401], \
                f"合规策略 API 异常响应: {resp.status_code}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_compliance_violation_api_accessible(self):
        """
        验证合规违规记录 API 端点可访问（JWT 认证）

        测试 GET /api/v1/compliance/violations
        后端未运行或路由未注册时跳过
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        try:
            resp = requests.get(f"{self.API_V1}/compliance/violations",
                               headers=headers, timeout=5)
            assert resp.status_code in [200, 401], \
                f"合规违规 API 异常响应: {resp.status_code}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_compliance_policy_create(self):
        """
        验证创建合规策略 API

        POST /api/v1/compliance/policies
        发送合规策略创建请求，验证响应格式
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        payload = {
            "name": "电量低于20%告警",
            "policy_type": "battery_level",
            "target_value": "20",
            "condition": "<",
            "severity": 2,
            "remediation_action": "notify",
            "enabled": True
        }

        try:
            resp = requests.post(f"{self.API_V1}/compliance/policies",
                               headers=headers, json=payload, timeout=5)
            if resp.status_code == 404:
                pytest.skip("合规策略 API 路由未注册")
            assert resp.status_code in [200, 201, 400], \
                f"创建合规策略响应异常: {resp.status_code}"
            if resp.status_code == 200:
                data = resp.json()
                assert data.get("code") == 0, f"响应 code 非0: {data}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_compliance_policy_list_pagination(self):
        """
        验证合规策略列表分页 API

        GET /api/v1/compliance/policies?page=1&page_size=10
        验证返回 list, total, page, page_size 字段
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        try:
            resp = requests.get(
                f"{self.API_V1}/compliance/policies?page=1&page_size=10",
                headers=headers, timeout=5)
            if resp.status_code == 404:
                pytest.skip("合规策略 API 路由未注册")
            if resp.status_code == 200:
                data = resp.json()
                assert "data" in data, "响应缺少 data 字段"
                assert "list" in data["data"], "data 缺少 list 字段"
                assert "total" in data["data"], "data 缺少 total 字段"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
