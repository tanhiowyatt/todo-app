package main

import (
	"fmt"
	"os"
	"sort"
	"time"
)

type Task struct {
	ID        int
	Text      string
	Completed bool
	Priority  int
	DueDate   time.Time
}

var tasks []Task

func addTask(text string, priority int, dueDate time.Time) {
	task := Task{
		ID:        len(tasks) + 1,
		Text:      text,
		Completed: false,
		Priority:  priority,
		DueDate:   dueDate,
	}
	tasks = append(tasks, task)
	fmt.Println("Задача добавлена:", task)
}

func deleteTask(id int) {
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Println("Задача удалена:", task)
			return
		}
	}
	fmt.Println("Задача с ID", id, "не найдена.")
}

func markTaskCompleted(id int) {
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Completed = true
			fmt.Println("Задача выполнена:", task)
			return
		}
	}
	fmt.Println("Задача с ID", id, "не найдена.")
}

func listTasks() {
	fmt.Println("Список задач:")
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].DueDate.Before(tasks[j].DueDate)
	})
	for _, task := range tasks {
		status := "Не выполнена"
		if task.Completed {
			status = "Выполнена"
		}
		fmt.Printf("%d. %s (Приоритет: %d, Статус: %s, Срок: %s)\n",
			task.ID, task.Text, task.Priority, status, task.DueDate.Format("02.01.2006"))
	}
}

func main() {
	for {
		fmt.Println("\nВыберите действие:")
		fmt.Println("1. Добавить задачу")
		fmt.Println("2. Удалить задачу")
		fmt.Println("3. Отметить задачу как выполненную")
		fmt.Println("4. Вывести список задач")
		fmt.Println("5. Выход")

		var choice int
		fmt.Print("Введите номер действия: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			fmt.Print("Введите текст задачи: ")
			var text string
			fmt.Scanln()
			fmt.Scan(&text)

			fmt.Print("Введите приоритет (целое число): ")
			var priority int
			fmt.Scan(&priority)

			fmt.Print("Введите срок выполнения (в формате ГГГГ-ММ-ДД): ")
			var dueDateInput string
			fmt.Scan(&dueDateInput)
			dueDate, err := time.Parse("2006-01-02", dueDateInput)
			if err != nil {
				fmt.Println("Ошибка ввода даты:", err)
				continue
			}

			addTask(text, priority, dueDate)
		case 2:
			fmt.Print("Введите ID задачи для удаления: ")
			var id int
			fmt.Scan(&id)
			deleteTask(id)
		case 3:
			fmt.Print("Введите ID задачи для отметки как выполненной: ")
			var id int
			fmt.Scan(&id)
			markTaskCompleted(id)
		case 4:
			listTasks()
		case 5:
			fmt.Println("Программа завершена.")
			os.Exit(0)
		default:
			fmt.Println("Неверный выбор. Пожалуйста, выберите снова.")
		}
	}
}
