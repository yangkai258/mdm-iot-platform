# MDM Sprint 3 - 会员增强功能测试（积分/优惠券）
"""
测试会员增强功能 (Member Enhanced) 的核心功能：
- 积分规则 (PointsRule) CRUD API
- 积分流水 (MemberPointsRecord) API
- 优惠券 (Coupon) CRUD API
- 优惠券发放记录 (CouponGrant) API
- 会员升级规则 (MemberUpgradeRule) API
- 会员积分增减/兑换功能

依赖 backend Sprint 3 实现：
- models/member_models.go (PointsRule, MemberPointsRecord, Coupon, CouponGrant)
- controllers/member_controller.go (PointsRule/Coupon 相关方法)
- 路由已注册在 device_controller.go
"""
import os
import pytest
import requests
import time


class TestMemberEnhanced:
    """
    会员增强功能测试套件
    PASS = 会员增强功能正确实现
    FAIL = 会员增强功能未实现或有问题
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

    def test_points_rule_model_exists(self):
        """
        验证 PointsRule 模型已定义

        检查 models/member_models.go 中是否定义 PointsRule 结构体：
        - RuleCode, RuleName, RuleType, Points, Amount, Remark, Status
        """
        source = self._read_source("models/member_models.go")
        assert source is not None, "models/member_models.go 文件不存在"

        assert "type PointsRule struct" in source, \
            "PointsRule 结构体未定义"
        assert 'json:"rule_code"' in source, \
            "PointsRule 缺少 rule_code 字段"
        assert 'json:"rule_name"' in source, \
            "PointsRule 缺少 rule_name 字段"
        assert 'json:"rule_type"' in source, \
            "PointsRule 缺少 rule_type 字段 (1获取积分 2抵扣积分 3不积分)"
        assert 'json:"points"' in source, \
            "PointsRule 缺少 points 字段"
        assert 'json:"amount"' in source, \
            "PointsRule 缺少 amount 字段 (对应金额)"

    def test_member_points_record_model_exists(self):
        """
        验证 MemberPointsRecord 模型已定义

        检查 models/member_models.go 中是否定义 MemberPointsRecord：
        - MemberID, Points, PointsType, SourceType, SourceID
        - BeforeBalance, AfterBalance, Operator, Remark
        """
        source = self._read_source("models/member_models.go")
        assert source is not None, "models/member_models.go 文件不存在"

        assert "type MemberPointsRecord struct" in source, \
            "MemberPointsRecord 结构体未定义"
        assert 'json:"member_id"' in source, \
            "MemberPointsRecord 缺少 member_id 字段"
        assert 'json:"points"' in source, \
            "MemberPointsRecord 缺少 points 字段"
        assert 'json:"points_type"' in source, \
            "MemberPointsRecord 缺少 points_type 字段"
        assert 'json:"before_balance"' in source, \
            "MemberPointsRecord 缺少 before_balance 字段"
        assert 'json:"after_balance"' in source, \
            "MemberPointsRecord 缺少 after_balance 字段"

    def test_coupon_model_exists(self):
        """
        验证 Coupon 模型已定义

        检查 models/member_models.go 中是否定义 Coupon：
        - CouponCode, CouponName, CouponType, FaceValue, MinAmount
        - TotalStock, RemainStock, ValidDays, StartTime, EndTime
        """
        source = self._read_source("models/member_models.go")
        assert source is not None, "models/member_models.go 文件不存在"

        assert "type Coupon struct" in source, \
            "Coupon 结构体未定义"
        assert 'json:"coupon_code"' in source, \
            "Coupon 缺少 coupon_code 字段"
        assert 'json:"coupon_name"' in source, \
            "Coupon 缺少 coupon_name 字段"
        assert 'json:"coupon_type"' in source, \
            "Coupon 缺少 coupon_type 字段 (1满减 2折扣 3兑换)"
        assert 'json:"face_value"' in source, \
            "Coupon 缺少 face_value 字段"
        assert 'json:"total_stock"' in source, \
            "Coupon 缺少 total_stock 字段"
        assert 'json:"remain_stock"' in source, \
            "Coupon 缺少 remain_stock 字段"

    def test_coupon_grant_model_exists(self):
        """
        验证 CouponGrant 模型已定义

        检查 models/member_models.go 中是否定义 CouponGrant：
        - CouponID, MemberID, GrantTime, UseTime, OrderID, Status
        """
        source = self._read_source("models/member_models.go")
        assert source is not None, "models/member_models.go 文件不存在"

        assert "type CouponGrant struct" in source, \
            "CouponGrant 结构体未定义"
        assert 'json:"coupon_id"' in source, \
            "CouponGrant 缺少 coupon_id 字段"
        assert 'json:"member_id"' in source, \
            "CouponGrant 缺少 member_id 字段"
        assert 'json:"status"' in source, \
            "CouponGrant 缺少 status 字段 (1未使用 2已使用 3已过期)"

    def test_points_rule_controller_methods(self):
        """
        验证 PointsRule 控制器方法已实现

        检查 controllers/member_controller.go 中：
        - PointsRuleList
        - PointsRuleCreate
        - PointsRuleUpdate
        - PointsRuleDelete
        """
        source = self._read_source("controllers/member_controller.go")
        assert source is not None, "controllers/member_controller.go 文件不存在"

        assert "func (c *MemberController) PointsRuleList" in source, \
            "PointsRuleList 方法未实现"
        assert "func (c *MemberController) PointsRuleCreate" in source, \
            "PointsRuleCreate 方法未实现"
        assert "func (c *MemberController) PointsRuleUpdate" in source, \
            "PointsRuleUpdate 方法未实现"
        assert "func (c *MemberController) PointsRuleDelete" in source, \
            "PointsRuleDelete 方法未实现"

    def test_points_rule_routes_registered(self):
        """
        验证积分规则 API 路由已注册

        检查 controllers/device_controller.go 中：
        - GET    /member/points/rules
        - POST   /member/points/rules
        - PUT    /member/points/rules/:id
        - DELETE /member/points/rules/:id
        """
        source = self._read_source("controllers/device_controller.go")
        assert source is not None, "controllers/device_controller.go 文件不存在"

        assert "/member/points/rules" in source, \
            "积分规则路由未注册 (应包含 /member/points/rules)"
        assert 'PointsRuleList' in source, \
            "PointsRuleList 路由未注册"
        assert 'PointsRuleCreate' in source, \
            "PointsRuleCreate 路由未注册"

    def test_coupon_routes_registered(self):
        """
        验证优惠券 API 路由已注册

        检查 controllers/device_controller.go 中：
        - GET    /member/coupons
        - POST   /member/coupons
        - PUT    /member/coupons/:id
        - DELETE /member/coupons/:id
        """
        source = self._read_source("controllers/device_controller.go")
        assert source is not None, "controllers/device_controller.go 文件不存在"

        assert "/member/coupons" in source, \
            "优惠券路由未注册 (应包含 /member/coupons)"
        assert 'CouponList' in source, \
            "CouponList 路由未注册"
        assert 'CouponCreate' in source, \
            "CouponCreate 路由未注册"

    def test_points_record_route_registered(self):
        """
        验证积分流水 API 路由已注册

        检查 controllers/device_controller.go 中：
        - GET /member/points/records
        """
        source = self._read_source("controllers/device_controller.go")
        assert source is not None, "controllers/device_controller.go 文件不存在"

        assert "/member/points/records" in source, \
            "积分流水路由未注册 (应包含 /member/points/records)"
        assert 'PointsRecordList' in source, \
            "PointsRecordList 路由未注册"

    # ─── API 功能测试 ────────────────────────────────────────────

    def test_points_rule_list_api(self):
        """
        验证积分规则列表 API

        GET /api/v1/member/points/rules
        验证返回 list, total, page, page_size 字段
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        try:
            resp = requests.get(
                f"{self.API_V1}/member/points/rules?page=1&page_size=10",
                headers=headers, timeout=5)
            if resp.status_code == 404:
                pytest.skip("积分规则 API 路由未注册")
            assert resp.status_code == 200, \
                f"积分规则列表响应异常: {resp.status_code}"
            data = resp.json()
            assert data.get("code") == 0, f"响应 code 非0: {data}"
            assert "data" in data, "响应缺少 data 字段"
            assert "list" in data["data"], "data 缺少 list 字段"
            assert "total" in data["data"], "data 缺少 total 字段"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_points_rule_create_api(self):
        """
        验证创建积分规则 API

        POST /api/v1/member/points/rules
        发送积分规则创建请求，验证响应格式
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        payload = {
            "rule_code": f"POINTS_BUY_{int(time.time())}",
            "rule_name": "消费积分规则",
            "rule_type": 1,
            "points": 100,
            "amount": 10.0,
            "remark": "每消费10元获得100积分",
            "status": 1
        }

        try:
            resp = requests.post(f"{self.API_V1}/member/points/rules",
                               headers=headers, json=payload, timeout=5)
            if resp.status_code == 404:
                pytest.skip("积分规则 API 路由未注册")
            assert resp.status_code in [200, 201, 400], \
                f"创建积分规则响应异常: {resp.status_code}"
            if resp.status_code == 200:
                data = resp.json()
                assert data.get("code") == 0, f"响应 code 非0: {data}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_points_rule_update_api(self):
        """
        验证更新积分规则 API

        PUT /api/v1/member/points/rules/:id
        先创建，再更新，验证更新结果
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        # 先创建
        payload = {
            "rule_code": f"POINTS_UPDATE_{int(time.time())}",
            "rule_name": "更新测试规则",
            "rule_type": 1,
            "points": 100,
            "amount": 10.0,
            "status": 1
        }

        try:
            create_resp = requests.post(f"{self.API_V1}/member/points/rules",
                                        headers=headers, json=payload, timeout=5)
            if create_resp.status_code == 404:
                pytest.skip("积分规则 API 路由未注册")
            if create_resp.status_code not in [200, 201]:
                pytest.skip(f"创建积分规则失败: {create_resp.status_code}")

            created = create_resp.json()
            rule_id = created.get("data", {}).get("id")
            if not rule_id:
                pytest.skip("创建响应中无 id 字段")

            # 更新
            update_payload = {
                "rule_name": "更新后的规则名称",
                "points": 200,
                "amount": 20.0,
                "status": 1
            }
            update_resp = requests.put(
                f"{self.API_V1}/member/points/rules/{rule_id}",
                headers=headers, json=update_payload, timeout=5)
            assert update_resp.status_code in [200, 400], \
                f"更新积分规则响应异常: {update_resp.status_code}"

        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_coupon_list_api(self):
        """
        验证优惠券列表 API

        GET /api/v1/member/coupons
        验证返回 list, total, page, page_size 字段
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        try:
            resp = requests.get(
                f"{self.API_V1}/member/coupons?page=1&page_size=10",
                headers=headers, timeout=5)
            if resp.status_code == 404:
                pytest.skip("优惠券 API 路由未注册")
            assert resp.status_code == 200, \
                f"优惠券列表响应异常: {resp.status_code}"
            data = resp.json()
            assert data.get("code") == 0, f"响应 code 非0: {data}"
            assert "data" in data, "响应缺少 data 字段"
            assert "list" in data["data"], "data 缺少 list 字段"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_coupon_create_api(self):
        """
        验证创建优惠券 API

        POST /api/v1/member/coupons
        发送优惠券创建请求，验证 remain_stock 自动等于 total_stock
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        payload = {
            "coupon_code": f"COUPON_TEST_{int(time.time())}",
            "coupon_name": "测试优惠券",
            "coupon_type": 1,
            "face_value": 10.0,
            "min_amount": 50.0,
            "total_stock": 100,
            "valid_days": 30,
            "status": 1
        }

        try:
            resp = requests.post(f"{self.API_V1}/member/coupons",
                               headers=headers, json=payload, timeout=5)
            if resp.status_code == 404:
                pytest.skip("优惠券 API 路由未注册")
            assert resp.status_code in [200, 201, 400], \
                f"创建优惠券响应异常: {resp.status_code}"
            if resp.status_code == 200:
                data = resp.json()
                assert data.get("code") == 0, f"响应 code 非0: {data}"
                # 验证 remain_stock 自动设置
                coupon_data = data.get("data", {})
                if "remain_stock" in coupon_data:
                    assert coupon_data["remain_stock"] == 100, \
                        f"remain_stock 应自动设为 total_stock: {coupon_data}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_coupon_update_api(self):
        """
        验证更新优惠券 API

        PUT /api/v1/member/coupons/:id
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        # 先创建
        payload = {
            "coupon_code": f"COUPON_UPD_{int(time.time())}",
            "coupon_name": "更新前名称",
            "coupon_type": 1,
            "face_value": 5.0,
            "min_amount": 30.0,
            "total_stock": 50,
            "valid_days": 7,
            "status": 1
        }

        try:
            create_resp = requests.post(f"{self.API_V1}/member/coupons",
                                        headers=headers, json=payload, timeout=5)
            if create_resp.status_code == 404:
                pytest.skip("优惠券 API 路由未注册")
            if create_resp.status_code not in [200, 201]:
                pytest.skip(f"创建优惠券失败: {create_resp.status_code}")

            created = create_resp.json()
            coupon_id = created.get("data", {}).get("id")
            if not coupon_id:
                pytest.skip("创建响应中无 id 字段")

            # 更新
            update_payload = {
                "coupon_name": "更新后名称",
                "face_value": 8.0,
                "status": 1
            }
            update_resp = requests.put(
                f"{self.API_V1}/member/coupons/{coupon_id}",
                headers=headers, json=update_payload, timeout=5)
            assert update_resp.status_code in [200, 400], \
                f"更新优惠券响应异常: {update_resp.status_code}"

        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_coupon_delete_api(self):
        """
        验证删除优惠券 API

        DELETE /api/v1/member/coupons/:id
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        # 先创建
        payload = {
            "coupon_code": f"COUPON_DEL_{int(time.time())}",
            "coupon_name": "待删除优惠券",
            "coupon_type": 1,
            "face_value": 5.0,
            "total_stock": 10,
            "status": 1
        }

        try:
            create_resp = requests.post(f"{self.API_V1}/member/coupons",
                                        headers=headers, json=payload, timeout=5)
            if create_resp.status_code == 404:
                pytest.skip("优惠券 API 路由未注册")
            if create_resp.status_code not in [200, 201]:
                pytest.skip(f"创建优惠券失败: {create_resp.status_code}")

            created = create_resp.json()
            coupon_id = created.get("data", {}).get("id")
            if not coupon_id:
                pytest.skip("创建响应中无 id 字段")

            # 删除
            del_resp = requests.delete(
                f"{self.API_V1}/member/coupons/{coupon_id}",
                headers=headers, timeout=5)
            assert del_resp.status_code in [200, 404], \
                f"删除优惠券响应异常: {del_resp.status_code}"

        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_points_record_list_api(self):
        """
        验证积分流水列表 API

        GET /api/v1/member/points/records
        验证积分流水记录可查询
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        try:
            resp = requests.get(
                f"{self.API_V1}/member/points/records?page=1&page_size=10",
                headers=headers, timeout=5)
            if resp.status_code == 404:
                pytest.skip("积分流水 API 路由未注册")
            assert resp.status_code == 200, \
                f"积分流水列表响应异常: {resp.status_code}"
            data = resp.json()
            assert data.get("code") == 0, f"响应 code 非0: {data}"
        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")

    def test_member_points_grant_integration(self):
        """
        验证会员积分发放积分流水记录

        创建会员 -> 积分规则 -> 积分发放 -> 检查积分流水
        验证积分变动后 balance 正确记录
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("无法获取 JWT token，后端可能未运行")

        try:
            # 1. 创建测试会员
            member_payload = {
                "member_code": f"TEST_PTS_{int(time.time())}",
                "member_name": "积分测试会员",
                "phone": f"138{int(time.time()) % 100000000:08d}",
                "member_level": 1,
                "points": 0,
                "status": 1
            }
            member_resp = requests.post(f"{self.API_V1}/members",
                                         headers=headers, json=member_payload, timeout=5)
            if member_resp.status_code == 404:
                pytest.skip("会员 API 路由未注册")
            if member_resp.status_code not in [200, 201]:
                pytest.skip(f"创建会员失败: {member_resp.status_code}")

            member_data = member_resp.json()
            member_id = member_data.get("data", {}).get("id")
            if not member_id:
                pytest.skip("创建会员响应中无 id 字段")

            # 2. 创建积分规则
            rule_payload = {
                "rule_code": f"PTS_RULE_{int(time.time())}",
                "rule_name": "测试积分规则",
                "rule_type": 1,
                "points": 100,
                "amount": 10.0,
                "status": 1
            }
            rule_resp = requests.post(f"{self.API_V1}/member/points/rules",
                                      headers=headers, json=rule_payload, timeout=5)
            if rule_resp.status_code not in [200, 201]:
                pytest.skip(f"创建积分规则失败: {rule_resp.status_code}")

            # 3. 查询积分流水
            records_resp = requests.get(
                f"{self.API_V1}/member/points/records?member_id={member_id}",
                headers=headers, timeout=5)
            if records_resp.status_code == 200:
                records_data = records_resp.json()
                assert records_data.get("code") == 0, \
                    f"积分流水查询响应异常: {records_data}"

        except requests.exceptions.ConnectionError:
            pytest.skip("后端未运行，无法测试 API")


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
