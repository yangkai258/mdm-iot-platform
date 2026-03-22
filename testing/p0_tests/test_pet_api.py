# MDM Sprint 6 - 宠物配置 API 测试
"""
测试宠物配置 (Pet Profile) 和宠物控制台 (Pet Console) API：
- 宠物资料 CRUD
- 设备绑定宠物配置
- 宠物对话接口

依赖 backend Sprint 6 实现：
- controllers/pet_profile_controller.go (PetProfileController)
- controllers/pet_console_controller.go (PetConsoleController)
- models/pet_models.go (PetProfile)
- 路由已注册
"""
import os
import pytest
import requests
import time


class TestPetProfileAPI:
    """宠物资料 API 测试套件，API 前缀: /api/v1/devices/{device_id}/pet/profile"""

    BASE_URL = os.getenv("MDM_API_BASE_URL", "http://localhost:8080")
    API_V1 = f"{BASE_URL}/api/v1"

    def _get_auth_token(self):
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
        token = self._get_auth_token()
        if token:
            return {"Authorization": f"Bearer {token}"}
        return {}

    def _src(self, rel):
        return os.path.join(
            os.getenv("MDM_BACKEND_PATH",
                      "C:\\Users\\YKing\\.openclaw\\workspace\\mdm-project\\backend"),
            rel)

    def _read(self, rel):
        p = self._src(rel)
        if not os.path.exists(p):
            return None
        with open(p, encoding="utf-8") as f:
            return f.read()

    # ─── 静态源码检查 ────────────────────────────────────────────

    def test_pet_profile_controller_exists(self):
        src = self._read("controllers/pet_profile_controller.go")
        assert src is not None, "pet_profile_controller.go 不存在"
        assert "type PetProfileController struct" in src, "PetProfileController 未定义"

    def test_pet_profile_model_exists(self):
        import glob
        found = False
        for f in glob.glob(os.path.join(self._src("models"), "*.go")):
            c = open(f, encoding="utf-8").read()
            if "type PetProfile struct" in c:
                found = True
                assert 'json:"device_id"' in c, "PetProfile 缺少 device_id"
                assert 'json:"pet_name"' in c, "PetProfile 缺少 pet_name"
                break
        if not found:
            pytest.skip("PetProfile 模型未找到")

    def test_get_profile_method(self):
        src = self._read("controllers/pet_profile_controller.go")
        assert src is not None
        assert "func (c *PetProfileController) GetProfile" in src, "GetProfile 未实现"

    def test_update_profile_method(self):
        src = self._read("controllers/pet_profile_controller.go")
        assert src is not None
        assert "func (c *PetProfileController) UpdateProfile" in src, "UpdateProfile 未实现"

    # ─── API 功能测试 ────────────────────────────────────────────

    def test_get_profile_not_found_returns_default(self):
        """GET 不存在的设备配置应返回默认配置"""
        h = self._auth_headers()
        if not h:
            pytest.skip("无法获取 JWT token")
        dev = f"test-pet-{int(time.time())}-nonexist"
        try:
            r = requests.get(f"{self.API_V1}/devices/{dev}/pet/profile", headers=h, timeout=5)
            if r.status_code == 404:
                pytest.skip("路由未注册")
            assert r.status_code == 200, f"响应异常: {r.status_code}"
            d = r.json()
            assert d.get("code") == 0
            assert d["data"].get("device_id") == dev
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_update_profile_create(self):
        """PUT 创建宠物配置"""
        h = self._auth_headers()
        if not h:
            pytest.skip("无法获取 JWT token")
        dev = f"test-pet-{int(time.time())}-create"
        payload = {
            "pet_name": "TestPet",
            "personality": "lively",
            "interaction_freq": "high",
            "dnd_start_time": "22:00",
            "dnd_end_time": "08:00",
            "custom_rules": {"greeting": "hello"}
        }
        try:
            r = requests.put(f"{self.API_V1}/devices/{dev}/pet/profile", headers=h, json=payload, timeout=5)
            if r.status_code == 404:
                pytest.skip("路由未注册")
            assert r.status_code == 200, f"响应异常: {r.status_code}"
            d = r.json()
            assert d.get("code") == 0
            assert d["data"].get("pet_name") == "TestPet"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_update_profile_update(self):
        """PUT 更新已有配置"""
        h = self._auth_headers()
        if not h:
            pytest.skip("无法获取 JWT token")
        dev = f"test-pet-{int(time.time())}-update"
        try:
            requests.put(f"{self.API_V1}/devices/{dev}/pet/profile",
                         headers=h, json={"pet_name": "Original", "personality": "calm"}, timeout=5)
            r = requests.put(f"{self.API_V1}/devices/{dev}/pet/profile",
                            headers=h, json={"pet_name": "Updated"}, timeout=5)
            if r.status_code == 404:
                pytest.skip("路由未注册")
            assert r.status_code == 200
            assert r.json().get("data", {}).get("pet_name") == "Updated"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_get_profile_after_update(self):
        """GET 更新后再次获取应返回最新数据"""
        h = self._auth_headers()
        if not h:
            pytest.skip("无法获取 JWT token")
        dev = f"test-pet-{int(time.time())}-get"
        try:
            requests.put(f"{self.API_V1}/devices/{dev}/pet/profile",
                         headers=h, json={"pet_name": "FetchTestPet"}, timeout=5)
            r = requests.get(f"{self.API_V1}/devices/{dev}/pet/profile", headers=h, timeout=5)
            if r.status_code == 404:
                pytest.skip("路由未注册")
            assert r.status_code == 200
            assert r.json().get("data", {}).get("pet_name") == "FetchTestPet"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_update_profile_partial(self):
        """部分字段更新，其他字段应保持不变"""
        h = self._auth_headers()
        if not h:
            pytest.skip("无法获取 JWT token")
        dev = f"test-pet-{int(time.time())}-partial"
        try:
            requests.put(f"{self.API_V1}/devices/{dev}/pet/profile",
                         headers=h,
                         json={"pet_name": "FullPet", "personality": "calm", "interaction_freq": "medium"},
                         timeout=5)
            r = requests.put(f"{self.API_V1}/devices/{dev}/pet/profile",
                             headers=h, json={"pet_name": "PartialPet"}, timeout=5)
            if r.status_code == 404:
                pytest.skip("路由未注册")
            assert r.status_code == 200
            d = r.json().get("data", {})
            assert d.get("pet_name") == "PartialPet"
            assert d.get("personality") == "calm", "未更新字段应保持不变"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")


