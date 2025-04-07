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

func (s *TaskService) AddTask(title, description string, priority int, dueDate time.Time, estimatedDuration int) (int, error) {
	var id int
	err := s.db.QueryRow(
		`INSERT INTO tasks 
        (title, description, done, priority, due_date, estimated_duration, created_at) 
        VALUES ($1, $2, false, $3, $4, $5, $6) 
        RETURNING id`,
		title, description, priority, dueDate, estimatedDuration, time.Now(),
	).Scan(&id)
	return id, err
}

func (s *TaskService) ListTasks() ([]model.Task, error) {
	rows, err := s.db.Query(`
        SELECT id, title, description, done, priority, due_date, estimated_duration, created_at, completed_at 
        FROM tasks 
        ORDER BY priority DESC, due_date ASC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		var completedAt sql.NullTime
		var dueDate sql.NullTime

		err := rows.Scan(
			&t.ID, &t.Title, &t.Description, &t.Done,
			&t.Priority, &dueDate, &t.EstimatedDuration,
			&t.CreatedAt, &completedAt,
		)
		if err != nil {
			return nil, err
		}

		if completedAt.Valid {
			t.CompletedAt = completedAt.Time
		}

		if dueDate.Valid {
			t.DueDate = dueDate.Time
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

func (s *TaskService) UpdatePriority(id, priority int) error {
	_, err := s.db.Exec(
		"UPDATE tasks SET priority = $1 WHERE id = $2",
		priority, id,
	)
	return err
}

func (s *TaskService) UpdateDueDate(id int, dueDate time.Time) error {
	_, err := s.db.Exec(
		"UPDATE tasks SET due_date = $1 WHERE id = $2",
		dueDate, id,
	)
	return err
}

func (s *TaskService) UpdateEstimatedDuration(id, duration int) error {
	_, err := s.db.Exec(
		"UPDATE tasks SET estimated_duration = $1 WHERE id = $2",
		duration, id,
	)
	return err
}
