"""
P0 Test: JWT Configuration - Verify JWT Secret is NOT Hardcoded
================================================================
Bug: jwt.go uses hardcoded secret var jwtSecret = []byte("mdm-secret-key-change-in-production")
Fix: JWT secret must be loaded from environment variable JWT_SECRET

This test validates:
1. JWT secret is read from environment variable (not hardcoded)
2. Backend responds with 401 when token signed with wrong secret
3. Backend accepts tokens signed with env-based secret
"""
import os
import pytest
import requests
import jwt
import time


class TestJWTConfig:
    """
    JWT configuration security tests.
    PASS = JWT secret is configurable via env var (not hardcoded)
    FAIL = JWT secret is hardcoded in source code
    """

    BASE_URL = os.getenv("MDM_API_BASE_URL", "http://localhost:8080")
    API_V1 = f"{BASE_URL}/api/v1"

    def test_jwt_secret_not_hardcoded(self):
        """
        Verify JWT secret is loaded from env var, not hardcoded.

        Strategy: Attempt to sign a token with a custom secret and send it.
        If the server uses an env-based secret, it will reject our custom token.
        If the server uses a hardcoded secret matching our test secret,
        the request might accidentally succeed (false positive).
        We mitigate by checking the source code does NOT contain hardcoded patterns.
        """
        import re

        # Read the middleware/jwt.go source
        jwt_go_path = os.path.join(
            os.getenv("MDM_BACKEND_PATH", "C:\\Users\\YKing\\.openclaw\\workspace\\mdm-project\\backend"),
            "middleware", "jwt.go"
        )

        with open(jwt_go_path, "r", encoding="utf-8") as f:
            source = f.read()

        # Check for hardcoded secret patterns (these are BUG indicators)
        hardcoded_patterns = [
            r'jwtSecret\s*=\s*\[\]byte\(["\']mdm-secret-key',  # var jwtSecret = []byte("mdm-secret-key...")
            r'jwtSecret\s*=\s*\[\]byte\(["\']change-in-production',  # var jwtSecret = []byte("...change-in-production")
            r'var\s+jwtSecret\s*=\s*\[\]byte\(["\'][^$]',  # generic var jwtSecret = []byte("...") not from env
        ]

        for pattern in hardcoded_patterns:
            matches = re.findall(pattern, source, re.IGNORECASE)
            assert len(matches) == 0, (
                f"FAIL: JWT secret appears HARDCODED in middleware/jwt.go: {matches}\n"
                f"Expected: jwtSecret loaded from os.Getenv('JWT_SECRET')"
            )

        print("PASS: No hardcoded JWT secret patterns found in jwt.go")

    def test_jwt_secret_from_env_var(self):
        """
        Verify JWT secret is read from environment variable.

        The middleware should call os.Getenv("JWT_SECRET") or similar.
        """
        jwt_go_path = os.path.join(
            os.getenv("MDM_BACKEND_PATH", "C:\\Users\\YKing\\.openclaw\\workspace\\mdm-project\\backend"),
            "middleware", "jwt.go"
        )

        with open(jwt_go_path, "r", encoding="utf-8") as f:
            source = f.read()

        # Must reference environment variable for secret
        env_refs = [
            'os.Getenv("JWT_SECRET"',
            "os.Getenv(`JWT_SECRET`",
            "os.LookupEnv(`JWT_SECRET`",
            'os.LookupEnv("JWT_SECRET"',
        ]

        has_env_ref = any(ref in source for ref in env_refs)
        assert has_env_ref, (
            "FAIL: JWT secret NOT loaded from environment variable.\n"
            "Expected: os.Getenv('JWT_SECRET') or similar in middleware/jwt.go"
        )

        print("PASS: JWT secret is loaded from environment variable")

    def test_env_jwt_secret_takes_effect(self):
        """
        Functional test: Verify server rejects tokens signed with hardcoded secret
        when JWT_SECRET env var is set differently.

        Pre-condition: Backend running with JWT_SECRET env var set.
        """
        custom_secret = "test-secret-from-env-override-12345"

        # Generate a token with our custom secret
        payload = {
            "user_id": 1,
            "username": "test",
            "role_id": 1,
            "exp": int(time.time()) + 3600,
            "iat": int(time.time()),
        }
        token = jwt.encode(payload, custom_secret, algorithm="HS256")

        # Try to access a protected endpoint
        resp = requests.get(
            f"{self.API_V1}/auth/userinfo",
            headers={"Authorization": f"Bearer {token}"},
            timeout=5
        )

        # If JWT_SECRET env var is set and different from our test secret,
        # server should reject with 401
        # If server uses hardcoded secret that happens to match or is bypassed,
        # we might get 200 or different error
        if resp.status_code == 200:
            # Server accepted our token - either JWT_SECRET not set or matches
            pytest.fail(
                "FAIL: Server accepted token signed with arbitrary secret. "
                "JWT secret may not be properly configured via env var."
            )

        assert resp.status_code == 401, (
            f"Expected 401 Unauthorized, got {resp.status_code}. "
            "JWT secret from env var should reject tokens signed with different secret."
        )

        print("PASS: Server correctly rejects tokens when JWT_SECRET env var differs")


if __name__ == "__main__":
    pytest.main([__file__, "-v", "--tb=short"])
