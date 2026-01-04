# 多阶段构建 - 前端构建阶段
FROM node:20-alpine AS frontend-builder

# 设置工作目录
WORKDIR /app

# 复制package.json和package-lock.json
COPY package*.json ./

# 复制前端源代码（包括vite.config.js）
COPY web/ ./web/

# 复制package.json到web目录（Vite需要）
RUN cp package.json web/package.json
COPY package-lock.json web/package-lock.json

# 在web目录下构建前端
WORKDIR /app/web
# 安装前端依赖
RUN npm install
RUN npm run build

# 清理缓存以减小镜像大小
RUN npm cache clean --force

# 返回/app目录
WORKDIR /app

# 多阶段构建 - 后端构建阶段
FROM golang:1.24-alpine AS backend-builder

# 设置工作目录
WORKDIR /app

# 配置Go模块代理（支持国内镜像）
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct
ENV GO111MODULE=on

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖（添加重试机制）
RUN go mod download || \
    (echo "首次下载失败，重试中..." && sleep 2 && go mod download) || \
    (echo "第二次重试..." && sleep 3 && go mod download) || \
    echo "警告：Go模块下载失败，依赖vendor目录"

# 复制源代码和可能的vendor目录
COPY . .

# 如果vendor目录存在，使用vendor模式构建
RUN if [ -d "vendor" ]; then \
        echo "使用vendor模式构建..."; \
        CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -mod=vendor -o main cmd/main.go; \
    else \
        echo "使用模块模式构建..."; \
        CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go; \
    fi

# 运行阶段
FROM alpine:latest

# 设置时区
ENV TZ=Asia/Shanghai

# 安装必要的依赖（使用阿里云镜像源加速）
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --no-cache add ca-certificates bash tzdata curl openssl netcat-openbsd && \
    cp /usr/share/zoneinfo/${TZ} /etc/localtime && \
    echo "${TZ}" > /etc/timezone && \
    apk del tzdata

# 创建应用用户
RUN addgroup -g 1000 appgroup && adduser -u 1000 -G appgroup -s /bin/sh -D appuser

# 创建存储目录
RUN mkdir -p /app/storage /app/backups /app/scripts && \
    chown -R appuser:appgroup /app

# 设置工作目录
WORKDIR /app

# 从构建阶段复制前端构建文件
COPY --from=frontend-builder /app/web/dist ./web/dist

# 从构建阶段复制后端二进制文件
COPY --from=backend-builder /app/main .

# 添加健康检查脚本
RUN echo '#!/bin/sh' > /app/scripts/healthcheck.sh && \
    echo 'curl -f http://localhost:${SERVER_PORT:-8080}/health/live || exit 1' >> /app/scripts/healthcheck.sh && \
    chmod +x /app/scripts/healthcheck.sh

# 添加启动脚本
RUN echo '#!/bin/sh' > /app/scripts/entrypoint.sh && \
    echo 'set -e' >> /app/scripts/entrypoint.sh && \
    echo 'echo "Starting AI Document Library Server..."' >> /app/scripts/entrypoint.sh && \
    echo 'echo "Server Port: ${SERVER_PORT:-8080}"' >> /app/scripts/entrypoint.sh && \
    echo 'echo "Storage Directory: ${STORAGE_DIR:-/app/storage}"' >> /app/scripts/entrypoint.sh && \
    echo 'echo "Database Host: ${DB_HOST:-localhost}"' >> /app/scripts/entrypoint.sh && \
    echo '' >> /app/scripts/entrypoint.sh && \
    echo '# 等待数据库就绪' >> /app/scripts/entrypoint.sh && \
    echo 'if [ -n "${DB_HOST}" ] && [ "${DB_HOST}" != "localhost" ]; then' >> /app/scripts/entrypoint.sh && \
    echo '  echo "Waiting for database connection..."' >> /app/scripts/entrypoint.sh && \
    echo '  timeout 60 sh -c "until nc -z ${DB_HOST} ${DB_PORT:-5432}; do sleep 1; done" || true' >> /app/scripts/entrypoint.sh && \
    echo '  echo "Database is ready!"' >> /app/scripts/entrypoint.sh && \
    echo 'fi' >> /app/scripts/entrypoint.sh && \
    echo '' >> /app/scripts/entrypoint.sh && \
    echo '# 启动应用' >> /app/scripts/entrypoint.sh && \
    echo 'echo "Starting application..."' >> /app/scripts/entrypoint.sh && \
    echo 'exec ./main' >> /app/scripts/entrypoint.sh && \
    chmod +x /app/scripts/entrypoint.sh

# 更改文件所有权
RUN chown -R appuser:appgroup /app

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
    CMD ["/bin/sh", "/app/scripts/healthcheck.sh"]

# 启动应用
CMD ["/bin/sh", "/app/scripts/entrypoint.sh"]