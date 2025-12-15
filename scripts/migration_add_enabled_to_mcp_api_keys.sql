-- 为mcp_api_keys表添加enabled字段
ALTER TABLE mcp_api_keys ADD COLUMN IF NOT EXISTS enabled BOOLEAN DEFAULT true;

-- 为现有记录设置默认值
UPDATE mcp_api_keys SET enabled = true WHERE enabled IS NULL;

-- 添加索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_mcp_api_keys_enabled ON mcp_api_keys(enabled);