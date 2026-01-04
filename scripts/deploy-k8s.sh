#!/bin/bash

# Kubernetes高可用部署脚本

set -e

# 配置变量
NAMESPACE=${NAMESPACE:-"ai-doc"}
IMAGE_TAG=${IMAGE_TAG:-"latest"}
REGISTRY=${REGISTRY:-"local"}
DEPLOYMENT_MODE=${DEPLOYMENT_MODE:-"all"}  # all, minimal, monitoring

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

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

# 检查Kubernetes集群
check_cluster() {
    log_info "检查Kubernetes集群连接..."
    if ! kubectl cluster-info &> /dev/null; then
        log_error "无法连接到Kubernetes集群"
        exit 1
    fi
    log_info "Kubernetes集群连接正常"
}

# 创建命名空间
create_namespace() {
    log_info "创建命名空间: ${NAMESPACE}"
    kubectl create namespace ${NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -
}

# 部署Secrets
deploy_secrets() {
    log_info "部署Secrets配置..."
    kubectl apply -f k8s/secrets.yaml -n ${NAMESPACE}
}

# 部署ConfigMaps
deploy_configmaps() {
    log_info "部署ConfigMaps配置..."
    kubectl apply -f k8s/configmap.yaml -n ${NAMESPACE}
}

# 部署PVC
deploy_pvc() {
    log_info "部署持久化存储卷声明..."
    kubectl apply -f k8s/pvc.yaml -n ${NAMESPACE}
}

# 部署PostgreSQL
deploy_postgres() {
    log_info "部署PostgreSQL数据库..."
    kubectl apply -f k8s/postgres.yaml -n ${NAMESPACE}
    
    log_info "等待PostgreSQL就绪..."
    kubectl wait --for=condition=ready pod -l app=postgres -n ${NAMESPACE} --timeout=300s
}

# 构建并推送镜像
build_and_push_image() {
    log_info "构建Docker镜像..."
    docker build -t ${REGISTRY}/ai-doc-backend:${IMAGE_TAG} .
    
    if [ "${REGISTRY}" != "local" ]; then
        log_info "推送镜像到注册表..."
        docker push ${REGISTRY}/ai-doc-backend:${IMAGE_TAG}
    fi
}

# 部署应用
deploy_app() {
    log_info "部署AI文档库应用..."
    kubectl apply -f k8s/deployment.yaml -n ${NAMESPACE}
    
    log_info "等待应用Pod就绪..."
    kubectl wait --for=condition=ready pod -l app=ai-doc-backend -n ${NAMESPACE} --timeout=300s
}

# 部署Ingress
deploy_ingress() {
    log_info "部署Ingress配置..."
    kubectl apply -f k8s/ingress.yaml -n ${NAMESPACE}
}

# 设置自动备份
setup_auto_backup() {
    log_info "设置自动备份CronJob..."
    
    cat > k8s/backup-cronjob.yaml <<EOF
apiVersion: batch/v1
kind: CronJob
metadata:
  name: auto-backup
  namespace: ${NAMESPACE}
spec:
  schedule: "0 2 * * *"  # 每天凌晨2点执行
  successfulJobsHistoryLimit: 3
  failedJobsHistoryLimit: 3
  jobTemplate:
    spec:
      template:
        spec:
          restartPolicy: OnFailure
          containers:
          - name: backup
            image: postgres:15-alpine
            command:
            - /bin/bash
            - -c
            - |
              PGPASSWORD=\${DB_PASSWORD} pg_dump -h \${DB_HOST} -U \${DB_USER} -d \${DB_NAME} | \
              gzip > /backups/backup_\$(date +%Y%m%d_%H%M%S).sql.gz && \
              find /backups -name "backup_*.sql.gz" -mtime +7 -delete
            env:
            - name: DB_HOST
              valueFrom:
                configMapKeyRef:
                  name: ai-doc-config
                  key: db-host
            - name: DB_PORT
              valueFrom:
                configMapKeyRef:
                  name: ai-doc-config
                  key: db-port
            - name: DB_USER
              valueFrom:
                secretKeyRef:
                  name: ai-doc-secrets
                  key: db-user
            - name: DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: ai-doc-secrets
                  key: db-password
            - name: DB_NAME
              valueFrom:
                configMapKeyRef:
                  name: ai-doc-config
                  key: db-name
            volumeMounts:
            - name: backups
              mountPath: /backups
          volumes:
          - name: backups
            persistentVolumeClaim:
              claimName: ai-doc-backups-pvc
EOF
    
    kubectl apply -f k8s/backup-cronjob.yaml -n ${NAMESPACE}
    log_info "自动备份已配置：每天凌晨2点执行"
}

# 显示部署状态
show_status() {
    log_info "部署状态："
    echo ""
    kubectl get all -n ${NAMESPACE}
    echo ""
    kubectl get pvc -n ${NAMESPACE}
    echo ""
    
    log_info "服务访问信息："
    echo ""
    kubectl get svc -n ${NAMESPACE}
    echo ""
    
    if command -v kubectl &> /dev/null; then
        INGRESS_HOST=$(kubectl get ingress ai-doc-ingress -n ${NAMESPACE} -o jsonpath='{.status.loadBalancer.ingress[0].host}' 2>/dev/null || echo "未配置")
        INGRESS_IP=$(kubectl get ingress ai-doc-ingress -n ${NAMESPACE} -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "未配置")
        
        if [ "${INGRESS_HOST}" != "未配置" ]; then
            log_info "Ingress Host: ${INGRESS_HOST}"
        elif [ "${INGRESS_IP}" != "未配置" ]; then
            log_info "Ingress IP: ${INGRESS_IP}"
        else
            log_info "使用 port-forward 访问服务："
            echo "  kubectl port-forward -n ${NAMESPACE} svc/ai-doc-backend-service 8080:8080"
        fi
    fi
}

# 主函数
# 部署所有配置
deploy_all_configs() {
    log_info "使用完整部署文件部署所有配置..."
    kubectl apply -f k8s/k8s-all.yaml
    
    log_info "等待所有Pod就绪..."
    kubectl wait --for=condition=ready pod -l app=postgres -n ${NAMESPACE} --timeout=300s
    kubectl wait --for=condition=ready pod -l app=redis -n ${NAMESPACE} --timeout=300s
    kubectl wait --for=condition=ready pod -l app=ai-doc-backend -n ${NAMESPACE} --timeout=300s
}

# 部署最小化配置
deploy_minimal_configs() {
    log_info "使用最小化部署方式..."
    
    check_cluster
    create_namespace
    deploy_secrets
    deploy_configmaps
    deploy_pvc
    deploy_postgres
    
    # 如果需要构建镜像
    if [ "${BUILD_IMAGE}" = "true" ]; then
        build_and_push_image
    fi
    
    deploy_app
    setup_auto_backup
    
    show_status
}

# 部署监控配置
deploy_monitoring() {
    log_info "部署监控组件..."
    
    # 部署Prometheus
    kubectl apply -f k8s/prometheus.yaml -n ${NAMESPACE}
    
    # 部署Grafana
    kubectl apply -f k8s/grafana.yaml -n ${NAMESPACE}
    
    log_info "监控组件部署完成"
    log_info "Prometheus: http://<node-ip>:9090"
    log_info "Grafana: http://<node-ip>:3001 (admin/admin)"
}

# 显示部署状态
show_deployment_help() {
    log_info "=========================================="
    log_info "部署使用说明"
    log_info "=========================================="
    echo ""
    log_info "完整部署（包含所有组件）:"
    echo "  NAMESPACE=ai-doc ./scripts/deploy-k8s.sh"
    echo ""
    log_info "最小化部署（仅核心组件）:"
    echo "  DEPLOYMENT_MODE=minimal NAMESPACE=ai-doc ./scripts/deploy-k8s.sh"
    echo ""
    log_info "包含监控的部署:"
    echo "  DEPLOYMENT_MODE=monitoring NAMESPACE=ai-doc ./scripts/deploy-k8s.sh"
    echo ""
    log_info "使用私有镜像仓库:"
    echo "  REGISTRY=your-registry.com NAMESPACE=ai-doc ./scripts/deploy-k8s.sh"
    echo ""
    log_info "构建并推送镜像（需要BUILD_IMAGE=true）:"
    echo "  BUILD_IMAGE=true REGISTRY=your-registry.com NAMESPACE=ai-doc ./scripts/deploy-k8s.sh"
    echo ""
    log_info "卸载部署:"
    echo "  NAMESPACE=ai-doc kubectl delete -f k8s/k8s-all.yaml"
    echo ""
    log_info "查看日志:"
    echo "  kubectl logs -n ai-doc deployment/ai-doc-backend -f"
    echo ""
    log_info "查看服务状态:"
    echo "  kubectl get all -n ai-doc"
    echo ""
    log_info "=========================================="
}

# 主函数
main() {
    log_info "=========================================="
    log_info "AI技术文档库 - Kubernetes部署脚本"
    log_info "=========================================="
    echo ""
    log_info "部署模式: ${DEPLOYMENT_MODE}"
    log_info "命名空间: ${NAMESPACE}"
    log_info "镜像标签: ${IMAGE_TAG}"
    log_info "镜像仓库: ${REGISTRY}"
    echo ""
    
    # 检查是否显示帮助
    if [ "$1" = "help" ] || [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
        show_deployment_help
        exit 0
    fi
    
    # 根据部署模式选择部署方式
    case ${DEPLOYMENT_MODE} in
        "all")
            log_info "开始完整部署..."
            check_cluster
            deploy_all_configs
            ;;
        "minimal")
            log_info "开始最小化部署..."
            deploy_minimal_configs
            ;;
        "monitoring")
            log_info "开始监控部署..."
            deploy_minimal_configs
            deploy_monitoring
            ;;
        *)
            log_error "未知的部署模式: ${DEPLOYMENT_MODE}"
            log_info "请使用: all, minimal 或 monitoring"
            show_deployment_help
            exit 1
            ;;
    esac
    
    show_status
    
    log_info "=========================================="
    log_info "部署完成！"
    log_info "=========================================="
    log_warn "请修改 k8s/secrets.yaml 中的默认密码和密钥"
    log_warn "生产环境请使用TLS证书配置HTTPS"
    echo ""
}

# 执行主函数
main "$@"