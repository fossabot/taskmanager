package taskmanager

import (
	"fmt"
	"sync"
	"time"
)

// WorkerPool - пул воркеров для обработки задач
type WorkerPool struct {
	tm                *Queue
	wg                sync.WaitGroup
	maxWorkers        int                // количество воркеров
	periodicityTicker *time.Ticker       // частота с которой воркер пул проверяет есть ли задачи в очереди
	closeTaskCh       chan struct{}      // канал для остановки пула воркеров
	taskCh            chan TaskInterface // канал с поступающими задачами
	quit              chan struct{}      // канал, после получения сигнала прекращает работу
}

// NewWorkerPool - конструктор для воркера задач
// maxWorkers - количество воркеров в пуле
// periodicity - частота с которой пул воркеров проверяет есть ли задачи в очереди
func NewWorkerPool(tm *Queue, maxWorkers int, periodicity time.Duration) *WorkerPool {
	return &WorkerPool{
		tm:                tm,
		maxWorkers:        maxWorkers,
		periodicityTicker: time.NewTicker(periodicity),
		closeTaskCh:       make(chan struct{}),
		taskCh:            make(chan TaskInterface),
		quit:              make(chan struct{}),
	}
}

// Запуск воркера для работы
func (w *WorkerPool) Run() {
	// заполняем канал задачами с определенной периодичностью, дабы не положить проц
	go func() {
		for {
			select {
			case <-w.closeTaskCh:
				close(w.taskCh)
				return
			case <-w.periodicityTicker.C:
				if task := w.tm.GetTask(); task != nil {
					w.taskCh <- task
				}
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
		task.Exec()
	}
	w.wg.Done()
}

// Shutdown - плавная остановка воркера
// воркер не остановится пока не выполнит все недоработанные задачи
// или не истечет тайм аут
func (w *WorkerPool) Shutdown(timeout time.Duration) error {
	// закрываем канал с задачами
	w.closeTaskCh <- struct{}{}
	// если воркеры закончили работу и остановились отправляем сообщение в канал ok
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
