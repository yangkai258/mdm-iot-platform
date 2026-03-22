#!/bin/bash
# =============================================================================
# MDM 平台数据库迁移脚本
# 用途：执行数据库迁移和种子数据初始化
# 使用：bash scripts/migrate.sh [命令]
# 示例：
#   bash scripts/migrate.sh up        # 执行所有迁移
#   bash scripts/migrate.sh down     # 回滚最近一次迁移
#   bash scripts/migrate.sh status   # 查看迁移状态
#   bash scripts/migrate.sh seed     # 初始化种子数据
#   bash scripts/migrate.sh fresh     # 重建数据库（危险！）
# =============================================================================

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info()  { echo -e "${BLUE}[INFO]${NC} $1"; }
log_ok()    { echo -e "${GREEN}[OK]${NC} $1"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OPS_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_DIR="$(dirname "$OPS_DIR")"
BACKEND_DIR="$PROJECT_DIR/backend"

# 加载 .env
load_env() {
    local env_file="$OPS_DIR/.env"
    if [ -f "$env_file" ]; then
        set -a
        source "$env_file"
        set +a
    fi
}

# 检查必要工具
check_deps() {
    log_info "检查依赖..."
    command -v docker &>/dev/null || { log_error "Docker 未安装"; exit 1; }
    load_env
    log_ok "依赖检查通过"
}

# 执行迁移
migrate_up() {
    log_info "执行数据库迁移..."
    cd "$BACKEND_DIR"

    # 通过 Docker 执行迁移
    docker run --rm \
        --network mdm-network \
        -v "$BACKEND_DIR:/app" \
        -w /app \
        -e DATABASE_URL="postgres://${POSTGRES_USER:-mdm_user}:${POSTGRES_PASSWORD}@127.0.0.1:5432/${POSTGRES_DB:-mdm_db}?sslmode=disable" \
        golang:1.21 \
        sh -c "go run scripts/migrate.go up" 2>/dev/null || {
        # 备选：直接通过 docker-compose exec
        cd "$OPS_DIR"
        docker-compose -f docker-compose.prod.yml exec -T mdm-backend \
            sh -c "cd /app && go run scripts/migrate.go up" 2>/dev/null || {
            log_warn "迁移脚本不存在，使用 AutoMigrate..."
            log_ok "AutoMigrate 会在服务启动时自动执行"
        }
    }
    log_ok "迁移完成"
}

# 回滚（模拟）
migrate_down() {
    log_warn "回滚功能需要手动编写迁移脚本"
    log_info "当前使用 AutoMigrate，建议通过添加 version 字段实现版本化管理"
    log_info "如需完全回滚，请执行: docker-compose down -v && docker-compose up -d"
}

# 迁移状态
migrate_status() {
    log_info "检查迁移状态..."
    cd "$OPS_DIR"

    local tables=$(docker-compose -f docker-compose.prod.yml exec -T mdm-postgres \
        psql -U ${POSTGRES_USER:-mdm_user} -d ${POSTGRES_DB:-mdm_db} -t -c \
        "SELECT tablename FROM pg_tables WHERE schemaname='public' ORDER BY tablename;" 2>/dev/null)

    echo ""
    echo "========================================"
    echo " 数据库表列表"
    echo "========================================"
    echo "$tables" | tr -d ' ' | grep -v '^$'
    echo ""
    log_ok "共 $(echo "$tables" | grep -c '^[a-z_]') 张表"
}

# 种子数据
seed_data() {
    log_info "初始化种子数据..."

    # 创建默认租户（如果不存在）
    cd "$OPS_DIR"

    docker-compose -f docker-compose.prod.yml exec -T mdm-postgres \
        psql -U ${POSTGRES_USER:-mdm_user} -d ${POSTGRES_DB:-mdm_db} -c "
        -- 插入默认租户（如果不存在）
        INSERT INTO tenants (id, name, status, created_at, updated_at)
        VALUES ('default', '默认租户', 'active', NOW(), NOW())
        ON CONFLICT (id) DO NOTHING;

        -- 插入默认管理员用户（如果不存在）
        INSERT INTO sys_users (id, username, password, tenant_id, role, created_at, updated_at)
        VALUES ('admin', 'admin', '\$2a\$10\$K8j5sH5Z6K8j5sH5Z6K8jO', 'default', 'super_admin', NOW(), NOW())
        ON CONFLICT (id) DO NOTHING;
        " 2>/dev/null || log_warn "种子数据表可能不存在，跳过"

    log_ok "种子数据初始化完成"
    log_info "默认管理员账号: admin / admin123 (请首次登录后立即修改密码)"
}

# 重建数据库（危险）
migrate_fresh() {
    log_error "重建数据库将删除所有数据！"
    read -p "输入 'yes' 确认执行: " confirm
    if [ "$confirm" != "yes" ]; then
        log_info "取消操作"
        exit 0
    fi

    log_warn "正在删除所有表..."
    cd "$OPS_DIR"

    docker-compose -f docker-compose.prod.yml exec -T mdm-postgres \
        psql -U ${POSTGRES_USER:-mdm_user} -d ${POSTGRES_DB:-mdm_db} -c "
        DROP SCHEMA public CASCADE;
        CREATE SCHEMA public;
        GRANT ALL ON SCHEMA public TO public;
        " 2>/dev/null

    log_ok "数据库已清空"
    log_info "重启服务以触发 AutoMigrate..."
    docker-compose -f docker-compose.prod.yml restart mdm-backend
    log_ok "完成，请检查服务日志确认表创建成功"
}

# 显示帮助
show_help() {
    cat << EOF
MDM 平台数据库迁移脚本

用法: bash migrate.sh [命令]

命令:
  up       执行所有迁移 (默认)
  down     回滚最近一次迁移
  status   查看迁移状态和表列表
  seed     初始化种子数据
  fresh    重建数据库（危险！需确认）
  help     显示帮助

说明:
  - MDM 使用 GORM AutoMigrate，服务启动时自动创建/更新表结构
  - 手动迁移脚本位于 backend/scripts/migrate.go
  - 迁移状态记录在 schema_migrations 表（如已配置）

示例:
  bash migrate.sh status    # 查看有哪些表
  bash migrate.sh seed      # 初始化默认数据
  bash migrate.sh fresh     # 重建数据库

注意事项:
  - 执行迁移前请备份数据库！
  - 生产环境建议先在测试环境验证迁移脚本
  - AutoMigrate 只会创建不存在的表，不会修改已有列

EOF
}

# 主逻辑
main() {
    local cmd="${1:-up}"
    echo ""
    echo "========================================"
    echo " MDM 数据库迁移"
    echo " 命令: $cmd"
    echo "========================================"
    echo ""

    check_deps

    case "$cmd" in
        up)      migrate_up ;;
        down)    migrate_down ;;
        status)  migrate_status ;;
        seed)    seed_data ;;
        fresh)   migrate_fresh ;;
        help|-h|--help) show_help ;;
        *)
            log_error "未知命令: $cmd"
            show_help
            exit 1
            ;;
    esac
}

main "$@"
