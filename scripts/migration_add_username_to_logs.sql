-- 添加username字段到log_entries表
-- 执行日期: 2026-01-03

-- 检查字段是否已存在，避免重复执行
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'log_entries' 
        AND column_name = 'username'
    ) THEN
        ALTER TABLE log_entries ADD COLUMN username VARCHAR(255);
        RAISE NOTICE 'Column "username" added to table "log_entries"';
    ELSE
        RAISE NOTICE 'Column "username" already exists in table "log_entries"';
    END IF;
END
$$;

-- 为username字段创建索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_log_entries_username ON log_entries(username);

-- 对于现有的记录，尝试根据user_id更新username（如果存在users表）
-- 这个操作可能需要时间，取决于数据量
DO $$
BEGIN
    -- 检查users表是否存在
    IF EXISTS (
        SELECT 1 
        FROM information_schema.tables 
        WHERE table_name = 'users'
    ) THEN
        -- 更新现有记录的username
        UPDATE log_entries l
        SET username = u.username
        FROM users u
        WHERE l.user_id = u.id::text AND l.username IS NULL OR l.username = '';
        
        RAISE NOTICE 'Updated username for existing log entries based on users table';
    END IF;
END
$$;

COMMENT ON COLUMN log_entries.username IS '用户名，用于日志显示';