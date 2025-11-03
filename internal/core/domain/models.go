package domain

type ClientConfig struct {
	AccessId  string // 访问ID
	AccessKey string // 访问密钥
}

// CertificateRequest 证书请求模型
type CertificateRequest struct {
	// 公钥
	PublicKey []byte
	// 主题信息
	Subject Subject
	// KEM公钥（可选，用于加密证书, 后量子专用）
	KemPublicKey []byte
	// CA证书ID（可选）
	CACertID uint
	// 证书模板ID（可选）
	TemplateID uint
	// 扩展属性
	Extensions []Extension
}

// Subject 证书主题
type Subject struct {
	CommonName         string
	Country            []string
	Organization       []string
	OrganizationalUnit []string
	Locality           []string
	Province           []string
}

// Extension 证书扩展
type Extension struct {
	ID       []int
	Critical bool
	Value    []byte
}

// CertificateResponse 证书响应模型
type CertificateResponse struct {
	// 签名证书
	SignerCert []byte
	// 加密证书（可选）
	EncCert []byte

	HyperSignerCert []byte

	HyperEncCert []byte
	// 加密的私钥信封数据（可选）
	EncryptedPrivateKey []byte
}

// PKIMessageRequest PKI消息请求
type PKIMessageRequest struct {
	// 消息体类型
	BodyType string
	// 证书请求（用于 ir/cr/kur 等）
	CertRequest *CertificateRequest
	// PKCS#10 请求（用于 p10cr）
	PKCS10Request []byte
	// 撤销请求（用于 rr）
	RevokeRequest *RevokeRequest
	// 事务ID（可选）
	TransactionID []byte
}

// RevokeRequest 撤销请求
type RevokeRequest struct {
	SerialNumber string
	Reason       string
}

// PKIMessageResponse PKI消息响应
type PKIMessageResponse struct {
	// 状态
	Status int32
	// 状态描述
	StatusString string
	// 证书响应
	CertResponse *CertificateResponse
	// 错误信息
	ErrorMessage string
}
