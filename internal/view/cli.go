package view

import (
	"fmt"

	"task-recommender/internal/model"
)

func PrintTaskList(tasks []model.Task) {
	if len(tasks) == 0 {
		fmt.Println("タスクがありません")
		return
	}

	fmt.Println("ID | タイトル | 説明 | 状態 | 作成日 | 完了日")
	fmt.Println("--------------------------------------------------")
	for _, t := range tasks {
		status := "未完了"
		completedAt := ""
		if t.Done {
			status = "完了"
			completedAt = t.CompletedAt.Format("2006-01-02 15:04:05")
		}
		fmt.Printf("%d | %s | %s | %s | %s | %s\n",
			t.ID, t.Title, t.Description, status,
			t.CreatedAt.Format("2006-01-02 15:04:05"), completedAt)
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

func PrintError(err error) {
	fmt.Printf("エラー: %v\n", err)
}
