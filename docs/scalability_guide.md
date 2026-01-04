# 系统扩展性指南

本指南说明如何使用AI技术文档库的横向扩展和存储扩展功能。

## 概述

AI技术文档库支持以下扩展性功能：

1. **服务横向扩展**：通过增加后端服务实例来提高系统处理能力
2. **存储扩展**：支持从本地存储无缝切换到分布式存储（S3/MinIO）

## 服务横向扩展

### 1. 启动多实例部署

系统默认部署一个后端实例。要启用横向扩展，需要启动多个后端实例。

#### 方式一：修改docker-compose.yml

编辑 [`docker-compose.yml`](../docker-compose.yml:1) 文件，取消注释backend2服务的配置：

```yaml
  # Go后端服务 - 扩展实例1（可选，用于横向扩展）
  backend2:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: ai-doc-backend-2
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=ai_doc_library
      - SERVER_PORT=8080
      - STORAGE_DIR=/app/storage
      - STORAGE_TYPE=local
      - NODE_ID=backend-2
    volumes:
      - ./storage:/app/storage
    depends_on:
      postgres:
        condition: service_started
    restart: unless-stopped
    networks:
      - ai-doc-network
```

#### 方式二：使用docker-compose scale命令

直接使用docker-compose的scale功能：

```bash
# 启动2个后端实例
docker-compose up -d --scale backend=2

# 启动3个后端实例
docker-compose up -d --scale backend=3
```

### 2. 配置Nginx负载均衡

编辑 [`nginx.conf`](../nginx.conf:1) 文件，在 `upstream ai_doc_backend` 部分添加更多后端实例：

```nginx
    upstream ai_doc_backend {
        least_conn;
        
        server backend:8080 max_fails=3 fail_timeout=30s weight=1;
        server ai-doc-backend-2:8080 max_fails=3 fail_timeout=30s weight=1;
        server ai-doc-backend-3:8080 max_fails=3 fail_timeout=30s weight=1;
        
        keepalive 32;
        keepalive_timeout 60s;
    }
```

### 3. 启动服务

```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看后端服务日志
docker-compose logs -f backend
```

### 4. 验证横向扩展

```bash
# 启动多个并发请求测试负载均衡
for i in {1..10}; do curl http://localhost/api/v1/documents; done
```

## 存储扩展

系统支持三种存储类型：本地存储、S3存储、MinIO存储。

### 1. 本地存储（默认）

本地存储是默认配置，无需额外配置。

**环境变量配置：**
```bash
STORAGE_TYPE=local
STORAGE_DIR=/app/storage
```

**优点：**
- 配置简单
- 无需额外服务
- 适合小规模部署

**缺点：**
- 无法跨实例共享存储
- 扩展性有限
- 单点故障风险

### 2. MinIO存储（推荐用于分布式部署）

MinIO是一个兼容S3 API的开源对象存储系统，适合内部部署。

#### 启动MinIO服务

MinIO服务已包含在 [`docker-compose.yml`](../docker-compose.yml:1) 中，直接启动即可：

```bash
docker-compose up -d minio
```

MinIO服务启动后，可通过以下地址访问：

- API端点：http://localhost:9000
- 管理控制台：http://localhost:9001
- 默认用户名：`minioadmin`
- 默认密码：`minioadmin123`

#### 配置后端使用MinIO

修改 [`docker-compose.yml`](../docker-compose.yml:1) 中backend服务的环境变量：

```yaml
  backend:
    environment:
      - STORAGE_TYPE=minio
      - STORAGE_DIR=/app/storage
      - MINIO_ENDPOINT=minio:9000
      - MINIO_ACCESS_KEY=minioadmin
      - MINIO_SECRET_KEY=minioadmin123
      - MINIO_BUCKET=ai-doc-library
      - MINIO_LOCATION=us-east-1
      - MINIO_USE_SSL=false
```

#### 创建存储桶

首次使用MinIO前，需要创建存储桶：

```bash
# 使用MinIO客户端（mc）创建存储桶
docker run --rm minio/mc \
  alias set minio http://localhost:9000 minioadmin minioadmin123
docker run --rm minio/mc \
  mb minio/ai-doc-library
```

或在管理控制台（http://localhost:9001）中手动创建。

### 3. AWS S3存储

适用于云环境，使用AWS S3作为对象存储。

**环境变量配置：**

```yaml
  backend:
    environment:
      - STORAGE_TYPE=s3
      - S3_REGION=us-east-1
      - S3_BUCKET=your-bucket-name
      - S3_ACCESS_KEY=your-access-key
      - S3_SECRET_KEY=your-secret-key
      - S3_ENDPOINT=https://s3.amazonaws.com
      - S3_DISABLE_SSL=false
```

