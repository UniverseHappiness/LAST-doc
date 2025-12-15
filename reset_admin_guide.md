# 重置管理员账户指南

## 方法一：使用psql命令行（如果数据库是PostgreSQL）

```bash
psql -h localhost -U postgres -d ai_doc_library -f scripts/reset_admin_password.sql
```

## 方法二：使用应用程序连接执行

如果应用程序正在运行，可以通过以下方式之一执行脚本：

1. **如果应用程序有数据库管理界面**：
   - 登录数据库管理界面
   - 执行以下SQL命令：
   ```sql
   DELETE FROM users WHERE username = 'admin';
   ```

2. **如果使用Docker**（但我们已经确认不是）：
   ```bash
   docker-compose exec postgres psql -U postgres -d ai_doc_library -f scripts/reset_admin_password.sql
   ```

## 执行脚本后

1. 重启应用程序，让它使用新的密码哈希重新创建管理员账户
2. 使用以下凭据登录：
   - 用户名：admin
   - 密码：admin123
   - 邮箱：admin@example.com

## 注意事项

- 执行脚本后，现有的admin账户将被删除
- 重启应用程序后，新的admin账户将使用正确的密码哈希
- 在生产环境中，请立即更改默认密码