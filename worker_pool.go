package taskmanager

import (
	"fmt"
	"sync"
	"time"
)

// Воркер для обработки задач
type WorkerPool struct {
	tm          *Queue
	wg          sync.WaitGroup
	maxWorkers  int                // количество воркеров
	periodicity time.Duration      // частота с которой воркер пул проверяет есть ли задачи в очереди
	closeTaskCh chan struct{}      // канал для остановки пула воркеров
	taskCh      chan TaskInterface // канал с поступающими задачами
	quit        chan struct{}      // канал, после получения сигнала прекращает работу
}

// Конструктор для воркера задач
func NewWorkerPool(tm *Queue, maxWorkers int, periodicity time.Duration) *WorkerPool {
	return &WorkerPool{
		tm:          tm,
		maxWorkers:  maxWorkers,
		periodicity: periodicity,
		closeTaskCh: make(chan struct{}),
		taskCh:      make(chan TaskInterface),
		quit:        make(chan struct{}),
	}
}

// Запуск воркера для работы
func (w *WorkerPool) Run() {
	// заполняем канал задачами с определенной периодичностью, дабы не положить проц
	go func() {
		ticker := time.NewTicker(w.periodicity)
		for {
			select {
			case <-w.closeTaskCh:
				close(w.taskCh)
				return
			case <-ticker.C:
				w.taskCh <- w.tm.GetTask()
			}
		}
	}()

	// запускаем пул воркеров
	w.wg.Add(w.maxWorkers)
	for i := 0; i < w.maxWorkers; i++ {
		go w.work()
	}
	<-w.quit
}

func (w *WorkerPool) work() {
	for task := range w.taskCh {
		if task != nil {
			task.Exec()
		}
	}
	w.wg.Done()
}

// плавная остановка воркера
// воркер не остановится пока не выполнит все недоработанные задачи
// или не истечет тайм аут
func (w *WorkerPool) Shutdown(timeout time.Duration) error {
	// закрываем канал с задачами
	w.closeTaskCh <- struct{}{}
	// если воркеры закончили работу и остановились отправляем сообщени в канал ok
	ok := make(chan struct{})
	go func() {
		w.wg.Wait()
		ok <- struct{}{}
	}()

	select {
	case <-ok:
		close(w.quit)
		return nil
	case <-time.After(timeout):
		close(w.quit)
		return fmt.Errorf(`timeout`)
	}
}
