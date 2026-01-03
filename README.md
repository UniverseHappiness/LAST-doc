# AI技术文档库 (LAST-doc)

一个基于Go后端、Vue前端和Python解析服务的智能文档管理系统，支持文档上传、解析、版本控制和元数据管理。

## 目录

- [项目概述](#项目概述)
- [功能特性](#功能特性)
- [架构设计](#架构设计)
- [技术栈](#技术栈)
- [快速开始](#快速开始)
- [部署说明](#部署说明)
- [API文档](#api文档)
- [MCP协议使用](#mcp协议使用)
- [开发指南](#开发指南)
- [测试](#测试)
- [贡献指南](#贡献指南)
- [许可证](#许可证)
- [更新说明](#更新说明)

## 项目概述

AI技术文档库是一个综合性的文档管理系统，旨在为企业或团队提供高效的文档存储、检索和管理解决方案。系统采用前后端分离架构，结合了现代Web技术栈和微服务设计，提供了强大的文档解析能力和灵活的版本控制功能。

## 功能特性

- **文档管理**：支持多种格式文档的上传、下载、预览和管理
- **文档解析**：自动解析PDF和DOCX文档内容，提取文本和元数据
- **版本控制**：支持文档版本管理，可查看和恢复历史版本
- **元数据管理**：支持自定义元数据，增强文档检索和分类能力
- **标签系统**：灵活的标签系统，便于文档分类和检索
- **智能检索**：支持关键词搜索、语义搜索和混合搜索，提供精准的文档检索能力
- **向量嵌入**：集成 OpenAI 兼容的 embedding 服务，支持语义向量化
- **MCP协议支持**：支持Model Context Protocol，便于与AI助手的集成
- **API密钥管理**：安全的API密钥创建和管理系统
- **缓存机制**：智能缓存策略，提高搜索性能和响应速度
- **权限管理**：基础的用户权限控制，确保文档安全
- **RESTful API**：标准化的API接口，便于集成和扩展

## 架构设计

系统采用微服务架构，主要包含以下组件：

1. **前端服务**：基于Vue.js 3和Bootstrap 5的Web应用
2. **后端服务**：Go语言实现的RESTful API服务
3. **文档解析服务**：Python实现的gRPC微服务，负责文档解析
4. **数据库**：PostgreSQL，存储文档元数据和内容
5. **文件存储**：本地文件系统存储文档文件
6. **反向代理**：Nginx，提供静态文件服务和API代理

### 系统架构图

```
┌─────────────┐     ┌──────────────┐     ┌─────────────────┐
│   Nginx     │────▶│   Go后端     │────▶│  PostgreSQL     │
│ (反向代理)  │     │   (API服务)  │     │   (数据库)      │
└─────────────┘     └──────────────┘     └─────────────────┘
       │                    │                     │
       │                    │                     │
       │                    ▼                     │
       │            ┌──────────────┐              │
       │            │ Python解析   │              │
       │            │ 服务(gRPC)   │              │
       │            └──────────────┘              │
       │                                          │
       ▼                                          │
┌─────────────┐                                  │
│   前端应用   │                                  │
│  (Vue.js)   │                                  │
└─────────────┘                                  │
       │                                          │
       └──────────────────────────────────────────┘
                        │
                        ▼
                ┌──────────────┐
                │  文件存储    │
                │ (本地存储)   │
                └──────────────┘
```

### 数据库设计

系统主要包含三个核心数据表：

1. **documents**：文档主表，存储文档基本信息
2. **document_versions**：文档版本表，管理文档历史版本
3. **document_metadata**：文档元数据表，存储扩展元数据信息

详细的表结构请参考 [`docs/draft.md`](docs/draft.md:1)。

## 技术栈

### 后端技术栈
- **编程语言**：Go 1.24
- **Web框架**：Gin
- **ORM**：GORM
- **数据库**：PostgreSQL
- **gRPC**：Google gRPC

### 前端技术栈
- **框架**：Vue.js 3
- **UI库**：Bootstrap 5
- **构建工具**：Vite
- **HTTP客户端**：Axios

### 文档解析服务
- **编程语言**：Python 3.8+
- **gRPC框架**：gRPC
- **文档解析**：
  - PDF：PyPDF2、pdfplumber
  - DOCX：python-docx

### 部署和运维
- **容器化**：Docker & Docker Compose
- **反向代理**：Nginx
- **进程管理**：systemd

## 快速开始

### 环境要求

- Go 1.24+
- Node.js 16+
- Python 3.8+
- PostgreSQL 15+
- Docker & Docker Compose (可选)

### 使用Docker Compose启动（推荐）

1. **克隆项目**
   ```bash
   git clone https://github.com/UniverseHappiness/LAST-doc.git
   cd LAST-doc
   ```

2. **启动所有服务**
   ```bash
   docker-compose up -d
   ```

3. **访问应用**
   - 前端界面：http://localhost
   - 后端API：http://localhost:8080

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

详细的部署指南请参考 [`DEPLOYMENT.md`](DEPLOYMENT.md:1)。

### Docker Compose部署

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 手动部署

手动部署各个组件的步骤请参考 [`DEPLOYMENT.md`](DEPLOYMENT.md:64) 中的"方式二：手动部署"部分。

### Python解析服务部署

Python解析服务的详细部署指南请参考 [`DEPLOYMENT.md`](DEPLOYMENT.md:347) 中的"Python解析服务部署"部分。

## API文档

### 基础信息

- **基础URL**：`/api/v1`
- **认证方式**：无（当前版本）
- **数据格式**：JSON

### 主要端点

#### 健康检查
- **GET** `/health` - 检查服务状态

#### 文档管理
- **GET** `/documents` - 获取文档列表
- **POST** `/documents` - 上传新文档
- **GET** `/documents/{id}` - 获取文档详情
- **PUT** `/documents/{id}` - 更新文档信息
- **DELETE** `/documents/{id}` - 删除文档

#### 文档版本
- **GET** `/documents/{id}/versions` - 获取文档版本列表
- **POST** `/documents/{id}/versions` - 创建新版本
- **GET** `/documents/{id}/versions/{version}` - 获取特定版本详情

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

### 检索功能

系统支持三种检索模式：

1. **关键词搜索**：基于BM25算法的传统文本搜索
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

- `OPENAI_API_KEY`：OpenAI API 密钥
- `OPENAI_MODEL`：使用的 embedding 模型（默认：text-embedding-ada-002）
- `OPENAI_BASE_URL`：自定义 API 基础 URL（用于兼容其他服务）

当未提供 API 密钥时，系统将使用模拟 embedding 服务，确保基本功能可用。

### 响应格式

成功响应：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    // 响应数据
  }
}
```

错误响应：
```json
{
  "code": 非0值,
  "message": "错误描述",
  "data": null
}
```

## MCP协议使用

AI技术文档库支持MCP（Model Context Protocol）协议，可以让AI助手（如Claude）直接访问和查询文档库中的内容。

### 获取API密钥

1. 登录AI技术文档库Web界面
2. 导航到"API密钥管理"页面
3. 点击"创建新密钥"按钮
4. 输入密钥名称并设置过期时间（可选）
5. 复制生成的API密钥并妥善保存

### 快速开始

#### 使用测试客户端

我们提供了一个Python测试客户端，方便您快速测试MCP功能：

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
- `types` (可选): 文档类型过滤器
- `version` (可选): 文档版本过滤器
- `limit` (可选): 返回结果数量限制

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
      "limit": 5
    }
  }
}
```

#### 2. get_document_content
获取指定文档的详细内容。

**参数:**
- `document_id` (必需): 文档ID
- `version` (可选): 文档版本

**示例:**
```json
{
  "jsonrpc": "2.0",
  "id": "2",
  "method": "tools/call",
  "params": {
    "name": "get_document_content",
    "arguments": {
      "document_id": "doc-123"
    }
  }
}
```

### 详细使用指南

更多详细的MCP使用说明，请参考：[MCP本地使用指南](docs/mcp_local_usage_guide.md)

### 安全注意事项

1. **保护API密钥**：不要在代码中硬编码API密钥
2. **权限控制**：确保API密钥只有必要的权限
3. **定期轮换**：定期更换API密钥以提高安全性

## 开发指南

### 项目结构

```
LAST-doc/
├── cmd/                    # Go应用程序入口
│   ├── main.go            # 主程序入口
│   └── main_test.go       # 主程序测试
├── internal/              # 内部包，不对外暴露
│   ├── handler/           # HTTP请求处理器
│   ├── model/             # 数据模型
│   ├── repository/        # 数据访问层
│   ├── router/            # 路由配置
│   └── service/           # 业务逻辑层
├── proto/                 # Protocol Buffers定义
├── python-parser-service/ # Python解析服务
│   ├── proto/             # gRPC协议定义
│   └── service/           # 解析服务实现
├── scripts/               # 数据库脚本
├── web/                   # 前端应用
│   └── src/               # 前端源码
├── docs/                  # 文档
├── docker-compose.yml     # Docker Compose配置
├── Dockerfile            # Docker镜像配置
├── nginx.conf            # Nginx配置
├── Makefile              # 构建脚本
└── README.md             # 项目说明
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
```

### 前端开发

#### 添加新页面

1. 在 [`web/src/views/`](web/src/views) 中创建新的Vue组件
2. 在 [`web/src/router/`](web/src/router) 中添加路由配置
3. 在 [`web/src/utils/documentService.js`](web/src/utils/documentService.js) 中添加API调用方法

#### 组件开发

遵循Vue.js 3的组合式API风格，使用Bootstrap 5组件库。

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

## 测试

### 后端测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/service

# 运行测试并生成覆盖率报告
go test -cover ./...
```

### 前端测试

```bash
cd web

# 运行测试
npm test

# 或使用Vite的测试插件
npm run test:unit
```

### 集成测试

```bash
# 使用测试脚本
./scripts/test-runner.sh
```

## 贡献指南

我们欢迎任何形式的贡献，包括但不限于：

- 修复bug
- 添加新功能
- 改进文档
- 提供建议

### 提交规范

请遵循以下提交信息格式：

```
<类型>(<范围>): <描述>

[可选的详细描述]

[可选的引用]
```

类型包括：
- feat：新功能
- fix：修复bug
- docs：文档更新
- style：代码格式调整
- refactor：重构
- test：测试相关
- chore：构建/工具相关

### 开发流程

1. Fork项目
2. 创建功能分支：`git checkout -b feature/your-feature-name`
3. 提交更改：`git commit -m "feat: add new feature"`
4. 推送分支：`git push origin feature/your-feature-name`
5. 创建Pull Request

### 代码规范

- Go代码遵循Go官方代码规范
- JavaScript代码遵循ESLint规范
- Python代码遵循PEP 8规范

## 更新说明

### v0.4.0 (2025-12-30)

#### 新增功能
- **搜索结果增强**：搜索文档时显示所属库信息，便于用户了解文档来源
- **按库获取文档**：新增 `get_documents_by_library` MCP工具，支持根据库名称获取文档列表
- **分页查询支持**：按库获取文档功能支持分页参数，可灵活控制返回结果数量

#### 技术改进
- **搜索结果模型扩展**：在 `SearchResult` 和 `MCPSearchDocument` 中添加 `Library` 字段
- **元数据提取优化**：从搜索索引元数据中提取库信息并显示在搜索结果中
- **MCP工具扩展**：新增工具支持按库筛选文档，提供更精准的文档检索能力

#### 使用示例

**搜索文档（显示所属库）：**
```json
{
  "query": "文档",
  "limit": 3
}
```

**按库获取文档列表：**
```json
{
  "library": "Eino",
  "page": 1,
  "size": 10
}
```

#### 文档更新
- 更新MCP工具文档，添加 `get_documents_by_library` 工具说明
- 完善README更新说明，记录v0.4.0版本功能

---

### v0.5.0 (2026-01-03)

#### 新增功能
- **系统监控功能**：完整的性能监控、日志管理和系统状态展示功能
- **性能报告API**：提供系统性能指标和统计数据的API接口
- **API网关功能**：实现请求路由、负载均衡和请求认证机制
- **MCP工具扩展**：新增按库筛选文档功能，提供更精准的文档检索能力
- **监控前端页面**：完善的系统监控界面，实时展示系统运行状态

#### 技术改进
- **性能监控服务**：实现请求统计、数据库查询监控和缓存性能跟踪
- **日志管理系统**：支持按时间和类型筛选日志，提供详细日志查看功能
- **网关路由机制**：优化API请求处理流程，支持负载均衡和认证
- **缓存性能优化**：实现智能缓存策略，提升系统响应速度
- **前端组件完善**：修复监控数据渲染问题，改善用户界面体验

#### 文档更新
- 更新MCP工具文档，添加 `get_documents_by_library` 工具说明
- 完善README更新说明，记录最新版本功能
- 添加系统监控使用文档和API文档

#### 已知问题
- 文档处理性能优化功能尚未实现
- 系统扩展性和可靠性功能待完善
- 部分高级部署配置需要进一步测试

---

### v0.4.0 (2025-12-16)

#### 新增功能
- **MCP工具扩展**：新增工具支持按库筛选文档，提供更精准的文档检索能力

#### 使用示例

**搜索文档（显示所属库）：**
```json
{
  "query": "文档",
  "limit": 3
}
```

**按库获取文档列表：**
```json
{
  "library": "Eino",
  "page": 1,
  "size": 10
}
```

#### 文档更新
- 更新MCP工具文档，添加 `get_documents_by_library` 工具说明
- 完善README更新说明，记录v0.4.0版本功能

---

### v0.3.0 (2025-12-15)

#### 新增功能
- **用户管理系统**：完整的用户注册、登录、认证功能
- **API密钥管理**：安全的API密钥创建、管理和删除功能
- **MCP协议支持**：实现Model Context Protocol，支持AI助手集成
- **用户界面优化**：重构MCP视图，移除冗余的API密钥管理功能
- **软删除机制**：API密钥采用软删除，提高数据安全性

#### 技术改进
- **认证中间件**：实现基于JWT的用户认证和授权
- **数据库迁移**：添加用户表和API密钥表
- **前端组件优化**：修复模态框显示问题，改善用户体验
- **API错误处理**：增强错误处理和调试日志

#### 文档更新
- 新增MCP本地使用指南
- 添加Python测试客户端示例
- 更新README文档，添加MCP协议使用说明

#### 已知问题
- 修改密码功能尚未完全实现
- MCP协议功能尚未经过全面测试
- 用户角色权限管理需要进一步完善

---

### v0.2.0 (2025-12-14)
- 新增 OpenAI 兼容的 embedding 服务
- 改进搜索功能，支持向量搜索和混合搜索
- 添加 embedding 服务接口和实现
- 支持 OpenAI API 密钥和模型配置
- 添加模拟 embedding 服务用于测试
- 优化搜索性能和准确性

### v0.1.0 (2025-12-06)
- 初始版本发布
- 实现基本的文档管理功能
- 支持PDF和DOCX等文档解析
- 实现文档版本控制
- 添加元数据管理功能
- 完善前端用户界面
- 提供完整的部署文档

### 计划中的功能
- [ ] 用户认证和权限管理
- [x] 文档全文搜索（向量搜索）
- [ ] 文档预览功能
- [ ] 批量操作支持
- [ ] 更多文档格式支持（PPT、XLS等）
- [ ] 文档协作功能
- [ ] API限流和安全增强
- [x] 性能优化和缓存机制

## 许可证

本项目采用 MIT 许可证。详情请参阅 [LICENSE](LICENSE) 文件。