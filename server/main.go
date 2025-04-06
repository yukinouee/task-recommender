package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"task-recommender/pkg/todo"

	"github.com/gin-gonic/gin"
)

var taskList todo.TaskList

func main() {
	// 本番環境用の設定
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())   // ログ出力を有効化
	r.Use(gin.Recovery()) // パニック時の回復を有効化

	// CORS設定
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// ヘルスチェックエンドポイント
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"port":   os.Getenv("PORT"),
		})
	})

	// タスク一覧の取得
	r.GET("/tasks", func(c *gin.Context) {
		c.JSON(http.StatusOK, taskList.Tasks)
	})

	// タスクの追加
	r.POST("/tasks", func(c *gin.Context) {
		var task todo.Task
		if err := c.BindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		id := taskList.AddTask(task.Name, task.Priority, task.DueDate, task.EstimatedDuration)
		c.JSON(http.StatusCreated, gin.H{"id": id})
	})

	// タスクの削除
	r.DELETE("/tasks/:id", func(c *gin.Context) {
		id := c.Param("id")
		if taskList.DeleteTask(id) {
			c.Status(http.StatusNoContent)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		}
	})

	// タスクの完了
	r.PUT("/tasks/:id/complete", func(c *gin.Context) {
		id := c.Param("id")
		for _, task := range taskList.Tasks {
			if task.ID == id {
				taskList.MarkComplete(task)
				c.Status(http.StatusNoContent)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
	})

	// タスクの永続化（ファイルへの保存）
	saveTasks := func() {
		data, err := json.Marshal(taskList)
		if err != nil {
			log.Printf("タスクの保存に失敗: %v", err)
			return
		}
		err = os.WriteFile("tasks.json", data, 0644)
		if err != nil {
			log.Printf("ファイルの書き込みに失敗: %v", err)
		}
	}

	// 保存されたタスクの読み込み
	if data, err := os.ReadFile("tasks.json"); err == nil {
		if err := json.Unmarshal(data, &taskList); err != nil {
			log.Printf("タスクの読み込みに失敗: %v", err)
		}
	}

	// サーバー終了時にタスクを保存
	defer saveTasks()

	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	fmt.Printf("Starting server on %s...\n", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}
