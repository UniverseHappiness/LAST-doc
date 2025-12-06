# AI文档库 - 部署指南

本文档提供了AI文档库项目的前后端分离部署方案。

## 架构概述

本项目采用前后端分离的架构，并包含专门的文档解析服务：

- **前端**：Vue.js + Bootstrap + Vite构建，通过Nginx提供静态文件服务
- **后端**：Go语言编写的RESTful API服务
- **数据库**：PostgreSQL
- **反向代理**：Nginx，用于静态文件服务和API代理
- **文档解析服务**：Python实现的gRPC服务，专门处理PDF和DOCX文档解析

## 部署方式

### 方式一：使用Docker Compose（推荐）

这是最简单的部署方式，适合开发和生产环境。

#### 前提条件

- Docker和Docker Compose已安装
- 服务器内存至少2GB

#### 部署步骤

1. **克隆项目**
   ```bash
   git clone <your-repo-url>
   cd ai-doc-library
   ```

2. **修改配置（可选）**
   
   如果需要修改默认配置，可以编辑以下文件：
   - `docker-compose.yml`：修改端口映射、环境变量等
   - `nginx.conf`：修改Nginx配置

3. **启动服务**
   ```bash
   # 构建并启动所有服务
   docker-compose up -d
   
   # 查看服务状态
   docker-compose ps
   
   # 查看日志
   docker-compose logs -f
   ```

4. **访问应用**
   
   - 前端界面：http://your-server-ip
   - API健康检查：http://your-server-ip/health
   - 直接访问后端API：http://your-server-ip:8080/health
   - Python解析服务：gRPC服务监听端口50051

5. **停止服务**
   ```bash
   docker-compose down
   ```

### 方式二：手动部署

如果您需要更灵活的部署方式，可以手动部署各个组件。

#### 1. 部署Go后端服务

```bash
# 1. 编译Go应用
go build -o bin/ai-doc-library cmd/main.go

# 2. 创建配置文件（可选，也可以使用环境变量）
cat > .env << EOF
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=ai_doc_library
SERVER_PORT=8080
STORAGE_DIR=./storage
EOF

# 3. 创建存储目录
mkdir -p storage

# 4. 运行应用
./bin/ai-doc-library
```

#### 2. 部署PostgreSQL数据库

```bash
# 使用Docker运行PostgreSQL
docker run -d \
  --name ai-doc-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=your_password \
  -e POSTGRES_DB=ai_doc_library \
  -p 5432:5432 \
  -v postgres_data:/var/lib/postgresql/data \
  postgres:15-alpine

# 或者直接在系统上安装PostgreSQL
# (略，请参考PostgreSQL官方文档)
```

#### 3. 部署Nginx前端服务

```bash
# 1. 安装Nginx
# Ubuntu/Debian
sudo apt update
sudo apt install nginx

# CentOS/RHEL
sudo yum install epel-release
sudo yum install nginx

# 2. 构建前端
cd /path/to/your/project
npm install
npm run build

# 3. 复制前端文件
sudo cp -r web/dist/* /usr/share/nginx/html/

# 4. 复制Nginx配置
sudo cp nginx.conf /etc/nginx/nginx.conf

# 5. 测试Nginx配置
sudo nginx -t

# 6. 重启Nginx
sudo systemctl restart nginx
```

#### 4. 部署Python解析服务（新增）

Python解析服务专门处理PDF和DOCX文档的解析，通过gRPC与Go后端通信。

