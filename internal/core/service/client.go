package service

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"time"

	pb "commonprotocol/pkimessage"
	"ee-sdk/internal/core/domain"
	"ee-sdk/internal/core/repository"

	utils "commonprotocol/utils"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/ploynomail/turingPQC/x509"
)

// PKIClient PKI客户端服务
type PKIClient struct {
	pkiRepo repository.PKIRepository
	keyRepo repository.KeyRepository
	config  *domain.ClientConfig
}

// NewPKIClient 创建PKI客户端
func NewPKIClient(pkiRepo repository.PKIRepository, keyRepo repository.KeyRepository, config *domain.ClientConfig) *PKIClient {
	return &PKIClient{
		pkiRepo: pkiRepo,
		keyRepo: keyRepo,
		config:  config,
	}
}

// RequestCertificate 请求证书（初始化请求 ir）
func (c *PKIClient) RequestCertificate(ctx context.Context, req *domain.CertificateRequest) (*domain.CertificateResponse, error) {
	// 构建证书请求消息
	certReqMsg, err := c.buildCertReqMessage(req)
	if err != nil {
		return nil, fmt.Errorf("build cert request message: %w", err)
	}
	h, err := c.buildPKIHeader()
	if err != nil {
		return nil, fmt.Errorf("build pki header: %w", err)
	}
	// 构建PKI消息
	pkiMsg := &pb.PKIMessage{
		Header: h,
		Body: &pb.PKIBody{
			Body: &pb.PKIBody_Ir{
				Ir: &pb.CertReqMessages{
					Messages: []*pb.CertReqMessage{certReqMsg},
				},
			},
		},
	}

	// 添加保护
	if err := c.addProtection(pkiMsg); err != nil {
		return nil, fmt.Errorf("add protection: %w", err)
	}

	// 发送消息
	respMsg, err := c.pkiRepo.SendPKIMessage(ctx, pkiMsg)
	if err != nil {
		return nil, fmt.Errorf("send PKI message: %w", err)
	}
	// 验证响应保护
	if err := c.verifyProtection(respMsg); err != nil {
		return nil, fmt.Errorf("verify response protection: %w", err)
	}

	// 解析响应
	return c.parseCertResponse(respMsg)
}

// RequestCertificateWithPKCS10 使用PKCS#10请求证书
func (c *PKIClient) RequestCertificateWithPKCS10(ctx context.Context, p10Request []byte) (*domain.CertificateResponse, error) {
	h, err := c.buildPKIHeader()
	if err != nil {
		return nil, fmt.Errorf("build pki header: %w", err)
	}
	// 构建PKI消息
	pkiMsg := &pb.PKIMessage{
		Header: h,
		Body: &pb.PKIBody{
			Body: &pb.PKIBody_P10Cr{
				P10Cr: p10Request,
			},
		},
	}

	// 添加保护
	if err := c.addProtection(pkiMsg); err != nil {
		return nil, fmt.Errorf("add protection: %w", err)
	}

	// 发送消息
	respMsg, err := c.pkiRepo.SendPKIMessage(ctx, pkiMsg)
	if err != nil {
		return nil, fmt.Errorf("send PKI message: %w", err)
	}

	// 验证响应保护
	if err := c.verifyProtection(respMsg); err != nil {
		return nil, fmt.Errorf("verify response protection: %w", err)
	}

	// 解析响应
	return c.parseCertResponse(respMsg)
}

// RevokeCertificate 撤销证书
func (c *PKIClient) RevokeCertificate(ctx context.Context, req *domain.RevokeRequest) error {
	// 构建撤销请求
	revDetails := &pb.RevDetails{
		CertDetails: &pb.CertTemplate{
			SerialNumber: []byte(req.SerialNumber),
		},
	}

	h, err := c.buildPKIHeader()
	if err != nil {
		return fmt.Errorf("build pki header: %w", err)
	}
	// 构建PKI消息
	pkiMsg := &pb.PKIMessage{
		Header: h,
		Body: &pb.PKIBody{
			Body: &pb.PKIBody_Rr{
				Rr: &pb.RevReqContent{
					RevDetails: []*pb.RevDetails{revDetails},
				},
			},
		},
	}

	// 添加保护
	if err := c.addProtection(pkiMsg); err != nil {
		return fmt.Errorf("add protection: %w", err)
	}

	// 发送消息
	respMsg, err := c.pkiRepo.SendPKIMessage(ctx, pkiMsg)
	if err != nil {
		return fmt.Errorf("send PKI message: %w", err)
	}

	// 验证响应保护
	if err := c.verifyProtection(respMsg); err != nil {
		return fmt.Errorf("verify response protection: %w", err)
	}

	// 检查响应状态
	if respMsg.Body.GetRp() == nil {
		return fmt.Errorf("invalid revoke response")
	}

	rp := respMsg.Body.GetRp()
	if len(rp.Status) == 0 || rp.Status[0].Status != 0 {
		return fmt.Errorf("revoke failed: %s", rp.Status[0].StatusString)
	}

	return nil
}

