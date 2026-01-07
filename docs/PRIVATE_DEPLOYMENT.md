# AI技术文档库 - 私有化部署配置指南

## 1. 私有化部署概述

### 1.1 适用场景

本私有化部署方案适用于以下场景：

- **企业内网环境**：无法访问外部网络的企业内部环境
- **数据安全要求高**：需要严格控制数据访问和传输
- **合规性要求**：需要满足特定行业的数据存储和处理规范
- **离线部署**：需要在完全离线的环境中部署和运行
- **自主可控**：需要完全掌控系统的运行和维护

### 1.2 私有化部署优势

- **数据安全**：所有数据存储在企业内部，不经过第三方
- **网络独立**：不依赖外部网络，安全可控
- **自主维护**：企业可以自主进行系统维护和升级
- **成本可控**：一次性部署，长期使用，成本可控
- **灵活定制**：可以根据企业需求进行定制化开发

## 2. 内网环境要求

### 2.1 硬件要求

| 组件 | 最小配置 | 推荐配置 | 说明               |
| ---- | -------- | -------- | ------------------ |
| CPU  | 4核      | 8核+     | 支持文档解析和搜索 |
| 内存 | 8GB      | 16GB+    | 支持并发访问和缓存 |
| 存储 | 100GB    | 500GB+   | 存储文档和数据库   |
| 网络 | 1Gbps    | 10Gbps   | 支持文档上传下载   |

### 2.2 软件要求

```
操作系统:
- Linux: CentOS 7+, Ubuntu 18.04+, Debian 10+
- Windows: Windows Server 2016+ (需安装WSL2或Docker Desktop)

运行时环境:
- Docker: 20.10+
- Docker Compose: 1.29+
- 或 Kubernetes: 1.20+ (K8s部署)

数据库:
- PostgreSQL: 15+ (推荐)
- 支持pgvector扩展（向量搜索）

可选组件:
- MinIO: 分布式对象存储
- Redis: 缓存服务
- Prometheus: 监控服务
- Grafana: 可视化工具
```

### 2.3 网络要求

```bash
# 必需端口
- 80: HTTP访问端口
- 443: HTTPS访问端口（可选）
- 8080: 后端API端口
- 8081: 后端备用端口（高可用）
- 5432: PostgreSQL数据库端口

# 可选端口
- 9000: MinIO API端口
- 9001: MinIO控制台端口
- 6379: Redis端口
- 9090: Prometheus端口
- 3001: Grafana端口
- 50051: gRPC解析服务端口

# 内网环境网络配置
- 确保所有服务在同一内网段
- 配置防火墙规则，只开放必要端口
- 配置负载均衡器（如Nginx、HAProxy）
```

## 3. 离线部署方案

### 3.1 准备离线安装包

```bash
# 1. 在有网络的环境中下载Docker镜像
docker pull postgres:15-alpine
docker pull nginx:alpine
docker pull redis:7-alpine
docker pull minio/minio:latest
docker pull prom/prometheus:latest
docker pull grafana/grafana:latest

# 2. 保存镜像为tar文件
docker save -o ai-doc-images.tar \
  postgres:15-alpine \
  nginx:alpine \
  redis:7-alpine \
  minio/minio:latest \
  prom/prometheus:latest \
  grafana/grafana:latest

# 3. 打包整个项目目录
tar -czf ai-doc-library.tar.gz \
  ai-doc-library/

# 4. 将文件传输到内网服务器
# 使用U盘、光盘、或安全的方式传输
```

### 3.2 在内网环境加载镜像

```bash
# 1. 加载Docker镜像
docker load -i ai-doc-images.tar

# 2. 解压项目文件
tar -xzf ai-doc-library.tar.gz

# 3. 进入项目目录
cd ai-doc-library

# 4. 构建应用镜像
docker build -t ai-doc-backend:latest .
```

### 3.3 使用Docker Compose离线部署

```bash
# 1. 使用本地镜像离线部署
# 修改docker-compose.yml，确保使用本地镜像

# 2. 启动所有服务
docker-compose up -d

# 3. 查看服务状态
docker-compose ps

# 4. 查看日志
docker-compose logs -f
```

## 4. 私有镜像仓库配置

