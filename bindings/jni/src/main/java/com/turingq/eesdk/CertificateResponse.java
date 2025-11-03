package com.turingq.eesdk;

/**
 * 证书响应
 */
public class CertificateResponse {
    private final byte[] signerCert;
    private final byte[] encCert;
    private final byte[] HyperSignCert;
    private final byte[] HyperEncCert;
    private final byte[] encryptedPrivateKey;
    
    public CertificateResponse(byte[] signerCert, byte[] encCert, byte[] HyperSignCert, byte[] HyperEncCert, byte[] encryptedPrivateKey) {
        this.signerCert = signerCert;
        this.encCert = encCert;
        this.HyperSignCert = HyperSignCert;
        this.HyperEncCert = HyperEncCert;
        this.encryptedPrivateKey = encryptedPrivateKey;
    }
    
    public byte[] getSignerCert() { return signerCert; }
    public byte[] getEncCert() { return encCert; }
    public byte[] getHyperSignCert() { return HyperSignCert; }
    public byte[] getHyperEncCert() { return HyperEncCert; }
    public byte[] getEncryptedPrivateKey() { return encryptedPrivateKey; }
}
