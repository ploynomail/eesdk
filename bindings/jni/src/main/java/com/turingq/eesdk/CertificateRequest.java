package com.turingq.eesdk;

/**
 * 证书请求
 */
public class CertificateRequest {
    private byte[] publicKey;
    private byte[] kemPublicKey;
    private int caCertId;
    private int templateId;
    private Subject subject;
    
    public CertificateRequest(byte[] publicKey, int caCertId, int templateId, Subject subject) {
        this.publicKey = publicKey;
        this.caCertId = caCertId;
        this.templateId = templateId;
        this.subject = subject;
    }
    
    public byte[] getPublicKey() { return publicKey; }
    public void setPublicKey(byte[] publicKey) { this.publicKey = publicKey; }
    
    public byte[] getKemPublicKey() { return kemPublicKey; }
    public void setKemPublicKey(byte[] kemPublicKey) { this.kemPublicKey = kemPublicKey; }
    
    public int getCaCertId() { return caCertId; }
    public void setCaCertId(int caCertId) { this.caCertId = caCertId; }
    
    public int getTemplateId() { return templateId; }
    public void setTemplateId(int templateId) { this.templateId = templateId; }
    
    public Subject getSubject() { return subject; }
    public void setSubject(Subject subject) { this.subject = subject; }
}