### 4.1 部署Harbor私有镜像仓库

```bash
# 使用Harbor部署私有镜像仓库
# 下载Harbor离线安装包
wget https://github.com/goharbor/harbor/releases/download/v2.8.0/harbor-offline-installer-v2.8.0.tgz

# 解压安装包
tar -xzf harbor-offline-installer-v2.8.0.tgz
cd harbor

# 配置Harbor
cp harbor.yml.tmpl harbor.yml
vi harbor.yml

# 修改以下配置
# hostname: harbor.yourcompany.com
# http.port: 80
# harbor_admin_password: YourStrongPassword
# data_volume: /data/harbor

# 安装Harbor
./install.sh

# 启动Harbor
docker-compose up -d
```

### 4.2 推送镜像到私有仓库

```bash
# 1. 登录私有镜像仓库
docker login harbor.yourcompany.com
# 用户名: admin
# 密码: YourStrongPassword

# 2. 标记镜像
docker tag ai-doc-backend:latest \
  harbor.yourcompany.com/ai-doc/ai-doc-backend:v1.0.0

# 3. 推送镜像
docker push harbor.yourcompany.com/ai-doc/ai-doc-backend:v1.0.0
```

### 4.3 配置Kubernetes使用私有仓库

```bash
# 1. 创建Docker registry secret
kubectl create secret docker-registry regcred \
  --docker-server=harbor.yourcompany.com \
  --docker-username=admin \
  --docker-password=YourStrongPassword \
  --docker-email=admin@yourcompany.com \
  -n ai-doc

# 2. 在Deployment中使用私有镜像
# 修改k8s/deployment.yaml中的imagePullSecrets
spec:
  template:
    spec:
      imagePullSecrets:
      - name: regcred
      containers:
      - name: backend
        image: harbor.yourcompany.com/ai-doc/ai-doc-backend:v1.0.0
```

## 5. 安全配置

### 5.1 生产环境密码配置

```bash
# 1. 生成强密码
# 使用openssl生成随机密码
openssl rand -base64 32

# 或使用pwgen工具
apt-get install pwgen
pwgen -s 32 1

# 2. 更新k8s/secrets.yaml中的密码
vi k8s/secrets.yaml

# 修改以下字段为强密码
apiVersion: v1
kind: Secret
metadata:
  name: ai-doc-secrets
  namespace: ai-doc
type: Opaque
stringData:
  db-password: "your-generated-strong-password-here"
  jwt-secret: "another-generated-strong-secret-here"
  postgres-replication-password: "replication-password-here"
```

### 5.2 TLS/SSL证书配置

```bash
# 1. 生成自签名证书（内网环境）
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout tls.key -out tls.crt \
  -subj "/CN=ai-doc.internal/O=YourCompany/C=CN"

# 2. 创建Kubernetes TLS secret
kubectl create secret tls ai-doc-tls-secret \
  --cert=tls.crt \
  --key=tls.key \
  -n ai-doc

# 3. 配置Ingress使用TLS
# 取消并修改k8s/ingress.yaml中的TLS配置注释
    tls:
    - hosts:
      - ai-doc.internal
      secretName: ai-doc-tls-secret
```

### 5.3 网络安全配置

```bash
# 1. 配置NetworkPolicy
cat > k8s/network-policy.yaml <<EOF
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: ai-doc-network-policy
  namespace: ai-doc
spec:
  podSelector:
    matchLabels:
      app: ai-doc-backend
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector: {}
    ports:
    - protocol: TCP
      port: 8080
  - from:
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
    ports:
    - protocol: TCP
      port: 8080
  egress:
  - to:
    - namespaceSelector: {}
      podSelector:
        matchLabels:
          app: postgres
    ports:
    - protocol: TCP
      port: 5432
  - to:
    - namespaceSelector: {}
      podSelector:
        matchLabels:
          app: redis
    ports:
    - protocol: TCP
      port: 6379
  - to:
    - namespaceSelector: {}
    ports:
    - protocol: TCP
      port: 53
  - to:
    - namespaceSelector: {}
    ports:
    - protocol: TCP
      port: 443
EOF

kubectl apply -f k8s/network-policy.yaml
```

## 6. 数据安全

### 6.1 数据库加密

