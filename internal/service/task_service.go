package service

import (
	"database/sql"
	"time"

	"task-recommender/internal/model"
)

type TaskService struct {
	db *sql.DB
}

func NewTaskService(db *sql.DB) *TaskService {
	return &TaskService{db: db}
}

func (s *TaskService) AddTask(title, description string) (int, error) {
	var id int
	err := s.db.QueryRow(
		"INSERT INTO tasks (title, description, done, created_at) VALUES ($1, $2, false, $3) RETURNING id",
		title, description, time.Now(),
	).Scan(&id)
	return id, err
}

func (s *TaskService) ListTasks() ([]model.Task, error) {
	rows, err := s.db.Query("SELECT id, title, description, done, created_at, completed_at FROM tasks ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		var completedAt sql.NullTime
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Done, &t.CreatedAt, &completedAt)
		if err != nil {
			return nil, err
		}
		if completedAt.Valid {
			t.CompletedAt = completedAt.Time
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (s *TaskService) CompleteTask(id int) error {
	_, err := s.db.Exec(
		"UPDATE tasks SET done = true, completed_at = $1 WHERE id = $2",
		time.Now(), id,
	)
	return err
}

func (s *TaskService) DeleteTask(id int) error {
	_, err := s.db.Exec("DELETE FROM tasks WHERE id = $1", id)
	return err
}
