<template>
  <div class="ai-format-view">
    <div class="header">
      <h1>AI友好格式</h1>
      <p>将文档转换为AI友好的格式，便于AI模型理解和使用</p>
    </div>

    <div class="content">
      <!-- 文档选择 -->
      <div class="document-selector">
        <h2>选择文档</h2>
        <div class="selection-controls">
          <div class="form-group">
            <label for="document-select">文档</label>
            <select id="document-select" v-model="selectedDocumentId" @change="onDocumentChange">
              <option value="">请选择文档</option>
              <option v-for="doc in documents" :key="doc.id" :value="doc.id">
                {{ doc.name }} ({{ doc.library }})
              </option>
            </select>
          </div>
          
          <div class="form-group" v-if="selectedDocumentId">
            <label for="version-select">版本</label>
            <select id="version-select" v-model="selectedVersion" @change="onVersionChange">
              <option value="">请选择版本</option>
              <option v-for="version in documentVersions" :key="version.id" :value="version.version">
                {{ version.version }} ({{ formatTime(version.created_at) }})
              </option>
            </select>
          </div>
        </div>
      </div>

      <!-- 格式选项 -->
      <div class="format-options" v-if="selectedDocumentId && selectedVersion">
        <h2>格式选项</h2>
        <div class="option-tabs">
          <button 
            class="tab-button" 
            :class="{ active: activeTab === 'structured' }"
            @click="activeTab = 'structured'"
          >
            结构化内容
          </button>
          <button 
            class="tab-button" 
            :class="{ active: activeTab === 'llm' }"
            @click="activeTab = 'llm'"
          >
            LLM优化格式
          </button>
          <button 
            class="tab-button" 
            :class="{ active: activeTab === 'multigranularity' }"
            @click="activeTab = 'multigranularity'"
          >
            多粒度表示
          </button>
          <button 
            class="tab-button" 
            :class="{ active: activeTab === 'context' }"
            @click="activeTab = 'context'"
          >
            上下文注入
          </button>
        </div>

        <!-- 结构化内容选项 -->
        <div v-if="activeTab === 'structured'" class="tab-content">
          <button class="action-button" @click="getStructuredContent" :disabled="loading">
            获取结构化内容
          </button>
        </div>

        <!-- LLM优化格式选项 -->
        <div v-if="activeTab === 'llm'" class="tab-content">
          <div class="form-group">
            <label for="max-tokens">最大Token数</label>
            <input 
              id="max-tokens" 
              type="number" 
              v-model="llmOptions.maxTokens" 
              min="100" 
              max="8000"
            />
          </div>
          
          <div class="form-group">
            <label>摘要级别</label>
            <div class="radio-group">
              <label>
                <input type="radio" v-model="llmOptions.summaryLevel" value="brief" />
                简要
              </label>
              <label>
                <input type="radio" v-model="llmOptions.summaryLevel" value="medium" />
                中等
              </label>
              <label>
                <input type="radio" v-model="llmOptions.summaryLevel" value="detailed" />
                详细
              </label>
            </div>
          </div>
          
          <div class="form-group">
            <label>
              <input type="checkbox" v-model="llmOptions.preserveCode" />
              保留代码示例
            </label>
          </div>
          
          <div class="form-group">
            <label>
              <input type="checkbox" v-model="llmOptions.includeMetadata" />
              包含元数据
            </label>
          </div>
          
          <button class="action-button" @click="generateLLMFormat" :disabled="loading">
            生成LLM优化格式
          </button>
        </div>

        <!-- 多粒度表示选项 -->
        <div v-if="activeTab === 'multigranularity'" class="tab-content">
          <button class="action-button" @click="generateMultiGranularity" :disabled="loading">
            生成多粒度表示
          </button>
        </div>

        <!-- 上下文注入选项 -->
        <div v-if="activeTab === 'context'" class="tab-content">
          <div class="form-group">
            <label for="context-query">查询内容</label>
            <textarea 
              id="context-query" 
              v-model="contextOptions.query" 
              placeholder="输入您的查询内容，系统将根据查询内容选择相关的文档片段"
              rows="4"
            ></textarea>
          </div>
          
          <div class="form-group">
            <label for="max-context-size">最大上下文大小</label>
            <input 
              id="max-context-size" 
              type="number" 
              v-model="contextOptions.maxContextSize" 
              min="100" 
              max="4000"
            />
          </div>
          
          <div class="form-group">
            <label>优先级</label>
            <div class="radio-group">
              <label>
                <input type="radio" v-model="contextOptions.priorityLevel" value="low" />
                低
              </label>
              <label>
                <input type="radio" v-model="contextOptions.priorityLevel" value="medium" />
                中
              </label>
              <label>
                <input type="radio" v-model="contextOptions.priorityLevel" value="high" />
                高
              </label>
            </div>
          </div>
          
          <div class="form-group">
            <label>输出格式</label>
            <div class="radio-group">
              <label>
                <input type="radio" v-model="contextOptions.format" value="markdown" />
                Markdown
              </label>
              <label>
                <input type="radio" v-model="contextOptions.format" value="json" />
                JSON
              </label>
              <label>
                <input type="radio" v-model="contextOptions.format" value="plain_text" />
                纯文本
              </label>
            </div>
          </div>
          
          <div class="form-group">
            <label>
              <input type="checkbox" v-model="contextOptions.includeCode" />
              包含代码示例
            </label>
          </div>
          
          <button class="action-button" @click="injectContext" :disabled="loading || !contextOptions.query">
            注入上下文
          </button>
        </div>
      </div>

      <!-- 结果展示 -->
      <div class="result-section" v-if="result">
        <div class="result-header">
          <h2>{{ resultTitle }}</h2>
          <button class="copy-button" @click="copyResult" title="复制结果">
            📋 复制
          </button>
        </div>
        
        <div class="result-content">
          <pre v-if="typeof result === 'string'">{{ result }}</pre>
          <div v-else-if="resultType === 'json'">
            <pre><code>{{ JSON.stringify(result, null, 2) }}</code></pre>
          </div>
          <div v-else-if="resultType === 'structured'">
            <div class="structured-result">
              <div class="section" v-for="segment in result.segments" :key="segment.id">
                <h3>{{ segment.title }}</h3>
                <p>{{ segment.content }}</p>
                <div class="annotations" v-if="segment.annotations && segment.annotations.length">
                  <h4>语义标注</h4>
                  <ul>
                    <li v-for="annotation in segment.annotations" :key="annotation.id">
                      {{ annotation.type }}: {{ annotation.value }}
                    </li>
                  </ul>
                </div>
              </div>
              
              <div class="section" v-if="result.codeExamples && result.codeExamples.length">
                <h3>代码示例</h3>
                <div v-for="example in result.codeExamples" :key="example.id" class="code-example">
                  <h4>{{ example.description }}</h4>
                  <pre><code :class="'language-' + example.language">{{ example.code }}</code></pre>
                </div>
              </div>
            </div>
          </div>
          <div v-else>
            <pre><code>{{ JSON.stringify(result, null, 2) }}</code></pre>
          </div>
        </div>
      </div>

      <!-- 加载状态 -->
      <div class="loading" v-if="loading">
        <div class="spinner"></div>
        <p>处理中，请稍候...</p>
      </div>

      <!-- 错误信息 -->
      <div class="error" v-if="error">
        <h3>错误</h3>
        <p>{{ error }}</p>
      </div>
    </div>
  </div>
