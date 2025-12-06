#!/bin/bash

# 自动生成测试脚本
# 该脚本使用gotests工具为项目自动生成测试框架

set -e  # 遇到错误时退出

echo "=========================================="
echo "AI技术文档库 - 自动生成测试"
echo "=========================================="

# 检查是否安装了gotests
if ! command -v gotests &> /dev/null; then
    echo "正在安装gotests工具..."
    go get github.com/cweill/gotests/...
    echo "✅ gotests工具安装成功"
fi

# 创建生成测试的目录
mkdir -p internal/generated_tests

# 生成服务层测试
echo ""
echo "生成服务层测试框架..."
echo "------------------------------------------"
if gotests -all -w -template_dir=test_templates internal/service/; then
    echo "✅ 服务层测试框架生成成功"
else
    echo "❌ 服务层测试框架生成失败"
fi

# 生成仓库层测试
echo ""
echo "生成仓库层测试框架..."
echo "------------------------------------------"
if gotests -all -w -template_dir=test_templates internal/repository/; then
    echo "✅ 仓库层测试框架生成成功"
else
    echo "❌ 仓库层测试框架生成失败"
fi

# 生成处理器层测试
echo ""
echo "生成处理器层测试框架..."
echo "------------------------------------------"
if gotests -all -w -template_dir=test_templates internal/handler/; then
    echo "✅ 处理器层测试框架生成成功"
else
    echo "❌ 处理器层测试框架生成失败"
fi

# 生成路由层测试
echo ""
echo "生成路由层测试框架..."
echo "------------------------------------------"
if gotests -all -w -template_dir=test_templates internal/router/; then
    echo "✅ 路由层测试框架生成成功"
else
    echo "❌ 路由层测试框架生成失败"
fi

echo ""
echo "=========================================="
echo "自动生成测试完成！"
echo "=========================================="
echo ""
echo "生成的测试文件位置："
echo "- 服务层测试：internal/service/*_test.go"
echo "- 仓库层测试：internal/repository/*_test.go"
echo "- 处理器层测试：internal/handler/*_test.go"
echo "- 路由层测试：internal/router/*_test.go"
echo ""
echo "注意：生成的测试框架需要手动添加测试逻辑和断言。"
echo "建议使用混合方式：基础框架使用自动生成，复杂业务逻辑手写。"
echo "模拟对象（mocks）需要手动创建，以确保测试的灵活性和准确性。"

echo ""
echo "=========================================="
echo "自动生成测试完成！"
echo "=========================================="
echo ""
echo "生成的测试文件位置："
echo "- 服务层测试：internal/service/*_test.go"
echo "- 仓库层测试：internal/repository/*_test.go"
echo "- 处理器层测试：internal/handler/*_test.go"
echo "- 路由层测试：internal/router/*_test.go"
echo "- 模拟对象：internal/mocks/"
echo ""
echo "注意：生成的测试框架需要手动添加测试逻辑和断言。"
echo "建议使用混合方式：基础框架使用自动生成，复杂业务逻辑手写。"