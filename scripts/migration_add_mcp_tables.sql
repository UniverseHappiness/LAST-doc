-- 创建MCP API密钥表
CREATE TABLE IF NOT EXISTS mcp_api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    key VARCHAR(255) NOT NULL UNIQUE,
    user_id VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NULL,
    last_used TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建MCP配置表
CREATE TABLE IF NOT EXISTS mcp_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    endpoint VARCHAR(255) NOT NULL,
    api_key VARCHAR(255) NOT NULL,
    enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_mcp_api_keys_key ON mcp_api_keys(key);
CREATE INDEX IF NOT EXISTS idx_mcp_api_keys_user_id ON mcp_api_keys(user_id);

-- 创建更新时间触发器函数
CREATE OR REPLACE FUNCTION update_mcp_tables_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 创建触发器
DROP TRIGGER IF EXISTS update_mcp_api_keys_updated_at ON mcp_api_keys;
CREATE TRIGGER update_mcp_api_keys_updated_at
    BEFORE UPDATE ON mcp_api_keys
    FOR EACH ROW
    EXECUTE FUNCTION update_mcp_tables_updated_at();

DROP TRIGGER IF EXISTS update_mcp_configs_updated_at ON mcp_configs;
CREATE TRIGGER update_mcp_configs_updated_at
    BEFORE UPDATE ON mcp_configs
    FOR EACH ROW
    EXECUTE FUNCTION update_mcp_tables_updated_at();

-- 插入默认MCP配置
INSERT INTO mcp_configs (name, description, endpoint, api_key, enabled)
VALUES (
    'AI技术文档库MCP服务',
    'AI技术文档库的MCP协议支持服务',
    'http://localhost:8080/mcp',
    'ai-doc-library-default-key',
    TRUE
) ON CONFLICT DO NOTHING;