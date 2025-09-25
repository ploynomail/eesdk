#!/usr/bin/env python3
"""
EE SDK Python使用示例
"""
import sys
import os
import json

# 添加SDK路径
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '../../bindings/python'))

from ee_sdk import EESdkClient

def main():
    print("=== EE SDK Python示例 ===")
    
    try:
        # 创建客户端
        client = EESdkClient(
            base_url="https://httpbin.org",
            # base_url="http://127.0.0.1:8099",
            api_key="your-api-key-here",
            timeout=30,
            retry_count=3
        )
        
        # 示例1: 健康检查
        print("\n=== 健康检查示例 ===")
        if client.health_check():
            print("服务健康状态正常")
        else:
            print("健康检查失败")
        
        # 示例2: GET请求
        print("\n=== GET请求示例 ===")
        try:
            response = client.execute_request("GET", "/get")
            print(f"状态码: {response['status_code']}")
            print(f"响应体: {response['body']}")
        except RuntimeError as e:
            print(f"GET请求失败: {e}")
        
        # 示例3: POST请求
        print("\n=== POST请求示例 ===")
        user_data = {
            "name": "张三",
            "email": "zhangsan@example.com",
            "age": 25
        }
        
        try:
            response = client.execute_request("POST", "/post", user_data)
            print(f"状态码: {response['status_code']}")
            
            # 如果响应是JSON，打印用户ID
            if isinstance(response['body'], dict) and 'id' in response['body']:
                print(f"创建的用户ID: {response['body']['id']}")
            else:
                print(f"响应体: {response['body']}")
        except RuntimeError as e:
            print(f"POST请求失败: {e}")
        
        # 示例4: PUT请求
        print("\n=== PUT请求示例 ===")
        update_data = {
            "name": "李四",
            "age": 28
        }
        
        try:
            response = client.execute_request("PUT", "/put", update_data)
            print(f"状态码: {response['status_code']}")
            print(f"响应体: {response['body']}")
        except RuntimeError as e:
            print(f"PUT请求失败: {e}")
        
        # 示例5: 错误处理
        print("\n=== 错误处理示例 ===")
        try:
            # 故意发送到不存在的端点
            response = client.execute_request("GET", "/get")
            print(f"意外成功: {response}")
        except RuntimeError as e:
            print(f"预期的错误: {e}")
    
    except Exception as e:
        print(f"初始化客户端失败: {e}")
        return 1
    
    print("\n示例执行完成")
    return 0

if __name__ == "__main__":
    sys.exit(main())
