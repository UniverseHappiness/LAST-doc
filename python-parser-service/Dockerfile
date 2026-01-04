FROM python:3.11-slim

# 设置工作目录
WORKDIR /app

# 使用国内镜像源加速包下载
RUN pip config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple

# 复制依赖文件
COPY requirements.txt .

# 安装 Python 依赖（不安装系统依赖）
RUN pip install --no-cache-dir -r requirements.txt

# 复制服务文件
COPY service/ ./service/
COPY proto/ ./proto/

# 生成 gRPC 代码
RUN python -m grpc_tools.protoc \
    -I./proto \
    --python_out=./ \
    --grpc_python_out=./ \
    proto/document_parser.proto

# 暴露端口
EXPOSE 50051

# 启动服务
CMD ["python", "-m", "service.server"]