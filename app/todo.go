package app

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Priority int

const (
	PriorityNone   Priority = 0
	PriorityHigh   Priority = 1
	PriorityMedium Priority = 2
	PriorityLow    Priority = 3
)

type TodoItem struct {
	Title    string   `json:"title"`
	Done     bool     `json:"done"`
	Priority Priority `json:"priority"`
}

func todoFilePath(username string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".aimssh_todos.json"
	}
	suffix := ""
	if username != "" {
		suffix = "_" + username
	}
	return filepath.Join(home, ".aimssh_todos"+suffix+".json")
}

func LoadTodos(username string) []TodoItem {
	data, err := os.ReadFile(todoFilePath(username))
	if err != nil {
		return []TodoItem{}
	}
	var todos []TodoItem
	if err := json.Unmarshal(data, &todos); err != nil {
		return []TodoItem{}
	}
	return todos
}

func SaveTodos(username string, todos []TodoItem) {
	data, err := json.Marshal(todos)
	if err != nil {
		return
	}
	_ = os.WriteFile(todoFilePath(username), data, 0644)
}
