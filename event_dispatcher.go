package taskmanager

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

// EventDispatcherInterface - интерфейс, который должен реализовывать EventDispatcher
// При реализации необходимо придерживаться First-in first-out
// Первый обработчик, повешенный на событие должен сработать первым
type EventDispatcherInterface interface {
	// OnEvent - вешаем обработчик на событие
	OnEvent(event Event, handler EventHandler)
	// EmitEvent - вызвать событие
	EmitEvent(event Event)
}

// EventDispatcher -реализация EventDispatcherInterface
type EventDispatcher struct {
	events map[Event][]EventHandler
}

// OnEvent - Функция, которая вешает обработчик на событие
func (d *EventDispatcher) OnEvent(event Event, handler EventHandler) {
	if d.events == nil {
		d.events = make(map[Event][]EventHandler)
	}
	d.events[event] = append(d.events[event], handler)
}

// EmitEvent вызывает событие, переданное в аргументах
func (d *EventDispatcher) EmitEvent(event Event) {
	for _, h := range d.events[event] {
		h()
	}
}
