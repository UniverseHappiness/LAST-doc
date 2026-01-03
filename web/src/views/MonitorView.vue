<template>
  <div class="monitor-view container-fluid">
    <h2 class="mb-4">系统监控</h2>

    <!-- 系统状态卡片 -->
    <div class="row mb-4">
      <div class="col-md-12">
        <div class="card">
          <div class="card-header d-flex justify-content-between align-items-center">
            <h5 class="mb-0">系统状态</h5>
            <span :class="`badge ${getStatusBadgeClass(systemStatus.overall_status)}`">
              {{ getStatusText(systemStatus.overall_status) }}
            </span>
          </div>
          <div class="card-body">
            <div class="row">
              <div class="col-md-3">
                <div class="status-item">
                  <small class="text-muted">CPU状态</small>
                  <div :class="`badge ${getStatusBadgeClass(systemStatus.cpu_status)}`">
                    {{ getStatusText(systemStatus.cpu_status) }}
                  </div>
                </div>
              </div>
              <div class="col-md-3">
                <div class="status-item">
                  <small class="text-muted">内存状态</small>
                  <div :class="`badge ${getStatusBadgeClass(systemStatus.memory_status)}`">
                    {{ getStatusText(systemStatus.memory_status) }}
                  </div>
                </div>
              </div>
              <div class="col-md-3">
                <div class="status-item">
                  <small class="text-muted">数据库状态</small>
                  <div :class="`badge ${getStatusBadgeClass(systemStatus.database_status)}`">
                    {{ getStatusText(systemStatus.database_status) }}
                  </div>
                </div>
              </div>
              <div class="col-md-3">
                <div class="status-item">
                  <small class="text-muted">服务状态</small>
                  <div :class="`badge ${getStatusBadgeClass(systemStatus.service_status)}`">
                    {{ getStatusText(systemStatus.service_status) }}
                  </div>
                </div>
              </div>
            </div>
            <div class="text-muted mt-2 small">
              最后更新: {{ formatTimestamp(systemStatus.timestamp) }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 当前指标卡片 -->
    <div class="row mb-4">
      <div class="col-md-12">
        <div class="card">
          <div class="card-header d-flex justify-content-between align-items-center">
            <h5 class="mb-0">当前系统指标</h5>
            <button class="btn btn-sm btn-primary" @click="refreshMetrics" :disabled="loading">
              <i class="bi bi-arrow-clockwise me-1"></i>刷新
            </button>
          </div>
          <div class="card-body" v-if="currentMetrics">
            <div class="row">
              <!-- CPU指标 -->
              <div class="col-md-4 mb-3">
                <div class="metric-card">
                  <h6 class="metric-title">CPU</h6>
                  <div class="metric-value">{{ currentMetrics.cpu_usage?.toFixed(2) || 0 }}%</div>
                  <div class="metric-details">
                    <small>核心数: {{ currentMetrics.cpu_cores || 0 }}</small><br>
                    <small>Goroutines: {{ currentMetrics.goroutine_count || 0 }}</small>
                  </div>
                </div>
              </div>

              <!-- 内存指标 -->
              <div class="col-md-4 mb-3">
                <div class="metric-card">
                  <h6 class="metric-title">内存</h6>
                  <div class="metric-value">{{ formatBytes(currentMetrics.memory_heap_alloc) }}</div>
                  <div class="metric-details">
                    <small>系统: {{ formatBytes(currentMetrics.memory_sys) }}</small><br>
                    <small>堆: {{ formatBytes(currentMetrics.memory_heap_sys) }}</small>
                  </div>
                </div>
              </div>

              <!-- GC指标 -->
              <div class="col-md-4 mb-3">
                <div class="metric-card">
                  <h6 class="metric-title">GC</h6>
                  <div class="metric-value">{{ currentMetrics.gc_num || 0 }}</div>
                  <div class="metric-details">
                    <small>暂停时间: {{ formatNanos(currentMetrics.gc_pause_total) }}</small><br>
                    <small>下次GC: {{ formatBytes(currentMetrics.gc_next) }}</small>
                  </div>
                </div>
              </div>

              <!-- 请求指标 -->
              <div class="col-md-4 mb-3">
                <div class="metric-card">
                  <h6 class="metric-title">请求</h6>
                  <div class="metric-value">{{ currentMetrics.request_count || 0 }}</div>
                  <div class="metric-details">
                    <small>错误: {{ currentMetrics.error_count || 0 }}</small><br>
                    <small>平均延迟: {{ currentMetrics.average_latency || 0 }}ms</small>
                  </div>
                </div>
              </div>

              <!-- 数据库指标 -->
              <div class="col-md-8 mb-3">
                <div class="metric-card">
                  <h6 class="metric-title">数据库连接</h6>
                  <div class="metric-value">{{ currentMetrics.db_connections || 0 }}</div>
                  <div class="metric-details">
                    <small>最大连接: {{ currentMetrics.db_max_open || 0 }}</small><br>
                    <small>使用中: {{ currentMetrics.db_in_use || 0 }} | 空闲: {{ currentMetrics.db_idle || 0 }}</small>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="card-body text-center" v-else>
            <div class="spinner-border text-primary" role="status">
              <span class="visually-hidden">加载中...</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 性能报告 -->
    <div class="row mb-4">
      <div class="col-md-12">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">性能报告（最近24小时）</h5>
          </div>
          <div class="card-body">
            <div class="row" v-if="performanceReport && performanceReport.length > 0">
              <div class="col-md-6 mb-3">
                <div class="chart-container">
                  <h6>CPU使用率</h6>
                  <canvas ref="cpuChart"></canvas>
                </div>
              </div>
              <div class="col-md-6 mb-3">
                <div class="chart-container">
                  <h6>内存使用</h6>
                  <canvas ref="memoryChart"></canvas>
                </div>
              </div>
              <div class="col-md-6 mb-3">
                <div class="chart-container">
                  <h6>请求数量</h6>
                  <canvas ref="requestChart"></canvas>
                </div>
              </div>
              <div class="col-md-6 mb-3">
                <div class="chart-container">
                  <h6>平均延迟</h6>
                  <canvas ref="latencyChart"></canvas>
                </div>
              </div>
            </div>
            <div class="alert alert-info" v-else>
              没有可用的性能报告数据。请等待指标收集或检查系统是否正常运行。
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 日志查看 -->
    <div class="row">
      <div class="col-md-12">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">系统日志</h5>
          </div>
          <div class="card-body">
            <!-- 日志过滤器 -->
            <div class="row mb-3">
              <div class="col-md-2">
                <select class="form-select form-select-sm" v-model="logFilter.level">
                  <option value="">全部级别</option>
                  <option value="debug">Debug</option>
                  <option value="info">Info</option>
                  <option value="warn">Warn</option>
                  <option value="error">Error</option>
                </select>
              </div>
              <div class="col-md-2">
                <select class="form-select form-select-sm" v-model="logFilter.service">
                  <option value="">全部服务</option>
                  <option value="document">文档服务</option>
                  <option value="search">搜索服务</option>
                  <option value="mcp">MCP服务</option>
                  <option value="user">用户服务</option>
                </select>
              </div>
              <div class="col-md-4">
                <input type="text" class="form-control form-control-sm" 
                       v-model="logFilter.message" 
                       placeholder="搜索消息">
              </div>
              <div class="col-md-2">
                <button class="btn btn-sm btn-primary w-100" @click="loadLogs">
                  <i class="bi bi-search me-1"></i>搜索
                </button>
              </div>
              <div class="col-md-2">
                <select class="form-select form-select-sm" v-model="logFilter.size" @change="loadLogs">
                  <option value="20">20条/页</option>
                  <option value="50">50条/页</option>
                  <option value="100">100条/页</option>
                </select>
              </div>
            </div>

            <!-- 日志列表 -->
            <div class="table-responsive">
              <table class="table table-sm table-striped">
                <thead>
                  <tr>
                    <th>时间</th>
                    <th>级别</th>
                    <th>服务</th>
                    <th>消息</th>
                    <th>用户</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="log in logs" :key="log.id">
                    <td small>{{ formatTimestamp(log.timestamp) }}</td>
                    <td>
                      <span :class="`badge ${getLogLevelClass(log.level)}`">
                        {{ log.level }}
                      </span>
                    </td>
                    <td>{{ log.service }}</td>
                    <td>{{ truncateMessage(log.message, 100) }}</td>
                    <td>{{ log.username || log.user_id || '-' }}</td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- 分页 -->
            <div class="d-flex justify-content-between align-items-center mt-3" v-if="logResponse">
              <small>共 {{ logResponse.total }} 条记录</small>
              <nav>
                <ul class="pagination pagination-sm mb-0">
                  <li class="page-item" :class="{ disabled: logResponse.page <= 1 }">
                    <a class="page-link" href="#" @click.prevent="changePage(logResponse.page - 1)">上一页</a>
                  </li>
                    <li class="page-item" :class="{ active: logResponse.page === page }" v-for="page in displayedPages" :key="page">
                    <a class="page-link" href="#" @click.prevent="changePage(page)">{{ page }}</a>
                  </li>
                  <li class="page-item" :class="{ disabled: logResponse.page >= logResponse.total_pages }">
                    <a class="page-link" href="#" @click.prevent="changePage(logResponse.page + 1)">下一页</a>
                  </li>
                </ul>
              </nav>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 数据清理 -->
    <div class="row mt-4 mb-4">
      <div class="col-md-12">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">数据清理</h5>
          </div>
          <div class="card-body">
            <div class="row align-items-center">
              <div class="col-md-8">
                <label class="form-label">保留天数:</label>
                <input type="number" class="form-control" v-model="retentionDays" min="1" max="365">
                <small class="text-muted">将删除 {{ retentionDays }} 天前的监控数据和日志</small>
              </div>
              <div class="col-md-4">
                <button class="btn btn-danger w-100" @click="cleanupOldData" :disabled="cleaning">
                  <i class="bi bi-trash me-1"></i>
                  {{ cleaning ? '清理中...' : '清理旧数据' }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import monitorService from '../utils/monitorService';
import Chart from 'chart.js/auto';
import authService from '../utils/authService';

export default {
  name: 'MonitorView',
  
  data() {
    return {
      loading: false,
      currentMetrics: null,
      systemStatus: {
        overall_status: 'healthy',
        cpu_status: 'healthy',
        memory_status: 'healthy',
        database_status: 'healthy',
        service_status: 'healthy',
        timestamp: new Date()
      },
      performanceReport: [],
      logs: [],
      logResponse: null,
      logFilter: {
        level: '',
        service: '',
        message: '',
        page:1,
        size: 20
      },
      retentionDays: 30,
      cleaning: false,
      charts: {},
      refreshInterval: null
    };
  },

  computed: {
    displayedPages() {
      if (!this.logResponse) return [];
      
      const totalPages = this.logResponse.total_pages;
      const currentPage = this.logResponse.page;
      const pages = [];
      
      // 显示最多5个页码
      let startPage = Math.max(1, currentPage - 2);
      let endPage = Math.min(totalPages, startPage + 4);
      
      startPage = Math.max(1, endPage - 4);
      
      for (let i = startPage; i <= endPage; i++) {
        pages.push(i);
      }
      
      return pages;
    }
  },

  mounted() {
    // 检查是否是管理员
    if (!authService.isAdmin()) {
      this.$router.push('/list');
      return;
    }
    
    this.loadData();
    
    // 每30秒自动刷新一次
    this.refreshInterval = setInterval(() => {
      this.refreshMetrics();
    }, 30000);
  },

  beforeUnmount() {
    // 清除定时器
    if (this.refreshInterval) {
      clearInterval(this.refreshInterval);
    }
    
    // 销毁图表
    Object.values(this.charts).forEach(chart => {
      if (chart) {
        chart.destroy();
      }
    });
  },

  methods: {
    async loadData() {
      await Promise.all([
        this.refreshMetrics(),
        this.loadPerformanceReport(),
        this.loadLogs()
      ]);
    },

    async refreshMetrics() {
      this.loading = true;
      
      try {
        const [metrics, status] = await Promise.all([
          monitorService.getCurrentMetrics(),
          monitorService.getSystemStatus()
        ]);
        
        this.currentMetrics = metrics;
        this.systemStatus = status;
      } catch (error) {
        console.error('刷新指标失败:', error);
        this.showError('刷新指标失败: ' + (error.response?.data?.error || error.message));
      } finally {
        this.loading = false;
      }
    },

    async loadPerformanceReport() {
      try {
        const startTime = new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString();
        const endTime = new Date().toISOString();
        
        const response = await monitorService.getPerformanceReport(startTime, endTime);
        
        console.log('性能报告响应结构:', {
          hasCurrent: !!response.current,
          hasHistory: !!response.history,
          historyLength: response.history ? response.history.length : 0,
          hasAverage: !!response.average,
          hasStatus: !!response.status
        });
        
        // 使用history数组作为performanceReport
        this.performanceReport = response.history || [];
        
        console.log('performanceReport数据点数:', this.performanceReport.length);
        
        // 渲染图表
        this.$nextTick(() => {
          if (response.history && response.history.length > 0) {
            this.renderCharts(response.history);
          }
        });
      } catch (error) {
        console.error('加载性能报告失败:', error);
        this.showError('加载性能报告失败');
      }
    },

    async loadLogs() {
      try {
        const response = await monitorService.getLogs(this.logFilter);
        this.logResponse = response;
        this.logs = response.logs;
      } catch (error) {
        console.error('加载日志失败:', error);
        this.showError('加载日志失败');
      }
    },

    changePage(page) {
      if (page < 1 || page > this.logResponse.total_pages) return;
      
      this.logFilter.page = page;
      this.loadLogs();
    },

    async cleanupOldData() {
      if (!confirm(`确定要删除 ${this.retentionDays} 天前的所有监控数据和日志吗？此操作不可恢复！`)) {
        return;
      }
      
      this.cleaning = true;
      
      try {
        await monitorService.cleanupOldData(this.retentionDays);
        this.showSuccess('旧数据清理成功');
        this.loadData();
      } catch (error) {
        console.error('清理失败:', error);
        this.showError('清理失败: ' + (error.response?.data?.error || error.message));
      } finally {
        this.cleaning = false;
      }
    },

    renderCharts(report) {
      // 检查响应数据是否存在且为数组
      if (!report || !Array.isArray(report)) {
        console.warn('性能报告数据为空或不是数组:', report);
        return;
      }
      
      if (report.length === 0) {
        console.warn('性能报告历史数据为空');
        return;
      }

      console.log('开始渲染图表，数据点数:', report.length);
      
      // 准备数据
      // 使用与formatTimestamp相同的逻辑，直接从时间字符串中提取部分，避免时区转换
      const labels = report.map(m => {
        if (!m.timestamp) return '-';
        
        let timestampStr = String(m.timestamp);
        
        // 如果是RFC3339格式（如 2026-01-03T21:14:11.013042Z）
        if (timestampStr.includes('T')) {
          const match = timestampStr.match(/(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2}):(\d{2})/);
          if (match) {
            const [, year, month, day, hours, minutes, seconds] = match;
            return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
          }
        }
        
        // 如果是数据库格式 (2026-01-03 21:02:50)，直接显示
        if (timestampStr.match(/^\d{4}-\d{2}-\d{2}[ T]\d{2}:\d{2}:\d{2}/)) {
          return timestampStr.replace(' ', ' ');
        }
        
        // 默认情况
        return timestampStr;
      });
      
      const cpuData = report.map(m => m.cpu_usage || 0);
      const memoryData = report.map(m => (m.memory_heap_alloc || 0) / 1024 / 1024);
      const requestData = report.map(m => m.request_count || 0);
      const latencyData = report.map(m => m.average_latency || 0);

      console.log('图表数据准备完成:', {
        labelsCount: labels.length,
        cpuDataCount: cpuData.length,
        memoryDataCount: memoryData.length
      });
      
      // 渲染CPU图表
      this.renderChart('cpuChart', 'line', {
        labels: labels,
        datasets: [{
          label: 'CPU使用率 (%)',
          data: cpuData,
          borderColor: 'rgb(75, 192, 192)',
          tension: 0.1
        }]
      });

      // 渲染内存图表
      this.renderChart('memoryChart', 'line', {
        labels: labels,
        datasets: [{
          label: '内存使用',
          data: memoryData,
          borderColor: 'rgb(54, 162, 235)',
          tension: 0.1
        }]
      });

      // 渲染请求图表
      this.renderChart('requestChart', 'line', {
        labels: labels,
        datasets: [{
          label: '请求数量',
          data: requestData,
          borderColor: 'rgb(255, 99, 132)',
          tension: 0.1
        }]
      });

      // 渲染延迟图表
      this.renderChart('latencyChart', 'line', {
        labels: labels,
        datasets: [{
          label: '平均延迟',
          data: latencyData,
          borderColor: 'rgb(153, 102, 255)',
          tension: 0.1
        }]
      });
      
      console.log('图表渲染完成');
    },

    renderChart(refName, type, data) {
      const canvas = this.$refs[refName];
      if (!canvas) return;
      
      // 销毁旧图表
      if (this.charts[refName]) {
        this.charts[refName].destroy();
      }
      
      this.charts[refName] = new Chart(canvas, {
        type: type,
        data: data,
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: {
            legend: {
              display: true,
              position: 'top'
            }
          },
          scales: {
            x: {
              display: true,
              grid: {
                display: false
              }
            },
            y: {
              beginAtZero: true,
              display: true,
              grid: {
                display: true,
                color: '#e0e0e0'
              }
            }
          },
          layout: {
            padding: {
              top: 10,
              right: 10,
              bottom: 10,
              left: 10
            }
          }
        }
      });
    },

    getStatusBadgeClass(status) {
      switch (status) {
        case 'healthy':
          return 'bg-success';
        case 'warning':
          return 'bg-warning';
        case 'critical':
          return 'bg-danger';
        default:
          return 'bg-secondary';
      }
    },

    getStatusText(status) {
      switch (status) {
        case 'healthy':
          return '健康';
        case 'warning':
          return '警告';
        case 'critical':
          return '严重';
        default:
          return '未知';
      }
    },

    getLogLevelClass(level) {
      switch (level) {
        case 'debug':
          return 'bg-secondary';
        case 'info':
          return 'bg-info';
        case 'warn':
          return 'bg-warning';
        case 'error':
          return 'bg-danger';
        default:
          return 'bg-secondary';
      }
    },

    formatTimestamp(timestamp) {
      if (!timestamp) return '-';
      
      // 转换为字符串
      let timestampStr = String(timestamp);
      
      // 如果是RFC3339格式（如 2026-01-03T21:14:11.013042Z 或 2026-01-03T21:14:11+08:00）
      // 直接提取日期和时间部分，忽略时区信息
      if (timestampStr.includes('T')) {
        // 从RFC3339格式中提取 YYYY-MM-DD HH:MM:SS
        // 格式1: 2026-01-03T21:14:11.013042Z
        // 格式2: 2026-01-03T21:14:11+08:00
        const match = timestampStr.match(/(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2}):(\d{2})/);
        if (match) {
          const [, year, month, day, hours, minutes, seconds] = match;
          // 直接返回提取的本地时间，不进行时区转换
          return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
        }
      }
      
      // 如果是数据库格式 (2026-01-03 21:02:50)，直接显示
      if (timestampStr.match(/^\d{4}-\d{2}-\d{2}[ T]\d{2}:\d{2}:\d{2}/)) {
        return timestampStr.replace(' ', ' ');
      }
      
      // 如果是Date对象，使用toLocaleString
      if (timestamp instanceof Date) {
        const date = timestamp;
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const day = String(date.getDate()).padStart(2, '0');
        const hours = String(date.getHours()).padStart(2, '0');
        const minutes = String(date.getMinutes()).padStart(2, '0');
        const seconds = String(date.getSeconds()).padStart(2, '0');
        return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
      }
      
      // 无法识别格式，返回原始字符串
      return timestampStr;
    },

    formatBytes(bytes) {
      if (!bytes) return '0 B';
      const sizes = ['B', 'KB', 'MB', 'GB'];
      const i = Math.floor(Math.log(bytes) / Math.log(1024));
      return (bytes / Math.pow(1024, i)).toFixed(2) + ' ' + sizes[i];
    },

    formatNanos(nanos) {
      if (!nanos) return '0ms';
      const ms = nanos / 1000000;
      if (ms >= 1000) {
        return (ms / 1000).toFixed(2) + 's';
      }
      return ms.toFixed(2) + 'ms';
    },

    truncateMessage(message, maxLength) {
      if (!message) return '';
      if (message.length <= maxLength) return message;
      return message.substring(0, maxLength) + '...';
    },

    showSuccess(message) {
      alert(message);
    },

    showError(message) {
      alert('错误: ' + message);
    }
  }
};
</script>

<style scoped>
.monitor-view {
  max-width: 1400px;
}

.status-item {
  text-align: center;
}

.metric-card {
  padding: 1rem;
  border: 1px solid #dee2e6;
  border-radius: 0.25rem;
  text-align: center;
}

.metric-title {
  color: #6c757d;
  margin-bottom: 0.5rem;
}

.metric-value {
  font-size: 1.5rem;
  font-weight: bold;
  color: #0d6efd;
  margin-bottom: 0.5rem;
}

.metric-details {
  font-size: 0.875rem;
  color: #6c757d;
}

.badge {
  padding: 0.5em 0.75em;
}

.chart-container {
  position: relative;
  height: 300px;
  width: 100%;
  padding: 10px;
}

canvas {
  position: relative;
  height: 100% !important;
  width: 100% !important;
}
</style>