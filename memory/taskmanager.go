package memory

import (
	"sync"

	"github.com/delgus/taskmanager"
)

// Queue реализует очередь с прироритетом
type Queue struct {
	queue queue
	mu    sync.Mutex
}

// AddTask - добавление задач c блокировкой для безопасного добавления в асинхронных потоках
func (q *Queue) AddTask(task taskmanager.Task) {
	q.mu.Lock()
	q.queue.push(task)
	q.mu.Unlock()
}

// GetTask - получение задачи из канала очереди с блокировкой для безопасного извлечения в асинхронных потоках
func (q *Queue) GetTask() taskmanager.Task {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.queue) > 0 {
		return q.queue.pop()
	}
	return nil
}

// Слайс для хранения очереди задач
type queue []taskmanager.Task

// добавление нового элемента, просеиваем двоичную кучу вверх
func (q *queue) push(t taskmanager.Task) {
	*q = append(*q, t)
	q.up()
}

// извлечение элемента с прросеиванием вниз
func (q *queue) pop() taskmanager.Task {
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
	events   map[taskmanager.Event][]taskmanager.EventHandler
	priority taskmanager.Priority
	handler  TaskHandler
}

// TaskHandler - тип для хэндлера с самой работой
type TaskHandler func() error

// NewTask - Конструктор для новой задачи
func NewTask(priority taskmanager.Priority, handler TaskHandler) *Task {
	task := &Task{
		events:   make(map[taskmanager.Event][]taskmanager.EventHandler),
		priority: priority,
		handler:  handler,
	}
	return task
}

// Priority  - возвращает приоритет задачи
func (t *Task) Priority() taskmanager.Priority {
	return t.priority
}

// Exec - Выполняем задачу
// При старте задачи вызываем событие BeforeExecEvent
// Если задача выполнена неуспешно - вызываем событие FailedEvent
// Если нужна более гибкая обработка ошибок - реализуем свой Task
// По окончании задачи, если она прошла успешно вызываем событие AfterExecEvent
func (t *Task) Exec() error {
	t.EmitEvent(taskmanager.BeforeExecEvent)
	if err := t.handler(); err != nil {
		t.EmitEvent(taskmanager.FailedEvent)
		return err
	}
	t.EmitEvent(taskmanager.AfterExecEvent)
	return nil
}

// OnEvent - Функция, которая вешает обработчик на событие
func (t *Task) OnEvent(event taskmanager.Event, handler taskmanager.EventHandler) {
	t.events[event] = append(t.events[event], handler)
}

// EmitEvent вызывает событие, переданное в аргументах
func (t *Task) EmitEvent(event taskmanager.Event) {
	for _, h := range t.events[event] {
		h()
	}
}
