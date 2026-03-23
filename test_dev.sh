#!/bin/sh
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6ImFkbWluIiwicm9sZV9pZCI6MCwidGVuYW50X2lkIjoiIiwiaXNfc3VwZXJfYWRtaW4iOmZhbHNlLCJleHAiOjE3NzQzNjc5MjMsImlhdCI6MTc0MjgxNTIzfQ.D-DH-mjMVp4dO7w5p5JNEqWRJcdRzPkgiXz0Wx2_BDI"
curl -v http://localhost:8080/api/v1/devices -H "Authorization: Bearer $TOKEN" 2>&1
