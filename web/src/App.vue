<template>
  <div id="app" class="container-fluid">
    <!-- 登录视图 -->
    <LoginView v-if="currentView === 'login'" @login-success="handleLoginSuccess" @show-register="currentView = 'register'" @show-forgot-password="currentView = 'forgot-password'" />
    
    <!-- 忘记密码视图 -->
    <ForgotPasswordView
      v-if="currentView === 'forgot-password'"
      @back-to-login="currentView = 'login'"
      @use-token="handleUseToken"
    />
    
    <!-- 重置密码视图 -->
    <ResetPasswordView v-if="currentView === 'reset-password'" :token="resetToken" @back-to-login="currentView = 'login'" />
    
    <!-- 注册视图 -->
    <RegisterView v-if="currentView === 'register'" @register-success="handleRegisterSuccess" @back-to-login="currentView = 'login'" />
    
    <div class="row" v-if="isAuthenticated">
      <!-- 侧边栏 -->
      <div class="col-md-3 sidebar p-4">
        <div class="d-flex justify-content-between align-items-center mb-4">
          <h4>AI技术文档库</h4>
          <div class="dropdown">
            <button class="btn btn-sm btn-outline-secondary dropdown-toggle" type="button" id="userDropdown" data-bs-toggle="dropdown" aria-expanded="false">
              <i class="bi bi-person-circle me-1"></i>{{ currentUser.username }}
            </button>
            <ul class="dropdown-menu" aria-labelledby="userDropdown">
              <li><a class="dropdown-item" href="#" @click="currentView = 'profile'">
                <i class="bi bi-person me-2"></i>个人资料
              </a></li>
              <li><a class="dropdown-item" href="#" @click="currentView = 'api-keys'">
                <i class="bi bi-key me-2"></i>API密钥
              </a></li>
              <li v-if="currentUser.role === 'admin'"><a class="dropdown-item" href="#" @click="currentView = 'user-management'">
                <i class="bi bi-people me-2"></i>用户管理
              </a></li>
              <li v-if="currentUser.role === 'admin'"><a class="dropdown-item" href="#" @click="currentView = 'monitor'">
                <i class="bi bi-activity me-2"></i>系统监控
              </a></li>
              <li><hr class="dropdown-divider"></li>
              <li><a class="dropdown-item" href="#" @click="logout">
                <i class="bi bi-box-arrow-right me-2"></i>退出登录
              </a></li>
            </ul>
          </div>
        </div>
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
               :class="{ active: currentView === 'search' }"
               @click="currentView = 'search'">
              <i class="bi bi-search me-2"></i>文档检索
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
          <li class="nav-item mb-2">
            <a class="nav-link"
               href="#"
               :class="{ active: currentView === 'ai-format' }"
               @click="currentView = 'ai-format'">
              <i class="bi bi-cpu me-2"></i>AI友好格式
            </a>
          </li>
          <li class="nav-item mb-2">
            <a class="nav-link"
               href="#"
               :class="{ active: currentView === 'mcp' }"
               @click="currentView = 'mcp'">
              <i class="bi bi-plug me-2"></i>MCP协议
            </a>
          </li>
          <li class="nav-item mb-2" v-if="currentUser.role === 'admin'">
            <a class="nav-link"
               href="#"
               :class="{ active: currentView === 'monitor' }"
               @click="currentView = 'monitor'">
              <i class="bi bi-activity me-2"></i>系统监控
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
            <label for="filterCategory" class="form-label">文档分类</label>
            <select class="form-select form-select-sm"
                    id="filterCategory"
                    v-model="filters.category">
              <option value="">全部分类</option>
              <option value="code">代码</option>
              <option value="document">文档</option>
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
          @navigate-to-upload="currentView = 'upload'"
        />
        
        <!-- 文档检索视图 -->
        <SearchView
          v-if="currentView === 'search'"
          @view-document="handleSearchViewDocument"
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
          @edit-version="editDocumentVersion"
          @add-version="addDocumentVersion"
          @delete-version="deleteDocumentVersion"
          @back="currentView = 'list'"
        />
        
        <!-- AI友好格式视图 -->
        <AIFormatView
          v-if="currentView === 'ai-format'"
        />
        
        <!-- MCP协议视图 -->
        <MCPView
          v-if="currentView === 'mcp'"
        />
        
        <!-- 系统监控视图（仅管理员） -->
        <MonitorView
          v-if="currentView === 'monitor' && currentUser.role === 'admin'"
        />
        
        <!-- 用户资料视图 -->
        <UserProfileView
          v-if="currentView === 'profile'"
          :current-user="currentUser"
          @profile-updated="handleProfileUpdated"
        />
        
        <!-- API密钥管理视图 -->
        <APIKeyManagementView
          v-if="currentView === 'api-keys'"
          :current-user="currentUser"
        />
        
        <!-- 用户管理视图（仅管理员） -->
        <UserManagementView
          v-if="currentView === 'user-management' && currentUser.role === 'admin'"
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
import SearchView from './views/SearchView.vue'
import AIFormatView from './views/AIFormatView.vue'
import MCPView from './views/MCPView.vue'
import LoginView from './views/LoginView.vue'
import RegisterView from './views/RegisterView.vue'
import ForgotPasswordView from './views/ForgotPasswordView.vue'
import ResetPasswordView from './views/ResetPasswordView.vue'
import UserProfileView from './views/UserProfileView.vue'
import APIKeyManagementView from './views/APIKeyManagementView.vue'
import UserManagementView from './views/UserManagementView.vue'
import MonitorView from './views/MonitorView.vue'
import axios from 'axios'
import { useDocumentService } from './utils/documentService'
import { useAuth } from './composables/useAuth'

