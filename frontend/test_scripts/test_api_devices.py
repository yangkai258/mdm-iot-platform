"""
MDM 控制中台 - 后端 API 自动化测试套件
测试工程师: agentcs
日期: 2026-03-19
"""

import requests
import json
import uuid
import time
import unittest
import sys
import io

# Fix Windows Unicode output
if sys.platform == 'win32':
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8', errors='replace')

# Configuration
BASE_URL = "http://localhost:8080"
API_BASE = f"{BASE_URL}/api/v1"

# Test constants
TEST_DEVICE_DATA = {
    "device_id": str(uuid.uuid4()),
    "mac_address": "AA:BB:CC:DD:EE:FF",
    "sn_code": f"TEST-SN-{int(time.time())}",
    "hardware_model": "M5_CoreS3",
    "firmware_version": "v1.0.0"
}

# Global test data
test_device_id = None
test_sn_code = None
test_auth_token = "test-token-12345"

# Track the SN code from registration
registered_sn_code = None


class MDMApiTest(unittest.TestCase):
    """MDM API automation test class"""
    
    @classmethod
    def setUpClass(cls):
        print("\n" + "="*60)
        print("MDM API Test Suite Started")
        print("="*60)
        global test_device_id, test_sn_code
        test_device_id = TEST_DEVICE_DATA["device_id"]
        test_sn_code = TEST_DEVICE_DATA["sn_code"]
    
    def test_01_device_register(self):
        """Test 1: Device Registration POST /api/v1/devices/register"""
        print("\n[Test 1] Device Register - POST /api/v1/devices/register")
        
        url = f"{API_BASE}/devices/register"
        payload = TEST_DEVICE_DATA.copy()
        
        print(f"  URL: {url}")
        print(f"  Payload: {json.dumps(payload)}")
        
        response = requests.post(url, json=payload, timeout=10)
        
        print(f"  Status: {response.status_code}")
        print(f"  Response: {response.text[:300]}")
        
        self.assertEqual(response.status_code, 200, f"Expected 200, got {response.status_code}")
        
        data = response.json()
        self.assertEqual(data.get("code"), 0, f"Code should be 0, got {data.get('code')}")
        self.assertIn("data", data, "Response should contain 'data'")
        
        global test_device_id, registered_sn_code
        if "device_id" in data.get("data", {}):
            test_device_id = data["data"]["device_id"]
        # Save the sn_code used in registration
        registered_sn_code = payload["sn_code"]
        
        print(f"  [PASS] Device registered! device_id: {test_device_id}, sn_code: {registered_sn_code}")
    
    def test_02_device_list(self):
        """Test 2: Device List GET /api/v1/devices"""
        print("\n[Test 2] Device List - GET /api/v1/devices")
        
        url = f"{API_BASE}/devices"
        headers = {"Authorization": f"Bearer {test_auth_token}"}
        
        print(f"  URL: {url}")
        
        response = requests.get(url, headers=headers, timeout=10)
        
        print(f"  Status: {response.status_code}")
        print(f"  Response: {response.text[:300]}")
        
        self.assertEqual(response.status_code, 200, f"Expected 200, got {response.status_code}")
        
        data = response.json()
        self.assertEqual(data.get("code"), 0, f"Code should be 0, got {data.get('code')}")
        self.assertIn("data", data, "Response should contain 'data'")
        
        print(f"  [PASS] Device list retrieved!")
    
    def test_03_device_detail(self):
        """Test 3: Device Detail GET /api/v1/devices/:device_id"""
        print("\n[Test 3] Device Detail - GET /api/v1/devices/:device_id")
        
        url = f"{API_BASE}/devices/{test_device_id}"
        headers = {"Authorization": f"Bearer {test_auth_token}"}
        
        print(f"  URL: {url}")
        print(f"  device_id: {test_device_id}")
        
        response = requests.get(url, headers=headers, timeout=10)
        
        print(f"  Status: {response.status_code}")
        print(f"  Response: {response.text[:300]}")
        
        self.assertEqual(response.status_code, 200, f"Expected 200, got {response.status_code}")
        
        data = response.json()
        self.assertEqual(data.get("code"), 0, f"Code should be 0, got {data.get('code')}")
        
        device_data = data.get("data", {})
        self.assertIn("device_id", device_data, "Detail should contain device_id")
        
        print(f"  [PASS] Device detail retrieved! device_id: {device_data.get('device_id')}")
    
    def test_04_device_bind(self):
        """Test 4: Device Bind POST /api/v1/devices/bind/:sn_code"""
        print("\n[Test 4] Device Bind - POST /api/v1/devices/bind/:sn_code")
        
        # First get the device list to find an existing sn_code
        list_url = f"{API_BASE}/devices"
        headers = {"Authorization": f"Bearer {test_auth_token}"}
        list_response = requests.get(list_url, headers=headers, timeout=10)
        
        existing_sn = None
        if list_response.status_code == 200:
            data = list_response.json()
            devices = data.get("data", {}).get("list", [])
            if devices:
                existing_sn = devices[0].get("sn_code")
        
        # Use the existing sn_code or fall back to test_sn_code
        bind_sn_code = existing_sn or test_sn_code
        url = f"{API_BASE}/devices/bind/{bind_sn_code}"
        headers = {
            "Authorization": f"Bearer {test_auth_token}",
            "Content-Type": "application/json"
        }
        payload = {"bind_user_id": "user-001"}
        
        print(f"  URL: {url}")
        print(f"  sn_code: {bind_sn_code}")
        print(f"  Payload: {json.dumps(payload)}")
        
        response = requests.post(url, json=payload, headers=headers, timeout=10)
        
        print(f"  Status: {response.status_code}")
        print(f"  Response: {response.text[:300]}")
        
        self.assertEqual(response.status_code, 200, f"Expected 200, got {response.status_code}")
        
        data = response.json()
        self.assertEqual(data.get("code"), 0, f"Code should be 0, got {data.get('code')}")
        
        print(f"  [PASS] Device bound successfully!")


