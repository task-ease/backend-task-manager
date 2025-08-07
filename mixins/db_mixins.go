package mixins

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func TXReturn(tx pgx.Tx, ctx context.Context, err error) error {
	if err != nil {
		rbErr := tx.Rollback(ctx)
		if rbErr != nil {
		}
		return err
	} else {
		return tx.Commit(ctx)
	}
}
