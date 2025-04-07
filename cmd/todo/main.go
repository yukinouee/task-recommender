package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"task-recommender/internal/controller"
	"task-recommender/internal/model"
	"task-recommender/internal/service"
	"task-recommender/internal/view"
	"task-recommender/pkg/db"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000" // Renderのデフォルトポート
	}

	// コマンドラインフラグの設定
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addTitle := addCmd.String("title", "", "タスクのタイトル")
	addDesc := addCmd.String("desc", "", "タスクの説明")
	addPriority := addCmd.Int("priority", 1, "タスクの優先度 (1=低, 2=中, 3=高)")
	addDueDate := addCmd.String("due", "", "タスクの期限 (YYYY-MM-DD形式)")
	addDuration := addCmd.Int("duration", 0, "タスクの見積所要時間（分）")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	completeCmd := flag.NewFlagSet("complete", flag.ExitOnError)
	completeID := completeCmd.Int("id", 0, "完了にするタスクのID")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteID := deleteCmd.Int("id", 0, "削除するタスクのID")

	priorityCmd := flag.NewFlagSet("priority", flag.ExitOnError)
	priorityID := priorityCmd.Int("id", 0, "優先度を変更するタスクのID")
	priorityValue := priorityCmd.Int("value", 1, "新しい優先度 (1=低, 2=中, 3=高)")

	dueDateCmd := flag.NewFlagSet("due", flag.ExitOnError)
	dueDateID := dueDateCmd.Int("id", 0, "期限を変更するタスクのID")
	dueDateValue := dueDateCmd.String("date", "", "新しい期限 (YYYY-MM-DD形式)")

	durationCmd := flag.NewFlagSet("duration", flag.ExitOnError)
	durationID := durationCmd.Int("id", 0, "見積時間を変更するタスクのID")
	durationValue := durationCmd.Int("minutes", 0, "新しい見積時間（分）")

	// データベース接続
	database, err := db.Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "データベース接続エラー: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	err = db.InitializeDatabase(database)
	if err != nil {
		fmt.Fprintf(os.Stderr, "データベース初期化エラー: %v\n", err)
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "タスク管理アプリケーションのAPIサーバーです")
	})

	fmt.Printf("サーバーを起動しています: 0.0.0.0:%s\n", port)
	http.ListenAndServe("0.0.0.0:"+port, nil)

	// サービスとコントローラーの初期化
	taskService := service.NewTaskService(database)
	taskController := controller.NewTaskController(taskService)

	// コマンドライン引数がない場合はヘルプを表示
	if len(os.Args) < 2 {
		fmt.Println("使用方法: todo [add|list|complete|delete|priority|due|duration] [オプション]")
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

		var dueDate time.Time
		if *addDueDate != "" {
			var err error
			dueDate, err = time.Parse("2006-01-02", *addDueDate)
			if err != nil {
				fmt.Println("期限の形式が正しくありません。YYYY-MM-DD形式で指定してください")
				os.Exit(1)
			}
		}

		id, err := taskController.AddTask(*addTitle, *addDesc, *addPriority, dueDate, *addDuration)
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

	case "priority":
		priorityCmd.Parse(os.Args[2:])
		if *priorityID <= 0 {
			fmt.Println("有効なIDを指定してください")
			os.Exit(1)
		}
		if *priorityValue < 1 || *priorityValue > 3 {
			fmt.Println("優先度は1から3の間で指定してください")
			os.Exit(1)
		}
		err := taskController.UpdatePriority(*priorityID, *priorityValue)
		if err != nil {
			view.PrintError(err)
			os.Exit(1)
		}
		view.PrintPriorityUpdated(*priorityID, *priorityValue)

	case "due":
		dueDateCmd.Parse(os.Args[2:])
		if *dueDateID <= 0 {
			fmt.Println("有効なIDを指定してください")
			os.Exit(1)
		}

		var dueDate time.Time
		if *dueDateValue != "" {
			var err error
			dueDate, err = time.Parse("2006-01-02", *dueDateValue)
			if err != nil {
				fmt.Println("期限の形式が正しくありません。YYYY-MM-DD形式で指定してください")
				os.Exit(1)
			}
		} else {
			fmt.Println("期限日を指定してください")
			os.Exit(1)
		}

		err := taskController.UpdateDueDate(*dueDateID, dueDate)
		if err != nil {
			view.PrintError(err)
			os.Exit(1)
		}
		view.PrintDueDateUpdated(*dueDateID, dueDate)

	case "duration":
		durationCmd.Parse(os.Args[2:])
		if *durationID <= 0 {
			fmt.Println("有効なIDを指定してください")
			os.Exit(1)
		}
		if *durationValue < 0 {
			fmt.Println("見積時間は0以上で指定してください")
			os.Exit(1)
		}
		err := taskController.UpdateEstimatedDuration(*durationID, *durationValue)
		if err != nil {
			view.PrintError(err)
			os.Exit(1)
		}
		view.PrintDurationUpdated(*durationID, *durationValue)

	default:
		fmt.Println("不明なコマンドです。add, list, complete, delete, priority, due, durationのいずれかを使用してください")
		os.Exit(1)
	}
}
