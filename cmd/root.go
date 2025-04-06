package cmd

import (
	"fmt"
	"os"
	"time"

	"task-recommender/pkg/todo"

	"github.com/spf13/cobra"
)

var (
	taskList *todo.TaskList
	apiURL   string
)

var rootCmd = &cobra.Command{
	Use:   "task-recommender",
	Short: "タスク管理CLIアプリケーション",
	Long: `タスク管理CLIアプリケーション

各パラメータのフォーマット:
  優先度: 0（低）、1（中）、2（高）
  期限: YYYY/MM/DD-HH:mm（例: 2024/03/21-15:04）
  推定所要時間: 分単位の整数（例: 30）

使用例:
  タスク追加:
    task-recommender add --name 会議の準備 --priority 2 --due 2024/03/21-15:04 --duration 30
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		taskList = todo.NewTaskList(apiURL)
		if err := taskList.FetchTasks(); err != nil {
			fmt.Printf("タスクの取得に失敗: %v\n", err)
		}
	},
}

var (
	taskName     string
	taskPriority string
	taskDue      string
	taskDuration int
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "タスクを追加",
	Run: func(cmd *cobra.Command, args []string) {
		dueDate, err := time.Parse("2006/01/02-15:04", taskDue)
		if err != nil {
			fmt.Printf("期限の形式が不正です: %v\n", err)
			return
		}

		id := taskList.AddTask(taskName, taskPriority, dueDate, taskDuration)
		fmt.Printf("タスクを追加しました（ID: %s）\n", id)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "タスク一覧を表示",
	Run: func(cmd *cobra.Command, args []string) {
		taskList.ShowTasks()
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "タスクを削除",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		if taskList.DeleteTask(id) {
			fmt.Println("タスクを削除しました")
		} else {
			fmt.Println("指定されたIDのタスクが見つかりません")
		}
	},
}

var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "タスクを完了にする",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		for _, task := range taskList.Tasks {
			if task.ID == id {
				taskList.MarkComplete(task)
				fmt.Println("タスクを完了にしました")
				return
			}
		}
		fmt.Println("指定されたIDのタスクが見つかりません")
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&apiURL, "api", "http://localhost:8080", "APIサーバーのURL")

	addCmd.Flags().StringVarP(&taskName, "name", "n", "", "タスク名")
	addCmd.Flags().StringVarP(&taskPriority, "priority", "p", "", "優先度（0: 低, 1: 中, 2: 高）")
	addCmd.Flags().StringVarP(&taskDue, "due", "d", "", "期限（YYYY/MM/DD-HH:mm）")
	addCmd.Flags().IntVarP(&taskDuration, "duration", "t", 0, "推定所要時間（分）")

	addCmd.MarkFlagRequired("name")
	addCmd.MarkFlagRequired("priority")
	addCmd.MarkFlagRequired("due")
	addCmd.MarkFlagRequired("duration")
}

func Execute() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(completeCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
