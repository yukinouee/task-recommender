package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"task-recommender/internal/controller"
	"task-recommender/internal/service"
	"task-recommender/pkg/db"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}

	// データベース接続
	database, err := db.Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "データベース接続エラー: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// データベース初期化
	err = db.InitializeDatabase(database)
	if err != nil {
		fmt.Fprintf(os.Stderr, "データベース初期化エラー: %v\n", err)
		os.Exit(1)
	}

	// サービスとコントローラーの初期化
	taskService := service.NewTaskService(database)
	taskController := controller.NewTaskController(taskService)

	// Webサーバーとして動作
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "タスク管理アプリケーションのAPIサーバーです")
	})

	// タスク一覧の取得と追加
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			tasks, err := taskController.ListTasks()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks)
			return
		}

		if r.Method == "POST" {
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

			id, err := taskController.AddTask(task.Title, task.Description, task.Priority, dueDate, task.EstimatedDuration)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]int{"id": id})
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// 個別のタスク操作
	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 3 {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		idStr := pathParts[2]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		// タスクの完了
		if len(pathParts) >= 4 && pathParts[3] == "complete" && r.Method == "PUT" {
			err := taskController.CompleteTask(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "completed"})
			return
		}

		// タスクの削除
		if r.Method == "DELETE" {
			err := taskController.DeleteTask(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
			return
		}

		// 優先度の更新
		if len(pathParts) >= 4 && pathParts[3] == "priority" && r.Method == "PUT" {
			var data struct {
				Priority int `json:"priority"`
			}

			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err := taskController.UpdatePriority(id, data.Priority)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "priority updated"})
			return
		}

		// 期限日の更新
		if len(pathParts) >= 4 && pathParts[3] == "due" && r.Method == "PUT" {
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

			err = taskController.UpdateDueDate(id, dueDate)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "due date updated"})
			return
		}

		// 見積時間の更新
		if len(pathParts) >= 4 && pathParts[3] == "duration" && r.Method == "PUT" {
			var data struct {
				Duration int `json:"duration"`
			}

			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err := taskController.UpdateEstimatedDuration(id, data.Duration)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "duration updated"})
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	fmt.Printf("サーバーを起動しています: 0.0.0.0:%s\n", port)
	http.ListenAndServe("0.0.0.0:"+port, nil)
}
