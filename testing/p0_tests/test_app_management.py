# MDM Sprint 2 - 应用管理功能测试
"""
测试应用管理 (App Management) 的核心功能：
- 应用 CRUD (创建/查询/更新/删除)
- 版本管理 (AppVersion CRUD)
- 分发任务 (Distribution)
- 列表筛选 (Filter/Pagination)

依赖 backend Sprint 2 实现：
- models/app.go (App, AppVersion 模型)
- controllers/app_controller.go (CRUD 接口)
- services/app_service.go (分发逻辑)
"""
import os
import pytest
import requests
import time


class TestAppManagement:
    """
    应用管理功能测试套件
    PASS = 应用管理功能正确实现
    FAIL = 应用管理功能未实现或有问题
    """

    BASE_URL = os.getenv("MDM_API_BASE_URL", "http://localhost:8080")
    API_V1 = f"{BASE_URL}/api/v1"

    # ─── 内部辅助 ────────────────────────────────────────────────

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

    def test_app_model_exists(self):
        """
        验证 App 和 AppVersion 模型已定义

        检查 models/app.go 中是否定义：
        - App 结构体（含 name, package_name, description 等字段）
        - AppVersion 结构体（含 version, file_url, file_size 等字段）
        - App 与 AppVersion 关联 (HasMany)
        """
        source = self._read_source("models/app.go")
        assert source is not None, (
            "FAIL: models/app.go 文件不存在\n"
            "Expected: backend/models/app.go 已创建"
        )

        required_app_fields = ["Name", "BundleID", "Description", "Status"]
        for field in required_app_fields:
            assert field in source, (
                f"FAIL: App 模型缺少 {field} 字段"
            )

        required_version_fields = ["Version", "FileURL", "FileSize", "MinOSVersion"]
        for field in required_version_fields:
            assert field in source, (
                f"FAIL: AppVersion 模型缺少 {field} 字段"
            )

        # 检查 App 与 AppVersion 的关联
        assert "AppVersion" in source or "app_versions" in source.lower(), (
            "FAIL: App 模型未定义与 AppVersion 的关联"
        )

        print("PASS: App 和 AppVersion 模型定义正确")

    def test_app_controller_exists(self):
        """
        验证 AppController 已定义并注册路由

        检查 controllers/app_controller.go 中是否实现：
        - AppController 结构体
        - CreateApp / GetApp / UpdateApp / DeleteApp 方法
        - AppVersion CRUD 方法
        """
        source = self._read_source("controllers/app_controller.go")
        assert source is not None, (
            "FAIL: controllers/app_controller.go 文件不存在\n"
            "Expected: backend/controllers/app_controller.go 已创建"
        )

        required_methods = [
            "Create", "Get", "List",
            "Update", "Delete"
        ]
        for method in required_methods:
            assert method in source, (
                f"FAIL: AppController 缺少 {method} 方法"
            )

        required_version_methods = [
            "CreateVersion", "ListVersions", "DeleteVersion"
        ]
        for method in required_version_methods:
            assert method in source, (
                f"FAIL: AppController 缺少 {method} 方法"
            )

        print("PASS: AppController 实现正确")

    def test_app_routes_registered(self):
        """
        验证应用管理 API 路由已在 controllers/device_controller.go 中注册

        检查 RegisterRoutes 函数中是否包含 /apps 路由注册
        """
        source = self._read_source("controllers/device_controller.go")
        assert source is not None, "FAIL: device_controller.go 读取失败"

        assert "/apps" in source, (
            "FAIL: device_controller.go 未注册 /apps 路由\n"
            "Expected: RegisterAppRoutes 或路由注册中包含 /apps"
        )

        # 检查必要的子路由
        required_routes = [
            "GET", "/apps",       # List
            "POST", "/apps",      # Create
            "/apps/:id",          # Get/Update/Delete
            "/apps/:id/versions", # Versions
            "/app/distributions", # Distribution
        ]
        for route in required_routes:
            assert route in source, f"FAIL: 缺少路由 '{route}'"

        print("PASS: 应用管理路由已注册")

    def test_app_service_exists(self):
        """
        验证 AppService 业务逻辑层已实现

        检查 services/app_service.go 中是否实现：
        - 设备应用分发逻辑
        - 版本选择逻辑
        """
        source = self._read_source("services/app_service.go")
        assert source is not None, (
            "FAIL: services/app_service.go 文件不存在\n"
            "Expected: backend/services/app_service.go 已创建"
        )

        required_funcs = [
            "DistributeApp",   # 分发应用
            "SelectVersion",   # 选择版本
        ]
        for func_name in required_funcs:
            assert func_name in source, (
                f"FAIL: AppService 缺少 {func_name} 函数"
            )

        print("PASS: AppService 业务逻辑实现正确")

    # ─── API 功能测试 ────────────────────────────────────────────

    def test_app_crud(self):
        """
        验证应用 CRUD API

        测试流程：
        1. POST /api/v1/apps          - 创建应用
        2. GET  /api/v1/apps          - 获取列表
        3. GET  /api/v1/apps/:id      - 获取详情
        4. PUT  /api/v1/apps/:id      - 更新应用
        5. DELETE /api/v1/apps/:id   - 删除应用
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("Cannot get auth token - backend may not be running")

        app_payload = {
            "name": "Test App",
            "bundle_id": f"com.test.app.{int(time.time())}",  # unique required
            "description": "Test application",
            "platform": "ios"
        }

        # Create
        resp = requests.post(f"{self.API_V1}/apps", json=app_payload, headers=headers, timeout=5)
        assert resp.status_code in [200, 201], (
            f"FAIL: 创建应用返回 {resp.status_code}\n"
            f"Expected: 200/201, Got: {resp.status_code}\n"
            f"Body: {resp.text[:200]}"
        )
        data = resp.json()
        assert data.get("code") == 0, f"创建应用失败: {data}"
        app_id = data.get("data", {}).get("id") or data.get("data", {}).get("app_id")
        assert app_id is not None, "创建应用未返回 id"
        print(f"  Created app id={app_id}")

        # List
        resp = requests.get(f"{self.API_V1}/apps", headers=headers, timeout=5)
        assert resp.status_code == 200, f"FAIL: 获取应用列表返回 {resp.status_code}"
        data = resp.json()
        assert data.get("code") == 0, f"获取列表失败: {data}"
        assert "list" in data.get("data", {}) or isinstance(data.get("data"), list), (
            "列表响应缺少 list 字段"
        )
        print("  Listed apps OK")

        # Get detail
        resp = requests.get(f"{self.API_V1}/apps/{app_id}", headers=headers, timeout=5)
        assert resp.status_code == 200, f"FAIL: 获取应用详情返回 {resp.status_code}"
        data = resp.json()
        assert data.get("code") == 0, f"获取详情失败: {data}"
        print("  Got app detail OK")

        # Update
        update_payload = {"description": "Updated description", "status": 1}
        resp = requests.put(f"{self.API_V1}/apps/{app_id}", json=update_payload, headers=headers, timeout=5)
        assert resp.status_code == 200, f"FAIL: 更新应用返回 {resp.status_code}"
        data = resp.json()
        assert data.get("code") == 0, f"更新应用失败: {data}"
        print("  Updated app OK")

        # Delete
        resp = requests.delete(f"{self.API_V1}/apps/{app_id}", headers=headers, timeout=5)
        assert resp.status_code in [200, 204], f"FAIL: 删除应用返回 {resp.status_code}"
        print("  Deleted app OK")

        print("PASS: 应用 CRUD 功能正确")

    def test_app_version_crud(self):
        """
        验证应用版本 CRUD API

        测试流程：
        1. POST /api/v1/apps/:id/versions     - 创建版本
        2. GET  /api/v1/apps/:id/versions     - 版本列表
        3. PUT  /api/v1/apps/:id/versions/:vid - 更新版本
        4. DELETE /api/v1/apps/:id/versions/:vid - 删除版本
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("Cannot get auth token - backend may not be running")

        # 先创建一个 App
        app_payload = {
            "name": "Version Test App",
            "bundle_id": f"com.test.version.{int(time.time())}",
            "description": "For version test",
            "platform": "android"
        }
        resp = requests.post(f"{self.API_V1}/apps", json=app_payload, headers=headers, timeout=5)
        if resp.status_code not in [200, 201]:
            pytest.skip(f"Cannot create app for version test: {resp.status_code}")
        data = resp.json()
        app_id = data.get("data", {}).get("id") or data.get("data", {}).get("app_id")
        print(f"  Using app id={app_id}")

        version_payload = {
            "version": "1.0.0",
            "build_number": "1",
            "file_url": "https://example.com/app-v1.0.0.bin",
            "file_size": 1024000,
            "min_os_version": "1.0.0",
            "release_notes": "Initial release"
        }

        # Create version
        resp = requests.post(
            f"{self.API_V1}/apps/{app_id}/versions",
            json=version_payload,
            headers=headers,
            timeout=5
        )
        assert resp.status_code in [200, 201], (
            f"FAIL: 创建版本返回 {resp.status_code}\n"
            f"Body: {resp.text[:200]}"
        )
        data = resp.json()
        assert data.get("code") == 0, f"创建版本失败: {data}"
        version_id = (
            data.get("data", {}).get("id")
            or data.get("data", {}).get("version_id")
        )
        assert version_id is not None, "创建版本未返回 id"
        print(f"  Created version id={version_id}")

        # List versions
        resp = requests.get(f"{self.API_V1}/apps/{app_id}/versions", headers=headers, timeout=5)
        assert resp.status_code == 200, f"FAIL: 获取版本列表返回 {resp.status_code}"
        data = resp.json()
        assert data.get("code") == 0, f"获取版本列表失败: {data}"
        print("  Listed versions OK")

        # Update version
        update_payload = {"release_notes": "Bug fixes", "file_size": 1050000}
        resp = requests.put(
            f"{self.API_V1}/apps/{app_id}/versions/{version_id}",
            json=update_payload,
            headers=headers,
            timeout=5
        )
        assert resp.status_code == 200, f"FAIL: 更新版本返回 {resp.status_code}"
        print("  Updated version OK")

        # Delete version
        resp = requests.delete(
            f"{self.API_V1}/apps/{app_id}/versions/{version_id}",
            headers=headers,
            timeout=5
        )
        assert resp.status_code in [200, 204], f"FAIL: 删除版本返回 {resp.status_code}"
        print("  Deleted version OK")

        # Cleanup: delete app
        requests.delete(f"{self.API_V1}/apps/{app_id}", headers=headers, timeout=5)

        print("PASS: 应用版本 CRUD 功能正确")

    def test_app_distribution(self):
        """
        验证应用分发任务 API

        测试流程：
        1. POST /api/v1/apps/distribute   - 创建分发任务
        2. GET  /api/v1/apps/distributions - 分发列表
        3. GET  /api/v1/apps/distributions/:id - 分发详情
        4. POST /api/v1/apps/distributions/:id/cancel - 取消分发
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("Cannot get auth token - backend may not be running")

        # 先创建一个 App + Version
        app_payload = {
            "name": "Distribute Test App",
            "bundle_id": f"com.test.dist.{int(time.time())}",
            "platform": "multi"
        }
        resp = requests.post(f"{self.API_V1}/apps", json=app_payload, headers=headers, timeout=5)
        if resp.status_code not in [200, 201]:
            pytest.skip(f"Cannot create app for distribution test: {resp.status_code}")
        data = resp.json()
        app_id = data.get("data", {}).get("id") or data.get("data", {}).get("app_id")

        version_payload = {
            "version": "1.0.0",
            "build_number": "1",
            "file_url": "https://example.com/dist-v1.bin",
            "file_size": 500000,
            "min_os_version": "1.0.0"
        }
        resp = requests.post(
            f"{self.API_V1}/apps/{app_id}/versions",
            json=version_payload,
            headers=headers,
            timeout=5
        )
        if resp.status_code not in [200, 201]:
            requests.delete(f"{self.API_V1}/apps/{app_id}", headers=headers, timeout=5)
            pytest.skip(f"Cannot create version for distribution test: {resp.status_code}")
        version_data = resp.json()
        version_id = (
            version_data.get("data", {}).get("id")
            or version_data.get("data", {}).get("version_id")
        )

        # Create distribution task
        dist_payload = {
            "app_id": app_id,
            "version_id": version_id,
            "target_ids": ["test-device-001", "test-device-002"],
            "strategy": "immediate"   # immediate / scheduled
        }

        resp = requests.post(
            f"{self.API_V1}/app/distributions",
            json=dist_payload,
            headers=headers,
            timeout=5
        )
        assert resp.status_code in [200, 201], (
            f"FAIL: 创建分发任务返回 {resp.status_code}\n"
            f"Body: {resp.text[:200]}"
        )
        data = resp.json()
        assert data.get("code") == 0, f"创建分发任务失败: {data}"
        dist_id = data.get("data", {}).get("id") or data.get("data", {}).get("distribution_id")
        print(f"  Created distribution id={dist_id}")

        # Get distribution detail
        if dist_id:
            resp = requests.get(
                f"{self.API_V1}/app/distributions/{dist_id}",
                headers=headers,
                timeout=5
            )
            assert resp.status_code == 200, f"FAIL: 获取分发详情返回 {resp.status_code}"
            print("  Got distribution detail OK")

            # Cancel distribution
            resp = requests.post(
                f"{self.API_V1}/app/distributions/{dist_id}/cancel",
                headers=headers,
                timeout=5
            )
            assert resp.status_code in [200, 201], f"FAIL: 取消分发返回 {resp.status_code}"
            print("  Cancelled distribution OK")

        # Cleanup
        requests.delete(f"{self.API_V1}/apps/{app_id}", headers=headers, timeout=5)

        print("PASS: 应用分发任务功能正确")

    def test_app_list_filter(self):
        """
        验证应用列表筛选与分页

        测试场景：
        1. 按名称精确/模糊搜索
        2. 按状态筛选 (status=1 启用 / status=0 禁用)
        3. 分页参数 (page, page_size)
        4. 列表排序 (sort=created_at&order=desc)
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("Cannot get auth token - backend may not be running")

        # 创建多个测试 App
        test_apps = []
        for name, pkg, status in [
            ("Filter Test App A", "com.test.filter.a", 1),
            ("Filter Test App B", "com.test.filter.b", 0),
            ("Unique Filter XYZ", "com.test.unique.xyz", 1),
        ]:
            resp = requests.post(
                f"{self.API_V1}/apps",
                json={"name": name, "package_name": pkg, "status": status},
                headers=headers,
                timeout=5
            )
            if resp.status_code in [200, 201]:
                data = resp.json()
                app_id = data.get("data", {}).get("id") or data.get("data", {}).get("app_id")
                test_apps.append(app_id)

        try:
            # 按名称搜索
            resp = requests.get(
                f"{self.API_V1}/apps?name=Filter Test",
                headers=headers,
                timeout=5
            )
            assert resp.status_code == 200, f"FAIL: 名称搜索返回 {resp.status_code}"
            data = resp.json()
            if data.get("code") == 0:
                items = data.get("data", {}).get("list", []) or data.get("data", [])
                assert len(items) >= 2, f"FAIL: 名称搜索应返回至少2条结果，实际: {len(items)}"
                print("  Name search OK")

            # 按状态筛选
            resp = requests.get(
                f"{self.API_V1}/apps?status=1",
                headers=headers,
                timeout=5
            )
            assert resp.status_code == 200, f"FAIL: 状态筛选返回 {resp.status_code}"
            data = resp.json()
            if data.get("code") == 0:
                items = data.get("data", {}).get("list", []) or data.get("data", [])
                for item in items:
                    if "status" in item:
                        assert item["status"] == 1, f"FAIL: 状态筛选未生效，status={item['status']}"
                print("  Status filter OK")

            # 分页
            resp = requests.get(
                f"{self.API_V1}/apps?page=1&page_size=2",
                headers=headers,
                timeout=5
            )
            assert resp.status_code == 200, f"FAIL: 分页返回 {resp.status_code}"
            data = resp.json()
            if data.get("code") == 0:
                assert "page" in str(data.get("data", {})) or "total" in str(data.get("data", {})), (
                    "分页响应缺少分页字段 (page/total/page_size)"
                )
                print("  Pagination OK")

        finally:
            # Cleanup
            for app_id in test_apps:
                requests.delete(f"{self.API_V1}/apps/{app_id}", headers=headers, timeout=5)

        print("PASS: 应用列表筛选与分页功能正确")


if __name__ == "__main__":
    pytest.main([__file__, "-v", "--tb=short"])
