// Package scheduler provides a priority-based task scheduling mechanism.
// Tasks with a larger absolute difference between GPU and CPU intensity
// are given higher priority to optimize execution.
package scheduler

import "math"

type PriorityQueue []*Task

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// Higher intensity difference gets higher priority
	diffI := math.Abs(pq[i].GPUIntensity - pq[i].CPUIntensity)
	diffJ := math.Abs(pq[j].GPUIntensity - pq[j].CPUIntensity)
	return diffI > diffJ // Max heap: larger difference = higher priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Task))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
