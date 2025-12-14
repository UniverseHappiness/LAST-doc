<template>
  <div class="search-view">
    <h2 class="mb-4">文档检索</h2>
    
    <!-- 搜索表单 -->
    <div class="card mb-4">
      <div class="card-body">
        <form @submit.prevent="performSearch">
          <div class="row g-3">
            <div class="col-md-8">
              <label for="searchQuery" class="form-label">搜索关键词</label>
              <input 
                type="text" 
                class="form-control" 
                id="searchQuery" 
                v-model="searchForm.query" 
                placeholder="输入要搜索的关键词或句子"
                required
              >
            </div>
            <div class="col-md-4">
              <label for="searchType" class="form-label">搜索类型</label>
              <select 
                class="form-select" 
                id="searchType" 
                v-model="searchForm.searchType"
              >
                <option value="keyword">关键词搜索</option>
                <option value="semantic">语义搜索</option>
                <option value="hybrid">混合搜索</option>
              </select>
            </div>
          </div>
          
          <!-- 高级搜索选项 -->
          <div class="mt-3">
            <button 
              type="button" 
              class="btn btn-sm btn-outline-secondary" 
              @click="showAdvancedOptions = !showAdvancedOptions"
            >
              {{ showAdvancedOptions ? '收起高级选项' : '展开高级选项' }}
            </button>
          </div>
          
          <div v-if="showAdvancedOptions" class="mt-3">
            <div class="row g-3">
              <div class="col-md-4">
                <label for="documentType" class="form-label">文档类型</label>
                <select
                  class="form-select"
                  id="documentType"
                  v-model="searchForm.filters.type"
                >
                  <option value="">全部类型</option>
                  <option value="markdown">Markdown</option>
                  <option value="pdf">PDF</option>
                  <option value="docx">DOCX</option>
                  <option value="swagger">Swagger</option>
                  <option value="openapi">OpenAPI</option>
                  <option value="java_doc">JavaDoc</option>
                </select>
              </div>
              <div class="col-md-4">
                <label for="documentCategory" class="form-label">文档分类</label>
                <select
                  class="form-select"
                  id="documentCategory"
                  v-model="searchForm.filters.category"
                >
                  <option value="">全部分类</option>
                  <option value="code">代码</option>
                  <option value="document">文档</option>
                </select>
              </div>
              <div class="col-md-4">
                <label for="library" class="form-label">所属库</label>
                <input
                  type="text"
                  class="form-control"
                  id="library"
                  v-model="searchForm.filters.library"
                  placeholder="指定库名称"
                >
              </div>
            </div>
            <div class="row g-3 mt-2">
              <div class="col-md-4">
                <label for="contentType" class="form-label">内容类型</label>
                <select
                  class="form-select"
                  id="contentType"
                  v-model="searchForm.filters.content_type"
                >
                  <option value="">全部类型</option>
                  <option value="text">文本内容</option>
                  <option value="code">代码内容</option>
                </select>
              </div>
              <div class="col-md-4">
                <label for="version" class="form-label">版本</label>
                <input
                  type="text"
                  class="form-control"
                  id="version"
                  v-model="searchForm.filters.version"
                  placeholder="指定版本号"
                >
              </div>
            </div>
          </div>
          
          <div class="mt-3">
            <button type="submit" class="btn btn-primary" :disabled="isSearching">
              <span v-if="isSearching" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
              搜索
            </button>
            <button type="button" class="btn btn-outline-secondary ms-2" @click="resetSearch">重置</button>
          </div>
        </form>
      </div>
    </div>
    
    <!-- 搜索结果 -->
    <div v-if="searchResults.length > 0" class="card">
      <div class="card-header d-flex justify-content-between align-items-center">
        <span>搜索结果 (共 {{ searchTotal }} 条)</span>
        <div>
          <select class="form-select form-select-sm" v-model="pageSize" @change="changePageSize">
            <option :value="10">每页 10 条</option>
            <option :value="20">每页 20 条</option>
            <option :value="50">每页 50 条</option>
          </select>
        </div>
      </div>
      <div class="card-body">
        <div class="list-group">
          <div v-for="result in searchResults" :key="result.id" class="list-group-item">
            <div class="d-flex justify-content-between">
              <div>
                <h5 class="mb-1">{{ result.metadata?.document_name || '无标题' }}</h5>
                <p class="mb-1 text-muted">{{ result.snippet }}</p>
                <div>
                  <span class="badge bg-secondary me-2">{{ result.content_type }}</span>
                  <span v-if="result.metadata?.document_type" class="badge bg-info me-2">类型: {{ getDocumentTypeText(result.metadata.document_type) }}</span>
                  <span v-if="result.metadata?.document_category" class="badge me-2" :class="getCategoryBadgeClass(result.metadata.document_category)">{{ getCategoryText(result.metadata.document_category) }}</span>
                  <span v-if="result.metadata?.document_library" class="badge bg-primary me-2">库: {{ result.metadata.document_library }}</span>
                  <span class="badge bg-light text-dark me-2">版本: {{ result.version }}</span>
                  <span class="badge bg-success">相关度: {{ (result.score * 100).toFixed(1) }}%</span>
                </div>
              </div>
              <div>
                <button class="btn btn-sm btn-outline-primary" @click="viewDocument(result)">查看文档</button>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 分页 -->
        <nav v-if="totalPages > 1" class="mt-4">
          <ul class="pagination justify-content-center">
            <li class="page-item" :class="{ disabled: currentPage === 1 }">
              <a class="page-link" href="#" @click.prevent="changePage(currentPage - 1)">上一页</a>
            </li>
            <li 
              v-for="page in paginationPages" 
              :key="page" 
              class="page-item" 
              :class="{ active: currentPage === page }"
            >
              <a class="page-link" href="#" @click.prevent="changePage(page)">{{ page }}</a>
            </li>
            <li class="page-item" :class="{ disabled: currentPage === totalPages }">
              <a class="page-link" href="#" @click.prevent="changePage(currentPage + 1)">下一页</a>
            </li>
          </ul>
        </nav>
      </div>
    </div>
    
    <!-- 无搜索结果 -->
    <div v-else-if="hasSearched && searchResults.length === 0" class="alert alert-info">
      未找到与搜索关键词相关的文档，请尝试使用其他关键词或调整搜索类型。
    </div>
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted } from 'vue'
import { useSearchService } from '../utils/searchService'

