package domain

import (
	"context"

	domainpb "github.com/hts0000/ops-backend/domain/api/gen/v1"
	"go.uber.org/zap"
)

type DomainManager interface {
	GetDomains(page, pageSize int64) (*domainpb.GetDomainsResponse, error)
}

type Service struct {
	Logger        *zap.Logger
	DomainManager DomainManager
	domainpb.UnimplementedDomainServiceServer
}

func (s *Service) GetDomains(ctx context.Context, req *domainpb.GetDomainsRequest) (*domainpb.GetDomainsResponse, error) {
	s.Logger.Info("get domains", zap.Int("page", int(req.Page)), zap.Int("pageSize", int(req.PageSize)))
	if req.Page == 0 {
		s.Logger.Warn("page is 0, set to 1")
		req.Page = 1
	}
	if req.PageSize == 0 {
		s.Logger.Warn("page size is 0, set to 10")
		req.PageSize = 10
	}
	return s.DomainManager.GetDomains(req.Page, req.PageSize)
}
