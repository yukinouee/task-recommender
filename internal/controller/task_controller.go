package controller

import (
	"time"

	"task-recommender/internal/service"
)

type TaskController struct {
	service *service.TaskService
}

func NewTaskController(service *service.TaskService) *TaskController {
	return &TaskController{service: service}
}

func (c *TaskController) AddTask(title, description string, priority int, dueDate time.Time, estimatedDuration int) (int, error) {
	return c.service.AddTask(title, description, priority, dueDate, estimatedDuration)
}

func (c *TaskController) ListTasks() (interface{}, error) {
	return c.service.ListTasks()
}

func (c *TaskController) CompleteTask(id int) error {
	return c.service.CompleteTask(id)
}

func (c *TaskController) DeleteTask(id int) error {
	return c.service.DeleteTask(id)
}

func (c *TaskController) UpdatePriority(id, priority int) error {
	return c.service.UpdatePriority(id, priority)
}

func (c *TaskController) UpdateDueDate(id int, dueDate time.Time) error {
	return c.service.UpdateDueDate(id, dueDate)
}

func (c *TaskController) UpdateEstimatedDuration(id, duration int) error {
	return c.service.UpdateEstimatedDuration(id, duration)
}
