package scheduler

import (
	"context"
	"testing"
	"time"
)

func Test_ScheduleTasks(t *testing.T) {
	s := NewScheduler(2)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// Add a CPU-bound task
	cpuTask := &Task{
		ID:           1,
		CPUIntensity: 2.0,
		GPUIntensity: 1.0,
	}
	s.AddTask(cpuTask)

	// Add a GPU-bound task
	gpuTask := &Task{
		ID:           2,
		CPUIntensity: 1.0,
		GPUIntensity: 2.0,
	}
	s.AddTask(gpuTask)

	// Start scheduling
	s.ScheduleTasks(ctx)

	// Wait for tasks to be scheduled
	time.Sleep(50 * time.Millisecond)

	// Check destinations
	tasks := s.Tasks()
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks to be scheduled, got %d", len(tasks))
	}

	for _, task := range tasks {
		if task.ID == 1 && task.Destination != "CPU" {
			t.Errorf("Expected CPU task destination 'CPU', got '%s'", task.Destination)
		}
		if task.ID == 2 && task.Destination != "GPU" {
			t.Errorf("Expected GPU task destination 'GPU', got '%s'", task.Destination)
		}
	}
}

func Test_ProcessTasks(t *testing.T) {
	s := NewScheduler(1)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	task := &Task{
		ID:           1,
		CPUIntensity: 1.0, // 100ms processing time
		GPUIntensity: 2.0,
	}
	s.cpuQueue <- task

	// Start a CPU worker
	s.wg.Add(1) // Explicitly add to WaitGroup since we're starting the goroutine manually
	go s.ProcessTasks(ctx, s.cpuQueue, false)

	// Wait for the task to be processed or timeout
	select {
	case <-time.After(200 * time.Millisecond): // Give enough time for processing (100ms + buffer)
		if task.TotalTime == 0 {
			t.Error("Task was not processed, TotalTime is still 0")
		}
	case <-ctx.Done():
		t.Error("Context canceled before task was processed")
	}

	// Cancel the context to stop the worker
	cancel()
	s.Wait() // Wait for the worker to finish
}
