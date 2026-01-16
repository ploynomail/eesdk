package main

/*
#include <stdlib.h>
#include <string.h>

// 配置结构体
typedef struct {
    const char* server_addr;
    const char* access_id;
    const char* access_key;
} ee_config_t;

// 证书主题
typedef struct {
    const char* common_name;
    const char** country;
    int country_len;
    const char** organization;
    int organization_len;
    const char** organizational_unit;
    int organizational_unit_len;
    const char** locality;
    int locality_len;
    const char** province;
    int province_len;
} ee_subject_t;

// 证书请求
typedef struct {
    const unsigned char* public_key;
    size_t public_key_len;
    const unsigned char* kem_public_key;
    size_t kem_public_key_len;
    unsigned int ca_cert_id;
    unsigned int template_id;
    ee_subject_t subject;
} ee_cert_request_t;

// 证书响应
typedef struct {
    unsigned char* signer_cert;
    size_t signer_cert_len;
    unsigned char* enc_cert;
    size_t enc_cert_len;
	unsigned char* hyper_signer_cert;
	size_t hyper_signer_cert_len;
	unsigned char* hyper_enc_cert;
	size_t hyper_enc_cert_len;
    unsigned char* encrypted_private_key;
    size_t encrypted_private_key_len;
    char* error_message;
} ee_cert_response_t;

// 撤销请求
typedef struct {
    const char* serial_number;
    const char* reason;
} ee_revoke_request_t;

// 通用响应
typedef struct {
    int success;
    char* error_message;
} ee_response_t;

// 客户端句柄
typedef struct ee_client ee_client_t;
*/
import "C"

import (
	"runtime"
	"sync"
	"unsafe"

	eesdk "github.com/ploynomail/eesdk/client"
)

// 全局客户端映射
var (
	clientMap   = make(map[uintptr]*eesdk.Client)
	clientMutex sync.RWMutex
	nextID      uintptr = 1
)

//export ee_client_create
func ee_client_create(config *C.ee_config_t) *C.ee_client_t {
	if config == nil {
		return nil
	}

	goConfig := &eesdk.EESdkConfig{
		ServerAddr: C.GoString(config.server_addr),
		AccessId:   C.GoString(config.access_id),
		AccessKey:  C.GoString(config.access_key),
	}

	goClient, err := eesdk.NewClient(goConfig)
	if err != nil {
		return nil
	}

	clientMutex.Lock()
	id := nextID
	nextID++
	clientMap[id] = goClient
	clientMutex.Unlock()

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
	if client, exists := clientMap[id]; exists {
		client.Close()
		delete(clientMap, id)
	}
	clientMutex.Unlock()
}

