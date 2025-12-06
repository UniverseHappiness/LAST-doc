<template>
  <div>
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h2>文档查看 - {{ document?.name }}</h2>
      <button class="btn btn-outline-secondary" @click="$emit('back')">返回</button>
    </div>
    
    <div v-if="document" class="card">
      <div class="card-header d-flex justify-content-between align-items-center">
        <div>
          <h5 class="mb-0">{{ document.name }}</h5>
          <p class="text-muted mb-0">{{ document.description || '暂无描述' }}</p>
        </div>
        <div>
          <span class="badge bg-secondary me-2">{{ document.type }}</span>
          <span class="badge bg-info me-2">{{ document.version }}</span>
          <span class="badge bg-primary me-2">{{ document.library }}</span>
          <span class="badge" :class="getStatusBadgeClass(document.status)">{{ getStatusText(document.status) }}</span>
        </div>
      </div>
      <div class="card-body">
        <!-- 文档操作按钮 -->
        <div class="mb-3">
          <button class="btn btn-sm btn-primary me-2" @click="downloadDocument">
            <i class="bi bi-download me-1"></i>下载
          </button>
          <button class="btn btn-sm btn-outline-primary me-2" @click="editDocument">
            <i class="bi bi-pencil me-1"></i>编辑
          </button>
        </div>
        
        <!-- 文档内容预览 -->
        <div v-if="document.content" class="document-content">
          <h6 class="mb-3">文档内容</h6>
          <div class="border rounded p-3 bg-light">
            <div v-if="document.type === 'markdown' || document.type === 'java_doc'" class="mb-0">
              <pre class="markdown-content">{{ document.content }}</pre>
            </div>
            <div v-else-if="document.type === 'swagger' || document.type === 'openapi'" class="mb-0">
              <pre class="json-content">{{ formatJsonContent(document.content) }}</pre>
            </div>
            <div v-else-if="document.type === 'pdf'" class="mb-0">
              <div class="pdf-preview">
                <p class="text-muted">PDF 文档预览</p>
                <div v-if="document.content && document.content.trim() !== ''">
                  <pre class="pdf-content">{{ document.content }}</pre>
                  <p class="text-info small mt-2">PDF 文档已提取文本内容，如需查看完整格式请下载文件</p>
                </div>
                <div v-else>
                  <p>PDF 文档内容无法直接预览，请点击下载按钮查看完整文档</p>
                </div>
              </div>
            </div>
            <div v-else-if="document.type === 'docx'" class="mb-0">
              <div class="docx-preview">
                <p class="text-muted">DOCX 文档预览</p>
                <div v-if="document.content && document.content.trim() !== ''">
                  <pre class="docx-content">{{ document.content }}</pre>
                  <p class="text-info small mt-2">DOCX 文档已提取文本内容，如需查看完整格式请下载文件</p>
                </div>
                <div v-else>
                  <p>DOCX 文档内容无法直接预览，请点击下载按钮查看完整文档</p>
                </div>
              </div>
            </div>
            <div v-else class="mb-0">
              <p>文档内容类型 {{ document.type }} 暂不支持预览</p>
              <p class="text-muted">请点击下载按钮查看完整文档</p>
            </div>
          </div>
        </div>
        <div v-else class="document-content">
          <h6 class="mb-3">文档内容</h6>
          <div class="border rounded p-3 bg-light">
            <p class="text-muted">暂无内容 - 状态: {{ getStatusText(document.status) }}</p>
            <p class="text-muted small">如果状态为"处理中"，请稍后刷新页面</p>
          </div>
        </div>
        
        <!-- 文档元数据 -->
        <div v-if="metadata" class="document-metadata mt-4">
          <h6 class="mb-3">文档元数据</h6>
          <div class="border rounded p-3 bg-light">
            <pre>{{ JSON.stringify(metadata, null, 2) }}</pre>
          </div>
        </div>
        
        <!-- 文档信息 -->
        <div class="document-info mt-4">
          <h6 class="mb-3">文档信息</h6>
          <table class="table table-sm">
            <tbody>
              <tr>
                <th scope="row" style="width: 120px;">文档ID</th>
                <td>{{ document.id }}</td>
              </tr>
              <tr>
                <th scope="row">文件路径</th>
                <td>{{ document.file_path }}</td>
              </tr>
              <tr>
                <th scope="row">文件大小</th>
                <td>{{ formatFileSize(document.file_size) }}</td>
              </tr>
              <tr>
                <th scope="row">标签</th>
                <td>
                  <span v-for="tag in document.tags" :key="tag" class="badge bg-light text-dark me-1">{{ tag }}</span>
                  <span v-if="!document.tags || document.tags.length === 0" class="text-muted">无标签</span>
                </td>
              </tr>
              <tr>
                <th scope="row">创建时间</th>
                <td>{{ formatDate(document.created_at) }}</td>
              </tr>
              <tr>
                <th scope="row">更新时间</th>
                <td>{{ formatDate(document.updated_at) }}</td>
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
      <p class="mt-2">正在加载文档内容...</p>
    </div>
    
    <!-- 错误状态 -->
    <div v-if="error" class="alert alert-danger" role="alert">
      {{ error }}
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'

export default {
  name: 'DocumentView',
  props: {
    documentId: {
      type: String,
      required: true
    }
  },
  emits: ['back', 'edit-document', 'download-document'],
  setup(props, { emit }) {
    const document = ref(null)
    const metadata = ref(null)
    const isLoading = ref(false)
    const error = ref(null)
    
    // 获取文档详情
    const fetchDocument = async () => {
      isLoading.value = true
      error.value = null
      
      try {
        const response = await fetch(`/api/v1/documents/${props.documentId}`)
        const result = await response.json()
        
        if (result.code === 200) {
          document.value = result.data
          // 获取文档元数据
          fetchDocumentMetadata()
        } else {
          error.value = result.message || '获取文档失败'
        }
      } catch (err) {
        console.error('获取文档失败:', err)
        error.value = '获取文档失败: ' + err.message
      } finally {
        isLoading.value = false
      }
    }
    
    // 获取文档元数据
    const fetchDocumentMetadata = async () => {
      try {
        const response = await fetch(`/api/v1/documents/${props.documentId}/metadata`)
        const result = await response.json()
        
        if (result.code === 200) {
          metadata.value = result.data
        }
      } catch (err) {
        console.error('获取文档元数据失败:', err)
      }
    }
    
    // 下载文档
    const downloadDocument = () => {
      emit('download-document', document.value)
    }
    
    // 编辑文档
    const editDocument = () => {
      emit('edit-document', document.value)
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
    
    // 格式化JSON内容
    const formatJsonContent = (content) => {
      try {
        return JSON.stringify(JSON.parse(content), null, 2)
      } catch (e) {
        console.error('JSON格式化失败:', e)
        return content
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
    
    // 组件挂载时获取文档详情
    onMounted(() => {
      fetchDocument()
    })
    
    return {
      document,
      metadata,
      isLoading,
      error,
      downloadDocument,
      editDocument,
      getStatusBadgeClass,
      getStatusText,
      formatJsonContent,
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

.document-content pre.pdf-content {
  background-color: #f8f9fa;
  border: 1px solid #dee2e6;
  padding: 15px;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.5;
}

.document-metadata pre {
  max-height: 300px;
  overflow-y: auto;
}

.pdf-preview {
  border: 1px solid #dee2e6;
  border-radius: 4px;
  padding: 15px;
  background-color: #f8f9fa;
}
</style>