</template>

<script>
import { documentService } from '../utils/documentService';

export default {
  name: 'AIFormatView',
  data() {
    return {
      documents: [],
      documentVersions: [],
      selectedDocumentId: '',
      selectedVersion: '',
      activeTab: 'structured',
      loading: false,
      error: null,
      result: null,
      resultType: 'json',
      resultTitle: '',
      llmOptions: {
        maxTokens: 4000,
        summaryLevel: 'medium',
        preserveCode: true,
        includeMetadata: true
      },
      contextOptions: {
        query: '',
        maxContextSize: 3000,
        priorityLevel: 'medium',
        format: 'markdown',
        includeCode: true
      }
    };
  },
  async created() {
    await this.loadDocuments();
  },
  methods: {
    async loadDocuments() {
      try {
        const response = await documentService.getDocuments();
        if (response.code === 200) {
          this.documents = response.data.items;
        } else {
          this.error = response.message || '加载文档列表失败';
        }
      } catch (err) {
        this.error = '加载文档列表失败: ' + err.message;
      }
    },
    
    async onDocumentChange() {
      if (this.selectedDocumentId) {
        await this.loadDocumentVersions();
      } else {
        this.documentVersions = [];
        this.selectedVersion = '';
      }
    },
    
    async loadDocumentVersions() {
      try {
        this.loading = true;
        const response = await documentService.getDocumentVersions(this.selectedDocumentId);
        if (response.code === 200) {
          this.documentVersions = response.data;
        } else {
          this.error = response.message || '加载文档版本失败';
        }
      } catch (err) {
        this.error = '加载文档版本失败: ' + err.message;
      } finally {
        this.loading = false;
      }
    },
    
    onVersionChange() {
      // 清除之前的结果
      this.result = null;
      this.error = null;
    },
    
    async getStructuredContent() {
      if (!this.selectedDocumentId || !this.selectedVersion) {
        this.error = '请先选择文档和版本';
        return;
      }

      try {
        this.loading = true;
        this.error = null;
        
        const response = await documentService.getStructuredContent(
          this.selectedDocumentId, 
          this.selectedVersion
        );
        
        if (response.code === 200) {
          this.result = response.data;
          this.resultType = 'structured';
          this.resultTitle = '结构化内容';
        } else {
          this.error = response.message || '获取结构化内容失败';
        }
      } catch (err) {
        this.error = '获取结构化内容失败: ' + err.message;
      } finally {
        this.loading = false;
      }
    },
    
    async generateLLMFormat() {
      if (!this.selectedDocumentId || !this.selectedVersion) {
        this.error = '请先选择文档和版本';
        return;
      }

      try {
        this.loading = true;
        this.error = null;
        
        const response = await documentService.generateLLMFormat(
          this.selectedDocumentId, 
          this.selectedVersion,
          this.llmOptions
        );
        
        if (response.code === 200) {
          this.result = response.data;
          this.resultType = 'json';
          this.resultTitle = 'LLM优化格式';
        } else {
          this.error = response.message || '生成LLM优化格式失败';
        }
      } catch (err) {
        this.error = '生成LLM优化格式失败: ' + err.message;
      } finally {
        this.loading = false;
      }
    },
    
    async generateMultiGranularity() {
      if (!this.selectedDocumentId || !this.selectedVersion) {
        this.error = '请先选择文档和版本';
        return;
      }

      try {
        this.loading = true;
        this.error = null;
        
        const response = await documentService.generateMultiGranularity(
          this.selectedDocumentId, 
          this.selectedVersion
        );
        
        if (response.code === 200) {
          this.result = response.data;
          this.resultType = 'json';
          this.resultTitle = '多粒度表示';
        } else {
          this.error = response.message || '生成多粒度表示失败';
        }
      } catch (err) {
        this.error = '生成多粒度表示失败: ' + err.message;
      } finally {
        this.loading = false;
      }
    },
    
    async injectContext() {
      if (!this.selectedDocumentId || !this.selectedVersion) {
        this.error = '请先选择文档和版本';
        return;
      }

      if (!this.contextOptions.query) {
        this.error = '请输入查询内容';
        return;
      }

      try {
        this.loading = true;
        this.error = null;
        
        const response = await documentService.injectContext(
          this.selectedDocumentId, 
          this.selectedVersion,
          this.contextOptions
        );
        
        if (response.code === 200) {
          this.result = response.data.formattedContext;
          this.resultType = 'string';
          this.resultTitle = '上下文注入结果';
        } else {
          this.error = response.message || '注入上下文失败';
        }
      } catch (err) {
        this.error = '注入上下文失败: ' + err.message;
      } finally {
        this.loading = false;
      }
    },
    
    copyResult() {
      if (!this.result) return;
      
      let textToCopy;
      if (typeof this.result === 'string') {
        textToCopy = this.result;
      } else {
        textToCopy = JSON.stringify(this.result, null, 2);
      }
      
      navigator.clipboard.writeText(textToCopy)
        .then(() => {
          // 可以显示一个提示消息
          const originalText = event.target.textContent;
          event.target.textContent = '✅ 已复制';
          setTimeout(() => {
            event.target.textContent = originalText;
          }, 2000);
        })
        .catch(err => {
          console.error('复制失败:', err);
        });
    },
    
    formatTime(timeString) {
      if (!timeString) return '';
      const date = new Date(timeString);
      return date.toLocaleString();
    }
  }
};
</script>