//export ee_request_certificate
func ee_request_certificate(clientPtr *C.ee_client_t, req *C.ee_cert_request_t) *C.ee_cert_response_t {
	if clientPtr == nil || req == nil {
		return createCertErrorResponse("客户端或请求为空")
	}

	id := uintptr(unsafe.Pointer(clientPtr))
	clientMutex.RLock()
	goClient, exists := clientMap[id]
	clientMutex.RUnlock()

	if !exists {
		return createCertErrorResponse("无效的客户端句柄")
	}

	// 转换请求
	certReq := &eesdk.CertificateRequest{
		CACertID:   uint(req.ca_cert_id),
		TemplateID: uint(req.template_id),
		Subject: eesdk.Subject{
			CommonName: C.GoString(req.subject.common_name),
		},
	}

	if req.public_key != nil && req.public_key_len > 0 {
		certReq.PublicKey = C.GoBytes(unsafe.Pointer(req.public_key), C.int(req.public_key_len))
	}

	if req.kem_public_key != nil && req.kem_public_key_len > 0 {
		certReq.KemPublicKey = C.GoBytes(unsafe.Pointer(req.kem_public_key), C.int(req.kem_public_key_len))
	}

	// 转换字符串数组
	if req.subject.country != nil && req.subject.country_len > 0 {
		certReq.Subject.Country = cStringArrayToGo(req.subject.country, req.subject.country_len)
	}
	if req.subject.organization != nil && req.subject.organization_len > 0 {
		certReq.Subject.Organization = cStringArrayToGo(req.subject.organization, req.subject.organization_len)
	}

	// 调用客户端
	resp, err := goClient.RequestCertificate(certReq)
	if err != nil {
		return createCertErrorResponse(err.Error())
	}

	// 创建C响应
	cResp := (*C.ee_cert_response_t)(C.malloc(C.size_t(unsafe.Sizeof(C.ee_cert_response_t{}))))
	if cResp == nil {
		return createCertErrorResponse("内存分配失败")
	}

	cResp.error_message = nil

	// 复制证书数据
	if len(resp.SignerCert) > 0 {
		cResp.signer_cert_len = C.size_t(len(resp.SignerCert))
		cResp.signer_cert = (*C.uchar)(C.malloc(cResp.signer_cert_len))
		C.memcpy(unsafe.Pointer(cResp.signer_cert), unsafe.Pointer(&resp.SignerCert[0]), cResp.signer_cert_len)
	}

	if len(resp.EncCert) > 0 {
		cResp.enc_cert_len = C.size_t(len(resp.EncCert))
		cResp.enc_cert = (*C.uchar)(C.malloc(cResp.enc_cert_len))
		C.memcpy(unsafe.Pointer(cResp.enc_cert), unsafe.Pointer(&resp.EncCert[0]), cResp.enc_cert_len)
	}
	if len(resp.HyperSignerCert) > 0 {
		cResp.hyper_signer_cert_len = C.size_t(len(resp.HyperSignerCert))
		cResp.hyper_signer_cert = (*C.uchar)(C.malloc(cResp.hyper_signer_cert_len))
		C.memcpy(unsafe.Pointer(cResp.hyper_signer_cert), unsafe.Pointer(&resp.HyperSignerCert[0]), cResp.hyper_signer_cert_len)
	}
	if len(resp.HyperEncCert) > 0 {
		cResp.hyper_enc_cert_len = C.size_t(len(resp.HyperEncCert))
		cResp.hyper_enc_cert = (*C.uchar)(C.malloc(cResp.hyper_enc_cert_len))
		C.memcpy(unsafe.Pointer(cResp.hyper_enc_cert), unsafe.Pointer(&resp.HyperEncCert[0]), cResp.hyper_enc_cert_len)
	}
	if len(resp.EncryptedPrivateKey) > 0 {
		cResp.encrypted_private_key_len = C.size_t(len(resp.EncryptedPrivateKey))
		cResp.encrypted_private_key = (*C.uchar)(C.malloc(cResp.encrypted_private_key_len))
		C.memcpy(unsafe.Pointer(cResp.encrypted_private_key), unsafe.Pointer(&resp.EncryptedPrivateKey[0]), cResp.encrypted_private_key_len)
	}

	return cResp
}

