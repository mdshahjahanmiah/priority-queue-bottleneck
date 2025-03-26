package utils

import (
	"fmt"
	"github.com/mdshahjahanmiah/priority-queue-bottleneck/scheduler"
	"github.com/olekukonko/tablewriter"
	"os"
)

// PrintFinalMetrics prints task data in a table format
func PrintFinalMetrics(tasks []*scheduler.Task) {
	printTaskTable(tasks)
}

// printTaskTable generates a table output of tasks
func printTaskTable(tasks []*scheduler.Task) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Task ID", "Task Name", "CPU Intensity", "GPU Intensity", "Destination"})

	for _, task := range tasks {
		table.Append([]string{
			fmt.Sprintf("%d", task.ID),
			task.Name,
			fmt.Sprintf("%.2f", task.CPUIntensity),
			fmt.Sprintf("%.2f", task.GPUIntensity),
			task.Destination,
		})
	}
	table.Render()
}