// ConfirmCertificate 确认证书
func (c *PKIClient) ConfirmCertificate(ctx context.Context, certHash []byte, certReqId int32) error {
	h, err := c.buildPKIHeader()
	if err != nil {
		return fmt.Errorf("build pki header: %w", err)
	}
	// 构建PKI消息
	pkiMsg := &pb.PKIMessage{
		Header: h,
		Body: &pb.PKIBody{
			Body: &pb.PKIBody_CertConf{
				CertConf: &pb.CertConfirmContent{
					CertStatus: []*pb.CertStatus{
						{
							CertHash:  certHash,
							CertReqId: certReqId,
							StatusInfo: &pb.PKIStatusInfo{
								Status: 0,
							},
						},
					},
				},
			},
		},
	}

	// 添加保护
	if err := c.addProtection(pkiMsg); err != nil {
		return fmt.Errorf("add protection: %w", err)
	}

	// 发送消息
	respMsg, err := c.pkiRepo.SendPKIMessage(ctx, pkiMsg)
	if err != nil {
		return fmt.Errorf("send PKI message: %w", err)
	}

	// 验证响应保护
	if err := c.verifyProtection(respMsg); err != nil {
		return fmt.Errorf("verify response protection: %w", err)
	}

	// 检查响应
	if respMsg.Body.GetConf() == nil {
		return fmt.Errorf("invalid confirm response")
	}

	return nil
}

// buildCertReqMessage 构建证书请求消息
func (c *PKIClient) buildCertReqMessage(req *domain.CertificateRequest) (*pb.CertReqMessage, error) {
	// 序列化公钥

	// 构建主题
	subject := pkix.Name{
		CommonName:         req.Subject.CommonName,
		Country:            req.Subject.Country,
		Organization:       req.Subject.Organization,
		OrganizationalUnit: req.Subject.OrganizationalUnit,
		Locality:           req.Subject.Locality,
		Province:           req.Subject.Province,
	}
	subjectBytes, err := asn1.Marshal(subject)
	if err != nil {
		return nil, fmt.Errorf("marshal subject: %w", err)
	}

	// 构建扩展（包含KEM公钥）
	var extensions [][]byte
	if len(req.KemPublicKey) > 0 {
		kemExt, err := c.buildKemExtension(req.KemPublicKey)
		if err != nil {
			return nil, fmt.Errorf("build kem extension: %w", err)
		}
		extensions = append(extensions, kemExt)
	}

	// 构建证书模板
	certTemplate := &pb.CertTemplate{
		Subject: subjectBytes,
		SubjectPublicKeyInfo: &pb.SubjectPublicKeyInfo{
			Algorithm: &pb.PKIAlgorithmIdentifier{
				AlgorithmOid: "1.2.156.10197.1.301", // TODO: 这里不正确
			},
			SubjectPublicKey: req.PublicKey,
		},
		Extensions: extensions,
	}

	// 构建控制信息
	controls := c.buildControls(req.CACertID, req.TemplateID)

	return &pb.CertReqMessage{
		CertReqId:    1,
		CertTemplate: certTemplate,
		Controls:     controls,
	}, nil
}

// buildKemExtension 构建KEM扩展
func (c *PKIClient) buildKemExtension(kemKey []byte) ([]byte, error) {
	ext := struct {
		Id       asn1.ObjectIdentifier
		Critical bool `asn1:"optional"`
		Value    []byte
	}{
		Id:    asn1.ObjectIdentifier{1, 2, 156, 10197, 1, 501}, // KEM公钥扩展OID
		Value: kemKey,
	}

	return asn1.Marshal(ext)
}

// buildControls 构建控制信息
func (c *PKIClient) buildControls(caCertID, templateID uint) *pb.Controls {
	if caCertID == 0 && templateID == 0 {
		return nil
	}

	var controls []*pb.AttributeTypeAndValue

	// CA证书ID
	if caCertID > 0 {
		atv := struct {
			Type  asn1.ObjectIdentifier
			Value int
		}{
			Type:  asn1.ObjectIdentifier{1, 2, 3, 4, 1},
			Value: int(caCertID),
		}
		atvBytes, _ := asn1.Marshal(atv)
		controls = append(controls, &pb.AttributeTypeAndValue{Raw: atvBytes})
	}

	// 模板ID
	if templateID > 0 {
		atv := struct {
			Type  asn1.ObjectIdentifier
			Value int
		}{
			Type:  asn1.ObjectIdentifier{1, 2, 3, 4, 2},
			Value: int(templateID),
		}
		atvBytes, _ := asn1.Marshal(atv)
		controls = append(controls, &pb.AttributeTypeAndValue{Raw: atvBytes})
	}

	if len(controls) == 0 {
		return nil
	}

	return &pb.Controls{Controls: controls}
}

