@echo off
chcp 65001 >nul
title MDM Backend Service
echo ========================================
echo MDM Backend 启动中...
echo ========================================

REM 设置环境变量
set JWT_SECRET=your-secret-key-change-in-production
set DATABASE_URL=postgres://mdm_user:mdm_password@localhost:5432/mdm_db?sslmode=disable
set REDIS_URL=redis://localhost:6379
set MQTT_BROKER=mqtt://localhost:1883
set CORS_ALLOWED_ORIGINS=http://localhost:3000,http://127.0.0.1:3000
set PORT=8085

REM 启动后端
cd /d "%~dp0"
start "MDM-Backend" mdm-backend.exe

echo Backend 启动完成!
echo 访问 http://localhost:8085/health 检查状态
pause
