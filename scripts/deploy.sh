#!/bin/bash

# AI文档库项目部署脚本
# 支持Docker Compose方式的完整部署

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置变量
PROJECT_NAME="ai-doc-library"
COMPOSE_FILE="docker-compose.yml"
BACKUP_DIR="./backups"
STORAGE_DIR="./storage"
LOGS_DIR="./logs"

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# 显示帮助信息
show_help() {
    cat << EOF
AI文档库项目部署脚本

用法: $0 [命令] [选项]

命令:
  deploy   部署项目（默认命令）
  start    启动服务
  stop     停止服务
  restart  重启服务
  status   查看服务状态
  logs     查看服务日志
  cleanup  清理所有数据和容器
  backup   备份数据库
  mirror   配置Docker镜像加速器
  help     显示帮助信息

选项:
  --no-build    跳过镜像构建
  --minimal     最小化部署（不启动监控服务）
  --with-monitor 启动监控服务（Prometheus + Grafana）

示例:
  $0 deploy               # 完整部署
  $0 deploy --minimal     # 最小化部署
  $0 start                # 启动服务
  $0 status               # 查看状态
  $0 logs backend         # 查看后端日志
  $0 mirror               # 配置Docker镜像加速器

网络问题排查:
  如果遇到镜像拉取失败，请使用以下命令配置镜像加速器:
  $0 mirror
  
  或手动配置 Docker 镜像加速器:
  1. 编辑 Docker 配置文件: /etc/docker/daemon.json
  2. 添加以下内容:
     {
       "registry-mirrors": [
         "https://docker.mirrors.ustc.edu.cn",
         "https://dockerhub.azk8s.cn",
         "https://hub-mirror.c.163.com"
       ]
     }
  3. 重启 Docker 服务: sudo systemctl restart docker

EOF
}

# 检查依赖
check_dependencies() {
    log_step "检查依赖..."
    
    # 检查Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker未安装，请先安装Docker"
        exit 1
    fi
    log_info "Docker版本: $(docker --version)"
    
    # 检查Docker Compose
    if ! docker compose version &> /dev/null; then
        log_error "Docker Compose未安装，请先安装Docker Compose"
        exit 1
    fi
    log_info "Docker Compose已就绪"
    
    # 检查端口占用
    check_ports
}

# 检查端口占用
check_ports() {
    log_step "检查端口占用..."
    
    PORTS=("80" "8080" "8081" "5432" "50051" "9000" "9001" "9090" "3001")
    OCCUPIED=()
    
    for port in "${PORTS[@]}"; do
        if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
            OCCUPIED+=("$port")
        fi
    done
    
    if [ ${#OCCUPIED[@]} -gt 0 ]; then
        log_warn "以下端口已被占用: ${OCCUPIED[*]}"
        read -p "是否继续? (y/n) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    else
        log_info "所有端口可用"
    fi
}

# 创建必要的目录
create_directories() {
    log_step "创建必要的目录..."
    
    if [ ! -d "$BACKUP_DIR" ]; then
        mkdir -p "$BACKUP_DIR"
        log_info "创建备份目录: $BACKUP_DIR"
    fi
    
    if [ ! -d "$STORAGE_DIR" ]; then
        mkdir -p "$STORAGE_DIR"
        log_info "创建存储目录: $STORAGE_DIR"
    fi
    
    if [ ! -d "$LOGS_DIR" ]; then
        mkdir -p "$LOGS_DIR/nginx"
        mkdir -p "$LOGS_DIR/backend"
        mkdir -p "$LOGS_DIR/backend2"
        log_info "创建日志目录: $LOGS_DIR"
    fi
    
    # 确保目录有正确的权限（忽略权限错误）
    chmod -R 755 "$BACKUP_DIR" "$STORAGE_DIR" "$LOGS_DIR" 2>/dev/null || true
    
    # 如果权限修改失败，尝试使用sudo
    if [ $? -ne 0 ]; then
        log_warn "部分目录权限修改失败，尝试使用sudo..."
        sudo chmod -R 755 "$BACKUP_DIR" "$STORAGE_DIR" "$LOGS_DIR" 2>/dev/null || \
            log_warn "权限修改失败，但这不影响脚本执行"
    fi
}

# 构建Docker镜像
build_images() {
    log_step "构建Docker镜像..."
    
    docker compose build
    log_info "镜像构建完成"
}

# 部署项目
deploy_project() {
    local skip_build=false
    local minimal=false
    local with_monitor=false
    
    # 解析参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            --no-build)
                skip_build=true
                shift
                ;;
            --minimal)
                minimal=true
                shift
                ;;
            --with-monitor)
                with_monitor=true
                shift
                ;;
            *)
                shift
                ;;
        esac
    done
    
    log_info "=========================================="
    log_info "开始部署AI文档库项目"
    log_info "=========================================="
    
    check_dependencies
    create_directories
    
    if [ "$skip_build" = false ]; then
        build_images
    fi
    
    log_step "启动服务..."
    
    # 根据参数选择部署模式
    if [ "$minimal" = true ]; then
        log_info "使用最小化部署模式"
        docker compose up -d postgres backend backend2 python-parser nginx
    elif [ "$with_monitor" = true ]; then
        log_info "使用完整部署模式（含监控）"
        docker compose --profile monitoring up -d
    else
        log_info "使用标准部署模式"
        docker compose up -d postgres backend backend2 python-parser nginx minio
    fi
    
    log_info "等待服务启动..."
    sleep 10
    
    show_status
    
    log_info "=========================================="
    log_info "部署完成！"
    log_info "=========================================="
    log_info "前端访问地址: http://localhost"
    log_info "后端API地址: http://localhost:8080"
    log_info "后端API地址2: http://localhost:8081"
    log_info "数据库端口: 5432"
    log_info "MinIO控制台: http://localhost:9001"
    log_info "=========================================="
    
    if [ "$with_monitor" = true ]; then
        log_info "监控系统:"
        log_info "  Prometheus: http://localhost:9090"
        log_info "  Grafana: http://localhost:3001 (admin/admin)"
        log_info "=========================================="
    fi
}

