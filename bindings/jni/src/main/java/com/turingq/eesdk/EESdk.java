package com.turingq.eesdk;

/**
 * EE SDK Java绑定
 */
public class EESdk {
    
    static {
        try {
            System.loadLibrary("ee_sdk_jni");
        } catch (UnsatisfiedLinkError e) {
            System.err.println("无法加载native库: " + e.getMessage());
            throw e;
        }
    }
    
    // Native方法声明
    private static native long createClient(String baseURL, String apiKey, int timeout, int retryCount);
    private static native void destroyClient(long clientPtr);
    private static native Response executeRequest(long clientPtr, String method, String path, String body);
    private static native boolean healthCheck(long clientPtr);
    
    private long clientPtr;
    
    /**
     * 构造函数
     * 
     * @param baseURL API基础URL
     * @param apiKey API密钥
     * @param timeout 超时时间（秒）
     * @param retryCount 重试次数
     */
    public EESdk(String baseURL, String apiKey, int timeout, int retryCount) {
        this.clientPtr = createClient(baseURL, apiKey, timeout, retryCount);
        if (this.clientPtr == 0) {
            throw new RuntimeException("无法创建SDK客户端");
        }
    }
    
    /**
     * 执行HTTP请求
     * 
     * @param method HTTP方法
     * @param path 请求路径
     * @param body 请求体（可选）
     * @return Response对象
     */
    public Response executeRequest(String method, String path, String body) {
        if (clientPtr == 0) {
            throw new IllegalStateException("客户端已被销毁");
        }
        return executeRequest(clientPtr, method, path, body);
    }
    
    /**
     * 执行HTTP请求（无请求体）
     * 
     * @param method HTTP方法
     * @param path 请求路径
     * @return Response对象
     */
    public Response executeRequest(String method, String path) {
        return executeRequest(method, path, null);
    }
    
    /**
     * 健康检查
     * 
     * @return 是否健康
     */
    public boolean healthCheck() {
        if (clientPtr == 0) {
            throw new IllegalStateException("客户端已被销毁");
        }
        return healthCheck(clientPtr);
    }
    
    /**
     * 销毁客户端
     */
    public void destroy() {
        if (clientPtr != 0) {
            destroyClient(clientPtr);
            clientPtr = 0;
        }
    }
    
    @Override
    protected void finalize() throws Throwable {
        destroy();
        super.finalize();
    }
}
