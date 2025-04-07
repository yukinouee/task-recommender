package main

import (
	"flag"
	"fmt"
	"os"
	"task-recommender/internal/model"

	"task-recommender/internal/controller"
	"task-recommender/internal/service"
	"task-recommender/internal/view"
	"task-recommender/pkg/db"
)

func main() {
	// コマンドラインフラグの設定
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addTitle := addCmd.String("title", "", "タスクのタイトル")
	addDesc := addCmd.String("desc", "", "タスクの説明")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	completeCmd := flag.NewFlagSet("complete", flag.ExitOnError)
	completeID := completeCmd.Int("id", 0, "完了にするタスクのID")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteID := deleteCmd.Int("id", 0, "削除するタスクのID")

	// データベース接続
	database, err := db.Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "データベース接続エラー: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// サービスとコントローラーの初期化
	taskService := service.NewTaskService(database)
	taskController := controller.NewTaskController(taskService)

	// コマンドライン引数がない場合はヘルプを表示
	if len(os.Args) < 2 {
		fmt.Println("使用方法: todo [add|list|complete|delete] [オプション]")
		os.Exit(1)
	}

	// サブコマンドに基づいて処理を分岐
	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		if *addTitle == "" {
			fmt.Println("タイトルは必須です")
			os.Exit(1)
		}
		id, err := taskController.AddTask(*addTitle, *addDesc)
		if err != nil {
			view.PrintError(err)
			os.Exit(1)
		}
		view.PrintTaskAdded(id, *addTitle)

	case "list":
		listCmd.Parse(os.Args[2:])
		tasks, err := taskController.ListTasks()
		if err != nil {
			view.PrintError(err)
			os.Exit(1)
		}
		view.PrintTaskList(tasks.([]model.Task))

	case "complete":
		completeCmd.Parse(os.Args[2:])
		if *completeID <= 0 {
			fmt.Println("有効なIDを指定してください")
			os.Exit(1)
		}
		err := taskController.CompleteTask(*completeID)
		if err != nil {
			view.PrintError(err)
			os.Exit(1)
		}
		view.PrintTaskCompleted(*completeID)

	case "delete":
		deleteCmd.Parse(os.Args[2:])
		if *deleteID <= 0 {
			fmt.Println("有効なIDを指定してください")
			os.Exit(1)
		}
		err := taskController.DeleteTask(*deleteID)
		if err != nil {
			view.PrintError(err)
			os.Exit(1)
		}
		view.PrintTaskDeleted(*deleteID)

	default:
		fmt.Println("不明なコマンドです。add, list, complete, deleteのいずれかを使用してください")
		os.Exit(1)
	}
}
