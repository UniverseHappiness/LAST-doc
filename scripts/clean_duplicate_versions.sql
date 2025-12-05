-- 清理重复的文档版本数据脚本

-- 首先查看重复的版本数据
SELECT 
    document_id, 
    version, 
    COUNT(*) as duplicate_count,
    ARRAY_AGG(id ORDER BY created_at DESC) as version_ids,
    ARRAY_AGG(created_at ORDER BY created_at DESC) as created_ats
FROM document_versions 
GROUP BY document_id, version 
HAVING COUNT(*) > 1;

-- 创建临时表存储需要保留的版本ID（每个document_id+version组合保留最新的一个）
CREATE TEMPORARY TABLE versions_to_keep AS
WITH ranked_versions AS (
    SELECT 
        id,
        document_id,
        version,
        ROW_NUMBER() OVER (PARTITION BY document_id, version ORDER BY created_at DESC) as rn
    FROM document_versions
)
SELECT id FROM ranked_versions WHERE rn = 1;

-- 删除重复的版本（保留最新的）
DELETE FROM document_versions 
WHERE id NOT IN (SELECT id FROM versions_to_keep);

-- 验证清理结果
SELECT 
    document_id, 
    version, 
    COUNT(*) as count_after_cleanup
FROM document_versions 
GROUP BY document_id, version 
HAVING COUNT(*) > 1;

-- 删除临时表
DROP TABLE versions_to_keep;