//export ee_request_certificate_pkcs10
func ee_request_certificate_pkcs10(clientPtr *C.ee_client_t, csr *C.uchar, csrLen C.size_t) *C.ee_cert_response_t {
	if clientPtr == nil || csr == nil {
		return createCertErrorResponse("客户端或CSR为空")
	}

	id := uintptr(unsafe.Pointer(clientPtr))
	clientMutex.RLock()
	goClient, exists := clientMap[id]
	clientMutex.RUnlock()

	if !exists {
		return createCertErrorResponse("无效的客户端句柄")
	}

	csrBytes := C.GoBytes(unsafe.Pointer(csr), C.int(csrLen))
	resp, err := goClient.RequestCertificateWithPKCS10(csrBytes)
	if err != nil {
		return createCertErrorResponse(err.Error())
	}

	// 创建C响应
	cResp := (*C.ee_cert_response_t)(C.malloc(C.size_t(unsafe.Sizeof(C.ee_cert_response_t{}))))
	if cResp == nil {
		return createCertErrorResponse("内存分配失败")
	}

	// 初始化所有字段为 NULL/0
	cResp.signer_cert = nil
	cResp.signer_cert_len = 0
	cResp.enc_cert = nil
	cResp.enc_cert_len = 0
	cResp.hyper_signer_cert = nil
	cResp.hyper_signer_cert_len = 0
	cResp.hyper_enc_cert = nil
	cResp.hyper_enc_cert_len = 0
	cResp.encrypted_private_key = nil
	cResp.encrypted_private_key_len = 0
	cResp.error_message = nil

	if len(resp.SignerCert) > 0 {
		cResp.signer_cert_len = C.size_t(len(resp.SignerCert))
		cResp.signer_cert = (*C.uchar)(C.malloc(cResp.signer_cert_len))
		C.memcpy(unsafe.Pointer(cResp.signer_cert), unsafe.Pointer(&resp.SignerCert[0]), cResp.signer_cert_len)
	}

	if len(resp.EncCert) > 0 {
		cResp.enc_cert_len = C.size_t(len(resp.EncCert))
		cResp.enc_cert = (*C.uchar)(C.malloc(cResp.enc_cert_len))
		C.memcpy(unsafe.Pointer(cResp.enc_cert), unsafe.Pointer(&resp.EncCert[0]), cResp.enc_cert_len)
	}

	if len(resp.HyperSignerCert) > 0 {
		cResp.hyper_signer_cert_len = C.size_t(len(resp.HyperSignerCert))
		cResp.hyper_signer_cert = (*C.uchar)(C.malloc(cResp.hyper_signer_cert_len))
		C.memcpy(unsafe.Pointer(cResp.hyper_signer_cert), unsafe.Pointer(&resp.HyperSignerCert[0]), cResp.hyper_signer_cert_len)
	}

	if len(resp.HyperEncCert) > 0 {
		cResp.hyper_enc_cert_len = C.size_t(len(resp.HyperEncCert))
		cResp.hyper_enc_cert = (*C.uchar)(C.malloc(cResp.hyper_enc_cert_len))
		C.memcpy(unsafe.Pointer(cResp.hyper_enc_cert), unsafe.Pointer(&resp.HyperEncCert[0]), cResp.hyper_enc_cert_len)
	}

	if len(resp.EncryptedPrivateKey) > 0 {
		cResp.encrypted_private_key_len = C.size_t(len(resp.EncryptedPrivateKey))
		cResp.encrypted_private_key = (*C.uchar)(C.malloc(cResp.encrypted_private_key_len))
		C.memcpy(unsafe.Pointer(cResp.encrypted_private_key), unsafe.Pointer(&resp.EncryptedPrivateKey[0]), cResp.encrypted_private_key_len)
	}

	return cResp
}

//export ee_revoke_certificate
func ee_revoke_certificate(clientPtr *C.ee_client_t, req *C.ee_revoke_request_t) *C.ee_response_t {
	if clientPtr == nil || req == nil {
		return createErrorResponse("客户端或请求为空")
	}

	id := uintptr(unsafe.Pointer(clientPtr))
	clientMutex.RLock()
	goClient, exists := clientMap[id]
	clientMutex.RUnlock()

	if !exists {
		return createErrorResponse("无效的客户端句柄")
	}

	serialNumber := C.GoString(req.serial_number)
	reason := C.GoString(req.reason)

	err := goClient.RevokeCertificate(serialNumber, reason)
	if err != nil {
		return createErrorResponse(err.Error())
	}

	cResp := (*C.ee_response_t)(C.malloc(C.size_t(unsafe.Sizeof(C.ee_response_t{}))))
	cResp.success = 1
	cResp.error_message = nil
	return cResp
}

