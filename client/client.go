package eesdk

import (
	"context"
	"crypto"
	"ee-sdk/internal/core/domain"
	"ee-sdk/internal/core/service"
	"fmt"
	"time"

	pb "commonprotocol/pkimessage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

// Client PKI客户端
type Client struct {
	conn      *grpc.ClientConn
	pkiClient *service.PKIClient
	config    *EESdkConfig
}

// EESdkConfig 客户端配置
type EESdkConfig struct {
	ServerAddr string // PKI服务器地址
	AccessId   string // 访问ID
	AccessKey  string // 访问密钥
}

// NewClient 创建新的PKI客户端
func NewClient(cfg *EESdkConfig) (*Client, error) {
	// 准备gRPC连接选项
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// 连接参数配置
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  1.0 * time.Second,
				Multiplier: 2.0,
				Jitter:     0.2,
				MaxDelay:   30 * time.Second,
			},
			MinConnectTimeout: 5 * time.Second,
		}),
	}

	// 建立gRPC连接
	conn, err := grpc.NewClient(
		cfg.ServerAddr,
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("dial server: %w", err)
	}

	// 创建PKI仓储
	pkiRepo := &grpcPKIRepository{
		client: pb.NewPKIServiceClient(conn),
	}

	// 创建密钥仓储
	keyRepo := &defaultKeyRepository{}

	// 创建PKI客户端服务
	pkiClient := service.NewPKIClient(pkiRepo, keyRepo, &domain.ClientConfig{
		AccessId:  cfg.AccessId,
		AccessKey: cfg.AccessKey,
	})

	return &Client{
		conn:      conn,
		pkiClient: pkiClient,
		config:    cfg,
	}, nil
}

// Close 关闭客户端连接
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// RequestCertificate 请求证书
func (c *Client) RequestCertificate(ctx context.Context, req *CertificateRequest) (*CertificateResponse, error) {
	// 转换请求
	domainReq := &domain.CertificateRequest{
		PublicKey:    req.PublicKey,
		KemPublicKey: req.KemPublicKey,
		CACertID:     req.CACertID,
		TemplateID:   req.TemplateID,
		Subject: domain.Subject{
			CommonName:         req.Subject.CommonName,
			Country:            req.Subject.Country,
			Organization:       req.Subject.Organization,
			OrganizationalUnit: req.Subject.OrganizationalUnit,
			Locality:           req.Subject.Locality,
			Province:           req.Subject.Province,
		},
	}

	// 调用服务
	resp, err := c.pkiClient.RequestCertificate(ctx, domainReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	return &CertificateResponse{
		SignerCert:          resp.SignerCert,
		EncCert:             resp.EncCert,
		EncryptedPrivateKey: resp.EncryptedPrivateKey,
		SerialNumber:        resp.SerialNumber,
		IssuedAt:            resp.IssuedAt,
		ExpiresAt:           resp.ExpiresAt,
	}, nil
}

// RequestCertificateWithPKCS10 使用PKCS#10请求证书
func (c *Client) RequestCertificateWithPKCS10(ctx context.Context, csrBytes []byte) (*CertificateResponse, error) {
	resp, err := c.pkiClient.RequestCertificateWithPKCS10(ctx, csrBytes)
	if err != nil {
		return nil, err
	}

	return &CertificateResponse{
		SignerCert:          resp.SignerCert,
		EncCert:             resp.EncCert,
		EncryptedPrivateKey: resp.EncryptedPrivateKey,
		SerialNumber:        resp.SerialNumber,
		IssuedAt:            resp.IssuedAt,
		ExpiresAt:           resp.ExpiresAt,
	}, nil
}

// GenerateCSR 生成证书签名请求
func (c *Client) GenerateCSR(req *CertificateRequest, privateKey crypto.Signer) ([]byte, error) {
	domainReq := &domain.CertificateRequest{
		PublicKey:    req.PublicKey,
		KemPublicKey: req.KemPublicKey,
		Subject: domain.Subject{
			CommonName:         req.Subject.CommonName,
			Country:            req.Subject.Country,
			Organization:       req.Subject.Organization,
			OrganizationalUnit: req.Subject.OrganizationalUnit,
			Locality:           req.Subject.Locality,
			Province:           req.Subject.Province,
		},
	}

	return c.pkiClient.GeneratePKCS10Request(domainReq, privateKey)
}

// RevokeCertificate 撤销证书
func (c *Client) RevokeCertificate(ctx context.Context, serialNumber string, reason string) error {
	return c.pkiClient.RevokeCertificate(ctx, &domain.RevokeRequest{
		SerialNumber: serialNumber,
		Reason:       reason,
	})
}

// ConfirmCertificate 确认证书
func (c *Client) ConfirmCertificate(ctx context.Context, certHash []byte, certReqId int32) error {
	return c.pkiClient.ConfirmCertificate(ctx, certHash, certReqId)
}
