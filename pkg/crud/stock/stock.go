package stock

import (
	"context"
	"fmt"

	_ "entgo.io/ent/dialect/sql" //nolint

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/stock-manager/pkg/db/ent"
	"github.com/NpoolPlatform/stock-manager/pkg/db/ent/stock"

	constant "github.com/NpoolPlatform/stock-manager/pkg/const"
	"github.com/NpoolPlatform/stock-manager/pkg/crud/entity"
	"github.com/NpoolPlatform/stock-manager/pkg/crud/tx"

	npool "github.com/NpoolPlatform/message/npool/stockmgr"

	"github.com/google/uuid"
)

type Stock struct {
	*entity.Entity
}

func New(ctx context.Context, _tx *ent.Tx) (*Stock, error) {
	e, err := entity.New(ctx, _tx)
	if err != nil {
		return nil, fmt.Errorf("fail create entity: %v", err)
	}

	return &Stock{
		Entity: e,
	}, nil
}

func (s *Stock) rowToObject(row *ent.Stock) *npool.Stock {
	return &npool.Stock{
		ID:        row.ID.String(),
		GoodID:    row.GoodID.String(),
		Total:     row.Total,
		InService: row.InService,
		Sold:      row.Sold,
	}
}

func (s *Stock) Create(ctx context.Context, in *npool.Stock) (*npool.Stock, error) {
	var info *ent.Stock
	var err error

	err = tx.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		info, err = s.Tx.Stock.Create().
			SetGoodID(uuid.MustParse(in.GetGoodID())).
			SetTotal(in.GetTotal()).
			SetInService(0).
			SetSold(0).
			Save(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail create stock: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *Stock) CreateBulk(ctx context.Context, in []*npool.Stock) ([]*npool.Stock, error) {
	rows := []*ent.Stock{}
	var err error

	err = tx.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		bulk := make([]*ent.StockCreate, len(in))
		for i, info := range in {
			bulk[i] = s.Tx.Stock.Create().
				SetGoodID(uuid.MustParse(info.GetGoodID())).
				SetTotal(info.GetTotal()).
				SetInService(0).
				SetSold(0)
		}
		rows, err = s.Tx.Stock.CreateBulk(bulk...).Save(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail create stocks: %v", err)
	}

	infos := []*npool.Stock{}
	for _, row := range rows {
		infos = append(infos, s.rowToObject(row))
	}

	return infos, nil
}

func (s *Stock) Update(ctx context.Context, in *npool.Stock) (*npool.Stock, error) {
	var info *ent.Stock
	var err error

	err = tx.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		info, err = s.Tx.Stock.UpdateOneID(uuid.MustParse(in.GetID())).
			SetInService(in.GetInService()).
			SetSold(in.GetSold()).
			Save(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail update stock: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *Stock) UpdateFields(ctx context.Context, id uuid.UUID, fields map[string]interface{}) (*npool.Stock, error) {
	var info *ent.Stock
	var err error

	err = tx.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		myTx := s.Tx.Stock.UpdateOneID(id)
		for k, v := range fields {
			switch k {
			case constant.StockFieldInService:
				myTx = myTx.SetInService(v.(uint32))
			case constant.StockFieldSold:
				myTx = myTx.SetSold(v.(uint32))
			default:
				return fmt.Errorf("invalid stock field")
			}
		}
		info, err = myTx.Save(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail update stock: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *Stock) AtomicUpdateFields(ctx context.Context, id uuid.UUID, fields map[string]interface{}) (*npool.Stock, error) {
	var info *ent.Stock

	err := tx.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		updater, err := s.Tx.Stock.Query().
			Where(stock.ID((id))).
			ForUpdate( /* sql.WithLockAction(sql.NoWait) */ ). //nolint
			Only(_ctx)
		if err != nil {
			return fmt.Errorf("fail lock stock: %v", err)
		}

		myTx := updater.Update()
		for k, v := range fields {
			switch k {
			case constant.StockFieldInService:
				myTx = myTx.SetInService(v.(uint32))
			case constant.StockFieldSold:
				myTx = myTx.SetSold(v.(uint32))
			default:
				return fmt.Errorf("invalid stock field")
			}
		}
		info, err = myTx.Save(_ctx)
		if err != nil {
			return fmt.Errorf("fail update stock: %v", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail atomic update stock: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *Stock) Row(ctx context.Context, id uuid.UUID) (*npool.Stock, error) {
	return nil, nil
}

func (s *Stock) Rows(ctx context.Context, conds map[string]*cruder.Cond, offset, limit uint32) ([]*npool.Stock, error) {
	return nil, nil
}

func (s *Stock) Count(ctx context.Context, conds map[string]*cruder.Cond) (uint32, error) {
	return 0, nil
}

func (s *Stock) Exist(ctx context.Context, id uuid.UUID) (bool, error) {
	return false, nil
}

func (s *Stock) ExistConds(ctx context.Context, conds map[string]*cruder.Cond) (bool, error) {
	return false, nil
}

func (s *Stock) Delete(ctx context.Context, id uuid.UUID) (*npool.Stock, error) {
	return nil, nil
}
