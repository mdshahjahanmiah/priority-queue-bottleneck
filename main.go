package main

import (
	"container/list"
	"context"
	"github.com/mdshahjahanmiah/priority-queue-bottleneck/scheduler"
	"github.com/mdshahjahanmiah/priority-queue-bottleneck/utils"
	"math/rand"
	"time"
)

func main() {
	list.New()

	ctx, cancel := context.WithCancel(context.Background())
	rand.Seed(time.Now().UnixNano())

	scheduler := scheduler.NewScheduler(10) // 2 workers for each queue
	scheduler.AddWorkers(ctx, 3, 3)         // Start workers

	// Generate and schedule tasks
	for i := 1; i <= 50; i++ {
		task := scheduler.CreateTask(i)
		scheduler.ScheduleTask(ctx, task)
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}

	cancel()
	scheduler.Wait()
	utils.PrintFinalMetrics(scheduler.Tasks())
}
