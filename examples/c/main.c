#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "libee_sdk.h"

void print_response(ee_response_t* response) {
    if (response->error_message) {
        printf("错误: %s\n", response->error_message);
    } else {
        printf("状态码: %d\n", response->status_code);
        printf("响应体长度: %zu\n", response->body_length);
        if (response->body) {
            printf("响应体: %s\n", response->body);
        }
    }
}

int main() {
    printf("=== EE SDK C语言示例 ===\n");
    
    // 创建配置 - 使用httpbin.org进行测试
    ee_config_t config = {
        .base_url = "https://httpbin.org",
        .api_key = "your-api-key-here",
        .timeout_seconds = 30,
        .retry_count = 3
    };
    
    // 创建客户端
    ee_client_t* client = ee_client_create(&config);
    if (!client) {
        printf("创建客户端失败\n");
        return 1;
    }
    
    // 示例1: 健康检查（使用httpbin的/get端点）
    printf("\n=== 健康检查示例 ===\n");
    int health_result = ee_health_check(client);
    if (health_result == 0) {
        printf("服务健康状态正常\n");
    } else {
        printf("健康检查失败\n");
    }
    
    // 示例2: GET请求
    printf("\n=== GET请求示例 ===\n");
    ee_response_t* response = ee_execute_request(client, "GET", "/get", NULL);
    if (response) {
        print_response(response);
        ee_response_free(response);
    }
    
    // 示例3: POST请求
    printf("\n=== POST请求示例 ===\n");
    // 使用非const字符串避免警告
    char json_body[] = "{\"name\":\"张三\",\"email\":\"zhangsan@example.com\",\"age\":25}";
    response = ee_execute_request(client, "POST", "/post", json_body);
    if (response) {
        print_response(response);
        ee_response_free(response);
    }
    
    // 示例4: PUT请求
    printf("\n=== PUT请求示例 ===\n");
    // 使用非const字符串避免警告
    char update_body[] = "{\"name\":\"李四\",\"age\":28}";
    response = ee_execute_request(client, "PUT", "/put", update_body);
    if (response) {
        print_response(response);
        ee_response_free(response);
    }
    
    // 清理资源
    ee_client_destroy(client);
    printf("\n客户端已销毁，示例执行完成\n");
    
    return 0;
}
