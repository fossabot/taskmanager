package main

import (
	"fmt"

	"github.com/delgus/taskmanager"
)

func main() {
	tq := new(taskmanager.Queue)
	//создаем задачи и добавляем в очередь
	//lowest task
	lTask := taskmanager.NewTask(taskmanager.LowestPriority, func() error {
		fmt.Println("i lowest! good work!")
		return nil
	})
	lTask.OnEvent(taskmanager.CreatedEvent, func() {
		fmt.Println("i lowest!i created!")
	})
	lTask.OnEvent(taskmanager.BeforeExecEvent, func() {
		fmt.Println("i lowest! i before execution!")
	})
	lTask.OnEvent(taskmanager.AfterExecEvent, func() {
		fmt.Println("i lowest! i after execution!")
	})
	lTask.EmitEvent(taskmanager.CreatedEvent)
	tq.AddTask(lTask)

	//highest task
	hTask := taskmanager.NewTask(taskmanager.HighestPriority, func() error {
		fmt.Println("i highest! good work!")
		return nil
	})
	hTask.OnEvent(taskmanager.CreatedEvent, func() {
		fmt.Println("i highest!i created!")
	})
	hTask.OnEvent(taskmanager.BeforeExecEvent, func() {
		fmt.Println("i highest! i before execution!")
	})
	hTask.OnEvent(taskmanager.AfterExecEvent, func() {
		fmt.Println("i highest! i after execution!")
	})
	hTask.EmitEvent(taskmanager.CreatedEvent)
	tq.AddTask(hTask)

	mTask := taskmanager.NewTask(taskmanager.MiddlePriority, func() error {
		fmt.Println("i middle! good work!")
		return nil
	})
	mTask.OnEvent(taskmanager.CreatedEvent, func() {
		fmt.Println("i middle! i created!")
	})
	mTask.OnEvent(taskmanager.BeforeExecEvent, func() {
		fmt.Println("i middle! i before execution!")
	})
	mTask.OnEvent(taskmanager.AfterExecEvent, func() {
		fmt.Println("i middle! i after execution!")
	})
	mTask.EmitEvent(taskmanager.CreatedEvent)
	tq.AddTask(mTask)

	//broken task
	bTask := taskmanager.NewTask(taskmanager.HighestPriority, func() error {
		return fmt.Errorf("i broken! sorry(")
	})
	bTask.OnEvent(taskmanager.CreatedEvent, func() {
		fmt.Println("i highest! i created!")
	})
	bTask.OnEvent(taskmanager.FailedEvent, func() {
		fmt.Println("i broke! sorry")
	})
	bTask.OnEvent(taskmanager.BeforeExecEvent, func() {
		fmt.Println("i highest! i before execution!")
	})
	bTask.OnEvent(taskmanager.AfterExecEvent, func() {
		fmt.Println("i highest! i after execution!")
	})
	bTask.EmitEvent(taskmanager.CreatedEvent)
	tq.AddTask(bTask)

	for {
		task := tq.GetTask()
		if task == nil {
			break
		}
		task.Exec()
	}
}
