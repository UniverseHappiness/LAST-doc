#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
测试数据可视化图表生成脚本
用于生成专业的测试结果图表
"""

import matplotlib.pyplot as plt
import numpy as np
from matplotlib import rcParams

# 设置中文字体支持
rcParams['font.sans-serif'] = ['DejaVu Sans', 'SimHei', 'Arial Unicode MS']
rcParams['axes.unicode_minus'] = False

# 设置图表样式
plt.style.use('seaborn-v0_8-darkgrid')

class ChartGenerator:
    """图表生成器"""
    
    def __init__(self, output_dir='./charts'):
        """初始化图表生成器"""
        self.output_dir = output_dir
        import os
        os.makedirs(output_dir, exist_ok=True)
    
    def bar_chart(self, data, labels, title, xlabel, ylabel, filename, threshold=None):
        """生成柱状图"""
        fig, ax = plt.subplots(figsize=(12, 6))
        
        bars = ax.bar(labels, data, color='#4CAF50', alpha=0.8, edgecolor='darkgreen', linewidth=2)
        
        # 在柱子上显示数值
        for bar in bars:
            height = bar.get_height()
            ax.text(bar.get_x() + bar.get_width()/2., height,
                   f'{height:.3f}',
                   ha='center', va='bottom', fontsize=10, fontweight='bold')
        
        # 添加阈值线
        if threshold:
            ax.axhline(y=threshold, color='red', linestyle='--', linewidth=2, 
                      label=f'Threshold: {threshold}', alpha=0.7)
            ax.legend()
        
        ax.set_title(title, fontsize=16, fontweight='bold', pad=20)
        ax.set_xlabel(xlabel, fontsize=12)
        ax.set_ylabel(ylabel, fontsize=12)
        ax.tick_params(axis='both', which='major', labelsize=10)
        
        plt.tight_layout()
        plt.savefig(f'{self.output_dir}/{filename}', dpi=300, bbox_inches='tight')
        plt.close()
        print(f'✓ 生成图表: {filename}')
    
    def line_chart(self, data_dict, title, xlabel, ylabel, filename, threshold=None):
        """生成折线图"""
        fig, ax = plt.subplots(figsize=(12, 6))
        
        colors = ['#4CAF50', '#2196F3', '#FF9800', '#9C27B0', '#F44336']
        
        for idx, (label, data) in enumerate(data_dict.items()):
            ax.plot(range(1, len(data)+1), data, 
                   marker='o', linewidth=2, markersize=8, 
                   color=colors[idx % len(colors)], label=label)
        
        # 添加阈值线
        if threshold:
            ax.axhline(y=threshold, color='red', linestyle='--', linewidth=2, 
                      label=f'Threshold: {threshold}', alpha=0.7)
        
        ax.set_title(title, fontsize=16, fontweight='bold', pad=20)
        ax.set_xlabel(xlabel, fontsize=12)
        ax.set_ylabel(ylabel, fontsize=12)
        ax.legend(fontsize=10, loc='best')
        ax.tick_params(axis='both', which='major', labelsize=10)
        
        plt.tight_layout()
        plt.savefig(f'{self.output_dir}/{filename}', dpi=300, bbox_inches='tight')
        plt.close()
        print(f'✓ 生成图表: {filename}')
    
    def pie_chart(self, data, labels, title, filename, explode=None):
        """生成饼图"""
        fig, ax = plt.subplots(figsize=(10, 8))
        
        colors = ['#4CAF50', '#2196F3', '#FF9800', '#9C27B0', '#F44336', '#00BCD4']
        
        if explode is None:
            explode = [0.05] * len(data)
        
        wedges, texts, autotexts = ax.pie(data, labels=labels, explode=explode,
                                          colors=colors, autopct='%1.1f%%',
                                          shadow=True, startangle=90,
                                          textprops={'fontsize': 11, 'fontweight': 'bold'})
        
        # 美化文本
        for autotext in autotexts:
            autotext.set_color('white')
            autotext.set_fontsize(12)
            autotext.set_fontweight('bold')
        
        ax.set_title(title, fontsize=16, fontweight='bold', pad=20)
        ax.axis('equal')
        
        plt.tight_layout()
        plt.savefig(f'{self.output_dir}/{filename}', dpi=300, bbox_inches='tight')
        plt.close()
        print(f'✓ 生成图表: {filename}')
    
    def radar_chart(self, data_dict, title, filename):
        """生成雷达图"""
        fig, ax = plt.subplots(figsize=(10, 8), subplot_kw=dict(projection='polar'))
        
        categories = list(data_dict.keys())
        N = len(categories)
        
        # 每个类别的角度
        angles = [n / float(N) * 2 * np.pi for n in range(N)]
        angles += angles[:1]
        
        colors = ['#4CAF50', '#2196F3', '#FF9800', '#9C27B0', '#F44336']
        
        for idx, (label, values) in enumerate(data_dict.items()):
            if isinstance(values, dict):
                values_list = list(values.values())
                labels = list(values.keys())
                
                values_list += values_list[:1]
                
                ax.plot(angles, values_list, 'o-', linewidth=2, 
                       color=colors[idx % len(colors)], label=label)
                ax.fill(angles, values_list, alpha=0.25, 
                       color=colors[idx % len(colors)])
        
        ax.set_xticks(angles[:-1])
        ax.set_xticklabels(categories, fontsize=11)
        ax.set_ylim(0, 100)
        ax.set_title(title, fontsize=16, fontweight='bold', pad=20)
        ax.legend(loc='upper right', bbox_to_anchor=(1.3, 1.1))
        ax.grid(True)
        
        plt.tight_layout()
        plt.savefig(f'{self.output_dir}/{filename}', dpi=300, bbox_inches='tight')
        plt.close()
        print(f'✓ 生成图表: {filename}')
    
    def comparison_bar_chart(self, data_dict, title, xlabel, ylabel, filename):
        """生成对比柱状图"""
        fig, ax = plt.subplots(figsize=(14, 7))
        
        categories = list(data_dict.keys())
        groups = list(data_dict[categories[0]].keys())
        
        x = np.arange(len(categories))
        width = 0.25
        
        colors = ['#4CAF50', '#2196F3', '#FF9800']
        
        for idx, group in enumerate(groups):
            values = [data_dict[cat][group] for cat in categories]
            offset = (idx - len(groups)/2 + 0.5) * width
            bars = ax.bar(x + offset, values, width, 
                         label=group, color=colors[idx % len(colors)],
                         alpha=0.8, edgecolor='black', linewidth=1)
            
            # 在柱子上显示数值
            for bar in bars:
                height = bar.get_height()
                ax.text(bar.get_x() + bar.get_width()/2., height,
                       f'{height:.3f}',
                       ha='center', va='bottom', fontsize=9)
        
        ax.set_title(title, fontsize=16, fontweight='bold', pad=20)
        ax.set_xlabel(xlabel, fontsize=12)
        ax.set_ylabel(ylabel, fontsize=12)
        ax.set_xticks(x)
        ax.set_xticklabels(categories, rotation=45, ha='right')
        ax.legend(fontsize=10)
        ax.tick_params(axis='both', which='major', labelsize=10)
        
        plt.tight_layout()
        plt.savefig(f'{self.output_dir}/{filename}', dpi=300, bbox_inches='tight')
        plt.close()
        print(f'✓ 生成图表: {filename}')
    
    def stacked_bar_chart(self, data_dict, title, xlabel, ylabel, filename):
        """生成堆叠柱状图"""
        fig, ax = plt.subplots(figsize=(12, 6))
        
        categories = list(data_dict.keys())
        groups = list(data_dict[categories[0]].keys())
        
        colors = ['#4CAF50', '#2196F3', '#FF9800', '#9C27B0', '#F44336']
        
        bottom = np.zeros(len(categories))
        
        for idx, group in enumerate(groups):
            values = [data_dict[cat][group] for cat in categories]
            ax.bar(categories, values, bottom=bottom, 
                  label=group, color=colors[idx % len(colors)], 
                  alpha=0.8, edgecolor='black', linewidth=1)
            bottom += values
        
        ax.set_title(title, fontsize=16, fontweight='bold', pad=20)
        ax.set_xlabel(xlabel, fontsize=12)
        ax.set_ylabel(ylabel, fontsize=12)
        ax.legend(fontsize=10)
        ax.tick_params(axis='both', which='major', labelsize=10)
        
        plt.tight_layout()
        plt.savefig(f'{self.output_dir}/{filename}', dpi=300, bbox_inches='tight')
        plt.close()
        print(f'✓ 生成图表: {filename}')


def generate_all_charts():
    """生成所有测试图表"""
    
    generator = ChartGenerator()
    
    print("=" * 60)
    print("开始生成测试数据可视化图表")
    print("=" * 60)
    print()
    
    # 1. 效果测试 - 性能对比柱状图
    print("1. 生成效果测试 - 性能对比 图表...")
    queries = ['Vue', 'PostgreSQL', 'Go gRPC', 'Docker', 'Nginx', 
              'Python', 'RESTful', '微服务', '监控', '高可用']
    traditional_times = [0.245, 0.312, 0.287, 0.263, 0.291, 0.278, 
                       0.254, 0.302, 0.269, 0.285]
    mcp_times = [0.238, 0.298, 0.275, 0.251, 0.283, 0.269, 
                 0.248, 0.294, 0.261, 0.277]
    
    generator.comparison_bar_chart(
        {'传统API': dict(zip(queries, traditional_times)), 
         'MCP工具': dict(zip(queries, mcp_times))},
        '传统API vs MCP工具性能对比',
        '查询类型', '响应时间 (秒)',
        'effectiveness_performance_comparison.png'
    )
    
    # 2. 效果测试 - 功能覆盖度对比
    print("2. 生成效果测试 - 功能覆盖度 图表...")
    categories = ['文档搜索', '文档获取', '库管理', '上下文理解', '结果整合']
    traditional_coverage = [60, 65, 40, 20, 30]
    mcp_coverage = [100, 100, 80, 80, 90]
    
    generator.comparison_bar_chart(
        {'传统方式': dict(zip(categories, traditional_coverage)), 
         'MCP方式': dict(zip(categories, mcp_coverage))},
        '功能覆盖度对比 (传统方式 vs MCP方式)',
        '功能类别', '覆盖度 (%)',
        'effectiveness_function_coverage.png'
    )
    
    # 3. 效果测试 - 综合评分雷达图
    print("3. 生成效果测试 - 综合评分雷达图...")
    effectiveness_data = {
        '传统方式': {'性能': 75, '效率': 40, '功能': 55, '体验': 35, '准确': 65},
        'MCP方式': {'性能': 78, '效率': 95, '功能': 88, '体验': 92, '准确': 93}
    }
    
    generator.radar_chart(
        effectiveness_data,
        '传统方式 vs MCP方式 综合评分对比',
        'effectiveness_radar_score.png'
    )
    
    # 4. 性能测试 - 文档上传性能
    print("4. 生成性能测试 - 文档上传 图表...")
    doc_sizes = ['1MB', '5MB', '10MB']
    upload_times = [0.690, 2.062, 3.141]
    
    generator.bar_chart(
        upload_times, doc_sizes,
        '不同大小文档上传性能',
        '文档大小', '耗时 (秒)',
        'performance_upload_times.png',
        threshold=60
    )
    
    # 5. 性能测试 - 文档检索性能
    print("5. 生成性能测试 - 文档检索 图表...")
    search_types = ['关键词搜索', '语义搜索', '混合搜索', '并发搜索(50)', '批量查询']
    search_times = [0.275, 0.315, 0.296, 0.387, 0.249]
    
    generator.bar_chart(
        search_times, search_types,
        '不同搜索类型性能对比',
        '搜索类型', '平均响应时间 (秒)',
        'performance_search_times.png',
        threshold=1.0
    )
    
    # 6. 性能测试 - 并发性能
    print("6. 生成性能测试 - 并发性能 图表...")
    concurrency_levels = [10, 20, 50]
    avg_response_times = [0.292, 0.315, 0.387]
    
    generator.line_chart(
        {'平均响应时间': avg_response_times},
        '并发搜索性能测试',
        '并发数', '平均响应时间 (秒)',
        'performance_concurrency.png',
        threshold=1.0
    )
    
    # 7. 性能测试 - 文档上传流程分解（堆叠柱状图）
    print("7. 生成性能测试 - 上传流程分解 图表...")
    upload_process = {
        '1MB': {'文件上传': 0.245, '文档解析': 0.183, '向量生成': 0.152, 
                '索引构建': 0.072, '元数据存储': 0.038},
        '5MB': {'文件上传': 0.892, '文档解析': 0.623, '向量生成': 0.325,
                '索引构建': 0.147, '元数据存储': 0.075},
        '10MB': {'文件上传': 1.523, '文档解析': 0.892, '向量生成': 0.456,
                 '索引构建': 0.198, '元数据存储': 0.072}
    }
    
    generator.stacked_bar_chart(
        upload_process,
        '文档上传流程耗时分解',
        '文档大小', '耗时 (秒)',
        'performance_upload_breakdown.png'
    )
    
    # 8. 综合性能评分
    print("8. 生成综合性能评分饼图...")
    score_data = [98.5, 1.5]
    score_labels = ['实际得分', '剩余分数']
    
    generator.pie_chart(
        score_data, score_labels,
        '综合性能评分占比',
        'overall_score.png',
        explode=[0.1, 0]
    )
    
    # 9. 性能要求满足情况
    print("9. 生成性能要求满足情况对比图...")
    metrics = ['文档更新', '文档检索']
    requirements = [60, 1.0]
    actual_performance = [1.964, 0.304]
    
    generator.comparison_bar_chart(
        {'性能要求': dict(zip(metrics, requirements)), 
         '实际性能': dict(zip(metrics, actual_performance))},
        '性能要求 vs 实际性能',
        '性能指标', '耗时 (秒)',
        'performance_requirements.png'
    )
    
    # 10. 性能提升百分比
    print("10. 生成性能提升百分比柱状图...")
    improvements = ['操作步骤减少', '平均耗时减少', '交互次数减少', 
                   '上下文理解提升', '查询准确率提升']
    improvement_values = [82.5, 82.5, 75, 350, 38]
    
    generator.bar_chart(
        improvement_values, improvements,
        '用户体验提升指标',
        '提升维度', '提升幅度 (%)',
        'user_experience_improvement.png'
    )
    
    print()
    print("=" * 60)
    print("✓ 所有图表生成完成！")
    print(f"图表保存位置: {generator.output_dir}/")
    print("=" * 60)


if __name__ == '__main__':
    try:
        generate_all_charts()
    except Exception as e:
        print(f"错误: {e}")
        print("提示: 确保已安装 matplotlib 库")
        print("安装命令: pip install matplotlib")