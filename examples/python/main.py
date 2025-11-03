#!/usr/bin/env python3
"""
EE SDK Python使用示例 - PKI证书管理
"""
import sys
import os
import base64

# 添加SDK路径
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '../../bindings/python'))

from ee_sdk import EESdkClient, Subject

def example_request_certificate(client):
    """示例1: 请求证书"""
    print("\n=== 请求证书示例 ===")
    
    # 模拟公钥和KEM公钥（实际应用中应该从密钥生成获取）
    public_key = b'\x04' + b'\x00' * 64  # 示例公钥
    kem_public_key_base64 = "HkSLs9sGqvNIUio0yLjE90wa7nVAkBKtyTjPd8ygSXURreDNVeRylEgQDydrqzgE9lZtLufJ/Wt5QuwWndaA/Et31WqgkWSQlnd1JHaQCXaFAZbIgxbKT3l90EYSY8eoYqFWTYZ7BLVy6jB3geU+hcYmr0IBiUQ2BFYFMYN79/Y8Z8IIB7E9UEgzrkmRRtYaLUUy1dpWo8JF+qiToqIRErMq/ulrazyog1mZeMGq4cRZFiNsOqRRfullrmAvpXM98quGlyAawZihiZWG+cJj1+W6iCUMB9qkkZWgoOZxuVinqAAnzceWyTYTE8tLUSuOQ4nLsyuCN5MY24CUdDluIRYmLtBwOZRT9OpfuJoNQ3dg9magj6IF0htNRms6FHluukQJcIw3eRURZVqLDvGG3TYxkJUcXusDcokvnqjPcvNY8hMxpNOTF/mVtdsL7exuG8oOi+UxzyBn2TkSROpjpwKgQRghp5WDMFtRP0QxkZc2iESBsDcXDihnPFSLgWwNxwPOxchKJ9kQx1QXrQoqMuXJXzM3Dsh3ZntWt2GNENMYVkcKAFQ/ZARs41glvMAWd4UCN2jIWuOGXnEzJQs0D0Y++7efMkEt1+k175p0wQpiE+iUCbSBEJILq5iM9RB2dTwQ8Zqdf+lYE0pTYqdURXdlCtIp+pdIyXVnAKNPbrHH2AbK8KyALmdih7Q+qHc4+hM3FMm6YLKS4UOit5McpKuvd+pwT5EvufB1DaG7+wpP/PYtg0CBXFtqnakSNeREHOl0NPijQzgt59VV1UN4yFwR82cbgQypfwd6q4V6+bMlNbqa4MtJAkSgbKRPRyAS4mh9oqUICQO/5ZPFwmGKZdsuipaZhRWkGSJaZxOhyAOyXLoOeTcXmbu2iFt8t9nIvzyDixO9TsQHI1M3hdXBrfpULWNEbzc1ugJB2ccQ0pm6dFkCl0oahCxXxUxFpHcUz8colDgaSXAgCzA9zAtbA9JmVnGXuWfPUZu0pQBZyYUr0ONKYnYOCvUOIOdJPLeqAJm96ABZCKZ/16sS8AGTHP2mmi2xxQ9IL6biZGiS0TI="
    kem_public_key = base64.b64decode(kem_public_key_base64)
    
    try:
        # 创建证书主题
        subject = Subject(
            common_name="test-device",
            country=["CN"],
            organization=["TuringQ"],
            organizational_unit=["Engineering"],
            locality=["Beijing"],
            province=["Beijing"]
        )
        
        # 请求证书
        response = client.request_certificate(
            public_key=public_key,
            kem_public_key=kem_public_key,
            ca_cert_id=1,
            template_id=1,
            subject=subject
        )
        
        print("证书请求成功!")
        print(f"签名证书长度: {len(response.signer_cert)} bytes")
        print(f"加密证书长度: {len(response.enc_cert)} bytes")
        print(f"加密私钥长度: {len(response.encrypted_private_key)} bytes")
        
        # 保存证书到文件
        with open('/tmp/signer_cert.der', 'wb') as f:
            f.write(response.signer_cert)
        with open('/tmp/enc_cert.der', 'wb') as f:
            f.write(response.enc_cert)
        print("证书已保存到 /tmp/signer_cert.der 和 /tmp/enc_cert.der")
        
    except RuntimeError as e:
        print(f"证书请求失败: {e}")

def example_request_certificate_with_pkcs10(client):
    """示例2: 使用PKCS#10 CSR请求证书"""
    print("\n=== 使用PKCS#10请求证书示例 ===")
    
    try:
        # 读取CSR文件
        csr_file = os.path.join(os.path.dirname(__file__), '../testdata/example.csr')
        with open(csr_file, 'rb') as f:
            csr_pem = f.read()
        
        # 使用CSR请求证书
        response = client.request_certificate_with_pkcs10(csr_pem)
        
        print("CSR证书请求成功!")
        print(f"签名证书长度: {len(response.signer_cert)} bytes")
        print(f"加密证书长度: {len(response.enc_cert)} bytes")
        
    except FileNotFoundError:
        print(f"CSR文件不存在: {csr_file}")
    except RuntimeError as e:
        print(f"CSR证书请求失败: {e}")

def example_revoke_certificate(client):
    """示例3: 撤销证书"""
    print("\n=== 撤销证书示例 ===")
    
    try:
        # 撤销证书（使用示例序列号）
        serial_number = "74"
        reason = "keyCompromise"
        
        success = client.revoke_certificate(serial_number, reason)
        
        if success:
            print(f"证书撤销成功! 序列号: {serial_number}, 原因: {reason}")
        else:
            print("证书撤销失败")
            
    except RuntimeError as e:
        print(f"证书撤销失败: {e}")

def example_confirm_certificate(client):
    """示例4: 确认证书"""
    print("\n=== 确认证书示例 ===")
    
    try:
        # 模拟证书哈希（实际应用中应该从证书计算）
        cert_hash = b'\x00' * 32  # 示例哈希
        cert_req_id = 123
        
        success = client.confirm_certificate(cert_hash, cert_req_id)
        
        if success:
            print(f"证书确认成功! 请求ID: {cert_req_id}")
        else:
            print("证书确认失败")
            
    except RuntimeError as e:
        print(f"证书确认失败: {e}")

def main():
    print("=== EE SDK Python PKI示例 ===")
    
    try:
        # 创建客户端
        client = EESdkClient(
            server_addr="localhost:9001",
            access_id="807e8150-da11-48a8-886e-54e9894e8003",
            access_key="HanBpT2gT31Gfn83R0dTjKu0O9mv55bGDyVtP3IWOW1="
        )
        
        print("SDK客户端创建成功")
        
        # 运行示例
        example_request_certificate(client)
        example_request_certificate_with_pkcs10(client)
        example_revoke_certificate(client)
        example_confirm_certificate(client)
        
    except Exception as e:
        print(f"初始化客户端失败: {e}")
        return 1
    
    print("\n所有示例执行完成")
    return 0

if __name__ == "__main__":
    sys.exit(main())