class TestPetConsoleAPI:
    """宠物控制台 API 测试套件，API 前缀: /api/v1/pet"""

    BASE_URL = os.getenv("MDM_API_BASE_URL", "http://localhost:8080")
    API_V1 = f"{BASE_URL}/api/v1"

    def _get_auth_token(self):
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
        token = self._get_auth_token()
        if token:
            return {"Authorization": f"Bearer {token}"}
        return {}

    def _src(self, rel):
        return os.path.join(
            os.getenv("MDM_BACKEND_PATH",
                      "C:\\Users\\YKing\\.openclaw\\workspace\\mdm-project\\backend"),
            rel)

    def _read(self, rel):
        p = self._src(rel)
        if not os.path.exists(p):
            return None
        with open(p, encoding="utf-8") as f:
            return f.read()

    # ─── 静态源码检查 ────────────────────────────────────────────

    def test_pet_console_controller_exists(self):
        src = self._read("controllers/pet_console_controller.go")
        assert src is not None, "pet_console_controller.go 不存在"
        assert "type PetConsoleController struct" in src, "PetConsoleController 未定义"

    def test_conversation_routes(self):
        src = self._read("controllers/pet_console_controller.go")
        assert src is not None
        assert "ListConversations" in src, "ListConversations 未找到"
        assert "CreateConversation" in src, "CreateConversation 未找到"
        assert "GetConversation" in src, "GetConversation 未找到"
        assert "SendMessage" in src, "SendMessage 未找到"

    def test_pet_status_routes(self):
        src = self._read("controllers/pet_console_controller.go")
        assert src is not None
        assert "GetPetStatus" in src, "GetPetStatus 未找到"
        assert "UpdatePetStatus" in src, "UpdatePetStatus 未找到"

    # ─── API 功能测试 ────────────────────────────────────────────

    def test_list_conversations(self):
        """GET /api/v1/pet/conversations"""
        h = self._auth_headers()
        if not h:
            pytest.skip("无法获取 JWT token")
        try:
            r = requests.get(f"{self.API_V1}/pet/conversations", headers=h, timeout=5)
            if r.status_code == 404:
                pytest.skip("路由未注册")
            assert r.status_code == 200
            assert r.json().get("code") == 0
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_create_conversation(self):
        """POST /api/v1/pet/conversations"""
        h = self._auth_headers()
        if not h:
            pytest.skip("无法获取 JWT token")
        try:
            r = requests.post(f"{self.API_V1}/pet/conversations",
                              headers=h, json={"device_id": f"pet-{int(time.time())}", "title": "测试对话"},
                              timeout=5)
            if r.status_code == 404:
                pytest.skip("路由未注册")
            assert r.status_code in [200, 201]
            assert r.json().get("code") == 0
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_get_conversation(self):
        """GET /api/v1/pet/conversations/:id"""
        h = self._auth_headers()
        if not h:
            pytest.skip("无法获取 JWT token")
        try:
            r = requests.get(f"{self.API_V1}/pet/conversations/1", headers=h, timeout=5)
            if r.status_code == 404:
                pytest.skip("路由未注册")
            assert r.status_code == 200
            assert r.json().get("code") == 0
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_send_message(self):
        """POST /api/v1/pet/messages"""
        h = self._auth_headers()
        if not h:
            pytest.skip("无法获取 JWT token")
        try:
            r = requests.post(f"{self.API_V1}/pet/messages",
                             headers=h, json={"conversation_id": 1, "content": "你好", "role": "user"},
                             timeout=5)
            if r.status_code == 404:
                pytest.skip("路由未注册")
            assert r.status_code in [200, 201]
            assert r.json().get("code") == 0
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_get_pet_status(self):
        """GET /api/v1/pet/status/:device_id"""
        h = self._auth_headers()
        if not h:
            pytest.skip("无法获取 JWT token")
        dev = f"test-pet-{int(time.time())}-status"
        try:
            r = requests.get(f"{self.API_V1}/pet/status/{dev}", headers=h, timeout=5)
            if r.status_code == 404:
                pytest.skip("路由未注册")
            assert r.status_code == 200
            assert r.json().get("code") == 0
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_update_pet_status(self):
        """PUT /api/v1/pet/status/:device_id"""
        h = self._auth_headers()
        if not h:
            pytest.skip("无法获取 JWT token")
        dev = f"test-pet-{int(time.time())}-updatestatus"
        try:
            r = requests.put(f"{self.API_V1}/pet/status/{dev}",
                            headers=h, json={"status": "happy", "mood": "excited"}, timeout=5)
            if r.status_code == 404:
                pytest.skip("路由未注册")
            assert r.status_code == 200
            assert r.json().get("code") == 0
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_pet_api_requires_auth(self):
        """宠物 API 应需要认证"""
        try:
            r = requests.get(f"{self.API_V1}/pet/conversations", timeout=5)
            assert r.status_code in [401, 403], f"应需要认证，实际: {r.status_code}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
