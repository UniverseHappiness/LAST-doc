# LAST-doc

**Local & Authoritative Source of Truth for Documentation**

LAST-doc 是一个面向企业内网环境的私有化 AI 技术文档库系统，专为 CoStrict 编程 IDE 设计。通过 MCP 协议实时注入精准、版本对齐的技术文档与代码示例，在无法访问互联网的隔离网络中，为开发者提供高质量上下文支持，显著提升 AI 生成代码的准确性，减少"幻觉"问题。

## 🌟 核心价值

将外部知识内化为企业可控资产，在保障安全合规的同时赋能智能编程。

### 创新点

- **离线优先**：在隔离网络环境中提供完整的 AI 辅助编程能力
- **版本感知**：精细管理技术栈版本，确保文档与代码版本对齐
- **多格式兼容**：支持 PDF、Markdown、Swagger、Javadoc 等多种文档格式
- **生态集成**：无缝对接 CoStrict IDE，填补内网 AI 编程辅助工具的空白

## ✨ 核心功能

### 1. 多格式文档支持

支持手动上传和管理多种技术文档格式：

- **PDF 文档**：技术书籍、API 参考手册
- **Markdown 文档**：项目文档、技术博客
- **Swagger/OpenAPI**：REST API 文档
- **Javadoc**：Java API 文档
- **其他格式**：支持扩展更多文档格式

### 2. 版本精细管理

按技术栈版本精细管理文档：

- 支持同一技术的多版本文档并存
- 确保开发者获取与项目版本匹配的文档
- 避免版本不一致导致的代码错误

### 3. MCP 协议集成

通过 MCP（Model Context Protocol）协议与 CoStrict IDE 无缝集成：

- 实时注入精准的技术文档
- 提供版本对齐的代码示例
- 增强 AI 生成代码的上下文理解能力

### 4. 离线环境支持

专为企业内网和隔离环境设计：

- 完全离线运行，无需互联网连接
- 保障企业数据安全和合规性
- 提供高质量的本地化知识库

## 🏗️ 架构设计

LAST-doc 采用模块化架构设计，主要包括：

1. **文档存储层**：管理各类格式文档的存储和索引
2. **版本管理层**：处理技术栈版本的精细化管理
3. **MCP 接口层**：提供标准化的 MCP 协议接口
4. **文档解析层**：支持多种文档格式的解析和转换
5. **检索增强层**：提供高效的文档检索和上下文提取

## 🚀 快速开始

### 环境要求

- Go 1.21+
- 企业内网环境
- CoStrict IDE（用于集成）

### 安装部署

```bash
# 克隆仓库
git clone https://github.com/UniverseHappiness/LAST-doc.git
cd LAST-doc

# 构建项目
go build -o last-doc

# 运行服务
./last-doc
```

### 配置说明

在 `.env` 文件中配置必要参数：

```env
# 服务端口
PORT=8080

# 文档存储路径（请根据实际部署环境配置）
DOC_STORAGE_PATH=/var/lib/last-doc/docs

# MCP 服务配置
MCP_ENABLED=true
MCP_PORT=9000
```

## 📚 使用场景

### 场景 1：企业内网 AI 编程

在无法访问外网的企业内网环境中，开发者可以通过 LAST-doc 获取：

- 内部技术规范文档
- 第三方库的 API 文档
- 最佳实践和代码示例

### 场景 2：版本敏感项目

在需要严格版本管理的项目中：

- 确保使用正确版本的 API 文档
- 避免因版本不一致导致的集成问题
- 提供版本迁移指南和差异对比

### 场景 3：合规性要求严格

在金融、政府等对数据安全要求严格的领域：

- 所有技术文档本地化存储
- 不依赖外部互联网服务
- 完全可控的知识库管理

## 🔧 文档管理

### 上传文档

```bash
# 通过 API 上传文档
curl -X POST http://localhost:8080/api/docs/upload \
  -F "file=@/path/to/document.pdf" \
  -F "type=pdf" \
  -F "tech_stack=spring-boot" \
  -F "version=3.2.0"
```

### 查询文档

```bash
# 查询指定技术栈和版本的文档
curl http://localhost:8080/api/docs?tech_stack=spring-boot&version=3.2.0
```

### 版本管理

```bash
# 列出所有支持的技术栈版本
curl http://localhost:8080/api/versions

# 添加新版本
curl -X POST http://localhost:8080/api/versions \
  -H "Content-Type: application/json" \
  -d '{"tech_stack":"spring-boot","version":"3.2.0"}'
```

## 🔌 CoStrict IDE 集成

### MCP 协议配置

在 CoStrict IDE 中配置 MCP 服务器地址：

```json
{
  "mcp": {
    "servers": [
      {
        "name": "LAST-doc",
        "url": "http://localhost:9000",
        "enabled": true
      }
    ]
  }
}
```

### 使用示例

在 CoStrict IDE 中使用 AI 辅助编程时，LAST-doc 会自动：

1. 识别当前项目的技术栈和版本
2. 从本地文档库检索相关文档
3. 通过 MCP 协议注入精准上下文
4. 提升 AI 生成代码的准确性

## 📖 支持的文档类型

| 文档类型 | 格式 | 用途 |
|---------|------|------|
| API 文档 | Swagger/OpenAPI, Javadoc, JSDoc | REST API、类库 API 参考 |
| 技术手册 | PDF, Markdown | 框架使用指南、技术书籍 |
| 代码示例 | Markdown, 代码文件 | 最佳实践、代码模板 |
| 规范文档 | PDF, Markdown, Word | 企业技术规范、编码标准 |

## 🛡️ 安全与合规

- **数据本地化**：所有文档存储在企业内网，不外传
- **访问控制**：支持用户认证和权限管理
- **审计日志**：记录所有文档访问和操作历史
- **版本追溯**：完整的文档版本历史和变更记录

## 🤝 贡献指南

欢迎贡献代码、文档或提出改进建议！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 许可证

本项目采用 [LICENSE](LICENSE) 中指定的许可证。

## 📬 联系方式

- 项目主页：[https://github.com/UniverseHappiness/LAST-doc](https://github.com/UniverseHappiness/LAST-doc)
- Issue 跟踪：[https://github.com/UniverseHappiness/LAST-doc/issues](https://github.com/UniverseHappiness/LAST-doc/issues)

## 🙏 致谢

感谢所有为 LAST-doc 项目做出贡献的开发者和用户！

---

**LAST-doc** - 让企业内网的 AI 编程更加智能和可靠！
