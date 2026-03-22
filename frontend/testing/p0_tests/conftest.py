# MDM P0 Tests - Shared Configuration
"""
Pytest configuration and shared fixtures for MDM P0 test suite.
"""
import os
import pytest

# Base URL for backend API
BASE_URL = os.getenv("MDM_API_BASE_URL", "http://localhost:8080")
API_V1 = f"{BASE_URL}/api/v1"

# MQTT Broker config
MQTT_BROKER = os.getenv("MQTT_BROKER", "tcp://localhost:1883")
MQTT_USERNAME = os.getenv("MQTT_USERNAME", "admin")
MQTT_PASSWORD = os.getenv("MQTT_PASSWORD", "public")


@pytest.fixture
def api_base_url():
    """Base URL for API requests."""
    return BASE_URL


@pytest.fixture
def api_v1_url():
    """API v1 base URL."""
    return API_V1


@pytest.fixture
def auth_token():
    """
    Returns a valid JWT token by logging in with test credentials.
    Requires MDM backend to be running with seeded test user.
    """
    import requests
    login_data = {
        "username": os.getenv("TEST_USERNAME", "admin"),
        "password": os.getenv("TEST_PASSWORD", "admin123")
    }
    resp = requests.post(f"{API_V1}/auth/login", json=login_data, timeout=5)
    if resp.status_code == 200:
        data = resp.json()
        if data.get("code") == 0:
            return data["data"]["token"]
    return None


@pytest.fixture
def auth_headers(auth_token):
    """Authorization headers with JWT token."""
    if auth_token:
        return {"Authorization": f"Bearer {auth_token}"}
    return {}
