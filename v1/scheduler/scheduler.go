package scheduler

import (
	"context"
	"log"
	"sync"
)

type scheduler struct {
	n int64 // max go routine number
	m int64 // max error number
}

func NewSchedular(goRoutineNum int64, maxErrNum int64) scheduler {
	return scheduler{goRoutineNum, maxErrNum}
}

func (s *scheduler) Execute(ctx context.Context, taskList []Task) {
	ctx, cancelFun := context.WithCancel(ctx)
	taskChan := make(chan Task)
	errChan := make(chan error)
	wg := sync.WaitGroup{}

	s.workerStarterStart(ctx, &wg, taskChan, errChan)
	s.taskSetterStart(ctx, taskList, taskChan)
	s.errorHandlerStart(cancelFun, errChan, s.m)

	wg.Wait()
	close(errChan)
}

// workerStarterStart
// Запускает n воркеров, которые получают задачи через канал taskChan
// Ошибки пишутся в канал errChan
func (s *scheduler) workerStarterStart(ctx context.Context, wg *sync.WaitGroup, taskChan chan Task, errChan chan error) {

	wg.Add(int(s.n))
	for i := 0; i < int(s.n); i++ {
		go func() {
			worker := NewWorker()
			worker.Start(ctx, taskChan, errChan)
			wg.Done()
		}()
	}
}

// taskSetterStart
// Публикует задачи в канал taskList
// Прерывает работу при завершении контекста
func (s *scheduler) taskSetterStart(ctx context.Context, taskList []Task, taskChan chan Task) {

	go func() {
		defer close(taskChan)
		for _, task := range taskList {
			select {
			case <-ctx.Done():
				return
			case taskChan <- task:
			}
		}
	}()
}

// errorHandlerStart
// Обрабатывает поступающией ошибки
func (s *scheduler) errorHandlerStart(cancelFun context.CancelFunc, errChan chan error, m int64) {
	go func() {
		var errNum int64

		for err := range errChan {
			errNum++
			if errNum >= m {
				log.Println("Max error level!", err)
				cancelFun()
			}
		}
	}()
}
