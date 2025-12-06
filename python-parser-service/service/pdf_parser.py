import os
import logging
from typing import Dict, Any, Tuple
import PyPDF2
import pdfplumber

logger = logging.getLogger(__name__)

class PDFParser:
    """PDF文档解析器"""
    
    def __init__(self):
        self.logger = logging.getLogger(__name__)
    
    def parse(self, file_path: str) -> Tuple[str, Dict[str, Any]]:
        """
        解析PDF文档
        
        Args:
            file_path: PDF文件路径
            
        Returns:
            Tuple[str, Dict[str, Any]]: (文本内容, 元数据)
        """
        try:
            self.logger.info(f"开始解析PDF文档: {file_path}")
            
            # 检查文件是否存在
            if not os.path.exists(file_path):
                raise FileNotFoundError(f"PDF文件不存在: {file_path}")
            
            # 获取文件大小
            file_size = os.path.getsize(file_path)
            self.logger.info(f"PDF文件大小: {file_size} 字节")
            
            # 使用PyPDF2提取文本
            content_text = ""
            page_count = 0
            
            with open(file_path, 'rb') as file:
                pdf_reader = PyPDF2.PdfReader(file)
                page_count = len(pdf_reader.pages)
                
                for page_num, page in enumerate(pdf_reader.pages):
                    try:
                        page_text = page.extract_text()
                        # 检查提取的文本是否包含非UTF-8字符
                        try:
                            page_text.encode('utf-8').decode('utf-8')
                            content_text += page_text + "\n"
                            self.logger.debug(f"已解析第 {page_num + 1} 页，文本长度: {len(page_text)}")
                        except UnicodeError as e:
                            self.logger.warning(f"第 {page_num + 1} 页包含非UTF-8字符: {e}")
                            # 尝试清理文本
                            clean_text = page_text.encode('utf-8', errors='replace').decode('utf-8', errors='replace')
                            content_text += clean_text + "\n"
                    except Exception as e:
                        self.logger.warning(f"解析第 {page_num + 1} 页时出错: {e}")
                        continue
            
            # 使用pdfplumber提取更详细的元数据
            metadata = self._extract_metadata_with_pdfplumber(file_path)
            
            # 添加基本元数据
            metadata.update({
                'file_size': str(file_size),
                'page_count': str(page_count),
                'content_length': str(len(content_text)),
                'parser': 'PyPDF2 + pdfplumber'
            })
            
            self.logger.info(f"PDF文档解析完成，总页数: {page_count}，内容长度: {len(content_text)}")
            
            return content_text, metadata
            
        except Exception as e:
            self.logger.error(f"解析PDF文档时发生错误: {e}")
            raise
    
    def _extract_metadata_with_pdfplumber(self, file_path: str) -> Dict[str, Any]:
        """使用pdfplumber提取元数据"""
        metadata = {}
        
        try:
            with pdfplumber.open(file_path) as pdf:
                # 提取文档元数据
                if pdf.metadata:
                    metadata.update({
                        'title': pdf.metadata.get('Title', ''),
                        'author': pdf.metadata.get('Author', ''),
                        'subject': pdf.metadata.get('Subject', ''),
                        'creator': pdf.metadata.get('Creator', ''),
                        'producer': pdf.metadata.get('Producer', ''),
                        'creation_date': str(pdf.metadata.get('CreationDate', '')),
                        'modification_date': str(pdf.metadata.get('ModDate', ''))
                    })
                
                # 统计表格数量
                table_count = 0
                for page in pdf.pages:
                    tables = page.extract_tables()
                    table_count += len(tables)
                
                metadata['table_count'] = str(table_count)
                
                # 统计图片数量（近似）
                image_count = 0
                for page in pdf.pages:
                    if page.images:
                        image_count += len(page.images)
                
                metadata['image_count'] = str(image_count)
                
        except Exception as e:
            self.logger.warning(f"使用pdfplumber提取元数据时出错: {e}")
        
        return metadata