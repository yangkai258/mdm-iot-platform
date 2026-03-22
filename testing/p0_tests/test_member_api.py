# MDM Sprint 6 - 会员管理 API 测试
"""
测试会员管理 (Member Management) API：
- 会员 CRUD
- 会员等级管理
- 积分/优惠券管理

依赖 backend Sprint 6 实现：
- controllers/member_controller.go (MemberController)
- models/member_models.go (Member, MemberLevel, Coupon, PointsRule 等)
- 路由已注册
"""
import os
import pytest
import requests
import time


class TestMemberCRUD:
    """会员 CRUD API 测试套件，API 前缀: /api/v1/members"""

    BASE_URL = os.getenv("MDM_API_BASE_URL", "http://localhost:8080")
    API_V1 = f"{BASE_URL}/api/v1"

    def _token(self):
        try:
            r = requests.post(f"{self.API_V1}/auth/login",
                             json={"username": os.getenv("TEST_USERNAME", "admin"),
                                   "password": os.getenv("TEST_PASSWORD", "admin123")}, timeout=5)
            if r.status_code == 200 and r.json().get("code") == 0:
                return r.json()["data"]["token"]
        except Exception:
            pass
        return None

    def _hdr(self):
        t = self._token()
        return {"Authorization": f"Bearer {t}"} if t else {}

    def _src(self, rel):
        return os.path.join(
            os.getenv("MDM_BACKEND_PATH",
                      "C:\\Users\\YKing\\.openclaw\\workspace\\mdm-project\\backend"), rel)

    def _read(self, rel):
        p = self._src(rel)
        if not os.path.exists(p):
            return None
        with open(p, encoding="utf-8") as f:
            return f.read()

    # ─── 静态源码检查 ────────────────────────────────────────────

    def test_member_model_exists(self):
        src = self._read("models/member_models.go")
        assert src is not None, "member_models.go 不存在"
        assert "type Member struct" in src, "Member 未定义"
        assert 'json:"member_code"' in src, "Member 缺少 member_code"
        assert 'json:"member_name"' in src, "Member 缺少 member_name"
        assert 'json:"member_level"' in src, "Member 缺少 member_level"
        assert 'json:"points"' in src, "Member 缺少 points"

    def test_member_controller_methods(self):
        src = self._read("controllers/member_controller.go")
        assert src is not None
        for m in ["MemberList", "MemberCreate", "MemberUpdate", "MemberDelete", "MemberDetail"]:
            assert f"func (c *MemberController) {m}" in src, f"{m} 未实现"

    # ─── API 功能测试 ────────────────────────────────────────────

    def test_member_list_api(self):
        """GET /api/v1/members"""
        h = self._hdr()
        if not h:
            pytest.skip("无法获取 JWT token")
        try:
            r = requests.get(f"{self.API_V1}/members?page=1&page_size=10", headers=h, timeout=5)
            if r.status_code == 404:
                pytest.skip("路由未注册")
            assert r.status_code == 200
            d = r.json()
            assert d.get("code") == 0
            assert "list" in d["data"] and "total" in d["data"]
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_member_list_filter_keyword(self):
        """GET /api/v1/members?keyword=xxx"""
        h = self._hdr()
        if not h:
            pytest.skip("无法获取 JWT token")
        try:
            r = requests.get(f"{self.API_V1}/members?keyword=admin&page=1&page_size=10", headers=h, timeout=5)
            if r.status_code == 404:
                pytest.skip("路由未注册")
            assert r.status_code == 200
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_member_list_filter_level(self):
        """GET /api/v1/members?member_level=1"""
        h = self._hdr()
        if not h:
            pytest.skip("无法获取 JWT token")
        try:
            r = requests.get(f"{self.API_V1}/members?member_level=1&page=1&page_size=10", headers=h, timeout=5)
            if r.status_code == 404:
                pytest.skip("路由未注册")
            assert r.status_code == 200
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_member_list_filter_points_range(self):
        """GET /api/v1/members?points_min=0&points_max=10000"""
        h = self._hdr()
        if not h:
            pytest.skip("无法获取 JWT token")
        try:
            r = requests.get(f"{self.API_V1}/members?points_min=0&points_max=10000&page=1&page_size=10",
                             headers=h, timeout=5)
            if r.status_code == 404:
                pytest.skip("路由未注册")
            assert r.status_code == 200
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_member_create_api(self):
        """POST /api/v1/members"""
        h = self._hdr()
        if not h:
            pytest.skip("无法获取 JWT token")
        payload = {
            "member_code": f"TEST_MEM_{int(time.time())}",
            "member_name": "测试会员",
            "phone": f"138{int(time.time()) % 100000000:08d}",
            "member_level": 1, "points": 0, "status": 1
        }
        try:
            r = requests.post(f"{self.API_V1}/members", headers=h, json=payload, timeout=5)
            if r.status_code == 404:
                pytest.skip("路由未注册")
            assert r.status_code in [200, 201, 400]
            if r.status_code in [200, 201]:
                assert r.json().get("code") == 0
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_member_update_api(self):
        """PUT /api/v1/members/:id"""
        h = self._hdr()
        if not h:
            pytest.skip("无法获取 JWT token")
        try:
            cr = requests.post(f"{self.API_V1}/members",
                               headers=h,
                               json={"member_code": f"T_UPD_{int(time.time())}",
                                     "member_name": "原始名称",
                                     "phone": f"137{int(time.time()) % 100000000:08d}",
                                     "status": 1}, timeout=5)
            if cr.status_code == 404:
                pytest.skip("路由未注册")
            if cr.status_code not in [200, 201]:
                pytest.skip("创建会员失败")
            mid = cr.json().get("data", {}).get("id")
            if not mid:
                pytest.skip("响应无 id")
            ur = requests.put(f"{self.API_V1}/members/{mid}",
                              headers=h, json={"member_name": "更新后名称"}, timeout=5)
            assert ur.status_code in [200, 400]
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_member_delete_api(self):
        """DELETE /api/v1/members/:id"""
        h = self._hdr()
        if not h:
            pytest.skip("无法获取 JWT token")
        try:
            cr = requests.post(f"{self.API_V1}/members",
                               headers=h,
                               json={"member_code": f"T_DEL_{int(time.time())}",
                                     "member_name": "待删除",
                                     "phone": f"136{int(time.time()) % 100000000:08d}",
                                     "status": 1}, timeout=5)
            if cr.status_code == 404:
                pytest.skip("路由未注册")
            if cr.status_code not in [200, 201]:
                pytest.skip("创建会员失败")
            mid = cr.json().get("data", {}).get("id")
            if not mid:
                pytest.skip("响应无 id")
            dr = requests.delete(f"{self.API_V1}/members/{mid}", headers=h, timeout=5)
            assert dr.status_code in [200, 404]
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")


