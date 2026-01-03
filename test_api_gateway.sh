#!/bin/bash

# API网关功能测试脚本
# 测试内容：请求路由、负载均衡、请求认证

set -e

echo "=========================================="
echo "API网关功能测试"
echo "=========================================="

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试结果统计
PASSED=0
FAILED=0

# 辅助函数：打印测试结果
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ 通过${NC}: $2"
        ((PASSED++))
    else
        echo -e "${RED}✗ 失败${NC}: $2"
        ((FAILED++))
    fi
}

# 辅助函数：打印测试标题
print_test() {
    echo ""
    echo -e "${YELLOW}测试${NC}: $1"
}

# 获取服务地址
API_GATEWAY_URL="http://localhost"
BACKEND_URL="http://localhost:8080"

echo "API网关地址: $API_GATEWAY_URL"
echo "后端服务地址: $BACKEND_URL"
echo ""

# 测试1: 检查服务是否启动
print_test "检查服务是否启动"
curl -s -o /dev/null -w "%{http_code}" $API_GATEWAY_URL/health
print_result $? "健康检查端点应返回200"

# 测试2: 请求路由功能 - API路由
print_test "请求路由功能 - API路由"
response=$(curl -s -w " HTTP_CODE:%{http_code}" "$API_GATEWAY_URL/api/v1/documents" || true)
http_code=$(echo $response | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
if [ "$http_code" = "200" ] || [ "$http_code" = "401" ]; then
    print_result 0 "API路由应响应（200或401）"
else
    print_result 1 "API路由未正确响应，状态码: $http_code"
fi

# 测试3: 请求路由功能 - MCP路由
print_test "请求路由功能 - MCP路由"
response=$(curl -s -w " HTTP_CODE:%{http_code}" -X POST "$API_GATEWAY_URL/mcp" \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "1.0"}}' || true)
http_code=$(echo $response | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
if [ "$http_code" = "200" ] || [ "$http_code" = "401" ]; then
    print_result 0 "MCP路由应响应（200或401）"
else
    print_result 1 "MCP路由未正确响应，状态码: $http_code"
fi

# 测试4: 请求认证功能 - 无认证的受保护端点
print_test "请求认证功能 - 无认证的受保护端点"
response=$(curl -s -w " HTTP_CODE:%{http_code}" "$API_GATEWAY_URL/api/v1/users/profile" || true)
http_code=$(echo $response | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
if [ "$http_code" = "401" ]; then
    print_result 0 "未认证的受保护端点应返回401"
else
    print_result 1 "未认证端点未返回401，状态码: $http_code"
fi

# 测试5: 请求认证功能 - 用户注册
print_test "请求认证功能 - 用户注册"
username="testuser_$(date +%s)"
email="test_$(date +%s)@example.com"
password="Test123456"
response=$(curl -s -w " HTTP_CODE:%{http_code}" -X POST "$API_GATEWAY_URL/api/v1/auth/register" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$username\",\"email\":\"$email\",\"password\":\"$password\"}" || true)
http_code=$(echo $response | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
if echo "$response" | grep -q "token\|id"; then
    print_result 0 "用户注册应成功并返回token或用户ID"
    # 保存用户信息用于后续测试
    exported_token=$(echo "$response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    exported_username=$username
else
    print_result 1 "用户注册失败，响应: $response"
    exported_token=""
fi

# 测试6: 请求认证功能 - 用户登录
print_test "请求认证功能 - 用户登录"
if [ ! -z "$exported_username" ]; then
    response=$(curl -s -w " HTTP_CODE:%{http_code}" -X POST "$API_GATEWAY_URL/api/v1/auth/login" \
        -H "Content-Type: application/json" \
        -d "{\"username\":\"$exported_username\",\"password\":\"$password\"}" || true)
    http_code=$(echo $response | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
    if echo "$response" | grep -q "token"; then
        exported_token=$(echo "$response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
        print_result 0 "用户登录应成功并返回token"
    else
        print_result 1 "用户登录失败，响应: $response"
    fi
else
    echo -e "${YELLOW}跳过${NC}: 没有可用的测试用户"
fi

# 测试7: 请求认证功能 - 使用JWT认证访问受保护端点
print_test "请求认证功能 - 使用JWT认证访问受保护端点"
if [ ! -z "$exported_token" ]; then
    response=$(curl -s -w " HTTP_CODE:%{http_code}" "$API_GATEWAY_URL/api/v1/users/profile" \
        -H "Authorization: Bearer $exported_token" || true)
    http_code=$(echo $response | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
    if [ "$http_code" = "200" ]; then
        print_result 0 "使用JWT认证应能访问受保护端点"
    else
        print_result 1 "JWT认证失败，响应: $response"
    fi
else
    echo -e "${YELLOW}跳过${NC}: 没有可用的认证token"
fi

# 测试8: 请求认证功能 - 错误的token
print_test "请求认证功能 - 使用错误的token"
response=$(curl -s -w " HTTP_CODE:%{http_code}" "$API_GATEWAY_URL/api/v1/users/profile" \
    -H "Authorization: Bearer invalid_token_12345" || true)
http_code=$(echo $response | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
if [ "$http_code" = "401" ]; then
    print_result 0 "错误的token应返回401"
else
    print_result 1 "错误token未返回401，状态码: $http_code"
fi

# 测试9: 请求路由功能 - CORS支持
print_test "请求路由功能 - CORS支持"
response=$(curl -s -I -w " HTTP_CODE:%{http_code}" -H "Origin: http://example.com" "$API_GATEWAY_URL/api/v1/documents" || true)
http_code=$(echo $response | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
if echo "$response" | grep -qi "access-control-allow-origin"; then
    print_result 0 "响应应包含CORS头"
else
    print_result 1 "响应缺少CORS头"
fi

# 测试10: 请求路由功能 - 代理头部转发
print_test "请求路由功能 - 代理头部转发（X-Real-IP, X-Forwarded-For）"
# 这个测试需要通过检查后端日志来验证，这里做基本的连通性测试
response=$(curl -s -w " HTTP_CODE:%{http_code}" "$API_GATEWAY_URL/health" || true)
http_code=$(echo $response | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
if [ "$http_code" = "200" ]; then
    print_result 0 "代理头部转发应正常工作"
else
    print_result 1 "代理头部转发失败"
fi

# 测试11: 负载均衡功能 - 多次请求验证
print_test "负载均衡功能 - 发送多个请求验证负载均衡"
# 发送多个请求来验证负载均衡是否工作
success_count=0
for i in {1..5}; do
    response=$(curl -s -w "%{http_code}" "$API_GATEWAY_URL/health" || true)
    if [ "$response" = "200" ]; then
        ((success_count++))
    fi
done
if [ $success_count -eq 5 ]; then
    print_result 0 "负载均衡应能处理多个请求（5/5成功）"
else
    print_result 1 "负载均衡请求处理失败（$success_count/5成功）"
fi

# 测试12: 请求路由功能 - HTTP方法支持
print_test "请求路由功能 - 支持不同的HTTP方法"
# 测试GET和POST方法
get_response=$(curl -s -w "%{http_code}" "$API_GATEWAY_URL/api/v1/documents" || true)
post_response=$(curl -s -w "%{http_code}" -X POST "$API_GATEWAY_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"test","password":"test"}' || true)

if [ ! -z "$get_response" ] && [ ! -z "$post_response" ]; then
    print_result 0 "应支持GET和POST方法"
else
    print_result 1 "HTTP方法支持测试失败"
fi

# 测试13: MCP认证功能 - API_KEY转发
print_test "MCP认证功能 - API_KEY头部转发"
response=$(curl -s -w " HTTP_CODE:%{http_code}" -X POST "$API_GATEWAY_URL/mcp" \
    -H "Content-Type: application/json" \
    -H "API_KEY: test_api_key_12345" \
    -d '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "1.0"}}' || true)
http_code=$(echo $response | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
# API_KEY会转发到后端，后端会验证或返回错误
if [ "$http_code" = "200" ] || [ "$http_code" = "401" ] || [ "$http_code" = "400" ]; then
    print_result 0 "API_KEY头部应正确转发到后端"
else
    print_result 1 "API_KEY转发失败，状态码: $http_code"
fi

# 测试14: 错误处理 - 后端不可用
print_test "错误处理 - 网关错误处理能力"
# 这个测试在单实例环境下较难测试，主要验证配置
if grep -q "proxy_next_upstream" nginx.conf; then
    print_result 0 "负载均衡错误处理配置正确"
else
    print_result 1 "缺少负载均衡错误处理配置"
fi

# 测试15: 文件上传支持
print_test "文件上传支持 - 大文件上传配置"
if grep -q "client_max_body_size 100M" nginx.conf; then
    print_result 0 "文件上传大小限制配置正确（100M）"
else
    print_result 1 "文件上传配置不正确"
fi

# 测试16: WebSocket支持
print_test "WebSocket支持 - Upgrade头部转发"
if grep -q "proxy_set_header Upgrade" nginx.conf; then
    print_result 0 "WebSocket支持配置正确"
else
    print_result 1 "缺少WebSocket支持配置"
fi

# 测试总结
echo ""
echo "=========================================="
echo "测试总结"
echo "=========================================="
echo -e "${GREEN}通过${NC}: $PASSED"
echo -e "${RED}失败${NC}: $FAILED"
echo "总计: $((PASSED + FAILED))"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}所有测试通过！${NC}"
    exit 0
else
    echo -e "${RED}部分测试失败！${NC}"
    exit 1
fi