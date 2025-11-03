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
 * Signature: (Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)J
 */
JNIEXPORT jlong JNICALL Java_com_turingq_eesdk_EESdk_createClient
  (JNIEnv *env, jclass cls, jstring serverAddr, jstring accessId, jstring accessKey) {
    
    const char *c_server_addr = (*env)->GetStringUTFChars(env, serverAddr, NULL);
    const char *c_access_id = (*env)->GetStringUTFChars(env, accessId, NULL);
    const char *c_access_key = (*env)->GetStringUTFChars(env, accessKey, NULL);
    
    if (c_server_addr == NULL || c_access_id == NULL || c_access_key == NULL) {
        if (c_server_addr) (*env)->ReleaseStringUTFChars(env, serverAddr, c_server_addr);
        if (c_access_id) (*env)->ReleaseStringUTFChars(env, accessId, c_access_id);
        if (c_access_key) (*env)->ReleaseStringUTFChars(env, accessKey, c_access_key);
        return 0;
    }
    
    ee_config_t config = {
        .server_addr = c_server_addr,
        .access_id = c_access_id,
        .access_key = c_access_key
    };
    
    ee_client_t* client = ee_client_create(&config);
    
    (*env)->ReleaseStringUTFChars(env, serverAddr, c_server_addr);
    (*env)->ReleaseStringUTFChars(env, accessId, c_access_id);
    (*env)->ReleaseStringUTFChars(env, accessKey, c_access_key);
    
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
 * Method:    requestCertificate
 * Signature: (JLcom/turingq/eesdk/CertificateRequest;)Lcom/turingq/eesdk/CertificateResponse;
 */
JNIEXPORT jobject JNICALL Java_com_turingq_eesdk_EESdk_requestCertificate
  (JNIEnv *env, jclass cls, jlong clientPtr, jobject request) {
    
    if (clientPtr == 0 || request == NULL) {
        return NULL;
    }
    
    ee_client_t* client = (ee_client_t*)(intptr_t)clientPtr;
    
    // 获取请求对象的字段
    jclass reqClass = (*env)->GetObjectClass(env, request);
    jfieldID publicKeyField = (*env)->GetFieldID(env, reqClass, "publicKey", "[B");
    jfieldID caCertIdField = (*env)->GetFieldID(env, reqClass, "caCertId", "I");
    jfieldID templateIdField = (*env)->GetFieldID(env, reqClass, "templateId", "I");
    jfieldID subjectField = (*env)->GetFieldID(env, reqClass, "subject", "Lcom/turingq/eesdk/Subject;");
    
    jbyteArray publicKey = (jbyteArray)(*env)->GetObjectField(env, request, publicKeyField);
    jint caCertId = (*env)->GetIntField(env, request, caCertIdField);
    jint templateId = (*env)->GetIntField(env, request, templateIdField);
    jobject subject = (*env)->GetObjectField(env, request, subjectField);
    
    // 准备C请求结构
    ee_cert_request_t c_req = {0};
    c_req.ca_cert_id = (unsigned int)caCertId;
    c_req.template_id = (unsigned int)templateId;
    
    // 转换公钥
    if (publicKey != NULL) {
        jsize keyLen = (*env)->GetArrayLength(env, publicKey);
        c_req.public_key_len = keyLen;
        c_req.public_key = (unsigned char*)(*env)->GetByteArrayElements(env, publicKey, NULL);
    }
    
    // 转换主题
    if (subject != NULL) {
        jclass subjectClass = (*env)->GetObjectClass(env, subject);
        jfieldID cnField = (*env)->GetFieldID(env, subjectClass, "commonName", "Ljava/lang/String;");
        jstring cn = (jstring)(*env)->GetObjectField(env, subject, cnField);
        c_req.subject.common_name = (*env)->GetStringUTFChars(env, cn, NULL);
    }
    
    // 调用C函数
    ee_cert_response_t* response = ee_request_certificate(client, &c_req);
    
    // 清理
    if (publicKey != NULL) {
        (*env)->ReleaseByteArrayElements(env, publicKey, (jbyte*)c_req.public_key, JNI_ABORT);
    }
    
    if (response == NULL || response->error_message != NULL) {
        if (response) ee_cert_response_free(response);
        return NULL;
    }
    
    // 创建Java响应对象
    jclass respClass = (*env)->FindClass(env, "com/turingq/eesdk/CertificateResponse");
    jmethodID constructor = (*env)->GetMethodID(env, respClass, "<init>", "([B[B[B)V");
    
    jbyteArray signerCert = (*env)->NewByteArray(env, response->signer_cert_len);
    (*env)->SetByteArrayRegion(env, signerCert, 0, response->signer_cert_len, (jbyte*)response->signer_cert);
    
    jbyteArray encCert = (*env)->NewByteArray(env, response->enc_cert_len);
    (*env)->SetByteArrayRegion(env, encCert, 0, response->enc_cert_len, (jbyte*)response->enc_cert);
    
    jbyteArray hyperSignerCert = (*env)->NewByteArray(env, response->hyper_signer_cert_len);
    (*env)->SetByteArrayRegion(env, hyperSignerCert, 0, response->hyper_signer_cert_len, (jbyte*)response->hyper_signer_cert);

    jbyteArray hyperEncCert = (*env)->NewByteArray(env, response->hyper_enc_cert_len);
    (*env)->SetByteArrayRegion(env, hyperEncCert, 0, response->hyper_enc_cert_len, (jbyte*)response->hyper_enc_cert);

    jbyteArray encKey = (*env)->NewByteArray(env, response->encrypted_private_key_len);
    (*env)->SetByteArrayRegion(env, encKey, 0, response->encrypted_private_key_len, (jbyte*)response->encrypted_private_key);
    
    jobject responseObj = (*env)->NewObject(env, respClass, constructor, signerCert, encCert, encKey);
    
    ee_cert_response_free(response);
    
    return responseObj;
}

/*
 * Class:     com_turingq_eesdk_EESdk
 * Method:    requestCertificateWithPKCS10
 * Signature: (J[B)Lcom/turingq/eesdk/CertificateResponse;
 */
JNIEXPORT jobject JNICALL Java_com_turingq_eesdk_EESdk_requestCertificateWithPKCS10
  (JNIEnv *env, jclass cls, jlong clientPtr, jbyteArray csr) {
    
    if (clientPtr == 0 || csr == NULL) {
        return NULL;
    }
    
    ee_client_t* client = (ee_client_t*)(intptr_t)clientPtr;
    
    // 获取CSR数据
    jsize csrLen = (*env)->GetArrayLength(env, csr);
    jbyte* csrData = (*env)->GetByteArrayElements(env, csr, NULL);
    
    if (csrData == NULL) {
        return NULL;
    }
    
    // 调用C函数 - 修正函数名
    ee_cert_response_t* response = ee_request_certificate_pkcs10(
        client, 
        (unsigned char*)csrData, 
        (size_t)csrLen
    );
    
    // 释放CSR数据
    (*env)->ReleaseByteArrayElements(env, csr, csrData, JNI_ABORT);
    
    if (response == NULL || response->error_message != NULL) {
        if (response) ee_cert_response_free(response);
        return NULL;
    }
    
    // 创建Java响应对象
    jclass respClass = (*env)->FindClass(env, "com/turingq/eesdk/CertificateResponse");
    jmethodID constructor = (*env)->GetMethodID(env, respClass, "<init>", "([B[B[B[B[B)V");
    
    jbyteArray signerCert = (*env)->NewByteArray(env, response->signer_cert_len);
    (*env)->SetByteArrayRegion(env, signerCert, 0, response->signer_cert_len, (jbyte*)response->signer_cert);
    
    jbyteArray encCert = (*env)->NewByteArray(env, response->enc_cert_len);
    (*env)->SetByteArrayRegion(env, encCert, 0, response->enc_cert_len, (jbyte*)response->enc_cert);
    
    jbyteArray hyperSignCert = (*env)->NewByteArray(env, response->hyper_signer_cert_len);
    (*env)->SetByteArrayRegion(env, hyperSignCert, 0, response->hyper_signer_cert_len, (jbyte*)response->hyper_signer_cert);
    
    jbyteArray hyperEncCert = (*env)->NewByteArray(env, response->hyper_enc_cert_len);
    (*env)->SetByteArrayRegion(env, hyperEncCert, 0, response->hyper_enc_cert_len, (jbyte*)response->hyper_enc_cert);
    
    jbyteArray encKey = (*env)->NewByteArray(env, response->encrypted_private_key_len);
    (*env)->SetByteArrayRegion(env, encKey, 0, response->encrypted_private_key_len, (jbyte*)response->encrypted_private_key);
    
    jobject responseObj = (*env)->NewObject(env, respClass, constructor, 
        signerCert, encCert, hyperSignCert, hyperEncCert, encKey);
    
    ee_cert_response_free(response);
    
    return responseObj;
}

/*
 * Class:     com_turingq_eesdk_EESdk
 * Method:    revokeCertificate
 * Signature: (JLjava/lang/String;Ljava/lang/String;)Z
 */
JNIEXPORT jboolean JNICALL Java_com_turingq_eesdk_EESdk_revokeCertificate
  (JNIEnv *env, jclass cls, jlong clientPtr, jstring serialNumber, jstring reason) {
    
    if (clientPtr == 0 || serialNumber == NULL || reason == NULL) {
        return JNI_FALSE;
    }
    
    ee_client_t* client = (ee_client_t*)(intptr_t)clientPtr;
    
    const char* c_serial = (*env)->GetStringUTFChars(env, serialNumber, NULL);
    const char* c_reason = (*env)->GetStringUTFChars(env, reason, NULL);
    
    if (c_serial == NULL || c_reason == NULL) {
        if (c_serial) (*env)->ReleaseStringUTFChars(env, serialNumber, c_serial);
        if (c_reason) (*env)->ReleaseStringUTFChars(env, reason, c_reason);
        return JNI_FALSE;
    }
    
    // 准备撤销请求结构体
    ee_revoke_request_t req = {
        .serial_number = c_serial,
        .reason = c_reason
    };
    
    // 调用C函数 - 返回类型是 ee_response_t*
    ee_response_t* result = ee_revoke_certificate(client, &req);
    
    (*env)->ReleaseStringUTFChars(env, serialNumber, c_serial);
    (*env)->ReleaseStringUTFChars(env, reason, c_reason);
    
    jboolean success = (result != NULL && result->error_message == NULL) ? JNI_TRUE : JNI_FALSE;
    if (result != NULL) {
        ee_response_free(result);
    }
    
    return success;
}

/*
 * Class:     com_turingq_eesdk_EESdk
 * Method:    confirmCertificate
 * Signature: (J[BI)Z
 */
JNIEXPORT jboolean JNICALL Java_com_turingq_eesdk_EESdk_confirmCertificate
  (JNIEnv *env, jclass cls, jlong clientPtr, jbyteArray certHash, jint certReqId) {
    
    if (clientPtr == 0 || certHash == NULL) {
        return JNI_FALSE;
    }
    
    ee_client_t* client = (ee_client_t*)(intptr_t)clientPtr;
    
    jsize hashLen = (*env)->GetArrayLength(env, certHash);
    jbyte* hashData = (*env)->GetByteArrayElements(env, certHash, NULL);
    
    if (hashData == NULL) {
        return JNI_FALSE;
    }
    
    // 调用C函数 - 传递哈希数据、长度和证书请求ID
    ee_response_t* result = ee_confirm_certificate(
        client, 
        (unsigned char*)hashData, 
        (size_t)hashLen,
        (unsigned int)certReqId
    );
    
    (*env)->ReleaseByteArrayElements(env, certHash, hashData, JNI_ABORT);
    
    jboolean success = (result != NULL && result->error_message == NULL) ? JNI_TRUE : JNI_FALSE;
    if (result != NULL) {
        ee_response_free(result);
    }
    
    return success;
}