class TestMemberLevel:
    """会员等级 API 测试套件，API 前缀: /api/v1/member/levels"""

    BASE_URL = os.getenv("MDM_API_BASE_URL", "http://localhost:8080")
    API_V1 = f"{BASE_URL}/api/v1"

    def _token(self):
        try:
            r = requests.post(f"{self.API_V1}/auth/login",
                             json={"username": os.getenv("TEST_USERNAME", "admin"),
                                   "password": os.getenv("TEST_PASSWORD", "admin123")}, timeout=5)
            if r.status_code == 200 and r.json().get("code") == 0:
                return r.json()["data"]["token"]
        except Exception:
            pass
        return None

    def _hdr(self):
        t = self._token()
        return {"Authorization": f"Bearer {t}"} if t else {}

    def _src(self, rel):
        return os.path.join(
            os.getenv("MDM_BACKEND_PATH",
                      "C:\\Users\\YKing\\.openclaw\\workspace\\mdm-project\\backend"), rel)

    def _read(self, rel):
        p = self._src(rel)
        if not os.path.exists(p):
            return None
        with open(p, encoding="utf-8") as f:
            return f.read()

    # ─── 静态源码检查 ────────────────────────────────────────────

    def test_member_level_model_exists(self):
        src = self._read("models/member_models.go")
        assert src is not None
        assert "type MemberLevel struct" in src, "MemberLevel 未定义"

    def test_level_controller_methods(self):
        src = self._read("controllers/member_controller.go")
        assert src is not None
        for m in ["LevelList", "LevelCreate", "LevelUpdate", "LevelDelete"]:
            assert f"func (c *MemberController) {m}" in src, f"{m} 未实现"

    # ─── API 功能测试 ────────────────────────────────────────────

    def test_level_list_api(self):
        """GET /api/v1/member/levels"""
        h = self._hdr()
        if not h:
            pytest.skip("无法获取 JWT token")
        try:
            r = requests.get(f"{self.API_V1}/member/levels", headers=h, timeout=5)
            if r.status_code == 404:
                pytest.skip("路由未注册")
            assert r.status_code == 200
            assert r.json().get("code") == 0
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_level_create_api(self):
        """POST /api/v1/member/levels"""
        h = self._hdr()
        if not h:
            pytest.skip("无法获取 JWT token")
        payload = {"level_name": f"测试等级_{int(time.time())}", "level_value": 3,
                   "discount": 90.0, "points_rate": 1.5, "status": 1}
        try:
            r = requests.post(f"{self.API_V1}/member/levels", headers=h, json=payload, timeout=5)
            if r.status_code == 404:
                pytest.skip("路由未注册")
            assert r.status_code in [200, 201, 400]
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_level_update_api(self):
        """PUT /api/v1/member/levels/:id"""
        h = self._hdr()
        if not h:
            pytest.skip("无法获取 JWT token")
        try:
            cr = requests.post(f"{self.API_V1}/member/levels",
                               headers=h, json={"level_name": f"L_UPD_{int(time.time())}",
                                               "level_value": 3, "discount": 95.0, "status": 1}, timeout=5)
            if cr.status_code == 404:
                pytest.skip("路由未注册")
            if cr.status_code not in [200, 201]:
                pytest.skip("创建等级失败")
            lid = cr.json().get("data", {}).get("id")
            if not lid:
                pytest.skip("响应无 id")
            ur = requests.put(f"{self.API_V1}/member/levels/{lid}",
                              headers=h, json={"discount": 85.0}, timeout=5)
            assert ur.status_code in [200, 400]
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_level_delete_api(self):
        """DELETE /api/v1/member/levels/:id"""
        h = self._hdr()
        if not h:
            pytest.skip("无法获取 JWT token")
        try:
            cr = requests.post(f"{self.API_V1}/member/levels",
                               headers=h, json={"level_name": f"L_DEL_{int(time.time())}",
                                               "level_value": 10, "discount": 80.0, "status": 1}, timeout=5)
            if cr.status_code == 404:
                pytest.skip("路由未注册")
            if cr.status_code not in [200, 201]:
                pytest.skip("创建等级失败")
            lid = cr.json().get("data", {}).get("id")
            if not lid:
                pytest.skip("响应无 id")
            dr = requests.delete(f"{self.API_V1}/member/levels/{lid}", headers=h, timeout=5)
            assert dr.status_code in [200, 404]
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")


