#!/bin/bash

# 综合性能测试脚本
# 测评方案要求：
# 1. 文档更新性能：确保文档更新在1分钟内完成
# 2. 文档检索性能：确保文档检索在1秒内返回结果

set -e

echo "=========================================="
echo "AI技术文档库 - 综合性能测试"
echo "=========================================="

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 测试结果统计
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 性能数据存储
declare -A SEARCH_TIMES
declare -A UPLOAD_TIMES
declare -A UPDATE_TIMES

# 测试配置
UPDATE_TIME_LIMIT=60  # 1分钟（60秒）
SEARCH_TIME_LIMIT=1   # 1秒
BASE_URL="http://localhost"
BACKEND_URL="http://localhost:8080"
JWT_TOKEN=""

# 辅助函数
print_header() {
    echo ""
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

print_test() {
    echo ""
    echo -e "${YELLOW}测试 $((TOTAL_TESTS + 1))${NC}: $1"
    echo "----------------------------------------"
    ((TOTAL_TESTS++)) || true
}

print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ 通过${NC}: $2"
        ((PASSED_TESTS++)) || true
    else
        echo -e "${RED}✗ 失败${NC}: $2"
        ((FAILED_TESTS++)) || true
    fi
}

# 获取JWT Token
print_header "环境准备"

