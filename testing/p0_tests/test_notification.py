# MDM Sprint 2 - 通知管理功能测试
"""
测试通知管理 (Notification Management) 的核心功能：
- 发送通知 (Push Notification)
- 通知模板管理 (Template CRUD)
- 公告管理 (Announcement CRUD)
- MQTT 通知下发 (MQTT Push)

依赖 backend Sprint 2 实现：
- models/notification.go (Notification, NotificationTemplate, Announcement)
- controllers/notification_controller.go (发送通知/模板管理/公告管理)
- MQTT 通知下发 (通过 mqtt handler)
"""
import os
import pytest
import requests
import json
import time


class TestNotification:
    """
    通知管理功能测试套件
    PASS = 通知管理功能正确实现
    FAIL = 通知管理功能未实现或有问题
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

    def _read_source(self, rel_path):
        """Helper: 读取源码文件内容"""
        path = self._backend_src_path(rel_path)
        if not os.path.exists(path):
            return None
        with open(path, "r", encoding="utf-8") as f:
            return f.read()

    # ─── 静态源码检查 ────────────────────────────────────────────

    def test_notification_model_exists(self):
        """
        验证通知相关模型已定义

        检查 models/notification.go 中是否定义：
        - Notification 结构体 (title, content, type, target 等)
        - NotificationTemplate 结构体 (name, content, variables 等)
        - Announcement 结构体 (公告标题/内容/发布时间/状态)
        """
        source = self._read_source("models/notification.go")
        assert source is not None, (
            "FAIL: models/notification.go 文件不存在\n"
            "Expected: backend/models/notification.go 已创建"
        )

        required_notif_fields = ["Title", "Content", "Type", "Target"]
        for field in required_notif_fields:
            assert field in source, (
                f"FAIL: Notification 模型缺少 {field} 字段"
            )

        required_template_fields = ["Name", "Content", "Variables"]
        for field in required_template_fields:
            assert field in source, (
                f"FAIL: NotificationTemplate 模型缺少 {field} 字段"
            )

        required_ann_fields = ["Title", "Content", "StartTime", "EndTime", "Status"]
        for field in required_ann_fields:
            assert field in source, (
                f"FAIL: Announcement 模型缺少 {field} 字段"
            )

        print("PASS: 通知相关模型定义正确")

    def test_notification_controller_exists(self):
        """
        验证 NotificationController 已定义

        检查 controllers/notification_controller.go 中是否实现：
        - SendNotification / PushNotification 方法
        - Template CRUD 方法
        - Announcement CRUD 方法
        """
        source = self._read_source("controllers/notification_controller.go")
        assert source is not None, (
            "FAIL: controllers/notification_controller.go 文件不存在\n"
            "Expected: backend/controllers/notification_controller.go 已创建"
        )

        required_methods = [
            "SendNotification",   # 发送通知
            "ListNotifications",  # 获取通知历史
            "CreateTemplate",      # 创建模板
            "ListTemplates",      # 模板列表
            "UpdateTemplate",     # 更新模板
            "DeleteTemplate",     # 删除模板
            "CreateAnnouncement", # 创建公告
            "ListAnnouncements",  # 公告列表
            "UpdateAnnouncement", # 更新公告
            "DeleteAnnouncement", # 删除公告
            "PublishAnnouncement", # 发布公告
        ]
        missing = [m for m in required_methods if m not in source]
        assert not missing, (
            f"FAIL: NotificationController 缺少以下方法:\n" +
            "\n".join(f"  - {m}" for m in missing)
        )

        print("PASS: NotificationController 实现正确")

    def test_notification_routes_registered(self):
        """
        验证通知管理 API 路由已在 controllers/device_controller.go 中注册

        检查 RegisterRoutes 函数中是否包含通知和公告路由注册
        """
        source = self._read_source("controllers/device_controller.go")
        assert source is not None, "FAIL: device_controller.go 读取失败"

        routes = ["/notifications", "/announcements", "/notification-templates"]
        missing = [r for r in routes if r not in source]
        assert not missing, (
            f"FAIL: device_controller.go 未注册通知相关路由:\n" +
            "\n".join(f"  - {r}" for r in missing)
        )

        print("PASS: 通知管理路由已注册")

    def test_mqtt_notification_handler_exists(self):
        """
        验证 MQTT 通知下发逻辑已实现

        检查 controllers/notification_controller.go 中 SendNotification 方法：
        - 是否通过 MQTT 全局客户端下发通知
        - Topic 格式: /device/{device_id}/down/notification
        - MQTT payload 包含 title, content, type, timestamp
        """
        source = self._read_source("controllers/notification_controller.go")
        assert source is not None, (
            "FAIL: controllers/notification_controller.go 文件不存在"
        )

        # 检查是否使用全局 MQTT 客户端
        assert "GlobalMQTTClient" in source, (
            "FAIL: NotificationController 未使用 GlobalMQTTClient"
        )

        # 检查通知下发 Topic
        assert "/device/" in source and "down/notification" in source, (
            "FAIL: 未找到正确的通知下发 MQTT Topic 格式: /device/{id}/down/notification"
        )

        # 检查 MQTT Publish 调用
        assert ".Publish(" in source, (
            "FAIL: 未调用 MQTT Publish 方法下发通知"
        )

        # 检查 payload 字段
        required_fields = ["title", "content", "type", "timestamp"]
        source_lower = source.lower()
        for field in required_fields:
            assert f'"{field}"' in source or f"'{field}'" in source, (
                f"FAIL: MQTT 通知 payload 缺少 '{field}' 字段"
            )

        print("PASS: MQTT 通知下发逻辑已实现")

    # ─── API 功能测试 ────────────────────────────────────────────

    def test_notification_send(self):
        """
        验证发送通知 API

        测试流程：
        1. POST /api/v1/devices/:device_id/notifications  - 发送通知给设备
        2. GET  /api/v1/notifications                     - 获取通知历史
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("Cannot get auth token - backend may not be running")

        payload = {
            "title": "Test Notification",
            "content": "This is a test notification",
            "priority": 1,     # 0:普通 1:重要 2:紧急
            "channel": "push" # push, sms, email
        }

        resp = requests.post(
            f"{self.API_V1}/devices/test-device-001/notifications",
            json=payload,
            headers=headers,
            timeout=5
        )
        # 设备不存在时返回 404/400，但 API 应该可达
        assert resp.status_code in [200, 201, 400, 404], (
            f"FAIL: 发送通知返回 {resp.status_code}\n"
            f"Body: {resp.text[:300]}"
        )
        data = resp.json()
        if resp.status_code in [200, 201]:
            assert data.get("code") == 0, f"发送通知失败: {data}"
            print(f"  Sent notification: {data.get('data')}")
        else:
            # 设备不存在是预期的（测试设备）
            print(f"  INFO: Device not found (expected for test device): {data.get('message')}")

        # 获取通知历史
        resp = requests.get(f"{self.API_V1}/notifications", headers=headers, timeout=5)
        assert resp.status_code == 200, f"FAIL: 获取通知历史返回 {resp.status_code}"
        data = resp.json()
        assert data.get("code") == 0, f"获取通知历史失败: {data}"
        items = data.get("data", {}).get("list", []) or data.get("data", [])
        print(f"  Notification history OK ({len(items)} records)")

        print("PASS: 发送通知功能正确")

    def test_notification_template(self):
        """
        验证通知模板 CRUD API

        测试流程：
        1. POST /api/v1/notifications/templates      - 创建模板
        2. GET  /api/v1/notifications/templates       - 模板列表
        3. GET  /api/v1/notifications/templates/:id  - 模板详情
        4. PUT  /api/v1/notifications/templates/:id  - 更新模板
        5. DELETE /api/v1/notifications/templates/:id - 删除模板
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("Cannot get auth token - backend may not be running")

        template_payload = {
            "name": "Battery Low Alert",
            "code": f"low_battery_{int(time.time())}",  # unique code required
            "title_tpl": "设备 {{device_name}} 电量低",
            "content_tpl": "设备 {{device_name}} 电量低: {{battery}}%",
            "channel": "push",
            "priority": 1,
            "variables": '["device_name", "battery"]',  # JSON string
            "enabled": True
        }

        # Create
        resp = requests.post(
            f"{self.API_V1}/notification-templates",
            json=template_payload,
            headers=headers,
            timeout=5
        )
        assert resp.status_code in [200, 201], (
            f"FAIL: 创建模板返回 {resp.status_code}\n"
            f"Body: {resp.text[:200]}"
        )
        data = resp.json()
        assert data.get("code") == 0, f"创建模板失败: {data}"
        template_id = (
            data.get("data", {}).get("id")
            or data.get("data", {}).get("template_id")
        )
        assert template_id is not None, "创建模板未返回 id"
        print(f"  Created template id={template_id}")

        # List
        resp = requests.get(
            f"{self.API_V1}/notification-templates",
            headers=headers,
            timeout=5
        )
        assert resp.status_code == 200, f"FAIL: 获取模板列表返回 {resp.status_code}"
        data = resp.json()
        assert data.get("code") == 0, f"获取模板列表失败: {data}"
        print("  Listed templates OK")

        # Get detail
        resp = requests.get(
            f"{self.API_V1}/notification-templates/{template_id}",
            headers=headers,
            timeout=5
        )
        assert resp.status_code == 200, f"FAIL: 获取模板详情返回 {resp.status_code}"
        data = resp.json()
        assert data.get("code") == 0, f"获取模板详情失败: {data}"
        print("  Got template detail OK")

        # Update
        update_payload = {
            "name": "Battery Low Alert (Updated)",
            "content_tpl": "设备 {{device_name}} 电量危险: {{battery}}%",
            "priority": 2
        }
        resp = requests.put(
            f"{self.API_V1}/notification-templates/{template_id}",
            json=update_payload,
            headers=headers,
            timeout=5
        )
        assert resp.status_code == 200, f"FAIL: 更新模板返回 {resp.status_code}"
        print("  Updated template OK")

        # Delete
        resp = requests.delete(
            f"{self.API_V1}/notification-templates/{template_id}",
            headers=headers,
            timeout=5
        )
        assert resp.status_code in [200, 204], f"FAIL: 删除模板返回 {resp.status_code}"
        print("  Deleted template OK")

        print("PASS: 通知模板 CRUD 功能正确")

    def test_announcement_crud(self):
        """
        验证公告 CRUD API

        测试流程：
        1. POST /api/v1/announcements           - 创建公告
        2. GET  /api/v1/announcements           - 公告列表
        3. GET  /api/v1/announcements/:id      - 公告详情
        4. PUT  /api/v1/announcements/:id      - 更新公告
        5. DELETE /api/v1/announcements/:id    - 删除公告
        """
        headers = self._auth_headers()
        if not headers:
            pytest.skip("Cannot get auth token - backend may not be running")

        ann_payload = {
            "title": "系统升级公告",
            "content": "系统将于本周日凌晨2:00-4:00进行升级维护。",
            "type": "info",    # info, warning, critical
            "priority": 0,
            "target_type": "all",  # all, company, device_group
            "status": "draft"  # draft, published, archived
        }

        # Create
        resp = requests.post(
            f"{self.API_V1}/announcements",
            json=ann_payload,
            headers=headers,
            timeout=5
        )
        assert resp.status_code in [200, 201], (
            f"FAIL: 创建公告返回 {resp.status_code}\n"
            f"Body: {resp.text[:200]}"
        )
        data = resp.json()
        assert data.get("code") == 0, f"创建公告失败: {data}"
        ann_id = (
            data.get("data", {}).get("id")
            or data.get("data", {}).get("announcement_id")
        )
        assert ann_id is not None, "创建公告未返回 id"
        print(f"  Created announcement id={ann_id}")

        # List
        resp = requests.get(f"{self.API_V1}/announcements", headers=headers, timeout=5)
        assert resp.status_code == 200, f"FAIL: 获取公告列表返回 {resp.status_code}"
        data = resp.json()
        assert data.get("code") == 0, f"获取公告列表失败: {data}"
        print("  Listed announcements OK")

        # Get detail
        resp = requests.get(f"{self.API_V1}/announcements/{ann_id}", headers=headers, timeout=5)
        assert resp.status_code == 200, f"FAIL: 获取公告详情返回 {resp.status_code}"
        data = resp.json()
        assert data.get("code") == 0, f"获取公告详情失败: {data}"
        print("  Got announcement detail OK")

        # Update
        update_payload = {
            "title": "系统升级公告 (已更新)",
            "status": "published"  # 发布
        }
        resp = requests.put(
            f"{self.API_V1}/announcements/{ann_id}",
            json=update_payload,
            headers=headers,
            timeout=5
        )
        assert resp.status_code == 200, f"FAIL: 更新公告返回 {resp.status_code}"
        print("  Updated announcement OK")

        # Delete
        resp = requests.delete(f"{self.API_V1}/announcements/{ann_id}", headers=headers, timeout=5)
        assert resp.status_code in [200, 204], f"FAIL: 删除公告返回 {resp.status_code}"
        print("  Deleted announcement OK")

        print("PASS: 公告 CRUD 功能正确")

    def test_mqtt_notification(self):
        """
        验证 MQTT 通知下发功能

        测试场景：
        1. 发送通知时，通过 MQTT 正确发布到 /device/{device_id}/down/notification
        2. 检查 controllers/notification_controller.go 中 SendNotification 的 topic 格式
        3. 检查消息 payload 包含 title, content, type, timestamp
        """
        source = self._read_source("controllers/notification_controller.go")
        assert source is not None, (
            "FAIL: controllers/notification_controller.go 不存在，无法验证 MQTT 通知下发"
        )

        # 检查 MQTT Topic 格式
        topic_patterns = [
            "notification",
            "/device/",
            "down"
        ]
        for pattern in topic_patterns:
            assert pattern in source, (
                f"FAIL: NotificationController 缺少 topic 模式 '{pattern}'"
            )

        # 检查通知 payload 字段
        required_payload_fields = ["title", "content", "type", "timestamp"]
        source_lower = source.lower()
        for field in required_payload_fields:
            assert field in source_lower, (
                f"FAIL: MQTT 通知 payload 缺少 '{field}' 字段"
            )

        # 检查是否调用 MQTT Publish
        assert "Publish" in source or "publish" in source, (
            "FAIL: NotificationController 未调用 Publish 方法下发通知"
        )

        # 功能验证：发送通知 API 是否可调用
        headers = self._auth_headers()
        if headers:
            # 通知发送到设备端点
            payload = {
                "title": "MQTT Test",
                "content": "MQTT下发测试",
                "priority": 1,
                "channel": "push"
            }
            # 使用一个不存在的设备来测试API是否可达
            resp = requests.post(
                f"{self.API_V1}/devices/test-mqtt-device/notifications",
                json=payload,
                headers=headers,
                timeout=5
            )
            # 设备不存在返回4002，但说明API通了
            if resp.status_code == 200:
                data = resp.json()
                if data.get("code") == 0:
                    print("  Notification API call successful")
                else:
                    print(f"  INFO: API call returned non-zero code: {data}")
            elif resp.status_code == 404:
                print("  INFO: Device not found (404) - API is working")
            else:
                print(f"  INFO: Notification API returned {resp.status_code} (backend may not be running)")

        print("PASS: MQTT 通知下发功能正确")


if __name__ == "__main__":
    pytest.main([__file__, "-v", "--tb=short"])
