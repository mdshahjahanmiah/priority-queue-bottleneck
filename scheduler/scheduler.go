package scheduler

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Task struct {
	ID           int
	Name         string
	CPUIntensity float64
	GPUIntensity float64
	TotalTime    time.Duration
	Destination  string
}

type Scheduler struct {
	cpuQueue chan *Task
	gpuQueue chan *Task
	wg       sync.WaitGroup
	tasks    []*Task // Store tasks for tabular output
}

func NewScheduler(queueSize int) *Scheduler {
	return &Scheduler{
		cpuQueue: make(chan *Task, queueSize),
		gpuQueue: make(chan *Task, queueSize),
		tasks:    []*Task{},
	}
}

func (s *Scheduler) CreateTask(id int) *Task {
	task := &Task{
		ID:           id,
		Name:         fmt.Sprintf("Task-%d", id),
		CPUIntensity: rand.Float64() * 10,
		GPUIntensity: rand.Float64() * 10,
	}
	return task
}

// ScheduleTask routes tasks to the appropriate queue
func (s *Scheduler) ScheduleTask(ctx context.Context, task *Task) {
	useGPU := task.GPUIntensity > task.CPUIntensity

	if useGPU {
		select {
		case s.gpuQueue <- task:
			task.Destination = "GPU"
		default:
			select {
			case s.cpuQueue <- task:
				task.Destination = "CPU (GPU full)"
			default:
				return
			}
		}
	} else {
		select {
		case s.cpuQueue <- task:
			task.Destination = "CPU"
		default:
			select {
			case s.gpuQueue <- task:
				task.Destination = "GPU (CPU full)"
			default:
				return
			}
		}
	}

	s.tasks = append(s.tasks, task)
}

// AddWorkers adds workers to process tasks
func (s *Scheduler) AddWorkers(ctx context.Context, cpuWorkers, gpuWorkers int) {
	for i := 0; i < cpuWorkers; i++ {
		s.wg.Add(1)
		go s.ProcessTasks(ctx, s.cpuQueue, false)
	}
	for i := 0; i < gpuWorkers; i++ {
		s.wg.Add(1)
		go s.ProcessTasks(ctx, s.gpuQueue, true)
	}
}

// ProcessTasks processes tasks in a given queue
func (s *Scheduler) ProcessTasks(ctx context.Context, queue chan *Task, isGPU bool) {
	defer s.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-queue:
			if !ok {
				return
			}

			if isGPU {
				time.Sleep(time.Duration(task.GPUIntensity*100) * time.Millisecond)
			} else {
				time.Sleep(time.Duration(task.CPUIntensity*100) * time.Millisecond)
			}
			task.TotalTime = time.Duration(task.CPUIntensity*100) * time.Millisecond
		}
	}
}

func (s *Scheduler) Tasks() []*Task {
	return s.tasks
}

// Wait waits for all tasks to complete
func (s *Scheduler) Wait() {
	s.wg.Wait()
}
