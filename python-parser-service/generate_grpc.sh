#!/bin/bash

# 生成Python gRPC代码
echo "生成Python gRPC代码..."

python -m grpc_tools.protoc \
    --proto_path=proto \
    --python_out=. \
    --grpc_python_out=. \
    proto/document_parser.proto

echo "Python gRPC代码生成完成"