export default {
  name: 'App',
  components: {
    DocumentListView,
    MonitorView,
    DocumentView,
    DocumentEditView,
    DocumentVersionView,
    UploadView,
    VersionManagementView,
    SearchView,
    AIFormatView,
    MCPView,
    LoginView,
    RegisterView,
    ForgotPasswordView,
    ResetPasswordView,
    UserProfileView,
    APIKeyManagementView,
    UserManagementView
  },
  setup() {
    // 处理使用令牌重置密码
    const handleUseToken = (token) => {
      resetToken.value = token;
      currentView.value = 'reset-password';
    };
    
    // 认证相关
    const {
      isAuthenticated,
      currentUser,
      login,
      logout,
      checkAuthStatus
    } = useAuth()
    
    const currentView = ref('login')
    const documents = ref([])
    const documentVersions = ref([])
    const selectedDocument = ref(null)
    const selectedDocumentId = ref(null)
    const selectedVersion = ref(null)
    const resetToken = ref('')
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
      category: '',
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
      category: '',
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
      updateDocumentVersion,
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
    const viewDocument = async (doc) => {
      console.log('DEBUG: 尝试查看文档', doc)
      selectedDocumentId.value = doc.id
      
      // 获取最新版本并设置为默认版本
      try {
        const latestVersion = await fetch(`/api/v1/documents/${doc.id}/versions/latest`)
        const result = await latestVersion.json()
        
        if (result.code === 200) {
          selectedVersion.value = result.data.version
          console.log('DEBUG: 获取到最新版本:', selectedVersion.value)
        }
      } catch (error) {
        console.error('获取最新版本失败:', error)
        // 如果获取失败，使用文档本身的版本
        selectedVersion.value = doc.version
      }
      
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
    
    // 编辑文档版本
    const editDocumentVersion = async (version) => {
      console.log('DEBUG: 尝试编辑文档版本', version)
      try {
        await updateDocumentVersion(version)
        // 重新获取文档版本列表
        await fetchDocumentVersions(version.document_id)
        alert('编辑版本成功')
      } catch (error) {
        console.error('编辑文档版本失败:', error)
        alert('编辑文档版本失败: ' + (error.response?.data?.message || error.message))
      }
    }
    
    // 添加文档新版本
    const addDocumentVersion = async (formData) => {
      console.log('DEBUG: 尝试添加文档新版本', formData)
      try {
        isLoading.value = true
        
        const response = await axios.post('/api/v1/documents', formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
        
        console.log('DEBUG: 添加新版本成功:', response.data)
        alert('添加新版本成功')
        
        // 重新获取文档版本列表
        await fetchDocumentVersions(selectedDocument.value.id)
      } catch (error) {
        console.error('添加文档新版本失败:', error)
        alert('添加文档新版本失败: ' + (error.response?.data?.message || error.message))
      } finally {
        isLoading.value = false
      }
    }
    
    // 处理文档检索中的查看文档
    const handleSearchViewDocument = async (data) => {
      console.log('DEBUG: 处理文档检索查看文档 - 完整数据:', data)
      console.log('DEBUG: 搜索结果中的版本号:', data.version, '文档ID:', data.documentId)
      
      selectedDocumentId.value = data.documentId
      
      let versionToUse = data.version
      
      // 如果版本是"latest"，需要先获取实际的最新版本号
      if (data.version === 'latest') {
        console.log('DEBUG: 检测到latest版本，尝试获取实际版本号')
        try {
          const response = await fetch(`/api/v1/documents/${data.documentId}/versions/latest`)
          const result = await response.json()
          
          if (result.code === 200 && result.data) {
            versionToUse = result.data.version
            console.log('DEBUG: 获取到实际版本号:', versionToUse)
            selectedVersion.value = versionToUse
          } else {
            console.error('DEBUG: 获取最新版本失败 - 响应:', result)
            alert('获取最新版本失败: ' + (result.message || '未知错误'))
            return
          }
        } catch (err) {
          console.error('获取最新版本失败:', err)
          alert('获取最新版本失败，请稍后重试')
          return
        }
      }
      
      // 如果版本为空或无效，使用最新版本
      if (!versionToUse || versionToUse === '') {
        console.log('DEBUG: 版本号为空，尝试获取最新版本')
        try {
          const response = await fetch(`/api/v1/documents/${data.documentId}/versions/latest`)
          const result = await response.json()
          
          if (result.code === 200 && result.data) {
            versionToUse = result.data.version
            console.log('DEBUG: 获取到最新版本号作为fallback:', versionToUse)
          }
        } catch (err) {
          console.error('获取最新版本失败:', err)
        }
      }
      
      // 验证版本是否有效
      try {
        const response = await fetch(`/api/v1/documents/${data.documentId}/versions/${versionToUse}`)
        const result = await response.json()
        console.log('DEBUG: 验证版本号 - 版本:', versionToUse, '响应:', result)
        
        if (result.code !== 200) {
          console.error('DEBUG: 版本号无效，尝试获取最新版本')
          alert(`版本号 ${versionToUse} 无效，将使用最新版本`)
          
          const latestResponse = await fetch(`/api/v1/documents/${data.documentId}/versions/latest`)
          const latestResult = await latestResponse.json()
          
          if (latestResult.code === 200 && latestResult.data) {
            versionToUse = latestResult.data.version
            console.log('DEBUG: 使用最新版本号:', versionToUse)
          }
        }
      } catch (err) {
        console.error('DEBUG: 验证版本失败:', err)
      }
      
      // 设置选中的版本号
      selectedVersion.value = versionToUse
      console.log('DEBUG: 最终选择的版本号:', selectedVersion.value)
      
      // 设置全局状态
      globalVersionState.documentId = data.documentId
      globalVersionState.version = versionToUse
      globalVersionState.isSet = true
      console.log('DEBUG: 全局状态已设置:', globalVersionState)
      
      // 切换到版本查看视图
      currentView.value = 'version-view'
      console.log('DEBUG: 切换到版本查看视图')
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
    
    // 处理登录成功
    const handleLoginSuccess = async (user) => {
      console.log('DEBUG: 收到登录成功事件', user)
      
      // 显示登录成功消息
      showNotification('登录成功！', 'success')
      
      // 强制触发认证状态更新
      await checkAuthStatus()
      
      // 使用一个技巧来强制Vue重新渲染
      // 创建一个临时的空对象来触发响应式更新
      const temp = {}
      temp.forceUpdate = Math.random()
      
      // 切换视图
      currentView.value = 'list'
      
      // 使用 nextTick 等待视图更新
      await nextTick()
      
      // 获取文档列表
      try {
        await fetchDocuments()
        console.log('DEBUG: 登录后文档列表获取成功，数量:', documents.value?.length)
      } catch (error) {
        console.error('登录后获取文档列表失败:', error)
      }
      
      console.log('DEBUG: 登录成功处理完成', {
        isAuthenticated: isAuthenticated.value,
        currentUser: currentUser.value,
        currentView: currentView.value,
        documentsLength: documents.value?.length
      })
    }
    
    // 处理注册成功
    const handleRegisterSuccess = (user) => {
      console.log('DEBUG: 注册成功', user)
      currentView.value = 'login'
      showNotification('注册成功！请登录', 'success')
    }
    
    // 处理登出
    const handleLogout = () => {
      logout()
      currentView.value = 'login'
      showNotification('已退出登录', 'info')
    }
    
    // 处理个人资料更新
    const handleProfileUpdated = (updatedUser) => {
      console.log('DEBUG: 个人资料更新成功', updatedUser)
      showNotification('个人资料更新成功！', 'success')
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
    
    // 初始化应用
    onMounted(async () => {
      // 清除可能存在的标记
      sessionStorage.removeItem('justLoggedIn')
      
      console.log('DEBUG: 应用初始化，检查认证状态')
      
      // 检查认证状态
      const isAuth = await checkAuthStatus()
      console.log('DEBUG: 初始认证状态', isAuth, '当前用户:', currentUser.value)
      
      if (isAuth) {
        console.log('DEBUG: 用户已认证，显示文档列表')
        currentView.value = 'list'
        try {
          await fetchDocuments()
          console.log('DEBUG: 文档列表获取成功，数量:', documents.value?.length)
        } catch (error) {
          console.error('获取文档列表失败:', error)
        }
      } else {
        console.log('DEBUG: 用户未认证，显示登录页面')
        currentView.value = 'login'
      }
    })
    
    return {
      currentView,
      documents,
      documentVersions,
      selectedDocument,
      selectedDocumentId,
      selectedVersion,
      selectedVersionComputed,
      resetToken,
      isLoading,
      isUploading,
      searchQuery,
      notification,
      globalVersionState,
      filters,
      pagination,
      uploadForm,
      isAuthenticated,
      currentUser,
      fetchDocuments,
      fetchDocumentVersions,
      uploadDocument,
      updateDocument,
      deleteDocument,
      deleteDocumentVersion,
      searchDocuments,
      applyFilters,
      changePage,
      getPaginationPages,
      viewDocument,
      viewDocumentVersions,
      viewDocumentVersion,
      handleSearchViewDocument,
      fileSelected,
      downloadDocument,
      downloadVersion,
      addDocumentVersion,
      editDocumentVersion,
      editDocument,
      handleDocumentUpdated,
      showNotification,
      hideNotification,
      dragOver,
      dragLeave,
      dropFile,
      handleLoginSuccess,
      handleRegisterSuccess,
      logout: handleLogout,
      handleProfileUpdated,
      handleUseToken
    }
  }
}
</script>