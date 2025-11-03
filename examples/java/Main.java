import com.turingq.eesdk.*;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.Base64;

/**
 * EE SDK Java使用示例 - PKI证书管理
 */
public class Main {
    
    /**
     * 示例1: 请求证书
     */
    private static void exampleRequestCertificate(EESdk client) {
        System.out.println("\n=== 请求证书示例 ===");
        
        try {
            // 模拟公钥和KEM公钥（实际应用中应该从密钥生成获取）
            byte[] publicKey = new byte[65];
            publicKey[0] = 0x04; // 未压缩点标识
            
            String kemPublicKeyBase64 = "HkSLs9sGqvNIUio0yLjE90wa7nVAkBKtyTjPd8ygSXURreDNVeRylEgQDydrqzgE9lZtLufJ/Wt5QuwWndaA/Et31WqgkWSQlnd1JHaQCXaFAZbIgxbKT3l90EYSY8eoYqFWTYZ7BLVy6jB3geU+hcYmr0IBiUQ2BFYFMYN79/Y8Z8IIB7E9UEgzrkmRRtYaLUUy1dpWo8JF+qiToqIRErMq/ulrazyog1mZeMGq4cRZFiNsOqRRfullrmAvpXM98quGlyAawZihiZWG+cJj1+W6iCUMB9qkkZWgoOZxuVinqAAnzceWyTYTE8tLUSuOQ4nLsyuCN5MY24CUdDluIRYmLtBwOZRT9OpfuJoNQ3dg9magj6IF0htNRms6FHluukQJcIw3eRURZVqLDvGG3TYxkJUcXusDcokvnqjPcvNY8hMxpNOTF/mVtdsL7exuG8oOi+UxzyBn2TkSROpjpwKgQRghp5WDMFtRP0QxkZc2iESBsDcXDihnPFSLgWwNxwPOxchKJ9kQx1QXrQoqMuXJXzM3Dsh3ZntWt2GNENMYVkcKAFQ/ZARs41glvMAWd4UCN2jIWuOGXnEzJQs0D0Y++7efMkEt1+k175p0wQpiE+iUCbSBEJILq5iM9RB2dTwQ8Zqdf+lYE0pTYqdURXdlCtIp+pdIyXVnAKNPbrHH2AbK8KyALmdih7Q+qHc4+hM3FMm6YLKS4UOit5McpKuvd+pwT5EvufB1DaG7+wpP/PYtg0CBXFtqnakSNeREHOl0NPijQzgt59VV1UN4yFwR82cbgQypfwd6q4V6+bMlNbqa4MtJAkSgbKRPRyAS4mh9oqUICQO/5ZPFwmGKZdsuipaZhRWkGSJaZxOhyAOyXLoOeTcXmbu2iFt8t9nIvzyDixO9TsQHI1M3hdXBrfpULWNEbzc1ugJB2ccQ0pm6dFkCl0oahCxXxUxFpHcUz8colDgaSXAgCzA9zAtbA9JmVnGXuWfPUZu0pQBZyYUr0ONKYnYOCvUOIOdJPLeqAJm96ABZCKZ/16sS8AGTHP2mmi2xxQ9IL6biZGiS0TI=";
            byte[] kemPublicKey = Base64.getDecoder().decode(kemPublicKeyBase64);
            
            // 创建证书主题
            Subject subject = new Subject("test-device");
            subject.setCountry(new String[]{"CN"});
            subject.setOrganization(new String[]{"TuringQ"});
            subject.setOrganizationalUnit(new String[]{"Engineering"});
            subject.setLocality(new String[]{"Beijing"});
            subject.setProvince(new String[]{"Beijing"});
            
            // 创建证书请求
            CertificateRequest request = new CertificateRequest(publicKey, 1, 1, subject);
            request.setKemPublicKey(kemPublicKey);
            
            // 请求证书
            CertificateResponse response = client.requestCertificate(request);
            
            System.out.println("证书请求成功!");
            System.out.println("签名证书长度: " + response.getSignerCert().length + " bytes");
            System.out.println("加密证书长度: " + response.getEncCert().length + " bytes");
            System.out.println("加密私钥长度: " + response.getEncryptedPrivateKey().length + " bytes");
            
            // 保存证书到文件
            Files.write(Paths.get("/tmp/signer_cert.der"), response.getSignerCert());
            Files.write(Paths.get("/tmp/enc_cert.der"), response.getEncCert());
            System.out.println("证书已保存到 /tmp/signer_cert.der 和 /tmp/enc_cert.der");
            
        } catch (Exception e) {
            System.err.println("证书请求失败: " + e.getMessage());
            e.printStackTrace();
        }
    }
    
