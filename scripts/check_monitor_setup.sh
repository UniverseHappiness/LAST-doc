#!/bin/bash

echo "=========================================="
echo "系统监控功能检查"
echo "=========================================="

# 检查后端文件
echo -e "\n1. 检查后端文件..."
FILES=(
  "internal/model/monitor.go"
  "internal/repository/monitor_repository.go"
  "internal/service/monitor_service.go"
  "internal/handler/monitor_handler.go"
)

for file in "${FILES[@]}"; do
  if [ -f "$file" ]; then
    echo "✅ $file 存在"
  else
    echo "❌ $file 不存在"
  fi
done

# 检查前端文件
echo -e "\n2. 检查前端文件..."
FRONTEND_FILES=(
  "web/src/utils/monitorService.js"
  "web/src/views/MonitorView.vue"
)

for file in "${FRONTEND_FILES[@]}"; do
  if [ -f "$file" ]; then
    echo "✅ $file 存在"
  else
    echo "❌ $file 不存在"
  fi
done

# 检查Chart.js依赖
echo -e "\n3. 检查Chart.js依赖..."
cd web
if npm list chart.js --depth=0 2>/dev/null | grep -q chart.js; then
  echo "✅ Chart.js 已安装"
else
  echo "❌ Chart.js 未安装，请运行: cd web && npm install"
fi
cd ..

# 检查App.vue中的导入和注册
echo -e "\n4. 检查App.vue配置..."
if grep -q "MonitorView" web/src/App.vue; then
  echo "✅ MonitorView 已导入"
else
  echo "❌ MonitorView 未导入"
fi

if grep -q "currentView === 'monitor'" web/src/App.vue; then
  echo "✅ 监控视图路由已配置"
else
  echo "❌ 监控视图路由未配置"
fi

# 检查路由配置
echo -e "\n5. 检查路由配置..."
if grep -q "/monitor" internal/router/router.go; then
  echo "✅ 监控路由已配置"
else
  echo "❌ 监控路由未配置"
fi

# 检查数据库表
echo -e "\n6. 检查数据库表..."
DB_CHECK=$(psql -h localhost -U postgres -d ai_doc_library -t -c "
  SELECT COUNT(*) FROM information_schema.tables 
  WHERE table_schema = 'public' 
  AND table_name IN ('system_metrics', 'log_entries')
" | tr -d ' ')

if [ "$DB_CHECK" = "2" ]; then
  echo "✅ 数据库表已创建（system_metrics, log_entries）"
else
  echo "⚠️  数据库表可能未完全创建，发现: $DB_CHECK 个表"
fi

echo -e "\n=========================================="
echo "检查完成"
echo "=========================================="
echo -e "\n访问系统监控页面："
echo "1. 使用管理员账户登录（admin / admin123）"
echo "2. 在侧边栏点击'系统监控'链接"
echo -e "\n如果前端服务未运行，请执行："
echo "cd web && npm run dev"