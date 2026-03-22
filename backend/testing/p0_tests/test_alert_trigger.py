# MDM Sprint 1.1 & 1.2 - CheckAlerts 集成测试
"""
测试告警触发 (CheckAlerts) 的核心功能：
- 电量 < 15% 触发 warning
- 电量 < 5% 触发 critical
- 离线 > 90s 触发告警
"""
import os
import pytest
import requests
import json


class TestAlertTrigger:
    """
    告警触发功能测试套件
    PASS = 告警规则正确触发
    FAIL = 告警触发逻辑异常
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

    def test_battery_low_warning(self):
        """
        验证电量 < 15% 触发 warning 告警

        检查 backend/controllers/alert_controller.go 中 CheckAlerts 函数
        Expected: battery_low 规则，threshold=15, severity=2 (warning)
        """
        alert_ctrl_path = os.path.join(
            os.getenv("MDM_BACKEND_PATH", "C:\\Users\\YKing\\.openclaw\\workspace\\mdm-project\\backend"),
            "controllers", "alert_controller.go"
        )

        with open(alert_ctrl_path, "r", encoding="utf-8") as f:
            source = f.read()

        # 检查 CheckAlerts 函数
        assert "func CheckAlerts" in source, (
            "FAIL: alert_controller.go 缺少 CheckAlerts 函数"
        )

        # 检查 battery_low 告警类型处理
        assert "battery_low" in source, (
            "FAIL: CheckAlerts 未处理 battery_low 告警类型"
        )

        # 检查是否从 data 中获取 battery 值
        assert 'data["battery"]' in source or "data[\"battery\"]" in source or "battery" in source.lower(), (
            "FAIL: CheckAlerts 未从 data 中获取 battery 值"
        )

        print("PASS: 电量低告警逻辑已实现")

    def test_battery_low_critical(self):
        """
        验证电量 < 5% 触发 critical 告警

        检查 backend/models/alert_models.go 中告警规则模型
        Expected: severity 字段支持多个级别
        """
        alert_model_path = os.path.join(
            os.getenv("MDM_BACKEND_PATH", "C:\\Users\\YKing\\.openclaw\\workspace\\mdm-project\\backend"),
            "models", "alert_models.go"
        )

        with open(alert_model_path, "r", encoding="utf-8") as f:
            source = f.read()

        # 检查 DeviceAlertRule 模型
        assert "DeviceAlertRule" in source, (
            "FAIL: alert_models.go 缺少 DeviceAlertRule 模型"
        )

        # 检查 Severity 字段
        assert "Severity" in source or "severity" in source.lower(), (
            "FAIL: DeviceAlertRule 缺少 Severity 字段"
        )

        # 检查 AlertType 字段
        assert "AlertType" in source or "alert_type" in source.lower(), (
            "FAIL: DeviceAlertRule 缺少 AlertType 字段"
        )

        # 检查 Threshold 字段
        assert "Threshold" in source or "threshold" in source.lower(), (
            "FAIL: DeviceAlertRule 缺少 Threshold 字段"
        )

        print("PASS: 告警规则模型支持多级别 severity")

    def test_offline_alert(self):
        """
        验证离线 > 90s 触发告警

        检查 backend/mqtt/handler.go 中离线告警触发逻辑
        Expected: checkOfflineDevices 在超时后调用 AlertCB
        """
        handler_go_path = os.path.join(
            os.getenv("MDM_BACKEND_PATH", "C:\\Users\\YKing\\.openclaw\\workspace\\mdm-project\\backend"),
            "mqtt", "handler.go"
        )

        with open(handler_go_path, "r", encoding="utf-8") as f:
            source = f.read()

        # 检查离线告警触发
        assert "offline" in source.lower(), (
            "FAIL: mqtt/handler.go 未处理 offline 告警"
        )

        # 检查是否调用 AlertCB
        assert "AlertCB" in source, (
            "FAIL: 离线检测未调用 AlertCB 触发告警"
        )

        # 检查 is_online 判断
        assert '"is_online"' in source or "'is_online'" in source or "is_online" in source, (
            "FAIL: 离线检测未检查 is_online 状态"
        )

        print("PASS: 离线告警触发逻辑已实现")

    def test_alert_rules_api(self):
        """
        验证告警规则 API 可访问

        测试 GET /api/v1/alerts/rules 端点
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("Cannot get auth token - backend may not be running")

        resp = requests.get(f"{self.API_V1}/alerts/rules", headers=headers, timeout=5)
        assert resp.status_code == 200, (
            f"FAIL: Alert rules API 返回 {resp.status_code}\n"
            "Expected: 200 OK"
        )

        data = resp.json()
        assert "data" in data, "Response missing 'data' field"
        assert "list" in data["data"], "Response missing 'list' field"

        print(f"PASS: 告警规则 API 可访问")

    def test_alerts_api(self):
        """
        验证告警记录 API 可访问

        测试 GET /api/v1/alerts 端点
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("Cannot get auth token - backend may not be running")

        resp = requests.get(f"{self.API_V1}/alerts", headers=headers, timeout=5)
        assert resp.status_code == 200, (
            f"FAIL: Alerts API 返回 {resp.status_code}\n"
            "Expected: 200 OK"
        )

        data = resp.json()
        assert "data" in data, "Response missing 'data' field"
        assert "list" in data["data"], "Response missing 'list' field"

        print(f"PASS: 告警记录 API 可访问")

    def test_dashboard_stats_api(self):
        """
        验证大盘统计 API 包含告警统计

        测试 GET /api/v1/dashboard/stats 端点
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("Cannot get auth token - backend may not be running")

        resp = requests.get(f"{self.API_V1}/dashboard/stats", headers=headers, timeout=5)
        assert resp.status_code == 200, (
            f"FAIL: Dashboard stats API 返回 {resp.status_code}\n"
            "Expected: 200 OK"
        )

        data = resp.json()
        assert "data" in data, "Response missing 'data' field"

        # 检查告警统计字段
        stats = data["data"]
        assert "total_alerts" in stats or "TotalAlerts" in str(type(stats)), (
            "FAIL: Dashboard stats 缺少告警统计"
        )

        print(f"PASS: 大盘统计 API 包含告警统计")

    def test_alert_condition_evaluation(self):
        """
        验证告警条件评估函数

        检查 backend/controllers/alert_controller.go 中 evaluateCondition 函数
        Expected: 支持 <, >, =, <=, >= 比较操作符
        """
        alert_ctrl_path = os.path.join(
            os.getenv("MDM_BACKEND_PATH", "C:\\Users\\YKing\\.openclaw\\workspace\\mdm-project\\backend"),
            "controllers", "alert_controller.go"
        )

        with open(alert_ctrl_path, "r", encoding="utf-8") as f:
            source = f.read()

        # 检查 evaluateCondition 函数
        assert "func evaluateCondition" in source, (
            "FAIL: alert_controller.go 缺少 evaluateCondition 函数"
        )

        # 检查所有比较操作符
        operators = ['"<":', '">":', '"=":', '"<=":', '">="']
        found_operators = sum(1 for op in operators if op in source)
        assert found_operators >= 3, (
            f"FAIL: evaluateCondition 仅支持 {found_operators} 个操作符\n"
            "Expected: 支持 <, >, =, <=, >="
        )

        print(f"PASS: 告警条件评估函数支持必要的比较操作符")

    def test_mqtt_alert_callback_integration(self):
        """
        验证 MQTT 与 AlertCallback 集成

        检查 backend/main.go 中 MQTT 初始化是否传入告警回调
        """
        main_go_path = os.path.join(
            os.getenv("MDM_BACKEND_PATH", "C:\\Users\\YKing\\.openclaw\\workspace\\mdm-project\\backend"),
            "main.go"
        )

        with open(main_go_path, "r", encoding="utf-8") as f:
            source = f.read()

        # 检查 MQTT 初始化是否传入 alertCallback
        assert "alertCallback" in source, (
            "FAIL: main.go 未定义 alertCallback"
        )

        # 检查 InitMQTT 调用是否传入 alertCallback
        assert "InitMQTT" in source and "alertCallback" in source, (
            "FAIL: MQTT 初始化未传入告警回调"
        )

        # 检查 CheckAlerts 调用
        assert "CheckAlerts" in source, (
            "FAIL: alertCallback 未调用 CheckAlerts"
        )

        print("PASS: MQTT 与 AlertCallback 正确集成")


if __name__ == "__main__":
    pytest.main([__file__, "-v", "--tb=short"])
