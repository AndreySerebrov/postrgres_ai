package scheduler

//go:generate mockgen -source=task.go -destination=task_mock.go -package scheduler

import "context"

// Наиболее абстрактная задача)
type Task interface {
	Do(ctx context.Context) error
}