// buildPKIHeader 构建PKI消息头
func (c *PKIClient) buildPKIHeader() (*pb.PKIHeader, error) {
	now := time.Now()
	nonce := make([]byte, 16)
	rand.Read(nonce)

	// 获取 SenderKid（从配置中获取）
	var senderKid []byte
	if c.config != nil && c.config.AccessId != "" {
		u, err := pb.UUIDToBytes(c.config.AccessId)
		if err != nil {
			return nil, err
		}
		senderKid = u
	}

	return &pb.PKIHeader{
		Pvno: 1,
		Sender: &pb.PKIGeneralName{
			Value: []byte("CN=EE Client"),
		},
		Recipient: &pb.PKIGeneralName{
			Value: []byte("CN=CA Server"),
		},
		SenderKid: senderKid,
		RecipKid:  senderKid,
		MessageTime: &timestamp.Timestamp{
			Seconds: now.Unix(),
			Nanos:   int32(now.Nanosecond()),
		},
		TransactionId: nonce,
		SenderNonce:   nonce,
	}, nil
}

// addProtection 添加消息保护
func (c *PKIClient) addProtection(msg *pb.PKIMessage) error {
	// 1. 生成 PBMParameter
	salt, err := pb.GenerateNonce()
	if err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	pbmParam := &pb.PBMParameter{
		Salt: salt,
		Owf: &pb.PKIAlgorithmIdentifier{
			AlgorithmOid: utils.OidSM3.String(),
		},
		IterationCount: 30, // 迭代次数
		Mac: &pb.PKIAlgorithmIdentifier{
			AlgorithmOid: utils.OidSM3WithKey.String(), // Updated to use the correct OID for HMAC
		},
	}

	// 2. 设置保护算法
	msg.Header.ProtectionAlg = &pb.PKIAlgorithmIdentifier{
		AlgorithmOid: utils.OidPasswordBasedMAC.String(), // id-PasswordBasedMAC OID
		Parameters: &pb.PKIAlgorithmIdentifier_PbmParameter{
			PbmParameter: pbmParam,
		},
	}

	// 3. 构建待保护的数据 (DER编码的 header + body)
	protectedPart := &pb.ProtectedPard{
		Header: msg.Header,
		Body:   msg.Body,
	}

	protectedData, err := pb.EncodeToDER(protectedPart)
	if err != nil {
		return fmt.Errorf("failed to encode protected part: %w", err)
	}

	// 4. 获取共享密钥（从配置中获取）
	sharedSecret := c.getSharedSecret()
	if sharedSecret == nil {
		return fmt.Errorf("shared secret not configured")
	}

	// 5. 计算 MAC
	mac := pb.CalculateMAC(protectedData, pbmParam, sharedSecret)

	// 6. 将 MAC 写入 Protection 字段
	msg.Protection = mac

	return nil
}

// verifyProtection 验证消息保护
func (c *PKIClient) verifyProtection(msg *pb.PKIMessage) error {
	// 1. 检查保护算法是否存在
	if msg.Header.ProtectionAlg == nil || msg.Header.ProtectionAlg.AlgorithmOid == "" {
		return fmt.Errorf("protection algorithm is missing")
	}

	// 2. 获取 PBMParameter
	pbmParam, ok := msg.Header.ProtectionAlg.Parameters.(*pb.PKIAlgorithmIdentifier_PbmParameter)
	if !ok || pbmParam.PbmParameter == nil {
		return fmt.Errorf("PBM parameter is missing or invalid")
	}

	// 3. 构建待保护的数据 (DER编码的 header + body)
	protectedPart := &pb.ProtectedPard{
		Header: msg.Header,
		Body:   msg.Body,
	}

	protectedData, err := pb.EncodeToDER(protectedPart)
	if err != nil {
		return fmt.Errorf("failed to encode protected part: %w", err)
	}

	// 4. 获取共享密钥
	sharedSecret := c.getSharedSecret()
	if sharedSecret == nil {
		return fmt.Errorf("shared secret not configured or shared secret is nil")
	}

	// 5. 验证 MAC
	if !pb.VerifyMAC(protectedData, msg.Protection, pbmParam.PbmParameter, sharedSecret) {
		return fmt.Errorf("HMAC verification failed")
	}

	return nil
}

