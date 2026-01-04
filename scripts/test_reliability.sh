#!/bin/bash

# 系统可靠性功能测试脚本

set -e

# 配置
API_URL=${API_URL:-"http://localhost:8080"}
BACKUP_DIR=${BACKUP_DIR:-"./backups"}

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 统计变量
PASSED=0
FAILED=0

# 测试函数
test_case() {
    local name=$1
    local cmd=$2
    
    echo -n "测试: $name ... "
    
    if eval "$cmd" > /dev/null 2>&1; then
        echo -e "${GREEN}通过${NC}"
        ((PASSED++))
        return 0
    else
        echo -e "${RED}失败${NC}"
        ((FAILED++))
        return 1
    fi
}

# 检查环境
check_env() {
    echo "=== 环境检查 ==="
    
    test_case "API服务运行中" "curl -f -s ${API_URL}/health/live"
    test_case "PostgreSQL可连接" "PGPASSWORD=postgres psql -h localhost -U postgres -d ai_doc_library -c 'SELECT 1'"
    test_case "备份目录存在" "test -d ${BACKUP_DIR}"
    
    echo ""
}

# 测试健康检查
test_health_check() {
    echo "=== 健康检查测试 ==="
    
    # 存活探针
    test_case "存活探针端点" "curl -f -s ${API_URL}/health/live"
    
    # 就绪探针
    test_case "就绪探针端点" "curl -f -s ${API_URL}/health/ready"
    
    # 完整健康检查
    test_case "完整健康检查端点" "curl -f -s ${API_URL}/health"
    
    # 断路器状态
    test_case "断路器状态端点" "curl -f -s ${API_URL}/health/circuit-breakers"
    
    # 验证健康状态JSON格式
    test_case "健康状态返回有效JSON" "curl -s ${API_URL}/health | jq .status"
    
    echo ""
}

# 测试备份功能
test_backup() {
    echo "=== 备份功能测试 ==="
    
    # 创建备份目录
    mkdir -p "${BACKUP_DIR}/database"
    
    # 测试 PostgreSQL 客户端
    test_case "pg_dump 命令可用" "which pg_dump"
    test_case "psql 命令可用" "which psql"
    
    # 测试备份脚本
    test_case "备份脚本可执行" "test -x scripts/backup-script.sh"
    
    # 创建测试数据库备份
    test_case "创建数据库备份" "PGPASSWORD=postgres pg_dump -h localhost -U postgres -d ai_doc_library > /tmp/test_backup.sql && test -s /tmp/test_backup.sql"
    
    # 清理测试文件
    rm -f /tmp/test_backup.sql
    
    echo ""
}

# 测试Docker Compose高可用
test_docker_ha() {
    echo "=== Docker Compose高可用测试 ==="
    
    test_case "docker-compose.yml存在" "test -f docker-compose.yml"
    test_case "docker-compose包含backend2服务" "grep -q 'backend2:' docker-compose.yml"
    test_case "docker-compose包含健康检查" "grep -q 'healthcheck:' docker-compose.yml"
    test_case "docker-compose包含Nginx负载均衡" "grep -q 'nginx:' docker-compose.yml"
    
    echo ""
}

# 测试Kubernetes配置
test_k8s_ha() {
    echo "=== Kubernetes高可用配置测试 ==="
    
    test_case "K8s部署配置存在" "test -f k8s/deployment.yaml"
    test_case "K8s配置HPA" "grep -q 'HorizontalPodAutoscaler' k8s/deployment.yaml"
    test_case "K8s配置健康检查" "grep -q 'livenessProbe' k8s/deployment.yaml"
    test_case "K8s配置就绪探针" "grep -q 'readinessProbe' k8s/deployment.yaml"
    test_case "K8s配置PVC" "test -f k8s/pvc.yaml"
    test_case "K8s配置Secrets" "test -f k8s/secrets.yaml"
    test_case "K8s配置ConfigMap" "test -f k8s/configmap.yaml"
    test_case "K8s配置PostgreSQL" "test -f k8s/postgres.yaml"
    test_case "K8s配置Ingress" "test -f k8s/ingress.yaml"
    
    echo ""
}

# 测试服务代码
test_service_code() {
    echo "=== 服务代码测试 ==="
    
    test_case "健康检查服务存在" "test -f internal/service/health_service.go"
    test_case "断路器服务存在" "test -f internal/service/circuit_breaker.go"
    test_case "备份服务存在" "test -f internal/service/backup_service.go"
    test_case "PostgreSQL备份服务存在" "test -f internal/service/postgres_backup.go"
    test_case "健康检查Handler存在" "test -f internal/handler/health_handler.go"
    test_case "备份Handler存在" "test -f internal/handler/backup_handler.go"
    
    echo ""
}

# 测试断路器功能
test_circuit_breaker() {
    echo "=== 断路器功能测试 ==="
    
    # 需要服务正在运行才能测试
    if curl -f -s ${API_URL}/health/live > /dev/null 2>&1; then
        echo "  需要服务实际运行以测试断路器功能"
        echo "  请参考 docs/reliability_guide.md 中的断路器使用示例"
    else
        echo "  跳过（服务未运行）"
    fi
    
    echo ""
}

# 生成测试报告
generate_report() {
    echo ""
    echo "=== 测试报告 ==="
    echo "通过: ${PASSED}"
    echo "失败: ${FAILED}"
    echo "总计: $((PASSED + FAILED))"
    
    if [ ${FAILED} -eq 0 ]; then
        echo -e "${GREEN}所有测试通过！${NC}"
        return 0
    else
        echo -e "${RED}有 ${FAILED} 个测试失败${NC}"
        return 1
    fi
}

# 主函数
main() {
    echo "==================================="
    echo "  系统可靠性功能测试"
    echo "==================================="
    echo ""
    
    check_env
    test_health_check
    test_backup
    test_docker_ha
    test_k8s_ha
    test_service_code
    test_circuit_breaker
    
    generate_report
}

# 执行测试
main "$@"