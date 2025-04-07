package api

import (
	"net/http"
	"strings"

	"task-recommender/internal/controller"
)

// SetupRouter ルーターを設定
func SetupRouter(taskController *controller.TaskController) http.Handler {
	mux := http.NewServeMux()

	// APIハンドラーの作成
	taskHandler := NewTaskHandler(taskController)

	// ルートパス
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(`
                <html>
                    <head>
                        <title>タスク管理アプリケーション</title>
                    </head>
                    <body>
                        <h1>タスク管理アプリケーションのAPIサーバー</h1>
                        <p><a href="/swagger/index.html">API Document by Swagger UI</a></p>
                    </body>
                </html>
            `))
			return
		}
		http.NotFound(w, r)
	})

	// タスク一覧の取得と追加
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.HandleListTasks(w, r)
		case http.MethodPost:
			taskHandler.HandleCreateTask(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// 個別のタスク操作
	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// 完了マーク: /tasks/{id}/complete
		if strings.HasSuffix(path, "/complete") {
			if r.Method == http.MethodPut {
				taskHandler.HandleCompleteTask(w, r)
				return
			}
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 優先度更新: /tasks/{id}/priority
		if strings.HasSuffix(path, "/priority") {
			if r.Method == http.MethodPut {
				taskHandler.HandleUpdatePriority(w, r)
				return
			}
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 期限日更新: /tasks/{id}/due
		if strings.HasSuffix(path, "/due") {
			if r.Method == http.MethodPut {
				taskHandler.HandleUpdateDueDate(w, r)
				return
			}
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 見積時間更新: /tasks/{id}/duration
		if strings.HasSuffix(path, "/duration") {
			if r.Method == http.MethodPut {
				taskHandler.HandleUpdateEstimatedDuration(w, r)
				return
			}
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// タスク削除: /tasks/{id}
		if r.Method == http.MethodDelete {
			taskHandler.HandleDeleteTask(w, r)
			return
		}

		http.Error(w, "Method not allowed or invalid path", http.StatusMethodNotAllowed)
	})

	// Swagger UI
	mux.Handle("/swagger/", http.StripPrefix("/swagger", NewSwaggerHandler()))

	return mux
}
