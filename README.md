# Priority Queue Bottleneck Scheduler

This project implements a priority queue bottleneck scheduler in Go. It schedules tasks based on their CPU and GPU intensity and processes them accordingly.

## Features

- Task creation with random CPU and GPU intensity.
- Scheduling tasks to CPU or GPU queues based on their intensity.
- Processing tasks with multiple workers.
- Displaying task metrics in a tabular format.

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/mdshahjahanmiah/priority-queue-bottleneck.git
    ```
2. Navigate to the project directory:
    ```sh
    cd priority-queue-bottleneck
    ```
3. Install dependencies:
    ```sh
    go mod tidy
    ```

## Usage

1. Run the scheduler:
    ```sh
    go run main.go
    ```

2. The program will create and schedule 50 tasks, process them, and print the final metrics in a table format.

## Project Structure

- `main.go`: Entry point of the application.
- `scheduler/queue.go`: Contains the `Scheduler` and `Task` structs and task scheduling logic.
- `scheduler/processor.go`: Contains the task processing logic and worker management.
- `utils.go`: Contains utility functions for printing task metrics.

## Dependencies

- [github.com/olekukonko/tablewriter](https://github.com/olekukonko/tablewriter): Used for printing task metrics in a table format.

## License

This project is licensed under the MIT License.