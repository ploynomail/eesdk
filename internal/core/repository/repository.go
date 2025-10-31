package repository

import (
	pb "commonprotocol/pkimessage"
	"context"
)

// PKIRepository PKI仓储接口
type PKIRepository interface {
	// SendPKIMessage 发送PKI消息
	SendPKIMessage(ctx context.Context, msg *pb.PKIMessage) (*pb.PKIMessage, error)
}

// KeyRepository 密钥仓储接口
type KeyRepository interface {
	// GenerateKeyPair 生成密钥对
	GenerateKeyPair() (privateKey, publicKey []byte, err error)

	// GenerateKemKeyPair 生成KEM密钥对
	GenerateKemKeyPair() (privateKey, publicKey []byte, err error)

	// Sign 签名
	Sign(privateKey, data []byte) ([]byte, error)

	// Verify 验证签名
	Verify(publicKey, data, signature []byte) error
}
