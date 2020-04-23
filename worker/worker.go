package worker

import (
	"fmt"
	"sync"
	"time"

	"github.com/delgus/taskmanager"
)

type Logger interface {
	Error(interface{})
}

// Pool - пул воркеров для обработки задач
type Pool struct {
	logger            Logger
	queue             taskmanager.Queue
	wg                sync.WaitGroup
	maxWorkers        int                   // количество воркеров
	periodicityTicker *time.Ticker          // частота с которой воркер пул проверяет есть ли задачи в очереди
	closeTaskCh       chan struct{}         // канал для остановки пула воркеров
	taskCh            chan taskmanager.Task // канал с поступающими задачами
	quit              chan struct{}         // канал, после получения сигнала прекращает работу
}

// NewPool - конструктор для воркера задач
// maxWorkers - количество воркеров в пуле
// periodicity - частота с которой пул воркеров проверяет есть ли задачи в очереди
func NewPool(queue taskmanager.Queue, maxWorkers int, periodicity time.Duration) *Pool {
	return &Pool{
		queue:             queue,
		maxWorkers:        maxWorkers,
		periodicityTicker: time.NewTicker(periodicity),
		closeTaskCh:       make(chan struct{}),
		taskCh:            make(chan taskmanager.Task),
		quit:              make(chan struct{}),
	}
}

// Run - запуск воркера для работы
func (w *Pool) Run() {
	// заполняем канал задачами с определенной периодичностью
	go func() {
		for {
			select {
			case <-w.closeTaskCh:
				close(w.taskCh)
				return
			case <-w.periodicityTicker.C:
				if task := w.queue.GetTask(); task != nil {
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

func (w *Pool) work() {
	for task := range w.taskCh {
		if err := task.Exec(); err != nil {
			w.logger.Error(err)
		}
	}
	w.wg.Done()
}

// Shutdown - плавная остановка воркера
// воркер не остановится пока не выполнит все недоработанные задачи
// или не истечет тайм аут
func (w *Pool) Shutdown(timeout time.Duration) error {
	w.closeTaskCh <- struct{}{}

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
		return fmt.Errorf(`taskmanager: timeout error`)
	}
}