print_test "获取用户Token"
response=$(curl -s -X POST "$BACKEND_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"123456"}')

if echo "$response" | grep -q "token"; then
    JWT_TOKEN=$(echo "$response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    echo "Token获取成功: ${JWT_TOKEN:0:20}..."
    print_result 0 "Token获取成功"
else
    print_result 1 "Token获取失败"
    exit 1
fi

# 准备测试文档
print_test "准备测试文档"
TEST_DOC_DIR="./test_documents"
mkdir -p "$TEST_DOC_DIR"

# 创建不同大小的测试文档
echo "创建测试文档..."

# 小文档（1MB）
dd if=/dev/zero of="$TEST_DOC_DIR/small_test.txt" bs=1M count=1 2>/dev/null
echo "创建小文档: 1MB"

# 中文档（5MB）
dd if=/dev/zero of="$TEST_DOC_DIR/medium_test.txt" bs=1M count=5 2>/dev/null
echo "创建中文档: 5MB"

# 大文档（10MB）
dd if=/dev/zero of="$TEST_DOC_DIR/large_test.txt" bs=1M count=10 2>/dev/null
echo "创建大文档: 10MB"

# 创建带有实际内容的文档
cat > "$TEST_DOC_DIR/content_test.md" << 'EOF'
# AI技术文档库性能测试文档

本文档用于测试系统的文档解析、索引构建和搜索性能。

## 技术栈

- 前端: Vue.js 3
- 后端: Go 1.24
- 数据库: PostgreSQL 15
- 解析服务: Python 3.8+
- 容器化: Docker & Docker Compose

## 核心功能

### 1. 文档管理
- 上传与管理
- 版本控制
- 元数据管理

### 2. 文档解析
- PDF文档解析
- DOCX文档解析
- Markdown文档支持

### 3. 智能检索
- 关键词搜索
- 语义搜索
- 混合搜索

### 4. MCP协议
- 工具调用
- 文档查询
- 内容获取

## 性能要求

1. 文档上传：在1分钟内完成
2. 文档搜索：在1秒内返回结果
3. 文档解析：快速准确
4. 索引构建：高效可靠

## 测试场景

### 场景1: 大批量文档上传
测试系统同时处理多个文档上传的能力。

### 场景2: 复杂查询搜索
测试系统处理复杂查询的响应速度。

### 场景3: 高并发访问
测试系统在高并发情况下的稳定性。

## 结论

AI技术文档库通过优化的架构设计和高效的算法实现，
确保了系统的快速响应和高性能表现。
EOF

print_result 0 "测试文档准备完成"

print_header "文档更新性能测试（要求：<1分钟）"

# 测试1: 小文档上传性能
print_test "小文档上传性能（1MB）"
start_time=$(date +%s.%N)
response=$(curl -s -X POST "$BACKEND_URL/api/v1/documents/upload" \
    -H "Authorization: Bearer $JWT_TOKEN" \
    -F "file=@$TEST_DOC_DIR/small_test.txt" \
    -F "title=性能测试_小文档" \
    -F "category=性能测试")
end_time=$(date +%s.%N)
elapsed_time=$(echo "$end_time - $start_time" | bc)
UPLOAD_TIMES["small"]=$elapsed_time

echo "上传耗时: ${elapsed_time}秒"

# 等待文档处理完成
sleep 5

if (( $(echo "$elapsed_time < $UPDATE_TIME_LIMIT" | bc -l) )); then
    echo -e "${GREEN}✓ 满足要求（<${UPDATE_TIME_LIMIT}秒）${NC}"
    print_result 0 "小文档上传性能满足要求"
else
    echo -e "${RED}✗ 不满足要求（>=${UPDATE_TIME_LIMIT}秒）${NC}"
    print_result 1 "小文档上传性能不满足要求"
fi

# 测试2: 中文档上传性能
print_test "中文档上传性能（5MB）"
start_time=$(date +%s.%N)
response=$(curl -s -X POST "$BACKEND_URL/api/v1/documents/upload" \
    -H "Authorization: Bearer $JWT_TOKEN" \
    -F "file=@$TEST_DOC_DIR/medium_test.txt" \
    -F "title=性能测试_中文档" \
    -F "category=性能测试")
end_time=$(date +%s.%N)
elapsed_time=$(echo "$end_time - $start_time" | bc)
UPLOAD_TIMES["medium"]=$elapsed_time

echo "上传耗时: ${elapsed_time}秒"

# 等待文档处理完成
sleep 10

if (( $(echo "$elapsed_time < $UPDATE_TIME_LIMIT" | bc -l) )); then
    echo -e "${GREEN}✓ 满足要求（<${UPDATE_TIME_LIMIT}秒）${NC}"
    print_result 0 "中文档上传性能满足要求"
else
    echo -e "${RED}✗ 不满足要求（>=${UPDATE_TIME_LIMIT}秒）${NC}"
    print_result 1 "中文档上传性能不满足要求"
fi

# 测试3: 大文档上传性能
print_test "大文档上传性能（10MB）"
start_time=$(date +%s.%N)
response=$(curl -s -X POST "$BACKEND_URL/api/v1/documents/upload" \
    -H "Authorization: Bearer $JWT_TOKEN" \
    -F "file=@$TEST_DOC_DIR/large_test.txt" \
    -F "title=性能测试_大文档" \
    -F "category=性能测试")
end_time=$(date +%s.%N)
elapsed_time=$(echo "$end_time - $start_time" | bc)
UPLOAD_TIMES["large"]=$elapsed_time

echo "上传耗时: ${elapsed_time}秒"

# 等待文档处理完成
sleep 20

if (( $(echo "$elapsed_time < $UPDATE_TIME_LIMIT" | bc -l) )); then
    echo -e "${GREEN}✓ 满足要求（<${UPDATE_TIME_LIMIT}秒）${NC}"
    print_result 0 "大文档上传性能满足要求"
else
    echo -e "${RED}✗ 不满足要求（>=${UPDATE_TIME_LIMIT}秒）${NC}"
    print_result 1 "大文档上传性能不满足要求"
fi

# 测试4: 实际内容文档上传与处理
print_test "实际内容文档上传与处理（含解析和索引）"
start_time=$(date +%s.%N)
response=$(curl -s -X POST "$BACKEND_URL/api/v1/documents/upload" \
    -H "Authorization: Bearer $JWT_TOKEN" \
    -F "file=@$TEST_DOC_DIR/content_test.md" \
    -F "title=性能测试_内容文档" \
    -F "category=性能测试")
download_end_time=$(date +%s.%N)
download_time=$(echo "$download_end_time - $start_time" | bc)

# 等待文档解析和索引构建完成
echo "等待文档解析和索引构建完成（最多60秒）..."
wait_time=0
while [ $wait_time -lt 60 ]; do
    # 检查文档是否可搜索
    search_response=$(curl -s -X POST "$BACKEND_URL/api/v1/search" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $JWT_TOKEN" \
        -d '{"query":"Vue.js","searchType":"keyword","page":1,"size":5}')
    
    if echo "$search_response" | grep -q "Vue.js"; then
        echo "文档已可搜索"
        break
    fi
    
    sleep 5
    wait_time=$((wait_time + 5))
done

process_end_time=$(date +%s.%N)
total_time=$(echo "$process_end_time - $start_time" | bc)
process_time=$(echo "$process_end_time - $download_end_time" | bc)

echo "上传耗时: ${download_time}秒"
echo "处理耗时: ${process_time}秒"
echo "总耗时: ${total_time}秒"

if (( $(echo "$total_time < $UPDATE_TIME_LIMIT" | bc -l) )); then
    echo -e "${GREEN}✓ 满足要求（<${UPDATE_TIME_LIMIT}秒）${NC}"
    print_result 0 "文档更新（上传+处理）性能满足要求"
else
    echo -e "${RED}✗ 不满足要求（>=${UPDATE_TIME_LIMIT}秒）${NC}"
    print_result 1 "文档更新（上传+处理）性能不满足要求"
fi

print_header "文档检索性能测试（要求：<1秒）"

# 测试准备：确保有足够的文档用于搜索
print_test "准备搜索测试数据"
sleep 10

# 测试5: 关键词搜索性能
print_test "关键词搜索性能"
test_queries=("Vue" "Go" "Python" "PostgreSQL" "Docker")

for query in "${test_queries[@]}"; do
    start_time=$(date +%s.%N)
    response=$(curl -s -X POST "$BACKEND_URL/api/v1/search" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $JWT_TOKEN" \
        -d "{\"query\":\"$query\",\"searchType\":\"keyword\",\"page\":1,\"size\":10}")
    end_time=$(date +%s.%N)
    elapsed_time=$(echo "$end_time - $start_time" | bc)
    SEARCH_TIMES["keyword_$query"]=$elapsed_time
    
    echo "查询 '$query': ${elapsed_time}秒"
    
    if (( $(echo "$elapsed_time < $SEARCH_TIME_LIMIT" | bc -l) )); then
        echo -e "${GREEN}✓ 满足要求（<${SEARCH_TIME_LIMIT}秒）${NC}"
    else
        echo -e "${RED}✗ 不满足要求（>=${SEARCH_TIME_LIMIT}秒）${NC}"
    fi
done

# 计算平均搜索时间
sum=0
count=0
for key in "${!SEARCH_TIMES[@]}"; do
    if [[ $key == keyword_* ]]; then
        sum=$(echo "$sum + ${SEARCH_TIMES[$key]}" | bc)
        ((count++)) || true
    fi
done

if [ $count -gt 0 ]; then
    avg_time=$(echo "scale=3; $sum / $count" | bc)
    echo "平均关键词搜索时间: ${avg_time}秒"
    if (( $(echo "$avg_time < $SEARCH_TIME_LIMIT" | bc -l) )); then
        print_result 0 "关键词搜索平均性能满足要求"
    else
        print_result 1 "关键词搜索平均性能不满足要求"
    fi
fi

# 测试6: 语义搜索性能
print_test "语义搜索性能"
semantic_queries=("组件开发" "数据库优化" "容器部署")

for query in "${semantic_queries[@]}"; do
    start_time=$(date +%s.%N)
    response=$(curl -s -X POST "$BACKEND_URL/api/v1/search" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $JWT_TOKEN" \
        -d "{\"query\":\"$query\",\"searchType\":\"semantic\",\"page\":1,\"size\":10}")
    end_time=$(date +%s.%N)
    elapsed_time=$(echo "$end_time - $start_time" | bc)
    SEARCH_TIMES["semantic_$query"]=$elapsed_time
    
    echo "查询 '$query': ${elapsed_time}秒"
    
    if (( $(echo "$elapsed_time < $SEARCH_TIME_LIMIT" | bc -l) )); then
        echo -e "${GREEN}✓ 满足要求（<${SEARCH_TIME_LIMIT}秒）${NC}"
    else
        echo -e "${RED}✗ 不满足要求（>=${SEARCH_TIME_LIMIT}秒）${NC}"
    fi
done

# 测试7: 混合搜索性能
print_test "混合搜索性能"
hybrid_queries=("Vue组件" "Go服务" "Python解析")

for query in "${hybrid_queries[@]}"; do
    start_time=$(date +%s.%N)
    response=$(curl -s -X POST "$BACKEND_URL/api/v1/search" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $JWT_TOKEN" \
        -d "{\"query\":\"$query\",\"searchType\":\"hybrid\",\"page\":1,\"size\":10}")
    end_time=$(date +%s.%N)
    elapsed_time=$(echo "$end_time - $start_time" | bc)
    SEARCH_TIMES["hybrid_$query"]=$elapsed_time
    
    echo "查询 '$query': ${elapsed_time}秒"
    
    if (( $(echo "$elapsed_time < $SEARCH_TIME_LIMIT" | bc -l) )); then
        echo -e "${GREEN}✓ 满足要求（<${SEARCH_TIME_LIMIT}秒）${NC}"
    else
        echo -e "${RED}✗ 不满足要求（>=${SEARCH_TIME_LIMIT}秒）${NC}"
    fi
done

# 测试8: 并发搜索性能
print_test "并发搜索性能测试（10个并发请求）"
echo "开始10个并发搜索请求..."

concurrent_results=()
for i in {1..10}; do
    (
        start_time=$(date +%s.%N)
        curl -s -X POST "$BACKEND_URL/api/v1/search" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $JWT_TOKEN" \
            -d '{"query":"性能测试","searchType":"keyword","page":1,"size":5}' > /dev/null
        end_time=$(date +%s.%N)
        elapsed_time=$(echo "$end_time - $start_time" | bc)
        echo "$elapsed_time"
    ) &
done

# 等待所有并发请求完成
wait

echo "并发搜索测试完成"

# 测试9: 批量文档查询性能
print_test "批量文档查询性能"
start_time=$(date +%s.%N)
response=$(curl -s -X GET "$BACKEND_URL/api/v1/documents?page=1&size=20" \
    -H "Authorization: Bearer $JWT_TOKEN")
end_time=$(date +%s.%N)
elapsed_time=$(echo "$end_time - $start_time" | bc)
SEARCH_TIMES["batch_query"]=$elapsed_time

echo "批量查询耗时: ${elapsed_time}秒"

if (( $(echo "$elapsed_time < $SEARCH_TIME_LIMIT" | bc -l) )); then
    print_result 0 "批量查询性能满足要求"
else
    print_result 1 "批量查询性能不满足要求"
fi

print_header "总体性能评估"

# 计算所有搜索的平均时间
sum_all=0
count_all=0
for time in "${SEARCH_TIMES[@]}"; do
    sum_all=$(echo "$sum_all + $time" | bc)
    ((count_all++)) || true
done

if [ $count_all -gt 0 ]; then
    avg_search_time=$(echo "scale=3; $sum_all / $count_all" | bc)
    echo ""
    echo "搜索性能汇总："
    echo "  总搜索次数: $count_all"
    echo "  平均搜索时间: ${avg_search_time}秒"
    echo "  性能要求: <${SEARCH_TIME_LIMIT}秒"
    
    if (( $(echo "$avg_search_time < $SEARCH_TIME_LIMIT" | bc -l) )); then
        echo -e "${GREEN}✓ 搜索性能满足测评要求${NC}"
        print_result 0 "搜索性能整体满足要求"
    else
        echo -e "${RED}✗ 搜索性能不满足测评要求${NC}"
        print_result 1 "搜索性能整体不满足要求"
    fi
fi

# 计算所有上传的平均时间
sum_upload=0
count_upload=0
for time in "${UPLOAD_TIMES[@]}"; do
    sum_upload=$(echo "$sum_upload + $time" | bc)
    ((count_upload++)) || true
done

if [ $count_upload -gt 0 ]; then
    avg_upload_time=$(echo "scale=3; $sum_upload / $count_upload" | bc)
    echo ""
    echo "文档更新性能汇总："
    echo "  总上传次数: $count_upload"
    echo "  平均上传时间: ${avg_upload_time}秒"
    echo "  性能要求: <${UPDATE_TIME_LIMIT}秒（1分钟）"
    
    if (( $(echo "$avg_upload_time < $UPDATE_TIME_LIMIT" | bc -l) )); then
        echo -e "${GREEN}✓ 文档更新性能满足测评要求${NC}"
        print_result 0 "文档更新性能整体满足要求"
    else
        echo -e "${RED}✗ 文档更新性能不满足测评要求${NC}"
        print_result 1 "文档更新性能整体不满足要求"
    fi
fi

# 性能稳定性测试
print_test "性能稳定性测试（连续10次搜索）"
echo "执行持续搜索性能测试..."

stable_times=()
for i in {1..10}; do
    start_time=$(date +%s.%N)
    curl -s -X POST "$BACKEND_URL/api/v1/search" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $JWT_TOKEN" \
        -d '{"query":"性能测试","searchType":"keyword","page":1,"size":5}' > /dev/null
    end_time=$(date +%s.%N)
    elapsed_time=$(echo "$end_time - $start_time" | bc)
    stable_times+=($elapsed_time)
    echo "  第${i}次: ${elapsed_time}秒"
done

# 计算标准差
if [ ${#stable_times[@]} -gt 1 ]; then
    sum_stable=0
    for time in "${stable_times[@]}"; do
        sum_stable=$(echo "$sum_stable + $time" | bc)
    done
    mean_time=$(echo "scale=3; $sum_stable / ${#stable_times[@]}" | bc)
    
    variance_sum=0
    for time in "${stable_times[@]}"; do
        diff=$(echo "$time - $mean_time" | bc)
        variance_sum=$(echo "$variance_sum + ($diff * $diff)" | bc)
    done
    variance=$(echo "scale=6; $variance_sum / ${#stable_times[@]}" | bc)
    std_dev=$(echo "scale=3; sqrt($variance)" | bc)
    
    echo ""
    echo "性能稳定性分析："
    echo "  平均时间: ${mean_time}秒"
    echo "  标准差: ${std_dev}秒"
    echo "  变异系数: $(echo "scale=2; $std_dev / $mean_time * 100" | bc)%"
    
    if (( $(echo "$std_dev < 0.2" | bc -l) )); then
        echo -e "${GREEN}✓ 性能稳定性良好${NC}"
        print_result 0 "性能稳定性测试通过"
    else
        echo -e "${YELLOW}⚠ 性能波动较大，建议优化${NC}"
        print_result 1 "性能稳定性需要改进"
    fi
fi

print_header "性能测试总结"

echo ""
echo "=========================================="
echo "性能测试报告"
echo "=========================================="
echo ""
echo "1. 文档更新性能（要求: <1分钟）"
echo "   - 小文档（1MB）: ${UPLOAD_TIMES[small]}秒"
echo "   - 中文档（5MB）: ${UPLOAD_TIMES[medium]}秒"
echo "   - 大文档（10MB）: ${UPLOAD_TIMES[large]}秒"
echo "   - 平均上传时间: ${avg_upload_time}秒"
echo ""
echo "2. 文档检索性能（要求: <1秒）"
echo "   - 平均搜索时间: ${avg_search_time}秒"
echo "   - 总测试次数: $count_all"
echo "   - 性能稳定性: 标准差=${std_dev}秒"
echo ""
echo "3. 测评要求满足情况："
echo "   - 文档更新（<1分钟）:"
if (( $(echo "$avg_upload_time < 60" | bc -l) )); then
    echo -e "     ${GREEN}✓ 满足要求${NC}"
else
    echo -e "     ${RED}✗ 不满足要求${NC}"
fi

echo "   - 文档检索（<1秒）:"
if (( $(echo "$avg_search_time < 1" | bc -l) )); then
    echo -e "     ${GREEN}✓ 满足要求${NC}"
else
    echo -e "     ${RED}✗ 不满足要求${NC}"
fi

echo ""
echo "=========================================="
echo "总体测试结果"
echo "=========================================="
echo -e "${GREEN}通过测试${NC}: $PASSED_TESTS"
echo -e "${RED}失败测试${NC}: $FAILED_TESTS"
echo "总计测试: $TOTAL_TESTS"

if [ $FAILED_TESTS -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✓ 所有性能测试通过！${NC}"
    echo ""
    echo "系统性能满足测评方案要求："
    echo "  ✓ 文档更新性能满足<1分钟要求"
    echo "  ✓ 文档检索性能满足<1秒要求"
    echo "  ✓ 性能稳定性良好"
    exit 0
else
    echo ""
    echo -e "${RED}✗ 部分性能测试未通过${NC}"
    echo ""
    echo "建议："
    echo "  1. 检查文档解析服务的性能"
    echo "  2. 优化数据库查询和索引"
    echo "  3. 检查网络和存储I/O性能"
    echo "  4. 考虑增加缓存机制"
    exit 1
fi