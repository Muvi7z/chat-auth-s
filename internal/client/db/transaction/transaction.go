package transaction

import (
	"context"
	"github.com/Muvi7z/chat-auth-s/internal/client/db"
	"github.com/Muvi7z/chat-auth-s/internal/client/db/pg"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type manager struct {
	db db.Transactor
}

func NewTransactionManager(db db.Transactor) db.TxManager {
	return &manager{db: db}
}

func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, fn db.Handler) (err error) {
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		return fn(ctx)
	}

	tx, err = m.db.BeginTx(ctx, opts)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	ctx = pg.MakeContextTx(ctx, tx)

	defer func() {
		if p := recover(); p != nil {
			err = errors.Errorf("panic: %v", p)
		}

		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = errors.Wrap(errRollback, "failed to rollback transaction")
			}

			return
		}

		if err == nil {
			err = tx.Commit(ctx)
			if err != nil {
				err = errors.Wrap(err, "failed to commit transaction")
			}
		}
	}()

	if err = fn(ctx); err != nil {
		err = errors.Wrapf(err, "failed to executing code inside transaction")
	}

	return err
}

func (m *manager) ReadCommited(ctx context.Context, f db.Handler) error {
	txOpt := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, txOpt, f)
}
