package main

import (
	"fmt"
	"time"
)

func main() {
	scheduler := NewScheduler(
		func(msg string) {
			fmt.Println(msg)
		},
	)

	task := &Task{
		Name:       "example",
		Interval:   time.Second * 5,
		MaxRetries: 3,
		Body: func() error {
			fmt.Println("Running task")
			return nil
		},
	}
	scheduler.AddTask(task)
	scheduler.Start()
	defer scheduler.Stop()

	time.Sleep(time.Minute)
}
