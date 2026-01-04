<template>
  <div>
    <div class="d-flex justify-content-between align-items-center mb-4">
      <div>
        <h2>版本管理 - {{ selectedDocument.name }}</h2>
        <p class="text-muted">以下列出了该文档的所有版本，您可以通过版本号区分不同的版本</p>
      </div>
      <div>
        <button class="btn btn-primary me-2" @click="openAddModal">
          <i class="bi bi-plus-circle me-1"></i>添加新版本
        </button>
        <button class="btn btn-outline-secondary" @click="$emit('back')">返回</button>
      </div>
    </div>
    
    <!-- 编辑版本模态框 -->
    <div class="modal" :class="{ show: showEditModal }" tabindex="-1" :style="{ display: showEditModal ? 'block' : 'none' }">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">编辑版本信息</h5>
            <button type="button" class="btn-close" @click="closeEditModal"></button>
          </div>
          <div class="modal-body">
            <div class="mb-3">
              <label for="editVersion" class="form-label">版本号</label>
              <input
                type="text"
                class="form-control"
                id="editVersion"
                v-model="editForm.version"
                placeholder="例如: 1.0.0">
            </div>
            <div class="mb-3">
              <label for="editDescription" class="form-label">版本描述</label>
              <textarea
                class="form-control"
                id="editDescription"
                v-model="editForm.description"
                rows="3"
                placeholder="描述此版本的主要更改和特性"></textarea>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="closeEditModal">取消</button>
            <button type="button" class="btn btn-primary" @click="saveEdit">保存</button>
          </div>
        </div>
      </div>
    </div>
    <div class="modal-backdrop" :class="{ show: showEditModal }" :style="{ display: showEditModal ? 'block' : 'none' }"></div>
    
    <!-- 添加新版本模态框 -->
    <div class="modal" :class="{ show: showAddModal }" tabindex="-1" :style="{ display: showAddModal ? 'block' : 'none' }">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">添加新版本</h5>
            <button type="button" class="btn-close" @click="closeAddModal"></button>
          </div>
          <div class="modal-body">
            <div class="mb-3">
              <label for="addVersion" class="form-label">版本号</label>
              <input
                type="text"
                class="form-control"
                id="addVersion"
                v-model="addForm.version"
                placeholder="例如: 1.0.0">
            </div>
            <div class="mb-3">
              <label for="addDescription" class="form-label">版本描述</label>
              <textarea
                class="form-control"
                id="addDescription"
                v-model="addForm.description"
                rows="3"
                placeholder="描述此版本的主要更改和特性"></textarea>
            </div>
            <div class="mb-3">
              <label for="addFile" class="form-label">选择文件</label>
              <input
                type="file"
                class="form-control"
                id="addFile"
                @change="handleFileChange"
                accept=".pdf,.docx,.md,.json">
              <div v-if="addForm.file" class="mt-2">
                <small class="text-muted">
                  已选择: {{ addForm.file.name }} ({{ formatFileSize(addForm.file.size) }})
                </small>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="closeAddModal">取消</button>
            <button type="button" class="btn btn-primary" @click="saveAdd" :disabled="!addForm.file">
              <i class="bi bi-upload me-1"></i>上传新版本
            </button>
          </div>
        </div>
      </div>
    </div>
    <div class="modal-backdrop" :class="{ show: showAddModal }" :style="{ display: showAddModal ? 'block' : 'none' }"></div>
    
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
              <span class="badge bg-secondary ms-1" v-if="currentVersion && currentVersion.version === version.version">当前版本</span>
            </div>
            <p class="text-muted small">创建时间: {{ formatDate(version.created_at) }}</p>
          </div>
          <div>
            <button class="btn btn-sm btn-outline-primary me-2" @click="$emit('view-version', version)">查看</button>
            <button class="btn btn-sm btn-outline-warning me-2" @click="openEditModal(version)">编辑</button>
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
import { ref, computed } from 'vue'

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
  emits: ['view-version', 'edit-version', 'add-version', 'delete-version', 'back'],
  setup(props, { emit }) {
    // 编辑模态框状态
    const showEditModal = ref(false)
    const editForm = ref({
      version: '',
      description: '',
      documentId: '',
      oldVersion: ''
    })
    
    // 添加新版本模态框状态
    const showAddModal = ref(false)
    const addForm = ref({
      version: '',
      description: '',
      file: null
    })
    
    // 打开添加模态框
    const openAddModal = () => {
      addForm.value = {
        version: '',
        description: '',
        file: null
      }
      showAddModal.value = true
    }
    
    // 关闭添加模态框
    const closeAddModal = () => {
      showAddModal.value = false
      addForm.value = {
        version: '',
        description: '',
        file: null
      }
    }
    
    // 处理文件选择
    const handleFileChange = (event) => {
      const file = event.target.files[0]
      if (file) {
        addForm.value.file = file
      }
    }
    
    // 保存添加
    const saveAdd = () => {
      if (!addForm.value.version.trim()) {
        alert('版本号不能为空')
        return
      }
      if (!addForm.value.file) {
        alert('请选择文件')
        return
      }
      
      const formData = new FormData()
      formData.append('file', addForm.value.file)
      formData.append('name', props.selectedDocument.name)
      formData.append('type', props.selectedDocument.type)
      formData.append('category', props.selectedDocument.category)
      formData.append('version', addForm.value.version)
      formData.append('library', props.selectedDocument.library)
      formData.append('description', addForm.value.description)
      formData.append('tags', props.selectedDocument.tags ? props.selectedDocument.tags.join(',') : '')
      
      emit('add-version', formData)
      closeAddModal()
    }
    
    // 打开编辑模态框
    const openEditModal = (version) => {
      editForm.value = {
        version: version.version,
        description: version.description || '',
        documentId: version.document_id,
        oldVersion: version.version
      }
      showEditModal.value = true
    }
    
    // 关闭编辑模态框
    const closeEditModal = () => {
      showEditModal.value = false
      editForm.value = {
        version: '',
        description: '',
        documentId: '',
        oldVersion: ''
      }
    }
    
    // 保存编辑
    const saveEdit = () => {
      if (!editForm.value.version.trim()) {
        alert('版本号不能为空')
        return
      }
      
      emit('edit-version', {
        document_id: editForm.value.documentId,
        version: editForm.value.oldVersion,
        newVersion: editForm.value.version,
        description: editForm.value.description
      })
      
      closeEditModal()
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
    
    // 计算当前版本（根据更新时间）
    const currentVersion = computed(() => {
      if (!props.documentVersions || props.documentVersions.length === 0) {
        return null
      }
      
      // 找到更新时间最新的版本，使用循环避免reduce的null问题
      if (props.documentVersions.length > 0) {
        let latestVersion = props.documentVersions[0]
        for (let i = 1; i < props.documentVersions.length; i++) {
          const version = props.documentVersions[i]
          const current = new Date(version.updated_at)
          const latest = new Date(latestVersion.updated_at)
          if (current > latest) {
            latestVersion = version
          }
        }
        return latestVersion
      }
      return null
    })
    
    return {
      showEditModal,
      editForm,
      openEditModal,
      closeEditModal,
      saveEdit,
      showAddModal,
      addForm,
      openAddModal,
      closeAddModal,
      handleFileChange,
      saveAdd,
      currentVersion,
      getStatusBadgeClass,
      getStatusText,
      formatFileSize,
      formatDate
    }
  }
}
</script>