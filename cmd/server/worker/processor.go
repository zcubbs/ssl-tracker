package worker

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/zcubbs/logwrapper/logger"
	db "github.com/zcubbs/ssl-tracker/cmd/server/db/sqlc"
	"github.com/zcubbs/x/mail"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server     *asynq.Server
	store      db.Store
	mailer     mail.Mailer
	attributes Attributes
	logger     logger.Logger
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, mailer mail.Mailer, attributes Attributes, logger logger.Logger) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			logger.Error("task processing error",
				"task", task.Type,
				"error", err,
				"payload", string(task.Payload()),
			)
		}),
		Logger: NewLogger(logger),
	})
	return &RedisTaskProcessor{server: server, store: store, mailer: mailer, attributes: attributes, logger: logger}
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, p.ProcessTaskSendVerifyEmail)

	return p.server.Start(mux)
}
