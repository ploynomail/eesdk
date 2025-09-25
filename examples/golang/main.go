package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"ee-sdk/src/client"
)

func main() {
	config := &client.Config{
		BaseURL:    "https://httpbin.org",
		APIKey:     "your-api-key-here",
		Timeout:    30 * time.Second,
		RetryCount: 3,
	}

	// 初始化客户端
	client := client.NewClient(config)

	ctx := context.Background()

	// 示例1: 健康检查
	fmt.Println("=== 健康检查示例 ===")
	if err := client.Health(ctx); err != nil {
		log.Printf("健康检查失败: %v", err)
	} else {
		fmt.Println("服务健康状态正常")
	}

	// 示例2: GET请求
	fmt.Println("\n=== GET请求示例 ===")
	resp, err := client.ExecuteRequest(ctx, "GET", "/get", nil)
	if err != nil {
		log.Printf("GET请求失败: %v", err)
	} else {
		fmt.Printf("状态码: %d\n", resp.StatusCode)
		fmt.Printf("响应体: %s\n", string(resp.Body))
	}

	// 示例3: POST请求
	fmt.Println("\n=== POST请求示例 ===")
	requestBody := map[string]interface{}{
		"name":  "张三",
		"email": "zhangsan@example.com",
		"age":   25,
	}

	resp, err = client.ExecuteRequest(ctx, "POST", "/post", requestBody)
	if err != nil {
		log.Printf("POST请求失败: %v", err)
	} else {
		fmt.Printf("状态码: %d\n", resp.StatusCode)

		// 解析JSON响应
		var result map[string]interface{}
		if err := json.Unmarshal(resp.Body, &result); err == nil {
			fmt.Printf("创建的用户ID: %v\n", result["id"])
		}
	}
}
