package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
}

type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

func NewSQLStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
