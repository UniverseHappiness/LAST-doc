#!/bin/bash

# AI技术文档库 - 部署配置验证脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 计数器
TOTAL_CHECKS=0
PASSED_CHECKS=0
FAILED_CHECKS=0

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[✓]${NC} $1"
    ((PASSED_CHECKS++))
    ((TOTAL_CHECKS++))
}

log_error() {
    echo -e "${RED}[✗]${NC} $1"
    ((FAILED_CHECKS++))
    ((TOTAL_CHECKS++))
}

log_warn() {
    echo -e "${YELLOW}[!]${NC} $1"
}

log_section() {
    echo ""
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo ""
}

# 检查文件是否存在
check_file_exists() {
    local file=$1
    if [ -f "$file" ]; then
        log_success "文件存在: $file"
        return 0
    else
        log_error "文件不存在: $file"
        return 1
    fi
}

# 检查目录是否存在
check_dir_exists() {
    local dir=$1
    if [ -d "$dir" ]; then
        log_success "目录存在: $dir"
        return 0
    else
        log_error "目录不存在: $dir"
        return 1
    fi
}

# 检查文件是否可读
check_file_readable() {
    local file=$1
    if [ -r "$file" ]; then
        log_success "文件可读: $file"
        return 0
    else
        log_error "文件不可读: $file"
        return 1
    fi
}

# 检查文件语法
check_dockerfile_syntax() {
    local dockerfile=$1
    if grep -q "^FROM" "$dockerfile" && grep -q "^RUN" "$dockerfile" && grep -q "^CMD" "$dockerfile"; then
        log_success "Dockerfile语法正确: $dockerfile"
        return 0
    else
        log_error "Dockerfile语法错误: $dockerfile"
        return 1
    fi
}

# 检查YAML语法
check_yaml_syntax() {
    local yaml_file=$1
    if command -v python3 &> /dev/null; then
        if python3 -c "import yaml; yaml.safe_load(open('$yaml_file'))" 2>/dev/null; then
            log_success "YAML语法正确: $yaml_file"
            return 0
        else
            log_error "YAML语法错误: $yaml_file"
            return 1
        fi
    else
        log_warn "未安装Python3，跳过YAML语法检查: $yaml_file"
        return 0
    fi
}

# 检查Shell脚本语法
check_shell_syntax() {
    local script=$1
    if bash -n "$script" 2>/dev/null; then
        log_success "Shell脚本语法正确: $script"
        return 0
    else
        log_error "Shell脚本语法错误: $script"
        return 1
    fi
}

