package todo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Priority          string    `json:"priority"`
	DueDate           time.Time `json:"dueDate"`
	EstimatedDuration int       `json:"estimatedDuration"`
	CreatedAt         time.Time `json:"createdAt"`
	CompletedAt       time.Time `json:"completedAt"`
	IsDone            bool      `json:"isDone"`
}

type TaskList struct {
	Tasks  []Task
	apiURL string
}

func NewTaskList(apiURL string) *TaskList {
	return &TaskList{
		Tasks:  make([]Task, 0),
		apiURL: apiURL,
	}
}

func (tl *TaskList) FetchTasks() error {
	resp, err := http.Get(tl.apiURL + "/tasks")
	if err != nil {
		return fmt.Errorf("APIリクエストに失敗: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("レスポンスの読み取りに失敗: %v", err)
	}

	var tasks []Task
	if err := json.Unmarshal(body, &tasks); err != nil {
		return fmt.Errorf("JSONのパースに失敗: %v", err)
	}

	tl.Tasks = tasks
	return nil
}

func (tl *TaskList) AddTask(name string, priority string, dueDate time.Time, estimatedDuration int) string {
	task := Task{
		ID:                uuid.New().String(),
		Name:              name,
		Priority:          priority,
		DueDate:           dueDate,
		EstimatedDuration: estimatedDuration,
		CreatedAt:         time.Now(),
		CompletedAt:       time.Time{},
		IsDone:            false,
	}

	jsonData, err := json.Marshal(task)
	if err != nil {
		fmt.Printf("JSONの生成に失敗: %v\n", err)
		return ""
	}

	resp, err := http.Post(tl.apiURL+"/tasks", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("APIリクエストに失敗: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	var result struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("レスポンスの解析に失敗: %v\n", err)
		return ""
	}

	tl.FetchTasks() // タスク一覧を更新
	return result.ID
}

func (tl *TaskList) DeleteTask(id string) bool {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/tasks/%s", tl.apiURL, id), nil)
	if err != nil {
		fmt.Printf("リクエストの作成に失敗: %v\n", err)
		return false
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("APIリクエストに失敗: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return false
	}

	tl.FetchTasks() // タスク一覧を更新
	return true
}

func (tl *TaskList) MarkComplete(task Task) {
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/tasks/%s/complete", tl.apiURL, task.ID), nil)
	if err != nil {
		fmt.Printf("リクエストの作成に失敗: %v\n", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("APIリクエストに失敗: %v\n", err)
		return
	}
	defer resp.Body.Close()

	tl.FetchTasks() // タスク一覧を更新
}

func (tl TaskList) ShowTasks() {
	// ヘッダーの出力
	fmt.Printf("%-20s %-8s %-19s %-8s %-5s\n",
		"タスク名", "優先度", "期限", "見積(分)", "完了")
	fmt.Println(strings.Repeat("-", 100))

	// タスクの出力
	for _, task := range tl.Tasks {
		status := "未"
		if task.IsDone {
			status = "済"
		}
		fmt.Printf("%-20s %-8s %-19s %-8d %-5s\n",
			task.Name,
			task.Priority,
			task.DueDate.Format("2006/01/02-15:04"),
			task.EstimatedDuration,
			status)
	}
}
