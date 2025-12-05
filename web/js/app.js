// 从全局 Vue 对象中解构出常用 API（这是通过 CDN 引入 Vue 时的用法）
const { createApp, ref, reactive, onMounted, computed } = Vue;
// 创建并挂载 Vue 应用
createApp({
    // setup() 是 Vue 3 Composition API 的入口函数，在组件创建前执行
    setup() {
        // 响应式数据
        const currentView = ref('list');
        const documents = ref([]);
        const documentVersions = ref([]);
        const selectedDocument = ref(null);
        const isLoading = ref(false);
        const isUploading = ref(false);
        const searchQuery = ref('');
        
        // 过滤器
        const filters = reactive({
            library: '',
            type: '',
            version: '',
            status: ''
        });
        
        // 分页
        const pagination = reactive({
            page: 1,
            size: 10,
            total: 0
        });
        
        // 上传表单
        const uploadForm = reactive({
            name: '',
            type: '',
            version: '',
            library: '',
            description: '',
            tags: '',
            file: null
        });
        
        // API基础URL
        const apiBase = '/api/v1';
        
        // 获取文档列表
        const fetchDocuments = async () => {
            isLoading.value = true;
            try {
                const params = {
                    page: pagination.page,
                    size: pagination.size
                };
                
                // 添加过滤条件
                if (filters.library) params.library = filters.library;
                if (filters.type) params.type = filters.type;
                if (filters.version) params.version = filters.version;
                if (filters.status) params.status = filters.status;
                
                const response = await axios.get(`${apiBase}/documents`, { params });
                documents.value = response.data.data.items;
                pagination.total = response.data.data.total;
            } catch (error) {
                console.error('获取文档列表失败:', error);
                alert('获取文档列表失败: ' + (error.response?.data?.message || error.message));
            } finally {
                isLoading.value = false;
            }
        };
        
        // 获取文档版本列表
        const fetchDocumentVersions = async (documentId) => {
            isLoading.value = true;
            try {
                const response = await axios.get(`${apiBase}/documents/${documentId}/versions`);
                documentVersions.value = response.data.data;
            } catch (error) {
                console.error('获取文档版本失败:', error);
                alert('获取文档版本失败: ' + (error.response?.data?.message || error.message));
            } finally {
                isLoading.value = false;
            }
        };
        
        // 上传文档
        const uploadDocument = async () => {
            if (!uploadForm.file) {
                alert('请选择文件');
                return;
            }
            
            isUploading.value = true;
            const formData = new FormData();
            formData.append('file', uploadForm.file);
            formData.append('name', uploadForm.name);
            formData.append('type', uploadForm.type);
            formData.append('version', uploadForm.version);
            formData.append('library', uploadForm.library);
            formData.append('description', uploadForm.description);
            formData.append('tags', uploadForm.tags);
            
            try {
                const response = await axios.post(`${apiBase}/documents`, formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data'
                    }
                });
                
                alert('上传成功');
                resetUploadForm();
                currentView.value = 'list';
                await fetchDocuments();
            } catch (error) {
                console.error('上传文档失败:', error);
                alert('上传文档失败: ' + (error.response?.data?.message || error.message));
            } finally {
                isUploading.value = false;
            }
        };
        
        // 删除文档
        const deleteDocument = async (doc) => {
            if (!confirm(`确定要删除文档 "${doc.name}" 吗？`)) {
                return;
            }
            
            isLoading.value = true;
            try {
                await axios.delete(`${apiBase}/documents/${doc.id}`);
                alert('删除成功');
                await fetchDocuments();
            } catch (error) {
                console.error('删除文档失败:', error);
                alert('删除文档失败: ' + (error.response?.data?.message || error.message));
            } finally {
                isLoading.value = false;
            }
        };
        
        // 删除文档版本
        const deleteDocumentVersion = async (version) => {
            if (!confirm(`确定要删除版本 "${version.version}" 吗？`)) {
                return;
            }
            
            isLoading.value = true;
            try {
                // 注意：这里需要后端提供删除版本API，目前暂时调用删除整个文档的API
                // 实际实现中应该有专门的删除版本API
                alert('删除版本功能待实现');
            } catch (error) {
                console.error('删除文档版本失败:', error);
                alert('删除文档版本失败: ' + (error.response?.data?.message || error.message));
            } finally {
                isLoading.value = false;
            }
        };
        
        // 查看文档
        const viewDocument = (doc) => {
            // 这里可以实现查看文档详情的逻辑
            alert(`查看文档: ${doc.name}`);
        };
        
        // 查看文档版本
        const viewDocumentVersions = async (doc) => {
            selectedDocument.value = doc;
            currentView.value = 'versions';
            await fetchDocumentVersions(doc.id);
        };
        
        // 查看特定版本的文档
        const viewDocumentVersion = (version) => {
            // 这里可以实现查看特定版本文档的逻辑
            alert(`查看版本: ${version.version}`);
        };
        
        // 搜索文档
        const searchDocuments = async () => {
            if (!searchQuery.value.trim()) {
                await fetchDocuments();
                return;
            }
            
            isLoading.value = true;
            try {
                const params = {
                    page: 1,
                    size: pagination.size,
                    name: searchQuery.value.trim()
                };
                
                const response = await axios.get(`${apiBase}/documents`, { params });
                documents.value = response.data.data.items;
                pagination.total = response.data.data.total;
                pagination.page = 1;
            } catch (error) {
                console.error('搜索文档失败:', error);
                alert('搜索文档失败: ' + (error.response?.data?.message || error.message));
            } finally {
                isLoading.value = false;
            }
        };
        
        // 应用过滤器
        const applyFilters = async () => {
            pagination.page = 1;
            await fetchDocuments();
        };
        
        // 切换页面
        const changePage = async (page) => {
            if (page < 1 || (page - 1) * pagination.size >= pagination.total) {
                return;
            }
            pagination.page = page;
            await fetchDocuments();
        };
        
        // 获取分页页码
        const getPaginationPages = () => {
            const totalPages = Math.ceil(pagination.total / pagination.size);
            const pages = [];
            const maxVisiblePages = 5;
            
            let startPage = Math.max(1, pagination.page - Math.floor(maxVisiblePages / 2));
            let endPage = Math.min(totalPages, startPage + maxVisiblePages - 1);
            
            if (endPage - startPage + 1 < maxVisiblePages) {
                startPage = Math.max(1, endPage - maxVisiblePages + 1);
            }
            
            for (let i = startPage; i <= endPage; i++) {
                pages.push(i);
            }
            
            return pages;
        };
        
        // 文件选择
        const fileSelected = (event) => {
            const file = event.target.files[0];
            if (file) {
                uploadForm.file = file;
            }
        };
        
        // 拖拽上传
        const dragOver = (event) => {
            event.target.classList.add('drag-over');
        };
        
        const dragLeave = (event) => {
            event.target.classList.remove('drag-over');
        };
        
        const dropFile = (event) => {
            event.target.classList.remove('drag-over');
            const file = event.dataTransfer.files[0];
            if (file) {
                uploadForm.file = file;
            }
        };
        
        // 重置上传表单
        const resetUploadForm = () => {
            uploadForm.name = '';
            uploadForm.type = '';
            uploadForm.version = '';
            uploadForm.library = '';
            uploadForm.description = '';
            uploadForm.tags = '';
            uploadForm.file = null;
        };
        
        // 获取状态徽章样式
        const getStatusBadgeClass = (status) => {
            switch (status) {
                case 'uploading': return 'bg-secondary';
                case 'processing': return 'bg-warning';
                case 'completed': return 'bg-success';
                case 'failed': return 'bg-danger';
                default: return 'bg-light text-dark';
            }
        };
        
        // 获取状态文本
        const getStatusText = (status) => {
            switch (status) {
                case 'uploading': return '上传中';
                case 'processing': return '处理中';
                case 'completed': return '已完成';
                case 'failed': return '失败';
                default: return status;
            }
        };
        
        // 格式化文件大小
        const formatFileSize = (bytes) => {
            if (bytes === 0) return '0 Bytes';
            
            const k = 1024;
            const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            
            return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
        };
        
        // 格式化日期
        const formatDate = (dateString) => {
            const date = new Date(dateString);
            return date.toLocaleString('zh-CN');
        };
        
        // 组件挂载时获取文档列表
        onMounted(() => {
            fetchDocuments();
        });
        
        return {
            currentView,
            documents,
            documentVersions,
            selectedDocument,
            isLoading,
            isUploading,
            searchQuery,
            filters,
            pagination,
            uploadForm,
            fetchDocuments,
            fetchDocumentVersions,
            uploadDocument,
            deleteDocument,
            deleteDocumentVersion,
            viewDocument,
            viewDocumentVersions,
            viewDocumentVersion,
            searchDocuments,
            applyFilters,
            changePage,
            getPaginationPages,
            fileSelected,
            dragOver,
            dragLeave,
            dropFile,
            getStatusBadgeClass,
            getStatusText,
            formatFileSize,
            formatDate
        };
    }
}).mount('#app');