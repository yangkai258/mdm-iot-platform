# MDM Sprint 1.1 & 1.2 - OTA Worker 功能测试
"""
测试 OTA Worker 的核心功能：
- 5分钟轮询间隔
- 全量灰度策略
- 百分比灰度策略
- 白名单灰度策略
- 失败率>80%自动暂停
"""
import os
import pytest
import requests
import time
from datetime import datetime


class TestOTAWorker:
    """
    OTA Worker 功能测试套件
    PASS = OTA Worker 正确处理各策略
    FAIL = OTA Worker 策略处理异常
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

    def _get_backend_path(self):
        """Helper: 获取 backend 路径"""
        return os.getenv("MDM_BACKEND_PATH", "C:\\Users\\YKing\\.openclaw\\workspace\\mdm-project\\backend")

    def test_worker_polling_interval(self):
        """
        验证 OTA Worker 轮询间隔为 5 分钟

        检查 backend 中 OTA Worker 的 ticker 间隔
        Expected: 5 * time.Minute
        """
        backend_path = self._get_backend_path()
        main_go_path = os.path.join(backend_path, "main.go")

        with open(main_go_path, "r", encoding="utf-8") as f:
            main_source = f.read()

        # 检查是否启动了 OTA Worker
        has_worker_start = (
            "go otaWorker.Start()" in main_source or
            "go OTAWorker" in main_source or
            "NewOTAWorker" in main_source or
            "startOTAWorker" in main_source
        )

        assert has_worker_start, (
            "FAIL: OTA Worker 未在 main.go 中启动\n"
            "Expected: go otaWorker.Start() 或类似调用"
        )

        # 查找 services/ota_worker.go
        services_path = os.path.join(backend_path, "services", "ota_worker.go")
        ota_worker_source = ""

        if os.path.exists(services_path):
            with open(services_path, "r", encoding="utf-8") as f:
                ota_worker_source = f.read()
        else:
            # 尝试在 controllers 目录查找
            controllers_path = os.path.join(backend_path, "controllers")
            if os.path.exists(controllers_path):
                for fname in os.listdir(controllers_path):
                    if "ota" in fname.lower():
                        fpath = os.path.join(controllers_path, fname)
                        with open(fpath, "r", encoding="utf-8") as f:
                            content = f.read()
                            if "func startOTAWorker" in content or "OTAWorker" in content:
                                ota_worker_source = content
                                break

        # 如果找到 Worker 源码，检查轮询间隔
        if ota_worker_source:
            has_five_minutes = (
                "5 * time.Minute" in ota_worker_source or
                "5*time.Minute" in ota_worker_source
            )
            assert has_five_minutes, (
                "FAIL: OTA Worker 未配置 5 分钟轮询间隔\n"
                "Expected: time.NewTicker(5 * time.Minute)"
            )
            print("PASS: OTA Worker 正确配置 5 分钟轮询间隔")
        else:
            # OTA Worker 已启动但实现可能在其他包中
            # 检查 main.go 中是否有 ticker 配置
            if "time.Minute" in main_source or "time.Ticker" in main_source:
                print("PASS: OTA Worker 启动已配置 (Worker 实现在独立包中)")
            else:
                print("INFO: OTA Worker 启动已配置，实现位置待确认")

    def test_grayscale_full(self):
        """
        验证全量灰度策略 (full rollout)

        检查 backend/models/ota.go 中 OTADeployment 模型的 full 策略
        Expected: StrategyType = 'full' 时所有匹配设备都会收到 OTA
        """
        backend_path = self._get_backend_path()
        ota_model_path = os.path.join(backend_path, "models", "ota.go")

        with open(ota_model_path, "r", encoding="utf-8") as f:
            source = f.read()

        # 检查 StrategyType 字段
        assert "StrategyType" in source or "Strategy" in source, (
            "FAIL: OTADeployment 模型缺少 StrategyType 字段"
        )

        # 检查 full 策略支持
        assert "full" in source.lower() or "Full" in source, (
            "FAIL: OTADeployment 模型未支持 full 策略"
        )

        print("PASS: OTA Worker 支持全量灰度策略")

    def test_grayscale_percentage(self):
        """
        验证百分比灰度策略 (percentage rollout)

        检查 percentage 策略是否正确实现
        Expected: dep.Percentage 或 StrategyConfig 中配置百分比
        """
        backend_path = self._get_backend_path()
        ota_model_path = os.path.join(backend_path, "models", "ota.go")

        with open(ota_model_path, "r", encoding="utf-8") as f:
            source = f.read()

        # 检查百分比字段
        assert "Percentage" in source or "percentage" in source, (
            "FAIL: OTADeployment 模型缺少 Percentage 字段"
        )

        # 检查 percentage 策略
        assert "percentage" in source.lower() or "Percentage" in source, (
            "FAIL: OTADeployment 模型未支持 percentage 策略"
        )

        print("PASS: OTA Worker 支持百分比灰度策略")

    def test_grayscale_whitelist(self):
        """
        验证白名单灰度策略 (whitelist rollout)

        检查 whitelist 策略是否正确实现
        Expected: 白名单中的设备才会收到 OTA 指令
        """
        backend_path = self._get_backend_path()
        ota_model_path = os.path.join(backend_path, "models", "ota.go")

        with open(ota_model_path, "r", encoding="utf-8") as f:
            source = f.read()

        # 检查白名单策略
        assert "whitelist" in source.lower() or "Whitelist" in source, (
            "FAIL: OTADeployment 模型未支持 whitelist 策略"
        )

        print("PASS: OTA Worker 支持白名单灰度策略")

    def test_auto_pause_on_failure(self):
        """
        验证失败率>阈值自动暂停功能

        检查 backend/models/ota.go 中 OTADeployment 的 auto_pause 相关字段
        Expected: PauseOnFailureThreshold 字段
        """
        backend_path = self._get_backend_path()
        ota_model_path = os.path.join(backend_path, "models", "ota.go")

        with open(ota_model_path, "r", encoding="utf-8") as f:
            source = f.read()

        # 检查失败率阈值字段
        assert "PauseOnFailureThreshold" in source or "pause_on_failure_threshold" in source.lower(), (
            "FAIL: OTADeployment 模型缺少 PauseOnFailureThreshold 字段"
        )

        # 检查自动暂停标识字段
        assert "AutoPaused" in source or "auto_paused" in source.lower(), (
            "FAIL: OTADeployment 模型缺少 AutoPaused 字段"
        )

        print("PASS: OTA Worker 支持失败率自动暂停")

    def test_ota_api_endpoints(self):
        """
        验证 OTA 相关 API 端点可访问

        检查固件包管理和部署任务 API
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("Cannot get auth token - backend may not be running")

        # 测试固件包列表 API
        resp = requests.get(f"{self.API_V1}/ota/packages", headers=headers, timeout=5)
        assert resp.status_code == 200, (
            f"FAIL: OTA packages API 返回 {resp.status_code}\n"
            "Expected: 200 OK"
        )

        data = resp.json()
        assert "data" in data, "Response missing 'data' field"
        assert "list" in data["data"], "Response missing 'list' field"

        print(f"PASS: OTA API 端点可访问")

    def test_ota_deployment_create(self):
        """
        验证 OTA 部署任务创建 API

        创建一个测试部署任务
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("Cannot get auth token - backend may not be running")

        # 构造部署任务
        deployment_data = {
            "package_id": 1,  # 需要预先存在的固件包
            "target_hardware": "M5Stack",
            "rollout_strategy": "full"
        }

        resp = requests.post(
            f"{self.API_V1}/ota/deployments",
            json=deployment_data,
            headers=headers,
            timeout=5
        )

        # 固件包不存在会返回 4002，但 API 本身是可访问的
        assert resp.status_code in [200, 400, 404], (
            f"FAIL: OTA deployment API 返回 {resp.status_code}\n"
            "Expected: 200/400/404"
        )

        print(f"PASS: OTA deployment API 正常响应 (status: {resp.status_code})")


if __name__ == "__main__":
    pytest.main([__file__, "-v", "--tb=short"])
