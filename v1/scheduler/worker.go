package scheduler

import (
	"context"
)

type Worker struct {
}

func NewWorker() *Worker {
	return &Worker{}
}

// Start
// Запускает цикл обработки задач, поступающих по каналу taskChan
func (w *Worker) Start(ctx context.Context, taskChan chan Task, errChan chan error) {

	for task := range taskChan {
		err := task.Do(ctx)
		if err != nil && err != context.Canceled {
			errChan <- err
		}
	}
}