```sql
-- 1. 启用PostgreSQL TDE（透明数据加密）
-- 需要编译PostgreSQL时启用TDE支持

-- 2. 使用pgcrypto扩展进行数据加密
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 3. 加密敏感数据示例
INSERT INTO api_keys (
    key_value, 
    key_name, 
    encrypted_description
) VALUES (
    'api-key-123',
    'Production Key',
    pgp_sym_encrypt('This is a description', 'encryption-key')
);

-- 4. 解密数据
SELECT 
    key_name,
    pgp_sym_decrypt(encrypted_description::bytea, 'encryption-key') as description
FROM api_keys;
```

### 6.2 文件存储加密

```bash
# 1. 创建加密的存储卷
# 使用LUKS加密Linux分区
cryptsetup -y -v luksFormat /dev/sdb1
cryptsetup luksOpen /dev/sdb1 encrypted_storage
mkfs.ext4 /dev/mapper/encrypted_storage
mount /dev/mapper/encrypted_storage /app/storage

# 2. 或使用eCryptfs加密目录
apt-get install ecryptfs-utils
mount -t ecryptfs /app/storage /app/storage

# 3. 配置自动挂载
vi /etc/fstab
# 添加:
# /dev/mapper/encrypted_storage /app/storage ext4 defaults 0 0
```

### 6.3 备份加密

```bash
# 1. 加密备份文件
# 修改scripts/backup-script.sh，在备份后加密
PGPASSWORD=$DB_PASSWORD pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME | \
  gzip | \
  openssl enc -aes-256-cbc -salt -pass pass:$BACKUP_ENCRYPTION_KEY \
  -out /backups/backup_$(date +%Y%m%d_%H%M%S).sql.gz.enc

# 2. 解密备份文件
openssl enc -d -aes-256-cbc -in backup.enc \
  -pass pass:$BACKUP_ENCRYPTION_KEY | \
  gunzip | \
  psql -h localhost -U postgres -d ai_doc_library
```

## 7. 网络配置

### 7.1 内网DNS配置

```bash
# 1. 配置内网DNS
# 编辑/etc/hosts或使用内部DNS服务器
vi /etc/hosts

# 添加以下记录
192.168.1.100  ai-doc.internal
192.168.1.100  api.ai-doc.internal
192.168.1.100  mcp.ai-doc.internal

# 或配置PowerDNS、BIND等DNS服务器
```

### 7.2 负载均衡配置

```nginx
# Nginx负载均衡配置
upstream ai_doc_backend {
    # 轮询算法
    server 192.168.1.101:8080 weight=1 max_fails=3 fail_timeout=30s;
    server 192.168.1.102:8080 weight=1 max_fails=3 fail_timeout=30s;
    server 192.168.1.103:8080 weight=1 max_fails=3 fail_timeout=30s;
  
    # 最少连接算法（推荐）
    least_conn;
  
    keepalive 32;
    keepalive_timeout 60s;
}

server {
    listen 80;
    server_name ai-doc.internal;
  
    # 前端静态文件
    location / {
        root /usr/share/nginx/html;
        try_files $uri $uri/ /index.html;
    }
  
    # API代理
    location /api/ {
        proxy_pass http://ai_doc_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
      
        # 超时配置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
      
        # 支持文件上传
        client_max_body_size 100M;
    }
  
    # MCP协议代理
    location /mcp {
        proxy_pass http://ai_doc_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header API_KEY $http_api_key;
    }
}
```

### 7.3 防火墙配置