```bash
# 1. 进入Python解析服务目录
cd python-parser-service

# 2. 安装依赖
python3 -m venv venv
source venv/bin/activate  # Linux/Mac
# 或 venv\\Scripts\\activate  # Windows
pip install -r requirements.txt

# 3. 生成gRPC代码
./generate_grpc.sh

# 4. 启动服务（推荐使用systemd管理）
# 创建systemd服务文件
sudo tee /etc/systemd/system/ai-doc-parser.service > /dev/null <<EOF
[Unit]
Description=AI文档库解析服务
After=network.target

[Service]
Type=simple
User=$USER
WorkingDirectory=/path/to/your/project/python-parser-service
Environment=PATH=/path/to/your/project/python-parser-service/venv/bin
ExecStart=/path/to/your/project/python-parser-service/venv/bin/python -m service.server
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# 5. 启动服务
sudo systemctl daemon-reload
sudo systemctl enable ai-doc-parser
sudo systemctl start ai-doc-parser

# 6. 检查服务状态
sudo systemctl status ai-doc-parser
sudo journalctl -u ai-doc-parser -f
```

## 环境变量配置

### 后端环境变量

| 变量名 | 默认值 | 描述 |
|--------|--------|------|
| DB_HOST | localhost | 数据库主机地址 |
| DB_PORT | 5432 | 数据库端口 |
| DB_USER | postgres | 数据库用户名 |
| DB_PASSWORD | postgres | 数据库密码 |
| DB_NAME | ai_doc_library | 数据库名称 |
| SERVER_PORT | 8080 | 后端服务端口 |
| STORAGE_DIR | ./storage | 文件存储目录 |
| PARSER_SERVICE_ADDR | localhost:50051 | Python解析服务地址 |

### 前端配置

前端API基础URL在 `web/src/utils/documentService.js` 中配置：

```javascript
const apiBase = '/api/v1';
```

在生产环境中，如果前端和后端部署在不同的域名下，需要修改为：

```javascript
const apiBase = 'https://your-api-domain.com/api/v1';
```

#### 前端构建

前端使用Vite构建，构建命令如下：

```bash
# 开发模式
npm run dev

# 生产构建
npm run build
```

构建后的文件位于 `web/dist` 目录。

## HTTPS配置（可选）

### 使用Let's Encrypt免费证书

1. **安装Certbot**
   ```bash
   # Ubuntu/Debian
   sudo apt install certbot python3-certbot-nginx
   
   # CentOS/RHEL
   sudo yum install certbot python3-certbot-nginx
   ```

2. **获取证书**
   ```bash
   sudo certbot --nginx -d your-domain.com
   ```

3. **自动续期**
   ```bash
   sudo crontab -e
   # 添加以下行
   0 12 * * * /usr/bin/certbot renew --quiet
   ```

### 修改Nginx配置

取消 `nginx.conf` 中HTTPS服务器配置部分的注释，并修改证书路径。

## 监控和日志

### 应用日志

- **后端日志**：通过Docker查看 `docker-compose logs backend`
- **Nginx日志**：`/var/log/nginx/access.log` 和 `/var/log/nginx/error.log`
- **数据库日志**：通过Docker查看 `docker-compose logs postgres`

### 健康检查

- 后端健康检查：`GET /health`
- 数据库连接检查：通过后端健康检查间接验证

## 备份和恢复

### 数据库备份

```bash
# 创建备份
docker exec ai-doc-postgres pg_dump -U postgres ai_doc_library > backup.sql

# 恢复备份
docker exec -i ai-doc-postgres psql -U postgres ai_doc_library < backup.sql
```

### 文件存储备份

```bash
# 备份存储目录
tar -czf storage_backup.tar.gz storage/

# 恢复存储目录
tar -xzf storage_backup.tar.gz
```

## 故障排除

### 常见问题

1. **前端无法访问**
   - 检查Nginx是否运行：`systemctl status nginx`
   - 检查防火墙设置：`sudo ufw status`
   - 检查Nginx配置：`nginx -t`

2. **API请求失败**
   - 检查后端服务是否运行：`docker-compose ps`
   - 检查后端日志：`docker-compose logs backend`
   - 检查网络连接：`curl http://localhost:8080/health`

3. **数据库连接失败**
   - 检查数据库服务是否运行：`docker-compose ps`
   - 检查数据库连接参数：`docker-compose exec backend env | grep DB_`
   - 检查数据库日志：`docker-compose logs postgres`

