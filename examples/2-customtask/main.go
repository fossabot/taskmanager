package main

import (
	"fmt"
	"time"

	"github.com/delgus/taskmanager"
)

// CustomTask пользовательский кастомный класс с высоким приоритетом
type CustomTask struct {
	// используем стандартную реализацию EventDispatcher для работы с событиями
	taskmanager.EventDispatcher
	isFailed bool
}

// NewCustomTask - конструктор для пользовательского класса
func NewCustomTask(isFailed bool) *CustomTask {
	task := &CustomTask{
		isFailed: isFailed,
	}
	return task
}

// Priority - реализует TaskInterface
func (t *CustomTask) Priority() taskmanager.Priority {
	return taskmanager.HighestPriority
}

// Exec - реализует TaskInterface
func (t *CustomTask) Exec() {
	t.EmitEvent(taskmanager.BeforeExecEvent)
	time.Sleep(time.Second * 1)
	if t.isFailed {
		t.EmitEvent(taskmanager.FailedEvent)
		return
	}
	t.EmitEvent(taskmanager.AfterExecEvent)
}

func main() {

	t1 := NewCustomTask(false)
	t1.OnEvent(taskmanager.BeforeExecEvent, func() {
		fmt.Println("executing custom task 1")
	})
	t1.OnEvent(taskmanager.FailedEvent, func() {
		fmt.Println(`oops task 1`)
	})
	t1.EmitEvent(taskmanager.CreatedEvent)

	t2 := NewCustomTask(true)
	t2.OnEvent(taskmanager.BeforeExecEvent, func() {
		fmt.Println("executing custom task 2")
	})
	t2.OnEvent(taskmanager.FailedEvent, func() {
		fmt.Println("oops task 2")
	})
	t2.EmitEvent(taskmanager.CreatedEvent)

	q := new(taskmanager.Queue)

	q.AddTask(t1)
	q.AddTask(t2)

	for {
		task := q.GetTask()
		if task == nil {
			break
		}
		task.Exec()
	}
}
