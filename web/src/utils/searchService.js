import axios from 'axios'

// API基础URL
const apiBase = '/api/v1'

/**
 * 搜索服务工具类
 */
export function useSearchService() {
  /**
   * 执行搜索
   * @param {Object} request 搜索请求对象
   * @returns {Promise} 搜索结果Promise
   */
  const searchDocuments = async (request) => {
    try {
      console.log('发送搜索请求:', request)
      const response = await axios.post(`${apiBase}/search`, request)
      console.log('搜索响应原始数据:', response)
      console.log('搜索响应数据:', response.data)
      return response.data
    } catch (error) {
      console.error('搜索失败:', error)
      console.error('错误响应:', error.response)
      throw error
    }
  }

  /**
   * 通过GET方式执行搜索
   * @param {Object} params 搜索参数
   * @returns {Promise} 搜索结果Promise
   */
  const searchDocumentsGet = async (params) => {
    try {
      const response = await axios.get(`${apiBase}/search`, { params })
      return response.data
    } catch (error) {
      console.error('搜索失败:', error)
      throw error
    }
  }

  /**
   * 构建文档索引
   * @param {string} documentId 文档ID
   * @param {string} version 版本号
   * @returns {Promise} 操作结果Promise
   */
  const buildDocumentIndex = async (documentId, version) => {
    try {
      const response = await axios.post(`${apiBase}/search/documents/${documentId}/versions/${version}/index`)
      return response.data
    } catch (error) {
      console.error('构建索引失败:', error)
      throw error
    }
  }

  /**
   * 获取索引状态
   * @param {string} documentId 文档ID
   * @returns {Promise} 索引状态Promise
   */
  const getIndexingStatus = async (documentId) => {
    try {
      const response = await axios.get(`${apiBase}/search/documents/${documentId}/index/status`)
      return response.data
    } catch (error) {
      console.error('获取索引状态失败:', error)
      throw error
    }
  }

  /**
   * 删除索引
   * @param {string} documentId 文档ID
   * @returns {Promise} 操作结果Promise
   */
  const deleteIndex = async (documentId) => {
    try {
      const response = await axios.delete(`${apiBase}/search/documents/${documentId}/index`)
      return response.data
    } catch (error) {
      console.error('删除索引失败:', error)
      throw error
    }
  }

  /**
   * 删除指定版本的索引
   * @param {string} documentId 文档ID
   * @param {string} version 版本号
   * @returns {Promise} 操作结果Promise
   */
  const deleteIndexByVersion = async (documentId, version) => {
    try {
      const response = await axios.delete(`${apiBase}/search/documents/${documentId}/versions/${version}/index`)
      return response.data
    } catch (error) {
      console.error('删除索引失败:', error)
      throw error
    }
  }

  return {
    searchDocuments,
    searchDocumentsGet,
    buildDocumentIndex,
    getIndexingStatus,
    deleteIndex,
    deleteIndexByVersion
  }
}