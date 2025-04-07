package main

import (
	"fmt"
	"net/http"
	"os"

	_ "task-recommender/docs"
	"task-recommender/internal/api"
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

	// ルーターの設定
	router := api.SetupRouter(taskController)

	fmt.Printf("サーバーを起動しています: 0.0.0.0:%s\n", port)
	http.ListenAndServe("0.0.0.0:"+port, router)
}
