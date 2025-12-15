#!/usr/bin/env python3
"""
简单的MCP客户端测试脚本
用于测试AI技术文档库的MCP功能
"""

import requests
import json
import sys
import argparse
from typing import Dict, Any, Optional

class MCPClient:
    """MCP客户端类"""
    
    def __init__(self, base_url: str, api_key: str):
        self.base_url = base_url.rstrip('/')
        self.api_key = api_key
        self.headers = {
            "API_KEY": api_key,
            "Content-Type": "application/json"
        }
    
    def test_connection(self) -> Dict[str, Any]:
        """测试MCP连接"""
        try:
            response = requests.get(
                f"{self.base_url}/api/v1/mcp/test",
                headers=self.headers
            )
            response.raise_for_status()
            return response.json()
        except requests.RequestException as e:
            return {"error": f"连接测试失败: {str(e)}"}
    
    def list_tools(self) -> Dict[str, Any]:
        """获取可用工具列表"""
        try:
            request_data = {
                "jsonrpc": "2.0",
                "id": "tools_list",
                "method": "tools/list",
                "params": {}
            }
            
            response = requests.post(
                f"{self.base_url}/mcp",
                headers=self.headers,
                json=request_data
            )
            response.raise_for_status()
            return response.json()
        except requests.RequestException as e:
            return {"error": f"获取工具列表失败: {str(e)}"}
    
    def search_documents(self, query: str, limit: int = 10, doc_types: Optional[list] = None) -> Dict[str, Any]:
        """搜索文档"""
        try:
            arguments = {
                "query": query,
                "limit": limit
            }
            
            if doc_types:
                arguments["types"] = doc_types
            
            request_data = {
                "jsonrpc": "2.0",
                "id": "search_docs",
                "method": "tools/call",
                "params": {
                    "name": "search_documents",
                    "arguments": arguments
                }
            }
            
            response = requests.post(
                f"{self.base_url}/mcp",
                headers=self.headers,
                json=request_data
            )
            response.raise_for_status()
            return response.json()
        except requests.RequestException as e:
            return {"error": f"搜索文档失败: {str(e)}"}
    
    def get_document_content(self, document_id: str, version: Optional[str] = None) -> Dict[str, Any]:
        """获取文档内容"""
        try:
            arguments = {
                "document_id": document_id
            }
            
            if version:
                arguments["version"] = version
            
            request_data = {
                "jsonrpc": "2.0",
                "id": "get_doc",
                "method": "tools/call",
                "params": {
                    "name": "get_document_content",
                    "arguments": arguments
                }
            }
            
            response = requests.post(
                f"{self.base_url}/mcp",
                headers=self.headers,
                json=request_data
            )
            response.raise_for_status()
            return response.json()
        except requests.RequestException as e:
            return {"error": f"获取文档内容失败: {str(e)}"}

def main():
    """主函数"""
    parser = argparse.ArgumentParser(description="MCP客户端测试工具")
    parser.add_argument("--url", default="http://localhost:8080", help="服务器URL (默认: http://localhost:8080)")
    parser.add_argument("--key", required=True, help="API密钥")
    parser.add_argument("--test", action="store_true", help="测试连接")
    parser.add_argument("--list-tools", action="store_true", help="列出可用工具")
    parser.add_argument("--search", help="搜索文档的关键词")
    parser.add_argument("--limit", type=int, default=10, help="搜索结果限制 (默认: 10)")
    parser.add_argument("--types", help="文档类型过滤，用逗号分隔 (如: pdf,markdown)")
    parser.add_argument("--get-doc", help="获取指定ID的文档内容")
    parser.add_argument("--version", help="文档版本 (与--get-doc一起使用)")
    
    args = parser.parse_args()
    
    # 创建MCP客户端
    client = MCPClient(args.url, args.key)
    
    # 执行相应的操作
    if args.test:
        print("测试MCP连接...")
        result = client.test_connection()
        print(json.dumps(result, indent=2, ensure_ascii=False))
    
    elif args.list_tools:
        print("获取可用工具列表...")
        result = client.list_tools()
        print(json.dumps(result, indent=2, ensure_ascii=False))
    
    elif args.search:
        print(f"搜索文档: {args.search}")
        doc_types = None
        if args.types:
            doc_types = [t.strip() for t in args.types.split(',')]
        
        result = client.search_documents(args.search, args.limit, doc_types)
        print(json.dumps(result, indent=2, ensure_ascii=False))
    
    elif args.get_doc:
        print(f"获取文档内容: {args.get_doc}")
        result = client.get_document_content(args.get_doc, args.version)
        print(json.dumps(result, indent=2, ensure_ascii=False))
    
    else:
        # 如果没有指定具体操作，执行基本测试
        print("执行基本MCP功能测试...")
        
        # 1. 测试连接
        print("\n1. 测试连接:")
        result = client.test_connection()
        print(json.dumps(result, indent=2, ensure_ascii=False))
        
        if "error" in result:
            print("连接失败，请检查服务器和API密钥")
            sys.exit(1)
        
        # 2. 获取工具列表
        print("\n2. 获取工具列表:")
        result = client.list_tools()
        print(json.dumps(result, indent=2, ensure_ascii=False))
        
        # 3. 搜索文档
        print("\n3. 搜索文档示例:")
        result = client.search_documents("测试", limit=3)
        print(json.dumps(result, indent=2, ensure_ascii=False))

if __name__ == "__main__":
    main()