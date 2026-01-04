# 系统可靠性功能实现指南

## 1. 概述

本文档描述AI技术文档库系统的可靠性功能实现，包括故障恢复、数据备份恢复机制和高可用性部署。

## 2. 故障恢复功能

### 2.1 健康检查机制

系统实现了多层次的健康检查机制：

#### 组件健康检查
- **数据库健康检查** (`DatabaseHealthCheck`)：定期检查数据库连接状态
- **存储健康检查** (`StorageHealthCheck`)：验证存储服务可用性

#### 探针接口
- **存活探针** (`/health/live`)：Kubernetes使用，检查容器是否存活
- **就绪探针** (`/health/ready`)：检查服务是否准备好接收流量
- **完整健康检查** (`/health`)：返回所有组件的详细健康状态

#### 健康状态
```json
{
  "status": "healthy",
  "timestamp": "2026-01-03T14:00:00Z",
  "components": [
    {
      "name": "database",
      "status": "healthy",
      "message": "数据库连接正常",
      "timestamp": "2026-01-03T14:00:00Z"
    },
    {
      "name": "storage",
      "status": "healthy",
      "message": "存储服务正常",
      "timestamp": "2026-01-03T14:00:00Z"
    }
  ],
  "uptime": "2h30m45s"
}
```

### 2.2 断路器机制

系统实现了断路器模式（Circuit Breaker Pattern）来防止级联故障：

#### 断路器状态
- **关闭**：正常请求状态
- **打开**：失败次数超过阈值，阻止请求
- **半开**：测试服务是否恢复

#### 使用示例
```go
// 获取断路器
cb := healthService.GetCircuitBreaker("database", 5, 30*time.Second)

// 执行带断路器保护的操作
err := cb.Execute(func() error {
    return db.Ping()
})

if err != nil {
    log.Println("操作失败或断路器打开")
}
```

## 3. 数据备份和恢复机制

### 3.1 备份服务架构

- **备份服务接口** (`BackupService`)：定义备份和恢复的统一接口
- **PostgreSQL备份** (`PostgreSQLBackup`)：实现数据库备份和恢复
- **备份管理**：通过API接口管理备份任务

### 3.2 备份类型

#### 数据库备份
```bash
# 手动触发数据库备份
curl -X POST "http://localhost:8080/api/v1/backup/create?type=database" \
  -H "Authorization: Bearer <token>"
```

#### 存储备份
```bash
# 手动触发存储备份
curl -X POST "http://localhost:8080/api/v1/backup/create?type=storage" \
  -H "Authorization: Bearer <token>"
```

#### 完整备份
```bash
# 创建完整备份（数据库+存储）
curl -X POST "http://localhost:8080/api/v1/backup/create?type=full" \
  -H "Authorization: Bearer <token>"
```

### 3.3 备份管理API

#### 获取备份列表
```bash
curl -X GET "http://localhost:8080/api/v1/backup/list" \
  -H "Authorization: Bearer <token>"
```

响应示例：
```json
{
  "count": 5,
  "backups": [
    {
      "id": "20260103_140000",
      "type": "full",
      "path": "/app/backups/full/20260103_140000",
      "size": 1073741824,
      "created_at": "2026-01-03T14:00:00Z"
    }
  ]
}
```

#### 恢复备份
```bash
curl -X POST "http://localhost:8080/api/v1/backup/restore/20260103_140000" \
  -H "Authorization: Bearer <token>"
```

#### 删除备份
```bash
curl -X DELETE "http://localhost:8080/api/v1/backup/20260103_140000" \
  -H "Authorization: Bearer <token>"
```

### 3.4 自动备份脚本

系统提供了自动化备份脚本 `scripts/backup-script.sh`：

```bash
# 数据库备份
./scripts/backup-script.sh backup

# 完整备份
./scripts/backup-script.sh full

# 恢复备份
./scripts/backup-script.sh restore /app/backups/database/backup_20260103_140000.sql.gz

# 清理旧备份
./scripts/backup-script.sh cleanup
```

