package eesdk

// CertificateRequest 证书请求（公开API）
type CertificateRequest struct {
	PublicKey       []byte
	KemPublicKey    []byte
	CACertID        uint
	TemplateID      uint
	Subject         Subject
	Extensions      []Extension
	ExtraExtensions []Extension
}

// Subject 证书主题（公开API）
type Subject struct {
	CommonName         string
	Country            []string
	Organization       []string
	OrganizationalUnit []string
	Locality           []string
	Province           []string
}

// Extension 证书扩展（公开API）
type Extension struct {
	ID       []int
	Critical bool
	Value    []byte
}

// CertificateResponse 证书响应（公开API）
type CertificateResponse struct {
	SignerCert          []byte
	EncCert             []byte
	HyperSignerCert     []byte
	HyperEncCert        []byte
	EncryptedPrivateKey []byte
}
