@echo off
chcp 65001 >nul
title MDM Frontend (Vite)
echo ========================================
echo MDM Frontend 启动中...
echo ========================================

cd /d "%~dp0..\frontend"
npm run dev

pause
