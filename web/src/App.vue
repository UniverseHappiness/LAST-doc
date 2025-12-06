<template>
  <div id="app" class="container-fluid">
    <div class="row">
      <!-- 侧边栏 -->
      <div class="col-md-3 sidebar p-4">
        <h4 class="mb-4">AI技术文档库</h4>
        <ul class="nav flex-column">
          <li class="nav-item mb-2">
            <a class="nav-link" 
               href="#" 
               :class="{ active: currentView === 'list' }"
               @click="currentView = 'list'">
              <i class="bi bi-file-earmark-text me-2"></i>文档列表
            </a>
          </li>
          <li class="nav-item mb-2">
            <a class="nav-link" 
               href="#" 
               :class="{ active: currentView === 'upload' }"
               @click="currentView = 'upload'">
              <i class="bi bi-cloud-upload me-2"></i>上传文档
            </a>
          </li>
        </ul>
        
        <!-- 过滤器 -->
        <div class="mt-4" v-if="currentView === 'list'">
          <h6>过滤器</h6>
          <div class="mb-3">
            <label for="filterLibrary" class="form-label">所属库</label>
            <input type="text" 
                   class="form-control form-control-sm" 
                   id="filterLibrary" 
                   v-model="filters.library" 
                   placeholder="输入库名称">
          </div>
          <div class="mb-3">
            <label for="filterType" class="form-label">文档类型</label>
            <select class="form-select form-select-sm" 
                    id="filterType" 
                    v-model="filters.type">
              <option value="">全部类型</option>
              <option value="markdown">Markdown</option>
              <option value="pdf">PDF</option>
              <option value="docx">DOCX</option>
              <option value="swagger">Swagger</option>
              <option value="openapi">OpenAPI</option>
              <option value="java_doc">JavaDoc</option>
            </select>
          </div>
          <div class="mb-3">
            <label for="filterVersion" class="form-label">版本</label>
            <input type="text" 
                   class="form-control form-control-sm" 
                   id="filterVersion" 
                   v-model="filters.version" 
                   placeholder="输入版本号">
          </div>
          <div class="mb-3">
            <label for="filterStatus" class="form-label">状态</label>
            <select class="form-select form-select-sm" 
                    id="filterStatus" 
                    v-model="filters.status">
              <option value="">全部状态</option>
              <option value="uploading">上传中</option>
              <option value="processing">处理中</option>
              <option value="completed">已完成</option>
              <option value="failed">失败</option>
            </select>
          </div>
          <button class="btn btn-sm btn-primary w-100" @click="applyFilters">应用过滤</button>
        </div>
      </div>
      
      <!-- 主内容区 -->
      <div class="col-md-9 main-content">
        <!-- 文档列表视图 -->
        <DocumentListView
          v-if="currentView === 'list'"
          :documents="documents"
          :pagination="pagination"
          :search-query="searchQuery"
          @search="searchDocuments"
          @change-page="changePage"
          @view-document="viewDocument"
          @view-versions="viewDocumentVersions"
          @delete-document="deleteDocument"
        />
        
        <!-- 文档查看视图 -->
        <DocumentView
          v-if="currentView === 'view'"
          :document-id="selectedDocumentId"
          @back="currentView = 'list'"
          @edit-document="editDocument"
          @download-document="downloadDocument"
        />
        
        <!-- 文档编辑视图 -->
        <DocumentEditView
          v-if="currentView === 'edit'"
          :document-id="selectedDocumentId"
          @back="currentView = 'list'"
          @document-updated="handleDocumentUpdated"
        />
        
        <!-- 上传文档视图 -->
        <UploadView 
          v-if="currentView === 'upload'"
          :is-uploading="isUploading"
          :upload-form="uploadForm"
          @upload="uploadDocument"
          @file-selected="fileSelected"
          @drag-over="dragOver"
          @drag-leave="dragLeave"
          @drop-file="dropFile"
          @back="currentView = 'list'"
        />
        
        <!-- 版本管理视图 -->
        <VersionManagementView
          v-if="currentView === 'versions' && selectedDocument"
          :selected-document="selectedDocument"
          :document-versions="documentVersions"
          @view-version="viewDocumentVersion"
          @delete-version="deleteDocumentVersion"
          @back="currentView = 'list'"
        />
        
        <!-- 文档版本查看视图 -->
        <DocumentVersionView
          v-if="currentView === 'version-view' && selectedDocumentId"
          :key="`version-${selectedDocumentId}-${selectedVersion || 'none'}`"
          :document-id="selectedDocumentId"
          :version="selectedVersion"
          :global-state="globalVersionState"
          @back="currentView = 'versions'"
          @download-version="downloadVersion"
        />
        <!-- 调试信息 -->
        <div v-if="currentView === 'version-view'" class="alert alert-info mt-3">
          <strong>调试信息:</strong><br>
          selectedDocumentId: {{ selectedDocumentId }}<br>
          selectedVersion: {{ selectedVersion }}<br>
          selectedVersionComputed: {{ selectedVersionComputed }}
        </div>
      </div>
    </div>
    
    <!-- 加载遮罩 -->
    <div class="overlay" :style="{ display: isLoading ? 'block' : 'none' }"></div>
    <div class="loading-spinner" :style="{ display: isLoading ? 'block' : 'none' }">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
    </div>
    
    <!-- 通知组件 -->
    <div v-if="notification.show"
         class="notification"
         :class="`notification-${notification.type}`"
         @click="hideNotification">
      <div class="notification-content">
        <span class="notification-message">{{ notification.message }}</span>
        <button class="notification-close" @click.stop="hideNotification">×</button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, onMounted, watch, computed, nextTick } from 'vue'
