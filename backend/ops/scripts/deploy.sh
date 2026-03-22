#!/bin/bash
# =============================================================================
# MDM 平台一键部署脚本
# 用途：生产环境一键部署（支持 Linux/macOS，Windows 请使用 PowerShell 版）
# 使用：bash scripts/deploy.sh [环境] [操作]
# 示例：bash scripts/deploy.sh prod up
#       bash scripts/deploy.sh prod down
#       bash scripts/deploy.sh prod restart
#       bash scripts/deploy.sh prod logs
#       bash scripts/deploy.sh prod migrate
# =============================================================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info()  { echo -e "${BLUE}[INFO]${NC} $1"; }
log_ok()    { echo -e "${GREEN}[OK]${NC} $1"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# 脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OPS_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_DIR="$(dirname "$OPS_DIR")"

# 环境参数
ENV="${1:-prod}"
ACTION="${2:-up}"

# 检查 Docker 和 Docker Compose
check_docker() {
    log_info "检查 Docker 环境..."
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose 未安装，请先安装 Docker Compose"
        exit 1
    fi
    DOCKER_VERSION=$(docker version --format '{{.Server.Version}}' 2>/dev/null || echo "0")
    if [ "$(printf '%s\n' "24.0.0" "$DOCKER_VERSION" | sort -V | head -n1)" != "24.0.0" ]; then
        log_warn "Docker 版本建议 >= 24.0.0，当前: $DOCKER_VERSION"
    fi
    log_ok "Docker 检查通过"
}

# 生成密钥
generate_secrets() {
    log_info "检查并生成必要密钥..."
    if [ -z "$JWT_SECRET" ]; then
        export JWT_SECRET=$(openssl rand -base64 32 2>/dev/null || head -c 32 /dev/urandom | base64)
        log_ok "JWT_SECRET 已生成"
    fi
    if [ -z "$POSTGRES_PASSWORD" ]; then
        export POSTGRES_PASSWORD=$(openssl rand -base64 24 2>/dev/null || head -c 24 /dev/urandom | base64)
        log_ok "POSTGRES_PASSWORD 已生成"
    fi
    if [ -z "$EMQX_ADMIN_PASSWORD" ]; then
        export EMQX_ADMIN_PASSWORD=$(openssl rand -base64 24 2>/dev/null || head -c 24 /dev/urandom | base64)
        log_ok "EMQX_ADMIN_PASSWORD 已生成"
    fi
}

# 检查 .env 文件
check_env_file() {
    local env_file="$OPS_DIR/.env"
    if [ ! -f "$env_file" ]; then
        log_warn ".env 文件不存在，创建默认配置..."
        generate_secrets
        cat > "$env_file" << EOF
# ================================================
# MDM 平台生产环境配置
# 生成时间: $(date '+%Y-%m-%d %H:%M:%S')
# ================================================

# --- 必须配置 ---
JWT_SECRET=$JWT_SECRET
POSTGRES_PASSWORD=$POSTGRES_PASSWORD
EMQX_ADMIN_PASSWORD=$EMQX_ADMIN_PASSWORD

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
EOF
        log_ok ".env 文件已创建: $env_file"
        log_warn "请编辑 .env 文件配置 CORS_ALLOWED_ORIGINS 和其他参数"
    else
        log_ok ".env 文件已存在"
    fi
}

# 创建数据目录
create_data_dirs() {
    log_info "创建数据目录..."
    mkdir -p "$OPS_DIR/data/postgres" \
             "$OPS_DIR/data/redis" \
             "$OPS_DIR/data/emqx" \
             "$OPS_DIR/data/nginx/logs" \
             "$OPS_DIR/data/apps"
    # 设置权限
    chmod -R 755 "$OPS_DIR/data"
    log_ok "数据目录创建完成"
}

# 拉取最新代码（仅 git 操作时）
git_pull() {
    if [ -d "$PROJECT_DIR/.git" ]; then
        log_info "拉取最新代码..."
        cd "$PROJECT_DIR"
        git pull origin main || git pull origin master
        cd "$OPS_DIR"
        log_ok "代码已更新"
    else
        log_warn "非 Git 仓库，跳过代码更新"
    fi
}

# 执行部署
do_deploy() {
    log_info "开始部署 [环境: $ENV]..."
    cd "$OPS_DIR"

    case "$ENV" in
        prod)
            docker-compose -f docker-compose.prod.yml up -d --build
            ;;
        dev)
            docker-compose up -d --build
            ;;
        *)
            log_error "未知环境: $ENV (仅支持 prod/dev)"
            exit 1
            ;;
    esac
    log_ok "部署完成"
}

