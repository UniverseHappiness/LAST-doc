<template>
  <div>
    <h2 class="mb-4">上传文档</h2>
    
    <div class="card">
      <div class="card-body">
        <form @submit.prevent="handleSubmit" id="uploadForm">
          <div class="mb-3">
            <label for="docName" class="form-label">文档名称 <span class="text-danger">*</span></label>
            <input type="text" 
                   class="form-control" 
                   id="docName" 
                   v-model="uploadForm.name" 
                   required>
          </div>
          
          <div class="row">
            <div class="col-md-6 mb-3">
              <label for="docType" class="form-label">文档类型 <span class="text-danger">*</span></label>
              <select class="form-select"
                      id="docType"
                      v-model="uploadForm.type"
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
              <label for="docCategory" class="form-label">文档分类 <span class="text-danger">*</span></label>
              <select class="form-select"
                      id="docCategory"
                      v-model="uploadForm.category"
                      required>
                <option value="">请选择文档分类</option>
                <option value="code">代码</option>
                <option value="document">文档</option>
              </select>
            </div>
          </div>
          
          <div class="mb-3">
            <label for="docVersion" class="form-label">版本 <span class="text-danger">*</span></label>
            <input type="text"
                   class="form-control"
                   id="docVersion"
                   v-model="uploadForm.version"
                   placeholder="例如: 1.0.0"
                   required>
          </div>
          
          <div class="mb-3">
            <label for="docLibrary" class="form-label">所属库 <span class="text-danger">*</span></label>
            <input type="text" 
                   class="form-control" 
                   id="docLibrary" 
                   v-model="uploadForm.library" 
                   placeholder="例如: react" 
                   required>
          </div>
          
          <div class="mb-3">
            <label for="docDescription" class="form-label">描述</label>
            <textarea class="form-control" 
                      id="docDescription" 
                      v-model="uploadForm.description" 
                      rows="3"></textarea>
          </div>
          
          <div class="mb-3">
            <label for="docTags" class="form-label">标签</label>
            <input type="text" 
                   class="form-control" 
                   id="docTags" 
                   v-model="uploadForm.tags" 
                   placeholder="用逗号分隔多个标签">
          </div>
          
          <div class="mb-3">
            <label class="form-label">上传文件 <span class="text-danger">*</span></label>
            <div class="upload-area" 
                 @click="$refs.fileInput.click()" 
                 @dragover.prevent="handleDragOver" 
                 @dragleave.prevent="handleDragLeave" 
                 @drop.prevent="handleDropFile">
              <input type="file"
                     ref="fileInput"
                     @change="handleFileSelected"
                     style="display: none;">
              <i class="bi bi-cloud-upload" style="font-size: 2rem;"></i>
              <p class="mt-2">点击或拖拽文件到此处上传</p>
              <p class="text-muted small">支持多种格式: Markdown, PDF, DOCX, Swagger, OpenAPI, JavaDoc</p>
              <div v-if="uploadForm.file" class="mt-2">
                <strong>已选择文件:</strong> {{ uploadForm.file.name }}
              </div>
            </div>
          </div>
          
          <div class="d-flex justify-content-between">
            <button type="button" class="btn btn-outline-secondary" @click="$emit('back')">返回</button>
            <button type="submit" class="btn btn-primary" :disabled="isUploading">
              <span v-if="isUploading" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
              上传
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, watch } from 'vue'

export default {
  name: 'UploadView',
  props: {
    isUploading: {
      type: Boolean,
      required: true
    },
    uploadForm: {
      type: Object,
      required: true
    }
  },
  emits: ['upload', 'file-selected', 'drag-over', 'drag-leave', 'drop-file', 'back'],
  setup(props, { emit }) {
    const fileInput = ref(null)
    
    // 处理表单提交
    const handleSubmit = () => {
      console.log('DEBUG: UploadView - handleSubmit 函数被调用')
      // 检查是否有文件（无论是通过拖拽还是点击选择）
      const hasFile = props.uploadForm.file !== null
      console.log('DEBUG: UploadView - 检查是否有文件:', hasFile)
      
      if (!hasFile) {
        console.error('DEBUG: UploadView - 没有选择文件，阻止表单提交')
        alert('请选择文件')
        return
      }
      
      console.log('DEBUG: UploadView - 文件已选择，允许表单提交')
      emit('upload')
    }
    
    // 处理文件选择
    const handleFileSelected = (event) => {
      emit('file-selected', event)
    }
    
    // 处理拖拽事件
    const handleDragOver = (event) => {
      emit('drag-over', event)
    }
    
    const handleDragLeave = (event) => {
      emit('drag-leave', event)
    }
    
    const handleDropFile = (event) => {
      console.log('DEBUG: UploadView - 拖拽文件事件触发', event)
      console.log('DEBUG: UploadView - 拖拽的文件列表:', event.dataTransfer.files)
      if (event.dataTransfer.files && event.dataTransfer.files.length > 0) {
        console.log('DEBUG: UploadView - 第一个文件信息:', {
          name: event.dataTransfer.files[0].name,
          size: event.dataTransfer.files[0].size,
          type: event.dataTransfer.files[0].type
        })
      }
      emit('drop-file', event)
    }
    
    return {
      fileInput,
      handleSubmit,
      handleFileSelected,
      handleDragOver,
      handleDragLeave,
      handleDropFile
    }
  }
}
</script>