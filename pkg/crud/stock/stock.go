package stock

import (
	"context"
	"fmt"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/stock-manager/pkg/db"
	"github.com/NpoolPlatform/stock-manager/pkg/db/ent"

	"github.com/NpoolPlatform/stock-manager/pkg/crud/tx"

	npool "github.com/NpoolPlatform/message/npool/stockmgr"

	"github.com/google/uuid"
)

type Stock struct {
	tx *ent.Tx
}

func New(ctx context.Context, _tx *ent.Tx) (*Stock, error) {
	if _tx != nil {
		return &Stock{
			tx: _tx,
		}, nil
	}

	cli, err := db.Client()
	if err != nil {
		return nil, fmt.Errorf("fail get db client: %v", err)
	}
	_tx, err = cli.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail get client transaction: %v", err)
	}

	return &Stock{
		tx: _tx,
	}, nil
}

func (s *Stock) rowToObject(row *ent.Stock) *npool.Stock {
	return &npool.Stock{
		ID:        row.ID.String(),
		GoodID:    row.GoodID.String(),
		InService: row.InService,
		Sold:      row.Sold,
	}
}

func (s *Stock) Create(ctx context.Context, in *npool.Stock) (*npool.Stock, error) {
	var info *ent.Stock
	var err error

	err = tx.WithTx(ctx, s.tx, func() error {
		info, err = s.tx.Stock.Create().
			SetGoodID(uuid.MustParse(in.GetGoodID())).
			SetInService(in.GetInService()).
			SetSold(in.GetSold()).
			Save(ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail create stock: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *Stock) CreateBulk(ctx context.Context, in []*npool.Stock) ([]*npool.Stock, error) {
	return nil, nil
}

func (s *Stock) Update(ctx context.Context, in *npool.Stock) (*npool.Stock, error) {
	return nil, nil
}

func (s *Stock) UpdateFields(ctx context.Context, id string, fields map[string]*npool.Stock) (*npool.Stock, error) {
	return nil, nil
}

func (s *Stock) AtomicInc(ctx context.Context, id string, fields []string) (*npool.Stock, error) {
	return nil, nil
}

func (s *Stock) AtomicSub(ctx context.Context, id string, fields []string) (*npool.Stock, error) {
	return nil, nil
}

func (s *Stock) AtomicSet(ctx context.Context, id string, fields map[string]*npool.Stock) (*npool.Stock, error) {
	return nil, nil
}

func (s *Stock) Row(ctx context.Context, id string) (*npool.Stock, error) {
	return nil, nil
}

func (s *Stock) Rows(ctx context.Context, conds map[string]cruder.Cond, offset, limit uint32) ([]*npool.Stock, error) {
	return nil, nil
}

func (s *Stock) Count(ctx context.Context, conds map[string]cruder.Cond) (uint32, error) {
	return 0, nil
}

func (s *Stock) Exist(ctx context.Context, conds map[string]cruder.Cond) (bool, error) {
	return false, nil
}

func (s *Stock) Delete(ctx context.Context, id string) (*npool.Stock, error) {
	return nil, nil
}
