<template>
  <div>
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h2>版本查看 - {{ documentVersion?.version }}</h2>
      <button class="btn btn-outline-secondary" @click="$emit('back')">返回</button>
    </div>
    
    <div v-if="documentVersion" class="card">
      <div class="card-header d-flex justify-content-between align-items-center">
        <div>
          <h5 class="mb-0">{{ documentVersion.version }}</h5>
          <p class="text-muted mb-0">{{ documentVersion.description || '暂无描述' }}</p>
        </div>
        <div>
          <span class="badge" :class="getStatusBadgeClass(documentVersion.status)">{{ getStatusText(documentVersion.status) }}</span>
          <span class="badge bg-light text-dark ms-1">{{ formatFileSize(documentVersion.file_size) }}</span>
        </div>
      </div>
      <div class="card-body">
        <!-- 文档操作按钮 -->
        <div class="mb-3">
          <button class="btn btn-sm btn-primary me-2" @click="downloadVersion">
            <i class="bi bi-download me-1"></i>下载
          </button>
        </div>
        
        <!-- 文档内容预览 -->
        <div class="document-content">
          <h6 class="mb-3">版本内容</h6>
          <div class="border rounded p-3 bg-light">
            <div v-if="!documentVersion.content || documentVersion.content === ''" class="mb-0">
              <p class="text-muted">暂无内容 - 状态: {{ getStatusText(documentVersion.status) }}</p>
              <p class="text-muted small">如果状态为"处理中"，请稍后刷新页面</p>
            </div>
            <div v-else>
              <pre v-if="documentType === 'markdown' || documentType === 'java_doc'" class="mb-0">{{ documentVersion.content }}</pre>
              <div v-else-if="documentType === 'swagger' || documentType === 'openapi'" class="mb-0">
                <pre>{{ JSON.stringify(JSON.parse(documentVersion.content), null, 2) }}</pre>
              </div>
              <div v-else class="mb-0">
                <p>文档内容解析中，请稍后再试...</p>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 版本信息 -->
        <div class="version-info mt-4">
          <h6 class="mb-3">版本信息</h6>
          <table class="table table-sm">
            <tbody>
              <tr>
                <th scope="row" style="width: 120px;">版本ID</th>
                <td>{{ documentVersion.id }}</td>
              </tr>
              <tr>
                <th scope="row">文档ID</th>
                <td>{{ documentVersion.document_id }}</td>
              </tr>
              <tr>
                <th scope="row">文件路径</th>
                <td>{{ documentVersion.file_path }}</td>
              </tr>
              <tr>
                <th scope="row">文件大小</th>
                <td>{{ formatFileSize(documentVersion.file_size) }}</td>
              </tr>
              <tr>
                <th scope="row">创建时间</th>
                <td>{{ formatDate(documentVersion.created_at) }}</td>
              </tr>
              <tr>
                <th scope="row">更新时间</th>
                <td>{{ formatDate(documentVersion.updated_at) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
    
    <!-- 加载状态 -->
    <div v-if="isLoading" class="text-center py-5">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">加载中...</span>
      </div>
      <p class="mt-2">正在加载版本内容...</p>
    </div>
    
    <!-- 错误状态 -->
    <div v-if="error" class="alert alert-danger" role="alert">
      {{ error }}
    </div>
  </div>
</template>

<script>
import { ref, onMounted, watch, watchEffect } from 'vue'

export default {
  name: 'DocumentVersionView',
  props: {
    documentId: {
      type: String,
      required: true
    },
    version: {
      type: String,
      required: false
    },
    globalState: {
      type: Object,
      required: false
    }
  },
  emits: ['back', 'download-version'],
  setup(props, { emit }) {
    console.log('DEBUG: DocumentVersionView组件初始化 - props:', props)
    const documentVersion = ref(null)
    const documentType = ref(null)
    const isLoading = ref(false)
    const error = ref(null)
    
    // 获取文档版本详情
    const fetchDocumentVersion = async () => {
      console.log('DEBUG: 开始执行fetchDocumentVersion函数 - props.documentId:', props.documentId, 'props.version:', props.version)
      
      // 尝试多种方式获取版本号
      let versionToUse = props.version
      
      // 如果props.version为空，尝试从全局状态获取
      if (!versionToUse && props.globalState && props.globalState.isSet) {
        console.log('DEBUG: props.version为空，尝试从全局状态获取版本号')
        versionToUse = props.globalState.version
        console.log('DEBUG: 从全局状态获取到版本号:', versionToUse)
      }
      
      // 如果仍然为空，尝试从URL参数获取
      if (!versionToUse) {
        console.log('DEBUG: props.version为空，尝试从URL参数获取版本号')
        const urlParams = new URLSearchParams(window.location.search)
        versionToUse = urlParams.get('version')
        console.log('DEBUG: 从URL参数获取到版本号:', versionToUse)
      }
      
      // 如果仍然为空，尝试从localStorage获取
      if (!versionToUse) {
        versionToUse = getVersionFromStorage()
        console.log('DEBUG: 从localStorage获取到版本号:', versionToUse)
      }
      
      if (!props.documentId || !versionToUse) {
        console.log('DEBUG: props.documentId或versionToUse为空，函数静默返回')
        return // 不显示错误，静默返回
      }
      
      isLoading.value = true
      error.value = null
      
      try {
        console.log('DEBUG: 前端开始获取文档版本 - 文档ID:', props.documentId, '版本:', versionToUse)
        const response = await fetch(`/api/v1/documents/${props.documentId}/versions/${versionToUse}`)
        const result = await response.json()
        
        console.log('DEBUG: 前端获取文档版本响应:', result)
        
        if (result.code === 200) {
          documentVersion.value = result.data
          console.log('DEBUG: 前端获取文档版本成功，数据:', documentVersion.value)
          console.log('DEBUG: 版本内容长度:', documentVersion.value.content ? documentVersion.value.content.length : '空')
          console.log('DEBUG: 版本状态:', documentVersion.value.status)
          console.log('DEBUG: 版本内容前100个字符:', documentVersion.value.content ? documentVersion.value.content.substring(0, 100) : '空')
          // 同时获取文档信息以确定文档类型
          fetchDocument()
        } else {
          console.error('DEBUG: 前端获取文档版本失败:', result.message)
          error.value = result.message || '获取文档版本失败'
        }
      } catch (err) {
        console.error('DEBUG: 前端获取文档版本异常:', err)
        error.value = '获取文档版本失败: ' + err.message
      } finally {
        isLoading.value = false
      }
    }
    
    // 尝试从localStorage获取版本号
    const getVersionFromStorage = () => {
      try {
        const storedVersion = localStorage.getItem('selectedVersion')
        console.log('DEBUG: 从localStorage获取版本号:', storedVersion)
        return storedVersion
      } catch (e) {
        console.error('DEBUG: 从localStorage获取版本号失败:', e)
        return null
      }
    }
    
    // 获取文档信息
    const fetchDocument = async () => {
      try {
        console.log('DEBUG: 前端开始获取文档信息 - 文档ID:', props.documentId)
        const response = await fetch(`/api/v1/documents/${props.documentId}`)
        const result = await response.json()
        
        console.log('DEBUG: 前端获取文档信息响应:', result)
        
        if (result.code === 200) {
          documentType.value = result.data.type
          console.log('DEBUG: 前端获取文档类型成功 - 类型:', documentType.value)
        } else {
          console.error('DEBUG: 前端获取文档类型失败:', result.message)
        }
      } catch (err) {
        console.error('DEBUG: 前端获取文档信息异常:', err)
      }
    }
    
    // 使用特定版本号获取文档版本
    const fetchDocumentVersionWithVersion = async (version) => {
      if (!props.documentId || !version) {
        console.log('DEBUG: fetchDocumentVersionWithVersion - documentId或version为空')
        return // 不显示错误，静默返回
      }
      
      isLoading.value = true
      error.value = null
      
      try {
        console.log('DEBUG: 前端开始获取文档版本（使用特定版本） - 文档ID:', props.documentId, '版本:', version)
        const response = await fetch(`/api/v1/documents/${props.documentId}/versions/${version}`)
        const result = await response.json()
        
        console.log('DEBUG: 前端获取文档版本响应（使用特定版本）:', result)
        
        if (result.code === 200) {
          documentVersion.value = result.data
          console.log('DEBUG: 前端获取文档版本成功（使用特定版本），数据:', documentVersion.value)
          console.log('DEBUG: 版本内容长度:', documentVersion.value.content ? documentVersion.value.content.length : '空')
          console.log('DEBUG: 版本状态:', documentVersion.value.status)
          console.log('DEBUG: 版本内容前100个字符:', documentVersion.value.content ? documentVersion.value.content.substring(0, 100) : '空')
          // 同时获取文档信息以确定文档类型
          fetchDocument()
        } else {
          console.error('DEBUG: 前端获取文档版本失败（使用特定版本）:', result.message)
          error.value = result.message || '获取文档版本失败'
        }
      } catch (err) {
        console.error('DEBUG: 前端获取文档版本异常（使用特定版本）:', err)
        error.value = '获取文档版本失败: ' + err.message
      } finally {
        isLoading.value = false
      }
    }
    
    // 下载版本
    const downloadVersion = () => {
      emit('download-version', documentVersion.value)
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
    
    // 组件挂载时获取文档版本详情
    onMounted(() => {
      fetchDocumentVersion()
    })
    
    // 使用watchEffect监听所有响应式数据
    watchEffect(() => {
      console.log('DEBUG: watchEffect触发 - props.documentId:', props.documentId, 'props.version:', props.version, 'props.globalState:', props.globalState)
      if (props.documentId && props.version) {
        console.log('DEBUG: watchEffect中检测到props都有值，重新获取数据')
        fetchDocumentVersion()
      } else if (props.documentId && !props.version) {
        console.log('DEBUG: 只有documentId没有version，尝试多种方式获取版本')
        
        // 首先尝试从全局状态获取
        if (props.globalState && props.globalState.isSet && props.globalState.version) {
          console.log('DEBUG: 从全局状态获取到版本:', props.globalState.version)
          fetchDocumentVersionWithVersion(props.globalState.version)
          return
        }
        
        // 尝试从URL获取版本号
        const urlParams = new URLSearchParams(window.location.search)
        const versionFromUrl = urlParams.get('version')
        if (versionFromUrl) {
          console.log('DEBUG: 从URL获取到版本:', versionFromUrl)
          // 直接使用获取到的版本号
          fetchDocumentVersionWithVersion(versionFromUrl)
        } else {
          // 如果URL中也没有，尝试从localStorage获取
          const versionFromStorage = getVersionFromStorage()
          if (versionFromStorage) {
            console.log('DEBUG: 从localStorage获取到版本:', versionFromStorage)
            fetchDocumentVersionWithVersion(versionFromStorage)
          }
        }
      }
    })
    
    // 监听props.version变化，当version从undefined变为有效值时重新获取数据
    watch(() => props.version, (newVersion, oldVersion) => {
      console.log('DEBUG: props.version变化 - 新版本:', newVersion, '旧版本:', oldVersion)
      if (newVersion) {
        console.log('DEBUG: props.version有值，重新获取数据')
        fetchDocumentVersion()
      }
    })
    
    // 监听props.documentId变化
    watch(() => props.documentId, (newDocumentId, oldDocumentId) => {
      console.log('DEBUG: props.documentId变化 - 新ID:', newDocumentId, '旧ID:', oldDocumentId)
      if (newDocumentId && props.version) {
        console.log('DEBUG: props.documentId有值且props.version也有值，重新获取数据')
        fetchDocumentVersion()
      }
    })
    
    return {
      documentVersion,
      documentType,
      isLoading,
      error,
      downloadVersion,
      getStatusBadgeClass,
      getStatusText,
      formatFileSize,
      formatDate
    }
  }
}
</script>

<style scoped>
.document-content pre {
  white-space: pre-wrap;
  word-wrap: break-word;
  max-height: 500px;
  overflow-y: auto;
}
</style>