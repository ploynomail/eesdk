package eesdk

import "fmt"

// defaultKeyRepository 默认密钥仓储实现
type defaultKeyRepository struct{}

func (r *defaultKeyRepository) GenerateKeyPair() (privateKey, publicKey []byte, err error) {
	// TODO: 实现密钥对生成
	return nil, nil, fmt.Errorf("not implemented")
}

func (r *defaultKeyRepository) GenerateKemKeyPair() (privateKey, publicKey []byte, err error) {
	// TODO: 实现KEM密钥对生成
	return nil, nil, fmt.Errorf("not implemented")
}

func (r *defaultKeyRepository) Sign(privateKey, data []byte) ([]byte, error) {
	// TODO: 实现签名
	return nil, fmt.Errorf("not implemented")
}

func (r *defaultKeyRepository) Verify(publicKey, data, signature []byte) error {
	// TODO: 实现验证
	return fmt.Errorf("not implemented")
}