```bash
# 1. 使用iptables配置防火墙
# 只允许必要的端口访问

# 清除现有规则
iptables -F
iptables -X

# 默认策略：拒绝所有连接
iptables -P INPUT DROP
iptables -P FORWARD DROP
iptables -P OUTPUT ACCEPT

# 允许本地回环
iptables -A INPUT -i lo -j ACCEPT

# 允许已建立的连接
iptables -A INPUT -m state --state ESTABLISHED,RELATED -j ACCEPT

# 允许SSH（仅限内网特定IP）
iptables -A INPUT -p tcp -s 192.168.1.0/24 --dport 22 -j ACCEPT

# 允许HTTP/HTTPS
iptables -A INPUT -p tcp --dport 80 -j ACCEPT
iptables -A INPUT -p tcp --dport 443 -j ACCEPT

# 允许后端API端口（仅限内网）
iptables -A INPUT -p tcp -s 192.168.1.0/24 --dport 8080:8081 -j ACCEPT

# 允许数据库端口（仅限应用服务器）
iptables -A INPUT -p tcp -s 192.168.1.101 --dport 5432 -j ACCEPT
iptables -A INPUT -p tcp -s 192.168.1.102 --dport 5432 -j ACCEPT

# 保存规则
iptables-save > /etc/iptables/rules.v4

# 2. 使用firewalld配置防火墙（CentOS/RHEL）
systemctl start firewalld
systemctl enable firewalld

# 添加服务
firewall-cmd --permanent --add-service=http
firewall-cmd --permanent --add-service=https
firewall-cmd --permanent --add-port=8080/tcp
firewall-cmd --permanent --add-port=5432/tcp

# 限制访问来源
firewall-cmd --permanent --add-rich-rule='rule family="ipv4" source address="192.168.1.0/24" port port="5432" protocol="tcp" accept'

# 重载配置
firewall-cmd --reload
```

## 8. 访问控制

### 8.1 用户权限管理

```bash
# 1. 创建不同级别的用户账号
# 系统管理员
INSERT INTO users (username, password_hash, role, created_at) 
VALUES (
    'admin', 
    'hashed_password_here', 
    'admin', 
    NOW()
);

# 普通用户
INSERT INTO users (username, password_hash, role, created_at) 
VALUES (
    'developer1', 
    'hashed_password_here', 
    'user', 
    NOW()
);

# 只读用户
INSERT INTO users (username, password_hash, role, created_at) 
VALUES (
    'viewer1', 
    'hashed_password_here', 
    'viewer', 
    NOW()
);
```

### 8.2 API访问控制

```bash
# 1. 生成API密钥
openssl rand -hex 32

# 2. 为不同用户创建API密钥
INSERT INTO api_keys (
    user_id, 
    key_value, 
    key_name, 
    permissions, 
    enabled
) VALUES (
    1,  -- admin user id
    'generated-api-key-here',
    'Admin Production Key',
    '["*"]',  -- 全部权限
    true
);

INSERT INTO api_keys (
    user_id, 
    key_value, 
    key_name, 
    permissions, 
    enabled
) VALUES (
    2,  -- developer user id
    'generated-api-key-here',
    'Developer API Key',
    '["docs:read", "docs:write", "search:read"]',  -- 有限权限
    true
);
```

### 8.3 IP白名单配置

```go
// 在后端代码中添加IP白名单检查
// internal/middleware/ip_whitelist.go

package middleware

import (
    "net"
    "strings"
    "github.com/gin-gonic/gin"
)

var allowedIPs = []string{
    "192.168.1.0/24",  // 内网网段
    "10.0.0.0/8",       // 私有网络
}

func isAllowedIP(ip string) bool {
    clientIP := net.ParseIP(ip)
    if clientIP == nil {
        return false
    }

    for _, allowedCIDR := range allowedIPs {
        _, allowedNetwork, err := net.ParseCIDR(allowedCIDR)
        if err != nil {
            continue
        }
        if allowedNetwork.Contains(clientIP) {
            return true
        }
    }
    return false
}

func IPWhitelist() gin.HandlerFunc {
    return func(c *gin.Context) {
        clientIP := c.ClientIP()
      
        // 检查X-Forwarded-For头（代理场景）
        if forwardedFor := c.GetHeader("X-Forwarded-For"); forwardedFor != "" {
            clientIP = strings.Split(forwardedFor, ",")[0]
        }

        if !isAllowedIP(clientIP) {
            c.JSON(403, gin.H{"error": "IP address not allowed"})
            c.Abort()
            return
        }

        c.Next()
    }
}
```

## 9. 备份恢复

### 9.1 自动备份配置

