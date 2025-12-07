-- 修复搜索索引表vector列类型的迁移脚本

-- 步骤1：创建临时列，使用正确的类型
ALTER TABLE search_indices ADD COLUMN IF NOT EXISTS vector_new JSONB;

-- 步骤2：将现有数据从旧列迁移到新列
UPDATE search_indices SET vector_new = CASE 
    WHEN vector IS NULL THEN '[]'::jsonb
    ELSE 
        -- 将float数组转换为JSONB数组
        (
            SELECT jsonb_agg(elem::float)
            FROM (
                SELECT unnest(vector) as elem
            ) t
        )
    END
WHERE vector_new IS NULL;

-- 步骤3：删除旧列
ALTER TABLE search_indices DROP COLUMN IF EXISTS vector;

-- 步骤4：重命名新列为原来的名称
ALTER TABLE search_indices RENAME COLUMN vector_new TO vector;

-- 步骤5：添加默认值
ALTER TABLE search_indices ALTER COLUMN vector SET DEFAULT '[]'::jsonb;

-- 步骤6：添加注释
COMMENT ON COLUMN search_indices.vector IS '语义向量，JSONB格式存储';