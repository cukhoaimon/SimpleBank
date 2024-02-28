package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	TaskSendVerifyEmail = "task:send_verify_email"
)

type PayloadVerifyEmail struct {
	Username string `json:"username"`
}

func (r RedisTaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadVerifyEmail, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("fail to marshal task payload: %v", payload)
	}
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)

	taskInfo, err := r.client.EnqueueContext(ctx, task)
	if err != nil {
		return err
	}

	log.Info().
		Str("task", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", taskInfo.Queue).
		Int("max_retry", taskInfo.MaxRetry).
		Msg("enqueue task")

	return nil
}

func (r RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadVerifyEmail

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("fail to unmarshal payload: %w", asynq.SkipRetry)
	}
	user, err := r.store.GetUser(ctx, payload.Username)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return fmt.Errorf("user does not exists: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("fail to get user: %w", asynq.SkipRetry)
	}

	//TODO send real email
	log.Info().
		Str("task", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg("processed task")
	return nil
}