```bash
# 1. 配置定时备份
crontab -e

# 添加以下定时任务
# 每天凌晨2点执行备份
0 2 * * * /path/to/scripts/backup-script.sh >> /var/log/backup.log 2>&1

# 每周日凌晨3点删除7天前的备份
0 3 * * 0 find /app/backups -name "backup_*.sql.gz" -mtime +7 -delete

# 2. 增量备份（使用pg_dump的--format选项）
pg_dump -h localhost -U postgres -d ai_doc_library \
  --format=directory \
  --file=/backups/dump_$(date +%Y%m%d)

# 3. 恢复备份
pg_restore -h localhost -U postgres -d ai_doc_library \
  /backups/dump_20231201
```

### 9.2 灾难恢复计划

```bash
# 1. 创建灾难恢复脚本
cat > /scripts/disaster-recovery.sh <<'EOF'
#!/bin/bash

set -e

echo "开始灾难恢复..."

# 1. 检查备份文件
LATEST_BACKUP=$(ls -t /backups/backup_*.sql.gz | head -1)
if [ -z "$LATEST_BACKUP" ]; then
    echo "错误: 未找到备份文件"
    exit 1
fi

echo "使用备份文件: $LATEST_BACKUP"

# 2. 停止应用服务
kubectl scale deployment ai-doc-backend --replicas=0 -n ai-doc

# 3. 恢复数据库
kubectl exec -it postgres-0 -n ai-doc -- bash -c "
    PGPASSWORD=$DB_PASSWORD psql -U postgres -d ai_doc_library < /dev/stdin
" < <(gunzip -c $LATEST_BACKUP)

# 4. 重启应用服务
kubectl scale deployment ai-doc-backend --replicas=3 -n ai-doc

# 5. 等待服务就绪
kubectl wait --for=condition=ready pod -l app=ai-doc-backend -n ai-doc --timeout=300s

echo "灾难恢复完成"
EOF

chmod +x /scripts/disaster-recovery.sh
```

## 10. 监控和日志

### 10.1 系统监控

```yaml
# Prometheus监控配置
cat > configs/prometheus.yml <<EOF
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'ai-doc-backend'
    kubernetes_sd_configs:
      - role: pod
        namespaces:
          names:
            - ai-doc
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_label_app]
        action: keep
        regex: ai-doc-backend
      - source_labels: [__meta_kubernetes_pod_name]
        target_label: pod

  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres:5432']

  - job_name: 'redis'
    static_configs:
      - targets: ['redis:6379']

  - job_name: 'nginx'
    static_configs:
      - targets: ['nginx:80']

alerting:
  alertmanagers:
    - static_configs:
        - targets: ['alertmanager:9093']

rule_files:
  - '/etc/prometheus/alerts.yml'
EOF
```

### 10.2 日志管理

```bash
# 1. 配置集中式日志收集
# 使用ELK Stack或Fluentd

# 2. 日志轮转配置
cat > /etc/logrotate.d/ai-doc <<EOF
/app/logs/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 0640 appuser appgroup
    sharedscripts
    postrotate
        docker-compose exec nginx nginx -s reopen
    endscript
}
EOF

# 3. 日志分析
# 查看错误日志
grep "ERROR" /app/logs/backend.log | tail -100

# 查看慢查询
grep "duration" /app/logs/backend.log | awk '$NF > 1000'

# 统计访问量
awk '{print $1}' /app/logs/nginx/access.log | sort | uniq -c | sort -rn | head -20
```

## 11. 性能优化

### 11.1 数据库优化

```sql
-- 1. 创建索引优化查询
CREATE INDEX CONCURRENTLY idx_documents_name ON documents(name);
CREATE INDEX CONCURRENTLY idx_documents_type ON documents(type);
CREATE INDEX CONCURRENTLY idx_documents_created_at ON documents(created_at DESC);

-- 2. 分析查询性能
EXPLAIN ANALYZE SELECT * FROM documents WHERE name LIKE '%keyword%';

-- 3. 优化查询
-- 使用全文搜索替代LIKE
CREATE INDEX idx_documents_name_gin ON documents USING gin(to_tsvector('english', name));

-- 4. 更新统计信息
ANALYZE documents;
VACUUM ANALYZE documents;
```

### 11.2 应用优化

