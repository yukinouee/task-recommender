package model

import (
	"time"
)

// @swagger:model Task
type Task struct {
	// @タスクのID
	// @example: 1
	ID int `json:"id"`

	// タスクのタイトル
	// @example: 牛乳を買う
	// @required: true
	Title string `json:"title"`

	// @タスクの説明
	// @example: スーパーで低脂肪牛乳を購入する
	Description string `json:"description"`

	// @タスクの完了状態
	// @example: false
	Done bool `json:"done"`

	// @タスクの優先度 (1=低, 2=中, 3=高)
	// @example: 2
	// @min: 1
	// @max: 3
	Priority int `json:"priority"`

	// @タスクの期限日
	// @example: 2023-12-31T00:00:00Z
	DueDate time.Time `json:"due_date"`

	// @タスクの見積所要時間（分）
	// @example: 30
	// @min: 0
	EstimatedDuration int `json:"estimated_duration"`

	// @タスクの作成日時
	// @example: 2023-01-01T10:00:00Z
	CreatedAt time.Time `json:"created_at"`

	// @タスクの完了日時
	// @example: 2023-01-02T15:30:00Z
	CompletedAt time.Time `json:"completed_at,omitempty"`
}
