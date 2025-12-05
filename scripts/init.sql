-- AI文档库数据库初始化脚本

-- 创建数据库（如果不存在）
-- 注意：在Docker环境中，数据库已经在docker-compose.yml中定义，这里不需要创建

-- 创建扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 插入测试数据（可选）
-- INSERT INTO documents (id, name, type, version, library, description, tags, status, created_at, updated_at) 
-- VALUES (
--     uuid_generate_v4(),
--     '示例文档',
--     'markdown',
--     '1.0.0',
--     '示例库',
--     '这是一个示例文档',
--     ARRAY['示例', '测试'],
--     'completed',
--     NOW(),
--     NOW()
-- );

-- 验证表是否创建成功
-- SELECT * FROM documents LIMIT 1;