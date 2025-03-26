package scheduler

import (
	"container/heap"
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

	priorityQueue PriorityQueue // Main queue of tasks
	tasks         []*Task       // Store tasks for tabular output for showcase

	wg sync.WaitGroup
	mu sync.Mutex
}

func NewScheduler(queueSize int) *Scheduler {
	return &Scheduler{
		priorityQueue: PriorityQueue{},
		cpuQueue:      make(chan *Task, queueSize),
		gpuQueue:      make(chan *Task, queueSize),
		tasks:         []*Task{},
	}
}

// CreateTask creates a new task with random CPU and GPU intensities
func (s *Scheduler) CreateTask(id int) *Task {
	task := &Task{
		ID:           id,
		Name:         fmt.Sprintf("Task-%d", id),
		CPUIntensity: rand.Float64() * 10,
		GPUIntensity: rand.Float64() * 10,
	}
	return task
}

// AddTask adds a task to the priority queue
func (s *Scheduler) AddTask(task *Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	heap.Push(&s.priorityQueue, task)
}

// ScheduleTasks continuously schedules tasks from the priority queue
func (s *Scheduler) ScheduleTasks(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				s.mu.Lock()
				if s.priorityQueue.Len() == 0 {
					s.mu.Unlock()
					time.Sleep(10 * time.Millisecond) // Avoid busy-waiting
					continue
				}
				task := heap.Pop(&s.priorityQueue).(*Task)
				s.mu.Unlock()

				// Schedule the task to the appropriate queue
				s.scheduleTask(task)
			}
		}
	}()
}

// scheduleTask routes a task to the appropriate queue based on the slide's rules
func (s *Scheduler) scheduleTask(task *Task) {
	fasterOnGPU := task.GPUIntensity > task.CPUIntensity

	if fasterOnGPU {
		// Try GPU first
		select {
		case s.gpuQueue <- task:
			task.Destination = "GPU"
		default:
			// GPU is full, check CPU
			select {
			case s.cpuQueue <- task:
				task.Destination = "CPU (GPU full)"
			default:
				// Both are full, put back in priority queue
				s.AddTask(task)
				return
			}
		}
	} else {
		// Try CPU first
		select {
		case s.cpuQueue <- task:
			task.Destination = "CPU"
		default:
			// CPU is full, check GPU
			select {
			case s.gpuQueue <- task:
				task.Destination = "GPU (CPU full)"
			default:
				// Both are full, put back in priority queue
				s.AddTask(task)
				return
			}
		}
	}

	s.mu.Lock()
	s.tasks = append(s.tasks, task)
	s.mu.Unlock()
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
				task.TotalTime = time.Duration(task.GPUIntensity*100) * time.Millisecond
			} else {
				time.Sleep(time.Duration(task.CPUIntensity*100) * time.Millisecond)
				task.TotalTime = time.Duration(task.CPUIntensity*100) * time.Millisecond
			}
		}
	}
}

// Tasks returns all tasks
func (s *Scheduler) Tasks() []*Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.tasks
}

// Wait waits for all tasks to complete
func (s *Scheduler) Wait() {
	s.wg.Wait()
}
