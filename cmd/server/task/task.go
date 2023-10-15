package task

import (
	db "github.com/zcubbs/tlz/cmd/server/db/sqlc"
	"github.com/zcubbs/tlz/cmd/server/logger"
)

var (
	log = logger.L()
)

type Task struct {
	store db.Store
}

func New(store db.Store) *Task {
	return &Task{
		store: store,
	}
}
