package taskmanager

import (
	"sync/atomic"
	"time"
)

// Воркер для обработки задач
type WorkerPool struct {
	tm              *Queue
	countGoroutines int           // количество горутин в которых обрабатываются задачи (количество одновременно обрабатываемых задач)
	periodicity     time.Duration // частота с которой воркер проверяет есть ли задачи в очереди
	quitCh          chan struct{} // канал для остановки пула воркеров
}

// Конструктор для воркера задач
func NewWorkerPool(tm *Queue, countGoroutines int, periodicity time.Duration) *WorkerPool {
	return &WorkerPool{
		tm:              tm,
		countGoroutines: countGoroutines,
		periodicity:     periodicity,
		quitCh:          make(chan struct{}),
	}
}

// Запуск воркера для работы
func (w *WorkerPool) Run() {
	// канал для отправки задач горутинам
	// канал небуферизирован чтоб получать только актуальную задачу
	tasker := make(chan TaskInterface)
	go func() {
		for range time.Tick(w.periodicity) {
			if task := w.tm.GetTask(); task != nil {
				tasker <- task
			}
		}
	}()

	// ограничиваем количество горутин чтоб не положить проц при большом количестве задач
	limiter := make(chan struct{}, w.countGoroutines)

	// количество задач находящихся в работе
	var countTasks int64
loop:
	for {
		select {
		case task := <-tasker:
			limiter <- struct{}{}
			atomic.AddInt64(&countTasks, 1)
			go func() {
				task.Exec()
				atomic.AddInt64(&countTasks, -1)
				<-limiter
			}()
		case <-w.quitCh:
			w.periodicity = time.Hour * 24 // увеличиваем частоту чтоб больше не доставать задачи из очереди
			for {
				//ждем когда выполнятся все задачи
				if atomic.LoadInt64(&countTasks) == 0 {
					// завершаем работу всего цикла
					break loop
				}
			}
		}
	}
}

// плавная остановка воркера
// воркер не остановится пока не выполнит все недоработанные задачи
func (w *WorkerPool) Stop() {
	w.quitCh <- struct{}{}
}
