# 测试覆盖率提升计划

**目标**: 将单元测试覆盖率从当前的~18%提升到80%以上

## 当前情况分析

| 层级 | 当前覆盖率 | 主要障碍 |
|-------|-----------|----------|
| internal/model | 96.5% | ✅ 已达标 |
| internal/middleware | 8.3% | 需要完整mock service |
| internal/service | 2.4% | 需要测试更多独立函数 |
| internal/repository | 0.5% | 需要数据库测试环境 |
| internal/handler | 0.0% | 需要mock多个service接口|
| internal/router | 0.0% | 需要mock handler/middleware |

## 实施方案

### 阶段一：独立函数测试（当前可立即执行，无新依赖）

**目标覆盖率提升**: +5-8%

#### 1.1 Service层独立函数测试
- ✅ 已完成：[`parser_service_utils_test.go`](internal/service/parser_service_utils_test.go)（min、元数据提取等）
- **待测试**:
  - `embedding_service.go` 的辅助函数
  - `cache_service.go` 的高频调用函数
  - `search_service.go` 的查询构造函数
  - `storage_service.go` 的路径处理函数

**预期覆盖率**: 从2.4% → 10-15%

#### 1.2 Handler层构造函数测试
虽然handler方法需要复杂mock，但handler的构造函数和简单逻辑可以测试：
```go
// 可测试的内容：
- NewXXXHandler()
- 简单的错误处理
- 响应格式化函数（如果有）
```

**预期覆盖率**: 从0.0% → 5-8%

### 阶段二：引入测试框架（需要修改go.mod，1-2小时）

**目标覆盖率提升**: +10-15%

#### 2.1 安装mock库
```bash
go get github.com/stretchr/testify
go get github.com/stretchr/testify/mock
```

#### 2.2 集成到go.mod
自动通过测试文件import使用。

#### 2.3 使用testify特性
- 使用`assert`和`require`替代手动判断
- 使用`suite`组织测试
- 使用`mock`创建灵活的mock对象

**预期覆盖率**: 全局提升5-10%

### 阶段三：核心业务逻辑测试（需要mock，2-3天）

**目标覆盖率提升**: +20-30%

#### 3.1 为Service层创建mock测试
使用testify/mockery后，可以轻松测试：
```go
// 示例：DocumentService的mock
mockRepo := &mocks.MockDocumentRepository{
    CreateFunc: func(ctx context.Context, doc *model.Document) (*model.Document, error) {
        return &model.Document{ID: uuid.New().String()}, nil
    },
}
```

可测试的service:
- `document_service.go` - 主要业务逻辑
- `search_service.go` - 搜索逻辑
- `user_service.go` - 用户管理逻辑

**预期覆盖率**: 从2.4% → 25-35%

#### 3.2 为Middleware层创建完整测试
使用service的mock，可以测试：
```go
func TestAuthMiddleware(t *testing.T) {
    mockUserService := &mocks.MockUserService{
        GetProfileFunc: func(ctx context.Context, id string) (*service.UserProfile, error) {
            return &service.UserProfile{}, nil
        },
    }
    
    mw := middleware.NewAuthMiddleware(mockUserService)
    
    // 测试各种场景：token无效、过期、权限等
}
```

**预期覆盖率**: 从8.3% → 40-60%

#### 3.3 为Handler层创建HTTP测试
使用httptest包（标准库，无需额外依赖）：
```go
func TestDocumentHandlerUpload(t *testing.T) {
    // 创建mock service
    // 创建测试HTTP请求
    // 验证响应
}
```

**预期覆盖率**: 从0.0% → 20-35%

### 阶段四：Repository层测试（需要test DB，3-4天）

#### 4.1 集成测试数据库
- 使用testcontainer/testfixtures
- 或者使用SQLite内存数据库进行测试
- 创建测试fixtures加载/清理函数

**预期覆盖率**: 从0.5% → 15-25%

## 具体实施步骤

### 立即可执行（无需依赖）

**今天可以完成（预计耗时2-3小时）：**

1. 为`cache_service.go`创建测试
   - 测试缓存key生成
   - 测试分页参数验证
   - 测试TTL计算

2. 为`embedding_service.go`创建测试
   - 测试embedding类型判断
   - 测试向量维度验证

3. 为`search_service.go`创建基础测试
   - 测试查询条件构建
   - 测试结果格式化

4. 为`health_service.go`创建测试
   - 测试组件健康检查逻辑
   - 测试状态聚合

**立即可以运行的命令：**
```bash
# 创建这3个测试文件后，执行：
go test -v -run "TestCache|TestEmbedding|TestSearch|TestHealth" ./internal/service/...
```

### 短期（1-2天，引入testify后）

1. ✅ 创建`go.mod`修改（如果需要）
2. ✅ 为`document_service`创建mock测试
3. ✅ 为`user_service`创建mock测试
4. ✅ 为`middleware`的logging和auth创建测试
5. ✅ 为5个主要handler创建基础HTTP测试

### 中期（3-5天，完善测试矩阵）

1. 为所有接口创建完整的mock
2. 使用testfixtures创建测试数据
3. 建立集成测试套件
4. 添加边界条件和错误场景测试

## 优先级排序

### 高优先级（立即执行）
1. ✅ Service层独立函数（已在进行）
2. ✅ Handler构造函数测试
3. ✅ 引入testify库

### 中优先级（1-2天内）
1. Service层业务逻辑测试（使用mock）
2. Middleware层完整测试（使用mock service）
3. Repository层基础测试（使用test DB）

### 低优先级（3-5天内）
1. Handler层完整HTTP测试
2. 集成测试和E2E测试
3. 性能测试和压力测试

## 预期最终覆盖率

| 层级 | 当前 | 阶段一后 | 引入testify后 | 完整测试后 |
|-------|------|-----------|-------------|--------------|
| internal/model | 96.5% | - | - | - |
| internal/service | 2.4% | 15% | 35% | 60%+ |
| internal/middleware | 8.3% | 15% | 45% | 70%+ |
| internal/repository | 0.5% | 5% | 20% | 50%+ |
| internal/handler | 0.0% | 5% | 25% | 60%+ |
| internal/router | 0.0% | 0% | 10% | 30%+ |
| **总体平均** | ~18% | ~8% | ~23% | ~38% |

**总覆盖率预期**: 约38%（阶段一）→ 约60%（testify）→ 约70%+（完整测试）

建议**快速路径**达到80%：
1. 完成阶段一（可立即达到38%）
2. 引入testify（提升到60%）
3. 集中测试5-7个关键service/handler（提升到80%+）
4. 其他层使用简化测试策略达到80%

## 成功标准

### 技术标准
- ✅ 所有测试用例能通过
- ✅ `make test` 无失败
- ✅ 整体覆盖率 ≥80%
- ✅ 核心业务逻辑覆盖率 ≥85%

### 质量标准
- ✅ 测试覆盖正常路径和边界条件
- ✅ 测试覆盖错误场景
- ✅ 测试代码清晰、可维护

---

**当前建议**：优先完成"阶段一"的4个测试文件，预计2-3小时完成，覆盖率可提升到38%左右。然后引入testify继续提升。

**需要我帮你立即执行吗？**
- 创建`cache_service_test.go`
- 创建`embedding_service_test.go`
- 创建`search_service_test.go`
- 创建`health_service_test.go`