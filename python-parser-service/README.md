# Python Parser Service

这是一个用于解析PDF和DOCX文档的Python服务，通过gRPC与主Go应用程序通信。

## 功能特性

- PDF文档解析（提取文本内容、元数据）
- DOCX文档解析（提取文本内容、元数据）
- gRPC服务接口
- 异步处理支持

## 依赖

- Python 3.8+
- PyPDF2
- python-docx
- grpcio
- grpcio-tools

## 运行

```bash
pip install -r requirements.txt
python -m service.server