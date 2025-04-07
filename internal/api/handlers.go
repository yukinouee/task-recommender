package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"task-recommender/internal/controller"
)

type TaskHandler struct {
	controller *controller.TaskController
}

func NewTaskHandler(controller *controller.TaskController) *TaskHandler {
	return &TaskHandler{controller: controller}
}

// swagger:route GET /tasks tasks listTasks
// タスク一覧を取得します
// responses:
//   200: body:[]Task

// HandleListTasks タスク一覧を取得するハンドラー
func (h *TaskHandler) HandleListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.controller.ListTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// swagger:route POST /tasks tasks createTask
// 新しいタスクを作成します
// responses:
//   201: body:map[string]int

// HandleCreateTask 新しいタスクを作成するハンドラー
func (h *TaskHandler) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var task struct {
		Title             string `json:"title"`
		Description       string `json:"description"`
		Priority          int    `json:"priority"`
		DueDate           string `json:"due_date"`
		EstimatedDuration int    `json:"estimated_duration"`
	}

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var dueDate time.Time
	if task.DueDate != "" {
		var err error
		dueDate, err = time.Parse("2006-01-02", task.DueDate)
		if err != nil {
			http.Error(w, "日付の形式が不正です。YYYY-MM-DD形式で指定してください", http.StatusBadRequest)
			return
		}
	}

	id, err := h.controller.AddTask(task.Title, task.Description, task.Priority, dueDate, task.EstimatedDuration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// swagger:route PUT /tasks/{id}/complete tasks completeTask
// タスクを完了としてマークします
// parameters:
//   + name: id
//     in: path
//     type: integer
//     required: true
// responses:
//   200: body:map[string]string

// HandleCompleteTask タスクを完了としてマークするハンドラー
func (h *TaskHandler) HandleCompleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.controller.CompleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "completed"})
}

// swagger:route DELETE /tasks/{id} tasks deleteTask
// タスクを削除します
// parameters:
//   + name: id
//     in: path
//     type: integer
//     required: true
// responses:
//   200: body:map[string]string

// HandleDeleteTask タスクを削除するハンドラー
func (h *TaskHandler) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.controller.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

// swagger:route PUT /tasks/{id}/priority tasks updateTaskPriority
// タスクの優先度を更新します
// parameters:
//   + name: id
//     in: path
//     type: integer
//     required: true
//   + name: body
//     in: body
//     required: true
//     schema:
//       type: object
//       required:
//         - priority
//       properties:
//         priority:
//           type: integer
// responses:
//   200: body:map[string]string

// HandleUpdatePriority タスクの優先度を更新するハンドラー
func (h *TaskHandler) HandleUpdatePriority(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var data struct {
		Priority int `json:"priority"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.controller.UpdatePriority(id, data.Priority)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "priority updated"})
}

// swagger:route PUT /tasks/{id}/due tasks updateTaskDueDate
// タスクの期限日を更新します
// parameters:
//   + name: id
//     in: path
//     type: integer
//     required: true
//   + name: body
//     in: body
//     required: true
//     schema:
//       type: object
//       required:
//         - due_date
//       properties:
//         due_date:
//           type: string
//           format: date
// responses:
//   200: body:map[string]string

// HandleUpdateDueDate タスクの期限日を更新するハンドラー
func (h *TaskHandler) HandleUpdateDueDate(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var data struct {
		DueDate string `json:"due_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dueDate, err := time.Parse("2006-01-02", data.DueDate)
	if err != nil {
		http.Error(w, "日付の形式が不正です。YYYY-MM-DD形式で指定してください", http.StatusBadRequest)
		return
	}

	err = h.controller.UpdateDueDate(id, dueDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "due date updated"})
}

// swagger:route PUT /tasks/{id}/duration tasks updateTaskDuration
// タスクの見積時間を更新します
// parameters:
//   + name: id
//     in: path
//     type: integer
//     required: true
//   + name: body
//     in: body
//     required: true
//     schema:
//       type: object
//       required:
//         - duration
//       properties:
//         duration:
//           type: integer
// responses:
//   200: body:map[string]string

// HandleUpdateEstimatedDuration タスクの見積時間を更新するハンドラー
func (h *TaskHandler) HandleUpdateEstimatedDuration(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var data struct {
		Duration int `json:"duration"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.controller.UpdateEstimatedDuration(id, data.Duration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "duration updated"})
}

// getIDFromPath URLパスからIDを抽出するヘルパー関数
func getIDFromPath(path string) (int, error) {
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return 0, fmt.Errorf("invalid path")
	}
	return strconv.Atoi(parts[2])
}