export default {
  name: 'SearchView',
  setup() {
    const searchForm = reactive({
      query: '',
      searchType: 'keyword',
      filters: {
        version: '',
        content_type: '',
        type: '',
        category: '',
        library: ''
      }
    })
    
    const showAdvancedOptions = ref(false)
    const isSearching = ref(false)
    const hasSearched = ref(false)
    const searchResults = ref([])
    const searchTotal = ref(0)
    const currentPage = ref(1)
    const pageSize = ref(10)
    
    const { searchDocuments } = useSearchService()
    
    // 计算总页数
    const totalPages = computed(() => {
      console.log('计算总页数 - searchTotal:', searchTotal.value, 'pageSize:', pageSize.value)
      if (!searchTotal.value || searchTotal.value <= 0) {
        return 0
      }
      return Math.ceil(searchTotal.value / pageSize.value)
    })
    
    // 获取分页页码
    const paginationPages = computed(() => {
      const pages = []
      const maxVisiblePages = 5
      let startPage = Math.max(1, currentPage.value - Math.floor(maxVisiblePages / 2))
      let endPage = Math.min(totalPages.value, startPage + maxVisiblePages - 1)
      
      if (endPage - startPage + 1 < maxVisiblePages) {
        startPage = Math.max(1, endPage - maxVisiblePages + 1)
      }
      
      for (let i = startPage; i <= endPage; i++) {
        pages.push(i)
      }
      
      return pages
    })
    
    // 执行搜索
    const performSearch = async () => {
      if (!searchForm.query.trim()) {
        alert('请输入搜索关键词')
        return
      }
      
      isSearching.value = true
      hasSearched.value = true
      currentPage.value = 1
      
      try {
        const request = {
          query: searchForm.query,
          searchType: searchForm.searchType,
          filters: { ...searchForm.filters },
          page: currentPage.value,
          size: pageSize.value
        }
        
        const response = await searchDocuments(request)
        console.log('搜索响应数据:', response)
        console.log('response.data:', response.data)
        console.log('response.data.items:', response.data?.items)
        console.log('response.data.total:', response.data?.total)
        
        // 添加数据验证
        if (!response.data) {
          console.error('错误：响应数据为空')
          alert('搜索失败：服务器返回空数据')
          return
        }
        
        if (!response.data.items) {
          console.error('错误：搜索结果项为空')
          searchResults.value = []
          searchTotal.value = 0
          return
        }
        
        searchResults.value = response.data.items
        searchTotal.value = response.data.total || 0
      } catch (error) {
        console.error('搜索失败:', error)
        alert('搜索失败: ' + (error.response?.data?.message || error.message))
      } finally {
        isSearching.value = false
      }
    }
    
    // 重置搜索
    const resetSearch = () => {
      searchForm.query = ''
      searchForm.searchType = 'keyword'
      searchForm.filters = {
        version: '',
        content_type: '',
        type: '',
        category: '',
        library: ''
      }
      searchResults.value = []
      searchTotal.value = 0
      currentPage.value = 1
      hasSearched.value = false
    }
    
    // 切换页码
    const changePage = async (page) => {
      if (page < 1 || page > totalPages.value || page === currentPage.value) {
        return
      }
      
      currentPage.value = page
      isSearching.value = true
      
      try {
        const request = {
          query: searchForm.query,
          searchType: searchForm.searchType,
          filters: { ...searchForm.filters },
          page: currentPage.value,
          size: pageSize.value
        }
        
        const response = await searchDocuments(request)
        console.log('分页搜索响应数据:', response)
        console.log('response.data:', response.data)
        console.log('response.data.items:', response.data?.items)
        console.log('response.data.total:', response.data?.total)
        
        // 添加数据验证
        if (!response.data) {
          console.error('错误：分页响应数据为空')
          alert('搜索失败：服务器返回空数据')
          return
        }
        
        if (!response.data.items) {
          console.error('错误：分页搜索结果项为空')
          searchResults.value = []
          searchTotal.value = 0
          return
        }
        
        searchResults.value = response.data.items
        searchTotal.value = response.data.total || 0
      } catch (error) {
        console.error('分页搜索失败:', error)
        alert('搜索失败: ' + (error.response?.data?.message || error.message))
      } finally {
        isSearching.value = false
      }
    }
    
    // 改变每页大小
    const changePageSize = async () => {
      currentPage.value = 1
      await changePage(1)
    }
    
    // 获取文档类型文本
    const getDocumentTypeText = (type) => {
      switch (type) {
        case 'markdown': return 'Markdown'
        case 'pdf': return 'PDF'
        case 'docx': return 'DOCX'
        case 'swagger': return 'Swagger'
        case 'openapi': return 'OpenAPI'
        case 'java_doc': return 'JavaDoc'
        default: return type
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
    
    // 查看文档
    const viewDocument = (result) => {
      // 这里可以实现跳转到文档详情页面的逻辑
      alert(`跳转到文档: ${result.document_id} (版本: ${result.version})`)
    }
    
    return {
      searchForm,
      showAdvancedOptions,
      isSearching,
      hasSearched,
      searchResults,
      searchTotal,
      currentPage,
      pageSize,
      totalPages,
      paginationPages,
      performSearch,
      resetSearch,
      changePage,
      changePageSize,
      viewDocument,
      getDocumentTypeText,
      getCategoryBadgeClass,
      getCategoryText
    }
  }
}
</script>

<style scoped>
.search-view {
  max-width: 1200px;
  margin: 0 auto;
}

.card {
  box-shadow: 0 0.125rem 0.25rem rgba(0, 0, 0, 0.075);
}

.list-group-item {
  border-left: none;
  border-right: none;
}

.list-group-item:first-child {
  border-top: none;
}

.list-group-item:last-child {
  border-bottom: none;
}

.badge {
  font-size: 0.75rem;
}
</style>