package taskmanager

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// Производительность при добавлении задач
func BenchmarkAddTask(b *testing.B) {
	tm := new(Queue)
	testTask := NewTask(HighestPriority, func() error {
		return nil
	})
	for i := 0; i < b.N; i++ {
		tm.AddTask(testTask)
	}
}

// Производительность при добавлении и чтении задач
func BenchmarkGetTask(b *testing.B) {
	tm := new(Queue)
	testTask := NewTask(HighestPriority, func() error {
		return nil
	})
	for i := 0; i < b.N; i++ {
		tm.AddTask(testTask)
	}
	for i := 0; i < b.N; i++ {
		tm.GetTask()
	}
}

//тест на выполнение задач по приоритету
func TestPriority(t *testing.T) {
	tasks := []*Task{
		NewTask(HighestPriority, func() error { return nil }),
		NewTask(HighPriority, func() error { return nil }),
		NewTask(MiddlePriority, func() error { return nil }),
		NewTask(LowPriority, func() error { return nil }),
		NewTask(LowestPriority, func() error { return nil }),
	}

	// перемешиваем задачи
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(tasks), func(i, j int) { tasks[i], tasks[j] = tasks[j], tasks[i] })

	q := new(Queue)
	for _, t := range tasks {
		q.AddTask(t)
	}

	highest := q.GetTask()
	priority := highest.Priority()
	if priority != HighestPriority {
		t.Errorf(`unexpected priority: expect %d get %d"`, HighestPriority, priority)
	}

	high := q.GetTask()
	priority = high.Priority()
	if priority != HighPriority {
		t.Errorf(`unexpected priority: expect %d get %d"`, HighPriority, priority)
	}

	middle := q.GetTask()
	priority = middle.Priority()
	if priority != MiddlePriority {
		t.Errorf(`unexpected priority: expect %d get %d"`, MiddlePriority, priority)
	}

	low := q.GetTask()
	priority = low.Priority()
	if priority != LowPriority {
		t.Errorf(`unexpected priority: expect %d get %d"`, LowPriority, priority)
	}

	lowest := q.GetTask()
	priority = lowest.Priority()
	if priority != LowestPriority {
		t.Errorf(`unexpected priority: expect %d get %d"`, LowestPriority, priority)
	}

}

// Если очередь пуста метод GetTask должен вернуть nil
// Если в очереди одно задание GetTask должен вернуть именно это задание
func TestGetTask(t *testing.T) {
	q := new(Queue)
	if q.GetTask() != nil {
		t.Error(`unexpected TaskInterface, expect nil`)
	}
	testTask := NewTask(HighestPriority, func() error { return nil })
	q.AddTask(testTask)
	taskFromQueue := q.GetTask()
	if testTask != taskFromQueue {
		t.Error(`one task from queue is not equal task put in queue `)
	}
}

// сколько положили в очередь задач, столько и должны получить
func TestCountTasks(t *testing.T) {
	q := new(Queue)

	tasksIn := 64
	tasks := []*Task{
		NewTask(HighestPriority, func() error { return nil }),
		NewTask(HighPriority, func() error { return nil }),
		NewTask(MiddlePriority, func() error { return nil }),
		NewTask(LowPriority, func() error { return nil }),
		NewTask(LowestPriority, func() error { return nil }),
	}

	for i := 0; i < tasksIn; i++ {
		// перемешиваем примеры задач
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(tasks), func(i, j int) { tasks[i], tasks[j] = tasks[j], tasks[i] })

		q.AddTask(tasks[0])
	}
	var tasksOut int
	for q.GetTask() != nil {
		tasksOut++
	}

	if tasksIn != tasksOut {
		t.Errorf(`unexpected out tasks - %d, expect - %d`, tasksOut, tasksIn)
	}
}

// В случае ошибки выполнения должен сработать хэндлер для события FailedEvent
func TestErrorTask(t *testing.T) {
	failFlag := false

	task := NewTask(HighestPriority, func() error { return fmt.Errorf(`test error`) })
	task.OnEvent(FailedEvent, func() {
		failFlag = true
	})
	task.Exec()

	if !failFlag {
		t.Errorf(`expected execution of handler for failed event!`)
	}
}

// Обращаемся асинхронно к незащищенным данным чтобы создать условия для race condition
func TestRaceCondition(t *testing.T) {
	q := new(Queue)
	go func() {
		q.AddTask(NewTask(HighestPriority, func() error { return nil }))
	}()
	q.GetTask()
}
