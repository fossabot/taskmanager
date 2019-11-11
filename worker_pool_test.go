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

	go func() {
		// ждем пока пул воркеров получит все задачи и останавливаем
		time.Sleep(time.Second * 1)
		worker.Stop()
	}()

	worker.Run()

	if workCounter != int64(countTasks) {
		t.Error(`не все задачи выполнились`)
	}
}
