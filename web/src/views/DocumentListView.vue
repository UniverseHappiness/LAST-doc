<template>
  <div>
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h2>文档列表</h2>
      <div>
        <input type="text" 
               class="form-control form-control-sm d-inline-block w-auto me-2" 
               v-model="localSearchQuery" 
               placeholder="搜索文档...">
        <button class="btn btn-sm btn-outline-secondary" @click="handleSearch">
          <i class="bi bi-search"></i>
        </button>
      </div>
    </div>
    
    <!-- 文档列表 -->
    <div v-if="documents.length > 0">
      <div v-for="doc in documents" :key="doc.id" class="document-item">
        <div class="d-flex justify-content-between align-items-start">
          <div>
            <h5>{{ doc.name }}</h5>
            <p class="text-muted mb-1">{{ doc.description || '暂无描述' }}</p>
            <div class="mb-2">
              <span class="badge bg-secondary document-badge">{{ doc.type }}</span>
              <span class="badge" :class="getCategoryBadgeClass(doc.category)">{{ getCategoryText(doc.category) }}</span>
              <span class="badge bg-info document-badge">当前版本: {{ doc.version }}</span>
              <span class="badge bg-primary document-badge">{{ doc.library }}</span>
              <span class="badge" :class="getStatusBadgeClass(doc.status)">{{ getStatusText(doc.status) }}</span>
              <span class="badge bg-warning document-badge" v-if="doc.version_count !== undefined">版本总数: {{ doc.version_count }}</span>
            </div>
            <div v-if="doc.tags && doc.tags.length > 0">
              <span v-for="tag in doc.tags" :key="tag" class="badge bg-light text-dark me-1">{{ tag }}</span>
            </div>
          </div>
          <div>
            <div class="btn-group" role="group">
              <button class="btn btn-sm btn-outline-primary" @click="$emit('view-document', doc)">查看</button>
              <button class="btn btn-sm btn-outline-secondary" @click="$emit('view-versions', doc)">版本管理</button>
              <button class="btn btn-sm btn-outline-danger" @click="$emit('delete-document', doc)">删除</button>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 空状态 -->
    <div v-else class="text-center py-5">
      <i class="bi bi-file-earmark-text" style="font-size: 3rem;"></i>
      <p class="mt-3">暂无文档</p>
      <button class="btn btn-primary" @click="$emit('navigate-to-upload')">上传文档</button>
    </div>
    
    <!-- 分页 -->
    <nav v-if="pagination.total > pagination.size" class="mt-4">
      <ul class="pagination justify-content-center">
        <li class="page-item" :class="{ disabled: pagination.page <= 1 }">
          <a class="page-link" href="#" @click.prevent="$emit('change-page', pagination.page - 1)">上一页</a>
        </li>
        <li v-for="page in getPaginationPages()" :key="page" class="page-item" :class="{ active: page === pagination.page }">
          <a class="page-link" href="#" @click.prevent="$emit('change-page', page)">{{ page }}</a>
        </li>
        <li class="page-item" :class="{ disabled: pagination.page * pagination.size >= pagination.total }">
          <a class="page-link" href="#" @click.prevent="$emit('change-page', pagination.page + 1)">下一页</a>
        </li>
      </ul>
    </nav>
  </div>
</template>

<script>
import { ref, watch } from 'vue'

export default {
  name: 'DocumentListView',
  props: {
    documents: {
      type: Array,
      required: true
    },
    pagination: {
      type: Object,
      required: true
    },
    searchQuery: {
      type: String,
      default: ''
    }
  },
  emits: ['search', 'change-page', 'view-document', 'view-versions', 'delete-document', 'navigate-to-upload'],
  setup(props, { emit }) {
    const localSearchQuery = ref(props.searchQuery)
    
    // 监听搜索查询变化
    watch(localSearchQuery, (newValue) => {
      if (!newValue.trim()) {
        emit('search')
      }
    })
    
    // 处理搜索
    const handleSearch = () => {
      emit('search', localSearchQuery.value)
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
    
    // 获取分类徽章样式
    const getCategoryBadgeClass = (category) => {
      switch (category) {
        case 'code': return 'bg-success'
        case 'document': return 'bg-info'
        default: return 'bg-light text-dark'
      }
    }
    
    // 获取分类文本
    const getCategoryText = (category) => {
      switch (category) {
        case 'code': return '代码'
        case 'document': return '文档'
        default: return category
      }
    }
    
    // 获取分页页码
    const getPaginationPages = () => {
      const totalPages = Math.ceil(props.pagination.total / props.pagination.size)
      const pages = []
      const maxVisiblePages = 5
      
      let startPage = Math.max(1, props.pagination.page - Math.floor(maxVisiblePages / 2))
      let endPage = Math.min(totalPages, startPage + maxVisiblePages - 1)
      
      if (endPage - startPage + 1 < maxVisiblePages) {
        startPage = Math.max(1, endPage - maxVisiblePages + 1)
      }
      
      for (let i = startPage; i <= endPage; i++) {
        pages.push(i)
      }
      
      return pages
    }
    
    return {
      localSearchQuery,
      handleSearch,
      getStatusBadgeClass,
      getStatusText,
      getCategoryBadgeClass,
      getCategoryText,
      getPaginationPages
    }
  }
}
</script>