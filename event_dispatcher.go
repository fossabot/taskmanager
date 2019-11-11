package taskmanager

// Event - тип для событий
type Event string

// EventHandler - тип для обработчика событий
type EventHandler func()

const (
	// События
	CreatedEvent    Event = "created"     // Создание
	BeforeExecEvent Event = "before_exec" // Начало выполнения
	AfterExecEvent  Event = "after_exec"  // Завершение выполнения
	FailedEvent     Event = "failed"      // ошибка выполнения
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
