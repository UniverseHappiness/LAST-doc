import axios from 'axios'

export function useDocumentService(state) {
  const {
    documents,
    documentVersions,
    pagination,
    filters,
    searchQuery,
    isLoading,
    isUploading,
    uploadForm,
    currentView
  } = state

  const apiBase = '/api/v1'

  // 获取文档列表
  const fetchDocuments = async () => {
    isLoading.value = true
    try {
      const params = {
        page: pagination.page,
        size: pagination.size
      }
      
      // 添加过滤条件
      if (filters.library) params.library = filters.library
      if (filters.type) params.type = filters.type
      if (filters.version) params.version = filters.version
      if (filters.status) params.status = filters.status
      
      const response = await axios.get(`${apiBase}/documents`, { params })
      documents.value = response.data.data.items
      pagination.total = response.data.data.total
    } catch (error) {
      console.error('获取文档列表失败:', error)
      alert('获取文档列表失败: ' + (error.response?.data?.message || error.message))
    } finally {
      isLoading.value = false
    }
  }

  // 获取文档版本列表
  const fetchDocumentVersions = async (documentId) => {
    isLoading.value = true
    try {
      const response = await axios.get(`${apiBase}/documents/${documentId}/versions`)
      documentVersions.value = response.data.data
    } catch (error) {
      console.error('获取文档版本失败:', error)
      alert('获取文档版本失败: ' + (error.response?.data?.message || error.message))
    } finally {
      isLoading.value = false
    }
  }

  // 获取最新版本
  const getLatestVersion = async (documentId) => {
    try {
      const response = await axios.get(`${apiBase}/documents/${documentId}/versions/latest`)
      return response.data.data
    } catch (error) {
      console.error('获取最新版本失败:', error)
      throw error
    }
  }

  // 上传文档
  const uploadDocument = async () => {
    console.log('DEBUG: documentService - uploadDocument 函数被调用')
    console.log('DEBUG: documentService - uploadForm.file 当前值:', uploadForm.file)
    
    if (!uploadForm.file) {
      console.error('DEBUG: documentService - 没有选择文件，uploadForm.file 为空')
      alert('请选择文件')
      return
    }
    
    console.log('DEBUG: documentService - 文件对象详细信息:', {
      name: uploadForm.file.name,
      size: uploadForm.file.size,
      type: uploadForm.file.type,
      lastModified: uploadForm.file.lastModified
    })
    
    isUploading.value = true
    const formData = new FormData()
    formData.append('file', uploadForm.file)
    formData.append('name', uploadForm.name)
    formData.append('type', uploadForm.type)
    formData.append('category', uploadForm.category)
    formData.append('version', uploadForm.version)
    formData.append('library', uploadForm.library)
    formData.append('description', uploadForm.description)
    formData.append('tags', uploadForm.tags)
    
    try {
      console.log('DEBUG: documentService - 前端开始上传文档 - 版本号:', uploadForm.version)
      console.log('DEBUG: documentService - 前端上传文档数据:', {
        name: uploadForm.name,
        type: uploadForm.type,
        version: uploadForm.version,
        library: uploadForm.library
      })
      
      // 添加表单数据的调试输出
      console.log('DEBUG: documentService - FormData中的type字段值:', uploadForm.type)
      console.log('DEBUG: documentService - FormData内容检查:')
      for (let [key, value] of formData.entries()) {
        if (key === 'file') {
          console.log(`  ${key}: [File object - name: ${value.name}, size: ${value.size}]`)
        } else {
          console.log(`  ${key}: ${value}`)
        }
      }
      
      const response = await axios.post(`${apiBase}/documents`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
      
      console.log('DEBUG: documentService - 前端上传文档成功:', response.data)
      alert('上传成功')
      resetUploadForm()
      currentView.value = 'list'
      await fetchDocuments()
    } catch (error) {
      console.error('DEBUG: documentService - 前端上传文档失败:', error)
      console.error('DEBUG: documentService - 前端上传文档失败详情:', {
        response: error.response?.data,
        status: error.response?.status,
        message: error.message
      })
      alert('上传文档失败: ' + (error.response?.data?.message || error.message))
    } finally {
      isUploading.value = false
    }
  }

  // 删除文档
  const deleteDocument = async (doc) => {
    if (!confirm(`确定要删除文档 "${doc.name}" 吗？`)) {
      return
    }
    
    isLoading.value = true
    try {
      await axios.delete(`${apiBase}/documents/${doc.id}`)
      alert('删除成功')
      await fetchDocuments()
    } catch (error) {
      console.error('删除文档失败:', error)
      alert('删除文档失败: ' + (error.response?.data?.message || error.message))
    } finally {
      isLoading.value = false
    }
  }

  // 删除文档版本
  const deleteDocumentVersion = async (version) => {
    console.log('DEBUG: 尝试删除文档版本', version)
    if (!confirm(`确定要删除版本 "${version.version}" 吗？`)) {
      return
    }
    
    isLoading.value = true
    try {
      await axios.delete(`${apiBase}/documents/${version.document_id}/versions/${version.version}`)
      alert('删除成功')
      // 重新获取文档版本列表
      await fetchDocumentVersions(version.document_id)
    } catch (error) {
      console.error('删除文档版本失败:', error)
      alert('删除文档版本失败: ' + (error.response?.data?.message || error.message))
    } finally {
      isLoading.value = false
    }
  }

  // 搜索文档
  const searchDocuments = async () => {
    if (!searchQuery.value.trim()) {
      await fetchDocuments()
      return
    }
    
    isLoading.value = true
    try {
      const params = {
        page: 1,
        size: pagination.size,
        name: searchQuery.value.trim()
      }
      
      const response = await axios.get(`${apiBase}/documents`, { params })
      documents.value = response.data.data.items
      pagination.total = response.data.data.total
      pagination.page = 1
    } catch (error) {
      console.error('搜索文档失败:', error)
      alert('搜索文档失败: ' + (error.response?.data?.message || error.message))
    } finally {
      isLoading.value = false
    }
  }

  // 应用过滤器
  const applyFilters = async () => {
    pagination.page = 1
    await fetchDocuments()
  }

  // 切换页面
  const changePage = async (page) => {
    if (page < 1 || (page - 1) * pagination.size >= pagination.total) {
      return
    }
    pagination.page = page
    await fetchDocuments()
  }

  // 更新文档
  const updateDocument = async (documentId, updates) => {
    isLoading.value = true
    try {
      const response = await axios.put(`${apiBase}/documents/${documentId}`, updates)
      alert('更新成功')
      return response.data
    } catch (error) {
      console.error('更新文档失败:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  // 获取分页页码
  const getPaginationPages = () => {
    const totalPages = Math.ceil(pagination.total / pagination.size)
    const pages = []
    const maxVisiblePages = 5
    
    let startPage = Math.max(1, pagination.page - Math.floor(maxVisiblePages / 2))
    let endPage = Math.min(totalPages, startPage + maxVisiblePages - 1)
    
    if (endPage - startPage + 1 < maxVisiblePages) {
      startPage = Math.max(1, endPage - maxVisiblePages + 1)
    }
    
    for (let i = startPage; i <= endPage; i++) {
      pages.push(i)
    }
    
    return pages
  }

  // 重置上传表单
  const resetUploadForm = () => {
    uploadForm.name = ''
    uploadForm.type = ''
    uploadForm.category = ''
    uploadForm.version = ''
    uploadForm.library = ''
    uploadForm.description = ''
    uploadForm.tags = ''
    uploadForm.file = null
  }

  // 获取状态徽章样式
  const getStatusBadgeClass = (status) => {
    switch (status) {
      case 'uploading': return 'bg-secondary'
      case 'processing': return 'bg-warning'
      case 'completed': return 'bg-success'
      case 'failed': return 'bg-danger'
      default: return 'bg-light text-dark'
    }
  }

  // 获取状态文本
  const getStatusText = (status) => {
    switch (status) {
      case 'uploading': return '上传中'
      case 'processing': return '处理中'
      case 'completed': return '已完成'
      case 'failed': return '失败'
      default: return status
    }
  }

  // 格式化文件大小
  const formatFileSize = (bytes) => {
    if (bytes === 0) return '0 Bytes'
    
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  // 格式化日期
  const formatDate = (dateString) => {
    const date = new Date(dateString)
    return date.toLocaleString('zh-CN')
  }

  // AI友好格式相关方法
  const getStructuredContent = async (documentId, version) => {
    try {
      const response = await axios.get(`${apiBase}/ai-format/documents/${documentId}/versions/${version}/structured`)
      return response.data
    } catch (error) {
      console.error('获取结构化内容失败:', error)
      throw error
    }
  }

  const generateLLMFormat = async (documentId, version, options) => {
    try {
      const response = await axios.post(`${apiBase}/ai-format/documents/${documentId}/versions/${version}/llm`, options)
      return response.data
    } catch (error) {
      console.error('生成LLM优化格式失败:', error)
      throw error
    }
  }

  const generateMultiGranularity = async (documentId, version) => {
    try {
      const response = await axios.get(`${apiBase}/ai-format/documents/${documentId}/versions/${version}/multigranularity`)
      return response.data
    } catch (error) {
      console.error('生成多粒度表示失败:', error)
      throw error
    }
  }

  const injectContext = async (documentId, version, options) => {
    try {
      const response = await axios.post(`${apiBase}/ai-format/documents/${documentId}/versions/${version}/context`, options)
      return response.data
    } catch (error) {
      console.error('注入上下文失败:', error)
      throw error
    }
  }

  const getAIFriendlyFormats = async (documentId, version, options = {}) => {
    try {
      const params = new URLSearchParams()
      if (options.type) params.append('type', options.type)
      if (options.maxTokens) params.append('max_tokens', options.maxTokens)
      if (options.includeCode !== undefined) params.append('include_code', options.includeCode)
      if (options.summaryLevel) params.append('summary_level', options.summaryLevel)
      
      const queryString = params.toString()
      const url = `${apiBase}/ai-format/documents/${documentId}/versions/${version}${queryString ? '?' + queryString : ''}`
      
      const response = await axios.get(url)
      return response.data
    } catch (error) {
      console.error('获取AI友好格式失败:', error)
      throw error
    }
  }

  // 导出传统的组合式API
  return {
    fetchDocuments,
    fetchDocumentVersions,
    getLatestVersion,
    uploadDocument,
    updateDocument,
    deleteDocument,
    deleteDocumentVersion,
    searchDocuments,
    applyFilters,
    changePage,
    getPaginationPages,
    getStatusBadgeClass,
    getStatusText,
    formatFileSize,
    formatDate,
    getStructuredContent,
    generateLLMFormat,
    generateMultiGranularity,
    injectContext,
    getAIFriendlyFormats
  }
}

// 导出独立的函数式API
export const documentService = {
  getDocuments: async (params = {}) => {
    try {
      const response = await axios.get('/api/v1/documents', { params })
      return response.data
    } catch (error) {
      console.error('获取文档列表失败:', error)
      throw error
    }
  },

  getDocumentVersions: async (documentId) => {
    try {
      const response = await axios.get(`/api/v1/documents/${documentId}/versions`)
      return response.data
    } catch (error) {
      console.error('获取文档版本失败:', error)
      throw error
    }
  },

  getDocumentByVersion: async (documentId, version) => {
    try {
      const response = await axios.get(`/api/v1/documents/${documentId}/versions/${version}`)
      return response.data
    } catch (error) {
      console.error('获取文档版本失败:', error)
      throw error
    }
  },

  getStructuredContent: async (documentId, version) => {
    try {
      const response = await axios.get(`/api/v1/ai-format/documents/${documentId}/versions/${version}/structured`)
      return response.data
    } catch (error) {
      console.error('获取结构化内容失败:', error)
      throw error
    }
  },

  generateLLMFormat: async (documentId, version, options) => {
    try {
      const response = await axios.post(`/api/v1/ai-format/documents/${documentId}/versions/${version}/llm`, options)
      return response.data
    } catch (error) {
      console.error('生成LLM优化格式失败:', error)
      throw error
    }
  },

  generateMultiGranularity: async (documentId, version) => {
    try {
      const response = await axios.get(`/api/v1/ai-format/documents/${documentId}/versions/${version}/multigranularity`)
      return response.data
    } catch (error) {
      console.error('生成多粒度表示失败:', error)
      throw error
    }
  },

  injectContext: async (documentId, version, options) => {
    try {
      const response = await axios.post(`/api/v1/ai-format/documents/${documentId}/versions/${version}/context`, options)
      return response.data
    } catch (error) {
      console.error('注入上下文失败:', error)
      throw error
    }
  },

  getAIFriendlyFormats: async (documentId, version, options = {}) => {
    try {
      const params = new URLSearchParams()
      if (options.type) params.append('type', options.type)
      if (options.maxTokens) params.append('max_tokens', options.maxTokens)
      if (options.includeCode !== undefined) params.append('include_code', options.includeCode)
      if (options.summaryLevel) params.append('summary_level', options.summaryLevel)
      
      const queryString = params.toString()
      const url = `/api/v1/ai-format/documents/${documentId}/versions/${version}${queryString ? '?' + queryString : ''}`
      
      const response = await axios.get(url)
      return response.data
    } catch (error) {
      console.error('获取AI友好格式失败:', error)
      throw error
    }
  }
}