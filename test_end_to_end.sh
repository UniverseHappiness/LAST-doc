#!/bin/bash

# 端到端测试脚本
# 测试完整流程：服务部署 → 增加技术文档库 → CoStrict插件配置MCP → 对话框发起提问

set -e

echo "=========================================="
echo "AI技术文档库 - 端到端测试"
echo "=========================================="

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试结果统计
TOTAL_STEPS=0
PASSED_STEPS=0
FAILED_STEPS=0

# 辅助函数
print_step() {
    echo ""
    echo -e "${YELLOW}步骤 $((TOTAL_STEPS + 1))${NC}: $1"
    echo "----------------------------------------"
}

print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ 通过${NC}: $2"
        ((PASSED_STEPS++))
    else
        echo -e "${RED}✗ 失败${NC}: $2"
        ((FAILED_STEPS++))
    fi
    ((TOTAL_STEPS++))
}

# 全局变量
BASE_URL="http://localhost"
BACKEND_URL="http://localhost:8080"
TEST_USERNAME="testuser_$(date +%s)"
TEST_EMAIL="test_$(date +%s)@example.com"
TEST_PASSWORD="Test123456"
ADMIN_USERNAME="admin"
ADMIN_PASSWORD="admin"
API_KEY=""
JWT_TOKEN=""
DOCUMENT_ID=""

# 测试1: 检查服务部署状态
print_step "检查服务部署状态"
echo "检查Docker容器状态..."

if docker compose ps | grep -q "Up"; then
    echo "Docker服务正在运行"
    docker compose ps
    print_result 0 "Docker容器已成功部署"
else
    print_result 1 "Docker容器未运行，尝试启动..."
    docker compose up -d
    sleep 30
    if docker compose ps | grep -q "Up"; then
        print_result 0 "Docker容器已成功启动"
    else
        print_result 1 "Docker容器启动失败"
        exit 1
    fi
fi

# 等待服务就绪
echo "等待服务就绪..."
sleep 10

# 测试2: 检查Nginx前端服务
print_step "检查Nginx前端服务"
response=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/")
if [ "$response" = "200" ]; then
    print_result 0 "Nginx前端服务正常"
else
    print_result 1 "Nginx前端服务异常，状态码: $response"
fi

# 测试3: 检查后端服务健康状态
print_step "检查后端服务健康状态"
response=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/health")
if [ "$response" = "200" ]; then
    print_result 0 "后端服务健康检查通过"
else
    print_result 1 "后端服务健康检查失败，状态码: $response"
fi

# 测试4: 管理员登录
print_step "管理员登录"
response=$(curl -s -X POST "$BACKEND_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$ADMIN_USERNAME\",\"password\":\"$ADMIN_PASSWORD\"}")

