// +build !codeanalysis

package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/stockmgr"
	crud "github.com/NpoolPlatform/stock-manager/pkg/crud/stock"

	"github.com/google/uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateStock(ctx context.Context, in *npool.CreateStockRequest) (*npool.CreateStockResponse, error) {
	if _, err := uuid.Parse(in.GetInfo().GetGoodID()); err != nil {
		logger.Sugar().Errorf("invalid request good id: %v", err)
		return &npool.CreateStockResponse{}, status.Error(codes.Internal, err.Error())
	}

	schema, err := crud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return &npool.CreateStockResponse{}, status.Error(codes.Internal, err.Error())
	}

	info, err := schema.Create(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorf("fail create stock: %v", err)
		return &npool.CreateStockResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateStockResponse{
		Info: info,
	}, nil
}

func (s *Server) CreateStocks(ctx context.Context, in *npool.CreateStocksRequest) (*npool.CreateStocksResponse, error) {
	for _, info := range in.GetInfos() {
		if _, err := uuid.Parse(info.GetGoodID()); err != nil {
			logger.Sugar().Errorf("invalid request good id: %v", err)
			return &npool.CreateStocksResponse{}, status.Error(codes.Internal, err.Error())
		}
	}

	schema, err := crud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return &npool.CreateStocksResponse{}, status.Error(codes.Internal, err.Error())
	}

	infos, err := schema.CreateBulk(ctx, in.GetInfos())
	if err != nil {
		logger.Sugar().Errorf("fail create stocks: %v", err)
		return &npool.CreateStocksResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateStocksResponse{
		Infos: infos,
	}, nil
}

func (s *Server) UpdateStock(ctx context.Context, in *npool.UpdateStockRequest) (*npool.UpdateStockResponse, error) {
	if _, err := uuid.Parse(in.GetInfo().GetGoodID()); err != nil {
		logger.Sugar().Errorf("invalid request good id: %v", err)
		return &npool.UpdateStockResponse{}, status.Error(codes.Internal, err.Error())
	}

	if _, err := uuid.Parse(in.GetInfo().GetID()); err != nil {
		logger.Sugar().Errorf("invalid stock id: %v", err)
		return &npool.UpdateStockResponse{}, status.Error(codes.Internal, err.Error())
	}

	schema, err := crud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail update schema entity: %v", err)
		return &npool.UpdateStockResponse{}, status.Error(codes.Internal, err.Error())
	}

	info, err := schema.Update(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorf("fail update stock: %v", err)
		return &npool.UpdateStockResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateStockResponse{
		Info: info,
	}, nil
}
