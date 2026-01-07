# AI技术文档库 (LAST-doc)

一个基于Go后端、Vue前端和Python解析服务的企业级智能文档管理系统，支持文档上传、解析、版本控制、元数据管理、智能检索和系统监控。

## 目录

- [项目概述](#项目概述)
- [功能特性](#功能特性)
- [架构设计](#架构设计)
- [技术栈](#技术栈)
- [快速开始](#快速开始)
- [部署说明](#部署说明)
- [API文档](#api文档)
- [MCP协议使用](#mcp协议使用)
- [监控和运维](#监控和运维)
- [测试指南](#测试指南)
- [开发指南](#开发指南)
- [贡献指南](#贡献指南)
- [许可证](#许可证)
- [更新说明](#更新说明)

## 项目概述

![项目展示](assets/main.gif)

AI技术文档库是一个综合性的企业级文档管理系统，采用前后端分离架构和微服务设计，提供高效的文档存储、检索、管理和监控解决方案。系统支持多格式文档解析、智能检索、MCP协议集成、用户权限管理、高可用部署和完整的系统监控功能。

### 核心价值

- 🚀 **高性能**: 双后端实例 + 负载均衡，支持高并发访问
- 🔍 **智能检索**: 支持关键词、语义和混合搜索，提供精准的文档检索能力
- 🔐 **安全可靠**: 完整的用户认证、权限控制、数据备份和高可用架构
- 📊 **全方位监控**: Prometheus + Grafana监控，实时掌握系统运行状态
- 🔌 **AI集成**: 支持MCP协议，便于与AI助手（如CoStrict IDE）深度集成
- 📦 **灵活部署**: 支持Docker Compose和Kubernetes两种部署方式

## 功能特性

### 文档管理
- **多格式支持**: 支持PDF、DOCX、Markdown等多种文档格式
- **文档解析**: 自动提取文档内容和元数据
- **版本控制**: 完整的文档版本管理，支持查看和恢复历史版本
- **元数据管理**: 支持自定义元数据，增强文档检索和分类能力
- **标签系统**: 灵活的标签系统，便于文档分类和检索

### 智能检索
- **关键词搜索**: 基于BM25算法的传统文本搜索
- **语义搜索**: 基于向量嵌入的语义相似度搜索
- **混合搜索**: 结合关键词和语义搜索的综合搜索
- **按库筛选**: 支持根据文档库进行精准筛选
- **搜索结果增强**: 显示文档所属库信息、相关度评分和内容片段

### AI集成
- **MCP协议支持**: 实现Model Context Protocol，便于与AI助手集成
- **API密钥管理**: 安全的API密钥创建、管理和删除功能
- **向量嵌入**: 集成OpenAI兼容的embedding服务，支持语义向量化
- **智能缓存**: 缓存搜索结果和embedding向量，提升响应速度

### 系统管理
- **用户管理**: 完整的用户注册、登录、认证和权限管理
- **系统监控**: 实时性能监控、日志管理和系统状态展示
- **数据备份**: 数据库备份和恢复机制
- **存储管理**: 支持本地存储和MinIO分布式存储

### 高可用特性
- **负载均衡**: Nginx反向代理和负载均衡
- **多实例部署**: 支持多个后端实例并行运行
- **健康检查**: 自动检测服务健康状态
- **故障恢复**: 自动重启失败服务
- **监控告警**: Prometheus + Grafana监控和告警

## 架构设计

### 系统架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                        用户访问层                                 │
│  ┌──────────────┐              ┌──────────────┐                 │
│  │  Web界面     │              │  MCP客户端    │                 │
│  │  (Vue.js)    │              │  (CoStrict)   │                 │
│  └──────────────┘              └──────────────┘                 │
└─────────────────────────────────────────────────────────────────┘
                                 ↓
┌─────────────────────────────────────────────────────────────────┐
│                       负载均衡层                                  │
│                      Nginx反向代理                                │
└─────────────────────────────────────────────────────────────────┘
                                 ↓
┌─────────────────────────────────────────────────────────────────┐
│                       应用服务层                                  │
│  ┌──────────────┐              ┌──────────────┐                 │
│  │  Backend-1   │              │  Backend-2   │                 │
│  │  (Go服务)    │              │  (Go服务)    │                 │
│  └──────────────┘              └──────────────┘                 │
│  ┌──────────────────────────────────────────────────┐          │
│  │         Python解析服务 (gRPC)                     │          │
│  │    PDF/DOCX文档解析                               │          │
│  └──────────────────────────────────────────────────┘          │
└─────────────────────────────────────────────────────────────────┘
                                 ↓
┌─────────────────────────────────────────────────────────────────┐
│                       数据存储层                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌─────────────┐         │
│  │  PostgreSQL  │  │  文件存储     │  │  MinIO      │         │
│  │   (数据库)   │  │  (本地存储)   │  │  (对象存储) │         │
│  └──────────────┘  └──────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────────┘
                                 ↓
┌─────────────────────────────────────────────────────────────────┐
│                       监控运维层                                  │
│  ┌──────────────┐              ┌──────────────┐                 │
│  │  Prometheus  │              │   Grafana     │                 │
│  │  (监控采集)   │              │  (可视化)     │                 │
│  └──────────────┘              └──────────────┘                 │
└─────────────────────────────────────────────────────────────────┘
```

### 组件说明

#### 前端服务
- **技术**: Vue.js 3 + Bootstrap 5 + Vite
- **职责**: 用户界面、交互逻辑、API调用
- **特性**: 响应式设计、Chart.js数据可视化

#### 后端服务
- **技术**: Go 1.24 + Gin + GORM
- **职责**: RESTful API、业务逻辑、数据访问
- **特性**: 双实例部署、负载均衡、健康检查

#### 文档解析服务
- **技术**: Python 3.8+ + gRPC
- **职责**: PDF和DOCX文档解析、内容提取
- **特性**: PyPDF2、pdfplumber、python-docx

#### 数据库
- **技术**: PostgreSQL 15 + pgvector
- **职责**: 数据存储、向量检索
- **特性**: 支持向量搜索、数据持久化

#### 反向代理
- **技术**: Nginx
- **职责**: 静态文件服务、API代理、负载均衡
- **特性**: 高性能、高并发处理

#### 监控系统
- **技术**: Prometheus + Grafana
- **职责**: 性能监控、数据可视化、告警
- **特性**: 实时监控、丰富的可视化面板

## 技术栈

### 后端技术栈
- **编程语言**: Go 1.24
- **Web框架**: Gin 1.9.1
- **ORM**: GORM 1.25.4
- **数据库**: PostgreSQL 15 (支持pgvector扩展)
- **gRPC**: Google gRPC 1.77.0
- **认证**: JWT (golang-jwt/jwt/v5)
- **向量嵌入**: go-openai 1.41.2
- **监控**: Prometheus client

### 前端技术栈
- **框架**: Vue.js 3.3.4
- **UI库**: Bootstrap 5.3.0
- **图标**: Bootstrap Icons 1.13.1
- **构建工具**: Vite 7.2.6
- **HTTP客户端**: Axios 1.13.2
- **可视化**: Chart.js 4.4.0

### 文档解析服务
- **编程语言**: Python 3.8+
- **gRPC框架**: gRPC
- **文档解析**:
  - PDF: PyPDF2、pdfplumber
  - DOCX: python-docx

### 部署和运维
- **容器化**: Docker & Docker Compose
- **编排**: Kubernetes (可选)
- **反向代理**: Nginx
- **监控**: Prometheus + Grafana
- **存储**: MinIO对象存储
- **进程管理**: systemd

## 快速开始

### 环境要求

- Go 1.24+
- Node.js 16+
- Python 3.8+
- PostgreSQL 15+
- Docker & Docker Compose (推荐)

### 使用Docker Compose启动（推荐）

#### 1. 克隆项目

```bash
git clone https://github.com/UniverseHappiness/LAST-doc.git
cd LAST-doc
```

#### 2. 快速启动

```bash
# 构建并启动所有服务
docker compose up -d

# 查看服务状态
docker compose ps

# 查看日志
docker compose logs -f
```

#### 3. 访问应用

- **前端界面**: http://localhost
- **后端API-1**: http://localhost:8080
- **后端API-2**: http://localhost:8080
- **监控面板**: http://localhost:9090 (Prometheus)
- **可视化**: http://localhost:3001 (Grafana，用户名/密码: admin/admin)
- **对象存储**: http://localhost:9000 (MinIO)

#### 4. 初始化账户

```bash
# 进入backend容器
docker compose exec backend sh

# 创建管理员账户
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "Admin@123"
  }'
```

### 手动启动

#### 1. 启动数据库

```bash
# 使用Docker启动PostgreSQL
docker run -d --name postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=ai_doc_library \
  -p 5432:5432 \
  postgres:15-alpine

# 初始化数据库
psql -h localhost -U postgres -d ai_doc_library -f scripts/init.sql
```

#### 2. 启动后端服务

```bash
# 编译Go应用
go build -o bin/ai-doc-library cmd/main.go

# 设置环境变量
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=ai_doc_library
export SERVER_PORT=8080
export STORAGE_DIR=./storage
export JWT_SECRET=your-secret-key-here
export ENABLE_HEALTH_CHECK=true
export GRPC_SERVER_HOST=localhost
export GRPC_SERVER_PORT=50051

# 创建存储目录
mkdir -p storage

# 启动服务
./bin/ai-doc-library
```

#### 3. 启动Python解析服务

```bash
cd python-parser-service

# 安装依赖
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt

# 生成gRPC代码
./generate_grpc.sh

# 启动服务
python -m service.server
```

#### 4. 构建和启动前端

```bash
cd web

# 安装依赖
npm install

# 开发模式启动
npm run dev

# 或构建生产版本
npm run build
```

## 部署说明

详细的部署指南请参考 [`DEPLOYMENT.md`](DEPLOYMENT.md) 和 [`docs/QUICK_START.md`](docs/QUICK_START.md)。

### Docker Compose部署

#### 标准部署

```bash
# 构建并启动所有服务
docker compose up -d

# 查看服务状态
docker compose ps

# 查看日志
docker compose logs -f

# 停止服务
docker compose down

# 停止并删除数据卷
docker compose down -v
```

#### 启用监控

```bash
# 启动监控服务（Prometheus + Grafana）
docker compose --profile monitoring up -d
```

### Kubernetes部署

项目提供完整的Kubernetes部署配置，位于 [`k8s/`](k8s/) 目录：

```bash
# 部署所有服务
./scripts/deploy-k8s.sh

# 或手动部署
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/postgres.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/ingress.yaml
```

### 部署对比

| 特性 | Docker Compose | Kubernetes |
|------|---------------|------------|
| 部署复杂度 | ⭐ 简单 | ⭐⭐⭐ 复杂 |
| 资源需求 | 低 | 中高 |
| 管理成本 | 低 | 中 |
| 扩展能力 | 手动 | 自动 |
| 高可用性 | 中 | 高 |
| 适用场景 | 开发/测试/小规模生产 | 大规模生产 |

## API文档

### 基础信息

- **基础URL**: `/api/v1`
- **认证方式**: JWT Token
- **数据格式**: JSON

### 主要端点

#### 健康检查
- **GET** `/health` - 检查服务健康状态
- **GET** `/metrics` - 获取Prometheus监控指标

#### 用户管理
- **POST** `/users/register` - 用户注册
- **POST** `/users/login` - 用户登录
- **GET** `/users/profile` - 获取用户信息
- **PUT** `/users/profile` - 更新用户信息

#### 文档管理
- **GET** `/documents` - 获取文档列表（支持分页、筛选）
- **POST** `/documents` - 上传新文档
- **GET** `/documents/{id}` - 获取文档详情
- **PUT** `/documents/{id}` - 更新文档信息
- **DELETE** `/documents/{id}` - 删除文档

#### 文档版本
- **GET** `/documents/{id}/versions` - 获取文档版本列表
- **POST** `/documents/{id}/versions` - 创建新版本
- **GET** `/documents/{id}/versions/{version}` - 获取特定版本详情
- **DELETE** `/documents/{id}/versions/{version}` - 删除特定版本

#### 文档元数据
- **GET** `/documents/{id}/metadata` - 获取文档元数据
- **POST** `/documents/{id}/metadata` - 添加文档元数据
- **PUT** `/documents/{id}/metadata/{key}` - 更新元数据
- **DELETE** `/documents/{id}/metadata/{key}` - 删除元数据

#### 文档检索
- **POST** `/search` - 执行文档搜索
- **GET** `/search` - 以GET方式执行文档搜索
- **POST** `/search/documents/{document_id}/versions/{version}/index` - 为指定文档版本构建搜索索引
- **GET** `/search/documents/{document_id}/index/status` - 获取文档的索引状态
- **DELETE** `/search/documents/{document_id}/index` - 删除指定文档的所有搜索索引
- **DELETE** `/search/documents/{document_id}/versions/{version}/index` - 删除指定文档版本的搜索索引
- **POST** `/search/clear-cache` - 清空搜索缓存

#### API密钥管理
- **GET** `/api-keys` - 获取用户的所有API密钥
- **POST** `/api-keys` - 创建新的API密钥
- **DELETE** `/api-keys/{id}` - 删除API密钥（软删除）
- **POST** `/api-keys/{id}/revoke` - 撤销API密钥

#### 系统监控
- **GET** `/monitor/metrics` - 获取系统性能指标
- **GET** `/monitor/logs` - 获取系统日志
- **POST** `/monitor/clear-cache` - 清空监控缓存

#### 数据备份
- **POST** `/backup/create` - 创建数据备份
- **GET** `/backup/list` - 列出所有备份
- **POST** `/backup/restore/{id}` - 恢复指定备份

### 检索功能

系统支持三种检索模式：

1. **关键词搜索**：基于关键词的传统文本搜索
2. **语义搜索**：基于向量的语义相似度搜索
3. **混合搜索**：结合关键词和语义搜索的综合搜索

#### 搜索请求格式

```json
{
  "query": "搜索关键词",
  "searchType": "keyword|semantic|hybrid",
  "page": 1,
  "size": 10,
  "filters": {
    "document_type": "pdf",
    "library": "技术文档"
  }
}
```

#### 搜索响应格式

```json
{
  "code": 200,
  "message": "搜索成功",
  "data": {
    "total": 100,
    "items": [
      {
        "id": "搜索结果ID",
        "document_id": "文档ID",
        "version": "文档版本",
        "content": "文档内容",
        "snippet": "内容片段",
        "score": 0.95,
        "content_type": "text",
        "section": "章节标题",
        "metadata": {
          "document_name": "文档名称",
          "document_type": "文档类型",
          "document_library": "文档库"
        }
      }
    ],
    "page": 1,
    "size": 10
  }
}
```

#### Embedding 服务配置

系统支持通过环境变量配置 OpenAI 兼容的 embedding 服务：

- `OPENAI_API_KEY`: OpenAI API 密钥
- `OPENAI_MODEL`: 使用的 embedding 模型（默认：text-embedding-ada-002）
- `OPENAI_BASE_URL`: 自定义 API 基础 URL（用于兼容其他服务）

当未提供 API 密钥时，系统将使用模拟 embedding 服务，确保基本功能可用。

## MCP协议使用

AI技术文档库支持MCP（Model Context Protocol）协议，可以让AI助手（如CoStrict IDE）直接访问和查询文档库中的内容。

### 获取API密钥

1. 登录AI技术文档库Web界面
2. 导航到"API密钥管理"页面
3. 点击"创建新密钥"按钮
4. 输入密钥名称并设置过期时间（可选）
5. 复制生成的API密钥并妥善保存

### 快速开始

#### 使用测试客户端

我们提供了Python测试客户端，方便您快速测试MCP功能：

```bash
# 安装依赖
pip install requests

# 运行测试（替换YOUR_API_KEY为您的实际API密钥）
python test_mcp_client.py --key YOUR_API_KEY --test
```

#### 基本配置

将以下配置添加到您的MCP客户端中：

```json
{
  "mcpServers": {
    "ai-doc-library": {
      "type": "streamable-http",
      "url": "http://localhost:8080/mcp",
      "headers": {
        "API_KEY": "your-api-key-here"
      }
    }
  }
}
```

### 可用的MCP工具

#### 1. search_documents
搜索技术文档，支持关键词和语义搜索。

**参数:**
- `query` (必需): 搜索查询关键词
- `types` (可选): 文档类型过滤器，如 ["pdf", "docx", "markdown"]
- `version` (可选): 文档版本过滤器
- `limit` (可选): 返回结果数量限制，默认为10
- `content_length` (可选): 每个搜索结果的内容片段最大字符数，默认为1000

**示例:**
```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "method": "tools/call",
  "params": {
    "name": "search_documents",
    "arguments": {
      "query": "Vue组件开发",
      "limit": 5,
      "content_length": 1000
    }
  }
}
```

#### 2. get_document_content
获取指定文档的详细内容。

**参数:**
- `document_id` (必需): 文档ID
- `version` (可选): 文档版本（推荐使用get_documents_by_library返回的版本ID）
- `start_position` (可选): 起始位置（字符位置）
- `end_position` (可选): 结束位置（字符位置）
- `query` (可选): 搜索关键词（用于定位内容位置）
- `content_length` (可选): 返回内容的最大字符数，默认为30000
- `smart_truncate` (可选): 是否启用智能截断模式，默认为true

**示例:**
```json
{
  "jsonrpc": "2.0",
  "id": "2",
  "method": "tools/call",
  "params": {
    "name": "get_document_content",
    "arguments": {
      "document_id": "doc-123",
      "content_length": 5000,
      "smart_truncate": true
    }
  }
}
```

#### 3. get_documents_by_library
根据所属库名称获取文档列表。

**参数:**
- `library` (必需): 库名称
- `page` (可选): 页码，默认为1
- `size` (可选): 每页数量，默认为10

**示例:**
```json
{
  "jsonrpc": "2.0",
  "id": "3",
  "method": "tools/call",
  "params": {
    "name": "get_documents_by_library",
    "arguments": {
      "library": "Eino",
      "page": 1,
      "size": 20
    }
  }
}
```

### 详细使用指南

更多详细的MCP使用说明，请参考：[MCP本地使用指南](docs/mcp_local_usage_guide.md)

### 安全注意事项

1. **保护API密钥**: 不要在代码中硬编码API密钥，使用环境变量或密钥管理系统
2. **权限控制**: 确保API密钥只有必要的权限，定期轮换密钥
3. **限流保护**: 合理设置API调用频率限制，防止滥用
4. **日志审计**: 记录所有API调用日志，便于安全审计

## 监控和运维

### 系统监控

系统集成了Prometheus + Grafana监控方案，提供全方位的系统监控能力。

#### Prometheus监控

- **访问地址**: http://localhost:9090
- **监控指标**:
  - 请求统计（请求总数、响应时间、错误率）
  - 数据库性能（查询时间、连接数、慢查询）
  - 缓存性能（命中率、缓存大小）
  - 系统资源（CPU、内存、磁盘IO）
  - 业务指标（文档数量、用户数量、搜索次数）

#### Grafana可视化

- **访问地址**: http://localhost:3001
- **默认账户**: admin / admin
- **可视化面板**:
  - 系统概览面板
  - 性能分析面板
  - 错误追踪面板
  - 业务统计面板

### 日志管理

- **日志路径**:
  - Nginx日志: `./logs/nginx/`
  - 后端日志: `./logs/backend/` 和 `./logs/backend2/`
  - PostgreSQL日志: 容器内部日志

- **日志级别**: 支持DEBUG、INFO、WARN、ERROR四种级别

### 数据备份

#### 创建备份

```bash
# 使用备份脚本
./scripts/backup-script.sh

# 或通过API
curl -X POST http://localhost:8080/api/v1/backup/create
```

#### 恢复备份

```bash
# 通过API恢复
curl -X POST http://localhost:8080/api/v1/backup/restore/{backup_id}
```

#### 定期备份

建议设置定期备份任务：

```bash
# 添加到crontab
0 2 * * * /path/to/scripts/backup-script.sh
```

### 健康检查

系统提供完整的健康检查机制：

```bash
# 检查服务状态
curl http://localhost:8080/health

# 检查后端-1
curl http://localhost:8080/health

# 检查后端-2
curl http://localhost:8081/health

# 检查数据库
docker compose exec postgres pg_isready -U postgres
```

### 性能优化

#### 数据库优化

- 使用索引加速查询
- 定期清理日志和历史数据
- 配置连接池优化连接管理

#### 缓存策略

- 启用Redis缓存热门数据
- 配置合理的缓存过期时间
- 实现多级缓存策略

#### 负载均衡

- 调整Nginx负载均衡算法
- 根据负载动态调整后端实例数量

## 测试指南

项目提供了完整的测试方案，涵盖单元测试、集成测试、效果测试和性能测试。

详细的测试指南请参考 [`TEST_GUIDE.md`](TEST_GUIDE.md)。

### 单元测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/model

# 生成覆盖率报告
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

# 使用Makefile
make test
```

### 集成测试

```bash
# 运行端到端测试
bash test_end_to_end.sh

# 运行API网关测试
bash test_api_gateway.sh

# 运行搜索功能测试
bash test_search_functionality.sh
```

### 性能测试

```bash
# 运行综合性能测试
bash test_performance_comprehensive.sh

# 运行可扩展性测试
bash test_scalability.sh

# 运行效果测试
bash test_effectiveness.sh
```

### 可靠性测试

```bash
# 运行可靠性测试
bash test_reliability.sh
```

### 测试覆盖率

当前测试覆盖率（部分达到测评要求）：

| 包路径 | 覆盖率 | 状态 |
|-------|--------|------|
| internal/model | 96.5% | ✅ 达标（>80%） |
| internal/middleware | 8.3% | ⚠️ 待提升 |
| internal/repository | 0.5% | ⚠️ 待提升 |
| internal/service | 0.8% | ⚠️ 待提升 |
| internal/handler | 0.0% | ⚠️ 待提升 |

详细的覆盖率报告请参考 [`TEST_COVERAGE_REPORT.md`](TEST_COVERAGE_REPORT.md)。

## 开发指南

### 项目结构

```
LAST-doc/
├── cmd/                      # Go应用程序入口
│   ├── main.go              # 主程序入口
│   └── main_test.go         # 主程序测试
├── internal/                 # 内部包，不对外暴露
│   ├── handler/             # HTTP请求处理器
│   ├── middleware/          # 中间件
│   ├── model/               # 数据模型
│   ├── repository/          # 数据访问层
│   ├── router/              # 路由配置
│   └── service/             # 业务逻辑层
├── proto/                    # Protocol Buffers定义
├── python-parser-service/    # Python解析服务
│   ├── proto/               # gRPC协议定义
│   └── service/             # 解析服务实现
├── scripts/                  # 数据库脚本和工具脚本
├── k8s/                      # Kubernetes部署配置
├── web/                      # 前端应用
│   ├── src/                 # 前端源码
│   │   ├── views/           # 页面组件
│   │   ├── composables/     # 组合式函数
│   │   └── utils/           # 工具函数
│   └── dist/                # 构建产物
├── docs/                     # 项目文档
├── configs/                  # 配置文件
├── logs/                     # 日志文件
├── storage/                  # 本地存储目录
├── backups/                  # 数据备份目录
├── docker compose.yml        # Docker Compose配置
├── Dockerfile               # Docker镜像配置
├── nginx.conf               # Nginx配置
├── Makefile                 # 构建脚本
├── go.mod                   # Go模块定义
└── README.md                # 项目说明
```

### 后端开发

#### 添加新的API端点

1. 在 [`internal/handler/`](internal/handler) 中添加处理器函数
2. 在 [`internal/router/router.go`](internal/router/router.go) 中注册路由
3. 在 [`internal/service/`](internal/service) 中实现业务逻辑
4. 在 [`internal/repository/`](internal/repository) 中实现数据访问
5. 在 [`internal/model/`](internal/model) 中定义数据模型

#### 数据库迁移

```bash
# 创建新的迁移脚本
# 将SQL文件添加到scripts/目录

# 应用迁移
psql -h localhost -U postgres -d ai_doc_library -f scripts/your_migration.sql

# 或使用Docker
docker compose exec postgres psql -U postgres -d ai_doc_library -f /docker-entrypoint-initdb.d/your_migration.sql
```

### 前端开发

#### 添加新页面

1. 在 [`web/src/views/`](web/src/views) 中创建新的Vue组件
2. 在 `web/src/router/index.js` 中添加路由配置
3. 在 [`web/src/utils/`](web/src/utils) 中添加API调用方法

#### 组件开发

遵循Vue.js 3的组合式API风格，使用Bootstrap 5组件库和Chart.js可视化库。

### Python解析服务开发

#### 添加新的文档格式支持

1. 在 [`python-parser-service/service/`](python-parser-service/service) 中创建新的解析器类
2. 在 [`python-parser-service/proto/`](python-parser-service/proto) 中更新Protocol Buffers定义
3. 在 [`python-parser-service/service/server.py`](python-parser-service/service/server.py) 中注册新解析器

#### 运行测试

```bash
cd python-parser-service

# 运行所有测试
python -m pytest

# 运行特定测试
python -m pytest test_service.py
```

### 代码规范

- Go代码遵循Go官方代码规范
- JavaScript代码遵循ESLint规范
- Python代码遵循PEP 8规范
- 使用 `make lint` 检查代码规范

### 构建命令

```bash
# 构建后端
make build

# 构建前端
cd web && npm run build

# 运行测试
make test

# 代码检查
make lint

# 生成覆盖率报告
make coverage
```

## 贡献指南

我们欢迎任何形式的贡献，包括但不限于：

- 🐛 修复bug
- ✨ 添加新功能
- 📖 改进文档
- 💡 提供建议
- 🧪 编写测试

### 提交规范

请遵循以下提交信息格式：

```
<类型>(<范围>): <描述>

[可选的详细描述]

[可选的引用]
```

类型包括：
- `feat`: 新功能
- `fix`: 修复bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建/工具相关
- `perf`: 性能优化
- `ci`: CI/CD相关

### 开发流程

1. Fork项目
2. 创建功能分支：`git checkout -b feature/your-feature-name`
3. 提交更改：`git commit -m "feat(scope): add new feature"`
4. 推送分支：`git push origin feature/your-feature-name`
5. 创建Pull Request

### 代码审查

所有代码提交都需要经过代码审查：
- 确保所有测试通过
- 代码符合项目规范
- 更新相关文档
- 无安全漏洞

## 许可证

本项目采用 MIT 许可证。详情请参阅 [LICENSE](LICENSE) 文件。

## 更新说明

### v0.6.0 (2026-01-06)

#### 新增功能
- **系统监控功能**: 完整的性能监控、日志管理和系统状态展示功能
- **性能报告API**: 提供系统性能指标和统计数据的API接口
- **API网关功能**: 实现请求路由、负载均衡和请求认证机制
- **监控前端页面**: 完善的系统监控界面，实时展示系统运行状态
- **双后端实例**: 支持多个后端实例并行运行，实现高可用
- **多数据库节点**: 支持主从复制，提高数据可靠性

#### 技术改进
- **性能监控服务**: 实现请求统计、数据库查询监控和缓存性能跟踪
- **日志管理系统**: 支持按时间和类型筛选日志，提供详细日志查看功能
- **网关路由机制**: 优化API请求处理流程，支持负载均衡和认证
- **缓存性能优化**: 实现智能缓存策略，提升系统响应速度
- **健康检查机制**: 完整的服务健康检查和自动恢复机制
- **Kubernetes部署**: 提供完整的K8s部署配置文件

#### 文档更新
- 完善README文档，新增监控和运维章节
- 添加Kubernetes部署指南
- 更新API文档，新增监控相关接口
- 添加测试指南，覆盖多种测试类型

---

### v0.5.0 (2026-01-03)

#### 新增功能
- **搜索结果增强**: 搜索文档时显示所属库信息，便于用户了解文档来源
- **按库获取文档**: 新增 `get_documents_by_library` MCP工具，支持根据库名称获取文档列表
- **分页查询支持**: 按库获取文档功能支持分页参数，可灵活控制返回结果数量

#### 技术改进
- **搜索结果模型扩展**: 在 `SearchResult` 和 `MCPSearchDocument` 中添加 `Library` 字段
- **元数据提取优化**: 从搜索索引元数据中提取库信息并显示在搜索结果中
- **MCP工具扩展**: 新增工具支持按库筛选文档，提供更精准的文档检索能力

---

### v0.4.0 (2025-12-15)

#### 新增功能
- **用户管理系统**: 完整的用户注册、登录、认证功能
- **API密钥管理**: 安全的API密钥创建、管理和删除功能
- **MCP协议支持**: 实现Model Context Protocol，支持AI助手集成
- **用户界面优化**: 重构MCP视图，移除冗余的API密钥管理功能
- **软删除机制**: API密钥采用软删除，提高数据安全性

#### 技术改进
- **认证中间件**: 实现基于JWT的用户认证和授权
- **数据库迁移**: 添加用户表和API密钥表
- **前端组件优化**: 修复模态框显示问题，改善用户体验
- **API错误处理**: 增强错误处理和调试日志

---

### v0.3.0 (2025-12-14)

#### 新增功能
- **OpenAI兼容的embedding服务**: 支持配置自定义embedding API
- **向量搜索**: 基于向量嵌入的语义搜索功能
- **混合搜索**: 结合关键词和语义搜索的综合搜索
- **搜索缓存**: 缓存搜索结果和embedding向量

#### 技术改进
- **embedding服务接口**: 添加OpenAI兼容的embedding服务接口
- **搜索性能优化**: 优化搜索算法和索引结构
- **缓存策略**: 实现智能缓存策略，提升搜索性能

---

### v0.2.0 (2025-12-06)

#### 新增功能
- **文档管理**: 支持多种格式文档的上传、下载、预览和管理
- **文档解析**: 自动解析PDF和DOCX文档内容，提取文本和元数据
- **版本控制**: 支持文档版本管理，可查看和恢复历史版本
- **元数据管理**: 支持自定义元数据，增强文档检索和分类能力
- **标签系统**: 灵活的标签系统，便于文档分类和检索
- **智能检索**: 支持关键词搜索、语义搜索和混合搜索

#### 技术架构
- **前后端分离**: Vue.js前端 + Go后端架构
- **微服务设计**: Python解析服务的gRPC微服务
- **数据库设计**: PostgreSQL + pgvector扩展支持向量搜索
- **容器化部署**: Docker + Docker Compose一键部署

---

### v0.1.0 (初始版本)

- 项目初始化
- 基础架构搭建

### 计划中的功能

- [ ] 文档协作功能
- [ ] 批量操作支持
- [ ] 更多文档格式支持（PPT、XLS等）
- [ ] 文档预览功能（在线预览）
- [ ] 高级权限管理（基于角色的权限控制）
- [ ] API限流和安全增强
- [ ] Elasticsearch集成（替代传统搜索）
- [ ] 移动端适配
- [ ] 国际化支持
- [ ] 插件系统

### 技术债务

- [ ] 提升单元测试覆盖率（handler、service、repository层）
- [ ] 优化数据库查询性能
- [ ] 实现更完善的错误处理机制
- [ ] 添加更多的集成测试
- [ ] 性能基准测试和优化
- [ ] 代码重构和模块化改进

## 常见问题

### 1. 如何重置管理员密码？

参考 [`docs/reset_admin_guide.md`](docs/reset_admin_guide.md)

### 2. 如何配置私有网络部署？

参考 [`docs/PRIVATE_DEPLOYMENT.md`](docs/PRIVATE_DEPLOYMENT.md)

### 3. 如何提高系统可靠性？

参考 [`docs/reliability_guide.md`](docs/reliability_guide.md)

### 4. 如何进行系统扩展？

参考 [`docs/scalability_guide.md`](docs/scalability_guide.md)

## 联系我们

- **项目地址**: https://github.com/UniverseHappiness/LAST-doc
- **问题反馈**: [GitHub Issues](https://github.com/UniverseHappiness/LAST-doc/issues)
- **文档**: [项目Wiki](https://github.com/UniverseHappiness/LAST-doc/tree/main/.cospec/wiki)

---

**感谢使用AI技术文档库！**