#!/bin/sh
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6ImFkbWluIiwicm9sZV9pZCI6MCwidGVuYW50X2lkIjoiIiwiaXNfc3VwZXJfYWRtaW4iOmZhbHNlLCJleHAiOjE3NzQzNjc5NzAsImlhdCI6MTc3NDI4MTU3MH0.URkp53kuWRIPL5U6D2dLYaXt4ylJuT0EXubWNDXfVE8"
AUTH="Authorization: Bearer $TOKEN"

echo "=== Testing POST /orders ==="
curl -s -w "\nHTTP_CODE:%{http_code}\n" -X POST "http://localhost:8080/api/v1/orders" -H "$AUTH" -H "Content-Type: application/json" -d '{"total_amount":100}' 

echo ""
echo "=== Testing POST /knowledge/query ==="
curl -s -w "\nHTTP_CODE:%{http_code}\n" -X POST "http://localhost:8080/api/v1/knowledge/query" -H "$AUTH" -H "Content-Type: application/json" -d '{"query":"test","top_k":5}'
