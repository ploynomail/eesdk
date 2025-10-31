package eesdk

import (
	pb "commonprotocol/pkimessage"
	"context"
)

// grpcPKIRepository gRPC PKI实现
type grpcPKIRepository struct {
	client pb.PKIServiceClient
}

func (r *grpcPKIRepository) SendPKIMessage(ctx context.Context, msg *pb.PKIMessage) (*pb.PKIMessage, error) {
	return r.client.SendPKIMessage(ctx, msg)
}
