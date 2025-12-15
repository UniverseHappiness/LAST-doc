-- 重置管理员账户脚本
-- 这个脚本会删除现有的admin账户，以便应用程序可以重新创建它

-- 删除现有的admin账户
DELETE FROM users WHERE username = 'admin';

-- 显示删除结果
DO $$
BEGIN
    IF FOUND THEN
        RAISE NOTICE '管理员账户已删除，请重启应用程序以创建新的管理员账户';
    ELSE
        RAISE NOTICE '未找到管理员账户';
    END IF;
END $$;