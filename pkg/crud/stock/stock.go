package stock

import (
	"context"
	"fmt"
	"time"

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

func (s *Stock) AddFields(ctx context.Context, id uuid.UUID, fields map[string]interface{}) (*npool.Stock, error) {
	var info *ent.Stock
	var err error

	err = tx.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		myTx := s.Tx.Stock.UpdateOneID(id)
		for k, v := range fields {
			increment, err := cruder.AnyTypeInt32(v)
			if err != nil {
				return fmt.Errorf("invalid value type: %v", err)
			}

			switch k {
			case constant.StockFieldInService:
				myTx = myTx.AddInService(increment)
			case constant.StockFieldSold:
				myTx = myTx.AddSold(increment)
			default:
				return fmt.Errorf("invalid stock field")
			}
		}

		info, err = myTx.Save(_ctx)
		if err != nil {
			return fmt.Errorf("fail update stock: %v", err)
		}

		fmt.Println(info)
		fmt.Printf("\n\n\n\n\n\n")

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail add stock fields: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *Stock) SubFields(ctx context.Context, id uuid.UUID, fields map[string]interface{}) (*npool.Stock, error) {
	var info *ent.Stock
	var err error

	err = tx.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		myTx := s.Tx.Stock.UpdateOneID(id)
		for k, v := range fields {
			increment, err := cruder.AnyTypeInt32(v)
			if err != nil {
				return fmt.Errorf("invalid value type: %v", err)
			}
			increment *= -1

			switch k {
			case constant.StockFieldInService:
				myTx = myTx.AddInService(increment)
			case constant.StockFieldSold:
				myTx = myTx.AddSold(increment)
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
	var info *ent.Stock
	var err error

	err = tx.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		info, err = s.Tx.Stock.Query().Where(stock.ID(id)).Only(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail get stock: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *Stock) queryFromConds(conds map[string]*cruder.Cond) (*ent.StockQuery, error) {
	stm := s.Tx.Stock.Query()
	for k, v := range conds {
		switch k {
		case constant.FieldID:
			id, err := cruder.AnyTypeUUID(v.Val)
			if err != nil {
				return nil, fmt.Errorf("invalid id: %v", err)
			}
			stm = stm.Where(stock.ID(id))
		case constant.StockFieldGoodID:
			id, err := cruder.AnyTypeUUID(v.Val)
			if err != nil {
				return nil, fmt.Errorf("invalid good id: %v", err)
			}
			stm = stm.Where(stock.GoodID(id))
		case constant.StockFieldTotal:
			value, err := cruder.AnyTypeUint32(v.Val)
			if err != nil {
				return nil, fmt.Errorf("invalid total value: %v", err)
			}
			switch v.Op {
			case cruder.EQ:
				stm = stm.Where(stock.TotalEQ(value))
			case cruder.GT:
				stm = stm.Where(stock.TotalGT(value))
			case cruder.LT:
				stm = stm.Where(stock.TotalLT(value))
			}
		case constant.StockFieldInService:
			value, err := cruder.AnyTypeUint32(v.Val)
			if err != nil {
				return nil, fmt.Errorf("invalid value type: %v", err)
			}
			switch v.Op {
			case cruder.EQ:
				stm = stm.Where(stock.InServiceEQ(value))
			case cruder.GT:
				stm = stm.Where(stock.InServiceGT(value))
			case cruder.LT:
				stm = stm.Where(stock.InServiceLT(value))
			}
		case constant.StockFieldSold:
			value, err := cruder.AnyTypeUint32(v.Val)
			if err != nil {
				return nil, fmt.Errorf("invalid value type: %v", err)
			}
			switch v.Op {
			case cruder.EQ:
				stm = stm.Where(stock.SoldEQ(value))
			case cruder.GT:
				stm = stm.Where(stock.SoldGT(value))
			case cruder.LT:
				stm = stm.Where(stock.SoldLT(value))
			}
		default:
			return nil, fmt.Errorf("invalid stock field")
		}
	}

	return stm, nil
}

func (s *Stock) Rows(ctx context.Context, conds map[string]*cruder.Cond, offset, limit int) ([]*npool.Stock, int, error) {
	rows := []*ent.Stock{}
	var total int
	var err error

	err = tx.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		stm, err := s.queryFromConds(conds)
		if err != nil {
			return fmt.Errorf("fail construct stm: %v", err)
		}

		total, err = stm.Count(_ctx)
		if err != nil {
			return fmt.Errorf("fail count stock: %v", err)
		}

		rows, err = stm.Order(ent.Desc(stock.FieldUpdatedAt)).Offset(offset).Limit(limit).All(_ctx)
		if err != nil {
			return fmt.Errorf("fail query stock: %v", err)
		}

		return nil
	})
	if err != nil {
		return nil, 0, fmt.Errorf("fail get stock: %v", err)
	}

	infos := []*npool.Stock{}
	for _, row := range rows {
		infos = append(infos, s.rowToObject(row))
	}

	return infos, total, nil
}

func (s *Stock) Count(ctx context.Context, conds map[string]*cruder.Cond) (uint32, error) {
	var err error
	var total int

	err = tx.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		stm, err := s.queryFromConds(conds)
		if err != nil {
			return fmt.Errorf("fail construct stm: %v", err)
		}

		total, err = stm.Count(_ctx)
		if err != nil {
			return fmt.Errorf("fail check stocks: %v", err)
		}

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("fail count stocks: %v", err)
	}

	return uint32(total), nil
}

func (s *Stock) Exist(ctx context.Context, id uuid.UUID) (bool, error) {
	var err error
	exist := false

	err = tx.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		exist, err = s.Tx.Stock.Query().Where(stock.ID(id)).Exist(_ctx)
		return err
	})
	if err != nil {
		return false, fmt.Errorf("fail check stock: %v", err)
	}

	return exist, nil
}

func (s *Stock) ExistConds(ctx context.Context, conds map[string]*cruder.Cond) (bool, error) {
	var err error
	exist := false

	err = tx.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		stm, err := s.queryFromConds(conds)
		if err != nil {
			return fmt.Errorf("fail construct stm: %v", err)
		}

		exist, err = stm.Exist(_ctx)
		if err != nil {
			return fmt.Errorf("fail check stocks: %v", err)
		}

		return nil
	})
	if err != nil {
		return false, fmt.Errorf("fail check stocks: %v", err)
	}

	return exist, nil
}

func (s *Stock) Delete(ctx context.Context, id uuid.UUID) (*npool.Stock, error) {
	var info *ent.Stock
	var err error

	err = tx.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		info, err = s.Tx.Stock.UpdateOneID(id).
			SetDeletedAt(uint32(time.Now().Unix())).
			Save(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail delete stock: %v", err)
	}

	return s.rowToObject(info), nil
}