# 查看服务状态
do_status() {
    cd "$OPS_DIR"
    echo ""
    echo "========================================"
    echo " MDM 服务状态"
    echo "========================================"
    docker-compose -f docker-compose.prod.yml ps
    echo ""
    echo "========================================"
    echo " 容器健康检查"
    echo "========================================"
    docker-compose -f docker-compose.prod.yml ps | grep -E "(Name|mdm-)" | awk '{print $1}' | tail -n +2 | while read container; do
        HEALTH=$(docker inspect --format='{{.State.Health.Status}}' "$container" 2>/dev/null || echo "none")
        STATUS=$(docker inspect --format='{{.State.Status}}' "$container" 2>/dev/null || echo "unknown")
        if [ "$HEALTH" = "healthy" ]; then
            echo -e "  $container: ${GREEN}healthy${NC} (running)"
        elif [ "$STATUS" = "running" ]; then
            echo -e "  $container: ${YELLOW}running${NC} (no health check)"
        else
            echo -e "  $container: ${RED}$STATUS${NC}"
        fi
    done
}

# 查看日志
do_logs() {
    cd "$OPS_DIR"
    local service="${3:-}"
    if [ -n "$service" ]; then
        docker-compose -f docker-compose.prod.yml logs -f --tail=100 "$service"
    else
        docker-compose -f docker-compose.prod.yml logs -f --tail=50
    fi
}

# 执行数据库迁移
do_migrate() {
    log_info "执行数据库迁移..."
    cd "$OPS_DIR"
    docker-compose -f docker-compose.prod.yml exec mdm-backend \
        sh -c "cd /app && go run scripts/migrate.go up" || \
    docker-compose -f docker-compose.prod.yml exec mdm-backend \
        sh -c "cd /app && echo 'Migrations handled by AutoMigrate on startup'"
    log_ok "迁移完成"
}

# 备份数据库
do_backup() {
    local backup_file="backup_$(date +%Y%m%d_%H%M%S).sql"
    log_info "备份数据库到: $backup_file"
    cd "$OPS_DIR"
    docker-compose -f docker-compose.prod.yml exec -T mdm-postgres \
        pg_dump -U mdm_user mdm_db > "$backup_file"
    log_ok "备份完成: $backup_file"
}

# 回滚
do_rollback() {
    log_warn "回滚操作将停止并删除当前容器"
    read -p "确认执行回滚? (yes/no): " confirm
    if [ "$confirm" != "yes" ]; then
        log_info "取消回滚"
        exit 0
    fi
    cd "$OPS_DIR"
    docker-compose -f docker-compose.prod.yml down
    # 回滚 Git 代码
    if [ -d "$PROJECT_DIR/.git" ]; then
        cd "$PROJECT_DIR"
        git reset --hard HEAD~1
        cd "$OPS_DIR"
    fi
    docker-compose -f docker-compose.prod.yml up -d --build
    log_ok "回滚完成"
}

# 清理
do_clean() {
    log_warn "清理操作将删除停止的容器、未使用的镜像和构建缓存"
    read -p "确认执行清理? (yes/no): " confirm
    if [ "$confirm" != "yes" ]; then
        log_info "取消清理"
        exit 0
    fi
    cd "$OPS_DIR"
    docker-compose -f docker-compose.prod.yml down --remove-orphans
    docker system prune -f --volumes
    log_ok "清理完成"
}

# 显示帮助
show_help() {
    cat << EOF
MDM 平台部署脚本

用法: bash deploy.sh [环境] [操作]

环境:
  prod    生产环境 (默认)
  dev     开发环境

操作:
  up          启动服务 (默认)
  down        停止服务
  restart     重启服务
  status      查看服务状态
  logs        查看日志 (可选: logs [服务名])
  migrate     执行数据库迁移
  backup      备份数据库
  rollback    回滚到上一版本
  clean       清理未使用资源
  pull        拉取最新代码并重启

示例:
  bash deploy.sh prod up           # 启动生产环境
  bash deploy.sh prod logs          # 查看所有日志
  bash deploy.sh prod logs backend # 查看后端日志
  bash deploy.sh prod migrate      # 执行迁移
  bash deploy.sh prod backup       # 备份数据库

前置条件:
  - Docker Engine 24.0+
  - Docker Compose v2.20+
  - .env 文件已配置 (首次运行自动创建)

EOF
}

# 主逻辑
main() {
    echo ""
    echo "========================================"
    echo " MDM 平台部署脚本"
    echo " 环境: $ENV | 操作: $ACTION"
    echo "========================================"
    echo ""

    case "$ACTION" in
        up|restart)
            check_docker
            check_env_file
            create_data_dirs
            git_pull
            do_deploy
            echo ""
            log_ok "服务已启动，请访问:"
            echo "  前端:        http://localhost:80"
            echo "  后端 API:    http://localhost:8080"
            echo "  EMQX Dash:   http://localhost:18083"
            ;;
        down)
            cd "$OPS_DIR"
            docker-compose -f docker-compose.prod.yml down
            log_ok "服务已停止"
            ;;
        status)
            do_status
            ;;
        logs)
            do_logs "$@"
            ;;
        migrate)
            check_docker
            do_migrate
            ;;
        backup)
            do_backup
            ;;
        rollback)
            do_rollback
            ;;
        clean)
            do_clean
            ;;
        pull)
            git_pull
            docker-compose -f docker-compose.prod.yml up -d --build
            log_ok "代码已更新并重启"
            ;;
        help|-h|--help)
            show_help
            ;;
        *)
            log_error "未知操作: $ACTION"
            show_help
            exit 1
            ;;
    esac
}

main "$@"
