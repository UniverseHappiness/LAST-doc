# AI技术文档库 - 快速部署指南

## 目录

- [1. 概述](#1-概述)
- [2. 环境准备](#2-环境准备)
- [3. Docker Compose部署（推荐）](#3-docker-compose部署推荐)
- [4. Kubernetes部署](#4-kubernetes部署)
- [5. 验证部署](#5-验证部署)
- [6. 常见问题](#6-常见问题)
- [7. 性能优化](#7-性能优化)
- [8. 监控和维护](#8-监控和维护)

## 1. 概述

### 1.1 项目简介

AI技术文档库是一个企业级技术文档管理系统，支持：

- ✅ 多格式文档上传和管理（docx、pdf、markdown等）
- ✅ 智能文档检索（关键词搜索、语义搜索）
- ✅ MCP协议支持，与CoStrict IDE集成
- ✅ 用户认证和权限管理
- ✅ 高可用部署和多实例负载均衡
- ✅ 数据备份和恢复
- ✅ 系统监控和日志管理

### 1.2 部署架构

```
┌─────────────────────────────────────────────────────────┐
│                    用户访问层                             │
│  ┌──────────────┐          ┌──────────────┐            │
│  │  Web界面     │          │  MCP客户端    │            │
│  └──────────────┘          └──────────────┘            │
└─────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────┐
│                    负载均衡层                             │
│                    Nginx反向代理                          │
└─────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────┐
│                    应用服务层                             │
│  ┌──────────────┐  ┌──────────────┐                    │
│  │  Backend-1   │  │  Backend-2   │                    │
│  └──────────────┘  └──────────────┘                    │
└─────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────┐
│                    数据存储层                             │
│  ┌──────────────┐  ┌──────────────┐  ┌─────────────┐ │
│  │  PostgreSQL  │  │  文件存储     │  │  Redis缓存   │ │
│  └──────────────┘  └──────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────┘
```

### 1.3 部署方式对比

| 特性 | Docker Compose | Kubernetes |
|------|---------------|------------|
| 部署复杂度 | ⭐ 简单 | ⭐⭐⭐ 复杂 |
| 资源需求 | 低 | 中 |
| 管理成本 | 低 | 中 |
| 扩展能力 | 手动 | 自动 |
| 高可用性 | 中 | 高 |
| 适用场景 | 开发/测试/小规模生产 | 大规模生产 |
| 学习曲线 | 平缓 | 陡峭 |

## 2. 环境准备

### 2.1 系统要求

#### 最低配置

| 组件 | CPU | 内存 | 存储 | 网络 |
|------|-----|------|------|------|
| 服务器 | 4核 | 8GB | 100GB | 1Gbps |

#### 推荐配置

| 组件 | CPU | 内存 | 存储 | 网络 |
|------|-----|------|------|------|
| 服务器 | 8核+ | 16GB+ | 500GB+ | 10Gbps |

### 2.2 软件要求

#### Docker Compose部署

```bash
# 检查系统
uname -a  # Linux x86_64

# 安装Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# 安装Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 验证安装
docker --version
docker-compose --version
```

#### Kubernetes部署

```bash
# 安装kubectl
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x kubectl
sudo mv kubectl /usr/local/bin/

# 或使用Minikube（用于本地测试）
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube

# 或使用k3s（轻量级Kubernetes）
curl -sfL https://get.k3s.io | sh -

# 验证安装
kubectl version --client
```

### 2.3 网络配置

```bash
# 检查网络连接
ping -c 4 8.8.8.8

# 检查端口占用
netstat -tuln | grep -E '(80|443|8080|5432)'

# 开放防火墙端口
# Ubuntu/Debian
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 8080/tcp

# CentOS/RHEL
sudo firewall-cmd --permanent --add-port=80/tcp
sudo firewall-cmd --permanent --add-port=443/tcp
sudo firewall-cmd --permanent --add-port=8080/tcp
sudo firewall-cmd --reload
```

## 3. Docker Compose部署（推荐）

### 3.1 快速开始

#### 步骤1: 克隆项目

```bash
# 克隆仓库
git clone <your-repo-url> ai-doc-library
cd ai-doc-library

# 或下载发布版本
wget <your-release-url>
tar -xzf ai-doc-library-v1.0.0.tar.gz
cd ai-doc-library
```

#### 步骤2: 配置环境变量

```bash
# 复制环境变量模板（如果有）
cp .env.example .env

# 编辑配置文件
vi .env

# 添加以下配置
# DB_HOST=postgres
# DB_PORT=5432
# DB_USER=postgres
# DB_PASSWORD=your-strong-password-here
# DB_NAME=ai_doc_library
# JWT_SECRET=your-jwt-secret-key-here
# SERVER_PORT=8080
```

#### 步骤3: 修改Docker Compose配置

```bash
# 编辑docker-compose.yml
vi docker-compose.yml

# 修改以下配置（可选）：
# 1. 端口映射
# 2. 存储卷路径
# 3. 资源限制
# 4. 环境变量
```

#### 步骤4: 构建并启动服务

```bash
# 构建镜像
docker-compose build

# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

#### 步骤5: 初始化数据库

```bash
# 等待数据库启动完成
docker-compose logs postgres | grep "database system is ready"

# 初始化数据库（如果需要）
docker-compose exec postgres psql -U postgres -d ai_doc_library -f /docker-entrypoint-initdb.d/init.sql

# 创建管理员账户
docker-compose exec backend ./main create-admin
# 或使用API创建
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "Admin@123"
  }'
```

### 3.2 服务管理

#### 查看服务状态

```bash
# 查看所有服务状态
docker-compose ps

# 查看特定服务日志
docker-compose logs backend
docker-compose logs nginx

# 实时查看日志
docker-compose logs -f backend

# 查看资源使用
docker stats
```

#### 停止和启动服务

```bash
# 停止所有服务
docker-compose stop

# 启动所有服务
docker-compose start

# 重启特定服务
docker-compose restart backend

# 停止并删除所有服务
docker-compose down

# 停止并删除所有服务和数据卷
docker-compose down -v
```

#### 更新服务

```bash
# 拉取最新代码
git pull origin main

# 重新构建镜像
docker-compose build

# 重启服务
docker-compose up -d

# 查看更新状态
docker-compose ps
```

### 3.3 配置选项

#### 端口配置

```yaml
# 在docker-compose.yml中修改端口映射
services:
  nginx:
    ports:
      - "8000:80"  # 将HTTP端口改为8000
      - "8443:443"  # 将HTTPS端口改为8443
  
  backend:
    ports:
      - "9000:8080"  # 将后端端口改为9000
  
  postgres:
    ports:
      - "5433:5432"  # 将数据库端口改为5433
```

#### 存储配置

```yaml
# 修改存储卷路径
volumes:
  postgres_data:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: /your/custom/path/postgres  # 自定义路径
  
  app_storage:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: /your/custom/path/storage  # 自定义路径
```

#### 资源限制

```yaml
# 添加资源限制
services:
  backend:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '1'
          memory: 1G
  
  postgres:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 4G
        reservations:
          cpus: '1'
          memory: 2G
```

### 3.4 高可用部署

#### 多实例部署

```yaml
# 在docker-compose.yml中已有配置
# backend - 主实例
# backend2 - 扩展实例

# 可以添加更多实例
backend3:
  build:
    context: .
    dockerfile: Dockerfile
  container_name: ai-doc-backend-3
  ports:
    - "8082:8080"
  environment:
    - DB_HOST=postgres
    - DB_PORT=5432
    - DB_USER=postgres
    - DB_PASSWORD=postgres
    - DB_NAME=ai_doc_library
    - SERVER_PORT=8080
    - NODE_ID=backend-3
  volumes:
    - ./storage:/app/storage
    - ./backups:/app/backups
  depends_on:
    postgres:
      condition: service_started
  restart: unless-stopped
```

#### 负载均衡配置

```nginx
# 在nginx.conf中已配置负载均衡
upstream ai_doc_backend {
    least_conn;  # 最少连接算法
    
    # 后端服务器列表
    server backend:8080 max_fails=3 fail_timeout=30s weight=1;
    server ai-doc-backend-2:8080 max_fails=3 fail_timeout=30s weight=1;
    server ai-doc-backend-3:8080 max_fails=3 fail_timeout=30s weight=1;
    
    keepalive 32;
    keepalive_timeout 60s;
}
```

## 4. Kubernetes部署

### 4.1 快速开始

#### 方式一: 使用完整部署清单

```bash
# 1. 构建镜像
docker build -t ai-doc-backend:latest .

# 2. 使用Kubernetes清单部署所有资源
kubectl apply -f k8s/k8s-all.yaml

# 3. 查看部署状态
kubectl get all -n ai-doc

# 4. 等待Pod就绪
kubectl wait --for=condition=ready pod -l app=ai-doc-backend -n ai-doc --timeout=300s
```

#### 方式二: 使用部署脚本

```bash
# 1. 使用脚本部署（完整部署）
NAMESPACE=ai-doc ./scripts/deploy-k8s.sh

# 2. 最小化部署（仅核心组件）
DEPLOYMENT_MODE=minimal NAMESPACE=ai-doc ./scripts/deploy-k8s.sh

# 3. 包含监控的部署
DEPLOYMENT_MODE=monitoring NAMESPACE=ai-doc ./scripts/deploy-k8s.sh

# 4. 查看帮助
./scripts/deploy-k8s.sh help
```

### 4.2 详细部署步骤

#### 步骤1: 准备命名空间

```bash
# 创建命名空间
kubectl create namespace ai-doc

# 设置默认命名空间
kubectl config set-context --current --namespace=ai-doc

# 验证命名空间
kubectl get namespaces
```

#### 步骤2: 配置密钥

```bash
# 生成强密码
DB_PASSWORD=$(openssl rand -base64 32)
JWT_SECRET=$(openssl rand -base64 32)

# 创建Secret
kubectl create secret generic ai-doc-secrets \
  --from-literal=db-password=$DB_PASSWORD \
  --from-literal=jwt-secret=$JWT_SECRET \
  --from-literal=db-user=postgres \
  -n ai-doc

# 验证Secret
kubectl describe secret ai-doc-secrets -n ai-doc
```

#### 步骤3: 配置ConfigMap

```bash
# 创建ConfigMap (k8s/configmap.yaml已包含完整配置)
kubectl apply -f k8s/configmap.yaml -n ai-doc

# 验证ConfigMap
kubectl describe configmap ai-doc-config -n ai-doc
```

#### 步骤4: 部署存储

```bash
# 创建持久化存储声明
kubectl apply -f k8s/pvc.yaml -n ai-doc

# 查看PVC状态
kubectl get pvc -n ai-doc

# 等待PVC绑定
kubectl wait --for=jsonpath='{.status.phase}'=Bound pvc/ai-doc-storage-pvc -n ai-doc --timeout=120s
```

#### 步骤5: 部署数据库

```bash
# 部署PostgreSQL
kubectl apply -f k8s/postgres.yaml -n ai-doc

# 查看数据库状态
kubectl get pods -l app=postgres -n ai-doc

# 等待数据库就绪
kubectl wait --for=condition=ready pod -l app=postgres -n ai-doc --timeout=300s

# 查看数据库日志
kubectl logs -f postgres-0 -n ai-doc
```

#### 步骤6: 部署应用

```bash
# 部署后端应用
kubectl apply -f k8s/deployment.yaml -n ai-doc

# 查看应用状态
kubectl get deployments -n ai-doc
kubectl get pods -l app=ai-doc-backend -n ai-doc

# 等待应用就绪
kubectl wait --for=condition=ready pod -l app=ai-doc-backend -n ai-doc --timeout=300s
```

#### 步骤7: 部署服务

```bash
# 部署服务
kubectl apply -f k8s/deployment.yaml -n ai-doc

# 查看服务
kubectl get svc -n ai-doc

# 测试服务访问
kubectl port-forward -n ai-doc svc/ai-doc-backend-service 8080:8080
```

#### 步骤8: 配置Ingress

```bash
# 部署Ingress
kubectl apply -f k8s/ingress.yaml -n ai-doc

# 查看Ingress
kubectl get ingress -n ai-doc

# 配置本地hosts（如果需要）
echo "$(kubectl get ingress ai-doc-ingress -n ai-doc -o jsonpath='{.status.loadBalancer.ingress[0].ip}') ai-doc.local" | sudo tee -a /etc/hosts
```

### 4.3 部署管理

#### 查看部署状态

```bash
# 查看所有资源
kubectl get all -n ai-doc

# 查看Pod详情
kubectl describe pod <pod-name> -n ai-doc

# 查看Pod日志
kubectl logs <pod-name> -n ai-doc

# 实时查看日志
kubectl logs -f <pod-name> -n ai-doc
```

#### 扩缩容

```bash
# 手动扩展Pod副本数
kubectl scale deployment ai-doc-backend --replicas=5 -n ai-doc

# 查看扩缩容状态
kubectl get pods -l app=ai-doc-backend -n ai-doc

# 使用HPA自动扩缩容
kubectl get hpa -n ai-doc

# 调整HPA配置
kubectl edit hpa ai-doc-backend-hpa -n ai-doc
```

#### 更新部署

```bash
# 滚动更新
kubectl set image deployment/ai-doc-backend \
  backend=ai-doc-backend:v2.0.0 \
  -n ai-doc

# 查看更新状态
kubectl rollout status deployment/ai-doc-backend -n ai-doc

# 查看更新历史
kubectl rollout history deployment/ai-doc-backend -n ai-doc

# 回滚到上一版本
kubectl rollout undo deployment/ai-doc-backend -n ai-doc

# 回滚到指定版本
kubectl rollout undo deployment/ai-doc-backend --to-revision=2 -n ai-doc
```

#### 删除部署

```bash
# 删除特定资源
kubectl delete deployment ai-doc-backend -n ai-doc
kubectl delete service ai-doc-backend-service -n ai-doc

# 删除命名空间及其所有资源
kubectl delete namespace ai-doc

# 使用清单删除所有资源
kubectl delete -f k8s/k8s-all.yaml
```

### 4.4 故障排查

#### Pod无法启动

```bash
# 查看Pod事件
kubectl describe pod <pod-name> -n ai-doc

# 查看Pod日志
kubectl logs <pod-name> -n ai-doc

# 查看之前的日志（如果Pod重启了）
kubectl logs <pod-name> -n ai-doc --previous

# 进入Pod调试
kubectl exec -it <pod-name> -n ai-doc -- /bin/sh
```

#### 服务无法访问

```bash
# 检查Service配置
kubectl describe svc ai-doc-backend-service -n ai-doc

# 检查Endpoints
kubectl get endpoints ai-doc-backend-service -n ai-doc

# 测试服务连接
kubectl run test-pod --image=busybox --rm -it --restart=Never -n ai-doc -- wget -O- http://ai-doc-backend-service:8080/health
```

#### 存储问题

```bash
# 查看PVC状态
kubectl describe pvc ai-doc-storage-pvc -n ai-doc

# 查看PV状态
kubectl get pv

# 查看存储类
kubectl get storageclass
```

## 5. 验证部署

### 5.1 健康检查

#### Docker Compose部署

```bash
# 检查所有容器状态
docker-compose ps

# 检查健康状态
curl http://localhost/health
curl http://localhost:8080/health/live
curl http://localhost:8080/health/ready

# 检查数据库连接
docker-compose exec postgres pg_isready -U postgres

# 检查Nginx状态
curl -I http://localhost
```

#### Kubernetes部署

```bash
# 检查Pod健康状态
kubectl get pods -n ai-doc

# 检查服务健康状态
kubectl get endpoints -n ai-doc

# 通过端口转发测试
kubectl port-forward -n ai-doc svc/ai-doc-backend-service 8080:8080
curl http://localhost:8080/health/live
```

### 5.2 功能测试

#### Web界面访问

```bash
# 测试前端界面
curl http://localhost

# 验证静态文件
curl http://localhost/index.html

# 测试API路由
curl http://localhost/api/v1/documents

# 验证CORS配置
curl -H "Origin: http://localhost" \
     -H "Access-Control-Request-Method: GET" \
     -X OPTIONS http://localhost/api/v1/documents
```

#### API功能测试

```bash
# 1. 用户注册
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Test@123"
  }'

# 2. 用户登录
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test@123"
  }'

# 3. 创建API密钥
TOKEN="your-jwt-token-here"
curl -X POST http://localhost:8080/api/v1/api-keys \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "key_name": "Test Key",
    "permissions": ["*"]
  }'

# 4. MCP协议测试
API_KEY="your-api-key-here"
curl -X POST http://localhost/mcp \
  -H "API-Key: $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "search",
    "query": "test query"
  }'
```

#### 性能测试

```bash
# 使用Apache Bench进行压力测试
ab -n 1000 -c 10 http://localhost:8080/health

# 使用wrk进行性能测试
wrk -t4 -c100 -d30s http://localhost:8080/health

# 查看响应时间
time curl http://localhost:8080/health

# 测试并发性能
for i in {1..100}; do
  curl http://localhost:8080/health &
done
wait
```

### 5.3 日志检查

#### Docker Compose

```bash
# 查看所有服务日志
docker-compose logs

# 查看特定服务日志
docker-compose logs backend
docker-compose logs nginx
docker-compose logs postgres

# 实时查看日志
docker-compose logs -f backend

# 查看最近100行日志
docker-compose logs --tail=100 backend

# 查看错误日志
docker-compose logs backend | grep ERROR
```

#### Kubernetes

```bash
# 查看Pod日志
kubectl logs <pod-name> -n ai-doc

# 实时查看日志
kubectl logs -f <pod-name> -n ai-doc

# 查看所有Pod的日志
kubectl logs -l app=ai-doc-backend -n ai-doc

# 查看错误日志
kubectl logs <pod-name> -n ai-doc | grep ERROR

# 查看最近日志
kubectl logs --tail=100 <pod-name> -n ai-doc
```

## 6. 常见问题

### 6.1 部署问题

#### 问题1: 端口已被占用

**错误信息**:
```
Error starting userland proxy: listen tcp 0.0.0.0:8080: bind: address already in use
```

**解决方案**:
```bash
# 查找占用端口的进程
lsof -i :8080
# 或
netstat -tuln | grep 8080

# 停止占用端口的进程
kill -9 <PID>

# 或修改docker-compose.yml中的端口映射
```

#### 问题2: 内存不足

**错误信息**:
```
no space left on device
```

**解决方案**:
```bash
# 检查磁盘空间
df -h

# 清理Docker缓存
docker system prune -a

# 清理未使用的镜像
docker image prune -a

# 检查Docker日志
du -sh /var/lib/docker

# 移动Docker数据目录（需要重启Docker）
```

#### 问题3: 数据库连接失败

**错误信息**:
```
could not connect to server: Connection refused
```

**解决方案**:
```bash
# 检查数据库是否运行
docker-compose ps postgres
# 或
kubectl get pods -l app=postgres -n ai-doc

# 检查数据库日志
docker-compose logs postgres
# 或
kubectl logs postgres-0 -n ai-doc

# 测试数据库连接
docker-compose exec postgres pg_isready -U postgres

# 检查网络配置
docker network inspect ai-doc-network
```

### 6.2 运行问题

#### 问题1: 服务响应缓慢

**诊断步骤**:
```bash
# 检查资源使用情况
docker stats
# 或
kubectl top pods -n ai-doc

# 检查数据库性能
docker-compose exec postgres psql -U postgres -d ai_doc_library -c "SELECT * FROM pg_stat_activity;"

# 检查慢查询
docker-compose logs postgres | grep "duration"

# 检查应用日志
docker-compose logs backend | grep "slow"
```

**优化方案**:
```bash
# 增加资源限制
# 在docker-compose.yml中添加
deploy:
  resources:
    limits:
      cpus: '2'
      memory: 2G

# 启用Redis缓存
# 在docker-compose.yml中取消Redis服务的注释

# 添加索引
docker-compose exec postgres psql -U postgres -d ai_doc_library -c "CREATE INDEX idx_documents_name ON documents(name);"

# 优化数据库查询时间
docker-compose exec postgres psql -U postgres -d ai_doc_library -c "ANALYZE;"
```

#### 问题2: 文件上传失败

**诊断步骤**:
```bash
# 检查Nginx配置
grep client_max_body_size nginx.conf

# 检查应用日志
docker-compose logs backend | grep ERROR

# 检查存储空间
df -h ./storage

# 测试文件上传权限
ls -la ./storage
```

**解决方案**:
```bash
# 增加上传限制
# 在nginx.conf中添加
client_max_body_size 200M;

# 修改Docker Compose配置
services:
  nginx:
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro

# 修改存储权限
chmod 777 ./storage
```

### 6.3 安全问题

#### 问题1: 默认密码未修改

**解决方案**:
```bash
# 生成强密码
openssl rand -base64 32

# 更新Docker Compose配置
vi docker-compose.yml
# 修改DB_PASSWORD和JWT_SECRET

# 更新Kubernetes Secret
kubectl create secret generic ai-doc-secrets \
  --from-literal=db-password=new-password \
  --from-literal=jwt-secret=new-secret \
  --dry-run=client -o yaml | kubectl apply -f -

# 重启服务
docker-compose restart
# 或
kubectl rollout restart deployment/ai-doc-backend -n ai-doc
```

#### 问题2: 证书过期

**解决方案**:
```bash
# 检查证书有效期
openssl x509 -in /path/to/cert.crt -noout -dates

# 更新证书
# 生成新证书
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout tls.key -out tls.crt

# 更新Kubernetes Secret
kubectl delete secret ai-doc-tls-secret -n ai-doc
kubectl create secret tls ai-doc-tls-secret \
  --cert=tls.crt \
  --key=tls.key \
  -n ai-doc

# 重启Ingress Controller
kubectl delete pods -l app=ingress-nginx -n ingress-nginx
```

## 7. 性能优化

### 7.1 数据库优化

```sql
-- 1. 创建索引
CREATE INDEX CONCURRENTLY idx_documents_name ON documents(name);
CREATE INDEX CONCURRENTLY idx_documents_type ON documents(type);
CREATE INDEX CONCURRENTLY idx_documents_created_at ON documents(created_at DESC);

-- 2. 启用查询缓存
ALTER DATABASE ai_doc_library SET shared_buffers = '256MB';
ALTER DATABASE ai_doc_library SET effective_cache_size = '1GB';

-- 3. 优化连接池
ALTER DATABASE ai_doc_library SET max_connections = 200;

-- 4. 定期维护
VACUUM ANALYZE documents;
REINDEX TABLE documents;
```

### 7.2 应用优化

```go
// 1. 启用数据库连接池
db.SetMaxOpenConns(100)
db.SetMaxIdleConns(10)
db.SetConnMaxLifetime(time.Hour)

// 2. 启用并发处理
go processDocuments()

// 3. 使用缓存
cache := cache.New(5*time.Minute, 10*time.Minute)
result, found := cache.Get("key")

// 4. 批量操作
db.CreateInBatches(documents, 100)
```

### 7.3 系统优化

```bash
# 1. 调整文件描述符限制
echo "* soft nofile 65536" >> /etc/security/limits.conf
echo "* hard nofile 65536" >> /etc/security/limits.conf

# 2. 调整内核参数
echo "net.ipv4.tcp_max_syn_backlog = 4096" >> /etc/sysctl.conf
echo "net.core.somaxconn = 4096" >> /etc/sysctl.conf
sysctl -p

# 3. 优化Docker
echo '{"log-driver":"json-file","log-opts":{"max-size":"10m","max-file":"3"}}' > /etc/docker/daemon.json
systemctl restart docker
```

## 8. 监控和维护

### 8.1 日志管理

```bash
# 1. 配置日志轮转
cat > /etc/logrotate.d/ai-doc <<EOF
/var/log/ai-doc/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 0640 appuser appgroup
}
EOF

# 2. 查看日志大小
du -sh /var/log/ai-doc/*

# 3. 清理旧日志
find /var/log/ai-doc -name "*.log" -mtime +7 -delete
```

### 8.2 备份策略

```bash
# 1. 配置自动备份
crontab -e
# 添加
0 2 * * * /app/scripts/backup-script.sh

# 2. 手动备份
kubectl exec -it postgres-0 -n ai-doc -- pg_dump -U postgres ai_doc_library > backup.sql

# 3. 恢复备份
cat backup.sql | kubectl exec -i postgres-0 -n ai-doc -- psql -U postgres ai_doc_library
```

### 8.3 监控告警

```yaml
# Prometheus告警规则
cat > alerts.yml <<EOF
groups:
- name: ai-doc-alerts
  rules:
  - alert: HighErrorRate
    expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
    for: 5m
    annotations:
      summary: "High error rate detected"
  
  - alert: HighLatency
    expr: histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m])) > 1
    for: 10m
    annotations:
      summary: "High latency detected"
  
  - alert: PodCrashLooping
    expr: rate(kube_pod_container_status_restarts_total[15m]) > 0
    for: 5m
    annotations:
      summary: "Pod is crash looping"
EOF

kubectl create configmap prometheus-alerts --from-file=alerts.yml -n ai-doc
```

## 9. 下一步

部署完成后，请参考以下文档：

- [私有化部署配置](PRIVATE_DEPLOYMENT.md) - 内网环境部署和安全配置
- [可靠性指南](reliability_guide.md) - 高可用和故障恢复
- [可扩展性指南](scalability_guide.md) - 水平扩展和性能优化
- [MCP本地使用指南](mcp_local_usage_guide.md) - 与CoStrict IDE集成

## 10. 技术支持

如需帮助，请联系：

- 📧 技术支持邮箱: support@example.com
- 📖 文档地址: http://localhost/docs
- 🐛 问题反馈: http://localhost/issues

---

**文档版本**: v1.0.0  
**最后更新**: 2026-01-03  
**维护团队**: AI技术文档库团队