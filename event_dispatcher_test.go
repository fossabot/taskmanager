package taskmanager

import "testing"

func TestEventDispatcher(t *testing.T) {
	ed := new(EventDispatcher)

	eventFlag := false
	ed.OnEvent(CreatedEvent, func() {
		eventFlag = true
	})

	ed.EmitEvent(BeforeExecEvent)

	if eventFlag {
		t.Errorf(`unexpected execution of handler`)
	}

	ed.EmitEvent(CreatedEvent)

	if !eventFlag {
		t.Errorf(`handler not execute`)
	}
}