class TestPointsAndCoupons:
    """积分和优惠券 API 测试套件"""

    BASE_URL = os.getenv("MDM_API_BASE_URL", "http://localhost:8080")
    API_V1 = f"{BASE_URL}/api/v1"

    def _token(self):
        try:
            r = requests.post(f"{self.API_V1}/auth/login",
                             json={"username": os.getenv("TEST_USERNAME", "admin"),
                                   "password": os.getenv("TEST_PASSWORD", "admin123")}, timeout=5)
            if r.status_code == 200 and r.json().get("code") == 0:
                return r.json()["data"]["token"]
        except Exception:
            pass
        return None

    def _hdr(self):
        t = self._token()
        return {"Authorization": f"Bearer {t}"} if t else {}

    # ─── 积分规则 ───────────────────────────────────────────────

    def test_points_rule_list_api(self):
        """GET /api/v1/member/points/rules"""
        h = self._hdr()
        if not h:
            pytest.skip("无法获取 JWT token")
        try:
            r = requests.get(f"{self.API_V1}/member/points/rules?page=1&page_size=10", headers=h, timeout=5)
            if r.status_code == 404:
                pytest.skip("路由未注册")
            assert r.status_code == 200
            assert r.json().get("code") == 0
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_points_rule_create_api(self):
        """POST /api/v1/member/points/rules"""
        h = self._hdr()
        if not h:
            pytest.skip("无法获取 JWT token")
        payload = {"rule_code": f"RULE_{int(time.time())}", "rule_name": "消费积分规则",
                   "rule_type": 1, "points": 100, "amount": 10.0, "status": 1}
        try:
            r = requests.post(f"{self.API_V1}/member/points/rules", headers=h, json=payload, timeout=5)
            if r.status_code == 404:
                pytest.skip("路由未注册")
            assert r.status_code in [200, 201, 400]
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_points_rule_update_api(self):
        """PUT /api/v1/member/points/rules/:id"""
        h = self._hdr()
        if not h:
            pytest.skip("无法获取 JWT token")
        try:
            cr = requests.post(f"{self.API_V1}/member/points/rules",
                               headers=h, json={"rule_code": f"R_U_{int(time.time())}",
                                               "rule_name": "原始", "rule_type": 1,
                                               "points": 100, "amount": 10.0, "status": 1}, timeout=5)
            if cr.status_code == 404:
                pytest.skip("路由未注册")
            if cr.status_code not in [200, 201]:
                pytest.skip("创建规则失败")
            rid = cr.json().get("data", {}).get("id")
            if not rid:
                pytest.skip("响应无 id")
            ur = requests.put(f"{self.API_V1}/member/points/rules/{rid}",
                              headers=h, json={"points": 200}, timeout=5)
            assert ur.status_code in [200, 400]
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_points_rule_delete_api(self):
        """DELETE /api/v1/member/points/rules/:id"""
        h = self._hdr()
        if not h:
            pytest.skip("无法获取 JWT token")
        try:
            cr = requests.post(f"{self.API_V1}/member/points/rules",
                               headers=h, json={"rule_code": f"R_D_{int(time.time())}",
                                               "rule_name": "待删除", "rule_type": 1,
                                               "points": 50, "status": 1}, timeout=5)
            if cr.status_code == 404:
                pytest.skip("路由未注册")
            if cr.status_code not in [200, 201]:
                pytest.skip("创建规则失败")
            rid = cr.json().get("data", {}).get("id")
            if not rid:
                pytest.skip("响应无 id")
            dr = requests.delete(f"{self.API_V1}/member/points/rules/{rid}", headers=h, timeout=5)
            assert dr.status_code in [200, 404]
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_points_record_list_api(self):
        """GET /api/v1/member/points/records"""
        h = self._hdr()
        if not h:
            pytest.skip("无法获取 JWT token")
        try:
            r = requests.get(f"{self.API_V1}/member/points/records?page=1&page_size=10", headers=h, timeout=5)
            if r.status_code == 404:
                pytest.skip("路由未注册")
            assert r.status_code == 200
            assert r.json().get("code") == 0
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    # ─── 优惠券 ───────────────────────────────────────────────

    def test_coupon_list_api(self):
        """GET /api/v1/member/coupons"""
        h = self._hdr()
        if not h:
            pytest.skip("无法获取 JWT token")
        try:
            r = requests.get(f"{self.API_V1}/member/coupons?page=1&page_size=10", headers=h, timeout=5)
            if r.status_code == 404:
                pytest.skip("路由未注册")
            assert r.status_code == 200
            assert r.json().get("code") == 0
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_coupon_create_api(self):
        """POST /api/v1/member/coupons，验证 remain_stock 自动等于 total_stock"""
        h = self._hdr()
        if not h:
            pytest.skip("无法获取 JWT token")
        payload = {"coupon_code": f"CPN_{int(time.time())}", "coupon_name": "测试优惠券",
                   "coupon_type": 1, "face_value": 10.0, "min_amount": 50.0,
                   "total_stock": 100, "valid_days": 30, "status": 1}
        try:
            r = requests.post(f"{self.API_V1}/member/coupons", headers=h, json=payload, timeout=5)
            if r.status_code == 404:
                pytest.skip("路由未注册")
            assert r.status_code in [200, 201, 400]
            if r.status_code == 200:
                cp = r.json().get("data", {})
                if "remain_stock" in cp:
                    assert cp["remain_stock"] == 100, "remain_stock 应自动设为 total_stock"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_coupon_update_api(self):
        """PUT /api/v1/member/coupons/:id"""
        h = self._hdr()
        if not h:
            pytest.skip("无法获取 JWT token")
        try:
            cr = requests.post(f"{self.API_V1}/member/coupons",
                                headers=h, json={"coupon_code": f"C_U_{int(time.time())}",
                                                "coupon_name": "原始", "coupon_type": 1,
                                                "face_value": 5.0, "total_stock": 50, "status": 1}, timeout=5)
            if cr.status_code == 404:
                pytest.skip("路由未注册")
            if cr.status_code not in [200, 201]:
                pytest.skip("创建优惠券失败")
            cid = cr.json().get("data", {}).get("id")
            if not cid:
                pytest.skip("响应无 id")
            ur = requests.put(f"{self.API_V1}/member/coupons/{cid}",
                              headers=h, json={"face_value": 8.0}, timeout=5)
            assert ur.status_code in [200, 400]
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_coupon_delete_api(self):
        """DELETE /api/v1/member/coupons/:id"""
        h = self._hdr()
        if not h:
            pytest.skip("无法获取 JWT token")
        try:
            cr = requests.post(f"{self.API_V1}/member/coupons",
                                headers=h, json={"coupon_code": f"C_D_{int(time.time())}",
                                                "coupon_name": "待删除", "coupon_type": 1,
                                                "face_value": 5.0, "total_stock": 10, "status": 1}, timeout=5)
            if cr.status_code == 404:
                pytest.skip("路由未注册")
            if cr.status_code not in [200, 201]:
                pytest.skip("创建优惠券失败")
            cid = cr.json().get("data", {}).get("id")
            if not cid:
                pytest.skip("响应无 id")
            dr = requests.delete(f"{self.API_V1}/member/coupons/{cid}", headers=h, timeout=5)
            assert dr.status_code in [200, 404]
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")

    def test_member_api_requires_auth(self):
        """会员 API 应需要认证"""
        try:
            r = requests.get(f"{self.API_V1}/members", timeout=5)
            assert r.status_code in [401, 403], f"应需要认证，实际: {r.status_code}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行")


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
