package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/delgus/taskmanager"
)

func main() {
	q := new(taskmanager.Queue)
	// имитируем асинхронное добавление задач
	stopAddTasks := make(chan struct{})
	go func() {
		ticker := time.NewTicker(time.Millisecond * 100)
		for {
			select {
			case <-ticker.C:
				task := taskmanager.NewTask(taskmanager.HighestPriority, func() error {
					time.Sleep(time.Second * 5)
					fmt.Println("i highest! good work!")
					return nil
				})
				task.OnEvent(taskmanager.CreatedEvent, func() {
					fmt.Println("i highest! i created!")
				})
				task.OnEvent(taskmanager.BeforeExecEvent, func() {
					fmt.Println("i highest! i before execution!")
				})
				task.OnEvent(taskmanager.AfterExecEvent, func() {
					fmt.Println("i highest! i after execution!")
				})
				task.EmitEvent(taskmanager.CreatedEvent)
				q.AddTask(task)
			case <-stopAddTasks:
				return
			}
		}
	}()

	// обрабатываем задачи в 10 потоков
	worker := taskmanager.NewWorkerPool(q, 10, time.Millisecond*50)

	// нажмите CTRL + C для остановки воркера
	// плавная остановка воркера при получение interrupt сигнала
	go func() {
		sigint := make(chan os.Signal, 1)
		// получаем interrupt сигнал из терминала
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		stopAddTasks <- struct{}{}
		if err := worker.Shutdown(time.Second * 5); err != nil {
			fmt.Println(`error by stopping:` + err.Error())
		}
		fmt.Println(`stopping worker pool`)
	}()

	worker.Run()

}
