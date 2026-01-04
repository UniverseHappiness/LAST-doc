package service

import (
	"context"
	"fmt"
	"os/exec"
)

// PostgreSQLBackup PostgreSQL数据库备份实现
type PostgreSQLBackup struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
}

// NewPostgreSQLBackup 创建PostgreSQL备份服务
func NewPostgreSQLBackup(host, port, user, password, dbname string) *PostgreSQLBackup {
	return &PostgreSQLBackup{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		dbname:   dbname,
	}
}

// DumpDatabase 备份数据库
func (p *PostgreSQLBackup) DumpDatabase(ctx context.Context, destFile string) error {
	args := []string{
		"-h", p.host,
		"-p", p.port,
		"-U", p.user,
		"-d", p.dbname,
		"-f", destFile,
		"--no-password",
	}

	cmd := exec.CommandContext(ctx, "pg_dump", args...)

	// 设置环境变量以避免密码提示
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", p.password))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pg_dump failed: %v, output: %s", err, string(output))
	}

	return nil
}

// RestoreDatabase 恢复数据库
func (p *PostgreSQLBackup) RestoreDatabase(ctx context.Context, sourceFile string) error {
	args := []string{
		"-h", p.host,
		"-p", p.port,
		"-U", p.user,
		"-d", p.dbname,
		"--no-password",
	}

	// 使用shell命令执行恢复
	args = append(args, "-f", sourceFile)
	cmd := exec.CommandContext(ctx, "psql", args...)

	// 设置环境变量以避免密码提示
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", p.password))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("psql restore failed: %v, output: %s", err, string(output))
	}

	return nil
}