#### 配置环境变量
```bash
# 备份目录
export BACKUP_DIR="/app/backups"
export DB_HOST="postgres"
export DB_PORT="5432"
export DB_USER="postgres"
export DB_PASSWORD="postgres"
export DB_NAME="ai_doc_library"
export RETENTION_DAYS=7  # 保留天数
```

## 4. 高可用性部署

### 4.1 Docker Compose高可用部署

#### 启动多实例部署
```bash
# 启动主实例 + 备用实例（2个后端服务）
docker-compose up -d

# 验证服务状态
docker-compose ps
```

#### 配置说明
- 多个后端实例通过nginx负载均衡
- PostgreSQL使用持久化存储
- 存储目录在多个实例间共享

### 4.2 Kubernetes高可用部署

#### 部署架构
```
                    +
                    | Ingress
                    v
            +---------------+
            |   Nginx LB    |
            +---------------+
                    |
        +-----------+-----------+
        |                       |
        v                       v
+--------------+      +--------------+
|  Backend-1   |      |  Backend-2   |
|  (Replica 1) |      |  (Replica 2) |
+--------------+      +--------------+
        |                       |
        +-----------+-----------+
                    |
                    v
            +---------------+
            |  PostgreSQL   |
            | (StatefulSet) |
            +---------------+
```

#### 部署步骤

1. **准备镜像**
```bash
# 构建镜像
docker build -t ai-doc-backend:latest .

# 或推送到镜像仓库
docker tag ai-doc-backend:latest registry.example.com/ai-doc-backend:latest
docker push registry.example.com/ai-doc-backend:latest
```

2. **配置Secrets**
编辑 `k8s/secrets.yaml`，修改默认密码：
```yaml
stringData:
  db-user: "postgres"
  db-password: "<强密码>"
  jwt-secret: "<强密钥>"
```

3. **部署到Kubernetes**
```bash
# 赋予脚本执行权限
chmod +x scripts/deploy-k8s.sh

# 执行部署
./scripts/deploy-k8s.sh
```

4. **验证部署**
```bash
# 查看Pod状态
kubectl get pods -n ai-doc

# 查看服务
kubectl get svc -n ai-doc

# 查看Ingress
kubectl get ingress -n ai-doc
```

#### 水平自动扩缩容 (HPA)
系统配置了HPA，可根据CPU和内存使用率自动扩缩容：

```yaml
minReplicas: 2
maxReplicas: 10
metrics:
- type: Resource
  resource:
    name: cpu
    target:
      averageUtilization: 70
- type: Resource
  resource:
    name: memory
    target:
      averageUtilization: 80
```

手动扩缩容：
```bash
# 扩容到5个副本
kubectl scale deployment ai-doc-backend --replicas=5 -n ai-doc

# 查看HPA状态
kubectl get hpa -n ai-doc
```

### 4.3 自动备份配置

Kubernetes部署包含自动备份CronJob，每天凌晨2点自动执行：

```bash
# 查看CronJob
kubectl get cronjob -n ai-doc

# 查看备份历史
kubectl get jobs -n ai-doc

# 手动触发备份
kubectl create job --from=cronjob/auto-backup manual-backup-$(date +%s) -n ai-doc
```

## 5. 故障恢复场景

### 5.1 数据库故障恢复

#### 场景1：数据库连接失败
1. 健康检查检测到数据库不可用
2. 断路器打开，阻止后续请求
3. 恢复后，自动重试连接
4. 断路器半开测试连接成功后关闭

#### 场景2：数据库损坏
1. 使用最新备份恢复数据库
```bash
./scripts/backup-script.sh restore /app/backups/database/backup_20260103_140000.sql.gz
```

2. 或通过API恢复
```bash
curl -X POST "http://localhost:8080/api/v1/backup/restore/20260103_140000" \
  -H "Authorization: Bearer <token>"
```

### 5.2 服务实例故障

#### Kubernetes自动恢复
- Pod崩溃时自动重启
- 不健康Pod自动替换
- 服务自动重新路由流量

