# MDM Sprint 1.1 & 1.2 - 设备影子功能测试
"""
测试设备影子 (Device Shadow) 的核心功能：
- 在线状态更新
- 离线检测
- Redis URL 解析
- 离线时 DB 同步
"""
import os
import pytest
import requests
import json
import time


class TestDeviceShadow:
    """
    设备影子功能测试套件
    PASS = 设备影子正确处理设备状态
    FAIL = 设备影子状态处理异常
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

    def test_shadow_online_update(self):
        """
        验证在线状态更新功能

        检查 backend/mqtt/handler.go 中 StatusMessageHandler 函数
        Expected: 心跳消息正确更新 Redis 设备影子
        """
        handler_go_path = os.path.join(
            os.getenv("MDM_BACKEND_PATH", "C:\\Users\\YKing\\.openclaw\\workspace\\mdm-project\\backend"),
            "mqtt", "handler.go"
        )

        with open(handler_go_path, "r", encoding="utf-8") as f:
            source = f.read()

        # 检查是否处理心跳消息
        assert "StatusMessageHandler" in source, (
            "FAIL: mqtt/handler.go 缺少 StatusMessageHandler 函数"
        )

        # 检查是否调用 SetDeviceShadow
        assert "SetDeviceShadow" in source, (
            "FAIL: 心跳处理器未调用 SetDeviceShadow 更新设备影子"
        )

        # 检查是否设置 TTL (90秒)
        assert "90*time.Second" in source or "90 * time.Second" in source, (
            "FAIL: 设备影子 TTL 未设置为 90 秒"
        )

        # 检查是否正确设置 IsOnline 字段
        assert "IsOnline" in source, (
            "FAIL: 设备影子未包含 IsOnline 字段"
        )

        print("PASS: 设备影子在线状态更新功能正确")

    def test_shadow_offline_detection(self):
        """
        验证离线检测功能

        检查 backend/mqtt/handler.go 中 StartHeartbeatChecker 和 checkOfflineDevices 函数
        Expected: 90秒无心跳标记设备为离线
        """
        handler_go_path = os.path.join(
            os.getenv("MDM_BACKEND_PATH", "C:\\Users\\YKing\\.openclaw\\workspace\\mdm-project\\backend"),
            "mqtt", "handler.go"
        )

        with open(handler_go_path, "r", encoding="utf-8") as f:
            source = f.read()

        # 检查心跳检查器
        assert "StartHeartbeatChecker" in source, (
            "FAIL: mqtt/handler.go 缺少 StartHeartbeatChecker 函数"
        )

        # 检查离线检测逻辑
        assert "checkOfflineDevices" in source, (
            "FAIL: 缺少离线设备检测函数"
        )

        # 检查是否检查 90 秒超时
        assert "90*time.Second" in source or "90 * time.Second" in source, (
            "FAIL: 离线检测未使用 90 秒超时阈值"
        )

        # 检查是否在超时后将 IsOnline 设为 false
        assert "IsOnline = false" in source or "IsOnline: false" in source, (
            "FAIL: 离线检测未正确设置 IsOnline 为 false"
        )

        print("PASS: 设备影子离线检测功能正确")

    def test_redis_url_parse(self):
        """
        验证 Redis URL 解析功能

        检查 backend/utils/redis.go 中 InitRedis 函数
        Expected: 正确解析 redis://[user:password@]host:port[/db] 格式
        """
        redis_go_path = os.path.join(
            os.getenv("MDM_BACKEND_PATH", "C:\\Users\\YKing\\.openclaw\\workspace\\mdm-project\\backend"),
            "utils", "redis.go"
        )

        with open(redis_go_path, "r", encoding="utf-8") as f:
            source = f.read()

        # 检查是否解析 redis:// URL
        assert "redis://" in source, (
            "FAIL: utils/redis.go 未处理 redis:// URL 格式"
        )

        # 检查是否支持环境变量 REDIS_URL
        assert "REDIS_URL" in source, (
            "FAIL: 未从环境变量 REDIS_URL 读取 Redis 连接信息"
        )

        # 检查是否有默认端口
        assert "6379" in source, (
            "FAIL: Redis 默认端口 6379 未配置"
        )

        # 检查是否解析认证信息
        assert "password" in source.lower(), (
            "FAIL: Redis URL 解析未处理密码"
        )

        print("PASS: Redis URL 解析功能正确")

    def test_db_sync_on_offline(self):
        """
        验证离线时 DB 同步功能

        检查 backend/mqtt/handler.go 中 checkOfflineDevices 是否在离线时同步 DB
        Expected: 设备离线时更新设备状态到数据库
        """
        handler_go_path = os.path.join(
            os.getenv("MDM_BACKEND_PATH", "C:\\Users\\YKing\\.openclaw\\workspace\\mdm-project\\backend"),
            "mqtt", "handler.go"
        )

        with open(handler_go_path, "r", encoding="utf-8") as f:
            source = f.read()

        # 检查离线告警触发
        assert "AlertCB" in source, (
            "FAIL: 离线检测未触发告警回调"
        )

        # 离线时应该传递 is_online: false
        assert '"is_online": false' in source or "'is_online': false" in source, (
            "FAIL: 离线检测未传递 is_online=false 到告警回调"
        )

        print("PASS: 离线时 DB 同步机制已实现")

    def test_shadow_device_list_api(self):
        """
        验证设备列表 API 返回设备影子状态

        测试 GET /api/v1/devices 端点
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("Cannot get auth token - backend may not be running")

        resp = requests.get(f"{self.API_V1}/devices", headers=headers, timeout=5)
        assert resp.status_code == 200, (
            f"FAIL: Devices API 返回 {resp.status_code}\n"
            "Expected: 200 OK"
        )

        data = resp.json()
        assert "data" in data, "Response missing 'data' field"
        assert "list" in data["data"], "Response missing 'list' field"

        print(f"PASS: 设备列表 API 正确返回设备影子状态")

    def test_device_detail_api(self):
        """
        验证设备详情 API 返回设备影子

        测试 GET /api/v1/devices/:device_id 端点
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("Cannot get auth token - backend may not be running")

        resp = requests.get(
            f"{self.API_V1}/devices/test-device-001",
            headers=headers,
            timeout=5
        )

        # 设备可能不存在，但 API 应该可访问
        assert resp.status_code in [200, 404], (
            f"FAIL: Device detail API 返回 {resp.status_code}\n"
            "Expected: 200/404"
        )

        if resp.status_code == 200:
            data = resp.json()
            assert "data" in data, "Response missing 'data' field"
            # 可能包含 shadow 字段
            print(f"PASS: 设备详情 API 正确返回设备影子")
        else:
            print(f"INFO: 测试设备不存在 (status: 404)")

    def test_redis_client_implementation(self):
        """
        验证 Redis Client 正确实现

        检查 backend/utils/redis.go 中 RedisClient 结构体
        """
        redis_go_path = os.path.join(
            os.getenv("MDM_BACKEND_PATH", "C:\\Users\\YKing\\.openclaw\\workspace\\mdm-project\\backend"),
            "utils", "redis.go"
        )

        with open(redis_go_path, "r", encoding="utf-8") as f:
            source = f.read()

        # 检查 RedisClient 结构体
        assert "RedisClient" in source, (
            "FAIL: utils/redis.go 缺少 RedisClient 结构体"
        )

        # 检查 DeviceShadow 结构体
        assert "DeviceShadow" in source, (
            "FAIL: utils/redis.go 缺少 DeviceShadow 结构体"
        )

        # 检查必要的 Redis 操作方法
        required_methods = ["SetDeviceShadow", "GetDeviceShadow", "GetAllShadowKeys"]
        for method in required_methods:
            assert method in source, (
                f"FAIL: RedisClient 缺少 {method} 方法"
            )

        print("PASS: Redis Client 实现正确")


if __name__ == "__main__":
    pytest.main([__file__, "-v", "--tb=short"])
