# 文档管理系统 - 测试指导文档

## 概述

本文档为 LAST-doc 文档管理系统提供了一套完整的测试机制说明，确保项目的可测试性得到显著提升。该系统是一个 Go 后端项目，包含文档处理、版本控制等功能。

## 测试机制说明

### 1. 现有测试机制

#### 1.1 单元测试
项目使用 Go 标准测试框架进行单元测试，现有测试文件位于 `internal/service/document_service_test.go`，包含以下测试用例：

- `TestDocumentService` - 测试文档服务的基本功能
- `TestStorageService` - 测试存储服务的基本功能
- `TestParserService` - 测试解析服务的基本功能
- `TestMarkdownParser` - 测试 Markdown 解析器的功能
- `TestPDFParser` - 测试 PDF 解析器的功能
- `TestDocxParser` - 测试 DOCX 解析器的功能
- `TestSwaggerParser` - 测试 Swagger 解析器的功能
- `TestOpenAPIParser` - 测试 OpenAPI 解析器的功能
- `TestJavaDocParser` - 测试 JavaDoc 解析器的功能

#### 1.2 测试覆盖率
项目使用 Go 内置的覆盖率工具生成测试覆盖率报告，可通过以下命令生成：

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

#### 1.3 测试脚本
我们提供了一个综合测试脚本 `scripts/test-runner.sh`，用于自动化执行所有测试相关任务：

```bash
./scripts/test-runner.sh
```

该脚本包含以下测试步骤：
1. 运行单元测试
2. 生成测试覆盖率报告
3. 检查代码格式
4. 检查代码规范（如果安装了 golangci-lint）
5. 运行构建测试

### 2. 测试机制总结

#### 2.1 测试架构
项目采用分层架构，测试覆盖以下层次：

1. **服务层测试** - 位于 `internal/service` 目录，测试核心业务逻辑
   - 文档服务测试
   - 存储服务测试
   - 解析服务测试

2. **处理器层测试** - 位于 `internal/handler` 目录，需要添加 HTTP 请求处理测试

3. **仓库层测试** - 位于 `internal/repository` 目录，需要添加数据库操作测试

4. **路由层测试** - 位于 `internal/router` 目录，需要添加路由配置测试

#### 2.2 测试覆盖率目标
当前测试覆盖率约为 9.1%（主要覆盖服务层），建议目标覆盖率：

- **服务层**：80% 以上
- **处理器层**：70% 以上
- **仓库层**：60% 以上
- **路由层**：50% 以上
- **整体覆盖率**：70% 以上

#### 2.3 测试类型建议

1. **单元测试** - 测试单个函数或方法的功能
2. **集成测试** - 测试多个组件之间的交互
3. **API 测试** - 测试 HTTP API 端点
4. **性能测试** - 测试系统在高负载下的表现

### 3. 测试执行指南

#### 3.1 本地测试执行

##### 运行所有测试
```bash
# 运行所有测试并生成覆盖率报告
go test -v -cover ./...

# 运行特定包的测试
go test -v ./internal/service/...

# 运行特定测试函数
go test -v -run TestDocumentService ./internal/service/
```

##### 生成覆盖率报告
```bash
# 生成覆盖率数据
go test -coverprofile=coverage.out ./...

# 生成 HTML 覆盖率报告
go tool cover -html=coverage.out -o coverage.html

# 查看函数级覆盖率
go tool cover -func=coverage.out
```

##### 使用测试脚本
```bash
# 给脚本执行权限（仅需一次）
chmod +x scripts/test-runner.sh

# 运行完整测试流程
./scripts/test-runner.sh
```

#### 3.2 CI/CD 集成

建议在 CI/CD 流程中添加以下测试步骤：

```yaml
name: Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - uses: actions/setup-go@v2
      with:
        go-version: 1.23
    - name: Run tests
      run: |
        go test -v -cover ./...
        go tool cover -func=coverage.out
    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v1
```

### 4. 测试扩展建议

#### 4.1 添加缺失的测试

##### 处理器层测试
建议为 `internal/handler/document_handler.go` 添加测试，使用 `httptest` 包模拟 HTTP 请求。

##### 仓库层测试
建议为 `internal/repository` 目录下的文件添加数据库测试，使用测试数据库或模拟对象。

##### 路由层测试
建议为 `internal/router/router.go` 添加路由配置测试。

#### 4.2 性能测试

建议添加性能测试，特别是对于文档上传和解析功能。

运行性能测试：
```bash
go test -bench=. ./...
```

#### 4.3 测试数据管理

建议创建测试数据工厂，用于生成一致的测试数据。

### 5. 测试最佳实践

#### 5.1 测试命名约定

- 测试文件名：`*_test.go`
- 测试函数名：`Test<FunctionName>` 或 `Test<FeatureName>`
- 基准测试函数名：`Benchmark<FunctionName>` 或 `Benchmark<FeatureName>`
- 示例测试函数名：`Example<FunctionName>` 或 `Example<FeatureName>`

#### 5.2 测试组织

- 使用子测试组织相关测试用例

#### 5.3 测试断言

使用 `testing` 包的断言方法或第三方断言库（如 testify）。

#### 5.4 Mock 和 Stub

使用接口和依赖注入来模拟外部依赖。

### 6. 测试工具和资源

#### 6.1 推荐工具

- **testify** - 提供丰富的断言和模拟功能
- **gomock** - 生成模拟代码的工具
- **httptest** - Go 标准库中的 HTTP 测试工具
- **sqlmock** - 数据库模拟库
- **ginkgo** - BDD 风格的测试框架

#### 6.2 有用资源

- [Go 测试文档](https://go.dev/doc/testing)
- [Go 测试博客](https://go.dev/blog/testing)
- [testify GitHub](https://github.com/stretchr/testify)
- [Go 测试最佳实践](https://github.com/golang/go/wiki/CodeReviewComments#testing)

## 总结

LAST-doc 文档管理系统已经建立了一套基础的测试机制，包括单元测试、覆盖率报告和自动化测试脚本。通过遵循本文档的指导，团队可以进一步完善测试覆盖，提高代码质量和系统稳定性。建议优先添加处理器层和仓库层的测试，然后逐步扩展到集成测试和性能测试，以达到推荐的测试覆盖率目标。