//export ee_confirm_certificate
func ee_confirm_certificate(clientPtr *C.ee_client_t, certHash *C.uchar, certHashLen C.size_t, certReqId C.int) *C.ee_response_t {
	if clientPtr == nil || certHash == nil {
		return createErrorResponse("客户端或证书哈希为空")
	}

	id := uintptr(unsafe.Pointer(clientPtr))
	clientMutex.RLock()
	goClient, exists := clientMap[id]
	clientMutex.RUnlock()

	if !exists {
		return createErrorResponse("无效的客户端句柄")
	}

	hash := C.GoBytes(unsafe.Pointer(certHash), C.int(certHashLen))
	err := goClient.ConfirmCertificate(hash, int32(certReqId))
	if err != nil {
		return createErrorResponse(err.Error())
	}

	cResp := (*C.ee_response_t)(C.malloc(C.size_t(unsafe.Sizeof(C.ee_response_t{}))))
	cResp.success = 1
	cResp.error_message = nil
	return cResp
}

//export ee_cert_response_free
func ee_cert_response_free(response *C.ee_cert_response_t) {
	if response == nil {
		return
	}

	if response.signer_cert != nil {
		C.free(unsafe.Pointer(response.signer_cert))
	}
	if response.enc_cert != nil {
		C.free(unsafe.Pointer(response.enc_cert))
	}
	if response.hyper_signer_cert != nil {
		C.free(unsafe.Pointer(response.hyper_signer_cert))
	}
	if response.hyper_enc_cert != nil {
		C.free(unsafe.Pointer(response.hyper_enc_cert))
	}
	if response.encrypted_private_key != nil {
		C.free(unsafe.Pointer(response.encrypted_private_key))
	}
	if response.error_message != nil {
		C.free(unsafe.Pointer(response.error_message))
	}
	C.free(unsafe.Pointer(response))
}

//export ee_response_free
func ee_response_free(response *C.ee_response_t) {
	if response == nil {
		return
	}

	if response.error_message != nil {
		C.free(unsafe.Pointer(response.error_message))
	}
	C.free(unsafe.Pointer(response))
}

// 辅助函数
func createCertErrorResponse(errorMsg string) *C.ee_cert_response_t {
	cResp := (*C.ee_cert_response_t)(C.malloc(C.size_t(unsafe.Sizeof(C.ee_cert_response_t{}))))
	if cResp == nil {
		return nil
	}

	// 初始化所有字段为 NULL/0
	cResp.signer_cert = nil
	cResp.signer_cert_len = 0
	cResp.enc_cert = nil
	cResp.enc_cert_len = 0
	cResp.hyper_signer_cert = nil
	cResp.hyper_signer_cert_len = 0
	cResp.hyper_enc_cert = nil
	cResp.hyper_enc_cert_len = 0
	cResp.encrypted_private_key = nil
	cResp.encrypted_private_key_len = 0

	if len(errorMsg) > 0 {
		errorBytes := []byte(errorMsg)
		cResp.error_message = (*C.char)(C.malloc(C.size_t(len(errorBytes) + 1)))
		C.memcpy(unsafe.Pointer(cResp.error_message), unsafe.Pointer(&errorBytes[0]), C.size_t(len(errorBytes)))
		*(*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(cResp.error_message)) + uintptr(len(errorBytes)))) = 0
	} else {
		cResp.error_message = nil
	}

	return cResp
}

func createErrorResponse(errorMsg string) *C.ee_response_t {
	cResp := (*C.ee_response_t)(C.malloc(C.size_t(unsafe.Sizeof(C.ee_response_t{}))))
	if cResp == nil {
		return nil
	}

	cResp.success = 0

	if len(errorMsg) > 0 {
		errorBytes := []byte(errorMsg)
		cResp.error_message = (*C.char)(C.malloc(C.size_t(len(errorBytes) + 1)))
		C.memcpy(unsafe.Pointer(cResp.error_message), unsafe.Pointer(&errorBytes[0]), C.size_t(len(errorBytes)))
		*(*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(cResp.error_message)) + uintptr(len(errorBytes)))) = 0
	}

	return cResp
}

func cStringArrayToGo(arr **C.char, length C.int) []string {
	if arr == nil || length == 0 {
		return nil
	}

	result := make([]string, int(length))
	slice := (*[1 << 28]*C.char)(unsafe.Pointer(arr))[:length:length]
	for i, str := range slice {
		result[i] = C.GoString(str)
	}
	return result
}

func main() {
	// CGO绑定需要main函数
}
