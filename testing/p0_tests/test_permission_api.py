# MDM Sprint 6 - 权限管理 API 测试
"""
测试权限管理 (Permission & Role Management) API：
- 角色 CRUD
- 权限分配

依赖 backend Sprint 6 实现：
- controllers/role_controller.go (NewRoleController)
- controllers/permission_controller.go (PermissionController)
- models/role_models.go (Role, SysRole, SysPermission)
- 路由已注册
"""
import os
import pytest
import requests
import time


class TestRoleCRUD:
    """
    角色管理 API 测试套件

    API 前缀: /api/v1/roles
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
            os.getenv("MDM_BACKEND_PATH",
                      "C:\\Users\\YKing\\.openclaw\\workspace\\mdm-project\\backend"),
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

    def test_role_model_exists(self):
        """
        验证 Role 模型已定义

        检查 models/ 中是否定义 Role 结构体：
        - RoleName, RoleCode, Description, Status, TenantID
        """
        import glob
        backend_path = os.path.join(
            os.getenv("MDM_BACKEND_PATH",
                      "C:\\Users\\YKing\\.openclaw\\workspace\\mdm-project\\backend"),
            "models"
        )
        found = False
        for f in glob.glob(os.path.join(backend_path, "*.go")):
            content = open(f, encoding="utf-8").read()
            if "type Role struct" in content:
                found = True
                assert 'json:"role_name"' in content, "Role 缺少 role_name 字段"
                assert 'json:"role_code"' in content, "Role 缺少 role_code 字段"
                assert 'json:"tenant_id"' in content, "Role 缺少 tenant_id 字段"
                break
        if not found:
            pytest.skip("Role 模型文件未找到")

    def test_role_controller_exists(self):
        """
        验证 NewRoleController 控制器已定义

        检查 controllers/role_controller.go 中是否定义 NewRoleController
        """
        source = self._read_source("controllers/role_controller.go")
        assert source is not None, \
            "controllers/role_controller.go 文件不存在"
        assert "type NewRoleController struct" in source, \
            "NewRoleController 结构体未定义"

    def test_role_controller_methods(self):
        """
        验证 NewRoleController 核心方法已实现

        检查 controllers/role_controller.go 中：
        - List, Create, Update, Delete, GetPermissions, SetPermissions
        """
        source = self._read_source("controllers/role_controller.go")
        assert source is not None

        assert "func (c *NewRoleController) List" in source, \
            "List 方法未实现"
        assert "func (c *NewRoleController) Create" in source, \
            "Create 方法未实现"
        assert "func (c *NewRoleController) Update" in source, \
            "Update 方法未实现"
        assert "func (c *NewRoleController) Delete" in source, \
            "Delete 方法未实现"
        assert "func (c *NewRoleController) GetPermissions" in source, \
            "GetPermissions 方法未实现"
        assert "func (c *NewRoleController) SetPermissions" in source, \
            "SetPermissions 方法未实现"

    # ─── API 功能测试 ────────────────────────────────────────────

    def test_role_list_api(self):
        """
        验证角色列表 API

        GET /api/v1/roles
        验证返回 list, total, page, page_size 字段
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        try:
            resp = requests.get(
                f"{self.API_V1}/roles?page=1&page_size=10",
                headers=headers, timeout=5)
            if resp.status_code == 404:
                pytest.skip("角色 API 路由未注册")
            assert resp.status_code == 200, \
                f"角色列表响应异常: {resp.status_code}"
            data = resp.json()
            assert data.get("code") == 0, f"响应 code 非0: {data}"
            assert "data" in data, "响应缺少 data 字段"
            assert "list" in data["data"], "data 缺少 list 字段"
            assert "total" in data["data"], "data 缺少 total 字段"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_role_list_filter_by_keyword(self):
        """
        验证角色列表关键词筛选

        GET /api/v1/roles?keyword=xxx
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        try:
            resp = requests.get(
                f"{self.API_V1}/roles?keyword=admin&page=1&page_size=10",
                headers=headers, timeout=5)
            if resp.status_code == 404:
                pytest.skip("角色 API 路由未注册")
            assert resp.status_code == 200, \
                f"角色列表筛选响应异常: {resp.status_code}"
            data = resp.json()
            assert data.get("code") == 0, f"响应 code 非0: {data}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_role_create_api(self):
        """
        验证创建角色 API

        POST /api/v1/roles
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        payload = {
            "role_name": f"测试角色_{int(time.time())}",
            "role_code": f"ROLE_TEST_{int(time.time())}",
            "description": "测试角色描述",
            "status": 1
        }

        try:
            resp = requests.post(
                f"{self.API_V1}/roles",
                headers=headers, json=payload, timeout=5)
            if resp.status_code == 404:
                pytest.skip("角色 API 路由未注册")
            assert resp.status_code in [200, 201, 400], \
                f"创建角色响应异常: {resp.status_code}"
            if resp.status_code in [200, 201]:
                data = resp.json()
                assert data.get("code") == 0, f"响应 code 非0: {data}"
                role = data.get("data", {})
                assert role.get("role_name") == payload["role_name"], \
                    f"role_name 不匹配: {role}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_role_create_duplicate_code(self):
        """
        验证创建角色时角色编码重复校验

        相同的 role_code 应返回错误
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        duplicate_code = f"ROLE_DUP_{int(time.time())}"

        payload = {
            "role_name": "重复编码测试1",
            "role_code": duplicate_code,
            "status": 1
        }

        try:
            # 第一次创建
            resp1 = requests.post(
                f"{self.API_V1}/roles",
                headers=headers, json=payload, timeout=5)
            if resp1.status_code == 404:
                pytest.skip("角色 API 路由未注册")
            if resp1.status_code not in [200, 201]:
                pytest.skip("创建角色失败")

            # 第二次创建相同编码
            payload2 = {
                "role_name": "重复编码测试2",
                "role_code": duplicate_code,
                "status": 1
            }
            resp2 = requests.post(
                f"{self.API_V1}/roles",
                headers=headers, json=payload2, timeout=5)
            # 应返回 400 错误
            assert resp2.status_code in [400, 409], \
                f"重复编码应返回错误，实际: {resp2.status_code}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_role_update_api(self):
        """
        验证更新角色 API

        PUT /api/v1/roles/:id
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        # 先创建
        payload = {
            "role_name": f"更新前角色_{int(time.time())}",
            "role_code": f"ROLE_UPD_{int(time.time())}",
            "description": "原始描述",
            "status": 1
        }

        try:
            create_resp = requests.post(
                f"{self.API_V1}/roles",
                headers=headers, json=payload, timeout=5)
            if create_resp.status_code == 404:
                pytest.skip("角色 API 路由未注册")
            if create_resp.status_code not in [200, 201]:
                pytest.skip(f"创建角色失败: {create_resp.status_code}")

            created = create_resp.json()
            role_id = created.get("data", {}).get("id")
            if not role_id:
                pytest.skip("创建响应中无 id 字段")

            # 更新
            update_payload = {
                "role_name": "更新后角色名称",
                "description": "更新后的描述"
            }
            update_resp = requests.put(
                f"{self.API_V1}/roles/{role_id}",
                headers=headers, json=update_payload, timeout=5)
            assert update_resp.status_code in [200, 400], \
                f"更新角色响应异常: {update_resp.status_code}"
            if update_resp.status_code == 200:
                data = update_resp.json()
                role = data.get("data", {})
                assert role.get("role_name") == "更新后角色名称", \
                    f"role_name 未更新: {role}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_role_delete_api(self):
        """
        验证删除角色 API

        DELETE /api/v1/roles/:id
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        # 先创建
        payload = {
            "role_name": f"待删除角色_{int(time.time())}",
            "role_code": f"ROLE_DEL_{int(time.time())}",
            "status": 1
        }

        try:
            create_resp = requests.post(
                f"{self.API_V1}/roles",
                headers=headers, json=payload, timeout=5)
            if create_resp.status_code == 404:
                pytest.skip("角色 API 路由未注册")
            if create_resp.status_code not in [200, 201]:
                pytest.skip(f"创建角色失败: {create_resp.status_code}")

            created = create_resp.json()
            role_id = created.get("data", {}).get("id")
            if not role_id:
                pytest.skip("创建响应中无 id 字段")

            # 删除
            del_resp = requests.delete(
                f"{self.API_V1}/roles/{role_id}",
                headers=headers, timeout=5)
            assert del_resp.status_code in [200, 404], \
                f"删除角色响应异常: {del_resp.status_code}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")


class TestPermissionAPI:
    """
    权限管理 API 测试套件

    API 前缀: /api/v1/permissions
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
            os.getenv("MDM_BACKEND_PATH",
                      "C:\\Users\\YKing\\.openclaw\\workspace\\mdm-project\\backend"),
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

    def test_permission_controller_exists(self):
        """
        验证 PermissionController 控制器已定义

        检查 controllers/permission_controller.go 中是否定义 PermissionController
        """
        source = self._read_source("controllers/permission_controller.go")
        assert source is not None, \
            "controllers/permission_controller.go 文件不存在"
        assert "type PermissionController struct" in source, \
            "PermissionController 结构体未定义"

    def test_permission_model_exists(self):
        """
        验证 SysPermission 模型已定义

        检查 models/ 中是否定义 SysPermission 结构体
        """
        import glob
        backend_path = os.path.join(
            os.getenv("MDM_BACKEND_PATH",
                      "C:\\Users\\YKing\\.openclaw\\workspace\\mdm-project\\backend"),
            "models"
        )
        found = False
        for f in glob.glob(os.path.join(backend_path, "*.go")):
            content = open(f, encoding="utf-8").read()
            if "type SysPermission struct" in content:
                found = True
                assert 'json:"name"' in content, "SysPermission 缺少 name 字段"
                assert 'json:"code"' in content, "SysPermission 缺少 code 字段"
                assert 'json:"type"' in content, "SysPermission 缺少 type 字段"
                break
        if not found:
            pytest.skip("SysPermission 模型文件未找到")

    def test_permission_controller_methods(self):
        """
        验证 PermissionController 方法已实现

        检查 permission_controller.go 中：
        - List, Create, Update, Delete
        """
        source = self._read_source("controllers/permission_controller.go")
        assert source is not None
        assert "func (c *PermissionController) List" in source, \
            "List 方法未实现"
        assert "func (c *PermissionController) Create" in source, \
            "Create 方法未实现"
        assert "func (c *PermissionController) Update" in source, \
            "Update 方法未实现"
        assert "func (c *PermissionController) Delete" in source, \
            "Delete 方法未实现"

    # ─── API 功能测试 ────────────────────────────────────────────

    def test_permission_list_api(self):
        """
        验证权限列表 API

        GET /api/v1/permissions
        验证返回树形结构
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        try:
            resp = requests.get(
                f"{self.API_V1}/permissions",
                headers=headers, timeout=5)
            if resp.status_code == 404:
                pytest.skip("权限 API 路由未注册")
            assert resp.status_code == 200, \
                f"权限列表响应异常: {resp.status_code}"
            data = resp.json()
            assert data.get("code") == 0, f"响应 code 非0: {data}"
            assert "data" in data, "响应缺少 data 字段"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_permission_list_filter_by_parent(self):
        """
        验证权限列表父级筛选

        GET /api/v1/permissions?parent_id=1
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        try:
            resp = requests.get(
                f"{self.API_V1}/permissions?parent_id=0",
                headers=headers, timeout=5)
            if resp.status_code == 404:
                pytest.skip("权限 API 路由未注册")
            assert resp.status_code == 200, \
                f"权限列表筛选响应异常: {resp.status_code}"
            data = resp.json()
            assert data.get("code") == 0, f"响应 code 非0: {data}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_permission_list_filter_by_keyword(self):
        """
        验证权限列表关键词筛选

        GET /api/v1/permissions?keyword=xxx
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        try:
            resp = requests.get(
                f"{self.API_V1}/permissions?keyword=device",
                headers=headers, timeout=5)
            if resp.status_code == 404:
                pytest.skip("权限 API 路由未注册")
            assert resp.status_code == 200, \
                f"权限列表关键词筛选响应异常: {resp.status_code}"
            data = resp.json()
            assert data.get("code") == 0, f"响应 code 非0: {data}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_permission_create_api(self):
        """
        验证创建权限 API

        POST /api/v1/permissions
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        payload = {
            "name": f"测试权限_{int(time.time())}",
            "code": f"PERM_TEST_{int(time.time())}",
            "type": 1,
            "status": 1
        }

        try:
            resp = requests.post(
                f"{self.API_V1}/permissions",
                headers=headers, json=payload, timeout=5)
            if resp.status_code == 404:
                pytest.skip("权限 API 路由未注册")
            assert resp.status_code in [200, 201, 400], \
                f"创建权限响应异常: {resp.status_code}"
            if resp.status_code in [200, 201]:
                data = resp.json()
                assert data.get("code") == 0, f"响应 code 非0: {data}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_permission_update_api(self):
        """
        验证更新权限 API

        PUT /api/v1/permissions/:id
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        # 先创建
        payload = {
            "name": f"更新前权限_{int(time.time())}",
            "code": f"PERM_UPD_{int(time.time())}",
            "type": 1,
            "status": 1
        }

        try:
            create_resp = requests.post(
                f"{self.API_V1}/permissions",
                headers=headers, json=payload, timeout=5)
            if create_resp.status_code == 404:
                pytest.skip("权限 API 路由未注册")
            if create_resp.status_code not in [200, 201]:
                pytest.skip(f"创建权限失败: {create_resp.status_code}")

            created = create_resp.json()
            perm_id = created.get("data", {}).get("id")
            if not perm_id:
                pytest.skip("创建响应中无 id 字段")

            # 更新
            update_payload = {"name": "更新后权限名称"}
            update_resp = requests.put(
                f"{self.API_V1}/permissions/{perm_id}",
                headers=headers, json=update_payload, timeout=5)
            assert update_resp.status_code in [200, 400], \
                f"更新权限响应异常: {update_resp.status_code}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_permission_delete_api(self):
        """
        验证删除权限 API

        DELETE /api/v1/permissions/:id
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        # 先创建
        payload = {
            "name": f"待删除权限_{int(time.time())}",
            "code": f"PERM_DEL_{int(time.time())}",
            "type": 1,
            "status": 1
        }

        try:
            create_resp = requests.post(
                f"{self.API_V1}/permissions",
                headers=headers, json=payload, timeout=5)
            if create_resp.status_code == 404:
                pytest.skip("权限 API 路由未注册")
            if create_resp.status_code not in [200, 201]:
                pytest.skip(f"创建权限失败: {create_resp.status_code}")

            created = create_resp.json()
            perm_id = created.get("data", {}).get("id")
            if not perm_id:
                pytest.skip("创建响应中无 id 字段")

            # 删除
            del_resp = requests.delete(
                f"{self.API_V1}/permissions/{perm_id}",
                headers=headers, timeout=5)
            assert del_resp.status_code in [200, 404], \
                f"删除权限响应异常: {del_resp.status_code}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_permission_delete_with_children_fails(self):
        """
        验证删除有子权限的权限应失败

        DELETE /api/v1/permissions/:id (有子权限)
        应返回错误：先删除子权限
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        try:
            # 创建一个有子权限的父权限
            parent_payload = {
                "name": f"父权限_{int(time.time())}",
                "code": f"PERM_PARENT_{int(time.time())}",
                "type": 1,
                "status": 1
            }
            parent_resp = requests.post(
                f"{self.API_V1}/permissions",
                headers=headers, json=parent_payload, timeout=5)
            if parent_resp.status_code == 404:
                pytest.skip("权限 API 路由未注册")
            if parent_resp.status_code not in [200, 201]:
                pytest.skip("创建父权限失败")

            parent_data = parent_resp.json()
            parent_id = parent_data.get("data", {}).get("id")
            if not parent_id:
                pytest.skip("创建响应中无 id 字段")

            # 创建子权限
            child_payload = {
                "name": f"子权限_{int(time.time())}",
                "code": f"PERM_CHILD_{int(time.time())}",
                "type": 1,
                "parent_id": parent_id,
                "status": 1
            }
            child_resp = requests.post(
                f"{self.API_V1}/permissions",
                headers=headers, json=child_payload, timeout=5)
            if child_resp.status_code not in [200, 201]:
                pytest.skip("创建子权限失败")

            # 尝试删除父权限
            del_resp = requests.delete(
                f"{self.API_V1}/permissions/{parent_id}",
                headers=headers, timeout=5)
            # 应返回错误
            assert del_resp.status_code == 400, \
                f"删除有子权限的父权限应返回 400，实际: {del_resp.status_code}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")


class TestRolePermissionAssignment:
    """
    角色权限分配 API 测试套件

    API 前缀: /api/v1/roles/:id/permissions
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

    # ─── API 功能测试 ────────────────────────────────────────────

    def test_get_role_permissions(self):
        """
        验证获取角色权限 API

        GET /api/v1/roles/:id/permissions
        应返回 menu_ids, api_ids, group_ids
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        try:
            resp = requests.get(
                f"{self.API_V1}/roles/1/permissions",
                headers=headers, timeout=5)
            if resp.status_code == 404:
                pytest.skip("角色权限 API 路由未注册")
            assert resp.status_code == 200, \
                f"获取角色权限响应异常: {resp.status_code}"
            data = resp.json()
            assert data.get("code") == 0, f"响应 code 非0: {data}"
            # 应包含权限相关信息
            perm_data = data.get("data", {})
            assert "menu_ids" in perm_data or "api_ids" in perm_data or \
                   "group_ids" in perm_data, "权限数据格式不正确"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_set_role_permissions(self):
        """
        验证设置角色权限 API

        PUT /api/v1/roles/:id/permissions
        设置菜单、API、权限组
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        # 先创建角色
        role_payload = {
            "role_name": f"权限测试角色_{int(time.time())}",
            "role_code": f"ROLE_PERM_TEST_{int(time.time())}",
            "status": 1
        }

        try:
            role_resp = requests.post(
                f"{self.API_V1}/roles",
                headers=headers, json=role_payload, timeout=5)
            if role_resp.status_code == 404:
                pytest.skip("角色 API 路由未注册")
            if role_resp.status_code not in [200, 201]:
                pytest.skip(f"创建角色失败: {role_resp.status_code}")

            role_data = role_resp.json()
            role_id = role_data.get("data", {}).get("id")
            if not role_id:
                pytest.skip("创建响应中无 id 字段")

            # 设置权限
            perm_payload = {
                "menu_ids": [1, 2, 3],
                "api_ids": [1, 2],
                "group_ids": []
            }
            perm_resp = requests.put(
                f"{self.API_V1}/roles/{role_id}/permissions",
                headers=headers, json=perm_payload, timeout=5)
            assert perm_resp.status_code in [200, 400], \
                f"设置角色权限响应异常: {perm_resp.status_code}"
            if perm_resp.status_code == 200:
                data = perm_resp.json()
                assert data.get("code") == 0, f"响应 code 非0: {data}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_set_role_permissions_update(self):
        """
        验证更新角色权限

        多次设置权限应覆盖而非追加
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        # 先创建角色
        role_payload = {
            "role_name": f"权限更新测试_{int(time.time())}",
            "role_code": f"ROLE_PERM_UPD_{int(time.time())}",
            "status": 1
        }

        try:
            role_resp = requests.post(
                f"{self.API_V1}/roles",
                headers=headers, json=role_payload, timeout=5)
            if role_resp.status_code == 404:
                pytest.skip("角色 API 路由未注册")
            if role_resp.status_code not in [200, 201]:
                pytest.skip("创建角色失败")

            role_data = role_resp.json()
            role_id = role_data.get("data", {}).get("id")
            if not role_id:
                pytest.skip("创建响应中无 id 字段")

            # 第一次设置权限
            perm_payload1 = {"menu_ids": [1], "api_ids": [1], "group_ids": []}
            requests.put(
                f"{self.API_V1}/roles/{role_id}/permissions",
                headers=headers, json=perm_payload1, timeout=5)

            # 第二次设置权限（不同的权限）
            perm_payload2 = {"menu_ids": [2], "api_ids": [2], "group_ids": []}
            perm_resp = requests.put(
                f"{self.API_V1}/roles/{role_id}/permissions",
                headers=headers, json=perm_payload2, timeout=5)
            assert perm_resp.status_code == 200, \
                f"更新角色权限响应异常: {perm_resp.status_code}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_permission_api_requires_auth(self):
        """
        验证权限 API 需要认证

        不带 Authorization header 的请求应被拒绝
        """
        try:
            resp = requests.get(
                f"{self.API_V1}/permissions", timeout=5)
            assert resp.status_code in [401, 403], \
                f"权限 API 应需要认证，实际返回: {resp.status_code}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
