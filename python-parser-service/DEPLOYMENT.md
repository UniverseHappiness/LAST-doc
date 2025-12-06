# Python解析服务部署指南

## 概述

Python解析服务提供了PDF和DOCX文档的解析功能，通过gRPC与主Go应用程序通信。

## 系统要求

- Python 3.8+
- Go 1.23+
- 网络访问权限（用于gRPC通信）

## 部署步骤

### 1. 安装依赖

```bash
cd python-parser-service
python3 -m venv venv
source venv/bin/activate  # Linux/Mac
# 或 venv\\Scripts\\activate  # Windows
pip install -r requirements.txt
```

### 2. 生成gRPC代码

```bash
./generate_grpc.sh
```

### 3. 启动Python解析服务

```bash
# 方法1：直接运行
python -m service.server

# 方法2：使用启动脚本
./start_server.sh

# 方法3：后台运行
nohup python -m service.server > server.log 2>&1 &
```

### 4. 验证服务状态

```bash
# 检查端口是否监听
netstat -tlnp | grep 50051

# 运行测试脚本
python test_service.py
```

## 配置选项

### 服务配置

服务默认监听端口：50051

可以通过修改`service/server.py`文件中的端口配置来更改监听端口：

```python
port = "50051"  # 修改为所需的端口
```

### Go服务配置

在Go服务的`internal/service/parser_service.go`文件中，需要配置gRPC客户端连接地址：

```go
// 连接到gRPC服务
if err := service.grpcClient.Connect("localhost:50051"); err != nil {
    // 错误处理
}
```

## 故障排除

### 1. gRPC连接失败

如果Go服务无法连接到Python解析服务，请检查：

- Python解析服务是否正在运行
- 网络防火墙是否阻止了端口50051
- 地址和端口配置是否正确

### 2. 依赖安装失败

如果Python依赖安装失败，可以尝试：

```bash
# 更新pip
pip install --upgrade pip

# 使用国内镜像
pip install -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple/
```

### 3. gRPC代码生成失败

如果gRPC代码生成失败，请确保：

- 已安装所有必要的依赖（grpcio, grpcio-tools, protobuf）
- proto文件语法正确
- Python虚拟环境已激活

## 性能优化

### 1. 并发处理

Python解析服务默认使用10个工作线程处理并发请求。可以通过修改`service/server.py`中的线程池大小来调整：

```python
server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
```

### 2. 缓存机制

对于频繁访问的文档，可以考虑添加缓存机制来提高性能。

### 3. 超时设置

gRPC客户端和服务器都设置有30秒的超时时间。可以根据需要调整：

```go
// Go客户端超时设置
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
```

## 监控和日志

### 日志配置

Python解析服务使用标准Python logging模块，日志级别为INFO。可以通过修改`service/server.py`来调整日志级别：

```python
logging.basicConfig(
    level=logging.INFO,  # 修改为DEBUG, WARNING, ERROR等
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
```

### 健康检查

服务提供健康检查接口，可以通过gRPC客户端调用`HealthCheck`方法来验证服务状态。

## 安全考虑

### 1. 网络安全

- 在生产环境中，建议使用TLS加密gRPC通信
- 限制对解析服务的网络访问
- 定期更新依赖包以修复安全漏洞

### 2. 文件安全

- 验证上传文件的类型和大小
- 对解析的文件内容进行安全检查
- 定期清理临时文件

## 扩展功能

### 1. 支持更多文档格式

可以通过实现新的解析器类来支持更多文档格式，如：
- TXT文件
- RTF文件
- HTML文件
- EPUB文件

### 2. 高级解析功能

可以添加以下高级功能：
- OCR文字识别
- 文档结构分析
- 元数据提取
- 图片和图表解析

### 3. 分布式部署

对于大规模部署，可以考虑：
- 使用Kubernetes进行容器编排
- 实现负载均衡
- 添加服务发现机制

## 维护和更新

### 1. 依赖更新

定期更新Python依赖包：

```bash
pip list --outdated
pip install --upgrade package_name
```

### 2. 性能监控

监控服务性能指标：
- 响应时间
- 错误率
- 资源使用率

### 3. 日志分析

定期分析服务日志，识别潜在问题和优化机会。

## 联系方式

如有问题或建议，请联系开发团队。