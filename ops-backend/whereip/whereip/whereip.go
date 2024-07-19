package whereip

import (
	"context"

	whereippb "github.com/hts0000/ops-backend/whereip/api/gen/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	Logger *zap.Logger
	whereippb.UnimplementedWhereipServiceServer
}

func (s *Service) GetIp(ctx context.Context, req *whereippb.GetIpRequest) (*whereippb.GetIpResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetIp not implemented")
}

func (s *Service) GetIps(ctx context.Context, req *whereippb.GetIpsRequest) (*whereippb.GetIpsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetIp not implemented")
}
