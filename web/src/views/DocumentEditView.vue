<template>
  <div>
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h2>编辑文档 - {{ editForm.name }}</h2>
      <button class="btn btn-outline-secondary" @click="$emit('back')">返回</button>
    </div>
    
    <div class="card">
      <div class="card-body">
        <form @submit.prevent="handleSubmit">
          <div class="mb-3">
            <label for="editDocName" class="form-label">文档名称 <span class="text-danger">*</span></label>
            <input type="text" 
                   class="form-control" 
                   id="editDocName" 
                   v-model="editForm.name" 
                   required>
          </div>
          
          <div class="row">
            <div class="col-md-6 mb-3">
              <label for="editDocType" class="form-label">文档类型 <span class="text-danger">*</span></label>
              <select class="form-select" 
                      id="editDocType" 
                      v-model="editForm.type" 
                      required>
                <option value="">请选择文档类型</option>
                <option value="markdown">Markdown</option>
                <option value="pdf">PDF</option>
                <option value="docx">DOCX</option>
                <option value="swagger">Swagger</option>
                <option value="openapi">OpenAPI</option>
                <option value="java_doc">JavaDoc</option>
              </select>
            </div>
            
            <div class="col-md-6 mb-3">
              <label for="editDocVersion" class="form-label">版本 <span class="text-danger">*</span></label>
              <input type="text" 
                     class="form-control" 
                     id="editDocVersion" 
                     v-model="editForm.version" 
                     placeholder="例如: 1.0.0" 
                     required>
            </div>
          </div>
          
          <div class="mb-3">
            <label for="editDocLibrary" class="form-label">所属库 <span class="text-danger">*</span></label>
            <input type="text" 
                   class="form-control" 
                   id="editDocLibrary" 
                   v-model="editForm.library" 
                   placeholder="例如: react" 
                   required>
          </div>
          
          <div class="mb-3">
            <label for="editDocDescription" class="form-label">描述</label>
            <textarea class="form-control" 
                      id="editDocDescription" 
                      v-model="editForm.description" 
                      rows="3"></textarea>
          </div>
          
          <div class="mb-3">
            <label for="editDocTags" class="form-label">标签</label>
            <input type="text" 
                   class="form-control" 
                   id="editDocTags" 
                   v-model="editForm.tags" 
                   placeholder="用逗号分隔多个标签">
          </div>
          
          <div class="mb-3">
            <label for="editDocStatus" class="form-label">状态</label>
            <select class="form-select" 
                    id="editDocStatus" 
                    v-model="editForm.status">
              <option value="uploading">上传中</option>
              <option value="processing">处理中</option>
              <option value="completed">已完成</option>
              <option value="failed">失败</option>
            </select>
          </div>
          
          <div class="mb-3">
            <div class="form-check">
              <input class="form-check-input" type="checkbox" id="createNewVersion" v-model="createNewVersion">
              <label class="form-check-label" for="createNewVersion">
                创建新版本
              </label>
            </div>
            <div v-if="createNewVersion" class="mt-2">
              <label for="newVersionNumber" class="form-label">新版本号</label>
              <input type="text"
                     class="form-control"
                     id="newVersionNumber"
                     v-model="newVersionNumber"
                     placeholder="例如: 1.0.1">
            </div>
          </div>
          
          <div class="d-flex justify-content-between">
            <button type="button" class="btn btn-outline-secondary" @click="$emit('back')">取消</button>
            <button type="submit" class="btn btn-primary" :disabled="isUpdating">
              <span v-if="isUpdating" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
              保存
            </button>
          </div>
        </form>
      </div>
    </div>
    
    <!-- 错误状态 -->
    <div v-if="error" class="alert alert-danger mt-3" role="alert">
      {{ error }}
    </div>
  </div>
</template>

<script>
import { ref, reactive, onMounted, watch } from 'vue'

export default {
  name: 'DocumentEditView',
  props: {
    documentId: {
      type: String,
      required: true
    }
  },
  emits: ['back', 'document-updated'],
  setup(props, { emit }) {
    const document = ref(null)
    const isLoading = ref(false)
    const isUpdating = ref(false)
    const error = ref(null)
    const createNewVersion = ref(false)
    const newVersionNumber = ref('')
    
    const editForm = reactive({
      name: '',
      type: '',
      version: '',
      library: '',
      description: '',
      tags: '',
      status: ''
    })
    
    // 获取文档详情
    const fetchDocument = async () => {
      isLoading.value = true
      error.value = null
      
      try {
        const response = await fetch(`/api/v1/documents/${props.documentId}`)
        const result = await response.json()
        
        if (result.code === 200) {
          document.value = result.data
          // 填充表单
          editForm.name = result.data.name
          editForm.type = result.data.type
          editForm.version = result.data.version
          editForm.library = result.data.library
          editForm.description = result.data.description || ''
          editForm.tags = result.data.tags ? result.data.tags.join(', ') : ''
          editForm.status = result.data.status
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
    
    // 处理表单提交
    const handleSubmit = async () => {
      isUpdating.value = true
      error.value = null
      
      try {
        const updates = {
          name: editForm.name,
          type: editForm.type,
          version: editForm.version,
          library: editForm.library,
          description: editForm.description,
          tags: editForm.tags.split(',').map(tag => tag.trim()).filter(tag => tag),
          status: editForm.status
        }
        
        // 这里应该使用 inject 获取 updateDocument 函数，但暂时还是使用 fetch
        const response = await fetch(`/api/v1/documents/${props.documentId}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(updates)
        })
        
        const result = await response.json()
        
        if (result.code === 200) {
          alert('更新成功')
          emit('document-updated')
          emit('back')
        } else {
          error.value = result.message || '更新文档失败'
        }
      } catch (err) {
        console.error('更新文档失败:', err)
        error.value = '更新文档失败: ' + err.message
      } finally {
        isUpdating.value = false
      }
    }
    
    // 组件挂载时获取文档详情
    onMounted(() => {
      fetchDocument()
    })
    
    return {
      document,
      isLoading,
      isUpdating,
      error,
      createNewVersion,
      newVersionNumber,
      editForm,
      handleSubmit
    }
  }
}
</script>