-- 清理孤立的搜索索引脚本
-- 孤立索引：在search_indices表中存在，但在document_versions表中不存在对应的记录

-- ============================================
-- 第一步：查看孤立的搜索索引
-- ============================================
SELECT 
    'Orphaned Search Indices' as analysis_type,
    si.id as search_index_id,
    si.document_id,
    si.version,
    si.content_type,
    SUBSTRING(si.content, 1, 100) as content_preview,
    si.created_at
FROM search_indices si
LEFT JOIN document_versions dv 
    ON si.document_id = dv.document_id 
    AND si.version = dv.version
WHERE dv.id IS NULL;

-- 统计孤立索引的数量
SELECT 
    COUNT(*) as orphaned_indices_count
FROM search_indices si
LEFT JOIN document_versions dv 
    ON si.document_id = dv.document_id 
    AND si.version = dv.version
WHERE dv.id IS NULL;

-- ============================================
-- 第二步：按document_id分组查看孤立索引
-- ============================================
SELECT 
    si.document_id,
    si.version,
    COUNT(*) as orphaned_indices_count,
    MIN(si.created_at) as first_created,
    MAX(si.created_at) as last_created
FROM search_indices si
LEFT JOIN document_versions dv 
    ON si.document_id = dv.document_id 
    AND si.version = dv.version
WHERE dv.id IS NULL
GROUP BY si.document_id, si.version
ORDER BY orphaned_indices_count DESC;

-- ============================================
-- 第三步：检查document_versions表中是否有该document_id的任何版本
-- ============================================
SELECT 
    si.id as search_index_id,
    si.document_id,
    si.version,
    CASE 
        WHEN EXISTS (SELECT 1 FROM document_versions WHERE document_id = si.document_id) 
        THEN 'Document exists but version not found'
        ELSE 'Document does not exist at all'
    END as issue_type,
    (SELECT COUNT(*) FROM document_versions WHERE document_id = si.document_id) as existing_versions_count
FROM search_indices si
LEFT JOIN document_versions dv 
    ON si.document_id = dv.document_id 
    AND si.version = dv.version
WHERE dv.id IS NULL
ORDER BY si.document_id, si.version;

-- ============================================
-- 第四步：备份数据（在删除前执行）
-- ============================================
-- 创建备份表
CREATE TABLE IF NOT EXISTS search_indices_backup (
    id TEXT,
    document_id TEXT,
    version TEXT,
    content TEXT,
    content_type TEXT,
    section TEXT,
    keywords TEXT,
    vector TEXT,
    metadata TEXT,
    score REAL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    backup_note TEXT
);

-- 备份孤立的索引
INSERT INTO search_indices_backup 
SELECT 
    *, 
    'Orphaned index cleaned on ' || NOW() as backup_note
FROM search_indices si
LEFT JOIN document_versions dv 
    ON si.document_id = dv.document_id 
    AND si.version = dv.version
WHERE dv.id IS NULL;

-- 验证备份数据
SELECT COUNT(*) as back_up_count FROM search_indices_backup 
WHERE backup_note LIKE '%Orphaned index cleaned%';

-- ============================================
-- 第五步：删除孤立的搜索索引
-- ============================================
-- 删除孤立的索引
DELETE FROM search_indices 
WHERE id IN (
    SELECT si.id
    FROM search_indices si
    LEFT JOIN document_versions dv 
        ON si.document_id = dv.document_id 
        AND si.version = dv.version
    WHERE dv.id IS NULL
);

-- ============================================
-- 第六步：验证清理结果
-- ============================================
-- 再次检查孤立索引（应该为0）
SELECT 
    COUNT(*) as remaining_orphaned_indices_count
FROM search_indices si
LEFT JOIN document_versions dv 
    ON si.document_id = dv.document_id 
    AND si.version = dv.version
WHERE dv.id IS NULL;

-- 显示当前search_indices表的总数
SELECT 'Total search indices after cleanup' as status, COUNT(*) as count FROM search_indices;

-- ============================================
-- 第七步：清理备份表（可选，确定不需要备份数据后执行）
-- ============================================
-- DROP TABLE IF EXISTS search_indices_backup;