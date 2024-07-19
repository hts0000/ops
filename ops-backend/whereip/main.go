package main

import (
	"log"

	"github.com/hts0000/ops-backend/shared/server"
	whereippb "github.com/hts0000/ops-backend/whereip/api/gen/v1"
	"github.com/hts0000/ops-backend/whereip/whereip"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	logger.Fatal("run grpc server failed", zap.Error(server.RunGRPCServer(&server.GRPCConfig{
		Name:   "whereip",
		Addr:   ":18084",
		Logger: logger,
		RegisterFunc: func(s *grpc.Server) {
			whereippb.RegisterWhereipServiceServer(s, &whereip.Service{
				Logger: logger,
			})
		},
	})))
}