<style scoped>
.ai-format-view {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.header {
  margin-bottom: 30px;
}

.header h1 {
  color: #2c3e50;
  margin-bottom: 10px;
}

.header p {
  color: #666;
  font-size: 16px;
}

.content {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.document-selector,
.format-options {
  background: #f9f9f9;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.document-selector h2,
.format-options h2 {
  margin-top: 0;
  margin-bottom: 20px;
  color: #2c3e50;
}

.selection-controls {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.form-group textarea {
  resize: vertical;
  min-height: 100px;
}

.option-tabs {
  display: flex;
  margin-bottom: 20px;
  border-bottom: 1px solid #ddd;
}

.tab-button {
  padding: 10px 15px;
  border: none;
  background: none;
  cursor: pointer;
  font-size: 16px;
  border-bottom: 3px solid transparent;
  transition: all 0.3s;
}

.tab-button:hover {
  background: #f0f0f0;
}

.tab-button.active {
  border-bottom-color: #3498db;
  color: #3498db;
}

.tab-content {
  padding-top: 20px;
}

.action-button {
  background: #3498db;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
  transition: background 0.3s;
}

.action-button:hover:not(:disabled) {
  background: #2980b9;
}

.action-button:disabled {
  background: #bdc3c7;
  cursor: not-allowed;
}

.radio-group {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
}

.radio-group label {
  display: flex;
  align-items: center;
  gap: 5px;
  cursor: pointer;
}

.result-section {
  background: #f9f9f9;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.result-header h2 {
  margin: 0;
  color: #2c3e50;
}

.copy-button {
  background: #27ae60;
  color: white;
  border: none;
  padding: 8px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background 0.3s;
}

.copy-button:hover {
  background: #219a52;
}

.result-content {
  max-height: 500px;
  overflow-y: auto;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 15px;
  background: white;
}

.result-content pre {
  white-space: pre-wrap;
  word-wrap: break-word;
  font-family: 'Courier New', Courier, monospace;
  font-size: 14px;
  line-height: 1.5;
}

.structured-result {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.section {
  border-bottom: 1px solid #eee;
  padding-bottom: 15px;
}

.section:last-child {
  border-bottom: none;
}

.section h3 {
  margin-top: 0;
  color: #2c3e50;
}

.section h4 {
  color: #34495e;
  margin-top: 15px;
  margin-bottom: 10px;
}

.code-example {
  margin-top: 15px;
}

.code-example pre {
  background: #f5f5f5;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 10px;
  overflow-x: auto;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #3498db;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 15px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error {
  background: #f8d7da;
  color: #721c24;
  padding: 15px;
  border-radius: 4px;
  border: 1px solid #f5c6cb;
}

.error h3 {
  margin-top: 0;
  margin-bottom: 10px;
}

.error p {
  margin: 0;
}

@media (max-width: 768px) {
  .selection-controls {
    flex-direction: column;
    gap: 10px;
  }
  
  .option-tabs {
    flex-wrap: wrap;
  }
  
  .tab-button {
    flex: 1;
    min-width: 120px;
    text-align: center;
  }
  
  .radio-group {
    flex-direction: column;
    gap: 10px;
  }
}
</style>