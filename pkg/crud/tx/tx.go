package tx

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/stock-manager/pkg/db/ent"
)

func WithTx(ctx context.Context, tx *ent.Tx, fn func() error) error {
	defer func() {
		if v := recover(); v != nil {
			err := tx.Rollback()
			if err != nil {
				logger.Sugar().Errorf("fail to rollback: %v", err)
			}
			panic(v)
		}
	}()
	if err := fn(); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("rolling back transaction: %v (%v)", err, rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %v", err)
	}
	return nil
}
