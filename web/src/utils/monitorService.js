// 系统监控服务
import axios from 'axios';

const API_BASE_URL = '/api/v1/monitor';

class MonitorService {
  // 获取当前系统指标
  async getCurrentMetrics() {
    try {
      const response = await axios.get(`${API_BASE_URL}/metrics/current`);
      return response.data.metrics;
    } catch (error) {
      throw error.response?.data || { error: '获取系统指标失败' };
    }
  }

  // 获取指标历史数据
  async getMetricsHistory(startTime, endTime) {
    try {
      const params = {};
      if (startTime) params.start_time = startTime;
      if (endTime) params.end_time = endTime;
      
      const response = await axios.get(`${API_BASE_URL}/metrics/history`, { params });
      return response.data.metrics;
    } catch (error) {
      throw error.response?.data || { error: '获取指标历史失败' };
    }
  }

  // 获取指标报告
  async getMetricsReport(duration = '1h') {
    try {
      const response = await axios.get(`${API_BASE_URL}/metrics/report`, {
        params: { duration }
      });
      return response.data;
    } catch (error) {
      throw error.response?.data || { error: '获取指标报告失败' };
    }
  }

  // 获取系统状态
  async getSystemStatus() {
    try {
      const response = await axios.get(`${API_BASE_URL}/status`);
      return response.data;
    } catch (error) {
      throw error.response?.data || { error: '获取系统状态失败' };
    }
  }

  // 获取日志列表
  async getLogs(filter = {}) {
    try {
      const response = await axios.get(`${API_BASE_URL}/logs`, { params: filter });
      return response.data;
    } catch (error) {
      throw error.response?.data || { error: '获取日志失败' };
    }
  }

  // 获取日志统计
  async getLogStats(startTime, endTime) {
    try {
      const params = {};
      if (startTime) params.start_time = startTime;
      if (endTime) params.end_time = endTime;
      
      const response = await axios.get(`${API_BASE_URL}/logs/stats`, { params });
      return response.data.stats;
    } catch (error) {
      throw error.response?.data || { error: '获取日志统计失败' };
    }
  }

  // 获取性能报告
  async getPerformanceReport(startTime, endTime) {
    try {
      const params = {};
      if (startTime) params.start_time = startTime;
      if (endTime) params.end_time = endTime;
      
      // 修复：使用正确的API路径 /metrics/report
      const response = await axios.get(`${API_BASE_URL}/metrics/report`, { params });
      return response.data;
    } catch (error) {
      throw error.response?.data || { error: '获取性能报告失败' };
    }
  }

  // 清理旧数据
  async cleanupOldData(retentionDays = 30) {
    try {
      const response = await axios.post(`${API_BASE_URL}/cleanup`, null, {
        params: { retention_days: retentionDays }
      });
      return response.data;
    } catch (error) {
      throw error.response?.data || { error: '清理旧数据失败' };
    }
  }
}

export default new MonitorService();