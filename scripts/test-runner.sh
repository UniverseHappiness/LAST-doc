#!/bin/bash

# 文档管理系统测试运行脚本
# 该脚本用于验证测试机制是否正常工作

set -e  # 遇到错误时退出

echo "=========================================="
echo "文档管理系统 - 测试机制验证"
echo "=========================================="

# 1. 运行单元测试
echo ""
echo "1. 运行单元测试..."
echo "------------------------------------------"
if go test -v -cover ./internal/...; then
    echo "✅ 单元测试执行成功"
else
    echo "❌ 单元测试执行失败"
    exit 1
fi

# 2. 运行所有测试并生成覆盖率报告
echo ""
echo "2. 生成测试覆盖率报告..."
echo "------------------------------------------"
if go test -coverprofile=coverage.out ./...; then
    echo "✅ 覆盖率数据生成成功"
    
    # 生成HTML覆盖率报告
    if go tool cover -html=coverage.out -o coverage.html; then
        echo "✅ HTML覆盖率报告生成成功"
    else
        echo "❌ HTML覆盖率报告生成失败"
    fi
else
    echo "❌ 覆盖率数据生成失败"
    exit 1
fi

# 3. 运行代码格式检查
echo ""
echo "3. 检查代码格式..."
echo "------------------------------------------"
if go fmt ./...; then
    echo "✅ 代码格式检查通过"
else
    echo "❌ 代码格式检查失败"
    exit 1
fi

# 4. 运行代码规范检查（如果golangci-lint已安装）
echo ""
echo "4. 检查代码规范..."
echo "------------------------------------------"
if command -v golangci-lint &> /dev/null; then
    if golangci-lint run; then
        echo "✅ 代码规范检查通过"
    else
        echo "❌ 代码规范检查失败"
        exit 1
    fi
else
    echo "⚠️ golangci-lint 未安装，跳过代码规范检查"
fi

# 5. 运行构建测试
echo ""
echo "5. 运行构建测试..."
echo "------------------------------------------"
if go build -o bin/ai-doc-library-test cmd/main.go; then
    echo "✅ 构建测试通过"
    rm -f bin/ai-doc-library-test  # 清理测试构建文件
else
    echo "❌ 构建测试失败"
    exit 1
fi

echo ""
echo "=========================================="
echo "所有测试验证完成！"
echo "=========================================="

# 显示覆盖率报告位置
if [ -f "coverage.html" ]; then
    echo "覆盖率报告已生成: coverage.html"
fi

echo ""
echo "如需查看详细覆盖率信息，请运行: go tool cover -func=coverage.out"