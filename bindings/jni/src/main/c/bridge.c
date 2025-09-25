#include <jni.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include "com_turingq_eesdk_EESdk.h"

// 包含C绑定头文件
#include "libee_sdk.h"

/*
 * Class:     com_turingq_eesdk_EESdk
 * Method:    createClient
 * Signature: (Ljava/lang/String;Ljava/lang/String;II)J
 */
JNIEXPORT jlong JNICALL Java_com_turingq_eesdk_EESdk_createClient
  (JNIEnv *env, jclass cls, jstring baseURL, jstring apiKey, jint timeout, jint retryCount) {
    
    // 转换Java字符串到C字符串
    const char *c_base_url = (*env)->GetStringUTFChars(env, baseURL, NULL);
    const char *c_api_key = (*env)->GetStringUTFChars(env, apiKey, NULL);
    
    if (c_base_url == NULL || c_api_key == NULL) {
        // 释放已分配的字符串
        if (c_base_url) (*env)->ReleaseStringUTFChars(env, baseURL, c_base_url);
        if (c_api_key) (*env)->ReleaseStringUTFChars(env, apiKey, c_api_key);
        return 0; // 返回NULL指针
    }
    
    // 创建配置结构体
    ee_config_t config = {
        .base_url = c_base_url,
        .api_key = c_api_key,
        .timeout_seconds = (int)timeout,
        .retry_count = (int)retryCount
    };
    
    // 调用C函数创建客户端
    ee_client_t* client = ee_client_create(&config);
    
    // 释放字符串资源
    (*env)->ReleaseStringUTFChars(env, baseURL, c_base_url);
    (*env)->ReleaseStringUTFChars(env, apiKey, c_api_key);
    
    // 返回客户端指针（转换为long）
    return (jlong)(intptr_t)client;
}

/*
 * Class:     com_turingq_eesdk_EESdk
 * Method:    destroyClient
 * Signature: (J)V
 */
JNIEXPORT void JNICALL Java_com_turingq_eesdk_EESdk_destroyClient
  (JNIEnv *env, jclass cls, jlong clientPtr) {
    
    if (clientPtr != 0) {
        ee_client_t* client = (ee_client_t*)(intptr_t)clientPtr;
        ee_client_destroy(client);
    }
}

/*
 * Class:     com_turingq_eesdk_EESdk
 * Method:    executeRequest
 * Signature: (JLjava/lang/String;Ljava/lang/String;Ljava/lang/String;)Lcom/turingq/eesdk/Response;
 */
JNIEXPORT jobject JNICALL Java_com_turingq_eesdk_EESdk_executeRequest
  (JNIEnv *env, jclass cls, jlong clientPtr, jstring method, jstring path, jstring body) {
    
    if (clientPtr == 0) {
        // 创建错误响应
        jclass responseClass = (*env)->FindClass(env, "com/turingq/eesdk/Response");
        jmethodID constructor = (*env)->GetMethodID(env, responseClass, "<init>", "(ILjava/lang/String;Ljava/lang/String;)V");
        jstring errorMsg = (*env)->NewStringUTF(env, "客户端指针为空");
        return (*env)->NewObject(env, responseClass, constructor, (jint)0, NULL, errorMsg);
    }
    
    ee_client_t* client = (ee_client_t*)(intptr_t)clientPtr;
    
    // 转换Java字符串到C字符串
    const char *c_method = (*env)->GetStringUTFChars(env, method, NULL);
    const char *c_path = (*env)->GetStringUTFChars(env, path, NULL);
    const char *c_body = NULL;
    
    if (body != NULL) {
        c_body = (*env)->GetStringUTFChars(env, body, NULL);
    }
    
    if (c_method == NULL || c_path == NULL) {
        // 清理资源并返回错误
        if (c_method) (*env)->ReleaseStringUTFChars(env, method, c_method);
        if (c_path) (*env)->ReleaseStringUTFChars(env, path, c_path);
        if (c_body) (*env)->ReleaseStringUTFChars(env, body, c_body);
        
        jclass responseClass = (*env)->FindClass(env, "com/turingq/eesdk/Response");
        jmethodID constructor = (*env)->GetMethodID(env, responseClass, "<init>", "(ILjava/lang/String;Ljava/lang/String;)V");
        jstring errorMsg = (*env)->NewStringUTF(env, "字符串转换失败");
        return (*env)->NewObject(env, responseClass, constructor, (jint)0, NULL, errorMsg);
    }
    
    // 调用C函数执行请求
    ee_response_t* response = ee_execute_request(client, (char*)c_method, (char*)c_path, (char*)c_body);
    
    // 释放字符串资源
    (*env)->ReleaseStringUTFChars(env, method, c_method);
    (*env)->ReleaseStringUTFChars(env, path, c_path);
    if (c_body) {
        (*env)->ReleaseStringUTFChars(env, body, c_body);
    }
    
    if (response == NULL) {
        jclass responseClass = (*env)->FindClass(env, "com/turingq/eesdk/Response");
        jmethodID constructor = (*env)->GetMethodID(env, responseClass, "<init>", "(ILjava/lang/String;Ljava/lang/String;)V");
        jstring errorMsg = (*env)->NewStringUTF(env, "请求执行失败");
        return (*env)->NewObject(env, responseClass, constructor, (jint)0, NULL, errorMsg);
    }
    
    // 创建Java Response对象
    jclass responseClass = (*env)->FindClass(env, "com/turingq/eesdk/Response");
    jmethodID constructor = (*env)->GetMethodID(env, responseClass, "<init>", "(ILjava/lang/String;Ljava/lang/String;)V");
    
    // 转换C字符串到Java字符串
    jstring javaBody = NULL;
    jstring javaErrorMsg = NULL;
    
    if (response->body && response->body_length > 0) {
        javaBody = (*env)->NewStringUTF(env, response->body);
    }
    
    if (response->error_message) {
        javaErrorMsg = (*env)->NewStringUTF(env, response->error_message);
    }
    
    // 创建Response对象
    jobject responseObj = (*env)->NewObject(env, responseClass, constructor, 
                                           (jint)response->status_code, javaBody, javaErrorMsg);
    
    // 释放C响应资源
    ee_response_free(response);
    
    return responseObj;
}

/*
 * Class:     com_turingq_eesdk_EESdk
 * Method:    healthCheck
 * Signature: (J)Z
 */
JNIEXPORT jboolean JNICALL Java_com_turingq_eesdk_EESdk_healthCheck
  (JNIEnv *env, jclass cls, jlong clientPtr) {
    
    if (clientPtr == 0) {
        return JNI_FALSE;
    }
    
    ee_client_t* client = (ee_client_t*)(intptr_t)clientPtr;
    int result = ee_health_check(client);
    
    return (result == 0) ? JNI_TRUE : JNI_FALSE;
}