# 启动服务
start_services() {
    log_step "启动服务..."
    docker compose start
    show_status
}

# 停止服务
stop_services() {
    log_step "停止服务..."
    docker compose stop
    log_info "服务已停止"
}

# 重启服务
restart_services() {
    log_step "重启服务..."
    docker compose restart
    show_status
}

# 查看服务状态
show_status() {
    log_info "=========================================="
    log_info "服务状态"
    log_info "=========================================="
    docker compose ps
    
    echo ""
    log_info "容器资源使用情况:"
    docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}" $(docker compose ps -q)
}

# 查看日志
view_logs() {
    local service=${1:-""}
    
    if [ -z "$service" ]; then
        log_info "查看所有服务日志..."
        docker compose logs -f
    else
        log_info "查看 $service 服务日志..."
        docker compose logs -f "$service"
    fi
}

# 清理项目
cleanup_project() {
    log_warn "警告：此操作将删除所有容器、卷和数据！"
    read -p "确认继续? (yes/no): " confirm
    
    if [ "$confirm" != "yes" ]; then
        log_info "操作已取消"
        return
    fi
    
    log_step "清理项目..."
    
    # 停止并删除容器
    docker compose down -v
    
    # 删除数据目录（可选）
    read -p "是否删除数据目录 $(pwd)/data? (y/n): " delete_data
    if [ "$delete_data" = "y" ]; then
        rm -rf data
        log_info "数据目录已删除"
    fi
    
    log_info "清理完成"
}

# 备份数据库
backup_database() {
    log_step "备份数据库..."
    
    local backup_file="$BACKUP_DIR/backup_$(date +%Y%m%d_%H%M%S).sql.gz"
    
    log_info "正在备份数据库到: $backup_file"
    
    docker compose exec -T postgres pg_dump -U postgres ai_doc_library | gzip > "$backup_file"
    
    if [ $? -eq 0 ]; then
        log_info "数据库备份成功: $backup_file"
        
        # 清理7天前的旧备份
        find "$BACKUP_DIR" -name "backup_*.sql.gz" -mtime +7 -delete
        log_info "已清理7天前的旧备份"
    else
        log_error "数据库备份失败"
        exit 1
    fi
}

# 配置Docker镜像加速器
configure_docker_mirror() {
    log_step "配置Docker镜像加速器..."
    
    local docker_config="/etc/docker/daemon.json"
    local backup_config="/etc/docker/daemon.json.backup"
    
    # 检查是否为root用户
    if [ "$EUID" -ne 0 ]; then
        log_error "此操作需要root权限，请使用 sudo"
        return 1
    fi
    
    # 备份现有配置
    if [ -f "$docker_config" ]; then
        log_info "备份现有配置到 $backup_config"
        cp "$docker_config" "$backup_config"
    fi
    
    # 创建新的配置
    log_info "创建新的Docker配置..."
    cat > "$docker_config" << 'EOF'
{
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn",
    "https://dockerhub.azk8s.cn",
    "https://hub-mirror.c.163.com",
    "https://mirror.ccs.tencentyun.com"
  ],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m",
    "max-file": "3"
  }
}
EOF
    
    # 重启Docker服务
    log_info "重启Docker服务..."
    systemctl daemon-reload
    systemctl restart docker
    
    # 验证配置
    if systemctl is-active --quiet docker; then
        log_info "Docker服务重启成功"
        log_info "镜像加速器配置已完成"
        docker info | grep -A 5 "Registry Mirrors"
    else
        log_error "Docker服务重启失败"
        # 恢复备份
        if [ -f "$backup_config" ]; then
            log_warn "恢复原有配置..."
            cp "$backup_config" "$docker_config"
            systemctl daemon-reload
            systemctl restart docker
        fi
        return 1
    fi
    
    log_info "建议清理Docker缓存以使用新的镜像源："
    echo "  docker system prune -a"
}

# 显示系统信息
show_system_info() {
    log_info "=========================================="
    log_info "系统信息"
    log_info "=========================================="
    log_info "操作系统: $(uname -s)"
    log_info "内核版本: $(uname -r)"
    log_info "Docker版本: $(docker --version)"
    log_info "Docker Compose版本: $(docker compose version)"
    
    echo ""
    log_info "Docker镜像加速器配置:"
    docker info 2>/dev/null | grep -A 5 "Registry Mirrors" || echo "未配置镜像加速器"
    
    echo ""
    log_info "Docker磁盘使用:"
    docker system df
    
    echo ""
    log_info "活跃容器:"
    docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
}

# 主函数
main() {
    local command=${1:-"deploy"}
    shift || true
    
    case "$command" in
        deploy)
            deploy_project "$@"
            ;;
        start)
            start_services
            ;;
        stop)
            stop_services
            ;;
        restart)
            restart_services
            ;;
        status)
            show_status
            ;;
        logs)
            view_logs "$@"
            ;;
        cleanup)
            cleanup_project
            ;;
        backup)
            backup_database
            ;;
        mirror)
            configure_docker_mirror
            ;;
        info)
            show_system_info
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            log_error "未知命令: $command"
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"