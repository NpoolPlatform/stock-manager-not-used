package api

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/stockmgr"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	npool.UnimplementedStockManagerServer
}

func Register(server grpc.ServiceRegistrar) {
	npool.RegisterStockManagerServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return npool.RegisterStockManagerHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
