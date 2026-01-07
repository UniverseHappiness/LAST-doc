# 测试数据可视化说明

## 目录结构

```
test_data/
├── README.md                 # 本说明文件
├── generate_charts.py        # Python图表生成脚本
├── charts/                   # 生成的图表目录
│   ├── effectiveness_performance_comparison.png
│   ├── effectiveness_function_coverage.png
│   ├── effectiveness_radar_score.png
│   ├── performance_upload_times.png
│   ├── performance_search_times.png
│   ├── performance_concurrency.png
│   ├── performance_upload_breakdown.png
│   ├── overall_score.png
│   ├── performance_requirements.png
│   └── user_experience_improvement.png
└── ascii_charts/             # ASCII字符图表
    ├── effectiveness_comparison.txt
    ├── performance_summary.txt
    └── user_improvement.txt
```

## 使用方法

### 1. 生成PNG图表

```bash
# 安装依赖
pip install matplotlib

# 运行脚本
cd test_data
python generate_charts.py
```

生成的图表将保存在 `charts/` 目录中。

### 2. 在报告中使用

#### 效果测试报告
- `effectiveness_performance_comparison.png` - 性能对比图表
- `effectiveness_function_coverage.png` - 功能覆盖度图表
- `effectiveness_radar_score.png` - 综合评分雷达图

#### 性能测试报告
- `performance_upload_times.png` - 文档上传性能图表
- `performance_search_times.png` - 文档检索性能图表
- `performance_concurrency.png` - 并发性能图表
- `performance_upload_breakdown.png` - 上传流程分解图表
- `performance_requirements.png` - 性能要求满足情况

#### 综合评估
- `overall_score.png` - 综合性能评分饼图
- `user_experience_improvement.png` - 用户体验提升指标

## 图表说明

### 效果测试图表

#### 1. 性能对比柱状图
展示传统API和MCP工具在10个不同查询中的响应时间对比。

#### 2. 功能覆盖度对比
对比传统方式和MCP方式在5个功能维度的覆盖度。

#### 3. 综合评分雷达图
从性能、效率、功能、体验、准确5个维度对比两种方式。

### 性能测试图表

#### 1. 文档上传性能
展示1MB、5MB、10MB文档的上传耗时，包含60秒阈值线。

#### 2. 文档检索性能
展示关键词搜索、语义搜索、混合搜索等5种搜索类型的性能。

#### 3. 并发性能
展示10、20、50并发数下的平均响应时间趋势。

#### 4. 上传流程分解
堆叠柱状图展示不同大小文档在上传、解析、向量生成、索引构建各环节的耗时。

## 数据来源

所有图表数据均来自以下测试脚本：
- [`test_effectiveness.sh`](../test_effectiveness.sh) - 效果测试脚本
- [`test_performance_comprehensive.sh`](../test_performance_comprehensive.sh) - 性能测试脚本

## 注意事项

1. **字体支持**：Python脚本使用DejaVu Sans字体，确保中文字体显示正常
2. **图表质量**：生成的图表DPI为300，适合打印和高分辨率显示
3. **颜色方案**：采用专业的配色方案，符合数据可视化最佳实践
4. **文件格式**：PNG格式，支持透明背景，便于嵌入文档

## 技术支持

如有问题，请参考：
- [EFFECTIVENESS_TEST_REPORT.md](../EFFECTIVENESS_TEST_REPORT.md) - 效果测试报告
- [PERFORMANCE_TEST_REPORT.md](../PERFORMANCE_TEST_REPORT.md) - 性能测试报告