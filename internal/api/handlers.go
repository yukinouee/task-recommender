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

// @Summary タスク一覧を取得
// @Description すべてのタスクの一覧を取得します
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {array} model.Task
// @Router /tasks [get]
func (h *TaskHandler) HandleListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.controller.ListTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// @Summary 新しいタスクを作成
// @Description タイトル、説明、優先度、期限日、見積時間を指定して新しいタスクを作成します
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body object true "タスク情報"
// @Success 201 {object} map[string]int
// @Failure 400 {object} string "不正なリクエスト"
// @Failure 500 {object} string "サーバーエラー"
// @Router /tasks [post]
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

// @Summary タスクを完了としてマーク
// @Description 指定されたIDのタスクを完了状態に更新します
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "タスクID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} string "不正なリクエスト"
// @Failure 500 {object} string "サーバーエラー"
// @Router /tasks/{id}/complete [put]
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

// @Summary タスクを削除
// @Description 指定されたIDのタスクを削除します
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "タスクID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} string "不正なリクエスト"
// @Failure 500 {object} string "サーバーエラー"
// @Router /tasks/{id} [delete]
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

// @Summary タスクの優先度を更新
// @Description 指定されたIDのタスクの優先度を更新します
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "タスクID"
// @Param priority body object true "優先度情報"
// @Success 200 {object} map[string]string
// @Failure 400 {object} string "不正なリクエスト"
// @Failure 500 {object} string "サーバーエラー"
// @Router /tasks/{id}/priority [put]
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

// @Summary タスクの期限日を更新
// @Description 指定されたIDのタスクの期限日を更新します
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "タスクID"
// @Param dueDate body object true "期限日情報"
// @Success 200 {object} map[string]string
// @Failure 400 {object} string "不正なリクエスト"
// @Failure 500 {object} string "サーバーエラー"
// @Router /tasks/{id}/due [put]
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

// @Summary タスクの見積時間を更新
// @Description 指定されたIDのタスクの見積時間を更新します
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "タスクID"
// @Param duration body object true "見積時間情報"
// @Success 200 {object} map[string]string
// @Failure 400 {object} string "不正なリクエスト"
// @Failure 500 {object} string "サーバーエラー"
// @Router /tasks/{id}/duration [put]
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
