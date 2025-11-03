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
    private static native long createClient(String serverAddr, String accessId, String accessKey);
    private static native void destroyClient(long clientPtr);
    private static native CertificateResponse requestCertificate(long clientPtr, CertificateRequest request);
    private static native CertificateResponse requestCertificateWithPKCS10(long clientPtr, byte[] csr);
    private static native boolean revokeCertificate(long clientPtr, String serialNumber, String reason);
    private static native boolean confirmCertificate(long clientPtr, byte[] certHash, int certReqId);
    
    private long clientPtr;
    
    /**
     * 构造函数
     * 
     * @param serverAddr 服务器地址
     * @param accessId 访问ID
     * @param accessKey 访问密钥
     */
    public EESdk(String serverAddr, String accessId, String accessKey) {
        this.clientPtr = createClient(serverAddr, accessId, accessKey);
        if (this.clientPtr == 0) {
            throw new RuntimeException("无法创建SDK客户端");
        }
    }
    
    /**
     * 请求证书
     * 
     * @param request 证书请求
     * @return CertificateResponse对象
     */
    public CertificateResponse requestCertificate(CertificateRequest request) {
        if (clientPtr == 0) {
            throw new IllegalStateException("客户端已被销毁");
        }
        return requestCertificate(clientPtr, request);
    }
    
    /**
     * 使用PKCS#10请求证书
     * 
     * @param csr PKCS#10 CSR
     * @return CertificateResponse对象
     */
    public CertificateResponse requestCertificateWithPKCS10(byte[] csr) {
        if (clientPtr == 0) {
            throw new IllegalStateException("客户端已被销毁");
        }
        return requestCertificateWithPKCS10(clientPtr, csr);
    }
    
    /**
     * 撤销证书
     * 
     * @param serialNumber 证书序列号
     * @param reason 撤销原因
     * @return 是否成功
     */
    public boolean revokeCertificate(String serialNumber, String reason) {
        if (clientPtr == 0) {
            throw new IllegalStateException("客户端已被销毁");
        }
        return revokeCertificate(clientPtr, serialNumber, reason);
    }
    
    /**
     * 确认证书
     * 
     * @param certHash 证书哈希
     * @param certReqId 证书请求ID
     * @return 是否成功
     */
    public boolean confirmCertificate(byte[] certHash, int certReqId) {
        if (clientPtr == 0) {
            throw new IllegalStateException("客户端已被销毁");
        }
        return confirmCertificate(clientPtr, certHash, certReqId);
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
