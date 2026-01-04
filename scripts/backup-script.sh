#!/bin/bash

# 自动化数据库备份脚本
# 支持定时备份和备份清理

set -e

# 配置变量
BACKUP_DIR=${BACKUP_DIR:-"/app/backups"}
DB_HOST=${DB_HOST:-"postgres"}
DB_PORT=${DB_PORT:-"5432"}
DB_USER=${DB_USER:-"postgres"}
DB_PASSWORD=${DB_PASSWORD:-"postgres"}
DB_NAME=${DB_NAME:-"ai_doc_library"}
RETENTION_DAYS=${RETENTION_DAYS:-7}  # 保留最近7天的备份

# 创建备份目录
mkdir -p "${BACKUP_DIR}/database"
mkdir -p "${BACKUP_DIR}/storage"

# 生成时间戳
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# 数据库备份函数
backup_database() {
    echo "开始备份数据库..."
    
    BACKUP_FILE="${BACKUP_DIR}/database/backup_${TIMESTAMP}.sql"
    
    PGPASSWORD="${DB_PASSWORD}" pg_dump \
        -h "${DB_HOST}" \
        -p "${DB_PORT}" \
        -U "${DB_USER}" \
        -d "${DB_NAME}" \
        -f "${BACKUP_FILE}" \
        --no-password
    
    if [ $? -eq 0 ]; then
        echo "数据库备份成功: ${BACKUP_FILE}"
        
        # 压缩备份文件
        gzip "${BACKUP_FILE}"
        echo "备份文件已压缩: ${BACKUP_FILE}.gz"
        
        # 记录备份信息
        echo "$(date +"%Y-%m-%d %H:%M:%S") | Database backup created: ${BACKUP_FILE}.gz" >> "${BACKUP_DIR}/backup.log"
    else
        echo "数据库备份失败"
        exit 1
    fi
}

# 清理旧备份函数
cleanup_old_backups() {
    echo "清理 ${RETENTION_DAYS} 天前的备份..."
    
    find "${BACKUP_DIR}/database" -type f -name "backup_*.sql.gz" -mtime +${RETENTION_DAYS} -delete
    find "${BACKUP_DIR}/storage" -type d -mtime +${RETENTION_DAYS} -exec rm -rf {} + 2>/dev/null || true
    
    echo "旧备份清理完成"
}

# 创建完整备份
create_full_backup() {
    echo "创建完整备份..."
    
    # 备份数据库
    backup_database
    
    # 创建存储备份标记
    FULL_BACKUP_DIR="${BACKUP_DIR}/full/backup_${TIMESTAMP}"
    mkdir -p "${FULL_BACKUP_DIR}"
    
    echo "Backup created at: $(date +"%Y-%m-%d %H:%M:%S")" > "${FULL_BACKUP_DIR}/backup_info.txt"
    echo "Storage type: ${STORAGE_TYPE:-local}" >> "${FULL_BACKUP_DIR}/backup_info.txt"
    
    # 复制数据库备份
    cp "${BACKUP_DIR}/database/backup_${TIMESTAMP}.sql.gz" "${FULL_BACKUP_DIR}/"
    
    echo "完整备份创建完成: ${FULL_BACKUP_DIR}"
    
    # 清理旧备份
    cleanup_old_backups
}

# 恢复备份函数
restore_backup() {
    BACKUP_FILE=$1
    
    if [ -z "${BACKUP_FILE}" ]; then
        echo "错误: 请指定备份文件路径"
        echo "使用方法: $0 restore <backup_file>"
        exit 1
    fi
    
    if [ ! -f "${BACKUP_FILE}" ]; then
        echo "错误: 备份文件不存在: ${BACKUP_FILE}"
        exit 1
    fi
    
    echo "开始恢复数据库备份: ${BACKUP_FILE}"
    
    # 解压备份文件
    TEMP_SQL="${BACKUP_FILE%.gz}"
    gunzip -c "${BACKUP_FILE}" > "${TEMP_SQL}"
    
    # 恢复数据库
    PGPASSWORD="${DB_PASSWORD}" psql \
        -h "${DB_HOST}" \
        -p "${DB_PORT}" \
        -U "${DB_USER}" \
        -d "${DB_NAME}" \
        -f "${TEMP_SQL}"
    
    if [ $? -eq 0 ]; then
        echo "数据库恢复成功"
        
        # 清理临时文件
        rm -f "${TEMP_SQL}"
        
        echo "$(date +"%Y-%m-%d %H:%M:%S") | Database restore completed from: ${BACKUP_FILE}" >> "${BACKUP_DIR}/backup.log"
    else
        echo "数据库恢复失败"
        rm -f "${TEMP_SQL}"
        exit 1
    fi
}

# 主函数
main() {
    case "${1:-backup}" in
        backup)
            backup_database
            cleanup_old_backups
            ;;
        full)
            create_full_backup
            ;;
        restore)
            restore_backup "${2}"
            ;;
        cleanup)
            cleanup_old_backups
            ;;
        *)
            echo "使用方法:"
            echo "  $0 backup        - 备份数据库"
            echo "  $0 full          - 创建完整备份"
            echo "  $0 restore <file> - 恢复备份"
            echo "  $0 cleanup       - 清理旧备份"
            exit 1
            ;;
    esac
}

main "$@"