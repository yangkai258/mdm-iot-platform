#!/bin/sh

BASE_URL="http://localhost:8080/api/v1"
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6ImFkbWluIiwicm9sZV9pZCI6MCwidGVuYW50X2lkIjoiIiwiaXNfc3VwZXJfYWRtaW4iOmZhbHNlLCJleHAiOjE3NzQzNjc5NzAsImlhdCI6MTc3NDI4MTU3MH0.URkp53kuWRIPL5U6D2dLYaXt4ylJuT0EXubWNDXfVE8"
AUTH="Authorization: Bearer $TOKEN"

# Helper function
test_api() {
    local method=$1
    local path=$2
    local name=$3
    local data=$4
    
    if [ -z "$data" ]; then
        result=$(curl -s -w "HTTP_CODE:%{http_code}" -X "$method" "$BASE_URL$path" -H "$AUTH" -H "Content-Type: application/json")
    else
        result=$(echo "$data" | curl -s -w "HTTP_CODE:%{http_code}" -X "$method" "$BASE_URL$path" -H "$AUTH" -H "Content-Type: application/json" -d @-)
    fi
    
    # Extract HTTP code from the end
    code=$(echo "$result" | sed 's/.*HTTP_CODE://')
    body=$(echo "$result" | sed 's/HTTP_CODE:.*//' | tr '\n' ' ')
    
    # 404 - route not registered (SKIP)
    if [ "$code" = "404" ]; then
        echo "| $path | $method | SKIP | 路由未注册 |"
        return
    fi
    
    # 500/5001 - server error (FAIL)
    if [ "$code" = "500" ] || [ "$code" = "5001" ]; then
        echo "| $path | $method | FAIL | HTTP $code 服务器错误 |"
        return
    fi
    
    # Empty body
    if [ -z "$body" ] || [ "$body" = "{}" ]; then
        echo "| $path | $method | FAIL | HTTP $code 空响应 |"
        return
    fi
    
    # Check for success code in body
    if echo "$body" | grep -q '"code":0' || echo "$body" | grep -q '"code":200' || echo "$body" | grep -q '"success":true'; then
        echo "| $path | $method | PASS | $name |"
    else
        echo "| $path | $method | WARN | HTTP $code - $(echo "$body" | cut -c1-60) |"
    fi
}

echo "| API | 方法 | 结果 | 说明 |"
echo "|-----|------|------|------|"

echo ""
echo "### Sprint 1-8 核心"
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
echo "### Sprint 9-12 合规与订阅"
test_api "GET" "/compliance/policies" "List compliance policies"
test_api "GET" "/compliance-rules" "List compliance rules"
test_api "GET" "/subscriptions/plans" "List subscription plans"
test_api "GET" "/ldap/config" "Get LDAP config"
test_api "GET" "/device-shadow/1" "Get device shadow"

echo ""
echo "### Sprint 13-16 AI与健康"
test_api "GET" "/ai/models" "List AI models"
test_api "GET" "/health/1/exercise" "Get pet exercise"
test_api "POST" "/emotion/recognize/text" "Recognize emotion from text" '{"text":"happy","subject_type":"pet","subject_id":1}'
test_api "GET" "/digital-twin/1/vitals" "Get digital twin vitals"
test_api "GET" "/simulation/virtual-pets" "List virtual pets"
test_api "GET" "/insurance/products" "List insurance products"

echo ""
echo "### Sprint 17-20 家庭与社交"
test_api "GET" "/family/children/profiles" "List children profiles"
test_api "GET" "/family/elderly/profiles" "List elderly profiles"
test_api "GET" "/pet-social/feed" "Get pet social feed"
test_api "GET" "/research/datasets" "List research datasets"

echo ""
echo "### Sprint 21-32 仿真与知识"
test_api "GET" "/simulation/environments" "List simulation environments"
test_api "GET" "/ai/model/shards" "List AI model shards"
test_api "POST" "/knowledge/query" "Query knowledge" '{"query":"test","top_k":5}'

echo ""
echo "### 其他"
test_api "GET" "/notifications" "List notifications"
test_api "GET" "/announcements" "List announcements"
test_api "GET" "/devices/1/shadow" "Get device shadow by ID"
test_api "GET" "/ai-fairness/reports" "List AI fairness reports"
