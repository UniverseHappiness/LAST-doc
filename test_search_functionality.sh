#!/bin/bash

# 测试文档检索功能脚本

echo "开始测试文档检索功能..."

# 检查是否已启动服务
if ! pgrep -f "last-doc" > /dev/null; then
    echo "启动服务..."
    SERVER_PORT=8081 go run cmd/main.go &
    sleep 5
fi

# 测试搜索API
echo "1. 测试搜索API..."

# 测试关键词搜索
echo "   测试关键词搜索..."
response=$(curl -s -X POST "http://localhost:8081/api/v1/search" \
    -H "Content-Type: application/json" \
    -d '{"query": "测试", "searchType": "keyword", "page": 1, "size": 10}')

echo "关键词搜索响应: $response"

# 测试语义搜索
echo "   测试语义搜索..."
response=$(curl -s -X POST "http://localhost:8081/api/v1/search" \
    -H "Content-Type: application/json" \
    -d '{"query": "测试", "searchType": "semantic", "page": 1, "size": 10}')

echo "语义搜索响应: $response"

# 测试混合搜索
echo "   测试混合搜索..."
response=$(curl -s -X POST "http://localhost:8081/api/v1/search" \
    -H "Content-Type: application/json" \
    -d '{"query": "测试", "searchType": "hybrid", "page": 1, "size": 10}')

echo "混合搜索响应: $response"

# 测试GET方式搜索
echo "2. 测试GET方式搜索..."
response=$(curl -s "http://localhost:8081/api/v1/search?query=测试&search_type=keyword&page=1&size=10")
echo "GET搜索响应: $response"

# 测试索引构建
echo "3. 测试索引构建..."
response=$(curl -s -X POST "http://localhost:8081/api/v1/search/documents/test-doc-id/versions/1.0.0/index")
echo "索引构建响应: $response"

# 测试索引状态
echo "4. 测试索引状态..."
response=$(curl -s "http://localhost:8081/api/v1/search/documents/test-doc-id/index/status")
echo "索引状态响应: $response"

# 测试性能
echo "5. 测试搜索性能..."
start_time=$(date +%s.%N)

# 执行10次搜索并计算平均时间
for i in {1..10}; do
    curl -s -X POST "http://localhost:8081/api/v1/search" \
        -H "Content-Type: application/json" \
        -d '{"query": "性能测试", "searchType": "keyword", "page": 1, "size": 10}' > /dev/null
done

end_time=$(date +%s.%N)
elapsed_time=$(echo "$end_time - $start_time" | bc)
avg_time=$(echo "scale=3; $elapsed_time / 10" | bc)

echo "平均搜索时间: ${avg_time}秒"

# 检查是否满足1秒内的要求
if (( $(echo "$avg_time < 1.0" | bc -l) )); then
    echo "✓ 搜索性能满足要求（小于1秒）"
else
    echo "✗ 搜索性能不满足要求（大于等于1秒）"
fi

echo "文档检索功能测试完成"