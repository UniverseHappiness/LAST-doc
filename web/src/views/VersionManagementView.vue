<template>
  <div>
    <div class="d-flex justify-content-between align-items-center mb-4">
      <div>
        <h2>版本管理 - {{ selectedDocument.name }}</h2>
        <p class="text-muted">以下列出了该文档的所有版本，您可以通过版本号区分不同的版本</p>
      </div>
      <button class="btn btn-outline-secondary" @click="$emit('back')">返回</button>
    </div>
    
    <div v-if="documentVersions.length > 0">
      <div class="alert alert-info">
        <i class="bi bi-info-circle me-2"></i>
        该文档共有 <strong>{{ documentVersions.length }}</strong> 个版本，每个版本代表文档在不同时间点的更新或修改
      </div>
      
      <div v-for="version in documentVersions" :key="version.id" class="document-item">
        <div class="d-flex justify-content-between align-items-center">
          <div>
            <h5>版本 {{ version.version }}</h5>
            <p class="text-muted mb-1">{{ version.description || '暂无描述' }}</p>
            <div class="mb-2">
              <span class="badge" :class="getStatusBadgeClass(version.status)">{{ getStatusText(version.status) }}</span>
              <span class="badge bg-light text-dark ms-1">{{ formatFileSize(version.file_size) }}</span>
              <span class="badge bg-secondary ms-1" v-if="version.version === selectedDocument.version">当前版本</span>
            </div>
            <p class="text-muted small">创建时间: {{ formatDate(version.created_at) }}</p>
          </div>
          <div>
            <button class="btn btn-sm btn-outline-primary me-2" @click="$emit('view-version', version)">查看</button>
            <button class="btn btn-sm btn-outline-danger" @click="$emit('delete-version', version)">删除</button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 空状态 -->
    <div v-else class="text-center py-5">
      <i class="bi bi-git-branch" style="font-size: 3rem;"></i>
      <p class="mt-3">暂无版本记录</p>
    </div>
  </div>
</template>

<script>
export default {
  name: 'VersionManagementView',
  props: {
    selectedDocument: {
      type: Object,
      required: true
    },
    documentVersions: {
      type: Array,
      required: true
    }
  },
  emits: ['view-version', 'delete-version', 'back'],
  setup(props, { emit }) {
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
    
    return {
      getStatusBadgeClass,
      getStatusText,
      formatFileSize,
      formatDate
    }
  }
}
</script>