    /**
     * 示例2: 使用PKCS#10 CSR请求证书
     */
    private static void exampleRequestCertificateWithPKCS10(EESdk client) {
        System.out.println("\n=== 使用PKCS#10请求证书示例 ===");
        
        try {
            // 读取CSR文件
            String csrPath = "../testdata/example.csr";
            System.out.println("读取CSR文件: " + csrPath);
            byte[] csrPem = Files.readAllBytes(Paths.get(csrPath));
            System.out.println("CSR文件大小: " + csrPem.length + " bytes");
            
            // 使用CSR请求证书
            System.out.println("调用 requestCertificateWithPKCS10...");
            CertificateResponse response = client.requestCertificateWithPKCS10(csrPem);
            
            System.out.println("CSR证书请求成功!");
            System.out.println("签名证书长度: " + response.getSignerCert().length + " bytes");
            System.out.println("加密证书长度: " + response.getEncCert().length + " bytes");
            System.out.println("加密私钥长度: " + response.getEncryptedPrivateKey().length + " bytes");
            System.out.println("HyperSign证书: " + response.getHyperSignCert().length + " bytes");
            System.out.println("HyperEnc证书: " + response.getHyperEncCert().length + " bytes");
            
        } catch (IOException e) {
            System.err.println("CSR文件读取失败: " + e.getMessage());
            e.printStackTrace();
        } catch (UnsatisfiedLinkError e) {
            System.err.println("JNI链接错误: " + e.getMessage());
            System.err.println("请检查 JNI 方法签名是否正确");
            e.printStackTrace();
        } catch (Exception e) {
            System.err.println("CSR证书请求失败: " + e.getMessage());
            e.printStackTrace();
        }
    }
    
    /**
     * 示例3: 撤销证书
     */
    private static void exampleRevokeCertificate(EESdk client) {
        System.out.println("\n=== 撤销证书示例 ===");
        
        try {
            // 撤销证书（使用示例序列号）
            String serialNumber = "74";
            String reason = "keyCompromise";
            
            boolean success = client.revokeCertificate(serialNumber, reason);
            
            if (success) {
                System.out.println("证书撤销成功! 序列号: " + serialNumber + ", 原因: " + reason);
            } else {
                System.out.println("证书撤销失败");
            }
            
        } catch (Exception e) {
            System.err.println("证书撤销失败: " + e.getMessage());
            e.printStackTrace();
        }
    }
    
    /**
     * 示例4: 确认证书
     */
    private static void exampleConfirmCertificate(EESdk client) {
        System.out.println("\n=== 确认证书示例 ===");
        
        try {
            // 模拟证书哈希（实际应用中应该从证书计算）
            byte[] certHash = new byte[32];
            int certReqId = 123;
            
            boolean success = client.confirmCertificate(certHash, certReqId);
            
            if (success) {
                System.out.println("证书确认成功! 请求ID: " + certReqId);
            } else {
                System.out.println("证书确认失败");
            }
            
        } catch (Exception e) {
            System.err.println("证书确认失败: " + e.getMessage());
            e.printStackTrace();
        }
    }
    
    public static void main(String[] args) {
        System.out.println("=== EE SDK Java PKI示例 ===");
        
        // 打印库加载信息
        System.out.println("Java Library Path: " + System.getProperty("java.library.path"));
        System.out.println("Loading native library...");
        
        EESdk client = null;
        try {
            // 创建客户端
            client = new EESdk(
                "localhost:9001",
                "807e8150-da11-48a8-886e-54e9894e8003",
                "HanBpT2gT31Gfn83R0dTjKu0O9mv55bGDyVtP3IWOW1="
            );
            
            System.out.println("SDK客户端创建成功");
            
            // 运行示例
            // exampleRequestCertificate(client);
            exampleRequestCertificateWithPKCS10(client);
            // exampleRevokeCertificate(client);
            // exampleConfirmCertificate(client);
            
        } catch (Exception e) {
            System.err.println("SDK使用过程中发生错误: " + e.getMessage());
            e.printStackTrace();
        } finally {
            // 清理资源
            if (client != null) {
                client.destroy();
            }
        }
        
        System.out.println("\n所有示例执行完成");
    }
}
