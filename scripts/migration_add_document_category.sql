-- 添加文档分类字段（如果不存在）
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'documents' AND column_name = 'category'
    ) THEN
        ALTER TABLE documents ADD COLUMN category TEXT;
    END IF;
END $$;

-- 为现有文档设置默认分类
UPDATE documents SET category = 'document' WHERE category IS NULL;

-- 设置分类字段为非空
ALTER TABLE documents ALTER COLUMN category SET NOT NULL;

-- 为分类字段创建索引
CREATE INDEX IF NOT EXISTS idx_documents_category ON documents(category);
$$;