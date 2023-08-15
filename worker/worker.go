package worker

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/hibiken/asynq"
	db "github.com/zcubbs/tlz/db/sqlc"
	"github.com/zcubbs/tlz/internal/util"
)

type Worker struct {
	TaskDistributor TaskDistributor
	store           db.Store
	redisOpt        asynq.RedisClientOpt
}

func New(cfg util.Config, store db.Store) *Worker {
	redisOpt := asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
	}
	return &Worker{
		TaskDistributor: NewRedisTaskDistributor(redisOpt),
		store:           store,
		redisOpt:        redisOpt,
	}
}

func (w *Worker) Run() {
	taskProcessor := NewRedisTaskProcessor(w.redisOpt, w.store)
	log.Info("✔️ starting task processor")

	err := taskProcessor.Start()
	if err != nil {
		log.Fatal("Cannot start task processor", "error", err)
	}
}
