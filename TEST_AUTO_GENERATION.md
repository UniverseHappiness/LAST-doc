# 自动生成测试使用指南

本文档介绍了 LAST-doc 项目中自动生成测试的工具和使用方法。

## 1. 工具概述

项目提供了自动生成测试的工具：

### 1.1 gotests
- **功能**：为 Go 源代码生成测试框架
- **适用场景**：为函数和方法生成基础测试结构
- **优势**：快速生成测试模板，减少重复工作

## 2. 使用方法

### 2.1 安装工具

如果尚未安装工具，可以运行以下命令：

```bash
# 安装 gotests
go get github.com/cweill/gotests/...
```

### 2.2 使用自动生成脚本

项目提供了便捷的自动生成脚本：

```bash
# 运行自动生成测试脚本
./scripts/generate-tests.sh
```

该脚本将：
1. 检查并安装必要的工具
2. 为所有服务层、仓库层、处理器层和路由层生成测试框架

### 2.3 手动使用工具

#### 2.3.1 生成特定文件的测试

```bash
# 为特定文件生成测试
gotests -all -w internal/service/document_service.go

# 为目录下所有文件生成测试
gotests -all -w internal/service/
```

## 3. 自定义测试模板

项目提供了自定义测试模板，位于 `test_templates/test_template.tmpl`。您可以修改此模板以适应项目的测试风格。

### 3.1 模板变量

- `{{.PackageName}}`：包名
- `{{.StructName}}`：结构体名
- `{{.MethodName}}`：方法名
- `{{.Imports}}`：必要的导入

### 3.2 自定义模板示例

```go
package {{.PackageName}}

import (
    "testing"
    "github.com/stretchr/testify/assert"
    {{if .Imports}}{{.Imports}}{{end}}
)

func Test{{.StructName}}_{{.MethodName}}(t *testing.T) {
    // 设置测试数据
    // TODO: 添加测试数据初始化代码
    
    // 调用被测试函数
    // TODO: 添加函数调用代码
    
    // 验证结果
    // TODO: 添加断言代码
    
    // 清理测试数据
    // TODO: 添加清理代码
}
```

## 4. 混合测试策略

项目采用混合测试策略：

### 4.1 自动生成部分
- 基础CRUD操作测试
- 简单的函数和方法测试
- 测试框架结构

### 4.2 手动编写部分
- 复杂业务逻辑测试
- 边界条件和异常场景测试
- 集成测试
- 性能测试
- 模拟对象（mocks）创建

### 4.3 示例工作流

1. 使用自动生成工具创建基础测试框架
2. 为简单测试用例添加断言和验证逻辑
3. 手动编写复杂业务逻辑的测试用例
4. 手动创建模拟对象隔离外部依赖
5. 运行测试并验证覆盖率

## 5. 手动创建模拟对象

由于项目不再使用自动生成的模拟对象，以下是手动创建模拟对象的示例：

### 5.1 创建模拟对象结构

```go
// MockDocumentService 是文档服务的模拟实现
type MockDocumentService struct {
    documents       map[string]*model.Document
    documentVersions map[string][]*model.DocumentVersion
    metadata        map[string]map[string]interface{}
}

func NewMockDocumentService() *MockDocumentService {
    return &MockDocumentService{
        documents:       make(map[string]*model.Document),
        documentVersions: make(map[string][]*model.DocumentVersion),
        metadata:        make(map[string]map[string]interface{}),
    }
}
```

### 5.2 实现接口方法

```go
func (m *MockDocumentService) UploadDocument(ctx context.Context, file *multipart.FileHeader, name, docType, version, library, description string, tags []string) (*model.Document, error) {
    docID := uuid.New().String()
    document := &model.Document{
        ID:          docID,
        Name:        name,
        Type:        model.DocumentType(docType),
        Version:     version,
        Tags:        tags,
        FilePath:    filepath.Join(os.TempDir(), file.Filename),
        FileSize:    file.Size,
        Status:      model.DocumentStatusCompleted,
        Description: description,
        Library:     library,
    }

    m.documents[docID] = document
    
    // 添加版本
    versionID := uuid.New().String()
    docVersion := &model.DocumentVersion{
        ID:          versionID,
        DocumentID:  docID,
        Version:     version,
        FilePath:    filepath.Join(os.TempDir(), file.Filename),
        FileSize:    file.Size,
        Status:      model.DocumentStatusCompleted,
        Description: description,
    }
    m.documentVersions[docID] = append(m.documentVersions[docID], docVersion)
    
    return document, nil
}

func (m *MockDocumentService) GetDocument(ctx context.Context, id string) (*model.Document, error) {
    if doc, exists := m.documents[id]; exists {
        return doc, nil
    }
    return nil, fmt.Errorf("文档不存在")
}
```

### 5.3 在测试中使用模拟对象

```go
func TestDocumentHandler_GetDocument(t *testing.T) {
    // 设置Gin为测试模式
    gin.SetMode(gin.TestMode)

    // 创建模拟文档服务
    mockService := NewMockDocumentService()
    handler := NewDocumentHandler(mockService)

    // 添加测试数据
    docID := uuid.New().String()
    doc := &model.Document{
        ID:          docID,
        Name:        "测试文档",
        Type:        model.DocumentTypeMarkdown,
        Version:     "1.0.0",
        Library:     "测试库",
        Description: "测试文档描述",
    }
    mockService.documents[docID] = doc

    // 创建测试路由
    router := gin.New()
    router.GET("/api/v1/documents/:id", handler.GetDocument)

    // 创建HTTP请求
    req, _ := http.NewRequest("GET", "/api/v1/documents/"+docID, nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 验证结果
    if w.Code != http.StatusOK {
        t.Errorf("预期状态码 %d, 实际 %d", http.StatusOK, w.Code)
    }
}
```

## 6. 最佳实践

### 5.1 测试命名约定
- 测试文件名：`*_test.go`
- 测试函数名：`Test<FunctionName>` 或 `Test<FeatureName>`
- 子测试使用 `t.Run` 组织

### 5.2 测试组织
```go
func TestDocumentService(t *testing.T) {
    tests := []struct {
        name    string
        want    interface{}
        wantErr bool
    }{
        // 测试用例
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 测试逻辑
        })
    }
}
```

### 5.3 使用断言库
项目推荐使用 `testify` 断言库，提高测试可读性：

```go
import "github.com/stretchr/testify/assert"

func TestSomething(t *testing.T) {
    result := DoSomething()
    
    assert.NotNil(t, result)
    assert.Equal(t, expected, result)
    assert.NoError(t, err)
}
```

## 6. 与现有测试集成

自动生成的测试将与项目中现有的测试共存：

### 6.1 运行所有测试
```bash
go test -v -cover ./...
```

### 6.2 运行特定包的测试
```bash
go test -v ./internal/service/...
```

### 6.3 生成覆盖率报告
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### 6.4 使用项目测试脚本
```bash
./scripts/test-runner.sh
```

## 7. 常见问题

### 7.1 工具未找到
确保已正确安装工具并将其添加到 PATH 中：

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### 7.2 模板未应用
确保模板文件位于正确路径：`test_templates/test_template.tmpl`

### 7.3 生成的测试需要修改
自动生成的测试只是基础框架，需要手动添加测试逻辑和断言。

## 8. 扩展阅读

- [gotests GitHub](https://github.com/cweill/gotests)
- [mockery GitHub](https://github.com/vektra/mockery)
- [testify GitHub](https://github.com/stretchr/testify)
- [Go 测试文档](https://go.dev/doc/testing)