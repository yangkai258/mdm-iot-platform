# =============================================================================
# MDM 平台一键部署脚本 (PowerShell 版)
# 用途：生产环境一键部署（Windows 专用）
# 使用：.\deploy.ps1 -Env prod -Action up
#       .\deploy.ps1 -Env prod -Action down
#       .\deploy.ps1 -Env prod -Action restart
#       .\deploy.ps1 -Env prod -Action logs
#       .\deploy.ps1 -Env prod -Action migrate
# =============================================================================

param(
    [ValidateSet("prod", "dev")]
    [string]$Env = "prod",

    [ValidateSet("up", "down", "restart", "status", "logs", "migrate", "backup", "rollback", "clean", "pull", "help")]
    [string]$Action = "up",

    [string]$Service = ""
)

$ErrorActionPreference = "Stop"

# 脚本所在目录
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$OpsDir = Split-Path -Parent $ScriptDir
$ProjectDir = Split-Path -Parent $OpsDir

# 颜色函数 (PowerShell 5+)
function Write-Info  { Write-Host "[INFO]  $args" -ForegroundColor Cyan }
function Write-Ok    { Write-Host "[OK]    $args" -ForegroundColor Green }
function Write-Warn  { Write-Host "[WARN]  $args" -ForegroundColor Yellow }
function Write-Err   { Write-Host "[ERROR] $args" -ForegroundColor Red }

# Docker 检查
function Test-Docker {
    Write-Info "检查 Docker 环境..."
    try {
        $dockerVersion = docker version --format '{{.Server.Version}}' 2>$null
        if (-not $dockerVersion) { throw "Docker not running" }
        Write-Ok "Docker 检查通过 (版本: $dockerVersion)"
    } catch {
        Write-Err "Docker 未安装或未运行，请先安装 Docker Desktop"
        exit 1
    }
}

