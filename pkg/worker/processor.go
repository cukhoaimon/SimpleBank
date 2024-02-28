package worker

import (
	"context"
	db "github.com/cukhoaimon/SimpleBank/internal/usecase/sqlc"
	"github.com/hibiken/asynq"
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
	server *asynq.Server
	store  db.Store
}

func (r RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, r.ProcessTaskSendVerifyEmail)

	return r.server.Start(mux)
}

func NewRedisTaskProcessor(redisOpts asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(
		redisOpts,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
		},
	)
	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}
