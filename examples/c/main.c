#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "libee_sdk.h"

/**
 * 示例1: 请求证书
 */
void example_request_certificate(ee_client_t* client) {
    printf("\n=== 请求证书示例 ===\n");
    
    // 模拟公钥（实际应用中应该从密钥生成获取）
    unsigned char public_key[65] = {0x04}; // 未压缩点标识
    
    // KEM公钥（Base64解码后的数据，这里简化处理）
    unsigned char kem_public_key[1184] = {0}; // 实际应该从Base64解码
    
    // 创建证书主题
    ee_subject_t subject = {
        .common_name = "test-device",
        .country = (const char*[]){"CN"},
        .country_len = 1,
        .organization = (const char*[]){"TuringQ"},
        .organization_len = 1,
        .organizational_unit = (const char*[]){"Engineering"},
        .organizational_unit_len = 1,
        .locality = (const char*[]){"Beijing"},
        .locality_len = 1,
        .province = (const char*[]){"Beijing"},
        .province_len = 1
    };
    
    // 创建证书请求
    ee_cert_request_t request = {
        .public_key = public_key,
        .public_key_len = sizeof(public_key),
        .kem_public_key = kem_public_key,
        .kem_public_key_len = sizeof(kem_public_key),
        .ca_cert_id = 1,
        .template_id = 1,
        .subject = subject
    };
    
    // 请求证书
    ee_cert_response_t* response = ee_request_certificate(client, &request);
    
    if (response) {
        if (response->error_message) {
            printf("证书请求失败: %s\n", response->error_message);
        } else {
            printf("证书请求成功!\n");
            printf("签名证书长度: %zu bytes\n", response->signer_cert_len);
            printf("加密证书长度: %zu bytes\n", response->enc_cert_len);
            printf("加密私钥长度: %zu bytes\n", response->encrypted_private_key_len);
            
            // 保存证书到文件
            FILE* f = fopen("/tmp/signer_cert.der", "wb");
            if (f) {
                fwrite(response->signer_cert, 1, response->signer_cert_len, f);
                fclose(f);
            }
            
            f = fopen("/tmp/enc_cert.der", "wb");
            if (f) {
                fwrite(response->enc_cert, 1, response->enc_cert_len, f);
                fclose(f);
            }
            printf("证书已保存到 /tmp/signer_cert.der 和 /tmp/enc_cert.der\n");
        }
        ee_cert_response_free(response);
    } else {
        printf("证书请求失败: 返回NULL\n");
    }
}

/**
 * 示例2: 使用PKCS#10 CSR请求证书
 */
void example_request_certificate_with_pkcs10(ee_client_t* client) {
    printf("\n=== 使用PKCS#10请求证书示例 ===\n");
    
    // 读取CSR文件
    FILE* f = fopen("../testdata/example.csr", "rb");
    if (!f) {
        printf("CSR文件不存在: ../testdata/example.csr\n");
        return;
    }
    
    fseek(f, 0, SEEK_END);
    long file_size = ftell(f);
    fseek(f, 0, SEEK_SET);
    
    unsigned char* csr_data = (unsigned char*)malloc(file_size);
    fread(csr_data, 1, file_size, f);
    fclose(f);
    
    // 使用CSR请求证书
    ee_cert_response_t* response = ee_request_certificate_pkcs10(client, csr_data, file_size);
    
    if (response) {
        if (response->error_message) {
            printf("CSR证书请求失败: %s\n", response->error_message);
        } else {
            printf("CSR证书请求成功!\n");
            printf("签名证书长度: %zu bytes\n", response->hyper_signer_cert_len);
            printf("加密证书长度: %zu bytes\n", response->hyper_enc_cert_len);
            printf("加密私钥长度: %zu bytes\n", response->encrypted_private_key_len);
        }
        ee_cert_response_free(response);
    } else {
        printf("CSR证书请求失败: 返回NULL\n");
    }
    
    free(csr_data);
}

/**
 * 示例3: 撤销证书
 */
void example_revoke_certificate(ee_client_t* client) {
    printf("\n=== 撤销证书示例 ===\n");
    
    // 创建撤销请求
    ee_revoke_request_t request = {
        .serial_number = "74",
        .reason = "keyCompromise"
    };
    
    // 撤销证书
    ee_response_t* response = ee_revoke_certificate(client, &request);
    
    if (response) {
        if (response->error_message) {
            printf("证书撤销失败: %s\n", response->error_message);
        } else if (response->success) {
            printf("证书撤销成功! 序列号: %s, 原因: %s\n", 
                   request.serial_number, request.reason);
        } else {
            printf("证书撤销失败\n");
        }
        ee_response_free(response);
    } else {
        printf("证书撤销失败: 返回NULL\n");
    }
}

/**
 * 示例4: 确认证书
 */
void example_confirm_certificate(ee_client_t* client) {
    printf("\n=== 确认证书示例 ===\n");
    
    // 模拟证书哈希（实际应用中应该从证书计算）
    unsigned char cert_hash[32] = {0};
    int cert_req_id = 123;
    
    // 确认证书
    ee_response_t* response = ee_confirm_certificate(client, cert_hash, sizeof(cert_hash), cert_req_id);
    
    if (response) {
        if (response->error_message) {
            printf("证书确认失败: %s\n", response->error_message);
        } else if (response->success) {
            printf("证书确认成功! 请求ID: %d\n", cert_req_id);
        } else {
            printf("证书确认失败\n");
        }
        ee_response_free(response);
    } else {
        printf("证书确认失败: 返回NULL\n");
    }
}

int main() {
    printf("=== EE SDK C语言 PKI示例 ===\n");
    
    // 创建配置
    ee_config_t config = {
        .server_addr = "localhost:9001",
        .access_id = "807e8150-da11-48a8-886e-54e9894e8003",
        .access_key = "HanBpT2gT31Gfn83R0dTjKu0O9mv55bGDyVtP3IWOW1="
    };
    
    // 创建客户端
    ee_client_t* client = ee_client_create(&config);
    if (!client) {
        printf("创建客户端失败\n");
        return 1;
    }
    
    printf("SDK客户端创建成功\n");
    
    // 运行示例
    // example_request_certificate(client);
    example_request_certificate_with_pkcs10(client);
    // example_revoke_certificate(client);
    // example_confirm_certificate(client);
    
    // 清理资源
    ee_client_destroy(client);
    printf("\n所有示例执行完成\n");
    
    return 0;
}
