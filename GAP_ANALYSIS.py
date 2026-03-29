#!/usr/bin/env python3
"""
Comprehensive Gap Analysis: Compare Controller Files vs Registered Routes
"""
import os
import re
from pathlib import Path

BACKEND = Path("C:/Users/YKing/.openclaw/workspace/mdm-project/backend")
CONTROLLERS = BACKEND / "controllers"
MAIN_GO = BACKEND / "main.go"
DEVICE_CTRL = CONTROLLERS / "device_controller.go"

# 1. Find all controller files
controller_files = []
for f in CONTROLLERS.glob("*_controller.go"):
    name = f.stem.replace("_controller", "")
    controller_files.append(name)

# 2. Extract all registered routes from main.go and device_controller.go
registered_routes = set()
registered_controllers = set()

files_to_scan = [MAIN_GO, DEVICE_CTRL]
for filepath in files_to_scan:
    content = filepath.read_text(encoding='utf-8')
    
    # Find api.Group or similar patterns
    # Pattern: api.GET("/xxx", ctrl.Method)
    route_pattern = re.compile(r'api\.(GET|POST|PUT|DELETE|PATCH|OPTIONS)\s*\(\s*"([^"]+)"')
    for match in route_pattern.finditer(content):
        route = match.group(2)
        registered_routes.add(route)
    
    # Find controller instantiation patterns
    ctrl_pattern = re.compile(r'(\w+Ctrl)\s*:?=\s*(?:&?\w+\{)?.*?\}', re.MULTILINE)
    for match in ctrl_pattern.finditer(content):
        ctrl_name = match.group(1).replace("Ctrl", "")
        registered_controllers.add(ctrl_name)

# 3. Find controllers that have files but NO routes registered
missing_routes = []
for ctrl in sorted(controller_files):
    if ctrl in ['pagination', 'helpers', 'device_crud']:
        continue
    # Check if controller name appears in registered routes or as a registered controller
    ctrl_underscore = ctrl.replace("_", "-")
    ctrl_hyphen = ctrl.replace("_", "-")
    
    # Look for routes that start with this controller's path
    has_route = any(r.startswith(f"/{ctrl_hyphen}") or r.startswith(f"/{ctrl_underscore}") for r in registered_routes)
    has_ctrl = ctrl in registered_controllers
    
    if not has_route and not has_ctrl:
        missing_routes.append(ctrl)

# 4. Find 404 routes from API test (controllers exist but route returns 404)
# Based on API test results
controllers_with_404 = [
    "action_learning", "action_library", "advanced", "ai (chat/models)", 
    "alert_history", "alert_rules", "alert_settings", "alerts_dedup", 
    "alerts_healing", "analytics", "batch", "card", "child_mode",
    "compliance", "content (review)", "coupon", "data_masking",
    "data_residency", "device_health", "device_monitor", "device_security",
    "device_shadow", "digital_twin", "emotion", "family_album",
    "flow", "gdpr", "health", "insurance", "integration",
    "map", "market", "member_enhanced", "member_profile", "miniapp",
    "miniclaw", "offline_sync", "ota_compatibility", "ota_sdk",
    "pet_profile", "pet_social", "platform_evo", "policy",
    "quota", "regions", "reports", "security_evo", "simulation",
    "smart_home", "subscription_gift", "timezones", "voice_emotion"
]

print("=" * 60)
print("GAP ANALYSIS: Controller Files vs Registered Routes")
print("=" * 60)

print(f"\nTotal Controller Files: {len(controller_files)}")
print(f"Total Registered Routes: {len(registered_routes)}")
print(f"Registered Controllers: {len(registered_controllers)}")

print("\n" + "=" * 60)
print("ISSUE CATEGORY 1: Controller files exist but NO routes registered")
print("=" * 60)
for ctrl in sorted(missing_routes):
    print(f"  [MISSING] {ctrl}_controller.go - NO ROUTES REGISTERED")

print("\n" + "=" * 60)
print("ISSUE CATEGORY 2: Routes registered but return 404/500")
print("=" * 60)
# These are controllers that ARE registered but their endpoints return errors
for ctrl in sorted(set(controllers_with_404)):
    print(f"  [ERROR] {ctrl} - returns 404/500")

print("\n" + "=" * 60)
print("SUMMARY")
print("=" * 60)
print(f"Controllers with missing routes: {len(missing_routes)}")
print(f"Controllers returning errors: {len(controllers_with_404)}")
print(f"Total issues: {len(missing_routes) + len(controllers_with_404)}")

# Priority list for fixing
print("\n" + "=" * 60)
print("PRIORITY FIX LIST (Grouped by module)")
print("=" * 60)
modules = {
    "AI": ["ai (chat/models)", "emotion", "voice_emotion"],
    "DEVICE": ["device_health", "device_monitor", "device_security", "device_shadow"],
    "PET": ["pet_profile", "pet_social", "action_library", "action_learning"],
    "ENTERPRISE": ["insurance", "billing", "analytics", "reports", "advanced"],
    "COMPLIANCE": ["security_evo", "compliance", "gdpr", "data_residency", "data_masking"],
    "FAMILY": ["family_album", "child_mode", "smart_home"],
    "CONTENT": ["content (review)", "market", "miniapp", "miniclaw"],
    "PLATFORM": ["ota_compatibility", "ota_sdk", "integration", "batch", "flow"],
    "OPERATIONS": ["subscription_gift", "quota", "regions", "timezones", "schedule"],
    "DIGITAL_TWIN": ["digital_twin", "health", "simulation", "map"],
    "ALERTS": ["alert_rules", "alert_settings", "alert_history", "alerts_dedup", "alerts_healing"],
    "MEMBER": ["member_profile", "member_enhanced", "card", "coupon"],
}

for module, controllers in modules.items():
    existing = [c for c in controllers if c in controllers_with_404 or any(c.replace(" (review)", "") in x or c.replace(" (chat/models)", "") in x for x in controllers_with_404)]
    if existing:
        print(f"\n[{module}]")
        for c in existing:
            print(f"   - {c}")
