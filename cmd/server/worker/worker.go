package worker

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zcubbs/tlz/cmd/server/config"
	db "github.com/zcubbs/tlz/cmd/server/db/sqlc"
	"github.com/zcubbs/tlz/cmd/server/logger"
	"github.com/zcubbs/x/mail"
)

type Worker struct {
	TaskDistributor TaskDistributor
	store           db.Store
	redisOpt        asynq.RedisClientOpt
	mailer          mail.Mailer
	attributes      Attributes
}

type Attributes struct {
	ApiDomainName string
}

var (
	log = logger.L()
)

func New(cfg config.Config, store db.Store, mailer mail.Mailer, attributes Attributes) *Worker {
	redisOpt := asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
	}
	return &Worker{
		TaskDistributor: NewRedisTaskDistributor(redisOpt),
		store:           store,
		redisOpt:        redisOpt,
		mailer:          mailer,
		attributes:      attributes,
	}
}

func (w *Worker) Run() {
	taskProcessor := NewRedisTaskProcessor(w.redisOpt, w.store, w.mailer, w.attributes, log)
	log.Info("✔️ starting task processor")

	err := taskProcessor.Start()
	if err != nil {
		log.Fatal("Cannot start task processor", "error", err)
	}
}
