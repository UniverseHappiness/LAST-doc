-- 为 document_versions 表添加文件类型字段
-- 这个迁移脚本添加了 type 字段，用于存储文档版本的文件类型

-- 添加字段
ALTER TABLE document_versions 
ADD COLUMN IF NOT EXISTS type TEXT NOT NULL DEFAULT 'markdown';

-- 更新现有记录的 type 字段，根据文件扩展名设置类型
UPDATE document_versions 
SET type = CASE 
    WHEN file_path LIKE '%.pdf' THEN 'pdf'
    WHEN file_path LIKE '%.docx' OR file_path LIKE '%.doc' THEN 'docx'
    WHEN file_path LIKE '%.json' OR file_path LIKE '%.yaml' OR file_path LIKE '%.yml' THEN
        CASE 
            WHEN content LIKE '%swagger%' OR content LIKE '%openapi%' THEN 'swagger'
            ELSE 'openapi'
        END
    WHEN file_path LIKE '%.html' OR file_path LIKE '%.htm' THEN 'java_doc'
    ELSE 'markdown'
END;

-- 添加索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_document_versions_type 
ON document_versions(type);

-- 验证迁移结果
SELECT type, COUNT(*) as count 
FROM document_versions 
GROUP BY type 
ORDER BY count DESC;