import DocumentListView from './views/DocumentListView.vue'
import DocumentView from './views/DocumentView.vue'
import DocumentEditView from './views/DocumentEditView.vue'
import DocumentVersionView from './views/DocumentVersionView.vue'
import UploadView from './views/UploadView.vue'
import VersionManagementView from './views/VersionManagementView.vue'
import { useDocumentService } from './utils/documentService'

export default {
  name: 'App',
  components: {
    DocumentListView,
    DocumentView,
    DocumentEditView,
    DocumentVersionView,
    UploadView,
    VersionManagementView
  },
  setup() {
    const currentView = ref('list')
    const documents = ref([])
    const documentVersions = ref([])
    const selectedDocument = ref(null)
    const selectedDocumentId = ref(null)
    const selectedVersion = ref(null)
    const isLoading = ref(false)
    const isUploading = ref(false)
    const searchQuery = ref('')
    const notification = ref({
      show: false,
      message: '',
      type: 'info' // info, success, warning, error
    })
    
    // 全局版本状态
    const globalVersionState = reactive({
      documentId: null,
      version: null,
      isSet: false
    })
    
    // 过滤器
    const filters = reactive({
      library: '',
      type: '',
      version: '',
      status: ''
    })
    
    // 分页
    const pagination = reactive({
      page: 1,
      size: 10,
      total: 0
    })
    
    // 上传表单
    const uploadForm = reactive({
      name: '',
      type: '',
      version: '',
      library: '',
      description: '',
      tags: '',
      file: null
    })
    
    // 使用文档服务
    const {
      fetchDocuments,
      fetchDocumentVersions,
      uploadDocument,
      updateDocument,
      deleteDocument,
      deleteDocumentVersion,
      searchDocuments,
      applyFilters,
      changePage,
      getPaginationPages
    } = useDocumentService({
      documents,
      documentVersions,
      pagination,
      filters,
      searchQuery,
      isLoading,
      isUploading,
      uploadForm,
      currentView
    })
    
    // 查看文档
    const viewDocument = (doc) => {
      console.log('DEBUG: 尝试查看文档', doc)
      selectedDocumentId.value = doc.id
      currentView.value = 'view'
    }
    
    // 查看文档版本
    const viewDocumentVersions = async (doc) => {
      selectedDocument.value = doc
      currentView.value = 'versions'
      await fetchDocumentVersions(doc.id)
    }
    
    // 查看特定版本的文档
    const viewDocumentVersion = (version) => {
      console.log('DEBUG: 尝试查看文档版本', version)
      if (!version || !version.document_id || !version.version) {
        console.error('版本信息不完整', version)
        alert('版本信息不完整，无法查看')
        return
      }
      selectedDocumentId.value = version.document_id
      selectedVersion.value = version.version
      console.log('DEBUG: 设置selectedVersion值为:', selectedVersion.value)
      
      // 将版本号保存到localStorage中
      try {
        localStorage.setItem('selectedVersion', version.version)
        console.log('DEBUG: 版本号已保存到localStorage:', version.version)
      } catch (e) {
        console.error('DEBUG: 保存版本号到localStorage失败:', e)
      }
      
      // 设置全局状态
      globalVersionState.documentId = version.document_id
      globalVersionState.version = version.version
      globalVersionState.isSet = true
      console.log('DEBUG: 设置全局状态:', globalVersionState)
      
      currentView.value = 'version-view'
      console.log('DEBUG: 切换到version-view视图')
      
      // 使用nextTick确保DOM更新
      nextTick(() => {
        console.log('DEBUG: nextTick回调 - selectedVersion.value:', selectedVersion.value)
      })
    }
    
    // 文件选择
    const fileSelected = (event) => {
      const file = event.target.files[0]
      if (file) {
        uploadForm.file = file
      }
    }
    
    // 下载文档
    const downloadDocument = async (doc) => {
      console.log('DEBUG: 尝试下载文档', doc)
      try {
        const response = await fetch(`/api/v1/documents/${doc.id}/download`)
        if (!response.ok) {
          throw new Error('下载失败')
        }
        
        const blob = await response.blob()
        const url = window.URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.style.display = 'none'
        a.href = url
        a.download = doc.name
        document.body.appendChild(a)
        a.click()
        window.URL.revokeObjectURL(url)
        document.body.removeChild(a)
      } catch (error) {
        console.error('下载文档失败:', error)
        alert('下载文档失败: ' + error.message)
      }
    }
    
    // 下载文档版本
    const downloadVersion = async (version) => {
      console.log('DEBUG: 尝试下载文档版本', version)
      try {
        const response = await fetch(`/api/v1/documents/${version.document_id}/versions/${version.version}/download`)
        if (!response.ok) {
          throw new Error('下载失败')
        }
        
        const blob = await response.blob()
        const url = window.URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.style.display = 'none'
        a.href = url
        a.download = `v${version.version}_${version.document_id}`
        document.body.appendChild(a)
        a.click()
        window.URL.revokeObjectURL(url)
        document.body.removeChild(a)
      } catch (error) {
        console.error('下载文档版本失败:', error)
        alert('下载文档版本失败: ' + error.message)
      }
    }
    
    // 编辑文档
    const editDocument = (doc) => {
      console.log('DEBUG: 尝试编辑文档', doc)
      selectedDocumentId.value = doc.id
      currentView.value = 'edit'
    }
    
    // 文档更新完成后的处理
    const handleDocumentUpdated = async () => {
      console.log('DEBUG: 文档更新完成，刷新文档列表')
      await fetchDocuments()
      // 如果当前有选中的文档，也需要刷新它
      if (selectedDocumentId.value) {
        // 这里可以添加获取单个文档的逻辑
      }
    }
    
    // 显示通知
    const showNotification = (message, type = 'info') => {
      notification.value = {
        show: true,
        message,
        type
      }
      
      // 3秒后自动隐藏
      setTimeout(() => {
        notification.value.show = false
      }, 3000)
    }
    
    // 隐藏通知
    const hideNotification = () => {
      notification.value.show = false
    }
    
    // 拖拽上传
    const dragOver = (event) => {
      event.target.classList.add('drag-over')
    }
    
    const dragLeave = (event) => {
      event.target.classList.remove('drag-over')
    }
    
    const dropFile = (event) => {
      console.log('DEBUG: App.vue - dropFile 函数被调用', event)
      event.target.classList.remove('drag-over')
      const file = event.dataTransfer.files[0]
      console.log('DEBUG: App.vue - 从拖拽事件中获取的文件:', file)
      if (file) {
        console.log('DEBUG: App.vue - 文件信息:', {
          name: file.name,
          size: file.size,
          type: file.type,
          lastModified: file.lastModified
        })
        uploadForm.file = file
        console.log('DEBUG: App.vue - 文件已设置到 uploadForm.file，当前值:', uploadForm.file)
      } else {
        console.error('DEBUG: App.vue - 拖拽事件中没有找到文件')
      }
    }
    
    // 监听selectedVersion的变化
    watch(() => selectedVersion.value, (newValue, oldValue) => {
      console.log('DEBUG: selectedVersion变化 - 新值:', newValue, '旧值:', oldValue)
    })
    
    // 监听currentView的变化，当切换到version-view时确保数据正确传递
    watch(() => currentView.value, (newValue, oldValue) => {
      console.log('DEBUG: currentView变化 - 新值:', newValue, '旧值:', oldValue)
      if (newValue === 'version-view') {
        console.log('DEBUG: 切换到version-view视图，selectedVersion值:', selectedVersion.value)
        // 不再强制重置selectedVersion，依靠全局状态传递版本号
      }
    })
    
    // 计算属性，确保响应式更新
    const selectedVersionComputed = computed(() => {
      console.log('DEBUG: selectedVersionComputed被调用，当前值:', selectedVersion.value)
      return selectedVersion.value
    })
    
    // 组件挂载时获取文档列表
    onMounted(() => {
      fetchDocuments()
    })
    
    return {
      currentView,
      documents,
      documentVersions,
      selectedDocument,
      selectedDocumentId,
      isLoading,
      isUploading,
      searchQuery,
      filters,
      pagination,
      uploadForm,
      fetchDocuments,
      uploadDocument,
      updateDocument,
      deleteDocument,
      deleteDocumentVersion,
      viewDocument,
      viewDocumentVersions,
      viewDocumentVersion,
      downloadDocument,
      downloadVersion,
      editDocument,
      handleDocumentUpdated,
      globalVersionState,
      notification,
      showNotification,
      hideNotification,
      searchDocuments,
      applyFilters,
      changePage,
      getPaginationPages,
      fileSelected,
      dragOver,
      dragLeave,
      dropFile
    }
  }
}
</script>