```go
// 1. 启用连接池
db.SetMaxOpenConns(100)
db.SetMaxIdleConns(10)
db.SetConnMaxLifetime(time.Hour)

// 2. 启用缓存
cache := cache.New(5*time.Minute, 10*time.Minute)

// 3. 批量操作
documents := []Document{}
// 批量插入
db.Create(&documents)

// 4. 使用goroutine并发处理
var wg sync.WaitGroup
for _, doc := range documents {
    wg.Add(1)
    go func(d Document) {
        defer wg.Done()
        processDocument(d)
    }(doc)
}
wg.Wait()
```

## 12. 故障排除

### 12.1 常见问题

```bash
# 问题1: 容器无法启动
docker logs ai-doc-backend-1
# 检查日志中的错误信息

# 问题2: 数据库连接失败
kubectl exec -it postgres-0 -n ai-doc -- psql -U postgres -d ai_doc_library
# 测试数据库连接

# 问题3: 服务无法访问
kubectl get pods -n ai-doc
kubectl describe pod <pod-name> -n ai-doc
# 检查Pod状态和事件

# 问题4: 磁盘空间不足
df -h
# 清理日志和备份文件
find /app/logs -name "*.log" -mtime +7 -delete
find /app/backups -name "backup_*" -mtime +30 -delete

# 问题5: 内存不足
kubectl top pods -n ai-doc
# 查看资源使用情况
kubectl top nodes
```

### 12.2 诊断脚本

```bash
# 创建诊断脚本
cat > /scripts/diagnose.sh <<'EOF'
#!/bin/bash

echo "========================================="
echo "AI技术文档库 - 系统诊断"
echo "========================================="
echo ""

echo "1. 检查Docker容器状态..."
docker ps -a
echo ""

echo "2. 检查磁盘使用情况..."
df -h
echo ""

echo "3. 检查内存使用情况..."
free -h
echo ""

echo "4. 检查最近的错误日志..."
grep "ERROR" /app/logs/backend.log | tail -20
echo ""

echo "5. 检查数据库连接..."
kubectl exec -it postgres-0 -n ai-doc -- psql -U postgres -c "SELECT 1" 2>/dev/null && echo "数据库连接正常" || echo "数据库连接失败"
echo ""

echo "6. 检查服务健康状态..."
curl -f http://localhost:8080/health/live && echo "后端服务正常" || echo "后端服务异常"
echo ""

echo "7. 检查Kubernetes Pod状态..."
kubectl get pods -n ai-doc
echo ""

echo "诊断完成"
EOF

chmod +x /scripts/diagnose.sh
```

## 13. 升级维护

### 13.1 版本升级

```bash
# 1. 备份当前版本
./scripts/backup-script.sh

# 2. 拉取新版本代码
git pull origin main

# 3. 构建新版本镜像
docker build -t ai-doc-backend:v2.0.0 .

# 4. 停止旧版本服务
docker-compose stop

# 5. 更新docker-compose.yml中的镜像版本
# image: ai-doc-backend:v2.0.0

# 6. 启动新版本服务
docker-compose up -d

# 7. 验证新版本
curl http://localhost:8080/health/live
```

### 13.2 滚动升级（Kubernetes）

```bash
# 使用Deployment进行滚动升级
kubectl set image deployment/ai-doc-backend \
  backend=ai-doc-backend:v2.0.0 \
  -n ai-doc

# 查看升级状态
kubectl rollout status deployment/ai-doc-backend -n ai-doc

# 如果升级失败，回滚到上一个版本
kubectl rollout undo deployment/ai-doc-backend -n ai-doc
```

## 14. 安全检查清单

在生产环境部署前，请确认以下安全检查项：

- [ ] 修改所有默认密码为强密码
- [ ] 配置TLS/SSL证书
- [ ] 配置防火墙规则，只开放必要端口
- [ ] 配置IP白名单，限制访问来源
- [ ] 启用数据库加密
- [ ] 配置备份加密
- [ ] 配置网络策略
- [ ] 启用审计日志
- [ ] 定期进行安全扫描
- [ ] 配置告警通知
- [ ] 制定灾难恢复计划
- [ ] 进行安全培训

## 15. 联系与支持

如需技术支持，请联系：

- 问题反馈: [GitHub · Where software is built](https://github.com/UniverseHappiness/LAST-doc/issues)

---

**文档版本**: v1.0.0
**最后更新**: 2026-01-03
