package taskmanager

import (
	"sync/atomic"
	"testing"
	"time"
)

// воркер должен получить все задачи из очереди и вызвать у них обработчик
func TestWorkerPool(t *testing.T) {
	// новая очередь задач
	q := new(Queue)

	var workCounter int64

	var countTasks = 5 //количество задач

	testTask := NewTask(HighestPriority, func() error {
		// добавляем атомарно в счетчик выполненую работу
		// чтобы избежать data race condition
		atomic.AddInt64(&workCounter, 1)
		time.Sleep(time.Second * 2)
		return nil
	})

	for i := 0; i < countTasks; i++ {
		q.AddTask(testTask)
	}

	worker := NewWorkerPool(q, 10, time.Millisecond)

	go worker.Run()

	// ждем пока пул воркеров получит все задачи и останавливаем
	time.Sleep(time.Second * 1)
	// таймаут устанавливаем 3 секунды, что больше чем любая из задач
	if err := worker.Shutdown(time.Second * 3); err != nil {
		t.Error(err)
	}

	if workCounter != int64(countTasks) {
		t.Error(`не все задачи выполнились`)
	}
}

func TestWorkerPool_Shutdown(t *testing.T) {
	// новая очередь задач
	q := new(Queue)

	testTask := NewTask(HighestPriority, func() error {
		// добавляем атомарно в счетчик выполненую работу
		// чтобы избежать data race condition
		time.Sleep(time.Second * 10)
		return nil
	})
	q.AddTask(testTask)
	workerPool := NewWorkerPool(q, 2, time.Millisecond)
	go workerPool.Run()
	time.Sleep(time.Second)
	if err := workerPool.Shutdown(time.Second); err == nil {
		t.Error(`expected timeout error`)
	}
}
