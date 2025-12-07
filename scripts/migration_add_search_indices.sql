-- 创建搜索索引表的迁移脚本

-- 创建扩展（如果不存在）
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 创建搜索索引表
CREATE TABLE IF NOT EXISTS search_indices (
    id VARCHAR(255) PRIMARY KEY,
    document_id VARCHAR(255) NOT NULL,
    version VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    content_type VARCHAR(50) NOT NULL,
    section VARCHAR(255),
    keywords TEXT,
    metadata JSONB,
    score FLOAT DEFAULT 0,
    vector FLOAT[],  -- 使用FLOAT数组类型存储向量
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_search_indices_document_id ON search_indices(document_id);
CREATE INDEX IF NOT EXISTS idx_search_indices_version ON search_indices(version);
CREATE INDEX IF NOT EXISTS idx_search_indices_content_type ON search_indices(content_type);
CREATE INDEX IF NOT EXISTS idx_search_indices_section ON search_indices(section);
CREATE INDEX IF NOT EXISTS idx_search_indices_created_at ON search_indices(created_at);

-- 创建全文搜索索引（PostgreSQL的全文搜索功能）
CREATE INDEX IF NOT EXISTS idx_search_indices_content_fts ON search_indices USING gin(to_tsvector('simple', content));

-- 创建关键词搜索索引（提高关键词搜索性能）
CREATE INDEX IF NOT EXISTS idx_search_indices_keywords ON search_indices USING gin(to_tsvector('simple', keywords));

-- 创建复合索引（提高过滤和排序性能）
CREATE INDEX IF NOT EXISTS idx_search_indices_doc_version ON search_indices(document_id, version);
CREATE INDEX IF NOT EXISTS idx_search_indices_doc_type_score ON search_indices(document_id, content_type, score);
CREATE INDEX IF NOT EXISTS idx_search_indices_created_score ON search_indices(created_at, score);

-- 创建复合唯一约束（确保同一文档版本的同一内容不会重复索引）
CREATE UNIQUE INDEX IF NOT EXISTS idx_search_indices_unique ON search_indices(document_id, version, content_type, section);

-- 添加外键约束
ALTER TABLE search_indices
ADD CONSTRAINT fk_search_indices_document_id
FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE CASCADE;

-- 添加注释
COMMENT ON TABLE search_indices IS '文档搜索索引表，存储文档内容的索引信息用于快速检索';
COMMENT ON COLUMN search_indices.id IS '索引唯一标识';
COMMENT ON COLUMN search_indices.document_id IS '关联的文档ID';
COMMENT ON COLUMN search_indices.version IS '文档版本号';
COMMENT ON COLUMN search_indices.content IS '索引的文档内容';
COMMENT ON COLUMN search_indices.content_type IS '内容类型（text, code等）';
COMMENT ON COLUMN search_indices.section IS '文档章节';
COMMENT ON COLUMN search_indices.keywords IS '关键词，用逗号分隔';
COMMENT ON COLUMN search_indices.metadata IS '额外的元数据，JSON格式';
COMMENT ON COLUMN search_indices.vector IS '向量表示，用于语义搜索';
COMMENT ON COLUMN search_indices.score IS '搜索相关度得分';
COMMENT ON COLUMN search_indices.created_at IS '创建时间';
COMMENT ON COLUMN search_indices.updated_at IS '更新时间';

-- 创建触发器函数，自动更新updated_at字段
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 创建触发器
DROP TRIGGER IF EXISTS update_search_indices_updated_at ON search_indices;
CREATE TRIGGER update_search_indices_updated_at
    BEFORE UPDATE ON search_indices
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();