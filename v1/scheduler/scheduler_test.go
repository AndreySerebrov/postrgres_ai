package scheduler

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func Test_10GoodTasks_3Threads(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	taskList := []Task{}
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		task := NewMockTask(mockCtrl)
		task.EXPECT().Do(gomock.Any()).Return(nil)

		taskList = append(taskList, task)
	}

	schldr := NewSchedular(3, 1)
	schldr.Execute(ctx, taskList)
}

func Test_10GoodTasks_30Threads(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	taskList := []Task{}
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		task := NewMockTask(mockCtrl)
		task.EXPECT().Do(gomock.Any()).Return(nil)

		taskList = append(taskList, task)
	}

	schldr := NewSchedular(30, 1)
	schldr.Execute(ctx, taskList)
}

func Test_3BadTasks_2GoodTask_2Threads(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	taskList := []Task{}
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		task := NewMockTask(mockCtrl)
		task.EXPECT().Do(gomock.Any()).Return(fmt.Errorf("some error"))
		taskList = append(taskList, task)
	}

	for i := 0; i < 2; i++ {
		task := NewMockTask(mockCtrl)
		task.EXPECT().Do(gomock.Any()).MinTimes(0).MaxTimes(1).Return(nil)
		taskList = append(taskList, task)
	}

	schldr := NewSchedular(2, 3)
	schldr.Execute(ctx, taskList)
}

func Test_3BadTasks_2GoodTask_2Threads_DifferentDuration(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	taskList := []Task{}
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		task := NewMockTask(mockCtrl)
		task.EXPECT().Do(gomock.Any()).DoAndReturn(
			func(ctx context.Context) error {
				select {
				case <-ctx.Done():
				case <-time.After(time.Millisecond * 100 * time.Duration(i)):
				}
				return fmt.Errorf("some error")
			},
		)

		taskList = append(taskList, task)
	}

	for i := 0; i < 2; i++ {
		task := NewMockTask(mockCtrl)
		task.EXPECT().Do(gomock.Any()).MinTimes(0).MaxTimes(1).DoAndReturn(
			func(ctx context.Context) error {
				select {
				case <-ctx.Done():
				case <-time.After(time.Millisecond * 100 * time.Duration(i)):
				}
				return nil
			},
		)

		taskList = append(taskList, task)
	}

	schldr := NewSchedular(2, 3)
	schldr.Execute(ctx, taskList)
}