**注意事项：**
- 需要先在AWS控制台创建S3存储桶
- 确保IAM用户有足够的权限访问S3
- 存储桶的访问权限需要正确配置

## 完整部署示例

### 示例1：最小化部署（单实例 + 本地存储）

```bash
docker-compose up -d
```

### 示例2：横向扩展部署（2个后端实例 + 本地存储）

```yaml
# docker-compose.yml
services:
  backend:
    # ... 配置 ...
    environment:
      - STORAGE_TYPE=local
      - NODE_ID=backend-1
  
  backend2:
    # ... 配置 ...
    environment:
      - STORAGE_TYPE=local
      - NODE_ID=backend-2
```

```bash
docker-compose up -d
```

### 示例3：分布式存储部署（单实例 + MinIO）

```yaml
# docker-compose.yml
services:
  backend:
    environment:
      - STORAGE_TYPE=minio
      - MINIO_ENDPOINT=minio:9000
      - MINIO_ACCESS_KEY=minioadmin
      - MINIO_SECRET_KEY=minioadmin123
      - MINIO_BUCKET=ai-doc-library
```

```bash
# 创建MinIO存储桶
docker run --rm minio/mc \
  alias set minio http://localhost:9000 minioadmin minioadmin123
docker run --rm minio/mc \
  mb minio/ai-doc-library

# 启动服务
docker-compose up -d
```

### 示例4：完整高可用部署（多实例 + MinIO）

```bash
# 1. 启动MinIO并创建存储桶
docker-compose up -d minio
docker run --rm minio/mc \
  alias set minio http://localhost:9000 minioadmin minioadmin123
docker run --rm minio/mc \
  mb minio/ai-doc-library

# 2. 启动多个后端实例
docker-compose up -d backend
docker-compose up -d backend2

# 3. 启动Nginx
docker-compose up -d nginx
```

## 性能优化建议

### 1. 负载均衡算法选择

在 [`nginx.conf`](../nginx.conf:1) 中，支持多种负载均衡算法：

- `least_conn`：最少连接数（推荐）
- `least_time`：最少响应时间（需要Nginx Plus）
- `ip_hash`：基于IP的哈希
- `random`：随机选择

```nginx
upstream ai_doc_backend {
    least_conn;  # 使用最少连接算法
    
    server backend:8080 weight=1;
    server backend2:8080 weight=1;
}
```

### 2. 健康检查配置

在 [`nginx.conf`](../nginx.conf:1) 中配置健康检查：

```nginx
server backend:8080 max_fails=3 fail_timeout=30s weight=1;
```

- `max_fails`：最大失败次数
- `fail_timeout`：失败后的超时时间

### 3. 连接池配置

```nginx
keepalive 32;
keepalive_timeout 60s;
```

根据实际负载调整连接池大小。

## 监控和日志

### 查看服务状态

```bash
docker-compose ps
```

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f backend
docker-compose logs -f nginx
```

### 监控系统扩展性

访问系统监控页面：http://localhost/monitor

查看：
- 系统性能指标
- 请求处理统计
- 存储使用情况

## 故障排查

### 问题1：后端实例无法启动

检查数据库连接：
```bash
docker-compose logs backend
```

确保postgres服务已启动：
```bash
docker-compose ps postgres
```

### 问题2：Nginx无法连接到后端

检查upstream配置：
```bash
docker-compose logs nginx
```

确保后端服务正常运行：
```bash
docker-compose ps backend
```

### 问题3：MinIO连接失败

检查MinIO服务状态：
```bash
docker-compose ps minio
docker-compose logs minio
```

检查网络连接：
```bash
docker-compose exec backend ping minio
```

### 问题4：存储桶不存在

手动创建存储桶：
```bash
docker run --rm minio/mc \
  alias set minio http://localhost:9000 minioadmin minioadmin123
docker run --rm minio/mc \
  mb minio/ai-doc-library
```

## 最佳实践

1. **生产环境建议使用MinIO或S3**：避免单点故障，支持横向扩展
2. **使用至少2个后端实例**：提高系统可用性
3. **配置适当的资源限制**：在docker-compose.yml中设置内存和CPU限制
4. **定期备份数据**：包括数据库和存储文件
5. **监控系统性能**：根据负载动态调整实例数量
6. **使用环境变量管理配置**：便于部署和维护
7. **配置日志收集**：便于故障排查和性能分析

## 参考资源

- [Docker Compose文档](https://docs.docker.com/compose/)
- [Nginx负载均衡文档](https://docs.nginx.com/nginx/admin-guide/load-balancer/)
- [MinIO官方文档](https://min.io/docs/minio/linux/index.html)
- [AWS S3文档](https://docs.aws.amazon.com/s3/)