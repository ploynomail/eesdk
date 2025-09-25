package main

/*
#include <stdlib.h>
#include <string.h>

// 配置结构体
typedef struct {
    const char* base_url;
    const char* api_key;
    int timeout_seconds;
    int retry_count;
} ee_config_t;

// 响应结构体
typedef struct {
    int status_code;
    char* body;
    size_t body_length;
    char* error_message;
} ee_response_t;

// 客户端句柄（不透明指针）
typedef struct ee_client ee_client_t;
*/
import "C"

import (
	"context"
	"ee-sdk/src/client"
	"runtime"
	"sync"
	"time"
	"unsafe"
)

// 全局客户端映射，避免指针问题
var (
	clientMap   = make(map[uintptr]*client.Client)
	clientMutex sync.RWMutex
	nextID      uintptr = 1
)

//export ee_client_create
func ee_client_create(config *C.ee_config_t) *C.ee_client_t {
	if config == nil {
		return nil
	}

	goConfig := &client.Config{
		BaseURL:    C.GoString(config.base_url),
		APIKey:     C.GoString(config.api_key),
		Timeout:    time.Duration(config.timeout_seconds) * time.Second,
		RetryCount: int(config.retry_count),
	}

	goClient := client.NewClient(goConfig)

	// 使用映射管理客户端，避免直接指针转换
	clientMutex.Lock()
	id := nextID
	nextID++
	clientMap[id] = goClient
	clientMutex.Unlock()

	// 防止Go对象被垃圾回收
	runtime.SetFinalizer(goClient, nil)

	return (*C.ee_client_t)(unsafe.Pointer(id))
}

//export ee_client_destroy
func ee_client_destroy(clientPtr *C.ee_client_t) {
	if clientPtr == nil {
		return
	}

	id := uintptr(unsafe.Pointer(clientPtr))

	clientMutex.Lock()
	delete(clientMap, id)
	clientMutex.Unlock()
}

//export ee_execute_request
func ee_execute_request(clientPtr *C.ee_client_t, method *C.char, path *C.char, body *C.char) *C.ee_response_t {
	if clientPtr == nil {
		return createErrorResponse("客户端为空")
	}

	id := uintptr(unsafe.Pointer(clientPtr))

	clientMutex.RLock()
	goClient, exists := clientMap[id]
	clientMutex.RUnlock()

	if !exists {
		return createErrorResponse("无效的客户端句柄")
	}

	// 安全地转换C字符串
	goMethod := ""
	goPath := ""
	var goBody interface{}

	if method != nil {
		goMethod = C.GoString(method)
	}
	if path != nil {
		goPath = C.GoString(path)
	}
	if body != nil && *body != 0 {
		goBody = C.GoString(body)
	}

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := goClient.ExecuteRequest(ctx, goMethod, goPath, goBody)
	if err != nil {
		return createErrorResponse(err.Error())
	}

	// 创建C响应结构体
	cResp := (*C.ee_response_t)(C.malloc(C.size_t(unsafe.Sizeof(C.ee_response_t{}))))
	if cResp == nil {
		return createErrorResponse("内存分配失败")
	}

	cResp.status_code = C.int(resp.StatusCode)
	cResp.error_message = nil

	// 安全地复制响应体数据
	if len(resp.Body) > 0 {
		cResp.body_length = C.size_t(len(resp.Body))
		cResp.body = (*C.char)(C.malloc(C.size_t(len(resp.Body) + 1)))
		if cResp.body != nil {
			// 使用C的内存复制函数，更安全
			C.memcpy(unsafe.Pointer(cResp.body), unsafe.Pointer(&resp.Body[0]), C.size_t(len(resp.Body)))
			// 添加null终止符
			*(*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(cResp.body)) + uintptr(len(resp.Body)))) = 0
		}
	} else {
		cResp.body = nil
		cResp.body_length = 0
	}

	return cResp
}

//export ee_health_check
func ee_health_check(clientPtr *C.ee_client_t) C.int {
	if clientPtr == nil {
		return -1
	}

	id := uintptr(unsafe.Pointer(clientPtr))

	clientMutex.RLock()
	goClient, exists := clientMap[id]
	clientMutex.RUnlock()

	if !exists {
		return -1
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := goClient.Health(ctx)
	if err != nil {
		return -1
	}
	return 0
}

//export ee_response_free
func ee_response_free(response *C.ee_response_t) {
	if response == nil {
		return
	}

	if response.body != nil {
		C.free(unsafe.Pointer(response.body))
	}
	if response.error_message != nil {
		C.free(unsafe.Pointer(response.error_message))
	}
	C.free(unsafe.Pointer(response))
}

// 辅助函数：创建错误响应
func createErrorResponse(errorMsg string) *C.ee_response_t {
	cResp := (*C.ee_response_t)(C.malloc(C.size_t(unsafe.Sizeof(C.ee_response_t{}))))
	if cResp == nil {
		return nil
	}

	cResp.status_code = 0
	cResp.body = nil
	cResp.body_length = 0

	// 安全地复制错误信息
	if len(errorMsg) > 0 {
		errorBytes := []byte(errorMsg)
		cResp.error_message = (*C.char)(C.malloc(C.size_t(len(errorBytes) + 1)))
		if cResp.error_message != nil {
			C.memcpy(unsafe.Pointer(cResp.error_message), unsafe.Pointer(&errorBytes[0]), C.size_t(len(errorBytes)))
			*(*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(cResp.error_message)) + uintptr(len(errorBytes)))) = 0
		}
	} else {
		cResp.error_message = nil
	}

	return cResp
}

func main() {
	// CGO绑定需要main函数，但实际不会执行
}
