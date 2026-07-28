package helper

import (
	"backend-task-manager/internal/repository"
	"backend-task-manager/mixins"
	"context"

	"github.com/jackc/pgx/v5"
)

func WithTx[T any](ctx context.Context, base *repository.BaseRepo, fn func(ctx context.Context, exec pgx.Tx) (T, error)) (res T, err error) {
	tx, err := base.BeginTx(ctx)
	if err != nil {
		return res, err
	}

	defer func() {
		_ = mixins.TXReturn(tx, ctx, err)
	}()

	return fn(ctx, tx)
}

func WithTxVoid(ctx context.Context, base *repository.BaseRepo, fn func(ctx context.Context, exec pgx.Tx) error) error {
	tx, err := base.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		_ = mixins.TXReturn(tx, ctx, err)
	}()

	return fn(ctx, tx)
}