# 主验证函数
main() {
    echo -e "${GREEN}"
    echo "=========================================="
    echo "AI技术文档库 - 部署配置验证"
    echo "=========================================="
    echo -e "${NC}"
    echo ""

    # 1. Docker配置验证
    log_section "1. Docker配置验证"
    check_file_exists "Dockerfile"
    check_file_exists "docker-compose.yml"
    check_file_readable "Dockerfile"
    check_file_readable "docker-compose.yml"
    check_dockerfile_syntax "Dockerfile"

    # 2. Kubernetes配置验证
    log_section "2. Kubernetes配置验证"
    check_dir_exists "k8s"
    
    # 检查各个Kubernetes配置文件
    check_file_exists "k8s/namespace.yaml"
    check_file_exists "k8s/deployment.yaml"
    check_file_exists "k8s/service.yaml"
    check_file_exists "k8s/configmap.yaml"
    check_file_exists "k8s/secrets.yaml"
    check_file_exists "k8s/pvc.yaml"
    check_file_exists "k8s/postgres.yaml"
    check_file_exists "k8s/redis.yaml"
    check_file_exists "k8s/ingress.yaml"
    check_file_exists "k8s/k8s-all.yaml"

    # 检查YAML语法
    for yaml_file in k8s/*.yaml; do
        check_yaml_syntax "$yaml_file"
    done

    # 3. Nginx配置验证
    log_section "3. Nginx配置验证"
    check_file_exists "nginx.conf"
    check_file_readable "nginx.conf"
    if command -v nginx &> /dev/null; then
        if nginx -t -c /dev/stdin 2>/dev/null < nginx.conf; then
            log_success "Nginx配置语法正确: nginx.conf"
        else
            log_warn "Nginx配置验证需要完整配置，跳过"
        fi
    else
        log_warn "未安装Nginx，跳过Nginx配置验证"
    fi

    # 4. 脚本验证
    log_section "4. 脚本验证"
    check_file_exists "scripts/deploy-k8s.sh"
    check_file_exists "scripts/backup-script.sh"
    check_shell_syntax "scripts/deploy-k8s.sh"
    if [ -f "scripts/backup-script.sh" ]; then
        check_shell_syntax "scripts/backup-script.sh"
    fi

    # 5. 配置内容验证
    log_section "5. 配置内容验证"

    # 检查Docker Compose配置
    if grep -q "version:" docker-compose.yml; then
        log_success "Docker Compose版本配置正确"
    else
        log_error "Docker Compose缺少版本配置"
    fi

    if grep -q "services:" docker-compose.yml; then
        log_success "Docker Compose服务配置正确"
    else
        log_error "Docker Compose缺少服务配置"
    fi

    # 检查必备服务
    for service in nginx backend postgres; do
        if grep -q "$service:" docker-compose.yml; then
            log_success "Docker Compose包含服务: $service"
        else
            log_error "Docker Compose缺少服务: $service"
        fi
    done

    # 检查Kubernetes配置
    if grep -q "apiVersion:" k8s/deployment.yaml; then
        log_success "Kubernetes Deployment配置正确"
    else
        log_error "Kubernetes Deployment配置错误"
    fi

    if grep -q "kind: ConfigMap" k8s/configmap.yaml; then
        log_success "Kubernetes ConfigMap配置正确"
    else
        log_error "Kubernetes ConfigMap配置错误"
    fi

    if grep -q "kind: Secret" k8s/secrets.yaml; then
        log_success "Kubernetes Secret配置正确"
    else
        log_error "Kubernetes Secret配置错误"
    fi

    # 6. 文档验证
    log_section "6. 文档验证"
    check_file_exists "DEPLOYMENT.md"
    check_file_exists "README.md"
    check_file_exists "docs/QUICK_START.md"
    check_file_exists "docs/PRIVATE_DEPLOYMENT.md"

    # 7. 安全配置验证
    log_section "7. 安全配置验证"
    
    # 检查是否有默认密码警告
    if grep -q "your-secret-key-change-in-production" k8s/secrets.yaml; then
        log_warn "生产环境请修改默认密钥配置"
    fi

    if grep -q "your-secret-key-change-in-production" docker-compose.yml; then
        log_warn "生产环境请修改默认密钥配置"
    fi

    # 检查健康检查配置
    if grep -q "healthcheck:" docker-compose.yml; then
        log_success "Docker Compose健康检查配置正确"
    else
        log_warn "Docker Compose缺少健康检查配置"
    fi

    # 8. 网络配置验证
    log_section "8. 网络配置验证"
    
    if grep -q "networks:" docker-compose.yml; then
        log_success "Docker Compose网络配置正确"
    else
        log_error "Docker Compose缺少网络配置"
    fi

    if grep -q "kind: Service" k8s/deployment.yaml || grep -q "kind: Service" k8s/k8s-all.yaml; then
        log_success "Kubernetes Service配置正确"
    else
        log_error "Kubernetes缺少Service配置"
    fi

    # 9. 存储配置验证
    log_section "9. 存储配置验证"
    
    if grep -q "volumes:" docker-compose.yml; then
        log_success "Docker Compose存储配置正确"
    else
        log_warn "Docker Compose缺少存储配置"
    fi

    if grep -q "kind: PersistentVolumeClaim" k8s/pvc.yaml; then
        log_success "Kubernetes PVC配置正确"
    else
        log_error "Kubernetes缺少PVC配置"
    fi

    # 10. 部署脚本验证
    log_section "10. 部署脚本验证"
    
    if grep -q "NAMESPACE" scripts/deploy-k8s.sh; then
        log_success "部署脚本包含NAMESPACE配置"
    else
        log_warn "部署脚本缺少NAMESPACE配置"
    fi

    if grep -q "kubectl apply" scripts/deploy-k8s.sh; then
        log_success "部署脚本包含kubectl apply命令"
    else
        log_error "部署脚本缺少kubectl apply命令"
    fi

    # 11. Makefile验证
    log_section "11. Makefile验证"
    check_file_exists "Makefile"
    if grep -q "docker-build" Makefile; then
        log_success "Makefile包含docker-build目标"
    fi

    # 显示验证结果
    log_section "验证结果汇总"
    echo ""
    echo -e "总检查项: ${TOTAL_CHECKS}"
    echo -e "${GREEN}通过: ${PASSED_CHECKS}${NC}"
    echo -e "${RED}失败: ${FAILED_CHECKS}${NC}"
    echo ""

    if [ $FAILED_CHECKS -eq 0 ]; then
        echo -e "${GREEN}==========================================${NC}"
        echo -e "${GREEN}✓ 所有配置验证通过！${NC}"
        echo -e "${GREEN}==========================================${NC}"
        echo ""
        echo "下一步操作："
        echo "1. Docker Compose部署: docker-compose up -d"
        echo "2. Kubernetes部署: ./scripts/deploy-k8s.sh"
        echo "3. 查看快速开始文档: docs/QUICK_START.md"
        return 0
    else
        echo -e "${RED}==========================================${NC}"
        echo -e "${RED}✗ 验证失败，请检查上述错误项${NC}"
        echo -e "${RED}==========================================${NC}"
        echo ""
        return 1
    fi
}

# 执行主函数
main "$@"