// getSharedSecret 获取共享密钥
func (c *PKIClient) getSharedSecret() []byte {
	if c.config == nil || c.config.AccessKey == "" {
		return nil
	}
	key, _ := base64.StdEncoding.DecodeString(c.config.AccessKey)
	return key
}

// parseCertResponse 解析证书响应
func (c *PKIClient) parseCertResponse(msg *pb.PKIMessage) (*domain.CertificateResponse, error) {
	// 获取响应体
	var certRepMsg *pb.CertRepMessage

	switch body := msg.Body.Body.(type) {
	case *pb.PKIBody_Ip:
		certRepMsg = body.Ip
	case *pb.PKIBody_Cp:
		certRepMsg = body.Cp
	case *pb.PKIBody_Kup:
		certRepMsg = body.Kup
	default:
		return nil, fmt.Errorf("unexpected response type")
	}

	if len(certRepMsg.Response) == 0 {
		return nil, fmt.Errorf("no certificate response")
	}

	resp := certRepMsg.Response[0]

	// 检查状态
	if resp.Status.Status != 0 {
		return nil, fmt.Errorf("certificate request failed: %s", resp.Status.StatusString)
	}

	// 解析证书
	certPair := resp.CertifiedKeyPair
	if certPair == nil || certPair.CertOrEncCert == nil {
		return nil, fmt.Errorf("no certificate in response")
	}

	var signerCert []byte
	switch cert := certPair.CertOrEncCert.Choice.(type) {
	case *pb.CertOrEncCert_Certificate:
		signerCert = c.extractCertBytes(cert.Certificate)
	default:
		return nil, fmt.Errorf("unexpected certificate type")
	}

	// 解析证书以获取序列号和有效期
	parsedCert, err := c.parseCertificate(signerCert)
	if err != nil {
		return nil, fmt.Errorf("parse certificate: %w", err)
	}

	result := &domain.CertificateResponse{
		SignerCert:   signerCert,
		SerialNumber: parsedCert.SerialNumber.String(),
		IssuedAt:     parsedCert.NotBefore,
		ExpiresAt:    parsedCert.NotAfter,
	}

	// 提取加密私钥
	if certPair.PrivateKey != nil {
		result.EncryptedPrivateKey = certPair.PrivateKey.EncValue
	}

	return result, nil
}

// extractCertBytes 提取证书字节
func (c *PKIClient) extractCertBytes(cert *pb.CMPCertificate) []byte {
	switch certType := cert.Cert.(type) {
	case *pb.CMPCertificate_X509V3PkCert:
		return certType.X509V3PkCert
	case *pb.CMPCertificate_Sm2Cert:
		return certType.Sm2Cert
	case *pb.CMPCertificate_Sm2HyperCert:
		return certType.Sm2HyperCert
	default:
		return nil
	}
}

// parseCertificate 解析证书
func (c *PKIClient) parseCertificate(certBytes []byte) (*x509.Certificate, error) {
	// 尝试直接解析
	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		// 尝试PEM格式
		block, _ := pem.Decode(certBytes)
		if block == nil {
			return nil, fmt.Errorf("invalid certificate format")
		}
		cert, err = x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
	}
	return cert, nil
}

// GeneratePKCS10Request 生成PKCS#10证书请求
func (c *PKIClient) GeneratePKCS10Request(req *domain.CertificateRequest, privateKey crypto.Signer) ([]byte, error) {
	// 构建主题
	subject := pkix.Name{
		CommonName:         req.Subject.CommonName,
		Country:            req.Subject.Country,
		Organization:       req.Subject.Organization,
		OrganizationalUnit: req.Subject.OrganizationalUnit,
		Locality:           req.Subject.Locality,
		Province:           req.Subject.Province,
	}

	// 构建扩展
	var extensions []pkix.Extension
	for _, ext := range req.Extensions {
		extensions = append(extensions, pkix.Extension{
			Id:       ext.ID,
			Critical: ext.Critical,
			Value:    ext.Value,
		})
	}

	// 添加KEM公钥扩展
	if len(req.KemPublicKey) > 0 {
		extensions = append(extensions, pkix.Extension{
			Id:    asn1.ObjectIdentifier{1, 2, 156, 10197, 1, 501},
			Value: req.KemPublicKey,
		})
	}

	// 创建证书请求模板
	template := &x509.CertificateRequest{
		Subject:    subject,
		Extensions: extensions,
	}

	// 生成CSR
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, template, privateKey)
	if err != nil {
		return nil, fmt.Errorf("create certificate request: %w", err)
	}

	return csrBytes, nil
}
