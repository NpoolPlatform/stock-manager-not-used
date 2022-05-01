package client

import (
	"context"
	"fmt"
	"time"

	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/stockmgr"

	servicename "github.com/NpoolPlatform/stock-manager/pkg/service-name"
)

func do(ctx context.Context, fn func(_ctx context.Context, cli npool.StockManagerClient) (cruder.Any, error)) (cruder.Any, error) {
	_ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	conn, err := grpc2.GetGRPCConn(servicename.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get stock connection: %v", err)
	}
	defer conn.Close()

	cli := npool.NewStockManagerClient(conn)

	return fn(_ctx, cli)
}

func GetStocks(ctx context.Context, conds cruder.FilterConds) ([]*npool.Stock, error) {
	infos, err := do(ctx, func(_ctx context.Context, cli npool.StockManagerClient) (cruder.Any, error) {
		resp, err := cli.GetStocks(ctx, &npool.GetStocksRequest{
			Conds: conds,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get stocks: %v", err)
		}
		return resp.Infos, nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail get stocks: %v", err)
	}
	return infos.([]*npool.Stock), nil
}

func GetStockOnly(ctx context.Context, conds cruder.FilterConds) (*npool.Stock, error) {
	info, err := do(ctx, func(_ctx context.Context, cli npool.StockManagerClient) (cruder.Any, error) {
		resp, err := cli.GetStockOnly(ctx, &npool.GetStockOnlyRequest{
			Conds: conds,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get stock: %v", err)
		}
		return resp.Info, nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail get stock: %v", err)
	}
	return info.(*npool.Stock), nil
}

func AddFields(ctx context.Context, id string, fields cruder.FilterFields) (*npool.Stock, error) {
	info, err := do(ctx, func(_ctx context.Context, cli npool.StockManagerClient) (cruder.Any, error) {
		resp, err := cli.AddStockFields(ctx, &npool.AddStockFieldsRequest{
			ID:     id,
			Fields: fields,
		})
		if err != nil {
			return nil, fmt.Errorf("fail add stock fields: %v", err)
		}
		return resp.Info, nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail add stock fields: %v", err)
	}
	return info.(*npool.Stock), nil
}