```bash
# 模拟Pod故障
kubectl delete pod -l app=ai-doc-backend -n ai-doc

# 观察自动恢复
kubectl get pods -w -n ai-doc
```

### 5.3 存储故障

#### 本地存储故障
1. 检查存储健康状态
```bash
curl -X GET "http://localhost:8080/health"
```

2. 如存储不可用，使用分布式存储（S3/MinIO）
3. 配置环境变量切换存储类型
```bash
export STORAGE_TYPE=minio
export MINIO_ENDPOINT=minio:9000
export MINIO_ACCESS_KEY=minioadmin
export MINIO_SECRET_KEY=minioadmin123
export MINIO_BUCKET=ai-doc
```

## 6. 监控和告警

### 6.1 健康监控指标

系统通过健康检查API提供以下指标：
- 组件状态（健康/不健康/降级）
- 系统运行时间
- 依赖服务状态

### 6.2 断路器监控

```bash
# 获取断路器状态
curl -X GET "http://localhost:8080/health/circuit-breakers"
```

### 6.3 日志监控

关键事件日志：
- 健康检查失败
- 断路器状态变化
- 备份操作
- 恢复操作

## 7. 最佳实践

### 7.1 备份策略
- **定期备份**：每天凌晨自动备份
- **多地备份**：将备份复制到异地存储
- **备份验证**：定期测试备份恢复流程
- **保留策略**：保留最近7天的备份

### 7.2 高可用配置
- **多实例部署**：至少2个后端实例
- **负载均衡**：使用Nginx或Kubernetes Service
- **健康检查**：配置合适的健康检查间隔
- **资源限制**：合理设置CPU和内存限制

### 7.3 故障演练
定期进行故障演练：
1. 模拟数据库故障
2. 模拟服务实例故障
3. 模拟存储故障
4. 演练备份恢复流程

## 8. 故障排查

### 8.1 常见问题

#### 问题1：健康检查失败
```bash
# 查看健康状态
curl -X GET "http://localhost:8080/health"

# 查看Pod日志
kubectl logs <pod-name> -n ai-doc

# 检查依赖服务
kubectl get pods -n ai-doc
```

#### 问题2：备份失败
```bash
# 查看备份脚本日志
./scripts/backup-script.sh backup 2>&1 | tee backup.log

# 检查磁盘空间
df -h

# 检查数据库连接
PGPASSWORD=password psql -h postgres -U postgres -d ai_doc_library -c "SELECT 1"
```

#### 问题3：恢复失败
```bash
# 检查备份文件完整性
gunzip -t /app/backups/database/backup_20260103_140000.sql.gz

# 检查备份内容
gunzip -c /app/backups/database/backup_20260103_140000.sql.gz | head -n 50

# 尝试恢复并查看详细错误
PGPASSWORD=password psql -h postgres -U postgres -d ai_doc_library \
  -f /app/backups/database/backup_20260103_140000.sql
```

## 9. API参考

### 9.1 健康检查API

| 端点 | 方法 | 描述 |
|------|------|------|
| `/health` | GET | 完整健康检查 |
| `/health/live` | GET | 存活探针 |
| `/health/ready` | GET | 就绪探针 |
| `/health/circuit-breakers` | GET | 断路器状态 |

### 9.2 备份管理API

| 端点 | 方法 | 描述 |
|------|------|------|
| `/api/v1/backup/create?type={type}` | POST | 创建备份 |
| `/api/v1/backup/list` | GET | 获取备份列表 |
| `/api/v1/backup/restore/{backupId}` | POST | 恢复备份 |
| `/api/v1/backup/{backupId}` | DELETE | 删除备份 |

## 10. 总结

系统可靠性功能实现了：
- ✅ 故障恢复：健康检查、断路器机制
- ✅ 数据备份恢复：自动备份、手动备份、备份管理API
- ✅ 高可用性部署：Docker Compose多实例、Kubernetes集群部署、自动扩缩容

通过这些功能，系统能够：
- 快速检测和响应故障
- 自动恢复服务实例
- 保护数据安全和完整性
- 在高负载下稳定运行