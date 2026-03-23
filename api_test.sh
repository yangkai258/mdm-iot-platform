#!/bin/sh

BASE_URL="http://localhost:8080/api/v1"
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6ImFkbWluIiwicm9sZV9pZCI6MCwidGVuYW50X2lkIjoiIiwiaXNfc3VwZXJfYWRtaW4iOmZhbHNlLCJleHAiOjE3NzQzNjc5MjMsImlhdCI6MTc0MjgxNTIzfQ.D-DH-mjMVp4dO7w5p5JNEqWRJcdRzPkgiXz0Wx2_BDI"
AUTH="Authorization: Bearer $TOKEN"

echo "=== MDM Backend API Test Report ==="
echo ""

# Helper function to test API
test_api() {
    local method=$1
    local path=$2
    local name=$3
    local data=$4
    
    if [ -z "$data" ]; then
        result=$(curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$path" -H "$AUTH" -H "Content-Type: application/json")
    else
        result=$(echo "$data" | curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$path" -H "$AUTH" -H "Content-Type: application/json" -d @-)
    fi
    
    code=$(echo "$result" | tail -1)
    body=$(echo "$result" | sed '$d')
    
    if [ "$code" = "200" ] || [ "$code" = "201" ]; then
        # Check if body contains success indicator
        if echo "$body" | grep -q '"code":0' || echo "$body" | grep -q '"code":200'; then
            echo "[PASS] $name ($method $path) - HTTP $code"
        else
            echo "[WARN] $name ($method $path) - HTTP $code (body may have error)"
            echo "       Body: $(echo "$body" | cut -c1-100)"
        fi
    elif [ "$code" = "404" ]; then
        echo "[SKIP] $name ($method $path) - HTTP 404 (route not registered)"
    elif [ "$code" = "401" ]; then
        echo "[FAIL] $name ($method $path) - HTTP 401 (unauthorized)"
    elif [ "$code" = "500" ] || [ "$code" = "5001" ]; then
        echo "[FAIL] $name ($method $path) - HTTP $code (server error)"
        echo "       Body: $(echo "$body" | cut -c1-150)"
    else
        # Check body for error
        if echo "$body" | grep -q '"code":0' || echo "$body" | grep -q '"code":200'; then
            echo "[PASS] $name ($method $path) - HTTP $code"
        elif echo "$body" | grep -q '"code":404'; then
            echo "[SKIP] $name ($method $path) - HTTP 404 (route not registered)"
        else
            echo "[INFO] $name ($method $path) - HTTP $code"
            echo "       Body: $(echo "$body" | cut -c1-100)"
        fi
    fi
}

echo "--- Sprint 1-8: Core APIs ---"
test_api "GET" "/devices" "List devices"
test_api "GET" "/members" "List members"
test_api "GET" "/alerts" "List alerts"
test_api "GET" "/member/cards" "List member cards"
test_api "GET" "/member/levels" "List member levels"
test_api "GET" "/member/tags" "List member tags"
test_api "GET" "/coupons" "List coupons"
test_api "GET" "/promotions" "List promotions"
test_api "GET" "/stores" "List stores"
test_api "POST" "/orders" "Create order" '{"total_amount":100}'
test_api "GET" "/knowledge" "List knowledge"

echo ""
echo "--- Sprint 9-12: Compliance & Subscriptions ---"
test_api "GET" "/compliance/policies" "List compliance policies"
test_api "GET" "/compliance-rules" "List compliance rules"
test_api "GET" "/subscriptions/plans" "List subscription plans"
test_api "GET" "/ldap/config" "Get LDAP config"
test_api "GET" "/device-shadow/1" "Get device shadow"

echo ""
echo "--- Sprint 13-16: AI & Health ---"
test_api "GET" "/ai/models" "List AI models"
test_api "GET" "/health/1/exercise" "Get pet exercise"
test_api "POST" "/emotion/recognize/text" "Recognize emotion from text" '{"text":"happy"}'
test_api "GET" "/digital-twin/1/vitals" "Get digital twin vitals"
test_api "GET" "/simulation/virtual-pets" "List virtual pets"
test_api "GET" "/insurance/products" "List insurance products"

echo ""
echo "--- Sprint 17-20: Family & Social ---"
test_api "GET" "/family/children/profiles" "List children profiles"
test_api "GET" "/family/elderly/profiles" "List elderly profiles"
test_api "GET" "/pet-social/feed" "Get pet social feed"
test_api "GET" "/research/datasets" "List research datasets"

echo ""
echo "--- Sprint 21-32: Simulation & Knowledge ---"
test_api "GET" "/simulation/environments" "List simulation environments"
test_api "GET" "/ai/model/shards" "List AI model shards"
test_api "POST" "/knowledge/query" "Query knowledge" '{"query":"test"}'

echo ""
echo "--- Other APIs ---"
test_api "GET" "/notifications" "List notifications"
test_api "GET" "/announcements" "List announcements"
test_api "GET" "/devices/1/shadow" "Get device shadow by ID"
test_api "GET" "/ai-fairness/reports" "List AI fairness reports"

echo ""
echo "=== Test Complete ==="
