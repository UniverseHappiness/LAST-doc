import logging
import grpc
import os
import sys
import locale
from concurrent import futures
from typing import Dict, Any

# 添加当前目录到Python路径
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

# 导入生成的gRPC代码
try:
    import document_parser_pb2 as pb2
    import document_parser_pb2_grpc as pb2_grpc
except ImportError:
    print("错误：请先生成gRPC代码")
    sys.exit(1)

# 导入解析器
from service.pdf_parser import PDFParser
from service.docx_parser import DOCXParser

# 配置日志 - 确保使用UTF-8编码
try:
    locale.setlocale(locale.LC_ALL, 'en_US.UTF-8')
except locale.Error:
    try:
        locale.setlocale(locale.LC_ALL, 'C.UTF-8')
    except locale.Error:
        pass

# 确保标准输出使用UTF-8编码
if hasattr(sys.stdout, 'reconfigure'):
    sys.stdout.reconfigure(encoding='utf-8')
if hasattr(sys.stderr, 'reconfigure'):
    sys.stderr.reconfigure(encoding='utf-8')

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.StreamHandler(sys.stdout),
        logging.FileHandler('parser_server.log', encoding='utf-8')
    ]
)
logger = logging.getLogger(__name__)

class DocumentParserServicer(pb2_grpc.DocumentParserServiceServicer):
    """文档解析服务实现"""
    
    def __init__(self):
        self.pdf_parser = PDFParser()
        self.docx_parser = DOCXParser()
        logger.info("文档解析服务初始化完成")
    
    def ParsePDF(self, request: pb2.ParsePDFRequest, context: grpc.ServicerContext) -> pb2.ParseDocumentResponse:
        """解析PDF文档"""
        try:
            # 添加工作目录诊断信息
            current_dir = os.getcwd()
            logger.info(f"收到PDF解析请求 - Python工作目录: {current_dir}")
            logger.info(f"收到PDF解析请求 - 文件路径: {request.file_path}")
            
            # 检查文件是否存在
            if not request.file_path:
                error_msg = "文件路径不能为空"
                logger.error(error_msg)
                return pb2.ParseDocumentResponse(
                    success=False,
                    error_message=error_msg
                )
            
            # 检查文件是否存在，并添加详细诊断信息
            file_exists = os.path.exists(request.file_path)
            logger.info(f"文件存在性检查 - 路径: {request.file_path}, 存在: {file_exists}")
            
            if not file_exists:
                # 添加更多诊断信息
                dir_exists = os.path.exists(os.path.dirname(request.file_path))
                parent_dir = os.path.dirname(os.path.dirname(request.file_path))
                parent_dir_exists = os.path.exists(parent_dir)
                logger.info(f"诊断信息 - 文件目录存在: {dir_exists}, 父目录存在: {parent_dir_exists}")
                logger.info(f"诊断信息 - 当前目录内容: {os.listdir('.')}")
                if dir_exists:
                    logger.info(f"诊断信息 - 文件目录内容: {os.listdir(os.path.dirname(request.file_path))}")
                
                error_msg = f"文件不存在: {request.file_path}"
                logger.error(error_msg)
                return pb2.ParseDocumentResponse(
                    success=False,
                    error_message=error_msg
                )
            
            # 解析PDF文档
            content, metadata = self.pdf_parser.parse(request.file_path)
            
            # 将Python字典转换为gRPC的map
            metadata_map = {str(k): str(v) for k, v in metadata.items()}
            
            logger.info(f"PDF解析成功，内容长度: {len(content)}")
            
            return pb2.ParseDocumentResponse(
                success=True,
                content=content,
                metadata=metadata_map
            )
            
        except Exception as e:
            error_msg = f"PDF解析失败: {str(e)}"
            logger.error(error_msg)
            return pb2.ParseDocumentResponse(
                success=False,
                error_message=error_msg
            )
    
    def ParseDOCX(self, request: pb2.ParseDOCXRequest, context: grpc.ServicerContext) -> pb2.ParseDocumentResponse:
        """解析DOCX文档"""
        try:
            # 添加工作目录诊断信息
            current_dir = os.getcwd()
            logger.info(f"收到DOCX解析请求 - Python工作目录: {current_dir}")
            logger.info(f"收到DOCX解析请求 - 文件路径: {request.file_path}")
            
            # 检查文件是否存在
            if not request.file_path:
                error_msg = "文件路径不能为空"
                logger.error(error_msg)
                return pb2.ParseDocumentResponse(
                    success=False,
                    error_message=error_msg
                )
            
            # 检查文件是否存在，并添加详细诊断信息
            file_exists = os.path.exists(request.file_path)
            logger.info(f"文件存在性检查 - 路径: {request.file_path}, 存在: {file_exists}")
            
            if not file_exists:
                # 添加更多诊断信息
                dir_exists = os.path.exists(os.path.dirname(request.file_path))
                parent_dir = os.path.dirname(os.path.dirname(request.file_path))
                parent_dir_exists = os.path.exists(parent_dir)
                logger.info(f"诊断信息 - 文件目录存在: {dir_exists}, 父目录存在: {parent_dir_exists}")
                logger.info(f"诊断信息 - 当前目录内容: {os.listdir('.')}")
                if dir_exists:
                    logger.info(f"诊断信息 - 文件目录内容: {os.listdir(os.path.dirname(request.file_path))}")
                
                error_msg = f"文件不存在: {request.file_path}"
                logger.error(error_msg)
                return pb2.ParseDocumentResponse(
                    success=False,
                    error_message=error_msg
                )
            
            # 解析DOCX文档
            content, metadata = self.docx_parser.parse(request.file_path)
            
            # 将Python字典转换为gRPC的map
            metadata_map = {str(k): str(v) for k, v in metadata.items()}
            
            logger.info(f"DOCX解析成功，内容长度: {len(content)}")
            
            return pb2.ParseDocumentResponse(
                success=True,
                content=content,
                metadata=metadata_map
            )
            
        except Exception as e:
            error_msg = f"DOCX解析失败: {str(e)}"
            logger.error(error_msg)
            return pb2.ParseDocumentResponse(
                success=False,
                error_message=error_msg
            )
    
    def HealthCheck(self, request: pb2.HealthCheckRequest, context: grpc.ServicerContext) -> pb2.HealthCheckResponse:
        """健康检查"""
        logger.info(f"收到健康检查请求，服务: {request.service}")
        
        try:
            # 检查解析器是否正常
            test_pdf_path = "/tmp/test.pdf"
            test_docx_path = "/tmp/test.docx"
            
            # 这里可以添加更多的健康检查逻辑
            # 例如检查依赖库是否正常，是否有足够的系统资源等
            
            return pb2.HealthCheckResponse(
                healthy=True,
                message="文档解析服务正常运行",
                version="1.0.0"
            )
            
        except Exception as e:
            logger.error(f"健康检查失败: {str(e)}")
            return pb2.HealthCheckResponse(
                healthy=False,
                message=f"健康检查失败: {str(e)}",
                version="1.0.0"
            )

def serve():
    """启动gRPC服务器"""
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    pb2_grpc.add_DocumentParserServiceServicer_to_server(
        DocumentParserServicer(), server
    )
    
    # 监听端口
    port = "50051"
    server.add_insecure_port(f"[::]:{port}")
    
    logger.info(f"启动gRPC服务器，监听端口: {port}")
    server.start()
    
    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        logger.info("收到终止信号，关闭服务器")
        server.stop(0)

if __name__ == "__main__":
    serve()