package taskmanager

type executor interface {
	Exec() error
}

// Priority - тип для приоритетов
type Priority int

const (
	// LowestPriority - самый низкий
	LowestPriority Priority = 1
	// LowPriority - низкий
	LowPriority Priority = 2
	// MiddlePriority - средний
	MiddlePriority Priority = 3
	// HighPriority - высокий
	HighPriority Priority = 4
	// HighestPriority - самый высокий
	HighestPriority Priority = 5
)

type prioritier interface {
	Priority() Priority
}

// Event - тип для событий
type Event string

// EventHandler - тип для обработчика событий
type EventHandler func()

const (
	// CreatedEvent - создание
	CreatedEvent Event = "created"
	// BeforeExecEvent - начало выполнения
	BeforeExecEvent Event = "before_exec"
	// AfterExecEvent - завершение выполнения
	AfterExecEvent Event = "after_exec"
	// FailedEvent - ошибка выполнения
	FailedEvent Event = "failed"
)

type eventer interface {
	OnEvent(event Event, handler EventHandler)
	EmitEvent(event Event)
}

// Task - интерфейс, который должна реализовывать задача
type Task interface {
	executor
	prioritier
	eventer
}

// Queue - интерфейс, который должна реализовывать очередь
type Queue interface {
	// AddTask - добавить задачу из очереди
	AddTask(task Task)
	// GetTask - получить задачу из очереди
	GetTask() Task
}
