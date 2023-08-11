package task

import db "github.com/zcubbs/tlz/db/sqlc"

type Task struct {
	store db.Store
}

func New(store db.Store) *Task {
	return &Task{
		store: store,
	}
}