### 性能优化

1. **Nginx优化**
   - 启用gzip压缩
   - 配置缓存头
   - 调整worker进程数

2. **数据库优化**
   - 添加适当的索引
   - 定期清理过期数据
   - 调整PostgreSQL配置参数

3. **Go应用优化**
   - 调整GOMAXPROCS
   - 使用连接池
   - 优化SQL查询

## 安全建议

1. **更改默认密码**
   - 修改数据库默认密码
   - 使用强密码

2. **网络隔离**
   - 使用防火墙限制访问
   - 只开放必要的端口

3. **定期更新**
   - 定期更新Docker镜像
   - 定期更新系统包

4. **监控**
    - 设置资源使用监控
    - 设置日志监控和告警

## Python解析服务部署

### 架构说明

Python解析服务作为独立的微服务运行，专门负责PDF和DOCX文档的解析。它通过gRPC协议与主Go应用程序通信，提供高性能的文档解析功能。

### 主要功能

- **PDF文档解析**：提取文本内容、元数据、页数、文件大小等信息
- **DOCX文档解析**：提取文本内容、元数据、段落、表格、图片等信息
- **gRPC服务接口**：提供标准化的远程调用接口
- **健康检查**：支持服务健康状态监控

### 部署要求

- Python 3.8+
- 网络访问权限（用于gRPC通信）
- 系统内存至少512MB

### 配置说明

#### 服务配置

- 默认监听端口：50051
- 支持并发处理：默认10个工作线程
- 超时设置：30秒

#### Go服务配置

在Go服务的`internal/service/parser_service.go`中配置gRPC客户端连接：

```go
// 连接到gRPC服务
if err := service.grpcClient.Connect("localhost:50051"); err != nil {
    // 错误处理和回退逻辑
}
```

### 运维管理

#### 服务管理

```bash
# 启动服务
sudo systemctl start ai-doc-parser

# 停止服务
sudo systemctl stop ai-doc-parser

# 重启服务
sudo systemctl restart ai-doc-parser

# 查看日志
sudo journalctl -u ai-doc-parser -f
```

#### 健康检查

可以通过以下方式检查Python解析服务状态：

```bash
# 检查端口监听
netstat -tlnp | grep 50051

# 检查服务状态
sudo systemctl status ai-doc-parser

# 检查进程
ps aux | grep "service.server"
```

### 故障排除

#### 常见问题

1. **gRPC连接失败**
   - 检查Python解析服务是否运行
   - 验证网络连接和端口访问
   - 检查防火墙设置

2. **解析服务启动失败**
   - 检查Python依赖是否安装正确
   - 验证gRPC代码是否生成成功
   - 查看详细错误日志

3. **性能问题**
   - 调整工作线程数量
   - 检查系统资源使用情况
   - 优化解析算法

#### 日志分析

Python解析服务的日志位置：
- systemd日志：`journalctl -u ai-doc-parser`
- 应用日志：在服务启动目录的server.log文件中

### 扩展和定制

#### 支持新文档格式

可以通过实现新的解析器类来支持更多文档格式：

```python
# 示例：添加TXT文档支持
class TXTParser:
    def parse(self, file_path):
        # 实现解析逻辑
        pass
```

#### 性能优化

- 实现文档缓存机制
- 添加异步处理队列
- 使用连接池管理gRPC连接

#### 安全增强

- 添加TLS加密gRPC通信
- 实现文件类型和大小验证
- 添加访问控制和认证机制

### 备份和恢复

#### 配置备份

```bash
# 备份Python解析服务配置
tar -czf parser-config-backup.tar.gz python-parser-service/
```

#### 恢复步骤

```bash
# 恢复配置
tar -xzf parser-config-backup.tar.gz

# 重新安装依赖
cd python-parser-service
source venv/bin/activate
pip install -r requirements.txt

# 重启服务
sudo systemctl restart ai-doc-parser
```