if echo "$response" | grep -q "token"; then
    JWT_TOKEN=$(echo "$response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    echo "获取管理员Token: ${JWT_TOKEN:0:20}..."
    print_result 0 "管理员登录成功"
else
    print_result 1 "管理员登录失败: $response"
fi

# 测试5: 检查PostgreSQL数据库
print_step "检查PostgreSQL数据库"
if docker exec ai-doc-postgres pg_isready -U postgres > /dev/null 2>&1; then
    print_result 0 "PostgreSQL数据库连接正常"
else
    print_result 1 "PostgreSQL数据库连接失败"
fi

# 测试6: 检查Python解析服务
print_step "检查Python解析服务"
if docker exec ai-doc-python-parser python -c "import socket; s=socket.socket(); s.connect(('localhost', 50051)); s.close()" 2>/dev/null; then
    print_result 0 "Python解析服务正常运行"
else
    print_result 1 "Python解析服务异常"
fi

# 测试7: 创建测试用户
print_step "创建测试用户"
response=$(curl -s -X POST "$BACKEND_URL/api/v1/auth/register" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$TEST_USERNAME\",\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\"}")

if echo "$response" | grep -q "token\|id\|success"; then
    echo "测试用户创建成功: $TEST_USERNAME"
    print_result 0 "测试用户创建成功"
else
    print_result 1 "测试用户创建失败: $response"
fi

# 测试8: 测试用户登录
print_step "测试用户登录"
response=$(curl -s -X POST "$BACKEND_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$TEST_USERNAME\",\"password\":\"$TEST_PASSWORD\"}")

if echo "$response" | grep -q "token"; then
    JWT_TOKEN=$(echo "$response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    echo "获取用户Token: ${JWT_TOKEN:0:20}..."
    print_result 0 "测试用户登录成功"
else
    print_result 1 "测试用户登录失败: $response"
fi

# 测试9: 创建API密钥（CoStrict MCP配置）
print_step "创建API密钥（MCP配置）"
response=$(curl -s -X POST "$BACKEND_URL/api/v1/mcp/api-keys" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $JWT_TOKEN" \
    -d "{\"name\":\"test_endpoint_key\",\"expires_in_days\":30}")

if echo "$response" | grep -q "api_key\|key\|token"; then
    # 尝试提取API密钥，不同的响应格式
    API_KEY=$(echo "$response" | grep -o '"api_key":"[^"]*"' | cut -d'"' -f4 2>/dev/null || echo "$response" | grep -o '"key":"[^"]*"' | cut -d'"' -f4 2>/dev/null || echo "$response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4 2>/dev/null)
    
    if [ -z "$API_KEY" ]; then
        # 如果json格式不同，尝试另一种方式
        API_KEY=$(echo "$response" | python3 -c "import sys, json; print(json.load(sys.stdin).get('api_key') or json.load(sys.stdin).get('key') or json.load(sys.stdin).get('token') or '')" 2>/dev/null)
    fi
    
    if [ -n "$API_KEY" ]; then
        echo "API密钥创建成功: ${API_KEY:0:20}..."
        print_result 0 "API密钥创建成功"
    else
        echo "API密钥响应: $response"
        print_result 1 "无法提取API密钥"
    fi
else
    echo "API密钥创建响应: $response"
    print_result 1 "API密钥创建失败"
fi

# 测试10: 检查API密钥列表
print_step "检查API密钥列表"
response=$(curl -s -X GET "$BACKEND_URL/api/v1/mcp/api-keys" \
    -H "Authorization: Bearer $JWT_TOKEN")

if echo "$response" | grep -q "api_keys\|keys\|data"; then
    echo "API密钥列表: $response"
    print_result 0 "API密钥列表获取成功"
else
    print_result 1 "API密钥列表获取失败"
fi

# 测试11: 准备测试文档
print_step "准备测试文档"
TEST_DOC_DIR="./test_documents"
mkdir -p "$TEST_DOC_DIR"

# 创建一个测试PDF文件（使用简单的文本文件模拟）
TEST_DOC="$TEST_DOC_DIR/test_document.pdf"
cat > "$TEST_DOC.txt" << 'EOF'
这是一个测试文档，用于演示AI技术文档库的功能。

文档内容：
1. 系统架构设计
   - 前端：Vue.js 3
   - 后端：Go 1.24
   - 数据库：PostgreSQL
   - 解析服务：Python

2. 核心功能
   - 文档上传与管理
   - 文档解析（PDF、DOCX）
   - 智能检索（关键词、语义、混合）
   - 版本控制
   - MCP协议支持

3. 技术特点
   - 微服务架构
   - RESTful API
   - 响应式设计
   - 高可用性
EOF

# 转换为PDF格式（如果系统有相关工具）
if command -v pandoc &> /dev/null; then
    pandoc "$TEST_DOC.txt" -o "$TEST_DOC"
    echo "使用pandoc创建PDF文档"
elif command -v wkhtmltopdf &> /dev/null; then
    echo "<html><body><pre>$(cat "$TEST_DOC.txt")</pre></body></html>" > "$TEST_DOC.html"
    wkhtmltopdf "$TEST_DOC.html" "$TEST_DOC"
    echo "使用wkhtmltopdf创建PDF文档"
else
    # 使用markdown文件作为替代
    mv "$TEST_DOC.txt" "$TEST_DOC.md"
    TEST_DOC="$TEST_DOC.md"
    echo "使用Markdown文档进行测试"
fi

print_result 0 "测试文档准备完成: $TEST_DOC"

# 测试12: 上传测试文档
print_step "上传测试文档"
if [ -f "$TEST_DOC" ]; then
    response=$(curl -s -X POST "$BACKEND_URL/api/v1/documents/upload" \
        -H "Authorization: Bearer $JWT_TOKEN" \
        -F "file=@$TEST_DOC" \
        -F "title=测试文档_端到端测试" \
        -F "category=技术文档" \
        -F "description=用于端到端测试的示例文档")
    
    echo "上传响应: $response"
    
    if echo "$response" | grep -q "id\|document_id\|success"; then
        # 尝试提取文档ID
        DOCUMENT_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4 2>/dev/null | head -1 || echo "$response" | grep -o '"document_id":"[^"]*"' | cut -d'"' -f4 2>/dev/null)
        
        if [ -n "$DOCUMENT_ID" ]; then
            echo "文档上传成功，ID: $DOCUMENT_ID"
            print_result 0 "文档上传成功"
        else
            print_result 1 "无法提取文档ID"
        fi
    else
        print_result 1 "文档上传失败"
    fi
else
    print_result 1 "测试文档文件不存在: $TEST_DOC"
fi

# 测试13: 等待文档解析和索引构建
print_step "等待文档解析和索引构建"
echo "等待30秒让文档完成解析和索引构建..."
sleep 30

# 测试14: 检查文档列表
print_step "检查文档列表"
response=$(curl -s -X GET "$BACKEND_URL/api/v1/documents" \
    -H "Authorization: Bearer $JWT_TOKEN")

if echo "$response" | grep -q "documents\|data"; then
    echo "文档数量: $(echo "$response" | grep -o '"id"' | wc -l)"
    print_result 0 "文档列表获取成功"
else
    print_result 1 "文档列表获取失败"
fi

# 测试15: 测试MCP连接（CoStrict插件配置验证）
print_step "测试MCP连接"
if [ -n "$API_KEY" ]; then
    response=$(curl -s -X GET "$BACKEND_URL/api/v1/mcp/test" \
        -H "API_KEY: $API_KEY")
    
    echo "MCP测试响应: $response"
    
    if echo "$response" | grep -q "success\|connected\|ready"; then
        print_result 0 "MCP连接测试成功"
    else
        print_result 1 "MCP连接测试失败"
    fi
else
    print_result 1 "API密钥为空，无法测试MCP连接"
fi

# 测试16: 获取MCP工具列表
print_step "获取MCP工具列表"
if [ -n "$API_KEY" ]; then
    response=$(curl -s -X POST "$BACKEND_URL/mcp" \
        -H "Content-Type: application/json" \
        -H "API_KEY: $API_KEY" \
        -d '{
            "jsonrpc": "2.0",
            "id": "tools_list",
            "method": "tools/list",
            "params": {}
        }')
    
    echo "MCP工具列表响应: ${response:0:200}..."
    
    if echo "$response" | grep -q "tools\|search_documents\|get_document_content"; then
        print_result 0 "MCP工具列表获取成功"
    else
        print_result 1 "MCP工具列表获取失败"
    fi
else
    print_result 1 "API密钥为空"
fi

# 测试17: 对话框提问1 - 搜索文档（CoStrict使用MCP工具）
print_step "对话框提问1 - 搜索文档（使用MCP）"
if [ -n "$API_KEY" ]; then
    response=$(curl -s -X POST "$BACKEND_URL/mcp" \
        -H "Content-Type: application/json" \
        -H "API_KEY: $API_KEY" \
        -d '{
            "jsonrpc": "2.0",
            "id": "search_test",
            "method": "tools/call",
            "params": {
                "name": "search_documents",
                "arguments": {
                    "query": "Vue组件开发",
                    "limit": 5
                }
            }
        }')
    
    echo "搜索结果响应: ${response:0:300}..."
    
    if echo "$response" | grep -q "result\|documents\|data"; then
        print_result 0 "MCP文档搜索成功"
    else
        print_result 1 "MCP文档搜索失败"
    fi
else
    print_result 1 "API密钥为空"
fi

# 测试18: 对话框提问2 - 获取文档内容
print_step "对话框提问2 - 获取文档内容"
if [ -n "$API_KEY" ] && [ -n "$DOCUMENT_ID" ]; then
    response=$(curl -s -X POST "$BACKEND_URL/mcp" \
        -H "Content-Type: application/json" \
        -H "API_KEY: $API_KEY" \
        -d "{
            \"jsonrpc\": \"2.0\",
            \"id\": \"get_content\",
            \"method\": \"tools/call\",
            \"params\": {
                \"name\": \"get_document_content\",
                \"arguments\": {
                    \"document_id\": \"$DOCUMENT_ID\"
                }
            }
        }")
    
    echo "文档内容响应: ${response:0:300}..."
    
    if echo "$response" | grep -q "content\|text\|data"; then
        print_result 0 "获取文档内容成功"
    else
        print_result 1 "获取文档内容失败"
    fi
else
    print_result 1 "API密钥或文档ID为空"
fi

# 测试19: 测试搜索API（直接调用）
print_step "测试搜索API（直接调用）"
response=$(curl -s -X POST "$BACKEND_URL/api/v1/search" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $JWT_TOKEN" \
    -d '{
        "query": "PostgreSQL",
        "searchType": "keyword",
        "page": 1,
        "size": 10
    }')

if echo "$response" | grep -q "results\|documents\|data"; then
    print_result 0 "搜索API调用成功"
else
    print_result 1 "搜索API调用失败"
fi

# 测试20: 测试用户资料获取
print_step "测试用户资料获取"
response=$(curl -s -X GET "$BACKEND_URL/api/v1/users/profile" \
    -H "Authorization: Bearer $JWT_TOKEN")

if echo "$response" | grep -q "username\|email\|id"; then
    echo "用户资料: ${response:0:200}..."
    print_result 0 "用户资料获取成功"
else
    print_result 1 "用户资料获取失败"
fi

# 测试21: 检查文档搜索性能
print_step "检查文档搜索性能（要求<1s）"
start_time=$(date +%s.%N)

response=$(curl -s -X POST "$BACKEND_URL/api/v1/search" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $JWT_TOKEN" \
    -d '{
        "query": "架构",
        "searchType": "keyword",
        "page": 1,
        "size": 5
    }')

end_time=$(date +%s.%N)
elapsed_time=$(echo "$end_time - $start_time" | bc)

echo "搜索耗时: ${elapsed_time}秒"

if (( $(echo "$elapsed_time < 1.0" | bc -l) )); then
    print_result 0 "搜索性能满足要求（<1秒）"
else
    print_result 1 "搜索性能不满足要求（>=1秒）"
fi

# 测试22: 测试文档上传性能
print_step "测试文档上传性能（要求<1min）"
start_time=$(date +%s.%N)

# 创建一个较大的测试文件
LARGE_TEST_DOC="$TEST_DOC_DIR/large_test.txt"
dd if=/dev/zero of="$LARGE_TEST_DOC" bs=1M count=10 2>/dev/null

response=$(curl -s -X POST "$BACKEND_URL/api/v1/documents/upload" \
    -H "Authorization: Bearer $JWT_TOKEN" \
    -F "file=@$LARGE_TEST_DOC" \
    -F "title=性能测试文档")
    # 块操作结束

end_time=$(date +%s.%N)
elapsed_time=$(echo "$end_time - $start_time" | bc)
elapsed_seconds=$(echo "$elapsed_time / 60" | bc)  # 转换为分钟

echo "上传耗时: ${elapsed_time}秒 (${elapsed_seconds}分钟)"

if (( $(echo "$elapsed_time < 60" | bc -l) )); then
    print_result 0 "文档上传性能满足要求（<1分钟）"
else
    print_result 1 "文档上传性能不满足要求（>=1分钟）"
fi

# 测试总结
echo ""
echo "=========================================="
echo "端到端测试总结"
echo "=========================================="
echo -e "${GREEN}通过步骤${NC}: $PASSED_STEPS"
echo -e "${RED}失败步骤${NC}: $FAILED_STEPS"
echo "总计步骤: $TOTAL_STEPS"

if [ $FAILED_STEPS -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✓ 所有测试通过！${NC}"
    echo ""
    echo "测试信息："
    echo "- 测试用户: $TEST_USERNAME"
    echo "- API密钥: ${API_KEY:0:20}..."
    echo "- JWT Token: ${JWT_TOKEN:0:20}..."
    echo "- 文档ID: $DOCUMENT_ID"
    echo ""
    echo "这些凭证可用于后续的有效性测试。"
    exit 0
else
    echo ""
    echo -e "${RED}✗ 部分测试失败，请检查日志${NC}"
    exit 1
fi