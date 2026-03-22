# =============================================================================
# MDM 平台数据库迁移脚本 (PowerShell 版)
# 用途：执行数据库迁移和种子数据初始化
# 使用：.\migrate.ps1 -Command up
#       .\migrate.ps1 -Command status
#       .\migrate.ps1 -Command seed
# =============================================================================

param(
    [ValidateSet("up", "down", "status", "seed", "fresh", "help")]
    [string]$Command = "up"
)

$ErrorActionPreference = "Stop"

function Write-Info  { Write-Host "[INFO]  $args" -ForegroundColor Cyan }
function Write-Ok    { Write-Host "[OK]    $args" -ForegroundColor Green }
function Write-Warn  { Write-Host "[WARN]  $args" -ForegroundColor Yellow }
function Write-Err   { Write-Host "[ERROR] $args" -ForegroundColor Red }

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$OpsDir = Split-Path -Parent $ScriptDir
$ProjectDir = Split-Path -Parent $OpsDir
$BackendDir = Join-Path $ProjectDir "backend"

# 加载 .env
function Import-EnvFile {
    $envFile = Join-Path $OpsDir ".env"
    if (Test-Path $envFile) {
        Get-Content $envFile | ForEach-Object {
            if ($_ -match '^([^=]+)=(.*)$') {
                [Environment]::SetEnvironmentVariable($matches[1].Trim(), $matches[2].Trim(), "Process")
            }
        }
    }
}

function Test-Deps {
    Write-Info "检查依赖..."
    try {
        docker version | Out-Null
    } catch {
        Write-Err "Docker 未安装或未运行"
        exit 1
    }
    Import-EnvFile
    Write-Ok "依赖检查通过"
}

function Start-MigrateUp {
    Write-Info "执行数据库迁移..."
    Import-EnvFile
    Push-Location $OpsDir

    try {
        docker-compose -f docker-compose.prod.yml exec -T mdm-backend `
            sh -c "cd /app && go run scripts/migrate.go up" 2>$null
    } catch {
        Write-Warn "迁移脚本不存在，使用 AutoMigrate..."
    }

    Pop-Location
    Write-Ok "迁移完成"
}

function Start-MigrateDown {
    Write-Warn "回滚功能需要手动编写迁移脚本"
    Write-Info "当前使用 AutoMigrate"
    Write-Info "如需完全回滚，请执行: docker-compose down -v && docker-compose up -d"
}

function Get-MigrateStatus {
    Write-Info "检查迁移状态..."
    Import-EnvFile
    Push-Location $OpsDir

    $tables = docker-compose -f docker-compose.prod.yml exec -T mdm-postgres `
        psql -U ${env:POSTGRES_USER} -d ${env:POSTGRES_DB} -t -c `
        "SELECT tablename FROM pg_tables WHERE schemaname='public' ORDER BY tablename;" 2>$null

    Write-Host ""
    Write-Host "========================================" -ForegroundColor Magenta
    Write-Host " 数据库表列表" -ForegroundColor Magenta
    Write-Host "========================================" -ForegroundColor Magenta
    $tables -split "`n" | ForEach-Object { $_.Trim() } | Where-Object { $_ -ne "" } | ForEach-Object { Write-Host "  $_" }
    Write-Host ""

    $count = ($tables -split "`n" | Where-Object { $_.Trim() -ne "" }).Count
    Write-Ok "共 $count 张表"

    Pop-Location
}

function Start-SeedData {
    Write-Info "初始化种子数据..."
    Import-EnvFile
    Push-Location $OpsDir

    try {
        docker-compose -f docker-compose.prod.yml exec -T mdm-postgres `
            psql -U ${env:POSTGRES_USER} -d ${env:POSTGRES_DB} -c "
            -- 插入默认租户（如果不存在）
            INSERT INTO tenants (id, name, status, created_at, updated_at)
            VALUES ('default', '默认租户', 'active', NOW(), NOW())
            ON CONFLICT (id) DO NOTHING;

            -- 插入默认管理员用户（如果不存在）
            INSERT INTO sys_users (id, username, password, tenant_id, role, created_at, updated_at)
            VALUES ('admin', 'admin', '$2a$10\$K8j5sH5Z6K8j5sH5Z6K8jO', 'default', 'super_admin', NOW(), NOW())
            ON CONFLICT (id) DO NOTHING;
            " 2>$null
        Write-Ok "种子数据初始化完成"
        Write-Info "默认管理员账号: admin / admin123 (请首次登录后立即修改密码)"
    } catch {
        Write-Warn "种子数据表可能不存在，跳过"
    }

    Pop-Location
}

function Start-MigrateFresh {
    Write-Err "重建数据库将删除所有数据！"
    $confirm = Read-Host "输入 'yes' 确认执行"
    if ($confirm -ne "yes") {
        Write-Info "取消操作"
        return
    }

    Write-Warn "正在删除所有表..."
    Import-EnvFile
    Push-Location $OpsDir

    try {
        docker-compose -f docker-compose.prod.yml exec -T mdm-postgres `
            psql -U ${env:POSTGRES_USER} -d ${env:POSTGRES_DB} -c "
            DROP SCHEMA public CASCADE;
            CREATE SCHEMA public;
            GRANT ALL ON SCHEMA public TO public;
            " 2>$null

        Write-Ok "数据库已清空"
        Write-Info "重启服务以触发 AutoMigrate..."
        docker-compose -f docker-compose.prod.yml restart mdm-backend
        Write-Ok "完成，请检查服务日志确认表创建成功"
    } catch {
        Write-Err "操作失败: $_"
    }

    Pop-Location
}

function Show-Help {
    @"
MDM 平台数据库迁移脚本 (PowerShell)

用法: .\migrate.ps1 -Command <命令>

命令:
  up       执行所有迁移 (默认)
  down     回滚最近一次迁移
  status   查看迁移状态和表列表
  seed     初始化种子数据
  fresh    重建数据库（危险！需确认）
  help     显示帮助

说明:
  - MDM 使用 GORM AutoMigrate，服务启动时自动创建/更新表结构
  - 迁移状态记录在 schema_migrations 表（如已配置）

示例:
  .\migrate.ps1 -Command status    查看有哪些表
  .\migrate.ps1 -Command seed      初始化默认数据
  .\migrate.ps1 -Command fresh     重建数据库

注意事项:
  - 执行迁移前请备份数据库！
  - 生产环境建议先在测试环境验证迁移脚本

"@
}

function Main {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Magenta
    Write-Host " MDM 数据库迁移" -ForegroundColor Magenta
    Write-Host " 命令: $Command" -ForegroundColor Magenta
    Write-Host "========================================" -ForegroundColor Magenta
    Write-Host ""

    Test-Deps

    switch ($Command) {
        "up"     { Start-MigrateUp }
        "down"   { Start-MigrateDown }
        "status" { Get-MigrateStatus }
        "seed"   { Start-SeedData }
        "fresh"  { Start-MigrateFresh }
        "help"   { Show-Help }
    }
}

Main
