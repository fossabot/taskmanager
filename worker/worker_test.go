package worker

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/delgus/taskmanager"
	"github.com/delgus/taskmanager/memory"
)

// воркер должен получить все задачи из очереди и вызвать у них обработчик
func TestWorkerPool(t *testing.T) {
	// новая очередь задач
	q := new(memory.Queue)

	var workCounter int64

	var countTasks = 5 // количество задач

	testTask := memory.NewTask(taskmanager.HighestPriority, func() error {
		// добавляем атомарно в счетчик выполненую работу
		// чтобы избежать data race condition
		atomic.AddInt64(&workCounter, 1)
		time.Sleep(time.Second * 2)
		return nil
	})

	for i := 0; i < countTasks; i++ {
		q.AddTask(testTask)
	}

	worker := NewPool(q, 10, time.Millisecond)

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
	q := new(memory.Queue)

	testTask := memory.NewTask(taskmanager.HighestPriority, func() error {
		time.Sleep(time.Second * 10)
		return nil
	})
	q.AddTask(testTask)
	workerPool := NewPool(q, 2, time.Millisecond)
	go workerPool.Run()
	time.Sleep(time.Second)
	if err := workerPool.Shutdown(time.Second); err == nil {
		t.Error(`expected timeout error`)
	}
}
