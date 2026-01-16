package main

import (
	"crypto/rand"
	"crypto/x509/pkix"
	eesdk "github.com/ploynomail/eesdk/client"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/open-quantum-safe/liboqs-go/oqs"
	"github.com/ploynomail/turingPQC/sm2"
	sm2dilithium2hybrid "github.com/ploynomail/turingPQC/sm2_dilithium2_hybrid"
	"github.com/ploynomail/turingPQC/x509"
	"github.com/tjfoc/gmsm/sm4"
)

var kemPublickKeyBytes []byte = make([]byte, 0)
var kemPrivateKeyBytes []byte = make([]byte, 0)

// func ReqKey() {
// 	kemPublickKeyBytes, err := base64.StdEncoding.DecodeString(kemPublickKeyBytes)
// 	if err != nil {
// 		fmt.Println("Error decoding KEM public key:", err)
// 		return
// 	}

// 	fmt.Println("EE SDK Golang Example")
// 	var cfg = &eesdk.EESdkConfig{
// 		ServerAddr: "localhost:9001",
// 		AccessId:   "807e8150-da11-48a8-886e-54e9894e8003",
// 		AccessKey:  "HanBpT2gT31Gfn83R0dTjKu0O9mv55bGDyVtP3IWOW1=",
// 	}

// 	client, err := eesdk.NewClient(cfg)
// 	if err != nil {
// 		fmt.Println("Error creating client:", err)
// 		return
// 	}
// 	fmt.Println("Client created successfully")

// 	privateKey, err := sm2dilithium2hybrid.GenerateKey(rand.Reader)
// 	if err != nil {
// 		fmt.Println("Error generating key:", err)
// 		return
// 	}

// 	publicKey := privateKey.Public()

// 	publicKeyDer, err := x509.MarshalPKIXPublicKey(publicKey)
// 	if publicKeyDer == nil || err != nil {
// 		fmt.Println("Error marshaling public key")
// 		return
// 	}

// 	certReq := &eesdk.CertificateRequest{
// 		PublicKey:       publicKeyDer,
// 		CACertID:        1,
// 		TemplateID:      1,
// 		KemPublicKey:    kemPublickKeyBytes,
// 		Subject:         eesdk.Subject{CommonName: "test-device"},
// 		Extensions:      []eesdk.Extension{},
// 		ExtraExtensions: []eesdk.Extension{},
// 	}

// 	fmt.Println("Requesting certificate...")

// 	rep, err := client.RequestCertificate(certReq)
// 	if err != nil {
// 		fmt.Println("Error requesting certificate:", err)
// 		return
// 	}
// 	fmt.Println("Certificate requested successfully:")
// 	fmt.Printf("Hyper Signer Cert Length: %d bytes\n", len(rep.HyperSignerCert))
// 	fmt.Printf("Hyper Enc Cert Length: %d bytes\n", len(rep.HyperEncCert))
// 	fmt.Printf("Encrypted Private Key Length: %d bytes\n", len(rep.EncryptedPrivateKey))
// 	defer client.Close()
// }

func GenerateKemKeyPair() {
	// 指定要使用的 KEM 算法，这里以 Kyber512 为例
	// 其他可选算法如 "Kyber768", "Kyber1024", 具体取决于 liboqs 的编译配置
	kemAlgorithm := "Kyber512"
	fmt.Printf("\n正在使用 KEM 算法: %s\n", kemAlgorithm)

	// 1. 接收端生成密钥对
	fmt.Println("\n*** 步骤 1: 接收端生成密钥对 ***")
	receiver := oqs.KeyEncapsulation{}
	// defer receiver.Clean() // 确保函数退出时释放原生资源

	// 初始化接收端的 KEM 上下文，使用指定的算法
	if err := receiver.Init(kemAlgorithm, nil); err != nil {
		log.Fatalf("初始化接收端 KEM 失败: %v", err)
	}

	// 生成公钥(pk)和私钥(sk)
	receiverPublicKey, err := receiver.GenerateKeyPair()
	if err != nil {
		log.Fatalf("生成密钥对失败: %v", err)
	}
	fmt.Printf("接收端公钥长度: %d 字节\n", len(receiverPublicKey))
	fmt.Printf("接收端公钥：%s\n", base64.StdEncoding.EncodeToString(receiverPublicKey))
	// 注意：私钥由 receiver 对象内部管理
	kemPublickKeyBytes = receiverPublicKey
	kemPrivateKeyBytes = receiver.ExportSecretKey()
}

func CreateCSR(id string) {

	// 生成私钥
	privateKeyPath := fmt.Sprintf("%s.key", id)
	privateKey, _ := sm2dilithium2hybrid.GenerateKey(rand.Reader)
	privateKeyDER, _ := x509.MarshalPKCS8PrivateKey(privateKey)
	p := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyDER})
	os.WriteFile(privateKeyPath, p, 0777)
	// 创建CSR模板
	template := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: id,
		},
		SignatureAlgorithm: x509.PureSm2Hybrid,
		ExtraExtensions: []pkix.Extension{
			{
				Id:       eesdk.OidExtensionKemPublicKey,
				Critical: false,
				Value:    kemPublickKeyBytes,
			},
		},
	}
	// 创建CSR
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, privateKey)
	if err != nil {
		panic(err)
	}

	// 将CSR序列化为PEM格式
	csrPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	})
	csrPath := "../testdata/example.csr"
	// 将CSR写入文件
	err = os.WriteFile(csrPath, csrPEM, 0644)
	if err != nil {
		panic(err)
	}
}