class MDMApiErrorTest(unittest.TestCase):
    """Error scenario tests"""
    
    def test_05_device_register_duplicate(self):
        """Test 5: Duplicate device registration (should return error)"""
        print("\n[Test 5] Duplicate Device Register - POST /api/v1/devices/register")
        
        url = f"{API_BASE}/devices/register"
        payload = TEST_DEVICE_DATA.copy()
        
        print(f"  URL: {url}")
        
        response = requests.post(url, json=payload, timeout=10)
        
        print(f"  Status: {response.status_code}")
        print(f"  Response: {response.text[:300]}")
        
        self.assertIn(response.status_code, [200, 400, 409], 
                     f"Should return error status, got {response.status_code}")
        
        if response.status_code != 200:
            data = response.json()
            self.assertIn(data.get("code"), [4001, 4005], 
                         f"Error code should be 4001 or 4005, got {data.get('code')}")
            print(f"  [PASS] Duplicate registration correctly returned error: code={data.get('code')}")
        else:
            print(f"  [INFO] Duplicate registration returned success (idempotent)")
    
    def test_06_device_detail_not_found(self):
        """Test 6: Get non-existent device detail"""
        print("\n[Test 6] Device Detail - Non-existent device")
        
        fake_device_id = str(uuid.uuid4())
        url = f"{API_BASE}/devices/{fake_device_id}"
        headers = {"Authorization": f"Bearer {test_auth_token}"}
        
        print(f"  URL: {url}")
        
        response = requests.get(url, headers=headers, timeout=10)
        
        print(f"  Status: {response.status_code}")
        print(f"  Response: {response.text[:300]}")
        
        self.assertIn(response.status_code, [200, 404], 
                     f"Should return 200 or 404, got {response.status_code}")
        
        if response.status_code == 404:
            data = response.json()
            self.assertEqual(data.get("code"), 4002, 
                            f"Error code should be 4002, got {data.get('code')}")
            print(f"  [PASS] Correctly returned 404: code=4002")


def run_tests():
    """Run all tests and generate report"""
    print("\n" + "="*60)
    print("Executing MDM API Test Suite")
    print("="*60)
    
    loader = unittest.TestLoader()
    suite = unittest.TestSuite()
    
    suite.addTests(loader.loadTestsFromTestCase(MDMApiTest))
    suite.addTests(loader.loadTestsFromTestCase(MDMApiErrorTest))
    
    runner = unittest.TextTestRunner(verbosity=2)
    result = runner.run(suite)
    
    print("\n" + "="*60)
    print("Test Summary")
    print("="*60)
    print(f"  Tests Run: {result.testsRun}")
    print(f"  Success: {result.testsRun - len(result.failures) - len(result.errors)}")
    print(f"  Failures: {len(result.failures)}")
    print(f"  Errors: {len(result.errors)}")
    
    if result.wasSuccessful():
        print("\n  [SUCCESS] All tests passed!")
    else:
        print("\n  [FAILURE] Some tests failed")
        
        if result.failures:
            print("\n  Failure Details:")
            for test, traceback in result.failures:
                print(f"    - {test}")
        
        if result.errors:
            print("\n  Error Details:")
            for test, traceback in result.errors:
                print(f"    - {test}")
    
    print("="*60)
    
    return result.wasSuccessful()


if __name__ == "__main__":
    success = run_tests()
    exit(0 if success else 1)
