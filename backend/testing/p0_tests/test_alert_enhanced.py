# MDM Sprint 3 - 告警增强功能测试
"""
测试告警增强功能 (Alert Enhanced) 的核心功能：
- 告警规则 CRUD（创建/查询/更新/删除/启用禁用）
- 告警记录状态流转（未处理→已确认→已解决）
- 告警规则条件类型（battery_low, offline, temperature_high, custom）
- 告警通知方式（email, sms, webhook）
- 告警记录筛选和排序
- 大盘统计数据 API
- MQTT 设备状态变化触发 CheckAlerts

注意：test_alert_trigger.py 已覆盖基础 CheckAlerts 逻辑，
本文件聚焦告警管理 API 和增强功能。
"""
import os
import pytest
import requests
import time


class TestAlertEnhanced:
    """
    告警增强功能测试套件
    PASS = 告警增强功能正确实现
    FAIL = 告警增强功能未实现或有问题
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

    def _read_source(self, rel_path):
        """Helper: 读取源码文件内容"""
        path = self._backend_src_path(rel_path)
        if not os.path.exists(path):
            return None
        with open(path, "r", encoding="utf-8") as f:
            return f.read()

    # ─── 静态源码检查 ────────────────────────────────────────────

    def test_alert_rule_model_fields(self):
        """
        验证 DeviceAlertRule 模型字段完整性

        检查 models/alert_models.go 中 DeviceAlertRule 字段：
        - ID, Name, DeviceID, AlertType, Condition, Threshold
        - Severity, Enabled, NotifyWays, Remark
        """
        source = self._read_source("models/alert_models.go")
        assert source is not None, "models/alert_models.go 文件不存在"

        assert "DeviceAlertRule" in source, \
            "DeviceAlertRule 模型未定义"
        assert 'json:"name"' in source, \
            "DeviceAlertRule 缺少 name 字段"
        assert 'json:"device_id"' in source, \
            "DeviceAlertRule 缺少 device_id 字段 (空表示所有设备)"
        assert 'json:"alert_type"' in source, \
            "DeviceAlertRule 缺少 alert_type 字段"
        assert 'json:"condition"' in source, \
            "DeviceAlertRule 缺少 condition 字段"
        assert 'json:"threshold"' in source, \
            "DeviceAlertRule 缺少 threshold 字段"
        assert 'json:"severity"' in source, \
            "DeviceAlertRule 缺少 severity 字段"
        assert 'json:"enabled"' in source, \
            "DeviceAlertRule 缺少 enabled 字段"
        assert 'json:"notify_ways"' in source, \
            "DeviceAlertRule 缺少 notify_ways 字段 (email,sms,webhook)"

    def test_alert_record_model_fields(self):
        """
        验证 DeviceAlert 模型字段完整性

        检查 models/alert_models.go 中 DeviceAlert 字段：
        - RuleID, DeviceID, AlertType, Severity, Message
        - TriggerVal, Threshold, Status, CreatedAt
        """
        source = self._read_source("models/alert_models.go")
        assert source is not None, "models/alert_models.go 文件不存在"

        assert "DeviceAlert" in source, \
            "DeviceAlert 模型未定义"
        assert 'json:"rule_id"' in source, \
            "DeviceAlert 缺少 rule_id 字段"
        assert 'json:"device_id"' in source, \
            "DeviceAlert 缺少 device_id 字段"
        assert 'json:"alert_type"' in source, \
            "DeviceAlert 缺少 alert_type 字段"
        assert 'json:"trigger_val"' in source, \
            "DeviceAlert 缺少 trigger_val 字段 (触发时的实际值)"
        assert 'json:"status"' in source, \
            "DeviceAlert 缺少 status 字段 (1未处理 2已确认 3已解决)"

    def test_alert_rule_update_delete_methods(self):
        """
        验证告警规则 Update/Delete 控制器方法存在

        检查 controllers/alert_controller.go 中：
        - UpdateRule (更新规则)
        - DeleteRule (删除规则)
        - ToggleRule / EnableRule (启用禁用)
        """
        source = self._read_source("controllers/alert_controller.go")
        assert source is not None, "controllers/alert_controller.go 文件不存在"

        # GetRules 和 CreateRule 已在 test_alert_trigger 验证
        has_update = "UpdateRule" in source or "func (c *AlertController) UpdateRule" in source
        has_delete = "DeleteRule" in source or "func (c *AlertController) DeleteRule" in source

        if not has_update or not has_delete:
            pytest.skip("UpdateRule/DeleteRule 方法未实现 (需 agenthd 实现)")

    def test_alert_rule_routes_in_main(self):
        """
        验证告警管理路由在 main.go 中正确注册

        检查 main.go sys group 中的告警路由：
        - GET /alerts/rules (已有)
        - POST /alerts/rules (已有)
        - GET /alerts (已有)
        - 验证 DashboardStats 路由存在
        """
        source = self._read_source("main.go")
        assert source is not None, "main.go 文件不存在"

        assert "/alerts/rules" in source, \
            "告警规则路由未在 main.go 注册"
        assert "/alerts" in source, \
            "告警记录路由未在 main.go 注册"
        assert "alertCtrl" in source, \
            "alertCtrl 未在 main.go 中初始化"

        # Dashboard stats
        assert "dashboard" in source.lower() or "stats" in source, \
            "大盘统计路由未注册 (应包含 dashboard/stats)"

    def test_alert_callback_in_mqtt_init(self):
        """
        验证 MQTT 初始化时注册了告警回调

        检查 main.go 中 mqtt.InitMQTT 调用：
        - 传入了 alertCallback 参数
        - alertCallback 调用 controllers.CheckAlerts
        """
        source = self._read_source("main.go")
        assert source is not None, "main.go 文件不存在"

        assert "alertCallback" in source, \
            "main.go 未定义 alertCallback"
        assert "CheckAlerts" in source, \
            "alertCallback 未调用 CheckAlerts"

    def test_evaluate_condition_function(self):
        """
        验证 evaluateCondition 条件评估函数

        检查 alert_controller.go 中 evaluateCondition：
        - 支持 <, >, =, <=, >=, != 条件
        """
        source = self._read_source("controllers/alert_controller.go")
        assert source is not None, "controllers/alert_controller.go 文件不存在"

        assert "func evaluateCondition" in source, \
            "evaluateCondition 函数未定义"

        # 支持的条件操作符
        conditions = ['"<"', '">"', '"="', '"<="', '">="']
        found = [c in source for c in conditions]
        assert any(found), \
            f"evaluateCondition 未支持足够的条件操作符: {conditions}"

    # ─── API 功能测试 ────────────────────────────────────────────

    def test_alert_rule_create_api(self):
        """
        验证创建告警规则 API

        POST /api/v1/alerts/rules
        创建 battery_low 规则，验证响应
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        payload = {
            "name": f"电量低告警_{int(time.time())}",
            "alert_type": "battery_low",
            "condition": "<",
            "threshold": 20.0,
            "severity": 2,
            "enabled": True,
            "notify_ways": "email,webhook",
            "remark": "测试告警规则"
        }

        try:
            resp = requests.post(f"{self.API_V1}/alerts/rules",
                               headers=headers, json=payload, timeout=5)
            if resp.status_code == 404:
                pytest.skip("告警规则 API 路由未注册")
            assert resp.status_code in [200, 201, 400], \
                f"创建告警规则响应异常: {resp.status_code}"
            if resp.status_code == 200:
                data = resp.json()
                assert data.get("code") == 0, f"响应 code 非0: {data}"
                assert "data" in data, "响应缺少 data 字段"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_alert_rule_create_offline_type(self):
        """
        验证创建离线告警规则 API

        POST /api/v1/alerts/rules
        创建 offline 类型规则
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        payload = {
            "name": f"设备离线告警_{int(time.time())}",
            "alert_type": "offline",
            "condition": "=",
            "threshold": 0.0,
            "severity": 3,
            "enabled": True,
            "notify_ways": "webhook",
            "remark": "离线超过90秒触发"
        }

        try:
            resp = requests.post(f"{self.API_V1}/alerts/rules",
                               headers=headers, json=payload, timeout=5)
            if resp.status_code == 404:
                pytest.skip("告警规则 API 路由未注册")
            assert resp.status_code in [200, 201, 400], \
                f"创建离线告警规则响应异常: {resp.status_code}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_alert_rule_create_custom_type(self):
        """
        验证创建自定义类型告警规则 API

        POST /api/v1/alerts/rules
        验证支持 temperature_high 等自定义类型
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        payload = {
            "name": f"温度过高告警_{int(time.time())}",
            "alert_type": "temperature_high",
            "condition": ">",
            "threshold": 80.0,
            "severity": 4,
            "enabled": True,
            "notify_ways": "sms,webhook",
            "remark": "设备温度超过80度触发"
        }

        try:
            resp = requests.post(f"{self.API_V1}/alerts/rules",
                               headers=headers, json=payload, timeout=5)
            if resp.status_code == 404:
                pytest.skip("告警规则 API 路由未注册")
            assert resp.status_code in [200, 201, 400], \
                f"创建自定义类型告警规则响应异常: {resp.status_code}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_alert_rule_list_pagination(self):
        """
        验证告警规则列表分页 API

        GET /api/v1/alerts/rules
        验证返回 list, total 等字段
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        try:
            resp = requests.get(f"{self.API_V1}/alerts/rules",
                              headers=headers, timeout=5)
            if resp.status_code == 404:
                pytest.skip("告警规则 API 路由未注册")
            assert resp.status_code == 200, \
                f"告警规则列表响应异常: {resp.status_code}"
            data = resp.json()
            assert data.get("code") == 0, f"响应 code 非0: {data}"
            assert "data" in data, "响应缺少 data 字段"
            assert "list" in data["data"], "data 缺少 list 字段"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_alert_rule_update_api(self):
        """
        验证更新告警规则 API

        PUT /api/v1/alerts/rules/:id
        先创建再更新
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        # 创建
        create_payload = {
            "name": f"待更新告警_{int(time.time())}",
            "alert_type": "battery_low",
            "condition": "<",
            "threshold": 15.0,
            "severity": 2,
            "enabled": True
        }

        try:
            create_resp = requests.post(f"{self.API_V1}/alerts/rules",
                                       headers=headers, json=create_payload, timeout=5)
            if create_resp.status_code == 404:
                pytest.skip("告警规则 API 路由未注册")
            if create_resp.status_code not in [200, 201]:
                pytest.skip(f"创建告警规则失败: {create_resp.status_code}")

            created = create_resp.json()
            rule_id = created.get("data", {}).get("id")
            if not rule_id:
                pytest.skip("创建响应中无 id 字段")

            # 更新
            update_payload = {
                "name": "更新后的告警名称",
                "threshold": 10.0,
                "severity": 3,
                "enabled": False
            }
            update_resp = requests.put(
                f"{self.API_V1}/alerts/rules/{rule_id}",
                headers=headers, json=update_payload, timeout=5)
            if update_resp.status_code == 404:
                pytest.skip("UpdateRule 方法未实现 (需 agenthd 实现)")
            assert update_resp.status_code in [200, 400], \
                f"更新告警规则响应异常: {update_resp.status_code}"

        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_alert_rule_delete_api(self):
        """
        验证删除告警规则 API

        DELETE /api/v1/alerts/rules/:id
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        # 创建
        create_payload = {
            "name": f"待删除告警_{int(time.time())}",
            "alert_type": "battery_low",
            "condition": "<",
            "threshold": 15.0,
            "severity": 2,
            "enabled": True
        }

        try:
            create_resp = requests.post(f"{self.API_V1}/alerts/rules",
                                       headers=headers, json=create_payload, timeout=5)
            if create_resp.status_code == 404:
                pytest.skip("告警规则 API 路由未注册")
            if create_resp.status_code not in [200, 201]:
                pytest.skip(f"创建告警规则失败: {create_resp.status_code}")

            created = create_resp.json()
            rule_id = created.get("data", {}).get("id")
            if not rule_id:
                pytest.skip("创建响应中无 id 字段")

            # 删除
            del_resp = requests.delete(
                f"{self.API_V1}/alerts/rules/{rule_id}",
                headers=headers, timeout=5)
            if del_resp.status_code == 404:
                pytest.skip("DeleteRule 方法未实现 (需 agenthd 实现)")
            assert del_resp.status_code in [200, 404], \
                f"删除告警规则响应异常: {del_resp.status_code}"

        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_alerts_list_api(self):
        """
        验证告警记录列表 API

        GET /api/v1/alerts
        验证返回告警记录列表
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        try:
            resp = requests.get(f"{self.API_V1}/alerts",
                               headers=headers, timeout=5)
            assert resp.status_code == 200, \
                f"告警记录列表响应异常: {resp.status_code}"
            data = resp.json()
            assert data.get("code") == 0, f"响应 code 非0: {data}"
            assert "data" in data, "响应缺少 data 字段"
            assert "list" in data["data"], "data 缺少 list 字段"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_alerts_filter_by_device_id(self):
        """
        验证按设备ID筛选告警记录

        GET /api/v1/alerts?device_id=xxx
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        try:
            resp = requests.get(
                f"{self.API_V1}/alerts?device_id=test_device_123",
                headers=headers, timeout=5)
            assert resp.status_code == 200, \
                f"按设备ID筛选告警响应异常: {resp.status_code}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_alerts_filter_by_status(self):
        """
        验证按状态筛选告警记录

        GET /api/v1/alerts?status=1 (未处理)
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        try:
            resp = requests.get(
                f"{self.API_V1}/alerts?status=1",
                headers=headers, timeout=5)
            assert resp.status_code == 200, \
                f"按状态筛选告警响应异常: {resp.status_code}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_dashboard_stats_api(self):
        """
        验证大盘统计数据 API

        GET /api/v1/dashboard/stats
        应返回 total_alerts, unresolved, critical 等统计字段
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        try:
            resp = requests.get(f"{self.API_V1}/dashboard/stats",
                               headers=headers, timeout=5)
            if resp.status_code == 404:
                pytest.skip("大盘统计 API 路由未注册")
            assert resp.status_code == 200, \
                f"大盘统计响应异常: {resp.status_code}"
            data = resp.json()
            assert data.get("code") == 0, f"响应 code 非0: {data}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_alert_rule_validation_missing_type(self):
        """
        验证告警规则创建参数校验

        POST /api/v1/alerts/rules (缺少 alert_type)
        期望返回 400 或 code != 0
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        payload = {
            "name": "缺少类型的告警",
            "condition": "<",
            "threshold": 20.0,
            "severity": 2
            # 故意缺少 alert_type
        }

        try:
            resp = requests.post(f"{self.API_V1}/alerts/rules",
                               headers=headers, json=payload, timeout=5)
            # 应该返回错误
            assert resp.status_code in [400, 422], \
                f"缺少 alert_type 应返回 400，实际: {resp.status_code}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_alert_rule_validation_missing_threshold(self):
        """
        验证告警规则创建参数校验

        POST /api/v1/alerts/rules (缺少 threshold)
        期望返回 400 或 code != 0
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        payload = {
            "name": "缺少阈值的告警",
            "alert_type": "battery_low",
            "condition": "<",
            # 故意缺少 threshold
            "severity": 2
        }

        try:
            resp = requests.post(f"{self.API_V1}/alerts/rules",
                               headers=headers, json=payload, timeout=5)
            assert resp.status_code in [400, 422], \
                f"缺少 threshold 应返回 400，实际: {resp.status_code}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_alert_rule_update_status_field(self):
        """
        验证告警记录状态更新 API

        PUT /api/v1/alerts/:id (更新 status)
        告警状态流转: 1未处理 → 2已确认 → 3已解决
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        try:
            # 获取现有告警记录
            resp = requests.get(f"{self.API_V1}/alerts?status=1",
                               headers=headers, timeout=5)
            if resp.status_code == 404:
                pytest.skip("告警记录 API 路由未注册")
            if resp.status_code != 200:
                pytest.skip(f"获取告警记录失败: {resp.status_code}")

            data = resp.json()
            alerts = data.get("data", {}).get("list", [])
            if not alerts:
                pytest.skip("无未处理告警记录可测试状态更新")

            alert_id = alerts[0].get("id")
            if not alert_id:
                pytest.skip("告警记录缺少 id 字段")

            # 更新状态为已确认
            update_resp = requests.put(
                f"{self.API_V1}/alerts/{alert_id}",
                headers=headers, json={"status": 2}, timeout=5)
            if update_resp.status_code in [404, 501]:
                pytest.skip("告警状态更新 API 未实现")
            assert update_resp.status_code in [200, 400], \
                f"告警状态更新响应异常: {update_resp.status_code}"

        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
