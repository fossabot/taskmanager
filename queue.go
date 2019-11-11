package taskmanager

import (
	"log"
	"sync"
)

// Тип для приоритетов
type Priority int

const (
	//Приоритеты задач
	LowestPriority  Priority = 1 // самый низкий
	LowPriority     Priority = 2 // низкий
	MiddlePriority  Priority = 3 // средний
	HighPriority    Priority = 4 // высокий
	HighestPriority Priority = 5 // самый высокий
)

// TaskInterface - интерфейс, который должна реализовывать задача
type TaskInterface interface {
	// Priority - метод возвращает приоритет задачи
	Priority() Priority
	// Exec - метод, где выполняется сама задача
	Exec()
}

// Интерфейс, который должна реализовывать очередь
type QueueInterface interface {
	// AddTask - добавить задачу из очереди
	AddTask(task TaskInterface)
	// GetTask - получить задачу из очереди
	GetTask() TaskInterface
}

// TaskQueue реализует очередь с прироритетом
type Queue struct {
	queue queue
	mu    sync.Mutex
}

// AddTask - добавление задач c блокировкой для безопасного добавления в асинхронных потоках
func (q *Queue) AddTask(task TaskInterface) {
	q.mu.Lock()
	q.queue.push(task)
	q.mu.Unlock()
}

// Получение задачи из канала очереди с блокировкой для безопасного извлечения в асинхронных потоках
func (q *Queue) GetTask() TaskInterface {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.queue) > 0 {
		return q.queue.pop()
	}
	return nil
}

// Слайс для хранения очереди задач
type queue []TaskInterface

// добавление нового элемента, просеиваем двоичную кучу вверх
func (q *queue) push(t TaskInterface) {
	*q = append(*q, t)
	q.up()
}

// извлечение элемента с прросеиванием вниз
func (q *queue) pop() TaskInterface {
	q.swap(0, len(*q)-1)
	q.down()

	old := *q
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // избегаем утечки памяти
	*q = old[0 : n-1]
	return item
}

func (q queue) less(i, j int) bool {
	return q[i].Priority() > q[j].Priority()
}

func (q queue) swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

// Восстановление пирамиды — операция up (реализация подсмотрена в container/heap)
func (q queue) up() {
	j := len(q) - 1 // последний вставленный элемент
	for {
		i := (j - 1) / 2 // родительский элемент
		if i == j || !q.less(j, i) {
			break
		}
		q.swap(i, j)
		j = i
	}
}

// Восстановление пирамиды - операция down (реализация подсмотрена в container/heap)
func (q queue) down() {
	n := len(q) - 1
	var i int
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && q.less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !q.less(j, i) {
			break
		}
		q.swap(i, j)
		i = j
	}
}

// Task - стандартная реализации задачи
type Task struct {
	EventDispatcherInterface // должна реализовывать интерфейс для работы с событиями
	priority                 Priority
	jobHandler               JobHandler
}

// JobHandler - тип для хэндлера с самой работой
type JobHandler func() error

// NewTask - Конструктор для новой задачи
func NewTask(priority Priority, handler JobHandler) *Task {
	task := &Task{
		EventDispatcherInterface: &EventDispatcher{},
		priority:                 priority,
		jobHandler:               handler,
	}
	return task
}

// Priority  - возвращает приоритет задачи
func (t *Task) Priority() Priority {
	return t.priority
}

// Exec - Выполняем задачу
// При старте задачи вызываем событие BeforeExecEvent
// Если задача выполнена неуспешно - вызываем событие FailedEvent
// Если нужна более гибкая обработка ошибок - реализуем свой Task
// По окончании задачи, если она прошла успешно вызываем событие AfterExecEvent
func (t *Task) Exec() {
	t.EmitEvent(BeforeExecEvent)
	if err := t.jobHandler(); err != nil {
		// логируем ошибку
		log.Println(err)
		t.EmitEvent(FailedEvent)
		// прерываем обработку задачи
		return
	}
	t.EmitEvent(AfterExecEvent)
}
