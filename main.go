package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Task struct
type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// File name
const fileName = "tasks.json"

// Load tasks from JSON
func loadTasks() ([]Task, error) {
	var tasks []Task
	file, err := os.Open(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil // file not found, return empty list
		}
		return nil, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&tasks)
	if err != nil {
		return []Task{}, nil // empty file
	}
	return tasks, nil
}

// Save tasks to JSON
func saveTasks(tasks []Task) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(tasks)
}

// Add a new task
func addTask(tasks []Task, description string) ([]Task, Task) {
	newTask := Task{
		ID:          len(tasks) + 1,
		Description: description,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	tasks = append(tasks, newTask)
	return tasks, newTask
}

// Update a task
func updateTask(tasks []Task, id int, newDesc, newStatus string) ([]Task, bool) {
	for i, t := range tasks {
		if t.ID == id {
			if newDesc != "" {
				tasks[i].Description = newDesc
			}
			if newStatus != "" {
				tasks[i].Status = newStatus
			}
			tasks[i].UpdatedAt = time.Now()
			return tasks, true
		}
	}
	return tasks, false
}

// Delete a task
func deleteTask(tasks []Task, id int) ([]Task, bool) {
	for i, t := range tasks {
		if t.ID == id {
			return append(tasks[:i], tasks[i+1:]...), true
		}
	}
	return tasks, false
}

// List tasks
func listTasks(tasks []Task) {
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}
	fmt.Println("📋 Task List:")
	for _, t := range tasks {
		fmt.Printf("[%d] %s | %s | Created: %s\n",
			t.ID, t.Description, t.Status, t.CreatedAt.Format("2006-01-02 15:04:05"))
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	tasks, _ := loadTasks()

	for i := 0; i < 2; i++ {
		fmt.Println("\n=== Task Manager ===")
		fmt.Println("1. Add Task")
		fmt.Println("2. Update Task")
		fmt.Println("3. Delete Task")
		fmt.Println("4. List Tasks")
		fmt.Println("5. Exit")
		fmt.Print("Choose an option: ")

		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)

		switch choiceStr {
		case "1":
			fmt.Print("Enter task description: ")
			desc, _ := reader.ReadString('\n')
			desc = strings.TrimSpace(desc)
			tasks, newTask := addTask(tasks, desc)
			saveTasks(tasks)
			fmt.Printf("✅ Task added successfully (ID: %d)\n", newTask.ID)

		case "2":
			listTasks(tasks)
			fmt.Print("Enter task ID to update: ")
			idStr, _ := reader.ReadString('\n')
			id, _ := strconv.Atoi(strings.TrimSpace(idStr))

			fmt.Print("New description (press Enter to skip): ")
			desc, _ := reader.ReadString('\n')
			desc = strings.TrimSpace(desc)

			fmt.Print("New status (pending/in progress/done): ")
			status, _ := reader.ReadString('\n')
			status = strings.TrimSpace(status)

			tasks, ok := updateTask(tasks, id, desc, status)
			if ok {
				saveTasks(tasks)
				fmt.Println("✅ Task updated successfully.")
			} else {
				fmt.Println("❌ Task not found.")
			}

		case "3":
			listTasks(tasks)
			fmt.Print("Enter task ID to delete: ")
			idStr, _ := reader.ReadString('\n')
			id, _ := strconv.Atoi(strings.TrimSpace(idStr))

			tasks, ok := deleteTask(tasks, id)
			if ok {
				saveTasks(tasks)
				fmt.Println("🗑️ Task deleted successfully.")
			} else {
				fmt.Println("❌ Task not found.")
			}

		case "4":
			listTasks(tasks)

		case "5":
			fmt.Println("👋 Exiting...")
			return

		default:
			fmt.Println("Invalid choice. Try again.")
		}
	}
}
