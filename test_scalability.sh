#!/bin/bash

# AI技术文档库 - 扩展性测试脚本

echo "========================================="
echo "AI技术文档库 - 扩展性功能测试"
echo "========================================="
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试结果统计
PASSED=0
FAILED=0

# 测试函数
test_case() {
    local test_name=$1
    local test_command=$2
    
    echo -n "测试: $test_name ... "
    
    if eval "$test_command" > /dev/null 2>&1; then
        echo -e "${GREEN}通过${NC}"
        ((PASSED++))
        return 0
    else
        echo -e "${RED}失败${NC}"
        echo "  执行命令: $test_command"
        ((FAILED++))
        return 1
    fi
}

echo "1. 检查存储服务实现"
echo "----------------------------------------"

# 检查存储服务接口
test_case "StorageService接口存在" "grep -q 'type StorageService interface' internal/service/storage_service.go"
test_case "LocalStorageService实现存在" "grep -q 'type LocalStorageService struct' internal/service/storage_service.go"
test_case "S3StorageService实现存在" "grep -q 'type S3StorageService struct' internal/service/s3_storage_service.go"

# 检查存储工厂
test_case "存储工厂函数存在" "grep -q 'func NewStorageService' internal/service/storage_factory.go"
test_case "环境变量配置函数存在" "grep -q 'func NewStorageServiceFromEnv' internal/service/storage_factory.go"

echo ""
echo "2. 检查docker-compose配置"
echo "----------------------------------------"

# 检查后端服务配置
test_case "主后端服务配置" "grep -q 'container_name: ai-doc-backend-1' docker-compose.yml"
test_case "存储类型环境变量" "grep -q 'STORAGE_TYPE' docker-compose.yml"
test_case "节点ID环境变量" "grep -q 'NODE_ID' docker-compose.yml"

# 检查MinIO服务
test_case "MinIO服务配置" "grep -q 'image: minio/minio' docker-compose.yml"
test_case "MinIO存储卷配置" "grep -q 'minio_data:' docker-compose.yml"

# 检查扩展实例配置（注释状态）
test_case "扩展实例配置模板" "grep -q 'backend2:' docker-compose.yml"

echo ""
echo "3. 检查Nginx负载均衡配置"
echo "----------------------------------------"

# 检查upstream配置
test_case "Nginx upstream配置" "grep -q 'upstream ai_doc_backend' nginx.conf"
test_case "最少连接算法" "grep -q 'least_conn' nginx.conf"
test_case "后端服务器配置" "grep -q 'server backend:8080' nginx.conf"
test_case "扩展服务器配置模板" "grep -q 'server ai-doc-backend-2:8080' nginx.conf"

# 检查健康检查
test_case "健康检查配置" "grep -q 'max_fails' nginx.conf"
test_case "失败超时配置" "grep -q 'fail_timeout' nginx.conf"

echo ""
echo "4. 检查编译状态"
echo "----------------------------------------"

# 检查Go编译
if [ -f "./main" ]; then
    echo -e "测试: 主程序编译 ... ${GREEN}通过${NC}"
    ((PASSED++))
else
    echo -e "测试: 主程序编译 ... ${RED}失败${NC} (main文件不存在)"
    ((FAILED++))
fi

# 检查依赖
test_case "AWS SDK依赖" "grep -q 'aws-sdk-go' go.mod"

echo ""
echo "5. 检查存储类型常量"
echo "----------------------------------------"

test_case "本地存储类型" "grep -q 'StorageTypeLocal' internal/service/storage_service.go"
test_case "S3存储类型" "grep -q 'StorageTypeS3' internal/service/storage_service.go"
test_case "MinIO存储类型" "grep -q 'StorageTypeMinIO' internal/service/storage_service.go"

echo ""
echo "6. 检查main.go存储服务集成"
echo "----------------------------------------"

test_case "main.go使用存储工厂" "grep -q 'NewStorageServiceFromEnv' cmd/main.go"
test_case "存储服务错误处理" "grep -q 'Failed to create storage service' cmd/main.go"
test_case "存储类型日志" "grep -q 'Using S3/MinIO storage service' cmd/main.go"

echo ""
echo "7. 检查文档"
echo "----------------------------------------"

test_case "扩展性指南文档" "[ -f 'docs/scalability_guide.md' ]"
test_case "文档包含本地存储说明" "grep -q '本地存储' docs/scalability_guide.md"
test_case "文档包含MinIO说明" "grep -q 'MinIO' docs/scalability_guide.md"
test_case "文档包含横向扩展说明" "grep -q '横向扩展' docs/scalability_guide.md"

echo ""
echo "8. 检查存储服务接口完整性"
echo "----------------------------------------"

test_case "SaveFile方法" "grep -q 'SaveFile.*context' internal/service/storage_service.go"
test_case "DeleteFile方法" "grep -q 'DeleteFile.*context' internal/service/storage_service.go"
test_case "GetFile方法" "grep -q 'GetFile.*context' internal/service/storage_service.go"
test_case "GenerateFilePath方法" "grep -q 'GenerateFilePath' internal/service/storage_service.go"
test_case "FileExists方法" "grep -q 'FileExists.*context' internal/service/storage_service.go"
test_case "CopyFile方法" "grep -q 'CopyFile.*context' internal/service/storage_service.go"
test_case "MoveFile方法" "grep -q 'MoveFile.*context' internal/service/storage_service.go"
test_case "GetFileSize方法" "grep -q 'GetFileSize.*context' internal/service/storage_service.go"
test_case "GetFileStream方法" "grep -q 'GetFileStream.*context' internal/service/storage_service.go"
test_case "SaveFileStream方法" "grep -q 'SaveFileStream.*context' internal/service/storage_service.go"

echo ""
echo "========================================="
echo "测试结果汇总"
echo "========================================="
echo -e "${GREEN}通过: $PASSED${NC}"
echo -e "${RED}失败: $FAILED${NC}"
echo "总计: $((PASSED + FAILED))"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}所有测试通过！${NC}"
    exit 0
else
    echo -e "${RED}存在失败的测试，请检查实现${NC}"
    exit 1
fi