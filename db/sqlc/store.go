package db

type Store interface {
	Querier
}

type SQLStore struct {
	*Queries
}

func NewSQLStore(db DBTX) Store {
	return &SQLStore{
		Queries: New(db),
	}
}
