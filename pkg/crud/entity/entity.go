package entity

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/stock-manager/pkg/db"
	"github.com/NpoolPlatform/stock-manager/pkg/db/ent"
)

type Entity struct {
	Tx *ent.Tx
}

func New(ctx context.Context, _tx *ent.Tx) (*Entity, error) {
	if _tx != nil {
		return &Entity{
			Tx: _tx,
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

	return &Entity{
		Tx: _tx,
	}, nil
}
