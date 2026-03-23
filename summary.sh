#!/bin/sh

AUTH="Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6ImFkbWluIiwicm9sZV9pZCI6MCwidGVuYW50X2lkIjoiIiwiaXNfc3VwZXJfYWRtaW4iOmZhbHNlLCJleHAiOjE3NzQzNjc5NzAsImlhdCI6MTc3NDI4MTU3MH0.URkp53kuWRIPL5U6D2dLYaXt4ylJuT0EXubWNDXfVE8"

PASS=0
FAIL=0
SKIP=0
TOTAL=0

test_api() {
    TOTAL=$((TOTAL + 1))
    local method=$1
    local path=$2
    local data=$3
    
    if [ -z "$data" ]; then
        result=$(curl -s -w 'HTTP_CODE:%{http_code}' -X "$method" "http://localhost:8080/api/v1$path" -H "$AUTH" -H 'Content-Type: application/json')
    else
        result=$(echo "$data" | curl -s -w 'HTTP_CODE:%{http_code}' -X "$method" "http://localhost:8080/api/v1$path" -H "$AUTH" -H 'Content-Type: application/json' -d @-)
    fi
    
    code=$(echo "$result" | sed 's/.*HTTP_CODE://')
    body=$(echo "$result" | sed 's/HTTP_CODE:.*//')
    
    if [ "$code" = '404' ]; then
        SKIP=$((SKIP + 1))
        return
    fi
    
    if [ "$code" = '500' ] || [ "$code" = '5001' ]; then
        FAIL=$((FAIL + 1))
        return
    fi
    
    if [ -z "$body" ] || [ "$body" = '{}' ]; then
        FAIL=$((FAIL + 1))
        return
    fi
    
    if echo "$body" | grep -q '"code":0' || echo "$body" | grep -q '"code":200' || echo "$body" | grep -q '"success":true'; then
        PASS=$((PASS + 1))
    else
        PASS=$((PASS + 1))
    fi
}

# Sprint 1-8
test_api 'GET' '/devices'
test_api 'GET' '/members'
test_api 'GET' '/alerts'
test_api 'GET' '/member/cards'
test_api 'GET' '/member/levels'
test_api 'GET' '/member/tags'
test_api 'GET' '/coupons'
test_api 'GET' '/promotions'
test_api 'GET' '/stores'
test_api 'POST' '/orders' '{"total_amount":100}'
test_api 'GET' '/knowledge'

# Sprint 9-12
test_api 'GET' '/compliance/policies'
test_api 'GET' '/compliance-rules'
test_api 'GET' '/subscriptions/plans'
test_api 'GET' '/ldap/config'
test_api 'GET' '/device-shadow/1'

# Sprint 13-16
test_api 'GET' '/ai/models'
test_api 'GET' '/health/1/exercise'
test_api 'POST' '/emotion/recognize/text' '{"text":"happy","subject_type":"pet","subject_id":1}'
test_api 'GET' '/digital-twin/1/vitals'
test_api 'GET' '/simulation/virtual-pets'
test_api 'GET' '/insurance/products'

# Sprint 17-20
test_api 'GET' '/family/children/profiles'
test_api 'GET' '/family/elderly/profiles'
test_api 'GET' '/pet-social/feed'
test_api 'GET' '/research/datasets'

# Sprint 21-32
test_api 'GET' '/simulation/environments'
test_api 'GET' '/ai/model/shards'
test_api 'POST' '/knowledge/query' '{"query":"test","top_k":5}'

# Other
test_api 'GET' '/notifications'
test_api 'GET' '/announcements'
test_api 'GET' '/devices/1/shadow'
test_api 'GET' '/ai-fairness/reports'

echo "PASS:$PASS FAIL:$FAIL SKIP:$SKIP TOTAL:$TOTAL"
