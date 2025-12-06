#!/bin/bash

# 启动Python解析服务

echo "启动Python解析服务..."

# 激活虚拟环境
source venv/bin/activate

# 设置PYTHONPATH
export PYTHONPATH=$PYTHONPATH:$(pwd)

# 启动gRPC服务器
python -m service.server

echo "Python解析服务已启动"