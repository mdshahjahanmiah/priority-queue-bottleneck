package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/mdshahjahanmiah/priority-queue-bottleneck/scheduler"
	"github.com/mdshahjahanmiah/priority-queue-bottleneck/utils"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	rand.Seed(time.Now().UnixNano())

	scheduler := scheduler.NewScheduler(10) // Queue size for CPU/GPU queues
	scheduler.AddWorkers(ctx, 3, 3)         // Start workers
	scheduler.ScheduleTasks(ctx)            // Start scheduling from priority queue

	// Generate and add tasks to the priority queue
	for i := 1; i <= 50; i++ {
		task := scheduler.CreateTask(i)
		scheduler.AddTask(task)
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}

	// Wait for a bit to let tasks process, then cancel
	time.Sleep(2 * time.Second)
	cancel()
	scheduler.Wait()
	utils.PrintFinalMetrics(scheduler.Tasks())
}
