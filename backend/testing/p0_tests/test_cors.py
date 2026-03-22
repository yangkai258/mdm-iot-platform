"""
P0 Test: CORS Configuration - Verify CORS Does NOT Allow All Origins
=====================================================================
Bug: main.go sets Access-Control-Allow-Origin to "*" (allow all)
Fix: CORS must whitelist specific domains only

This test validates:
1. CORS header is NOT "*"
2. CORS only allows configured/trusted origins
3. Pre-flight OPTIONS requests return proper CORS headers
"""
import os
import pytest
import requests


class TestCORSConfig:
    """
    CORS security configuration tests.
    PASS = CORS restricts origins to whitelisted domains
    FAIL = CORS allows all origins ("*")
    """

    BASE_URL = os.getenv("MDM_API_BASE_URL", "http://localhost:8080")
    API_V1 = f"{BASE_URL}/api/v1"

    def test_cors_header_not_wildcard(self):
        """
        Verify Access-Control-Allow-Origin is NOT set to "*".

        Wildcard CORS allows any website to make authenticated requests.
        """
        resp = requests.get(f"{self.API_V1}/auth/login", timeout=5)

        allow_origin = resp.headers.get("Access-Control-Allow-Origin", "")

        assert allow_origin != "*", (
            "FAIL: CORS Access-Control-Allow-Origin is set to '*' (WILDCARD).\n"
            "This allows ANY website to make authenticated requests to the API.\n"
            "Expected: Specific domain(s) like 'https://mdm.example.com'"
        )

        print(f"PASS: Access-Control-Allow-Origin is restricted to: {allow_origin}")

    def test_cors_header_is_configured(self):
        """
        Verify CORS headers are properly set (not missing).
        """
        resp = requests.get(f"{self.API_V1}/auth/login", timeout=5)

        allow_origin = resp.headers.get("Access-Control-Allow-Origin")
        allow_methods = resp.headers.get("Access-Control-Allow-Methods")
        allow_headers = resp.headers.get("Access-Control-Allow-Headers")

        # At minimum, allow_origin must be set
        assert allow_origin is not None and allow_origin != "", (
            "FAIL: Access-Control-Allow-Origin header is missing or empty.\n"
            "CORS headers must be properly configured."
        )

        print(f"PASS: CORS headers configured - Origin: {allow_origin}, "
              f"Methods: {allow_methods}, Headers: {allow_headers}")

    def test_cors_preflight_options(self):
        """
        Verify pre-flight OPTIONS requests return proper CORS headers.
        """
        headers = {
            "Origin": "http://localhost:3000",
            "Access-Control-Request-Method": "POST",
            "Access-Control-Request-Headers": "Content-Type,Authorization",
        }

        resp = requests.options(f"{self.API_V1}/auth/login", headers=headers, timeout=5)

        allow_origin = resp.headers.get("Access-Control-Allow-Origin", "")

        assert allow_origin != "*", (
            "FAIL: Pre-flight CORS Access-Control-Allow-Origin is '*' (WILDCARD)."
        )

        # Should return 200 or 204 for preflight
        assert resp.status_code in (200, 204), (
            f"Expected 200/204 for pre-flight OPTIONS, got {resp.status_code}"
        )

        print(f"PASS: Pre-flight CORS handled correctly, allowed origin: {allow_origin}")

    def test_cors_with_different_origin_rejected(self):
        """
        Verify that origins NOT in the whitelist are NOT reflected in
        Access-Control-Allow-Origin (browser will block the request).
        """
        # Send request from an untrusted origin
        resp = requests.get(
            f"{self.API_V1}/auth/login",
            headers={"Origin": "https://evil-website.com"},
            timeout=5
        )

        allow_origin = resp.headers.get("Access-Control-Allow-Origin", "")

        # If allow_origin is specific and not *, then it should either:
        # 1. Match our trusted origin (not the untrusted one), OR
        # 2. Be empty/missing (origin not allowed)
        # In either case, it should NOT reflect the attacker's origin

        assert allow_origin != "https://evil-website.com", (
            "FAIL: CORS incorrectly reflects untrusted origin in "
            "Access-Control-Allow-Origin header."
        )

        print(f"PASS: Untrusted origin not reflected in CORS response")

    def test_cors_config_file_exists(self):
        """
        Verify CORS configuration is managed via config/env, not hardcoded in main.go.

        Check that the CORS middleware reads from environment or config file.
        """
        import re

        main_go_path = os.path.join(
            os.getenv("MDM_BACKEND_PATH", "C:\\Users\\YKing\\.openclaw\\workspace\\mdm-project\\backend"),
            "main.go"
        )

        with open(main_go_path, "r", encoding="utf-8") as f:
            source = f.read()

        # Check for hardcoded wildcard in CORS context
        hardcoded_wildcard = 'Access-Control-Allow-Origin", "*"'

        assert hardcoded_wildcard not in source, (
            "FAIL: CORS Allow-Origin '*' is hardcoded in main.go.\n"
            "Expected: CORS origin loaded from config or environment variable."
        )

        # Check for proper CORS config patterns
        config_patterns = [
            "CORS_ALLOWED_ORIGINS",
            "CORS_ORIGINS",
            "ALLOWED_ORIGINS",
            "Getenv",
        ]

        has_config = any(pattern in source for pattern in config_patterns)
        assert has_config, (
            "FAIL: No evidence of CORS configuration from env/config.\n"
            "Expected: CORS origins loaded from environment or config."
        )

        print("PASS: CORS configuration is externalized (not hardcoded)")


if __name__ == "__main__":
    pytest.main([__file__, "-v", "--tb=short"])
