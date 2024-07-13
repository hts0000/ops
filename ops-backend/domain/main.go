package main

import (
	"log"
	"os"

	"github.com/hts0000/ops-backend/domain/alicloud"
	domainpb "github.com/hts0000/ops-backend/domain/api/gen/v1"
	"github.com/hts0000/ops-backend/domain/domain"
	"github.com/hts0000/ops-backend/shared/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	aKey, hasKey := os.LookupEnv("ALICLOUD_ACCESS_KEY")
	if !hasKey {
		logger.Fatal("cannot get access key")
	}

	aSecret, hasSecret := os.LookupEnv("ALICLOUD_ACCESS_SECRET")
	if !hasSecret {
		logger.Fatal("cannot get access secret")
	}

	dm, err := alicloud.NewDomainManager(
		alicloud.WithAccessKey(aKey),
		alicloud.WithAccessSecret(aSecret),
		alicloud.WithDomain("miniso.com"),
		alicloud.WithRegionID("cn-shenzhen"),
	)
	if err != nil {
		logger.Fatal("cannot create domain manager", zap.Error(err))
	}

	logger.Fatal("run grpc server failed", zap.Error(server.RunGRPCServer(&server.GRPCConfig{
		Name:   "domain",
		Addr:   ":18083",
		Logger: logger,
		RegisterFunc: func(s *grpc.Server) {
			domainpb.RegisterDomainServiceServer(s, &domain.Service{
				Logger:        logger,
				DomainManager: dm,
			})
		},
	})))
}
