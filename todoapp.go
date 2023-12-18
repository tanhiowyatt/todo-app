package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
)

type Task struct {
	ID        int
	Text      string
	Completed bool
	Priority  int
	DueDate   time.Time
}

var tasks []Task

const filename = "tasks.json"

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
	saveTasksToFile()
}

func editTask(id int, text string, priority int, dueDate time.Time) {
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Text = text
			tasks[i].Priority = priority
			tasks[i].DueDate = dueDate
			fmt.Println("Задача изменена:", tasks[i])
			saveTasksToFile()
			return
		}
	}
	fmt.Println("Задача с ID", id, "не найдена.")
}

func deleteTask(id int) {
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Println("Задача удалена:", task)
			saveTasksToFile()
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
			saveTasksToFile()
			return
		}
	}
	fmt.Println("Задача с ID", id, "не найдена.")
}

func listTasks(showCompleted bool) {
	fmt.Println("Список задач:")
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].DueDate.Before(tasks[j].DueDate)
	})
	for _, task := range tasks {
		if !showCompleted && task.Completed {
			continue
		}
		status := "Не выполнена"
		if task.Completed {
			status = "Выполнена"
		}
		fmt.Printf("%d. %s (Приоритет: %d, Статус: %s, Срок: %s)\n",
			task.ID, task.Text, task.Priority, status, task.DueDate.Format("02.01.2006"))
	}
}

func saveTasksToFile() {
	data, err := json.Marshal(tasks)
	if err != nil {
		fmt.Println("Ошибка при маршалинге данных:", err)
		return
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
	}
}

func loadTasksFromFile() {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Ошибка при чтении из файла:", err)
		return
	}

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		fmt.Println("Ошибка при демаршалинге данных:", err)
		return
	}
}

func parseDueDateInput(input string) (time.Time, error) {
	formats := []string{"02.01.2006", "2006-01-02", "01/02/2006", "01-02-2006"}
	for _, format := range formats {
		parsedDate, err := time.Parse(format, input)
		if err == nil {
			return parsedDate, nil
		}
	}
	return time.Time{}, fmt.Errorf("неверный формат даты")
}

func searchTasks(keywordsInput string) {
	keywords := strings.Fields(keywordsInput)
	matchedTasks := make([]Task, 0)

	for _, task := range tasks {
		taskTextLower := strings.ToLower(task.Text)
		for _, keyword := range keywords {
			keywordLower := strings.ToLower(keyword)
			if strings.Contains(taskTextLower, keywordLower) {
				matchedTasks = append(matchedTasks, task)
				break
			}
		}
	}

	if len(matchedTasks) > 0 {
		fmt.Println("Найденные задачи:")
		for _, task := range matchedTasks {
			status := "Не выполнена"
			if task.Completed {
				status = "Выполнена"
			}
			fmt.Printf("%d. %s (Приоритет: %d, Статус: %s, Срок: %s)\n",
				task.ID, task.Text, task.Priority, status, task.DueDate.Format("02.01.2006"))
		}
	} else {
		fmt.Println("Нет совпадений для введенных ключевых слов.")
	}
}

func main() {
	loadTasksFromFile()

	for {
		prompt := &survey.Select{
			Message: "Выберите действие:",
			Options: []string{
				"Добавить задачу",
				"Редактировать задачу",
				"Удалить задачу",
				"Отметить задачу как выполненную",
				"Вывести список задач",
				"Вывести список выполненных задач",
				"Поиск по ключевым словам",
				"Сохранить и выйти",
			},
		}
		var choice string
		survey.AskOne(prompt, &choice)

		switch choice {
		case "Добавить задачу":
			var text, dueDateInput string
			var priority int

			survey.AskOne(&survey.Input{Message: "Введите текст задачи:"}, &text)
			survey.AskOne(&survey.Input{Message: "Введите приоритет (целое число):"}, &priority)
			survey.AskOne(&survey.Input{Message: "Введите срок выполнения (в формате DD.MM.YYYY):"}, &dueDateInput)

			dueDate, err := parseDueDateInput(dueDateInput)
			if err != nil {
				fmt.Println("Ошибка ввода даты:", err)
				continue
			}

			addTask(text, priority, dueDate)

		case "Редактировать задачу":
			var id, priority int
			var text, dueDateInput string

			survey.AskOne(&survey.Input{Message: "Введите ID задачи для редактирования:"}, &id)
			survey.AskOne(&survey.Input{Message: "Введите новый текст задачи:"}, &text)
			survey.AskOne(&survey.Input{Message: "Введите новый приоритет (целое число):"}, &priority)
			survey.AskOne(&survey.Input{Message: "Введите новый срок выполнения (в формате DD.MM.YYYY):"}, &dueDateInput)

			dueDate, err := parseDueDateInput(dueDateInput)
			if err != nil {
				fmt.Println("Ошибка ввода даты:", err)
				continue
			}

			editTask(id, text, priority, dueDate)

		case "Удалить задачу":
			var id int
			survey.AskOne(&survey.Input{Message: "Введите ID задачи для удаления:"}, &id)
			deleteTask(id)

		case "Отметить задачу как выполненную":
			var id int
			survey.AskOne(&survey.Input{Message: "Введите ID задачи для отметки как выполненной:"}, &id)
			markTaskCompleted(id)

		case "Вывести список задач":
			listTasks(false)

		case "Вывести список выполненных задач":
			listTasks(true)

		case "Поиск по ключевым словам":
			var keywordsInput string
			survey.AskOne(&survey.Input{Message: "Введите ключевые слова для поиска:"}, &keywordsInput)
			searchTasks(keywordsInput)

		case "Сохранить и выйти":
			fmt.Println("Сохранение данных и завершение программы.")
			saveTasksToFile()
			os.Exit(0)

		default:
			fmt.Println("Неверный выбор. Пожалуйста, выберите снова.")
		}
	}
}
