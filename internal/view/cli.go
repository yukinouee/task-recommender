package view

import (
	"fmt"
	"time"

	"task-recommender/internal/model"
)

func PrintTaskList(tasks []model.Task) {
	if len(tasks) == 0 {
		fmt.Println("タスクがありません")
		return
	}

	fmt.Println("ID | 優先度 | タイトル | 説明 | 期限 | 見積時間(分) | 状態 | 作成日 | 完了日")
	fmt.Println("---------------------------------------------------------------------------------")
	for _, t := range tasks {
		status := "未完了"
		completedAt := ""
		if t.Done {
			status = "完了"
			completedAt = t.CompletedAt.Format("2006-01-02 15:04:05")
		}

		dueDate := ""
		if !t.DueDate.IsZero() {
			dueDate = t.DueDate.Format("2006-01-02")
		}

		// 優先度を文字列に変換
		priorityStr := "低"
		if t.Priority == 2 {
			priorityStr = "中"
		} else if t.Priority >= 3 {
			priorityStr = "高"
		}

		fmt.Printf("%d | %s | %s | %s | %s | %d | %s | %s | %s\n",
			t.ID, priorityStr, t.Title, t.Description, dueDate, t.EstimatedDuration,
			status, t.CreatedAt.Format("2006-01-02 15:04:05"), completedAt)
	}
}

func PrintTaskAdded(id int, title string) {
	fmt.Printf("タスク追加: ID=%d, タイトル=%s\n", id, title)
}

func PrintTaskCompleted(id int) {
	fmt.Printf("タスク完了: ID=%d\n", id)
}

func PrintTaskDeleted(id int) {
	fmt.Printf("タスク削除: ID=%d\n", id)
}

func PrintPriorityUpdated(id int, priority int) {
	priorityStr := "低"
	if priority == 2 {
		priorityStr = "中"
	} else if priority >= 3 {
		priorityStr = "高"
	}
	fmt.Printf("優先度更新: ID=%d, 優先度=%s\n", id, priorityStr)
}

func PrintDueDateUpdated(id int, dueDate time.Time) {
	fmt.Printf("期限日更新: ID=%d, 期限日=%s\n", id, dueDate.Format("2006-01-02"))
}

func PrintDurationUpdated(id int, duration int) {
	fmt.Printf("見積時間更新: ID=%d, 見積時間=%d分\n", id, duration)
}

func PrintError(err error) {
	fmt.Printf("エラー: %v\n", err)
}