# .env 文件检查
function Test-EnvFile {
    $envFile = Join-Path $OpsDir ".env"
    if (-not (Test-Path $envFile)) {
        Write-Warn ".env 文件不存在，创建默认配置..."

        # 生成随机密钥
        $jwtSecret = [Convert]::ToBase64String((1..32 | ForEach-Object { Get-Random -Maximum 256 }) )
        $postgresPwd = [Convert]::ToBase64String((1..24 | ForEach-Object { Get-Random -Maximum 256 }) )
        $emqxPwd = [Convert]::ToBase64String((1..24 | ForEach-Object { Get-Random -Maximum 256 }) )

        $envContent = @"
# ================================================
# MDM 平台生产环境配置
# 生成时间: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')
# ================================================

# --- 必须配置 ---
`$JWT_SECRET=$jwtSecret
POSTGRES_PASSWORD=$postgresPwd
EMQX_ADMIN_PASSWORD=$emqxPwd

# --- CORS 白名单（生产环境必须配置具体域名）---
CORS_ALLOWED_ORIGINS=https://mdm.yourdomain.com

# --- 数据库 ---
POSTGRES_USER=mdm_user
POSTGRES_DB=mdm_db

# --- EMQX ---
EMQX_HOST=emqx
EMQX_ADMIN_USER=admin

# --- API Base URL ---
VITE_API_BASE_URL=https://api.mdm.yourdomain.com

# --- 日志级别 ---
GIN_MODE=release
LOG_LEVEL=info
LOG_FORMAT=json

# --- 告警通知（可选）---
# SMTP_HOST=smtp.example.com
# SMTP_PORT=587
# SMTP_USER=alerts@example.com
# SMTP_PASSWORD=your-smtp-password
# SMTP_FROM=noreply@mdm.example.com
# SMTP_USE_TLS=true
# ALERT_ADMIN_EMAIL=admin@example.com
# WEBHOOK_URL=https://hooks.example.com/alert
# WEBHOOK_TOKEN=your-webhook-secret
"@
        Set-Content -Path $envFile -Value $envContent -Encoding UTF8
        Write-Ok ".env 文件已创建: $envFile"
        Write-Warn "请编辑 .env 文件配置 CORS_ALLOWED_ORIGINS 和其他参数"
    } else {
        Write-Ok ".env 文件已存在"
    }
}

# 创建数据目录
function New-DataDirs {
    Write-Info "创建数据目录..."
    $dirs = @(
        "data\postgres",
        "data\redis",
        "data\emqx",
        "data\nginx\logs",
        "data\apps"
    )
    foreach ($dir in $dirs) {
        $path = Join-Path $OpsDir $dir
        if (-not (Test-Path $path)) {
            New-Item -ItemType Directory -Force -Path $path | Out-Null
        }
    }
    Write-Ok "数据目录创建完成"
}

# Git Pull
function Update-GitCode {
    $gitDir = Join-Path $ProjectDir ".git"
    if (Test-Path $gitDir) {
        Write-Info "拉取最新代码..."
        Push-Location $ProjectDir
        git pull origin main 2>$null -or (git pull origin master 2>$null)
        Pop-Location
        Write-Ok "代码已更新"
    }
}

# 部署
function Start-Deploy {
    Write-Info "开始部署 [环境: $Env]..."
    Push-Location $OpsDir

    $composeFile = if ($Env -eq "prod") { "docker-compose.prod.yml" } else { "docker-compose.yml" }
    docker-compose -f $composeFile up -d --build

    Pop-Location
    Write-Ok "部署完成"
}

# 服务状态
function Get-ServiceStatus {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Magenta
    Write-Host " MDM 服务状态" -ForegroundColor Magenta
    Write-Host "========================================" -ForegroundColor Magenta

    Push-Location $OpsDir
    docker-compose -f docker-compose.prod.yml ps
    Pop-Location

    Write-Host ""
    Write-Host "========================================" -ForegroundColor Magenta
    Write-Host " 容器健康检查" -ForegroundColor Magenta
    Write-Host "========================================" -ForegroundColor Magenta

    $containers = docker-compose -f (Join-Path $OpsDir "docker-compose.prod.yml") ps -q
    foreach ($containerId in $containers) {
        $name = docker inspect --format='{{.Name}}' $containerId 2>$null | ForEach-Object { $_.TrimStart('/') }
        $health = docker inspect --format='{{.State.Health.Status}}' $containerId 2>$null
        $status = docker inspect --format='{{.State.Status}}' $containerId 2>$null

        if ($health -eq "healthy") {
            Write-Host "  $name : running (healthy)" -ForegroundColor Green
        } elseif ($status -eq "running") {
            Write-Host "  $name : running" -ForegroundColor Yellow
        } else {
            Write-Host "  $name : $status" -ForegroundColor Red
        }
    }
}

# 日志
function Show-Logs {
    Push-Location $OpsDir
    if ($Service) {
        docker-compose -f docker-compose.prod.yml logs -f --tail=100 $Service
    } else {
        docker-compose -f docker-compose.prod.yml logs -f --tail=50
    }
    Pop-Location
}

# 迁移
function Start-Migrate {
    Write-Info "执行数据库迁移..."
    Push-Location $OpsDir
    docker-compose -f docker-compose.prod.yml exec mdm-backend sh -c "cd /app && go run scripts/migrate.go up" 2>$null
    if (-not $?) {
        docker-compose -f docker-compose.prod.yml exec mdm-backend sh -c "echo 'Migrations handled by AutoMigrate on startup'"
    }
    Pop-Location
    Write-Ok "迁移完成"
}

# 备份
function Start-Backup {
    $timestamp = Get-Date -Format 'yyyyMMdd_HHmmss'
    $backupFile = "backup_$timestamp.sql"
    Write-Info "备份数据库到: $backupFile"
    Push-Location $OpsDir
    docker-compose -f docker-compose.prod.yml exec -T mdm-postgres pg_dump -U mdm_user mdm_db | Set-Content $backupFile -Encoding UTF8
    Pop-Location
    Write-Ok "备份完成: $backupFile"
}

# 回滚
function Start-Rollback {
    Write-Warn "回滚操作将停止并删除当前容器"
    $confirm = Read-Host "确认执行回滚? (yes/no)"
    if ($confirm -ne "yes") {
        Write-Info "取消回滚"
        return
    }

    Push-Location $OpsDir
    docker-compose -f docker-compose.prod.yml down

    # 回滚 Git 代码
    $gitDir = Join-Path $ProjectDir ".git"
    if (Test-Path $gitDir) {
        Push-Location $ProjectDir
        git reset --hard HEAD~1
        Pop-Location
    }

    docker-compose -f docker-compose.prod.yml up -d --build
    Pop-Location
    Write-Ok "回滚完成"
}

# 清理
function Clear-Docker {
    Write-Warn "清理操作将删除停止的容器、未使用的镜像和构建缓存"
    $confirm = Read-Host "确认执行清理? (yes/no)"
    if ($confirm -ne "yes") {
        Write-Info "取消清理"
        return
    }

    Push-Location $OpsDir
    docker-compose -f docker-compose.prod.yml down --remove-orphans
    docker system prune -f
    Pop-Location
    Write-Ok "清理完成"
}

# 帮助
function Show-Help {
    @"
MDM 平台部署脚本 (PowerShell)

用法: .\deploy.ps1 -Env prod -Action up

参数:
  -Env     环境: prod (默认) | dev
  -Action  操作: up | down | restart | status | logs | migrate | backup | rollback | clean | pull | help
  -Service 服务名 (可选，用于 logs)

示例:
  .\deploy.ps1 -Env prod -Action up           启动生产环境
  .\deploy.ps1 -Env prod -Action status       查看服务状态
  .\deploy.ps1 -Env prod -Action logs         查看所有日志
  .\deploy.ps1 -Env prod -Action logs -Service mdm-backend  查看后端日志
  .\deploy.ps1 -Env prod -Action migrate      执行迁移
  .\deploy.ps1 -Env prod -Action backup       备份数据库

前置条件:
  - Docker Desktop 已安装并运行
  - .env 文件已配置 (首次运行自动创建)

"@
}

# 主逻辑
function Main {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Magenta
    Write-Host " MDM 平台部署脚本" -ForegroundColor Magenta
    Write-Host " 环境: $Env | 操作: $Action" -ForegroundColor Magenta
    Write-Host "========================================" -ForegroundColor Magenta
    Write-Host ""

    switch ($Action) {
        "up" {
            Test-Docker
            Test-EnvFile
            New-DataDirs
            Update-GitCode
            Start-Deploy
            Write-Host ""
            Write-Ok "服务已启动，请访问:"
            Write-Host "  前端:        http://localhost:80"
            Write-Host "  后端 API:    http://localhost:8080"
            Write-Host "  EMQX Dash:   http://localhost:18083"
        }
        "down" {
            Push-Location $OpsDir
            docker-compose -f docker-compose.prod.yml down
            Pop-Location
            Write-Ok "服务已停止"
        }
        "restart" {
            Test-Docker
            Push-Location $OpsDir
            docker-compose -f docker-compose.prod.yml down
            docker-compose -f docker-compose.prod.yml up -d --build
            Pop-Location
            Write-Ok "重启完成"
        }
        "status" { Get-ServiceStatus }
        "logs"   { Show-Logs }
        "migrate" {
            Test-Docker
            Start-Migrate
        }
        "backup" { Start-Backup }
        "rollback" { Start-Rollback }
        "clean"   { Clear-Docker }
        "pull" {
            Update-GitCode
            Push-Location $OpsDir
            docker-compose -f docker-compose.prod.yml up -d --build
            Pop-Location
            Write-Ok "代码已更新并重启"
        }
        "help" { Show-Help }
    }
}

Main
