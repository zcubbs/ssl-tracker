package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	db "github.com/zcubbs/ssl-tracker/cmd/server/db/sqlc"
	"github.com/zcubbs/x/mail"
	"github.com/zcubbs/x/random"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

func (d *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	task := asynq.NewTask(
		TaskSendVerifyEmail,
		jsonPayload,
		opts...,
	)
	info, err := d.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info("enqueued task",
		"task", task.Type(),
		"payload", string(task.Payload()),
		"queue", info.Queue,
		"max_retry", info.MaxRetry,
	)

	return nil
}

func (p *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	user, err := p.store.GetUser(ctx, payload.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user not found: %w", asynq.SkipRetry)
		}

		return fmt.Errorf("failed to get user: %w", err)
	}

	verifyEmail, err := p.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		UserID:     user.ID,
		Email:      user.Email,
		SecretCode: random.String(32),
	})
	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	log.Info("processed task",
		"task", task.Type(),
		"payload", payload,
		"email", user.Email,
	)

	verifyUrl := fmt.Sprintf("%s/api/v1/verify_email?email_id=%s&secret_code=%s",
		p.attributes.ApiDomainName, verifyEmail.ID.String(), verifyEmail.SecretCode)

	content := fmt.Sprintf(`Hello %s, <br>
	Welcome to TLZ! <br>
	please verify your email address by clicking
	this <a href="%s">link</a>`,
		user.Username,
		verifyUrl,
	)

	if err := p.mailer.SendMail(mail.Mail{
		To:            []string{user.Email},
		Cc:            nil,
		Bcc:           nil,
		Subject:       "Welcome to TLZ",
		Content:       content,
		AttachedFiles: nil,
	}); err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	return nil
}
