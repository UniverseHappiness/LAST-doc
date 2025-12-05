# 多阶段构建 - 前端构建阶段
FROM node:18-alpine AS frontend-builder

# 设置工作目录
WORKDIR /app

# 复制package.json和package-lock.json
COPY package*.json ./

# 安装依赖
RUN npm ci

# 复制前端源代码
COPY web/ ./web/

# 构建前端
RUN npm run build

# 多阶段构建 - 后端构建阶段
FROM golang:1.23-alpine AS backend-builder

# 设置工作目录
WORKDIR /app

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# 运行阶段
FROM alpine:latest

# 安装ca证书和bash
RUN apk --no-cache add ca-certificates bash

# 创建应用用户
RUN addgroup -g 1000 appgroup && adduser -u 1000 -G appgroup -s /bin/sh -D appuser

# 创建存储目录
RUN mkdir -p /app/storage && chown -R appuser:appgroup /app

# 设置工作目录
WORKDIR /app

# 从构建阶段复制前端构建文件
COPY --from=frontend-builder /app/web/dist ./web/dist

# 从构建阶段复制后端二进制文件
COPY --from=backend-builder /app/main .

# 更改文件所有权
RUN chown -R appuser:appgroup /app

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]