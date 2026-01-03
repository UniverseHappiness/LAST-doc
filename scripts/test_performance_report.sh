#!/bin/bash

echo "=========================================="
echo "测试性能报告API"
echo "=========================================="

# 获取管理员token
echo -e "\n1. 获取管理员token..."
TOKEN=$(curl -s -X POST "http://localhost:8080/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "❌ 获取token失败"
  exit 1
fi

echo "✅ Token获取成功"

# 测试性能报告（24小时）
echo -e "\n2. 测试性能报告API（24小时）..."
START_TIME=$(date -u -d '24 hours ago' +%Y-%m-%dT%H:%M:%SZ)
END_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)

echo "开始时间: $START_TIME"
echo "结束时间: $END_TIME"

REPORT_RESPONSE=$(curl -s -X GET "http://localhost:8080/api/v1/monitor/performance" \
  -H "Authorization: Bearer ${TOKEN}" \
  -G -d "start_time=${START_TIME}" \
  -d "end_time=${END_TIME}")

echo -e "\n性能报告响应:"
echo "$REPORT_RESPONSE"

# 检查响应结构
echo -e "\n3. 检查响应结构..."
if echo "$REPORT_RESPONSE" | grep -q "cpu_usage"; then
  DATA_COUNT=$(echo "$REPORT_RESPONSE" | grep -o '"cpu_usage":\[[^]]*' | grep -o ',[^]]*' | head -1 | cut -d',' -f2 | wc -w)
  echo "✅ 响应包含cpu_usage，数据点数量: $(echo "$DATA_COUNT" | tr -d ' ')"

  # 提取CPU数据
  CPU_DATA=$(echo "$REPORT_RESPONSE" | grep -o '"cpu_usage":\[[^]]*' | head -1)
  echo "CPU数据示例: $CPU_DATA"
else
  echo "❌ 响应不包含cpu_usage字段"
fi

# 测试指标历史API
echo -e "\n4. 测试指标历史API..."
HISTORY_RESPONSE=$(curl -s -X GET "http://localhost:8080/api/v1/monitor/metrics/history" \
  -H "Authorization: Bearer ${TOKEN}" \
  -G -d "start_time=${START_TIME}" \
  -d "end_time=${END_TIME}")

echo -e "\n指标历史响应（前500字符）:"
echo "$HISTORY_RESPONSE" | head -c 500

# 检查指标历史
if echo "$HISTORY_RESPONSE" | grep -q "metrics"; then
  METRICS_COUNT=$(echo "$HISTORY_RESPONSE" | grep -o '"count":[^,}]*' | cut -d':' -f2)
  echo "✅ 指标历史查询成功，总记录数: $METRICS_COUNT"
else
  echo "❌ 指标历史查询失败"
fi

# 测试指标报告API
echo -e "\n5. 测试指标报告API..."
METRICS_REPORT=$(curl -s -X GET "http://localhost:8080/api/v1/monitor/metrics/report?duration=24h" \
  -H "Authorization: Bearer ${TOKEN}")

echo -e "\n指标报告响应:"
echo "$METRICS_REPORT"

echo -e "\n=========================================="
echo "测试完成"
echo "=========================================="