func ReqKeyWithCSR() {
	// read csr
	csrPemBytes, err := os.ReadFile("../testdata/example.csr")
	if err != nil {
		fmt.Println("Error reading CSR file:", err)
		return
	}

	fmt.Println("EE SDK Golang Example - Request Key with CSR")
	var cfg = &eesdk.EESdkConfig{
		ServerAddr: "localhost:9001",
		AccessId:   "807e8150-da11-48a8-886e-54e9894e8003",
		AccessKey:  "HanBpT2gT31Gfn83R0dTjKu0O9mv55bGDyVtP3IWOW1=",
	}

	client, err := eesdk.NewClient(cfg)
	if err != nil {
		fmt.Println("Error creating client:", err)
		return
	}
	fmt.Println("Client created successfully")

	defer client.Close()
	fmt.Println("Requesting certificate with CSR...")
	rep, err := client.RequestCertificateWithPKCS10(csrPemBytes)
	if err != nil {
		fmt.Println("Error requesting certificate:", err)
		return
	}

	fmt.Println("Certificate requested successfully:")
	fmt.Printf("Signer Cert Length: %d bytes\n", len(rep.SignerCert))
	fmt.Printf("Enc Cert Length: %d bytes\n", len(rep.EncCert))
	fmt.Printf("Hyper Signer Cert Length: %d bytes\n", len(rep.HyperSignerCert))
	SavePemToFile("hyper_signer_cert.pem", rep.HyperSignerCert)
	fmt.Printf("Hyper Enc Cert Length: %d bytes\n", len(rep.HyperEncCert))
	SavePemToFile("hyper_enc_cert.pem", rep.HyperEncCert)
	fmt.Printf("Encrypted Private Key Length: %d bytes\n", len(rep.EncryptedPrivateKey))

	fmt.Println("Certificates and Encrypted Private Key received:")
	fmt.Printf("Encrypted Private Key (Base64): %s\n", base64.StdEncoding.EncodeToString(rep.EncryptedPrivateKey))
	DeCryptPrivateKey(rep.EncryptedPrivateKey)
}

func DeCryptPrivateKey(encryptedPrivateKeyBytes []byte) {
	var encryptedValue eesdk.EncryptedValue
	asn1.Unmarshal(encryptedPrivateKeyBytes, &encryptedValue)

	// 使用 liboqs 进行解密
	kem := oqs.KeyEncapsulation{}
	defer kem.Clean()

	// 初始化 KEM 上下文，使用与加密时相同的算法
	kemAlgorithm := "Kyber512"
	if err := kem.Init(kemAlgorithm, nil); err != nil {
		log.Fatalf("初始化 KEM 失败: %v", err)
	}

	// 设置接收端的私钥
	if err := kem.Init(kemAlgorithm, kemPrivateKeyBytes); err != nil {
		log.Fatalf("设置私钥失败: %v", err)
	}

	// 解密对称密钥
	symmetricKey, err := kem.DecapSecret(encryptedValue.EncSymmKey.Bytes)
	if err != nil {
		log.Fatalf("解密对称密钥失败: %v", err)
	}
	fmt.Printf("解密得到的对称密钥长度: %d 字节\n", len(symmetricKey))

	key := symmetricKey[:sm4.BlockSize]
	// 使用对称密钥解密私钥数据
	ecbDec, err := sm4.Sm4Ecb(key, encryptedValue.EncValue.Bytes, false)
	if err != nil {
		log.Fatalf("使用对称密钥解密私钥数据失败: %v", err)
	}
	var sm2PrivateKey sm2.PrivateKey = sm2.PrivateKey{
		D: new(big.Int).SetBytes(ecbDec),
	}
	sm2PrivateKey.PublicKey = sm2.PublicKey{
		Curve: sm2.P256Sm2(),
	}
	sm2PrivateKey.PublicKey.X, sm2PrivateKey.PublicKey.Y = sm2.P256Sm2().ScalarBaseMult(sm2PrivateKey.D.Bytes())
	fmt.Printf("解密得到的SM2私钥D值: %s\n", sm2PrivateKey.D)
	pemBytes, err := x509.MarshalPKCS8PrivateKey(&sm2PrivateKey)
	if err != nil {
		log.Fatalf("将SM2私钥编码为PKCS#8格式失败: %v", err)
	}
	pemBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pemBytes,
	}
	pemData := pem.EncodeToMemory(pemBlock)
	SavePemToFile("decrypted_sm2_private_key.pem", pemData)
	fmt.Printf("解密得到的SM2私钥PEM格式:\n%s\n", string(pemData))
}

func RevokeCertificate() {
	fmt.Println("EE SDK Golang Example")
	var cfg = &eesdk.EESdkConfig{
		ServerAddr: "localhost:9001",
		AccessId:   "807e8150-da11-48a8-886e-54e9894e8003",
		AccessKey:  "HanBpT2gT31Gfn83R0dTjKu0O9mv55bGDyVtP3IWOW1=",
	}

	client, err := eesdk.NewClient(cfg)
	if err != nil {
		fmt.Println("Error creating client:", err)
		return
	}
	fmt.Println("Client created successfully")
	client.RevokeCertificate("74", "keyCompromise")
}

func SavePemToFile(filename string, pemBytes []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(pemBytes)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// ReqKey()
	GenerateKemKeyPair()
	CreateCSR("openvpn")
	ReqKeyWithCSR()
	// RevokeCertificate()
}
